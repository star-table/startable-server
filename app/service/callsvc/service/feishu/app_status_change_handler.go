package callsvc

import (
	"github.com/star-table/startable-server/common/core/errs"
	"github.com/star-table/startable-server/common/core/util/json"
)

type AppStatusChangeHandler struct{}

type AppStatusChangeReq struct {
	Ts    string `json:"ts"`
	Uuid  string `json:"uuid"`
	Token string `json:"token"`
	Type  string `json:"type"`

	Event AppStatusChangeReqData `json:"event"`
}

type AppStatusChangeReqData struct {
	AppId     string `json:"app_id"`
	TenantKey string `json:"tenant_key"`
	Type      string `json:"type"`
	// 应用状态 start_by_tenant: 租户启用; stop_by_tenant: 租户停用; stop_by_platform: 平台停用
	// https://open.feishu.cn/document/uAjLw4CM/ukTMukTMukTM/application-v6/event/app-enabled-or-disabled
	Status string `json:"status"`
	// 仅status=start_by_tenant时有此字段
	Operator AppStatusChangeEventOperator `json:"operator"`
}

type AppStatusChangeEventOperator struct {
	OpenId  string `json:"open_id"`
	UserId  string `json:"user_id"` // 仅自建应用才会返回
	UnionId string `json:"union_id"`
}

const StartByTenant = "start_by_tenant"

// 飞书处理
func (AppStatusChangeHandler) Handle(data string) (string, errs.SystemErrorInfo) {
	req := &AppStatusChangeReq{}
	_ = json.FromJson(data, req)

	eventJson, err := json.ToJson(data)
	log.Infof("飞书应用启动/停用通知 eventJson: %s, err: %v", eventJson, err)

	tenantKey := req.Event.TenantKey
	if req.Event.Status == StartByTenant {
		initParam := InstallInfo{
			InstallerOpenId: req.Event.Operator.OpenId,
			IsAdmin:         true,
		}
		err := FsInit(tenantKey, initParam)
		if err != nil {
			log.Error(err)
			return "", err
		}
	}

	return "ok", nil
}
