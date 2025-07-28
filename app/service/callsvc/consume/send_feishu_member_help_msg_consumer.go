package callsvc

import (
	"gitea.bjx.cloud/allstar/polaris-backend/facade/msgfacade"
	"github.com/star-table/startable-server/common/core/config"
	"github.com/star-table/startable-server/common/core/errors"
	"github.com/star-table/startable-server/common/core/errs"
	"github.com/star-table/startable-server/common/core/model"
	"github.com/star-table/startable-server/common/core/util/json"
	"github.com/star-table/startable-server/common/core/util/stack"
	"github.com/star-table/startable-server/common/library/mq"
	"github.com/star-table/startable-server/common/model/bo"
)

func SendFeishuHelpMsgConsumer() {

	log.Infof("mq消息-推送飞书欢迎语消费者启动成功")

	info := config.GetHelpMessagePushToFeiShu()

	client := *mq.GetMQClient()
	_ = client.ConsumeMessage(info.Topic, info.GroupId, func(message *model.MqMessageExt) errors.SystemErrorInfo {
		defer func() {
			if r := recover(); r != nil {
				log.Error(errs.BuildSystemErrorInfoWithPanicRecover(r, stack.GetStack()))
			}
		}()

		log.Infof("mq消息-推送飞书欢迎语消费信息 topic %s, value %s", message.Topic, message.Body)

		infoBo := &bo.FeishuHelpObjectBo{}
		err := json.FromJson(message.Body, infoBo)
		if err != nil {
			log.Error(err)
			return errs.JSONConvertError
		}
		//pushErr := feishu.PushBotHelpNotice(infoBo.TenantKey, "", infoBo.OpenIds...)
		//pushErr := feishu.PushWelcomeMsgToOrgMember(infoBo.TenantKey, infoBo.OpenIds...)
		pushErr := feishu.PushBotToInstallerRemindMsg(infoBo.TenantKey, infoBo.OpenIds)
		if pushErr != nil {
			log.Error(pushErr)
			return pushErr
		}

		return nil
	}, func(message *model.MqMessageExt) {
		//失败的消息处理回调
		mqMessage := message.MqMessage
		log.Infof("mq消息消费失败-推送飞书欢迎语消费信息 topic %s, value %s", message.Topic, message.Body)

		msgErr := msgfacade.InsertMqConsumeFailMsgRelaxed(mqMessage, 0, 0)
		if msgErr != nil {
			log.Errorf("mq消息消费失败，入表失败，消息内容：%s, 失败信息: %v", json.ToJsonIgnoreError(mqMessage), msgErr)
		}
	})
}
