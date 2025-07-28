package callsvc

import (
	"fmt"

	platform_sdk "gitea.bjx.cloud/allstar/platform-sdk"
	sdk_const "gitea.bjx.cloud/allstar/platform-sdk/consts"
	"github.com/star-table/startable-server/common/core/util/json"
	"github.com/star-table/startable-server/common/core/util/strs"
)

//var log = logger.GetDefaultLogger()

type SuiteTicketHandler struct{}

type SuiteTicketCallbackMsg struct {
	SyncAction  string `json:"syncAction"`
	SuiteTicket string `json:"suiteTicket"`
}

func (SuiteTicketHandler) Handle(data req_vo.DingEventBizData) error {
	msg := &SuiteTicketCallbackMsg{}
	json.FromJson(data.BizData, msg)

	err := platform_sdk.GetPlatformInfo(sdk_const.SourceChannelDingTalk).Cache.SetTicket(msg.SuiteTicket)
	if err != nil {
		log.Error(strs.ObjectToString(err))
		return err
	}
	fmt.Println("callback -> ", data)
	return nil
}
