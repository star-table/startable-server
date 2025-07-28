package callsvc

import (
	sdk_const "gitea.bjx.cloud/allstar/platform-sdk/consts"
	"github.com/star-table/startable-server/common/core/errs"
	"github.com/star-table/startable-server/common/core/util/json"
)

type UserUpdateHandler struct {
	base.CallBackBase
}

type UserUpdateReq struct {
	Ts    string `json:"ts"`
	Uuid  string `json:"uuid"`
	Token string `json:"token"`
	Type  string `json:"type"`

	Event UserUpdateReqData `json:"event"`
}

type UserUpdateReqData struct {
	AppId     string `json:"app_id"`
	TenantKey string `json:"tenant_key"`
	Type      string `json:"type"`

	OpenId     string `json:"open_id"`
	EmployeeId string `json:"employee_id"`
}

// 飞书处理
// 由于需求方的要求是部门信息只能手动同步。因此，自动同步时，只同步头像和姓名信息
func (u UserUpdateHandler) Handle(data string) (string, errs.SystemErrorInfo) {
	userLeaveReq := &UserUpdateReq{}
	_ = json.FromJson(data, userLeaveReq)
	log.Infof("用户信息更改推送 %s", data)
	tenantKey := userLeaveReq.Event.TenantKey
	openId := userLeaveReq.Event.OpenId

	log.Infof("event: userUpdate, tenantKey %s, openId %s", tenantKey, openId)
	err := u.UserUpdate(sdk_const.SourceChannelFeishu, tenantKey, openId, "")
	if err != nil {
		return "", err
	}

	return "ok", nil
}
