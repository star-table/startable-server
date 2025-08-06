package api

import (
	"net/http"
	"github.com/gin-gonic/gin"
	orgsvcService "github.com/star-table/startable-server/app/service/orgsvc/service"
	"github.com/star-table/startable-server/common/core/errs"
	"github.com/star-table/startable-server/common/core/util/copyer"
	"github.com/star-table/startable-server/common/model/vo"
	"github.com/star-table/startable-server/common/model/vo/orgvo"
	"github.com/star-table/startable-server/common/core/logger"
)

// GetOrgBoList 获取组织业务对象列表
func GetOrgBoList(c *gin.Context) {
	res, err := orgsvcService.GetOrgBoList()
	if err != nil {
		response := orgvo.GetOrgBoListRespVo{Err: vo.NewErr(errs.BuildSystemErrorInfo(errs.ServerError, err))}
		c.JSON(http.StatusOK, response)
		return
	}
	response := orgvo.GetOrgBoListRespVo{Err: vo.NewErr(nil), OrganizationBoList: res}
	c.JSON(http.StatusOK, response)
}

// SaveOrgSummaryTableAppId 保存组织汇总表应用ID
func SaveOrgSummaryTableAppId(c *gin.Context) {
	var req orgvo.SaveOrgSummaryTableAppIdReqVo
	if err := c.ShouldBindJSON(&req); err != nil {
		response := orgvo.VoidRespVo{Err: vo.NewErr(errs.BuildSystemErrorInfo(errs.ParamError, err))}
		c.JSON(http.StatusOK, response)
		return
	}
	
	res, err := orgsvcService.SaveOrgSummaryTableAppId(req.OrgId, req.UserId, req.Input)
	if err != nil {
		response := orgvo.VoidRespVo{Err: vo.NewErr(errs.BuildSystemErrorInfo(errs.ServerError, err))}
		c.JSON(http.StatusOK, response)
		return
	}
	response := orgvo.VoidRespVo{
		Err: vo.NewErr(nil),
		Data: orgvo.VoidRespVoData{
			IsOk: res,
		},
	}
	c.JSON(http.StatusOK, response)
}

// CreateOrg 创建组织
func CreateOrg(c *gin.Context) {
	var req orgvo.CreateOrgReqVo
	if err := c.ShouldBindJSON(&req); err != nil {
		response := orgvo.CreateOrgRespVo{Err: vo.NewErr(errs.BuildSystemErrorInfo(errs.ParamError, err))}
		c.JSON(http.StatusOK, response)
		return
	}
	
	sourceChannel := ""
	sourcePlatform := ""
	if req.Data.CreateOrgReq.SourceChannel != nil {
		sourceChannel = *req.Data.CreateOrgReq.SourceChannel
	}
	if req.Data.CreateOrgReq.SourcePlatform != nil {
		sourcePlatform = *req.Data.CreateOrgReq.SourcePlatform
	}
	orgId, err := orgsvcService.CreateOrgBase(req, sourceChannel, sourcePlatform)
	if err != nil {
		logger.Error(err)
		response := orgvo.CreateOrgRespVo{Err: vo.NewErr(errs.BuildSystemErrorInfo(errs.ServerError, err)), Data: orgvo.CreateOrgRespVoData{OrgId: orgId}}
		c.JSON(http.StatusOK, response)
		return
	}

	if req.Data.ImportSampleData == 1 {
		//初始化示例数据
		err1 := orgsvcService.NewbieGuideInit(orgId, req.UserId, sourceChannel)
		if err1 != nil {
			logger.Error(err1)
		}
	}
	response := orgvo.CreateOrgRespVo{Err: vo.NewErr(nil), Data: orgvo.CreateOrgRespVoData{OrgId: orgId}}
	c.JSON(http.StatusOK, response)
}

// PrepareTransferOrg 准备转移组织
func PrepareTransferOrg(c *gin.Context) {
	var req orgvo.PrepareTransferOrgReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response := orgvo.PrepareTransferOrgRespVo{Err: vo.NewErr(errs.BuildSystemErrorInfo(errs.ParamError, err))}
		c.JSON(http.StatusOK, response)
		return
	}
	
	res, err := orgsvcService.PrepareTransferOrg(req.Input.OriginOrgId, req.Input.ToOrgId)
	if err != nil {
		response := orgvo.PrepareTransferOrgRespVo{Err: vo.NewErr(errs.BuildSystemErrorInfo(errs.ServerError, err))}
		c.JSON(http.StatusOK, response)
		return
	}
	response := orgvo.PrepareTransferOrgRespVo{
		Err:  vo.NewErr(nil),
		Data: res,
	}
	c.JSON(http.StatusOK, response)
}

// TransferOrg 转移组织
func TransferOrg(c *gin.Context) {
	var req orgvo.PrepareTransferOrgReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response := vo.CommonRespVo{Err: vo.NewErr(errs.BuildSystemErrorInfo(errs.ParamError, err))}
		c.JSON(http.StatusOK, response)
		return
	}
	
	err := orgsvcService.TransferOrg(req.Input.OriginOrgId, req.Input.ToOrgId)
	if err != nil {
		response := vo.CommonRespVo{Err: vo.NewErr(errs.BuildSystemErrorInfo(errs.ServerError, err))}
		c.JSON(http.StatusOK, response)
		return
	}
	response := vo.CommonRespVo{
		Err: vo.NewErr(nil),
		Void: &vo.Void{
			ID: req.OrgId,
		},
	}
	c.JSON(http.StatusOK, response)
}

// UserOrganizationList 用户组织列表
func UserOrganizationList(c *gin.Context) {
	var req orgvo.UserOrganizationListReqVo
	if err := c.ShouldBindJSON(&req); err != nil {
		response := orgvo.UserOrganizationListRespVo{Err: vo.NewErr(errs.BuildSystemErrorInfo(errs.ParamError, err))}
		c.JSON(http.StatusOK, response)
		return
	}
	
	res, err := orgsvcService.UserOrganizationList(req.UserId)
	if err != nil {
		response := orgvo.UserOrganizationListRespVo{Err: vo.NewErr(errs.BuildSystemErrorInfo(errs.ServerError, err))}
		c.JSON(http.StatusOK, response)
		return
	}
	response := orgvo.UserOrganizationListRespVo{Err: vo.NewErr(nil), UserOrganizationListResp: res}
	c.JSON(http.StatusOK, response)
}

// SwitchUserOrganization 切换用户组织
func SwitchUserOrganization(c *gin.Context) {
	var req orgvo.SwitchUserOrganizationReqVo
	if err := c.ShouldBindJSON(&req); err != nil {
		response := orgvo.SwitchUserOrganizationRespVo{Err: vo.NewErr(errs.BuildSystemErrorInfo(errs.ParamError, err))}
		c.JSON(http.StatusOK, response)
		return
	}
	
	err := orgsvcService.SwitchUserOrganization(req.OrgId, req.UserId, req.Token)
	if err != nil {
		response := orgvo.SwitchUserOrganizationRespVo{Err: vo.NewErr(errs.BuildSystemErrorInfo(errs.ServerError, err))}
		c.JSON(http.StatusOK, response)
		return
	}
	response := orgvo.SwitchUserOrganizationRespVo{Err: vo.NewErr(nil), OrgId: req.OrgId}
	c.JSON(http.StatusOK, response)
}

// UpdateOrganizationSetting 更新组织设置
func UpdateOrganizationSetting(c *gin.Context) {
	var req orgvo.UpdateOrganizationSettingReqVo
	if err := c.ShouldBindJSON(&req); err != nil {
		response := orgvo.UpdateOrganizationSettingRespVo{Err: vo.NewErr(errs.BuildSystemErrorInfo(errs.ParamError, err))}
		c.JSON(http.StatusOK, response)
		return
	}
	
	res, err := orgsvcService.UpdateOrganizationSetting(req)
	if err != nil {
		response := orgvo.UpdateOrganizationSettingRespVo{Err: vo.NewErr(errs.BuildSystemErrorInfo(errs.ServerError, err))}
		c.JSON(http.StatusOK, response)
		return
	}
	response := orgvo.UpdateOrganizationSettingRespVo{Err: vo.NewErr(nil), OrgId: res}
	c.JSON(http.StatusOK, response)
}

// UpdateOrgRemarkSetting org 表的 remark 存放着一些配置，该接口用户更新配置中的部分项
func UpdateOrgRemarkSetting(c *gin.Context) {
	var req orgvo.UpdateOrgRemarkSettingReqVo
	if err := c.ShouldBindJSON(&req); err != nil {
		response := vo.BoolRespVo{Err: vo.NewErr(errs.BuildSystemErrorInfo(errs.ParamError, err))}
		c.JSON(http.StatusOK, response)
		return
	}
	
	isOk, err := orgsvcService.UpdateOrgRemarkSetting(req)
	if err != nil {
		response := vo.BoolRespVo{Err: vo.NewErr(errs.BuildSystemErrorInfo(errs.ServerError, err))}
		c.JSON(http.StatusOK, response)
		return
	}
	response := vo.BoolRespVo{
		Err:    vo.NewErr(nil),
		IsTrue: isOk,
	}
	c.JSON(http.StatusOK, response)
}

// OrganizationInfo 组织信息
func OrganizationInfo(c *gin.Context) {
	var req orgvo.OrganizationInfoReqVo
	if err := c.ShouldBindJSON(&req); err != nil {
		response := orgvo.OrganizationInfoRespVo{Err: vo.NewErr(errs.BuildSystemErrorInfo(errs.ParamError, err))}
		c.JSON(http.StatusOK, response)
		return
	}
	
	res, err := orgsvcService.OrganizationInfo(req)
	if err != nil {
		response := orgvo.OrganizationInfoRespVo{Err: vo.NewErr(errs.BuildSystemErrorInfo(errs.ServerError, err))}
		c.JSON(http.StatusOK, response)
		return
	}
	response := orgvo.OrganizationInfoRespVo{Err: vo.NewErr(nil), OrganizationInfo: res}
	c.JSON(http.StatusOK, response)
}

// GetOrgOutInfoByOutOrgId 根据外部组织ID获取组织外部信息
func GetOrgOutInfoByOutOrgId(c *gin.Context) {
	var req orgvo.GetOutOrgInfoByOutOrgIdReqVo
	if err := c.ShouldBindJSON(&req); err != nil {
		response := orgvo.GetOutOrgInfoByOutOrgIdRespVo{Err: vo.NewErr(errs.BuildSystemErrorInfo(errs.ParamError, err))}
		c.JSON(http.StatusOK, response)
		return
	}
	
	res, err := orgsvcService.GetOrgOutInfoByOutOrgId(req.OrgId, req.OutOrgId)
	if err != nil {
		response := orgvo.GetOutOrgInfoByOutOrgIdRespVo{Err: vo.NewErr(errs.BuildSystemErrorInfo(errs.ServerError, err))}
		c.JSON(http.StatusOK, response)
		return
	}
	response := orgvo.GetOutOrgInfoByOutOrgIdRespVo{Err: vo.NewErr(nil), Data: res}
	c.JSON(http.StatusOK, response)
}

// ScheduleOrganizationPageList 计划组织分页列表
func ScheduleOrganizationPageList(c *gin.Context) {
	var req orgvo.ScheduleOrganizationPageListReqVo
	if err := c.ShouldBindJSON(&req); err != nil {
		response := orgvo.ScheduleOrganizationPageListRespVo{Err: vo.NewErr(errs.BuildSystemErrorInfo(errs.ParamError, err))}
		c.JSON(http.StatusOK, response)
		return
	}
	
	res, err := orgsvcService.ScheduleOrganizationPageList(req)
	if err != nil {
		response := orgvo.ScheduleOrganizationPageListRespVo{Err: vo.NewErr(errs.BuildSystemErrorInfo(errs.ServerError, err))}
		c.JSON(http.StatusOK, response)
		return
	}
	response := orgvo.ScheduleOrganizationPageListRespVo{Err: vo.NewErr(nil), ScheduleOrganizationPageListResp: res}
	c.JSON(http.StatusOK, response)
}

// GetOrgIdListBySourceChannel 通过来源获取组织id列表
func GetOrgIdListBySourceChannel(c *gin.Context) {
	var req orgvo.GetOrgIdListBySourceChannelReqVo
	if err := c.ShouldBindJSON(&req); err != nil {
		response := orgvo.GetOrgIdListBySourceChannelRespVo{Err: vo.NewErr(errs.BuildSystemErrorInfo(errs.ParamError, err))}
		c.JSON(http.StatusOK, response)
		return
	}
	
	res, err := orgsvcService.GetOrgIdListBySourceChannel(req.Input.SourceChannels, req.Page, req.Size, req.Input.IsPaid)
	if err != nil {
		response := orgvo.GetOrgIdListBySourceChannelRespVo{Err: vo.NewErr(errs.BuildSystemErrorInfo(errs.ServerError, err))}
		c.JSON(http.StatusOK, response)
		return
	}
	response := orgvo.GetOrgIdListBySourceChannelRespVo{Err: vo.NewErr(nil), Data: orgvo.GetOrgIdListBySourceChannelRespData{OrgIds: res}}
	c.JSON(http.StatusOK, response)
}

// GetOrgIdListByPage 分页获取组织 id 列表
func GetOrgIdListByPage(c *gin.Context) {
	var req orgvo.GetOrgIdListByPageReqVo
	if err := c.ShouldBindJSON(&req); err != nil {
		response := orgvo.GetOrgIdListBySourceChannelRespVo{Err: vo.NewErr(errs.BuildSystemErrorInfo(errs.ParamError, err))}
		c.JSON(http.StatusOK, response)
		return
	}
	
	res, err := orgsvcService.GetOrgIdListByPage(req.Input, req.Page, req.Size)
	if err != nil {
		response := orgvo.GetOrgIdListBySourceChannelRespVo{Err: vo.NewErr(errs.BuildSystemErrorInfo(errs.ServerError, err))}
		c.JSON(http.StatusOK, response)
		return
	}
	response := orgvo.GetOrgIdListBySourceChannelRespVo{Err: vo.NewErr(nil), Data: orgvo.GetOrgIdListBySourceChannelRespData{OrgIds: res}}
	c.JSON(http.StatusOK, response)
}

// GetOrgBoListByPage 分页获取组织列表
func GetOrgBoListByPage(c *gin.Context) {
	var req orgvo.GetOrgIdListByPageReqVo
	if err := c.ShouldBindJSON(&req); err != nil {
		response := orgvo.GetOrgBoListByPageRespVo{Err: vo.NewErr(errs.BuildSystemErrorInfo(errs.ParamError, err))}
		c.JSON(http.StatusOK, response)
		return
	}
	
	resData, err := orgsvcService.GetOrgBoListByPage(req.Input, req.Page, req.Size)
	if err != nil {
		response := orgvo.GetOrgBoListByPageRespVo{Err: vo.NewErr(errs.BuildSystemErrorInfo(errs.ServerError, err))}
		c.JSON(http.StatusOK, response)
		return
	}
	response := orgvo.GetOrgBoListByPageRespVo{Err: vo.NewErr(nil), Data: resData}
	c.JSON(http.StatusOK, response)
}

// GetOrgConfig 获取组织配置
func GetOrgConfig(c *gin.Context) {
	var req orgvo.GetOrgConfigReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response := orgvo.GetOrgConfigResp{Err: vo.NewErr(errs.BuildSystemErrorInfo(errs.ParamError, err))}
		c.JSON(http.StatusOK, response)
		return
	}
	
	res, err := orgsvcService.GetOrgConfig(req.OrgId)
	respData := vo.OrgConfig{}
	copyer.Copy(res, &respData)

	if err != nil {
		response := orgvo.GetOrgConfigResp{Err: vo.NewErr(errs.BuildSystemErrorInfo(errs.ServerError, err))}
		c.JSON(http.StatusOK, response)
		return
	}
	response := orgvo.GetOrgConfigResp{Err: vo.NewErr(nil), Data: &respData}
	c.JSON(http.StatusOK, response)
}

// GetOrgConfigRest 获取组织配置，用于 restful 路由
func GetOrgConfigRest(c *gin.Context) {
	var req orgvo.GetOrgConfigReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response := orgvo.GetOrgConfigRestResp{Err: vo.NewErr(errs.BuildSystemErrorInfo(errs.ParamError, err))}
		c.JSON(http.StatusOK, response)
		return
	}
	
	res, err := orgsvcService.GetOrgConfig(req.OrgId)
	if err != nil {
		response := orgvo.GetOrgConfigRestResp{Err: vo.NewErr(errs.BuildSystemErrorInfo(errs.ServerError, err))}
		c.JSON(http.StatusOK, response)
		return
	}
	response := orgvo.GetOrgConfigRestResp{Err: vo.NewErr(nil), Data: res}
	c.JSON(http.StatusOK, response)
}

// GetOutOrgInfoByOrgIdBatch 批量根据组织ID获取外部组织信息
func GetOutOrgInfoByOrgIdBatch(c *gin.Context) {
	var req orgvo.GetOutOrgInfoByOrgIdBatchReqVo
	if err := c.ShouldBindJSON(&req); err != nil {
		response := orgvo.GetOutOrgInfoByOrgIdBatchRespVo{Err: vo.NewErr(errs.BuildSystemErrorInfo(errs.ParamError, err))}
		c.JSON(http.StatusOK, response)
		return
	}
	
	res, err := orgsvcService.GetOutOrgInfoByOrgIdBatch(req.OrgIds)
	if err != nil {
		response := orgvo.GetOutOrgInfoByOrgIdBatchRespVo{Err: vo.NewErr(errs.BuildSystemErrorInfo(errs.ServerError, err))}
		c.JSON(http.StatusOK, response)
		return
	}
	response := orgvo.GetOutOrgInfoByOrgIdBatchRespVo{Err: vo.NewErr(nil), Data: res}
	c.JSON(http.StatusOK, response)
}

// ApplyScopes 申请授权
func ApplyScopes(c *gin.Context) {
	var req orgvo.ApplyScopesReqVo
	if err := c.ShouldBindJSON(&req); err != nil {
		response := orgvo.ApplyScopesRespVo{Err: vo.NewErr(errs.BuildSystemErrorInfo(errs.ParamError, err))}
		c.JSON(http.StatusOK, response)
		return
	}
	
	res, err := orgsvcService.ApplyScopes(req.OrgId, req.UserId)
	if err != nil {
		response := orgvo.ApplyScopesRespVo{Err: vo.NewErr(errs.BuildSystemErrorInfo(errs.ServerError, err))}
		c.JSON(http.StatusOK, response)
		return
	}
	response := orgvo.ApplyScopesRespVo{Err: vo.NewErr(nil), Data: res}
	c.JSON(http.StatusOK, response)
}

// CheckSpecificScope 检查特定权限
func CheckSpecificScope(c *gin.Context) {
	var req orgvo.CheckSpecificScopeReqVo
	if err := c.ShouldBindJSON(&req); err != nil {
		response := orgvo.CheckSpecificScopeRespVo{Err: vo.NewErr(errs.BuildSystemErrorInfo(errs.ParamError, err))}
		c.JSON(http.StatusOK, response)
		return
	}
	
	res, err := orgsvcService.CheckSpecificScope(req.OrgId, req.UserId, req.PowerFlag)
	if err != nil {
		response := orgvo.CheckSpecificScopeRespVo{Err: vo.NewErr(errs.BuildSystemErrorInfo(errs.ServerError, err))}
		c.JSON(http.StatusOK, response)
		return
	}
	response := orgvo.CheckSpecificScopeRespVo{Err: vo.NewErr(nil), Data: res}
	c.JSON(http.StatusOK, response)
}

// CheckAndSetSuperAdmin 检查组织拥有者是否是超管，如果不是，则设置为超管。
func CheckAndSetSuperAdmin(c *gin.Context) {
	var req orgvo.CheckAndSetSuperAdminReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response := vo.BoolRespVo{Err: vo.NewErr(errs.BuildSystemErrorInfo(errs.ParamError, err))}
		c.JSON(http.StatusOK, response)
		return
	}
	
	res, err := orgsvcService.CheckAndSetSuperAdmin(req.OrgId)
	if err != nil {
		response := vo.BoolRespVo{Err: vo.NewErr(errs.BuildSystemErrorInfo(errs.ServerError, err))}
		c.JSON(http.StatusOK, response)
		return
	}
	response := vo.BoolRespVo{
		Err:    vo.NewErr(nil),
		IsTrue: res,
	}
	c.JSON(http.StatusOK, response)
}

// StopThirdIntegration 停止第三方集成
func StopThirdIntegration(c *gin.Context) {
	var req vo.CommonReqVo
	if err := c.ShouldBindJSON(&req); err != nil {
		response := vo.CommonRespVo{Err: vo.NewErr(errs.BuildSystemErrorInfo(errs.ParamError, err))}
		c.JSON(http.StatusOK, response)
		return
	}
	
	err := orgsvcService.StopThirdIntegration(req.OrgId, req.UserId, req.SourceChannel)
	if err != nil {
		response := vo.CommonRespVo{Err: vo.NewErr(errs.BuildSystemErrorInfo(errs.ServerError, err))}
		c.JSON(http.StatusOK, response)
		return
	}
	response := vo.CommonRespVo{
		Err:  vo.NewErr(nil),
		Void: &vo.Void{ID: req.OrgId},
	}
	c.JSON(http.StatusOK, response)
}

// SendCardToAdminForUpgrade 成员点击升级极星时，需要发送卡片给管理员
func SendCardToAdminForUpgrade(c *gin.Context) {
	var req orgvo.SendCardToAdminForUpgradeReqVo
	if err := c.ShouldBindJSON(&req); err != nil {
		response := vo.BoolRespVo{Err: vo.NewErr(errs.BuildSystemErrorInfo(errs.ParamError, err))}
		c.JSON(http.StatusOK, response)
		return
	}
	
	err := orgsvcService.SendCardToAdminForUpgrade(req.OrgId, req.UserId, req.SourceChannel)
	response := vo.BoolRespVo{
		Err:    vo.NewErr(errs.BuildSystemErrorInfo(errs.ServerError, err)),
		IsTrue: err == nil,
	}
	if err == nil {
		response.Err = vo.NewErr(nil)
	}
	c.JSON(http.StatusOK, response)
}

// GetOrgInfo 获取组织信息
func GetOrgInfo(c *gin.Context) {
	var req orgvo.GetOrgInfoReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response := orgvo.GetOrgInfoResp{Err: vo.NewErr(errs.BuildSystemErrorInfo(errs.ParamError, err))}
		c.JSON(http.StatusOK, response)
		return
	}
	
	res, err := orgsvcService.GetOrgInfoBo(req.OrgId)
	if err != nil {
		response := orgvo.GetOrgInfoResp{Err: vo.NewErr(errs.BuildSystemErrorInfo(errs.ServerError, err))}
		c.JSON(http.StatusOK, response)
		return
	}
	response := orgvo.GetOrgInfoResp{
		Err:  vo.NewErr(nil),
		Data: res,
	}
	c.JSON(http.StatusOK, response)
}
