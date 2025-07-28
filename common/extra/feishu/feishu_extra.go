package feishu

import (
	"fmt"

	sdk_const "gitea.bjx.cloud/allstar/platform-sdk/consts"

	"gitea.bjx.cloud/allstar/feishu-sdk-golang/core/model/vo"
	"gitea.bjx.cloud/allstar/feishu-sdk-golang/sdk"
	platform_sdk "gitea.bjx.cloud/allstar/platform-sdk"
	"github.com/star-table/startable-server/common/core/consts"
	"github.com/star-table/startable-server/common/core/errs"
	"github.com/star-table/startable-server/common/core/logger"
	"github.com/star-table/startable-server/common/core/util/queue"
	"github.com/star-table/startable-server/common/core/util/slice"
)

var log = logger.GetDefaultLogger()

func GetTenant(tenantKey string) (*sdk.Tenant, errs.SystemErrorInfo) {
	if tenantKey == "" {
		return nil, errs.FeiShuClientTenantError
	}
	client, err := platform_sdk.GetClient(sdk_const.SourceChannelFeishu, tenantKey)
	if err != nil {
		log.Errorf("[GetTenant] GetClient corpId:%v, err:%v", tenantKey, err)
		return nil, errs.BuildSystemErrorInfo(errs.PlatFormOpenApiCallError, err)
	}

	return client.GetOriginClient().(*sdk.Tenant), nil
}

/*
*
获取部门列表
*/
func GetDeptList(tenantKey string) ([]vo.DepartmentRestInfoVo, errs.SystemErrorInfo) {
	return GetDeptListWithDepId(tenantKey, "0")
}

func GetDeptListWithDepId(tenantKey string, depId string) ([]vo.DepartmentRestInfoVo, errs.SystemErrorInfo) {
	var respVo *vo.GetDepartmentSimpleListV2RespVo = nil
	var err1 error = nil
	var hasMore = true

	var pageSize = 100
	var pageToken = ""

	depList := make([]vo.DepartmentRestInfoVo, 0)
	for {
		if !hasMore {
			break
		}
		tenant, err := GetTenant(tenantKey)
		if err != nil {
			log.Error(err)
			return nil, err
		}
		respVo, err1 = tenant.GetDepartmentSimpleListV2(depId, pageToken, pageSize, true)
		if err1 != nil {
			log.Error(err1)
			return nil, errs.BuildSystemErrorInfo(errs.FeiShuOpenApiCallError)
		}
		if respVo.Code != 0 {
			log.Error(respVo.Msg)
			return nil, errs.BuildSystemErrorInfoWithMessage(errs.FeiShuOpenApiCallError, respVo.Msg)
		}

		hasMore = respVo.Data.HasMore
		pageToken = respVo.Data.PageToken

		depList = append(depList, respVo.Data.DepartmentInfos...)
	}

	return depList, nil
}

func GetScopeDeps(tenantKey string) ([]vo.DepartmentRestInfoVo, errs.SystemErrorInfo) {
	tenant, err := GetTenant(tenantKey)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	resp, scopeErr := tenant.GetScopeV2()
	if scopeErr != nil {
		log.Error(scopeErr)
		return nil, errs.BuildSystemErrorInfo(errs.FeiShuOpenApiCallError)
	}
	if resp.Code != 0 {
		log.Error(resp.Msg)
		return nil, errs.BuildSystemErrorInfo(errs.FeiShuOpenApiCallError)
	}

	depList := make([]vo.DepartmentRestInfoVo, 0)

	deps := resp.Data.AuthedDepartments
	if deps != nil && len(deps) > 0 {
		if deps[0] == "0" {
			depList, err = GetDeptList(tenantKey)
			if err != nil {
				log.Error(err)
				return nil, errs.BuildSystemErrorInfo(errs.FeiShuOpenApiCallError)
			}
			if resp.Code != 0 {
				log.Error(resp.Msg)
				return nil, errs.BuildSystemErrorInfo(errs.FeiShuOpenApiCallError)
			}
		} else {
			//父部门
			parentDeptList, err := tenant.GetDepartmentInfoBatch(deps)
			if err != nil {
				log.Error(err)
				return nil, errs.BuildSystemErrorInfo(errs.FeiShuOpenApiCallError)
			}
			if parentDeptList.Code != 0 {
				log.Error(resp.Msg)
				return nil, errs.BuildSystemErrorInfo(errs.FeiShuOpenApiCallError)
			}
			for _, depDetail := range parentDeptList.Data.Departments {
				depList = append(depList, vo.DepartmentRestInfoVo{
					Id:       depDetail.Id,
					Name:     depDetail.Name,
					ParentId: depDetail.ParentId,
				})
			}

			for _, dep := range deps {
				depDetails, err := GetDeptListWithDepId(tenantKey, dep)
				if err != nil {
					log.Error(err)
					return nil, errs.BuildSystemErrorInfo(errs.FeiShuOpenApiCallError)
				}
				if resp.Code != 0 {
					log.Error(resp.Msg)
					return nil, errs.BuildSystemErrorInfo(errs.FeiShuOpenApiCallError)
				}
				if depDetails != nil && len(depDetails) > 0 {
					for _, depDetail := range depDetails {
						depList = append(depList, vo.DepartmentRestInfoVo{
							Id:       depDetail.Id,
							Name:     depDetail.Name,
							ParentId: depDetail.ParentId,
						})
					}
				}
			}
		}
	}

	//做去重
	depMap := map[string]vo.DepartmentRestInfoVo{}
	for _, dep := range depList {
		depMap[dep.Id] = dep
	}
	depList = make([]vo.DepartmentRestInfoVo, 0)
	for _, dep := range depMap {
		depList = append(depList, dep)
	}

	//增加根部门
	//depList = append(depList, vo.DepartmentRestInfoVo{
	//	Id: "0",
	//	Name: "飞书平台组织",
	//	ParentId: "Root_Department_Identification",
	//})
	return depList, nil
}

// isNeedLimit 目前主要用于组织初始化，只初始化前面的10000用户
func GetDeptUserInfosByDeptIds(tenantKey string, deptIds []string, isNeedLimit bool) ([]vo.UserRestInfoVo, errs.SystemErrorInfo) {
	userInfos := make([]vo.UserRestInfoVo, 0)
	if deptIds == nil || len(deptIds) == 0 {
		return userInfos, nil
	}

	interrupted := false

	q := queue.GetQueue()
	for _, deptId := range deptIds {
		q.Push(deptId)
	}
	for {
		obj, err := q.Pop()
		if err != nil {
			break
		}
		deptId := obj.(string)

		var hasMore = true
		var pageSize = 100
		var pageToken = ""
		for {
			if !hasMore {
				break
			}
			tenant, err := GetTenant(tenantKey)
			if err != nil {
				log.Error(err)
				return nil, err
			}
			respVo, err1 := tenant.GetDepartmentUserListV2(deptId, pageToken, pageSize, true)
			if err1 != nil {
				log.Error(err1)
				continue
			}
			if respVo.Code == 40162 {
				log.Errorf("tenant.GetDepartmentUserListV2 err: %s", respVo.Msg)
				childs, err := GetDeptListWithDepId(tenantKey, deptId)
				if err != nil {
					log.Error(err)
					return nil, err
				}
				if len(childs) > 0 {
					for _, child := range childs {
						q.Push(child.Id)
					}
				}
				break
			}
			hasMore = respVo.Data.HasMore
			pageToken = respVo.Data.PageToken
			userList := respVo.Data.UserList

			if userList != nil && len(userList) > 0 {
				userInfos = append(userInfos, userList...)
			}

			if isNeedLimit && len(userInfos) > 10000 {
				//加了限制，只取前10000条处理
				interrupted = true
				break
			}
		}
		if interrupted {
			break
		}
	}

	return userInfos, nil
}

func GetScopeOpenIds(tenantKey string) ([]string, errs.SystemErrorInfo) {
	tenant, err := GetTenant(tenantKey)
	if err != nil {
		log.Error(err)
		return nil, err
	}

	resp, scopeErr := tenant.GetScopeV2()
	if scopeErr != nil {
		log.Error(scopeErr)
		return nil, errs.BuildSystemErrorInfo(errs.FeiShuOpenApiCallError)
	}
	if resp.Code != 0 {
		log.Error(resp.Msg)
		return nil, errs.BuildSystemErrorInfo(errs.FeiShuOpenApiCallError)
	}

	//获取当前授权范围内的所有用户
	scopeUsers := resp.Data.AuthedUsers
	scopeOpenIds := make([]string, 0)
	if scopeUsers != nil && len(scopeUsers) > 0 {
		for _, scopeUser := range scopeUsers {
			scopeOpenIds = append(scopeOpenIds, scopeUser.OpenId)
		}
	}

	fsUserInfos, err := GetDeptUserInfosByDeptIds(tenantKey, resp.Data.AuthedDepartments, false)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	for _, fsUserInfo := range fsUserInfos {
		scopeOpenIds = append(scopeOpenIds, fsUserInfo.OpenId)
	}
	scopeOpenIds = slice.SliceUniqueString(scopeOpenIds)
	fmt.Println(len(scopeOpenIds))
	return scopeOpenIds, nil
}

func GetScopeOpenIdsLimit(tenantKey string) ([]string, errs.SystemErrorInfo) {
	tenant, err := GetTenant(tenantKey)
	if err != nil {
		log.Error(err)
		return nil, err
	}

	resp, scopeErr := tenant.GetScopeV2()
	if scopeErr != nil {
		log.Error(scopeErr)
		return nil, errs.BuildSystemErrorInfo(errs.FeiShuOpenApiCallError)
	}
	if resp.Code != 0 {
		log.Error(resp.Msg)
		return nil, errs.BuildSystemErrorInfo(errs.FeiShuOpenApiCallError)
	}

	//获取当前授权范围内的所有用户
	scopeUsers := resp.Data.AuthedUsers
	scopeOpenIds := make([]string, 0)
	if scopeUsers != nil && len(scopeUsers) > 0 {
		for _, scopeUser := range scopeUsers {
			scopeOpenIds = append(scopeOpenIds, scopeUser.OpenId)
		}
	}

	fsUserInfos, err := GetDeptUserInfosByDeptIds(tenantKey, resp.Data.AuthedDepartments, true)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	for _, fsUserInfo := range fsUserInfos {
		scopeOpenIds = append(scopeOpenIds, fsUserInfo.OpenId)
	}

	//获取管理员用户，没获取到就不处理，不用中断
	adminUserResp, adminUserErr := tenant.AdminUserList()
	if adminUserErr != nil {
		log.Error(adminUserErr)
	} else {
		if adminUserResp.Code != 0 {
			log.Error(adminUserResp)
		} else {
			for _, data := range adminUserResp.Data.UserList {
				if data.OpenId != "" {
					scopeOpenIds = append(scopeOpenIds, data.OpenId)
				}
			}
		}
	}

	scopeOpenIds = slice.SliceUniqueString(scopeOpenIds)

	return scopeOpenIds, nil
}

// 获取应用（极星）的授权状态，检查是否有需要的权限。
func GetAppScopes(tenantKey string) (*vo.GetScopesResp, errs.SystemErrorInfo) {
	tenant, err := GetTenant(tenantKey)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	resp, respErr := tenant.GetScopes()
	if respErr != nil {
		log.Error(respErr)
		return nil, errs.FeiShuOpenApiCallError
	}
	return resp, nil
}

// 判断是否有高级权限(仅用于判断通讯录权限)
func JudgeIsInAppScopes(tenantKey string) (bool, errs.SystemErrorInfo) {
	resp, err := GetAppScopes(tenantKey)
	if err != nil {
		return false, err
	}

	for _, item := range resp.Data.Scopes {
		//if item.ScopeName == scopeName && item.GrantStatus == 1 {
		//	return true, nil
		//}
		if ok, _ := slice.Contain([]string{consts.ScopeNameContactAccessAsApp, consts.ScopeNameContactReadOnly, consts.ScopeNameContactReadOnlyAsApp}, item.ScopeName); ok {
			if item.GrantStatus == 1 {
				return true, nil
			}
		}
		//if item.GrantStatus == 2 {
		//	return false, nil
		//}
	}
	return false, nil
}

func GetTenantInfo(tenantKey string) (*vo.OrgInfoData, errs.SystemErrorInfo) {
	tenant, err := GetTenant(tenantKey)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	resp, orgInfoErr := tenant.OrgInfo()
	if orgInfoErr != nil {
		log.Error(orgInfoErr)
		return nil, errs.FeiShuOpenApiCallError
	}

	if resp.Code != 0 {
		log.Errorf("[GetTenantInfo] failed: code:%d, message:%s, tenantKey:%s", resp.Code, resp.Msg, tenantKey)
		if resp.Code == 1184001 || resp.Code == 99991672 {
			return nil, errs.FeiShuGetInfoNoPermission
		} else {
			return nil, errs.FeiShuOpenApiCallError
		}
	}

	return &resp.Data, nil
}
