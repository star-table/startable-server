package callsvc

import (
	sdk_const "gitea.bjx.cloud/allstar/platform-sdk/consts"
	"github.com/star-table/startable-server/common/core/errs"
	"github.com/star-table/startable-server/common/core/util/json"
)

type UserAddHandler struct {
	base.CallBackBase
}

type UserAddReq struct {
	Ts    string `json:"ts"`
	Uuid  string `json:"uuid"`
	Token string `json:"token"`
	Type  string `json:"type"`

	Event UserAddReqData `json:"event"`
}

type UserAddReqData struct {
	AppId     string `json:"app_id"`
	TenantKey string `json:"tenant_key"`
	Type      string `json:"type"`

	OpenId     string `json:"open_id"`
	EmployeeId string `json:"employee_id"`
}

// 飞书处理
func (u UserAddHandler) Handle(data string) (string, errs.SystemErrorInfo) {
	userAddReq := &UserAddReq{}
	_ = json.FromJson(data, userAddReq)
	log.Infof("用户入职推送 %s", data)
	tenantKey := userAddReq.Event.TenantKey
	openId := userAddReq.Event.OpenId

	log.Infof("tenantKey %s, openId %s", tenantKey, openId)

	err := u.AddFsUser(tenantKey, openId)
	if err != nil {
		log.Error(err)
		return "err", err
	}

	return "ok", nil
}

func (u UserAddHandler) AddFsUser(tenantKey string, openId string) errs.SystemErrorInfo {
	err := u.UserAdd(sdk_const.SourceChannelFeishu, tenantKey, openId)
	if err != nil {
		return err
	}

	return nil
}
