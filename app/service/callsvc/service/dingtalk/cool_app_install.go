package callsvc

import (
	"github.com/star-table/startable-server/common/core/util/json"
	"github.com/star-table/startable-server/common/model/vo/orgvo"
)

type CoolAppInstallInfo struct {
	Operator               string `json:"Operator"`
	OpenConversationCorpId string `json:"OpenConversationCorpId"`
	OperateTime            string `json:"OperateTime"`
	OpenConversationId     string `json:"OpenConversationId"`
	SyncAction             string `json:"syncAction"`
	CoolAppCode            string `json:"CoolAppCode"`
	RobotCode              string `json:"RobotCode"`
	SyncSeq                string `json:"syncSeq"`
}

type CoolAppInstallHandler struct {
}

func (d CoolAppInstallHandler) Handle(data req_vo.DingEventBizData) error {
	info := &CoolAppInstallInfo{}
	errJson := json.FromJson(data.BizData, info)
	if errJson != nil {
		log.Errorf("[CoolAppInstallHandler] handle json err:%v", errJson)
		return errJson
	}

	err := orgfacade.CreateCoolApp(orgvo.CreateCoolAppReq{
		Input: &orgvo.CreateCoolAppData{
			OpenConversationId: info.OpenConversationId,
			CorpId:             data.CorpId,
			RobotCode:          info.RobotCode},
	})

	if err.Failure() {
		log.Errorf("[CoolAppInstallHandler] CreateCoolApp err:%v", err.Error())
	}

	return nil
}
