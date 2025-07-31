package orgsvc

import (
	"strings"

	"github.com/star-table/startable-server/app/service/orgsvc/domain"
	"github.com/star-table/startable-server/common/core/consts"
	"github.com/star-table/startable-server/common/core/errs"
	"github.com/star-table/startable-server/common/core/util/json"
	"github.com/star-table/startable-server/common/core/util/slice"
	"github.com/star-table/startable-server/common/core/util/str"
	"github.com/star-table/startable-server/common/model/bo"
	"github.com/star-table/startable-server/common/model/vo"
	"github.com/star-table/startable-server/common/model/vo/orgvo"
	"github.com/star-table/startable-server/common/model/vo/projectvo"
)

// 融合极星-获取个人的项目权限
func GetPersonalPermissionInfoForFuse(orgId, userId int64, projectId, issueId *int64, sourceChannel string) (map[string]interface{}, errs.SystemErrorInfo) {
	var projectAuthInfo *bo.ProjectAuthBo
	result := make(map[string]interface{}, 0)
	groupOpList := make(map[string][]string, 0)
	if projectId != nil && *projectId > 0 {
		// 传入了项目id，则表示查询项目的权限信息
		//获取项目信息
		projectInfo := projectfacade.GetCacheProjectInfo(projectvo.GetCacheProjectInfoReqVo{
			ProjectId: *projectId,
			OrgId:     orgId,
		})
		if projectInfo.Failure() {
			log.Error(projectInfo.Error())
			return nil, projectInfo.Error()
		}
		projectAuthInfo = projectInfo.ProjectCacheBo
		// 如果当前用户是项目所有者，则有权限。
		if projectInfo.ProjectCacheBo.Owner == userId {
			return consts.GetPermissionForPro(), nil
			// isProjectOwner = true
		}
		//获取项目权限项
		optAuthResp := permissionfacade.GetAppAuth(orgId, projectAuthInfo.AppId, 0, userId)
		if optAuthResp.Failure() {
			log.Error(optAuthResp.Message)
			return nil, optAuthResp.Error()
		}
		log.Infof("optAuthResp: %s", json.ToJsonIgnoreError(optAuthResp))
		// ["Permission.Org.ProjectObjectType-Modify", "Permission.Org.ProjectObjectType-Create", "Permission.Org.ProjectObjectType-Delete", "Permission.Pro.Config-View", "Permission.Pro.Config-Modify,Bind,Unbind", "Permission.Pro.Config-Filing,UnFiling", "Permission.Pro.Config-ModifyStatus", "Permission.Pro.Config-ModifyField", "Permission.Pro.Iteration-Modify", "Permission.Pro.Iteration-Create", "Permission.Pro.Iteration-Delete", "Permission.Pro.Iteration-ModifyStatus", "Permission.Pro.Iteration-Bind,Unbind", "Permission.Pro.Issue.4-Modify,Bind,Unbind", "Permission.Pro.Issue.4-Create", "Permission.Pro.Issue.4-Delete", "Permission.Pro.Issue.4-ModifyStatus", "Permission.Pro.Issue.4-Comment", "Permission.Pro.Role-Modify", "Permission.Pro.Role-Create", "Permission.Pro.Role-Delete", "Permission.Pro.Role-ModifyPermission", "Permission.Pro.File-Download", "Permission.Pro.File-Modify", "Permission.Pro.File-Delete", "Permission.Pro.File-CreateFolder", "Permission.Pro.File-ModifyFolder", "Permission.Pro.File-DeleteFolder", "Permission.Pro.Tag-Delete", "Permission.Pro.Tag-Remove", "Permission.Pro.Tag-Modify", "Permission.Pro.Attachment-Download", "Permission.Pro.Attachment-Delete", "Permission.Pro.Member-Bind", "Permission.Pro.Member-Unbind"]
		optAuthArr := optAuthResp.Data.OptAuth
		if (optAuthResp.Data.HasAppRootPermission) ||
			len(optAuthArr) == 1 && optAuthArr[0] == "*" {
			return consts.GetPermissionForPro(), nil
		}
		// 没有角色的特殊判断 注释掉，这部分的逻辑千源接口已经实现了。
		// 由于还存在未返回的情况，因此这里还是做下兜底：没有权限时，返回编辑者权限。
		if len(optAuthArr) < 1 && (optAuthResp.Data.LangCode == nil || *optAuthResp.Data.LangCode == "") {
			// 返回默认的角色权限项：项目成员
			return consts.GetPermissionDefaultOperationForPro(), nil
		}
		for _, item := range optAuthArr {
			infos := strings.Split(item, "-")
			if len(infos) < 2 {
				continue
			}
			code := GetPermissionCodeMap(infos[0])
			ops := strings.Split(infos[1], ",")
			if _, ok := groupOpList[code]; ok {
				groupOpList[code] = append(groupOpList[code], ops...)
			} else {
				groupOpList[code] = ops
			}
		}
		for k, item := range groupOpList {
			result[k] = item
		}

		return result, nil
	} else if (projectId != nil && *projectId == 0) && (issueId != nil && *issueId > 0) {
		// 针对"未放入项目"的项目处理
		// 1.管理员有权限；2.任务相关人有权限。
		//获取组织权限项，判断是否是管理员
		manageAuthInfoResp := userfacade.GetUserAuthority(orgId, userId)
		if manageAuthInfoResp.Failure() {
			log.Error(manageAuthInfoResp.Message)
			return nil, manageAuthInfoResp.Error()
		}
		manageAuthInfo := manageAuthInfoResp.Data
		if manageAuthInfo.IsSysAdmin || manageAuthInfo.IsSubAdmin {
			return consts.GetPermissionForNoProIssue(), nil
		}

		// 是否是任务相关人
		resp := projectfacade.CheckIsIssueRelatedPeople(projectvo.CheckIsIssueMemberReqVo{
			Input: vo.CheckIsIssueMemberReq{
				IssueID: *issueId,
				UserID:  userId,
			},
			UserId: userId,
			OrgId:  orgId,
		})
		if resp.Failure() {
			return result, resp.Error()
		}
		if resp.Data.IsTrue {
			return consts.GetPermissionForNoProIssue(), nil
		}
		return result, nil
	} else if projectId != nil && *projectId == -1 {
		// project 为 -1 表示汇总表。获取**汇总表**的应用权限。
		// 查询组织的汇总表 appId，然后查询其权限数据
		orgBo, infoErr := domain.GetOrgBoById(orgId)
		if infoErr != nil {
			log.Error(infoErr)
			return nil, infoErr
		}
		orgRemarkJson := orgBo.Remark
		orgRemarkObj := &orgvo.OrgRemarkConfigType{}
		if len(orgRemarkJson) > 0 {
			oriErr := json.FromJson(orgRemarkJson, orgRemarkObj)
			if oriErr != nil {
				log.Error(oriErr)
				return nil, errs.BuildSystemErrorInfo(errs.JSONConvertError, oriErr)
			}
		}
		optAuthResp := permissionfacade.GetAppAuth(orgId, orgRemarkObj.OrgSummaryTableAppId, 0, userId)
		if optAuthResp.Failure() {
			log.Error(optAuthResp.Message)
			return nil, optAuthResp.Error()
		}
		log.Infof("summary table app optAuthResp: %s", json.ToJsonIgnoreError(optAuthResp))
		if optAuthResp.Data.HasAppRootPermission {
			return consts.GetPermissionForPro(), nil
		}
		groupOpList = AssemblyOptAuthData(optAuthResp.Data.OptAuth, groupOpList)
		for k, item := range groupOpList {
			result[k] = item
		}

		return result, nil
	} else {
		// 组织管理组-获取组织权限项
		manageAuthInfoResp := userfacade.GetUserAuthority(orgId, userId)
		if manageAuthInfoResp.Failure() {
			log.Error(manageAuthInfoResp.Message)
			return nil, manageAuthInfoResp.Error()
		}
		manageAuthInfo := manageAuthInfoResp.Data
		optAuthArr := manageAuthInfoResp.Data.OptAuth
		if manageAuthInfo.IsSysAdmin || (len(optAuthArr) == 1 && optAuthArr[0] == "*") {
			log.Infof("权限校验成功，用户 %d 是组织 %d 的超级管理员", userId, orgId)
			// 组装成极星需要的 map 结构
			return consts.GetPermissionForOrg(), nil
		}
		// 由于权限被合并为无码的三个权限项，因此需要对这三个特殊的权限项做映射，映射为极星的权限项
		// * “编辑组织设置”、“编辑/审核成员状态”、“邀请成员”、“加入/解除角色成员”、“查看成员列表”统一为“通讯录管理-组织架构” —— `"1"`
		// * “创建角色”、“编辑角色”、“删除角色”统一为“通讯录管理-角色架构” —— `"2"`
		// * “创建项目”、“新增/管理自定义字段”统一为“可添加/删除应用” —— `"3"`
		for _, item := range optAuthArr {
			if isOk, _ := slice.Contain([]string{"1", "2", "3"}, item); isOk {
				tmpGroupOpList := LessCodeOp2PolarisOpArr(item)
				for k, v := range tmpGroupOpList {
					// 如果有值，则与已有的 groupOpList 合并。
					if _, ok := groupOpList[k]; ok {
						groupOpList[k] = append(groupOpList[k], v...)
					} else {
						groupOpList[k] = v
					}
				}
				continue
			}
			infos := strings.Split(item, "-")
			if len(infos) < 2 {
				// 如果没有 `-`，则可能是结果已经是转化好了的（如：xxx.xxx.xxxConfig.Modify），此时获取权限项的前缀部分，去匹配分类
				opCateStr := GetCateOfOperation(item)
				opSuffix := GetOpSuffixOfOperation(item)
				if opCateStr != "" {
					if _, ok := groupOpList[opCateStr]; ok {
						groupOpList[opCateStr] = append(groupOpList[opCateStr], opSuffix)
					} else {
						groupOpList[opCateStr] = []string{opSuffix}
					}
				}
				continue
			}
			code := GetPermissionCodeMap(infos[0])
			ops := strings.Split(infos[1], ",")
			if _, ok := groupOpList[code]; ok {
				groupOpList[code] = append(groupOpList[code], ops...)
			} else {
				groupOpList[code] = ops
			}
		}

		for k, item := range groupOpList {
			result[k] = item
		}

		return result, nil
	}
}

// GetOpSuffixOfOperation 通过转化好的操作码，匹配它的权限的 suffix 值
// 如：`xxx.Config.Modify` 则得到 Modify
func GetOpSuffixOfOperation(op string) string {
	suffix := str.Substr(op, strings.LastIndex(op, ".")+1, len(op))
	return suffix
}

// 无码的**管理组**-权限项限映射为极星的权限项
func LessCodeOp2PolarisOpArr(op string) map[string][]string {
	groupOpList := make(map[string][]string, 0)
	switch op {
	case "1": // 通讯录管理-组织架构
		groupOpList["OrgConfig"] = []string{
			"Modify",
			"Transfer",
		}
		groupOpList["User"] = []string{
			"ModifyStatus",
			"Invite",
			"Bind",
			"Unbind",
			"Watch",
			"ModifyDepartment",
		}
	case "2": // 通讯录管理-角色架构
		groupOpList["Role"] = []string{
			"Create",
			"Modify",
			"Delete",
		}
		groupOpList["RoleGroup"] = []string{
			"Create",
			"Modify",
			"Delete",
		}
	case "3": // 可添加/删除应用
		groupOpList["Project"] = []string{
			"Create",
			"Attention",
			"UnAttention",
			"ModifyField",
		}
	default:
	}
	return groupOpList
}

func AssemblyOptAuthData(optAuthArr []string, groupOpList map[string][]string) map[string][]string {
	for _, item := range optAuthArr {
		infos := strings.Split(item, "-")
		if len(infos) < 2 {
			continue
		}
		code := GetPermissionCodeMap(infos[0])
		ops := strings.Split(infos[1], ",")
		if _, ok := groupOpList[code]; ok {
			groupOpList[code] = append(groupOpList[code], ops...)
		} else {
			groupOpList[code] = ops
		}
	}
	return groupOpList
}

// GetCateOfOperation 通过转化好的操作码，匹配它所在的分类组
// 如：`Permission.Pro.View-ManagePrivate` 拿到 Permission.Pro.View，去匹配分类
// 如：`Permission.Pro.View.ManagePrivate` 拿到 Permission.Pro.View，去匹配分类
func GetCateOfOperation(op string) string {
	prev := ""
	// 优先尝试使用 - 切割
	if strings.IndexAny(op, "-") != -1 {
		info := strings.Split(op, "-")
		if len(info) > 0 {
			prev = info[0]
		}
	} else {
		prev = str.Substr(op, 0, strings.LastIndex(op, "."))
	}
	cate := GetPermissionCodeMap(prev)
	return cate
}

// 权限 langCode 到 code 的映射，如：Permission.Org.ProjectObjectType->ProjectObjectType
func GetPermissionCodeMap(langCode string) string {
	m1 := map[string]string{
		// 项目权限 code 映射
		"Permission.Org.ProjectObjectType": "ProjectObjectType",
		"Permission.Pro.Config":            "ProConfig",
		"Permission.Pro.Iteration":         "Iteration",
		"Permission.Pro.Issue.4":           "Issue",
		"Permission.Pro.Role":              "Role",
		"Permission.Pro.File":              "File",
		"Permission.Pro.Tag":               "Tag",
		"Permission.Pro.Attachment":        "Attachment",
		"Permission.Pro.Member":            "Member",
		"Permission.Pro.View":              "View", // 新增分组-视图

		// 组织权限 code 映射
		"Permission.Org.Config":     "OrgConfig",
		"Permission.Org.User":       "User",
		"Permission.Org.Role":       "Role",
		"Permission.Org.Project":    "Project",
		"Permission.Org.Department": "Department",
		"Permission.Org.InviteUser": "InviteUser",        // 新增的分组。权限管理分组有所调整，产品：子龙
		"Permission.Org.AdminGroup": "AdminGroup",        // 新增的分组。
		"MenuPermission.Org":        "MenuPermissionOrg", // [菜单权限项] 管理组菜单权限项
		"MenuPermission.Pro":        "MenuPermissionPro", // [菜单权限项] 应用菜单权限项

		"Permission.Org.PersonInfo": "PersonInfo", // 个人信息管理
		"Permission.Org.AddUser":    "AddUser",    // 添加成员
	}
	res := ""
	if s, ok := m1[langCode]; ok {
		res = s
	}
	return res
}
