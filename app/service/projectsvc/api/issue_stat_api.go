package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/star-table/startable-server/app/service/projectsvc/service"
	"github.com/star-table/startable-server/common/core/errs"
	"github.com/star-table/startable-server/common/model/vo"
	"github.com/star-table/startable-server/common/model/vo/projectvo"
)

func IssueAssignRank(c *gin.Context) {
	var reqVo projectvo.IssueAssignRankReqVo
	if err := c.ShouldBindJSON(&reqVo); err != nil {
		c.JSON(http.StatusBadRequest, projectvo.IssueAssignRankRespVo{Err: vo.NewErr(errs.ReqParamsValidateError), IssueAssignRankResp: nil})
		return
	}

	input := reqVo.Input
	projectId := input.ProjectID
	rankTop := 5
	if input.RankTop != nil {
		rt := *input.RankTop
		if rt >= 1 && rt <= 100 {
			rankTop = rt
		}
	}
	res, err := service.IssueAssignRank(reqVo.OrgId, projectId, rankTop)
	c.JSON(http.StatusOK, projectvo.IssueAssignRankRespVo{Err: vo.NewErr(err), IssueAssignRankResp: res})
}

func IssueAndProjectCountStat(c *gin.Context) {
	var reqVo projectvo.IssueAndProjectCountStatReqVo
	if err := c.ShouldBindJSON(&reqVo); err != nil {
		c.JSON(http.StatusBadRequest, projectvo.IssueAndProjectCountStatRespVo{Err: vo.NewErr(errs.ReqParamsValidateError), Data: nil})
		return
	}

	res, err := service.IssueAndProjectCountStat(reqVo)
	c.JSON(http.StatusOK, projectvo.IssueAndProjectCountStatRespVo{Err: vo.NewErr(err), Data: res})
}

func IssueDailyPersonalWorkCompletionStat(c *gin.Context) {
	var reqVo projectvo.IssueDailyPersonalWorkCompletionStatReqVo
	if err := c.ShouldBindJSON(&reqVo); err != nil {
		c.JSON(http.StatusBadRequest, projectvo.IssueDailyPersonalWorkCompletionStatRespVo{Err: vo.NewErr(errs.ReqParamsValidateError), Data: nil})
		return
	}

	res, err := service.IssueDailyPersonalWorkCompletionStat(reqVo)
	c.JSON(http.StatusOK, projectvo.IssueDailyPersonalWorkCompletionStatRespVo{Err: vo.NewErr(err), Data: res})
}

//func (PostGreeter) GetIssueCountByStatus(req projectvo.GetIssueCountByStatusReqVo) projectvo.GetIssueCountByStatusRespVo {
//	res, err := service.GetIssueCountByStatus(req.OrgId, req.ProjectId, req.StatusId)
//	return projectvo.GetIssueCountByStatusRespVo{
//		Err: vo.NewErr(err),
//		Data: projectvo.GetIssueCountByStatusData{
//			Count: res,
//		},
//	}
//}

func AuthProject(c *gin.Context) {
	var req projectvo.AuthProjectReqVo
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, vo.CommonRespVo{Err: vo.NewErr(errs.ReqParamsValidateError), Void: nil})
		return
	}

	err := service.AuthProject(req.OrgId, req.UserId, req.ProjectId, req.Path, req.Operation)
	c.JSON(http.StatusOK, vo.CommonRespVo{
		Err:  vo.NewErr(err),
		Void: nil,
	})
}
