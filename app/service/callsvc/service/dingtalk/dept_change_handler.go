package callsvc

import (
	"fmt"

	sdk_const "gitea.bjx.cloud/allstar/platform-sdk/consts"

	"github.com/star-table/startable-server/common/core/errs"
	"github.com/star-table/startable-server/common/core/util/json"
)

type DeptChangeHandler struct {
	base.CallBackBase
}

type DeptChangeReq struct {
	Ts    string `json:"ts"`
	Uuid  string `json:"uuid"`
	Token string `json:"token"`
	Type  string `json:"type"`

	Event DeptChangeReqData `json:"event"`
}

type DeptChangeReqData struct {
	AppId            string         `json:"app_id"`
	TenantKey        string         `json:"tenant_key"`
	Type             string         `json:"type"`
	OpenDepartmentId string         `json:"open_department_id"`
	Department       DepartmentInfo `json:"department"`
}

type DepartmentInfo struct {
	CustomId string `json:"custom_id"`
	OpenId   string `json:"open_id"`
}

// 飞书处理
func (d DeptChangeHandler) Handle(data string) (string, errs.SystemErrorInfo) {
	deptReq := &DeptChangeReq{}
	_ = json.FromJson(data, deptReq)
	log.Infof("部门变更 %s", data)
	//tenantKey := deptReq.Header.TenantKey

	err := d.HandleDeptChange(sdk_const.SourceChannelFeishu, deptReq.Event.TenantKey, deptReq.Event.Type,
		deptReq.Event.Department.CustomId)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("部门变更操作成功！！类型：%s", deptReq.Event.Type), nil
}
