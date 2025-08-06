package api

import (
	"net/http"
	"github.com/gin-gonic/gin"
	"github.com/star-table/startable-server/app/service/projectsvc/service"
	"github.com/star-table/startable-server/common/model/vo"
	"github.com/star-table/startable-server/common/model/vo/projectvo"
)

func IterationStats(c *gin.Context) {
	var reqVo projectvo.IterationStatsReqVo
	if err := c.ShouldBindJSON(&reqVo); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	res, err := service.IterationStats(reqVo.OrgId, reqVo.Page, reqVo.Size, reqVo.Input)
	response := projectvo.IterationStatsRespVo{Err: vo.NewErr(err), IterationStats: res}
	c.JSON(http.StatusOK, response)
}
