package api

import (
	"github.com/star-table/startable-server/app/service/projectsvc/service"
	"github.com/star-table/startable-server/common/model/vo"
	"github.com/star-table/startable-server/common/model/vo/projectvo"
)

func (PostGreeter) ProjectIssueRelatedStatus(reqVo projectvo.ProjectIssueRelatedStatusReqVo) projectvo.ProjectIssueRelatedStatusRespVo {
	res, err := service.ProjectIssueRelatedStatus(reqVo.OrgId, reqVo.Input)
	return projectvo.ProjectIssueRelatedStatusRespVo{Err: vo.NewErr(err), ProjectIssueRelatedStatus: res}
}
