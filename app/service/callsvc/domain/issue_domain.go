package callsvc

import (
	tableV1 "gitea.bjx.cloud/LessCode/interface/golang/table/v1"
	"github.com/star-table/startable-server/app/facade/formfacade"
	"github.com/star-table/startable-server/app/facade/projectfacade"
	"github.com/star-table/startable-server/common/core/businees"
	"github.com/star-table/startable-server/common/core/consts"
	"github.com/star-table/startable-server/common/core/errs"
	"github.com/star-table/startable-server/common/core/util/copyer"
	"github.com/star-table/startable-server/common/core/util/json"
	"github.com/star-table/startable-server/common/extra/lc_helper"
	"github.com/star-table/startable-server/common/model/bo"
	"github.com/star-table/startable-server/common/model/vo/formvo"
	"github.com/star-table/startable-server/common/model/vo/projectvo"
	"github.com/spf13/cast"
)

// GetIssueByChatId 通过任务群聊的 chatId 查询对应的任务
func GetIssueByChatId(orgId int64, outChatId string, needDel bool) ([]bo.IssueBo, errs.SystemErrorInfo) {
	resList := make([]bo.IssueBo, 0)
	condition := &tableV1.Condition{
		Type: tableV1.ConditionType_and,
		Conditions: []*tableV1.Condition{
			&tableV1.Condition{
				Type:   tableV1.ConditionType_equal,
				Value:  json.ToJsonIgnoreError([]interface{}{orgId}),
				Column: lc_helper.ConvertToCondColumn(consts.BasicFieldOrgId),
			},
		},
	}
	if !needDel {
		condition.Conditions = append(condition.Conditions, &tableV1.Condition{Column: lc_helper.ConvertToCondColumn(consts.BasicFieldRecycleFlag),
			Type: tableV1.ConditionType_equal, Value: json.ToJsonIgnoreError([]interface{}{consts.AppIsNoDelete})})
	}
	condition.Conditions = append(condition.Conditions, &tableV1.Condition{Column: lc_helper.ConvertToCondColumn(consts.BasicFieldOutChatId),
		Type: tableV1.ConditionType_equal, Value: json.ToJsonIgnoreError([]interface{}{outChatId})})

	lessResp := projectfacade.GetIssueRowList(projectvo.IssueRowListReq{
		OrgId:  orgId,
		UserId: 0,
		Input: &tableV1.ListRawRequest{
			DbType:    tableV1.DbType_slave1,
			Condition: condition,
		},
	})
	if lessResp.Failure() {
		log.Errorf("[GetIssueByChatId] GetIssueRowList err:%v, orgId:%v, outChatId:%v", lessResp.Error(), orgId, outChatId)
		return nil, lessResp.Error()
	}

	for _, issueInfo := range lessResp.Data {
		issueBo, errSys := ConvertIssueDataToIssueBo(issueInfo)
		if errSys != nil {
			log.Errorf("[GetIssueByChatId] ConvertIssueDataToIssueBo err: %v", errSys)
			return nil, errSys
		}
		resList = append(resList, *issueBo)
	}
	issueIds := make([]int64, 0, len(resList))
	for _, issue := range resList {
		issueIds = append(issueIds, issue.Id)
	}
	issueDelFlagMap := make(map[int64]int, len(issueIds))
	if len(issueIds) > 0 {
		// 通过查询取出删除状态
		proResp := projectfacade.GetLcIssueInfoBatch(projectvo.GetLcIssueInfoBatchReqVo{
			OrgId:    orgId,
			IssueIds: issueIds,
		})
		if proResp.Failure() {
			log.Errorf("[GetIssueByChatId] GetSimpleIssueInfoBatch err: %v", proResp.Error())
			return nil, proResp.Error()
		}
		for _, issue := range proResp.Data {
			issueDelFlagMap[issue.Id] = issue.IsDelete
		}
	}
	for i, item := range resList {
		resList[i].IsDelete = issueDelFlagMap[item.Id]
	}

	return resList, nil
}

func ConvertIssueDataToIssueBo(data map[string]interface{}) (*bo.IssueBo, errs.SystemErrorInfo) {
	issueBo := &bo.IssueBo{}
	err := copyer.Copy(data, issueBo)
	if err != nil {
		log.Errorf("[ConvertIssueDataToIssueBo] json转换错误, err:%v", err)
		return nil, errs.JSONConvertError
	}
	issueBo.Id = cast.ToInt64(data[consts.BasicFieldIssueId])
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
	issueBo.LessData = data

	return issueBo, nil
}

func SaveIssueForClearChatId(orgId int64, issue *bo.IssueBo) errs.SystemErrorInfo {
	updForm := make([]map[string]interface{}, 0, 1)
	updForm = append(updForm, map[string]interface{}{
		consts.BasicFieldIssueId:   issue.Id,
		consts.BasicFieldOutChatId: "",
	})
	appId, err := GetOrgSummaryAppId(orgId)
	if err != nil {
		log.Errorf("[SaveIssueForClearChatId] GetOrgSummaryAppId err: %v, issueId: %d", err, issue.Id)
		return err
	}
	lcReq := formvo.LessUpdateIssueReq{
		OrgId:   orgId,
		AppId:   appId,
		TableId: issue.TableId,
		UserId:  0,
		Form:    updForm,
	}
	// 更新任务信息到无码
	resp := formfacade.LessUpdateIssue(lcReq)
	if resp.Failure() {
		log.Errorf("[SaveIssueForClearChatId] LessUpdateIssue: %v, issueId: %d", resp.Error(), issue.Id)
		return resp.Error()
	}

	return nil
}
