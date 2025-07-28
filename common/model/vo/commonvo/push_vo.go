package commonvo

import (
	pushV1 "gitea.bjx.cloud/LessCode/interface/golang/push/v1"
	"github.com/star-table/startable-server/common/model/vo"
)

type GenerateCardResp struct {
	vo.Err
	Data *pushV1.GenerateCardReply `json:"data"`
}

type PushMqttResp struct {
	vo.Err
	Data *pushV1.PushMqttReply `json:"data"`
}

type GenerateMqttKeyResp struct {
	vo.Err
	Data *pushV1.GenerateMqttKeyReply `json:"data"`
}
