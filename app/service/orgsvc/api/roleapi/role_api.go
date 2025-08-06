package api

import (
	"net/http"
	"github.com/gin-gonic/gin"
	orgsvcService "github.com/star-table/startable-server/app/service/orgsvc/service"
	"github.com/star-table/startable-server/common/model/vo"
	"github.com/star-table/startable-server/common/model/vo/rolevo"
	"github.com/star-table/startable-server/common/core/errs"
	"github.com/star-table/startable-server/common/core/logger"
)

func CreateRole(c *gin.Context) {
	var req rolevo.CreateOrgReqVo
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Error("CreateRole bind request failed", err)
		response := vo.CommonRespVo{Err: vo.NewErr(errs.BuildSystemErrorInfo(errs.ReqParamsValidateError, err))}
		c.JSON(http.StatusOK, response)
		return
	}

	res, err := orgsvcService.CreateRole(req.OrgId, req.UserId, req.Input)
	response := vo.CommonRespVo{Err: vo.NewErr(err), Void: res}
	c.JSON(http.StatusOK, response)
}

func DelRole(c *gin.Context) {
	var req rolevo.DelRoleReqVo
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Error("DelRole bind request failed", err)
		response := vo.CommonRespVo{Err: vo.NewErr(errs.BuildSystemErrorInfo(errs.ReqParamsValidateError, err))}
		c.JSON(http.StatusOK, response)
		return
	}

	err := orgsvcService.DelRole(req.OrgId, req.UserId, req.Input)
	response := vo.CommonRespVo{Void: &vo.Void{ID: 0}, Err: vo.NewErr(err)}
	c.JSON(http.StatusOK, response)
}

func UpdateRole(c *gin.Context) {
	var req rolevo.UpdateRoleReqVo
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Error("UpdateRole bind request failed", err)
		response := vo.CommonRespVo{Err: vo.NewErr(errs.BuildSystemErrorInfo(errs.ReqParamsValidateError, err))}
		c.JSON(http.StatusOK, response)
		return
	}

	id, err := orgsvcService.UpdateRole(req.OrgId, req.UserId, req.Input)
	response := vo.CommonRespVo{Void: &vo.Void{ID: id}, Err: vo.NewErr(err)}
	c.JSON(http.StatusOK, response)
}

func ClearUserRoleList(c *gin.Context) {
	var req rolevo.ClearUserRoleReqVo
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Error("ClearUserRoleList bind request failed", err)
		response := vo.CommonRespVo{Err: vo.NewErr(errs.BuildSystemErrorInfo(errs.ReqParamsValidateError, err))}
		c.JSON(http.StatusOK, response)
		return
	}

	err := orgsvcService.ClearUserRoleList(req.OrgId, req.UserIds, req.ProjectId)
	response := vo.CommonRespVo{Void: &vo.Void{ID: 0}, Err: vo.NewErr(err)}
	c.JSON(http.StatusOK, response)
}
