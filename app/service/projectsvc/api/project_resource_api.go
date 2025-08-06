package api

import (
	"net/http"
	"github.com/gin-gonic/gin"
	"github.com/star-table/startable-server/app/service/projectsvc/service"
	"github.com/star-table/startable-server/common/model/vo"
	"github.com/star-table/startable-server/common/model/vo/projectvo"
)

func CreateProjectResource(c *gin.Context) {
	var reqVo projectvo.CreateProjectResourceReqVo
	if err := c.ShouldBindJSON(&reqVo); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	res, err := service.CreateProjectResource(reqVo.OrgId, reqVo.UserId, reqVo.Input)
	response := vo.CommonRespVo{Err: vo.NewErr(err), Void: res}
	c.JSON(http.StatusOK, response)
}

func UpdateProjectResourceName(c *gin.Context) {
	var reqVo projectvo.UpdateProjectResourceNameReqVo
	if err := c.ShouldBindJSON(&reqVo); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	res, err := service.UpdateProjectResourceName(reqVo.OrgId, reqVo.UserId, reqVo.Input)
	response := vo.CommonRespVo{Err: vo.NewErr(err), Void: res}
	c.JSON(http.StatusOK, response)
}

func UpdateProjectFileResource(c *gin.Context) {
	var reqVo projectvo.UpdateProjectFileResourceReqVo
	if err := c.ShouldBindJSON(&reqVo); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	res, err := service.UpdateProjectFileResource(reqVo.OrgId, reqVo.UserId, reqVo.Input)
	response := vo.CommonRespVo{Err: vo.NewErr(err), Void: res}
	c.JSON(http.StatusOK, response)
}

func UpdateProjectResourceFolder(c *gin.Context) {
	var reqVo projectvo.UpdateProjectResourceFolderReqVo
	if err := c.ShouldBindJSON(&reqVo); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	res, err := service.UpdateProjectResourceFolder(reqVo.OrgId, reqVo.UserId, reqVo.Input)
	response := projectvo.UpdateProjectResourceFolderRespVo{Err: vo.NewErr(err), Output: res}
	c.JSON(http.StatusOK, response)
}

func DeleteProjectResource(c *gin.Context) {
	var reqVo projectvo.DeleteProjectResourceReqVo
	if err := c.ShouldBindJSON(&reqVo); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	res, err := service.DeleteProjectResource(reqVo.OrgId, reqVo.UserId, reqVo.Input)
	response := projectvo.DeleteProjectResourceRespVo{Err: vo.NewErr(err), Output: res}
	c.JSON(http.StatusOK, response)
}

func GetProjectResource(c *gin.Context) {
	var reqVo projectvo.GetProjectResourceReqVo
	if err := c.ShouldBindJSON(&reqVo); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	res, err := service.GetProjectResource(reqVo.OrgId, reqVo.UserId, reqVo.Page, reqVo.Size, reqVo.Input)
	response := projectvo.GetProjectResourceResVo{Err: vo.NewErr(err), Output: res}
	c.JSON(http.StatusOK, response)
}

func GetProjectResourceInfo(c *gin.Context) {
	var reqVo projectvo.GetProjectResourceInfoReqVo
	if err := c.ShouldBindJSON(&reqVo); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	res, err := service.GetProjectResourceInfo(reqVo.OrgId, reqVo.UserId, reqVo.Input)
	response := projectvo.GetProjectResourceInfoRespVo{Err: vo.NewErr(err), Output: res}
	c.JSON(http.StatusOK, response)
}
