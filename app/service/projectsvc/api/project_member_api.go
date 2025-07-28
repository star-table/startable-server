package api

import (
	"github.com/star-table/startable-server/app/service/projectsvc/service"
	"github.com/star-table/startable-server/common/model/vo"
	"github.com/star-table/startable-server/common/model/vo/projectvo"
)

func (PostGreeter) GetProjectRelationUserIds(req projectvo.GetProjectRelationUserIdsReq) projectvo.GetProjectRelationUserIdsResp {
	res, err := service.GetProjectRelationUserIds(req.ProjectId, req.RelationType)
	return projectvo.GetProjectRelationUserIdsResp{
		Err:     vo.NewErr(err),
		UserIds: res,
	}
}

func (PostGreeter) ProjectMemberIdList(req projectvo.ProjectMemberIdListReq) projectvo.ProjectMemberIdListResp {
	res, err := service.ProjectMemberIdList(req.OrgId, req.ProjectId, req.Data)
	return projectvo.ProjectMemberIdListResp{
		Err:  vo.NewErr(err),
		Data: res,
	}
}

func (PostGreeter) GetProjectMemberIds(req projectvo.GetProjectMemberIdsReqVo) projectvo.GetProjectMemberIdsResp {
	res, err := service.GetProjectMemberIds(req.OrgId, req.Input.ProjectId, req.Input.IncludeAdmin)
	return projectvo.GetProjectMemberIdsResp{
		Err:  vo.NewErr(err),
		Data: res,
	}
}

func (PostGreeter) GetTrendsMembers(req projectvo.GetTrendListMembersReqVo) projectvo.GetTrendListMembersResp {
	res, err := service.GetTrendsMembers(req.OrgId, req.UserId, req.Input.ProjectId)
	return projectvo.GetTrendListMembersResp{
		Err:  vo.NewErr(err),
		Data: res,
	}
}
