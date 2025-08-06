package api

import (
	"net/http"
	"github.com/gin-gonic/gin"
	commonService "github.com/star-table/startable-server/app/service/commonsvc/service"
	"github.com/star-table/startable-server/common/model/vo"
	"github.com/star-table/startable-server/common/model/vo/commonvo"
)

// AreaLinkageList 获取地区联动列表
func AreaLinkageList(c *gin.Context) {
	var req commonvo.AreaLinkageListReqVo
	if err := c.ShouldBindJSON(&req); err != nil {
		response := commonvo.AreaLinkageListRespVo{Err: vo.NewErr(err)}
		c.JSON(http.StatusOK, response)
		return
	}
	
	res, err := commonService.AreaLinkageList(req.Input)
	response := commonvo.AreaLinkageListRespVo{Err: vo.NewErr(err), AreaLinkageListResp: res}
	c.JSON(http.StatusOK, response)
}

// AreaInfo 获取地区信息
func AreaInfo(c *gin.Context) {
	var req commonvo.AreaInfoReqVo
	if err := c.ShouldBindJSON(&req); err != nil {
		response := commonvo.AreaInfoRespVo{Err: vo.NewErr(err)}
		c.JSON(http.StatusOK, response)
		return
	}
	
	res, err := commonService.OrgAreaInfo(req)
	response := commonvo.AreaInfoRespVo{Err: vo.NewErr(err), AreaInfoResp: res}
	c.JSON(http.StatusOK, response)
}
