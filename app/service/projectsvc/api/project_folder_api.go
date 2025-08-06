package api

import (
	"net/http"
	"github.com/gin-gonic/gin"
	"github.com/star-table/startable-server/app/service/projectsvc/service"
	"github.com/star-table/startable-server/common/model/vo"
	"github.com/star-table/startable-server/common/model/vo/projectvo"
)

func CreateProjectFolder(c *gin.Context) {
	var reqVo projectvo.CreateProjectFolderReqVo
	if err := c.ShouldBindJSON(&reqVo); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	res, err := service.CreateProjectFolder(reqVo.OrgId, reqVo.UserId, reqVo.Input)
	response := vo.CommonRespVo{Err: vo.NewErr(err), Void: res}
	c.JSON(http.StatusOK, response)
}

func UpdateProjectFolder(c *gin.Context) {
	var reqVo projectvo.UpdateProjectFolderReqVo
	if err := c.ShouldBindJSON(&reqVo); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	res, err := service.UpdateProjectFolder(reqVo.OrgId, reqVo.UserId, reqVo.Input)
	response := projectvo.UpdateProjectFolderRespVo{Err: vo.NewErr(err), Output: res}
	c.JSON(http.StatusOK, response)
}

func DeleteProjectFolder(c *gin.Context) {
	var reqVo projectvo.DeleteProjectFolerReqVo
	if err := c.ShouldBindJSON(&reqVo); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	res, err := service.DeleteProjectFolder(reqVo.OrgId, reqVo.UserId, reqVo.Input)
	response := projectvo.DeleteProjectFolerRespVo{Err: vo.NewErr(err), Output: res}
	c.JSON(http.StatusOK, response)
}

func GetProjectFolder(c *gin.Context) {
	var reqVo projectvo.GetProjectFolderReqVo
	if err := c.ShouldBindJSON(&reqVo); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	res, err := service.GetProjectFolder(reqVo.OrgId, reqVo.UserId, reqVo.Page, reqVo.Size, reqVo.Input)
	response := projectvo.GetProjectFolderRespVo{Err: vo.NewErr(err), Output: res}
	c.JSON(http.StatusOK, response)
}
