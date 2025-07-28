package callsvc

import (
	"github.com/star-table/startable-server/common/core/errs"
)

type IssueDealTimeHandler struct{}

func (i IssueDealTimeHandler) Handle(cardReq CardReq) (string, errs.SystemErrorInfo) {
	//log.Infof("[IssueDealTimeHandler] handle info:%s", json.ToJsonIgnoreError(cardReq))
	//
	//value := cardReq.Action.Value
	//
	//// 操作人的 openId
	//opUserOpenId := cardReq.OpenId
	//issueId := value.IssueId
	//orgId := value.OrgId
	//optionsDate := cardReq.Action.Option
	//action := value.Action
	//cardTitle := value.CardTitle
	//
	//projectId := value.ProjectId
	//if value.ProjectId == 0 {
	//	projectId = int64(-1)
	//}
	//
	//appId := value.AppId
	//appIdInt64 := int64(-1)
	//if appId != "" {
	//	parseInt, err := strconv.ParseInt(appId, 10, 64)
	//	if err != nil {
	//		log.Errorf("[IssueDealTimeHandler] Handle ParseInt error:%v", err)
	//		return "", errs.TypeConvertError
	//	}
	//	appIdInt64 = parseInt
	//}
	//
	//tableId := int64(-1)
	//if value.TableId != "" {
	//	parseInt, err := strconv.ParseInt(value.TableId, 10, 64)
	//	if err != nil {
	//		log.Errorf("[IssueDealTimeHandler] Handle ParseInt error:%v", err)
	//		return "", errs.TypeConvertError
	//	}
	//	tableId = parseInt
	//}
	//
	//// 查询操作人信息
	//resp := orgfacade.GetBaseUserInfoByEmpId(orgvo.GetBaseUserInfoByEmpIdReqVo{
	//	OrgId: orgId,
	//	EmpId: opUserOpenId,
	//})
	//if resp.Failure() {
	//	return "", errs.BuildSystemErrorInfo(errs.UserNotFoundError, resp.Error())
	//}
	//opUserInfo := resp.BaseUserInfo
	//
	//// 查询出任务数据
	//issuesResp := projectfacade.GetIssueInfoList(projectvo.IssueInfoListReqVo{
	//	UserId:   opUserInfo.UserId,
	//	OrgId:    orgId,
	//	IssueIds: []int64{issueId},
	//})
	//if issuesResp.Failure() {
	//	log.Errorf("[IssueDealTimeHandler] handle GetIssueInfoList error:%v, orgId:%d, issueId:%d",
	//		issuesResp.Error(), orgId, issueId)
	//}
	//issueInfo := issuesResp.IssueInfos[0]
	//projectName := consts.CardDefaultIssueProjectName
	//tableName := consts.DefaultTableName
	//if projectId > 0 {
	//	projectResp := projectfacade.ProjectInfo(projectvo.ProjectInfoReqVo{
	//		Input:         vo.ProjectInfoReq{ProjectID: projectId},
	//		OrgId:         orgId,
	//		UserId:        opUserInfo.UserId,
	//		SourceChannel: sdk_const.SourceChannelFeishu,
	//	})
	//	if projectResp.Failure() {
	//		log.Errorf("[IssueDealTimeHandler] handle ProjectInfo error:%v, orgId:%d, issueId:%d",
	//			projectResp.Error(), orgId, issueId)
	//		return "", projectResp.Error()
	//	}
	//	projectName = projectResp.ProjectInfo.Name
	//}
	//if tableId > 0 {
	//	tableResp := projectfacade.GetTable(projectvo.GetTableInfoReq{
	//		OrgId:  orgId,
	//		UserId: opUserInfo.UserId,
	//		Input:  &tableV1.ReadTableRequest{TableId: tableId},
	//	})
	//	if tableResp.Failure() {
	//		log.Errorf("[IssueDealTimeHandler] handle GetTable error:%v, orgId:%d, issueId:%d",
	//			tableResp.Error(), orgId, issueId)
	//		return "", tableResp.Error()
	//	}
	//	tableName = tableResp.NewData.Table.Name
	//
	//}
	//
	//// 查询负责人信息
	//userInfoBatch := orgfacade.GetBaseUserInfoBatch(orgvo.GetBaseUserInfoBatchReqVo{
	//	OrgId:   orgId,
	//	UserIds: issueInfo.Owners,
	//})
	//if userInfoBatch.Failure() {
	//	log.Errorf("[IssueDealTimeHandler] handle GetBaseUserInfoBatch error:%v, orgId:%d, issueId:%d",
	//		userInfoBatch.Error(), orgId, issueId)
	//	return "", userInfoBatch.Error()
	//}
	//ownerNames := []string{}
	//for _, owner := range userInfoBatch.BaseUserInfos {
	//	ownerNames = append(ownerNames, owner.Name)
	//}
	//
	//parentTitle := consts.CardDefaultIssueParentTitleForUpdateIssue
	//if issueInfo.ParentID > 0 {
	//	// 查父记录
	//	parents := projectfacade.GetIssueInfoList(projectvo.IssueInfoListReqVo{
	//		UserId:   opUserInfo.UserId,
	//		OrgId:    orgId,
	//		IssueIds: []int64{issueInfo.ParentID},
	//	})
	//	if issuesResp.Failure() {
	//		log.Errorf("[IssueDealTimeHandler] handle GetIssueInfoList error:%v, orgId:%d, parentId:%d",
	//			parents.Error(), orgId, issueInfo.ParentID)
	//	}
	//	if len(parents.IssueInfos) > 0 {
	//		parentTitle = parents.IssueInfos[0].Title
	//	}
	//
	//}
	//
	//tableColumns, err := domain.GetTableColumnsMap(orgId, tableId, nil)
	//if err != nil {
	//	log.Errorf("[IssueDealTimeHandler]GetTableColumnsMap 获取表头失败 org:%d, proj:%d, table:%d, err: %v",
	//		orgId, projectId, tableId, err)
	//	return "", nil
	//}
	//headers := make(map[string]lc_table.LcCommonField, 0)
	//copyer.Copy(tableColumns, &headers)
	//
	//issueLinks := projectfacade.GetIssueLinks(projectvo.GetIssueLinksReqVo{
	//	SourceChannel: sdk_const.SourceChannelFeishu,
	//	OrgId:         orgId,
	//	IssueId:       issueId,
	//}).NewData
	//
	//if (action == consts.FsCardActionUpdatePlanEndTime || action == consts.FsCardActionUpdatePlanStartTime) && optionsDate != "" {
	//	targetTime, parseErr := time.Parse(consts.AppTimeFormatYYYYMMDDHHmmTimezone, optionsDate)
	//	if parseErr != nil {
	//		log.Error(parseErr)
	//		return "", errs.BuildSystemErrorInfo(errs.DateParseError)
	//	}
	//	targetTypesTime := types.Time(targetTime)
	//	form := make([]map[string]interface{}, 0)
	//	data := make(map[string]interface{})
	//	data[consts.BasicFieldId] = issueId
	//
	//	if action == consts.FsCardActionUpdatePlanStartTime {
	//		data[consts.BasicFieldPlanStartTime] = targetTypesTime
	//	}
	//	if action == consts.FsCardActionUpdatePlanEndTime {
	//		data[consts.BasicFieldPlanEndTime] = targetTypesTime
	//	}
	//
	//	form = append(form, data)
	//
	//	updateRespVo := projectfacade.BatchUpdateIssue(projectvo.BatchUpdateIssueReqVo{
	//		UserId:    opUserInfo.UserId,
	//		OrgId:     orgId,
	//		AppId:     appIdInt64,
	//		ProjectId: projectId,
	//		TableId:   tableId,
	//		NewData:      form,
	//	})
	//	if updateRespVo.Failure() {
	//		log.Errorf("[IssueDealTimeHandler] handle BatchUpdateIssue error:%v, orgId:%d, issueId:%d",
	//			updateRespVo.Error(), orgId, issueId)
	//		// 如果没有编辑权限就返回相应的卡片提示
	//		//if updateRespVo.Code == errs.NoOperationPermissionForIssueUpdate.Code() {
	//
	//		tipsCard := projectvo.FsIssueNoPermissionTipsCard{
	//			OrgId:             orgId,
	//			UserId:            opUserInfo.UserId,
	//			IssueId:           issueId,
	//			ParentId:          issueInfo.ParentID,
	//			ProjectId:         projectId,
	//			ProjectTypeId:     issueInfo.ProjectObjectTypeID,
	//			AppId:             appIdInt64,
	//			TableId:           tableId,
	//			ProjectName:       projectName,
	//			TableName:         tableName,
	//			OperateUserName:   opUserInfo.Name,
	//			IssueOwnerName:    ownerNames,
	//			Title:             issueInfo.Title,
	//			ParentTitle:       parentTitle,
	//			ColumnDisplayName: consts.Owner,
	//			CardTitle:         cardTitle,
	//			TableColumnMap:    headers,
	//			PlanStartTime:     issueInfo.PlanStartTime,
	//			PlanEndTime:       issueInfo.PlanEndTime,
	//			IssueLinks:        issueLinks,
	//		}
	//		respCard := card.GetFsCardAddMemberNoPermission(tipsCard, updateRespVo.Message)
	//		feiShuCard := card.GenerateFeiShuCard(respCard)
	//		return json.ToJsonIgnoreError(feiShuCard), nil
	//	}
	//
	//	planStartTime := issueInfo.PlanStartTime
	//	planEndTime := issueInfo.PlanEndTime
	//	if action == consts.FsCardActionUpdatePlanStartTime {
	//		planStartTime = targetTypesTime
	//	}
	//	if action == consts.FsCardActionUpdatePlanEndTime {
	//		planEndTime = targetTypesTime
	//	}
	//	issueTrendsBo := &bo.IssueTrendsBo{
	//		OrgId:   orgId,
	//		IssueId: issueId,
	//		TableId: tableId,
	//	}
	//	issueCard := &projectvo.BaseInfoBoxForIssueCard{
	//		ProjectAuthBo: bo.ProjectAuthBo{
	//			Id:    projectId,
	//			Name:  projectName,
	//			AppId: appIdInt64,
	//		},
	//		Table: projectvo.TableMetaData{
	//			Name: tableName,
	//		},
	//		OperateUser: *opUserInfo,
	//		IssueInfo: bo.IssueBo{
	//			Title:         issueInfo.Title,
	//			ParentId:      issueInfo.ParentID,
	//			PlanStartTime: planStartTime,
	//			PlanEndTime:   planEndTime,
	//		},
	//		ParentIssue:       bo.IssueBo{Title: parentTitle},
	//		IssueOwnerNameArr: ownerNames,
	//		OperateUserId:     opUserInfo.UserId,
	//		IssueInfoUrl:      issueLinks.SideBarLink,
	//		IssuePcUrl:        issueLinks.Link,
	//		TableColumnMap:    headers,
	//	}
	//	columnDisplayName := headers[consts.BasicFieldOwnerId].AliasTitle
	//	if headers[consts.BasicFieldOwnerId].AliasTitle == "" {
	//		columnDisplayName = headers[consts.BasicFieldOwnerId].Label
	//	}
	//	respCard := card.GetCardIssueAddMember(issueTrendsBo, issueCard, columnDisplayName)
	//	feiShuCard := card.GenerateFeiShuCard(respCard)
	//	return json.ToJsonIgnoreError(feiShuCard), nil
	//}

	return "", nil
}
