package service

import (
	"fmt"
	"sort"
	"strconv"
	"strings"
	"time"

	msgPb "gitea.bjx.cloud/LessCode/interface/golang/msg/v1"
	tablePb "gitea.bjx.cloud/LessCode/interface/golang/table/v1"
	"github.com/star-table/startable-server/app/facade/common/report"
	"github.com/star-table/startable-server/app/facade/formfacade"
	"github.com/star-table/startable-server/app/facade/idfacade"
	"github.com/star-table/startable-server/app/facade/tablefacade"
	"github.com/star-table/startable-server/app/service/projectsvc/domain"
	"github.com/star-table/startable-server/common/core/consts"
	"github.com/star-table/startable-server/common/core/errs"
	"github.com/star-table/startable-server/common/core/threadlocal"
	"github.com/star-table/startable-server/common/core/util/asyn"
	"github.com/star-table/startable-server/common/core/util/copyer"
	"github.com/star-table/startable-server/common/core/util/json"
	"github.com/star-table/startable-server/common/core/util/jsonx"
	"github.com/star-table/startable-server/common/core/util/slice"
	slice2 "github.com/star-table/startable-server/common/core/util/slice"
	"github.com/star-table/startable-server/common/extra/lc_helper"
	"github.com/star-table/startable-server/common/library/db/mysql"
	"github.com/star-table/startable-server/common/model/bo"
	"github.com/star-table/startable-server/common/model/vo"
	"github.com/star-table/startable-server/common/model/vo/commonvo"
	"github.com/star-table/startable-server/common/model/vo/datacenter"
	"github.com/star-table/startable-server/common/model/vo/formvo"
	"github.com/star-table/startable-server/common/model/vo/projectvo"
	"github.com/star-table/startable-server/common/model/vo/uservo"
	"github.com/spf13/cast"
	"upper.io/db.v3/lib/sqlbuilder"
)

func JudgeProjectFiling(orgId, projectId int64) errs.SystemErrorInfo {
	projectInfo, err := domain.LoadProjectAuthBo(orgId, projectId)
	if err != nil {
		log.Error(err)
		return err
	}
	if projectInfo.IsFilling == consts.ProjectIsFiling {
		return errs.ProjectIsFilingYet
	}

	return nil
}

func JudgeProjectFilingByIssueId(orgId, issueId int64) errs.SystemErrorInfo {
	issueBo, err := domain.GetIssueInfoLc(orgId, 0, issueId)
	if err != nil {
		log.Error(err)
		return err
	}
	projectId := issueBo.ProjectId
	if projectId != 0 {
		projectInfo, err := domain.LoadProjectAuthBo(orgId, projectId)
		if err != nil {
			log.Error(err)
			return err
		}
		if projectInfo.IsFilling == consts.ProjectIsFiling {
			return errs.ProjectIsFilingYet
		}
	}

	return nil
}

func DeleteIssueWithoutAuth(reqVo projectvo.DeleteIssueReqVo) (*vo.Issue, errs.SystemErrorInfo) {
	input := reqVo.Input
	orgId := reqVo.OrgId
	currentUserId := reqVo.UserId
	sourceChannel := reqVo.SourceChannel

	issueBos, err2 := domain.GetIssueInfosLc(orgId, reqVo.UserId, []int64{input.ID})
	if err2 != nil {
		log.Error(err2)
		return nil, errs.BuildSystemErrorInfo(errs.IssueDomainError, err2)
	}
	if len(issueBos) < 1 {
		log.Errorf("[DeleteIssueWithoutAuth] not found issue issueId:%v", input.ID)
		return nil, errs.IssueNotExist
	}
	issueBo := issueBos[0]
	err := domain.AuthIssue(orgId, currentUserId, issueBo, consts.RoleOperationPathOrgProIssueT, consts.OperationProIssue4Delete)
	if err != nil {
		log.Error(err)
		return nil, errs.BuildSystemErrorInfo(errs.Unauthorized, err)
	}

	childIds, err3 := domain.DeleteIssue(issueBo, currentUserId, sourceChannel, input.TakeChildren)
	if err3 != nil {
		log.Error(err3)
		return nil, errs.BuildSystemErrorInfo(errs.IssueDomainError, err3)
	}

	// 事件上报
	asyn.Execute(func() {
		// 拿任务信息
		var deletedIssueIds []int64
		deletedIssueIds = append(deletedIssueIds, input.ID)
		deletedIssueIds = append(deletedIssueIds, childIds...)
		columnIds := []string{
			lc_helper.ConvertToFilterColumn(consts.BasicFieldTitle),
			lc_helper.ConvertToFilterColumn(consts.BasicFieldAppId),
			lc_helper.ConvertToFilterColumn(consts.BasicFieldProjectId),
			lc_helper.ConvertToFilterColumn(consts.BasicFieldTableId),
		}
		data, errSys := domain.GetIssueInfosMapLcByIssueIds(orgId, currentUserId, deletedIssueIds, columnIds...)
		if errSys != nil {
			return
		}

		openTraceId, _ := threadlocal.Mgr.GetValue(consts.JaegerContextTraceKey)
		openTraceIdStr := cast.ToString(openTraceId)

		for _, d := range data {
			dataId := cast.ToInt64(d[consts.BasicFieldId])
			issueId := cast.ToInt64(d[consts.BasicFieldIssueId])
			appId := cast.ToInt64(d[consts.BasicFieldAppId])
			projectId := cast.ToInt64(d[consts.BasicFieldProjectId])
			tableId := cast.ToInt64(d[consts.BasicFieldTableId])
			e := &commonvo.DataEvent{
				OrgId:     orgId,
				AppId:     appId,
				ProjectId: projectId,
				TableId:   tableId,
				DataId:    dataId,
				IssueId:   issueId,
				UserId:    currentUserId,
			}
			report.ReportDataEvent(msgPb.EventType_DataDeleted, openTraceIdStr, e)
		}

	})

	result := &vo.Issue{}
	copyErr := copyer.Copy(issueBo, result)
	if copyErr != nil {
		log.Errorf("copyer.Copy(): %q\n", copyErr)
	}
	result.IssueIds = append(childIds, input.ID)
	return result, nil
}

func DeleteIssue(reqVo projectvo.DeleteIssueReqVo) (*vo.Issue, errs.SystemErrorInfo) {
	input := reqVo.Input
	orgId := reqVo.OrgId
	currentUserId := reqVo.UserId

	issueBo, err2 := domain.GetIssueInfoLc(orgId, 0, input.ID)
	if err2 != nil {
		log.Error(err2)
		return nil, errs.BuildSystemErrorInfo(errs.IssueDomainError, err2)
	}

	err := domain.AuthIssue(orgId, currentUserId, issueBo, consts.RoleOperationPathOrgProIssueT, consts.OperationProIssue4Delete)
	if err != nil {
		log.Error(err)
		return nil, errs.BuildSystemErrorInfo(errs.Unauthorized, err)
	}
	return DeleteIssueWithoutAuth(reqVo)
}

func DeleteIssueBatch(reqVo projectvo.DeleteIssueBatchReqVo) (*vo.DeleteIssueBatchResp, errs.SystemErrorInfo) {
	input := reqVo.Input
	orgId := reqVo.OrgId
	currentUserId := reqVo.UserId
	sourceChannel := reqVo.SourceChannel

	//查找目标任务以及所有的子任务
	issueAndChildrenIds, err := domain.GetIssueAndChildrenIds(orgId, input.Ids)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	issueAndChildren, err := domain.GetIssueInfosLc(orgId, currentUserId, issueAndChildrenIds)
	if err != nil {
		log.Errorf("[DeleteIssueBatch] GetIssueInfosLc err: %v", err)
		return nil, errs.IssueNotExist
	}

	allNeedAuth := []*bo.IssueBo{}
	allNeedAuthIssueIds := make([]int64, 0)
	for _, child := range issueAndChildren {
		if ok, _ := slice.Contain(input.Ids, child.Id); ok {
			if child.ParentId == 0 {
				allNeedAuth = append(allNeedAuth, child)
				allNeedAuthIssueIds = append(allNeedAuthIssueIds, child.Id)
			} else if ok1, _ := slice.Contain(input.Ids, child.ParentId); !ok1 {
				//找出需要判断权限的(只要父任务不在所选范围内，就表示他是父任务)
				allNeedAuth = append(allNeedAuth, child)
				allNeedAuthIssueIds = append(allNeedAuthIssueIds, child.Id)
			}
		}
	}

	notAuthPassIssues := []*bo.IssueBo{}
	noAuthIds := []int64{}
	trulyIssueBos := []*bo.IssueBo{}

	var notAuthPath []string
	for _, issueBo := range allNeedAuth {
		//权限判断
		err := domain.AuthIssueWithAppId(orgId, currentUserId, issueBo, consts.RoleOperationPathOrgProIssueT, consts.OperationProIssue4Delete, reqVo.InputAppId)
		if err != nil {
			//log.Error(err)
			notAuthPassIssues = append(notAuthPassIssues, issueBo)
			noAuthIds = append(noAuthIds, issueBo.Id)
			notAuthPath = append(notAuthPath, fmt.Sprintf("%s%d,", issueBo.Path, issueBo.Id))
		} else {
			trulyIssueBos = append(trulyIssueBos, issueBo)
		}
	}

	//移除所有不通过权限校验的以及他们的子任务
	allAuthedIssues := make([]*bo.IssueBo, 0)
	if len(notAuthPath) > 0 {
		for _, child := range issueAndChildren {
			if ok, _ := slice.Contain(noAuthIds, child.Id); ok {
				continue
			}
			pass := true
			if ok, _ := slice.Contain(noAuthIds, child.ParentId); ok {
				for _, s := range notAuthPath {
					if strings.Contains(child.Path, s) {
						pass = false
						break
					}
				}
			}

			if pass {
				allAuthedIssues = append(allAuthedIssues, child)
			}
		}
	} else {
		allAuthedIssues = issueAndChildren
	}
	tableId, tableIdErr := strconv.ParseInt(input.TableID, 10, 64)
	if tableIdErr != nil {
		log.Errorf("[DeleteIssueBatch] tableId参数错误：%v, inputTableId:%s", tableIdErr, input.TableID)
		return nil, errs.ParamError
	}
	deletedIssueIds, err3 := domain.DeleteIssueBatch(orgId, allAuthedIssues, input.Ids, currentUserId, sourceChannel, input.ProjectID, tableId)
	if err3 != nil {
		log.Error(err3)
		return nil, errs.BuildSystemErrorInfo(errs.IssueDomainError, err3)
	}

	// 事件上报
	asyn.Execute(func() {
		// 拿任务信息
		columnIds := []string{
			lc_helper.ConvertToFilterColumn(consts.BasicFieldTitle),
			lc_helper.ConvertToFilterColumn(consts.BasicFieldAppId),
			lc_helper.ConvertToFilterColumn(consts.BasicFieldProjectId),
			lc_helper.ConvertToFilterColumn(consts.BasicFieldTableId),
		}
		data, errSys := domain.GetIssueInfosMapLcByIssueIds(orgId, currentUserId, deletedIssueIds, columnIds...)
		if errSys != nil {
			return
		}

		issueProjectId := int64(0)
		openTraceId, _ := threadlocal.Mgr.GetValue(consts.JaegerContextTraceKey)
		openTraceIdStr := cast.ToString(openTraceId)

		for _, d := range data {
			dataId := cast.ToInt64(d[consts.BasicFieldId])
			issueId := cast.ToInt64(d[consts.BasicFieldIssueId])
			appId := cast.ToInt64(d[consts.BasicFieldAppId])
			projectId := cast.ToInt64(d[consts.BasicFieldProjectId])
			tableId := cast.ToInt64(d[consts.BasicFieldTableId])
			e := &commonvo.DataEvent{
				OrgId:     orgId,
				AppId:     appId,
				ProjectId: projectId,
				TableId:   tableId,
				DataId:    dataId,
				IssueId:   issueId,
				UserId:    currentUserId,
			}
			issueProjectId = projectId
			report.ReportDataEvent(msgPb.EventType_DataDeleted, openTraceIdStr, e)
		}

		domain.UpdateDingTopCard(orgId, issueProjectId)
	})

	resBo := bo.DeleteIssueBatchBo{
		NoAuthIssues:         notAuthPassIssues,
		RemainChildrenIssues: []*bo.IssueBo{},
	}

	res := &vo.DeleteIssueBatchResp{}
	for _, child := range allAuthedIssues {
		if ok, _ := slice.Contain(input.Ids, child.Id); ok {
			resBo.SuccessIssues = append(resBo.SuccessIssues, child)
		}
	}

	_ = copyer.Copy(resBo, res)

	return res, nil
}

//func MoveIssue(orgId, userId, appId int64, input vo.MoveIssueReq) (*vo.Void, errs.SystemErrorInfo) {
//	//获取原本的Issue
//	issueBo, errSys := GetIssueInfo(orgId, userId, input.ID)
//	if errSys != nil {
//		log.Errorf("[MoveIssue] GetIssueInfo err: %v issueId: %v", errSys, input.ID)
//		return nil, errSys
//	}
//
//	//查找目标任务以及所有的子任务
//	allMoveIssueIds, errSys := domain.GetIssueAndChildrenIds(orgId, input.ChildrenIds)
//	if errSys != nil {
//		log.Errorf("[MoveIssue] GetIssueAndChildrenIds err: %v", errSys)
//		return nil, errSys
//	}
//	allMoveIssueIds = append(allMoveIssueIds, input.ID)
//
//	issueAndChildrenIds, errSys := domain.GetIssueAndChildrenIds(orgId, []int64{input.ID})
//	if errSys != nil {
//		log.Errorf("[MoveIssue] GetIssueAndChildrenIds err: %v", errSys)
//		return nil, errSys
//	}
//
//	batchInput := vo.MoveIssueBatchReq{
//		Ids:           allMoveIssueIds,
//		FromProjectID: issueBo.ProjectId,
//		FromTableID:   cast.ToString(issueBo.TableId),
//		ProjectID:     input.ProjectID,
//		TableID:       input.TableID,
//		ChooseField:   input.ChooseField,
//	}
//	issueTitleMap := map[int64]string{}
//	issueTitleMap[input.ID] = input.Title
//	_, errSys = moveIssueBatch(orgId, userId, batchInput, issueAndChildrenIds, issueTitleMap)
//	return &vo.Void{ID: input.ID}, errSys
//}

// moveIssueBatchInner 批量移动任务
func moveIssueBatchInner(orgId, operatorId int64, allData []map[string]interface{}, allIssueIds []int64,
	fromTableId int64, fromProjectId int64, targetTableId int64, targetProjectId int64,
	fieldMapping map[string]string, isCreateMissingSelectOptions bool, minOrder float64) ([]int64, errs.SystemErrorInfo) {
	allMoveData := make([]map[string]interface{}, 0) // 除去鉴权失败的剩余所有符合条件的有效移动任务
	allMoveIds := make([]int64, 0)
	allRemainIds := make([]int64, 0)
	for _, data := range allData {
		issueId := cast.ToInt64(data[consts.BasicFieldIssueId])
		if ok, _ := slice.Contain(allIssueIds, issueId); ok {
			allMoveData = append(allMoveData, data)
			allMoveIds = append(allMoveIds, issueId)
		} else {
			allRemainIds = append(allRemainIds, issueId)
		}
	}
	if len(allMoveIds) == 0 {
		return allMoveIds, nil
	}

	var fromProjectBo *bo.ProjectBo
	var targetProjectBo *bo.ProjectBo
	var fromTable *projectvo.TableMetaData
	var targetTable *projectvo.TableMetaData
	var relatingColumnIds []string
	var errSys errs.SystemErrorInfo
	rootPath := "0,"

	isChangeProject := fromProjectId != targetProjectId
	isChangeTable := fromTableId != targetTableId
	if !isChangeProject && !isChangeTable {
		return []int64{}, nil
	}

	// from/target project
	if fromProjectId > 0 {
		fromProjectBo, errSys = domain.GetProjectSimple(orgId, fromProjectId)
		if errSys != nil {
			log.Error(errSys)
			return nil, errSys
		}
	}
	if targetProjectId == fromProjectId {
		targetProjectBo = fromProjectBo
	} else if targetProjectId != fromProjectId && targetProjectId > 0 {
		targetProjectBo, errSys = domain.GetProjectSimple(orgId, targetProjectId)
		if errSys != nil {
			log.Error(errSys)
			return nil, errSys
		}
	}

	// from/target table
	var fromTableSchema, targetTableSchema map[string]*projectvo.TableColumnData
	if fromTableId > 0 {
		fromTable, errSys = domain.GetTableByTableId(orgId, operatorId, fromTableId)
		if errSys != nil {
			log.Error(errSys)
			return nil, errSys
		}
	}
	fromTableSchema, errSys = domain.GetTableColumnsMap(orgId, fromTableId, nil, true)
	if errSys != nil {
		log.Error(errSys)
		return nil, errSys
	}
	for columnId, column := range fromTableSchema {
		if column.Field.Type == tablePb.ColumnType_relating.String() ||
			column.Field.Type == tablePb.ColumnType_baRelating.String() ||
			column.Field.Type == tablePb.ColumnType_singleRelating.String() {
			relatingColumnIds = append(relatingColumnIds, columnId)
		}
	}
	if targetTableId == fromTableId {
		targetTable = fromTable
		targetTableSchema = fromTableSchema
	} else {
		if targetTableId > 0 {
			targetTable, errSys = domain.GetTableByTableId(orgId, operatorId, targetTableId)
			if errSys != nil {
				log.Error(errSys)
				return nil, errSys
			}
		}

		targetTableSchema, errSys = domain.GetTableColumnsMap(orgId, targetTableId, nil, true)
		if errSys != nil {
			log.Error(errSys)
			return nil, errSys
		}
	}

	// target app
	var fromAppId int64
	var targetAppId int64
	var summaryAppId int64
	var targetProjectTypeId int64
	if fromProjectBo != nil {
		fromAppId = fromProjectBo.AppId
	}
	if targetProjectBo != nil {
		targetAppId = targetProjectBo.AppId
		targetProjectTypeId = targetProjectBo.ProjectTypeId
	}
	summaryAppId, errSys = domain.GetOrgSummaryAppId(orgId)
	if errSys != nil {
		log.Error(errSys)
		return nil, errSys
	}
	if fromAppId == 0 {
		fromAppId = summaryAppId
	}
	if targetAppId == 0 {
		targetAppId = summaryAppId
	}

	// 拿目标表格最大order
	var orderOffset float64
	if targetTableId > 0 {
		targetMaxOrder, errSys := domain.GetTableIssueMaxOrderOfTable(orgId, operatorId, targetTableId)
		if errSys != nil {
			log.Error(errSys)
			return nil, errSys
		}
		orderOffset = targetMaxOrder - minOrder
	}

	nowTimeStr := time.Now().Format(consts.AppTimeFormat)

	changeToRootIds := make([]int64, 0) // 需要转化为顶级父任务的
	targetChangePath := make(map[string]string)
	fromChangePath := make(map[string]string)

	for _, data := range allData {
		issueId := cast.ToInt64(data[consts.BasicFieldIssueId])
		parentId := cast.ToInt64(data[consts.BasicFieldParentId])
		path := cast.ToString(data[consts.BasicFieldPath])
		if ok, _ := slice.Contain(allMoveIds, issueId); ok {
			// 需要移动的任务
			// 是子任务
			if parentId != 0 {
				// 如果父任务不在需要移动的里面，该任务需要转为顶级任务并移动到新项目表
				if ok, _ = slice.Contain(allMoveIds, parentId); !ok {
					changeToRootIds = append(changeToRootIds, issueId)
					targetChangePath[fmt.Sprintf("%s%d,", path, issueId)] = fmt.Sprintf("0,%d,", issueId)
				}
			}
		} else {
			// 不需要移动的任务
			// 是子任务
			if parentId != 0 {
				// 如果父任务在需要移动的里面，该任务需要转为顶级任务并保留在原项目表
				if ok, _ = slice.Contain(allMoveIds, parentId); ok {
					changeToRootIds = append(changeToRootIds, issueId)
					fromChangePath[fmt.Sprintf("%s%d,", path, issueId)] = fmt.Sprintf("0,%d,", issueId)
				}
			}
		}
	}

	var updateBatchReqs []*formvo.LessUpdateIssueBatchReq

	// 需要移动的任务涉及的关联引用的任务ID，刷新updateTime...
	if len(relatingColumnIds) > 0 {
		var allRelatingIssueIds []int64
		for _, data := range allData {
			issueId := cast.ToInt64(data[consts.BasicFieldIssueId])
			if ok, _ := slice.Contain(allMoveIds, issueId); ok {
				// 需要移动的任务
				for _, columnId := range relatingColumnIds {
					if v, ok := data[columnId]; ok {
						oldRelating := &bo.RelatingIssue{}
						jsonx.Copy(v, oldRelating)
						allRelatingIssueIds = append(allRelatingIssueIds, slice2.StringToInt64Slice(oldRelating.LinkTo)...)
						allRelatingIssueIds = append(allRelatingIssueIds, slice2.StringToInt64Slice(oldRelating.LinkFrom)...)
					}
				}
			}
		}
		allRelatingIssueIds = slice.SliceUniqueInt64(allRelatingIssueIds)
		var relatingIssueIds []int64
		for _, id := range allRelatingIssueIds {
			if ok, _ := slice.Contain(allMoveIds, id); !ok {
				relatingIssueIds = append(relatingIssueIds, id)
			}
		}
		if len(relatingIssueIds) > 0 {
			updateBatchReqs = append(updateBatchReqs, &formvo.LessUpdateIssueBatchReq{
				OrgId:  orgId,
				AppId:  summaryAppId, // 目标项目appId
				UserId: operatorId,
				Condition: vo.LessCondsData{
					Type: consts.ConditionAnd,
					Conds: []*vo.LessCondsData{
						{
							Type:   consts.ConditionEqual,
							Value:  orgId,
							Column: lc_helper.ConvertToCondColumn(consts.BasicFieldOrgId),
						},
						{
							Type:   consts.ConditionIn,
							Values: relatingIssueIds,
							Column: lc_helper.ConvertToCondColumn(consts.BasicFieldIssueId),
						},
					},
				},
				Sets: []datacenter.Set{
					{
						Column:          lc_helper.ConvertToCondColumn(consts.BasicFieldUpdateTime),
						Value:           "NOW()",
						Type:            consts.SetTypeNormal,
						Action:          consts.SetActionSet,
						WithoutPretreat: true,
					},
				},
			})
		}
	}

	// 刷parentId path
	if len(changeToRootIds) > 0 {
		// 移动到目标项目表里的任务，刷path
		for oldPath, newPath := range targetChangePath {
			updateBatchReqs = append(updateBatchReqs, &formvo.LessUpdateIssueBatchReq{
				OrgId:  orgId,
				AppId:  targetAppId, // 目标项目appId
				UserId: operatorId,
				Condition: vo.LessCondsData{
					Type: consts.ConditionAnd,
					Conds: []*vo.LessCondsData{
						{
							Type:   consts.ConditionEqual,
							Value:  orgId,
							Column: lc_helper.ConvertToCondColumn(consts.BasicFieldOrgId),
						},
						{
							Type:   consts.ConditionIn,
							Values: allMoveIds,
							Column: lc_helper.ConvertToCondColumn(consts.BasicFieldIssueId),
						},
					},
				},
				Sets: []datacenter.Set{
					{
						Column:          lc_helper.ConvertToCondColumn(consts.BasicFieldPath),
						Value:           fmt.Sprintf("regexp_replace(\"%s\",'%s','%s')", consts.BasicFieldPath, oldPath, newPath),
						Type:            consts.SetTypeNormal,
						Action:          consts.SetActionSet,
						WithoutPretreat: true,
					},
				},
			})
		}
		// 留在原项目表里的任务，刷path
		for oldPath, newPath := range fromChangePath {
			updateBatchReqs = append(updateBatchReqs, &formvo.LessUpdateIssueBatchReq{
				OrgId:  orgId,
				AppId:  fromAppId, // 原项目appId
				UserId: operatorId,
				Condition: vo.LessCondsData{
					Type: consts.ConditionAnd,
					Conds: []*vo.LessCondsData{
						{
							Type:   consts.ConditionEqual,
							Value:  orgId,
							Column: lc_helper.ConvertToCondColumn(consts.BasicFieldOrgId),
						},
						{
							Type:   consts.ConditionIn,
							Values: allRemainIds,
							Column: lc_helper.ConvertToCondColumn(consts.BasicFieldIssueId),
						},
					},
				},
				Sets: []datacenter.Set{
					{
						Column:          lc_helper.ConvertToCondColumn(consts.BasicFieldPath),
						Value:           fmt.Sprintf("regexp_replace(\"%s\",'%s','%s')", consts.BasicFieldPath, oldPath, newPath),
						Type:            consts.SetTypeNormal,
						Action:          consts.SetActionSet,
						WithoutPretreat: true,
					},
				},
			})
		}
	}

	// 字段映射
	validFieldMapping := domain.ValidFieldMapping(orgId, fromProjectId, fromTableId, targetProjectId, targetTableId,
		fromTableSchema, targetTableSchema, fieldMapping)
	if targetProjectId > 0 {
		validFieldMapping[consts.BasicFieldOrder] = []string{consts.BasicFieldOrder} // 移动任务的时候 order需要保留
	}
	selectOptions := make(map[string]map[string]*projectvo.ColumnSelectOption)
	var allNewMoveData []map[string]interface{}
	if len(validFieldMapping) > 0 {
		for _, data := range allMoveData {
			newData := make(map[string]interface{})
			copyer.Copy(data, &newData)
			domain.ApplyFieldMappingCreateMissingSelectOptions(data, newData, validFieldMapping, selectOptions,
				fromTableSchema, targetTableSchema, targetProjectTypeId, isCreateMissingSelectOptions)
			// 更新order值，使之成为目标表格的顶部
			if targetTableId > 0 && orderOffset > 0 {
				order := cast.ToFloat64(newData[consts.BasicFieldOrder])
				newData[consts.BasicFieldOrder] = order + orderOffset
			}
			newData[consts.BasicFieldId] = data[consts.BasicFieldIssueId] // 设置更新主键
			allNewMoveData = append(allNewMoveData, newData)
		}
	}

	// 不移动项目表的任务
	var noMoveTableUpdateFrom []map[string]interface{}
	for _, id := range changeToRootIds {
		noMoveTableUpdateFrom = append(noMoveTableUpdateFrom, map[string]interface{}{
			consts.BasicFieldIssueId:  id,
			consts.BasicFieldParentId: 0,
			consts.BasicFieldPath:     rootPath,
			consts.BasicFieldMoveTime: nowTimeStr,
			consts.BasicFieldUpdator:  operatorId,
		})
	}
	// 移动项目表的任务
	var moveTableUpdateForm []map[string]interface{}
	for _, id := range allMoveIds {
		moveTableUpdateForm = append(moveTableUpdateForm, map[string]interface{}{
			consts.BasicFieldIssueId:   id,
			consts.BasicFieldMoveTime:  nowTimeStr,
			consts.BasicFieldUpdator:   operatorId,
			consts.BasicFieldAppId:     cast.ToString(targetAppId),
			consts.BasicFieldProjectId: targetProjectId,
			consts.BasicFieldTableId:   cast.ToString(targetTableId),
		})
	}

	transErr := mysql.TransX(func(tx sqlbuilder.Tx) error {
		// 无码更新: 原项目表中刷新parentId/path等信息
		if len(noMoveTableUpdateFrom) > 0 {
			resp := formfacade.LessUpdateIssue(formvo.LessUpdateIssueReq{
				AppId:   fromAppId,
				OrgId:   orgId,
				UserId:  operatorId,
				TableId: fromTableId,
				Form:    noMoveTableUpdateFrom,
			})
			if resp.Failure() {
				log.Error(resp.Error())
				return resp.Error()
			}
		}
		// 无码更新：目标项目表刷projectId/tableId，这样update最终在正确的table表头中执行
		if len(moveTableUpdateForm) > 0 {
			resp := formfacade.LessUpdateIssue(formvo.LessUpdateIssueReq{
				AppId:   targetAppId,
				OrgId:   orgId,
				UserId:  operatorId,
				TableId: targetTableId,
				Form:    moveTableUpdateForm,
			})
			if resp.Failure() {
				log.Error(resp.Error())
				return resp.Error()
			}
		}
		// 重新设置数据字段映射之后的值
		if len(allNewMoveData) > 0 {
			req := &projectvo.BatchUpdateIssueReqInnerVo{
				OrgId:         orgId,
				UserId:        operatorId,
				AppId:         targetAppId,
				ProjectId:     targetProjectId,
				TableId:       targetTableId,
				Data:          allNewMoveData,
				SelectOptions: selectOptions,
			}
			if len(validFieldMapping) > 1 {
				req.IsFullSet = true
			}
			errSys = SyncBatchUpdateIssue(req, true, &projectvo.TriggerBy{TriggerBy: consts.TriggerByMove})
			if errSys != nil {
				log.Error(errSys)
				//return errSys // 已经移动过去了，这里如果失败也只能忽略，相当于移动成功但字段映射值未成功
			}
		}
		// 无码批量更新
		for _, req := range updateBatchReqs {
			batchResp := formfacade.LessUpdateIssueBatchRaw(req)
			if batchResp.Failure() {
				log.Error(batchResp.Error())
				return batchResp.Error()
			}
		}

		// 删除工时
		if _, ok := validFieldMapping[consts.BasicFieldWorkHour]; !ok {
			errSys = domain.UpdateWorkHourProjectId(orgId, fromProjectId, targetProjectId, allMoveIds, true)
			if errSys != nil {
				log.Error(errSys)
				return errSys
			}
		}
		return nil
	})
	if transErr != nil {
		log.Error(transErr)
		return nil, errs.MysqlOperateError
	}

	// 保存字段映射
	domain.SaveFieldMappingCache(orgId, fromTableId, targetTableId, fieldMapping)

	asyn.Execute(func() {
		if isChangeProject {
			// 更新飞书日历
			calendarErr := domain.SwitchCalendar(orgId, fromProjectId, allMoveIds, operatorId, targetProjectId)
			if calendarErr != nil {
				log.Error(calendarErr)
				return
			}
		}
		if isChangeProject || isChangeTable {
			pushType := consts.PushTypeUpdateIssueProjectTable
			for _, data := range allMoveData {
				issueBo, errSys := domain.ConvertIssueDataToIssueBo(data)
				if errSys != nil {
					log.Error(errSys)
					continue
				}
				var changeList []bo.TrendChangeListBo
				var oldProjectName, newProjectName string
				if fromProjectBo != nil {
					oldProjectName = fromProjectBo.Name
				}
				if targetProjectBo != nil {
					newProjectName = targetProjectBo.Name
				}
				changeList = append(changeList, bo.TrendChangeListBo{
					Field:     consts.BasicFieldProjectId,
					FieldName: consts.Project,
					OldValue:  oldProjectName,
					NewValue:  newProjectName,
				})
				var oldTableName, newTableName string
				if fromTable != nil {
					oldTableName = fromTable.Name
				}
				if targetTable != nil {
					newTableName = targetTable.Name
				}
				changeList = append(changeList, bo.TrendChangeListBo{
					Field:     consts.BasicFieldTableId,
					FieldName: consts.Table,
					OldValue:  oldTableName,
					NewValue:  newTableName,
				})

				oldValue := bo.ProjectTableBo{}
				oldValue.ProjectId = fromProjectId
				oldValue.TableId = fromTableId

				newValue := bo.ProjectTableBo{}
				newValue.ProjectId = targetProjectId
				newValue.TableId = targetTableId

				issueTrendsBo := &bo.IssueTrendsBo{
					PushType:      pushType,
					OrgId:         issueBo.OrgId,
					OperatorId:    operatorId,
					DataId:        issueBo.DataId,
					IssueId:       issueBo.Id,
					ParentIssueId: issueBo.ParentId,
					ProjectId:     targetProjectId,
					TableId:       targetTableId,
					PriorityId:    issueBo.PriorityId,
					IssueTitle:    issueBo.Title,
					ParentId:      issueBo.ParentId,
					OldValue:      json.ToJsonIgnoreError(oldValue),
					NewValue:      json.ToJsonIgnoreError(newValue),
					Ext: bo.TrendExtensionBo{
						ObjName:    issueBo.Title,
						ChangeList: changeList,
					},
				}
				domain.PushIssueTrends(issueTrendsBo)

				//asyn.Execute(func() {
				//	PushIssueThirdPlatformNotice(issueTrendsBo)
				//})
				//asyn.Execute(func() {
				//	//推送群聊卡片
				//	PushInfoToChat(issueBo.OrgId, targetProjectId, issueTrendsBo)
				//})
			}
		}
	})
	return allMoveIds, nil
}

func MoveIssueBatch(orgId, userId int64, input *projectvo.LcMoveIssuesInput) ([]int64, errs.SystemErrorInfo) {
	allIds := input.IssueIds
	var errSys errs.SystemErrorInfo

	//查找目标任务以及所有的子任务
	allIdsWithChildren, errSys := domain.GetIssueAndChildrenIds(orgId, input.IssueIds)
	if errSys != nil {
		log.Errorf("[MoveIssueBatch] GetIssueAndChildrenIds err: %v", errSys)
		return nil, errSys
	}

	if input.IsIncludeChildren {
		allIds = allIdsWithChildren
	}

	// 获取所有原始任务数据
	allData, errSys := domain.GetIssueInfosMapLcByIssueIds(orgId, userId, allIdsWithChildren)
	if errSys != nil {
		log.Errorf("[MoveIssueBatch] GetIssueInfosMapLcByIssueIds err: %v", errSys)
		return nil, errSys
	}
	if len(allData) == 0 {
		return nil, errs.IssueNotExist
	}
	var allIssues []*bo.IssueBo
	for _, data := range allData {
		issue, errSys := domain.ConvertIssueDataToIssueBo(data)
		if errSys != nil {
			log.Errorf("[MoveIssueBatch] ConvertIssueDataToIssueBo failed, orgId: %v, data: %v, err: %v",
				orgId, json.ToJsonIgnoreError(data), errSys)
			return nil, errSys
		}
		allIssues = append(allIssues, issue)
	}

	var projectId int64
	if input.AppId > 0 {
		projectId, errSys = domain.GetProjectIdByAppId(orgId, input.AppId)
		if errSys != nil {
			log.Errorf("[MoveIssueBatch] GetProjectIdByAppId, orgId: %v, appId: %v, err: %v", orgId, input.AppId, errSys)
			return nil, errSys
		}
	}

	fromProjectId := cast.ToInt64(allData[0][consts.BasicFieldProjectId])
	fromTableId := cast.ToInt64(allData[0][consts.BasicFieldTableId])
	targetProjectId := projectId
	targetTableId := input.TableId
	if fromProjectId == targetProjectId && fromTableId == targetTableId {
		log.Infof("[MoveIssueBatch] same project/table no need to move, orgId: %v, projectId: %v, tableId: %v",
			orgId, targetProjectId, targetTableId)
		return nil, nil
	}

	var allNeedAuthIssues []*bo.IssueBo
	var allAuthedData []map[string]interface{}
	var noAuthPath []string
	var noAuthIssueIds []int64

	for _, issue := range allIssues {
		// 如果父任务为0，或者父任务不在传入的任务id里面，这些任务将作为父任务进行权限校验
		if issue.ParentId == int64(0) {
			allNeedAuthIssues = append(allNeedAuthIssues, issue)
		} else {
			if ok, _ := slice.Contain(input.IssueIds, issue.ParentId); !ok {
				allNeedAuthIssues = append(allNeedAuthIssues, issue)
			}
		}
	}

	// 权限判断
	if fromProjectId != targetProjectId && targetProjectId > 0 {
		errSys = AuthProject(orgId, userId, targetProjectId, consts.RoleOperationPathOrgProIssueT, consts.OperationProIssue4Create)
		if errSys != nil {
			log.Errorf("[MoveIssueBatch] AuthProject err: %v", errSys)
			return nil, errs.BuildSystemErrorInfo(errs.Unauthorized, errSys)
		}
	}
	// TODO: 优化成批量
	for _, issueBo := range allNeedAuthIssues {
		errSys = domain.AuthIssueWithAppId(orgId, userId, issueBo, consts.RoleOperationPathOrgProIssueT, consts.OperationProIssue4Modify, 0)
		if errSys != nil {
			log.Infof("[MoveIssueBatch] AuthIssueWithAppId err: %v", errSys)
			noAuthIssueIds = append(noAuthIssueIds, issueBo.Id)
			noAuthPath = append(noAuthPath, fmt.Sprintf("%s%d,", issueBo.Path, issueBo.Id))
		}
	}

	// 移除所有不通过权限校验的以及他们的子任务
	if len(noAuthPath) > 0 {
		for _, data := range allData {
			issueId := cast.ToInt64(data[consts.BasicFieldIssueId])
			path := cast.ToString(data[consts.BasicFieldPath])

			if ok, _ := slice.Contain(noAuthIssueIds, issueId); ok {
				continue
			}
			pass := true
			for _, s := range noAuthPath {
				if strings.Contains(path, s) {
					pass = false
					break
				}
			}
			if pass {
				allAuthedData = append(allAuthedData, data)
			}
		}
	} else {
		allAuthedData = allData
	}

	// 找出最小order
	gotMinOrder := false
	var minOrder float64
	for _, data := range allAuthedData {
		order := cast.ToFloat64(data[consts.BasicFieldOrder])
		if !gotMinOrder {
			minOrder = order
			gotMinOrder = true
		} else if order < minOrder {
			minOrder = order
		}
	}

	// 更新任务的项目对象类型
	allMovedIssueIds, errSys := moveIssueBatchInner(orgId, userId, allAuthedData, allIds,
		fromTableId, fromProjectId, targetTableId, targetProjectId, input.FieldMapping, input.IsCreateMissingSelectOptions, minOrder)
	if errSys != nil {
		log.Infof("[MoveIssueBatch] moveIssueBatchInner err: %v", errSys)
		return nil, errSys
	}

	// 上报事件
	asyn.Execute(func() {
		ReportMoveEvents(orgId, userId, allIssues, allMovedIssueIds, fromTableId, targetTableId)
	})

	return allMovedIssueIds, nil
}

func GetFieldMapping(orgId, fromTableId, toTableId int64) map[string]string {
	return domain.GetFieldMappingCache(orgId, fromTableId, toTableId)
}

func GetIssueInfo(orgId, userId int64, id int64) (*bo.IssueBo, errs.SystemErrorInfo) {
	bos, err := domain.GetIssueInfosLc(orgId, userId, []int64{id})
	if err != nil {
		return nil, err
	}
	if len(bos) == 0 {
		return nil, errs.IssueNotExist
	}
	return bos[0], nil
}

func GetIssueInfoList(orgId, userId int64, ids []int64) ([]vo.Issue, errs.SystemErrorInfo) {
	result, err := domain.GetIssueInfosLc(orgId, userId, ids)
	if err != nil {
		return nil, err
	}
	issueResp := make([]vo.Issue, 0, len(result))
	err1 := copyer.Copy(result, &issueResp)
	if err1 != nil {
		log.Errorf("copyer.Copy: %q\n", err1)
		return nil, errs.BuildSystemErrorInfo(errs.ObjectCopyError, err1)
	}
	// 赋值负责人 ids
	for i, issue := range result {
		issueResp[i].Owners = issue.OwnerIdI64
	}

	return issueResp, nil
}

func GetIssueInfoListByDataIds(orgId, userId int64, dataIds []int64) ([]vo.Issue, errs.SystemErrorInfo) {
	result, err := domain.GetIssueInfosLcByDataIds(orgId, userId, dataIds)
	if err != nil {
		return nil, err
	}
	issueResp := &[]vo.Issue{}
	err1 := copyer.Copy(result, issueResp)
	if err1 != nil {
		log.Errorf("copyer.Copy: %q\n", err1)
		return nil, errs.BuildSystemErrorInfo(errs.ObjectCopyError, err1)
	}

	return *issueResp, nil
}

func copyIssueBatchInner(orgId int64, userId int64, newProjectId, newTableId int64, allIssueIds []int64,
	allCopyData []map[string]interface{}, fieldMapping map[string]string, triggerBy *projectvo.TriggerBy,
	isSyncCreate, isStaticCopy, isCreateMissingSelectOptions bool, beforeDataId, afterDataId int64) (
	[]map[string]interface{}, map[string]*uservo.MemberDept, map[string]map[string]interface{}, errs.SystemErrorInfo) {
	// 获取project
	project, errSys := domain.GetProjectSimple(orgId, newProjectId)
	if errSys != nil {
		log.Errorf("[copyIssueBatchInner] GetProjectSimple err: %v", errSys)
		return nil, nil, nil, errSys
	}

	// 所有数据按order降序排列
	sort.Slice(allCopyData, func(i, j int) bool {
		orderI := cast.ToInt64(allCopyData[i][consts.BasicFieldOrder])
		orderJ := cast.ToInt64(allCopyData[j][consts.BasicFieldOrder])
		return orderI > orderJ
	})

	// 生成issue ids
	issueIds, errSys := idfacade.ApplyMultiplePrimaryIdRelaxed(consts.TableIssue, len(allCopyData))
	if errSys != nil {
		log.Errorf("[copyIssueBatchInner] ApplyMultiplePrimaryIdRelaxed failed, count: %v, err: %v",
			len(allCopyData), errSys)
		return nil, nil, nil, errSys
	}

	// 生成issue codes
	preCode := consts.NoProjectPreCode
	if project.PreCode != "" {
		preCode = project.PreCode
	} else {
		preCode = fmt.Sprintf("$%d", newProjectId)
	}
	issueCodes, errSys := idfacade.ApplyMultipleIdRelaxed(orgId, preCode, "", int64(len(allCopyData)))
	if errSys != nil {
		log.Errorf("[copyIssueBatchInner] ApplyMultipleIdRelaxed failed, count: %v, err: %v",
			len(allCopyData), errSys)
		return nil, nil, nil, errSys
	}

	// 找出父子任务关系
	parentChildrenMapping := map[int64][]int64{}       // 父任务->子任务列表
	allData := map[int64]map[string]interface{}{}      // 任务id->任务数据
	allParentData := make([]map[string]interface{}, 0) // 父任务
	data := make([]map[string]interface{}, 0)          // 待创建的任务数据
	issueIdMapping := make(map[int64]int64, 0)         // oldIssueId->newIssueId
	//isSameProject := false
	isSameTable := false
	var oldProjectId, oldTableId int64
	for i, d := range allCopyData {
		issueId := cast.ToInt64(d[consts.BasicFieldIssueId])
		parentId := cast.ToInt64(d[consts.BasicFieldParentId])
		projectId := cast.ToInt64(d[consts.BasicFieldProjectId])
		tableId := cast.ToInt64(d[consts.BasicFieldTableId])
		oldProjectId = projectId
		oldTableId = tableId

		allData[issueId] = d

		if i == 0 {
			//isSameProject = projectId == newProjectId
			isSameTable = tableId == newTableId
		}

		// 赋予新issueId/code
		newIssueId := issueIds.Ids[i].Id
		d[consts.TempFieldIssueId] = newIssueId
		issueIdMapping[issueId] = newIssueId
		d[consts.TempFieldCode] = issueCodes.Ids[i].Code

		// 与父任务断开的子任务也视为父任务
		if has, _ := slice.Contain(allIssueIds, parentId); !has {
			// 父任务
			allParentData = append(allParentData, d)

			if !isSameTable {
				d[consts.BasicFieldParentId] = 0
				d[consts.BasicFieldPath] = "0,"
			}
		} else {
			// 子任务
			parentChildrenMapping[parentId] = append(parentChildrenMapping[parentId], issueId)
		}
	}

	// 获取表头
	oldTableColumns, errSys := domain.GetTableColumnsMap(orgId, oldTableId, nil, true)
	if errSys != nil {
		log.Errorf("[copyIssueBatchInner] GetTableColumnsMap failed, org:%d table:%d, err: %v",
			orgId, oldTableId, errSys)
		return nil, nil, nil, errSys
	}
	newTableColumns := oldTableColumns
	if !isSameTable {
		newTableColumns, errSys = domain.GetTableColumnsMap(orgId, newTableId, nil, true)
		if errSys != nil {
			log.Errorf("[copyIssueBatchInner] GetTableColumnsMap failed, org:%d table:%d, err: %v",
				orgId, newTableId, errSys)
			return nil, nil, nil, errSys
		}
	}
	// 找出关联/前后置/单向关联列
	relatingColumnIds := []string{consts.BasicFieldRelating, consts.BasicFieldBaRelating}
	for columnId, column := range newTableColumns {
		if columnId != consts.BasicFieldRelating && columnId != consts.BasicFieldBaRelating {
			if column.Field.Type == consts.LcColumnFieldTypeRelating ||
				column.Field.Type == consts.LcColumnFieldTypeSingleRelating {
				relatingColumnIds = append(relatingColumnIds, columnId)
			}
		}
	}

	// 检查有效的字段映射
	validFieldMapping := domain.ValidFieldMapping(orgId, oldProjectId, oldTableId, newProjectId, newTableId,
		oldTableColumns, newTableColumns, fieldMapping)

	var totalCount int64 // 任务总数（包括子任务）
	selectOptions := make(map[string]map[string]*projectvo.ColumnSelectOption)

	// 处理父任务
	for _, d := range allParentData {
		issueId := cast.ToInt64(d[consts.BasicFieldIssueId])
		parentData := domain.BuildIssueDataByCopy(d, nil, validFieldMapping, issueIdMapping,
			relatingColumnIds, selectOptions, oldTableColumns, newTableColumns, project.ProjectTypeId,
			isStaticCopy, isCreateMissingSelectOptions, false, false)

		if childrenIds, ok := parentChildrenMapping[issueId]; ok {
			// 处理子任务（递归）
			count, err := domain.BuildIssueChildDataByCopy(parentData, childrenIds, allData, parentChildrenMapping,
				nil, validFieldMapping, issueIdMapping, relatingColumnIds, selectOptions, oldTableColumns, newTableColumns,
				project.ProjectTypeId, isStaticCopy, isCreateMissingSelectOptions, false, false)
			if err != nil {
				log.Errorf("[copyIssueBatchInner] BuildIssueChildDataByCopy err: %v", err)
				return nil, nil, nil, err
			}
			totalCount += count
		}

		totalCount += 1
		data = append(data, parentData)
	}

	//// TODO: 工时 - 暂不能支持无码化复制
	//if util.FieldInUpdate(input.ChooseField, consts.BasicFieldWorkHour) {
	//	// 查出原来工时记录，传给createIssue创建
	//	issueWorkHours, errSys := domain.GetCopyIssueWorkHours(orgId, input.ProjectID, input.OldIssueID)
	//	if errSys != nil {
	//		log.Errorf("[copyIssue]GetCopyIssueWorkHours err:%v", errSys)
	//		return nil, errSys
	//	}
	//	createIssueReq.LessCreateIssueReq[consts.CopyIssueWorkHour] = issueWorkHours
	//	// 无码更新值
	//	createIssueReq.LessCreateIssueReq[consts.BasicFieldWorkHour] = data[consts.BasicFieldWorkHour]
	//}

	// 保存字段映射
	domain.SaveFieldMappingCache(orgId, oldTableId, newTableId, fieldMapping)

	if !isSyncCreate {
		// 异步批量创建
		req := &projectvo.BatchCreateIssueReqVo{
			OrgId:  orgId,
			UserId: userId,
			Input: &projectvo.BatchCreateIssueInput{
				AppId:         project.AppId,
				ProjectId:     project.Id,
				TableId:       newTableId,
				Data:          data,
				BeforeDataId:  beforeDataId,
				AfterDataId:   afterDataId,
				IsIdGenerated: true,
				SelectOptions: selectOptions,
			},
		}
		AsyncBatchCreateIssue(req, true, triggerBy, "")
		return nil, nil, nil, nil
	} else {
		// 同步批量创建
		req := &projectvo.BatchCreateIssueReqVo{
			OrgId:  orgId,
			UserId: userId,
			Input: &projectvo.BatchCreateIssueInput{
				AppId:         project.AppId,
				ProjectId:     project.Id,
				TableId:       newTableId,
				Data:          data,
				BeforeDataId:  beforeDataId,
				AfterDataId:   afterDataId,
				IsIdGenerated: true,
				SelectOptions: selectOptions,
			},
		}
		return SyncBatchCreateIssue(req, true, triggerBy)
	}
}

func CopyIssueBatch(orgId int64, userId int64, input *projectvo.LcCopyIssuesInput) (
	[]map[string]interface{}, map[string]*uservo.MemberDept, map[string]map[string]interface{}, errs.SystemErrorInfo) {
	projectId, errSys := domain.GetProjectIdByAppId(orgId, input.AppId)
	if errSys != nil {
		log.Errorf("[CopyIssueBatch] GetProjectIdByAppId, orgId: %v, appId: %v, err: %v", orgId, input.AppId, errSys)
		return nil, nil, nil, errSys
	}
	// 判断权限 - 在目标项目创建任务的权限
	errSys = domain.AuthProject(orgId, userId, projectId, consts.RoleOperationPathOrgProIssueT, consts.OperationProIssue4Create)
	if errSys != nil {
		log.Errorf("[CopyIssueBatch] AuthProject err: %v", errSys)
		return nil, nil, nil, errSys
	}

	allIds := input.IssueIds

	// 获取所有子任务id + 子任务的子任务id
	if input.IsIncludeChildren {
		allIds, errSys = domain.GetIssueAndChildrenIds(orgId, allIds)
		if errSys != nil {
			log.Errorf("[CopyIssueBatch] GetIssueAndChildrenIds err: %v", errSys)
			return nil, nil, nil, errSys
		}
	}

	// 判断权限 - 付费任务数
	errSys = domain.AuthPayTask(orgId, consts.FunctionTaskLimit, len(allIds))
	if errSys != nil {
		log.Errorf("[CopyIssueBatch] AuthPayTask err: %v", errSys)
		return nil, nil, nil, errSys
	}

	// 获取所有原始任务数据
	allCopyData, errSys := domain.GetIssueInfosMapLcByIssueIds(orgId, userId, allIds)
	if errSys != nil {
		log.Errorf("[CopyIssueBatch] GetIssueInfosMapLcByIssueIds err: %v", errSys)
		return nil, nil, nil, errSys
	}

	return copyIssueBatchInner(orgId, userId, projectId, input.TableId, allIds, allCopyData, input.FieldMapping,
		&projectvo.TriggerBy{TriggerBy: consts.TriggerByCopy},
		len(input.IssueIds) == 1, input.IsStaticCopy, input.IsCreateMissingSelectOptions, input.BeforeDataId, input.AfterDataId)
}

func CopyIssueBatchWithData(orgId int64, userId int64, newProjectId, newTableId int64, allIssueIds []int64,
	allCopyData []map[string]interface{}, triggerBy *projectvo.TriggerBy,
	isSyncCreate, isStaticCopy, isCreateMissingSelectOptions, isInnerSuper bool) (
	[]map[string]interface{}, map[string]*uservo.MemberDept, map[string]map[string]interface{}, errs.SystemErrorInfo) {
	if len(allIssueIds) == 0 {
		return nil, nil, nil, errs.IssueIdsNotBeenChosen
	}

	// 判断权限 - 在目标项目创建任务的权限
	if !isInnerSuper {
		errSys := domain.AuthProject(orgId, userId, newProjectId, consts.RoleOperationPathOrgProIssueT, consts.OperationProIssue4Create)
		if errSys != nil {
			log.Errorf("[CopyIssueBatchWithData] AuthProject err: %v", errSys)
			return nil, nil, nil, errSys
		}
	}

	// 判断付费限制 - 付费任务数
	errSys := domain.AuthPayTask(orgId, consts.FunctionTaskLimit, len(allIssueIds))
	if errSys != nil {
		log.Errorf("[CopyIssueBatch] AuthPayTask err: %v", errSys)
		return nil, nil, nil, errSys
	}

	// TODO: field mapping from input
	return copyIssueBatchInner(orgId, userId, newProjectId, newTableId, allIssueIds, allCopyData,
		nil, triggerBy, isSyncCreate, isStaticCopy, isCreateMissingSelectOptions, 0, 0)
}

// ViewAuditIssue 查看审核任务
func ViewAuditIssue(orgId, userId, issueId int64) errs.SystemErrorInfo {
	issueInfos, sysErr := domain.GetIssueInfosLc(orgId, userId, []int64{issueId})
	if sysErr != nil {
		log.Error(sysErr)
		return sysErr
	}
	if len(issueInfos) < 1 {
		log.Errorf("[ViewAuditIssue] not found issue issueId: %v", issueId)
		return errs.IssueNotExist
	}
	issueInfo := issueInfos[0]

	// 判断当前操作人是否是任务确认人
	if ok, _ := slice.Contain(issueInfo.AuditorIdsI64, userId); !ok {
		return nil
	}

	summaryAppId, sysErr := domain.GetOrgSummaryAppId(orgId)
	if sysErr != nil {
		log.Error(sysErr)
		return sysErr
	}

	auditStatus, ok := issueInfo.AuditStatusDetail[cast.ToString(userId)]
	if !ok || auditStatus == consts.AuditStatusNotView {
		// 更新无码
		sysErr = domain.UpdateLcIssueAuditStatusDetailByUser(orgId, summaryAppId, userId, issueId, consts.AuditStatusView)
		if sysErr != nil {
			log.Error(sysErr)
			return sysErr
		}
	}

	return nil
}

// WithdrawIssue 撤回审核任务
func WithdrawIssue(orgId, userId int64, issueId int64) errs.SystemErrorInfo {
	issueInfos, sysErr := domain.GetIssueInfosLc(orgId, userId, []int64{issueId})
	if sysErr != nil {
		log.Error(sysErr)
		return sysErr
	}
	if len(issueInfos) < 1 {
		log.Errorf("[WithdrawIssue] not found issue issueId: %v", issueId)
		return errs.IssueNotExist
	}
	issueInfo := issueInfos[0]

	// 判断当前操作人是否是任务负责人
	if ok, _ := slice.Contain(issueInfo.OwnerIdI64, userId); !ok {
		return errs.OnlyOwnerCanWithdrawIssue
	}

	// 未完成的任务无需撤回
	if issueInfo.IssueStatusType != consts.StatusTypeComplete {
		return errs.NotFinishIssue
	}

	// 已经通过审核了
	if issueInfo.AuditStatus == consts.AuditStatusPass {
		return errs.IssueIsAuditPass
	}

	allStatusMap, sysErr := domain.GetIssueAllStatus(orgId, []int64{issueInfo.ProjectId}, []int64{issueInfo.TableId})
	if sysErr != nil {
		log.Error(sysErr)
		return sysErr
	}
	issueStatusId := issueInfo.Status
	if allStatus, ok := allStatusMap[issueInfo.TableId]; ok {
		for _, statusInfo := range allStatus {
			if statusInfo.Type == consts.StatusTypeNotStart {
				issueStatusId = statusInfo.ID
			}
		}
	} else {
		return errs.BuildSystemErrorInfo(errs.IssueStatusNotExist)
	}

	//transErr := mysql.TransX(func(tx sqlbuilder.Tx) error {
	//// 更新无码
	//sysErr = domain.UpdateLcIssueAuditStatusDetailByUser(orgId, appId, userId, issueId, consts.AuditStatusNotView)
	//if sysErr != nil {
	//	log.Error(sysErr)
	//	return sysErr
	//}

	//_, err := mysql.TransUpdateSmartWithCond(tx, consts.TableIssueRelation, db.Cond{
	//	consts.TcRelationType: consts.IssueRelationTypeAuditor,
	//	consts.TcOrgId:        orgId,
	//	consts.TcIssueId:      issueId,
	//}, mysql.Upd{
	//	consts.TcStatus:  consts.AuditStatusNotView,
	//	consts.TcUpdator: userId,
	//})
	//if err != nil {
	//	log.Error(err)
	//	return err
	//}

	// 清空审批状态
	d := make(map[string]interface{})
	d[consts.BasicFieldId] = issueId
	d[consts.BasicFieldAuditStatus] = consts.AuditStatusNotView
	d[consts.BasicFieldAuditStatusDetail] = map[string]int{}
	d[consts.BasicFieldIssueStatus] = issueStatusId
	d[consts.BasicFieldIssueStatusType] = consts.StatusTypeNotStart //撤回到未完成状态

	// 更新无码
	pushType := consts.PushTypeWithdrawIssue
	sysErr = BatchUpdateIssue(&projectvo.BatchUpdateIssueReqInnerVo{
		OrgId:         orgId,
		UserId:        userId,
		AppId:         issueInfo.AppId,
		ProjectId:     issueInfo.ProjectId,
		TableId:       issueInfo.TableId,
		Data:          []map[string]interface{}{d},
		TrendPushType: &pushType,
	}, true, &projectvo.TriggerBy{
		TriggerBy: consts.TriggerByNormal,
	})
	if sysErr != nil {
		log.Error(sysErr)
		return sysErr
	}

	//return nil
	//})
	//if transErr != nil {
	//	log.Error(transErr)
	//	return errs.MysqlOperateError
	//}

	////动态
	//asyn.Execute(func() {
	//	issueTrendsBo := bo.IssueTrendsBo{
	//		PushType:      consts.PushTypeWithdrawIssue,
	//		OrgId:         orgId,
	//		UserId:    userId,
	//		IssueId:       issueId,
	//		ParentIssueId: issueInfo.ParentId,
	//		ProjectId:     issueInfo.ProjectId,
	//		PriorityId:    issueInfo.PriorityId,
	//		ParentId:      issueInfo.ParentId,
	//
	//		IssueTitle:    issueInfo.Title,
	//		IssueStatusId: issueStatusId,
	//		TableId:       issueInfo.TableId,
	//	}
	//
	//	asyn.Execute(func() {
	//		domain.PushIssueTrends(issueTrendsBo)
	//	})
	//
	//	asyn.Execute(func() {
	//		PushModifyIssueNotice(issueInfo.OrgId, issueInfo.ProjectId, issueInfo.Id, userId)
	//	})
	//})

	return nil
}

func ReportMoveEvents(orgId, userId int64, oldIssues []*bo.IssueBo, issueIds []int64, fromTableId, targetTableId int64) {
	var errSys errs.SystemErrorInfo

	// 获取表头
	var oldTableColumns, newTableColumns map[string]*projectvo.TableColumnData
	oldTableColumns, errSys = domain.GetTableColumnsMap(orgId, fromTableId, nil, true)
	if errSys != nil {
		log.Errorf("[ReportMoves] 获取表头失败 org:%d table:%d, err: %v", orgId, fromTableId, errSys)
		return
	}
	if targetTableId != fromTableId {
		newTableColumns, errSys = domain.GetTableColumnsMap(orgId, targetTableId, nil, true)
		if errSys != nil {
			log.Errorf("[ReportMoves] 获取表头失败 org:%d table:%d, err: %v", orgId, targetTableId, errSys)
			return
		}
	} else {
		newTableColumns = oldTableColumns
	}

	var userIds, deptIds []int64

	// issueIds 里面的任务可能移动也可能没移动，但这里统一都拿出来，然后再判断是否有变化
	newDatas, errSys := domain.GetIssueInfosMapLcByIssueIds(orgId, userId, issueIds)
	if errSys != nil {
		log.Errorf("[ReportMoves] GetIssueInfosMapLcByIssueIds org:%d, issues:%q err: %v", orgId, issueIds, errSys)
		return
	}
	oldIssueMap := make(map[int64]*bo.IssueBo)
	for _, issue := range oldIssues {
		oldIssueMap[issue.Id] = issue
	}
	for _, newData := range newDatas {
		uIds, dIds := domain.CollectUserDeptIds(newData, newTableColumns)
		userIds = append(userIds, uIds...)
		deptIds = append(deptIds, dIds...)
	}
	userDepts := domain.AssembleUserDeptsByIds(orgId, userIds, deptIds)

	for _, newData := range newDatas {
		issueId := cast.ToInt64(newData[consts.BasicFieldIssueId])
		if oldIssue, ok := oldIssueMap[issueId]; ok {
			domain.ReportMoveEvent(userId, oldIssue, newData, userDepts)
		}
	}
}

// ChangeParentIssue 变更父任务（现已不支持跨表）
func ChangeParentIssue(orgId, userId int64, input vo.ChangeParentIssueReq) (*vo.Void, errs.SystemErrorInfo) {
	//获取原本的Issue
	issueBo, errSys := domain.GetIssueInfoLc(orgId, userId, input.IssueID)
	if errSys != nil {
		log.Error(errSys)
		return nil, errs.BuildSystemErrorInfo(errs.IssueNotExist, errSys)
	}

	errSys = domain.AuthIssue(orgId, userId, issueBo, consts.RoleOperationPathOrgProIssueT, consts.OperationProIssue4Modify)
	if errSys != nil {
		log.Error(errSys)
		return nil, errs.BuildSystemErrorInfo(errs.Unauthorized, errSys)
	}

	if issueBo.ParentId == input.ParentID {
		return &vo.Void{ID: input.IssueID}, nil
	}

	//获取目标父任务
	parentIssueBo, errSys := domain.GetIssueInfoLc(orgId, userId, input.ParentID)
	if errSys != nil {
		log.Error(errSys)
		return nil, errSys
	}

	//先转化为目标任务的子任务
	allIssueIds, oldIssues, errSys := domain.ChangeIssueChildren(userId, issueBo, input.ParentID, parentIssueBo.Path)
	if errSys != nil {
		log.Error(errSys)
		return nil, errSys
	}
	targetTableId := parentIssueBo.TableId

	// 上报事件
	asyn.Execute(func() {
		ReportMoveEvents(orgId, userId, oldIssues, allIssueIds, issueBo.TableId, targetTableId)
	})

	return &vo.Void{
		ID: input.IssueID,
	}, nil
}

// ConvertIssueToParent 转化为父任务（现已不支持跨表）
func ConvertIssueToParent(orgId, userId int64, input vo.ConvertIssueToParentReq) (*vo.Void, errs.SystemErrorInfo) {
	//获取原本的Issue
	issueBo, errSys := domain.GetIssueInfoLc(orgId, userId, input.ID)
	if errSys != nil {
		log.Error(errSys)
		return nil, errs.BuildSystemErrorInfo(errs.IssueNotExist, errSys)
	}

	errSys = domain.AuthIssue(orgId, userId, issueBo, consts.RoleOperationPathOrgProIssueT, consts.OperationProIssue4Modify)
	if errSys != nil {
		log.Error(errSys)
		return nil, errs.BuildSystemErrorInfo(errs.Unauthorized, errSys)
	}

	//先转化为父任务
	allIssueIds, oldIssues, errSys := domain.ChangeIssueChildren(userId, issueBo, 0, "")
	if errSys != nil {
		log.Error(errSys)
		return nil, errSys
	}
	targetTableId := cast.ToInt64(input.TableID)

	// 上报事件
	asyn.Execute(func() {
		ReportMoveEvents(orgId, userId, oldIssues, allIssueIds, issueBo.TableId, targetTableId)
	})

	return &vo.Void{
		ID: input.ID,
	}, nil
}

func GetIssueRowList(orgId, userId int64, input *tablePb.ListRawRequest) ([]map[string]interface{}, errs.SystemErrorInfo) {
	rawListResp, err := tablefacade.ListRawExpand(orgId, userId, input)
	if err != nil {
		log.Errorf("[GetIssueRowList] ListRawExpand err:%v, orgId:%v, userId:%v, ListRawRequest:%v",
			err, orgId, userId, input)
		return nil, err
	}
	return rawListResp.Data, nil
}
