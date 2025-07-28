package service

import (
	"github.com/star-table/startable-server/app/facade/orgfacade"
	"github.com/star-table/startable-server/app/facade/userfacade"
	"github.com/star-table/startable-server/app/service/projectsvc/domain"
	"github.com/star-table/startable-server/common/core/consts"
	"github.com/star-table/startable-server/common/core/errs"
	"github.com/star-table/startable-server/common/core/util/slice"
	"github.com/star-table/startable-server/common/model/vo"
	"github.com/star-table/startable-server/common/model/vo/orgvo"
	"github.com/star-table/startable-server/common/model/vo/projectvo"
)

func GetProjectRelationUserIds(projectId int64, relationType *int64) ([]int64, errs.SystemErrorInfo) {
	typeArr := []int64{consts.ProjectRelationTypeOwner, consts.ProjectRelationTypeParticipant}
	if relationType != nil {
		if *relationType == 1 {
			//负责人
			typeArr = []int64{consts.ProjectRelationTypeOwner}
		} else if *relationType == 2 {
			//参与人
			typeArr = []int64{consts.ProjectRelationTypeParticipant}
		}
	}
	relationBo, err := domain.GetProjectRelationByType(projectId, typeArr)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	//去重所有的用户id
	distinctUserIds, _ := DistinctUserIds(relationBo)

	return distinctUserIds, nil
}

func ProjectMemberIdList(orgId int64, projectId int64, input *projectvo.ProjectMemberIdListReqData) (*vo.ProjectMemberIDListResp, errs.SystemErrorInfo) {
	if input == nil {
		input = &projectvo.ProjectMemberIdListReqData{}
	}
	relations, err := domain.GetProjectMembers(orgId, projectId, []int64{consts.ProjectRelationTypeDepartmentParticipant, consts.ProjectRelationTypeParticipant, consts.ProjectRelationTypeOwner})
	if err != nil {
		log.Error(err)
		return nil, err
	}
	appId := int64(-1)
	if projectId > 0 {
		projectBo, err := domain.GetProjectSimple(orgId, projectId)
		if err != nil {
			log.Errorf("[GetProjectMemberIds]GetProjectSimple err:%v, orgId:%v, projectId:%v", err, orgId, projectId)
			return nil, err
		}
		appId = projectBo.AppId
	}
	adminUserIds, err := GetAdminOfProject(orgId, appId)
	if err != nil {
		log.Error(err)
		return nil, err
	}

	var deptIds, userIds []int64
	for _, relation := range relations {
		if relation.RelationType == consts.ProjectRelationTypeDepartmentParticipant {
			deptIds = append(deptIds, relation.RelationId)
		} else {
			userIds = append(userIds, relation.RelationId)
		}
	}

	if input.IncludeAdmin == 1 {
		// 追加上管理员的 id
		userIds = append(userIds, adminUserIds...)
	}

	return &vo.ProjectMemberIDListResp{
		DepartmentIds: slice.SliceUniqueInt64(deptIds),
		UserIds:       slice.SliceUniqueInt64(userIds),
	}, nil
}

// 获取项目的管理员列表
func GetAdminOfProject(orgId, appId int64) ([]int64, errs.SystemErrorInfo) {
	adminUserResp := userfacade.GetUsersCouldManage(orgId, appId)
	if adminUserResp.Failure() {
		log.Error(adminUserResp.Error())
		return nil, adminUserResp.Error()
	}
	adminUserIds := make([]int64, 0)
	for _, user := range adminUserResp.Data.List {
		adminUserIds = append(adminUserIds, user.Id)
	}
	return adminUserIds, nil
}

// GetProjectMemberIds 获取项目成员ids列表，rest api调用
func GetProjectMemberIds(orgId, projectId int64, includeAdmin int) (*projectvo.GetProjectMemberIdsData, errs.SystemErrorInfo) {
	relations, err := domain.GetProjectMembers(orgId, projectId, []int64{consts.ProjectRelationTypeDepartmentParticipant, consts.ProjectRelationTypeParticipant, consts.ProjectRelationTypeOwner})
	if err != nil {
		log.Errorf("[GetProjectMemberIds] GetProjectMembers err:%v, orgId:%v, projectId:%v", err, orgId, projectId)
		return nil, err
	}
	var deptIds, userIds []int64
	for _, relation := range relations {
		if relation.RelationType == consts.ProjectRelationTypeDepartmentParticipant {
			deptIds = append(deptIds, relation.RelationId)
		} else {
			userIds = append(userIds, relation.RelationId)
		}
	}

	appId := int64(-1)
	if projectId > 0 {
		projectBo, err := domain.GetProjectSimple(orgId, projectId)
		if err != nil {
			log.Errorf("[GetProjectMemberIds]GetProjectSimple err:%v, orgId:%v, projectId:%v", err, orgId, projectId)
			return nil, err
		}
		appId = projectBo.AppId
	}

	if includeAdmin == consts.IncludeAdmin {
		// 包含管理员
		adminUserIds, err := GetAdminOfProject(orgId, appId)
		if err != nil {
			log.Errorf("[GetProjectMemberIds] GetAdminOfProject err:%v, orgId:%v, projectId:%v", err, orgId, projectId)
			return nil, err
		}
		userIds = append(userIds, adminUserIds...)
	}

	return &projectvo.GetProjectMemberIdsData{
		DepartmentIds: slice.SliceUniqueInt64(deptIds),
		UserIds:       slice.SliceUniqueInt64(userIds),
	}, nil
}

// 超管-选择所有项目或者某个项目都显示所有成员
// 项目管理员-选择管理的项目显示项目所有成员
// 如果是团队成员，只可以选择当前用户
func GetTrendsMembers(orgId, userId, projectId int64) (*projectvo.GetTrendListMembersData, errs.SystemErrorInfo) {
	userIds := []int64{}
	userIds = append(userIds, userId)

	appId := int64(0)
	if projectId > 0 {
		appIdResp, err := domain.GetAppIdFromProjectId(orgId, projectId)
		if err != nil {
			log.Errorf("[GetTrendsMembers] domain.GetAppIdFromProjectId err:%v, orgId:%v, userId:%v, projectId:%v",
				err, orgId, userId, projectId)
			return nil, err
		}
		appId = appIdResp
	}

	projectAdminIds, err := domain.GetProjectAdminIds(orgId, appId)
	if err != nil {
		log.Errorf("[GetTrendsMembers] domain.GetProjectAdminIds err:%v, orgId:%v, userId:%v, projectId:%v",
			err, orgId, userId, projectId)
		return nil, err
	}
	if ok, er := slice.Contain(projectAdminIds, userId); er == nil && ok {
		// 是管理员
		projectMemberData, err := GetProjectMemberIds(orgId, projectId, consts.IncludeAdmin)
		if err != nil {
			log.Errorf("[GetTrendsMembers] GetProjectMemberIds err:%v, orgId:%v, userId:%v, projectId:%v",
				err, orgId, userId, projectId)
			return nil, err
		}
		userIds = append(userIds, projectMemberData.UserIds...)

		if len(projectMemberData.DepartmentIds) > 0 {
			resp := orgfacade.GetUserIdsByDeptIds(&orgvo.GetUserIdsByDeptIdsReq{
				OrgId:   orgId,
				DeptIds: projectMemberData.DepartmentIds,
			})
			if resp.Failure() {
				log.Errorf("[GetTrendsMembers] orgfacade.GetUserIdsByDeptIds err:%v, orgId:%v, userId:%v, projectId:%v",
					resp.Error(), orgId, userId, projectId)
				return nil, resp.Error()
			}
			userIds = append(userIds, resp.Data.UserIds...)
		}
	}

	userIds = slice.SliceUniqueInt64(userIds)

	userInfoBatch := orgfacade.GetBaseUserInfoBatch(orgvo.GetBaseUserInfoBatchReqVo{
		OrgId:   orgId,
		UserIds: userIds,
	})

	if userInfoBatch.Failure() {
		log.Errorf("[GetTrendsMembers] GetBaseUserInfoBatch err:%v, orgId:%v, userId:%v, projectId:%v",
			userInfoBatch.Error(), orgId, userId, projectId)
		return nil, userInfoBatch.Error()
	}

	users := []projectvo.UserInfo{}
	for _, u := range userInfoBatch.BaseUserInfos {
		users = append(users, projectvo.UserInfo{
			Id:   u.UserId,
			Name: u.Name,
		})
	}

	return &projectvo.GetTrendListMembersData{Users: users}, nil
}
