package service

import (
	"fmt"

	tablev1 "gitea.bjx.cloud/LessCode/interface/golang/table/v1"
	"github.com/star-table/startable-server/app/facade/permissionfacade"
	"github.com/star-table/startable-server/app/facade/tablefacade"
	consts2 "github.com/star-table/startable-server/app/service/projectsvc/consts"
	"github.com/star-table/startable-server/app/service/projectsvc/domain"
	"github.com/star-table/startable-server/app/service/projectsvc/po"
	"github.com/star-table/startable-server/common/core/consts"
	"github.com/star-table/startable-server/common/core/errs"
	"github.com/star-table/startable-server/common/core/util/asyn"
	"github.com/star-table/startable-server/common/core/util/format"
	"github.com/star-table/startable-server/common/core/util/json"
	"github.com/star-table/startable-server/common/extra/lc_helper"
	"github.com/star-table/startable-server/common/library/cache"
	"github.com/star-table/startable-server/common/model/bo"
	"github.com/star-table/startable-server/common/model/vo/projectvo"
	"github.com/spf13/cast"
)

func CreateTable(req projectvo.CreateTableReq) (*projectvo.CreateTableReply, errs.SystemErrorInfo) {
	// 判断名字合法性
	if !format.VerifyTableNameFormat(req.Input.Name) {
		return nil, errs.InvalidProjectTableName
	}

	var projectTypeId, projectId int64
	if req.Input.SummaryFlag != consts.SummaryFlagFolder {
		// 获取projectId
		projectInfo, err := domain.GetProjectInfoByAppId(req.Input.AppId)
		projectId = projectInfo.Id
		projectTypeId = projectInfo.ProjectTypeId
		//_, projectId, err := domain.GetProjectIdByAppId(fmt.Sprintf("%d", req.Input.AppId))
		if err != nil {
			log.Error(err)
			return nil, errs.BuildSystemErrorInfo(errs.ProjectNotExist, err)
		}

		// 权限校验
		if err := domain.AuthProject(req.OrgId, req.UserId, projectId, consts.RoleOperationPathOrgProjectObjectType, consts.OperationProProjectTableCreate); err != nil {
			log.Error(err)
			return nil, errs.BuildSystemErrorInfo(errs.Unauthorized, err)
		}

		existErr := domain.CheckTableName(req.OrgId, req.UserId, req.Input.AppId, req.Input.Name)
		if existErr != nil {
			log.Errorf(existErr.Error())
			return nil, existErr
		}
	} else {
		projectTypeId = consts.ProjectTypeCommon2022V47
	}

	// 获取basicFields
	basicFields, err2 := domain.GetBasicFields(int(projectTypeId))
	if err2 != nil {
		log.Errorf("[CreateTable] GetBasicFields获取basicFields失败")
		return nil, err2
	}

	// 创建表
	req.Input.IsNeedColumn = true
	req.Input.BasicColumns = basicFields
	// 根据不同的项目类型，对应的表不合并一些列字段
	notNeedSummeryColumnIds := domain.GetNoNeedColumnByProjectType(projectTypeId)

	req.Input.NotNeedSummeryColumnIds = notNeedSummeryColumnIds
	if projectTypeId == consts.ProjectTypeEmpty {
		req.Input.Columns = []interface{}{lc_helper.GetSelectColumn(), lc_helper.GetMultiSelectColumn(), lc_helper.GetDocumentColumn()}
	}

	tableResp := tablefacade.CreateTable(req)
	if tableResp.Failure() {
		log.Errorf("[CreateTable] failed, orgId:%v, userId:%v, err:%v", req.OrgId, req.UserId, tableResp.Error())
		return nil, errs.BuildSystemErrorInfoWithMessage(errs.TableDomainError, "tablefacade.CreateTable error")
	}
	//创建默认视图
	viewErr := domain.CreateProjectDefaultView(req.OrgId, projectId, req.Input.AppId, projectTypeId, &tableResp.Data.Table.TableId, req.Input.SummaryFlag == consts.SummaryFlagProject || req.Input.SummaryFlag == consts.SummaryFlagFolder)
	if viewErr != nil {
		log.Errorf("[CreateTable] err: %v", viewErr)
		return nil, viewErr
	}
	asyn.Execute(func() {
		// 创建表对应的群聊配置
		if tableResp.Data.Table.TableId > 0 && projectId > 0 {
			// 检查是否需要增加该表的群聊配置。
			// 规则：如果现有的表（不包含此次新加的表）全部都有群聊配置，则视为全选了项目下的所有表，即新增的表也需要增加群聊配置。
			if err := domain.CreateGroupChatSettingForNewTable(req.OrgId, projectId, tableResp.Data.Table.TableId, req.UserId); err != nil {
				log.Errorf("[CreateTable] CreateGroupChatSettingForOneTable err: %v", err)
			}
		}
	})

	return tableResp.Data, nil
}

func RenameTable(req projectvo.RenameTableReq) (*projectvo.TableMetaData, errs.SystemErrorInfo) {
	// 获取projectId
	projectId, err := domain.GetProjectIdByAppId(req.OrgId, req.Input.AppId)
	if err != nil {
		return nil, err
	}
	tableInfo, err := domain.GetTableByTableId(req.OrgId, req.UserId, req.Input.TableId)
	if err != nil {
		return nil, err
	}
	if tableInfo.SummaryFlag != consts.SummaryFlagFolder {
		//用户角色权限校验。只有具备“管理表头字段”操作项的用户才能更新/创建/删除
		err = domain.AuthProject(req.OrgId, req.UserId, projectId, consts.RoleOperationPathOrgProjectObjectType, consts.OperationProProjectTableModify)
		if err != nil {
			log.Error(err)
			return nil, errs.BuildSystemErrorInfo(errs.Unauthorized, err)
		}

		// 校验表名是否重复
		existErr := domain.CheckTableName(req.OrgId, req.UserId, req.Input.AppId, req.Input.Name)
		if existErr != nil {
			log.Errorf(existErr.Error())
			return nil, existErr
		}
	}

	// 校验名字合法性
	isNameRight := format.VerifyTableNameFormat(req.Input.Name)
	if !isNameRight {
		return nil, errs.InvalidProjectTableName
	}

	// 更新表名
	tableResp := tablefacade.RenameTable(req)
	if tableResp.Failure() {
		log.Errorf("[RenameTable] failed, orgId:%v, userId:%v, err:%v", req.OrgId, req.UserId, tableResp.Error())
		return nil, errs.BuildSystemErrorInfoWithMessage(errs.TableDomainError, "tablefacade.RenameTable error")
	}
	// 通知飞书
	// asyn.Execute(func() {
	// 	domain.PushMessageToFeishuShortcut(bo.ShortcutPushBo{
	// 		TriggerType:         consts.FsTriggerDoProjectObjectType,
	// 		EventType:           []string{consts.FsEventUpdateProjectObjectTypeName},
	// 		OrgId:               req.OrgId,
	// 		ProjectId:           projectId,
	// 		ProjectObjectTypeId: req.Input.TableId,
	// 		IssueId:             0,
	// 		Operator:            req.UserId,
	// 	})
	// })

	return tableResp.Data, nil
}

func DeleteTable(req projectvo.DeleteTableReq) (*projectvo.TableMetaData, errs.SystemErrorInfo) {
	tableInfo, err := domain.GetTableByTableId(req.OrgId, req.UserId, req.Input.TableId)
	if err != nil {
		return nil, err
	}

	var projectInfo *po.PpmProProject
	if tableInfo.SummaryFlag != consts.SummaryFlagFolder {
		// 获取projectId
		projectInfo, err = domain.GetProjectInfoByAppId(req.Input.AppId)
		if err != nil {
			log.Error(err)
			return nil, errs.BuildSystemErrorInfo(errs.ProjectNotExist, err)
		}
		// 用户角色权限校验
		err = domain.AuthProject(req.OrgId, req.UserId, projectInfo.Id, consts.RoleOperationPathOrgProjectObjectType, consts.OperationProProjectTableDelete)
		if err != nil {
			log.Error(err)
			return nil, errs.BuildSystemErrorInfo(errs.Unauthorized, err)
		}

		// 如果所在应用有异步任务在执行，则暂不允许删表操作
		isExecuting := domain.CheckAsyncTaskIsRunning(req.OrgId, req.Input.AppId, req.Input.TableId)
		if isExecuting {
			return nil, errs.DenyDeleteTableWhenAsyncTask
		}
	}

	//var issueIds []int64
	if projectInfo != nil {
		_, err = domain.DeleteTableIssues(req.OrgId, req.UserId, projectInfo.Id, req.Input.AppId, req.Input.TableId, projectInfo.TemplateFlag)
		if err != nil {
			return nil, err
		}
	}

	deleteTableResp := tablefacade.DeleteTable(projectvo.DeleteTableReq{
		OrgId:  req.OrgId,
		UserId: req.UserId,
		Input:  &tablev1.DeleteTableRequest{AppId: req.Input.AppId, TableId: req.Input.TableId},
	})
	if deleteTableResp.Failure() {
		log.Errorf("[DeleteTable] failed, orgId:%v, userId:%v, tableId:%v, err:%v", req.Input.AppId, req.UserId, req.Input.TableId, deleteTableResp.Error())
		return nil, errs.BuildSystemErrorInfoWithMessage(errs.TableDomainError, "tablefacade.DeleteTable error")
	}

	if projectInfo != nil {
		asyn.Execute(func() {
			// 删除表后，对应的群聊配置也需要删除
			if err := domain.DeleteGroupChatSettingForOneTable(req.OrgId, projectInfo.Id, req.Input.TableId, req.UserId); err != nil {
				log.Errorf("[DeleteTable] DeleteGroupChatSettingForOneTable err: %v", err)
			}

			//// 删除动态
			//resp := trendsfacade.DeleteTrends(trendsvo.DeleteTrendsReq{
			//	OrgId: req.OrgId,
			//	Input: trendsvo.DeleteTrends{
			//		ProjectId: projectInfo.Id,
			//		IssueIds:  issueIds,
			//	},
			//})
			//if resp.Failure() {
			//	log.Errorf("[DeleteTable] DeleteTrends err:%v, orgId:%v, issueIds:%v", resp.Error(), req.OrgId, issueIds)
			//	return
			//}
		})
	}

	return &projectvo.TableMetaData{AppId: req.Input.AppId, TableId: req.Input.TableId}, nil
}

func SetAutoSchedule(req projectvo.SetAutoScheduleReq) (*projectvo.TableAutoSchedule, errs.SystemErrorInfo) {
	// 获取app权限
	optAuthResp := permissionfacade.GetAppAuth(req.OrgId, req.Input.AppId, req.Input.TableId, req.UserId)
	if optAuthResp.Failure() {
		log.Errorf("[SetAutoSchedule] orgId:%d, userId:%d, appId:%v, GetAppAuth failure:%v", req.OrgId, req.UserId, req.Input.AppId, optAuthResp.Error())
		return nil, optAuthResp.Error()
	}
	appAuthInfo := optAuthResp.Data
	isAdmin := appAuthInfo.HasAppRootPermission || appAuthInfo.SysAdmin || appAuthInfo.OrgOwner || appAuthInfo.AppOwner
	if !isAdmin {
		return nil, errs.NoOperationPermissionForProject
	}

	// 设置自动排期
	resp := tablefacade.SetAutoSchedule(projectvo.SetAutoScheduleReq{
		OrgId:  req.OrgId,
		UserId: req.UserId,
		Input: &tablev1.SetAutoScheduleRequest{
			AppId:            req.Input.AppId,
			TableId:          req.Input.TableId,
			AutoScheduleFlag: req.Input.AutoScheduleFlag,
		},
	})
	if resp.Failure() {
		log.Errorf("[SetAutoSchedule] failed, orgId:%v, userId:%v, tableId:%v, err:%v", req.OrgId, req.UserId, req.Input.TableId, resp.Error())
		return nil, errs.BuildSystemErrorInfoWithMessage(errs.TableDomainError, "tablefacade.SetAutoSchedule error")
	}

	return &projectvo.TableAutoSchedule{
		TableId:          req.Input.TableId,
		AutoScheduleFlag: req.Input.AutoScheduleFlag,
	}, nil
}

func GetTable(req projectvo.GetTableInfoReq) (*projectvo.ReadTableReply, errs.SystemErrorInfo) {
	// 直接调用go-table
	tableResp := tablefacade.ReadOneTable(req)
	if tableResp.Failure() {
		log.Errorf("[GetTable] failed, orgId:%v, userId:%v, err:%v", req.OrgId, req.UserId, tableResp.Error())
		return nil, errs.BuildSystemErrorInfoWithMessage(errs.TableDomainError, "tablefacade.ReadOneTable error")
	}
	return tableResp.Data, nil
}

func GetTables(orgId, userId, appId int64) (*projectvo.TableData, errs.SystemErrorInfo) {
	tablesResp := tablefacade.ReadTables(projectvo.GetTablesReqVo{
		OrgId:  orgId,
		UserId: userId,
		Input:  &tablev1.ReadTablesRequest{AppId: appId},
	})
	if tablesResp.Failure() {
		log.Errorf("[GetTables] failed, orgId:%v, userId:%v, err:%v", orgId, userId, tablesResp.Error())
		return nil, errs.TablesNotExist
	}

	return tablesResp.Data, nil
}

func GetTablesByApps(orgId, userId int64, appIds []int64) (*projectvo.ReadTablesByAppsData, errs.SystemErrorInfo) {
	tablesResp := tablefacade.ReadTablesByApps(projectvo.ReadTablesByAppsReqVo{
		OrgId:  orgId,
		UserId: userId,
		Input:  &tablev1.ReadTablesByAppsRequest{AppIds: appIds},
	})
	if tablesResp.Failure() {
		log.Errorf("[GetTables] failed, orgId:%v, userId:%v, err:%v", orgId, userId, tablesResp.Error())
		return nil, errs.TablesNotExist
	}

	return tablesResp.Data, nil
}

func GetTablesByOrg(orgId, userId int64) (*projectvo.TableData, errs.SystemErrorInfo) {
	tablesResp := tablefacade.GetTablesByOrg(projectvo.GetTablesByOrgReq{
		OrgId:  orgId,
		UserId: userId,
		Input:  &tablev1.ReadOrgTablesRequest{},
	})
	if tablesResp.Failure() {
		log.Errorf("[GetTablesByOrg] failed, orgId:%v, userId:%v, err:%v", orgId, userId, tablesResp.Error())
		return nil, errs.TablesNotExist
	}
	return tablesResp.Data, nil
}

func GetBigTableModeConfig(req projectvo.GetBigTableModeConfigReqVo) (*projectvo.GetBigTableModeConfigData, errs.SystemErrorInfo) {
	configJson, errCache := cache.HGet(fmt.Sprintf(consts2.CacheBigTableModeConfig, req.OrgId), cast.ToString(req.Input.TableId))
	if errCache != nil {
		log.Errorf("[GetBigTableModeConfig] cache err:%v", errCache)
		return nil, errs.RedisOperateError
	}

	if configJson == "" {
		return &projectvo.GetBigTableModeConfigData{IsBigTableMode: false, IsCanSwitch: true}, nil
	}

	config := bo.BigTableModeConfig{}
	errJson := json.FromJson(configJson, &config)
	if errJson != nil {
		log.Errorf("[GetBigTableModeConfig] cache err:%v", errJson)
		return nil, errs.JSONConvertError
	}

	return &projectvo.GetBigTableModeConfigData{IsBigTableMode: config.IsBigTableMode, IsCanSwitch: config.IsCanSwitch}, nil
}

func SwitchBigTableMode(req projectvo.SwitchBigTableModeReqVo) errs.SystemErrorInfo {
	return switchBigTableMode(req.OrgId, req.Input.TableId, &req.Input.IsBigTableMode, nil)
}

func switchBigTableMode(orgId, tableId int64, isBigTableMode, isCanSwitch *bool) errs.SystemErrorInfo {
	key := fmt.Sprintf(consts2.CacheBigTableModeConfig, orgId)
	fieldKey := cast.ToString(tableId)
	configJson, errCache := cache.HGet(key, fieldKey)
	if errCache != nil {
		log.Errorf("[SwitchBigTableMode] cache err:%v", errCache)
		return errs.RedisOperateError
	}
	config := &bo.BigTableModeConfig{IsCanSwitch: true}
	if configJson != "" {
		errJson := json.FromJson(configJson, &config)
		if errJson != nil {
			log.Errorf("[SwitchBigTableMode] cache err:%v", errJson)
			return errs.JSONConvertError
		}
	}

	if isBigTableMode != nil {
		config.IsBigTableMode = *isBigTableMode
	}
	if isCanSwitch != nil {
		config.IsCanSwitch = *isCanSwitch
	}

	value := json.ToJsonIgnoreError(config)
	errCache = cache.HSet(key, fieldKey, value)
	if errCache != nil {
		log.Errorf("[SwitchBigTableMode] cache err:%v", errCache)
		return errs.RedisOperateError
	}

	return nil
}
