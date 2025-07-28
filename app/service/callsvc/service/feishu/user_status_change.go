package callsvc

import (
	sdk_const "gitea.bjx.cloud/allstar/platform-sdk/consts"
	"github.com/star-table/startable-server/app/facade/orgfacade"
	"github.com/star-table/startable-server/common/core/consts"
	"github.com/star-table/startable-server/common/core/errs"
	"github.com/star-table/startable-server/common/core/util/json"
	"github.com/star-table/startable-server/common/model/bo"
	"github.com/star-table/startable-server/common/model/vo/orgvo"
)

// 用户状态发生改变回调
// 用户状态的改变，目前场景主要是：
// 	* 飞书管理后台，暂停用户、恢复用户

type UserStatusChangeHandler struct{}

type UserStatusChangeReq struct {
	Ts    string `json:"ts"`
	Uuid  string `json:"uuid"`
	Token string `json:"token"`
	Type  string `json:"type"`

	Event UserStatusChangeReqEvent `json:"event"`
}

type UserStatusChangeReqEvent struct {
	AppID         string        `json:"app_id"`
	BeforeStatus  BeforeStatus  `json:"before_status"`
	ChangeTime    string        `json:"change_time"`
	CurrentStatus CurrentStatus `json:"current_status"`
	OpenID        string        `json:"open_id"`
	TenantKey     string        `json:"tenant_key"`
	Type          string        `json:"type"`
	UnionID       string        `json:"union_id"`
}

type BeforeStatus struct {
	IsActive   bool `json:"is_active"`
	IsFrozen   bool `json:"is_frozen"`
	IsResigned bool `json:"is_resigned"`
}

type CurrentStatus struct {
	IsActive   bool `json:"is_active"`
	IsFrozen   bool `json:"is_frozen"`
	IsResigned bool `json:"is_resigned"`
}

func (UserStatusChangeHandler) Handle(data string) (string, errs.SystemErrorInfo) {
	eventBody := &UserStatusChangeReq{}
	_ = json.FromJson(data, eventBody)
	tenantKey := eventBody.Event.TenantKey
	openId := eventBody.Event.OpenID
	log.Infof("管理员在企业管理后台变更用户状态，tenantKey: %s, openId: %s, reqBody: %s", tenantKey, openId, data)

	return HandleUserStatusChange(&eventBody.Event)
}

func HandleUserStatusChange(reqEvent *UserStatusChangeReqEvent) (string, errs.SystemErrorInfo) {
	//获取组织信息
	orgInfoResp := orgfacade.GetBaseOrgInfoByOutOrgId(orgvo.GetBaseOrgInfoByOutOrgIdReqVo{OutOrgId: reqEvent.TenantKey})
	if orgInfoResp.Failure() {
		log.Error(orgInfoResp.Message)
		return "err", orgInfoResp.Error()
	}
	orgId := orgInfoResp.BaseOrgInfo.OrgId
	//获取用户信息
	userInfoResp := orgfacade.GetBaseUserInfoByEmpId(orgvo.GetBaseUserInfoByEmpIdReqVo{
		OrgId: orgId,
		EmpId: reqEvent.OpenID,
	})
	if userInfoResp.Failure() {
		log.Error(userInfoResp.Message)
		return "err", userInfoResp.Error()
	}

	changeType := 0
	if reqEvent.CurrentStatus.IsFrozen != reqEvent.BeforeStatus.IsFrozen {
		if reqEvent.CurrentStatus.IsFrozen {
			changeType = consts.OrgMemberChangeTypeDisable
		} else {
			changeType = consts.OrgMemberChangeTypeEnable
		}
	}
	if changeType == 0 {
		return "err", errs.BuildSystemErrorInfoWithMessage(errs.FeiShuEventNotSupport, "不支持的用户状态变更类型。")
	}
	domain.PushOrgMemberChange(bo.OrgMemberChangeBo{
		ChangeType:    changeType,
		OrgId:         orgId,
		UserId:        userInfoResp.BaseUserInfo.UserId,
		OpenId:        reqEvent.OpenID,
		SourceChannel: sdk_const.SourceChannelFeishu,
	})

	return "ok", nil
}
