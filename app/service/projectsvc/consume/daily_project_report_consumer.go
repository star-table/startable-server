package consume

import (
	platform_sdk "gitea.bjx.cloud/allstar/platform-sdk"
	sdkConsts "gitea.bjx.cloud/allstar/platform-sdk/consts"
	sdk_vo "gitea.bjx.cloud/allstar/platform-sdk/vo"

	"github.com/star-table/startable-server/common/extra/card"
	"github.com/star-table/startable-server/common/model/vo/commonvo"

	"github.com/star-table/startable-server/app/facade/msgfacade"
	"github.com/star-table/startable-server/common/core/config"
	consts2 "github.com/star-table/startable-server/common/core/consts"
	"github.com/star-table/startable-server/common/core/errors"
	"github.com/star-table/startable-server/common/core/errs"
	"github.com/star-table/startable-server/common/core/model"
	"github.com/star-table/startable-server/common/core/threadlocal"
	"github.com/star-table/startable-server/common/core/util/json"
	"github.com/star-table/startable-server/common/core/util/stack"
	"github.com/star-table/startable-server/common/library/mq"
	"github.com/star-table/startable-server/common/model/bo"
	"github.com/jtolds/gls"
)

func DailyProjectReportMsgConsumer() {
	log.Infof("mq消息-每日项目日报Msg消费者启动成功")

	dailyProjectReportMsgTopicConfig := config.GetMqDailyProjectReportMsgTopicConfig()
	client := *mq.GetMQClient()

	_ = client.ConsumeMessage(dailyProjectReportMsgTopicConfig.Topic, dailyProjectReportMsgTopicConfig.GroupId, func(message *model.MqMessageExt) errors.SystemErrorInfo {
		defer func() {
			if r := recover(); r != nil {
				log.Error(errs.BuildSystemErrorInfoWithPanicRecover(r, stack.GetStack()))
			}
		}()

		msgBo := &bo.DailyProjectReportMsgBo{}
		err := json.FromJson(message.Body, msgBo)
		if err != nil {
			log.Error(err)
			return errs.BuildSystemErrorInfo(errs.JSONConvertError)
		}

		if msgBo.SourceChannel == "" || msgBo.OutOrgId == "" || len(msgBo.OpenIds) == 0 {
			return nil
		}

		threadlocal.Mgr.SetValues(gls.Values{consts2.TraceIdKey: msgBo.ScheduleTraceId}, func() {
			log.Infof("[DailyProjectReportMsgConsumer] mq消息-每日项目日报msg-信息详情 topic %s, value %s", message.Topic, message.Body)
			SendCardDailyProjectReport(msgBo)
		})

		return nil
	}, func(message *model.MqMessageExt) {
		//失败的消息处理回调
		mqMessage := message.MqMessage

		log.Infof("mq消息消费失败-动态-信息详情 topic %s, value %s", message.Topic, message.Body)

		msgBo := &bo.DailyProjectReportMsgBo{}
		err := json.FromJson(message.Body, msgBo)
		if err != nil {
			log.Error(err)
			return
		}
		msgErr := msgfacade.InsertMqConsumeFailMsgRelaxed(mqMessage, int(msgBo.PushType), msgBo.OrgId)
		if msgErr != nil {
			log.Errorf("mq消息消费失败，入表失败，消息内容：%s, 失败信息: %v", json.ToJsonIgnoreError(mqMessage), msgErr)
		}
	})
}

// SendCardDailyProjectReport 发送“项目日报”卡片
func SendCardDailyProjectReport(msg *bo.DailyProjectReportMsgBo) {
	sdk, err := platform_sdk.GetClient(msg.SourceChannel, "")
	if err != nil {
		log.Errorf("[SendCardDailyProjectReport] platform_sdk.GetClient, sourceChannel: %v, orgId: %d, err: %v", msg.SourceChannel, msg.OrgId, err)
		return
	}
	host := ""
	serverCommon := config.GetConfig().ServerCommon
	if serverCommon != nil {
		if msg.SourceChannel == sdkConsts.SourceChannelWeixin {
			host = serverCommon.WeiXinHost
		} else {
			host = serverCommon.Host
		}
	}
	link, sErr := sdk.GetProjectLink(&sdk_vo.GetProjectLinkReq{
		AppId:         msg.AppId,
		ProjectId:     msg.ProjectId,
		ProjectTypeId: msg.ProjectTypeId,
		Host:          host,
		CorpId:        msg.OutOrgId,
	})
	if sErr != nil {
		log.Errorf("[SendCardDailyProjectReport] platform_sdk.GetProjectLink, err: %v", sErr)
		return
	}

	// 项目日报卡片
	projectCard := card.GetDailyProjectCard(msg, link.Url)
	errSys := card.PushCard(msg.OrgId, &commonvo.PushCard{
		OrgId:         msg.OrgId,
		OutOrgId:      msg.OutOrgId,
		SourceChannel: msg.SourceChannel,
		OpenIds:       msg.OpenIds,
		CardMsg:       projectCard,
	})

	if errSys != nil {
		log.Errorf("[SendCardDailyProjectReport] err:%v", errSys)
	}
}
