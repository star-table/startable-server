package service

import (
	"github.com/star-table/startable-server/app/facade/resourcefacade"
	"github.com/star-table/startable-server/app/service/projectsvc/domain"
	"github.com/star-table/startable-server/common/core/errs"
	"github.com/star-table/startable-server/common/core/util"
	"github.com/star-table/startable-server/common/model/bo"
	"github.com/star-table/startable-server/common/model/vo/projectvo"
	"github.com/star-table/startable-server/common/model/vo/resourcevo"
)

// PayLimitNum 组织下已经使用的资源数量统计
func PayLimitNum(orgId int64) (*projectvo.PayLimitNumRespData, errs.SystemErrorInfo) {
	return domain.PayLimitNum(orgId)
}

func GetProjectStatistics(orgId, userId int64, req projectvo.GetProjectStatisticsReq) ([]*projectvo.GetProjectStatisticsData, errs.SystemErrorInfo) {
	params := req.ExtraInfoReq
	projectIds := make([]int64, 0, len(params))
	resourceIds := make([]int64, 0, len(params))
	for _, info := range params {
		projectIds = append(projectIds, info.ProjectId)
		resourceIds = append(resourceIds, info.ResourceId)
	}
	summaryAppId, err := domain.GetOrgSummaryAppId(orgId)
	if err != nil {
		log.Errorf("[GetProjectStatistics]GetOrgSummaryAppId err:%v, orgId:%v", err, orgId)
		return nil, err
	}

	statForProjects, err := domain.GetIssueFinishedStatForProjects(orgId, userId, summaryAppId, projectIds)
	if err != nil {
		log.Errorf("[GetProjectStatistics]GetIssueFinishedStatForProjects err:%v, orgId:%v, userId:%v, projectId:%v",
			err, orgId, userId, projectIds)
		return nil, err
	}
	resourceBosMap := make(map[int64]bo.ResourceBo)
	if len(resourceIds) > 0 {
		resourceResp := resourcefacade.GetResourceById(
			resourcevo.GetResourceByIdReqVo{GetResourceByIdReqBody: resourcevo.GetResourceByIdReqBody{ResourceIds: resourceIds}},
		)
		if resourceResp.Failure() {
			log.Errorf("[GetProjectStatistics] GetResourceById err:%v, orgId:%v", resourceResp.Error(), orgId)
			return nil, resourceResp.Error()
		}
		for _, res := range resourceResp.ResourceBos {
			resourceBosMap[res.Id] = res
		}
	}

	res := []*projectvo.GetProjectStatisticsData{}

	for _, info := range params {
		data := &projectvo.GetProjectStatisticsData{
			ProjectId:       info.ProjectId,
			CoverResourceId: info.ResourceId,
		}
		if stat, ok := statForProjects[info.ProjectId]; ok {
			data.AllIssues = stat.All
			data.FinishIssues = stat.Finish
			data.OverdueIssues = stat.Overdue
			data.UnfinishedIssues = stat.RelateUnfinish
		}
		if resourceBo, ok := resourceBosMap[info.ResourceId]; ok {
			coverUrl := util.JointUrl(resourceBo.Host, resourceBo.Path)
			resourcePath := util.GetCompressedPath(coverUrl, resourceBo.Type)
			data.CoverResourcePath = resourcePath
		}
		res = append(res, data)
	}
	return res, nil
}
