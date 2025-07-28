package callsvc

import (
	"github.com/star-table/startable-server/common/core/errs"
	"github.com/star-table/startable-server/common/core/util/json"
)

type AddUserToChatHandler struct{}

type ChangeUserFromChatReq struct {
	Ts    string `json:"ts"`
	Uuid  string `json:"uuid"`
	Token string `json:"token"`
	Type  string `json:"type"`

	Event ChangeUserFromChatData `json:"event"`
}

type ChangeUserFromChatData struct {
	AppId     string             `json:"app_id"`
	ChatId    string             `json:"chat_id"`
	Operator  OperatorInfo       `json:"operator"`
	Tenantkey string             `json:"tenant_key"`
	Type      string             `json:"type"`
	Users     []ChangedUsersInfo `json:"users"`
}

type OperatorInfo struct {
	OpenId string `json:"open_id"`
	UserId string `json:"user_id"`
}

type ChangedUsersInfo struct {
	Name   string `json:"name"`
	OpenId string `json:"open_id"`
	UserId string `json:"user_id"`
}

func (AddUserToChatHandler) Handle(data string) (string, errs.SystemErrorInfo) {
	userAddReq := &ChangeUserFromChatReq{}
	_ = json.FromJson(data, userAddReq)
	log.Infof("用户进群出群事件通知 %s", data)
	tenantKey := userAddReq.Event.Tenantkey
	openId := userAddReq.Event.Operator.OpenId

	log.Infof("tenantKey %s, openId %s", tenantKey, openId)

	err := UserAddHandler{}.AddFsUser(tenantKey, openId)
	if err != nil {
		log.Error(err)
		return "err", err
	}

	return "ok", nil
}
