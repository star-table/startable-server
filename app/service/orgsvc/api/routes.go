package api

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	orgsvcService "github.com/star-table/startable-server/app/service/orgsvc/service"
	"github.com/star-table/startable-server/common/core/errs"
	"github.com/star-table/startable-server/common/model/vo"
	"github.com/star-table/startable-server/common/model/vo/orgvo"
	"github.com/star-table/startable-server/common/model/vo/rolevo"
	"github.com/star-table/startable-server/common/core/logger"
)

// RegisterRoutes 注册orgsvc的所有路由
// serviceName: 服务名称，如 "orgsvc"
// version: API版本，如 "v1"
func RegisterRoutes(r *gin.Engine, serviceName, version string) {
	// 构建API路径前缀，模拟mvc.BuildApi的逻辑
	apiPrefix := fmt.Sprintf("/api/%s/%s", serviceName, version)
	
	// 创建路由组
	apiGroup := r.Group(apiPrefix)
	
	// 注册GET路由
	apiGroup.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok", "message": "orgsvc service is running"})
	})
	apiGroup.GET("/getOrgBoList", func(c *gin.Context) {
		res, err := orgsvcService.GetOrgBoList()
		if err != nil {
			response := orgvo.GetOrgBoListRespVo{Err: vo.NewErr(errs.BuildSystemErrorInfo(errs.ServerError, err))}
			c.JSON(http.StatusOK, response)
			return
		}
		response := orgvo.GetOrgBoListRespVo{Err: vo.NewErr(nil), OrganizationBoList: res}
		c.JSON(http.StatusOK, response)
	})
	// TODO: 其他GET接口需要在对应文件中实现后再添加
	
	// 注册POST路由
	// org_api.go 接口
	apiGroup.POST("/saveOrgSummaryTableAppId", SaveOrgSummaryTableAppId)
	apiGroup.POST("/createOrg", CreateOrg)
	apiGroup.POST("/prepareTransferOrg", PrepareTransferOrg)
	apiGroup.POST("/transferOrg", TransferOrg)
	apiGroup.POST("/userOrganizationList", UserOrganizationList)
	apiGroup.POST("/switchUserOrganization", SwitchUserOrganization)
	apiGroup.POST("/updateOrganizationSetting", UpdateOrganizationSetting)
	apiGroup.POST("/updateOrgRemarkSetting", UpdateOrgRemarkSetting)
	apiGroup.POST("/organizationInfo", OrganizationInfo)
	apiGroup.POST("/getOrgOutInfoByOutOrgId", GetOrgOutInfoByOutOrgId)
	apiGroup.POST("/scheduleOrganizationPageList", ScheduleOrganizationPageList)
	apiGroup.POST("/getOrgIdListBySourceChannel", GetOrgIdListBySourceChannel)
	apiGroup.POST("/getOrgIdListByPage", GetOrgIdListByPage)
	apiGroup.POST("/getOrgBoListByPage", GetOrgBoListByPage)
	apiGroup.POST("/getOutOrgInfoByOrgIdBatch", GetOutOrgInfoByOrgIdBatch)
	apiGroup.POST("/applyScopes", ApplyScopes)
	apiGroup.POST("/checkSpecificScope", CheckSpecificScope)
	apiGroup.POST("/stopThirdIntegration", StopThirdIntegration)
	apiGroup.POST("/sendCardToAdminForUpgrade", SendCardToAdminForUpgrade)
	apiGroup.POST("/getOrgInfo", GetOrgInfo)
	
	// login_api.go 接口
	apiGroup.POST("/userLogin", UserLogin)
	apiGroup.POST("/userQuit", UserQuit)
	apiGroup.POST("/joinOrgByInviteCode", JoinOrgByInviteCode)
	apiGroup.POST("/exchangeShareToken", ExchangeShareToken)
	
	// department_api.go 接口
	apiGroup.POST("/departments", Departments)
	apiGroup.POST("/departmentMembers", DepartmentMembers)
	apiGroup.POST("/createDepartment", CreateDepartment)
	apiGroup.POST("/updateDepartment", UpdateDepartment)
	apiGroup.POST("/deleteDepartment", DeleteDepartment)
	apiGroup.POST("/allocateDepartment", AllocateDepartment)
	apiGroup.POST("/setUserDepartmentLevel", SetUserDepartmentLevel)
	apiGroup.POST("/departmentMembersList", DepartmentMembersList)
	apiGroup.POST("/verifyDepartments", VerifyDepartments)
	apiGroup.POST("/getDeptByIds", GetDeptByIds)
	
	// invite_api.go 接口
	apiGroup.POST("/inviteUserByPhones", InviteUserByPhones)
	apiGroup.POST("/exportInviteTemplate", ExportInviteTemplate)
	apiGroup.POST("/importMembers", ImportMembers)
	apiGroup.POST("/addUser", AddUser)
	
	// roleapi/auth_api.go 接口
	apiGroup.POST("/authenticate", Authenticate)
	apiGroup.POST("/roleUserRelation", RoleUserRelation)
	apiGroup.POST("/removeRoleUserRelation", RemoveRoleUserRelation)
	apiGroup.POST("/roleInit", RoleInit)
	apiGroup.POST("/removeRoleDepartmentRelation", RemoveRoleDepartmentRelation)
	
	// TODO: 继续添加其他文件的接口路由
	// 已重构完成的文件：
	// ✅ org_api.go - 完整重构
	// ✅ login_api.go - 完整重构
	// ✅ user_api.go - 部分重构 (PersonalInfo)
	// ✅ department_api.go - 完整重构
	// ✅ invite_api.go - 完整重构
	// ✅ roleapi/auth_api.go - 完整重构
	
	// 待重构的文件：
	// - user_api.go (剩余方法)
	// - org_member_api.go
	// - roleapi/role_api.go
	// - roleapi/user_role_api.go
	// - roleapi/permission_operation_api.go
	// - account_api.go
	// - 其他api文件
}