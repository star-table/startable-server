package callsvc

import (
	"fmt"

	sdk_const "gitea.bjx.cloud/allstar/platform-sdk/consts"
	"github.com/star-table/startable-server/app/facade/orgfacade"
	"github.com/star-table/startable-server/common/core/consts"
	"github.com/star-table/startable-server/common/core/errs"
	"github.com/star-table/startable-server/common/core/util/json"
	"github.com/star-table/startable-server/common/extra/card"
	"github.com/star-table/startable-server/common/model/vo"
	"github.com/star-table/startable-server/common/model/vo/commonvo"
	"github.com/star-table/startable-server/common/model/vo/orgvo"
)

type NoticeSubscribe struct{}

func (NoticeSubscribe) Handle(cardReq CardReq) (string, errs.SystemErrorInfo) {
	log.Infof("[NoticeSubscribe.Handle] 退订消息详情：%s", json.ToJsonIgnoreError(cardReq))
	input := cardReq.Action.Value
	orgId := input.OrgId
	openId := cardReq.OpenId
	if orgId == int64(0) {
		return "", nil
	}

	userBaseInfoResp := orgfacade.GetBaseUserInfoByEmpId(orgvo.GetBaseUserInfoByEmpIdReqVo{
		OrgId: orgId,
		EmpId: openId,
	})
	if userBaseInfoResp.Failure() {
		log.Error(userBaseInfoResp.Error())
		return "", nil
	}
	userBaseInfo := userBaseInfoResp.BaseUserInfo
	userId := userBaseInfo.UserId

	userConfig, err := orgfacade.GetUserConfigInfoRelaxed(orgId, userId)
	if err != nil {
		log.Errorf("[NoticeSubscribe.Handle] err: %v", err)
		return "", nil
	}

	typeName := ""
	req := vo.UpdateUserConfigReq{}
	switch input.Action {
	case consts.FsCardUnsubscribePersonalReport:
		//退订个人日报
		if userConfig.DailyReportMessageStatus == 2 {
			log.Errorf("[NoticeSubscribe.Handle][FsCardUnsubscribePersonalReport] 用户 %d 未开启个人日报推送", userId)
			return "", nil
		}
		typeName = "个人日报"
		req.DailyReportMessageStatus = 2
	case consts.FsCardUnsubscribeProjectReport:
		//退订项目日报
		if userConfig.DailyProjectReportMessageStatus == 2 {
			log.Errorf("[NoticeSubscribe.Handle][FsCardUnsubscribeProjectReport] 用户 %d 未开启项目日报推送", userId)
			return "", nil
		}
		typeName = "项目日报"
		req.DailyProjectReportMessageStatus = 2
	case consts.FsCardUnsubscribeIssueRemind:
		//退订任务提醒
		if userConfig.RemindMessageStatus == 2 && userConfig.ModifyMessageStatus == 2 && userConfig.RelationMessageStatus == 2 && userConfig.CommentAtMessageStatus == 2 {
			log.Errorf("[NoticeSubscribe.Handle][FsCardUnsubscribeIssueRemind] 用户 %d 未开启任务提醒推送", userId)
			return "", nil
		}
		typeName = "任务提醒"
		req.RemindMessageStatus = 2
		req.ModifyMessageStatus = 2
		req.RelationMessageStatus = 2
		req.CommentAtMessageStatus = 2
	case consts.FsCardUnsubscribeMyOwn: // 退订我负责的
		if userConfig.OwnerRangeStatus == consts.CardSubscribeSwitchOff {
			log.Errorf("[NoticeSubscribe.Handle][FsCardUnsubscribeMyOwn] 用户未开启. userId: %d", userId)
			return "", nil
		}
		req.OwnerRangeStatus = 2
	case consts.FsCardUnsubscribeMyCollaborate:
		if userConfig.OwnerRangeStatus == consts.CardSubscribeSwitchOff {
			log.Errorf("[NoticeSubscribe.Handle][FsCardUnsubscribeMyCollaborate] 用户未开启. userId: %d", userId)
			return "", nil
		}
		req.CollaborateMessageStatus = 2
	default:
		log.Errorf("[NoticeSubscribe.Handle] 不支持的退订类型：%s", input.Action)
		return "", nil
	}

	//更新个人通知配置
	resp := orgfacade.UpdateUserConfig(orgvo.UpdateUserConfigReqVo{
		UpdateUserConfigReq: req,
		OrgId:               orgId,
		UserId:              userId,
	})
	if resp.Failure() {
		log.Error(resp.Failure())
		return "", nil
	}

	content := fmt.Sprintf(consts.FsCardSubscribeNotice, typeName, consts.BotPrivateInsNameOfSetting)
	cardMsg := card.GetFsCardSubscribeNotice(content)

	errSys := card.PushCard(orgId, &commonvo.PushCard{
		OrgId:         orgId,
		OutOrgId:      userBaseInfo.OutOrgId,
		SourceChannel: sdk_const.SourceChannelFeishu,
		OpenIds:       []string{userBaseInfo.OutUserId},
		CardMsg:       cardMsg,
	})
	if errSys != nil {
		log.Errorf("[NoticeSubscribe.Handle] err:%v", errSys)
		return "", errSys
	}
	return "", nil
}
