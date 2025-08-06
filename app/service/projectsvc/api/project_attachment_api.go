package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/star-table/startable-server/app/service/projectsvc/service"
	"github.com/star-table/startable-server/common/model/vo/projectvo"
)

// DeleteProjectAttachment 删除项目附件
func DeleteProjectAttachment(c *gin.Context) {
	var reqVo projectvo.DeleteProjectAttachmentReqVo
	if err := c.ShouldBindJSON(&reqVo); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	res, err := service.DeleteProjectAttachment(reqVo.OrgId, reqVo.UserId, reqVo.Input)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, projectvo.DeleteProjectAttachmentRespVo{Output: res})
}

// GetProjectAttachment 获取项目附件列表
func GetProjectAttachment(c *gin.Context) {
	var reqVo projectvo.GetProjectAttachmentReqVo
	if err := c.ShouldBindJSON(&reqVo); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if reqVo.Size == 0 {
		reqVo.Size = 20
	}

	res, err := service.GetProjectAttachment(reqVo.OrgId, reqVo.UserId, reqVo.Page, reqVo.Size, reqVo.Input)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, projectvo.GetProjectAttachmentRespVo{Output: res})
}

// GetProjectAttachmentInfo 获取项目附件信息
func GetProjectAttachmentInfo(c *gin.Context) {
	var reqVo projectvo.GetProjectAttachmentInfoReqVo
	if err := c.ShouldBindJSON(&reqVo); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	res, err := service.GetProjectAttachmentInfo(reqVo.OrgId, reqVo.UserId, reqVo.Input)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, projectvo.GetProjectAttachmentInfoRespVo{Output: res})
}
