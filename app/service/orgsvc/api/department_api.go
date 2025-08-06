package api

import (
	"net/http"
	"github.com/gin-gonic/gin"
	orgsvcService "github.com/star-table/startable-server/app/service/orgsvc/service"
	"github.com/star-table/startable-server/common/core/errs"
	"github.com/star-table/startable-server/common/model/vo"
	"github.com/star-table/startable-server/common/model/vo/orgvo"
)

// Departments 获取部门列表
func Departments(c *gin.Context) {
	var req orgvo.DepartmentsReqVo
	if err := c.ShouldBindJSON(&req); err != nil {
		response := orgvo.DepartmentsRespVo{Err: vo.NewErr(errs.BuildSystemErrorInfo(errs.ParamError, err))}
		c.JSON(http.StatusOK, response)
		return
	}
	
	page := req.Page
	size := req.Size
	params := req.Params
	orgId := req.OrgId

	pageA := uint(0)
	sizeA := uint(0)
	if page != nil && size != nil && *page > 0 && *size > 0 {
		pageA = uint(*page)
		sizeA = uint(*size)
	}
	res, err := orgsvcService.Departments(pageA, sizeA, params, orgId)
	if err != nil {
		response := orgvo.DepartmentsRespVo{Err: vo.NewErr(errs.BuildSystemErrorInfo(errs.ServerError, err))}
		c.JSON(http.StatusOK, response)
		return
	}
	response := orgvo.DepartmentsRespVo{Err: vo.NewErr(nil), DepartmentList: res}
	c.JSON(http.StatusOK, response)
}

// DepartmentMembers 获取部门成员
func DepartmentMembers(c *gin.Context) {
	var req orgvo.DepartmentMembersReqVo
	if err := c.ShouldBindJSON(&req); err != nil {
		response := orgvo.DepartmentMembersRespVo{Err: vo.NewErr(errs.BuildSystemErrorInfo(errs.ParamError, err))}
		c.JSON(http.StatusOK, response)
		return
	}
	
	res, err := orgsvcService.DepartmentMembers(req.Params, req.OrgId)
	if err != nil {
		response := orgvo.DepartmentMembersRespVo{Err: vo.NewErr(errs.BuildSystemErrorInfo(errs.ServerError, err))}
		c.JSON(http.StatusOK, response)
		return
	}
	response := orgvo.DepartmentMembersRespVo{Err: vo.NewErr(nil), DepartmentMemberInfos: res}
	c.JSON(http.StatusOK, response)
}

// CreateDepartment 创建部门
func CreateDepartment(c *gin.Context) {
	var req orgvo.CreateDepartmentReqVo
	if err := c.ShouldBindJSON(&req); err != nil {
		response := vo.CommonRespVo{Err: vo.NewErr(errs.BuildSystemErrorInfo(errs.ParamError, err))}
		c.JSON(http.StatusOK, response)
		return
	}
	
	res, err := orgsvcService.CreateDepartment(req.Data, req.OrgId, req.CurrentUserId)
	if err != nil {
		response := vo.CommonRespVo{Err: vo.NewErr(errs.BuildSystemErrorInfo(errs.ServerError, err))}
		c.JSON(http.StatusOK, response)
		return
	}
	response := vo.CommonRespVo{Err: vo.NewErr(nil), Void: res}
	c.JSON(http.StatusOK, response)
}

// UpdateDepartment 更新部门
func UpdateDepartment(c *gin.Context) {
	var req orgvo.UpdateDepartmentReqVo
	if err := c.ShouldBindJSON(&req); err != nil {
		response := vo.CommonRespVo{Err: vo.NewErr(errs.BuildSystemErrorInfo(errs.ParamError, err))}
		c.JSON(http.StatusOK, response)
		return
	}
	
	res, err := orgsvcService.UpdateDepartment(req.Data, req.OrgId, req.CurrentUserId)
	if err != nil {
		response := vo.CommonRespVo{Err: vo.NewErr(errs.BuildSystemErrorInfo(errs.ServerError, err))}
		c.JSON(http.StatusOK, response)
		return
	}
	response := vo.CommonRespVo{Err: vo.NewErr(nil), Void: res}
	c.JSON(http.StatusOK, response)
}

// DeleteDepartment 删除部门
func DeleteDepartment(c *gin.Context) {
	var req orgvo.DeleteDepartmentReqVo
	if err := c.ShouldBindJSON(&req); err != nil {
		response := vo.CommonRespVo{Err: vo.NewErr(errs.BuildSystemErrorInfo(errs.ParamError, err))}
		c.JSON(http.StatusOK, response)
		return
	}
	
	res, err := orgsvcService.DeleteDepartment(req.Data, req.OrgId, req.CurrentUserId)
	if err != nil {
		response := vo.CommonRespVo{Err: vo.NewErr(errs.BuildSystemErrorInfo(errs.ServerError, err))}
		c.JSON(http.StatusOK, response)
		return
	}
	response := vo.CommonRespVo{Err: vo.NewErr(nil), Void: res}
	c.JSON(http.StatusOK, response)
}

// AllocateDepartment 分配部门
func AllocateDepartment(c *gin.Context) {
	var req orgvo.AllocateDepartmentReqVo
	if err := c.ShouldBindJSON(&req); err != nil {
		response := vo.CommonRespVo{Err: vo.NewErr(errs.BuildSystemErrorInfo(errs.ParamError, err))}
		c.JSON(http.StatusOK, response)
		return
	}
	
	err := orgsvcService.AllocateDepartment(req.Data, req.OrgId, req.CurrentUserId)
	if err != nil {
		response := vo.CommonRespVo{Err: vo.NewErr(errs.BuildSystemErrorInfo(errs.ServerError, err))}
		c.JSON(http.StatusOK, response)
		return
	}
	response := vo.CommonRespVo{Err: vo.NewErr(nil)}
	c.JSON(http.StatusOK, response)
}

// SetUserDepartmentLevel 设置用户部门级别
func SetUserDepartmentLevel(c *gin.Context) {
	var req orgvo.SetUserDepartmentLevelReqVo
	if err := c.ShouldBindJSON(&req); err != nil {
		response := vo.CommonRespVo{Err: vo.NewErr(errs.BuildSystemErrorInfo(errs.ParamError, err))}
		c.JSON(http.StatusOK, response)
		return
	}
	
	err := orgsvcService.SetUserDepartmentLevel(req.OrgId, req.CurrentUserId, req.Data)
	if err != nil {
		response := vo.CommonRespVo{Err: vo.NewErr(errs.BuildSystemErrorInfo(errs.ServerError, err))}
		c.JSON(http.StatusOK, response)
		return
	}
	response := vo.CommonRespVo{Err: vo.NewErr(nil)}
	c.JSON(http.StatusOK, response)
}

// DepartmentMembersList 部门成员列表
func DepartmentMembersList(c *gin.Context) {
	var req orgvo.DepartmentMembersListReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response := orgvo.DepartmentMembersListResp{Err: vo.NewErr(errs.BuildSystemErrorInfo(errs.ParamError, err))}
		c.JSON(http.StatusOK, response)
		return
	}
	
	page := req.Page
	size := req.Size

	pageA := 0
	sizeA := 0
	if page != nil && size != nil && *page > 0 && *size > 0 {
		pageA = *page
		sizeA = *size
	}
	var name *string
	var userIds []int64
	var excludeProjectId int64 = 0
	var relationType int64 = 0
	if req.Params != nil {
		if req.Params.Name != nil {
			name = req.Params.Name
		}
		userIds = req.Params.UserIds
		if req.Params.ExcludeProjectID != nil {
			excludeProjectId = *req.Params.ExcludeProjectID
		}
		if req.Params.RelationType != nil {
			relationType = *req.Params.RelationType
		}
	}
	res, err := orgsvcService.DepartmentMembersList(req.OrgId, req.SourceChannel, name, pageA, sizeA, userIds, req.IgnoreDelete, excludeProjectId, relationType)
	if err != nil {
		response := orgvo.DepartmentMembersListResp{Err: vo.NewErr(errs.BuildSystemErrorInfo(errs.ServerError, err))}
		c.JSON(http.StatusOK, response)
		return
	}

	response := orgvo.DepartmentMembersListResp{
		Err:  vo.NewErr(nil),
		Data: res,
	}
	c.JSON(http.StatusOK, response)
}

// VerifyDepartments 验证部门
func VerifyDepartments(c *gin.Context) {
	var req orgvo.VerifyDepartmentsReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response := vo.BoolRespVo{Err: vo.NewErr(errs.BuildSystemErrorInfo(errs.ParamError, err))}
		c.JSON(http.StatusOK, response)
		return
	}
	
	isTrue := orgsvcService.VerifyDepartments(req.OrgId, req.DepartmentIds)
	response := vo.BoolRespVo{
		Err:    vo.NewErr(nil),
		IsTrue: isTrue,
	}
	c.JSON(http.StatusOK, response)
}

// GetDeptByIds 根据ID获取部门
func GetDeptByIds(c *gin.Context) {
	var req orgvo.GetDeptByIdsReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response := orgvo.GetDeptByIdsResp{Err: vo.NewErr(errs.BuildSystemErrorInfo(errs.ParamError, err))}
		c.JSON(http.StatusOK, response)
		return
	}
	
	res, err := orgsvcService.GetDeptByIds(req.OrgId, req.DeptIds)
	if err != nil {
		response := orgvo.GetDeptByIdsResp{Err: vo.NewErr(errs.BuildSystemErrorInfo(errs.ServerError, err))}
		c.JSON(http.StatusOK, response)
		return
	}
	response := orgvo.GetDeptByIdsResp{
		Err:  vo.NewErr(nil),
		List: res,
	}
	c.JSON(http.StatusOK, response)
}
