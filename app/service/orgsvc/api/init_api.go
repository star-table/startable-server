package orgsvc

import (
	sdk_const "gitea.bjx.cloud/allstar/platform-sdk/consts"
	"github.com/star-table/startable-server/app/service"
	"github.com/star-table/startable-server/common/core/consts"
	"github.com/star-table/startable-server/common/core/errs"
	"github.com/star-table/startable-server/common/core/util/uuid"
	"github.com/star-table/startable-server/common/library/cache"
	"github.com/star-table/startable-server/common/model/vo"
	"github.com/star-table/startable-server/common/model/vo/orgvo"
	"github.com/google/martian/log"
)

// 飞书初始化调用
func (PostGreeter) InitOrg(reqVo orgvo.InitOrgReqVo) orgvo.OrgInitRespVo {
	respVo := orgvo.OrgInitRespVo{}
	tenantKey := reqVo.InitOrg.OutOrgId
	//加锁
	lockKey := consts.NewFeiShuCorpInitKey + tenantKey
	uuid := uuid.NewUuid()
	suc, err := cache.TryGetDistributedLock(lockKey, uuid)
	log.Infof("准备获取分布式锁 %v", suc)
	if err != nil {
		log.Error(err)
		return orgvo.OrgInitRespVo{
			Err:   vo.NewErr(errs.RedisOperateError),
			OrgId: 0,
		}
	}
	if suc {
		log.Infof("获取分布式锁成功 %v", suc)
		defer func() {
			if _, lockErr := cache.ReleaseDistributedLock(lockKey, uuid); lockErr != nil {
				log.Error(lockErr)
			}
		}()
		//如果已经有可用团队，直接正常进去(防止一个人停在注册页面，重复创建)
		orgOutInfo, _ := domain.GetOrgOutInfoByTenantKey(tenantKey)
		if orgOutInfo != nil {
			if reqVo.InitOrg.PermanentCode != "" {
				_ = domain.UpdateOrgOutInfoPermanentCode(orgOutInfo.OrgId, reqVo.InitOrg.PermanentCode)
			}
			respVo.OrgId = orgOutInfo.OrgId
			return respVo
		}

		orgId, err := service.InitOrg(reqVo.InitOrg)
		if err != nil {
			log.Error(err)
			respVo.Err = vo.NewErr(err)
			return respVo
		}
		respVo.OrgId = orgId
		if !respVo.Failure() {
			//初始化示例数据
			orgConfig, configErr := service.GetOrgInfoBo(respVo.OrgId)
			if configErr != nil {
				log.Error(configErr)
			} else {
				if orgConfig.Owner != 0 {
					err1 := service.NewbieGuideInit(respVo.OrgId, orgConfig.Owner, reqVo.InitOrg.SourceChannel)
					if err1 != nil {
						log.Error(err1)
					}
				}
			}

			//组织付费等级确认（可能付费之后再初始化组织的）
			if reqVo.InitOrg.SourceChannel == sdk_const.SourceChannelFeishu {
				payLevelErr := domain.SetFsOrgPayLevel(reqVo.InitOrg.OutOrgId)
				if payLevelErr != nil {
					log.Error(payLevelErr)
				}
			}

			if reqVo.InitOrg.SourceChannel == sdk_const.SourceChannelDingTalk {
				log.Infof("ding update orgConfig, orgId:%v, outOrgId:%v", orgId, reqVo.InitOrg.OutOrgId)
				// 钉钉  有可能先初始化组织了，再更新订单、付费等级等信息
				errDing := service.UpdateDingOrgConfig(orgId, reqVo.InitOrg.OutOrgId)
				if errDing != nil {
					log.Error(errDing)
				}
			}

			if reqVo.InitOrg.SourceChannel == sdk_const.SourceChannelWeixin {
				log.Infof("weixin update orgConfig, orgId:%v, outOrgId:%v", orgId, reqVo.InitOrg.OutOrgId)
				payLevelErr := domain.SetWeiXinOrgPayLevel(reqVo.InitOrg.OutOrgId)
				if payLevelErr != nil {
					log.Error(payLevelErr)
				}
			}

		}
	}

	//asyn.Execute(func() {
	//	time.Sleep(6 * time.Second)
	//	//向成员发送欢迎语加入mq
	//	sendHelpErr := service.SendFeishuMemberHelpMsg(reqVo.InitOrg.OutOrgId, orgId, reqVo.InitOrg.OutOrgOwnerId)
	//	if sendHelpErr != nil {
	//		log.Error(sendHelpErr)
	//		return
	//	}
	//})

	return respVo
}

// 发送飞书帮助信息
func (PostGreeter) SendFeishuMemberHelpMsg(reqVo orgvo.SendFeishuMemberHelpMsgReqVo) vo.VoidErr {
	err := service.SendFeishuMemberHelpMsg(reqVo.TenantKey, reqVo.OrgId, reqVo.OwnerOpenId)
	return vo.VoidErr{Err: vo.NewErr(err)}
}

func (GetGreeter) ScheduleOrgUseMobileAndEmail(reqVo orgvo.ScheduleOrgUseMobileAndEmailReqVo) vo.VoidErr {
	err := domain.ScheduleOrgUserMobileAndEmail(reqVo.TenantKey)
	return vo.VoidErr{Err: vo.NewErr(err)}
}
