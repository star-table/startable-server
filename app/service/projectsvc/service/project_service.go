package service

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"sync"
	"time"

	msgPb "gitea.bjx.cloud/LessCode/interface/golang/msg/v1"
	sdk_const "gitea.bjx.cloud/allstar/platform-sdk/consts"
	"github.com/star-table/startable-server/app/facade/common/report"
	"github.com/star-table/startable-server/common/core/threadlocal"
	"github.com/star-table/startable-server/common/model/vo/commonvo"

	tablePb "gitea.bjx.cloud/LessCode/interface/golang/table/v1"
	"github.com/star-table/startable-server/app/facade/appfacade"
	"github.com/star-table/startable-server/app/facade/formfacade"
	"github.com/star-table/startable-server/app/facade/idfacade"
	"github.com/star-table/startable-server/app/facade/orgfacade"
	"github.com/star-table/startable-server/app/facade/resourcefacade"
	"github.com/star-table/startable-server/app/facade/userfacade"
	projectsvcConsts "github.com/star-table/startable-server/app/service/projectsvc/consts"
	"github.com/star-table/startable-server/app/service/projectsvc/dao"
	"github.com/star-table/startable-server/app/service/projectsvc/domain"
	"github.com/star-table/startable-server/app/service/projectsvc/po"
	"github.com/star-table/startable-server/common/core/consts"
	"github.com/star-table/startable-server/common/core/errs"
	"github.com/star-table/startable-server/common/core/types"
	"github.com/star-table/startable-server/common/core/util"
	"github.com/star-table/startable-server/common/core/util/asyn"
	"github.com/star-table/startable-server/common/core/util/cond"
	"github.com/star-table/startable-server/common/core/util/copyer"
	"github.com/star-table/startable-server/common/core/util/format"
	"github.com/star-table/startable-server/common/core/util/json"
	lang2 "github.com/star-table/startable-server/common/core/util/lang"
	"github.com/star-table/startable-server/common/core/util/maps"
	"github.com/star-table/startable-server/common/core/util/slice"
	"github.com/star-table/startable-server/common/core/util/stack"
	"github.com/star-table/startable-server/common/core/util/strs"
	"github.com/star-table/startable-server/common/core/util/times"
	"github.com/star-table/startable-server/common/extra/lc_helper"
	"github.com/star-table/startable-server/common/library/db/mysql"
	"github.com/star-table/startable-server/common/model/bo"
	"github.com/star-table/startable-server/common/model/vo"
	"github.com/star-table/startable-server/common/model/vo/appvo"
	"github.com/star-table/startable-server/common/model/vo/formvo"
	"github.com/star-table/startable-server/common/model/vo/orgvo"
	"github.com/star-table/startable-server/common/model/vo/projectvo"
	"github.com/star-table/startable-server/common/model/vo/resourcevo"
	"github.com/star-table/startable-server/common/model/vo/uservo"
	"github.com/spf13/cast"
	"gopkg.in/fatih/set.v0"
	"upper.io/db.v3"
	"upper.io/db.v3/lib/sqlbuilder"
)

func Projects(reqVo projectvo.ProjectsRepVo) (*vo.ProjectList, errs.SystemErrorInfo) {
	orgId := reqVo.OrgId
	currentUserId := reqVo.UserId
	page := reqVo.Page
	size := reqVo.Size
	params := reqVo.ProjectExtraBody.Params
	order := reqVo.ProjectExtraBody.Order
	input := reqVo.ProjectExtraBody.Input
	projectTypeIds := reqVo.ProjectExtraBody.ProjectTypeIds

	log.Infof(consts.UserLoginSentence, currentUserId, orgId)

	var joinParams db.Cond
	var joinErr errs.SystemErrorInfo
	if len(params) == 0 {
		joinParams, joinErr = GetProjectCondAssemblyByInput(input, currentUserId, orgId, projectTypeIds)
	} else {
		joinParams, joinErr = GetProjectCondAssemblyByParam(params, currentUserId, orgId)
	}
	if joinErr != nil {
		log.Error(joinErr)
		return nil, joinErr
	}
	joinParams[consts.TcOrgId] = orgId
	joinParams[consts.TcTemplateFlag] = consts.TemplateFalse
	var union *db.Union = nil
	if input != nil && input.Name != nil {
		name := strings.ToLower(*input.Name)
		union = db.Or(db.Cond{
			consts.TcName: db.Like("%" + name + "%"),
		}).Or(db.Cond{
			consts.TcPreCode: db.Like("%" + name + "%"),
		}).Or(db.Cond{
			consts.TcName: db.Like("%" + *input.Name + "%"),
		}).Or(db.Cond{
			consts.TcPreCode: db.Like("%" + *input.Name + "%"),
		})
	}
	//获取我收藏的项目 已经不在极星维护了，下面这一段没啥用了

	//获取项目列表
	var totalNumberOfEntries int64
	entities, totalNumberOfEntries, err := domain.GetProjectList(currentUserId, joinParams, union, order, size, page)

	if err != nil {
		return nil, errs.BuildSystemErrorInfo(errs.ProjectDomainError, err)
	}

	resultList := &[]*vo.Project{}
	copyer.Copy(entities, resultList)

	// 获取冗余信息
	if !reqVo.ProjectExtraBody.NoNeedRedundancyInfo {
		if errSys := getRedundancyInfo(resultList, orgId, reqVo.SourceChannel, currentUserId); errSys != nil {
			return nil, errs.BuildSystemErrorInfo(errs.ProjectDomainError)
		}
	}

	result := &vo.ProjectList{}
	result.Total = totalNumberOfEntries
	result.List = *resultList
	return result, nil
}

// 获取冗余信息
func getRedundancyInfo(resultList *[]*vo.Project, orgId int64, sourceChannel string, currentUserId int64) errs.SystemErrorInfo {
	resourceIds := []int64{}
	creatorIds := []int64{}
	projectIds := []int64{}
	//agileProjectIds := []int64{}
	appIds := []int64{}

	if len(*resultList) != 0 {
		for _, v := range *resultList {
			resourceIds = append(resourceIds, v.ResourceID)
			creatorIds = append(creatorIds, v.Creator)
			projectIds = append(projectIds, v.ID)
			appId, parseErr := strconv.ParseInt(v.AppID, 10, 64)
			if parseErr != nil {
				log.Errorf("[getRedundancyInfo] appId convert err:%v, orgId:%v, appId:%v", parseErr, orgId, appId)
				return errs.TypeConvertError
			}
			if appId != 0 {
				appIds = append(appIds, appId)
			}
		}

		type midType struct {
			creatorInfo           map[int64]bo.UserIDInfoBo
			ownerInfo             map[int64][]bo.UserIDInfoBo
			participantInfo       map[int64][]bo.UserIDInfoBo
			followerInfo          map[int64][]bo.UserIDInfoBo
			resourceByPath        map[int64]bo.ResourceBo
			issueStat             map[int64]bo.IssueStatistic
			projectTypeLocalCache maps.LocalMap
			projectDetailById     maps.LocalMap
			iterationStat         map[int64]bo.IssueStatistic

			appIcon map[int64]string
		}

		handlerFuncList := make([]func(midInfo *midType) errs.SystemErrorInfo, 0, 8)
		midInfo := &midType{}

		handlerFuncList = append(handlerFuncList, func(midInfo *midType) errs.SystemErrorInfo {
			creatorIds = slice.SliceUniqueInt64(creatorIds)
			ownerInfo, participantInfo, followerInfo, creatorInfo, err := domain.GetProjectMemberInfo(projectIds, orgId, creatorIds)
			if err != nil {
				log.Error(err)
				return err
			}
			midInfo.creatorInfo = creatorInfo
			midInfo.ownerInfo = ownerInfo
			midInfo.followerInfo = followerInfo
			midInfo.participantInfo = participantInfo
			return nil
		})

		handlerFuncList = append(handlerFuncList, func(midInfo *midType) errs.SystemErrorInfo {
			//资源列表
			resourceIds = slice.SliceUniqueInt64(resourceIds)
			resourcesRespVo := resourcefacade.GetResourceById(resourcevo.GetResourceByIdReqVo{GetResourceByIdReqBody: resourcevo.GetResourceByIdReqBody{ResourceIds: resourceIds}})
			if resourcesRespVo.Failure() {
				log.Error(resourcesRespVo.Message)
				return resourcesRespVo.Error()
			}
			resourceEntities := resourcesRespVo.ResourceBos
			resourceByPath := map[int64]bo.ResourceBo{}
			for _, v := range resourceEntities {
				resourceByPath[v.Id] = v
			}
			midInfo.resourceByPath = resourceByPath
			return nil
		})

		handlerFuncList = append(handlerFuncList, func(midInfo *midType) errs.SystemErrorInfo {
			//获取projectType localCache
			projectTypeList, err := domain.GetProjectTypeList(orgId)
			if err != nil {
				log.Error(err)
				return err
			}
			projectTypeLocalCache := maps.NewMap("Id", projectTypeList)
			midInfo.projectTypeLocalCache = projectTypeLocalCache
			return nil
		})

		handlerFuncList = append(handlerFuncList, func(midInfo *midType) errs.SystemErrorInfo {
			//获取项目详情
			projectDetails, detailErr := domain.GetProjectDetails(orgId, projectIds)
			if detailErr != nil {
				log.Error(detailErr)
				return detailErr
			}
			projectDetailById := maps.NewMap("ProjectId", projectDetails)
			midInfo.projectDetailById = projectDetailById
			return nil
		})

		handlerFuncList = append(handlerFuncList, func(midInfo *midType) errs.SystemErrorInfo {
			resp := appfacade.GetAppInfoList(appvo.GetAppInfoListReq{
				OrgId:  orgId,
				AppIds: appIds,
			})
			if resp.Failure() {
				log.Error(resp.Error())
				return resp.Error()
			}
			iconMap := make(map[int64]string, 0)
			for _, datum := range resp.Data {
				iconMap[datum.Id] = datum.Icon
			}
			midInfo.appIcon = iconMap
			return nil
		})

		var wg sync.WaitGroup
		wg.Add(len(handlerFuncList))
		businessErrChan := make(chan errs.SystemErrorInfo, 100)
		defer close(businessErrChan)
		for _, handlerFunc := range handlerFuncList {
			currentFunc := handlerFunc
			go func() {
				defer func() {
					if r := recover(); r != nil {
						log.Error(errs.BuildSystemErrorInfoWithPanicRecover(r, stack.GetStack()))
					}
				}()
				defer wg.Done()
				if err := currentFunc(midInfo); err != nil {
					businessErrChan <- err
				}
			}()
		}
		wg.Wait()
		select {
		case err := <-businessErrChan:
			log.Errorf("[getRedundancyInfo] handle handlerFuncList err: %v", err)
			return err
		default:
		}

		dealResultList(resultList, midInfo.ownerInfo, midInfo.participantInfo, midInfo.followerInfo, midInfo.resourceByPath,
			midInfo.issueStat, midInfo.projectTypeLocalCache, midInfo.projectDetailById, midInfo.iterationStat, midInfo.appIcon)
		addCreatorInfo(resultList, midInfo.creatorInfo)
	}

	return nil
}

func GetProjectCondAssemblyByParam(params map[string]interface{}, currentUserId int64, orgId int64) (db.Cond, errs.SystemErrorInfo) {
	var relationType interface{}
	if _, ok := params["relation_type"]; ok {
		if val, ok := params["relation_type"].(map[string]interface{}); ok {
			if val["type"] != nil && val["value"] != nil {
				relationType = val["value"]
			}
		} else {
			relationType = params["relation_type"]
		}
		delete(params, "relation_type")
	}
	var relateType int64 = 0

	converRelateType(&relateType, relationType)

	condParam, err := cond.HandleParams(params)
	if err != nil {
		return nil, errs.BuildSystemErrorInfo(errs.ConditionHandleError, err)
	}
	switch relateType {
	case 0:
		//所有
	case 1:
		//我发起的
		condParam[consts.TcCreator] = currentUserId
	case 2:
		//我负责的
		condParam[consts.TcOwner] = db.Eq(currentUserId)
	case 3:
		//我参与的
		need, err := domain.GetParticipantMembers(orgId, currentUserId)
		if err != nil {
			return nil, errs.BuildSystemErrorInfo(errs.ConditionHandleError, err)
		}
		condParam[consts.TcId+" In"] = db.In(need)
	}
	//默认查询没有被删除的
	condParam[consts.TcIsDelete] = consts.AppIsNoDelete

	if params["id"] != nil {
		condParam[consts.TcId] = db.In(db.Raw("select p.id as id from ppm_pro_project p where p.id = ? and (p.public_status = 1 or p.id in (SELECT DISTINCT pr.project_id FROM ppm_pro_project_relation pr WHERE pr.relation_id = ? AND relation_type in (1,2) AND pr.is_delete = 2)) and p.is_delete = 2", params["id"], currentUserId))
	} else {
		condParam[consts.TcId] = db.In(db.Raw("select p.id as id from ppm_pro_project p where (p.public_status = 1 or p.id in (SELECT DISTINCT pr.project_id FROM ppm_pro_project_relation pr WHERE pr.relation_id = ? AND relation_type in (1,2) AND pr.is_delete = 2)) and p.is_delete = 2", currentUserId))
	}

	return condParam, nil
}

func converRelateType(relateType *int64, relationType interface{}) {
	*relateType = cast.ToInt64(relationType)
}

// 状态类型,1未开始2进行中3已完成4未完成
func condStatusAssembly(cond db.Cond, orgId int64, status int) errs.SystemErrorInfo {
	var statusIds []int64 = nil
	if status == 4 {
		statusIds = append(statusIds, consts.StatusRunning.ID, consts.StatusNotStart.ID)
	} else if status == 3 {
		statusIds = append(statusIds, consts.StatusComplete.ID)
	} else if status == 1 {
		statusIds = append(statusIds, consts.StatusNotStart.ID)
	} else if status == 2 {
		statusIds = append(statusIds, consts.StatusRunning.ID)
	}
	cond[consts.TcStatus+" in"] = db.In(statusIds)
	return nil
}

func GetProjectCondAssemblyByInput(input *vo.ProjectsReq, currentUserId int64, orgId int64, projectTypeIds []int64) (db.Cond, errs.SystemErrorInfo) {
	condParam := make(db.Cond)

	if input == nil {
		input = &vo.ProjectsReq{}
	}

	//拿到当前用户的管理员flag
	adminFlag := uservo.GetUserAuthorityData{
		OrgId:                orgId,
		UserId:               0,
		IsOrgOwner:           true,
		IsSysAdmin:           true,
		IsSubAdmin:           true,
		HasDeptOptAuth:       true,
		HasRoleDeptAuth:      true,
		HasAppPackageOptAuth: true,
	}
	if currentUserId != 0 {
		manageAuthInfoResp := userfacade.GetUserAuthority(orgId, currentUserId)
		if manageAuthInfoResp.Failure() {
			log.Error(manageAuthInfoResp.Message)
			return condParam, manageAuthInfoResp.Error()
		}
		adminFlag = manageAuthInfoResp.Data
	}

	deptResp := orgfacade.GetUserDeptIdsWithParentId(orgvo.GetUserDeptIdsWithParentIdReq{
		OrgId:  orgId,
		UserId: currentUserId,
	})
	if deptResp.Failure() {
		log.Error(deptResp.Error())
		return nil, deptResp.Error()
	}
	allDeptIds := []int64{0}
	allDeptIds = append(allDeptIds, deptResp.Data.DeptIds...)

	//获取协作人项目
	orgInfoResp := orgfacade.OrganizationInfo(orgvo.OrganizationInfoReqVo{
		OrgId:  orgId,
		UserId: 0,
	})
	if orgInfoResp.Failure() {
		log.Error(orgInfoResp.Error())
		return nil, orgInfoResp.Error()
	}
	orgRemarkObj := &orgvo.OrgRemarkConfigType{}
	oriErr := json.FromJson(orgInfoResp.OrganizationInfo.Remark, orgRemarkObj)
	if oriErr != nil {
		log.Errorf("[GetProjectCondAssemblyByInput] orgId: %d, err: %v", orgId, oriErr)
		return nil, errs.BuildSystemErrorInfo(errs.JSONConvertError, errs.OrgInitError)
	}

	allValues := []string{fmt.Sprintf("U_%d", currentUserId)}
	for _, id := range allDeptIds {
		allValues = append(allValues, fmt.Sprintf("D_%d", id))
	}

	req := &tablePb.ListRawRequest{
		FilterColumns: []string{
			lc_helper.ConvertToFilterColumn(consts.BasicFieldProjectId),
		},
		Condition: &tablePb.Condition{
			Type: tablePb.ConditionType_and,
			Conditions: domain.GetNoRecycleCondition(
				domain.GetRowsCondition("collaborators && ARRAY['"+strings.Join(allValues, "','")+"']", tablePb.ConditionType_raw_sql, nil, allValues),
			),
		},
		Groups: []string{
			// group by 用别名也可
			lc_helper.ConvertToFilterColumn(consts.BasicFieldProjectId),
		},
	}
	issueResp, err := domain.GetRawRows(orgId, currentUserId, req)
	if err != nil {
		log.Error(err)
		return nil, err
	}

	var allCollaborateProjectIds []int64
	for _, m := range issueResp.Data {
		if projectIdInterface, ok := m[consts.BasicFieldProjectId]; ok {
			projectId := cast.ToInt64(projectIdInterface)
			allCollaborateProjectIds = append(allCollaborateProjectIds, projectId)
		}
	}
	log.Infof("[GetProjectList] orgId: %v, userId: %v, allCollaborateProjectIds: %v", orgId, currentUserId, allCollaborateProjectIds)

	err2 := dealInputRelateType(condParam, input, currentUserId, orgId, adminFlag, allCollaborateProjectIds)
	if err2 != nil {
		return nil, errs.BuildSystemErrorInfo(errs.ConditionHandleError, err2)
	}

	//默认查询没有被删除的
	condParam[consts.TcIsDelete] = consts.AppIsNoDelete
	if input.ID != nil {
		condParam[consts.TcId] = input.ID
	}
	if input.ProjectIds != nil {
		condParam[consts.TcId] = db.In(input.ProjectIds)
	}
	if input.AppIds != nil {
		condParam[consts.TcAppId] = db.In(input.AppIds)
	}
	//if input.Name != nil {
	//	condParam[consts.TcName] = db.Like("%" + *input.Name + "%")
	//}
	if input.Status != nil {
		condParam[consts.TcStatus] = input.Status
	}
	if input.StatusType != nil {
		err := condStatusAssembly(condParam, orgId, *input.StatusType)
		if err != nil {
			log.Error(err)
			return nil, errs.BuildSystemErrorInfo(errs.ConditionHandleError, err)
		}
	}
	if len(input.OwnerIds) > 0 {
		condParam[consts.TcOwner+" "] = db.In(input.OwnerIds)
	}
	if len(input.CreatorIds) > 0 {
		condParam[consts.TcCreator+" "] = db.In(input.CreatorIds)
	}
	if input.Owner != nil {
		condParam[consts.TcOwner] = input.Owner
	}
	if input.ProjectTypeID != nil {
		condParam[consts.TcProjectTypeId] = input.ProjectTypeID
	}
	if len(projectTypeIds) != 0 {
		condParam[consts.TcProjectTypeId] = db.In(projectTypeIds)
	}

	if input.IsFiling != nil {
		if *input.IsFiling == 1 || *input.IsFiling == 2 {
			condParam[consts.TcIsFiling] = input.IsFiling
		}
	} else {
		//默认查未归档
		condParam[consts.TcIsFiling] = consts.AppIsNotFilling
	}
	if input.PriorityID != nil {
		condParam[consts.TcPriorityId] = input.PriorityID
	}
	if input.PlanEndTime != nil {
		condParam[consts.TcPlanEndTime] = db.Lte(input.PlanEndTime.String())
	}
	if input.PlanStartTime != nil {
		condParam[consts.TcPlanStartTime] = db.Gte(input.PlanStartTime.String())
	}

	args := []interface{}{orgId}
	sql := "select p.id as id from ppm_pro_project p where p.org_id = ?"
	if input.ID != nil {
		sql += " and p.id = ?"
		args = append(args, *input.ID)
	}

	//不是超级管理员莫得私有项目查看权
	//if !adminFlag.IsSysAdmin && !adminFlag.IsSubAdmin {
	//	sql += " and (p.public_status = 1 or p.id in ? or p.id in (SELECT DISTINCT pr.project_id FROM ppm_pro_project_relation pr WHERE ((pr.relation_id = ? AND relation_type in (1,2)) or (pr.relation_type=25 and pr.relation_id in ?)) AND pr.is_delete = 2))"
	//	args = append(args, allCollaborateProjectIds, currentUserId, allDeptIds)
	//}

	condParam[consts.TcId] = db.In(db.Raw(sql, args...))

	if len(input.Participants) > 0 || len(input.ParticipantDeptIds) > 0 {
		hasManager := false
		if len(input.Participants) > 0 {
			//判断成员中是否有管理员
			resp := userfacade.GetUsersCouldManage(orgId, -1)
			if resp.Failure() {
				log.Error(resp.Error())
				return nil, resp.Error()
			}
			for _, infoBo := range resp.Data.List {
				if ok, _ := slice.Contain(input.Participants, infoBo.Id); ok {
					hasManager = true
					break
				}
			}
		}

		if !hasManager {
			var deptIds []int64
			deptIds = append(deptIds, 0)
			if len(input.ParticipantDeptIds) > 0 {
				deptIds = append(deptIds, input.ParticipantDeptIds...)
			}
			if len(input.Participants) > 0 {
				//查找参与者对应的部门
				deptMapResp := orgfacade.GetUserDeptIdsBatch(&orgvo.GetUserDeptIdsBatchReq{
					OrgId:   orgId,
					UserIds: input.Participants,
				})
				if deptMapResp.Failure() {
					log.Error(deptMapResp.Error())
					return nil, deptMapResp.Error()
				}
				for _, int64s := range deptMapResp.Data.Data {
					deptIds = append(deptIds, int64s...)
				}
			}
			deptIds = slice.SliceUniqueInt64(deptIds)
			var sql interface{}
			if len(input.Participants) > 0 && len(deptIds) > 0 {
				sql = db.Raw("select distinct(project_id) as id from ppm_pro_project_relation where is_delete = 2 and ((relation_type=2 and relation_id in ?) or (relation_type=25 and relation_id in ?))", input.Participants, deptIds)
			} else if len(input.Participants) > 0 && len(deptIds) == 0 {
				sql = db.Raw("select distinct(project_id) as id from ppm_pro_project_relation where is_delete = 2 and relation_type = 2 and relation_id in ?", input.Participants)
			} else if len(input.Participants) == 0 && len(deptIds) > 0 {
				sql = db.Raw("select distinct(project_id) as id from ppm_pro_project_relation where is_delete = 2 and relation_type = 25 and relation_id in ?", deptIds)
			}
			condParam[consts.TcId+" "] = db.In(sql)
		}
	}
	if len(input.Followers) > 0 {
		idStr := strings.Replace(strings.Trim(fmt.Sprint(input.Followers), "[]"), " ", ",", -1)
		sql := "select distinct(project_id) as id from ppm_pro_project_relation where is_delete = 2 and relation_type = 3 and relation_id in (" + idStr + ")"
		condParam[consts.TcId+"  "] = db.In(db.Raw(sql))
	}
	if input.IsMember != nil && *input.IsMember == 1 {
		//我是项目成员
		//查找我参与的部门
		if !adminFlag.IsSysAdmin && !adminFlag.IsSubAdmin {
			deptResp := orgfacade.GetUserDeptIds(&orgvo.GetUserDeptIdsReq{
				OrgId:  orgId,
				UserId: currentUserId,
			})
			if deptResp.Failure() {
				log.Error(deptResp.Error())
				return nil, deptResp.Error()
			}
			allDeptIds := []int64{0}
			allDeptIds = append(allDeptIds, deptResp.Data.DeptIds...)
			sql := db.Raw("select distinct(project_id) as id from ppm_pro_project_relation where is_delete = 2 and relation_type in (1, 2, 3) and relation_id = ?", currentUserId)
			if len(allDeptIds) > 0 {
				sql = db.Raw("select distinct(project_id) as id from ppm_pro_project_relation where is_delete = 2 and ((relation_type in (1, 2, 3) and relation_id = ?) or (relation_type=25 and relation_id in ?))", currentUserId, allDeptIds)
			}
			condParam[consts.TcId+"   "] = db.In(sql)
		}
	}

	return condParam, nil
}

func dealInputRelateType(condParam db.Cond, input *vo.ProjectsReq, currentUserId int64, orgId int64, adminFlag uservo.GetUserAuthorityData, allCollaborateProjectIds []int64) errs.SystemErrorInfo {
	if input.RelateType != nil {
		//查找我参与的部门
		//deptResp := orgfacade.GetUserDeptIds(orgvo.GetUserDeptIdsReq{
		//	OrgId:  orgId,
		//	UserId: currentUserId,
		//})
		//if deptResp.Failure() {
		//	log.Error(deptResp.Error())
		//	return deptResp.Error()
		//}
		//userAuthorityResp := userfacade.GetUserAuthority(orgId, currentUserId)
		//if userAuthorityResp.Failure() {
		//	log.Error(userAuthorityResp.Error())
		//	return userAuthorityResp.Error()
		//}
		userDeptIdsWithParentIdResp := orgfacade.GetUserDeptIdsWithParentId(orgvo.GetUserDeptIdsWithParentIdReq{
			OrgId:  orgId,
			UserId: currentUserId,
		})
		if userDeptIdsWithParentIdResp.Failure() {
			log.Errorf("[dealInputRelateType] orgfacade.GetUserDeptIdsWithParentId err:%v, orgId:%v, userId:%v",
				userDeptIdsWithParentIdResp.Error(), orgId, currentUserId)
			return userDeptIdsWithParentIdResp.Error()
		}
		allDeptIds := []int64{0}
		allDeptIds = append(allDeptIds, userDeptIdsWithParentIdResp.Data.DeptIds...)
		switch *input.RelateType {
		case 0, 4:
			// 如果是普通管理员
			//if adminFlag.IsSubAdmin {
			//	if len(adminFlag.ManageApps) > 0 && adminFlag.ManageApps[0] != -1 {
			//		// 转换为projectIds
			//		proIds, err := domain.GetProjectIdsByAppIds(orgId, adminFlag.ManageApps)
			//		if err != nil {
			//			log.Errorf("[dealInputRelateType] err:%v, orgId:%v", err, orgId)
			//			return err
			//		}}
			//		condParam[consts.TcId+" In"] = proIds
			//	}
			//}
			if adminFlag.IsSubAdmin {
				// 如果是普通管理员，查可以管理哪些应用
				if len(adminFlag.ManageApps) > 0 && adminFlag.ManageApps[0] != -1 {
					proIds, err := domain.GetProjectIdsByAppIds(orgId, adminFlag.ManageApps)
					if err != nil {
						log.Errorf("[dealInputRelateType] err:%v, orgId:%v", err, orgId)
						return err
					}
					allCollaborateProjectIds = append(allCollaborateProjectIds, proIds...)
				}
			}

			if !adminFlag.IsSysAdmin && !(adminFlag.IsSubAdmin && len(adminFlag.ManageApps) > 0 && adminFlag.ManageApps[0] == -1) {
				sql := db.Raw("select distinct(project_id) as id from ppm_pro_project_relation where org_id = ? and is_delete = 2 and ((relation_type in (1,2) and relation_id = ?) or (relation_type=25 and relation_id in ?))", orgId, currentUserId, allDeptIds)
				if len(allCollaborateProjectIds) > 0 {
					sql = db.Raw("select distinct(project_id) as id from ppm_pro_project_relation where (org_id = ? and is_delete = 2 and ((relation_type in (1,2) and relation_id = ?) or (relation_type=25 and relation_id in ?)) or project_id in ?)", orgId, currentUserId, allDeptIds, slice.SliceUniqueInt64(allCollaborateProjectIds))
				}
				condParam[consts.TcId+" In"] = db.In(sql)
			}

			//case 1:
			//	//我发起的
			//	condParam[consts.TcCreator] = currentUserId
			//case 2:
			//	//我负责的
			//	condParam[consts.TcOwner] = db.Eq(currentUserId)
			//case 3:
			//	//我参与的
			//	if !adminFlag.IsSysAdmin && !adminFlag.IsSubAdmin {
			//		sql := db.Raw("select distinct(project_id) as id from ppm_pro_project_relation where org_id = ? and is_delete = 2 and relation_type = ? and relation_id = ?", orgId, consts.IssueRelationTypeParticipant, currentUserId)
			//		if len(allDeptIds) > 0 {
			//			sql = db.Raw("select distinct(project_id) as id from ppm_pro_project_relation where org_id = ? and is_delete = 2 and ((relation_type = ? and relation_id = ?) or (relation_type=25 and relation_id in ?))", orgId, consts.IssueRelationTypeParticipant, currentUserId, allDeptIds)
			//		}
			//		condParam[consts.TcId+" In"] = db.In(sql)
			//	}
			//case 4:
			//	//我参与的和我负责的 + 和我协作的
			//	if !adminFlag.IsSysAdmin && !adminFlag.IsSubAdmin {
			//		sql := db.Raw("select distinct(project_id) as id from ppm_pro_project_relation where org_id = ? and is_delete = 2 and ((relation_type in (1,2) and relation_id = ?) or (relation_type=25 and relation_id in ?))", orgId, currentUserId, allDeptIds)
			//		//condParam[consts.TcId+" In"] = db.Or(db.In(sql))
			//		if len(allCollaborateProjectIds) > 0 {
			//			sql = db.Raw("select distinct(project_id) as id from ppm_pro_project_relation where (org_id = ? and is_delete = 2 and ((relation_type in (1,2) and relation_id = ?) or (relation_type=25 and relation_id in ?) or project_id in ?))", orgId, currentUserId, allDeptIds, allCollaborateProjectIds)
			//		}
			//		condParam[consts.TcId+" In"] = db.In(sql)
			//	}
			//case 5:
			//	//我关注的
			//	//获取我收藏的项目
			//	startProject, err := domain.GetProjectRelationByCond(db.Cond{
			//		consts.TcOrgId:        orgId,
			//		consts.TcRelationType: consts.IssueRelationTypeStar,
			//		consts.TcRelationId:   currentUserId,
			//		consts.TcIsDelete:     consts.AppIsNoDelete,
			//	})
			//	if err != nil {
			//		log.Error(err)
			//		return errs.MysqlOperateError
			//	}
			//	need := []int64{}
			//	for _, relationBo := range *startProject {
			//		need = append(need, relationBo.ProjectId)
			//	}
			//	condParam[consts.TcId+" In"] = db.In(need)
		}
	}
	return nil
}

func addCreatorInfo(resultList *[]*vo.Project, creatorInfo map[int64]bo.UserIDInfoBo) {
	for k, v := range *resultList {
		if _, ok := creatorInfo[v.Creator]; ok {
			creatorInfoModel := &vo.UserIDInfo{}
			copyer.Copy(creatorInfo[v.Creator], creatorInfoModel)
			(*resultList)[k].CreatorInfo = creatorInfoModel
		}
	}
}
func dealResultList(resultList *[]*vo.Project, ownerInfo map[int64][]bo.UserIDInfoBo, participantInfo map[int64][]bo.UserIDInfoBo,
	followerInfo map[int64][]bo.UserIDInfoBo, resourceByPath map[int64]bo.ResourceBo, issueStat map[int64]bo.IssueStatistic,
	projectTypeLocalMap, projectDetailById maps.LocalMap,
	iterationStat map[int64]bo.IssueStatistic, appIconMap map[int64]string) {
	// 先拿出所有的 projectId
	orgId := int64(0)
	projectIds := make([]int64, 0)
	for _, v := range *resultList {
		projectIds = append(projectIds, v.ID)
		if orgId < 1 {
			orgId = v.OrgID
		}
	}

	//statusBosMap, err := domain.GetProjectStatusBatch(orgId, projectIds)
	//if err != nil {
	//	log.Error(err)
	//	return
	//}
	proStatusBoList, _ := domain.GetProjectStatus(orgId, 0)

	for k, v := range *resultList {
		(*resultList)[k].OwnerInfo = &vo.UserIDInfo{}
		(*resultList)[k].OwnersInfo = []*vo.UserIDInfo{}
		if _, ok := ownerInfo[v.ID]; ok {
			ownerInfoModel := &[]*vo.UserIDInfo{}
			copyer.Copy(ownerInfo[v.ID], ownerInfoModel)
			if *ownerInfoModel != nil && len(*ownerInfoModel) > 0 {
				(*resultList)[k].OwnerInfo = (*ownerInfoModel)[0]
				(*resultList)[k].OwnersInfo = *ownerInfoModel
			}
		}
		if _, ok := participantInfo[v.ID]; ok {
			participantInfoModel := &[]*vo.UserIDInfo{}
			copyer.Copy(participantInfo[v.ID], participantInfoModel)
			(*resultList)[k].MemberInfo = *participantInfoModel
		} else {
			(*resultList)[k].MemberInfo = []*vo.UserIDInfo{}
		}
		if _, ok := followerInfo[v.ID]; ok {
			followerInfoModel := &[]*vo.UserIDInfo{}
			copyer.Copy(followerInfo[v.ID], followerInfoModel)
			(*resultList)[k].FollowerInfo = *followerInfoModel
		} else {
			(*resultList)[k].FollowerInfo = []*vo.UserIDInfo{}
		}
		if _, ok := resourceByPath[v.ResourceID]; ok {
			resource := resourceByPath[v.ResourceID]
			coverUrl := util.JointUrl(resource.Host, resource.Path)

			thumbnailUrl := util.GetCompressedPath(coverUrl, resource.Type)
			(*resultList)[k].ResourcePath = thumbnailUrl
			(*resultList)[k].ResourceCompressedPath = thumbnailUrl
		}
		if _, ok := issueStat[v.ID]; ok {
			(*resultList)[k].AllIssues = issueStat[v.ID].All
			(*resultList)[k].FinishIssues = issueStat[v.ID].Finish
			(*resultList)[k].OverdueIssues = issueStat[v.ID].Overdue
			(*resultList)[k].RelateUnfinish = issueStat[v.ID].RelateUnfinish
		}
		if _, ok := iterationStat[v.ID]; ok {
			temp := vo.IterationStatSimple{
				ID:            iterationStat[v.ID].IterationId,
				Name:          iterationStat[v.ID].IterationName,
				AllIssues:     iterationStat[v.ID].All,
				OverdueIssues: iterationStat[v.ID].Overdue,
				FinishIssues:  iterationStat[v.ID].Finish,
			}
			(*resultList)[k].IterationStat = &temp
		} else {
			if _, ok := issueStat[v.ID]; ok {
				(*resultList)[k].IterationStat = &vo.IterationStatSimple{
					ID:            0,
					Name:          "",
					AllIssues:     issueStat[v.ID].All,
					OverdueIssues: issueStat[v.ID].Overdue,
					FinishIssues:  issueStat[v.ID].Finish,
				}
			} else {
				(*resultList)[k].IterationStat = &vo.IterationStatSimple{}
			}
		}
		if times.GetUnixTime(*v.PlanStartTime) <= 0 {
			(*resultList)[k].PlanStartTime = nil
		}
		if times.GetUnixTime(*v.PlanEndTime) <= 0 {
			(*resultList)[k].PlanEndTime = nil
		}
		if projectTypeInterface, ok := projectTypeLocalMap[v.ProjectTypeID]; ok {
			projectType := projectTypeInterface.(bo.ProjectTypeBo)
			(*resultList)[k].ProjectTypeName = projectType.Name
			(*resultList)[k].ProjectTypeLangCode = projectType.LangCode
		}
		if projectDetailInterface, ok := projectDetailById[v.ID]; ok {
			projectDetail := projectDetailInterface.(bo.ProjectDetailBo)
			(*resultList)[k].IsSyncOutCalendar = projectDetail.IsSyncOutCalendar
		}

		isOtherLang := lang2.IsEnglish()

		//获取项目状态，先批量查询，再根据 projectId 获取
		var allProjectStatus = &[]bo.CacheProcessStatusBo{}
		allProjectStatus = &proStatusBoList

		lang := lang2.GetLang()
		otherLanguageMap := make(map[string]string, 0)
		if tmpMap, ok1 := consts.LANG_ISSUE_STAT_DESC_MAP[lang]; ok1 {
			otherLanguageMap = tmpMap
		}
		statusInfo := []*vo.HomeIssueStatusInfo{}
		//项目状态去除未开始
		statusNeedUpdate := false
		var processingStatus int64
		var statusIds []int64
		for _, val := range *allProjectStatus {
			if val.StatusId == v.Status {
				(*resultList)[k].StatusType = val.StatusType
			}
			statusIds = append(statusIds, val.StatusId)
			if val.StatusType == consts.StatusTypeNotStart {
				if val.StatusId == v.Status {
					statusNeedUpdate = true
				}
				continue
			}
			if val.StatusType == consts.StatusTypeRunning {
				processingStatus = val.StatusId
			}
			if isOtherLang {
				if tmpVal, ok2 := otherLanguageMap[val.Name]; ok2 {
					val.Name = tmpVal
				}
			}
			displayName := val.DisplayName
			info := vo.HomeIssueStatusInfo{
				Type:        val.StatusType,
				ID:          val.StatusId,
				Name:        val.Name,
				DisplayName: &displayName,
				BgStyle:     val.BgStyle,
				FontStyle:   val.FontStyle,
			}
			statusInfo = append(statusInfo, &info)
		}
		(*resultList)[k].AllStatus = statusInfo
		//如果项目是未开始则改为进行中
		if ok, _ := slice.Contain(statusIds, v.Status); !ok {
			statusNeedUpdate = true
		}
		if statusNeedUpdate && processingStatus != 0 {
			(*resultList)[k].Status = processingStatus
			(*resultList)[k].StatusType = consts.StatusTypeRunning
		}

		appId, appIdErr := strconv.ParseInt(v.AppID, 10, 64)
		if appIdErr == nil && appId != 0 {
			if icon, ok := appIconMap[appId]; ok {
				(*resultList)[k].Icon = icon
			}
		}
	}
}

func CreateProject(reqVo projectvo.CreateProjectReqVo) (*vo.Project, errs.SystemErrorInfo) {
	isError, err := checkAuth(&reqVo.UserId, &reqVo.OrgId)
	if isError {
		return nil, errs.BuildSystemErrorInfo(errs.TokenAuthError, err)
	}
	return CreateProjectWithoutAuth(reqVo)
}

func CreateProjectWithoutAuth(reqVo projectvo.CreateProjectReqVo) (*vo.Project, errs.SystemErrorInfo) {
	orgId := reqVo.OrgId
	currentUserId := reqVo.UserId
	input := reqVo.Input
	sourceChannel := reqVo.SourceChannel

	authFunctionErr := domain.AuthPayProjectNum(orgId, consts.FunctionProjectCreate)
	if authFunctionErr != nil {
		log.Error(authFunctionErr)
		return nil, authFunctionErr
	}

	id, err := idfacade.ApplyPrimaryIdRelaxed(consts.TableProject)
	if err != nil {
		return nil, errs.BuildSystemErrorInfo(errs.ApplyIdError, err)
	}

	// 融合极星-创建项目，需要在无码系统中创建一个应用与之对应
	// 查询组织对应的汇总表的 appId
	org := orgfacade.OrganizationInfo(orgvo.OrganizationInfoReqVo{
		OrgId:  orgId,
		UserId: currentUserId,
	})
	remarkStr := org.OrganizationInfo.Remark
	remarkObj := orgvo.OrgRemarkConfigType{}
	if err := json.FromJson(remarkStr, &remarkObj); err != nil {
		log.Errorf("[CreateProjectWithoutAuth] orgId: %d, err: %v", orgId, authFunctionErr)
		return nil, authFunctionErr
	}

	//日历兼容老的
	if input.SyncCalendarStatusList == nil {
		if input.IsSyncOutCalendar != nil && *input.IsSyncOutCalendar == 1 {
			ownerSync := consts.IsSyncOutCalendarForOwner
			followerSync := consts.IsSyncOutCalendarForFollower
			subSync := consts.IsSyncOutCalendarForSubCalendar
			input.SyncCalendarStatusList = []*int{&ownerSync, &followerSync, &subSync}
		} else {
			input.SyncCalendarStatusList = []*int{}
		}
	}

	blankString := consts.BlankString
	var zero int64 = 0

	checkErr := assignmentInput(zero, &input, blankString)
	if checkErr != nil {
		log.Error(checkErr)
		return nil, checkErr
	}

	//兼容老的项目负责人
	allOwnerIds := []int64{}
	if input.Owner != int64(0) {
		allOwnerIds = append(allOwnerIds, input.Owner)
	}
	if input.OwnerIds != nil && len(input.OwnerIds) > 0 {
		allOwnerIds = append(allOwnerIds, input.OwnerIds...)
	}
	if len(allOwnerIds) == 0 {
		return nil, errs.NeedProjectOwner
	}
	entity := &bo.ProjectBo{
		Id:            id,
		OrgId:         orgId,
		Code:          *input.Code,
		Name:          input.Name,
		PreCode:       *input.PreCode,
		Owner:         allOwnerIds[0], //随便取一个，后续用不着这个
		OwnerIds:      allOwnerIds,
		PriorityId:    *input.PriorityID,
		PublicStatus:  input.PublicStatus,
		ProjectTypeId: *input.ProjectTypeID,
		IsFiling:      2,
		Remark:        *input.Remark,
		Creator:       currentUserId,
		CreateTime:    types.NowTime(),
		Updator:       currentUserId,
		UpdateTime:    types.NowTime(),
		Version:       1,
		IsDelete:      consts.AppIsNoDelete,
	}
	initStatusErr := initProjectTypeAndProcessStatus(orgId, entity)
	if initStatusErr != nil {
		return nil, initStatusErr
	}

	isRepeat, err := checkRepeat(err, input, orgId, entity)
	if isRepeat {
		return nil, errs.BuildSystemErrorInfo(errs.TokenAuthError, err)
	}

	//创建project
	appId, newProEntity, addedMemberIds, err := domain.CreateProject(*entity, orgId, currentUserId, input, remarkObj)
	if err != nil {
		log.Errorf("[CreateProjectWithoutAuth] err: %v", err)
		return nil, errs.BuildSystemErrorInfo(errs.ProjectDomainError, err)
	}
	(*entity).AppId = appId

	result := &vo.Project{}
	err1 := copyer.Copy(newProEntity, result)
	if err1 != nil {
		log.Errorf("[CreateProjectWithoutAuth] copyer.Copy err: %v", err)
		return nil, errs.BuildSystemErrorInfo(errs.ObjectCopyError, err)
	}
	//创建人默认设置成项目的管理员
	asyn.Execute(func() {

	})

	//创建默认四个视图
	viewErr := domain.CreateProjectDefaultView(orgId, id, appId, *input.ProjectTypeID, nil, false)
	if viewErr != nil {
		log.Errorf("[CreateProjectWithoutAuth] err: %v", viewErr)
		return nil, viewErr
	}

	//创建日历
	go func() {
		defer func() {
			if r := recover(); r != nil {
				log.Error(errs.BuildSystemErrorInfoWithPanicRecover(r, stack.GetStack()))
			}
		}()
		// 将日历同步的状态值转换一下
		syncFlag := domain.TransferSyncOutCalendarStatusIntoOne(input.SyncCalendarStatusList)
		domain.CreateCalendar(&syncFlag, orgId, id, currentUserId, addedMemberIds)
	}()

	// 项目数据同步到无码
	//if err := CreateProjectInLessCode(orgId, remarkObj.ProjectFormAppId, entity, input); err != nil {
	//	log.Errorf("[CreateProjectWithoutAuth] err: %v, orgId： %d", err, orgId)
	//	return nil, err
	//}

	asyn.Execute(func() {
		ext := bo.TrendExtensionBo{}
		ext.ObjName = input.Name
		projectTrendsBo := bo.ProjectTrendsBo{
			PushType:            consts.PushTypeCreateProject,
			OrgId:               orgId,
			ProjectId:           id,
			OperatorId:          currentUserId,
			BeforeChangeMembers: []int64{},
			AfterChangeMembers:  addedMemberIds,
			NewValue:            json.ToJsonIgnoreError(entity),
			Ext:                 ext,
			SourceChannel:       sourceChannel,
		}
		domain.PushProjectTrends(projectTrendsBo)
		// 更新一下缓存
		domain.ClearSomeCache(projectsvcConsts.CacheBaseProjectInfo, map[string]interface{}{
			consts.CacheKeyOrgIdConstName:     orgId,
			consts.CacheKeyProjectIdConstName: entity.Id,
		})
	})
	asyn.Execute(func() {
		PushAddProjectNotice(orgId, result.ID, currentUserId)
	})
	asyn.Execute(func() {
		//新建项目不创建群聊
		//chatId := ""
		//if input.IsCreateFsChat == nil || *input.IsCreateFsChat != 2 {
		//	chatId, err = AddChatFromPlatform(orgId, currentUserId, addedMemberIds, result.ID, input.Name, input.Remark, allOwnerIds, input.MemberForDepartmentID, sourceChannel)
		//	if err != nil {
		//		log.Errorf("[CreateProjectWithoutAuth] AddChat err:%v", err)
		//	}
		//}
		//orgResp := orgfacade.GetBaseOrgInfo(orgvo.GetBaseOrgInfoReqVo{OrgId: orgId})
		//if orgResp.Failure() {
		//	log.Errorf("[CreateProjectWithoutAuth] GetBaseOrgInfo failed: %v, orgId: %d", orgResp.Error(), orgId)
		//	return
		//}

		//sourceChannel = orgResp.BaseOrgInfo.SourceChannel
		//if ok, _ := slice.Contain([]string{sdk_const.SourceChannelFeishu, sdk_const.SourceChannelWeixin,
		//	sdk_const.SourceChannelDingTalk}, sourceChannel); ok {
		//	// 创建群聊推送
		//	e := &commonvo.AppEvent{
		//		OrgId:     orgId,
		//		AppId:     appId,
		//		ProjectId: id,
		//		UserId:    input.Owner,
		//		App:       nil,
		//		Project:   nil,
		//		Chat:      chatId,
		//	}
		//
		//	openTraceId, _ := threadlocal.Mgr.GetValue(consts.JaegerContextTraceKey)
		//	openTraceIdStr := cast.ToString(openTraceId)
		//
		//	report.ReportAppEvent(msgPb.EventType_AppChatCreated, openTraceIdStr, e)
		//}

	})

	// 任务状态改造后，创建敏捷项目无需再更新项目对应的应用的 form 头。因为使用的是 tables 模式。
	// deleted code

	result.AppID = fmt.Sprintf("%d", appId)
	return result, nil
}

// CreateProjectInLessCode 创建项目时，将新项目同步到无码
func CreateProjectInLessCode(orgId int64, proFormAppId int64, newProjectObj *bo.ProjectBo, input vo.CreateProjectReq) errs.SystemErrorInfo {
	formData := map[string]interface{}{}
	// 查询详情，组装要新增的数据
	curDetail, err := domain.GetProjectDetailByProjectIdBo(newProjectObj.Id, orgId)
	if err != nil {
		log.Error(err)
		return err
	}
	copyErr := copyer.Copy(newProjectObj, &formData)
	if copyErr != nil {
		log.Errorf("[CreateProjectInLessCode] error: %s", copyErr)
		return errs.ObjectCopyError
	}
	formData["orgId"] = orgId
	formData[consts.ProBasicFieldName] = newProjectObj.Name
	formData[consts.ProBasicFieldProId] = newProjectObj.Id
	formData[consts.ProBasicFieldAppId] = newProjectObj.AppId
	formData[consts.ProBasicFieldCode] = newProjectObj.Code
	formData[consts.ProBasicFieldPreCode] = newProjectObj.PreCode
	formData[consts.ProBasicFieldProTypeId] = newProjectObj.ProjectTypeId
	formData[consts.ProBasicFieldPriorityId] = newProjectObj.PriorityId
	formData[consts.ProBasicFieldPlanStartTime] = newProjectObj.PlanStartTime
	formData[consts.ProBasicFieldPlanEndTime] = newProjectObj.PlanEndTime
	formData[consts.ProBasicFieldPublicStatus] = newProjectObj.PublicStatus
	formData[consts.ProBasicFieldTemplateFlag] = newProjectObj.TemplateFlag
	formData[consts.ProBasicFieldResource] = newProjectObj.ResourceId // 取出 resource 存入地址
	formData[consts.ProBasicFieldIsFiling] = newProjectObj.IsFiling
	formData[consts.ProBasicFieldRemark] = newProjectObj.Remark

	formData[consts.ProBasicFieldOwnerIds] = []string{fmt.Sprintf("U_%d", input.Owner)}
	proParticipantIdStrArr := make([]string, 0)
	for _, tmpUserId := range input.MemberIds {
		proParticipantIdStrArr = append(proParticipantIdStrArr, fmt.Sprintf("U_%d", tmpUserId))
	}
	if len(input.MemberForDepartmentID) > 0 {
		for _, deptId := range input.MemberForDepartmentID {
			proParticipantIdStrArr = append(proParticipantIdStrArr, fmt.Sprintf("D_%d", deptId))
		}
	}
	// 如果是成员全选，则用 `D_0` 表示
	if input.IsAllMember != nil && *input.IsAllMember {
		proParticipantIdStrArr = append(proParticipantIdStrArr, fmt.Sprintf("D_%d", 0))
	}
	formData[consts.ProBasicFieldParticipantIds] = proParticipantIdStrArr
	formData[consts.ProBasicFieldOutCalendar] = make([]string, 0) // 日历创建是异步的，因此这里暂时设为空切片
	formData[consts.ProBasicFieldOutCalendarSettings] = domain.TransferSyncOutCalendarStatusIntoOne(input.SyncCalendarStatusList)
	formData[consts.ProBasicFieldOutChat] = make([]string, 0) // 群聊创建是异步的 这里暂时设为空切片
	formData[consts.ProBasicFieldIsEnableWorkHours] = curDetail.IsEnableWorkHours

	formData[consts.ProBasicFieldCreateTime] = newProjectObj.CreateTime
	formData[consts.ProBasicFieldCreator] = newProjectObj.Creator
	formData[consts.ProBasicFieldUpdateTime] = newProjectObj.UpdateTime
	formData[consts.ProBasicFieldUpdator] = newProjectObj.Updator
	formData[consts.ProBasicFieldVersion] = newProjectObj.Version
	formData[consts.ProBasicFieldIsDelete] = newProjectObj.IsDelete
	delete(formData, "id")

	insertResp := formfacade.LessCreateIssue(formvo.LessCreateIssueReq{
		AppId:       proFormAppId,
		OrgId:       orgId,
		UserId:      newProjectObj.Creator,
		RedirectIds: []int64{proFormAppId}, // 这个参数是指什么
		Import:      true,
		Form:        []map[string]interface{}{formData},
	})
	if insertResp.Failure() {
		log.Error(insertResp.Error())
		return insertResp.Error()
	}

	return nil
}

func initProjectTypeAndProcessStatus(orgId int64, entity *bo.ProjectBo) errs.SystemErrorInfo {
	projectTypeId, status, err := domain.GetTypeAndStatus(entity.ProjectTypeId)
	if err != nil {
		return errs.BuildSystemErrorInfo(errs.ProjectDomainError, err)
	}
	entity.ProjectTypeId = projectTypeId
	entity.Status = status

	return nil
}

// 校验重复
func checkRepeat(err error, input vo.CreateProjectReq, orgId int64, entity *bo.ProjectBo) (isRepeatError bool, repeatErr errs.SystemErrorInfo) {
	//_, err = domain.JudgeRepeatProjectName(&input.Name, orgId, nil)
	//if err != nil {
	//	return true, errs.BuildSystemErrorInfo(errs.ProjectDomainError, err)
	//}
	if input.PreCode != nil && *input.PreCode != "" {
		_, err = domain.JudgeRepeatProjectPreCode(input.PreCode, orgId, nil)
		if err != nil {
			return true, errs.BuildSystemErrorInfo(errs.ProjectDomainError, err)
		}
	}

	if ok, _ := slice.Contain([]int{consts.PublicProject, consts.PrivateProject}, entity.PublicStatus); !ok {
		return true, errs.BuildSystemErrorInfo(errs.ProjectDomainError, errors.New("项目可见性选择有误"))
	}
	//entity.PlanStartTime input.PlanStartTime的地址
	if input.PlanStartTime != nil && input.PlanStartTime.IsNotNull() {

		(*entity).PlanStartTime = *input.PlanStartTime
	} else {

		PlanStartTime := types.Time(consts.BlankTimeObject)

		(*entity).PlanStartTime = PlanStartTime
	}

	//entity.planEndTime的指针变量等于 input.PlanEndTime的地址
	if input.PlanEndTime != nil && input.PlanEndTime.IsNotNull() {

		(*entity).PlanEndTime = *input.PlanEndTime
	} else {
		BlankTime := types.Time(consts.BlankTimeObject)

		(*entity).PlanEndTime = BlankTime
	}

	if time.Time(entity.PlanEndTime).After(consts.BlankTimeObject) && time.Time(entity.PlanStartTime).After(time.Time(entity.PlanEndTime)) {
		return true, errs.BuildSystemErrorInfo(errs.CreateProjectTimeError)
	}
	return false, nil
}

// 校验权限
func checkAuth(currentUserId *int64, orgId *int64) (isError bool, error error) {
	// 更换权限的 operation code，以适配鉴权。
	err := domain.AuthOrg(*orgId, *currentUserId, consts.RoleOperationPathOrgProject, consts.OperationOrgProjectCreate)
	if err != nil {
		log.Error(err)
		return true, errs.BuildSystemErrorInfo(errs.Unauthorized, err)
	}

	return false, nil
}

func assignmentInput(zero int64, input *vo.CreateProjectReq, blankString string) errs.SystemErrorInfo {
	if strings.Trim(input.Name, " ") == "" {
		return errs.ProjectNameEmpty
	}

	isNameRight := format.VerifyProjectNameFormat(input.Name)
	if !isNameRight {
		log.Error(errs.InvalidProjectNameError)
		return errs.InvalidProjectNameError
	}

	if input.Code == nil {
		input.Code = &blankString
	}
	if strs.Len(*input.Code) > 64 {
		return errs.BuildSystemErrorInfo(errs.ProjectCodeLenError)
	}

	if input.PreCode == nil {
		input.PreCode = &blankString
	} else {
		isPreCodeRight := format.VerifyProjectPreviousCodeFormat(*input.PreCode)
		if !isPreCodeRight {
			log.Error(errs.InvalidProjectPreCodeError)
			return errs.InvalidProjectPreCodeError
		}
	}

	if input.PriorityID == nil {
		input.PriorityID = &zero
	}
	if input.ProjectTypeID == nil {
		input.ProjectTypeID = &zero
	}
	if input.Remark == nil {
		input.Remark = &blankString
	}
	isRemarkRight := format.VerifyProjectRemarkFormat(*input.Remark)
	if !isRemarkRight {
		log.Error(errs.InvalidProjectRemarkError)
		return errs.InvalidProjectRemarkError
	}

	return nil
}

func UpdateProject(reqVo projectvo.UpdateProjectReqVo) (*vo.Project, errs.SystemErrorInfo) {
	currentUserId := reqVo.UserId
	orgId := reqVo.OrgId
	input := reqVo.Input

	log.Infof(consts.UserLoginSentence, currentUserId, orgId)

	//修改成员的权限另外获取
	if (util.FieldInUpdate(input.UpdateFields, "memberIds") && input.MemberIds != nil) || (util.FieldInUpdate(input.UpdateFields, "memberForDepartmentId") && input.MemberForDepartmentID != nil) {
		if input.IsAllMember != nil && *input.IsAllMember == true {
			resp := orgfacade.GetOrgUserIds(orgvo.GetOrgUserIdsReq{
				OrgId: orgId,
			})

			if resp.Failure() {
				log.Error(resp.Error())
				return nil, resp.Error()
			}
			tempMemberIds := []int64{}
			for _, info := range resp.Data {
				tempMemberIds = append(tempMemberIds, info)
			}
			input.MemberIds = tempMemberIds
		}
		input.MemberIds = slice.SliceUniqueInt64(input.MemberIds)
		//获取项目成员和成员部门
		projectRelation, relationErr := domain.GetProjectMembers(orgId, input.ID, []int64{consts.ProjectRelationTypeParticipant, consts.ProjectRelationTypeDepartmentParticipant})
		if relationErr != nil {
			log.Error(relationErr)
			return nil, relationErr
		}
		memberIds := set.New(set.ThreadSafe)
		departmentIds := set.New(set.ThreadSafe)
		for _, relationBo := range projectRelation {
			if relationBo.RelationType == consts.ProjectRelationTypeDepartmentParticipant {
				departmentIds.Add(relationBo.RelationId)
			} else {
				memberIds.Add(relationBo.RelationId)
			}
		}
		afterIds := set.New(set.ThreadSafe)
		for _, id := range input.MemberIds {
			afterIds.Add(id)
		}
		afterDepartmentIds := set.New(set.ThreadSafe)
		for _, i2 := range input.MemberForDepartmentID {
			afterDepartmentIds.Add(i2)
		}
		if set.Difference(memberIds, afterIds).Size() > 0 || set.Difference(departmentIds, afterDepartmentIds).Size() > 0 {
			err := domain.AuthProject(orgId, currentUserId, input.ID, consts.RoleOperationPathOrgProMember, consts.OperationProMemberUnbind)
			if err != nil {
				log.Error(err)
				return nil, errs.BuildSystemErrorInfo(errs.Unauthorized, err)
			}
		}

		if set.Difference(afterIds, memberIds).Size() > 0 || set.Difference(afterDepartmentIds, departmentIds).Size() > 0 {
			err := domain.AuthProject(orgId, currentUserId, input.ID, consts.RoleOperationPathOrgProMember, consts.OperationProMemberBind)
			if err != nil {
				log.Error(err)
				return nil, errs.BuildSystemErrorInfo(errs.Unauthorized, err)
			}
		}
	} else {
		err := domain.AuthProject(orgId, currentUserId, input.ID, consts.RoleOperationPathOrgProProConfig, consts.OperationProConfigModify)
		if err != nil {
			log.Error(err)
			return nil, errs.BuildSystemErrorInfo(errs.Unauthorized, err)
		}
	}

	return UpdateProjectWithoutAuth(reqVo)
}

func UpdateProjectWithoutAuth(reqVo projectvo.UpdateProjectReqVo) (*vo.Project, errs.SystemErrorInfo) {
	currentUserId := reqVo.UserId
	orgId := reqVo.OrgId
	input := reqVo.Input
	sourceChannel := reqVo.SourceChannel

	log.Infof(consts.UserLoginSentence, currentUserId, orgId)

	//日历兼容老的
	if input.SyncCalendarStatusList == nil {
		if input.IsSyncOutCalendar != nil && *input.IsSyncOutCalendar == 1 {
			ownerSync := consts.IsSyncOutCalendarForOwner
			followerSync := consts.IsSyncOutCalendarForFollower
			subSync := consts.IsSyncOutCalendarForSubCalendar
			input.SyncCalendarStatusList = []*int{&ownerSync, &followerSync, &subSync}
		} else {
			input.SyncCalendarStatusList = []*int{}
		}
	}

	originProjectInfo, err := domain.GetProjectInfoWithOwnerIds(input.ID, orgId)
	if err != nil {
		log.Error(err)
		return nil, errs.BuildSystemErrorInfo(errs.MysqlOperateError, errors.New("项目不存在或已被删除"))
	}
	if originProjectInfo.IsFiling == consts.ProjectIsFiling {
		return nil, errs.ProjectIsFilingYet
	}
	entity := &bo.UpdateProjectBo{}
	newValue := &bo.ProjectBo{}
	copyer.Copy(input, entity)
	copyer.Copy(originProjectInfo, newValue)

	old := &map[string]interface{}{}
	new := &map[string]interface{}{}
	changeList := []bo.TrendChangeListBo{}

	upd, err := domain.UpdateProjectCondAssembly(*entity, orgId, old, new, originProjectInfo, &changeList)
	if err != nil {
		log.Error(err)
		return nil, errs.BuildSystemErrorInfo(errs.ProjectDomainError, err)
	}

	if util.FieldInUpdate(input.UpdateFields, "resourcePath") && util.FieldInUpdate(input.UpdateFields, "resourceType") && input.ResourcePath != nil && input.ResourceType != nil {
		//资源列表
		resourcesRespVo := resourcefacade.GetResourceById(resourcevo.GetResourceByIdReqVo{GetResourceByIdReqBody: resourcevo.GetResourceByIdReqBody{ResourceIds: []int64{originProjectInfo.ResourceId}}})
		if resourcesRespVo.Failure() {
			log.Error(resourcesRespVo.Error())
			return nil, resourcesRespVo.Error()
		}
		oldResourcePath := ""
		if len(resourcesRespVo.ResourceBos) > 0 {
			oldResourcePath = resourcesRespVo.ResourceBos[0].Host + resourcesRespVo.ResourceBos[0].Path
		}
		changeList = append(changeList, bo.TrendChangeListBo{
			Field:     "resourcePath",
			FieldName: consts.ProjectResourcePath,
			NewValue:  *input.ResourcePath,
			OldValue:  oldResourcePath,
		})
	}

	err = domain.UpdateProject(orgId, currentUserId, upd, entity, originProjectInfo.AppId)
	if err != nil {
		log.Error(err)
		return nil, errs.BuildSystemErrorInfo(errs.ProjectDomainError, err)
	}

	err1 := domain.RefreshProjectAuthBo(orgId, input.ID)
	if err1 != nil {
		log.Error(err)
		return nil, errs.BuildSystemErrorInfo(errs.ProjectDomainError, err1)
	}

	pushType := consts.PushTypeUpdateProject

	// 查询旧的日历信息。
	oldProjectCalendarInfo, err := domain.GetProjectCalendarInfo(orgId, input.ID)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	// 更新存在于 pro detail 表中的数据的一些信息：同步日历、创建群聊、隐私状态等
	detailUpdLc := mysql.Upd{}
	updateErr := updateProjectWithDetail(input, orgId, currentUserId, detailUpdLc)
	if updateErr != nil {
		return nil, updateErr
	}

	//更新日历
	go func() {
		defer func() {
			if r := recover(); r != nil {
				log.Error(errs.BuildSystemErrorInfoWithPanicRecover(r, stack.GetStack()))
			}
		}()
		domain.UpdateCalendar(input, orgId, currentUserId, oldProjectCalendarInfo, originProjectInfo)
	}()

	// 项目更新同步到无码（todo 无码项目表）
	//lcUpd := slice2.CaseCamelCopy(upd)
	//detailUpdLc = slice2.CaseCamelCopy(detailUpdLc)
	//for k, v := range detailUpdLc {
	//	lcUpd[k] = v
	//}
	//proFormAppId, err := domain.GetProjectFormAppId(orgId)
	//if err != nil {
	//	log.Error(err)
	//	return nil, err
	//}
	//resp := formfacade.LessUpdateIssue(formvo.LessUpdateIssueReq{
	//	AppId:  proFormAppId,
	//	OrgId:  orgId,
	//	UserId: currentUserId,
	//	Form:   []map[string]interface{}{lcUpd},
	//})
	//if resp.Failure() {
	//	log.Error(resp.Error())
	//	return nil, resp.Error()
	//}

	//异步添加项目动态信息
	asyn.Execute(func() {
		ext := bo.TrendExtensionBo{}
		ext.ObjName = originProjectInfo.Name
		ext.ChangeList = changeList
		domain.PushProjectTrends(bo.ProjectTrendsBo{
			PushType:   pushType,
			OrgId:      orgId,
			ProjectId:  input.ID,
			OperatorId: currentUserId,
			OldValue:   json.ToJsonIgnoreError(old),
			NewValue:   json.ToJsonIgnoreError(new),
			Ext:        ext,

			SourceChannel: sourceChannel,
		})
	})
	asyn.Execute(func() {
		PushModifyProjectNotice(orgId, input.ID, currentUserId)
	})
	asyn.Execute(func() {

		chatId, err := domain.GetProjectMainChatId(orgId, input.ID)
		if err != nil {
			log.Error(err)
			return
		}
		if util.FieldInUpdate(input.UpdateFields, "isCreateFsChat") && input.IsCreateFsChat != nil {
			//更新配置表及关联关系
			updateSettingErr := UpdateProjectDetailChatSetting(orgId, input.ID, *input.IsCreateFsChat)
			if updateSettingErr != nil {
				log.Error(updateSettingErr)
				return
			}
			if *input.IsCreateFsChat == 2 {
				//关闭群聊
				if chatId != "" {
					//删除群聊
					DeleteFsChat(orgId, input.ID, chatId)
					chatId = ""
				}
			} else {
				//开启群聊
				if chatId != "" {
					//更新群聊名称
					if util.FieldInUpdate(input.UpdateFields, "name") && input.Name != nil {
						domain.UpdateChatTitle(orgId, input.ID, *input.Name)
					}
				} else {
					//查找项目成员 todo 后续增加部门人员
					participants, participantsErr := domain.GetProjectRelationByType(input.ID, []int64{consts.ProjectRelationTypeParticipant, consts.ProjectRelationTypeDepartmentParticipant})
					if participantsErr != nil {
						log.Error(participantsErr)
					}
					participantIds := make([]int64, 0)
					participantDeptIds := make([]int64, 0)
					for _, bo := range *participants {
						if bo.RelationType == consts.ProjectRelationTypeParticipant {
							participantIds = append(participantIds, bo.RelationId)
						} else {
							participantDeptIds = append(participantDeptIds, bo.RelationId)
						}
					}
					projectInfo, err := domain.GetProject(orgId, input.ID)
					if err != nil {
						log.Error(err)
						return
					}
					chatId, err = AddFsChat(orgId, currentUserId, participantIds, input.ID, projectInfo.Name, nil, projectInfo.OwnerIds, participantDeptIds)
					if err != nil {
						log.Errorf("[UpdateProjectWithoutAuth] AddFsChat err:%v", err)
					}
				}

				orgResp := orgfacade.GetBaseOrgInfo(orgvo.GetBaseOrgInfoReqVo{OrgId: orgId})
				if orgResp.Failure() {
					log.Errorf("[UpdateProjectWithoutAuth] GetBaseOrgInfo failed: %v, orgId: %d", orgResp.Error(), orgId)
					return
				}
				sourceChannel = orgResp.BaseOrgInfo.SourceChannel
				//只有开启群聊才发起推送
				if ok, _ := slice.Contain([]string{sdk_const.SourceChannelFeishu, sdk_const.SourceChannelWeixin,
					sdk_const.SourceChannelDingTalk}, sourceChannel); ok {
					e := &commonvo.AppEvent{
						OrgId:     orgId,
						AppId:     newValue.AppId,
						ProjectId: input.ID,
						UserId:    newValue.Updator,
						App:       nil,
						Project:   nil,
						Chat:      chatId,
					}

					openTraceId, _ := threadlocal.Mgr.GetValue(consts.JaegerContextTraceKey)
					openTraceIdStr := cast.ToString(openTraceId)

					report.ReportAppEvent(msgPb.EventType_AppChatCreated, openTraceIdStr, e)
				}
			}
			err := domain.ClearProjectMainChatCache(orgId, input.ID)
			if err != nil {
				log.Error(err)
				return
			}
		} else {
			//更新群聊名称
			if util.FieldInUpdate(input.UpdateFields, "name") && input.Name != nil {
				domain.UpdateChatTitle(orgId, input.ID, *input.Name)
			}
		}

	})
	asyn.Execute(func() {
		eventType := []string{}
		if util.FieldInUpdate(input.UpdateFields, "name") {
			eventType = append(eventType, consts.FsEventUpdateProjectName)
		}
		if util.FieldInUpdate(input.UpdateFields, "remark") {
			eventType = append(eventType, consts.FsEventUpdateProjectRemark)
		}
		//if util.FieldInUpdate(input.UpdateFields, "owner") {
		//	eventType = append(eventType, consts.FsEventUpdateProjectOwner)
		//}
		//if addedMembersSet.Size() != 0 {
		//	eventType = append(eventType, consts.FsEventAddProjectMember)
		//}
		//if deletedMembersSet.Size() != 0 {
		//	eventType = append(eventType, consts.FsEventDeleteProjectMember)
		//}

		// if len(eventType) > 0 {
		// 	domain.PushMessageToFeishuShortcut(bo.ShortcutPushBo{
		// 		TriggerType:         consts.FsTriggerUpdateProject,
		// 		EventType:           eventType,
		// 		OrgId:               orgId,
		// 		ProjectId:           input.ID,
		// 		ProjectObjectTypeId: 0,
		// 		IssueId:             0,
		// 		Operator:            currentUserId,
		// 	})
		// }
	})

	if input.Name != nil && *input.Name != "" {
		domain.UpdateDingTopCard(orgId, input.ID)
	}

	return &vo.Project{ID: input.ID}, nil
}

// 更新日历、开启项目群聊设置等
func updateProjectWithDetail(input vo.UpdateProjectReq, orgId int64, currentUserId int64, detailUpdLc mysql.Upd) errs.SystemErrorInfo {
	if !util.FieldInUpdate(input.UpdateFields, "syncCalendarStatusList") &&
		!util.FieldInUpdate(input.UpdateFields, "isSyncOutCalendar") &&
		!util.FieldInUpdate(input.UpdateFields, "isCreateFsChat") {
		return nil
	}
	projectDetail, err := domain.GetProjectDetailByProjectIdBo(input.ID, orgId)
	if err != nil {
		log.Error(err)
		return err
	}
	updateData := &bo.ProjectDetailBo{
		Id:      projectDetail.Id,
		Updator: currentUserId,
	}
	if util.FieldInUpdate(input.UpdateFields, "syncCalendarStatusList") || util.FieldInUpdate(input.UpdateFields, "isSyncOutCalendar") {
		if input.SyncCalendarStatusList != nil {
			syncFlag := domain.TransferSyncOutCalendarStatusIntoOne(input.SyncCalendarStatusList)
			if projectDetail.IsSyncOutCalendar != syncFlag {
				updateData.IsSyncOutCalendar = syncFlag
				detailUpdLc[consts.ProBasicFieldOutCalendarSettings] = syncFlag
			}
		}
	}

	updateErr := domain.UpdateProjectDetail(updateData)
	if updateErr != nil {
		log.Error(updateErr)
		return updateErr
	}
	cacheErr := domain.DeleteProjectCalendarInfo(orgId, input.ID)
	if cacheErr != nil {
		log.Error(cacheErr)
		return cacheErr
	}

	return nil
}

func UpdateProjectStatus(reqVo projectvo.UpdateProjectStatusReqVo) (*vo.Void, errs.SystemErrorInfo) {
	input := reqVo.Input
	orgId := reqVo.OrgId
	currentUserId := reqVo.UserId
	sourceChannel := reqVo.SourceChannel

	projectId := input.ProjectID

	projectBo, err1 := domain.GetProjectSimple(orgId, projectId)
	if err1 != nil {
		log.Error(err1)
		return nil, errs.BuildSystemErrorInfo(errs.ProjectNotExist)
	}
	if projectBo.IsFiling == consts.ProjectIsFiling {
		return nil, errs.ProjectIsFilingYet
	}

	err := domain.AuthProject(orgId, currentUserId, projectId, consts.RoleOperationPathOrgProProConfig, consts.OperationProConfigModifyStatus)
	if err != nil {
		log.Error(err)
		return nil, errs.BuildSystemErrorInfo(errs.Unauthorized, err)
	}

	err2 := domain.UpdateProjectStatus(*projectBo, input.NextStatusID)
	if err2 != nil {
		log.Error(err2)
		return nil, errs.BuildSystemErrorInfo(errs.ProjectDomainError, err2)
	}

	refreshErr := domain.RefreshProjectAuthBo(orgId, projectId)
	if refreshErr != nil {
		log.Error(refreshErr)
		return nil, errs.BuildSystemErrorInfo(errs.ProjectDomainError, refreshErr)
	}

	asyn.Execute(func() {
		operateObjProperty := consts.TrendsOperObjPropertyNameStatus
		oldValueMap := map[string]interface{}{
			operateObjProperty: projectBo.Status,
		}
		newValueMap := map[string]interface{}{
			operateObjProperty: input.NextStatusID,
		}

		//状态列表
		statusList := consts.ProjectStatusList

		change := bo.TrendChangeListBo{
			Field:     "status",
			FieldName: consts.Status,
		}
		for _, v := range statusList {
			if v.ID == projectBo.Status {
				change.OldValue = v.Name
			} else if v.ID == input.NextStatusID {
				change.NewValue = v.Name
			}
		}
		changeList := []bo.TrendChangeListBo{}
		changeList = append(changeList, change)
		ext := bo.TrendExtensionBo{
			ObjName:    projectBo.Name,
			ChangeList: changeList,
		}

		domain.PushProjectTrends(bo.ProjectTrendsBo{
			PushType:           consts.PushTypeUpdateProjectStatus,
			OrgId:              orgId,
			ProjectId:          projectId,
			OperatorId:         currentUserId,
			OperateObjProperty: operateObjProperty,
			OldValue:           json.ToJsonIgnoreError(oldValueMap),
			NewValue:           json.ToJsonIgnoreError(newValueMap),
			Ext:                ext,

			SourceChannel: sourceChannel,
		})
	})

	asyn.Execute(func() {
		PushModifyProjectNotice(orgId, projectId, currentUserId)
	})
	return &vo.Void{
		ID: projectId,
	}, nil
}

func ProjectInfo(orgId int64, userId int64, input vo.ProjectInfoReq) (*vo.ProjectInfo, errs.SystemErrorInfo) {
	projectId := input.ProjectID
	projectBo, err1 := domain.GetProjectSimple(orgId, projectId)
	if err1 != nil {
		log.Error(err1)
		return nil, errs.BuildSystemErrorInfo(errs.ProjectNotExist)
	}

	result := &vo.ProjectInfo{}
	err2 := copyer.Copy(projectBo, result)
	if err2 != nil {
		log.Error(err2)
		return nil, errs.BuildSystemErrorInfo(errs.ObjectCopyError)
	}

	resourcesRespVo := resourcefacade.GetResourceById(resourcevo.GetResourceByIdReqVo{GetResourceByIdReqBody: resourcevo.GetResourceByIdReqBody{ResourceIds: []int64{projectBo.ResourceId}}})
	if resourcesRespVo.Failure() {
		log.Error(resourcesRespVo.Message)
		return nil, resourcesRespVo.Error()
	}
	resourceEntities := resourcesRespVo.ResourceBos

	for _, v := range resourceEntities {
		result.ResourcePath = v.Host + v.Path
	}

	//项目相关人员
	ownerInfo, participantInfo, followerInfo, creatorInfo, err := domain.GetProjectMemberInfo([]int64{projectId}, orgId, []int64{projectBo.Creator})
	if err != nil {
		return nil, errs.BuildSystemErrorInfo(errs.ProjectDomainError, err)
	}
	result.OwnerInfo = &vo.UserIDInfo{}
	result.OwnersInfo = []*vo.UserIDInfo{}
	if _, ok := ownerInfo[projectId]; ok {
		ownerInfoModel := &[]*vo.UserIDInfo{}
		copyer.Copy(ownerInfo[projectId], ownerInfoModel)
		if ownerInfoModel != nil && len(*ownerInfoModel) > 0 {
			result.OwnerInfo = (*ownerInfoModel)[0]
			result.OwnersInfo = *ownerInfoModel
		}
	}
	if _, ok := participantInfo[projectId]; ok {
		participantInfoModel := &[]*vo.UserIDInfo{}
		copyer.Copy(participantInfo[projectId], participantInfoModel)
		result.MemberInfo = *participantInfoModel
	}
	if _, ok := followerInfo[projectId]; ok {
		followerInfoModel := &[]*vo.UserIDInfo{}
		copyer.Copy(followerInfo[projectId], followerInfoModel)
		result.FollowerInfo = *followerInfoModel
	}
	if _, ok := creatorInfo[projectBo.Creator]; ok {
		creatorInfoModel := &vo.UserIDInfo{}
		copyer.Copy(creatorInfo[projectBo.Creator], creatorInfoModel)
		result.CreatorInfo = creatorInfoModel
	}

	//获取项目成员部门
	departments, departmentsErr := domain.GetProjectMemberDepartmentsInfo(orgId, projectId)
	if departmentsErr != nil {
		log.Error(departmentsErr)
		return nil, departmentsErr
	}
	departmentVos := &[]*vo.DepartmentSimpleInfo{}
	copyer.Copy(departments, departmentVos)
	result.MemberDepartmentInfo = *departmentVos

	//获取项目状态
	allProjectStatus, err := domain.GetProjectStatus(orgId, projectId)
	if err != nil {
		return nil, errs.BuildSystemErrorInfo(errs.CacheProxyError, err)
	}
	// 检查用户使用的语言包
	lang := lang2.GetLang()
	isOtherLang := lang2.IsEnglish()
	statusInfo := []*vo.HomeIssueStatusInfo{}

	//项目状态去除未开始
	statusNeedUpdate := false
	var processingStatus int64
	var statusIds []int64
	for i, v := range allProjectStatus {
		statusIds = append(statusIds, v.StatusId)
		if v.StatusType == consts.StatusTypeNotStart {
			if v.StatusId == result.Status {
				statusNeedUpdate = true
			}
			continue
		}
		if v.StatusType == consts.StatusTypeRunning {
			processingStatus = v.StatusId
		}
		if isOtherLang {
			if tmpMap, ok1 := consts.LANG_PROJECT_STATUS_MAP[lang]; ok1 {
				if tmpVal, ok2 := tmpMap[v.Name]; ok2 {
					v.Name = tmpVal
				}
			}
		}
		info := vo.HomeIssueStatusInfo{
			Type:        v.StatusType,
			ID:          v.StatusId,
			Name:        v.Name,
			DisplayName: &allProjectStatus[i].Name,
			BgStyle:     v.BgStyle,
			FontStyle:   v.FontStyle,
		}
		statusInfo = append(statusInfo, &info)
	}
	if ok, _ := slice.Contain(statusIds, result.Status); !ok {
		statusNeedUpdate = true
	}
	result.AllStatus = statusInfo
	if processingStatus != 0 && statusNeedUpdate {
		result.Status = processingStatus
	}
	projectDetail, err := domain.GetProjectDetailByProjectIdBo(projectId, orgId)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	// 将状态值转化为对应的单个状态值集合。
	result.SyncCalendarStatusList = domain.TransferSyncOutCalendarStatus(projectDetail.IsSyncOutCalendar)
	result.IsSyncOutCalendar = projectDetail.IsSyncOutCalendar

	//fsChatOpenFlag, _ := domain.CheckProFsChatSetIsOpen(orgId, projectDetail.ProjectId)
	// 配置中可能没有值，此时查找群聊关联是否存在
	//if fsChatOpenFlag == 0 {
	//	chatId, chatIdErr := domain.GetProjectMainChatId(orgId, projectId)
	//	if chatIdErr != nil {
	//		log.Error(chatIdErr)
	//		return nil, chatIdErr
	//	}
	//	if chatId != "" {
	//		result.IsCreateFsChat = 1
	//	} else {
	//		result.IsCreateFsChat = 2
	//	}
	//} else {
	//	result.IsCreateFsChat = fsChatOpenFlag
	//}
	chatId, chatIdErr := domain.GetProjectMainChatId(orgId, projectId)
	if chatIdErr != nil {
		log.Error(chatIdErr)
		return nil, chatIdErr
	}
	if chatId != "" {
		result.IsCreateFsChat = 1
	} else {
		result.IsCreateFsChat = 2
	}

	////获取我收藏的项目
	//startProject, err := domain.GetProjectRelationByCond(db.Cond{
	//	consts.TcOrgId:        orgId,
	//	consts.TcRelationType: consts.IssueRelationTypeStar,
	//	consts.TcRelationId:   userId,
	//	consts.TcIsDelete:     consts.AppIsNoDelete,
	//	consts.TcProjectId:    projectId,
	//})
	//if err != nil {
	//	log.Error(err)
	//	return nil, err
	//}
	//if startProject == nil || len(*startProject) == 0 {
	//	result.IsStar = 0
	//} else {
	//	result.IsStar = 1
	//}

	resp := appfacade.GetAppInfoList(appvo.GetAppInfoListReq{
		OrgId:  orgId,
		AppIds: []int64{projectBo.AppId},
	})
	if resp.Failure() {
		log.Error(resp.Error())
		return nil, resp.Error()
	}
	if len(resp.Data) > 0 {
		result.Icon = resp.Data[0].Icon
		result.ParentID = cast.ToString(resp.Data[0].ParentId)
	}

	return result, nil
}

// 通过项目类型langCode获取项目列表
func GetProjectBoListByProjectTypeLangCode(orgId int64, projectTypeLangCode *string) ([]bo.ProjectBo, errs.SystemErrorInfo) {
	return domain.GetProjectBoListByProjectTypeLangCode(orgId, projectTypeLangCode)
}

func GetSimpleProjectInfo(orgId int64, ids []int64) (*[]vo.Project, errs.SystemErrorInfo) {
	list, err := domain.GetProjectBoList(orgId, ids)
	if err != nil {
		return nil, err
	}
	projectVo := &[]vo.Project{}
	copyErr := copyer.Copy(list, projectVo)
	if err != nil {
		log.Error(err)
		return nil, errs.BuildSystemErrorInfo(errs.ObjectCopyError, copyErr)
	}

	return projectVo, nil
}

// 批量获取 project detail 信息
func GetProjectDetails(orgId int64, ids []int64) ([]bo.ProjectDetailBo, errs.SystemErrorInfo) {
	list, err := domain.GetProjectDetails(orgId, ids)
	if err != nil {
		return nil, err
	}
	detailBoList := make([]bo.ProjectDetailBo, 0)
	copyErr := copyer.Copy(list, &detailBoList)
	if err != nil {
		log.Error(err)
		return nil, errs.BuildSystemErrorInfo(errs.ObjectCopyError, copyErr)
	}
	for i, item := range detailBoList {
		detailBoList[i].SyncOutCalendarStatusArr = domain.TransferSyncOutCalendarStatus(item.IsSyncOutCalendar)
	}

	return detailBoList, nil
}

func GetProjectRelation(projectId int64, relationType []int64) ([]projectvo.ProjectRelationList, errs.SystemErrorInfo) {
	bos, err := domain.GetProjectRelationByType(projectId, relationType)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	res := []projectvo.ProjectRelationList{}
	for _, v := range *bos {
		res = append(res, projectvo.ProjectRelationList{
			Id:           v.Id,
			RelationId:   v.RelationId,
			RelationType: v.RelationType,
		})
	}

	return res, nil
}

func GetProjectRelationBatch(orgId int64, input *projectvo.GetProjectRelationBatchData) ([]bo.ProjectRelationBo, errs.SystemErrorInfo) {
	list, err := domain.GetProjectRelationByCond(db.Cond{
		consts.TcOrgId:        orgId,
		consts.TcProjectId:    input.ProjectIds,
		consts.TcRelationType: db.In(input.RelationTypes),
		consts.TcIsDelete:     consts.AppIsNoDelete,
	})
	if err != nil {
		log.Error(err)
		return nil, err
	}
	boList := make([]bo.ProjectRelationBo, 0)
	if err := copyer.Copy(list, &boList); err != nil {
		log.Errorf("[ObjectCopyError] copy err: %s", err)
		return nil, errs.BuildSystemErrorInfo(errs.ObjectCopyError, err)
	}

	return boList, nil
}

func ArchiveProject(orgId, userId, projectId int64) (*vo.Void, errs.SystemErrorInfo) {
	err := domain.AuthProject(orgId, userId, projectId, consts.RoleOperationPathOrgProProConfig, consts.OperationProConfigFiling)
	if err != nil {
		log.Error(err)
		return nil, errs.BuildSystemErrorInfo(errs.Unauthorized, err)
	}

	projectInfo, err := domain.GetProject(orgId, projectId)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	if projectInfo.IsFiling == consts.ProjectIsFiling {
		return nil, errs.ProjectIsFilingYet
	}

	_, updateErr := dao.UpdateProjectByOrg(projectId, orgId, mysql.Upd{
		consts.TcIsFiling: consts.ProjectIsFiling,
	})
	if updateErr != nil {
		return nil, errs.BuildSystemErrorInfo(errs.ProjectDomainError, updateErr)
	}

	err1 := domain.RefreshProjectAuthBo(orgId, projectId)
	if err1 != nil {
		log.Error(err)
		return nil, errs.BuildSystemErrorInfo(errs.ProjectDomainError, err1)
	}

	asyn.Execute(func() {
		PushModifyProjectNotice(orgId, projectId, userId)
	})
	return &vo.Void{
		ID: projectId,
	}, nil
}

func DeleteProject(orgId, userId, projectId int64) (*vo.Void, errs.SystemErrorInfo) {
	err := domain.AuthProject(orgId, userId, projectId, consts.RoleOperationPathOrgProProConfig, consts.OperationProConfigDelete)
	if err != nil {
		log.Error(err)
		return nil, errs.BuildSystemErrorInfo(errs.Unauthorized, err)
	}
	return DeleteProjectWithoutAuth(orgId, userId, projectId)
}

// DeleteProjectWithoutAuth 删除项目
func DeleteProjectWithoutAuth(orgId, userId, projectId int64) (*vo.Void, errs.SystemErrorInfo) {
	projectInfo, infoErr := domain.GetProjectSimple(orgId, projectId)
	if infoErr != nil {
		log.Error(infoErr)
		return nil, infoErr
	}
	// 如果有异步任务在执行，则暂不允许删应用操作
	isExecuting := domain.CheckAsyncTaskIsRunning(orgId, projectInfo.AppId, 0)
	if isExecuting {
		return nil, errs.DenyDeleteProWhenAsyncTask
	}

	updateErr := mysql.TransX(func(tx sqlbuilder.Tx) error {
		//删除项目
		_, updateProjectErr := dao.UpdateProjectByOrg(projectId, orgId, mysql.Upd{
			consts.TcIsDelete: consts.AppIsDeleted,
		}, tx)
		if updateProjectErr != nil {
			log.Error(updateProjectErr)
			return updateProjectErr
		}
		//删除项目下的任务
		//_, updateIssueErr := mysql.TransUpdateSmartWithCond(tx, consts.TableIssue, db.Cond{
		//	consts.TcIsDelete:  consts.AppIsNoDelete,
		//	consts.TcProjectId: projectId,
		//	consts.TcOrgId:     orgId,
		//}, mysql.Upd{
		//	consts.TcIsDelete: consts.AppIsDeleted,
		//})
		//if updateIssueErr != nil {
		//	log.Error(updateIssueErr)
		//	return updateIssueErr
		//}
		//删除迭代
		_, updateIterationErr := mysql.TransUpdateSmartWithCond(tx, consts.TableIteration, db.Cond{
			consts.TcIsDelete:  consts.AppIsNoDelete,
			consts.TcProjectId: projectId,
			consts.TcOrgId:     orgId,
		}, mysql.Upd{
			consts.TcIsDelete: consts.AppIsDeleted,
		})
		if updateIterationErr != nil {
			log.Error(updateIterationErr)
			return updateIterationErr
		}
		//删除无码表单
		if projectInfo.AppId > 0 {
			tables, err := domain.GetAppTableList(orgId, projectInfo.AppId)
			if err == nil {
				for _, table := range tables {
					issueIds, err := domain.DeleteTableIssues(orgId, userId, projectId, projectInfo.AppId, table.TableId, projectInfo.TemplateFlag)
					log.Infof("[DeleteProjectWithoutAuth] issueIds:%v", issueIds)
					if err != nil {
						log.Errorf("[DeleteProjectWithoutAuth] DeleteTableIssues err: %v, projectId: %d", err, projectId)
					}
				}
			}

			lcResp := appfacade.DeleteLessCodeApp(&appvo.DeleteLessCodeAppReq{
				AppId:  projectInfo.AppId,
				OrgId:  orgId,
				UserId: userId,
			})
			if lcResp.Failure() {
				log.Error(lcResp.Error())
				return lcResp.Error()
			}
		}

		return nil
	})
	if updateErr != nil {
		return nil, errs.DeleteProjectErr
	}

	err1 := domain.RefreshProjectAuthBo(orgId, projectId)
	if err1 != nil {
		log.Error(err1)
		return nil, errs.BuildSystemErrorInfo(errs.ProjectDomainError, err1)
	}

	asyn.Execute(func() {
		ext := bo.TrendExtensionBo{}
		ext.ObjName = projectInfo.Name
		projectTrendsBo := bo.ProjectTrendsBo{
			PushType:   consts.PushTypeDeleteProject,
			OrgId:      orgId,
			ProjectId:  projectId,
			OperatorId: userId,
			Ext:        ext,
		}
		domain.PushProjectTrends(projectTrendsBo)
	})
	asyn.Execute(func() {
		chatId, err := domain.GetProjectMainChatId(orgId, projectId)
		if err != nil {
			log.Error(err)
			return
		}
		if chatId != "" {
			err := DeleteFsChat(orgId, projectId, chatId)
			if err != nil {
				log.Errorf("[DeleteProjectWithoutAuth] failed:%v, chatId: %s", err, chatId)
				return
			}
		}
		// 删除该项目在其他群聊中的关联
		if err := domain.DeleteGroupChatSettingsByProjectId(orgId, projectId); err != nil {
			log.Errorf("[DeleteProjectWithoutAuth] DeleteGroupChatSettingsByProjectId err: %v, chatId: %s", err, chatId)
			return
		}
	})
	asyn.Execute(func() {
		PushDeleteProjectNotice(orgId, projectId, userId)
		// 删除项目后，需要删除对应的日历、所有任务对应的日程
		if err := domain.DeleteOneCalendar(orgId, projectId); err != nil {
			log.Errorf("删除项目时，后续处理，删除日历异常：%v", err)
		}
	})
	asyn.Execute(func() {
		resp := orgfacade.DeleteAppUserViewLocation(orgvo.UpdateUserViewLocationReq{
			OrgId:  orgId,
			UserId: userId,
			AppId:  projectInfo.AppId,
		})
		if resp.Failure() {
			log.Errorf("[DeleteProjectWithoutAuth] delete userLocation err:%v, orgId:%v, userId:%v, appId:%v",
				resp.Error(), orgId, userId, projectInfo.AppId)
		}

		//trendsResp := trendsfacade.DeleteTrends(trendsvo.DeleteTrendsReq{
		//	OrgId: orgId,
		//	Input: trendsvo.DeleteTrends{
		//		ProjectId: projectId,
		//	},
		//})
		//if trendsResp.Failure() {
		//	log.Errorf("[DeleteProjectWithoutAuth] DeleteTrends err:%v, orgId:%v, projectId:%v", resp.Error(), orgId, projectId)
		//	return
		//}
	})

	domain.DeleteDingCoolApp(orgId, projectId)

	return &vo.Void{
		ID: projectId,
	}, nil
}

func CancelArchivedProject(orgId, userId, projectId int64) (*vo.Void, errs.SystemErrorInfo) {
	err := domain.AuthProjectWithOutArchivedCheck(orgId, userId, projectId, consts.RoleOperationPathOrgProProConfig, consts.OperationProConfigUnFiling)
	if err != nil {
		log.Error(err)
		return nil, errs.BuildSystemErrorInfo(errs.Unauthorized, err)
	}

	_, err = domain.GetProjectSimple(orgId, projectId)
	if err != nil {
		log.Error(err)
		return nil, err
	}

	_, updateErr := dao.UpdateProjectByOrg(projectId, orgId, mysql.Upd{
		consts.TcIsFiling: consts.ProjectIsNotFiling,
	})
	if updateErr != nil {
		return nil, errs.BuildSystemErrorInfo(errs.ProjectDomainError, updateErr)
	}

	err1 := domain.RefreshProjectAuthBo(orgId, projectId)
	if err1 != nil {
		log.Error(err)
		return nil, errs.BuildSystemErrorInfo(errs.ProjectDomainError, err1)
	}

	asyn.Execute(func() {
		PushModifyProjectNotice(orgId, projectId, userId)
	})
	return &vo.Void{
		ID: projectId,
	}, nil
}

func GetProjectInfoByOrgIds(orgIds []int64) ([]projectvo.GetProjectInfoListByOrgIdsRespVo, errs.SystemErrorInfo) {

	bos, err := domain.GetProjectInfoByOrgIds(orgIds)

	if err != nil {
		return nil, err
	}

	result := []projectvo.GetProjectInfoListByOrgIdsRespVo{}

	for _, value := range bos {

		vo := projectvo.GetProjectInfoListByOrgIdsRespVo{
			OrgId:     value.OrgId,
			ProjectId: value.Id,
			Owner:     value.Owner,
		}
		result = append(result, vo)
	}

	return result, nil

}

func GetCacheProjectInfo(reqVo projectvo.GetCacheProjectInfoReqVo) (*bo.ProjectAuthBo, errs.SystemErrorInfo) {
	projectAuthBo, err := domain.LoadProjectAuthBo(reqVo.OrgId, reqVo.ProjectId)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	return projectAuthBo, nil
}

func OrgProjectMembers(input projectvo.OrgProjectMemberReqVo) (*projectvo.OrgProjectMemberRespVo, errs.SystemErrorInfo) {
	relationBo, err := domain.GetProjectRelationByType(input.ProjectId, []int64{consts.ProjectRelationTypeOwner, consts.ProjectRelationTypeParticipant})
	if err != nil {
		log.Error(err)
		return nil, err
	}

	//去重所有的用户id
	distinctUserIds, err := DistinctUserIds(relationBo)
	if err != nil {
		log.Error(err)
		return nil, err
	}

	userInfos, err := orgfacade.GetBaseUserInfoBatchRelaxed(input.OrgId, distinctUserIds)
	if err != nil {
		log.Error(err)
		return nil, err
	}

	allMembers := &[]projectvo.OrgProjectMemberVo{}
	_ = copyer.Copy(userInfos, allMembers)

	var ownerId int64
	//参与者
	participants := make(map[int64]bool)
	//关注者
	follower := make(map[int64]bool)

	//用户id分组
	for _, value := range *relationBo {
		if value.RelationType == consts.ProjectRelationTypeOwner {
			ownerId = value.RelationId
			continue
		}

		if value.RelationType == consts.ProjectRelationTypeParticipant {
			participants[value.RelationId] = true
			continue
		}

		//if value.RelationType == consts.IssueRelationTypeFollower {
		//	follower[value.RelationId] = true
		//	continue
		//}
	}
	//返回结果数组
	participantsMemberList := make([]projectvo.OrgProjectMemberVo, 0)
	followerMemberList := make([]projectvo.OrgProjectMemberVo, 0)

	ownerMember := projectvo.OrgProjectMemberVo{}

	for _, member := range *allMembers {
		//拥有者
		if member.UserId == ownerId {
			ownerMember = member
		}

		if _, ok := participants[member.UserId]; ok {
			participantsMemberList = append(participantsMemberList, member)
		}

		if _, ok := follower[member.UserId]; ok {
			followerMemberList = append(followerMemberList, member)
		}
	}

	return &projectvo.OrgProjectMemberRespVo{
		Owner:        ownerMember,
		Participants: participantsMemberList,
		Follower:     followerMemberList,
		AllMembers:   *allMembers,
	}, nil

}

func DistinctUserIds(relationBo *[]bo.ProjectRelationBo) ([]int64, errs.SystemErrorInfo) {

	rlen := len(*relationBo)

	if relationBo == nil || rlen < 1 {
		return nil, errs.BuildSystemErrorInfo(errs.ProjectRelationNotExist)
	}

	allUserIds := make([]int64, rlen)

	for i := 0; i < rlen; i++ {
		allUserIds[i] = (*relationBo)[i].RelationId
	}

	//去重
	uniqueInt64 := slice.SliceUniqueInt64(allUserIds)

	return uniqueInt64, nil
}

// 通过群聊的id，获取对应项目的id
func GetProjectIdByChatId(input projectvo.GetProjectIdByChatIdReqVo) (*projectvo.GetProjectIdByChatIdResp, errs.SystemErrorInfo) {
	result := &projectvo.GetProjectIdByChatIdResp{}
	orgId, projectId, err := domain.GetProjectIdByOpenChatId(input.OpenChatId)
	if err != nil {
		log.Error(err)
		return result, errs.BuildSystemErrorInfo(errs.ParamError, err)
	}
	result.OrgId = orgId
	result.ProjectId = projectId

	return result, nil
}

// GetProjectIdsByChatId 通过群聊id，查询该群聊关联的项目 id 列表
func GetProjectIdsByChatId(orgId int64, input projectvo.GetProjectIdsByChatIdReqVo) (*projectvo.GetProjectIdsByChatIdRespData, errs.SystemErrorInfo) {
	result := &projectvo.GetProjectIdsByChatIdRespData{}
	projectIds, err := domain.GetProjectIdsByChatId(orgId, input.OpenChatId)
	if err != nil {
		log.Errorf("[GetProjectIdsByChatId] err: %v, chatId: %s", err, input.OpenChatId)
		return result, errs.BuildSystemErrorInfo(errs.ParamError, err)
	}
	result.ProjectIds = projectIds

	return result, nil
}

// OpenPriorityList openAPI 获取项目的优先级列表
func OpenPriorityList(orgId int64) (*projectvo.OpenSomeAttrListRespVoData, errs.SystemErrorInfo) {
	resList := make([]*projectvo.OpenSomeAttrListRespVoDataItem, 0)
	//priorities, err := domain.GetPriorityListByType(orgId, consts.PriorityTypeIssue)
	//if err != nil {
	//	log.Error(err)
	//	return nil, err
	//}
	//
	//for _, item := range *priorities {
	//	resList = append(resList, &projectvo.OpenSomeAttrListRespVoDataItem{
	//		Name: item.Name,
	//		Id:   item.Id,
	//	})
	//}
	return &projectvo.OpenSomeAttrListRespVoData{
		List: resList,
	}, nil
}

// OpenGetProjectObjectTypeList openAPI 获取项目的任务类型列表，缺陷、需求、任务等
func OpenGetProjectObjectTypeList(reqVo projectvo.OpenPriorityListReqVo) (*projectvo.OpenSomeAttrListRespVoData, errs.SystemErrorInfo) {
	resItemList := make([]*projectvo.OpenSomeAttrListRespVoDataItem, 0)
	//resp, err := ProjectObjectTypesWithProject(reqVo.OrgId, reqVo.ProjectId)
	//if err != nil {
	//	log.Error(err)
	//	return nil, err
	//}
	//for _, item := range resp.List {
	//	// 迭代不算是项目中的任务类型，只需需求、任务、缺陷。
	//	if item.LangCode == consts.ProjectObjectTypeLangCodeIteration {
	//		continue
	//	}
	//	resItemList = append(resItemList, &projectvo.OpenSomeAttrListRespVoDataItem{
	//		Id:   item.ID,
	//		Name: item.Name,
	//	})
	//}

	return &projectvo.OpenSomeAttrListRespVoData{
		List: resItemList,
	}, nil
}

// OpenGetIterationList openAPI 获取项目的迭代列表
func OpenGetIterationList(reqVo projectvo.OpenGetIterationListReqVo) (*projectvo.OpenGetIterationListRespVoData, errs.SystemErrorInfo) {
	resItemList := make([]*projectvo.OpenGetIterationListRespVoDataItem, 0)
	resp, err := IterationList(reqVo.OrgId, 1, 20000, &vo.IterationListReq{
		ProjectID: &reqVo.ProjectId,
	})
	if err != nil {
		log.Error(err)
		return nil, err
	}
	for _, item := range resp.List {
		resItemList = append(resItemList, &projectvo.OpenGetIterationListRespVoDataItem{
			Id:   item.ID,
			Name: item.Name,
		})
	}
	return &projectvo.OpenGetIterationListRespVoData{
		List: resItemList,
	}, nil
}

// OpenGetDemandSourceList openAPI 获取需求来源列表
func OpenGetDemandSourceList(reqVo projectvo.OpenGetDemandSourceListReqVo) (*projectvo.OpenSomeAttrListRespVoData, errs.SystemErrorInfo) {
	// resItemList := make([]*projectvo.OpenSomeAttrListRespVoDataItem, 0)
	// resp, err := IssueSourceList(reqVo.OrgId, 1, 20000, vo.IssueSourcesReq{
	// 	ProjectID: reqVo.ProjectId,
	// })
	// if err != nil {
	// 	log.Error(err)
	// 	return nil, err
	// }
	// for _, item := range resp.List {
	// 	resItemList = append(resItemList, &projectvo.OpenSomeAttrListRespVoDataItem{
	// 		Id:   item.ID,
	// 		Name: item.Name,
	// 	})
	// }
	// return &projectvo.OpenSomeAttrListRespVoData{
	// 	List: resItemList,
	// }, nil
	return nil, nil
}

// OpenGetPropertyList openAPI 获取严重程度列表
func OpenGetPropertyList(reqVo projectvo.OpenGetPropertyListReqVo) (*projectvo.OpenSomeAttrListRespVoData, errs.SystemErrorInfo) {
	//resItemList := make([]*projectvo.OpenSomeAttrListRespVoDataItem, 0)
	//resp, err := IssuePropertyList(reqVo.OrgId, vo.IssuePropertysReq{
	//	ProjectID: reqVo.ProjectId,
	//})
	//if err != nil {
	//	log.Error(err)
	//	return nil, err
	//}
	//for _, item := range resp.List {
	//	resItemList = append(resItemList, &projectvo.OpenSomeAttrListRespVoDataItem{
	//		Id:   item.ID,
	//		Name: item.Name,
	//	})
	//}
	//return &projectvo.OpenSomeAttrListRespVoData{
	//	List: resItemList,
	//}, nil
	return nil, nil
}

func GetSimpleProjectsByOrgId(orgId int64) ([]bo.ProjectBo, errs.SystemErrorInfo) {
	infoPos := &[]po.PpmProProject{}
	err := mysql.SelectAllByCond(consts.TableProject, db.Cond{
		consts.TcOrgId:    orgId,
		consts.TcIsDelete: consts.AppIsNoDelete,
	}, infoPos)

	if err != nil {
		log.Error(err)
		return nil, errs.MysqlOperateError
	}

	res := &[]bo.ProjectBo{}
	copyErr := copyer.Copy(infoPos, res)
	if copyErr != nil {
		log.Error(copyErr)
		return nil, errs.ObjectCopyError
	}

	return *res, nil
}

// 这个方法仅仅是给PA私有化上报用的，其他地方没有用到
func GetOrgIssueAndProjectCount(orgId int64) (*projectvo.GetOrgIssueAndProjectCountRespData, errs.SystemErrorInfo) {
	//issueCount, mysqlErr := mysql.SelectCountByCond(consts.TableIssue, db.Cond{
	//	consts.TcOrgId:    orgId,
	//	consts.TcIsDelete: consts.AppIsNoDelete,
	//})
	//if mysqlErr != nil {
	//	log.Error(mysqlErr)
	//	return nil, errs.MysqlOperateError
	//}
	//projectCount, mysqlErr := mysql.SelectCountByCond(consts.TableProject, db.Cond{
	//	consts.TcOrgId:    orgId,
	//	consts.TcIsDelete: consts.AppIsNoDelete,
	//})
	//if mysqlErr != nil {
	//	log.Error(mysqlErr)
	//	return nil, errs.MysqlOperateError
	//}
	//return &projectvo.GetOrgIssueAndProjectCountRespData{
	//	IssueCount:   issueCount,
	//	ProjectCount: projectCount,
	//}, nil
	return nil, nil
}

// QueryProcessForAsyncTask 异步任务查询进度条
func QueryProcessForAsyncTask(orgId int64, input *projectvo.QueryProcessForAsyncTaskReqVoData) (projectvo.AsyncTask, errs.SystemErrorInfo) {
	result := projectvo.AsyncTask{}
	// 根据 taskId 去 redis 查询已处理的任务数量
	taskInfo, err := domain.GetAsyncTask(orgId, input.TaskId)
	if err != nil {
		log.Errorf("[QueryProcessForAsyncTask] GetAsyncTask err: %v, orgId: %d, taskId: %s", err, orgId, input.TaskId)
		return result, err
	}
	if taskInfo.Total < 1 {
		// 和前端约定：没有异步任务时，暂时视为进度已达 100%
		taskInfo.Total = 0
		taskInfo.Processed = 0
		taskInfo.PercentVal = 100
		return *taskInfo, err
	}
	// 如果有错误码，则直接报错返回
	if taskInfo.ErrCode != -1 {
		err := errs.GetResultCodeInfoByCode(taskInfo.ErrCode)
		if err.Code() == errs.SystemError.Code() {
			err = errs.ImportIssueCellInvalid
		}
		log.Errorf("[QueryProcessForAsyncTask] GetAsyncTask err: %v, orgId: %d, taskId: %s", err, orgId, input.TaskId)
		return result, err
	}
	if taskInfo.PercentVal != 100 && time.Since(time.Unix(taskInfo.StartTimestamp, 0)).Seconds() > consts.AsyncTaskProcessTimeout {
		err := errs.ImportIssueFailed
		log.Errorf("[QueryProcessForAsyncTask] GetAsyncTask err: %v, orgId: %d, taskId: %s", err, orgId, input.TaskId)
		return result, err
	}

	return *taskInfo, nil
}
