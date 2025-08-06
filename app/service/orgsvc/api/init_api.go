package api

import (
	"net/http"

	sdk_const "gitea.bjx.cloud/allstar/platform-sdk/consts"
	"github.com/gin-gonic/gin"
	"github.com/star-table/startable-server/common/core/consts"
	"github.com/star-table/startable-server/common/core/errs"
	"github.com/star-table/startable-server/common/core/logger"
	"github.com/star-table/startable-server/common/core/util/uuid"
	"github.com/star-table/startable-server/common/library/cache"
	"github.com/star-table/startable-server/common/model/vo"
	"github.com/star-table/startable-server/common/model/vo/orgvo"
	orgsvcService "github.com/star-table/startable-server/app/service/orgsvc/service"
	orgsvcDomain "github.com/star-table/startable-server/app/service/orgsvc/domain"
)

// 飞书初始化调用
func InitOrg(c *gin.Context) {
	var reqVo orgvo.InitOrgReqVo
	if err := c.ShouldBindJSON(&reqVo); err != nil {
		c.JSON(http.StatusOK, orgvo.OrgInitRespVo{Err: vo.NewErr(errs.BuildSystemErrorInfo(errs.ReqParamsValidateError, err))})
		return
	}

	respVo := orgvo.OrgInitRespVo{}
	tenantKey := reqVo.InitOrg.OutOrgId
	//加锁
	lockKey := consts.NewFeiShuCorpInitKey + tenantKey
	uuid := uuid.NewUuid()
	suc, err := cache.TryGetDistributedLock(lockKey, uuid)
	logger.Infof("准备获取分布式锁 %v", suc)
	if err != nil {
		logger.Error("InitOrg TryGetDistributedLock error", logger.ErrorField(err))
		c.JSON(http.StatusOK, orgvo.OrgInitRespVo{
			Err:   vo.NewErr(errs.RedisOperateError),
			OrgId: 0,
		})
		return
	}
	if suc {
		logger.Infof("获取分布式锁成功 %v", suc)
		defer func() {
			if _, lockErr := cache.ReleaseDistributedLock(lockKey, uuid); lockErr != nil {
				logger.Error("InitOrg ReleaseDistributedLock error", logger.ErrorField(lockErr))
			}
		}()
		//如果已经有可用团队，直接正常进去(防止一个人停在注册页面，重复创建)
		orgOutInfo, _ := orgsvcDomain.GetOrgOutInfoByTenantKey(tenantKey)
		if orgOutInfo != nil {
			if reqVo.InitOrg.PermanentCode != "" {
				_ = orgsvcDomain.UpdateOrgOutInfoPermanentCode(orgOutInfo.OrgId, reqVo.InitOrg.PermanentCode)
			}
			respVo.OrgId = orgOutInfo.OrgId
			c.JSON(http.StatusOK, respVo)
			return
		}

		orgId, err := orgsvcService.InitOrg(reqVo.InitOrg)
		if err != nil {
			logger.Error("InitOrg service.InitOrg error", logger.ErrorField(err))
			respVo.Err = vo.NewErr(err)
			c.JSON(http.StatusOK, respVo)
			return
		}
		respVo.OrgId = orgId
		if !respVo.Failure() {
			//初始化示例数据
			orgConfig, configErr := orgsvcService.GetOrgInfoBo(respVo.OrgId)
			if configErr != nil {
				logger.Error("InitOrg GetOrgInfoBo error", logger.ErrorField(configErr))
			} else {
				if orgConfig.Owner != 0 {
					err1 := orgsvcService.NewbieGuideInit(respVo.OrgId, orgConfig.Owner, reqVo.InitOrg.SourceChannel)
					if err1 != nil {
						logger.Error("InitOrg NewbieGuideInit error", logger.ErrorField(err1))
					}
				}
			}

			//组织付费等级确认（可能付费之后再初始化组织的）
			if reqVo.InitOrg.SourceChannel == sdk_const.SourceChannelFeishu {
				payLevelErr := orgsvcDomain.SetFsOrgPayLevel(reqVo.InitOrg.OutOrgId)
				if payLevelErr != nil {
					logger.Error("InitOrg SetFsOrgPayLevel error", logger.ErrorField(payLevelErr))
				}
			}

			if reqVo.InitOrg.SourceChannel == sdk_const.SourceChannelDingTalk {
				logger.Infof("ding update orgConfig, orgId:%v, outOrgId:%v", orgId, reqVo.InitOrg.OutOrgId)
				// 钉钉  有可能先初始化组织了，再更新订单、付费等级等信息
				errDing := orgsvcService.UpdateDingOrgConfig(orgId, reqVo.InitOrg.OutOrgId)
				if errDing != nil {
					logger.Error("InitOrg UpdateDingOrgConfig error", logger.ErrorField(errDing))
				}
			}

			if reqVo.InitOrg.SourceChannel == sdk_const.SourceChannelWeixin {
				logger.Infof("weixin update orgConfig, orgId:%v, outOrgId:%v", orgId, reqVo.InitOrg.OutOrgId)
				payLevelErr := orgsvcDomain.SetWeiXinOrgPayLevel(reqVo.InitOrg.OutOrgId)
				if payLevelErr != nil {
					logger.Error("InitOrg SetWeiXinOrgPayLevel error", logger.ErrorField(payLevelErr))
				}
			}

		}
	}

	//asyn.Execute(func() {
	//	time.Sleep(6 * time.Second)
	//	//向成员发送欢迎语加入mq
	//	sendHelpErr := orgsvcService.SendFeishuMemberHelpMsg(reqVo.InitOrg.OutOrgId, orgId, reqVo.InitOrg.OutOrgOwnerId)
	//	if sendHelpErr != nil {
	//		logger.Error(sendHelpErr)
	//		return
	//	}
	//})

	c.JSON(http.StatusOK, respVo)
}

// 发送飞书帮助信息
func SendFeishuMemberHelpMsg(c *gin.Context) {
	var reqVo orgvo.SendFeishuMemberHelpMsgReqVo
	if err := c.ShouldBindJSON(&reqVo); err != nil {
		c.JSON(http.StatusOK, vo.VoidErr{Err: vo.NewErr(errs.BuildSystemErrorInfo(errs.ReqParamsValidateError, err))})
		return
	}

	err := orgsvcService.SendFeishuMemberHelpMsg(reqVo.TenantKey, reqVo.OrgId, reqVo.OwnerOpenId)
	if err != nil {
		logger.Error("SendFeishuMemberHelpMsg error", logger.ErrorField(err))
	}
	c.JSON(http.StatusOK, vo.VoidErr{Err: vo.NewErr(err)})
}

func ScheduleOrgUseMobileAndEmail(c *gin.Context) {
	var reqVo orgvo.ScheduleOrgUseMobileAndEmailReqVo
	if err := c.ShouldBindJSON(&reqVo); err != nil {
		c.JSON(http.StatusOK, vo.VoidErr{Err: vo.NewErr(errs.BuildSystemErrorInfo(errs.ReqParamsValidateError, err))})
		return
	}

	err := orgsvcDomain.ScheduleOrgUserMobileAndEmail(reqVo.TenantKey)
	if err != nil {
		logger.Error("ScheduleOrgUseMobileAndEmail error", logger.ErrorField(err))
	}
	c.JSON(http.StatusOK, vo.VoidErr{Err: vo.NewErr(err)})
}
