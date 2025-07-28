package callsvc

import (
	"gitea.bjx.cloud/allstar/polaris-backend/facade/msgfacade"
	"github.com/star-table/startable-server/common/core/config"
	"github.com/star-table/startable-server/common/core/errors"
	"github.com/star-table/startable-server/common/core/errs"
	"github.com/star-table/startable-server/common/core/model"
	"github.com/star-table/startable-server/common/core/util/json"
	"github.com/star-table/startable-server/common/core/util/stack"
	"github.com/star-table/startable-server/common/core/util/uuid"
	"github.com/star-table/startable-server/common/library/mq"
	"github.com/star-table/startable-server/common/model/bo"
)

func FeiShuCallBackUesrChangeConsumer() {

	log.Infof("mq消息-飞书回调员工变更消费者启动成功")

	fsCallbackInfo := config.GetFeishuCallBackUserChange()

	client := *mq.GetMQClient()
	_ = client.ConsumeMessage(fsCallbackInfo.Topic, fsCallbackInfo.GroupId, func(message *model.MqMessageExt) errors.SystemErrorInfo {
		defer func() {
			if r := recover(); r != nil {
				log.Error(errs.BuildSystemErrorInfoWithPanicRecover(r, stack.GetStack()))
			}
		}()

		log.Infof("mq消息-飞书回调员工变更信息 topic %s, value %s", message.Topic, message.Body)

		fsCallBackBo := &bo.FeiShuCallBackMqBo{}
		err := json.FromJson(message.Body, fsCallBackBo)
		if err != nil {
			log.Error(err)
			return errs.JSONConvertError
		}

		reqEvent := fsCallBackBo.Type
		handler, ok := feishu.FsCallBackHandlers[reqEvent]
		if !ok {
			log.Infof("%s事件不在正常处理范围内", reqEvent)
			return nil
		}

		msgId := uuid.NewUuid()
		log.Infof("飞书回调员工变更异步处理准备msgId:%s, msgBody: %s", msgId, message.Body)
		_, handleErr := handler.Handle(fsCallBackBo.Data)
		if handleErr != nil {
			log.Infof("飞书回调员工变更异步处理失败msgId:%s, msgBody: %s", msgId, message.Body)
			log.Error(handleErr)
			return handleErr
		}

		return nil
	}, func(message *model.MqMessageExt) {
		//失败的消息处理回调
		mqMessage := message.MqMessage
		log.Infof("mq消息消费失败-飞书回调消费信息 topic %s, value %s", message.Topic, message.Body)

		msgErr := msgfacade.InsertMqConsumeFailMsgRelaxed(mqMessage, 0, 0)
		if msgErr != nil {
			log.Errorf("mq消息消费失败，入表失败，消息内容：%s, 失败信息: %v", json.ToJsonIgnoreError(mqMessage), msgErr)
		}
	})
}
