package callsvc

/// 卡片交互事件回调：任务催办卡片交互
// https://open.feishu.cn/document/ukTMukTMukTM/uYjNwUjL2YDM14iN2ATN

import (
	"strconv"

	pushV1 "gitea.bjx.cloud/LessCode/interface/golang/push/v1"
	tableV1 "gitea.bjx.cloud/LessCode/interface/golang/table/v1"
	"github.com/star-table/startable-server/common/core/util/json"
	"github.com/star-table/startable-server/common/core/util/slice"
	sdk_const "gitea.bjx.cloud/allstar/platform-sdk/consts"
	"github.com/star-table/startable-server/common/core/consts"
	"github.com/star-table/startable-server/common/core/errs"
	"github.com/star-table/startable-server/common/extra/card"
	"github.com/star-table/startable-server/common/model/bo"
	"github.com/star-table/startable-server/common/model/vo/orgvo"
	"github.com/star-table/startable-server/common/model/vo/projectvo"
	"github.com/star-table/startable-server/app/facade/orgfacade"
	"github.com/star-table/startable-server/app/facade/projectfacade"
	"github.com/star-table/startable-server/app/facade/pushfacade"
)
)

type FsIssueUrgeHandler struct{}

func (FsIssueUrgeHandler) Handle(cardReq CardReq) (string, errs.SystemErrorInfo) {
	input := cardReq.Action.Value
	optionVal := cardReq.Action.Option
	log.Infof("FsIssueUrgeHandler cardReq: %s", json.ToJsonIgnoreError(cardReq))
	if input.OrgId == 0 || input.UserId == 0 {
		log.Errorf("FsIssueUrgeHandler 卡片通知，编辑任务回调，任务id %d, 用户或组织信息未指定", input.IssueId)
		return "", errs.BuildSystemErrorInfo(errs.ReqParamsValidateError)
	}
	// 1.用户再页面点击催办->2.发送卡片给被催办人->3.被催办人点击回应：处理/暂不处理->4.催办人收到卡片展示被催者是否处理 & 被催人卡片状态变更。
	// 这里执行第 4 步。
	// 4.1 被催者点击卡片交互后，给催办人回复卡片
	resp := projectfacade.GetIssueInfoList(projectvo.IssueInfoListReqVo{
		OrgId:    input.OrgId,
		UserId:   input.UserId,
		IssueIds: []int64{input.IssueId},
	})
	if resp.Failure() {
		log.Error(resp.Error())
		return "", resp.Error()
	}
	if len(resp.IssueInfos) < 1 {
		errMsg := "| 任务未找到[1]"
		log.Errorf("%s", errMsg)
		return "", errs.BuildSystemErrorInfoWithMessage(errs.IssueNotExist, errMsg)
	}
	issue := resp.IssueInfos[0]
	projectName := consts.CardDefaultIssueProjectName
	appId := int64(0)
	projectTypeId := int64(0)
	if issue.ProjectID != 0 {
		resp2 := projectfacade.GetSimpleProjectInfo(projectvo.GetSimpleProjectInfoReqVo{
			OrgId: input.OrgId,
			Ids:   []int64{issue.ProjectID},
		})
		if resp2.Failure() {
			log.Error(resp2.Error())
			return "", resp2.Error()
		}
		if resp2.Data == nil || len(*resp2.Data) < 1 {
			errMsg := "| 任务未找到[2]"
			log.Error(errMsg)
			return "", errs.BuildSystemErrorInfoWithMessage(errs.ProjectNotExist, errMsg)
		}
		project := (*resp2.Data)[0]
		projectName = project.Name
		projectTypeId = project.ProjectTypeID
		appIdConv, err := strconv.ParseInt(project.AppID, 10, 64)
		if err != nil {
			log.Error(err)
			return "", errs.TypeConvertError
		}
		appId = appIdConv
	}

	parentTitle := ""
	if issue.ParentID != 0 {
		parentIssueResp := projectfacade.GetIssueInfoList(projectvo.IssueInfoListReqVo{
			OrgId:    input.OrgId,
			UserId:   input.UserId,
			IssueIds: []int64{issue.ParentID},
		})
		if parentIssueResp.Failure() {
			log.Error(parentIssueResp.Error())
			return "", parentIssueResp.Error()
		}
		parentIssue := parentIssueResp.IssueInfos[0]
		parentTitle = parentIssue.Title
	}

	tableId, errSys := strconv.ParseInt(issue.TableID, 10, 64)
	if errSys != nil {
		log.Errorf("FsIssueUrgeHandler tableId err:%v, tableIdStr:%s", errSys, issue.TableID)
		return "", errs.ParamTableIdIsMust
	}

	tableName := consts.DefaultTableName
	if tableId != 0 {
		tableInfoResp := projectfacade.GetTable(projectvo.GetTableInfoReq{
			OrgId:  issue.OrgID,
			UserId: issue.Owner,
			Input:  &tableV1.ReadTableRequest{TableId: tableId},
		})
		if tableInfoResp.Failure() {
			log.Errorf("[FsIssueUrgeHandler.Handle] projectfacade.GetTable err:%v,  tableId:%d", tableInfoResp.Error(), tableId)
			return "", tableInfoResp.Error()
		}
		tableName = tableInfoResp.Data.Table.Name
	}

	// 查询负责人和催办人
	userRelationIds := []int64{}
	usersMap := map[int64]bo.BaseUserInfoBo{}
	userRelationIds = append(userRelationIds, issue.Owners...)
	userRelationIds = append(userRelationIds, input.UserId)
	userRelationIds = slice.SliceUniqueInt64(userRelationIds)

	userInfoBatch := orgfacade.GetBaseUserInfoBatch(orgvo.GetBaseUserInfoBatchReqVo{
		OrgId:   input.OrgId,
		UserIds: userRelationIds,
	})

	if userInfoBatch.Failure() {
		log.Errorf("[FsIssueUrgeHandler]orgfacade.GetBaseUserInfoBatch err:%v, orgId:%d, ",
			userInfoBatch.Error(), input.OrgId)
	}

	for _, u := range userInfoBatch.BaseUserInfos {
		usersMap[u.UserId] = u
	}

	// 负责人信息
	ownerInfos := []*bo.BaseUserInfoBo{}
	for _, owner := range issue.Owners {
		ownerInfo := usersMap[owner]
		ownerInfos = append(ownerInfos, &ownerInfo)
	}

	// 催办人信息
	urgeUserName := usersMap[input.UserId].Name
	urgeOutOrgId := usersMap[input.UserId].OutOrgId
	urgeOutUserId := usersMap[input.UserId].OutUserId

	// 被催办者的信息
	beUrgeResp := orgfacade.GetBaseUserInfoByEmpId(orgvo.GetBaseUserInfoByEmpIdReqVo{
		OrgId: issue.OrgID,
		EmpId: cardReq.OpenId,
	})
	if beUrgeResp.Failure() {
		log.Error(beUrgeResp.Error())
		return "", beUrgeResp.Error()
	}
	beUrgedUser := beUrgeResp.BaseUserInfo

	issueLinks := projectfacade.GetIssueLinks(projectvo.GetIssueLinksReqVo{
		SourceChannel: sdk_const.SourceChannelFeishu,
		OrgId:         input.OrgId,
		IssueId:       input.IssueId,
	})
	if issueLinks.Failure() {
		log.Errorf("[FsIssueUrgeHandler]GetIssueLinks err:%v, issueId:%d", issueLinks.Error(), input.IssueId)
		return "", issueLinks.Error()
	}
	columnsResp := projectfacade.GetOneTableColumns(projectvo.GetTableColumnReq{
		OrgId:     input.OrgId,
		UserId:    input.UserId,
		ProjectId: issue.ProjectID,
		TableId:   tableId,
	})
	if columnsResp.Failure() {
		log.Errorf("[PushCreateIssueSuccessNotice] err:%v", columnsResp.Error())
		return "", columnsResp.Error()
	}

	columnMap := make(map[string]*projectvo.TableColumnData)
	for _, column := range columnsResp.Data.Columns {
		columnMap[column.Name] = column
	}

	cardSt := projectvo.UrgeReplyCard{
		OrgId:           issue.OrgID,
		OutOrgId:        urgeOutOrgId,
		OpenIds:         []string{urgeOutUserId},
		ProjectName:     projectName,
		TableName:       tableName,
		OperateUserName: urgeUserName,
		BeUrgedUserName: beUrgedUser.Name,
		UrgeText:        input.FsCardValueUrgeIssueUrgeText,
		ReplyMsg:        optionVal,
		ParentTitle:     parentTitle,
		AppId:           appId,
		ProjectTypeId:   projectTypeId,
		Issue:           issue,
		OwnerInfos:      ownerInfos,
		IssueLinks:      issueLinks.Data,
		TableColumn:     columnMap,
	}

	errCard := card.SendCardBeUrgeIssue(sdk_const.SourceChannelFeishu, cardSt)
	if errCard != nil {
		log.Errorf("FsIssueUrgeHandler-Handle 发送消息异常 err:%v", errSys)
		return "", errCard
	}

	replyCardMeta := card.GetFsCardAfterReplyUrgeIssue(cardSt)
	feiShuCardReply := pushfacade.GenerateCard(&pushV1.GenerateCardReq{
		SourceChannel: sdk_const.SourceChannelFeishu,
		Card:          replyCardMeta,
	})
	if feiShuCardReply.Failure() {
		log.Errorf("FsIssueUrgeHandler-Handle pushfacade.GenerateCard err:%v", feiShuCardReply.Error())
		return "", feiShuCardReply.Error()
	}
	return feiShuCardReply.Data.Card, nil
}
