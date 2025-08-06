package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/star-table/startable-server/app/service/projectsvc/service"
	"github.com/star-table/startable-server/common/core/errs"
	"github.com/star-table/startable-server/common/core/util/json"
	"github.com/star-table/startable-server/common/model/vo"
	"github.com/star-table/startable-server/common/model/vo/projectvo"
)

func LcViewStatForAll(c *gin.Context) {
	var reqVo projectvo.LcViewStatReqVo
	if err := c.ShouldBindJSON(&reqVo); err != nil {
		// Replaced: c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		c.JSON(http.StatusBadRequest, projectvo.LcViewStatRespVo{Err: vo.NewErr(errs.ReqParamsValidateError)})
		return
	}

	viewStats, err := service.LcViewStatForAll(reqVo.OrgId, reqVo.UserId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, projectvo.LcViewStatRespVo{Err: vo.NewErr(err), Data: viewStats})
}

func LcHomeIssues(c *gin.Context) {
	var reqVo projectvo.HomeIssuesReqVo
	if err := c.ShouldBindJSON(&reqVo); err != nil {
		// Replaced: c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		c.String(http.StatusBadRequest, json.ToJsonIgnoreError(&projectvo.LcHomeIssuesRespVo{Err: vo.NewErr(errs.ReqParamsValidateError), Data: "{}"}))
		return
	}

	if reqVo.Size > 30000 {
		reqVo.Size = 30000
	}
	res, err := service.LcHomeIssuesForProject(reqVo.OrgId, reqVo.UserId, reqVo.Page, reqVo.Size, reqVo.Input, false)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.String(http.StatusOK, res.Data)
}

func LcHomeIssuesForAll(c *gin.Context) {
	var reqVo projectvo.HomeIssuesReqVo
	if err := c.ShouldBindJSON(&reqVo); err != nil {
		// Replaced: c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		c.String(http.StatusBadRequest, json.ToJsonIgnoreError(&projectvo.LcHomeIssuesRespVo{Err: vo.NewErr(errs.ReqParamsValidateError), Data: "{}"}))
		return
	}

	if reqVo.Size > 30000 {
		reqVo.Size = 30000
	}
	res, err := service.LcHomeIssuesForAll(reqVo.OrgId, reqVo.UserId, reqVo.Page, reqVo.Size, reqVo.Input, false)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.String(http.StatusOK, res.Data)
}

func LcHomeIssuesForIssue(c *gin.Context) {
	var reqVo projectvo.IssueDetailReqVo
	if err := c.ShouldBindJSON(&reqVo); err != nil {
		// Replaced: c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		c.JSON(http.StatusBadRequest, projectvo.IssueDetailRespVo{Err: vo.NewErr(errs.ReqParamsValidateError)})
		return
	}

	res, err := service.LcHomeIssuesForIssue(reqVo.OrgId, reqVo.UserId, reqVo.AppId, reqVo.TableId, reqVo.IssueId, reqVo.TodoId, false)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, projectvo.IssueDetailRespVo{Err: vo.NewErr(nil), Data: res})
}

func IssueStatusTypeStat(c *gin.Context) {
	var reqVo projectvo.IssueStatusTypeStatReqVo
	if err := c.ShouldBindJSON(&reqVo); err != nil {
		// Replaced: c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		c.JSON(http.StatusBadRequest, projectvo.IssueStatusTypeStatRespVo{Err: vo.NewErr(errs.ReqParamsValidateError)})
		return
	}

	res, err := service.IssueStatusTypeStat(reqVo.OrgId, reqVo.UserId, reqVo.Input)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, projectvo.IssueStatusTypeStatRespVo{Err: vo.NewErr(err), IssueStatusTypeStat: res})
}

func IssueStatusTypeStatDetail(c *gin.Context) {
	var reqVo projectvo.IssueStatusTypeStatReqVo
	if err := c.ShouldBindJSON(&reqVo); err != nil {
		// Replaced: c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		c.JSON(http.StatusBadRequest, projectvo.IssueStatusTypeStatDetailRespVo{Err: vo.NewErr(errs.ReqParamsValidateError)})
		return
	}

	res, err := service.IssueStatusTypeStatDetail(reqVo.OrgId, reqVo.UserId, reqVo.Input)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, projectvo.IssueStatusTypeStatDetailRespVo{Err: vo.NewErr(err), IssueStatusTypeStatDetail: res})
}

func GetLcIssueInfoBatch(c *gin.Context) {
	var reqVo projectvo.GetLcIssueInfoBatchReqVo
	if err := c.ShouldBindJSON(&reqVo); err != nil {
		// Replaced: c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		c.JSON(http.StatusBadRequest, projectvo.GetLcIssueInfoBatchRespVo{Err: vo.NewErr(errs.ReqParamsValidateError)})
		return
	}

	res, err := service.GetLcIssueInfoBatch(reqVo.OrgId, reqVo.IssueIds)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, projectvo.GetLcIssueInfoBatchRespVo{
		Err:  vo.NewErr(err),
		Data: res,
	})
}

func IssueListStat(c *gin.Context) {
	var reqVo projectvo.IssueListStatReq
	if err := c.ShouldBindJSON(&reqVo); err != nil {
		// Replaced: c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		c.JSON(http.StatusBadRequest, projectvo.IssueListStatResp{Err: vo.NewErr(errs.ReqParamsValidateError)})
		return
	}

	res, err := service.IssueListStat(reqVo.OrgId, reqVo.UserId, reqVo.Input.ProjectID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, projectvo.IssueListStatResp{Err: vo.NewErr(err), Data: res})
}

func IssueListSimpleByDataIds(c *gin.Context) {
	var reqVo projectvo.GetIssueListSimpleByDataIdsReqVo
	if err := c.ShouldBindJSON(&reqVo); err != nil {
		// Replaced: c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		c.JSON(http.StatusBadRequest, projectvo.SimpleIssueListRespVo{Err: vo.NewErr(errs.ReqParamsValidateError)})
		return
	}

	res, err := service.IssueListSimpleByDataIds(reqVo.OrgId, reqVo.UserId, reqVo.Input.AppId, reqVo.Input.DataIds)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, projectvo.SimpleIssueListRespVo{
		Err:  vo.NewErr(err),
		Data: res,
	})
}
