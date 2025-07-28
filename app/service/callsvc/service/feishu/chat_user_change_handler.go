package callsvc

import (
	"github.com/star-table/startable-server/common/core/errs"
)

type ChatUserChangeHandler struct{}

type ChatUserChangeReq struct {
	Ts    string `json:"ts"`
	Uuid  string `json:"uuid"`
	Token string `json:"token"`
	Type  string `json:"type"`

	Event ChatUserChangeData `json:"event"`
}

type ChatUserChangeData struct {
	AppId     string `json:"app_id"`
	TenantKey string `json:"tenant_key"`
	Type      string `json:"type"`

	OpenId     string `json:"open_id"`
	EmployeeId string `json:"employee_id"`
}

// Handle 飞书事件 - 撤销拉用户进群
// https://open.feishu.cn/document/uAjLw4CM/ukTMukTMukTM/reference/im-v1/chat-member-user/events/withdrawn
func (ChatUserChangeHandler) Handle(data string) (string, errs.SystemErrorInfo) {
	//userAddReq := &UserAddReq{}
	//_ = json.FromJson(data, userAddReq)
	//log.Infof("用户入职推送 %s", data)
	//tenantKey := userAddReq.CustomEvent.TenantKey
	//openId := userAddReq.CustomEvent.OpenId
	//
	//log.Infof("tenantKey %s, openId %s", tenantKey, openId)
	//
	//err := addFsUser(tenantKey, openId)
	//if err != nil{
	//	log.Error(err)
	//	return "err", err
	//}

	return "ok", nil
}
