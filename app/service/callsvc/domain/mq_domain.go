package callsvc

import (
	"github.com/star-table/startable-server/app/facade/msgfacade"
	"github.com/star-table/startable-server/common/core/config"
	"github.com/star-table/startable-server/common/core/logger"
	"github.com/star-table/startable-server/common/core/model"
	"github.com/star-table/startable-server/common/core/util/json"
	"github.com/star-table/startable-server/common/core/util/uuid"
	"github.com/star-table/startable-server/common/model/bo"
)

var log = logger.GetDefaultLogger()

func PushOrgMemberChange(orgMemberChangeBo bo.OrgMemberChangeBo) {
	message, err := json.ToJson(orgMemberChangeBo)
	if err != nil {
		log.Error(err)
	}
	mqMessage := &model.MqMessage{
		Topic:          config.GetMqOrgMemberChangeConfig().Topic,
		Keys:           uuid.NewUuid(),
		Body:           message,
		DelayTimeLevel: 3,
	}
	msgErr := msgfacade.PushMsgToMqRelaxed(*mqMessage, 0, orgMemberChangeBo.OrgId)
	if msgErr != nil {
		log.Errorf("mq消息推送失败，入表失败，消息内容：%s, 失败信息: %v", json.ToJsonIgnoreError(mqMessage), msgErr)
	}
}

func PushFeiShuCallBack(fsCallBack bo.FeiShuCallBackMqBo) {
	message, err := json.ToJson(fsCallBack)
	if err != nil {
		log.Error(err)
	}
	mqMessage := &model.MqMessage{
		Topic:          config.GetMqFeiShuCallBackTopicConfig().Topic,
		Keys:           uuid.NewUuid(),
		Body:           message,
		DelayTimeLevel: 3,
	}
	msgErr := msgfacade.PushMsgToMqRelaxed(*mqMessage, 0, 0)
	if msgErr != nil {
		log.Errorf("mq消息推送失败，入表失败，消息内容：%s, 失败信息: %v", json.ToJsonIgnoreError(mqMessage), msgErr)
	}
}

func PushFeiShuCallBackMsg(fsCallBack bo.FeiShuCallBackMqBo) {
	message, err := json.ToJson(fsCallBack)
	if err != nil {
		log.Error(err)
	}
	mqMessage := &model.MqMessage{
		Topic:          config.GetFeishuCallBackMsg().Topic,
		Keys:           uuid.NewUuid(),
		Body:           message,
		DelayTimeLevel: 3,
	}
	msgErr := msgfacade.PushMsgToMqRelaxed(*mqMessage, 0, 0)
	if msgErr != nil {
		log.Errorf("mq消息推送失败，入表失败，消息内容：%s, 失败信息: %v", json.ToJsonIgnoreError(mqMessage), msgErr)
	}
}

func PushFeiShuCallBackOrder(fsCallBack bo.FeiShuCallBackMqBo) {
	message, err := json.ToJson(fsCallBack)
	if err != nil {
		log.Error(err)
	}
	mqMessage := &model.MqMessage{
		Topic:          config.GetFeishuCallBackOrder().Topic,
		Keys:           uuid.NewUuid(),
		Body:           message,
		DelayTimeLevel: 3,
	}
	msgErr := msgfacade.PushMsgToMqRelaxed(*mqMessage, 0, 0)
	if msgErr != nil {
		log.Errorf("mq消息推送失败，入表失败，消息内容：%s, 失败信息: %v", json.ToJsonIgnoreError(mqMessage), msgErr)
	}
}

func PushFeiShuCallBackUserChange(fsCallBack bo.FeiShuCallBackMqBo) {
	message, err := json.ToJson(fsCallBack)
	if err != nil {
		log.Error(err)
	}
	mqMessage := &model.MqMessage{
		Topic:          config.GetFeishuCallBackUserChange().Topic,
		Keys:           uuid.NewUuid(),
		Body:           message,
		DelayTimeLevel: 3,
	}
	msgErr := msgfacade.PushMsgToMqRelaxed(*mqMessage, 0, 0)
	if msgErr != nil {
		log.Errorf("mq消息推送失败，入表失败，消息内容：%s, 失败信息: %v", json.ToJsonIgnoreError(mqMessage), msgErr)
	}
}

func PushFeiShuCallBackContactScopeChange(fsCallBack bo.FeiShuCallBackMqBo) {
	message, err := json.ToJson(fsCallBack)
	if err != nil {
		log.Error(err)
	}
	mqMessage := &model.MqMessage{
		Topic:          config.GetFeishuCallBackContactScopeChange().Topic,
		Keys:           uuid.NewUuid(),
		Body:           message,
		DelayTimeLevel: 3,
	}
	msgErr := msgfacade.PushMsgToMqRelaxed(*mqMessage, 0, 0)
	if msgErr != nil {
		log.Errorf("mq消息推送失败，入表失败，消息内容：%s, 失败信息: %v", json.ToJsonIgnoreError(mqMessage), msgErr)
	}
}
