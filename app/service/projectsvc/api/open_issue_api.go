package api

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/star-table/startable-server/app/facade/orgfacade"
	"github.com/star-table/startable-server/app/service/projectsvc/service"
	"github.com/star-table/startable-server/common/core/errs"
	"github.com/star-table/startable-server/common/core/util/format"
	"github.com/star-table/startable-server/common/core/util/json"
	"github.com/star-table/startable-server/common/model/vo"
	"github.com/star-table/startable-server/common/model/vo/projectvo"
)

func OpenLcHomeIssues(c *gin.Context) {
	req := projectvo.HomeIssuesReqVo{}
	if err := c.ShouldBindJSON(&req); err != nil {
		// Replaced: c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		c.String(http.StatusOK, json.ToJsonIgnoreError(&projectvo.LcHomeIssuesRespVo{Err: vo.NewErr(errs.BuildSystemErrorInfo(errs.ReqParamsValidateError, err)), Data: "{}"}))
		return
	}

	if req.Size > 30000 {
		req.Size = 30000
	}
	res, err := service.LcHomeIssuesForProject(req.OrgId, req.UserId, req.Page, req.Size, req.Input, true)
	if err != nil {
		// Replaced: c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()}); return
		c.String(http.StatusOK, json.ToJsonIgnoreError(&projectvo.LcHomeIssuesRespVo{Err: vo.NewErr(err), Data: "{}"}))
		return
	}

	c.String(http.StatusOK, res.Data)
}

func OpenLcHomeIssuesForAll(c *gin.Context) {
	req := projectvo.HomeIssuesReqVo{}
	if err := c.ShouldBindJSON(&req); err != nil {
		// Replaced: c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		c.String(http.StatusOK, json.ToJsonIgnoreError(&projectvo.LcHomeIssuesRespVo{Err: vo.NewErr(errs.BuildSystemErrorInfo(errs.ReqParamsValidateError, err)), Data: "{}"}))
		return
	}

	if req.Size > 30000 {
		req.Size = 30000
	}
	res, err := service.LcHomeIssuesForAll(req.OrgId, req.UserId, req.Page, req.Size, req.Input, true)
	if err != nil {
		// Replaced: c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()}); return
		c.String(http.StatusOK, json.ToJsonIgnoreError(&projectvo.LcHomeIssuesRespVo{Err: vo.NewErr(err), Data: "{}"}))
		return
	}

	c.String(http.StatusOK, res.Data)
}

func OpenIssueInfo(c *gin.Context) {
	req := projectvo.IssueInfoReqVo{}
	if err := c.ShouldBindJSON(&req); err != nil {
		// Replaced: c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		c.JSON(http.StatusOK, &projectvo.IssueDetailRespVo{Err: vo.NewErr(errs.BuildSystemErrorInfo(errs.ReqParamsValidateError, err))})
		return
	}

	data, err := service.LcHomeIssuesForIssue(req.OrgId, req.UserId, -1, -1, req.IssueID, -1, true)
	if err != nil {
		// Replaced: c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()}); return
	}
	c.JSON(http.StatusOK, &projectvo.IssueDetailRespVo{Err: vo.NewErr(err), Data: data})
}

func OpenIssueStatusTypeStat(c *gin.Context) {
	req := projectvo.IssueStatusTypeStatReqVo{}
	if err := c.ShouldBindJSON(&req); err != nil {
		// Replaced: c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		c.JSON(http.StatusOK, projectvo.IssueStatusTypeStatRespVo{Err: vo.NewErr(errs.BuildSystemErrorInfo(errs.ReqParamsValidateError, err))})
		return
	}

	res, err := service.IssueStatusTypeStat(req.OrgId, req.UserId, req.Input)
	if err != nil {
		// Replaced: c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()}); return
	}
	c.JSON(http.StatusOK, projectvo.IssueStatusTypeStatRespVo{Err: vo.NewErr(err), IssueStatusTypeStat: res})
}

func OpenCreateIssue(c *gin.Context) {
	req := projectvo.CreateIssueReqVo{}
	if err := c.ShouldBindJSON(&req); err != nil {
		// Replaced: c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		c.JSON(http.StatusOK, projectvo.IssueRespVo{Err: vo.NewErr(errs.BuildSystemErrorInfo(errs.ReqParamsValidateError, err))})
		return
	}

	c.JSON(http.StatusOK, projectvo.IssueRespVo{})
}

func OpenUpdateIssue(c *gin.Context) {
	var reqVo projectvo.OpenUpdateIssueReqVo
	if err := c.ShouldBindJSON(&reqVo); err != nil {
		// Replaced: c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		c.JSON(http.StatusBadRequest, projectvo.UpdateIssueRespVo{Err: vo.NewErr(errs.ReqParamsValidateError)})
		return
	}

	result := projectvo.UpdateIssueRespVo{Err: vo.NewErr(nil), UpdateIssue: &vo.UpdateIssueResp{}}
	c.JSON(http.StatusOK, result)
}

func OpenDeleteIssue(c *gin.Context) {
	var reqVo projectvo.OpenDeleteIssueReqVo
	if err := c.ShouldBindJSON(&reqVo); err != nil {
		// Replaced: c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		c.JSON(http.StatusBadRequest, projectvo.DeleteIssueRespVo{Err: vo.NewErr(errs.ReqParamsValidateError)})
		return
	}

	// 校验操作人是否存在
	_, userErr := orgfacade.GetBaseUserInfoRelaxed(reqVo.OrgId, reqVo.UserId)
	if userErr != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": userErr.Error()})
		return
	}

	judgeErr := service.JudgeProjectFilingByIssueId(reqVo.OrgId, reqVo.Data.IssueId)
	if judgeErr != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": judgeErr.Error()})
		return
	}
	_, err := service.DeleteIssueWithoutAuth(projectvo.DeleteIssueReqVo{
		Input: vo.DeleteIssueReq{
			ID: reqVo.Data.IssueId,
		},
		OrgId:  reqVo.OrgId,
		UserId: reqVo.UserId,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, projectvo.DeleteIssueRespVo{})
}

func OpenIssueList(c *gin.Context) {
	var reqVo projectvo.OpenIssueListReqVo
	if err := c.ShouldBindJSON(&reqVo); err != nil {
		// Replaced: c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		c.JSON(http.StatusBadRequest, projectvo.LcHomeIssuesRespVo{Err: vo.NewErr(errs.ReqParamsValidateError), Data: "{}"})
		return
	}

	res, err := service.LcHomeIssuesForProject(reqVo.OrgId, reqVo.UserId, reqVo.Page, reqVo.Size, reqVo.Data, false)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, projectvo.LcHomeIssuesRespVo{Err: vo.NewErr(nil), Data: json.ToJsonIgnoreError(map[string]interface{}{
		"code":    0,
		"success": true,
		"data":    res,
	})})
}

func OpenCreateIssueComment(c *gin.Context) {
	var reqVo projectvo.CreateIssueCommentReqVo
	if err := c.ShouldBindJSON(&reqVo); err != nil {
		// Replaced: c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		c.JSON(http.StatusBadRequest, vo.CommonRespVo{Err: vo.NewErr(errs.ReqParamsValidateError), Void: nil})
		return
	}

	reqVo.Input.Comment = strings.TrimSpace(reqVo.Input.Comment)
	isCommentRight := format.VerifyIssueCommenFormat(reqVo.Input.Comment)
	if !isCommentRight {
		c.JSON(http.StatusBadRequest, gin.H{"error": "评论格式验证失败"})
		return
	}

	res, err := service.CreateIssueComment(reqVo.OrgId, reqVo.UserId, reqVo.Input)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, vo.CommonRespVo{Err: vo.NewErr(err), Void: res})
}
