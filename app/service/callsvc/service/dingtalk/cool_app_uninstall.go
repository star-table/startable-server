package callsvc

import (
	"github.com/star-table/startable-server/common/core/util/json"
	"github.com/star-table/startable-server/common/model/vo/orgvo"
)

type CoolAppUninstallHandler struct {
}

func (d CoolAppUninstallHandler) Handle(data req_vo.DingEventBizData) error {
	info := &CoolAppInstallInfo{}
	errJson := json.FromJson(data.BizData, info)
	if errJson != nil {
		log.Errorf("[CoolAppUninstallHandler] handle json err:%v", errJson)
		return errJson
	}

	err := orgfacade.DeleteCoolApp(orgvo.DeleteCoolAppReq{Input: &orgvo.CoolAppBaseData{OpenConversationId: info.OpenConversationId}})

	if err.Failure() {
		log.Errorf("[CoolAppUninstallHandler] DeleteCoolApp err:%v", err.Error())
	}

	return nil
}
