package api

import (
	"github.com/star-table/startable-server/app/service/projectsvc/domain"
	"github.com/star-table/startable-server/common/model/vo"
	"github.com/star-table/startable-server/common/model/vo/projectvo"
)

func (PostGreeter) GetIssueLinks(reqVo projectvo.GetIssueLinksReqVo) projectvo.GetIssueLinksRespVo {
	return projectvo.GetIssueLinksRespVo{
		Err:  vo.NewErr(nil),
		Data: domain.GetIssueLinks(reqVo.SourceChannel, reqVo.OrgId, reqVo.IssueId),
	}
}
