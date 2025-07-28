package callsvc

import (
	sdk_const "gitea.bjx.cloud/allstar/platform-sdk/consts"

	"github.com/star-table/startable-server/common/core/errs"
	"github.com/star-table/startable-server/common/core/util/json"
)

// https://open.feishu.cn/document/uAjLw4CM/ukTMukTMukTM/reference/contact-v3/scope/events/updated

type ContactScopeChangeV3Handler struct {
	base.CallBackBase
}

type ContactScopeChangeV3Req struct {
	req_vo.FeiShuBaseEventV2
	Event ContactScopeChangeV3Event `json:"event"`
}

type ContactScopeChangeV3Event struct {
	Added   AddedOrRemove `json:"added"`
	Removed AddedOrRemove `json:"removed"`
}

type AddedOrRemove struct {
	Departments []Departments `json:"departments"`
	Users       []Users       `json:"users"`
	UserGroups  []UserGroup   `json:"user_groups"`
}

type Departments struct {
	Name               string           `json:"name"`
	I18NName           req_vo.I18NName  `json:"i18n_name"`
	ParentDepartmentID string           `json:"parent_department_id"`
	DepartmentID       string           `json:"department_id"`
	OpenDepartmentID   string           `json:"open_department_id"`
	LeaderUserID       string           `json:"leader_user_id"`
	ChatID             string           `json:"chat_id"`
	Order              string           `json:"order"`
	UnitIds            []string         `json:"unit_ids"`
	MemberCount        int              `json:"member_count"`
	Status             DepartmentStatus `json:"status"`
}

type Users struct {
	UnionID      string                               `json:"union_id"`
	UserID       string                               `json:"user_id"`
	OpenID       string                               `json:"open_id"`
	Name         string                               `json:"name"`
	EnName       string                               `json:"en_name"`
	Email        string                               `json:"email"`
	Mobile       string                               `json:"mobile"`
	Gender       int                                  `json:"gender"`
	Avatar       req_vo.EventV3ReqEventUserItemAvatar `json:"avatar"`
	Status       UserStatus                           `json:"status"`
	LeaderUserId string                               `json:"leader_user_id"`
	IsFrozen     bool                                 `json:"is_frozen"`
}

type UserGroup struct {
	UserGroupId string `json:"user_group_id"`
	Name        string `json:"name"`
	Type        int    `json:"type"`
	MemberCount int    `json:"member_count"`
	Status      int    `json:"status"`
}

type UserStatus struct {
	IsFrozen    bool `json:"is_frozen"`
	IsResigned  bool `json:"is_resigned"`
	IsActivated bool `json:"is_activated"`
	IsExited    bool `json:"is_exited"`
	IsUnJoin    bool `json:"is_unjoin"`
}

type DepartmentStatus struct {
	IsDeleted bool `json:"is_deleted"`
}

// Handle 飞书处理
func (c ContactScopeChangeV3Handler) Handle(data string) (string, errs.SystemErrorInfo) {
	contactScopeChangeReq := &req_vo.EventV3Req{}
	_ = json.FromJson(data, contactScopeChangeReq)
	log.Infof("管理员在企业管理后台变更权限范围 %s", data)
	tenantKey := contactScopeChangeReq.Header.TenantKey
	log.Infof("tenantKey %s", tenantKey)

	err := c.HandleDeptScopeSyncV3(tenantKey)
	if err != nil {
		log.Error(err)
		return "error", err
	}

	return "ok", nil
}

// HandleDeptScopeSyncV3 处理部门授权范围同步
func (c ContactScopeChangeV3Handler) HandleDeptScopeSyncV3(tenantKey string) errs.SystemErrorInfo {
	return c.ChangeAuthScope(tenantKey, sdk_const.SourceChannelFeishu)
}
