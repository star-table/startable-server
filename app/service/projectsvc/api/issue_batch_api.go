package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/star-table/startable-server/app/service/projectsvc/service"
	"github.com/star-table/startable-server/common/core/consts"
	"github.com/star-table/startable-server/common/core/errs"
	"github.com/star-table/startable-server/common/model/vo"
	"github.com/star-table/startable-server/common/model/vo/projectvo"
)

func BatchCreateIssue(c *gin.Context) {
	var reqVo projectvo.BatchCreateIssueReqVo
	if err := c.ShouldBindJSON(&reqVo); err != nil {
		// Replaced: c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		c.JSON(http.StatusBadRequest, projectvo.BatchCreateIssueRespVo{Err: vo.NewErr(errs.ReqParamsValidateError)})
		return
	}

	data, userDepts, relateData, err := service.BatchCreateIssue(&reqVo, false, &projectvo.TriggerBy{
		TriggerBy: consts.TriggerByNormal,
	}, "", 0)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, projectvo.BatchCreateIssueRespVo{Data: data, UserDepts: userDepts, RelateData: relateData, Err: vo.NewErr(err)})
}

func BatchUpdateIssue(c *gin.Context) {
	var reqVo projectvo.BatchUpdateIssueReqVo
	if err := c.ShouldBindJSON(&reqVo); err != nil {
		// Replaced: c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		c.JSON(http.StatusBadRequest, vo.VoidErr{Err: vo.NewErr(errs.ReqParamsValidateError)})
		return
	}

	err := service.BatchUpdateIssue(&projectvo.BatchUpdateIssueReqInnerVo{
		OrgId:        reqVo.OrgId,
		UserId:       reqVo.UserId,
		AppId:        reqVo.Input.AppId,
		ProjectId:    reqVo.Input.ProjectId,
		TableId:      reqVo.Input.TableId,
		Data:         reqVo.Input.Data,
		BeforeDataId: reqVo.Input.BeforeDataId,
		AfterDataId:  reqVo.Input.AfterDataId,
		TodoId:       reqVo.Input.TodoId,
		TodoOp:       reqVo.Input.TodoOp,
		TodoMsg:      reqVo.Input.TodoMsg,
	}, false, &projectvo.TriggerBy{
		TriggerBy: consts.TriggerByNormal,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, vo.VoidErr{Err: vo.NewErr(err)})
}

func BatchAuditIssue(c *gin.Context) {
	var reqVo projectvo.BatchAuditIssueReqVo
	if err := c.ShouldBindJSON(&reqVo); err != nil {
		// Replaced: c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		c.JSON(http.StatusBadRequest, vo.DataRespVo{Err: vo.NewErr(errs.ReqParamsValidateError)})
		return
	}

	issueIds, err := service.BatchAuditIssue(&reqVo)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, vo.DataRespVo{Data: issueIds, Err: vo.NewErr(err)})
}

func BatchUrgeIssue(c *gin.Context) {
	var reqVo projectvo.BatchUrgeIssueReqVo
	if err := c.ShouldBindJSON(&reqVo); err != nil {
		// Replaced: c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		c.JSON(http.StatusBadRequest, vo.DataRespVo{Err: vo.NewErr(errs.ReqParamsValidateError)})
		return
	}

	issueIds, err := service.BatchUrgeIssue(&reqVo)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, vo.DataRespVo{Data: issueIds, Err: vo.NewErr(err)})
}
