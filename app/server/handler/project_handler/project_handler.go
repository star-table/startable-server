package project_handler

import (
	"strconv"

	encoding "github.com/star-table/startable-server/go-common/pkg/encoding"

	permissionV1 "gitea.bjx.cloud/LessCode/interface/golang/permission/v1"
	tablev1 "gitea.bjx.cloud/LessCode/interface/golang/table/v1"

	"github.com/gin-gonic/gin"
	"github.com/star-table/startable-server/app/facade/projectfacade"
	"github.com/star-table/startable-server/app/facade/tablefacade"
	"github.com/star-table/startable-server/app/server/handler"
	"github.com/star-table/startable-server/common/core/errs"
	"github.com/star-table/startable-server/common/model/vo/projectvo"
)

type projectHandlers struct{}

var ProjectHandler projectHandlers

// ExportSameNameUserDept
// @Security ApiKeyAuth
// @Summary 导出同名的部门和用户列表，导出为 excel
// @Description 导出同名的部门和用户列表，导出为 excel
// @Tags 项目
// @accept application/json
// @Produce application/json
// @Success 200 {object} vo.ExportIssueTemplateResp
// @Failure 400
// @Router /api/rest/project/export-same-user-dept [post]
func (p projectHandlers) ExportSameNameUserDept(c *gin.Context) {
	cacheUserInfo, err := handler.GetCacheUserInfo(c)
	if err != nil {
		handler.Fail(c, err)
		return
	}
	resp := projectfacade.ExportUserOrDeptSameNameList(projectvo.ExportUserOrDeptSameNameListReqVo{
		OrgId:     cacheUserInfo.OrgId,
		UserId:    cacheUserInfo.UserId,
		ProjectId: 0,
	})
	if resp.Failure() {
		handler.Fail(c, resp.Error())
	} else {
		handler.Success(c, resp.Data)
	}
}

func (p projectHandlers) GetTableColumns(c *gin.Context) {
	cacheUserInfo, err := handler.GetCacheUserInfo(c)
	if err != nil {
		handler.Fail(c, err)
		return
	}

	var req projectvo.GetTableColumnReq
	err1 := p.unmarshal(c, &req)
	if err1 != nil {
		handler.Fail(c, errs.BuildSystemErrorInfoWithMessage(errs.ReqParamsValidateError, err1.Error()))
		return
	}

	projectId, err1 := strconv.ParseInt(c.Param("projectId"), 10, 64)
	if err1 != nil {
		handler.Fail(c, errs.ReqParamsValidateError)
		return
	}
	req.ProjectId = projectId

	req.OrgId = cacheUserInfo.OrgId
	req.UserId = cacheUserInfo.UserId

	columnsResp := projectfacade.GetOneTableColumns(req)
	if columnsResp.Failure() {
		handler.Fail(c, columnsResp.Error())
		return
	} else {
		handler.Success(c, columnsResp.Data)
	}
}

func (p projectHandlers) GetTablesColumns(c *gin.Context) {
	cacheUserInfo, err := handler.GetCacheUserInfo(c)
	if err != nil {
		handler.Fail(c, err)
		return
	}

	input := &projectvo.TablesColumnsInput{}
	err1 := p.unmarshal(c, input)
	if err1 != nil {
		handler.Fail(c, errs.BuildSystemErrorInfoWithMessage(errs.ReqParamsValidateError, err1.Error()))
		return
	}

	columnsResp := projectfacade.GetTablesColumns(projectvo.GetTablesColumnsReq{
		OrgId:  cacheUserInfo.OrgId,
		UserId: cacheUserInfo.UserId,
		Input:  input,
	})
	if columnsResp.Failure() {
		handler.Fail(c, columnsResp.Error())
		return
	} else {
		handler.Success(c, columnsResp.Data)
	}
}

//func (p projectHandlers) GetTablesColumnsList(c *gin.Context) {
//	cacheUserInfo, err := handler.GetCacheUserInfo(c)
//	if err != nil {
//		handler.Fail(c, err)
//		return
//	}
//
//	var req projectvo.GetTablesColumnsReq
//	err1 := p.unmarshal(c,&req)
//	if err1 != nil {
//		handler.Fail(c, errs.BuildSystemErrorInfoWithMessage(errs.ReqParamsValidateError, err1.Error()))
//		return
//	}
//
//
//}

func (p projectHandlers) CreateColumn(c *gin.Context) {
	cacheUserInfo, err := handler.GetCacheUserInfo(c)
	if err != nil {
		handler.Fail(c, err)
		return
	}
	var inputReqVo projectvo.CreateColumnRequest
	err1 := p.unmarshal(c, &inputReqVo)
	if err1 != nil {
		handler.Fail(c, errs.BuildSystemErrorInfoWithMessage(errs.ReqParamsValidateError, err1.Error()))
		return
	}
	req := projectvo.CreateColumnReqVo{
		SourceChannel: cacheUserInfo.SourceChannel,
		OrgId:         cacheUserInfo.OrgId,
		UserId:        cacheUserInfo.UserId,
		Input:         &inputReqVo,
	}
	resp := projectfacade.CreateColumn(req)
	if resp.Failure() {
		handler.Fail(c, resp.Error())
	} else {
		handler.Success(c, resp.Data)
	}
}

func (p projectHandlers) CopyColumn(c *gin.Context) {
	cacheUserInfo, err := handler.GetCacheUserInfo(c)
	if err != nil {
		handler.Fail(c, err)
		return
	}
	var inputReqVo projectvo.CopyColumnRequest
	err1 := p.unmarshal(c, &inputReqVo)
	if err1 != nil {
		handler.Fail(c, errs.BuildSystemErrorInfoWithMessage(errs.ReqParamsValidateError, err1.Error()))
		return
	}
	req := projectvo.CopyColumnReqVo{
		OrgId:  cacheUserInfo.OrgId,
		UserId: cacheUserInfo.UserId,
		Input:  &inputReqVo,
	}
	resp := projectfacade.CopyColumn(req)
	if resp.Failure() {
		handler.Fail(c, resp.Error())
	} else {
		handler.Success(c, resp.Data)
	}
}

func (p projectHandlers) UpdateColumn(c *gin.Context) {
	cacheUserInfo, err := handler.GetCacheUserInfo(c)
	if err != nil {
		handler.Fail(c, err)
		return
	}
	var inputReqVo projectvo.UpdateColumnReqVoInput
	err1 := p.unmarshal(c, &inputReqVo)
	if err1 != nil {
		handler.Fail(c, errs.BuildSystemErrorInfoWithMessage(errs.ReqParamsValidateError, err1.Error()))
		return
	}
	req := projectvo.UpdateColumnReqVo{
		OrgId:         cacheUserInfo.OrgId,
		UserId:        cacheUserInfo.UserId,
		SourceChannel: cacheUserInfo.SourceChannel,
		Input:         &inputReqVo,
	}
	resp := projectfacade.UpdateColumn(req)
	if resp.Failure() {
		handler.Fail(c, resp.Error())
	} else {
		handler.Success(c, resp.Data)
	}
}

func (p projectHandlers) UpdateColumnDescription(c *gin.Context) {
	cacheUserInfo, err := handler.GetCacheUserInfo(c)
	if err != nil {
		handler.Fail(c, err)
		return
	}
	var inputReqVo projectvo.UpdateColumnDescriptionRequest
	err1 := p.unmarshal(c, &inputReqVo)
	if err1 != nil {
		handler.Fail(c, errs.BuildSystemErrorInfoWithMessage(errs.ReqParamsValidateError, err1.Error()))
		return
	}
	req := projectvo.UpdateColumnDescriptionReqVo{
		OrgId:         cacheUserInfo.OrgId,
		UserId:        cacheUserInfo.UserId,
		SourceChannel: cacheUserInfo.SourceChannel,
		Input:         &inputReqVo,
	}
	resp := projectfacade.UpdateColumnDescription(req)
	if resp.Failure() {
		handler.Fail(c, resp.Error())
	} else {
		handler.Success(c, resp.Data)
	}
}

func (p projectHandlers) DeleteColumn(c *gin.Context) {
	cacheUserInfo, err := handler.GetCacheUserInfo(c)
	if err != nil {
		handler.Fail(c, err)
		return
	}
	var inputReqVo projectvo.DeleteColumnRequest
	err1 := p.unmarshal(c, &inputReqVo)
	if err1 != nil {
		handler.Fail(c, errs.BuildSystemErrorInfoWithMessage(errs.ReqParamsValidateError, err1.Error()))
		return
	}
	req := projectvo.DeleteColumnReqVo{
		OrgId:         cacheUserInfo.OrgId,
		UserId:        cacheUserInfo.UserId,
		SourceChannel: cacheUserInfo.SourceChannel,
		Input:         &inputReqVo,
	}
	resp := projectfacade.DeleteColumn(req)
	if resp.Failure() {
		handler.Fail(c, resp.Error())
	} else {
		handler.Success(c, resp.Data)
	}
}

func (p projectHandlers) CreateTable(c *gin.Context) {
	cacheUserInfo, err := handler.GetCacheUserInfo(c)
	if err != nil {
		handler.Fail(c, err)
		return
	}
	var inputReqVo projectvo.CreateTableData
	err1 := p.unmarshal(c, &inputReqVo)
	if err1 != nil {
		handler.Fail(c, errs.BuildSystemErrorInfoWithMessage(errs.ReqParamsValidateError, err1.Error()))
		return
	}

	req := projectvo.CreateTableReq{
		OrgId:  cacheUserInfo.OrgId,
		UserId: cacheUserInfo.UserId,
		Input:  &inputReqVo,
	}

	resp := projectfacade.CreateTable(req)
	if resp.Failure() {
		handler.Fail(c, resp.Error())
	} else {
		handler.Success(c, resp.Data)
	}
}

func (p projectHandlers) RenameTable(c *gin.Context) {
	cacheUserInfo, err := handler.GetCacheUserInfo(c)
	if err != nil {
		handler.Fail(c, err)
		return
	}
	var inputReqVo tablev1.RenameTableRequest
	err1 := p.unmarshal(c, &inputReqVo)
	if err1 != nil {
		handler.Fail(c, errs.BuildSystemErrorInfoWithMessage(errs.ReqParamsValidateError, err1.Error()))
		return
	}

	req := projectvo.RenameTableReq{
		OrgId:  cacheUserInfo.OrgId,
		UserId: cacheUserInfo.UserId,
		Input:  &inputReqVo,
	}

	resp := projectfacade.RenameTable(req)
	if resp.Failure() {
		handler.Fail(c, resp.Error())
	} else {
		handler.Success(c, resp.Data)
	}
}

func (p projectHandlers) DeleteTable(c *gin.Context) {
	cacheUserInfo, err := handler.GetCacheUserInfo(c)
	if err != nil {
		handler.Fail(c, err)
		return
	}
	var inputReqVo tablev1.DeleteTableRequest
	err1 := p.unmarshal(c, &inputReqVo)
	if err1 != nil {
		handler.Fail(c, errs.BuildSystemErrorInfoWithMessage(errs.ReqParamsValidateError, err1.Error()))
		return
	}
	req := projectvo.DeleteTableReq{
		OrgId:  cacheUserInfo.OrgId,
		UserId: cacheUserInfo.UserId,
		Input:  &inputReqVo,
	}
	resp := projectfacade.DeleteTable(req)
	if resp.Failure() {
		handler.Fail(c, resp.Error())
	} else {
		handler.Success(c, resp.Data)
	}
}

func (p projectHandlers) SetTableAutoSchedule(c *gin.Context) {
	cacheUserInfo, err := handler.GetCacheUserInfo(c)
	if err != nil {
		handler.Fail(c, err)
		return
	}
	var inputReqVo tablev1.SetAutoScheduleRequest
	err1 := p.unmarshal(c, &inputReqVo)
	if err1 != nil {
		handler.Fail(c, errs.BuildSystemErrorInfoWithMessage(errs.ReqParamsValidateError, err1.Error()))
		return
	}
	req := projectvo.SetAutoScheduleReq{
		OrgId:  cacheUserInfo.OrgId,
		UserId: cacheUserInfo.UserId,
		Input:  &inputReqVo,
	}
	resp := projectfacade.SetAutoSchedule(req)
	if resp.Failure() {
		handler.Fail(c, resp.Error())
	} else {
		handler.Success(c, resp.Data)
	}
}

func (p projectHandlers) GetOneTable(c *gin.Context) {
	cacheUserInfo, err := handler.GetCacheUserInfo(c)
	if err != nil {
		handler.Fail(c, err)
		return
	}
	var inputReqVo tablev1.ReadTableRequest
	err1 := p.unmarshal(c, &inputReqVo)
	if err1 != nil {
		handler.Fail(c, errs.BuildSystemErrorInfoWithMessage(errs.ReqParamsValidateError, err1.Error()))
		return
	}
	req := projectvo.GetTableInfoReq{
		OrgId:  cacheUserInfo.OrgId,
		UserId: cacheUserInfo.UserId,
		Input:  &inputReqVo,
	}
	resp := projectfacade.GetTable(req)
	if resp.Failure() {
		handler.Fail(c, resp.Error())
	} else {
		handler.Success(c, resp.Data)
	}
}

func (p projectHandlers) GetTableList(c *gin.Context) {
	cacheUserInfo, err := handler.GetCacheUserInfo(c)
	if err != nil {
		handler.Fail(c, err)
		return
	}
	var inputReqVo tablev1.ReadTablesRequest
	err1 := p.unmarshal(c, &inputReqVo)
	if err1 != nil {
		handler.Fail(c, errs.BuildSystemErrorInfoWithMessage(errs.ReqParamsValidateError, err1.Error()))
		return
	}
	req := projectvo.GetTablesReqVo{
		OrgId:  cacheUserInfo.OrgId,
		UserId: cacheUserInfo.UserId,
		Input:  &inputReqVo,
	}
	resp := projectfacade.GetTables(req)
	if resp.Failure() {
		handler.Fail(c, resp.Error())
	} else {
		handler.Success(c, resp.Data)
	}
}

func (p projectHandlers) GetTablesByApps(c *gin.Context) {
	cacheUserInfo, err := handler.GetCacheUserInfo(c)
	if err != nil {
		handler.Fail(c, err)
		return
	}
	var inputReqVo tablev1.ReadTablesByAppsRequest
	err1 := p.unmarshal(c, &inputReqVo)
	if err1 != nil {
		handler.Fail(c, errs.BuildSystemErrorInfoWithMessage(errs.ReqParamsValidateError, err1.Error()))
		return
	}
	req := projectvo.ReadTablesByAppsReqVo{
		OrgId:  cacheUserInfo.OrgId,
		UserId: cacheUserInfo.UserId,
		Input:  &inputReqVo,
	}
	resp := projectfacade.GetTablesByApps(req)
	if resp.Failure() {
		handler.Fail(c, resp.Error())
	} else {
		handler.Success(c, resp.Data)
	}
}

func (p projectHandlers) GetTableListByOrg(c *gin.Context) {
	cacheUserInfo, err := handler.GetCacheUserInfo(c)
	if err != nil {
		handler.Fail(c, err)
		return
	}
	req := projectvo.GetTablesByOrgReq{
		OrgId:  cacheUserInfo.OrgId,
		UserId: cacheUserInfo.UserId,
		Input:  nil,
	}
	resp := projectfacade.GetTablesByOrg(req)
	if resp.Failure() {
		handler.Fail(c, resp.Error())
	} else {
		handler.Success(c, resp.Data)
	}
}

func (p projectHandlers) GetMenu(c *gin.Context) {
	cacheUserInfo, err := handler.GetCacheUserInfo(c)
	if err != nil {
		handler.Fail(c, err)
		return
	}

	var inputReqVo projectvo.GetMenuReq
	err1 := p.unmarshal(c, &inputReqVo)
	if err1 != nil {
		handler.Fail(c, errs.BuildSystemErrorInfoWithMessage(errs.ReqParamsValidateError, err1.Error()))
		return
	}
	req := projectvo.GetMenuReqVo{
		OrgId: cacheUserInfo.OrgId,
		AppId: inputReqVo.AppId,
	}
	resp := projectfacade.GetMenu(req)
	if resp.Failure() {
		handler.Fail(c, resp.Error())
	} else {
		handler.Success(c, resp.Data)
	}
}

func (p projectHandlers) SaveMenu(c *gin.Context) {
	cacheUserInfo, err := handler.GetCacheUserInfo(c)
	if err != nil {
		handler.Fail(c, err)
		return
	}
	var inputReqVo projectvo.SaveMenuData
	err2 := p.unmarshal(c, &inputReqVo)
	if err2 != nil {
		handler.Fail(c, errs.BuildSystemErrorInfoWithMessage(errs.ReqParamsValidateError, err2.Error()))
		return
	}
	req := projectvo.SaveMenuReqVo{
		OrgId: cacheUserInfo.OrgId,
		Input: projectvo.SaveMenuData{
			AppId:  inputReqVo.AppId,
			Config: inputReqVo.Config,
		},
	}
	resp := projectfacade.SaveMenu(req)
	if resp.Failure() {
		handler.Fail(c, resp.Error())
	} else {
		handler.Success(c, resp.Data)
	}
}

func (p projectHandlers) unmarshal(c *gin.Context, v interface{}) error {
	bts, err := c.GetRawData()
	if err != nil {
		return err
	}

	return encoding.GetJsonCodec().Unmarshal(bts, v)
}

func (p projectHandlers) GetCollaborators(c *gin.Context) {
	cacheUserInfo, errSys := handler.GetCacheUserInfo(c)
	if errSys != nil {
		handler.Fail(c, errSys)
		return
	}

	var inputReqVo permissionV1.GetCollaboratorsRequest
	err := p.unmarshal(c, &inputReqVo)
	if err != nil {
		handler.Fail(c, errs.BuildSystemErrorInfoWithMessage(errs.ReqParamsValidateError, err.Error()))
		return
	}

	collaborators, errSys := tablefacade.GetAppCollaborators(cacheUserInfo.OrgId, cacheUserInfo.UserId, inputReqVo.AppId)
	if errSys != nil {
		handler.Fail(c, errSys)
		return
	} else {
		handler.Success(c, collaborators)
	}
}

// QueryProcessForAsyncTask 查询异步任务的进度。比如：导入任务、应用模板等
func (p projectHandlers) QueryProcessForAsyncTask(c *gin.Context) {
	cacheUserInfo, err := handler.GetCacheUserInfo(c)
	if err != nil {
		handler.Fail(c, err)
		return
	}
	inputReq := projectvo.QueryProcessForAsyncTaskReqVoData{}
	err2 := p.unmarshal(c, &inputReq)
	if err2 != nil {
		handler.Fail(c, errs.BuildSystemErrorInfoWithMessage(errs.ReqParamsValidateError, err2.Error()))
		return
	}
	req := projectvo.QueryProcessForAsyncTaskReqVo{
		OrgId:  cacheUserInfo.OrgId,
		UserId: cacheUserInfo.UserId,
		Input:  inputReq,
	}
	resp := projectfacade.QueryProcessForAsyncTask(req)
	if resp.Failure() {
		handler.Fail(c, resp.Error())
	} else {
		handler.Success(c, resp.Data)
	}
}

// PayLimitNum 组织下项目等资源统计信息，用于未付费资源限制
func (p projectHandlers) PayLimitNum(c *gin.Context) {
	cacheUserInfo, err := handler.GetCacheUserInfo(c)
	if err != nil {
		handler.Fail(c, err)
		return
	}
	resp := projectfacade.PayLimitNumForRest(projectvo.PayLimitNumReq{
		OrgId: cacheUserInfo.OrgId,
	})
	if resp.Failure() {
		handler.Fail(c, resp.Error())
	} else {
		handler.Success(c, resp.Data)
	}
}

func (p projectHandlers) GetProjectMemberIds(c *gin.Context) {
	cacheUserInfo, err := handler.GetCacheUserInfo(c)
	if err != nil {
		handler.Fail(c, err)
		return
	}
	inputReq := projectvo.GetProjectMemberIdsReq{}
	err2 := p.unmarshal(c, &inputReq)
	if err2 != nil {
		handler.Fail(c, errs.BuildSystemErrorInfoWithMessage(errs.ReqParamsValidateError, err2.Error()))
		return
	}
	req := projectvo.GetProjectMemberIdsReqVo{
		OrgId:  cacheUserInfo.OrgId,
		UserId: cacheUserInfo.UserId,
		Input:  inputReq,
	}
	resp := projectfacade.GetProjectMemberIds(req)
	if resp.Failure() {
		handler.Fail(c, resp.Error())
	} else {
		handler.Success(c, resp.Data)
	}
}

// CheckIsShowProChatIcon 检查是否展示项目群聊 icon
func (p projectHandlers) CheckIsShowProChatIcon(c *gin.Context) {
	cacheUserInfo, err := handler.GetCacheUserInfo(c)
	if err != nil {
		handler.Fail(c, err)
		return
	}
	inputReq := projectvo.CheckIsShowProChatIconReqData{}
	err2 := p.unmarshal(c, &inputReq)
	if err2 != nil {
		handler.Fail(c, errs.BuildSystemErrorInfoWithMessage(errs.ReqParamsValidateError, err2.Error()))
		return
	}
	req := projectvo.CheckIsShowProChatIconReq{
		OrgId:         cacheUserInfo.OrgId,
		UserId:        cacheUserInfo.UserId,
		SourceChannel: cacheUserInfo.SourceChannel,
		Input:         inputReq,
	}
	resp := projectfacade.CheckIsShowProChatIcon(req)
	if resp.Failure() {
		handler.Fail(c, resp.Error())
	} else {
		handler.Success(c, resp.Data)
	}
}

func (p projectHandlers) FilterForTrendsMembers(c *gin.Context) {
	cacheUserInfo, err := handler.GetCacheUserInfo(c)
	if err != nil {
		handler.Fail(c, err)
		return
	}
	inputReq := projectvo.GetTrendListMembersReq{}
	err2 := p.unmarshal(c, &inputReq)
	if err2 != nil {
		handler.Fail(c, errs.BuildSystemErrorInfoWithMessage(errs.ReqParamsValidateError, err2.Error()))
		return
	}
	req := projectvo.GetTrendListMembersReqVo{
		OrgId:  cacheUserInfo.OrgId,
		UserId: cacheUserInfo.UserId,
		Input:  inputReq,
	}

	resp := projectfacade.GetTrendsMembers(req)
	if resp.Failure() {
		handler.Fail(c, resp.Error())
	} else {
		handler.Success(c, resp.Data)
	}
}

func (p projectHandlers) FilterProjectStatistics(c *gin.Context) {
	cacheUserInfo, err := handler.GetCacheUserInfo(c)
	if err != nil {
		handler.Fail(c, err)
		return
	}
	inputReq := projectvo.GetProjectStatisticsReq{}
	err2 := p.unmarshal(c, &inputReq)
	if err2 != nil {
		handler.Fail(c, errs.BuildSystemErrorInfoWithMessage(errs.ReqParamsValidateError, err2.Error()))
		return
	}
	req := projectvo.GetProjectStatisticsReqVo{
		OrgId:  cacheUserInfo.OrgId,
		UserId: cacheUserInfo.UserId,
		Input:  inputReq,
	}
	resp := projectfacade.GetProjectStatistics(req)
	if resp.Failure() {
		handler.Fail(c, resp.Error())
	} else {
		handler.Success(c, resp.Data)
	}
}

//func (p projectHandlers) GetShareViewInfoByKey(c *gin.Context) {
//	inputReq := &projectvo.GetShareViewInfoByKeyReq{}
//	err2 := p.unmarshal(c, inputReq)
//	if err2 != nil {
//		handler.Fail(c, errs.BuildSystemErrorInfoWithMessage(errs.ReqParamsValidateError, err2.Error()))
//		return
//	}
//	resp := projectfacade.GetShareViewInfoByKey(inputReq)
//	if resp.Failure() {
//		handler.Fail(c, resp.Error())
//	} else {
//		handler.Success(c, resp.Data)
//	}
//}
