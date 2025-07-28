package api

import (
	"github.com/star-table/startable-server/app/service/projectsvc/service"
	"github.com/star-table/startable-server/common/core/util/json"
	"github.com/star-table/startable-server/common/model/vo"
	"github.com/star-table/startable-server/common/model/vo/projectvo"
)

func (PostGreeter) LcViewStatForAll(reqVo *projectvo.LcViewStatReqVo) *projectvo.LcViewStatRespVo {
	viewStats, err := service.LcViewStatForAll(reqVo.OrgId, reqVo.UserId)
	return &projectvo.LcViewStatRespVo{Err: vo.NewErr(err), Data: viewStats}
}

func (PostGreeter) LcHomeIssues(reqVo *projectvo.HomeIssuesReqVo) string {
	if reqVo.Size > 30000 {
		reqVo.Size = 30000
	}
	res, err := service.LcHomeIssuesForProject(reqVo.OrgId, reqVo.UserId, reqVo.Page, reqVo.Size, reqVo.Input, false)
	if err != nil {
		return json.ToJsonIgnoreError(&projectvo.LcHomeIssuesRespVo{Err: vo.NewErr(err), Data: "{}"})
	}

	return res.Data
}

func (PostGreeter) LcHomeIssuesForAll(reqVo *projectvo.HomeIssuesReqVo) string {
	if reqVo.Size > 30000 {
		reqVo.Size = 30000
	}
	res, err := service.LcHomeIssuesForAll(reqVo.OrgId, reqVo.UserId, reqVo.Page, reqVo.Size, reqVo.Input, false)
	if err != nil {
		return json.ToJsonIgnoreError(&projectvo.LcHomeIssuesRespVo{Err: vo.NewErr(err), Data: "{}"})
	}

	return res.Data
}

func (PostGreeter) LcHomeIssuesForIssue(reqVo *projectvo.IssueDetailReqVo) *projectvo.IssueDetailRespVo {
	res, err := service.LcHomeIssuesForIssue(reqVo.OrgId, reqVo.UserId, reqVo.AppId, reqVo.TableId, reqVo.IssueId, reqVo.TodoId, false)
	if err != nil {
		return &projectvo.IssueDetailRespVo{Err: vo.NewErr(err)}
	}
	return &projectvo.IssueDetailRespVo{Err: vo.NewErr(err), Data: res}
}

func (PostGreeter) IssueStatusTypeStat(reqVo projectvo.IssueStatusTypeStatReqVo) projectvo.IssueStatusTypeStatRespVo {
	res, err := service.IssueStatusTypeStat(reqVo.OrgId, reqVo.UserId, reqVo.Input)
	return projectvo.IssueStatusTypeStatRespVo{Err: vo.NewErr(err), IssueStatusTypeStat: res}
}

func (PostGreeter) IssueStatusTypeStatDetail(reqVo projectvo.IssueStatusTypeStatReqVo) projectvo.IssueStatusTypeStatDetailRespVo {
	res, err := service.IssueStatusTypeStatDetail(reqVo.OrgId, reqVo.UserId, reqVo.Input)
	return projectvo.IssueStatusTypeStatDetailRespVo{Err: vo.NewErr(err), IssueStatusTypeStatDetail: res}
}

func (PostGreeter) GetLcIssueInfoBatch(reqVo projectvo.GetLcIssueInfoBatchReqVo) projectvo.GetLcIssueInfoBatchRespVo {
	res, err := service.GetLcIssueInfoBatch(reqVo.OrgId, reqVo.IssueIds)
	return projectvo.GetLcIssueInfoBatchRespVo{
		Err:  vo.NewErr(err),
		Data: res,
	}
}

func (PostGreeter) IssueListStat(reqVo projectvo.IssueListStatReq) projectvo.IssueListStatResp {
	res, err := service.IssueListStat(reqVo.OrgId, reqVo.UserId, reqVo.Input.ProjectID)
	return projectvo.IssueListStatResp{Err: vo.NewErr(err), Data: res}
}

func (PostGreeter) IssueListSimpleByDataIds(reqVo projectvo.GetIssueListSimpleByDataIdsReqVo) projectvo.SimpleIssueListRespVo {
	res, err := service.IssueListSimpleByDataIds(reqVo.OrgId, reqVo.UserId, reqVo.Input.AppId, reqVo.Input.DataIds)
	return projectvo.SimpleIssueListRespVo{
		Err:  vo.NewErr(err),
		Data: res,
	}
}
