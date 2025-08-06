package domain

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	automationPb "gitea.bjx.cloud/LessCode/interface/golang/automation/v1"
	msgPb "gitea.bjx.cloud/LessCode/interface/golang/msg/v1"
	pushPb "gitea.bjx.cloud/LessCode/interface/golang/push/v1"
	tablePb "gitea.bjx.cloud/LessCode/interface/golang/table/v1"
	"github.com/star-table/startable-server/app/facade/common"
	"github.com/star-table/startable-server/app/facade/commonfacade"
	"github.com/star-table/startable-server/app/facade/formfacade"
	"github.com/star-table/startable-server/app/facade/orgfacade"
	"github.com/star-table/startable-server/app/facade/resourcefacade"
	sconsts "github.com/star-table/startable-server/app/service/projectsvc/consts"
	"github.com/star-table/startable-server/common/core/businees"
	"github.com/star-table/startable-server/common/core/config"
	"github.com/star-table/startable-server/common/core/consts"
	"github.com/star-table/startable-server/common/core/errs"
	"github.com/star-table/startable-server/common/core/threadlocal"
	"github.com/star-table/startable-server/common/core/util/asyn"
	"github.com/star-table/startable-server/common/core/util/copyer"
	"github.com/star-table/startable-server/common/core/util/json"
	"github.com/star-table/startable-server/common/core/util/jsonx"
	"github.com/star-table/startable-server/common/core/util/slice"
	slice2 "github.com/star-table/startable-server/common/core/util/slice"
	"github.com/star-table/startable-server/common/core/util/str"
	"github.com/star-table/startable-server/common/extra/card"
	"github.com/star-table/startable-server/common/extra/lc_helper"
	"github.com/star-table/startable-server/common/library/cache"
	"github.com/star-table/startable-server/common/library/db/mysql"
	"github.com/star-table/startable-server/common/model/bo"
	"github.com/star-table/startable-server/common/model/bo/status"
	"github.com/star-table/startable-server/common/model/vo"
	"github.com/star-table/startable-server/common/model/vo/commonvo"
	"github.com/star-table/startable-server/common/model/vo/datacenter"
	"github.com/star-table/startable-server/common/model/vo/formvo"
	"github.com/star-table/startable-server/common/model/vo/lc_table"
	"github.com/star-table/startable-server/common/model/vo/orgvo"
	"github.com/star-table/startable-server/common/model/vo/projectvo"
	"github.com/star-table/startable-server/common/model/vo/resourcevo"
	"github.com/star-table/startable-server/common/model/vo/uservo"
	"github.com/spf13/cast"
	"upper.io/db.v3"
	"upper.io/db.v3/lib/sqlbuilder"
)

func GetCalendarInfo(orgId, projectId int64) (isOk bool, info bo.CalendarInfo) {
	isOk = false
	projectCalendarInfo, err := GetProjectCalendarInfo(orgId, projectId)
	if err != nil {
		log.Error(err)
		return
	}
	// 判断略显麻烦，先去掉。
	//if !CheckSyncOutCalendarValid(projectCalendarInfo.IsSyncOutCalendar) {
	//	log.Error("项目设置导出日历的配置不合法。")
	//	return
	//}
	if projectCalendarInfo.CalendarId == consts.BlankString {
		log.Info("无对应项目日历或未设置导出日历")
		return
	}
	orgBaseInfo, err := orgfacade.GetBaseOrgInfoRelaxed(orgId)
	if err != nil {
		log.Errorf("组织外部信息不存在 %v", err)
		return
	}
	isOk = true
	info.OrgId = orgId
	info.OutOrgId = orgBaseInfo.OutOrgId
	info.CalendarId = projectCalendarInfo.CalendarId
	info.Creator = projectCalendarInfo.Creator
	info.SourceChannel = orgBaseInfo.SourceChannel
	info.SyncCalendarFlag = projectCalendarInfo.IsSyncOutCalendar
	return
}

// GetIssueFinishedStatForProjects 获取多个项目下任务完成进度。逾期任务数、已完成未完成任务数等
// 参考 `SelectIssueAssignCount` 函数
func GetIssueFinishedStatForProjects(orgId int64, curUserId int64, summaryAppId int64, projectIds []int64) (map[int64]bo.IssueStatistic, errs.SystemErrorInfo) {
	issueStat := map[int64]bo.IssueStatistic{}
	// 从无码接口查询两次：
	// 查询逾期未完成任务数、已完成任务数，并按项目分组
	issueFinishedOrNotRes, err := StatisticProIssueFinishedOrNot(orgId, curUserId, summaryAppId, projectIds)
	if err != nil {
		log.Errorf("[GetIssueFinishedStatForProjects] orgId:%d, StatisticProIssueFinishedOrNot err: %v", orgId, err)
		return issueStat, err
	}
	dataJsonStr := json.ToJsonIgnoreError(issueFinishedOrNotRes)
	statistic := []bo.IssueStatistic{}
	json.FromJson(dataJsonStr, &statistic)
	for _, v := range statistic {
		issueStat[v.ProjectId] = v
	}

	// 查询我关注或者我负责的未完成的任务数，按项目分组
	//relateRes, err := StatisticProIssuesRelateSomeone(orgId, curUserId, summaryAppId, projectIds)
	//if err != nil {
	//	log.Errorf("[GetIssueFinishedStatForProjects] orgId:%d, GetIssueFinishedStatForProjects err: %v", orgId, err)
	//	return issueStat, err
	//}
	//for _, total := range relateRes {
	//	unfinish := total.Total
	//	if temp, ok := issueStat[total.ProjectId]; ok {
	//		// 我相关的未完成的任务统计
	//		temp.RelateUnfinish = unfinish
	//		issueStat[total.ProjectId] = temp
	//	}
	//}
	//log.Infof("[GetIssueFinishedStatForProjects] issueStat: %s", json.ToJsonIgnoreError(issueStat))

	return issueStat, nil
}

// StatisticProIssueFinishedOrNot 查询逾期未完成任务数、已完成任务数，并按项目分组
func StatisticProIssueFinishedOrNot(orgId int64, curUserId int64, summaryAppId int64, projectIds []int64) ([]map[string]interface{}, errs.SystemErrorInfo) {
	nowTimeDateStr := time.Now().Format(consts.AppTimeFormat)
	finishedTypeIds := []int64{consts.StatusTypeComplete}
	unfinishedTypeIds := []int64{consts.StatusTypeNotStart, consts.StatusTypeRunning}
	// data::jsonb->statusType 方式的查询，会将数值转为字符串，因此，in 查询后的列举值必须是带引号的字符串，如：`in ('7', '8')`
	finishedTypeIdsStr := "'" + str.Int64Implode(finishedTypeIds, "','") + "'"
	unfinishedTypeIdsStr := "'" + str.Int64Implode(unfinishedTypeIds, "','") + "'"

	condition := &tablePb.Condition{
		Type: tablePb.ConditionType_and,
		Conditions: GetNoRecycleCondition(
			GetRowsCondition(consts.BasicFieldProjectId, tablePb.ConditionType_in, nil, projectIds),
			//GetRowsCondition(consts.BasicFieldPlanEndTime, tablePb.ConditionType_between, nil, []interface{}{consts.BlankTime, nowTimeDateStr}),
		),
	}

	req := &tablePb.ListRawRequest{
		DbType: tablePb.DbType_slave1,
		FilterColumns: []string{
			" count(*) as \"all\" ",
			fmt.Sprintf("\"%s\"", consts.BasicFieldProjectId),
			fmt.Sprintf(" sum(case when \"data\" :: jsonb -> '%s' in (%s) and \"data\" :: jsonb ->> '%s' >'%s' and \"data\" :: jsonb ->> '%s' <'%s' then 1 else 0 end) \"overdue\" ",
				consts.BasicFieldIssueStatusType, unfinishedTypeIdsStr, consts.BasicFieldPlanEndTime, "1970-01-01 00:00:00", consts.BasicFieldPlanEndTime, nowTimeDateStr),
			fmt.Sprintf(" sum(case when \"data\" :: jsonb -> '%s' in (%s) then 1 else 0 end) \"finish\" ", consts.BasicFieldIssueStatusType, finishedTypeIdsStr),
		},
		Condition: condition,
		Groups: []string{
			// group by 用别名也可
			fmt.Sprintf("\"%s\"", consts.BasicFieldProjectId),
		},
		Page: 1,
		Size: int32(len(projectIds) + 10),
	}

	lessResp, err := GetRawRows(orgId, curUserId, req)
	if err != nil {
		log.Errorf("[GetIssueFinishedStatForProjects] orgId:%d, err: %v", orgId, err.Error())
		return nil, err
	}

	return lessResp.Data, nil
}

// GetIssueIdList 通过无码查询任务 id 数据
func GetIssueIdList(orgId int64, conds []*tablePb.Condition, page, size int64) ([]*bo.IssueIdData, errs.SystemErrorInfo) {
	list, err := GetIssueInfosMapLc(orgId, 0, &tablePb.Condition{
		Type:       tablePb.ConditionType_and,
		Conditions: conds,
	}, []string{lc_helper.ConvertToFilterColumn(consts.BasicFieldIssueId)}, page, size)
	if err != nil {
		return nil, err
	}

	issueIds := make([]*bo.IssueIdData, 0, len(list))
	if len(list) == 0 {
		return issueIds, nil
	}

	err2 := copyer.Copy(list, &issueIds)
	if err2 != nil {
		log.Errorf("[GetIssueIdList] json转换错误, err:%v", err2)
		return nil, errs.JSONConvertError
	}
	return issueIds, nil
}

// GetIssueInfoByIterationForLc 统计迭代下完成的、逾期的任务数，按迭代分组，通过无码接口查询。
func GetIssueInfoByIterationForLc(iterationIds []int64, orgId int64, projectId int64) (map[int64]bo.IssueStatistic, errs.SystemErrorInfo) {
	issueStat := map[int64]bo.IssueStatistic{}
	statistic := make([]*bo.IssueStatistic, 0)
	finishedTypeIds := []int64{consts.StatusTypeComplete}
	unfinishedTypeIds := []int64{consts.StatusTypeNotStart, consts.StatusTypeRunning}
	// data::jsonb->statusType 方式的查询，会将数值转为字符串，因此，in 查询后的列举值必须是带引号的字符串，如：`in ('7', '8')`
	finishedTypeIdsStr := "'" + str.Int64Implode(finishedTypeIds, "','") + "'"
	unfinishedTypeIdsStr := "'" + str.Int64Implode(unfinishedTypeIds, "','") + "'"
	nowTimeDateStr := time.Now().Format(consts.AppTimeFormat)

	condition := &tablePb.Condition{
		Type:       tablePb.ConditionType_and,
		Conditions: GetNoRecycleCondition(),
	}

	if projectId != 0 {
		condition.Conditions = append(condition.Conditions, GetRowsCondition(consts.BasicFieldProjectId, tablePb.ConditionType_equal, projectId, nil))
		condition.Conditions = append(condition.Conditions, GetRowsCondition(consts.BasicFieldIterationId, tablePb.ConditionType_in, nil, iterationIds))
	} else {
		condition.Conditions = append(condition.Conditions, GetRowsCondition(consts.BasicFieldIterationId, tablePb.ConditionType_in, nil, iterationIds))
	}

	req := &tablePb.ListRawRequest{
		DbType: tablePb.DbType_slave1,
		FilterColumns: []string{
			"count(*) as all",
			fmt.Sprintf("\"data\" :: jsonb -> '%s' \"%s\"", consts.BasicFieldIterationId, consts.BasicFieldIterationId),
			fmt.Sprintf(" sum(case when \"data\" :: jsonb -> '%s' in (%s) and \"data\" :: jsonb ->> '%s' >'%s' and \"data\" :: jsonb ->> '%s' <'%s' then 1 else 0 end) \"overdue\" ",
				consts.BasicFieldIssueStatusType, unfinishedTypeIdsStr, consts.BasicFieldPlanEndTime, consts.BlankTime, consts.BasicFieldPlanEndTime, nowTimeDateStr),
			fmt.Sprintf(" sum(case when \"data\" :: jsonb -> '%s' in (%s) then 1 else 0 end) \"finish\" ", consts.BasicFieldIssueStatusType, finishedTypeIdsStr),
		},
		Condition: condition,
		Groups: []string{
			"\"" + consts.BasicFieldIterationId + "\"",
		},
	}

	lessResp, err := GetRawRows(orgId, 0, req)
	if err != nil {
		log.Errorf("[GetIssueInfoByIterationForLc] orgId:%d, projectId:%d, LessIssueList failure:%v", orgId, projectId, err)
		return nil, err
	}

	err2 := copyer.Copy(lessResp.Data, &statistic)
	if err2 != nil {
		log.Errorf("[GetIssueInfoByIterationForLc] orgId:%d, projectId:%d, Copy failure:%v", orgId, projectId, err2)
		return nil, errs.BuildSystemErrorInfo(errs.ObjectCopyError, err2)
	}
	for _, v := range statistic {
		issueStat[v.IterationId] = *v
	}

	return issueStat, nil
}

// GetIssueAndChildrenIds 查询任务以及所有子任务的任务ID
func GetIssueAndChildrenIds(orgId int64, issueIds []int64) ([]int64, errs.SystemErrorInfo) {
	if len(issueIds) == 0 {
		return nil, nil
	}

	allIds := issueIds

	// 首先查找有子任务的任务，即父任务
	list, errSys := GetIssueInfosMapLc(orgId, 0,
		GetRowsCondition(consts.BasicFieldParentId, tablePb.ConditionType_in, nil, issueIds),
		[]string{"DISTINCT " + lc_helper.ConvertToFilterColumn(consts.BasicFieldParentId)}, -1, -1, true)
	if errSys != nil {
		return nil, errSys
	}
	if len(list) > 0 {
		parentIds := make([]*bo.ParentIdData, 0, len(list))
		if err := copyer.Copy(list, &parentIds); err != nil {
			log.Errorf("[GetIssueAndChildrenIds] json转换错误, err:%v", err)
			return nil, errs.JSONConvertError
		}

		// 根据path模糊查找所有子任务
		conditions := make([]*tablePb.Condition, 0, len(parentIds)+1)
		conditions = append(conditions)
		for _, parentId := range parentIds {
			conditions = append(conditions, GetRowsCondition(consts.BasicFieldPath, tablePb.ConditionType_like, fmt.Sprintf(",%d,", parentId.Id), nil))
		}
		issueIdList, errSys := GetIssueIdList(orgId, []*tablePb.Condition{
			GetRowsCondition(consts.BasicFieldRecycleFlag, tablePb.ConditionType_equal, 2, nil),
			{Type: tablePb.ConditionType_or, Conditions: conditions}}, -1, -1)
		if errSys != nil {
			log.Errorf("[GetIssueAndChildrenIds] GetIssueIdList err:%v, orgId:%v, issueIds:%v", errSys, orgId, issueIds)
			return nil, errSys
		}

		for _, child := range issueIdList {
			allIds = append(allIds, child.Id)
		}
	}
	allIds = slice.SliceUniqueInt64(allIds)
	return allIds, nil
}

func DeleteIssue(issueBo *bo.IssueBo, operatorId int64, sourceChannel string, takeChildren *bool) ([]int64, errs.SystemErrorInfo) {
	orgId := issueBo.OrgId
	issueId := issueBo.Id

	var beforeParticipantIds []int64
	beforeFollowerIds := issueBo.FollowerIdsI64
	//beforeParticipantIds, err1 := GetIssueRelationIdsByRelateType(orgId, issueId, consts.IssueRelationTypeParticipant)
	//if err1 != nil {
	//	log.Error(err1)
	//	return nil, errs.BuildSystemErrorInfo(errs.IssueDomainError, err1)
	//}
	//beforeFollowerIds, err1 := GetIssueRelationIdsByRelateType(orgId, issueId, consts.IssueRelationTypeFollower)
	//if err1 != nil {
	//	log.Error(err1)
	//	return nil, errs.BuildSystemErrorInfo(errs.IssueDomainError, err1)
	//}
	appId, appIdErr := GetAppIdFromProjectId(orgId, issueBo.ProjectId)
	if appIdErr != nil {
		log.Error(appIdErr)
		return nil, appIdErr
	}

	var updateBatchReqs []*formvo.LessUpdateIssueBatchReq

	condition := &tablePb.Condition{
		Type:       tablePb.ConditionType_and,
		Conditions: GetNoRecycleCondition(GetRowsCondition(consts.BasicFieldPath, tablePb.ConditionType_like, fmt.Sprintf(",%d,", issueId), nil)),
	}
	filterColumns := []string{
		lc_helper.ConvertToFilterColumn(consts.BasicFieldIssueId),
	}
	issueChildInfos, errSys := GetIssueInfosMapLc(orgId, operatorId, condition, filterColumns, -1, -1)
	if errSys != nil {
		log.Errorf("[DeleteIssue] GetIssueInfosMapLc err:%v, orgId:%v, issueId:%v", errSys, orgId, issueId)
		return nil, errSys
	}
	issueChildIds := make([]int64, 0, len(issueChildInfos))
	for _, data := range issueChildInfos {
		issueChildIds = append(issueChildIds, cast.ToInt64(data[consts.BasicFieldIssueId]))
	}

	summaryAppId, errSys := GetOrgSummaryAppId(orgId)
	if errSys != nil {
		log.Error(errSys)
		return nil, errSys
	}

	// 拉取原始数据
	condition = &tablePb.Condition{
		Type:       tablePb.ConditionType_and,
		Conditions: GetNoRecycleCondition(GetRowsCondition(consts.BasicFieldIssueId, tablePb.ConditionType_equal, issueId, nil)),
	}
	issueDatas, errSys := GetIssueInfosMapLc(orgId, operatorId, condition, nil, -1, -1)
	if errSys != nil || len(issueDatas) == 0 {
		log.Errorf("[DeleteIssue] GetIssueInfosMapLc err:%v, orgId:%v, issueId:%v", errSys, orgId, issueId)
		return nil, errSys
	}
	issueData := issueDatas[0]

	// 获取表头
	fromTableSchema, errSys := GetTableColumnsMap(orgId, issueBo.TableId, nil, true)
	if errSys != nil {
		log.Error(errSys)
		return nil, errSys
	}
	var allRelatingIssueIds []int64
	for columnId, column := range fromTableSchema {
		if column.Field.Type == tablePb.ColumnType_relating.String() ||
			column.Field.Type == tablePb.ColumnType_baRelating.String() ||
			column.Field.Type == tablePb.ColumnType_singleRelating.String() {
			if v, ok := issueData[columnId]; ok {
				oldRelating := &bo.RelatingIssue{}
				jsonx.Copy(v, oldRelating)
				allRelatingIssueIds = append(allRelatingIssueIds, slice2.StringToInt64Slice(oldRelating.LinkTo)...)
				allRelatingIssueIds = append(allRelatingIssueIds, slice2.StringToInt64Slice(oldRelating.LinkFrom)...)
			}
		}
	}
	allRelatingIssueIds = slice.SliceUniqueInt64(allRelatingIssueIds)
	if len(allRelatingIssueIds) > 0 {
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
						Values: allRelatingIssueIds,
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

	conn, err := mysql.GetConnect()
	if err != nil {
		log.Errorf(consts.DBOpenErrorSentence, err)
		return nil, errs.BuildSystemErrorInfo(errs.MysqlOperateError, err)
	}
	tx, err := conn.NewTx(nil)
	if err != nil {
		log.Errorf(consts.TxOpenErrorSentence, err)
		return nil, errs.BuildSystemErrorInfo(errs.MysqlOperateError, err)
	}
	defer mysql.Close(conn, tx)

	allIssueId := []int64{issueId}
	childIssueIds := []int64{}
	if takeChildren == nil || *takeChildren {
		// 默认删除子任务
		//for _, issue := range *issueChild {
		//	allIssueId = append(allIssueId, issue.Id)
		//	childIssueIds = append(childIssueIds, issue.Id)
		//}
		allIssueId = append(allIssueId, issueChildIds...)
		childIssueIds = append(childIssueIds, issueChildIds...)
	} else {
		newPath := "0,"

		// 如果不删除子任务，就把顶级子任务变为父任务
		likeCond := fmt.Sprintf("%s%d,%s", issueBo.Path, issueId, "%")
		//_, err = mysql.TransUpdateSmartWithCond(tx, consts.TableIssue, db.Cond{
		//	consts.TcPath:     db.Like(likeCond),
		//	consts.TcOrgId:    orgId,
		//	consts.TcParentId: issueId,
		//}, mysql.Upd{
		//	consts.TcParentId: 0,
		//	consts.TcPath:     db.Raw(fmt.Sprintf("replace(`%s`, '%s', '%s')", consts.TcPath, issueBo.Path, newPath)),
		//	consts.TcUpdator:  operatorId,
		//})
		//if err != nil {
		//	log.Error(err)
		//	return nil, errs.BuildSystemErrorInfo(errs.MysqlOperateError, err)
		//}

		// 子任务的子任务，更新path
		//_, err = mysql.TransUpdateSmartWithCond(tx, consts.TableIssue, db.Cond{
		//	consts.TcPath:     db.Like(likeCond),
		//	consts.TcOrgId:    orgId,
		//	consts.TcParentId: db.NotEq(issueId),
		//}, mysql.Upd{
		//	consts.TcPath:    db.Raw(fmt.Sprintf("replace(`%s`, '%s', '%s')", consts.TcPath, issueBo.Path, newPath)),
		//	consts.TcUpdator: operatorId,
		//})
		//if err != nil {
		//	log.Error(err)
		//	return nil, errs.BuildSystemErrorInfo(errs.MysqlOperateError, err)
		//}

		// 无码的 刷顶级子任务的parentId
		updateBatchReqs = append(updateBatchReqs, &formvo.LessUpdateIssueBatchReq{
			OrgId:  issueBo.OrgId,
			AppId:  appId,
			UserId: operatorId,
			Condition: vo.LessCondsData{
				Type: consts.ConditionAnd,
				Conds: []*vo.LessCondsData{
					{
						Type:   consts.ConditionEqual,
						Value:  issueBo.OrgId,
						Column: lc_helper.ConvertToCondColumn(consts.BasicFieldOrgId),
					},
					{
						Type:   consts.ConditionEqual,
						Value:  issueId,
						Column: lc_helper.ConvertToCondColumn(consts.BasicFieldParentId),
					},
				},
			},
			Sets: []datacenter.Set{
				{
					Column:          lc_helper.ConvertToCondColumn(consts.BasicFieldParentId),
					Value:           0,
					Type:            consts.SetTypeNormal,
					Action:          consts.SetActionSet,
					WithoutPretreat: false,
				},
			},
		})
		// 无码的 刷子任务的path
		updateBatchReqs = append(updateBatchReqs, &formvo.LessUpdateIssueBatchReq{
			OrgId:  issueBo.OrgId,
			AppId:  appId,
			UserId: operatorId,
			Condition: vo.LessCondsData{
				Type: consts.ConditionAnd,
				Conds: []*vo.LessCondsData{
					{
						Type:   consts.ConditionEqual,
						Value:  issueBo.OrgId,
						Column: lc_helper.ConvertToCondColumn(consts.BasicFieldOrgId),
					},
					{
						Type:   consts.ConditionLike,
						Value:  likeCond,
						Column: lc_helper.ConvertToCondColumn(consts.BasicFieldPath),
					},
				},
			},
			Sets: []datacenter.Set{
				{
					Column:          lc_helper.ConvertToCondColumn(consts.BasicFieldPath),
					Value:           fmt.Sprintf("regexp_replace(\"%s\",'%s','%s')", consts.BasicFieldPath, issueBo.Path, newPath),
					Type:            consts.SetTypeNormal,
					Action:          consts.SetActionSet,
					WithoutPretreat: true,
				},
			},
		})

	}

	recycleVersionId, versionErr := AddRecycleRecord(orgId, operatorId, issueBo.ProjectId, []int64{issueId}, consts.RecycleTypeIssue, tx)
	if versionErr != nil {
		log.Error(versionErr)
		return nil, versionErr
	}

	// 删除任务的关联实体前，进行一些操作。尝试删除日程
	if err := DeleteCalendarEventBatch(orgId, issueBo.ProjectId, []int64{issueBo.Id}); err != nil {
		errMsg := fmt.Sprintf("DeleteIssue 删除任务，后续处理，删除日程异常：%v", err)
		log.Error(errMsg)
	}

	err3 := DeleteAllIssueRelation(tx, operatorId, orgId, allIssueId, recycleVersionId)
	if err3 != nil {
		log.Error(err3)
		return nil, errs.BuildSystemErrorInfo(errs.IssueDomainError, err3)
	}

	// 更新无码
	for _, req := range updateBatchReqs {
		resp := formfacade.LessUpdateIssueBatchRaw(req)
		if resp.Failure() {
			log.Error(resp.Error())
			return nil, resp.Error()
		}
	}

	//无码删除
	resp := formfacade.LessRecycleIssue(formvo.LessRecycleIssueReq{
		AppId:    appId,
		OrgId:    orgId,
		UserId:   operatorId,
		IssueIds: allIssueId,
		TableId:  issueBo.TableId,
	})
	if resp.Failure() {
		log.Error(resp.Error())
		return nil, resp.Error()
	}

	err = tx.Commit()
	if err != nil {
		return nil, errs.BuildSystemErrorInfo(errs.MysqlOperateError, err)
	}

	//文件关联
	resourceRelationResp := resourcefacade.DeleteResourceRelation(resourcevo.DeleteResourceRelationReqVo{
		OrgId:  orgId,
		UserId: operatorId,
		Input: resourcevo.DeleteResourceRelationData{
			ProjectId:        issueBo.ProjectId,
			IssueIds:         allIssueId,
			RecycleVersionId: recycleVersionId,
			SourceTypes:      []int{consts.OssPolicyTypeIssueResource, consts.OssPolicyTypeLesscodeResource},
		},
	})
	if resourceRelationResp.Failure() {
		log.Error(resourceRelationResp.Error())
		return nil, resourceRelationResp.Error()
	}
	// 工时的删除
	if err := DeleteWorkHourForIssues(orgId, []int64{issueId}, recycleVersionId); err != nil {
		log.Errorf("[DeleteIssue] issueId: %d, DeleteWorkHourForIssues err: %v", issueId, err)
		return nil, err
	}

	asyn.Execute(func() {
		ownerIds, errSys := businees.LcMemberToUserIdsWithError(issueBo.OwnerId)
		if errSys != nil {
			log.Errorf("原始的数据:%v, err:%v", issueBo.OwnerId, err)
			return
		}
		blank := []int64{}
		issueTrendsBo := &bo.IssueTrendsBo{
			PushType:                 consts.PushTypeDeleteIssue,
			OrgId:                    orgId,
			OperatorId:               operatorId,
			DataId:                   issueBo.DataId,
			IssueId:                  issueBo.Id,
			ParentIssueId:            issueBo.ParentId,
			ProjectId:                issueBo.ProjectId,
			TableId:                  issueBo.TableId,
			PriorityId:               issueBo.PriorityId,
			ParentId:                 issueBo.ParentId,
			IssueTitle:               issueBo.Title,
			IssueStatusId:            issueBo.Status,
			BeforeOwner:              ownerIds,
			AfterOwner:               blank,
			BeforeChangeFollowers:    beforeFollowerIds,
			AfterChangeFollowers:     blank,
			BeforeChangeParticipants: beforeParticipantIds,
			AfterChangeParticipants:  blank,
			IssueChildren:            childIssueIds,
			SourceChannel:            sourceChannel,
		}
		asyn.Execute(func() {
			PushIssueTrends(issueTrendsBo)
		})
		asyn.Execute(func() {
			PushIssueThirdPlatformNotice(issueTrendsBo)
		})
	})

	return childIssueIds, nil
}

func DeleteIssueBatch(orgId int64, issueAndChildrenBos []*bo.IssueBo, allIssueId []int64, operatorId int64, sourceChannel string, projectId int64, tableId int64) ([]int64, errs.SystemErrorInfo) {
	needChangeParentIds := make([]int64, 0) //需要转化为独立任务（父任务id转为0）
	needChangePath := make([]string, 0)
	pathIssueIdMap := make(map[string]int64, 0)

	allInputIssues := make([]*bo.IssueBo, 0)
	allRemainIssues := make([]*bo.IssueBo, 0)
	issueBos := make([]*bo.IssueBo, 0) //需要移动的所有父任务（包含变更后成为父任务的）
	allParentIds := make([]int64, 0)
	couldDeleteIssueIds := make([]int64, 0)

	var relatingColumnIds []string
	tableSchema, errSys := GetTableColumnsMap(orgId, tableId, nil, true)
	if errSys != nil {
		log.Error(errSys)
		return nil, errSys
	}
	for columnId, column := range tableSchema {
		if column.Field.Type == tablePb.ColumnType_relating.String() ||
			column.Field.Type == tablePb.ColumnType_baRelating.String() ||
			column.Field.Type == tablePb.ColumnType_singleRelating.String() {
			relatingColumnIds = append(relatingColumnIds, columnId)
		}
	}
	appId, appIdErr := GetAppIdFromProjectId(orgId, projectId)
	if appIdErr != nil {
		log.Error(appIdErr)
		return nil, appIdErr
	}

	// issueAndChildrenBos 中可能包含只有父任务的id，也可能只有单独的子任务id。因为可以只删除父任务，而不删除子任务。即存在以下几种情况：
	// 	* 父任务、子任务一起删除
	// 	* 只有子任务删除
	// 	* 只有父任务删除
	// 		* 这种情况需要**注意**，删除父任务之后，需要将其子任务的 parentId 改为 0；

	// issueAndChildrenBos 是**有权限**删除的任务列表，如果它为空，则表示没有能被删除的任务。此时可以直接返回。
	if len(issueAndChildrenBos) < 1 {
		log.Infof("[DeleteIssueBatch] 没有可以删除的任务。orgId: %d", orgId)
		return nil, nil
	}
	for _, issueBo := range issueAndChildrenBos {
		if ok, _ := slice.Contain(allIssueId, issueBo.Id); ok {
			allInputIssues = append(allInputIssues, issueBo)
			couldDeleteIssueIds = append(couldDeleteIssueIds, issueBo.Id)
			if issueBo.ParentId != 0 {
				if ok1, _ := slice.Contain(allIssueId, issueBo.ParentId); !ok1 {
					issueBos = append(issueBos, issueBo)
					allParentIds = append(allParentIds, issueBo.Id)
				}
			} else {
				issueBos = append(issueBos, issueBo)
				allParentIds = append(allParentIds, issueBo.Id)
			}
		} else {
			allRemainIssues = append(allRemainIssues, issueBo)
			//需要转化为父任务的逻辑（当前任务不需要被删除，但是父任务需要删除）
			if ok, _ := slice.Contain(allIssueId, issueBo.ParentId); ok {
				needChangeParentIds = append(needChangeParentIds, issueBo.Id)
				//curPath := fmt.Sprintf("%s%d,", issueBo.Path, issueBo.Id)
				curPath := issueBo.Path
				needChangePath = append(needChangePath, curPath)
				pathIssueIdMap[curPath] = issueBo.Id
			}
		}
	}
	if len(couldDeleteIssueIds) < 1 {
		log.Info("没有可以删除的任务。")
		return nil, nil
	}

	var updateBatchReqs []*formvo.LessUpdateIssueBatchReq

	// 需要移动的任务涉及的关联引用的任务ID
	if len(relatingColumnIds) > 0 {
		var allRelatingIssueIds []int64
		for _, issueBo := range allInputIssues {
			for _, columnId := range relatingColumnIds {
				if v, ok := issueBo.LessData[columnId]; ok {
					oldRelating := &bo.RelatingIssue{}
					jsonx.Copy(v, oldRelating)
					allRelatingIssueIds = append(allRelatingIssueIds, slice2.StringToInt64Slice(oldRelating.LinkTo)...)
					allRelatingIssueIds = append(allRelatingIssueIds, slice2.StringToInt64Slice(oldRelating.LinkFrom)...)
				}
			}
		}
		allRelatingIssueIds = slice.SliceUniqueInt64(allRelatingIssueIds)
		var relatingIssueIds []int64
		for _, id := range allRelatingIssueIds {
			if ok, _ := slice.Contain(couldDeleteIssueIds, id); !ok {
				relatingIssueIds = append(relatingIssueIds, id)
			}
		}
		if len(relatingIssueIds) > 0 {
			updateBatchReqs = append(updateBatchReqs, &formvo.LessUpdateIssueBatchReq{
				OrgId:  orgId,
				AppId:  appId,
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

	//查找只需要修改path的子任务
	changePathIssueMap := make(map[string][]int64, 0)
	childIssueIds := map[int64][]int64{}
	changePathIssueBoMap := make(map[string][]*bo.IssueBo, 0)
	needChangePathMap := make(map[string]int, 0)
	for _, s := range needChangePath {
		needChangePathMap[s] = strings.Count(s, ",")
	}
	//从最长的path去匹配，用最长的去替换path
	for _, issueBo := range allInputIssues {
		if ok1, _ := slice.Contain(allIssueId, issueBo.ParentId); ok1 {
			max := 0
			curPath := ""
			for s, i := range needChangePathMap {
				if strings.Contains(issueBo.Path, s) {
					if i > max {
						curPath = s
						max = i
					}
				}
			}
			if curPath != "" {
				changePathIssueMap[curPath] = append(changePathIssueMap[curPath], issueBo.Id)
				changePathIssueBoMap[curPath] = append(changePathIssueBoMap[curPath], issueBo)
				if parentId, ok := pathIssueIdMap[curPath]; ok {
					childIssueIds[parentId] = append(childIssueIds[parentId], issueBo.Id)
				}
			}
		}
	}

	for _, issueBo := range allRemainIssues {
		if ok1, _ := slice.Contain(allIssueId, issueBo.ParentId); !ok1 {
			max := 0
			curPath := ""
			for s, i := range needChangePathMap {
				if strings.Contains(issueBo.Path, s) {
					if i > max {
						curPath = s
						max = i
					}
				}
			}
			if curPath != "" {
				changePathIssueMap[curPath] = append(changePathIssueMap[curPath], issueBo.Id)
			}
		}
	}
	log.Infof("[DeleteIssueBatch] needChangePathMap: %v, changePathIssueMap: %v", json.ToJsonIgnoreError(needChangePathMap), json.ToJsonIgnoreError(changePathIssueMap))

	childrenMap := map[int64][]*bo.IssueBo{}
	for s, int64s := range changePathIssueBoMap {
		idArr := strings.Split(s[0:len(s)-1], ",")
		if len(idArr) >= 2 { // fix panic
			parentId, err := strconv.ParseInt(idArr[len(idArr)-2], 10, 64)
			if err != nil {
				log.Error(err)
				return nil, errs.ParamError
			}
			childrenMap[parentId] = int64s
		}
	}

	conn, err := mysql.GetConnect()
	if err != nil {
		log.Errorf(consts.DBOpenErrorSentence, err)
		return nil, errs.BuildSystemErrorInfo(errs.MysqlOperateError, err)
	}
	tx, err := conn.NewTx(nil)
	if err != nil {
		log.Errorf(consts.TxOpenErrorSentence, err)
		return nil, errs.BuildSystemErrorInfo(errs.MysqlOperateError, err)
	}
	defer mysql.Close(conn, tx)

	newPath := "0,"
	updForm := make([]map[string]interface{}, 0)
	if len(needChangeParentIds) > 0 {
		//_, err := mysql.TransUpdateSmartWithCond(tx, consts.TableIssue, db.Cond{
		//	consts.TcOrgId: orgId,
		//	consts.TcId:    db.In(needChangeParentIds),
		//}, mysql.Upd{
		//	consts.TcParentId: 0,
		//	consts.TcPath:     newPath,
		//})
		//if err != nil {
		//	log.Error(err)
		//	return nil, errs.MysqlOperateError
		//}
		// 刷无码
		for _, id := range needChangeParentIds {
			updForm = append(updForm, map[string]interface{}{
				consts.BasicFieldIssueId:  id,
				consts.BasicFieldParentId: 0,
				consts.BasicFieldPath:     newPath,
			})
		}

		for i, s := range changePathIssueMap {
			//_, err := mysql.TransUpdateSmartWithCond(tx, consts.TableIssue, db.Cond{
			//	consts.TcOrgId: orgId,
			//	consts.TcId:    db.In(s),
			//}, mysql.Upd{
			//	consts.TcPath: db.Raw(fmt.Sprintf("replace(`%s`, '%s', '%s')", consts.TcPath, i, newPath)),
			//})
			//if err != nil {
			//	log.Error(err)
			//	return nil, errs.MysqlOperateError
			//}

			// 刷无码
			updateBatchReqs = append(updateBatchReqs, &formvo.LessUpdateIssueBatchReq{
				OrgId:  orgId,
				AppId:  appId,
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
							Values: s,
							Column: lc_helper.ConvertToCondColumn(consts.BasicFieldIssueId),
						},
					},
				},
				Sets: []datacenter.Set{
					{
						Column:          lc_helper.ConvertToCondColumn(consts.BasicFieldPath),
						Value:           fmt.Sprintf("regexp_replace(\"%s\",'%s','%s')", consts.BasicFieldPath, i, newPath),
						Type:            consts.SetTypeNormal,
						Action:          consts.SetActionSet,
						WithoutPretreat: true,
					},
				},
			})
		}
	}

	// 无码更新
	resp := formfacade.LessUpdateIssue(formvo.LessUpdateIssueReq{
		AppId:   appId,
		OrgId:   orgId,
		UserId:  operatorId,
		TableId: tableId,
		Form:    updForm,
	})
	if resp.Failure() {
		log.Error(resp.Error())
		return nil, resp.Error()
	}

	// 无码更新
	for _, req := range updateBatchReqs {
		batchResp := formfacade.LessUpdateIssueBatchRaw(req)
		if batchResp.Failure() {
			log.Error(batchResp.Error())
			return nil, batchResp.Error()
		}
	}

	recycleVersionId, versionErr := AddRecycleRecord(orgId, operatorId, projectId, allParentIds, consts.RecycleTypeIssue, tx)
	if versionErr != nil {
		log.Error(versionErr)
		return nil, versionErr
	}
	err2 := deleteIssuesAndRelation(orgId, projectId, appId, tableId, operatorId, recycleVersionId, couldDeleteIssueIds, tx)
	if err2 != nil {
		return nil, err2
	}

	err = tx.Commit()
	if err != nil {
		return nil, errs.BuildSystemErrorInfo(errs.MysqlOperateError, err)
	}

	//文件关联
	resourceRelationResp := resourcefacade.DeleteResourceRelation(resourcevo.DeleteResourceRelationReqVo{
		OrgId:  orgId,
		UserId: operatorId,
		Input: resourcevo.DeleteResourceRelationData{
			ProjectId:        projectId,
			IssueIds:         couldDeleteIssueIds,
			RecycleVersionId: recycleVersionId,
			SourceTypes:      []int{consts.OssPolicyTypeIssueResource, consts.OssPolicyTypeLesscodeResource},
		},
	})
	if resourceRelationResp.Failure() {
		log.Error(resourceRelationResp.Error())
		return nil, resourceRelationResp.Error()
	}
	// 删除工时
	if err := DeleteWorkHourForIssues(orgId, couldDeleteIssueIds, recycleVersionId); err != nil {
		log.Errorf("[DeleteIssueBatch] orgId: %d, err: %v", orgId, err)
		return nil, err
	}

	asyn.Execute(func() {
		//beforeFollowerIds, err1 := GetRelationInfoByIssueIds(allParentIds, []int{consts.IssueRelationTypeFollower})
		//if err1 != nil {
		//	log.Error(err1)
		//	return
		//}
		//beforeOwnerIds, err1 := GetRelationInfoByIssueIds(allParentIds, []int{consts.IssueRelationTypeOwner})
		//if err1 != nil {
		//	log.Error(err1)
		//	return
		//}
		//ownerMap := map[int64][]int64{}
		//for _, ownerId := range beforeOwnerIds {
		//	ownerMap[ownerId.IssueId] = append(ownerMap[ownerId.IssueId], ownerId.RelationId)
		//}
		//followerMap := map[int64][]int64{}
		//for _, id := range beforeFollowerIds {
		//	followerMap[id.IssueId] = append(followerMap[id.IssueId], id.RelationId)
		//}
		for _, issueBo := range issueBos {
			issueTrendsBo := &bo.IssueTrendsBo{
				PushType:      consts.PushTypeDeleteIssue,
				OrgId:         orgId,
				OperatorId:    operatorId,
				DataId:        issueBo.DataId,
				IssueId:       issueBo.Id,
				ParentIssueId: issueBo.ParentId,
				ProjectId:     issueBo.ProjectId,
				TableId:       tableId,
				PriorityId:    issueBo.PriorityId,
				ParentId:      issueBo.ParentId,
				IssueTitle:    issueBo.Title,
				IssueStatusId: issueBo.Status,
				SourceChannel: sourceChannel,
			}
			//if _, ok := ownerMap[issueBo.Id]; ok {
			//	issueTrendsBo.BeforeOwner = ownerMap[issueBo.Id]
			//}
			if _, ok := childIssueIds[issueBo.Id]; ok {
				issueTrendsBo.IssueChildren = childIssueIds[issueBo.Id]
			}
			//if _, ok := followerMap[issueBo.Id]; ok {
			//	issueTrendsBo.BeforeChangeFollowers = followerMap[issueBo.Id]
			//}

			asyn.Execute(func() {
				PushIssueTrends(issueTrendsBo)
			})
			asyn.Execute(func() {
				PushIssueThirdPlatformNotice(issueTrendsBo)
			})
		}
	})

	return couldDeleteIssueIds, nil
}

func deleteIssuesAndRelation(orgId, projectId, appId, tableId, operatorId, recycleVersionId int64, issueIds []int64, tx sqlbuilder.Tx) errs.SystemErrorInfo {
	//_, err := mysql.TransUpdateSmartWithCond(tx, consts.TableIssue, db.Cond{
	//	consts.TcOrgId:    orgId,
	//	consts.TcId:       db.In(issueIds),
	//	consts.TcIsDelete: consts.AppIsNoDelete,
	//}, mysql.Upd{
	//	consts.TcUpdator:  operatorId,
	//	consts.TcVersion:  recycleVersionId,
	//	consts.TcIsDelete: consts.AppIsDeleted,
	//})
	//if err != nil {
	//	return errs.BuildSystemErrorInfo(errs.MysqlOperateError, err)
	//}

	//issueDetail := &po.PpmPriIssueDetail{}
	//_, err = mysql.TransUpdateSmartWithCond(tx, issueDetail.TableName(), db.Cond{
	//	consts.TcOrgId:    orgId,
	//	consts.TcIssueId:  db.In(issueIds),
	//	consts.TcIsDelete: consts.AppIsNoDelete,
	//}, mysql.Upd{
	//	consts.TcUpdator:  operatorId,
	//	consts.TcVersion:  recycleVersionId,
	//	consts.TcIsDelete: consts.AppIsDeleted,
	//})
	//if err != nil {
	//	return errs.BuildSystemErrorInfo(errs.MysqlOperateError, err)
	//}
	// 删除关联关系前，进行一些操作。尝试删除日程
	if err := DeleteCalendarEventBatch(orgId, projectId, issueIds); err != nil {
		errMsg := fmt.Sprintf("DeleteIssueBatch 批量删除任务，后续处理，删除日程异常：%v", err)
		log.Info(errMsg)
		// 如果是“查询日历信息异常。”，则表示项目配置中没有设置日历相关的配置，此时不做处理。
		if err.Error() != "查询日历信息异常。" {
			log.Error(errMsg)
			commonfacade.DingTalkInfo(commonvo.DingTalkInfoReqVo{
				LogData: map[string]interface{}{
					"app": config.GetApplication().Name,
					"env": config.GetEnv(),
					"msg": errMsg,
					"tag": "【告警日志】",
					"err": err.Error(),
				},
			})
		}
	}

	err3 := DeleteAllIssueRelation(tx, operatorId, orgId, issueIds, recycleVersionId)
	if err3 != nil {
		log.Error(err3)
		return errs.BuildSystemErrorInfo(errs.IssueDomainError, err3)
	}

	//无码回收
	recycleResp := formfacade.LessRecycleIssue(formvo.LessRecycleIssueReq{
		AppId:    appId,
		OrgId:    orgId,
		UserId:   operatorId,
		IssueIds: issueIds,
		TableId:  tableId,
	})
	if recycleResp.Failure() {
		log.Error(recycleResp.Error())
	}

	return nil
}

// DeleteTableIssues 删除一个table，把所有任务都删除，不进入回收站
func DeleteTableIssues(orgId, operatorId, projectId, appId, tableId int64, templateFlag int) ([]int64, errs.SystemErrorInfo) {
	conditions := []*tablePb.Condition{
		GetRowsCondition(consts.BasicFieldTableId, tablePb.ConditionType_equal, cast.ToString(tableId), nil),
	}
	list, err := GetIssueIdList(orgId, conditions, 0, 0)
	if err != nil {
		return nil, err
	}
	issueIds := make([]int64, 0, len(list))

	if len(list) == 0 {
		return issueIds, nil
	}

	for _, data := range list {
		issueIds = append(issueIds, data.Id)
	}
	// 删除附件和任务关系、删除附件关联关系、删除回收站附件
	errSys := DeleteAttachmentsForOneTable(orgId, operatorId, projectId, tableId)
	if errSys != nil {
		log.Errorf("[DeleteTableIssues] DeleteAttachmentsForOneTable err:%v, orgId:%v, userId:%v, tableId:%v",
			errSys, orgId, operatorId, tableId)
		return nil, errSys
	}

	err2 := mysql.TransX(func(tx sqlbuilder.Tx) error {
		if len(issueIds) > 0 {
			err := deleteIssuesAndRelation(orgId, projectId, appId, tableId, operatorId, 0, issueIds, tx)
			if err != nil {
				return err
			}

			err = deleteRecycleRecord(orgId, operatorId, projectId, tableId, issueIds, tx)
			if err != nil {
				return err
			}
		}

		return nil
	})

	if err2 != nil {
		log.Errorf("[DeleteTableIssues] err:%v", err2)
		return nil, errs.MysqlOperateError
	}

	if templateFlag == consts.TemplateFalse {
		errCache := SetDeletedIssuesNum(orgId, int64(len(issueIds)))
		if errCache != nil {
			log.Errorf("[DeleteTableIssues]SetDeletedIssuesNum err:%v", errCache)
		}
	}

	return issueIds, nil
}

func deleteRecycleRecord(orgId, userId, projectId, tableId int64, issueIds []int64, tx sqlbuilder.Tx) errs.SystemErrorInfo {
	_, recordsErr := mysql.TransUpdateSmartWithCond(tx, consts.TableRecycleBin, db.Cond{
		consts.TcIsDelete:     consts.AppIsNoDelete,
		consts.TcProjectId:    projectId,
		consts.TcRelationType: consts.RecycleTypeIssue,
		consts.TcRelationId:   db.In(issueIds),
		consts.TcOrgId:        orgId,
	}, mysql.Upd{
		consts.TcIsDelete: consts.AppIsDeleted,
		consts.TcUpdator:  userId,
	})
	if recordsErr != nil {
		log.Errorf("[deleteRecycleRecord] UpdateSmartWithCond failed:%v", recordsErr)
		return errs.MysqlOperateError
	}

	return nil
}

func ConvertStatusInfoToHomeIssueStatusInfo(statusInfo bo.CacheProcessStatusBo) *status.StatusInfoBo {
	homeStatusInfo := &status.StatusInfoBo{}
	homeStatusInfo.ID = statusInfo.StatusId
	homeStatusInfo.Name = statusInfo.Name
	homeStatusInfo.Type = statusInfo.StatusType
	homeStatusInfo.BgStyle = statusInfo.BgStyle
	homeStatusInfo.FontStyle = statusInfo.FontStyle
	homeStatusInfo.Sort = statusInfo.Sort
	return homeStatusInfo
}

func GetHomeProjectInfoBo(orgId, projectId int64) (*bo.HomeIssueProjectInfoBo, errs.SystemErrorInfo) {
	projectCacheInfo, err := LoadProjectAuthBo(orgId, projectId)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	projectInfo := ConvertProjectCacheInfoToHomeIssueProjectInfo(*projectCacheInfo)
	return projectInfo, nil
}

func ConvertProjectCacheInfoToHomeIssueProjectInfo(projectCacheInfo bo.ProjectAuthBo) *bo.HomeIssueProjectInfoBo {
	projectInfo := &bo.HomeIssueProjectInfoBo{}
	projectInfo.Name = projectCacheInfo.Name
	projectInfo.ID = projectCacheInfo.Id
	projectInfo.IsFilling = projectCacheInfo.IsFilling
	projectInfo.ProjectTypeId = projectCacheInfo.ProjectType
	return projectInfo
}

func NoticeIssueAudit(orgId int64, userId int64, orgBaseInfo *bo.BaseOrgInfoBo, issueBo *bo.IssueBo, notNeedNoticeIds []int64) {
	if orgBaseInfo.OutOrgId == "" {
		return
	}

	var allUserIds []int64
	for _, auditorId := range issueBo.AuditorIdsI64 {
		if ok, _ := slice.Contain(notNeedNoticeIds, auditorId); ok {
			continue
		}
		allUserIds = append(allUserIds, auditorId)
	}

	allUserIds = slice.SliceUniqueInt64(allUserIds)
	if len(allUserIds) == 0 {
		return
	}

	// 查询用户信息
	resp := orgfacade.GetBaseUserInfoBatch(orgvo.GetBaseUserInfoBatchReqVo{
		OrgId:   orgId,
		UserIds: allUserIds,
	})
	if resp.Failure() {
		log.Errorf("[NoticeIssueAudit] GetBaseUserInfoBatch err: %v", resp.Error())
		return
	}
	if len(resp.BaseUserInfos) == 0 {
		return
	}
	userInfoMap := map[int64]*bo.BaseUserInfoBo{}
	for i, info := range resp.BaseUserInfos {
		userInfoMap[info.UserId] = &resp.BaseUserInfos[i]
	}

	// 机器人将信息回复给用户
	var allOutUserIds []string
	for _, id := range allUserIds {
		if userId == id {
			//自己操作的就不需要推送给自己了
			continue
		}
		if info, ok := userInfoMap[id]; ok {
			if info.OutUserId != "" {
				allOutUserIds = append(allOutUserIds, info.OutUserId)
			}
		}
	}
	if len(allOutUserIds) == 0 {
		return
	}
	cd, err := getAuditIssueNoticeCard(userId, issueBo)
	if err != nil {
		log.Errorf("[NoticeIssueAudit] getAuditIssueNoticeCard err: %v", err)
		return
	}
	cardMsg := &commonvo.PushCard{
		OrgId:         orgId,
		OutOrgId:      orgBaseInfo.OutOrgId,
		SourceChannel: orgBaseInfo.SourceChannel,
		OpenIds:       allOutUserIds,
		CardMsg:       cd,
	}
	errSys := card.PushCard(orgId, cardMsg)
	if errSys != nil {
		log.Errorf("[NoticeIssueAudit] PushCard err:%v", errSys)
		return
	}
	//card.SendCardMeta(orgBaseInfo.SourceChannel, orgBaseInfo.OutOrgId, cd, allOutUserIds)
}

// getAuditIssueNoticeCard 状态变更后，发起审批（符合条件时）
// starterId 审批发起人 id
func getAuditIssueNoticeCard(starterId int64, issue *bo.IssueBo) (*pushPb.TemplateCard, errs.SystemErrorInfo) {
	if issue.Title == "" {
		return nil, errs.CardTitleEmpty
	}
	infoBox, err := getCardInfoForStartAudit(starterId, issue)
	if err != nil {
		log.Errorf("[getAuditIssueNoticeCard] err: %v, issueId: %d", err, issue.Id)
		return nil, err
	}

	return card.GetCardAuditorWhenUpdateIssue(infoBox)
}

// getCardInfoForStartAudit 获取并向 cardInfo 设置组装卡片时的信息
// starter 审批发起人的 id
func getCardInfoForStartAudit(starter int64, issue *bo.IssueBo) (*projectvo.BaseInfoBoxForIssueCard, errs.SystemErrorInfo) {
	infoObj := &projectvo.BaseInfoBoxForIssueCard{}

	infoObj.IssueInfo = *issue
	infoObj.OperateUserId = starter
	// 查询任务负责人
	userIds := make([]int64, 0, len(infoObj.IssueInfo.OwnerIdI64)+1)
	userIds = append(userIds, starter)
	userIds = append(userIds, infoObj.IssueInfo.OwnerIdI64...)
	userInfoArr, err := orgfacade.GetBaseUserInfoBatchRelaxed(issue.OrgId, userIds)
	if err != nil {
		log.Errorf("[GetAndSetCardInfoForUrgeAudit] 查询组织 %d 用户信息出现异常 %v", issue.OrgId, err)
		return infoObj, err
	}
	userMap := make(map[int64]bo.BaseUserInfoBo, len(userInfoArr))
	for _, user := range userInfoArr {
		userMap[user.UserId] = user
	}
	if opUser, ok := userMap[starter]; ok {
		infoObj.OperateUser = opUser
	}
	//ownerUserNameArr := make([]string, 0)
	ownerInfos := make([]*bo.BaseUserInfoBo, 0)
	for _, uid := range infoObj.IssueInfo.OwnerIdI64 {
		if user, ok := userMap[uid]; ok {
			//ownerUserNameArr = append(ownerUserNameArr, user.Name)
			ownerInfos = append(ownerInfos, &user)
		}
	}
	infoObj.OwnerInfos = ownerInfos

	projectInfo, err1 := LoadProjectAuthBo(issue.OrgId, issue.ProjectId)
	if err1 != nil {
		log.Error(err1)
		return infoObj, err1
	}
	infoObj.ProjectAuthBo = *projectInfo

	tableColumns, err := GetTableColumnsMap(issue.OrgId, infoObj.IssueInfo.TableId, nil)
	if err != nil {
		log.Errorf("[GetAndSetCardInfoForUrgeAudit] 获取表头失败 org:%d proj:%d table:%d starter: %d, err: %v",
			issue.OrgId, issue.ProjectId, infoObj.IssueInfo.TableId, starter, err)
		return infoObj, err
	}
	infoObj.ProjectTableColumn = tableColumns
	headers := make(map[string]lc_table.LcCommonField, 0)
	copyer.Copy(tableColumns, &headers)
	infoObj.TableColumnMap = headers

	// 查询父任务信息
	if issue.ParentId > 0 {
		parents, err := GetIssueInfosLc(issue.OrgId, 0, []int64{issue.ParentId})
		if err != nil {
			log.Errorf("[getCardInfoForStartAudit] GetIssueInfosLc err: %v, parentId: %d", err, issue.ParentId)
			return nil, err
		}
		if len(parents) > 0 {
			parent := parents[0]
			infoObj.ParentIssue = *parent
		}
	}

	// 查询任务所在表名
	if infoObj.IssueInfo.TableId != 0 {
		tableInfo, err := GetTableByTableId(issue.OrgId, starter, infoObj.IssueInfo.TableId)
		if err != nil {
			log.Errorf("[GetAndSetCardInfoForUrgeAudit] GetTableByTableId err:%v, userId:%d, tableId:%d", err, starter,
				infoObj.IssueInfo.TableId)
			return infoObj, err
		}
		infoObj.IssueTableInfo = *tableInfo
	} else {
		err := errs.InvalidTableId
		log.Errorf("[GetAndSetCardInfoForUrgeAudit] err: %v, issueId: %d", err, issue.Id)
		return infoObj, err
	}

	orgBaseResp, errSys := orgfacade.GetBaseOrgInfoRelaxed(issue.OrgId)
	if errSys != nil {
		log.Errorf("[getCardInfoForStartAudit] err:%v, orgId:%v", errSys, issue.OrgId)
		return infoObj, errSys
	}
	infoObj.SourceChannel = orgBaseResp.SourceChannel

	// 查看详情、PC端查看等按钮的 url 链接
	links := GetIssueLinks(orgBaseResp.SourceChannel, issue.OrgId, issue.Id)
	infoObj.IssueInfoUrl = links.SideBarLink
	infoObj.IssuePcUrl = links.Link

	return infoObj, nil
}

// GetTodoCard 构造待办卡片
func GetTodoCard(orgId, userId int64, sourceChannel string, todo *automationPb.Todo) (*pushPb.TemplateCard, errs.SystemErrorInfo) {
	cardInfo := &projectvo.TodoCard{
		TodoType: todo.Type,
	}

	// issue
	issue, errSys := GetIssueInfoLc(orgId, userId, todo.IssueId)
	if errSys != nil {
		log.Errorf("[GetTodoCard] GetIssueInfoLc failed. org:%d user:%d issue:%d, err: %v", orgId, userId, todo.IssueId, errSys)
		return nil, errSys
	}
	cardInfo.Issue = issue

	triggerUserId := todo.TriggerUserId

	var userIds []int64
	userIds = append(userIds, triggerUserId)
	userIds = append(userIds, cardInfo.Issue.OwnerIdI64...)
	userInfoArr, err := orgfacade.GetBaseUserInfoBatchRelaxed(issue.OrgId, userIds)
	if err != nil {
		log.Errorf("[GetTodoCard] 查询组织 %d 用户信息出现异常 %v", issue.OrgId, err)
		return nil, err
	}
	userMap := make(map[int64]*bo.BaseUserInfoBo, len(userInfoArr))
	for i, user := range userInfoArr {
		userMap[user.UserId] = &userInfoArr[i]
	}
	if u, ok := userMap[triggerUserId]; ok {
		cardInfo.Trigger = u
	}

	// owners
	owners := make([]*bo.BaseUserInfoBo, 0)
	for _, uid := range cardInfo.Issue.OwnerIdI64 {
		if user, ok := userMap[uid]; ok {
			owners = append(owners, user)
		}
	}
	cardInfo.Owners = owners

	// project
	project, errSys := GetProject(issue.OrgId, issue.ProjectId)
	if errSys != nil {
		log.Errorf("[GetTodoCard] GetProject failed. org:%d project:%d, err: %v", issue.OrgId, issue.ProjectId, errSys)
		return nil, errSys
	}
	cardInfo.Project = project

	// tableColumns
	tableColumns, err := GetTableColumnsMap(issue.OrgId, issue.TableId, nil)
	if err != nil {
		log.Errorf("[GetTodoCard] 获取表头失败 org:%d proj:%d table:%d triggerUserId: %d, err: %v",
			issue.OrgId, issue.ProjectId, cardInfo.Issue.TableId, triggerUserId, err)
		return nil, err
	}
	cardInfo.TableColumns = tableColumns

	// table
	if cardInfo.Issue.TableId != 0 {
		tableInfo, err := GetTableByTableId(issue.OrgId, triggerUserId, cardInfo.Issue.TableId)
		if err != nil {
			log.Errorf("[GetTodoCard] GetTableByTableId err:%v, userId:%d, tableId:%d", err, triggerUserId,
				cardInfo.Issue.TableId)
			return nil, err
		}
		cardInfo.Table = tableInfo
	}

	// 查看详情、PC端查看等按钮的 url 链接
	links := GetTodoLinks(sourceChannel, issue.OrgId, cast.ToInt64(todo.Id))
	cardInfo.TodoSideBarUrl = links.SideBarLink
	cardInfo.TodoUrl = links.Link

	return card.GetCardTodo(cardInfo)
}

func ChangeIssueChildren(operatorId int64, issueBo *bo.IssueBo, parentId int64, parentPath string) ([]int64, []*bo.IssueBo, errs.SystemErrorInfo) {
	if issueBo.ParentId == parentId {
		return nil, nil, nil
	}
	fromParentId := issueBo.ParentId
	targetParentId := parentId

	//判断父任务和当前任务是否有交叉（即：父任务是不是当前任务的下级）
	parentArr := strings.Split(parentPath, ",")
	if ok, _ := slice.Contain(parentArr, strconv.FormatInt(issueBo.Id, 10)); ok {
		return nil, nil, errs.ParentIsChildIssue
	}

	//childPos := &[]po.PpmPriIssue{}
	////查询是否存在子任务
	//err := mysql.SelectAllByCond(consts.TableIssue, db.Cond{
	//	consts.TcPath:     db.Like(fmt.Sprintf("%s%d,%s", issueBo.Path, issueBo.Id, "%")),
	//	consts.TcIsDelete: consts.AppIsNoDelete,
	//	consts.TcOrgId:    issueBo.OrgId,
	//}, childPos)
	//if err != nil {
	//	log.Error(err)
	//	return nil, nil, errs.MysqlOperateError
	//}
	//
	////子任务集合
	//childrenIds := []int64{}
	//if len(*childPos) > 0 {
	//	for _, data := range *childPos {
	//		childrenIds = append(childrenIds, data.Id)
	//	}
	//}

	filterColumns := []string{
		lc_helper.ConvertToFilterColumn(consts.BasicFieldIssueId),
		lc_helper.ConvertToFilterColumn(consts.TcPath),
	}
	condition := &tablePb.Condition{Type: tablePb.ConditionType_and}
	condition.Conditions = GetNoRecycleCondition(
		GetRowsCondition(consts.TcPath, tablePb.ConditionType_like, fmt.Sprintf(",%d,", issueBo.Id), nil))
	childInfos, errSys := GetIssueInfosMapLc(issueBo.OrgId, operatorId, condition, filterColumns, -1, -1)
	if errSys != nil {
		log.Errorf("[ChangeIssueChildren] GetIssueInfosMapLc err:%v, orgId:%v, issueId:%v", errSys, issueBo.OrgId, issueBo.Id)
		return nil, nil, errSys
	}
	var issueBos []*bo.IssueBo
	childrenIds := make([]int64, 0, len(childInfos))
	for _, data := range childInfos {
		issue, errSys := ConvertIssueDataToIssueBo(data)
		if errSys != nil {
			log.Errorf("[ChangeIssueChildren]ConvertIssueDataToIssueBo err:%v, orgId:%v, issueId:%v", errSys, issueBo.OrgId, issueBo.Id)
			return nil, nil, errSys
		}
		childrenIds = append(childrenIds, issue.Id)
		issueBos = append(issueBos, issue)
	}

	newPath := "0,"
	if parentId != 0 {
		newPath = parentPath + strconv.FormatInt(parentId, 10) + ","
	}

	//判断层级是否超出9级
	if parentId != 0 {
		if strings.Count(newPath, ",") > consts.IssueLevel {
			return nil, nil, errs.IssueLevelOutLimit
		}

		//转化后相差的等级
		changeLevel := strings.Count(newPath, ",") - strings.Count(issueBo.Path, ",")
		for _, issue := range issueBos {
			if strings.Count(issue.Path, ",")+changeLevel > consts.IssueLevel {
				return nil, nil, errs.IssueLevelOutLimit
			}
		}
	}

	appId, appIdErr := GetAppIdFromProjectId(issueBo.OrgId, issueBo.ProjectId)
	if appIdErr != nil {
		log.Error(appIdErr)
		return nil, nil, appIdErr
	}

	// 变更之前先查询老数据信息
	allIssueIds := append(childrenIds, issueBo.Id)
	oldIssueBos, errSys := GetIssueInfosLc(issueBo.OrgId, operatorId, allIssueIds, getMoveRelateColumnIds()...)
	if errSys != nil {
		log.Errorf("[ChangeIssueChildren] GetIssueInfosLc err: %v", errSys)
		return nil, nil, errSys
	}

	transErr := mysql.TransX(func(tx sqlbuilder.Tx) error {
		// 更新父任务
		upd := mysql.Upd{
			consts.TcParentId: parentId,
			consts.TcPath:     newPath,
			consts.TcUpdator:  operatorId,
		}
		//_, updateErr := mysql.UpdateSmartWithCond(consts.TableIssue, db.Cond{
		//	consts.TcId: issueBo.Id,
		//}, upd)
		//if updateErr != nil {
		//	log.Error(updateErr)
		//	return errs.MysqlOperateError
		//}

		// 更新子任务
		//_, updateErr = mysql.UpdateSmartWithCond(consts.TableIssue, db.Cond{
		//	consts.TcId: db.In(childrenIds),
		//}, mysql.Upd{
		//	consts.TcPath:    db.Raw("replace(`path`, \"" + issueBo.Path + "\", '" + newPath + "')"),
		//	consts.TcUpdator: operatorId,
		//})
		//if updateErr != nil {
		//	log.Error(updateErr)
		//	return errs.MysqlOperateError
		//}

		// 无码更新: 父任务的Path和ParentId
		updData := slice2.CaseCamelCopy(upd)
		updData[consts.BasicFieldIssueId] = issueBo.Id
		resp := formfacade.LessUpdateIssue(formvo.LessUpdateIssueReq{
			AppId:   appId,
			OrgId:   issueBo.OrgId,
			UserId:  operatorId,
			TableId: issueBo.TableId,
			Form: []map[string]interface{}{
				updData,
			},
		})
		if resp.Failure() {
			log.Error(resp.Error())
			return resp.Error()
		}

		// 无码更新：子任务的Path
		batchResp := formfacade.LessUpdateIssueBatchRaw(&formvo.LessUpdateIssueBatchReq{
			OrgId:  issueBo.OrgId,
			AppId:  appId,
			UserId: operatorId,
			Condition: vo.LessCondsData{
				Type: consts.ConditionAnd,
				Conds: []*vo.LessCondsData{
					{
						Type:   consts.ConditionEqual,
						Value:  issueBo.OrgId,
						Column: lc_helper.ConvertToCondColumn(consts.BasicFieldOrgId),
					},
					{
						Type:   consts.ConditionIn,
						Values: childrenIds,
						Column: lc_helper.ConvertToCondColumn(consts.BasicFieldIssueId),
					},
				},
			},
			Sets: []datacenter.Set{
				{
					Column:          lc_helper.ConvertToCondColumn(consts.BasicFieldPath),
					Value:           fmt.Sprintf("regexp_replace(\"%s\",'%s','%s')", consts.BasicFieldPath, issueBo.Path, newPath),
					Type:            consts.SetTypeNormal,
					Action:          consts.SetActionSet,
					WithoutPretreat: true,
				},
			},
		})
		if batchResp.Failure() {
			log.Error(batchResp.Error())
			return batchResp.Error()
		}

		return nil
	})
	if transErr != nil {
		log.Error(transErr)
		return nil, nil, errs.MysqlOperateError
	}

	issueBo.ParentId = parentId
	issueBo.Path = newPath

	// 保存动态
	asyn.Execute(func() {
		var oldName, newName string
		if fromParentId > 0 {
			issues, errSys := GetIssueInfosLc(issueBo.OrgId, operatorId, []int64{fromParentId}, lc_helper.ConvertToFilterColumn(consts.BasicFieldTitle))
			if errSys != nil {
				log.Error(errSys)
				return
			}
			if len(issues) > 0 {
				oldName = issues[0].Title
				if strings.Trim(oldName, " ") == "" {
					oldName = "未命名记录"
				}
			}
		}
		if targetParentId > 0 {
			issues, errSys := GetIssueInfosLc(issueBo.OrgId, operatorId, []int64{targetParentId}, lc_helper.ConvertToFilterColumn(consts.BasicFieldTitle))
			if errSys != nil {
				log.Error(errSys)
				return
			}
			if len(issues) > 0 {
				newName = issues[0].Title
				if strings.Trim(newName, " ") == "" {
					newName = "未命名记录"
				}
			}
		}
		pushType := consts.PushTypeUpdateIssueParent
		var changeList []bo.TrendChangeListBo
		changeList = append(changeList, bo.TrendChangeListBo{
			Field:     consts.BasicFieldParentId,
			FieldName: consts.Parent,
			OldValue:  oldName,
			NewValue:  newName,
		})

		issueTrendsBo := &bo.IssueTrendsBo{
			PushType:      pushType,
			OrgId:         issueBo.OrgId,
			OperatorId:    operatorId,
			DataId:        issueBo.DataId,
			IssueId:       issueBo.Id,
			ParentIssueId: issueBo.ParentId,
			ProjectId:     issueBo.ProjectId,
			TableId:       issueBo.TableId,
			PriorityId:    issueBo.PriorityId,
			IssueTitle:    issueBo.Title,
			ParentId:      issueBo.ParentId,
			OldValue:      cast.ToString(fromParentId),
			NewValue:      cast.ToString(targetParentId),
			Ext: bo.TrendExtensionBo{
				ObjName:    issueBo.Title,
				ChangeList: changeList,
			},
		}
		PushIssueTrends(issueTrendsBo)
	})

	return allIssueIds, oldIssueBos, nil
}

func GetIssueAuthBo(issueBo *bo.IssueBo, currentUserId int64) (*bo.IssueAuthBo, errs.SystemErrorInfo) {
	issueId := issueBo.Id

	//ownerIds, err2 := GetIssueRelationIdsByRelateType(issueBo.OrgId, issueId, consts.IssueRelationTypeOwner)
	//if err2 != nil {
	//	log.Error(err2)
	//	return nil, err2
	//}
	issueAuthBo := &bo.IssueAuthBo{
		Id:        issueId,
		Owner:     issueBo.OwnerIdI64,
		Creator:   issueBo.Creator,
		ProjectId: issueBo.ProjectId,
		Status:    issueBo.Status,
		TableId:   issueBo.TableId,
	}

	// 2019-12-24-nico: 产品想要支持父任务负责人操作所有子任务
	if issueBo.ParentId != 0 {
		issueBos, err := GetIssueInfosLc(issueBo.OrgId, currentUserId, []int64{issueBo.ParentId})
		if err != nil {
			log.Errorf("[GetIssueAuthBo] err: %v, parentId: %d", err, issueBo.ParentId)
			return issueAuthBo, err
		}
		if len(issueBos) <= 0 {
			return issueAuthBo, err
		}
		parentIssueBo := issueBos[0]
		parentOwnerIds, errSys := businees.LcMemberToUserIdsWithError(parentIssueBo.OwnerId)
		if errSys != nil {
			log.Errorf("[GetIssueAuthBo] parent issue ownerId: %v, err: %v", parentIssueBo.OwnerId, errSys)
			return nil, errSys
		}

		if ok, _ := slice.Contain(parentOwnerIds, currentUserId); ok {
			issueAuthBo.Owner = append(issueAuthBo.Owner, currentUserId)
		}
	}

	return issueAuthBo, nil
}

// GetIssueInfosMapLcByIssueIds 从无码获取多条数据
func GetIssueInfosMapLcByIssueIds(orgId, userId int64, issueIds []int64, filterColumns ...string) ([]map[string]interface{}, errs.SystemErrorInfo) {
	condition := GetRowsCondition(consts.BasicFieldIssueId, tablePb.ConditionType_in, nil, issueIds)
	return GetIssueInfosMapLc(orgId, userId, condition, filterColumns, 1, 2000)
}

// GetIssueInfosMapLcByDataIds 从无码获取多条数据
func GetIssueInfosMapLcByDataIds(orgId, userId int64, dataIds []int64, filterColumns ...string) ([]map[string]interface{}, errs.SystemErrorInfo) {
	condition := GetRowsCondition(consts.BasicFieldId, tablePb.ConditionType_in, nil, dataIds)
	return GetIssueInfosMapLc(orgId, userId, condition, filterColumns, 1, 2000)
}

func filterColumnUnique(s []string, noExtraIdFlag ...bool) []string {
	if s == nil {
		return nil
	}
	res := make([]string, 0)
	exist := make(map[string]bool)
	for _, s2 := range s {
		if _, ok := exist[s2]; ok {
			continue
		}
		res = append(res, s2)
		exist[s2] = true
	}
	if len(noExtraIdFlag) == 0 {
		if _, ok := exist[consts.BasicFieldId]; !ok {
			res = append(res, lc_helper.ConvertToFilterColumn(consts.BasicFieldId))
		}
		if _, ok := exist[consts.BasicFieldIssueId]; !ok {
			res = append(res, lc_helper.ConvertToFilterColumn(consts.BasicFieldIssueId))
		}
	}
	return res
}

func GetTableIssueMaxOrder(orgId, userId int64) (float64, errs.SystemErrorInfo) {
	lessResp, err := GetRawRows(orgId, userId, &tablePb.ListRawRequest{
		FilterColumns: []string{lc_helper.ConvertToFilterColumn(consts.BasicFieldOrder)},
		Orders: []*tablePb.Order{
			{
				Asc:    false,
				Column: lc_helper.ConvertToCondColumn(consts.BasicFieldOrder),
			},
		},
		Page: 1,
		Size: 1,
	})
	if err != nil {
		log.Errorf("[GetTableIssueLastOrder] LessIssueRawList failed, orgId: %d,  err: %v", orgId, err.Error())
		return 0, err
	}
	if len(lessResp.Data) == 0 {
		return 0, nil
	}

	return cast.ToFloat64(lessResp.Data[0][consts.BasicFieldOrder]), nil
}

func GetTableIssueMaxOrderOfTable(orgId, userId, tableId int64) (float64, errs.SystemErrorInfo) {
	lessResp, err := GetRawRows(orgId, userId, &tablePb.ListRawRequest{
		Condition:     GetRowsCondition(consts.BasicFieldTableId, tablePb.ConditionType_equal, tableId, nil),
		FilterColumns: []string{lc_helper.ConvertToFilterColumn(consts.BasicFieldOrder)},
		Orders: []*tablePb.Order{
			{
				Asc:    false,
				Column: lc_helper.ConvertToCondColumn(consts.BasicFieldOrder),
			},
		},
		Page: 1,
		Size: 1,
	})
	if err != nil {
		log.Errorf("[GetTableIssueLastOrder] LessIssueRawList failed, orgId: %d,  err: %v", orgId, err.Error())
		return 0, err
	}
	if len(lessResp.Data) == 0 {
		return 0, nil
	}

	return cast.ToFloat64(lessResp.Data[0][consts.BasicFieldOrder]), nil
}

func GetIssueInfosMapLc(orgId, userId int64, condition *tablePb.Condition, filterColumns []string, page, size int64, noExtraIdFlag ...bool) ([]map[string]interface{}, errs.SystemErrorInfo) {
	filterColumns = filterColumnUnique(filterColumns, noExtraIdFlag...)
	req := &tablePb.ListRawRequest{
		FilterColumns: filterColumns,
		Condition:     condition,
		Page:          int32(page),
		Size:          int32(size),
	}

	reply, err := GetRawRows(orgId, userId, req)
	if err != nil {
		return nil, err
	}

	return reply.Data, nil
}

func ConvertIssueDataToIssueBo(data map[string]interface{}) (*bo.IssueBo, errs.SystemErrorInfo) {
	issueBo := &bo.IssueBo{}
	err := copyer.Copy(data, issueBo)
	if err != nil {
		log.Errorf("[ConvertIssueDataToIssueBo] json转换错误, err:%v", err)
		return nil, errs.JSONConvertError
	}
	issueBo.Id = cast.ToInt64(data[consts.BasicFieldIssueId])
	//issueBo.IssueId = cast.ToInt64(data[consts.BasicFieldIssueId])
	issueBo.DataId = cast.ToInt64(data[consts.BasicFieldId])
	issueBo.TableId = cast.ToInt64(data[consts.BasicFieldTableId])
	issueBo.Status = cast.ToInt64(data[consts.BasicFieldIssueStatus])
	issueBo.AppId = cast.ToInt64(data[consts.BasicFieldAppId])
	if len(issueBo.Path) == 0 {
		issueBo.Path = "0,"
	}
	issueBo.OwnerIdI64 = businees.LcMemberToUserIds(issueBo.OwnerId)
	issueBo.AuditorIdsI64 = businees.LcMemberToUserIds(issueBo.AuditorIds)
	issueBo.FollowerIdsI64 = businees.LcMemberToUserIds(issueBo.FollowerIds)
	if issueBo.AuditStatusDetail == nil {
		issueBo.AuditStatusDetail = make(map[string]int)
	}
	issueBo.LessData = data
	return issueBo, nil
}

func GetIssueInfoLc(orgId, userId int64, issueId int64, filterColumns ...string) (*bo.IssueBo, errs.SystemErrorInfo) {
	issueLcDatas, errSys := GetIssueInfosMapLcByIssueIds(orgId, userId, []int64{issueId}, filterColumns...)
	if errSys != nil {
		log.Errorf("[GetIssueInfosLc] GetIssueInfosMapLcByIssueIds err: %v, issueId: %v", errSys, issueId)
		return nil, errSys
	}

	if len(issueLcDatas) < 1 {
		log.Errorf("[GetIssueInfosLc] GetIssueInfosMapLcByIssueIds err: %v, issueId: %v", errSys, issueId)
		return nil, errs.IssueNotExist
	}

	var issueBo *bo.IssueBo
	for _, data := range issueLcDatas {
		issueBo, errSys = ConvertIssueDataToIssueBo(data)
		if errSys != nil {
			log.Errorf("[GetIssueInfosLc] ConvertIssueDataToIssueBo err: %v, issueId: %v", errSys, issueId)
			return nil, errSys
		}
		break
	}
	return issueBo, nil
}

func GetIssueInfosLc(orgId, userId int64, issueIds []int64, filterColumns ...string) ([]*bo.IssueBo, errs.SystemErrorInfo) {
	issueLcDatas, errSys := GetIssueInfosMapLcByIssueIds(orgId, userId, issueIds, filterColumns...)
	if errSys != nil {
		log.Errorf("[GetIssueInfosLc] GetIssueInfosMapLcByIssueIds err: %v, issueIds: %v", errSys,
			json.ToJsonIgnoreError(issueIds))
		return nil, errSys
	}

	var issueBos []*bo.IssueBo
	for _, data := range issueLcDatas {
		issueBo, errSys := ConvertIssueDataToIssueBo(data)
		if errSys != nil {
			log.Errorf("[GetIssueInfosLc] ConvertIssueDataToIssueBo err: %v, issueIds: %s", errSys,
				json.ToJsonIgnoreError(issueIds))
			return nil, errSys
		}
		issueBos = append(issueBos, issueBo)
	}
	return issueBos, nil
}

func GetIssueInfosLcByDataIds(orgId, userId int64, dataIds []int64) ([]*bo.IssueBo, errs.SystemErrorInfo) {
	issueLcDatas, errSys := GetIssueInfosMapLcByDataIds(orgId, userId, dataIds)
	if errSys != nil {
		return nil, errSys
	}

	var issueBos []*bo.IssueBo
	for _, data := range issueLcDatas {
		issueBo, errSys := ConvertIssueDataToIssueBo(data)
		if errSys != nil {
			return nil, errSys
		}
		issueBos = append(issueBos, issueBo)
	}
	return issueBos, nil
}

func getMoveRelateColumnIds() []string {
	return []string{
		lc_helper.ConvertToFilterColumn(consts.BasicFieldOrgId),
		lc_helper.ConvertToFilterColumn(consts.BasicFieldAppId),
		lc_helper.ConvertToFilterColumn(consts.BasicFieldProjectId),
		lc_helper.ConvertToFilterColumn(consts.BasicFieldTableId),
		lc_helper.ConvertToFilterColumn(consts.BasicFieldParentId),
		lc_helper.ConvertToFilterColumn(consts.BasicFieldPath),
		lc_helper.ConvertToFilterColumn(consts.BasicFieldOrder),
		lc_helper.ConvertToFilterColumn(consts.BasicFieldTitle),
	}
}

// getCopyIssueDeleteColumnIds 复制任务的时候，必须清空以便重新设置正确值的字段
func getCopyIssueDeleteColumnIds() []string {
	// 暂不支持工时
	deleteColumnIds := []string{consts.ProBasicFieldWorkHour}
	deleteColumnIds = append(deleteColumnIds, consts.BasicFieldId)
	deleteColumnIds = append(deleteColumnIds, consts.BasicFieldOrgId)
	deleteColumnIds = append(deleteColumnIds, consts.BasicFieldAppId)
	deleteColumnIds = append(deleteColumnIds, consts.BasicFieldProjectId)
	deleteColumnIds = append(deleteColumnIds, consts.BasicFieldTableId)
	deleteColumnIds = append(deleteColumnIds, consts.BasicFieldIssueId)
	deleteColumnIds = append(deleteColumnIds, consts.BasicFieldCode)
	deleteColumnIds = append(deleteColumnIds, consts.BasicFieldOrder)
	deleteColumnIds = append(deleteColumnIds, consts.BasicFieldRecycleFlag)
	deleteColumnIds = append(deleteColumnIds, consts.BasicFieldCollaborators)
	return deleteColumnIds
}

func ReportMoveEvent(userId int64, oldIssue *bo.IssueBo, newData map[string]interface{}, userDepts map[string]*uservo.MemberDept) {
	var columnIds []string
	dataId := cast.ToInt64(newData[consts.BasicFieldId])
	issueId := cast.ToInt64(newData[consts.BasicFieldIssueId])
	orgId := cast.ToInt64(newData[consts.BasicFieldOrgId])
	appId := cast.ToInt64(newData[consts.BasicFieldAppId])
	projectId := cast.ToInt64(newData[consts.BasicFieldProjectId])
	tableId := cast.ToInt64(newData[consts.BasicFieldTableId])
	parentId := cast.ToInt64(newData[consts.BasicFieldParentId])
	path := cast.ToString(newData[consts.BasicFieldPath])
	order := cast.ToInt64(newData[consts.BasicFieldOrder])
	if oldIssue.ProjectId != projectId || oldIssue.AppId != appId {
		columnIds = append(columnIds, consts.BasicFieldProjectId)
		columnIds = append(columnIds, consts.BasicFieldAppId)
	}
	if oldIssue.TableId != tableId {
		columnIds = append(columnIds, consts.BasicFieldTableId)
	}
	if oldIssue.ParentId != parentId {
		columnIds = append(columnIds, consts.BasicFieldParentId)
	}
	if oldIssue.Path != path {
		columnIds = append(columnIds, consts.BasicFieldPath)
	}
	if oldIssue.Order != order {
		columnIds = append(columnIds, consts.BasicFieldOrder)
	}
	if len(columnIds) == 0 {
		return
	}

	oldData := oldIssue.LessData

	// 转成前端识别的模式。。
	AssembleDataIds(oldData)
	AssembleDataIds(newData)

	e := &commonvo.DataEvent{
		OrgId:          orgId,
		AppId:          appId,
		ProjectId:      projectId,
		TableId:        tableId,
		DataId:         dataId,
		IssueId:        issueId,
		UserId:         userId,
		Old:            oldData,
		New:            newData,
		UpdatedColumns: columnIds,
		UserDepts:      userDepts,
	}
	openTraceId, _ := threadlocal.Mgr.GetValue(consts.JaegerContextTraceKey)
	openTraceIdStr := cast.ToString(openTraceId)
	common.ReportDataEvent(msgPb.EventType_DataMoved, openTraceIdStr, e)
}

func getRelatingAppTableIdByProps(props map[string]interface{}) (appId int64, tableId int64) {
	if v, ok := props["relateAppId"]; ok {
		appId = cast.ToInt64(v)
	}
	if v, ok := props["relateTableId"]; ok {
		tableId = cast.ToInt64(v)
	}
	return
}

// ValidFieldMapping 检查字段隐射的合法性
func ValidFieldMapping(orgId, fromProjectId, fromTableId, toProjectId, toTableId int64,
	fromTableColumns, toTableColumns map[string]*projectvo.TableColumnData, fieldMapping map[string]string) map[string][]string {
	finalFieldMapping := make(map[string][]string)
	for toColumnId, fromColumnId := range fieldMapping {
		fromColumn, ok1 := fromTableColumns[fromColumnId]
		if !ok1 {
			log.Infof("[ValidFieldMapping] from column has no header, orgId: %v, tableId: %v, column: %v", orgId, fromTableId, fromColumnId)
			continue
		}
		toColumn, ok2 := toTableColumns[toColumnId]
		if !ok2 {
			log.Infof("[ValidFieldMapping] to column has no header, orgId: %v, tableId: %v, column: %v", orgId, toTableId, toColumnId)
			continue
		}

		// 工时暂不支持
		if fromColumn.Field.Type == consts.BasicFieldWorkHour || toColumn.Field.Type == consts.BasicFieldWorkHour {
			log.Infof("[ValidFieldMapping] workHour is not supported, orgId: %v, fromTableId: %v, fromColumn: %v, toTableId: %v, toColumn: %v",
				orgId, fromTableId, fromColumnId, toTableId, toColumnId)
			continue
		}

		// 暂时只支持同类型字段映射
		if fromColumn.Field.Type != toColumn.Field.Type {
			log.Infof("[ValidFieldMapping] columns' type are different, orgId: %v, fromTableId: %v, fromColumn: %v, toTableId: %v, toColumn: %v",
				orgId, fromTableId, fromColumnId, toTableId, toColumnId)
			continue
		}

		// 迭代只能项目内映射，并且只能迭代映射迭代
		if (fromColumnId == consts.BasicFieldIterationId || toColumnId == consts.BasicFieldIterationId) &&
			(fromProjectId != toProjectId || fromColumnId != toColumnId) {
			log.Infof("[ValidFieldMapping] not same project/column, cannot support iterationId, orgId: %v, fromProjectId: %v, fromColumn: %v, toProjectId: %v, toColumn: %v",
				orgId, fromProjectId, fromColumnId, toProjectId, toColumnId)
			continue
		}

		switch fromColumn.Field.Type {
		case consts.LcColumnFieldTypeRelating, consts.LcColumnFieldTypeSingleRelating:
			fromRelatingAppId, fromRelatingTableId := getRelatingAppTableIdByProps(fromColumn.Field.Props)
			toRelatingAppId, toRelatingTableId := getRelatingAppTableIdByProps(toColumn.Field.Props)
			if fromRelatingAppId != toRelatingAppId || fromRelatingTableId != toRelatingTableId {
				log.Infof("[ValidFieldMapping] relating column not mapping, orgId: %v, fromTableId: %v, fromColumn: %v, toTableId: %v, toColumn: %v",
					orgId, fromTableId, fromColumnId, toTableId, toColumnId)
				continue
			}
		}

		if _, ok := finalFieldMapping[fromColumnId]; ok {
			finalFieldMapping[fromColumnId] = append(finalFieldMapping[fromColumnId], toColumnId)
		} else {
			finalFieldMapping[fromColumnId] = []string{toColumnId}
		}
	}

	// 如果映射了issueStatus，那么把auditStatus和auditStatusDetail也带过去
	if _, ok := finalFieldMapping[consts.BasicFieldIssueStatus]; ok {
		finalFieldMapping[consts.BasicFieldAuditStatus] = []string{consts.BasicFieldAuditStatus}
		finalFieldMapping[consts.BasicFieldAuditStatusDetail] = []string{consts.BasicFieldAuditStatusDetail}
	}

	log.Infof("[ValidFieldMapping] final: %v", json.ToJsonIgnoreError(finalFieldMapping))

	return finalFieldMapping
}

// ApplyFieldMappingCreateMissingSelectOptions 字段映射 && 表头选项补齐信息收集
func ApplyFieldMappingCreateMissingSelectOptions(data, newData map[string]interface{}, fieldMapping map[string][]string,
	selectOptions map[string]map[string]*projectvo.ColumnSelectOption,
	oldTableColumns, newTableColumns map[string]*projectvo.TableColumnData,
	newProjectTypeId int64, isCreateMissingSelectOptions bool) {
	if fieldMapping == nil || len(fieldMapping) == 0 {
		return
	}

	// 先删一遍没有被映射的字段
	for oldK, _ := range data {
		if oldK == consts.TempFieldIssueId ||
			oldK == consts.TempFieldCode {
			continue
		}
		if _, ok := fieldMapping[oldK]; !ok {
			delete(newData, oldK)
		}
	}

	for oldK, _ := range data {
		if oldK == consts.TempFieldIssueId ||
			oldK == consts.TempFieldCode {
			continue
		}
		if v, ok := newData[oldK]; ok {
			if newKs, ok := fieldMapping[oldK]; ok {
				oldKRemains := false
				for _, newK := range newKs {
					newData[newK] = v
					if newK == oldK { // 映射到同名字段，老的字段就需要保留
						oldKRemains = true
					}

					// 找出需创建的选项值
					if isCreateMissingSelectOptions {
						oldColumn, ok1 := oldTableColumns[oldK]
						newColumn, ok2 := newTableColumns[newK]
						if ok1 && ok2 {
							if (oldColumn.Field.Type == consts.LcColumnFieldTypeSelect || oldColumn.Field.Type == consts.LcColumnFieldTypeMultiSelect ||
								(newProjectTypeId != consts.ProjectTypeNormalId && oldColumn.Field.Type == consts.LcColumnFieldTypeGroupSelect)) &&
								(newColumn.Field.Type == consts.LcColumnFieldTypeSelect || newColumn.Field.Type == consts.LcColumnFieldTypeMultiSelect ||
									(newProjectTypeId != consts.ProjectTypeNormalId && newColumn.Field.Type == consts.LcColumnFieldTypeGroupSelect)) {
								oldProps := &projectvo.FormConfigColumnFieldMultiselectProps{}
								newProps := &projectvo.FormConfigColumnFieldMultiselectProps{}
								copyer.Copy(oldColumn.Field.Props, oldProps)
								copyer.Copy(newColumn.Field.Props, newProps)
								var vIds []interface{}
								var oldOptions, newOptions []projectvo.ColumnSelectOption
								if len(oldProps.Select.Options) > 0 {
									oldOptions = oldProps.Select.Options
									vIds = append(vIds, v)
								}
								if len(oldProps.Multiselect.Options) > 0 {
									oldOptions = oldProps.Multiselect.Options
									is := cast.ToSlice(v)
									for _, i := range is {
										vIds = append(vIds, i)
									}
								}
								if len(oldProps.GroupSelect.Options) > 0 {
									oldOptions = oldProps.GroupSelect.Options
									vIds = append(vIds, v)
								}
								if len(newProps.Select.Options) > 0 {
									newOptions = newProps.Select.Options
								}
								if len(newProps.Multiselect.Options) > 0 {
									newOptions = newProps.Multiselect.Options
								}
								if len(newProps.GroupSelect.Options) > 0 {
									newOptions = newProps.GroupSelect.Options
								}
								for _, vId := range vIds {
									found := false
									for _, newOpt := range newOptions {
										if vId == newOpt.Id { // 找到新的表头选项，不需要补
											found = true
											break
										}
									}
									if !found {
										for i, oldOpt := range oldOptions {
											oldId := cast.ToString(oldOpt.Id)
											if vId == oldOpt.Id { // 找到老的表头选项
												if so, ok := selectOptions[newK]; ok {
													so[oldId] = &oldOptions[i]
												} else {
													selectOptions[newK] = map[string]*projectvo.ColumnSelectOption{
														oldId: &oldOptions[i],
													}
												}
												break
											}
										}
									}
								}
							}
						}
					}
				}

				if !oldKRemains {
					delete(newData, oldK)
				}
			}
		}
	}
}

func GetFieldMappingCache(orgId, fromTableId, toTableId int64) map[string]string {
	key := fmt.Sprintf("%d_%d_%d", orgId, fromTableId, toTableId)
	fieldMapping := make(map[string]string)
	data, err := cache.HGet(sconsts.CacheFieldMapping, key)
	if err != nil {
		return fieldMapping
	}
	_ = json.FromJson(data, &fieldMapping)
	return fieldMapping
}

func SaveFieldMappingCache(orgId, fromTableId, toTableId int64, fieldMapping map[string]string) {
	key := fmt.Sprintf("%d_%d_%d", orgId, fromTableId, toTableId)
	_ = cache.HSet(sconsts.CacheFieldMapping, key, json.ToJsonIgnoreError(fieldMapping))
}

// BuildIssueDataByCopy 通过复制一条任务数据构造一条新任务数据
// data 复制数据源
// iterationMapping 迭代值映射，传nil则保留原值(同一个项目中复制可以沿用迭代值，跨项目则需要映射)
// fieldMapping 字段映射，将其中的字段key值映射成新的，删掉不存在的key值。传nil则全量复制并且不进行映射
// (注：同一张表中复制任务可以全量复制不用映射，应用模板时由于目标表格复制了原表的表头，也可以全量复制）
func BuildIssueDataByCopy(data map[string]interface{},
	iterationMapping map[int64]int64, fieldMapping map[string][]string, issueIdMapping map[int64]int64,
	relatingColumnIds []string, selectOptions map[string]map[string]*projectvo.ColumnSelectOption,
	oldTableColumns, newTableColumns map[string]*projectvo.TableColumnData,
	newProjectTypeId int64,
	isStaticCopy, isCreateMissingSelectOptions, isCreateTemplate, isUploadTemplate bool) map[string]interface{} {
	newData := make(map[string]interface{})
	copyer.Copy(data, &newData)

	// 处理迭代（迭代没有无码化，导致这里需要映射）
	if iterationMapping != nil {
		if iterationIdI, ok := newData[consts.BasicFieldIterationId]; ok {
			if iterationId, ok := iterationMapping[cast.ToInt64(iterationIdI)]; ok {
				newData[consts.BasicFieldIterationId] = iterationId
			} else {
				delete(newData, consts.BasicFieldIterationId)
			}
		}
	}

	// 字段映射 && 表头选项补齐信息收集
	ApplyFieldMappingCreateMissingSelectOptions(data, newData, fieldMapping, selectOptions, oldTableColumns, newTableColumns,
		newProjectTypeId, isCreateMissingSelectOptions)

	// 删除确定目前不支持复制的字段，以及必须需重新设置值的元信息
	deleteColumnIds := getCopyIssueDeleteColumnIds()
	for _, id := range deleteColumnIds {
		delete(newData, id)
	}

	// 动态复制or静态复制
	if !isStaticCopy {
		for _, columnId := range relatingColumnIds {
			if relatingI, ok := newData[columnId]; ok && relatingI != nil {
				relating := &bo.RelatingIssue{}
				if err := jsonx.Copy(relatingI, relating); err != nil {
					log.Errorf("[BuildIssueDataByCopy] parse relating column failed, issueId: %v, columnId: %v, value: %v, err: %v",
						newData[consts.BasicFieldId], columnId, json.ToJsonIgnoreError(relatingI), err)
					delete(newData, columnId)
				} else {
					var linkTo, linkFrom []string
					for _, id := range relating.LinkTo {
						if newIssueId, ok := issueIdMapping[cast.ToInt64(id)]; ok {
							linkTo = append(linkTo, cast.ToString(newIssueId))
						} else if !isUploadTemplate {
							linkTo = append(linkTo, id)
						}
					}
					for _, id := range relating.LinkFrom {
						if newIssueId, ok := issueIdMapping[cast.ToInt64(id)]; ok {
							linkFrom = append(linkFrom, cast.ToString(newIssueId))
						} else if !isUploadTemplate && !isCreateTemplate {
							linkFrom = append(linkFrom, id)
						}
					}
					relating.LinkTo = linkTo
					relating.LinkFrom = linkFrom
					newData[columnId] = relating
				}
			}
		}
	}

	newData[consts.TempFieldChildren] = []map[string]interface{}{}
	newData[consts.TempFieldChildrenCount] = 0
	return newData
}

// BuildIssueChildDataByCopy 复制一条任务的子任务，来构造新任务的子任务
// parentData 需复制子任务的父任务
// childrenIds 需复制的子任务id
// allData 任务数据id->数据
// parentChildrenMapping 父任务id->子任务id mapping
// iterationMapping 迭代值映射mapping。传nil则保留原值(同一个项目中复制可以沿用迭代值，跨项目则需要映射)
// fieldMapping 字段映射，将其中的字段key值映射成新的，删掉不存在的key值。传nil则全量复制并且不进行映射
// (注：同一张表中复制任务可以全量复制不用映射，应用模板时由于目标表格复制了原表的表头，也可以全量复制）
func BuildIssueChildDataByCopy(parentData map[string]interface{}, childrenIds []int64, allData map[int64]map[string]interface{},
	parentChildrenMapping map[int64][]int64, iterationMapping map[int64]int64, fieldMapping map[string][]string, issueIdMapping map[int64]int64,
	relatingColumnIds []string, selectOptions map[string]map[string]*projectvo.ColumnSelectOption,
	oldTableColumns, newTableColumns map[string]*projectvo.TableColumnData, newProjectTypeId int64,
	isStaticCopy, isCreateMissingSelectOptions, isCreateTemplate, isUploadTemplate bool) (int64, errs.SystemErrorInfo) {
	var count int64
	children := parentData[consts.TempFieldChildren].([]map[string]interface{})

	for _, childId := range childrenIds {
		data, ok := allData[childId]
		if !ok {
			continue
		}
		newData := BuildIssueDataByCopy(data, iterationMapping, fieldMapping, issueIdMapping,
			relatingColumnIds, selectOptions, oldTableColumns, newTableColumns, newProjectTypeId,
			isStaticCopy, isCreateMissingSelectOptions, isCreateTemplate, isUploadTemplate)

		children = append(children, newData)
		count += 1

		if cIds, ok := parentChildrenMapping[childId]; ok {
			ct, err := BuildIssueChildDataByCopy(newData, cIds, allData, parentChildrenMapping,
				iterationMapping, fieldMapping, issueIdMapping, relatingColumnIds,
				selectOptions, oldTableColumns, newTableColumns, newProjectTypeId,
				isStaticCopy, isCreateMissingSelectOptions, isCreateTemplate, isUploadTemplate)
			if err != nil {
				return 0, err
			}
			count += ct
		}
	}
	parentData[consts.TempFieldChildren] = children
	parentData[consts.TempFieldChildrenCount] = count

	return count, nil
}
