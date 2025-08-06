package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/star-table/startable-server/common/core/errs"
	"github.com/star-table/startable-server/common/core/logger"
	"github.com/star-table/startable-server/common/model/vo"
	"github.com/star-table/startable-server/common/model/vo/rolevo"
	orgsvcService "github.com/star-table/startable-server/app/service/orgsvc/service"
)

func GetOrgRoleUser(c *gin.Context) {
	var req rolevo.GetOrgRoleUserReqVo
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusOK, rolevo.GetOrgRoleUserRespVo{Err: vo.NewErr(errs.BuildSystemErrorInfo(errs.ReqParamsValidateError, err))})
		return
	}

	res, err := orgsvcService.GetOrgRoleUser(req.OrgId, req.ProjectId)
	if err != nil {
		logger.Error("GetOrgRoleUser error", logger.ErrorField(err))
	}
	c.JSON(http.StatusOK, rolevo.GetOrgRoleUserRespVo{Err: vo.NewErr(err), Data: res})
}

func GetOrgAdminUser(c *gin.Context) {
	var req rolevo.GetOrgAdminUserReqVo
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusOK, rolevo.GetOrgAdminUserRespVo{Err: vo.NewErr(errs.BuildSystemErrorInfo(errs.ReqParamsValidateError, err))})
		return
	}

	res, err := orgsvcService.GetOrgAdminUser(req.OrgId)
	if err != nil {
		logger.Error("GetOrgAdminUser error", logger.ErrorField(err))
	}
	c.JSON(http.StatusOK, rolevo.GetOrgAdminUserRespVo{Err: vo.NewErr(err), Data: res})
}

func GetOrgAdminUserBatch(c *gin.Context) {
	var req rolevo.GetOrgAdminUserBatchReqVo
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusOK, rolevo.GetOrgAdminUserBatchRespVo{Err: vo.NewErr(errs.BuildSystemErrorInfo(errs.ReqParamsValidateError, err))})
		return
	}

	res, err := orgsvcService.GetOrgAdminUserBatch(req.Input.OrgIds)
	if err != nil {
		logger.Error("GetOrgAdminUserBatch error", logger.ErrorField(err))
	}
	c.JSON(http.StatusOK, rolevo.GetOrgAdminUserBatchRespVo{Err: vo.NewErr(err), Data: res})
}

func UpdateUserOrgRole(c *gin.Context) {
	var reqVo rolevo.UpdateUserOrgRoleReqVo
	if err := c.ShouldBindJSON(&reqVo); err != nil {
		c.JSON(http.StatusOK, vo.CommonRespVo{Err: vo.NewErr(errs.BuildSystemErrorInfo(errs.ReqParamsValidateError, err))})
		return
	}

	res, err := orgsvcService.UpdateUserOrgRole(reqVo)
	if err != nil {
		logger.Error("UpdateUserOrgRole error", logger.ErrorField(err))
	}
	c.JSON(http.StatusOK, vo.CommonRespVo{Err: vo.NewErr(err), Void: res})
}

func GetOrgRoleList(c *gin.Context) {
	var reqVo rolevo.GetOrgRoleListReqVo
	if err := c.ShouldBindJSON(&reqVo); err != nil {
		c.JSON(http.StatusOK, rolevo.GetOrgRoleListRespVo{Err: vo.NewErr(errs.BuildSystemErrorInfo(errs.ReqParamsValidateError, err))})
		return
	}

	res, err := orgsvcService.GetOrgRoleList(reqVo.OrgId)
	if err != nil {
		logger.Error("GetOrgRoleList error", logger.ErrorField(err))
	}
	c.JSON(http.StatusOK, rolevo.GetOrgRoleListRespVo{Err: vo.NewErr(err), Data: res})
}

func GetUserAdminFlag(c *gin.Context) {
	var reqVo rolevo.GetUserAdminFlagReqVo
	if err := c.ShouldBindJSON(&reqVo); err != nil {
		c.JSON(http.StatusOK, rolevo.GetUserAdminFlagRespVo{Err: vo.NewErr(errs.BuildSystemErrorInfo(errs.ReqParamsValidateError, err))})
		return
	}

	res, err := orgsvcService.GetUserAdminFlag(reqVo.OrgId, reqVo.UserId)
	if err != nil {
		logger.Error("GetUserAdminFlag error", logger.ErrorField(err))
	}
	c.JSON(http.StatusOK, rolevo.GetUserAdminFlagRespVo{Err: vo.NewErr(err), Data: res})
}

func GetProjectRoleList(c *gin.Context) {
	var reqVo rolevo.GetProjectRoleListReqVo
	if err := c.ShouldBindJSON(&reqVo); err != nil {
		c.JSON(http.StatusOK, rolevo.GetOrgRoleListRespVo{Err: vo.NewErr(errs.BuildSystemErrorInfo(errs.ReqParamsValidateError, err))})
		return
	}

	res, err := orgsvcService.GetProjectRoleList(reqVo.OrgId, reqVo.ProjectId)
	if err != nil {
		logger.Error("GetProjectRoleList error", logger.ErrorField(err))
	}
	c.JSON(http.StatusOK, rolevo.GetOrgRoleListRespVo{Err: vo.NewErr(err), Data: res})
}

func UpdateOrgOwner(c *gin.Context) {
	var req rolevo.UpdateOrgOwnerReqVo
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusOK, vo.CommonRespVo{Err: vo.NewErr(errs.BuildSystemErrorInfo(errs.ReqParamsValidateError, err))})
		return
	}

	res, err := orgsvcService.UpdateOrgAdmin(req.UserId, req.OrgId, req.OldOwnerId, req.NewOwnerId)
	if err != nil {
		logger.Error("UpdateOrgOwner error", logger.ErrorField(err))
	}
	c.JSON(http.StatusOK, vo.CommonRespVo{Err: vo.NewErr(err), Void: res})
}

func UpdateUserOrgRoleBatch(c *gin.Context) {
	var reqVo rolevo.UpdateUserOrgRoleBatchReqVo
	if err := c.ShouldBindJSON(&reqVo); err != nil {
		c.JSON(http.StatusOK, vo.CommonRespVo{Err: vo.NewErr(errs.BuildSystemErrorInfo(errs.ReqParamsValidateError, err))})
		return
	}

	res, err := orgsvcService.UpdateUserOrgRoleBatch(reqVo)
	if err != nil {
		logger.Error("UpdateUserOrgRoleBatch error", logger.ErrorField(err))
	}
	c.JSON(http.StatusOK, vo.CommonRespVo{Err: vo.NewErr(err), Void: res})
}
