package projectvo

import "github.com/star-table/startable-server/common/model/vo"

type ProjectIssueRelatedStatusRespVo struct {
	vo.Err
	ProjectIssueRelatedStatus []*vo.HomeIssueStatusInfo `json:"data"`
}

type ProjectIssueRelatedStatusReqVo struct {
	Input vo.ProjectIssueRelatedStatusReq `json:"input"`
	OrgId int64                           `json:"orgId"`
}
