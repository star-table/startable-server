package service

import (
	"sort"

	"github.com/star-table/startable-server/common/model/vo/appvo"

	"github.com/star-table/startable-server/app/facade/appfacade"

	"github.com/star-table/startable-server/app/facade/userfacade"

	"github.com/spf13/cast"

	tablev1 "gitea.bjx.cloud/LessCode/interface/golang/table/v1"
	"github.com/star-table/startable-server/app/facade/tablefacade"
	"github.com/star-table/startable-server/app/service/projectsvc/domain"
	"github.com/star-table/startable-server/common/core/consts"
	"github.com/star-table/startable-server/common/core/errs"
	"github.com/star-table/startable-server/common/core/util/asyn"
	"github.com/star-table/startable-server/common/core/util/copyer"
	"github.com/star-table/startable-server/common/core/util/slice"
	"github.com/star-table/startable-server/common/model/bo"
	"github.com/star-table/startable-server/common/model/vo"
	"github.com/star-table/startable-server/common/model/vo/lc_table"
	"github.com/star-table/startable-server/common/model/vo/projectvo"
	"upper.io/db.v3"
)

func CreateColumn(req projectvo.CreateColumnReqVo) (*projectvo.CreateColumnReply, errs.SystemErrorInfo) {
	var err errs.SystemErrorInfo
	req.Input.TableId, err = checkTableAuth(req.OrgId, req.UserId, req.Input.AppId, req.Input.ProjectId, req.Input.TableId)
	if err != nil {
		return nil, err
	}

	columnInfo := tablefacade.CreateColumn(req)
	if columnInfo.Code != 0 {
		log.Errorf("[tablefacade.CreateColumn] err: %v", columnInfo.Error())
		return nil, errs.BuildSystemErrorInfoWithMessage(errs.TableDomainError, columnInfo.Message)
	}

	asyn.Execute(func() {
		recordColumnChangeTrends(req.OrgId, req.UserId, req.Input.ProjectId, req.Input.TableId, req.Input.Column.Name, "", "")
	})

	return columnInfo.Data, nil
}

func CopyColumn(req projectvo.CopyColumnReqVo) (*tablev1.CopyColumnReply, errs.SystemErrorInfo) {
	var err errs.SystemErrorInfo
	req.Input.TableId, err = checkTableAuth(req.OrgId, req.UserId, req.Input.AppId, req.Input.ProjectId, req.Input.TableId)
	if err != nil {
		return nil, err
	}

	columnInfo := tablefacade.CopyColumn(req)
	if columnInfo.Code != 0 {
		log.Errorf("[tablefacade.CopyColumn] err: %v", columnInfo.Error())
		return nil, errs.BuildSystemErrorInfoWithMessage(errs.TableDomainError, columnInfo.Message)
	}

	asyn.Execute(func() {
		recordColumnChangeTrends(req.OrgId, req.UserId, req.Input.ProjectId, req.Input.TableId, columnInfo.Data.CreateColumnId, "", "")
	})

	return columnInfo.Data, nil
}

func UpdateColumn(req projectvo.UpdateColumnReqVo) (*projectvo.UpdateColumnReply, errs.SystemErrorInfo) {
	params := req.Input
	if params.Column == nil {
		return nil, errs.ParamError
	}

	var err errs.SystemErrorInfo
	req.Input.TableId, err = checkTableAuth(req.OrgId, req.UserId, req.Input.AppId, req.Input.ProjectId, req.Input.TableId)
	if err != nil {
		return nil, err
	}

	//如果是更改状态需要判断有没有删除选项，删除了则需要判断该状态下是否有任务，有的话不允许更新
	if params.Column.Name == consts.BasicFieldIssueStatus {
		oldStatusList, oldStatusListErr := domain.GetTableStatus(req.OrgId, params.TableId)
		if oldStatusListErr != nil {
			log.Errorf("[UpdateColumn]GetTableStatus failed:%v, orgId:%d, tableId:%d", oldStatusListErr, req.OrgId, params.TableId)
			return nil, oldStatusListErr
		}

		newStatusList := domain.GetStatusListFromStatusColumn(params.Column)
		var newStatusIds []int64
		for _, infoBo := range newStatusList {
			newStatusIds = append(newStatusIds, infoBo.ID)
		}
		var deleteIds []int64
		for _, infoBo := range oldStatusList {
			if ok, _ := slice.Contain(newStatusIds, infoBo.ID); !ok {
				deleteIds = append(deleteIds, infoBo.ID)
			}
		}
		if len(deleteIds) > 0 {
			listReq := &tablev1.ListRawRequest{
				FilterColumns: []string{
					"count(*) as count",
				},
				Condition: &tablev1.Condition{
					Type: tablev1.ConditionType_and,
					Conditions: domain.GetNoRecycleCondition(
						domain.GetRowsCondition(consts.BasicFieldTableId, tablev1.ConditionType_equal, cast.ToString(params.TableId), nil),
						domain.GetRowsCondition(consts.BasicFieldIssueStatus, tablev1.ConditionType_in, nil, deleteIds),
					),
				},
				TableId: params.TableId,
			}

			lessResp, err := domain.GetRawRows(req.OrgId, req.UserId, listReq)
			if err != nil {
				log.Errorf("[UpdateColumn] orgId:%d, projectId:%d, LessIssueList failure:%v", req.OrgId, params.ProjectId, err.Error())
				return nil, err
			}
			issueAssignCountBos := []bo.CommonCountBo{}
			err2 := copyer.Copy(lessResp.Data, &issueAssignCountBos)
			if err2 != nil {
				log.Errorf("[UpdateColumn] orgId:%d, projectId:%d, Copy failure:%v", req.OrgId, params.ProjectId, err2)
				return nil, errs.ObjectCopyError
			}
			if len(issueAssignCountBos) > 0 && issueAssignCountBos[0].Count > 0 {
				return nil, errs.RemainIssuesInStatus
			}
		}
	}

	// 校验 && 任务栏
	updateColumnResp := tablefacade.UpdateColumn(req)
	if updateColumnResp.Code != 0 {
		log.Errorf("[tablefacade.UpdateColumn] err: %v", updateColumnResp.Error())
		return nil, errs.BuildSystemErrorInfoWithMessage(errs.TableDomainError, updateColumnResp.Message)
	}

	// 推送
	asyn.Execute(func() {
		recordColumnChangeTrends(req.OrgId, req.UserId, params.ProjectId, req.Input.TableId, "", req.Input.Column.Name, "")

		// 打开/关闭协作人字段，都需要将该字段的人添加/踢出 群
		if req.Input.Column.Field.Type == consts.LcColumnFieldTypeMember {
			domain.DealCollaboratorsColumnsToIssueChat(req.OrgId, req.Input.TableId, req.Input.Column, consts.UpdateColumn)
		}
	})

	return updateColumnResp.Data, nil
}

func DeleteColumn(req projectvo.DeleteColumnReqVo) (*tablev1.DeleteColumnReply, errs.SystemErrorInfo) {
	params := req.Input
	//默认字段不允许删除
	if params.ColumnId == consts.BasicFieldIssueStatus {
		return nil, errs.DefaultFieldError
	}

	var err errs.SystemErrorInfo
	req.Input.TableId, err = checkTableAuth(req.OrgId, req.UserId, req.Input.AppId, req.Input.ProjectId, req.Input.TableId)
	if err != nil {
		return nil, err
	}

	// 附件字段删除，把回收站的附件数据删除，相关的关联关系删除
	errSys := domain.DeleteAttachmentsForColumn(req.OrgId, req.UserId, params.ProjectId, params.TableId, params.ColumnId)
	if errSys != nil {
		// 这里不return，附件字段删除 相关的数据处理失败，不应影响删除列的主操作
		log.Errorf("[DeleteColumn] DeleteAttachmentsForColumn err:%v, orgId:%d, userId:%d, projectId:%d, tableId:%d, columnId:%s",
			errSys, req.OrgId, req.UserId, params.ProjectId, params.TableId, params.ColumnId)
	}

	deleteColumnResp := tablefacade.DeleteColumn(req)
	if deleteColumnResp.Code != 0 {
		log.Errorf("[tablefacade.DeleteColumn] err: %v", deleteColumnResp.Error())
		return nil, errs.BuildSystemErrorInfoWithMessage(errs.TableDomainError, deleteColumnResp.Message)
	}

	asyn.Execute(func() {
		recordColumnChangeTrends(req.OrgId, req.UserId, params.ProjectId, req.Input.TableId, "", "", params.ColumnId)
	})

	return deleteColumnResp.Data, nil
}

func UpdateColumnDesc(req projectvo.UpdateColumnDescriptionReqVo) (*tablev1.UpdateColumnDescriptionReply, errs.SystemErrorInfo) {
	var err errs.SystemErrorInfo
	req.Input.TableId, err = checkTableAuth(req.OrgId, req.UserId, req.Input.AppId, req.Input.ProjectId, req.Input.TableId)
	if err != nil {
		return nil, err
	}

	updateColumnDescResp := tablefacade.UpdateColumnDescription(req)
	if updateColumnDescResp.Code != 0 {
		log.Errorf("[tablefacade.UpdateColumnDesc] err: %v", updateColumnDescResp.Error())
		return nil, errs.BuildSystemErrorInfoWithMessage(errs.TableDomainError, updateColumnDescResp.Message)
	}

	return updateColumnDescResp.Data, nil

}

func recordColumnChangeTrends(orgId, userId, projectId, tableId int64, createColumn, updateColumn, deleteColumn string) {
	ext := bo.TrendExtensionBo{ProjectObjectTypeId: tableId}
	added := make([]string, 0)
	deleted := make([]string, 0)
	updated := make([]string, 0)

	added = append(added, createColumn)
	updated = append(updated, updateColumn)
	deleted = append(deleted, deleteColumn)

	ext.AddedFormFields = added
	ext.DeletedFormFields = deleted
	ext.UpdatedFormFields = updated

	projectTrendsBo := bo.ProjectTrendsBo{
		PushType:   consts.PushTypeUpdateFormField,
		OrgId:      orgId,
		ProjectId:  projectId,
		OperatorId: userId,
		Ext:        ext,
	}
	domain.PushProjectTrends(projectTrendsBo)
}

func GetOneTableColumns(req projectvo.GetTableColumnReq) (*projectvo.TableColumnsTable, errs.SystemErrorInfo) {
	// 获取汇总表
	if req.TableId <= 0 {
		return getSummaryTableColumns(req)
	}

	// 获取普通表
	return getTableColumns(req)
}

// addIterationOptions 添加迭代选项
func addIterationOptions(orgId, projectId int64, column *projectvo.TableColumnData) errs.SystemErrorInfo {
	iterationList, _, err := domain.GetIterationBoList(0, 0, db.Cond{
		consts.TcIsDelete:  consts.AppIsNoDelete,
		consts.TcProjectId: projectId,
		consts.TcOrgId:     orgId,
	}, nil)
	if err != nil {
		log.Error(err)
		return err
	}
	options := make([]*lc_table.LcOptions, 0, len(*iterationList)+1)
	options = append(options, &lc_table.LcOptions{
		Id:    0,
		Value: "未规划",
	})
	for _, iterationBo := range *iterationList {
		options = append(options, &lc_table.LcOptions{
			Id:    iterationBo.Id,
			Value: iterationBo.Name,
		})
	}

	if column.Field.Props != nil && column.Field.Props["select"] != nil {
		if m, ok := column.Field.Props["select"].(map[string]interface{}); ok {
			m["options"] = options
		}
	}

	return nil
}

func getSummaryTableColumns(req projectvo.GetTableColumnReq) (*projectvo.TableColumnsTable, errs.SystemErrorInfo) {
	columns, err := getTableColumnsByOrgId(req.OrgId, req.UserId)
	if err != nil {
		log.Errorf("[getSummaryTableColumns] 获取汇总表appID异常: %v", err)
		return nil, err
	}

	if req.NotAllIssue == 1 {
		return columns, nil
	}
	resp := tablefacade.ReadTableSchemasByOrgId(projectvo.GetTableSchemasByOrgIdReq{
		OrgId:     req.OrgId,
		UserId:    req.UserId,
		ColumnIds: []string{consts.BasicFieldIssueStatus},
	})
	if resp.Failure() {
		log.Errorf("[addIssueStatusOptions] tablefacade.ReadTableSchemasByOrgId failed, orgId:%v, err:%v", req.OrgId, resp.Error())
		return nil, errs.BuildSystemErrorInfoWithMessage(errs.TableDomainError, resp.Message)
	}

	tablesMap, err := addIssueStatusOptions(resp.Data.Tables, columns.Columns)
	if err != nil {
		return nil, err
	}
	err = addProjectOptions(req.OrgId, req.UserId, columns.Columns, tablesMap, nil)

	return columns, err
}

func getTableColumnsByOrgId(orgId, userId int64) (*projectvo.TableColumnsTable, errs.SystemErrorInfo) {
	summaryAppId, err := domain.GetOrgSummaryAppId(orgId)
	if err != nil {
		log.Errorf("[GetOneTableColumns] 获取汇总表appID异常: %v", err)
		return nil, err
	}
	tablesColumns, err := domain.GetTableColumnByAppId(orgId, userId, summaryAppId, nil)
	if err != nil {
		log.Errorf("[GetTablesColumnsList].GetTableColumnByAppId failed, orgId:%v, appId:%v, err:%v", orgId, summaryAppId, err)
		return nil, err
	}
	if len(tablesColumns.Tables) == 0 {
		return nil, errs.TableNotExist
	}

	return tablesColumns.Tables[0], nil
}

// 根据tableId获取表
func getTableColumns(req projectvo.GetTableColumnReq) (*projectvo.TableColumnsTable, errs.SystemErrorInfo) {
	columnsResp, err := domain.GetTablesColumnsByTableIds(req.OrgId, req.UserId, []int64{req.TableId}, nil, true)
	if err != nil {
		log.Errorf("[getTableColumns] failed, err:%v", err)
		return nil, err
	}

	if len(columnsResp.Tables) == 0 {
		return nil, errs.TableNotExist
	}

	columns := columnsResp.Tables[0]
	var projectInfo *bo.ProjectBo
	if req.ProjectId > 0 {
		projectInfo, err = domain.GetProjectSimple(req.OrgId, req.ProjectId)
		if err != nil {
			log.Errorf("[getTableColumns] GetProjectSimple failed, orgId:%v, projectId:%v, err:%v", req.OrgId, req.ProjectId, err)
			return nil, errs.ProjectNotExist
		}

		if projectInfo.ProjectTypeId == consts.ProjectTypeAgileId {
			for _, column := range columns.Columns {
				if column.Name == consts.BasicFieldIterationId {
					err = addIterationOptions(req.OrgId, req.ProjectId, column)
					break
				}
			}
		}
		columns.Columns = domain.RemoveColumnByProjectType(columns.Columns, projectInfo.ProjectTypeId)
	}

	table, err := domain.GetTableByTableId(req.OrgId, req.UserId, req.TableId)
	if err != nil {
		return nil, err
	}

	// 如果是文件夹类型，则嵌入相关project和tableId表头数据
	if table.SummaryFlag == consts.SummaryFlagFolder {
		childAppsResp := appfacade.GetAppList(appvo.GetAppListReq{
			OrgId:    req.OrgId,
			Type:     consts.LcAppTypeForPolaris,
			ParentId: table.AppId,
		})
		if childAppsResp.Failure() {
			return nil, childAppsResp.Error()
		}
		if len(childAppsResp.Data) > 0 {
			appIds := make([]int64, 0, len(childAppsResp.Data))
			for _, datum := range childAppsResp.Data {
				appIds = append(appIds, datum.Id)
			}

			tableIds, tablesMap, err := getAppsTableInfo(req.OrgId, req.UserId, appIds)
			if err != nil {
				return nil, err
			}

			statusColumnsResp := tablefacade.GetTableColumns(projectvo.GetTableColumnsReq{
				OrgId:  req.OrgId,
				UserId: req.UserId,
				Input: &tablev1.ReadTableSchemasRequest{
					TableIds:  tableIds,
					ColumnIds: []string{consts.BasicFieldIssueStatus},
				},
			})
			if statusColumnsResp.Failure() {
				log.Errorf("[getTableColumns] GetTableColumns failed, orgId:%v, projectId:%v, err:%v", req.OrgId, req.ProjectId, statusColumnsResp.Error())
				return nil, statusColumnsResp.Error()
			}
			for _, t := range statusColumnsResp.Data.Tables {
				temp := tablesMap[t.TableId]
				t.Name = temp.Name
				t.AppId = temp.AppId
				t.SummaryFlag = temp.SummaryFlag
			}

			err = addProjectIdColumn(req.OrgId, req.UserId, appIds, statusColumnsResp.Data.Tables, columns)
			if err != nil {
				return nil, err
			}
		}
	} else if table.SummaryFlag == consts.SummaryFlagProject && projectInfo != nil {
		appColumnsResp, err := domain.GetTableColumnByAppId(req.OrgId, req.UserId, projectInfo.AppId, []string{consts.BasicFieldIssueStatus})
		if err != nil {
			log.Errorf("[getTableColumns] GetTableColumnByAppId failed, err:%v", err)
			return nil, err
		}
		err = addTableInfo(req.OrgId, req.UserId, []int64{projectInfo.AppId}, appColumnsResp.Tables)
		if err != nil {
			return nil, err
		}

		err = addProjectIdColumn(req.OrgId, req.UserId, []int64{projectInfo.AppId}, appColumnsResp.Tables, columns)
		if err != nil {
			return nil, err
		}

	}

	return columns, err
}

func addProjectIdColumn(orgId, userId int64, appIds []int64, tables []*projectvo.TableColumnsTable,
	columns *projectvo.TableColumnsTable) errs.SystemErrorInfo {

	tableOptionsMap, err := addIssueStatusOptions(tables, columns.Columns)
	if err != nil {
		return err
	}

	projectColumn := &projectvo.TableColumnData{
		Name: consts.BasicFieldProjectId,
		Field: projectvo.TableColumnField{
			Type:       "",
			CustomType: "",
			DataType:   "",
			Props:      nil,
		},
	}
	err = addProjectOptions(orgId, userId, []*projectvo.TableColumnData{projectColumn}, tableOptionsMap, appIds)
	if err != nil {
		return err
	}
	columns.Columns = append(columns.Columns, projectColumn)

	return nil
}

func addTableInfo(orgId, userId int64, appIds []int64, tables []*projectvo.TableColumnsTable) errs.SystemErrorInfo {
	_, tablesMap, err := getAppsTableInfo(orgId, userId, appIds)
	if err != nil {
		return err
	}

	for _, columnsTable := range tables {
		temp := tablesMap[columnsTable.TableId]
		columnsTable.Name = temp.Name
		columnsTable.AppId = temp.AppId
		columnsTable.SummaryFlag = temp.SummaryFlag
	}

	return nil
}

func getAppsTableInfo(orgId, userId int64, appIds []int64) ([]int64, map[int64]*tablev1.TableMeta, errs.SystemErrorInfo) {
	tablesResp := tablefacade.GetTablesOfApp(projectvo.GetTablesOfAppReq{
		OrgId:  orgId,
		UserId: userId,
		Input:  &tablev1.ReadTablesByAppsRequest{AppIds: appIds},
	})
	if tablesResp.Failure() {
		return nil, nil, tablesResp.Error()
	}

	tableIds := make([]int64, 0, len(tablesResp.Data.AppsTables))
	tablesMap := make(map[int64]*tablev1.TableMeta, len(tablesResp.Data.AppsTables))
	for _, appsTable := range tablesResp.Data.AppsTables {
		for _, temp := range appsTable.Tables {
			tableIds = append(tableIds, temp.TableId)
			tablesMap[temp.TableId] = temp
		}
	}

	return tableIds, tablesMap, nil
}

// addProjectOptions 添加projectId的选项，加入所有project
func addProjectOptions(orgId, userId int64, columns []*projectvo.TableColumnData, tablesMap map[int64][]*lc_table.LcOptions, appIds []int64) errs.SystemErrorInfo {
	isFilling := 3
	projectList, err := Projects(projectvo.ProjectsRepVo{
		Page: 0,
		Size: 0,
		ProjectExtraBody: projectvo.ProjectExtraBody{
			Input: &vo.ProjectsReq{
				IsFiling: &isFilling,
				AppIds:   appIds,
			},
			NoNeedRedundancyInfo: true,
			ProjectTypeIds:       []int64{consts.ProjectTypeNormalId, consts.ProjectTypeAgileId, consts.ProjectTypeCommon2022V47},
		},
		OrgId:         orgId,
		UserId:        userId,
		SourceChannel: "",
	})
	if err != nil {
		log.Error(err)
		return err
	}

	var projectOptions []lc_table.LcProjectOptions
	for _, project := range projectList.List {
		child := tablesMap[cast.ToInt64(project.AppID)]
		if child == nil {
			child = make([]*lc_table.LcOptions, 0)
		}
		projectOptions = append(projectOptions, lc_table.LcProjectOptions{
			LcOptions: lc_table.LcOptions{
				Id:    project.ID,
				Value: project.Name,
				Key:   consts.BasicFieldProjectId,
			},
			Children: child,
		})
	}

	for _, column := range columns {
		if column.Name == consts.BasicFieldProjectId {
			column.Field.Type = consts.DataTypeCascader
			column.Field.CustomType = consts.DataTypeCascader
			column.Label = "所属项目与表"
			column.Field.Props = map[string]interface{}{
				consts.DataTypeCascader: map[string]interface{}{"options": projectOptions},
			}
		}
	}

	return nil
}

// 添加所有表头的issueStatus选项
func addIssueStatusOptions(tables []*projectvo.TableColumnsTable, columns []*projectvo.TableColumnData) (map[int64][]*lc_table.LcOptions, errs.SystemErrorInfo) {
	tablesMap := make(map[int64][]*lc_table.LcOptions, len(tables))
	issueStatusOptions := make([]interface{}, 0, len(tables))
	for _, table := range tables {
		if table.SummaryFlag != consts.SummaryFlagNormal && table.SummaryFlag != 0 {
			continue
		}
		tablesMap[table.AppId] = append(tablesMap[table.AppId],
			&lc_table.LcOptions{Id: cast.ToString(table.TableId), Value: table.Name, Key: consts.BasicFieldTableId})

		if len(table.Columns) >= 1 && table.Columns[0] != nil && table.Columns[0].Field.Props["groupSelect"] != nil {
			options := lc_table.ExchangeToLcGroupSelectOptions(table.Columns[0].Field.Props, cast.ToString(table.TableId))
			if len(options) > 0 {
				for i := range options {
					issueStatusOptions = append(issueStatusOptions, options[i])
				}
			}
		}
	}

	for _, table := range tablesMap {
		sort.Slice(table, func(i, j int) bool {
			return cast.ToString(table[i].Id) < cast.ToString(table[j].Id)
		})
	}

	for _, column := range columns {
		if column.Name == consts.BasicFieldIssueStatus {
			if column.Field.Props != nil && column.Field.Props["groupSelect"] != nil {
				if m, ok := column.Field.Props["groupSelect"].(map[string]interface{}); ok {
					m["options"] = issueStatusOptions
					break
				}
			}
		}
	}

	return tablesMap, nil
}

func GetTablesColumns(req projectvo.GetTablesColumnsReq) (*projectvo.TablesColumnsRespData, errs.SystemErrorInfo) {
	columnsResp, err := domain.GetTablesColumnsByTableIds(req.OrgId, req.UserId, req.Input.TableIds, req.Input.ColumnIds, false)
	if err != nil {
		log.Errorf("[GetOneTableColumns] failed, err:%v", err)
		return nil, err
	}

	if len(columnsResp.Tables) == 0 {
		return nil, errs.TableNotExist
	}

	return &projectvo.TablesColumnsRespData{Tables: columnsResp.Tables}, nil
}

func checkTableAuth(orgId, userId, appId, projectId, tableId int64) (int64, errs.SystemErrorInfo) {
	if tableId == 0 {
		manageAuthInfoResp := userfacade.GetUserAuthority(orgId, userId)
		if manageAuthInfoResp.Failure() {
			log.Errorf("[CreateColumn] orgId:%v, curUserId:%v, error:%v", orgId, userId, manageAuthInfoResp.Error())
			return 0, manageAuthInfoResp.Error()
		}
		isSysAdmin := manageAuthInfoResp.Data.IsSysAdmin
		isOrgOwner := manageAuthInfoResp.Data.IsOrgOwner
		isSubAdmin := manageAuthInfoResp.Data.IsSubAdmin

		if !isSysAdmin && !isOrgOwner && !isSubAdmin {
			return 0, errs.NoOperationPermissionForProject
		}

		// 找汇总表的状态
		summeryTableResp := tablefacade.GetSummeryTableId(projectvo.GetSummeryTableIdReqVo{
			OrgId:  orgId,
			UserId: 0,
			Input:  &tablev1.ReadSummeryTableIdRequest{},
		})
		if summeryTableResp.Failure() {
			log.Errorf("[CreateColumn] failed, orgId: %d, err: %v", orgId, summeryTableResp.Error())
			return 0, summeryTableResp.Error()
		}
		return summeryTableResp.Data.TableId, nil
	} else {
		//判断项目权限
		authErr := domain.AuthProjectWithAppId(orgId, userId, projectId, consts.RoleOperationPathOrgProProConfig, consts.OperationProConfigModifyField, appId)
		if authErr != nil {
			log.Errorf("[CreateColumn.AuthProjectWithAppId], err: %v", authErr)
			return 0, authErr
		}
	}

	return tableId, nil
}
