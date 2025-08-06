package api

import (
	"net/http"
	"github.com/gin-gonic/gin"
	msgPb "gitea.bjx.cloud/LessCode/interface/golang/msg/v1"
	"github.com/star-table/startable-server/app/facade/common"
	"github.com/star-table/startable-server/common/model/vo"
	"github.com/star-table/startable-server/common/model/vo/commonvo"
)

func ReportAppEvent(c *gin.Context) {
	var req commonvo.ReportAppEventReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	common.ReportAppEvent(msgPb.EventType(req.EventType), req.TraceId, req.AppEvent)
	response := &vo.CommonRespVo{}
	c.JSON(http.StatusOK, response)
}

func ReportTableEvent(c *gin.Context) {
	var req commonvo.ReportTableEventReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	common.ReportTableEvent(msgPb.EventType(req.EventType), req.TraceId, req.TableEvent)
	response := &vo.CommonRespVo{}
	c.JSON(http.StatusOK, response)
}
