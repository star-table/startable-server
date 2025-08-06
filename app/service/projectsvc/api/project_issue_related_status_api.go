package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/star-table/startable-server/app/service/projectsvc/service"
	"github.com/star-table/startable-server/common/model/vo/projectvo"
)

// ProjectIssueRelatedStatus 获取项目任务相关状态
func ProjectIssueRelatedStatus(c *gin.Context) {
	var reqVo projectvo.ProjectIssueRelatedStatusReqVo
	if err := c.ShouldBindJSON(&reqVo); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	res, err := service.ProjectIssueRelatedStatus(reqVo.OrgId, reqVo.Input)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, projectvo.ProjectIssueRelatedStatusRespVo{
		ProjectIssueRelatedStatus: res,
	})
}
