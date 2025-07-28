package callsvc

import (
	sdk_const "gitea.bjx.cloud/allstar/platform-sdk/consts"
	"github.com/star-table/startable-server/common/core/errs"
	"github.com/star-table/startable-server/common/core/util/json"
)

type UserLeaveHandler struct {
	base.CallBackBase
}

type UserLeaveReq struct {
	Ts    string `json:"ts"`
	Uuid  string `json:"uuid"`
	Token string `json:"token"`
	Type  string `json:"type"`

	Event UserLeaveReqData `json:"event"`
}

type UserLeaveReqData struct {
	AppId     string `json:"app_id"`
	TenantKey string `json:"tenant_key"`
	Type      string `json:"type"`

	OpenId     string `json:"open_id"`
	EmployeeId string `json:"employee_id"`
}

// 飞书处理
func (u UserLeaveHandler) Handle(data string) (string, errs.SystemErrorInfo) {
	userLeaveReq := &UserLeaveReq{}
	_ = json.FromJson(data, userLeaveReq)
	log.Infof("用户离职推送 %s", data)
	tenantKey := userLeaveReq.Event.TenantKey
	openId := userLeaveReq.Event.OpenId

	log.Infof("tenantKey %s, openId %s", tenantKey, openId)

	err := u.UserLeave(sdk_const.SourceChannelFeishu, tenantKey, openId)
	if err != nil {
		return "", err
	}

	return "ok", nil
}
