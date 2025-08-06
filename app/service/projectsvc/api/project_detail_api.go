package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/star-table/startable-server/app/service/projectsvc/service"
	"github.com/star-table/startable-server/common/model/vo/projectvo"
)

// ProjectDetail 获取项目详情
func ProjectDetail(c *gin.Context) {
	var reqVo projectvo.ProjectDetailReqVo
	if err := c.ShouldBindJSON(&reqVo); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	res, err := service.ProjectDetail(reqVo.OrgId, reqVo.ProjectId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, projectvo.ProjectDetailRespVo{
		ProjectDetail: res,
	})
}

// CreateProjectDetail 创建项目详情
func CreateProjectDetail(c *gin.Context) {
	var reqVo projectvo.CreateProjectDetailReqVo
	if err := c.ShouldBindJSON(&reqVo); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	res, err := service.CreateProjectDetail(reqVo.UserId, reqVo.Input)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": res})
}

// UpdateProjectDetail 更新项目详情
func UpdateProjectDetail(c *gin.Context) {
	var reqVo projectvo.UpdateProjectDetailReqVo
	if err := c.ShouldBindJSON(&reqVo); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	res, err := service.UpdateProjectDetail(reqVo.OrgId, reqVo.UserId, reqVo.Input)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": res})
}

// DeleteProjectDetail 删除项目详情
func DeleteProjectDetail(c *gin.Context) {
	var reqVo projectvo.DeleteProjectDetailReqVo
	if err := c.ShouldBindJSON(&reqVo); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	res, err := service.DeleteProjectDetail(reqVo.UserId, reqVo.Input)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": res})
}
