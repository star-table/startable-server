package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/star-table/startable-server/app/service/projectsvc/service"
	"github.com/star-table/startable-server/common/model/vo/projectvo"
)

// SaveMenu 保存菜单配置
func SaveMenu(c *gin.Context) {
	var req projectvo.SaveMenuReqVo
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	res, err := service.SaveMenu(req.OrgId, req.Input)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, projectvo.SaveMenuRespVo{
		Data: res,
	})
}

// GetMenu 获取菜单配置
func GetMenu(c *gin.Context) {
	var req projectvo.GetMenuReqVo
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	res, err := service.GetMenu(req.OrgId, req.AppId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, projectvo.GetMenuRespVo{
		Data: res,
	})
}
