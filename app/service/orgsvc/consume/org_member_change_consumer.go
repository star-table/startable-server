package orgsvc

import (
	platform_sdk "gitea.bjx.cloud/allstar/platform-sdk"
	"github.com/star-table/startable-server/app/facade/msgfacade"
	"github.com/star-table/startable-server/common/core/config"
	"github.com/star-table/startable-server/common/core/consts"
	"github.com/star-table/startable-server/common/core/errors"
	"github.com/star-table/startable-server/common/core/errs"
	"github.com/star-table/startable-server/common/core/logger"
	"github.com/star-table/startable-server/common/core/model"
	"github.com/star-table/startable-server/common/core/util/asyn"
	"github.com/star-table/startable-server/common/core/util/json"
	"github.com/star-table/startable-server/common/core/util/slice"
	"github.com/star-table/startable-server/common/core/util/stack"
	"github.com/star-table/startable-server/common/core/util/strs"
	"github.com/star-table/startable-server/common/library/cache"
	"github.com/star-table/startable-server/common/library/mq"
	"github.com/star-table/startable-server/common/model/bo"
)

var log = logger.GetDefaultLogger()

// 暂时忽略处理的组织
var ignoreOrgIds = []int64{
	16317,
}

func OrgMemberChangeConsume() {

	log.Infof("mq消息-组织成员变动消费者启动成功")

	orgMemberChangeTopicConfig := config.GetMqOrgMemberChangeConfig()

	client := *mq.GetMQClient()
	_ = client.ConsumeMessage(orgMemberChangeTopicConfig.Topic, orgMemberChangeTopicConfig.GroupId, func(message *model.MqMessageExt) errors.SystemErrorInfo {
		defer func() {
			if r := recover(); r != nil {
				log.Error(errs.BuildSystemErrorInfoWithPanicRecover(r, stack.GetStack()))
			}
		}()

		log.Infof("mq消息-组织成员变动消费信息 topic %s, value %s", message.Topic, message.Body)

		orgMemberChange := &bo.OrgMemberChangeBo{}
		err := json.FromJson(message.Body, orgMemberChange)
		if err != nil {
			log.Error(err)
			return errs.JSONConvertError
		}

		orgId := orgMemberChange.OrgId

		var businessErr errs.SystemErrorInfo = nil

		changeType := orgMemberChange.ChangeType
		//业务处理
		switch changeType {
		//禁用
		case consts.OrgMemberChangeTypeDisable:
			////暂时先不禁用
			//return nil
			businessErr = domain.ModifyOrgMemberStatus(orgId, []int64{orgMemberChange.UserId}, consts.AppStatusHidden, 0)
		//启用
		case consts.OrgMemberChangeTypeEnable:
			businessErr = domain.ModifyOrgMemberStatus(orgId, []int64{orgMemberChange.UserId}, consts.AppStatusEnable, 0)
		//添加用户
		case consts.OrgMemberChangeTypeAdd, consts.OrgMemberChangeTypeAddDisable:
			err := addUser(orgMemberChange)
			if err != nil {
				return err
			}
		case consts.OrgMemberChangeTypeRemove:
			businessErr = domain.RemoveOrgMember(orgId, []int64{orgMemberChange.UserId}, 0)
		case consts.OrgMemberChangeTypeUpdate: // 更新用户信息的回调处理
			err := updateUser(orgMemberChange)
			if err != nil {
				log.Error(err)
				return err
			}
		}

		if businessErr != nil {
			log.Error(businessErr)
		}

		//在并发操作时，有几率更新失败，所以忽略异常
		return nil
	}, func(message *model.MqMessageExt) {
		//失败的消息处理回调
		mqMessage := message.MqMessage

		log.Infof("mq消息消费失败-动态-信息详情 topic %s, value %s", message.Topic, message.Body)

		orgMemberChange := &bo.OrgMemberChangeBo{}
		err := json.FromJson(message.Body, orgMemberChange)
		if err != nil {
			log.Error(err)
			return
		}

		msgErr := msgfacade.InsertMqConsumeFailMsgRelaxed(mqMessage, 0, orgMemberChange.OrgId)
		if msgErr != nil {
			log.Errorf("mq消息消费失败，入表失败，消息内容：%s, 失败信息: %v", json.ToJsonIgnoreError(mqMessage), msgErr)
		}
	})
}

func addUser(orgMemberChange *bo.OrgMemberChangeBo) errors.SystemErrorInfo {
	orgOutInfo, err := domain.GetBaseOrgInfo(orgMemberChange.OrgId)
	if err != nil {
		log.Error(err)
		return err
	}
	//判断是否在白名单里面，如果在的话就不处理
	key := consts.CacheFeishuNotSyncUser
	value, err1 := cache.Get(key)
	if err1 != nil {
		log.Error(err)
		return errs.CacheProxyError
	}
	if value != "" {
		fsUserOrgIdWhiteList := &[]string{}
		err1 = json.FromJson(value, fsUserOrgIdWhiteList)
		if err != nil {
			log.Info(strs.ObjectToString(err))
			return errs.JSONConvertError
		}
		// 当前企业在白名单当中，则不处理
		if ok, _ := slice.Contain(*fsUserOrgIdWhiteList, orgOutInfo.OutOrgId); ok {
			return nil
		}
	}

	// FsAuth 中会有兜底初始化飞书用户的逻辑。因此当执行初始化时，需要新增「推送事件到前端」的消息队列。
	baseUserInfo, err := domain.PlatformAuth(orgOutInfo.SourceChannel, orgOutInfo.OutOrgId, orgMemberChange.OpenId, "", orgMemberChange.DeptIds...)
	if err != nil {
		log.Error(err)
		return err
	}
	if baseUserInfo.OrgUserIsDelete == consts.AppIsDeleted {
		inDisabled := orgMemberChange.ChangeType == consts.OrgMemberChangeTypeAddDisable
		log.Infof("添加用户是否被禁用 %v", inDisabled)
		err = domain.AddOrgMember(baseUserInfo.OrgId, baseUserInfo.UserId, 0, false, inDisabled)
		if err != nil {
			log.Error(err)
			return err
		}
	}
	//~~推送飞书帮助信息~~ 整改：组织增加用户后，不再向新加的用户推送帮助消息。
	//domain.FeishuMemberHelpMsgToMq(bo.FeishuHelpObjectBo{OpenIds:[]string{orgMemberChange.OpenId}, TenantKey:orgOutInfo.OutOrgId, OrgId:orgId})

	return nil
}

func updateUser(orgMemberChange *bo.OrgMemberChangeBo) errs.SystemErrorInfo {
	orgOutInfo, err := domain.GetBaseOrgInfo(orgMemberChange.OrgId)
	if err != nil {
		log.Error(err)
		return err
	}
	baseInfo, err := domain.GetBaseUserInfoByEmpId(orgMemberChange.OrgId, orgMemberChange.OpenId)
	if err != nil {
		return err
	}

	client, sdkErr := platform_sdk.GetClient(orgMemberChange.SourceChannel, orgOutInfo.OutOrgId)
	if sdkErr != nil {
		log.Errorf("[updateUser] platform_sdk.GetClient err: %v", sdkErr)
		return errs.BuildSystemErrorInfo(errs.PlatFormOpenApiCallError, sdkErr)
	}
	scopeUser, err := domain.GetPlatformUserDetailInfo(client, baseInfo.OutUserId, orgMemberChange.DeptIds...)
	if err != nil {
		return err
	}

	// 更新用户信息
	businessErr := domain.UpdateUserInfoWithPlatformUserInfo(orgMemberChange.OrgId, orgMemberChange.UserId, scopeUser,
		orgMemberChange.NewOutUserId, orgMemberChange.SourceChannel)
	if businessErr != nil {
		log.Error(businessErr)
		return businessErr
	}
	// 通过 userId 查询用户部门等信息
	userDeptIdMap, oriErr := domain.GetUserDepartmentIdMap(orgMemberChange.OrgId, []int64{orgMemberChange.UserId})
	if oriErr != nil {
		log.Error(oriErr)
		return errs.BuildSystemErrorInfo(errs.MysqlOperateError, oriErr)
	}
	currentUserDeptId := int64(0)
	if tmpDeptId, ok := userDeptIdMap[orgMemberChange.UserId]; ok {
		currentUserDeptId = tmpDeptId
	}
	asyn.Execute(func() {
		domain.PushUpdateOrgMemberNotice(orgMemberChange.OrgId, currentUserDeptId, []int64{orgMemberChange.UserId}, 0)
	})

	return nil
}
