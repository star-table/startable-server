package callsvc

/// 卡片交互事件回调：任务审批卡片交互

import (
	"strconv"

	pushPb "gitea.bjx.cloud/LessCode/interface/golang/push/v1"
	tableV1 "gitea.bjx.cloud/LessCode/interface/golang/table/v1"
	sdk_const "gitea.bjx.cloud/allstar/platform-sdk/consts"
	"github.com/star-table/startable-server/app/facade/orgfacade"
	"github.com/star-table/startable-server/app/facade/projectfacade"
	"github.com/star-table/startable-server/app/facade/pushfacade"
	"github.com/star-table/startable-server/common/core/consts"
	"github.com/star-table/startable-server/common/core/errs"
	"github.com/star-table/startable-server/common/core/util/json"
	"github.com/star-table/startable-server/common/extra/card"
	"github.com/star-table/startable-server/common/model/bo"
	"github.com/star-table/startable-server/common/model/vo"
	"github.com/star-table/startable-server/common/model/vo/orgvo"
	"github.com/star-table/startable-server/common/model/vo/projectvo"
)

type FsIssueAuditHandler struct{}

func (FsIssueAuditHandler) Handle(cardReq CardReq) (string, errs.SystemErrorInfo) {
	input := cardReq.Action.Value
	log.Infof("[FsIssueAuditHandler.Handle] cardReq: %s", json.ToJsonIgnoreError(cardReq))
	if input.OrgId == 0 || input.UserId == 0 {
		log.Errorf("[FsIssueAuditHandler.Handle] 卡片通知，编辑任务回调，任务id %d, 用户或组织信息未指定", input.IssueId)
		return "", errs.BuildSystemErrorInfo(errs.ReqParamsValidateError)
	}
	auditStatus := input.FsCardValueAuditIssueResult
	// 1.用户再页面更新任务为待确认->2.触发发送卡片给确认人->3.确认人点击审批通过/驳回->4.任务状态修改人收到卡片展示审批结果 & 确认人卡片状态变更。
	// 这里执行第 4 步：收到确认的审批结果后进行一些操作
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
		log.Errorf("[FsIssueAuditHandler.Handle] get issue err: %s", errMsg)
		return "", errs.BuildSystemErrorInfoWithMessage(errs.IssueNotExist, errMsg)
	}
	issue := resp.IssueInfos[0]
	projectName := consts.CardDefaultIssueProjectName
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
			log.Errorf("[FsIssueAuditHandler.Handle] err: %s", errMsg)
			return "", errs.BuildSystemErrorInfoWithMessage(errs.ProjectNotExist, errMsg)
		}
		project := (*resp2.Data)[0]
		projectName = project.Name
	}
	tableId, tableIdErr := strconv.ParseInt(issue.TableID, 10, 64)
	if tableIdErr != nil {
		log.Errorf("[FsIssueAuditHandler.Handle] ParseInt failed:%v", tableIdErr)
		return "", errs.ParamTableIdIsMust
	}
	tableInfoResp := projectfacade.GetTable(projectvo.GetTableInfoReq{
		OrgId:  issue.OrgID,
		UserId: issue.Owner,
		Input:  &tableV1.ReadTableRequest{TableId: tableId},
	})
	if tableInfoResp.Failure() {
		log.Errorf("[FsIssueAuditHandler.Handle] projectfacade.GetTable err:%v,  tableId:%d", tableInfoResp.Error(), tableId)
		return "", tableInfoResp.Error()
	}

	tableInfo := tableInfoResp.Data.Table

	// 负责人
	users, err := orgfacade.GetBaseUserInfoBatchRelaxed(issue.OrgID, issue.Owners)
	if err != nil {
		log.Errorf("[FsIssueAuditHandler.Handle] GetBaseUserInfoBatch err: %v， orgId: %d, userIds: %s", err, issue.OrgID, json.ToJsonIgnoreError(issue.Owners))
		return "", err
	}
	ownerInfos := make([]*bo.BaseUserInfoBo, 0, len(users))
	for _, u := range users {
		ownerInfo := u
		ownerInfos = append(ownerInfos, &ownerInfo)
	}

	// 审批人信息
	auditorResp := orgfacade.GetBaseUserInfoByEmpId(orgvo.GetBaseUserInfoByEmpIdReqVo{
		OrgId: input.OrgId,
		EmpId: cardReq.OpenId,
	})
	if auditorResp.Failure() {
		log.Errorf("[FsIssueAuditHandler.Handle] GetBaseUserInfoByEmpId err: %v, openId: %s", auditorResp.Error(), cardReq.OpenId)
		return "", auditorResp.Error()
	}
	auditor := auditorResp.BaseUserInfo

	// 链接
	issueLinks := projectfacade.GetIssueLinks(projectvo.GetIssueLinksReqVo{
		SourceChannel: sdk_const.SourceChannelFeishu,
		OrgId:         input.OrgId,
		IssueId:       input.IssueId,
	})
	if issueLinks.Failure() {
		log.Errorf("[FsIssueAuditHandler.Handle]GetIssueLinks err:%v", issueLinks.Error())
		return "", issueLinks.Error()
	}

	columnsResp := projectfacade.GetOneTableColumns(projectvo.GetTableColumnReq{
		OrgId:     input.OrgId,
		UserId:    input.UserId,
		ProjectId: issue.ProjectID,
		TableId:   tableId,
	})
	if columnsResp.Failure() {
		log.Errorf("[FsIssueAuditHandler.Handle] err:%v", columnsResp.Error())
		return "", columnsResp.Error()
	}

	columnMap := make(map[string]*projectvo.TableColumnData)
	for _, column := range columnsResp.Data.Columns {
		columnMap[column.Name] = column
	}

	// 审批发起人信息
	// respUser := orgfacade.GetBaseUserInfo(orgvo.GetBaseUserInfoReqVo{
	//     UserId:        input.UserId,
	//     OrgId:         input.OrgId,
	//     SourceChannel: sdk_const.SourceChannelFeishu,
	// })
	// if respUser.Failure() {
	//     log.Error(respUser.Error())
	//     return "", respUser.Error()
	// }
	// auditStarter := respUser.BaseUserInfo

	// 审批，该逻辑中包含向审批发起人发送审批结果卡片
	doAuditResp := projectfacade.AuditIssue(projectvo.AuditIssueReq{
		OrgId:  input.OrgId,
		UserId: auditor.UserId,
		Params: vo.AuditIssueReq{
			IssueID: input.IssueId,
			Status:  auditStatus,
		},
	})
	if doAuditResp.Failure() {
		log.Errorf("[FsIssueAuditHandler.Handle] err: %v, issueId: %d, auditor: %d", doAuditResp.Error(),
			input.IssueId, auditor.UserId)
		// 如果符合审批的条件，则返回不符合的卡片
		//if doAuditResp.Code == errs.CannotAuditTwice.Code() || doAuditResp.Code == errs.NotNeedAuditIssueNow.Code() {
		//	log.Infof("[FsIssueAuditHandler.Handle] NotNeedAuditIssueNow issueId： %d", input.IssueId)
		//	replyCard := card.ReplyAuditorNotNeedAudit(projectName, tableInfo.Name, ownerInfos, issue, issueLinks.NewData)
		//	cardResp := pushfacade.GenerateCard(&pushPb.GenerateCardReq{
		//		SourceChannel: SourceChannel,
		//		Card:          replyCard,
		//	})
		//	if cardResp.Failure() {
		//		log.Errorf("[FsIssueAuditHandler.Handle] err: %v, issueId: %d, auditor: %d", cardResp.Error(),
		//			input.IssueId, auditor.UserId)
		//		return "", nil
		//	}
		//
		//	return cardResp.NewData.Card, nil
		//}

		return "", doAuditResp.Error()
	}

	// 向发起人（starter）发送审批结果卡片
	// auditResultCardMsg, err := GetCardForStarter(projectName, *auditor, issue, auditStatus)
	// if err != nil {
	//     log.Errorf("[FsIssueAuditHandler.Handle] GetCardForStarter err: %v", err)
	//     return "", err
	// }
	// tenant, err := feishu.GetTenant(auditStarter.OutOrgId)
	// if err != nil {
	//     log.Errorf("[FsIssueAuditHandler.Handle] GetTenant err: %v", err)
	//     return "", err
	// }
	// fsResp, oriErr := tenant.SendMessage(fsvo.MsgVo{
	//     OpenId:  auditStarter.OutUserId,
	//     MsgType: "interactive",
	//     Card:    auditResultCardMsg,
	// })
	// if oriErr != nil {
	//     log.Errorf("[FsIssueAuditHandler.Handle] send card to audit starter SendMessage err: %v", oriErr)
	//     return "", errs.BuildSystemErrorInfo(errs.FeiShuOpenApiCallError, oriErr)
	// }
	// if fsResp.Code != 0 {
	//     log.Errorf("[FsIssueAuditHandler.Handle] send card to audit starter code err: %s",
	//        json.ToJsonIgnoreError(fsResp))
	//     return "", errs.BuildSystemErrorInfoWithMessage(errs.FeiShuOpenApiCallError, fsResp.Msg)
	// }

	if doAuditResp.Void.ID == 0 {
		// 已经审批过了，没有需要审批的记录
		replyCard := card.ReplyAuditorNotNeedAudit(projectName, tableInfo.Name, ownerInfos, issue, issueLinks.Data, columnMap)
		cardResp := pushfacade.GenerateCard(&pushPb.GenerateCardReq{
			SourceChannel: SourceChannel,
			Card:          replyCard,
		})
		if cardResp.Failure() {
			log.Errorf("[FsIssueAuditHandler.Handle] err: %v, issueId: %d, auditor: %d", cardResp.Error(),
				input.IssueId, auditor.UserId)
			return "", nil
		}

		return cardResp.Data.Card, nil
	}

	// 审批人审批完，更新他所看到的卡片状态
	replyCard := card.ReplyIssueAuditCard(projectName, tableInfo.Name, ownerInfos, issue, auditStatus, issueLinks.Data, columnMap)

	cardResp := pushfacade.GenerateCard(&pushPb.GenerateCardReq{
		SourceChannel: SourceChannel,
		Card:          replyCard,
	})
	if cardResp.Failure() {
		log.Errorf("[FsIssueAuditHandler.Handle] err: %v, issueId: %d, auditor: %d", cardResp.Error(),
			input.IssueId, auditor.UserId)
		return "", nil
	}
	return cardResp.Data.Card, nil
}

//func GetOwnerNameStr(orgId int64, userIds []int64) (string, errs.SystemErrorInfo) {
//	users, err := orgfacade.GetBaseUserInfoBatchRelaxed(orgId, userIds)
//	if err != nil {
//		log.Errorf("[GetOwnerNameStr] err: %v， orgId: %d, userIds: %s", err, orgId, json.ToJsonIgnoreError(userIds))
//		return "", err
//	}
//	userNames := make([]string, 0, len(users))
//	for _, user := range users {
//		userNames = append(userNames, user.Name)
//	}
//
//	return strings.Join(userNames, "，"), nil
//}
//
//// GetCardForStarter 向发起人发送审批结果卡片
//func GetCardForStarter(projectName string, auditor bo.BaseUserInfoBo, issue vo.Issue, auditStatus int) (*fsvo.Card, errs.SystemErrorInfo) {
//	issueLinks := projectfacade.GetIssueLinks(projectvo.GetIssueLinksReqVo{
//		SourceChannel: "",
//		OrgId:         issue.OrgID,
//		IssueId:       issue.ID,
//	}).NewData
//	tableId, tableIdErr := strconv.ParseInt(issue.TableID, 10, 64)
//	if tableIdErr != nil {
//		log.Errorf("[GetCardForStarter] ParseInt failed:%v", tableIdErr)
//		return nil, errs.ParamTableIdIsMust
//	}
//	elements := []interface{}{
//		fsvo.CardElementContentModule{
//			Tag: "div",
//			Fields: []fsvo.CardElementField{
//				{
//					Text: fsvo.CardElementText{
//						Tag:     "lark_md",
//						Content: fmt.Sprintf("**标题** %s", issue.Title),
//					},
//				},
//			},
//		},
//	}
//	tableId, oriErr := strconv.ParseInt(issue.TableID, 10, 64)
//	if oriErr != nil {
//		log.Errorf("[GetCardForStarter] err: %v, issueId: %d", oriErr, issue.ID)
//	}
//	tableInfo, err := domain.GetTableByTableId(issue.OrgID, issue.Owner, tableId)
//	if err != nil {
//		log.Errorf("[GetCardForStarter] GetTableByTableId err:%v,  tableId:%d", err, tableId)
//		return nil, err
//	}
//	ownerNameStr, err := GetOwnerNameStr(issue.OrgID, issue.Owners)
//	if err != nil {
//		log.Errorf("[GetCardForStarter] GetOwnerNameStr err: %v", err)
//		return nil, err
//	}
//
//	elements = append(elements, fsvo.CardElementContentModule{
//		Tag: "div",
//		Fields: []fsvo.CardElementField{
//			{
//				Text: fsvo.CardElementText{
//					Tag:     "lark_md",
//					Content: fmt.Sprintf("**项目/表** %s/%s", projectName, tableInfo.Name),
//				},
//			},
//		},
//	}, fsvo.CardElementContentModule{
//		Tag: "div",
//		Fields: []fsvo.CardElementField{
//			{
//				Text: fsvo.CardElementText{
//					Tag:     "lark_md",
//					Content: fmt.Sprintf("**负责人** %s", ownerNameStr),
//				},
//			},
//		},
//	}, fsvo.CardElementActionModule{
//		Tag: "action",
//		Actions: []interface{}{
//			fsvo.ActionButton{
//				Tag: "button",
//				Text: fsvo.CardElementText{
//					Tag:     "plain_text",
//					Content: consts.CardButtonTextForViewDetail,
//				},
//				Url:  issueLinks.SideBarLink,
//				Type: consts.FsCardButtonColorPrimary,
//			},
//			fsvo.ActionButton{
//				Tag: "button",
//				Text: fsvo.CardElementText{
//					Tag:     "plain_text",
//					Content: consts.CardButtonTextForViewInsideApp,
//				},
//				Url:  issueLinks.Link,
//				Type: consts.FsCardButtonColorDefault,
//			},
//		},
//	})
//	titlePart := ""
//	cardTitleColor := consts.FsCardTitleColorBlue
//	switch auditStatus {
//	case consts.AuditStatusPass:
//		titlePart = auditor.Name + "通过了审批"
//	case consts.AuditStatusReject:
//		titlePart = auditor.Name + "驳回了审批"
//		cardTitleColor = consts.FsCardTitleColorOrange
//	}
//	cardMsg := &fsvo.Card{
//		Header: &fsvo.CardHeader{
//			Title: &fsvo.CardHeaderTitle{
//				Tag:     "plain_text",
//				Content: titlePart,
//			},
//			Template: cardTitleColor,
//		},
//		Elements: elements,
//	}
//
//	return cardMsg, nil
//}
