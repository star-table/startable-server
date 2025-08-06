package api

import (
	"net/http"
	"github.com/gin-gonic/gin"
	idsvcService "github.com/star-table/startable-server/app/service/idsvc/service"
	"github.com/star-table/startable-server/common/model/vo"
	"github.com/star-table/startable-server/common/model/vo/idvo"
)

// ApplyPrimaryId 申请主键ID
func ApplyPrimaryId(c *gin.Context) {
	var req idvo.ApplyPrimaryIdReqVo
	if err := c.ShouldBindJSON(&req); err != nil {
		response := idvo.ApplyPrimaryIdRespVo{Err: vo.NewErr(err)}
		c.JSON(http.StatusOK, response)
		return
	}
	
	id, err := idsvcService.ApplyPrimaryId(req.Code)
	response := idvo.ApplyPrimaryIdRespVo{Id: id, Err: vo.NewErr(err)}
	c.JSON(http.StatusOK, response)
}

// ApplyCode 申请编码
func ApplyCode(c *gin.Context) {
	var req idvo.ApplyCodeReqVo
	if err := c.ShouldBindJSON(&req); err != nil {
		response := idvo.ApplyCodeRespVo{Err: vo.NewErr(err)}
		c.JSON(http.StatusOK, response)
		return
	}
	
	code, err := idsvcService.ApplyCode(req.OrgId, req.Code, req.PreCode)
	response := idvo.ApplyCodeRespVo{Code: code, Err: vo.NewErr(err)}
	c.JSON(http.StatusOK, response)
}

// ApplyMultipleId 申请多个ID
func ApplyMultipleId(c *gin.Context) {
	var req idvo.ApplyMultipleIdReqVo
	if err := c.ShouldBindJSON(&req); err != nil {
		response := idvo.ApplyMultipleIdRespVo{Err: vo.NewErr(err)}
		c.JSON(http.StatusOK, response)
		return
	}
	
	idCodes, err := idsvcService.ApplyMultipleIdExtra(req.OrgId, req.Code, req.PreCode, req.Count)
	response := idvo.ApplyMultipleIdRespVo{IdCodes: idCodes, Err: vo.NewErr(err)}
	c.JSON(http.StatusOK, response)
}

// ApplyMultiplePrimaryId 申请多个主键ID
func ApplyMultiplePrimaryId(c *gin.Context) {
	var req idvo.ApplyMultiplePrimaryIdReqVo
	if err := c.ShouldBindJSON(&req); err != nil {
		response := idvo.ApplyMultipleIdRespVo{Err: vo.NewErr(err)}
		c.JSON(http.StatusOK, response)
		return
	}
	
	idCodes, err := idsvcService.ApplyMultiplePrimaryId(req.Code, req.Count)
	response := idvo.ApplyMultipleIdRespVo{IdCodes: idCodes, Err: vo.NewErr(err)}
	c.JSON(http.StatusOK, response)
}

// ApplyMultiplePrimaryIdByCodes 根据编码申请多个主键ID
func ApplyMultiplePrimaryIdByCodes(c *gin.Context) {
	var req idvo.ApplyMultiplePrimaryIdByCodesReqVo
	if err := c.ShouldBindJSON(&req); err != nil {
		response := idvo.ApplyMultipleIdRespByCodesVo{Err: vo.NewErr(err)}
		c.JSON(http.StatusOK, response)
		return
	}
	
	idCodes, err := idsvcService.ApplyMultiplePrimaryIdByCodes(req)
	response := idvo.ApplyMultipleIdRespByCodesVo{CodesIds: idCodes, Err: vo.NewErr(err)}
	c.JSON(http.StatusOK, response)
}
