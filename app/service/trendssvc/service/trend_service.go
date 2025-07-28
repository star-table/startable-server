package trendssvc

import (
	"strconv"
	"time"

	"github.com/spf13/cast"

	tableV1 "gitea.bjx.cloud/LessCode/interface/golang/table/v1"
	"github.com/star-table/startable-server/app/facade/orgfacade"
	"github.com/star-table/startable-server/app/facade/projectfacade"
	"github.com/star-table/startable-server/app/facade/tablefacade"
	"github.com/star-table/startable-server/common/core/consts"
	"github.com/star-table/startable-server/common/core/errs"
	"github.com/star-table/startable-server/common/core/logger"
	"github.com/star-table/startable-server/common/core/types"
	"github.com/star-table/startable-server/common/core/util/copyer"
	"github.com/star-table/startable-server/common/core/util/json"
	"github.com/star-table/startable-server/common/core/util/lang"
	"github.com/star-table/startable-server/common/core/util/slice"
	"github.com/star-table/startable-server/common/language/english"
	"github.com/star-table/startable-server/common/library/db/mysql"
	"github.com/star-table/startable-server/common/model/bo"
	"github.com/star-table/startable-server/common/model/bo/status"
	"github.com/star-table/startable-server/common/model/vo"
	"github.com/star-table/startable-server/common/model/vo/projectvo"
	"github.com/star-table/startable-server/common/model/vo/trendsvo"
	"upper.io/db.v3"
	"upper.io/db.v3/lib/sqlbuilder"
)

var log = logger.GetDefaultLogger()

type Ext struct {
	ObjName string
}

func TrendList(orgId, currentUserId int64, input *vo.TrendReq) (*vo.TrendsList, errs.SystemErrorInfo) {

	log.Infof(consts.UserLoginSentence, currentUserId, orgId)

	trendBo := &bo.TrendsQueryCondBo{
		LastTrendID: input.LastTrendID,
		OrgId:       orgId,
		ObjId:       input.ObjID,
		ObjType:     input.ObjType,
		OperId:      input.OperID,
		StartTime:   input.StartTime,
		EndTime:     input.EndTime,
		Type:        input.Type,
		Page:        input.Page,
		Size:        input.Size,
		OrderType:   input.OrderType,
	}
	result, err := domain.QueryTrends(trendBo)
	if err != nil {
		log.Error(err)
		return nil, errs.BuildSystemErrorInfo(errs.TrendDomainError, err)
	}

	resultModel := &vo.TrendsList{}
	copyer.Copy(result, resultModel)
	creatorIds := []int64{}
	for _, v := range resultModel.List {
		creatorIds = append(creatorIds, v.Creator)
	}
	userInfo, err := orgfacade.GetBaseUserInfoBatchRelaxed(orgId, creatorIds)
	if err != nil {
		return nil, err
	}
	userInfoById := map[int64]bo.BaseUserInfoBo{}
	for _, v := range userInfo {
		userInfoById[v.UserId] = v
	}
	issueIds := []int64{}
	projectIds := []int64{}

	statusIds := []int64{}
	iterationIds := []int64{}
	customFields := []int64{}

	//2评论 5审批 6流转（仅限项目类型负责人）
	if *trendBo.Type != 2 && *trendBo.Type != 5 && *trendBo.Type != 6 {
		for _, v := range resultModel.List {
			if v.Module3 == consts.TrendsModuleIssue && v.Module3Id != 0 {
				issueIds = append(issueIds, v.Module3Id)
			}
			if v.Module2 == consts.TrendsModuleProject && v.Module2Id != 0 && v.Module3 == consts.BlankString {
				projectIds = append(projectIds, v.Module2Id)
			}
			if ok, _ := slice.Contain([]string{
				consts.TrendsRelationTypeAddCustomField,
				consts.TrendsRelationTypeDeleteCustomField,
				consts.TrendsRelationTypeUseCustomField,
				consts.TrendsRelationTypeForbidCustomField,
				consts.TrendsRelationTypeUpdateCustomField,
				consts.TrendsRelationTypeUseOrgCustomField,
			}, v.RelationType); ok {
				ext := &vo.TrendExtension{}
				err := json.FromJson(v.Ext, ext)
				if v.RelationType == consts.TrendsRelationTypeUseOrgCustomField {
					if err == nil {
						for _, id := range ext.FieldIds {
							customFields = append(customFields, *id)
						}
					}
				} else {
					customFields = append(customFields, v.RelationObjID)
				}
			}
		}
	}

	//获取最新任务信息
	issueInfos := projectfacade.GetLcIssueInfoBatch(projectvo.GetLcIssueInfoBatchReqVo{OrgId: orgId, IssueIds: issueIds})
	if issueInfos.Failure() {
		log.Error(issueInfos.Error())
		return nil, issueInfos.Error()
	}
	issueInfosMap := make(map[int64]*bo.IssueBo, len(issueInfos.Data))
	for _, data := range issueInfos.Data {
		issueInfosMap[data.Id] = data
	}

	//获取任务类型
	tableIds := []int64{}
	proIds := []int64{}
	issueTypeMap := map[int64]int64{} //任务=>类型
	for _, issue := range issueInfos.Data {
		tableIds = append(tableIds, issue.TableId)
		proIds = append(proIds, issue.ProjectId)
		issueTypeMap[issue.Id] = issue.TableId
	}

	getSimpleProjectInfo := projectfacade.GetSimpleProjectInfo(projectvo.GetSimpleProjectInfoReqVo{
		OrgId: orgId,
		Ids:   proIds,
	})
	if getSimpleProjectInfo.Failure() {
		log.Error(getSimpleProjectInfo.Error())
		return nil, getSimpleProjectInfo.Error()
	}
	appIds := []int64{}
	for _, info := range *getSimpleProjectInfo.Data {
		appId, errPhase := strconv.ParseInt(info.AppID, 10, 64)
		if errPhase != nil {
			return nil, errs.TypeConvertError
		}
		appIds = append(appIds, appId)
	}

	issueTypeInfo := map[int64]string{}
	if len(appIds) > 0 {
		tablesResp := projectfacade.GetTablesByApps(projectvo.ReadTablesByAppsReqVo{
			OrgId:  orgId,
			UserId: currentUserId,
			Input:  &tableV1.ReadTablesByAppsRequest{AppIds: appIds},
		})
		if tablesResp.Failure() {
			log.Errorf("[TrendList] projectfacade.GetTablesByApps err:%v, orgId:%v, appIds:%v",
				tablesResp.Error(), orgId, appIds)
			return nil, tablesResp.Error()
		}
		for _, tables := range tablesResp.Data.AppsTables {
			for _, ta := range tables.Tables {
				issueTypeInfo[ta.TableId] = ta.Name
			}
		}
	}

	resIssueTypeNameMap := map[int64]string{} //任务=>类型名称
	for i, i2 := range issueTypeMap {
		if temp, ok := issueTypeInfo[i2]; ok {
			resIssueTypeNameMap[i] = temp
		} else {
			resIssueTypeNameMap[i] = "任务"
		}
	}

	//获取最新项目信息
	projectResp := projectfacade.GetSimpleProjectInfo(projectvo.GetSimpleProjectInfoReqVo{OrgId: orgId, Ids: slice.SliceUniqueInt64(projectIds)})
	if projectResp.Failure() {
		log.Error(projectResp.Error())
		return nil, projectResp.Error()
	}
	projectInfo := map[int64]string{}
	for _, v := range *projectResp.Data {
		projectInfo[v.ID] = v.Name
	}

	//获取状态信息
	statusMap := map[int64]string{}
	if len(statusIds) > 0 {
		statusList, err2 := GetTableStatus(orgId, tableIds)
		if err2 != nil {
			log.Errorf(err2.Error())
			return nil, err2
		}
		for _, bos := range statusList {
			statusMap[bos.ID] = bos.Name
		}
	}

	//获取迭代信息
	iterationMap := map[int64]string{}
	if len(iterationIds) > 0 {
		allIterationInfoResp := projectfacade.GetSimpleIterationInfo(projectvo.GetSimpleIterationInfoReqVo{
			OrgId:        orgId,
			IterationIds: iterationIds,
		})
		if allIterationInfoResp.Failure() {
			log.Error(allIterationInfoResp.Error())
			return nil, allIterationInfoResp.Error()
		}
		for _, datum := range allIterationInfoResp.Data {
			iterationMap[datum.Id] = datum.Name
		}
	}

	var lastId int64
	//真正属于流转的列表
	TransferRecordList := []*vo.Trend{}
	for k, v := range resultModel.List {
		(resultModel.List)[k].ObjIsDelete = false
		if _, ok := userInfoById[v.Creator]; ok {
			(resultModel.List)[k].CreatorInfo = assemblyUserIdInfo(userInfoById[v.Creator])
		}
		ext := &vo.TrendExtension{}
		err := json.FromJson(v.Ext, ext)
		if err == nil {
			if (*ext).ObjName == nil {
				(resultModel.List)[k].OperObjName = consts.BlankString
			} else {
				(resultModel.List)[k].OperObjName = *ext.ObjName
			}
			(resultModel.List)[k].Extension = ext
		}

		if *trendBo.Type == 6 && len(v.Extension.ChangeList) > 0 {
			//判断是否是任务流转记录
			isTransferRecord := false
			for _, list := range v.Extension.ChangeList {
				if *list.Field == "ownerId" || *list.Field == "issueStatus" {
					isTransferRecord = true
					TransferRecordList = append(TransferRecordList, (resultModel.List)[k])
				}
			}
			if isTransferRecord == false {
				continue
			}
		}

		if consts.TrendsRelationTypeUpdateFormField == v.RelationType {
			(resultModel.List)[k].Extension.AddedFormFields = ext.AddedFormFields
			(resultModel.List)[k].Extension.DeletedFormFields = ext.DeletedFormFields
			(resultModel.List)[k].Extension.UpdatedFormFields = ext.UpdatedFormFields
		}

		//如果是任务，获取最新任务名称
		if v.Module3 == consts.TrendsModuleIssue && v.Module3Id != 0 {
			if issueBo, ok := issueInfosMap[v.Module3Id]; ok {
				(resultModel.List)[k].OperObjName = issueBo.Title
				if issueBo.IsDelete == consts.AppIsDeleted {
					(resultModel.List)[k].ObjIsDelete = true
				}
			} else {
				//没查到的意味着已删除
				(resultModel.List)[k].ObjIsDelete = true
			}
			if _, ok := resIssueTypeNameMap[v.Module3Id]; ok {
				temp := resIssueTypeNameMap[v.Module3Id]
				(resultModel.List)[k].Extension.IssueType = &temp
			}
		}

		//如果是项目，获取最新项目名称
		if v.Module2 == consts.TrendsModuleProject && v.Module2Id != 0 && v.Module3 == consts.BlankString {
			if _, ok := projectInfo[v.Module2Id]; ok {
				(resultModel.List)[k].OperObjName = projectInfo[v.Module2Id]
			} else {
				(resultModel.List)[k].ObjIsDelete = true
			}
		}

		if v.RelationType == consts.TrendsRelationTypeCreateIssueComment {
			if *v.NewValue != "" {
				(resultModel.List)[k].Comment = v.NewValue
			}

		}
		if len(v.Extension.ChangeList) > 0 {
			for i, list := range v.Extension.ChangeList {
				if *list.Field == "planEndTime" {
					name := consts.PlanEndTime
					(resultModel.List)[k].Extension.ChangeList[i].FieldName = &name
				}
			}
			if lang.IsEnglish() {
				for i, list := range v.Extension.ChangeList {
					if list.FieldName == nil {
						continue
					}
					if enName, ok := english.TrendsLang[*list.FieldName]; ok {
						(resultModel.List)[k].Extension.ChangeList[i].FieldName = &enName
					}
				}
			}
		}

		//如果查询的是任务（多加个判断是否是子任务处理）
		if input.ObjType != nil && *input.ObjType == consts.TrendsOperObjectTypeIssue && *input.ObjID == v.RelationObjID {
			if v.RelationType == consts.TrendsRelationTypeCreateIssue {
				(resultModel.List)[k].RelationType = consts.TrendsRelationTypeCreateChildIssue
			} else if v.RelationType == consts.TrendsRelationTypeDeleteIssue {
				(resultModel.List)[k].RelationType = consts.TrendsRelationTypeDeleteChildIssue
			} else if v.RelationType == consts.TrendsRelationTypeCreateChildIssue {
				if (resultModel.List)[k].NewValue != nil {
					(resultModel.List)[k].OperObjID = cast.ToInt64(*(resultModel.List)[k].NewValue)
				}
			} else if v.RelationType == consts.TrendsRelationTypeDeleteChildIssue {
				if (resultModel.List)[k].OldValue != nil {
					(resultModel.List)[k].OperObjID = cast.ToInt64(*(resultModel.List)[k].OldValue)
				}
			}
		}
		lastId = v.ID
	}
	resultModel.LastTrendID = lastId
	//最终记录处理
	if *trendBo.Type == 6 {
		resultModel.List = TransferRecordList
	}
	return resultModel, nil
}

func assemblyUserIdInfo(baseUserInfo bo.BaseUserInfoBo) *vo.UserIDInfo {
	return &vo.UserIDInfo{
		UserID: baseUserInfo.UserId,
		Name:   baseUserInfo.Name,
		Avatar: baseUserInfo.Avatar,
		EmplID: baseUserInfo.OutUserId,
	}
}

func AddIssueTrends(issueTrendsBo bo.IssueTrendsBo) {
	domain.AddIssueTrends(issueTrendsBo)
}

// 添加项目趋势
func AddProjectTrends(projectTrendsBo bo.ProjectTrendsBo) {
	domain.AddProjectTrends(projectTrendsBo)
}

func AddOrgTrends(orgTrendsBo bo.OrgTrendsBo) {
	domain.AddOrgTrends(orgTrendsBo)
}

func CreateTrends(trendsBo *bo.TrendsBo, tx ...sqlbuilder.Tx) (*int64, errs.SystemErrorInfo) {
	return domain.CreateTrends(trendsBo, tx...)
}

func GetProjectLatestUpdateTime(orgId int64, projectIds []int64) (map[int64]types.Time, errs.SystemErrorInfo) {
	return domain.GetProjectLatestUpdateTime(orgId, projectIds)
}

func GetIssueCommentCount(orgId int64, issueIds []int64) (map[int64]int64, errs.SystemErrorInfo) {
	return domain.GetIssueCommentCount(orgId, issueIds)
}

func GetRecentChangeTime(orgIds []int64) (map[int64]time.Time, errs.SystemErrorInfo) {
	pos := &[]po.PpmTreTrends{}
	conn, err := mysql.GetConnect()
	if err != nil {
		return nil, errs.MysqlOperateError
	}
	infoErr := conn.Select(db.Raw("max(create_time) as create_time, org_id")).From(consts.TableTrends).Where(db.Cond{
		consts.TcOrgId:    db.In(orgIds),
		consts.TcIsDelete: consts.AppIsNoDelete,
	}).GroupBy(consts.TcOrgId).All(pos)
	if infoErr != nil {
		log.Error(infoErr)
		return nil, errs.MysqlOperateError
	}

	res := map[int64]time.Time{}
	for _, trends := range *pos {
		res[trends.OrgId] = trends.CreateTime
	}

	return res, nil
}

func GetLatestOperatedProjectId(orgId int64, userId int64, projectIds []int64) (int64, errs.SystemErrorInfo) {
	info := &po.PpmTreTrends{}

	conn, err := mysql.GetConnect()
	if err != nil {
		return 0, errs.MysqlOperateError
	}
	infoErr := conn.Collection(consts.TableTrends).Find(db.Cond{
		consts.TcOrgId:     orgId,
		consts.TcModule2:   consts.TrendsModuleProject,
		consts.TcModule2Id: db.In(projectIds),
		consts.TcCreator:   userId,
	}).OrderBy("create_time desc").One(info)
	if infoErr != nil {
		if infoErr == db.ErrNoMoreRows {
			return 0, nil
		} else {
			log.Error(infoErr)
			return 0, errs.MysqlOperateError
		}
	}

	return info.Module2Id, nil
}

func GetStatusListFromStatusColumn(column *projectvo.TableColumnData) []status.StatusInfoBo {
	statusBoList := make([]status.StatusInfoBo, 0)
	field := column.Field
	// 任务状态是分组单选类型
	if field.Type != tableV1.ColumnType_groupSelect.String() {
		return statusBoList
	}
	fieldProps := projectvo.FormConfigColumnFieldGroupSelectProps{}
	propsJson := json.ToJsonIgnoreError(field.Props)
	json.FromJson(propsJson, &fieldProps)
	for _, option := range fieldProps.GroupSelect.Options {
		statusBoList = append(statusBoList, status.StatusInfoBo{
			ID:          option.ID,
			Name:        option.Value,
			DisplayName: option.Value,
			Type:        option.ParentID,
		})
	}

	return statusBoList
}

func GetTableStatus(orgId int64, tableIds []int64) ([]status.StatusInfoBo, errs.SystemErrorInfo) {
	tableColumnsResp := tablefacade.GetTableColumns(projectvo.GetTableColumnsReq{
		OrgId:  orgId,
		UserId: 0,
		Input: &tableV1.ReadTableSchemasRequest{
			TableIds:  tableIds,
			ColumnIds: []string{consts.BasicFieldIssueStatus},
		},
	})
	if tableColumnsResp.Failure() {
		log.Errorf("[GetTableStatus] failed, orgId: %d, tableIds: %v, err: %v", orgId, tableIds, tableColumnsResp.Error())
		return nil, tableColumnsResp.Error()
	}
	for _, tableSchema := range tableColumnsResp.Data.Tables {
		// curTableId := tableSchema.TableId
		columns := tableSchema.Columns
		columnMap := GetColumnMapByColumnsForTableMode(columns)
		if statusColumn, ok := columnMap[consts.BasicFieldIssueStatus]; ok {
			return GetStatusListFromStatusColumn(statusColumn), nil
		}
	}
	return []status.StatusInfoBo{}, nil
}

func GetColumnMapByColumnsForTableMode(columns []*projectvo.TableColumnData) map[string]*projectvo.TableColumnData {
	columnMap := make(map[string]*projectvo.TableColumnData, len(columns))
	for _, item := range columns {
		columnMap[item.Name] = item
	}

	return columnMap
}

// 删除任务时删除trends
func DeleteTrends(req trendsvo.DeleteTrendsReq) errs.SystemErrorInfo {
	return domain.DeleteTrends(req.OrgId, req.Input.ProjectId, req.Input.IssueIds)
}
