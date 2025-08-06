package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/star-table/startable-server/app/service/projectsvc/domain"
	"github.com/star-table/startable-server/common/core/errs"

	"github.com/star-table/startable-server/common/model/vo"
	"github.com/star-table/startable-server/common/model/vo/projectvo"
)

func GetIssueLinks(c *gin.Context) {
	req := projectvo.GetIssueLinksReqVo{}
	if err := c.ShouldBindJSON(&req); err != nil {
		// Replaced: c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		c.JSON(http.StatusOK, projectvo.GetIssueLinksRespVo{Err: vo.NewErr(errs.BuildSystemErrorInfo(errs.ReqParamsValidateError, err))})
		return
	}

	c.JSON(http.StatusOK, projectvo.GetIssueLinksRespVo{
		Err:  vo.NewErr(nil),
		Data: domain.GetIssueLinks(req.SourceChannel, req.OrgId, req.IssueId),
	})
}
