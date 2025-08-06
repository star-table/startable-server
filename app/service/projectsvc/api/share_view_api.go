package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/star-table/startable-server/app/service/projectsvc/service"
	"github.com/star-table/startable-server/common/model/vo/projectvo"
)

// GetShareViewInfo 获取共享视图信息
func GetShareViewInfo(c *gin.Context) {
	var req projectvo.GetShareViewInfoReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	info, err := service.GetShareViewInfo(&req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, &projectvo.GetShareViewInfoResp{
		Data: info,
	})
}

// GetShareViewInfoByKey 通过Key获取共享视图信息
func GetShareViewInfoByKey(c *gin.Context) {
	var req projectvo.GetShareViewInfoByKeyReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	info, err := service.GetShareViewInfoByKey(&req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, &projectvo.GetShareViewInfoResp{
		Data: info,
	})
}

// CreateShareView 创建共享视图
func CreateShareView(c *gin.Context) {
	var req projectvo.CreateShareViewReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	info, err := service.CreateShareView(&req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, &projectvo.GetShareViewInfoResp{
		Data: info,
	})
}

// DeleteShareView 删除共享视图
func DeleteShareView(c *gin.Context) {
	var req projectvo.DeleteShareViewReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := service.DeleteShareView(&req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true})
}

// ResetShareKey 重置共享密钥
func ResetShareKey(c *gin.Context) {
	var req projectvo.ResetShareKeyReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	info, err := service.ResetShareKey(&req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, &projectvo.GetShareViewInfoResp{
		Data: info,
	})
}

// UpdateSharePassword 更新共享密码
func UpdateSharePassword(c *gin.Context) {
	var req projectvo.UpdateSharePasswordReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := service.UpdateSharePassword(&req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true})
}

// UpdateShareConfig 更新共享配置
func UpdateShareConfig(c *gin.Context) {
	var req projectvo.UpdateShareConfigReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := service.UpdateShareConfig(&req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true})
}

// CheckShareViewPassword 检查共享视图密码
func CheckShareViewPassword(c *gin.Context) {
	var req projectvo.CheckShareViewPasswordReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	info, err := service.CheckShareViewPassword(req.Input)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, &projectvo.CheckShareViewPasswordResp{
		Data: info,
	})
}
