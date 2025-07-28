package callsvc

import (
	platform_sdk "gitea.bjx.cloud/allstar/platform-sdk"
	sdk_const "gitea.bjx.cloud/allstar/platform-sdk/consts"
	"github.com/star-table/startable-server/common/core/errs"
	"github.com/star-table/startable-server/common/core/util/json"
)

type AppTicketHandler struct{}

type AppTicketReq struct {
	Ts    string `json:"ts"`
	Uuid  string `json:"uuid"`
	Token string `json:"token"`
	Type  string `json:"type"`

	Event AppTicketReqData `json:"event"`
}

type AppTicketReqData struct {
	AppId     string `json:"app_id"`
	AppTicket string `json:"app_ticket"`
	Type      string `json:"type"`
}

func (AppTicketHandler) Handle(data string) (string, errs.SystemErrorInfo) {
	appTicketReq := &AppTicketReq{}
	_ = json.FromJson(data, appTicketReq)

	appId := appTicketReq.Event.AppId
	appTicket := appTicketReq.Event.AppTicket

	err := platform_sdk.GetPlatformInfo(sdk_const.SourceChannelFeishu).Cache.SetTicket(appTicket)
	if err != nil {
		log.Errorf("appTicket 保存失败 appId: %s, appTicket %s, err %v", appId, appTicket, err)
		return "", nil
	}
	return "ok", nil
}
