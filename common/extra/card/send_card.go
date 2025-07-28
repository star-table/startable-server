package card

import (
	pushV1 "gitea.bjx.cloud/LessCode/interface/golang/push/v1"
	sdk_const "gitea.bjx.cloud/allstar/platform-sdk/consts"
	"github.com/star-table/startable-server/common/core/config"
	"github.com/star-table/startable-server/common/core/errs"
	"github.com/star-table/startable-server/common/core/model"
	"github.com/star-table/startable-server/common/core/util/json"
	"github.com/star-table/startable-server/common/library/mq"
	"github.com/star-table/startable-server/common/model/vo/commonvo"
	"github.com/spf13/cast"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/anypb"
)

func SendCardMeta(sourceChannel string, corpId string, meta *commonvo.CardMeta, openIds []string) {
	if sourceChannel == sdk_const.SourceChannelFeishu {
		SendCardMetaFeiShu(corpId, meta, openIds)
	} else if sourceChannel == sdk_const.SourceChannelWeixin {
		SendCardMetaWeiXin(corpId, meta, openIds)
	} else if sourceChannel == sdk_const.SourceChannelDingTalk {
		SendCardMetaDingTalk(corpId, meta, openIds)
	}
}

func PushCard(orgId int64, pushCard *commonvo.PushCard) errs.SystemErrorInfo {
	client := *mq.GetMQClient()
	if client == nil {
		log.Errorf("[PushCard] can not get kafka client, msg: %v", json.ToJsonIgnoreError(pushCard))
		return errs.KafkaMqSendMsgError
	}
	topic := config.GetCardPushConfig().Topic
	log.Infof("[PushCard] topic: %s, orgId: %d, msg: %v", topic, orgId, json.ToJsonIgnoreError(pushCard))

	repushTimes := 3
	reconsumeTimes := 0

	pushCardReq := &pushV1.PushCardReq{
		OrgId:         orgId,
		SourceChannel: pushCard.SourceChannel,
		OutOrgId:      pushCard.OutOrgId,
		Card:          pushCard.CardMsg,
		ToOpenIds:     pushCard.OpenIds,
		ToDeptIds:     pushCard.DeptIds,
		ToChatIds:     pushCard.ChatIds,
	}

	any, err := anypb.New(pushCardReq)
	if err != nil {
		return errs.ProtobufMarshalError
	}
	push := &pushV1.PushMessage{
		Type:    pushV1.PushType_Card,
		Payload: any,
	}
	body, err := proto.Marshal(push)
	if err != nil {
		return errs.ProtobufMarshalError
	}

	msg := &model.MqMessage{
		Topic:          topic,
		Keys:           cast.ToString(orgId),
		Body:           string(body),
		RePushTimes:    &repushTimes,
		ReconsumeTimes: &reconsumeTimes,
	}
	_, sysErr := client.PushMessage(msg)
	return sysErr
}
