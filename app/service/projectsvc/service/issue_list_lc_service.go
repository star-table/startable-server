package service

import (
	"bytes"
	"fmt"
	"strconv"
	"strings"
	"time"

	"gitea.bjx.cloud/LessCode/go-common/utils/unsafe"
	"github.com/star-table/startable-server/app/facade/appfacade"
	"github.com/star-table/startable-server/app/facade/orgfacade"
	"github.com/star-table/startable-server/app/facade/permissionfacade"
	"github.com/star-table/startable-server/app/facade/userfacade"
	sconsts "github.com/star-table/startable-server/app/service/projectsvc/consts"
	"github.com/star-table/startable-server/app/service/projectsvc/domain"
	"github.com/star-table/startable-server/app/service/projectsvc/po"
	"github.com/star-table/startable-server/common/core/consts"
	"github.com/star-table/startable-server/common/core/errs"
	"github.com/star-table/startable-server/common/core/util/json"
	"github.com/star-table/startable-server/common/core/util/slice"
	"github.com/star-table/startable-server/common/extra/lc_helper"
	"github.com/star-table/startable-server/common/library/db/mysql"
	"github.com/star-table/startable-server/common/model/bo"
	"github.com/star-table/startable-server/common/model/vo"
	"github.com/star-table/startable-server/common/model/vo/appvo"
	"github.com/star-table/startable-server/common/model/vo/formvo"
	"github.com/star-table/startable-server/common/model/vo/orgvo"
	"github.com/star-table/startable-server/common/model/vo/permissionvo"
	"github.com/star-table/startable-server/common/model/vo/permissionvo/appauth"
	"github.com/star-table/startable-server/common/model/vo/projectvo"
	"github.com/spf13/cast"
	"upper.io/db.v3"
)

func MirrorStat(orgId, userId int64, appIdsStr []string) (map[int64]int64, errs.SystemErrorInfo) {
	var appIds []int64
	copyErr := json.FromJson(json.ToJsonIgnoreError(appIdsStr), &appIds)
	if copyErr != nil {
		log.Error(copyErr)
		return nil, errs.ParamError
	}
	res := make(map[int64]int64)
	appListResp := appfacade.GetAppInfoList(appvo.GetAppInfoListReq{
		OrgId:  orgId,
		AppIds: appIds,
	})
	if appListResp.Failure() {
		log.Error(appListResp.Error())
		return nil, appListResp.Error()
	}

	allViewIds := make([]int64, 0)
	viewProjectMap := make(map[int64]int64, 0)
	viewAppMap := make(map[int64][]int64, 0) //一个视图对应多个app快捷方式
	for _, datum := range appListResp.Data {
		allViewIds = append(allViewIds, datum.MirrorViewId)
		viewProjectMap[datum.MirrorViewId] = datum.ProjectId
		viewAppMap[datum.MirrorViewId] = append(viewAppMap[datum.MirrorViewId], datum.Id)
	}
	if len(allViewIds) == 0 {
		return res, nil
	}
	var viewList []po.LcAppView
	err := mysql.SelectAllByCond(consts.TableLcAppView, db.Cond{
		consts.TcOrgId:   orgId,
		consts.TcDelFlag: consts.AppIsNoDelete,
		consts.TcId:      db.In(allViewIds),
	}, &viewList)
	if err != nil {
		log.Error(err)
		return nil, errs.MysqlOperateError
	}

	//处理参数
	for _, view := range viewList {
		configReq := vo.AppViewConfig{}
		err1 := json.FromJson(view.Config, &configReq)
		if err1 != nil {
			log.Error(err1)
			return nil, errs.JSONConvertError
		}

		params := &vo.HomeIssueInfoReq{}
		if projectId, ok := viewProjectMap[view.ID]; ok {
			if projectId >= 0 {
				params.ProjectID = &projectId
			}
		}

		//拼接参数
		var lessDataArr = []*vo.LessCondsData{
			{
				Type:   "equal",
				Value:  consts.DeleteFlagNotDel,
				Column: consts.BasicFieldRecycleFlag,
			},
		}
		//类型(between,equal,gt,gte,in,like,lt,lte,not_in,not_like,not_null,is_null,all_in,values_in)
		if configReq.RealCondition != nil {
			lessDataArr = append(lessDataArr, configReq.RealCondition)
		}
		if configReq.TableId > 0 {
			lessDataArr = append(lessDataArr, &vo.LessCondsData{
				Type:   "equal",
				Value:  configReq.TableId,
				Column: consts.BasicFieldTableId,
			})
		}
		params.LessConds = &vo.LessCondsData{
			Type:  "and",
			Conds: lessDataArr,
		}
		//排序
		//if req.Orders != nil {
		//	params.LessOrder = req.Orders
		//}

		//归档的问题稍后再说 todo
		//isFiling := c.Query("isFiling")
		//isFilingInt, isFilingIntErr:= strconv.Atoi(isFiling)
		//if isFilingIntErr != nil {
		//	log.Error(isFilingIntErr)
		//} else {
		//	params.IsFiling = &isFilingInt
		//}
		count, countErr := handleCount(orgId, userId, view.AppId, 0, 0, params)
		if countErr != nil {
			log.Error(countErr)
			return nil, countErr
		}

		if appIds, ok := viewAppMap[view.ID]; ok {
			for _, id := range appIds {
				res[id] = count
			}
		}
	}

	return res, nil
}

func handleCount(orgId, currentUserId, appId int64, page int, size int, input *vo.HomeIssueInfoReq) (int64, errs.SystemErrorInfo) {
	projectBo := &bo.ProjectBo{}
	// 获取项目信息
	if input.ProjectID != nil && *input.ProjectID > 0 {
		tmpProBo, err := domain.GetProject(orgId, *input.ProjectID)
		if err != nil {
			return 0, errs.BuildSystemErrorInfo(errs.ProjectDomainError, err)
		}
		projectBo = tmpProBo
	}

	orgInfoResp := orgfacade.OrganizationInfo(orgvo.OrganizationInfoReqVo{
		OrgId:  orgId,
		UserId: currentUserId,
	})
	if orgInfoResp.Failure() {
		log.Error(orgInfoResp.Error())
		return 0, orgInfoResp.Error()
	}
	orgRemarkObj := &orgvo.OrgRemarkConfigType{}
	oriErr := json.FromJson(orgInfoResp.OrganizationInfo.Remark, orgRemarkObj)
	if oriErr != nil {
		log.Error(oriErr)
		return 0, errs.BuildSystemErrorInfo(errs.JSONConvertError, oriErr)
	}

	ownerIsContain, _ := slice.Contain(projectBo.OwnerIds, currentUserId)
	isAdmin := domain.CheckIsAdmin(orgId, currentUserId, appId) || (currentUserId == projectBo.Owner || ownerIsContain)
	if input.ProjectID != nil {
		if !isAdmin && projectBo.TemplateFlag != consts.TemplateTrue {
			input.LessConds.Conds = append(input.LessConds.Conds, domain.GetCollaborateDataQueryConds(orgId, currentUserId, []int64{}))
		}
	}

	////根据所有符合条件的id去查询无码
	//var allUsefulIssueIds []interface{}
	//issueInfoMap := map[int64]bo.IssueBo{}
	//for _, issueBo := range *issueBos {
	//	allUsefulIssueIds = append(allUsefulIssueIds, issueBo.Id)
	//	issueInfoMap[issueBo.Id] = issueBo
	//}
	lessReq := *input.LessConds
	if input.ProjectID != nil && *input.ProjectID != 0 {
		projectInfo, err := domain.GetProject(orgId, *input.ProjectID)
		if err != nil {
			log.Error(err)
			return 0, err
		}

		if projectInfo.ProjectTypeId == consts.ProjectTypeNormalId {
			lessReq = *domain.ConvertIssueStatusFilterReqForNormalProject(&lessReq)
		}
	} else {
		if input.ProjectID == nil {
			lessReq = domain.ConvertIssueStatusFilterReqForAll(lessReq)
		}
	}

	lessReqParam := formvo.LessIssueListReq{
		FilterColumns: []string{" count(*) as \"all\" "},
		Condition:     lessReq,
		AppId:         orgRemarkObj.OrgSummaryTableAppId,
		OrgId:         orgId,
		UserId:        currentUserId,
		Page:          int64(page),
		Size:          int64(size),
	}
	if input.LessOrder != nil {
		lessReqParam.Orders = input.LessOrder
	}
	//暂时只在单项目查询的时候传入，因为综合查询是不需要展示自定义字段的
	if input.ProjectID != nil {
		lessReqParam.RedirectIds = []int64{projectBo.AppId}
	}
	if input.TableID != nil {
		tableId, tableIdErr := strconv.ParseInt(*input.TableID, 10, 64)
		if tableIdErr != nil {
			log.Errorf("[handleCount] tableId err:%v. orgId:%d, userId:%d", tableIdErr, orgId, currentUserId)
			return 0, errs.InvalidTableId
		}
		lessReqParam.TableId = tableId
	}
	lessResp := domain.GetRowsExpand(&lessReqParam)
	if lessResp.Failure() {
		log.Error(lessResp.Error())
		return 0, lessResp.Error()
	}

	count := int64(0)
	if len(lessResp.List) > 0 {
		count = cast.ToInt64(lessResp.List[0]["all"])
	}

	return count, nil
}

func LcViewStatForAll(orgId, userId int64) ([]*projectvo.LcViewStatVo, errs.SystemErrorInfo) {
	// 拿汇总表AppId
	orgInfoResp := orgfacade.OrganizationInfo(orgvo.OrganizationInfoReqVo{
		OrgId:  orgId,
		UserId: userId,
	})
	if orgInfoResp.Failure() {
		log.Errorf("[MirrorStatForAll] orgId:%v, userId:%v, OrganizationInfo failure:%v", orgId, userId, orgInfoResp.Error())
		return nil, orgInfoResp.Error()
	}
	orgRemarkObj := orgvo.OrgRemarkConfigType{}
	if err := json.FromJson(orgInfoResp.OrganizationInfo.Remark, &orgRemarkObj); err != nil {
		log.Errorf("[MirrorStatForAll] orgId:%v, userId:%v, OrgInfo FromJson failure:%v", orgId, userId, err)
		return nil, errs.BuildSystemErrorInfo(errs.JSONConvertError, err)
	}
	appId := orgRemarkObj.OrgSummaryTableAppId

	// 获取用户部门信息
	resp := orgfacade.GetUserDeptIdsWithParentId(orgvo.GetUserDeptIdsWithParentIdReq{
		OrgId:  orgId,
		UserId: userId,
	})
	if resp.Failure() {
		log.Errorf("[MirrorStatForAll] orgId:%v, userId:%v, GetUserDeptIdsWithParentId failure:%v", orgId, userId, resp.Error())
		return nil, resp.Error()
	}
	deptIdsResp := resp
	log.Infof("[MirrorStatForAll] orgId:%d, userId:%d, deptIds:%v", orgId, userId, deptIdsResp.Data.DeptIds)

	// 是否管理员。
	// 是否可以有权限管理所有项目
	manageAuthInfoResp := userfacade.GetUserAuthority(orgId, userId)
	if manageAuthInfoResp.Failure() {
		log.Errorf("[MirrorStatForAll] orgId:%v, userId:%v, error:%v", orgId, userId, manageAuthInfoResp.Error())
		return nil, manageAuthInfoResp.Error()
	}
	isSysAdmin := manageAuthInfoResp.Data.IsSysAdmin
	isOrgOwner := manageAuthInfoResp.Data.IsOrgOwner
	isSubAdmin := manageAuthInfoResp.Data.IsSubAdmin
	manageApps := manageAuthInfoResp.Data.ManageApps
	isAdmin := isSysAdmin || isOrgOwner || (isSubAdmin && len(manageApps) > 0 && manageApps[0] == -1)
	log.Infof("[MirrorStatForAll] orgId:%d, userId:%d, isAdmin:%v", orgId, userId, isAdmin)

	// 权限验证
	appAuthsInfo := permissionfacade.GetAppAuths(permissionvo.GetAppAuthBatchReq{OrgId: orgId, UserId: userId, AppIds: []int64{}})
	if appAuthsInfo.Failure() {
		log.Errorf("[MirrorStatForAll] orgId:%d, userId:%d, GetAppAuths failure:%v", orgId, userId, appAuthsInfo.Error())
		return nil, appAuthsInfo.Error()
	}
	log.Infof("[MirrorStatForAll] orgId:%d, userId:%d, app auths: %v", orgId, userId, appAuthsInfo)

	// 所有没删除且未归档的项目，且不是空项目
	projectIds, errSys := domain.GetAllNotFillingAnEmptyProjectIds(orgId, false)
	if errSys != nil {
		log.Errorf("[MirrorStatForAll] orgId:%d, userId:%d, GetAllProjectIds failure:%v", orgId, userId, errSys)
		return nil, errSys
	}

	// 获取我参与的所有项目
	myProjectIds, errSys := domain.GetAllMyProjectIdsWithDeptIds(orgId, userId, deptIdsResp.Data.DeptIds, true)
	if errSys != nil {
		log.Errorf("[MirrorStatForAll] orgId:%d, userId:%d, deptIds:%v, GetAllMyProjectIdsWithDeptIds failure:%v", orgId, userId, deptIdsResp.Data.DeptIds, errSys)
		return nil, errSys
	}
	myProjectIds = filterEmptyProject(projectIds, myProjectIds)
	if len(manageApps) > 0 && manageApps[0] != -1 {
		proIds, errSys := domain.GetProjectIdsByAppIds(orgId, manageApps)
		if errSys != nil {
			log.Errorf("[MirrorStatForAll] err:%v, orgId:%v, userId:%v", errSys, orgId, userId)
			return nil, errSys
		}
		myProjectIds = append(myProjectIds, proIds...)
	}
	myProjectIds = slice.SliceUniqueInt64(myProjectIds)

	// 拿任务模块视图
	viewResp := appfacade.GetAppViewList(&appvo.GetAppViewListReq{
		OrgId:  orgId,
		AppId:  appId,
		UserId: userId,
	})
	if viewResp.Failure() {
		log.Error(viewResp.Error())
		return nil, viewResp.Error()
	}

	count := 0
	var viewStats []*projectvo.LcViewStatVo
	for _, view := range viewResp.Data {
		viewConfig := &vo.AppViewConfig{}
		err := json.FromJson(cast.ToString(view.Config), viewConfig)
		if err != nil {
			log.Error(err)
			return nil, errs.JSONConvertError
		}

		// 构造筛选条件
		lessReq := vo.LessCondsData{Type: "and"}

		// 筛选条件：前端传入的条件
		if viewConfig.RealCondition != nil {
			lessReq.Conds = append(lessReq.Conds, viewConfig.RealCondition)
			changeInMapCondition(&lessReq)
		}

		lessReq.Conds = append(lessReq.Conds, &vo.LessCondsData{
			Type:   "equal",
			Value:  orgId,
			Column: consts.BasicFieldOrgId,
		})

		lessReq.Conds = append(lessReq.Conds, &vo.LessCondsData{
			Type:   "equal",
			Value:  2,
			Column: consts.BasicFieldRecycleFlag,
		})

		if viewConfig.TableId > 0 {
			lessReq.Conds = append(lessReq.Conds, &vo.LessCondsData{
				Type:   "equal",
				Value:  viewConfig.TableId,
				Column: consts.BasicFieldTableId,
			})
		}

		if isAdmin {
			// 筛选条件：管理员只需排除已删除的项目（因历史数据删除项目的时候没有删除无码任务，所以需要通过极星的有效project id去筛选）
			// log.Infof("[LcHomeIssuesForAll] orgId:%d, userId:%d, projectIds:%v", orgId, currentUserId, projectIds)
			lessReq.Conds = append(lessReq.Conds, domain.GetProjectIdsQueryConds(projectIds))
		} else {
			innerConds1 := []*vo.LessCondsData{
				// 筛选条件：我参与的所有项目
				domain.GetProjectIdsQueryConds(myProjectIds),
			}
			// 筛选条件：我有权限的表格
			tableCond := domain.GetTableIdsQueryCondsByAppAuths(appAuthsInfo.Data)
			if tableCond != nil {
				innerConds1 = append(innerConds1, tableCond)
			}

			innerConds2 := []*vo.LessCondsData{
				// 排除已删除的项目
				domain.GetProjectIdsQueryConds(projectIds),
				// 筛选条件：我协作的任务
				domain.GetCollaborateDataQueryConds(orgId, userId, deptIdsResp.Data.DeptIds),
			}
			lessReq.Conds = append(lessReq.Conds, &vo.LessCondsData{
				Type: "or",
				Conds: []*vo.LessCondsData{
					{
						Type:  "and",
						Conds: innerConds1,
					},
					{
						Type:  "and",
						Conds: innerConds2,
					},
				},
			})
		}

		// 筛选条件：处理任务状态转换
		lessReq = domain.ConvertIssueStatusFilterReqForAll(lessReq)

		log.Infof("[MirrorStatForAll] orgId:%d, userId:%d, LessIssueListReq conds:%v", orgId, userId, lessReq.Conds)

		// 拉取数据
		lessReqParam := formvo.LessIssueListReq{
			Condition: lessReq,
			AppId:     appId,
			OrgId:     orgId,
			UserId:    userId,
			Page:      int64(1),
			Size:      int64(1),
		}
		lessReqParam.FilterColumns = []string{"count(1) as count"}

		log.Infof("[MirrorStatForAll] orgId:%d, userId:%d, LessIssueListReq req: %s", orgId, userId, json.ToJsonIgnoreError(lessReqParam))
		lessResp := domain.GetRowsExpand(&lessReqParam)
		if lessResp.Failure() {
			log.Errorf("[MirrorStatForAll] orgId:%d, userId:%d, LessIssueList failure:%v", orgId, userId, lessResp.Error())
			return nil, lessResp.Error()
		}

		if len(lessResp.List) > 0 {
			viewStats = append(viewStats, &projectvo.LcViewStatVo{
				Id:    cast.ToInt64(view.ID),
				Name:  view.ViewName,
				Total: cast.ToInt64(lessResp.List[0]["count"]),
			})
		}

		count += 1
		if count >= 6 {
			break
		}
	}

	return viewStats, nil
}

// LcHomeIssuesForAll 基于组织下所有项目查询任务列表
// @param isInnerSuper 是否是 openAPI或特定需要超管权限的内部请求，如果是，则不验证权限
func LcHomeIssuesForAll(orgId, currentUserId int64, page int, size int, input *projectvo.HomeIssueInfoReq, isInnerSuper bool) (*projectvo.LcHomeIssuesRespVo, errs.SystemErrorInfo) {
	var inputAppId int64
	var projectIds []int64
	var err error
	var errSys errs.SystemErrorInfo

	// 主页任务菜单 前端必须传汇总表ID过来
	if input.MenuAppID != nil {
		inputAppId, err = strconv.ParseInt(*input.MenuAppID, 10, 64)
		if err != nil {
			log.Errorf("[LcHomeIssuesForAll] orgId:%v, curUserId:%v, appId:%v, appId failure:%v", orgId, currentUserId, *input.MenuAppID, err)
			return nil, errs.ParamError
		}
	}
	// 没拿到汇总表的APP ID的情况，后端兜个底
	if inputAppId <= 0 {
		if isInnerSuper {
			return nil, errs.ParamError
		} else {
			orgInfoResp := orgfacade.OrganizationInfo(orgvo.OrganizationInfoReqVo{
				OrgId:  orgId,
				UserId: currentUserId,
			})
			if orgInfoResp.Failure() {
				log.Errorf("[LcHomeIssuesForAll] orgId:%v, curUserId:%v, OrganizationInfo failure:%v", orgId, currentUserId, orgInfoResp.Error())
				return nil, orgInfoResp.Error()
			}
			orgRemarkObj := orgvo.OrgRemarkConfigType{}
			if err = json.FromJson(orgInfoResp.OrganizationInfo.Remark, &orgRemarkObj); err != nil {
				log.Errorf("[LcHomeIssuesForAll] orgId:%v, curUserId:%v, OrgInfo FromJson failure:%v", orgId, currentUserId, err)
				return nil, errs.BuildSystemErrorInfo(errs.JSONConvertError, err)
			}
			inputAppId = orgRemarkObj.OrgSummaryTableAppId
		}
	}

	log.Infof("[LcHomeIssuesForAll] orgId:%d, userId:%d, appId:%v", orgId, currentUserId, inputAppId)

	// 获取用户部门信息
	var deptIdsResp *orgvo.GetUserDeptIdsWithParentIdResp
	if !isInnerSuper {
		resp := orgfacade.GetUserDeptIdsWithParentId(orgvo.GetUserDeptIdsWithParentIdReq{
			OrgId:  orgId,
			UserId: currentUserId,
		})
		//resp := orgfacade.GetUserDeptIds(orgvo.GetUserDeptIdsReq{
		//	OrgId:  orgId,
		//	UserId: currentUserId,
		//})
		if resp.Failure() {
			log.Errorf("[LcHomeIssuesForAll] orgId:%v, curUserId:%v, GetUserDeptIds failure:%v", orgId, currentUserId, resp.Error())
			return nil, resp.Error()
		}
		deptIdsResp = resp
		log.Infof("[LcHomeIssuesForAll] orgId:%d, userId:%d, deptIds:%v", orgId, currentUserId, deptIdsResp.Data.DeptIds)
	}

	// 是否管理员。isInnerSuper 为 true 表示来源为 openAPI，无需权限校验
	// 是否可以有权限管理所有项目
	manageAuthInfoResp := userfacade.GetUserAuthority(orgId, currentUserId)
	if manageAuthInfoResp.Failure() {
		log.Errorf("[LcHomeIssuesForAll] orgId:%v, curUserId:%v, error:%v", orgId, currentUserId, manageAuthInfoResp.Error())
		return nil, manageAuthInfoResp.Error()
	}
	isSysAdmin := manageAuthInfoResp.Data.IsSysAdmin
	isOrgOwner := manageAuthInfoResp.Data.IsOrgOwner
	isSubAdmin := manageAuthInfoResp.Data.IsSubAdmin
	manageApps := manageAuthInfoResp.Data.ManageApps
	isCanManageAllApps := isSysAdmin || isOrgOwner || (isSubAdmin && len(manageApps) > 0 && manageApps[0] == -1)
	isAdmin := isInnerSuper || isCanManageAllApps
	//isAdmin := isInnerSuper || domain.CheckIsAdmin(orgId, currentUserId)
	log.Infof("[LcHomeIssuesForAll] orgId:%d, userId:%d, isAdmin:%v", orgId, currentUserId, isAdmin)

	// 权限验证
	var appAuthsInfo *permissionvo.GetAppAuthBatchResp
	if !isInnerSuper {
		appAuthsInfo = permissionfacade.GetAppAuths(permissionvo.GetAppAuthBatchReq{OrgId: orgId, UserId: currentUserId, AppIds: []int64{}})
		if appAuthsInfo.Failure() {
			log.Errorf("[LcHomeIssuesForAll] orgId:%d, userId:%d, GetAppAuths failure:%v", orgId, currentUserId, appAuthsInfo.Error())
			return nil, appAuthsInfo.Error()
		}
		//log.Infof("[LcHomeIssuesForAll] orgId:%d, userId:%d, app auths: %v", orgId, currentUserId, appAuthsInfo)
	}

	needDeleteData, err2 := setIsNeedRefresh(input, 0)
	if err2 != nil {
		log.Errorf("[LcHomeIssuesForAll] setIsNeedRefresh orgId:%d, userId:%d, LessIssueList failure:%v", orgId, currentUserId, err)
		return nil, err2
	}

	// 构造筛选条件
	lessReq := vo.LessCondsData{Type: "and"}

	// 筛选条件：前端传入的条件
	if input.LessConds != nil {
		changeInMapCondition(input.LessConds)
		lessReq = *input.LessConds
	}

	lessReq.Conds = append(lessReq.Conds, getRefreshCondition(input)...)

	if input.TableID != nil && *input.TableID != "0" {
		lessReq.Conds = append(lessReq.Conds, &vo.LessCondsData{
			Type:   "equal",
			Value:  *input.TableID,
			Column: consts.BasicFieldTableId,
		})
	}

	// 所有没删除且未归档的项目，且不是空项目
	projectIds, errSys = domain.GetAllNotFillingAnEmptyProjectIds(orgId, input.IsNeedAll)
	if errSys != nil {
		log.Errorf("[LcHomeIssuesForAll] orgId:%d, userId:%d, GetAllProjectIds failure:%v", orgId, currentUserId, errSys)
		return nil, errSys
	}

	if isAdmin {
		// 筛选条件：管理员只需排除已删除的项目（因历史数据删除项目的时候没有删除无码任务，所以需要通过极星的有效project id去筛选）
		// log.Infof("[LcHomeIssuesForAll] orgId:%d, userId:%d, projectIds:%v", orgId, currentUserId, projectIds)
		lessReq.Conds = append(lessReq.Conds, domain.GetProjectIdsQueryConds(projectIds))
	} else {
		// 获取我参与的所有项目
		myProjectIds, errSys := domain.GetAllMyProjectIdsWithDeptIds(orgId, currentUserId, deptIdsResp.Data.DeptIds, true)
		if errSys != nil {
			log.Errorf("[LcHomeIssuesForAll] orgId:%d, userId:%d, deptIds:%v, GetAllMyProjectIdsWithDeptIds failure:%v", orgId, currentUserId, deptIdsResp.Data.DeptIds, errSys)
			return nil, errSys
		}
		myProjectIds = filterEmptyProject(projectIds, myProjectIds)

		if len(manageApps) > 0 && manageApps[0] != -1 {
			proIds, errSys := domain.GetProjectIdsByAppIds(orgId, manageApps)
			if errSys != nil {
				log.Errorf("[LcHomeIssuesForAll] err:%v, orgId:%v, userId:%v", errSys, orgId, currentUserId)
				return nil, errSys
			}
			myProjectIds = append(myProjectIds, proIds...)
		}
		myProjectIds = slice.SliceUniqueInt64(myProjectIds)
		innerConds1 := []*vo.LessCondsData{
			// 筛选条件：我参与的所有项目
			domain.GetProjectIdsQueryConds(myProjectIds),
		}
		// 筛选条件：我有权限的表格
		tableCond := domain.GetTableIdsQueryCondsByAppAuths(appAuthsInfo.Data)
		if tableCond != nil {
			innerConds1 = append(innerConds1, tableCond)
		}

		innerConds2 := []*vo.LessCondsData{
			// 排除已删除的项目
			domain.GetProjectIdsQueryConds(projectIds),
			// 筛选条件：我协作的任务
			domain.GetCollaborateDataQueryConds(orgId, currentUserId, deptIdsResp.Data.DeptIds),
		}
		lessReq.Conds = append(lessReq.Conds, &vo.LessCondsData{
			Type: "or",
			Conds: []*vo.LessCondsData{
				&vo.LessCondsData{
					Type:  "and",
					Conds: innerConds1,
				},
				&vo.LessCondsData{
					Type:  "and",
					Conds: innerConds2,
				},
			},
		})
	}

	// 筛选条件：处理任务状态转换
	lessReq = domain.ConvertIssueStatusFilterReqForAll(lessReq)

	log.Infof("[LcHomeIssuesForAll] orgId:%d, userId:%d, LessIssueListReq conds:%v", orgId, currentUserId, lessReq.Conds)

	addOrder(input, isInnerSuper)

	// 拉取数据
	lessReqParam := formvo.LessIssueListReq{
		Condition:      lessReq,
		AppId:          inputAppId,
		OrgId:          orgId,
		UserId:         currentUserId,
		Page:           int64(page),
		Size:           int64(size),
		Export:         true,
		Orders:         input.LessOrder,
		NeedTotal:      input.IsNeedTotal,
		NeedDeleteData: needDeleteData,
		NeedChangeId:   true,
	}
	if input.FilterColumns != nil && len(input.FilterColumns) > 0 {
		filterColumns := []string{
			lc_helper.ConvertToFilterColumn(consts.BasicFieldId),
			lc_helper.ConvertToFilterColumn(consts.BasicFieldIssueId),
			lc_helper.ConvertToFilterColumn(consts.BasicFieldUpdateTime),
		}
		for _, column := range input.FilterColumns {
			filterColumns = append(filterColumns, lc_helper.ConvertToFilterColumn(column))
		}
		lessReqParam.FilterColumns = slice.SliceUniqueString(filterColumns)
	} else {
		lessReqParam.FilterColumns = []string{
			lc_helper.ConvertToFilterColumn(consts.BasicFieldId),
			lc_helper.ConvertToFilterColumn(consts.BasicFieldIssueId),
			lc_helper.ConvertToFilterColumn(consts.BasicFieldTitle),
			lc_helper.ConvertToFilterColumn(consts.BasicFieldCode),
			lc_helper.ConvertToFilterColumn(consts.BasicFieldOwnerId),
			lc_helper.ConvertToFilterColumn(consts.BasicFieldAuditStatus),
			lc_helper.ConvertToFilterColumn(consts.BasicFieldAuditStatusDetail),
			lc_helper.ConvertToFilterColumn(consts.BasicFieldIssueStatus),
			lc_helper.ConvertToFilterColumn(consts.BasicFieldIssueStatusType),
			lc_helper.ConvertToFilterColumn(consts.BasicFieldProjectObjectTypeId),
			lc_helper.ConvertToFilterColumn(consts.BasicFieldPlanStartTime),
			lc_helper.ConvertToFilterColumn(consts.BasicFieldPlanEndTime),
			lc_helper.ConvertToFilterColumn(consts.BasicFieldAppId),
			lc_helper.ConvertToFilterColumn(consts.BasicFieldProjectId),
			lc_helper.ConvertToFilterColumn(consts.BasicFieldTableId),
			lc_helper.ConvertToFilterColumn(consts.BasicFieldRemark),
			lc_helper.ConvertToFilterColumn(consts.BasicFieldFollowerIds),
			lc_helper.ConvertToFilterColumn(consts.BasicFieldAuditorIds),
			lc_helper.ConvertToFilterColumn(consts.BasicFieldParentId),
			lc_helper.ConvertToFilterColumn(consts.BasicFieldRelating),
			lc_helper.ConvertToFilterColumn(consts.BasicFieldBaRelating),
			lc_helper.ConvertToFilterColumn(consts.BasicFieldUpdateTime),
			lc_helper.ConvertToFilterColumn(consts.BasicFieldRecycleFlag),
			lc_helper.ConvertToFilterColumn(consts.BasicFieldCreator),
		}
		columns, err := getTableColumnsByOrgId(orgId, currentUserId)
		if err != nil {
			log.Errorf("[LcHomeIssuesForAll] getTableColumnsByOrgId orgId:%d, userId:%d, LessIssueList failure:%v", orgId, currentUserId, err)
		}
		for _, column := range columns.Columns {
			if column.SummaryFlag == consts.SummaryFlagAll {
				lessReqParam.FilterColumns = append(lessReqParam.FilterColumns, lc_helper.ConvertToFilterColumn(column.Name))
			}
		}
	}
	log.Infof("[LcHomeIssuesForAll] orgId:%d, userId:%d, LessIssueListReq req: %s", orgId, currentUserId, json.ToJsonIgnoreError(lessReqParam))
	lessResp := domain.GetRows(&lessReqParam)
	if lessResp.Failure() {
		log.Errorf("[LcHomeIssuesForAll] orgId:%d, userId:%d, LessIssueList failure:%v", orgId, currentUserId, lessResp.Error())
		return nil, lessResp.Error()
	}

	if int(lessResp.Data.Count) < size {
		input.IsRefreshAll = false
	}

	return issueListToJsonString(lessResp, input), nil
}

func issueListToJsonString(reply *projectvo.ListRowsReply, input *projectvo.HomeIssueInfoReq) *projectvo.LcHomeIssuesRespVo {
	startTime := time.Now()
	resp := &projectvo.LcHomeIssuesRespVo{}
	lastUpdateTime := input.LastUpdateTime
	if reply.Data.LastUpdateTime > lastUpdateTime {
		lastUpdateTime = reply.Data.LastUpdateTime
	}
	str := fmt.Sprintf(`{"code":0,"data":{"total":%d,"lastUpdateTime":"%s","lastPermissionUpdateTime":"%s", "userDepts":%s, "isRefreshAll":%v,"list":`,
		reply.Data.Count, lastUpdateTime, input.LastPermissionUpdateTime, json.ToJsonIgnoreError(reply.UserDepts), input.IsRefreshAll)
	buffer := bytes.Buffer{}
	buffer.Grow(len(reply.Data.Data) + len(reply.Data.RelateData) + len(str) + 15)
	buffer.WriteString(str)
	buffer.Write(reply.Data.Data)
	if len(reply.Data.RelateData) > 0 {
		buffer.WriteString(`, "relateData":`)
		buffer.Write(reply.Data.RelateData)
	}
	buffer.WriteString("}}")
	resp.Data = unsafe.BytesString(buffer.Bytes())

	log.Infof("[issueListToJsonString] cost time:%v", time.Since(startTime).Seconds())

	return resp
}

func getRefreshCondition(input *projectvo.HomeIssueInfoReq) []*vo.LessCondsData {
	// 筛选条件：未删除的任务
	conditions := make([]*vo.LessCondsData, 0, 3)
	if input.IsRefreshAll {
		conditions = append(conditions, &vo.LessCondsData{
			Type:   "equal",
			Value:  2,
			Column: consts.BasicFieldRecycleFlag,
		})
	} else {
		if input.LastUpdateTime != "" {
			conditions = append(conditions, &vo.LessCondsData{
				Type:   "gt",
				Value:  input.LastUpdateTime,
				Column: consts.BasicFieldUpdateTime,
			})
		}
	}

	return conditions
}

// 因为allProjectIds里面没有empty项目了，所以不在里面的就是empty项目
func filterEmptyProject(allProjectIds, myProjectIds []int64) []int64 {
	idsMap := make(map[int64]int64, len(allProjectIds))
	for _, id := range allProjectIds {
		idsMap[id] = id
	}
	newIds := make([]int64, 0, len(myProjectIds))
	for _, id := range myProjectIds {
		if idsMap[id] != 0 {
			newIds = append(newIds, id)
		}
	}

	return newIds
}

// 如果是in和not_in，将map的条件转换为单个条件逻辑
func changeInMapCondition(conds *vo.LessCondsData) {
	for i, cond := range conds.Conds {
		if cond.Column == consts.BasicFieldCreator {
			if is, ok := cond.Values.([]interface{}); ok {
				for j := range is {
					is[j] = strings.Replace(cast.ToString(is[j]), "U_", "", 1)
				}
			}
		}
		if cond.Type == "in" || cond.Type == "not_in" {
			if values, ok := cond.Values.([]interface{}); ok && len(values) > 0 {
				newCond := vo.LessCondsData{Type: "or"}
				if cond.Type == "not_in" {
					newCond = vo.LessCondsData{Type: "and"}
				}
				for _, value := range values {
					if m, ok2 := value.(map[string]interface{}); ok2 && len(m) > 0 {
						andCond := vo.LessCondsData{Type: "and"}
						for s, v := range m {
							if _, ok3 := v.([]interface{}); ok3 {
								andCond.Conds = append(andCond.Conds, &vo.LessCondsData{
									Type:   cond.Type,
									Values: v,
									Column: s,
								})
							} else {
								andCond.Conds = append(andCond.Conds, &vo.LessCondsData{
									Type:   cond.Type,
									Values: []interface{}{v},
									Column: s,
								})
							}
						}
						newCond.Conds = append(newCond.Conds, &andCond)
					}
				}
				if len(newCond.Conds) > 0 {
					conds.Conds[i] = &newCond
				}
			}
		} else if cond.Type == "and" || cond.Type == "or" {
			changeInMapCondition(cond)
		}
	}
}

// LcHomeIssuesForProject 项目下的任务筛选
// @param isInnerSuper 如果是 openAPI或特定需要超管权限的内部请求，则无需校验具体权限，并且不受字段权限控制（可以看到所有字段）
func LcHomeIssuesForProject(orgId, currentUserId int64, page int, size int, input *projectvo.HomeIssueInfoReq, isInnerSuper bool) (*projectvo.LcHomeIssuesRespVo, errs.SystemErrorInfo) {
	var inputAppId, summaryAppId int64
	var err error
	var sysErr errs.SystemErrorInfo
	var projectId int64

	// 必须传tableId
	if input.TableID == nil {
		log.Errorf("[LcHomeIssuesForProject] no tableId. orgId:%d, userId:%d", orgId, currentUserId)
		return nil, errs.InvalidTableId
	}
	tableId, tableIdErr := strconv.ParseInt(*input.TableID, 10, 64)
	if tableIdErr != nil {
		log.Errorf("[LcHomeIssuesForProject] tableId err:%v. orgId:%d, userId:%d", tableIdErr, orgId, currentUserId)
		return nil, errs.InvalidTableId
	}

	// 获取项目信息
	projectBo := &bo.ProjectBo{}
	if input.ProjectID != nil && *input.ProjectID > 0 {
		projectId = *input.ProjectID
		projectBo, sysErr = domain.GetProjectSimple(orgId, projectId)
		if sysErr != nil {
			return nil, errs.BuildSystemErrorInfo(errs.ProjectDomainError, sysErr)
		}
	} else {
		// 没项目信息，转为查所有项目
		return LcHomeIssuesForAll(orgId, currentUserId, page, size, input, isInnerSuper)
	}

	// 获取APP信息
	if input.MenuAppID != nil && *input.MenuAppID != "" {
		inputAppId, err = strconv.ParseInt(*input.MenuAppID, 10, 64)
		if err != nil {
			log.Errorf("[LcHomeIssuesForProject] MenuAppID %v ParseInt err:%v, orgId:%d, userId:%d", input.MenuAppID, err, orgId, currentUserId)
			return nil, errs.ParamError
		}
	}

	//// 获取表头信息
	//tableColumnMap, sysErr := domain.GetTableColumnsMap(orgId, tableId, nil, true)
	//if sysErr != nil {
	//	log.Errorf("[LcHomeIssuesForProject] 获取表头失败 org:%d app:%d proj:%d table:%d user:%d, err: %v",
	//		orgId, inputAppId, projectId, tableId, currentUserId, sysErr)
	//	return nil, errs.BuildSystemErrorInfo(errs.ProjectDomainError, sysErr)
	//}
	//var tableColumns []*projectvo.TableColumnData
	//for _, column := range tableColumnMap {
	//	tableColumns = append(tableColumns, column)
	//}

	// 注意：
	// projectBo.AppId 是原始项目的AppId
	// inputAppId 可能是项目创建的镜像AppId，也可能是原始项目的AppId

	// 获取汇总表信息
	orgInfoResp := orgfacade.OrganizationInfo(orgvo.OrganizationInfoReqVo{
		OrgId:  orgId,
		UserId: currentUserId,
	})
	if orgInfoResp.Failure() {
		log.Errorf("[LcHomeIssuesForProject] OrganizationInfo err:%v. orgId:%d, userId:%d", orgInfoResp.Error(), orgId, currentUserId)
		return nil, orgInfoResp.Error()
	}
	orgRemarkObj := &orgvo.OrgRemarkConfigType{}
	if err = json.FromJson(orgInfoResp.OrganizationInfo.Remark, orgRemarkObj); err != nil {
		log.Errorf("[LcHomeIssuesForProject] OrganizationInfo Remark %s FromJson err:%v. orgId:%d, userId:%d", orgInfoResp.OrganizationInfo.Remark, err, orgId, currentUserId)
		return nil, errs.BuildSystemErrorInfo(errs.JSONConvertError, err)
	}
	summaryAppId = orgRemarkObj.OrgSummaryTableAppId

	log.Infof("[LcHomeIssuesForProject] orgId:%d, userId:%d, projectId:%d, projAppId:%v, inputAppId:%v, summaryAppId: %d",
		orgId, currentUserId, projectId, projectBo.AppId, inputAppId, summaryAppId)

	// 获取用户部门信息
	//var deptIdsResp *orgvo.GetUserDeptIdsResp
	//if !isInnerSuper {
	//	resp := orgfacade.GetUserDeptIds(orgvo.GetUserDeptIdsReq{
	//		OrgId:  orgId,
	//		UserId: currentUserId,
	//	})
	//	if resp.Failure() {
	//		log.Errorf("[LcHomeIssuesForProject] orgId:%v, curUserId:%v, GetUserDeptIds failure:%v", orgId, currentUserId, resp.Error())
	//		return nil, resp.Error()
	//	}
	//	deptIdsResp = &resp
	//	log.Infof("[LcHomeIssuesForProject] orgId:%d, userId:%d, deptIds:%v", orgId, currentUserId, deptIdsResp.Data.DeptIds)
	//}

	// 获取app权限
	isAdmin := false
	var appAuthInfo *appauth.GetAppAuthData
	if isInnerSuper {
		isAdmin = isInnerSuper
	} else {
		optAuthResp := permissionfacade.GetAppAuth(orgId, projectBo.AppId, tableId, currentUserId)
		if optAuthResp.Failure() {
			log.Errorf("[LcHomeIssuesForProject] orgId:%d, userId:%d, appId:%v, GetAppAuthWithoutCollaborator failure:%v", orgId, currentUserId, projectBo.AppId, optAuthResp.Error())
			return nil, optAuthResp.Error()
		}
		appAuthInfo = &optAuthResp.Data
		isAdmin = appAuthInfo.HasAppRootPermission || appAuthInfo.SysAdmin || appAuthInfo.OrgOwner || appAuthInfo.AppOwner || currentUserId == projectBo.Owner
		log.Infof("[LcHomeIssuesForProject] orgId:%d, userId:%d, isAdmin:%v, app auth: %v", orgId, currentUserId, isAdmin, appAuthInfo)
	}

	needDeleteData, err2 := setIsNeedRefresh(input, inputAppId)
	if err2 != nil {
		log.Errorf("[LcHomeIssuesForProject] setIsNeedRefresh orgId:%d, userId:%d, LessIssueList failure:%v", orgId, currentUserId, err)
		return nil, err2
	}

	// 构造筛选条件
	lessReq, err2 := getConditionForProject(orgId, currentUserId, inputAppId, input, isAdmin, isInnerSuper, projectBo)
	if err2 != nil {
		return nil, err2
	}

	log.Infof("[LcHomeIssuesForProject] orgId:%d, userId:%d, LessIssueListReq conds:%v", orgId, currentUserId, lessReq.Conds)

	addOrder(input, isInnerSuper)

	// 拉取数据
	lessReqParam := formvo.LessIssueListReq{
		Condition:   *lessReq,
		AppId:       summaryAppId,
		RedirectIds: []int64{projectBo.AppId},
		OrgId:       orgId,
		UserId:      currentUserId,
		Page:        int64(page),
		Size:        int64(size),
		Export:      true,
		Orders:      input.LessOrder,
		//Orders:         orders,
		TableId:        tableId,
		NeedRefColumn:  true,
		NeedDeleteData: needDeleteData,
		NeedChangeId:   true,
		NeedTotal:      input.IsNeedTotal,
	}
	if input.FilterColumns != nil && len(input.FilterColumns) > 0 {
		filterColumns := []string{
			lc_helper.ConvertToFilterColumn(consts.BasicFieldId),
			lc_helper.ConvertToFilterColumn(consts.BasicFieldIssueId),
			lc_helper.ConvertToFilterColumn(consts.BasicFieldOrgId),
			lc_helper.ConvertToFilterColumn(consts.BasicFieldAppId),
			lc_helper.ConvertToFilterColumn(consts.BasicFieldTableId),
			lc_helper.ConvertToFilterColumn(consts.BasicFieldAuditStatusDetail),
			lc_helper.ConvertToFilterColumn(consts.BasicFieldUpdateTime),
		}
		for _, column := range input.FilterColumns {
			filterColumns = append(filterColumns, lc_helper.ConvertToFilterColumn(column))
		}
		lessReqParam.FilterColumns = slice.SliceUniqueString(filterColumns)
	}
	log.Infof("[LcHomeIssuesForProject] orgId:%d, userId:%d, LessIssueListReq req: %s", orgId, currentUserId, json.ToJsonIgnoreError(lessReqParam))
	// 鉴权时，检查是否是系统管理员
	appAuthInfoStr := ""
	if appAuthInfo != nil && !appAuthInfo.HasAppRootPermission {
		tableIdStr := cast.ToString(tableId)
		if !appAuthInfo.HasAllFieldAuthOfTable(tableIdStr) {
			appAuthInfoStr = json.ToJsonIgnoreError(appAuthInfo)
		}
	}

	lessResp, _ := domain.GetIssueListWithRefColumn(lessReqParam, nil, appAuthInfoStr)
	if lessResp.Failure() {
		log.Errorf("[LcHomeIssuesForProject] orgId:%d, userId:%d, LessIssueList failure:%v", orgId, currentUserId, lessResp.Error())
		return nil, lessResp.Error()
	}

	lastUpdateTime := input.LastUpdateTime
	if lessResp.Data.LastUpdateTime > lastUpdateTime {
		lastUpdateTime = lessResp.Data.LastUpdateTime
	}

	if int(lessResp.Data.Count) < size {
		input.IsRefreshAll = false
	}

	return issueListToJsonString(lessResp, input), nil
}

func addOrder(input *projectvo.HomeIssueInfoReq, isInnerSuper bool) {
	if input.IsRefreshAll {
		if isInnerSuper {
			hasOrderOrder := false
			hasIdOrder := false
			for _, order := range input.LessOrder {
				if order.Column == consts.BasicFieldOrder {
					hasOrderOrder = true
				}
				if order.Column == consts.BasicFieldId {
					hasIdOrder = true
				}
			}
			if !hasOrderOrder {
				input.LessOrder = append(input.LessOrder, &vo.LessOrder{Asc: false, Column: consts.BasicFieldOrder})
			}
			if !hasIdOrder {
				input.LessOrder = append(input.LessOrder, &vo.LessOrder{Asc: false, Column: consts.BasicFieldId})
			}
		} else {
			hasIdOrder := false
			for _, order := range input.LessOrder {
				if order.Column == consts.BasicFieldId {
					hasIdOrder = true
				}
			}
			if !hasIdOrder {
				input.LessOrder = append(input.LessOrder, &vo.LessOrder{Asc: false, Column: consts.BasicFieldId})
			}
		}
	}
}

// setIsNeedRefresh 判断是否需要大刷还是增量获取数据
func setIsNeedRefresh(input *projectvo.HomeIssueInfoReq, appId int64) (bool, errs.SystemErrorInfo) {
	if input.LastUpdateTime == "" || input.LastPermissionUpdateTime == "" {
		input.IsRefreshAll = true
	}
	if input.IsRefreshAll {
		input.LastPermissionUpdateTime = cast.ToString(time.Now().Unix())
		return false, nil
	}

	if appId > 0 {
		resp := permissionfacade.GetPermissionUpdateTime(permissionvo.GetPermissionUpdateTimeReq{AppId: appId})
		if resp.Failure() {
			log.Errorf("[setIsNeedRefresh] GetPermissionUpdateTime appId:%v, err:%v", appId, resp.Error())
			return false, resp.Error()
		}
		// 如果权限修改了，则需要大刷
		if resp.Data.UpdateTime > cast.ToInt64(input.LastPermissionUpdateTime) {
			input.LastPermissionUpdateTime = cast.ToString(resp.Data.UpdateTime)
			input.IsRefreshAll = true
			return false, nil
		}
	}

	return true, nil
}

func getConditionForProject(orgId, currentUserId, inputAppId int64, input *projectvo.HomeIssueInfoReq, isAdmin, isInnerSuper bool, projectBo *bo.ProjectBo) (*vo.LessCondsData, errs.SystemErrorInfo) {
	// 构造筛选条件
	lessReq := &vo.LessCondsData{Type: "and"}
	// 筛选条件：前端传入的条件，先忽略
	if input.LessConds != nil {
		if input.LessConds.Type == "and" {
			lessReq = input.LessConds
		} else {
			lessReq.Conds = append(lessReq.Conds, input.LessConds)
		}
	}

	lessReq.Conds = append(lessReq.Conds, getRefreshCondition(input)...)

	lessReq.Conds = append(lessReq.Conds, &vo.LessCondsData{
		Type:   "equal",
		Value:  *input.TableID,
		Column: consts.BasicFieldTableId,
	})

	if !isAdmin && projectBo.TemplateFlag != consts.TemplateTrue {
		// 非公开项目：如果是我参与的项目展示所有的，如果不是，则只展示我协作的任务
		isMemberResp := appfacade.IsAppMember(appvo.IsAppMemberReq{
			AppId:  inputAppId,
			OrgId:  orgId,
			UserId: currentUserId,
		})
		if isMemberResp.Failure() {
			log.Error(isMemberResp.Error())
			return nil, isMemberResp.Error()
		}

		if !isMemberResp.Data {
			// 如果不是项目成员，走协作人逻辑
			deptIdsWithParentIdResp := orgfacade.GetUserDeptIdsWithParentId(orgvo.GetUserDeptIdsWithParentIdReq{
				OrgId:  orgId,
				UserId: currentUserId,
			})
			if deptIdsWithParentIdResp.Failure() {
				log.Errorf("[LcHomeIssuesForProject] GetUserDeptIdsWithParentId err:%v, orgId:%v, userId:%v, appId:%v",
					deptIdsWithParentIdResp.Error(), orgId, currentUserId, projectBo.AppId)
				return nil, deptIdsWithParentIdResp.Error()
			}
			deptIds := deptIdsWithParentIdResp.Data.DeptIds
			// 筛选条件：我协作的任务
			lessReq.Conds = append(lessReq.Conds, domain.GetCollaborateDataQueryConds(orgId, currentUserId, deptIds))
		}
	}

	// 筛选条件：处理任务状态
	if projectBo.ProjectTypeId == consts.ProjectTypeNormalId {
		lessReq = domain.ConvertIssueStatusFilterReqForNormalProject(lessReq)
	}

	return lessReq, nil
}

// LcHomeIssuesForProject 项目下的任务筛选
// @param isInnerSuper 如果是 openAPI，则无需校验具体权限，并且不受字段权限控制（可以看到所有字段）
func LcHomeIssuesForIssue(orgId, currentUserId, appId, tableId, issueId, todoId int64, isInnerSuper bool) (*projectvo.IssueDetailInfo, errs.SystemErrorInfo) {
	var errSys errs.SystemErrorInfo

	// 先拿issue的信息，因为可能没拿到appId和tableId
	issueBo, errSys := domain.GetIssueInfoLc(orgId, currentUserId, issueId)
	if errSys != nil {
		log.Errorf("[LcHomeIssuesForIssue] GetIssueInfosLc err:%v, orgId:%v, issueId:%v", errSys, orgId, issueId)
		return nil, errSys
	}
	if tableId <= 0 {
		tableId = issueBo.TableId
	}

	// 拿table的信息
	var tableMetaData *projectvo.TableMetaData
	if tableId > 0 {
		tableMetaData, errSys = domain.GetTableByTableId(orgId, currentUserId, tableId)
		if errSys != nil {
			return nil, errs.TableNotExist
		}
	}

	// 再拿project信息
	var projectBo *bo.ProjectBo
	if issueBo.ProjectId > 0 {
		projectBo, errSys = domain.GetProjectSimple(orgId, issueBo.ProjectId)
		if errSys != nil {
			return nil, errs.IssueNotExist
		}
		if appId <= 0 {
			appId = projectBo.AppId
		}
	}
	// 拿不到appId就用汇总表的
	if appId <= 0 {
		orgInfoResp := orgfacade.OrganizationInfo(orgvo.OrganizationInfoReqVo{
			OrgId:  orgId,
			UserId: currentUserId,
		})
		if orgInfoResp.Failure() {
			log.Errorf("[LcHomeIssuesForIssue] orgId:%v, curUserId:%v, OrganizationInfo failure:%v", orgId, currentUserId, orgInfoResp.Error())
			return nil, orgInfoResp.Error()
		}
		orgRemarkObj := orgvo.OrgRemarkConfigType{}
		if err := json.FromJson(orgInfoResp.OrganizationInfo.Remark, &orgRemarkObj); err != nil {
			log.Errorf("[LcHomeIssuesForIssue] orgId:%v, curUserId:%v, OrgInfo FromJson failure:%v", orgId, currentUserId, err)
			return nil, errs.BuildSystemErrorInfo(errs.JSONConvertError, err)
		}
		appId = orgRemarkObj.OrgSummaryTableAppId
	}

	// 获取用户部门信息
	//deptIdsResp := orgfacade.GetUserDeptIds(orgvo.GetUserDeptIdsReq{
	//	OrgId:  orgId,
	//	UserId: currentUserId,
	//})
	//if deptIdsResp.Failure() {
	//	log.Errorf("[LcHomeIssuesForIssue] orgId:%v, curUserId:%v, GetUserDeptIds failure:%v", orgId, currentUserId, deptIdsResp.Error())
	//	return nil, deptIdsResp.Error()
	//}
	//log.Infof("[LcHomeIssuesForIssue] orgId:%d, userId:%d", orgId, currentUserId)

	// 是否管理员。isInnerSuper 为 true 表示来源为 openAPI，无需权限校验
	isAdmin := isInnerSuper || domain.CheckIsAdmin(orgId, currentUserId, appId) || todoId > 0 // 先简单处理，有todoId则不校验权限
	log.Infof("[LcHomeIssuesForIssue] orgId:%d, userId:%d, isAdmin:%v", orgId, currentUserId, isAdmin)

	// 构造筛选条件
	lessReq := domain.GetIssueIdDetailQueryConds(issueId, issueBo.Path)

	if !isAdmin && projectBo != nil {
		// 判断是否是项目成员
		isMemberResp := appfacade.IsAppMember(appvo.IsAppMemberReq{
			AppId:  appId,
			OrgId:  orgId,
			UserId: currentUserId,
		})
		if isMemberResp.Failure() {
			log.Errorf("[LcHomeIssuesForIssue] IsAppMember err:%v, orgId:%v, userId:%v, issueId:%v",
				isMemberResp.Error(), orgId, currentUserId, issueId)
			return nil, isMemberResp.Error()
		}
		// 不是项目成员，需要走协作人
		if !isMemberResp.Data {
			deptIdsWithParentIdResp := orgfacade.GetUserDeptIdsWithParentId(orgvo.GetUserDeptIdsWithParentIdReq{
				OrgId:  orgId,
				UserId: currentUserId,
			})
			if deptIdsWithParentIdResp.Failure() {
				log.Errorf("[LcHomeIssuesForIssue] orgId:%v, curUserId:%v, GetUserDeptIds failure:%v", orgId, currentUserId, deptIdsWithParentIdResp.Error())
				return nil, deptIdsWithParentIdResp.Error()
			}
			deptIds := deptIdsWithParentIdResp.Data.DeptIds
			// 筛选条件：我协作的任务
			lessReq.Conds = append(lessReq.Conds, domain.GetCollaborateDataQueryConds(orgId, currentUserId, deptIds))
		}
		// 获取我参与的所有项目
		//myProjectIds, errSys := domain.GetAllMyProjectIdsWithDeptIds(orgId, currentUserId, deptIdsResp.Data.DeptIds, false)
		//if errSys != nil {
		//	log.Errorf("[LcHomeIssuesForIssue] orgId:%d, userId:%d, deptIds:%v, GetAllMyProjectIdsWithDeptIds failure:%v", orgId, currentUserId, deptIdsResp.Data.DeptIds, errSys)
		//	return nil, errSys
		//}
		// 不是我参与的项目，判断协作人权限
		//if has, _ := slice.Contain(myProjectIds, projectBo.Id); !has {
		//	// 筛选条件：我协作的任务
		//	lessReq.Conds = append(lessReq.Conds, domain.GetCollaborateDataQueryConds(orgId, currentUserId, deptIdsResp.Data.DeptIds))
		//}
	}

	// 筛选条件：处理任务状态转换
	lessReq = domain.ConvertIssueStatusFilterReqForAll(lessReq)

	// 拉取数据
	lessReqParam := formvo.LessIssueListReq{
		Condition:     lessReq,
		AppId:         appId,
		OrgId:         orgId,
		UserId:        currentUserId,
		Page:          0,
		Size:          0,
		TableId:       tableId,
		NeedRefColumn: true,
		AggNoLimit:    true,
		NeedChangeId:  true,
	}
	log.Infof("[LcHomeIssuesForIssue] orgId:%d, userId:%d, LessIssueListReq req: %s", orgId, currentUserId, json.ToJsonIgnoreError(lessReqParam))
	lessResp := domain.GetRowsExpand(&lessReqParam)
	if lessResp.Failure() {
		log.Errorf("[LcHomeIssuesForIssue] orgId:%d, userId:%d, LessIssueList failure:%v", orgId, currentUserId, lessResp.Error())
		return nil, lessResp.Error()
	}

	for _, info := range lessResp.List {
		// 审批确认项目，给请求的任务本体补上催办时间，确认人的状态
		if projectBo != nil && cast.ToInt64(info[consts.BasicFieldId]) == issueId {
			// 组装：催办时间
			if lastUrgeTimeForOwner, err := domain.GetLastUrgeIssue(orgId, projectBo.Id, issueId, sconsts.CacheUrgeIssue); err == nil {
				info["lastUrgeTimeForOwner"] = lastUrgeTimeForOwner
			}
			if lastUrgeTimeForAuditor, err := domain.GetLastUrgeIssue(orgId, projectBo.Id, issueId, sconsts.CacheUrgeIssueAudit); err == nil {
				info["lastUrgeTimeForAuditor"] = lastUrgeTimeForAuditor
			}
		}
	}

	var projectInfo *projectvo.ProjectMetaInfo
	if projectBo != nil {
		projectInfo = &projectvo.ProjectMetaInfo{
			Id:            int(projectBo.Id),
			Name:          projectBo.Name,
			ProjectTypeId: int(projectBo.ProjectTypeId),
			IsFilling:     projectBo.IsFiling,
		}
	}
	var tableInfo *projectvo.TableSimpleInfo
	if tableMetaData != nil {
		tableInfo = &projectvo.TableSimpleInfo{
			Id:   cast.ToString(tableMetaData.TableId),
			Name: tableMetaData.Name,
		}
	}
	return &projectvo.IssueDetailInfo{
		Project:   projectInfo,
		Table:     tableInfo,
		Data:      lessResp.List,
		UserDepts: lessResp.UserDepts,
	}, nil
}
