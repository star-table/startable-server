package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/star-table/startable-server/app/service/projectsvc/service"
	"github.com/star-table/startable-server/common/core/util/copyer"
	"github.com/star-table/startable-server/common/model/vo"
	"github.com/star-table/startable-server/common/model/vo/projectvo"
)

// PayLimitNum 获取付费限制数量
func PayLimitNum(c *gin.Context) {
	var req projectvo.PayLimitNumReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	res, err := service.PayLimitNum(req.OrgId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	respData := vo.PayLimitNumResp{}
	copyer.Copy(res, &respData)

	c.JSON(http.StatusOK, projectvo.PayLimitNumResp{
		Data: &respData,
	})
}

// PayLimitNumForRest 获取付费限制数量（REST版本）
func PayLimitNumForRest(c *gin.Context) {
	var req projectvo.PayLimitNumReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	res, err := service.PayLimitNum(req.OrgId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, projectvo.PayLimitNumForRestResp{
		Data: res,
	})
}

// GetProjectStatistics 获取项目统计信息
func GetProjectStatistics(c *gin.Context) {
	var req projectvo.GetProjectStatisticsReqVo
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	res, err := service.GetProjectStatistics(req.OrgId, req.UserId, req.Input)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, projectvo.GetProjectStatisticsResp{
		Data: res,
	})
}
