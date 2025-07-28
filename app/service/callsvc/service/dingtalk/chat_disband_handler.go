package callsvc

import (
	"gitea.bjx.cloud/allstar/polaris-backend/facade/orgfacade"
	"github.com/star-table/startable-server/common/core/errs"
	"github.com/star-table/startable-server/common/core/util/json"
	"github.com/star-table/startable-server/common/model/vo/orgvo"
	"github.com/star-table/startable-server/common/model/vo/projectvo"
)

type ChatDisbandHandler struct{}

type ChatDisbandReq struct {
	Ts    string `json:"ts"`
	Uuid  string `json:"uuid"`
	Token string `json:"token"`
	Type  string `json:"type"`

	Event ChatDisbandData `json:"event"`
}

type ChatDisbandData struct {
	AppId     string       `json:"app_id"`
	TenantKey string       `json:"tenant_key"`
	Type      string       `json:"type"`
	ChatId    string       `json:"chat_id"`
	Operator  OperatorInfo `json:"operator"`
}

// Handle 飞书处理
func (ChatDisbandHandler) Handle(data string) (string, errs.SystemErrorInfo) {
	reqInfo := &ChatDisbandReq{}
	_ = json.FromJson(data, reqInfo)
	log.Infof("群解散推送： %s", data)
	tenantKey := reqInfo.Event.TenantKey
	openId := reqInfo.Event.Operator.OpenId

	log.Infof("tenantKey %s, openId %s", tenantKey, openId)

	if tenantKey == "" {
		return "ok", nil
	}
	orgInfo := orgfacade.GetBaseOrgInfoByOutOrgId(orgvo.GetBaseOrgInfoByOutOrgIdReqVo{OutOrgId: tenantKey})
	if orgInfo.Failure() {
		log.Error(orgInfo.Error())
		return "err", orgInfo.Error()
	}
	resp := projectfacade.FsChatDisbandCallback(projectvo.FsChatDisbandCallbackReq{
		ChatId: reqInfo.Event.ChatId,
		OrgId:  orgInfo.BaseOrgInfo.OrgId,
	})
	if resp.Failure() {
		log.Error(resp.Error())
		return "err", resp.Error()
	}

	return "ok", nil
}
