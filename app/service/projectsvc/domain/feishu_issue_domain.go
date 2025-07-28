package domain

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	automationPb "gitea.bjx.cloud/LessCode/interface/golang/automation/v1"
	pushPb "gitea.bjx.cloud/LessCode/interface/golang/push/v1"
	pushV1 "gitea.bjx.cloud/LessCode/interface/golang/push/v1"
	sdk_const "gitea.bjx.cloud/allstar/platform-sdk/consts"
	"github.com/star-table/startable-server/app/facade/orgfacade"
	"github.com/star-table/startable-server/app/facade/projectfacade"
	consts2 "github.com/star-table/startable-server/app/service/projectsvc/consts"
	"github.com/star-table/startable-server/common/core/businees"
	"github.com/star-table/startable-server/common/core/consts"
	"github.com/star-table/startable-server/common/core/errs"
	"github.com/star-table/startable-server/common/core/util"
	"github.com/star-table/startable-server/common/core/util/copyer"
	"github.com/star-table/startable-server/common/core/util/json"
	"github.com/star-table/startable-server/common/core/util/slice"
	int642 "github.com/star-table/startable-server/common/core/util/slice/int64"
	"github.com/star-table/startable-server/common/core/util/str"
	"github.com/star-table/startable-server/common/extra/card"
	"github.com/star-table/startable-server/common/extra/lc_helper"
	"github.com/star-table/startable-server/common/model/bo"
	vo1 "github.com/star-table/startable-server/common/model/vo"
	"github.com/star-table/startable-server/common/model/vo/commonvo"
	"github.com/star-table/startable-server/common/model/vo/lc_table"
	"github.com/star-table/startable-server/common/model/vo/orgvo"
	"github.com/star-table/startable-server/common/model/vo/projectvo"
	"github.com/spf13/cast"
)

// PushInfoToChat 任务更新后，发送群聊卡片
func PushInfoToChat(orgId int64, projectId int64, info *bo.IssueTrendsBo, sourceChanel string) errs.SystemErrorInfo {
	// 群聊卡片目前只支持飞书
	if sourceChanel != sdk_const.SourceChannelFeishu {
		return nil
	}
	log.Infof("[PushInfoToChat] 推送群聊卡片开始--------------%v", json.ToJsonIgnoreError(info))

	//判断是否要发送通知
	settingBoArr, err := GetProTableChatSettingBoArrByTableId(orgId, projectId, info.TableId)
	if len(settingBoArr) < 1 {
		//没有群聊关联
		log.Infof("[PushInfoToChat] 没有对应的表群聊设置，无需推送, proId: %d", projectId)
		return nil
	}
	// 过滤需要推送的群聊配置。如果配置的结果是无需推送，则过滤掉
	settingBoArr = FilterProTableChatSetting(settingBoArr, info.PushType)

	baseOrgInfo, err := orgfacade.GetBaseOrgInfoRelaxed(orgId)
	if err != nil {
		log.Error(err)
		return err
	}
	if baseOrgInfo.OutOrgId == "" {
		return nil
	}

	msg, err := pushChatMsg(sourceChanel, orgId, info)
	if err != nil {
		if err == errs.CardTitleEmpty || err == errs.CardColumnEmpty || err == errs.CardNoNeedPush {
			return nil
		}
		log.Errorf("[PushInfoToChat] orgId:%v, info:%s pushChatMsg error:%v", orgId, json.ToJsonIgnoreError(info), err)
		return err
	}
	if msg == nil {
		return nil
	}
	chatIds := []string{}
	for _, ch := range settingBoArr {
		chatIds = append(chatIds, ch.ChatId)
	}
	pushCard := &commonvo.PushCard{
		OrgId:         orgId,
		OutOrgId:      baseOrgInfo.OutOrgId,
		SourceChannel: sourceChanel,
		ChatIds:       chatIds,
		CardMsg:       msg,
	}

	errSys := card.PushCard(orgId, pushCard)
	if errSys != nil {
		log.Errorf("[PushInfoToChat] err:%v", errSys)
		return errSys
	}
	//tenant, tenantErr := feishu.GetTenant(baseOrgInfo.OutOrgId)
	//if tenantErr != nil {
	//	log.Error(tenantErr)
	//	return nil
	//}
	//fsCard := card.GenerateFeiShuCard(msg)
	//for _, tmpSetting := range settingBoArr {
	//	resp, respErr := tenant.SendMessage(vo.MsgVo{
	//		ChatId:  tmpSetting.ChatId,
	//		MsgType: "interactive",
	//		Card:    fsCard,
	//	})
	//	if respErr != nil {
	//		log.Errorf("[PushInfoToChat] err: %v", respErr)
	//		return nil
	//	}
	//
	//	if resp.Code != 0 {
	//		log.Infof("[PushInfoToChat] toId:%v, msg:%v, 发送飞书卡片通知失败:%v", tmpSetting.ChatId, msg, resp.Msg)
	//		return errs.FeiShuOpenApiCallError
	//	}
	//}

	return nil
}

func FilterProTableChatSetting(settingBoArr []bo.GetFsTableSettingOfProjectChatBo,
	pushType consts.IssueNoticePushType,
) []bo.GetFsTableSettingOfProjectChatBo {
	resultArr := make([]bo.GetFsTableSettingOfProjectChatBo, 0)
	for _, tmpSetting := range settingBoArr {
		tmpNeedSendFlag := false
		switch pushType {
		case consts.PushTypeCreateIssue:
			if tmpSetting.CreateIssue != 1 {
				continue
			} else {
				tmpNeedSendFlag = true
			}
		case consts.PushTypeIssueComment:
			if tmpSetting.CreateIssueComment != 1 {
				continue
			} else {
				tmpNeedSendFlag = true
			}
		case consts.PushTypeUpdateIssue:
			// 只要是更新任务，都会推送，因此只要开启了更新任务推送，则推送卡片
			if tmpSetting.UpdateIssueCase != 1 {
				continue
			} else {
				tmpNeedSendFlag = true
			}
		case consts.PushTypeUploadResource, consts.PushTypeReferResource:
			// 暂时忽略，上传附件会再次调用 update，从而产生 consts.PushTypeUpdateIssue 调用
		}
		if tmpNeedSendFlag {
			resultArr = append(resultArr, tmpSetting)
		}
	}

	return resultArr
}

// FilterProTableChatSettingForManyChat 为多个群聊配置过滤
func FilterProTableChatSettingForManyChat(settingBoArr []bo.GetFsTableSettingOfProjectChatBo,
	pushType consts.IssueNoticePushType,
) []bo.GetFsTableSettingOfProjectChatBo {
	resultArr := make([]bo.GetFsTableSettingOfProjectChatBo, 0)

	for _, tmpSetting := range settingBoArr {
		tmpNeedSendFlag := false
		switch pushType {
		case consts.PushTypeCreateIssue:
			if tmpSetting.CreateIssue != 1 {
				continue
			} else {
				tmpNeedSendFlag = true
			}
		case consts.PushTypeIssueComment:
			if tmpSetting.CreateIssueComment != 1 {
				continue
			} else {
				tmpNeedSendFlag = true
			}
		case consts.PushTypeUpdateIssue:
			// 只要是更新任务，都会推送，因此只要开启了更新任务推送，则推送卡片
			if tmpSetting.UpdateIssueCase != 1 {
				continue
			} else {
				tmpNeedSendFlag = true
			}
		case consts.PushTypeUploadResource, consts.PushTypeReferResource:
			// 暂时忽略，上传附件会再次调用 update，从而产生 consts.PushTypeUpdateIssue 调用
		}
		if tmpNeedSendFlag {
			resultArr = append(resultArr, tmpSetting)
		}
	}

	return resultArr
}

func assemblyIssueUpdateAndWorkHour(pushType consts.IssueNoticePushType, orgId, tableId int64, info *bo.IssueTrendsBo, card *pushPb.TemplateCard) errs.SystemErrorInfo {
	// 获取表头
	tableColumns, errSys := GetTableColumnsMap(orgId, tableId, nil)
	if errSys != nil {
		log.Errorf("[assemblyIssueUpdateAndWorkHour] 获取表头失败 orgId:%d, tableId:%d, err: %v", orgId, tableId, errSys)
		return errSys
	}
	headers := make(map[string]lc_table.LcCommonField, 0)
	errCopy := copyer.Copy(tableColumns, &headers)
	if errCopy != nil {
		return errs.BuildSystemErrorInfo(errs.ObjectCopyError, errCopy)
	}
	tableColumnMetas := make(map[string]lc_table.LcCommonField, 0)
	for columnId, tableColumn := range headers {
		if tableColumn.Field.Props.PushMsg {
			tableColumnMetas[columnId] = tableColumn
		}
	}
	if len(tableColumnMetas) < 1 {
		//log.Errorf("[assemblyIssueUpdateAndWorkHour]没有需要推送的字段, tableId:%v", tableId)
		return errs.CardColumnEmpty
	}

	switch pushType {
	case consts.PushTypeUpdateIssue:
		AssemblyFsCardForUpdateIssue(tableColumnMetas, info, card)
	case consts.PushTypeCreateWorkHour, consts.PushTypeUpdateWorkHour, consts.PushTypeDeleteWorkHour:
		AssemblyFsCardForUpdateIssueWorkHour(tableColumnMetas, info, card)
	}

	if card.Divs == nil || len(card.Divs) < 1 {
		return errs.CardNoNeedPush
	}

	return nil
}

// 任务新增/更新后，群聊卡片的组装
func pushChatMsg(sourceChannel string, orgId int64, info *bo.IssueTrendsBo) (*pushPb.TemplateCard, errs.SystemErrorInfo) {
	issueBos, errSys := GetIssueInfosLc(orgId, 0, []int64{info.IssueId})
	if errSys != nil {
		log.Errorf("[pushChatMsg]GetIssueInfosLc err:%v, issueId:%v", errSys, info.IssueId)
		return nil, errSys
	}
	if len(issueBos) < 1 {
		log.Errorf("[pushChatMsg] not fount issue issueId:%v", info.IssueId)
		return nil, errs.IssueNotExist
	}
	issueInfo := issueBos[0]
	//任务标题为空，则不需要推送
	if issueInfo.Title == "" {
		return nil, errs.CardTitleEmpty
	}

	cardMsg := &pushPb.TemplateCard{}

	// 项目 表
	projectName := consts.CardDefaultIssueProjectName
	projectTypeId := int64(0)
	if issueInfo.ProjectId > 0 {
		projectInfo, errSys := GetProjectSimple(orgId, issueInfo.ProjectId)
		if errSys != nil {
			log.Errorf("[pushChatMsg]GetProjectSimple err:%v, projectId:%v", errSys, issueInfo.ProjectId)
			return nil, errSys
		}
		projectName = projectInfo.Name
		projectTypeId = projectInfo.ProjectTypeId
	}

	tableName := ""
	if issueInfo.TableId != 0 {
		tableInfo, errSys := GetTableByTableId(orgId, info.OperatorId, issueInfo.TableId)
		if errSys != nil {
			log.Errorf("[pushChatMsg]GetTableInfo failed:%v, userId:%d, tableId:%d", errSys, info.OperatorId, issueInfo.TableId)
			return nil, errSys
		}
		//if projectName != consts.CardDefaultIssueProjectName {
		//	contentProject = fmt.Sprintf(consts.CardTablePro, projectName, tableInfo.Name)
		//} else {
		//	contentProject = consts.CardDefaultIssueProjectName
		//}
		tableName = tableInfo.Name
	}

	userIds := []int64{}
	userIds = append(userIds, info.OperatorId)
	userIds = append(userIds, issueInfo.OwnerIdI64...)
	userIds = slice.SliceUniqueInt64(userIds)

	userInfoBatch := orgfacade.GetBaseUserInfoBatch(orgvo.GetBaseUserInfoBatchReqVo{
		OrgId:   orgId,
		UserIds: userIds,
	})
	if userInfoBatch.Failure() {
		log.Errorf("[pushChatMsg]GetBaseUserInfoBatch err:%v, userIds:%v", userInfoBatch.Error(), userIds)
		return nil, userInfoBatch.Error()
	}
	userMap := map[int64]bo.BaseUserInfoBo{}
	for _, u := range userInfoBatch.BaseUserInfos {
		userMap[u.UserId] = u
	}

	// 操作人
	operatorName := userMap[info.OperatorId].Name

	// 负责人
	ownerInfos := make([]*bo.BaseUserInfoBo, 0)
	for _, ownerId := range issueInfo.OwnerIdI64 {
		if v, ok := userMap[ownerId]; ok {
			ownerInfos = append(ownerInfos, &v)
		}
	}

	// 父记录
	parentTitle := ""
	if issueInfo.ParentId > 0 {
		data, errSys := GetIssueInfosMapLcByIssueIds(orgId, info.OperatorId, []int64{issueInfo.ParentId}, lc_helper.ConvertToFilterColumn(consts.BasicFieldTitle))
		if errSys != nil {
			return nil, errSys
		}
		if len(data) < 1 {
			return nil, errs.IssueNotExist
		}
		parentTitle = cast.ToString(data[0][consts.BasicFieldTitle])

	}

	taskTypeName := "记录"
	title := card.GetOperationName(info, operatorName) + " "

	switch info.PushType {
	case consts.PushTypeCreateIssue:
		title += "创建了新" + taskTypeName
	case consts.PushTypeUpdateIssue, consts.PushTypeCreateWorkHour, consts.PushTypeUpdateWorkHour, consts.PushTypeDeleteWorkHour:
		title += "更新了" + taskTypeName
		errSys := assemblyIssueUpdateAndWorkHour(info.PushType, orgId, issueInfo.TableId, info, cardMsg)
		if errSys != nil {
			//log.Errorf("[pushChatMsg] err:%v", errSys)
			return nil, errSys
		}
	case consts.PushTypeIssueComment:
		title += "评论了" + taskTypeName
		speakerNameDesc := fmt.Sprintf("%s:", operatorName)
		comment := GetFeiShuAtContent(orgId, info.Ext.CommentBo.Content)
		// 不截取@成员的字数
		//words := consts.GetCardCommentRealWords(comment)
		//comment = consts.TruncateText(words, consts.GrouChatIssueChangeDescLimitPerLine)
		if info.Ext.ResourceInfo != nil && len(info.Ext.ResourceInfo) > 0 {
			comment += " [附件]"
		}
		//commentCutter := str.NewCommentCutter(comment)
		//comment = commentCutter.CutComment(consts.GroupChatIssueCommentLimit)
		comment = CutComment(comment, consts.GroupChatIssueCommentLimit)
		cardMsg.Divs = append(cardMsg.Divs, &pushPb.CardTextDiv{
			Fields: []*pushPb.CardTextField{
				&pushPb.CardTextField{
					Key:   speakerNameDesc,
					Value: fmt.Sprintf(consts.CardDoubleQuotationMarks, comment),
				},
			},
		})
	default:
		return nil, nil
	}

	if title == operatorName+" " {
		errSys = errs.BuildSystemErrorInfoWithMessage(errs.NotSupportTypeForPushFeishuChat, fmt.Sprintf("pushType: %d", info.PushType))
		//log.Errorf("不支持的推送类型，param: %s", json.ToJsonIgnoreError(info))
		return nil, errSys
	}

	issueLinks := GetIssueLinks("", orgId, info.IssueId)

	tableColumns, errSys := GetTableColumnsMap(orgId, issueInfo.TableId, []string{consts.BasicFieldTitle, consts.BasicFieldOwnerId})
	if errSys != nil {
		log.Errorf("[pushChatMsg] GetTableColumnConfig err:%v, orgId:%v, issueId:%v", errSys, orgId, info.IssueId)
		return nil, errSys
	}

	cardCommonMsg := card.GenerateCard(&commonvo.GenerateCard{
		OrgId:             orgId,
		CardLevel:         pushPb.CardElementLevel_Default,
		CardTitle:         title,
		IssueTitle:        issueInfo.Title,
		ParentId:          issueInfo.ParentId,
		ParentTitle:       parentTitle,
		ProjectId:         issueInfo.ProjectId,
		ProjectTypeId:     projectTypeId,
		ProjectName:       projectName,
		TableName:         tableName,
		OwnerInfos:        ownerInfos,
		Links:             issueLinks,
		IsNeedUnsubscribe: false,
		SpecialDiv:        cardMsg.Divs,
		TableColumnInfo:   tableColumns,
	})

	log.Infof("[pushChatMsg] cardInfo:%s", json.ToJsonIgnoreError(cardCommonMsg))

	return cardCommonMsg, nil
}

// 任务新增/更新后，群聊卡片的组装
//func dealPushMsg(sourceChannel string, orgId int64, info bo.IssueTrendsBo) (*commonvo.CardMeta, errs.SystemErrorInfo) {
//	issueLinks := GetIssueLinks(sourceChannel, orgId, info.IssueId)
//	operatorBaseInfo, err := orgfacade.GetBaseUserInfoRelaxed(orgId, info.OperatorId)
//	if err != nil {
//		log.Error(err)
//		return nil, err
//	}
//	issueBos, err2 := GetIssueInfosLc(orgId, 0, []int64{info.IssueId})
//	if err2 != nil {
//		log.Error(err2)
//		return nil, errs.BuildSystemErrorInfo(errs.IssueDomainError, err2)
//	}
//	if len(issueBos) < 1 {
//		log.Errorf("[dealPushMsg] not found issue issueId:%v", info.IssueId)
//		return nil, errs.IssueNotExist
//	}
//	issueInfo := issueBos[0]
//	//任务标题为空，则不需要推送
//	if issueInfo.Title == "" {
//		return nil, errs.CardTitleEmpty
//	}
//	projectName := consts.CardDefaultIssueProjectName
//	if issueInfo.ProjectId > 0 {
//		project, err := GetProjectSimple(orgId, issueInfo.ProjectId)
//		if err != nil {
//			log.Errorf("[dealPushMsg] err: %v, projectId:%v", err, issueInfo.ProjectId)
//			return nil, err
//		}
//		projectName = project.Name
//	}
//	taskTypeName := "记录"
//	// "⏰ " +
//	title := operatorBaseInfo.Name + " "
//	//markdown := fmt.Sprintf("%s%s", consts.CardElementTitle+consts.FsCard4Tab, info.IssueTitle)
//
//	cardMeta := &commonvo.CardMeta{
//		IsWide: false,
//		Level:  consts.CardLevelInfo,
//	}
//
//	cardDivs := []*commonvo.CardDiv{}
//	cardDivs = append(cardDivs, &commonvo.CardDiv{
//		Fields: []*commonvo.CardField{
//			&commonvo.CardField{
//				Key:   consts.CardElementTitle,
//				Value: info.IssueTitle,
//			},
//		},
//	})
//
//	if issueInfo.TableId != 0 {
//		tableInfo, err := GetTableByTableId(orgId, info.OperatorId, issueInfo.TableId)
//		if err != nil {
//			log.Errorf("[dealPushMsg]GetTableInfo failed:%v, userId:%d, tableId:%d", err, info.OperatorId, issueInfo.TableId)
//			return nil, err
//		}
//		contentProject := ""
//		if projectName != consts.CardDefaultIssueProjectName {
//			contentProject = fmt.Sprintf(consts.CardTablePro, projectName, tableInfo.Name)
//		} else {
//			contentProject = consts.CardDefaultIssueProjectName
//		}
//		cardDivs = append(cardDivs, &commonvo.CardDiv{
//			Fields: []*commonvo.CardField{
//				&commonvo.CardField{
//					Key:   consts.CardElementProjectTable,
//					Value: contentProject,
//				},
//			},
//		})
//	}
//	if issueInfo.ParentId != 0 {
//		parentTitle := consts.CardDefaultRelationIssueTitle
//		parentInfo, err := GetIssueBo(orgId, issueInfo.ParentId)
//		if err != nil {
//			log.Error(err)
//			return nil, err
//		}
//		if parentInfo.Title != "" {
//			parentTitle = parentInfo.Title
//		}
//		// 飞书群聊，更新任务时，暂时不用展示父任务信息。
//		if info.PushType != consts.PushTypeUpdateIssue {
//			cardDivs = append(cardDivs, &commonvo.CardDiv{
//				Fields: []*commonvo.CardField{
//					&commonvo.CardField{
//						Key:   consts.CardElementParent,
//						Value: parentTitle,
//					},
//				},
//			})
//		}
//	}
//	ownerInfos, err := orgfacade.GetBaseUserInfoBatchRelaxed(orgId, businees.LcMemberToUserIds(issueInfo.OwnerId))
//	if err != nil {
//		log.Error(err)
//		return nil, err
//	}
//	ownerSlice := []string{}
//	for _, ownerInfo := range ownerInfos {
//		ownerSlice = append(ownerSlice, ownerInfo.Name)
//	}
//	ownerDisplayName := strings.Join(ownerSlice, "，")
//	if ownerDisplayName == "" {
//		ownerDisplayName = consts.CardDefaultOwnerNameForUpdateIssue
//	}
//	// 负责人
//	cardDivs = append(cardDivs, &commonvo.CardDiv{
//		Fields: []*commonvo.CardField{
//			&commonvo.CardField{
//				Key:   consts.CardElementOwner,
//				Value: ownerDisplayName,
//			},
//		},
//	})
//
//	cardMeta.Divs = cardDivs
//
//	tableColumns, errSys := GetTableColumnsMap(info.OrgId, info.TableId, nil)
//	if errSys != nil {
//		log.Errorf("[dealPushMsg] 获取表头失败 org:%d proj:%d table:%d user:%d, err: %v",
//			info.OrgId, info.ProjectId, info.TableId, info.OperatorId, errSys)
//		return nil, errSys
//	}
//	headers := make(map[string]lc_table.LcCommonField, 0)
//	errCopy := copyer.Copy(tableColumns, &headers)
//	if errCopy != nil {
//		return nil, errs.BuildSystemErrorInfo(errs.ObjectCopyError, errCopy)
//	}
//
//	tableColumnMetas := make(map[string]lc_table.LcCommonField, 0)
//
//	for columnId, tableColumn := range headers {
//		if tableColumn.Field.Props.PushMsg {
//			tableColumnMetas[columnId] = tableColumn
//		}
//	}
//
//	switch info.PushType {
//	case consts.PushTypeCreateIssue:
//		title += "创建了新" + taskTypeName
//	case consts.PushTypeUpdateIssue, consts.PushTypeCreateWorkHour, consts.PushTypeUpdateWorkHour, consts.PushTypeDeleteWorkHour:
//		title += "更新了" + taskTypeName
//		//elements, err = assemblyIssueUpdateAndWorkHour(info.PushType, info.OrgId, info.TableId, info, elements)
//		//if err != nil {
//		//log.Errorf("[dealPushMsg] err:%v", err)
//		//return nil, err
//		//}
//		//elements, _, _ = AssemblyFsCardForUpdateIssue(headers, &info, elements)
//	//case consts.PushTypeUpdateIssueStatus:
//	//	title += "更新了" + taskTypeName
//	//	elements = AssemblyFsCardForUpdateIssueStatus(headers, &info, elements)
//	//case consts.PushTypeCreateWorkHour, consts.PushTypeUpdateWorkHour, consts.PushTypeDeleteWorkHour: // 工时更新的卡片推送
//	//	title += "更新了" + taskTypeName
//	//	elements, _ = AssemblyFsCardForUpdateIssueWorkHour(headers, &info, elements)
//	case consts.PushTypeIssueComment:
//		title += "评论了" + taskTypeName
//		speakerNameDesc := fmt.Sprintf("%s:", operatorBaseInfo.Name)
//		comment := GetFeiShuAtContent(orgId, info.Ext.CommentBo.Content)
//		// 不截取@成员的字数
//		//words := consts.GetCardCommentRealWords(comment)
//		//comment = consts.TruncateText(words, consts.GrouChatIssueChangeDescLimitPerLine)
//		if info.Ext.ResourceInfo != nil && len(info.Ext.ResourceInfo) > 0 {
//			comment += " [附件]"
//		}
//		//commentCutter := str.NewCommentCutter(comment)
//		//comment = commentCutter.CutComment(consts.GroupChatIssueCommentLimit)
//		comment = CutComment(comment, consts.GroupChatIssueCommentLimit)
//		cardMeta.Divs = append(cardMeta.Divs, &commonvo.CardDiv{
//			Fields: []*commonvo.CardField{
//				&commonvo.CardField{
//					Key:   speakerNameDesc,
//					Value: fmt.Sprintf(consts.CardDoubleQuotationMarks, comment),
//				},
//			},
//		})
//	}
//	if title == operatorBaseInfo.Name+" " {
//		err = errs.BuildSystemErrorInfoWithMessage(errs.NotSupportTypeForPushFeishuChat, fmt.Sprintf("pushType: %d", info.PushType))
//		//log.Errorf("不支持的推送类型，param: %s", json.ToJsonIgnoreError(info))
//		return nil, err
//	}
//
//	cardMeta.ActionMarkdowns = []string{fmt.Sprintf(consts.MarkdownLink, consts.CardButtonTextForViewDetail, issueLinks.Link)}
//	cardMeta.FsActionElements = []interface{}{
//		vo.CardElementActionModule{
//			Tag: "action",
//			Actions: []interface{}{
//				vo.ActionButton{
//					Tag: "button",
//					Text: vo.CardElementText{
//						Tag:     "plain_text",
//						Content: consts.CardButtonTextForViewDetail,
//					},
//					Url:  issueLinks.SideBarLink,
//					Type: consts.FsCardButtonColorPrimary,
//				},
//				vo.ActionButton{
//					Tag: "button",
//					Text: vo.CardElementText{
//						Tag:     "plain_text",
//						Content: consts.CardButtonTextForViewInsideApp,
//					},
//					Url:  issueLinks.Link,
//					Type: consts.FsCardButtonColorDefault,
//				},
//			},
//		},
//	}
//
//	cardMeta.Title = title
//
//	log.Infof("[dealPushMsg] cardInfo:%s", json.ToJsonIgnoreError(cardMeta))
//	return cardMeta, nil
//}

func AssemblyFsCardForUpdateIssue(headers map[string]lc_table.LcCommonField, issueTrendsBo *bo.IssueTrendsBo, card *pushPb.TemplateCard) ([]string, []string) {
	add, notAdd := AssemblyUpdateContentForUpdateIssue(headers, issueTrendsBo, card)
	return add, notAdd
}

func AssemblyUpdateContentForUpdateIssue(headers map[string]lc_table.LcCommonField, issueTrendsBo *bo.IssueTrendsBo, card *pushPb.TemplateCard) ([]string, []string) {
	var err errs.SystemErrorInfo

	// 添加成员的更新和非成员添加的更新，用于后续去重判断
	addMemberColumnNames := make([]string, 0)
	notAddMemberColumnNames := make([]string, 0)

	changeMap := map[string]bo.TrendChangeListBo{}
	// Ext.ChangeList 中的都是自定义列
	for _, listBo := range issueTrendsBo.Ext.ChangeList {
		changeMap[listBo.Field] = listBo
	}

	// 组装附件的变更内容
	tmpChangedColumns, err := assemblyDocumentMsgForUpdateIssue(headers, issueTrendsBo, card)
	if err != nil {
		log.Errorf("[AssemblyUpdateContentForUpdateIssue] org:%d, proj:%d table:%d user:%d, assemblyUpdateIssueMsgForSpecialColumn err: %v",
			issueTrendsBo.OrgId, issueTrendsBo.ProjectId, issueTrendsBo.TableId, issueTrendsBo.OperatorId, err)
	}
	notAddMemberColumnNames = append(notAddMemberColumnNames, tmpChangedColumns...)

	// 组装图片的变更内容
	tmpChangedColumns, err = assemblyImageMsg(headers, issueTrendsBo, card)
	if err != nil {
		log.Errorf("[AssemblyUpdateContentForUpdateIssue] org:%d, proj:%d table:%d user:%d, dealPictureMsg err: %v",
			issueTrendsBo.OrgId, issueTrendsBo.ProjectId, issueTrendsBo.TableId, issueTrendsBo.OperatorId, err)
	}
	notAddMemberColumnNames = append(notAddMemberColumnNames, tmpChangedColumns...)

	tmpChangedColumns, err = assemblyRelationIssueMsg(headers, issueTrendsBo, card)
	if err != nil {
		log.Errorf("[AssemblyUpdateContentForUpdateIssue] org:%d, proj:%d table:%d user:%d, assemblyRelationIssueMsg err: %v",
			issueTrendsBo.OrgId, issueTrendsBo.ProjectId, issueTrendsBo.TableId, issueTrendsBo.OperatorId, err)
	}
	notAddMemberColumnNames = append(notAddMemberColumnNames, tmpChangedColumns...)

	tmpAddMemberColumns, err := assemblyMembersMsg(headers, issueTrendsBo, card)
	if err != nil {
		log.Errorf("[AssemblyUpdateContentForUpdateIssue] org:%d, proj:%d table:%d user:%d, assemblyMembersMsg err: %v",
			issueTrendsBo.OrgId, issueTrendsBo.ProjectId, issueTrendsBo.TableId, issueTrendsBo.OperatorId, err)
	}
	addMemberColumnNames = append(addMemberColumnNames, tmpAddMemberColumns...)

	tmpChangedColumns = AssemblyFsCardForUpdateIssueWorkHour(headers, issueTrendsBo, card)
	notAddMemberColumnNames = append(notAddMemberColumnNames, tmpChangedColumns...)

	// 其他字段的处理
	for _, changeItem := range changeMap {
		// 前面已经处理过负责人、关注人、确认人、工时等，这里就跳过。
		if exist, _ := slice.Contain([]string{
			consts.BasicFieldOwnerId, consts.BasicFieldFollowerIds, consts.BasicFieldAuditorIds, consts.ProBasicFieldWorkHour,
		},
			changeItem.Field); changeItem.IsForWorkHour || exist {
			continue
		}
		// 收集新增成员的列，用于后续判断是否只更改了成员列
		if changeItem.FieldType == consts.LcColumnFieldTypeMember {
			addMemberColumnNames = append(addMemberColumnNames, changeItem.Field)
		} else {
			notAddMemberColumnNames = append(notAddMemberColumnNames, changeItem.Field)
		}
		AssemblyCardForUpdateIssueOtherColumn(headers, &changeItem, card)
	}
	addMemberColumnNames = append(addMemberColumnNames, tmpAddMemberColumns...)
	addMemberColumnNames = slice.SliceUniqueString(addMemberColumnNames)
	notAddMemberColumnNames = slice.SliceUniqueString(notAddMemberColumnNames)

	// addMemberColumnNames 该请求中，新增成员的成员字段名
	// notAddMemberColumnNames 该请求中，编辑列的列名（非成员列）
	return addMemberColumnNames, notAddMemberColumnNames
}

// AssemblyFsCardForUpdateIssueWorkHour 任务工时更新后的卡片组装
func AssemblyFsCardForUpdateIssueWorkHour(headers map[string]lc_table.LcCommonField, issueTrendsBo *bo.IssueTrendsBo, cardMeta *pushPb.TemplateCard) []string {
	changedColumns := make([]string, 0)
	changeDescBuilder := strings.Builder{}
	curColumn, ok := headers[consts.ProBasicFieldWorkHour]
	if !ok {
		log.Infof("[AssemblyFsCardForUpdateIssueWorkHour] 表头信息中没有工时字段信息 issueId: %d", issueTrendsBo.IssueId)
		return changedColumns
	}
	defaultOldWorkerName, defaultNewWorkerName := "佚名", "佚名"
	defaultWorkerName := "佚名"
	workHourTypeDesc := "预估工时"
	// 更新工时比较特殊，实际上是将一条工时记录中，若干个值的变更保存在 `issueTrendsBo.Ext.ChangeList` 中
	if !issueTrendsBo.UpdateWorkHour {
		return changedColumns
	}
	// 如果是更新工时，则需要遍历完 issueTrendsBo.Ext.ChangeList 后，进行一次聚合
	isUpdateWorkHour := false
	// 工时变更的类型需要通过 trendsBo 中的 newValue 和 oldValue 来判定
	workHourChangeType := ""
	if len(issueTrendsBo.BeforeWorkHourIds) == 0 && len(issueTrendsBo.AfterWorkHourIds) > 0 {
		workHourChangeType = "create"
	} else if len(issueTrendsBo.BeforeWorkHourIds) > 0 && len(issueTrendsBo.AfterWorkHourIds) == 0 {
		workHourChangeType = "delete"
	} else {
		// 都不为空
		workHourChangeType = "update"
		isUpdateWorkHour = true
	}
	if changeDescBuilder.Len() > 0 {
		changeDescBuilder.WriteString("\n")
	}
	oldVal, newVal := issueTrendsBo.OldValue, issueTrendsBo.NewValue
	writeStr := ""
	switch workHourChangeType {
	case "create":
		workHourObj := vo1.OneWorkHourRecord{}
		json.FromJson(newVal, &workHourObj)
		userInfos, err := orgfacade.GetBaseUserInfoBatchRelaxed(issueTrendsBo.OrgId, []int64{workHourObj.Worker.UserID})
		if err != nil {
			log.Errorf("[AssemblyFsCardForUpdateIssueWorkHour] err: %v", err)
			return changedColumns
		}
		if len(userInfos) > 0 {
			defaultWorkerName = userInfos[0].Name
		}
		if workHourObj.Type == consts2.WorkHourTypeActual {
			workHourTypeDesc = "实际工时"
		}
		dateRangeDesc := RenderWorkHourTimeRangeDesc(workHourObj)
		// 添加了 xx 06月10日-06月11日的预估工时 3小时
		writeStr = fmt.Sprintf(consts.CardWorkHourAddDesc, workHourTypeDesc, defaultWorkerName, workHourObj.NeedTime,
			dateRangeDesc)
		break
	case "update":
		isUpdateWorkHour = true
		break
	case "delete":
		workHourObj := vo1.OneWorkHourRecord{}
		json.FromJson(oldVal, &workHourObj)
		userInfos, err := orgfacade.GetBaseUserInfoBatchRelaxed(issueTrendsBo.OrgId, []int64{workHourObj.Worker.UserID})
		if err != nil {
			log.Errorf("[AssemblyFsCardForUpdateIssueWorkHour] err: %v", err)
			return changedColumns
		}
		if len(userInfos) > 0 {
			defaultWorkerName = userInfos[0].Name
		}
		dateRangeDesc := RenderWorkHourTimeRangeDesc(workHourObj)
		if workHourObj.Type == consts2.WorkHourTypeActual {
			workHourTypeDesc = "实际工时"
		}
		writeStr = fmt.Sprintf(consts.CardWorkHourDelDesc, workHourTypeDesc, defaultWorkerName, workHourObj.NeedTime,
			dateRangeDesc)
		break
	}
	if len(writeStr) > 0 {
		changeDescBuilder.WriteString(writeStr)
	}

	if isUpdateWorkHour {
		writeStr := ""
		oldVal, newVal := issueTrendsBo.OldValue, issueTrendsBo.NewValue
		oldWorkHourObj := vo1.OneWorkHourRecord{}
		newWorkHourObj := vo1.OneWorkHourRecord{}
		json.FromJson(oldVal, &oldWorkHourObj)
		json.FromJson(newVal, &newWorkHourObj)
		userInfos, err := orgfacade.GetBaseUserInfoBatchRelaxed(issueTrendsBo.OrgId, []int64{oldWorkHourObj.Worker.UserID, newWorkHourObj.Worker.UserID})
		if err != nil {
			log.Errorf("[AssemblyFsCardForUpdateIssueWorkHour] err: %v", err)
			return changedColumns
		}
		userMap := make(map[int64]bo.BaseUserInfoBo, 2)
		for _, item := range userInfos {
			userMap[item.UserId] = item
		}
		if user, ok := userMap[oldWorkHourObj.Worker.UserID]; ok {
			defaultOldWorkerName = user.Name
		}
		if user, ok := userMap[newWorkHourObj.Worker.UserID]; ok {
			defaultNewWorkerName = user.Name
		}

		if oldWorkHourObj.Type == consts2.WorkHourTypeActual {
			workHourTypeDesc = "实际工时"
		}
		oldDateRangeDesc := RenderWorkHourTimeRangeDesc(oldWorkHourObj)
		newDateRangeDesc := RenderWorkHourTimeRangeDesc(newWorkHourObj)
		oldDesc := fmt.Sprintf("%s %s小时 %s", defaultOldWorkerName, oldWorkHourObj.NeedTime, oldDateRangeDesc)
		newDesc := fmt.Sprintf("%s %s小时 %s", defaultNewWorkerName, newWorkHourObj.NeedTime, newDateRangeDesc)
		// 修改了预估工时：苏三 1小时 01月01日 **→** 张六 2小时 01月02日
		writeStr = fmt.Sprintf(consts.CardWorkHourModify, workHourTypeDesc, oldDesc, newDesc)
		if len(writeStr) > 0 {
			changeDescBuilder.WriteString(writeStr)
		}
	}
	if changeDescBuilder.Len() > 0 {
		columnDisplayName := str.TruncateColumnName(curColumn)
		// markdownTile := RenderChangeDescForUpdateIssueOneDescInFsCard(columnDisplayName, changeDescBuilder.String())
		//tabStr := GetTabStrForColumnName(columnDisplayName)
		markdownTile := fmt.Sprintf("%s", changeDescBuilder.String())
		cardMeta.Divs = append(cardMeta.Divs, &pushPb.CardTextDiv{
			Fields: []*pushPb.CardTextField{
				&pushPb.CardTextField{
					Key:   columnDisplayName,
					Value: markdownTile,
				},
			},
		})
		changedColumns = append(changedColumns, consts.ProBasicFieldWorkHour)
	}

	return changedColumns
}

// RenderWorkHourTimeRangeDesc 展示工时信息时的日期文本展示
func RenderWorkHourTimeRangeDesc(workHour vo1.OneWorkHourRecord) string {
	startTime := time.Unix(workHour.StartTime, 0)
	defaultDateFormat := consts.FsCardWorkHourTimeRangeFormat
	// 06月10日-06月11日
	dateDesc := ""
	if workHour.StartTime == 0 && workHour.EndTime == 0 {
		// do nothing
	} else if workHour.StartTime == 0 && workHour.EndTime > 0 {
		endTime := time.Unix(workHour.EndTime, 0)
		dateDesc = fmt.Sprintf("%s", endTime.Format(defaultDateFormat))
	} else if workHour.StartTime > 0 && workHour.EndTime == 0 {
		dateDesc = fmt.Sprintf("%s", startTime.Format(defaultDateFormat))
	} else {
		// 如果起止时间相对跨年了，则展示出年份
		endTime := time.Unix(workHour.EndTime, 0)
		if startTime.Year() != endTime.Year() {
			defaultDateFormat = consts.FsCardWorkHourTimeRangeFormatWithYear
		}
		dateDesc = fmt.Sprintf("%s-%s", startTime.Format(defaultDateFormat), endTime.Format(defaultDateFormat))
	}

	return dateDesc
}

func GetFeiShuAtContent(orgId int64, content string) string {
	userIds := util.GetCommentAtUserIds(content)
	//userId -> openId
	openIdMap := map[string]string{}
	if len(userIds) > 0 {
		userInfos, err := orgfacade.GetBaseUserInfoBatchRelaxed(orgId, userIds)
		if err != nil {
			log.Error(err)
		} else {
			for _, userInfo := range userInfos {
				openIdMap[strconv.FormatInt(userInfo.UserId, 10)] = userInfo.OutUserId
			}
		}
	}

	return util.RenderCommentContentToMarkDownWithOpenIdMap(content, false, openIdMap)
}

// PushAuditToAuditStarter 飞书卡片-发送审批结果给发起人
func PushAuditToAuditStarter(orgId, userId int64, issueInfo *bo.IssueBo, auditStatus int, comment string, attachments []*vo1.AttachmentSimpleInfo) {
	if issueInfo.Title == "" {
		return
	}
	starterId := int64(0)
	// 只想审批发起人发送审批结果的卡片
	if uid, ok := issueInfo.LessData[consts.BasicFieldAuditStarter]; ok {
		starterId = cast.ToInt64(uid)
	}
	if starterId < 1 {
		log.Errorf("[PushAuditToIssueOwner] 发送审批结果给审批发起人，但审批发起人的 id 不合法。starterId: %d", starterId)
		return
	}
	if userId == starterId {
		log.Infof("[PushAuditToIssueOwner] 审批人就是审批发起人，无需发送卡片。")
		return
	}
	ownerIds := businees.LcMemberToUserIds(issueInfo.OwnerId)
	// 操作人和任务负责人
	userIds := []int64{
		userId, starterId,
	}
	userIds = append(userIds, ownerIds...)
	userIds = slice.SliceUniqueInt64(userIds)
	userResp := orgfacade.GetBaseUserInfoBatch(orgvo.GetBaseUserInfoBatchReqVo{
		OrgId:   orgId,
		UserIds: userIds,
	})
	if userResp.Failure() {
		log.Error(userResp.Error())
		return
	}
	orgResp := orgfacade.GetBaseOrgInfo(orgvo.GetBaseOrgInfoReqVo{
		OrgId: orgId,
	})
	if orgResp.Failure() {
		log.Error(orgResp.Error())
		return
	}
	// 表名
	tableName := consts.DefaultTableName
	if issueInfo.TableId > 0 {
		tableInfo, err := GetTableByTableId(issueInfo.OrgId, userId, issueInfo.TableId)
		if err != nil {
			log.Errorf("[PushAuditToIssueOwner] GetTableByTableId err:%v, tableId:%d", err, issueInfo.TableId)
			return
		}
		tableName = tableInfo.Name
	}
	projectName := consts.CardDefaultIssueProjectName
	if issueInfo.ProjectId != 0 {
		projectInfo, err := GetProjectSimple(orgId, issueInfo.ProjectId)
		if err != nil {
			log.Error(err)
			return
		}
		projectName = projectInfo.Name
	}

	ownerInfos := make([]*bo.BaseUserInfoBo, 0)
	operatorInfo := bo.BaseUserInfoBo{}
	usersMap := make(map[int64]bo.BaseUserInfoBo, len(userIds))
	for _, user := range userResp.BaseUserInfos {
		usersMap[user.UserId] = user
	}
	for _, owner := range ownerIds {
		if ownerInfo, ok := usersMap[owner]; ok {
			ownerInfos = append(ownerInfos, &ownerInfo)
		}
	}
	if info, ok := usersMap[userId]; ok {
		operatorInfo = info
	}

	statusName := "通过"
	if auditStatus == consts.AuditStatusReject {
		statusName = "驳回"
	}

	parentTitle := ""
	if issueInfo.ParentId > 0 {
		parentIssueInfo, err := GetIssueInfosLc(orgId, 0, []int64{issueInfo.ParentId})
		if err != nil {
			log.Errorf("[PushAuditToAuditStarter] GetIssueInfosLc err:%v, orgId:%d, issueParentId:%d", err, orgId, issueInfo.ParentId)
		}
		if len(parentIssueInfo) < 1 {
			log.Errorf("[PushAuditToAuditStarter] not found issue parentId:%v", issueInfo.ParentId)
			return
		}
		parent := parentIssueInfo[0]
		parentTitle = parent.Title
		if parentTitle == "" {
			parentTitle = consts.CardDefaultRelationIssueTitle
		}
	}

	tableColumns, errSys := GetTableColumnsMap(orgId, issueInfo.TableId, []string{consts.BasicFieldTitle, consts.BasicFieldOwnerId})
	if errSys != nil {
		log.Errorf("[pushChatMsg] GetTableColumnConfig err:%v, orgId:%v, issueId:%v", errSys, orgId, issueInfo.Id)
		return
	}

	cardDivs := []*pushPb.CardTextDiv{}

	if comment != "" {
		cardDivs = append(cardDivs, &pushPb.CardTextDiv{
			Fields: []*pushPb.CardTextField{
				&pushPb.CardTextField{
					Key:   consts.CardElementEvaluate,
					Value: comment,
				},
			},
		})
	}
	if attachments != nil && len(attachments) > 0 {
		attachmentName := ""
		for _, attachment := range attachments {
			if attachment.Name != "" {
				attachmentName += attachment.Name + " "
			}
		}
		cardDivs = append(cardDivs, &pushPb.CardTextDiv{
			Fields: []*pushPb.CardTextField{
				&pushPb.CardTextField{
					Key:   consts.CardElementAttachment,
					Value: attachmentName,
				},
			},
		})
	}

	links := GetIssueLinks("", orgId, issueInfo.Id)
	title := fmt.Sprintf("%s%s了审批", operatorInfo.Name, statusName)
	// 卡片标题颜色
	headTitleColorLevel := pushPb.CardElementLevel_Default
	if auditStatus == consts.AuditStatusPass {
		headTitleColorLevel = pushPb.CardElementLevel_Default
	} else if auditStatus == consts.AuditStatusReject {
		headTitleColorLevel = pushPb.CardElementLevel_Important
	}

	sendUserOpenIds := make([]string, 0)
	if user, ok := usersMap[starterId]; ok {
		sendUserOpenIds = append(sendUserOpenIds, user.OutUserId)
	}

	cardMsg := card.GenerateCard(&commonvo.GenerateCard{
		OrgId:             orgId,
		CardTitle:         title,
		CardLevel:         headTitleColorLevel,
		ParentTitle:       parentTitle,
		ProjectName:       projectName,
		TableName:         tableName,
		IssueTitle:        issueInfo.Title,
		ProjectId:         issueInfo.ProjectId,
		ParentId:          issueInfo.ParentId,
		OwnerInfos:        ownerInfos,
		Links:             links,
		IsNeedUnsubscribe: false,
		SpecialDiv:        cardDivs,
		TableColumnInfo:   tableColumns,
	})

	errSys = card.PushCard(orgId, &commonvo.PushCard{
		OrgId:         orgId,
		OutOrgId:      orgResp.BaseOrgInfo.OutOrgId,
		SourceChannel: orgResp.BaseOrgInfo.SourceChannel,
		OpenIds:       sendUserOpenIds,
		CardMsg:       cardMsg,
	})
	if errSys != nil {
		log.Errorf("[PushAuditToAuditStarter] err:%v", errSys)
		return
	}

}

func PushCardTodoAuditDone(orgId, userId int64, issueInfo *bo.IssueBo, todoType automationPb.TodoType, msg string) {
	if issueInfo.Title == "" {
		return
	}
	starterId := int64(0)
	// 只想审批发起人发送审批结果的卡片
	if uid, ok := issueInfo.LessData[consts.BasicFieldAuditStarter]; ok {
		starterId = cast.ToInt64(uid)
	}
	if starterId < 1 {
		log.Errorf("[PushAuditToIssueOwner] 发送审批结果给审批发起人，但审批发起人的 id 不合法。starterId: %d", starterId)
		return
	}
	if userId == starterId {
		log.Infof("[PushAuditToIssueOwner] 审批人就是审批发起人，无需发送卡片。")
		return
	}
	ownerIds := businees.LcMemberToUserIds(issueInfo.OwnerId)
	// 操作人和任务负责人
	userIds := []int64{
		userId, starterId,
	}
	userIds = append(userIds, ownerIds...)
	userIds = slice.SliceUniqueInt64(userIds)
	userResp := orgfacade.GetBaseUserInfoBatch(orgvo.GetBaseUserInfoBatchReqVo{
		OrgId:   orgId,
		UserIds: userIds,
	})
	if userResp.Failure() {
		log.Error(userResp.Error())
		return
	}
	orgResp := orgfacade.GetBaseOrgInfo(orgvo.GetBaseOrgInfoReqVo{
		OrgId: orgId,
	})
	if orgResp.Failure() {
		log.Error(orgResp.Error())
		return
	}
	// 表名
	tableName := consts.DefaultTableName
	if issueInfo.TableId > 0 {
		tableInfo, err := GetTableByTableId(issueInfo.OrgId, userId, issueInfo.TableId)
		if err != nil {
			log.Errorf("[PushAuditToIssueOwner] GetTableByTableId err:%v, tableId:%d", err, issueInfo.TableId)
			return
		}
		tableName = tableInfo.Name
	}
	projectName := consts.CardDefaultIssueProjectName
	if issueInfo.ProjectId != 0 {
		projectInfo, err := GetProjectSimple(orgId, issueInfo.ProjectId)
		if err != nil {
			log.Error(err)
			return
		}
		projectName = projectInfo.Name
	}

	ownerInfos := make([]*bo.BaseUserInfoBo, 0)
	operatorInfo := bo.BaseUserInfoBo{}
	usersMap := make(map[int64]bo.BaseUserInfoBo, len(userIds))
	for _, user := range userResp.BaseUserInfos {
		usersMap[user.UserId] = user
	}
	for _, owner := range ownerIds {
		if ownerInfo, ok := usersMap[owner]; ok {
			ownerInfos = append(ownerInfos, &ownerInfo)
		}
	}
	if info, ok := usersMap[userId]; ok {
		operatorInfo = info
	}

	statusName := "通过"
	//if auditStatus == consts.AuditStatusReject {
	//	statusName = "驳回"
	//}

	parentTitle := ""
	if issueInfo.ParentId > 0 {
		parentIssueInfo, err := GetIssueInfosLc(orgId, 0, []int64{issueInfo.ParentId})
		if err != nil {
			log.Errorf("[PushAuditToAuditStarter] GetIssueInfosLc err:%v, orgId:%d, issueParentId:%d", err, orgId, issueInfo.ParentId)
		}
		if len(parentIssueInfo) < 1 {
			log.Errorf("[PushAuditToAuditStarter] not found issue parentId:%v", issueInfo.ParentId)
			return
		}
		parent := parentIssueInfo[0]
		parentTitle = parent.Title
		if parentTitle == "" {
			parentTitle = consts.CardDefaultRelationIssueTitle
		}
	}

	tableColumns, errSys := GetTableColumnsMap(orgId, issueInfo.TableId, []string{consts.BasicFieldTitle, consts.BasicFieldOwnerId})
	if errSys != nil {
		log.Errorf("[pushChatMsg] GetTableColumnConfig err:%v, orgId:%v, issueId:%v", errSys, orgId, issueInfo.Id)
		return
	}

	cardDivs := []*pushPb.CardTextDiv{}

	if msg != "" {
		cardDivs = append(cardDivs, &pushPb.CardTextDiv{
			Fields: []*pushPb.CardTextField{
				&pushPb.CardTextField{
					Key:   consts.CardElementEvaluate,
					Value: msg,
				},
			},
		})
	}

	links := GetIssueLinks("", orgId, issueInfo.Id)
	title := fmt.Sprintf("%s%s了审批", operatorInfo.Name, statusName)
	// 卡片标题颜色
	headTitleColorLevel := pushPb.CardElementLevel_Default
	//if auditStatus == consts.AuditStatusPass {
	//	headTitleColorLevel = pushPb.CardElementLevel_Default
	//} else if auditStatus == consts.AuditStatusReject {
	//	headTitleColorLevel = pushPb.CardElementLevel_Important
	//}

	sendUserOpenIds := make([]string, 0)
	if user, ok := usersMap[starterId]; ok {
		sendUserOpenIds = append(sendUserOpenIds, user.OutUserId)
	}

	cardMsg := card.GenerateCard(&commonvo.GenerateCard{
		OrgId:             orgId,
		CardTitle:         title,
		CardLevel:         headTitleColorLevel,
		ParentTitle:       parentTitle,
		ProjectName:       projectName,
		TableName:         tableName,
		IssueTitle:        issueInfo.Title,
		ProjectId:         issueInfo.ProjectId,
		ParentId:          issueInfo.ParentId,
		OwnerInfos:        ownerInfos,
		Links:             links,
		IsNeedUnsubscribe: false,
		SpecialDiv:        cardDivs,
		TableColumnInfo:   tableColumns,
	})

	errSys = card.PushCard(orgId, &commonvo.PushCard{
		OrgId:         orgId,
		OutOrgId:      orgResp.BaseOrgInfo.OutOrgId,
		SourceChannel: orgResp.BaseOrgInfo.SourceChannel,
		OpenIds:       sendUserOpenIds,
		CardMsg:       cardMsg,
	})
	if errSys != nil {
		log.Errorf("[PushAuditToAuditStarter] err:%v", errSys)
		return
	}

}

func RenderChangeDescForUpdateIssueInCard(oldDesc, newDesc string) string {
	// 排除 `：/~/**→**` 等 8 个非内容字符
	limitCharNum := consts.GrouChatIssueChangeDescLimitPerLine - 8
	if (str.RuneLen(oldDesc) + str.RuneLen(newDesc)) > limitCharNum {
		halfNum := limitCharNum / 2
		if str.RuneLen(oldDesc) > halfNum {
			oldDesc = str.Substr(oldDesc, 0, halfNum) + consts.CardIssueColumnNameOverflow
		}
		if str.RuneLen(newDesc) > halfNum {
			newDesc = str.Substr(newDesc, 0, halfNum) + consts.CardIssueColumnNameOverflow
		}
	}
	// eg：修改了 任务状态 进行中 -> 已完成
	// fmt.Sprintf("%s：~~%s~~ **→** %s", displayName, oldDesc, newDesc)
	oneLineStr := fmt.Sprintf(consts.CardCommonChangeDesc, oldDesc, newDesc)

	return oneLineStr
}

// 计算制表符
func GetTabStrForColumnName(columnName string) string {
	return consts.GetTabCharacter(columnName)
}

func RenderChangeDescForUpdateIssueOneDescInCard(actionStr, changeDesc string) string {
	if str.RuneLen(changeDesc) > consts.GrouChatIssueChangeDescLimitPerLine {
		changeDesc = str.Substr(changeDesc, 0, consts.GrouChatIssueChangeDescLimitPerLine) + consts.CardIssueChangeDescTextOverflow
	}
	oneLineStr := fmt.Sprintf(actionStr, changeDesc)
	return oneLineStr

}

// 任务更新时，飞书卡片内容组装-附件部分
func assemblyDocumentMsgForUpdateIssue(headers map[string]lc_table.LcCommonField, info *bo.IssueTrendsBo, cardMeta *pushPb.TemplateCard) ([]string, errs.SystemErrorInfo) {
	changeColumns := make([]string, 0)
	// 附件是否有更新的
	beforeChangeDocuments := info.BeforeChangeDocuments
	afterChangeDocuments := info.AfterChangeDocuments

	log.Infof("[assemblyDocumentMsgForUpdateIssue], headerMap:%s, beforeChangeDocuments:%v, afterChangeDocuments:%v",
		json.ToJsonIgnoreError(headers), beforeChangeDocuments, afterChangeDocuments)

	beforeColumns := make([]string, 0)
	afterColumns := make([]string, 0)
	for columnKey, _ := range beforeChangeDocuments {
		beforeColumns = append(beforeColumns, columnKey)
	}
	for columnKey, _ := range afterChangeDocuments {
		afterColumns = append(afterColumns, columnKey)
	}
	// 列的新增/删除
	_, addColumns, delColumns := int642.CompareSliceAddDelString(afterColumns, beforeColumns)
	updateColumns := int642.StringArrIntersect(afterColumns, beforeColumns)
	if len(addColumns) > 0 {
		for _, columnKey := range addColumns {
			curColumn, ok := headers[columnKey]
			if !ok {
				break
			}
			displayName := str.TruncateColumnName(curColumn)
			afterVal, _ := afterChangeDocuments[columnKey]
			afterDocs := make(map[string]lc_table.LcDocumentValue)
			err := copyer.Copy(afterVal, &afterDocs)
			if err != nil {
				log.Errorf("[dealUpdateIssueMsgForSpecialColumn] err: %v", err)
				return changeColumns, errs.BuildSystemErrorInfo(errs.ObjectCopyError, err)
			}
			addDocNames := []string{}
			for _, resource := range afterDocs {
				addDocNames = append(addDocNames, resource.Name)
			}
			writeStr := strings.Join(addDocNames, "、")
			markdownTile := RenderChangeDescForUpdateIssueOneDescInCard(consts.CardDocAndImageAddDesc, writeStr)
			cardMeta.Divs = append(cardMeta.Divs, &pushPb.CardTextDiv{
				Fields: []*pushPb.CardTextField{
					&pushPb.CardTextField{
						Key:   displayName,
						Value: markdownTile,
					},
				},
			})
			changeColumns = append(changeColumns, columnKey)
		}
	}
	if len(delColumns) > 0 {
		for _, columnKey := range delColumns {
			curColumn, ok := headers[columnKey]
			if !ok {
				break
			}
			displayName := str.TruncateColumnName(curColumn)
			beforeVal, _ := beforeChangeDocuments[columnKey]
			beforeDocs := make(map[string]lc_table.LcDocumentValue)
			err := copyer.Copy(beforeVal, &beforeDocs)
			if err != nil {
				log.Errorf("[dealUpdateIssueMsgForSpecialColumn] err: %v", err)
				return changeColumns, errs.BuildSystemErrorInfo(errs.ObjectCopyError, err)
			}
			delDocNames := []string{}
			for _, resource := range beforeDocs {
				delDocNames = append(delDocNames, resource.Name)
			}
			writeStr := strings.Join(delDocNames, "、")
			markdownTile := RenderChangeDescForUpdateIssueOneDescInCard(consts.CardDocAndImageDelDesc, writeStr)
			cardMeta.Divs = append(cardMeta.Divs, &pushPb.CardTextDiv{
				Fields: []*pushPb.CardTextField{
					&pushPb.CardTextField{
						Key:   displayName,
						Value: markdownTile,
					},
				},
			})

			changeColumns = append(changeColumns, columnKey)
		}
	}
	if len(updateColumns) > 0 {
		for _, columnKey := range updateColumns {
			afterVal, _ := afterChangeDocuments[columnKey]
			beforeVal, _ := beforeChangeDocuments[columnKey]
			beforeDocs := make(map[string]lc_table.LcDocumentValue)
			err := copyer.Copy(beforeVal, &beforeDocs)
			if err != nil {
				log.Errorf("[assemblyDocumentMsgForUpdateIssue] err: %v", err)
				return changeColumns, errs.BuildSystemErrorInfo(errs.ObjectCopyError, err)
			}
			afterDocs := make(map[string]lc_table.LcDocumentValue)
			err = copyer.Copy(afterVal, &afterDocs)
			if err != nil {
				log.Errorf("[assemblyDocumentMsgForUpdateIssue] err: %v", err)
				return changeColumns, errs.BuildSystemErrorInfo(errs.ObjectCopyError, err)
			}
			beforeResourceIds := make([]string, 0)
			afterResourceIds := make([]string, 0)
			for resourceId, _ := range beforeDocs {
				beforeResourceIds = append(beforeResourceIds, resourceId)
			}
			for resourceId, _ := range afterDocs {
				afterResourceIds = append(afterResourceIds, resourceId)
			}
			isSame, add, del := int642.CompareSliceAddDelString(afterResourceIds, beforeResourceIds)
			if isSame {
				break
			}
			curColumn, ok := headers[columnKey]
			if !ok {
				break
			}
			displayName := str.TruncateColumnName(curColumn)
			if len(add) > 0 {
				addResourceNames := []string{}
				for _, resourceId := range add {
					tmpResource, ok := afterDocs[resourceId]
					if !ok {
						continue
					}
					addResourceNames = append(addResourceNames, tmpResource.Name)
				}
				//writeStr := fmt.Sprintf("添加了 “%s”", strings.Join(addResourceNames, "、"))
				markdownTile := RenderChangeDescForUpdateIssueOneDescInCard(consts.CardDocAndImageAddDesc, strings.Join(addResourceNames, "、"))
				cardMeta.Divs = append(cardMeta.Divs, &pushPb.CardTextDiv{
					Fields: []*pushPb.CardTextField{
						&pushPb.CardTextField{
							Key:   displayName,
							Value: markdownTile,
						},
					},
				})

			}
			if len(del) > 0 {
				log.Infof("[assemblyDocumentMsgForUpdateIssue] del:%v", del)
				delResourceNames := []string{}
				for _, resourceId := range del {
					tmpResource, ok := beforeDocs[resourceId]
					if !ok {
						continue
					}
					delResourceNames = append(delResourceNames, tmpResource.Name)

				}
				delWords := strings.Join(delResourceNames, "、")
				markdownTile := RenderChangeDescForUpdateIssueOneDescInCard(consts.CardDocAndImageDelDesc, delWords)
				//writeStr := fmt.Sprintf("删除了 ~~“%s”~~", strings.Join(delResourceNames, "、"))
				//markdownTile := RenderChangeDescForUpdateIssueOneDescInFsCard(displayName, writeStr, true)
				cardMeta.Divs = append(cardMeta.Divs, &pushPb.CardTextDiv{
					Fields: []*pushPb.CardTextField{
						&pushPb.CardTextField{
							Key:   displayName,
							Value: markdownTile,
						},
					},
				})
			}

			changeColumns = append(changeColumns, columnKey)
		}
	}

	return changeColumns, nil
}

// 任务更新时，飞书卡片内容组装-图片部分
func assemblyImageMsg(headerMap map[string]lc_table.LcCommonField, info *bo.IssueTrendsBo, cardMeta *pushPb.TemplateCard) ([]string, errs.SystemErrorInfo) {
	log.Infof("[assemblyImageMsg] headerMap:%s, beforeChangeImages:%v, afterChangeImages:%v", json.ToJsonIgnoreError(headerMap), info.BeforeChangeImages, info.AfterChangeImages)
	changedColumns := make([]string, 0)
	// 处理图片
	beforeChangeImages := info.BeforeChangeImages
	afterChangeImages := info.AfterChangeImages

	newImages := make([]lc_table.LcDocumentValue, 0)
	oldImages := make([]lc_table.LcDocumentValue, 0)

	for _, v := range beforeChangeImages {
		err := copyer.Copy(v, &oldImages)
		if err != nil {
			return changedColumns, errs.ObjectCopyError
		}
	}
	for _, v := range afterChangeImages {
		err := copyer.Copy(v, &newImages)
		if err != nil {
			return changedColumns, errs.ObjectCopyError
		}
	}
	if (beforeChangeImages == nil && afterChangeImages == nil) || len(newImages) == len(oldImages) {
		return changedColumns, nil
	}

	imageColumns := make([]string, 0)
	for columnKey, _ := range beforeChangeImages {
		imageColumns = append(imageColumns, columnKey)
	}
	for columnKey, _ := range afterChangeImages {
		imageColumns = append(imageColumns, columnKey)
	}
	imageColumns = slice.SliceUniqueString(imageColumns)

	afterImages := map[int64]lc_table.LcDocumentValue{}
	beforeImages := map[int64]lc_table.LcDocumentValue{}
	for _, columnKey := range imageColumns {
		beforeResourceIds := make([]int64, 0)
		afterResourceIds := make([]int64, 0)
		for _, v := range oldImages {
			beforeResourceIds = append(beforeResourceIds, v.Id)
			beforeImages[v.Id] = v
		}
		for _, v := range newImages {
			afterResourceIds = append(afterResourceIds, v.Id)
			afterImages[v.Id] = v
		}
		isSame, add, del := int642.CompareSliceAddDelInt64(afterResourceIds, beforeResourceIds)
		if isSame {
			break
		}
		curColumn, ok := headerMap[columnKey]
		if !ok {
			break
		}
		displayName := str.TruncateColumnName(curColumn)

		if len(add) > 0 {
			addResourceNames := []string{}
			for _, resourceId := range add {
				tmpResource, ok := afterImages[resourceId]
				if !ok {
					continue
				}
				addResourceNames = append(addResourceNames, tmpResource.Name)
			}
			//writeStr := fmt.Sprintf("添加了 “%s”", strings.Join(addResourceNames, "、"))
			markdownTile := RenderChangeDescForUpdateIssueOneDescInCard(consts.CardDocAndImageAddDesc, strings.Join(addResourceNames, "、"))
			cardMeta.Divs = append(cardMeta.Divs, &pushPb.CardTextDiv{
				Fields: []*pushPb.CardTextField{
					&pushPb.CardTextField{
						Key:   displayName,
						Value: markdownTile,
					},
				},
			})
		}

		if len(del) > 0 {
			delResourceNames := []string{}
			for _, resourceId := range del {
				tmpResource, ok := beforeImages[resourceId]
				if !ok {
					continue
				}
				delResourceNames = append(delResourceNames, tmpResource.Name)

			}
			markdownTile := RenderChangeDescForUpdateIssueOneDescInCard(consts.CardDocAndImageDelDesc, strings.Join(delResourceNames, "、"))
			cardMeta.Divs = append(cardMeta.Divs, &pushPb.CardTextDiv{
				Fields: []*pushPb.CardTextField{
					&pushPb.CardTextField{
						Key:   displayName,
						Value: markdownTile,
					},
				},
			})
		}

		changedColumns = append(changedColumns, columnKey)
	}

	return changedColumns, nil
}

// 更新任务时，组装前后置/关联的更新内容（飞书卡片）
func assemblyRelationIssueMsg(headerMap map[string]lc_table.LcCommonField, issueTrendsBo *bo.IssueTrendsBo, cardMeta *pushPb.TemplateCard) ([]string, errs.SystemErrorInfo) {
	changedColumns := make([]string, 0)
	relatingList := []map[string]*bo.RelatingChangeBo{issueTrendsBo.RelatingChange, issueTrendsBo.SingleRelatingChange, issueTrendsBo.BaRelatingChange}
	var allRelateIssues []int64

	for _, relating := range relatingList {
		for _, changeBo := range relating {
			allRelateIssues = append(allRelateIssues, changeBo.LinkToDel...)
			allRelateIssues = append(allRelateIssues, changeBo.LinkToAdd...)
			allRelateIssues = append(allRelateIssues, changeBo.LinkFromDel...)
			allRelateIssues = append(allRelateIssues, changeBo.LinkFromAdd...)
		}
	}

	allRelateIssues = slice.SliceUniqueInt64(allRelateIssues)
	if len(allRelateIssues) == 0 {
		return changedColumns, nil
	}

	// 关联、前后置字段无码化之后，这里传过来的id已经是dataid，而非issueid
	issueInfo := projectfacade.GetIssueInfoList(projectvo.IssueInfoListReqVo{
		OrgId:    issueTrendsBo.OrgId,
		UserId:   issueTrendsBo.OperatorId,
		IssueIds: allRelateIssues,
	})
	if issueInfo.Failure() {
		log.Errorf("[assemblyRelationIssueMsg] err: %v, orgId: %d", issueInfo.Error(), issueTrendsBo.OrgId)
		return changedColumns, issueInfo.Error()
	}
	issueInfoById := map[int64]vo1.Issue{}
	for _, v := range issueInfo.IssueInfos {
		issueInfoById[cast.ToInt64(v.ID)] = v
	}

	relatingDesc := &bo.RelatingDesc{
		AddToDesc:   consts.CardRelatingAddDesc,
		DelToDesc:   consts.CardRelatingDelDesc,
		AddFromDesc: consts.CardRelatingAddDesc,
		DelFromDesc: consts.CardRelatingDelDesc,
	}
	baRelatingDesc := &bo.RelatingDesc{
		AddToDesc:   consts.CardBaRelatingLinkToAddDesc,
		DelToDesc:   consts.CardBaRelatingLinkToDelDesc,
		AddFromDesc: consts.CardBaRelatingLinkFromAddDesc,
		DelFromDesc: consts.CardBaRelatingLinkFromDelDesc,
	}

	relatingDescList := []*bo.RelatingDesc{relatingDesc, relatingDesc, baRelatingDesc}
	for i, relatingChange := range relatingList {
		for columnId, changeBo := range relatingChange {
			if relatingColumn, ok := headerMap[columnId]; ok {
				displayName := str.TruncateColumnName(relatingColumn)
				if len(changeBo.LinkToDel) > 0 { // 删除关联
					addRelatingCard(changeBo.LinkToDel, displayName, relatingDescList[i].DelToDesc, cardMeta, issueInfoById)
					changedColumns = append(changedColumns, columnId)
				}
				if len(changeBo.LinkToAdd) > 0 { // 新增关联
					addRelatingCard(changeBo.LinkToAdd, displayName, relatingDescList[i].AddToDesc, cardMeta, issueInfoById)
					changedColumns = append(changedColumns, columnId)
				}
				if len(changeBo.LinkFromAdd) > 0 {
					addRelatingCard(changeBo.LinkFromAdd, displayName, relatingDescList[i].AddFromDesc, cardMeta, issueInfoById)
					changedColumns = append(changedColumns, columnId)
				}
				if len(changeBo.LinkFromDel) > 0 {
					addRelatingCard(changeBo.LinkFromDel, displayName, relatingDescList[i].DelFromDesc, cardMeta, issueInfoById)
					changedColumns = append(changedColumns, columnId)
				}
			}
		}
	}

	return changedColumns, nil
}

func addRelatingCard(ids []int64, displayName, desc string, cardMeta *pushPb.TemplateCard, issueInfoById map[int64]vo1.Issue) {
	issueTitles := make([]string, 0, len(ids))
	for _, dataId := range ids {
		tmpIssue, ok1 := issueInfoById[dataId]
		if !ok1 {
			continue
		}
		if tmpIssue.Title == "" {
			tmpIssue.Title = consts.CardDefaultRelationIssueTitle
		}
		issueTitles = append(issueTitles, tmpIssue.Title)
	}
	valStr := strings.Join(issueTitles, "、")
	markdownTile := RenderChangeDescForUpdateIssueOneDescInCard(desc, valStr)
	cardMeta.Divs = append(cardMeta.Divs, &pushPb.CardTextDiv{
		Fields: []*pushPb.CardTextField{
			&pushPb.CardTextField{
				Key:   displayName,
				Value: markdownTile,
			},
		},
	})
}

// 任务更新时，飞书卡片内容组装-成员变更部分
func assemblyMembersMsg(headMap map[string]lc_table.LcCommonField, issueTrendsBo *bo.IssueTrendsBo, cardMeta *pushPb.TemplateCard) ([]string, errs.SystemErrorInfo) {
	addMemberColumns := make([]string, 0, 3)
	beforeChangeFollowers := slice.SliceUniqueInt64(issueTrendsBo.BeforeChangeFollowers)
	afterChangeFollowers := slice.SliceUniqueInt64(issueTrendsBo.AfterChangeFollowers)
	beforeChangeAuditors := slice.SliceUniqueInt64(issueTrendsBo.BeforeChangeAuditors)
	afterChangeAuditors := slice.SliceUniqueInt64(issueTrendsBo.AfterChangeAuditors)
	beforeChangeOwners := slice.SliceUniqueInt64(issueTrendsBo.BeforeOwner)
	afterChangeOwners := slice.SliceUniqueInt64(issueTrendsBo.AfterOwner)

	if beforeChangeFollowers == nil {
		beforeChangeFollowers = []int64{}
	}
	if afterChangeFollowers == nil {
		afterChangeFollowers = []int64{}
	}
	if beforeChangeAuditors == nil {
		beforeChangeAuditors = []int64{}
	}
	if afterChangeAuditors == nil {
		afterChangeAuditors = []int64{}
	}
	if beforeChangeOwners == nil {
		beforeChangeOwners = []int64{}
	}
	if afterChangeOwners == nil {
		afterChangeOwners = []int64{}
	}

	var allRelateUserIds []int64
	allRelateUserIds = append(allRelateUserIds, beforeChangeFollowers...)
	allRelateUserIds = append(allRelateUserIds, afterChangeFollowers...)
	allRelateUserIds = append(allRelateUserIds, beforeChangeAuditors...)
	allRelateUserIds = append(allRelateUserIds, afterChangeAuditors...)
	allRelateUserIds = append(allRelateUserIds, beforeChangeOwners...)
	allRelateUserIds = append(allRelateUserIds, afterChangeOwners...)

	allRelateUserIds = slice.SliceUniqueInt64(allRelateUserIds)
	// 查出用户信息
	userInfos, err := orgfacade.GetBaseUserInfoBatchRelaxed(issueTrendsBo.OrgId, allRelateUserIds)
	if err != nil {
		log.Errorf("[assemblyMembersMsg]获取用户相关信息失败:%v, orgId:%d, userIds:%v", err, issueTrendsBo.OrgId, allRelateUserIds)
		return addMemberColumns, err
	}
	userIdsMap := map[int64]bo.SimpleUserInfoBo{}
	for _, v := range userInfos {
		userIdsMap[v.UserId] = bo.SimpleUserInfoBo{
			Id:     v.UserId,
			Name:   v.Name,
			Avatar: v.Avatar,
		}
	}

	// 关注人
	beforeFollowers := []string{}
	afterFollowers := []string{}
	if v, ok := headMap[consts.BasicFieldFollowerIds]; ok {
		if !int642.CompareSliceInt64(beforeChangeFollowers, afterChangeFollowers) {
			if len(beforeChangeFollowers) == 0 {
				beforeFollowers = append(beforeFollowers, consts.CardDefaultOwnerNameForUpdateIssue)
			}
			for _, f := range beforeChangeFollowers {
				beforeFollowers = append(beforeFollowers, userIdsMap[f].Name)
			}

			if len(afterChangeFollowers) == 0 {
				afterFollowers = append(afterFollowers, consts.CardDefaultOwnerNameForUpdateIssue)
			}

			for _, f := range afterChangeFollowers {
				afterFollowers = append(afterFollowers, userIdsMap[f].Name)
			}
			displayName := str.TruncateColumnName(v)
			markdownTile := RenderChangeDescForUpdateIssueInCard(strings.Join(beforeFollowers, "，"), strings.Join(afterFollowers, "，"))
			cardMeta.Divs = append(cardMeta.Divs, &pushPb.CardTextDiv{
				Fields: []*pushPb.CardTextField{
					&pushPb.CardTextField{
						Key:   displayName,
						Value: markdownTile,
					},
				},
			})

			if addUserIds := int642.ArrayDiff(afterChangeFollowers, beforeChangeFollowers); len(addUserIds) > 0 {
				addMemberColumns = append(addMemberColumns, consts.BasicFieldFollowerIds)
			}
		}
	}

	beforeOwners := []string{}
	afterOwners := []string{}
	if v, ok := headMap[consts.BasicFieldOwnerId]; ok {
		if !int642.CompareSliceInt64(beforeChangeOwners, afterChangeOwners) {
			if len(beforeChangeOwners) == 0 {
				beforeOwners = append(beforeOwners, consts.CardDefaultOwnerNameForUpdateIssue)
			}
			for _, f := range beforeChangeOwners {
				beforeOwners = append(beforeOwners, userIdsMap[f].Name)
			}

			if len(afterChangeOwners) == 0 {
				afterOwners = append(afterOwners, consts.CardDefaultOwnerNameForUpdateIssue)
			}

			for _, f := range afterChangeOwners {
				afterOwners = append(afterOwners, userIdsMap[f].Name)
			}

			displayName := str.TruncateColumnName(v)
			markdownTile := RenderChangeDescForUpdateIssueInCard(strings.Join(beforeOwners, "，"), strings.Join(afterOwners, "，"))
			cardMeta.Divs = append(cardMeta.Divs, &pushPb.CardTextDiv{
				Fields: []*pushPb.CardTextField{
					&pushPb.CardTextField{
						Key:   displayName,
						Value: markdownTile,
					},
				},
			})

			if addUserIds := int642.ArrayDiff(afterChangeOwners, beforeChangeOwners); len(addUserIds) > 0 {
				addMemberColumns = append(addMemberColumns, consts.BasicFieldOwnerId)
			}
		}
	}

	beforeAuditors := []string{}
	afterAuditors := []string{}
	if v, ok := headMap[consts.BasicFieldAuditorIds]; ok {
		if !int642.CompareSliceInt64(beforeChangeAuditors, afterChangeAuditors) {
			if len(beforeChangeAuditors) == 0 {
				beforeAuditors = append(beforeAuditors, consts.CardDefaultOwnerNameForUpdateIssue)
			}
			for _, f := range beforeChangeAuditors {
				beforeAuditors = append(beforeAuditors, userIdsMap[f].Name)
			}

			if len(afterChangeAuditors) == 0 {
				afterAuditors = append(afterAuditors, consts.CardDefaultOwnerNameForUpdateIssue)
			}

			for _, f := range afterChangeAuditors {
				afterAuditors = append(afterAuditors, userIdsMap[f].Name)
			}
			displayName := str.TruncateColumnName(v)
			markdownTile := RenderChangeDescForUpdateIssueInCard(strings.Join(beforeAuditors, "，"), strings.Join(afterAuditors, "，"))
			cardMeta.Divs = append(cardMeta.Divs, &pushPb.CardTextDiv{
				Fields: []*pushPb.CardTextField{
					&pushPb.CardTextField{
						Key:   displayName,
						Value: markdownTile,
					},
				},
			})
			if addUserIds := int642.ArrayDiff(afterChangeAuditors, beforeChangeAuditors); len(addUserIds) > 0 {
				addMemberColumns = append(addMemberColumns, consts.BasicFieldAuditorIds)
			}
		}
	}

	return addMemberColumns, nil
}

// 拼装卡片内容，包括标题、起止时间和一些自定义字段 AssemblyFsCardForUpdateIssueWorkHour
func AssemblyCardForUpdateIssueOtherColumn(headers map[string]lc_table.LcCommonField, changeItem *bo.TrendChangeListBo, cardMeta *pushV1.TemplateCard) {
	curColumn, ok := headers[changeItem.Field]
	if !ok {
		//log.Infof("[AssemblyCardForUpdateIssueOtherColumn] could not match column. columnKey: %s", changeItem.Field)
		return
	}
	displayName := curColumn.AliasTitle
	if displayName == "" {
		displayName = curColumn.Label
	}
	if str.RuneLen(displayName) > consts.ColumnAliasLimit {
		displayName = str.Substr(displayName, 0, consts.ColumnAliasLimit) + consts.CardIssueColumnNameOverflow
	}

	oldValStr := consts.CardDefaultOwnerNameForUpdateIssue
	newValStr := consts.CardDefaultOwnerNameForUpdateIssue
	if strings.Trim(changeItem.OldValue, " \n\t") != "" {
		oldValStr = changeItem.OldValue
	}
	if strings.Trim(changeItem.NewValue, " \n\t") != "" {
		newValStr = changeItem.NewValue
	}

	// 时间去除默认的空时间
	if changeItem.FieldType == consts.LcColumnFieldTypeDatepicker {
		if oldValStr == consts.BlankTime || oldValStr == consts.BlankEmptyTime {
			oldValStr = consts.CardDefaultOwnerNameForUpdateIssue
		}
		if newValStr == consts.BlankTime || newValStr == consts.BlankEmptyTime {
			newValStr = consts.CardDefaultOwnerNameForUpdateIssue
		}
	}
	// 富文本去除标签
	if changeItem.FieldType == consts.LcColumnFieldTypeRichText {
		oldValStr = str.TrimHtml(oldValStr)
		newValStr = str.TrimHtml(newValStr)
		// 去重空格和空行
		oldValStr = str.ReplaceWhiteSpaceCharToSpace(oldValStr)
		newValStr = str.ReplaceWhiteSpaceCharToSpace(newValStr)
		if strings.Trim(oldValStr, " \n\t") == "" {
			oldValStr = consts.CardIssueChangeDescForNotSupportShow
		}
		if strings.Trim(newValStr, " \n\t") == "" {
			newValStr = consts.CardIssueChangeDescForNotSupportShow
		}
	}

	markdownTile := RenderChangeDescForUpdateIssueInCard(oldValStr, newValStr)
	cardMeta.Divs = append(cardMeta.Divs, &pushV1.CardTextDiv{
		Fields: []*pushV1.CardTextField{
			&pushV1.CardTextField{
				Key:   displayName,
				Value: markdownTile,
			},
		},
	})
	//cardElements.MarkdownElements = append(cardElements.MarkdownElements, markdownTile)
	//cardElements.FsElements = append(cardElements.FsElements, vo.CardElementContentModule{
	//	Tag: "div",
	//	Fields: []vo.CardElementField{
	//		{
	//			Text: vo.CardElementText{
	//				Tag:     "lark_md",
	//				Content: markdownTile,
	//			},
	//		},
	//	},
	//})
}
