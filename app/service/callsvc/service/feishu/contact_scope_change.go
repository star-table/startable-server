package callsvc

import (
	"time"

	sdk_const "gitea.bjx.cloud/allstar/platform-sdk/consts"
	"github.com/star-table/startable-server/app/facade/orgfacade"
	"github.com/star-table/startable-server/common/core/consts"
	"github.com/star-table/startable-server/common/core/errs"
	"github.com/star-table/startable-server/common/core/util/json"
	"github.com/star-table/startable-server/common/core/util/slice"
	"github.com/star-table/startable-server/common/extra/feishu"
	"github.com/star-table/startable-server/common/model/bo"
	"github.com/star-table/startable-server/common/model/vo"
	"github.com/star-table/startable-server/common/model/vo/orgvo"
)

// https://open.feishu.cn/document/ukTMukTMukTM/uETNz4SM1MjLxUzM//event/scope-change
type ContactScopeChangeHandler struct{}

type ContactScopeChangeReq struct {
	Ts    string `json:"ts"`
	Uuid  string `json:"uuid"`
	Token string `json:"token"`
	Type  string `json:"type"`

	Event ContactScopeChangeReqData `json:"event"`
}

type ContactScopeChangeReqData struct {
	AppId     string `json:"app_id"`
	TenantKey string `json:"tenant_key"`
	Type      string `json:"type"`
}

// 飞书处理
func (ContactScopeChangeHandler) Handle(data string) (string, errs.SystemErrorInfo) {
	contactScopeChangeReq := &ContactScopeChangeReq{}
	_ = json.FromJson(data, contactScopeChangeReq)
	log.Infof("授权范围变更事件1.0版本, 管理员在企业管理后台变更权限范围 %s", data)
	tenantKey := contactScopeChangeReq.Event.TenantKey
	log.Infof("tenantKey %s", tenantKey)

	err := HandleDeptScopeSync(tenantKey)
	if err != nil {
		log.Error(err)
		return "error", err
	}

	return "ok", nil
}

var ignoreOrgIds = []int64{
	//16317,
}

// 处理部门授权范围同步
func HandleDeptScopeSync(tenantKey string) errs.SystemErrorInfo {
	defer func() {
		if err := feishu.ClearFsScopeCache(tenantKey); err != nil {
			log.Error(err)
		}
	}()

	//获取组织信息
	orgInfoResp := orgfacade.GetBaseOrgInfoByOutOrgId(orgvo.GetBaseOrgInfoByOutOrgIdReqVo{OutOrgId: tenantKey})
	if orgInfoResp.Failure() {
		log.Error(orgInfoResp.Message)
		return orgInfoResp.Error()
	}
	orgId := orgInfoResp.BaseOrgInfo.OrgId
	if ok, _ := slice.Contain(ignoreOrgIds, orgId); ok {
		log.Infof("[通讯录变更]当前组织暂时不处理，组织id:%v", orgId)
		return nil
	}

	scopeOpenIds, err1 := feishu.GetScopeOpenIds(tenantKey)
	if err1 != nil {
		log.Error(err1)
		return err1
	}

	//获取组织下所有的用户(目测量不会很大，且单条记录数据量也不大，这里直接全量查)
	userInfoResp := orgfacade.GetOrgUserInfoListBySourceChannel(orgvo.GetOrgUserInfoListBySourceChannelReq{
		SourceChannel: sdk_const.SourceChannelFeishu,
		OrgId:         orgId,
		Page:          -1,
		Size:          -1,
	})
	if userInfoResp.Failure() {
		log.Error(userInfoResp.Message)
		return userInfoResp.Error()
	}

	scopeOpenIdMap := map[string]bool{}
	for _, scopeOpenId := range scopeOpenIds {
		scopeOpenIdMap[scopeOpenId] = true
	}

	log.Infof("授权范围变动 %s", json.ToJsonIgnoreError(scopeOpenIdMap))

	//获取之前所有的用户
	userOpenIdMap := map[string]bo.OrgUserInfo{}
	for _, userInfo := range userInfoResp.Data.List {
		userOpenIdMap[userInfo.OutUserId] = userInfo
	}

	//要禁用的用户id
	disableUserIds := make([]int64, 0)
	//要启用的用户id
	enableUserIds := make([]int64, 0)
	//要新增的用户openId
	addOpenIds := make([]string, 0)

	for scopeOpenId, _ := range scopeOpenIdMap {
		if scopeOpenId == "" {
			continue
		}
		userInfo, ok := userOpenIdMap[scopeOpenId]
		if !ok {
			addOpenIds = append(addOpenIds, scopeOpenId)
		} else {
			if userInfo.OrgUserStatus == consts.AppStatusDisabled || userInfo.OrgUserStatus == consts.AppStatusHidden {
				enableUserIds = append(enableUserIds, userInfo.UserId)
			}
		}
	}

	for userOpenId, userInfo := range userOpenIdMap {
		if _, ok := scopeOpenIdMap[userOpenId]; !ok {
			disableUserIds = append(disableUserIds, userInfo.UserId)
		}
	}

	log.Infof("禁用/隐藏的用户 %s", json.ToJsonIgnoreError(disableUserIds))
	log.Infof("启用的用户 %s", json.ToJsonIgnoreError(enableUserIds))
	log.Infof("新增的用户 %s", json.ToJsonIgnoreError(addOpenIds))

	if len(disableUserIds) > 0 {
		for _, uid := range disableUserIds {
			go domain.PushOrgMemberChange(bo.OrgMemberChangeBo{
				OrgId:         orgId,
				ChangeType:    consts.OrgMemberChangeTypeDisable,
				UserId:        uid,
				SourceChannel: sdk_const.SourceChannelFeishu,
			})
			time.Sleep(10 * time.Millisecond)
		}
	}

	if len(enableUserIds) > 0 {
		for _, uid := range enableUserIds {
			go domain.PushOrgMemberChange(bo.OrgMemberChangeBo{
				OrgId:         orgId,
				ChangeType:    consts.OrgMemberChangeTypeEnable,
				UserId:        uid,
				SourceChannel: sdk_const.SourceChannelFeishu,
			})
			time.Sleep(10 * time.Millisecond)
		}
	}

	if len(addOpenIds) > 0 {
		for _, openId := range addOpenIds {
			go domain.PushOrgMemberChange(bo.OrgMemberChangeBo{
				OrgId:         orgId,
				ChangeType:    consts.OrgMemberChangeTypeAdd,
				OpenId:        openId,
				SourceChannel: sdk_const.SourceChannelFeishu,
			})
			time.Sleep(10 * time.Millisecond)
		}
	}

	// 如果组织拥有者不是超管，则把拥有者设置为超管
	// 考虑到组织拥有者可能已离职，这段逻辑先注释。
	//resp1 := orgfacade.CheckAndSetSuperAdmin(orgvo.CheckAndSetSuperAdminReq{
	//	OrgId:  orgId,
	//	UserId: -1,
	//})
	//if resp1.Failure() {
	//	log.Error(resp1.Error())
	//	return resp1.Error()
	//}

	//同步部门信息
	resp := orgfacade.ChangeDeptScope(vo.CommonReqVo{OrgId: orgId})
	if resp.Failure() {
		log.Error(resp.Error())
		return resp.Error()
	}

	errCache := base.DelPayRangeInfoCache(orgId, tenantKey)
	if errCache != nil {
		log.Errorf("[HandleDeptScopeSync] cache del err:%v, orgId:%v", errCache, orgId)
	}

	return nil
}
