package callsvc

import (
	"sync"

	"gitea.bjx.cloud/allstar/polaris-backend/facade/msgfacade"
	"github.com/star-table/startable-server/common/core/util/asyn"

	"github.com/star-table/startable-server/common/core/config"
	"github.com/star-table/startable-server/common/core/errors"
	"github.com/star-table/startable-server/common/core/errs"
	"github.com/star-table/startable-server/common/core/model"
	"github.com/star-table/startable-server/common/core/util/json"
	"github.com/star-table/startable-server/common/core/util/stack"
	"github.com/star-table/startable-server/common/core/util/uuid"
	"github.com/star-table/startable-server/common/library/mq"
)

func DingTalkCallBackConsumer() {
	// 启动5个消费协程
	maxCount := 5
	log.Infof("[DingTalkCallBackConsumer] mq消息-钉钉回调消费者启动成功")
	callbackInfo := config.GetDingTalkCallBackTopic()

	g := sync.WaitGroup{}
	g.Add(maxCount)
	for i := 0; i < maxCount; i++ {
		asyn.Execute(func() {
			consumeMessage(callbackInfo)
			g.Done()
		})
	}
	g.Wait()
}

func consumeMessage(info config.TopicConfigInfo) {
	client := *mq.GetMQClient()
	_ = client.ConsumeMessage(info.Topic, info.GroupId, func(message *model.MqMessageExt) errors.SystemErrorInfo {
		defer func() {
			if r := recover(); r != nil {
				log.Error(errs.BuildSystemErrorInfoWithPanicRecover(r, stack.GetStack()))
			}
		}()

		log.Infof("[DingTalkCallBackConsumer] mq消息-钉钉回调消费信息 topic %s, value %s", message.Topic, message.Body)

		callBackBo := &req_vo.DingTalkCallBackMqVo{}
		err := json.FromJson(message.Body, callBackBo)
		if err != nil {
			log.Error(err)
			return errs.JSONConvertError
		}

		reqEvent := callBackBo.Type
		if callBackBo.Data.BizType == consts.EventDingAppTrialBizType {
			reqEvent = consts.EventDingOrderTrail
		}

		handler, ok := dingtalk.DingTalkHandlers[reqEvent]
		if !ok {
			log.Infof("[DingTalkCallBackConsumer] %s事件不在正常处理范围内", reqEvent)
			return nil
		}

		msgId := uuid.NewUuid()
		log.Infof("[DingTalkCallBackConsumer] 钉钉回调消息异步处理准备msgId:%s, msgBody: %s", msgId, message.Body)
		handleErr := handler.Handle(callBackBo.Data)
		if handleErr != nil {
			log.Errorf("[DingTalkCallBackConsumer] 钉钉回调消息异步处理失败msgId:%s, msgBody: %s, err:%v", msgId, message.Body, handleErr)
			return errs.BuildSystemErrorInfoWithMessage(errs.DingTalkOpenApiCallError, handleErr.Error())
		}

		return nil
	}, func(message *model.MqMessageExt) {
		//失败的消息处理回调
		mqMessage := message.MqMessage
		log.Infof("[DingTalkCallBackConsumer] mq消息消费失败-钉钉回调消费信息 topic %s, value %s", message.Topic, message.Body)

		msgErr := msgfacade.InsertMqConsumeFailMsgRelaxed(mqMessage, 0, 0)
		if msgErr != nil {
			log.Errorf("[DingTalkCallBackConsumer] mq消息消费失败，入表失败，消息内容：%s, 失败信息: %v", json.ToJsonIgnoreError(mqMessage), msgErr)
		}
	})
}
