package api

import (
	"net/http"
	"github.com/gin-gonic/gin"
	orgsvcService "github.com/star-table/startable-server/app/service/orgsvc/service"
	"github.com/star-table/startable-server/common/core/errs"
	"github.com/star-table/startable-server/common/model/vo"
	"github.com/star-table/startable-server/common/model/vo/rolevo"
)

// Authenticate 权限认证
func Authenticate(c *gin.Context) {
	var req rolevo.AuthenticateReqVo
	if err := c.ShouldBindJSON(&req); err != nil {
		response := vo.VoidErr{Err: vo.NewErr(errs.BuildSystemErrorInfo(errs.ParamError, err))}
		c.JSON(http.StatusOK, response)
		return
	}
	
	err := orgsvcService.Authenticate(req.OrgId, req.UserId, req.AuthInfoReqVo.ProjectAuthInfo, req.AuthInfoReqVo.IssueAuthInfo, req.Path, req.Operation, req.AuthInfoReqVo.Fields)
	if err != nil {
		response := vo.VoidErr{Err: vo.NewErr(errs.BuildSystemErrorInfo(errs.ServerError, err))}
		c.JSON(http.StatusOK, response)
		return
	}
	response := vo.VoidErr{Err: vo.NewErr(nil)}
	c.JSON(http.StatusOK, response)
}

// RoleUserRelation 角色用户关系
func RoleUserRelation(c *gin.Context) {
	var req rolevo.RoleUserRelationReqVo
	if err := c.ShouldBindJSON(&req); err != nil {
		response := vo.VoidErr{Err: vo.NewErr(errs.BuildSystemErrorInfo(errs.ParamError, err))}
		c.JSON(http.StatusOK, response)
		return
	}
	
	err := orgsvcService.RoleUserRelation(req.OrgId, req.UserId, req.RoleId)
	if err != nil {
		response := vo.VoidErr{Err: vo.NewErr(errs.BuildSystemErrorInfo(errs.ServerError, err))}
		c.JSON(http.StatusOK, response)
		return
	}
	response := vo.VoidErr{Err: vo.NewErr(nil)}
	c.JSON(http.StatusOK, response)
}

// RemoveRoleUserRelation 移除角色用户关系
func RemoveRoleUserRelation(c *gin.Context) {
	var req rolevo.RemoveRoleUserRelationReqVo
	if err := c.ShouldBindJSON(&req); err != nil {
		response := vo.VoidErr{Err: vo.NewErr(errs.BuildSystemErrorInfo(errs.ParamError, err))}
		c.JSON(http.StatusOK, response)
		return
	}
	
	err := orgsvcService.RemoveRoleUserRelation(req)
	if err != nil {
		response := vo.VoidErr{Err: vo.NewErr(errs.BuildSystemErrorInfo(errs.ServerError, err))}
		c.JSON(http.StatusOK, response)
		return
	}
	response := vo.VoidErr{Err: vo.NewErr(nil)}
	c.JSON(http.StatusOK, response)
}

// RoleInit 角色初始化
func RoleInit(c *gin.Context) {
	var req rolevo.RoleInitReqVo
	if err := c.ShouldBindJSON(&req); err != nil {
		response := rolevo.RoleInitRespVo{Err: vo.NewErr(errs.BuildSystemErrorInfo(errs.ParamError, err))}
		c.JSON(http.StatusOK, response)
		return
	}
	
	roleInitResp, err := orgsvcService.RoleInit(req.OrgId)
	if err != nil {
		response := rolevo.RoleInitRespVo{Err: vo.NewErr(errs.BuildSystemErrorInfo(errs.ServerError, err))}
		c.JSON(http.StatusOK, response)
		return
	}
	response := rolevo.RoleInitRespVo{
		RoleInitResp: roleInitResp,
		Err: vo.NewErr(nil),
	}
	c.JSON(http.StatusOK, response)
}

// ChangeDefaultRole 更改默认角色
func ChangeDefaultRole(c *gin.Context) {
	err := orgsvcService.ChangeDefaultRole()
	if err != nil {
		response := vo.CommonRespVo{Err: vo.NewErr(errs.BuildSystemErrorInfo(errs.ServerError, err))}
		c.JSON(http.StatusOK, response)
		return
	}
	response := vo.CommonRespVo{
		Err:  vo.NewErr(nil),
		Void: nil,
	}
	c.JSON(http.StatusOK, response)
}

// RemoveRoleDepartmentRelation 移除角色部门关系
func RemoveRoleDepartmentRelation(c *gin.Context) {
	var req rolevo.RemoveRoleDepartmentRelationReqVo
	if err := c.ShouldBindJSON(&req); err != nil {
		response := vo.VoidErr{Err: vo.NewErr(errs.BuildSystemErrorInfo(errs.ParamError, err))}
		c.JSON(http.StatusOK, response)
		return
	}
	
	err := orgsvcService.RemoveRoleDepartmentRelation(req)
	if err != nil {
		response := vo.VoidErr{Err: vo.NewErr(errs.BuildSystemErrorInfo(errs.ServerError, err))}
		c.JSON(http.StatusOK, response)
		return
	}
	response := vo.VoidErr{Err: vo.NewErr(nil)}
	c.JSON(http.StatusOK, response)
}
