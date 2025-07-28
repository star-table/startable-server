package card

import (
	"fmt"
	"strconv"
	"strings"

	pushPb "gitea.bjx.cloud/LessCode/interface/golang/push/v1"
	sdk_const "gitea.bjx.cloud/allstar/platform-sdk/consts"
	"github.com/star-table/startable-server/common/core/config"
	"github.com/star-table/startable-server/common/core/consts"
	"github.com/star-table/startable-server/common/core/errs"
	"github.com/star-table/startable-server/common/core/util/json"
	"github.com/star-table/startable-server/common/core/util/str"
	"github.com/star-table/startable-server/common/extra/feishu"
	"github.com/star-table/startable-server/common/model/bo"
	"github.com/star-table/startable-server/common/model/vo"
	"github.com/star-table/startable-server/common/model/vo/commonvo"
	"github.com/star-table/startable-server/common/model/vo/projectvo"
	"github.com/spf13/cast"
	"google.golang.org/protobuf/types/known/structpb"
)

// GetCardUrgeIssue 飞书催办-催办负责人，卡片信息拼凑
func GetCardUrgeIssue(sourceChannel string, card *projectvo.UrgeIssueCard) (*pushPb.TemplateCard, errs.SystemErrorInfo) {
	log.Infof("[GetCardUrgeIssue] UrgeIssueCard: %s", json.ToJsonIgnoreError(card.IssueBo))
	issue := card.IssueBo

	cardDivs := []*pushPb.CardTextDiv{}
	// 催办
	cardDivs = append(cardDivs, &pushPb.CardTextDiv{
		Fields: []*pushPb.CardTextField{
			&pushPb.CardTextField{
				Key:   fmt.Sprintf(consts.CardUrgePeople, card.OperateUserName),
				Value: fmt.Sprintf(consts.CardDoubleQuotationMarks, card.UrgeText),
			},
		},
	})

	cardActionModules := []*pushPb.CardActionModule{}

	if sourceChannel == sdk_const.SourceChannelFeishu {
		// 你的回复
		cardDivs = append(cardDivs, &pushPb.CardTextDiv{
			Fields: []*pushPb.CardTextField{
				&pushPb.CardTextField{
					Key:   consts.CardElementReply,
					Value: "",
				},
			},
		})
		// select框：
		values := map[string]interface{}{
			consts.FsCardValueIssueId:  card.IssueId,
			consts.FsCardValueCardType: consts.FsCardTypeIssueUrge,
			consts.FsCardValueOrgId:    card.OrgId,
			consts.FsCardValueUserId:   card.OperateUserId,
			//consts.FsCardValueUrgeIssueBeUrgedUserId: beUrgedUserId,
			consts.FsCardValueUrgeIssueUrgeText: card.UrgeText,
		}
		newStruct, err := structpb.NewStruct(values)
		if err != nil {
			return nil, errs.BuildSystemErrorInfo(errs.TypeConvertError, err)
		}
		cardActionModules = []*pushPb.CardActionModule{
			&pushPb.CardActionModule{
				Type: pushPb.CardActionModuleType_Action,
				Actions: []*pushPb.CardAction{
					&pushPb.CardAction{
						Type:    pushPb.CardActionType_Select,
						Text:    consts.CardReplySelectFirst,
						Values:  newStruct,
						Options: consts.GetUrgeMsgOptions(),
					},
				},
			},
		}
	}

	cardMsg := GenerateCard(&commonvo.GenerateCard{
		OrgId:             card.OrgId,
		CardTitle:         fmt.Sprintf(consts.CardTitleUrgeIssue, card.OperateUserName),
		CardLevel:         pushPb.CardElementLevel_Important,
		IssueTitle:        issue.Title,
		ParentId:          issue.ParentId,
		ProjectId:         issue.ProjectId,
		ProjectTypeId:     card.ProjectTypeId,
		ParentTitle:       card.ParentTitle,
		ProjectName:       card.ProjectName,
		TableName:         card.TableName,
		OwnerInfos:        card.OwnerInfos,
		Links:             card.IssueLinks,
		IsNeedUnsubscribe: false,
		SpecialDiv:        cardDivs,
		SpecialAction:     cardActionModules,
		TableColumnInfo:   card.TableColumn,
	})

	return cardMsg, nil
}

// 被催者点击卡片交互后，给催办人回复卡片
func GetCardReplyUrgeIssue(card projectvo.UrgeReplyCard) (*pushPb.TemplateCard, errs.SystemErrorInfo) {
	issue := card.Issue

	// 催办
	cardDivs := []*pushPb.CardTextDiv{
		&pushPb.CardTextDiv{
			Fields: []*pushPb.CardTextField{
				&pushPb.CardTextField{
					Key:   fmt.Sprintf(consts.CardUrgePeople, card.OperateUserName),
					Value: fmt.Sprintf(consts.CardDoubleQuotationMarks, card.UrgeText),
				},
			},
		},
	}
	// 回复
	cardDivs = append(cardDivs, &pushPb.CardTextDiv{
		Fields: []*pushPb.CardTextField{
			&pushPb.CardTextField{
				Key:   fmt.Sprintf(consts.CardReplyBeUrged, card.BeUrgedUserName),
				Value: fmt.Sprintf(consts.CardDoubleQuotationMarks, card.ReplyMsg),
			},
		},
	})

	cardMsg := GenerateCard(&commonvo.GenerateCard{
		OrgId:             card.OrgId,
		CardTitle:         fmt.Sprintf(consts.CardReplyUrgeTitle, card.BeUrgedUserName),
		CardLevel:         pushPb.CardElementLevel_Default,
		IssueTitle:        issue.Title,
		ParentId:          issue.ParentID,
		ProjectId:         issue.ProjectID,
		ProjectTypeId:     card.ProjectTypeId,
		ParentTitle:       card.ParentTitle,
		ProjectName:       card.ProjectName,
		TableName:         card.TableName,
		OwnerInfos:        card.OwnerInfos,
		Links:             card.IssueLinks,
		IsNeedUnsubscribe: false,
		SpecialDiv:        cardDivs,
		TableColumnInfo:   card.TableColumn,
	})

	return cardMsg, nil
}

// GetCardAuditorWhenUpdateIssue 催办确认人卡片-修改任务状态时，如果状态为待确认，则向确认人推送卡片
func GetCardAuditorWhenUpdateIssue(infoObj *projectvo.BaseInfoBoxForIssueCard) (*pushPb.TemplateCard, errs.SystemErrorInfo) {
	if infoObj.IssueInfo.Title == "" {
		return nil, errs.CardTitleEmpty
	}

	cardMsg := &pushPb.TemplateCard{}
	cardMsg.Title = &pushPb.CardTitle{
		Level: pushPb.CardElementLevel_Default,
		Title: fmt.Sprintf(consts.CardAuditorTitle, infoObj.OperateUser.Name),
	}

	contentStr := ""
	if infoObj.IssueInfo.ProjectId == 0 {
		contentStr = consts.CardDefaultIssueProjectName
	} else {
		contentStr = fmt.Sprintf(consts.CardTablePro, infoObj.ProjectAuthBo.Name, infoObj.IssueTableInfo.Name)
	}

	// 父记录
	parentTitle := consts.CardDefaultIssueParentTitleForUpdateIssue
	if infoObj.IssueInfo.ParentId != 0 {
		if infoObj.ParentIssue.Title != "" {
			parentTitle = infoObj.ParentIssue.Title
		}
	}

	// 负责人
	ownerSlice := []string{}
	for _, ownerInfo := range infoObj.OwnerInfos {
		ownerSlice = append(ownerSlice, ownerInfo.Name)
	}
	ownerDisplayName := strings.Join(ownerSlice, "，")
	if ownerDisplayName == "" {
		ownerDisplayName = consts.CardDefaultOwnerNameForUpdateIssue
	}

	columnTitle, columnOwnerId := GetCardTitleAndOwnerColumn(infoObj.ProjectTableColumn)

	cardMsg.Divs = append(cardMsg.Divs, &pushPb.CardTextDiv{
		Fields: []*pushPb.CardTextField{
			&pushPb.CardTextField{
				Key:   columnTitle,
				Value: infoObj.IssueInfo.Title,
			},
		},
	})
	cardMsg.Divs = append(cardMsg.Divs, &pushPb.CardTextDiv{
		Fields: []*pushPb.CardTextField{
			&pushPb.CardTextField{
				Key:   consts.CardElementProjectTable,
				Value: contentStr,
			},
		},
	})
	if infoObj.IssueInfo.ParentId != 0 {
		cardMsg.Divs = append(cardMsg.Divs, &pushPb.CardTextDiv{
			Fields: []*pushPb.CardTextField{
				&pushPb.CardTextField{
					Key:   consts.CardElementParent,
					Value: parentTitle,
				},
			},
		})
	}
	cardMsg.Divs = append(cardMsg.Divs, &pushPb.CardTextDiv{
		Fields: []*pushPb.CardTextField{
			&pushPb.CardTextField{
				Key:   columnOwnerId,
				Value: ownerDisplayName,
			},
		},
	})

	cardActions := []*pushPb.CardAction{}
	if infoObj.SourceChannel == sdk_const.SourceChannelFeishu {
		mPass := map[string]interface{}{
			consts.FsCardValueCardType:            consts.FsCardTypeIssueAudit,
			consts.FsCardValueOrgId:               infoObj.IssueInfo.OrgId,
			consts.FsCardValueUserId:              infoObj.OperateUserId,
			consts.FsCardValueAppId:               strconv.FormatInt(infoObj.ProjectAuthBo.AppId, 10),
			consts.FsCardValueTableId:             strconv.FormatInt(infoObj.IssueInfo.TableId, 10),
			consts.FsCardValueIssueId:             infoObj.IssueInfo.Id,
			consts.FsCardValueAuditIssueResult:    consts.AuditStatusPass,
			consts.FsCardValueAuditIssueStarterId: infoObj.OperateUserId,
		}
		passValue, err := structpb.NewStruct(mPass)
		if err != nil {
			log.Errorf("[GetCardAuditorWhenUpdateIssue] err:%v", err)
			return nil, errs.BuildSystemErrorInfo(errs.TypeConvertError, err)
		}

		mReject := map[string]interface{}{
			consts.FsCardValueCardType:            consts.FsCardTypeIssueAudit,
			consts.FsCardValueOrgId:               infoObj.IssueInfo.OrgId,
			consts.FsCardValueUserId:              infoObj.OperateUserId,
			consts.FsCardValueAppId:               strconv.FormatInt(infoObj.ProjectAuthBo.AppId, 10),
			consts.FsCardValueTableId:             strconv.FormatInt(infoObj.IssueInfo.TableId, 10),
			consts.FsCardValueIssueId:             infoObj.IssueInfo.Id,
			consts.FsCardValueAuditIssueResult:    consts.AuditStatusReject,
			consts.FsCardValueAuditIssueStarterId: infoObj.OperateUserId,
		}
		rejectValue, err := structpb.NewStruct(mReject)
		if err != nil {
			log.Errorf("[GetCardAuditorWhenUpdateIssue] err:%v", err)
			return nil, errs.BuildSystemErrorInfo(errs.TypeConvertError, err)
		}

		cardActions = []*pushPb.CardAction{
			&pushPb.CardAction{
				Level:  pushPb.CardElementLevel_Default,
				Type:   pushPb.CardActionType_ButtonCallback,
				Text:   consts.CardAuditPass,
				Values: passValue,
			},
			&pushPb.CardAction{
				Level:  pushPb.CardElementLevel_Important,
				Type:   pushPb.CardActionType_ButtonCallback,
				Text:   consts.CardAuditReject,
				Values: rejectValue,
			},
		}
	}

	// 查看详情按钮
	cardActions = append(cardActions, &pushPb.CardAction{
		Level: pushPb.CardElementLevel_Default,
		Type:  pushPb.CardActionType_ButtonJumpUrl,
		Text:  consts.CardButtonTextForViewDetail,
		Url:   infoObj.IssueInfoUrl,
	})
	cardActions = append(cardActions, &pushPb.CardAction{
		Level: pushPb.CardElementLevel_NotImportant,
		Type:  pushPb.CardActionType_ButtonJumpUrl,
		Text:  consts.CardButtonTextForViewInsideApp,
		Url:   infoObj.IssuePcUrl,
	})

	cardMsg.ActionModules = append(cardMsg.ActionModules, &pushPb.CardActionModule{
		Type:    pushPb.CardActionModuleType_Action,
		Actions: cardActions,
	})

	return cardMsg, nil
}

// 添加成员的卡片
func GetCardIssueAddMember(issueTrendsBo *bo.IssueTrendsBo, infoObj *projectvo.BaseInfoBoxForIssueCard, columnDisplayName string) *pushPb.TemplateCard {
	meta := GenerateCard(&commonvo.GenerateCard{
		OrgId:         issueTrendsBo.OrgId,
		CardTitle:     fmt.Sprintf(consts.CardIssueAddMember, GetOperationName(issueTrendsBo, infoObj.OperateUser.Name), columnDisplayName),
		CardLevel:     pushPb.CardElementLevel_Default,
		IssueTitle:    infoObj.IssueInfo.Title,
		ParentId:      infoObj.IssueInfo.ParentId,
		ProjectId:     infoObj.IssueInfo.ProjectId,
		ProjectTypeId: infoObj.ProjectAuthBo.ProjectTypeId,
		ProjectName:   infoObj.ProjectAuthBo.Name,
		TableName:     infoObj.IssueTableInfo.Name,
		ParentTitle:   infoObj.ParentIssue.Title,
		OwnerInfos:    infoObj.OwnerInfos,
		Links: &projectvo.IssueLinks{
			Link:        infoObj.IssuePcUrl,
			SideBarLink: infoObj.IssueInfoUrl,
		},
		IsNeedUnsubscribe: true,
		TableColumnInfo:   infoObj.ProjectTableColumn,
	})

	return meta
}

// GetOperationName 获取操作人的名字
func GetOperationName(issueTrendsBo *bo.IssueTrendsBo, name string) string {
	automation := issueTrendsBo.Ext.AutomationInfo
	if automation != nil && automation.WorkflowName != "" {
		if automation.TriggerUser != nil {
			return fmt.Sprintf(consts.CardDataEvent, automation.TriggerUser.Name, automation.WorkflowName)
		}

		return fmt.Sprintf(consts.CardCronTrigger, automation.WorkflowName)
	}

	return name
}

// 即将逾期的卡片 飞书和其他平台不同
func GetCardBeOverdue(card *projectvo.CardBeOverdue) *pushPb.TemplateCard {

	cardMsg := &pushPb.TemplateCard{}
	cardMsg.Title = &pushPb.CardTitle{
		Level: pushPb.CardElementLevel_Important,
		Title: consts.CardTitlePlanEndRemind,
	}

	columnTitle, columnOwnerId := GetCardTitleAndOwnerColumn(card.TableColumn)

	cardMsg.Divs = append(cardMsg.Divs, &pushPb.CardTextDiv{
		Fields: []*pushPb.CardTextField{
			&pushPb.CardTextField{
				Key:   columnTitle,
				Value: card.IssueTitle,
			},
		},
	})
	cardMsg.Divs = append(cardMsg.Divs, &pushPb.CardTextDiv{
		Fields: []*pushPb.CardTextField{
			&pushPb.CardTextField{
				Key:   consts.CardElementProjectTable,
				Value: card.ContentProject,
			},
		},
	})

	if card.ParentTitle != consts.CardDefaultRelationIssueTitle {
		cardMsg.Divs = append(cardMsg.Divs, &pushPb.CardTextDiv{
			Fields: []*pushPb.CardTextField{
				&pushPb.CardTextField{
					Key:   consts.CardElementParent,
					Value: card.ParentTitle,
				},
			},
		})
	}
	cardMsg.Divs = append(cardMsg.Divs, &pushPb.CardTextDiv{
		Fields: []*pushPb.CardTextField{
			&pushPb.CardTextField{
				Key:   columnOwnerId,
				Value: card.OwnerNameStr,
			},
		},
	})

	if card.SourceChannel == sdk_const.SourceChannelFeishu {
		cardMsg.Divs = append(cardMsg.Divs, &pushPb.CardTextDiv{
			Fields: []*pushPb.CardTextField{
				&pushPb.CardTextField{
					Key: consts.CardPlanEndTime,
				},
			},
		})
		m := map[string]interface{}{
			consts.FsCardValueIssueId:   card.IssueId,
			consts.FsCardValueCardType:  consts.FsCardTypeIssueRemind,
			consts.FsCardValueAction:    consts.FsCardActionUpdatePlanEndTime,
			consts.FsCardValueOrgId:     card.OrgId,
			consts.FsCardValueUserId:    card.OwnerId,
			consts.FsCardValueTableId:   card.TableId,
			consts.FsCardValueAppId:     card.AppId,
			consts.FsCardDatePickerInit: card.PlanEndTime,
		}
		values, err := structpb.NewStruct(m)
		if err != nil {
			log.Errorf("[GetCardBeOverdue] err:%v", err)
			return nil
		}
		cardMsg.ActionModules = append(cardMsg.ActionModules, &pushPb.CardActionModule{
			Type: pushPb.CardActionModuleType_Action,
			Actions: []*pushPb.CardAction{
				&pushPb.CardAction{
					Level:  pushPb.CardElementLevel_Default,
					Type:   pushPb.CardActionType_DatePicker,
					Text:   consts.CardModifyDeadline,
					Values: values,
				},
			},
		})
		m2 := map[string]interface{}{
			consts.FsCardValueIssueId:    card.IssueId,
			consts.FsCardValueCardType:   consts.FsCardTypeIssueRemind,
			consts.FsCardValueOrgId:      card.OrgId,
			consts.FsCardValueUserId:     card.OwnerId,
			consts.FsCardValueAppId:      card.AppId,
			consts.FsCardValueTableId:    card.TableId,
			consts.FsCardValueStatusType: consts.StatusTypeComplete,
			consts.FsCardValueAction:     consts.FsCardActionUpdateIssueStatus,
		}
		values2, err := structpb.NewStruct(m2)
		if err != nil {
			log.Errorf("[GetCardBeOverdue] err:%v", err)
			return nil
		}
		cardMsg.ActionModules = append(cardMsg.ActionModules, &pushPb.CardActionModule{
			Type: pushPb.CardActionModuleType_Action,
			Actions: []*pushPb.CardAction{
				&pushPb.CardAction{
					Level:  pushPb.CardElementLevel_Default,
					Type:   pushPb.CardActionType_ButtonCallback,
					Text:   consts.CardCompleteTask,
					Values: values2,
				},
				&pushPb.CardAction{
					Level: pushPb.CardElementLevel_Default,
					Type:  pushPb.CardActionType_ButtonJumpUrl,
					Text:  consts.CardButtonTextForViewDetail,
					Url:   card.IssueLinks.SideBarLink,
				},
				&pushPb.CardAction{
					Level: pushPb.CardElementLevel_NotImportant,
					Type:  pushPb.CardActionType_ButtonJumpUrl,
					Text:  consts.CardButtonTextForViewInsideApp,
					Url:   card.IssueLinks.Link,
				},
			},
		})
	} else {
		cardMsg.Divs = append(cardMsg.Divs, &pushPb.CardTextDiv{
			Fields: []*pushPb.CardTextField{
				&pushPb.CardTextField{
					Key:   consts.CardPlanEndTime,
					Value: card.PlanEndTime,
				},
			},
		})
		cardMsg.ActionModules = append(cardMsg.ActionModules, &pushPb.CardActionModule{
			Type: pushPb.CardActionModuleType_Action,
			Actions: []*pushPb.CardAction{
				&pushPb.CardAction{
					Level: pushPb.CardElementLevel_Default,
					Type:  pushPb.CardActionType_ButtonJumpUrl,
					Text:  consts.CardButtonTextForViewDetail,
					Url:   card.IssueLinks.SideBarLink,
				},
				&pushPb.CardAction{
					Level: pushPb.CardElementLevel_NotImportant,
					Type:  pushPb.CardActionType_ButtonJumpUrl,
					Text:  consts.CardButtonTextForViewInsideApp,
					Url:   card.IssueLinks.Link,
				},
			},
		})
	}

	if !card.HasPermission {
		cardMsg.Divs = append(cardMsg.Divs, &pushPb.CardTextDiv{
			Fields: []*pushPb.CardTextField{
				&pushPb.CardTextField{
					Key:   card.Tips,
					Value: "",
				},
			},
		})
	}

	return cardMsg
}

// GetCardImportIssue 导入/应用模板时提示创建任务的卡片
func GetCardImportIssue(infoBox *projectvo.AsyncTaskCard, issueNum int) *pushPb.TemplateCard {
	cardMsg := &pushPb.TemplateCard{}
	cardMsg.Title = &pushPb.CardTitle{
		Level: pushPb.CardElementLevel_Default,
		Title: fmt.Sprintf("%s批量创建了%d条新记录", infoBox.User.Name, issueNum),
	}

	titleKey := consts.CardElementProjectTable
	titleValue := infoBox.Project.Name + "/" + infoBox.Table.Name
	if infoBox.IsApplyTemplate {
		titleKey = consts.CardElementProject
		titleValue = infoBox.Project.Name
	}
	cardMsg.Divs = append(cardMsg.Divs, &pushPb.CardTextDiv{
		Fields: []*pushPb.CardTextField{
			&pushPb.CardTextField{
				Key:   titleKey,
				Value: titleValue,
			},
		},
	})
	cardMsg.ActionModules = append(cardMsg.ActionModules, &pushPb.CardActionModule{
		Type: pushPb.CardActionModuleType_Action,
		Actions: []*pushPb.CardAction{
			&pushPb.CardAction{
				Level: pushPb.CardElementLevel_Default,
				Type:  pushPb.CardActionType_ButtonJumpUrl,
				Text:  consts.CardButtonTextForViewDetail,
				Url:   infoBox.Link,
			},
		},
	})
	return cardMsg
	//meta := &commonvo.CardMeta{}
	//meta.IsWide = false
	//meta.Level = consts.CardLevelInfo
	//meta.Title = fmt.Sprintf("%s批量创建了%d条新记录", infoBox.User.Name, issueNum)
	//
	//eleTitle := consts.CardElementProjectTable
	//eleContent := fmt.Sprintf(consts.CardTablePro, infoBox.Project.Name, infoBox.Table.Name)
	//if infoBox.IsApplyTemplate {
	//	eleTitle = consts.CardElementProject
	//	eleContent = infoBox.Project.Name
	//}
	//
	//meta.Divs = append(meta.Divs, &commonvo.CardDiv{
	//	Fields: []*commonvo.CardField{
	//		&commonvo.CardField{
	//			Key:   eleTitle,
	//			Value: eleContent,
	//		},
	//	},
	//})
	//
	//meta.ActionMarkdowns = []string{fmt.Sprintf(consts.MarkdownLink, consts.CardButtonTextForViewDetail, infoBox.Link)}
	//meta.FsActionElements = []interface{}{
	//	fsSdkVo.CardElementActionModule{
	//		Tag:    "action",
	//		Layout: "bisected",
	//		Actions: []interface{}{
	//			fsSdkVo.ActionButton{
	//				Tag: "button",
	//				Text: fsSdkVo.CardElementText{
	//					Tag:     "plain_text",
	//					Content: consts.CardButtonTextForViewDetail,
	//				},
	//				Url:  infoBox.Link,
	//				Type: consts.FsCardButtonColorPrimary,
	//			},
	//		},
	//	},
	//}
	//return meta
}

// ReplyAuditorNotNeedAudit 向审批者发送无需审批的卡片，因为任务状态不符合审批条件
func ReplyAuditorNotNeedAudit(projectName, tableName string, ownerInfos []*bo.BaseUserInfoBo, issue vo.Issue,
	issueLinks *projectvo.IssueLinks, tableColumnMap map[string]*projectvo.TableColumnData) *pushPb.TemplateCard {

	cardMsg := GenerateCard(&commonvo.GenerateCard{
		OrgId:             issue.OrgID,
		CardTitle:         consts.CardTitleAuditAlready,
		CardLevel:         pushPb.CardElementLevel_Default,
		ProjectName:       projectName,
		TableName:         tableName,
		IssueTitle:        issue.Title,
		ProjectId:         issue.ProjectID,
		OwnerInfos:        ownerInfos,
		Links:             issueLinks,
		IsNeedUnsubscribe: false,
		TableColumnInfo:   tableColumnMap,
	})

	return cardMsg
}

// ReplyIssueAuditCard 审批通过卡片返回
func ReplyIssueAuditCard(projectName, tableName string, ownerInfos []*bo.BaseUserInfoBo, issue vo.Issue, auditStatus int,
	issueLinks *projectvo.IssueLinks, tableColumnMap map[string]*projectvo.TableColumnData) *pushPb.TemplateCard {
	commonCard := &commonvo.GenerateCard{
		OrgId:             issue.OrgID,
		ProjectName:       projectName,
		TableName:         tableName,
		IssueTitle:        issue.Title,
		ProjectId:         issue.ProjectID,
		OwnerInfos:        ownerInfos,
		Links:             issueLinks,
		IsNeedUnsubscribe: false,
		TableColumnInfo:   tableColumnMap,
	}
	titlePart := ""
	cardTitleColor := pushPb.CardElementLevel_Default
	switch auditStatus {
	case consts.AuditStatusPass:
		titlePart = "✅ 你已通过了审批"
	case consts.AuditStatusReject:
		titlePart = "你已驳回了审批"
		cardTitleColor = pushPb.CardElementLevel_Important
	}
	commonCard.CardTitle = titlePart
	commonCard.CardLevel = cardTitleColor
	cardMsg := GenerateCard(commonCard)
	return cardMsg
}

// GetFsCardForUpgradeNotifyAdmin 获取 fs卡片，通知组织管理员，有成员要升级极星版本
func GetFsCardForUpgradeNotifyAdmin(opUserOpenId, opUserName string) *pushPb.TemplateCard {
	atSomeoneStr := str.RenderAtSomeoneStr(opUserOpenId, opUserName)
	// targetProVersionName := businees.GetOrgPayVersionName(targetPayLevel)
	cardMsg := &pushPb.TemplateCard{}
	cardMsg.Title = &pushPb.CardTitle{
		Level: pushPb.CardElementLevel_Default,
		Title: consts.FsCardUpgradeNoticeTitle,
	}
	cardMsg.Divs = append(cardMsg.Divs, &pushPb.CardTextDiv{
		Fields: []*pushPb.CardTextField{
			&pushPb.CardTextField{
				Key:   fmt.Sprintf(consts.FsCardUpgradeNoticeContent, atSomeoneStr),
				Value: "",
			},
		},
	})
	cardMsg.ActionModules = append(cardMsg.ActionModules, &pushPb.CardActionModule{
		Type: pushPb.CardActionModuleType_Action,
		Actions: []*pushPb.CardAction{
			&pushPb.CardAction{
				Level: pushPb.CardElementLevel_Default,
				Type:  pushPb.CardActionType_ButtonJumpUrl,
				Text:  consts.FsCardButtonTextForGotoSetting,
				Url:   fmt.Sprintf(feishu.FsLinkForSettingPolarisApp, config.GetConfig().FeiShu.AppId),
			},
		},
	})

	return cardMsg
}

// ChatCreateIssueFailNotice 群聊中创建任务失败时，推送卡片以提示
func ChatCreateIssueFailNotice(errMsg string) *pushPb.TemplateCard {
	card := &pushPb.TemplateCard{}
	card.Title = &pushPb.CardTitle{
		Level: pushPb.CardElementLevel_Default,
		Title: consts.CardIssueCreateErr,
	}

	card.Divs = append(card.Divs, &pushPb.CardTextDiv{
		Fields: []*pushPb.CardTextField{
			&pushPb.CardTextField{
				Key:   consts.CardIssueCreateErr,
				Value: "",
			},
		},
	})
	if errMsg == consts.FsCardTextIssueMaxLimit {
		// 如果是极星不是企业版导致的任务数限制导致的问题，则卡片显示特定的按钮
		card.ActionModules = append(card.ActionModules, &pushPb.CardActionModule{
			Type: pushPb.CardActionModuleType_Action,
			Actions: []*pushPb.CardAction{
				&pushPb.CardAction{
					Level: pushPb.CardElementLevel_Default,
					Type:  pushPb.CardActionType_ButtonJumpUrl,
					Text:  consts.FsCardButtonTextForGotoSetting,
					Url:   fmt.Sprintf(feishu.FsLinkForSettingPolarisApp, config.GetConfig().FeiShu.AppId),
				},
			},
		})
	}
	return card
}

// 分享任务的卡片消息
func GetIssueShareCardMsg(cardTitle, issueTitle, issueCode, projectName, tableName, ownerNameStr, issueInfoUrl, issuePcUrl string) *pushPb.TemplateCard {
	cardMsg := &pushPb.TemplateCard{}
	cardMsg.Title = &pushPb.CardTitle{
		Title: cardTitle,
	}
	cardMsg.Divs = append(cardMsg.Divs, &pushPb.CardTextDiv{
		Fields: []*pushPb.CardTextField{
			&pushPb.CardTextField{
				Key:   consts.CardElementTitle,
				Value: issueTitle,
			},
		},
	})
	cardMsg.Divs = append(cardMsg.Divs, &pushPb.CardTextDiv{
		Fields: []*pushPb.CardTextField{
			&pushPb.CardTextField{
				Key:   consts.CardElementCode,
				Value: issueCode,
			},
		},
	})
	contentStr := ""
	if projectName == consts.CardDefaultIssueProjectName {
		contentStr = consts.CardDefaultIssueProjectName
	} else {
		contentStr = fmt.Sprintf(consts.CardTablePro, projectName, tableName)
	}
	cardMsg.Divs = append(cardMsg.Divs, &pushPb.CardTextDiv{
		Fields: []*pushPb.CardTextField{
			&pushPb.CardTextField{
				Key:   consts.CardElementProjectTable,
				Value: contentStr,
			},
		},
	})
	cardMsg.Divs = append(cardMsg.Divs, &pushPb.CardTextDiv{
		Fields: []*pushPb.CardTextField{
			&pushPb.CardTextField{
				Key:   consts.CardElementOwner,
				Value: ownerNameStr,
			},
		},
	})
	cardAction := []*pushPb.CardAction{}
	cardAction = append(cardAction, &pushPb.CardAction{
		Level: pushPb.CardElementLevel_Default,
		Text:  consts.CardButtonTextForViewDetail,
		Url:   issueInfoUrl,
	})
	cardAction = append(cardAction, &pushPb.CardAction{
		Level: pushPb.CardElementLevel_NotImportant,
		Text:  consts.CardButtonTextForViewInsideApp,
		Url:   issuePcUrl,
	})

	cardMsg.ActionModules = []*pushPb.CardActionModule{
		&pushPb.CardActionModule{
			Type:    pushPb.CardActionModuleType_Action,
			Actions: cardAction,
		},
	}
	return cardMsg
}

// 个人日报
func GetDailyIssueReportCard(orgId, overdueSum, dueOfTodaySum, beAboutToOverdueSum, pendingSum int64, link string) *pushPb.TemplateCard {
	cardMsg := &pushPb.TemplateCard{}
	cardMsg.Title = &pushPb.CardTitle{
		Level: pushPb.CardElementLevel_Default,
		Title: consts.DailyPersonalIssueReportTitle,
	}
	cardMsg.Divs = append(cardMsg.Divs, &pushPb.CardTextDiv{
		Fields: []*pushPb.CardTextField{
			&pushPb.CardTextField{
				Key:   consts.DailyOverdueCountTitle,
				Value: strconv.FormatInt(overdueSum, 10),
			},
			&pushPb.CardTextField{
				Key:   consts.DailyDueOfTodayTitle,
				Value: strconv.FormatInt(dueOfTodaySum, 10),
			},
			&pushPb.CardTextField{
				Key:   consts.BeAboutToOverdueTitle,
				Value: strconv.FormatInt(beAboutToOverdueSum, 10),
			},
			&pushPb.CardTextField{
				Key:   consts.PendingTitle,
				Value: strconv.FormatInt(pendingSum, 10),
			},
		},
	})

	cardMsg.ActionModules = append(cardMsg.ActionModules, &pushPb.CardActionModule{
		Type: pushPb.CardActionModuleType_Action,
		Actions: []*pushPb.CardAction{
			&pushPb.CardAction{
				Type: pushPb.CardActionType_ButtonJumpUrl,
				Text: consts.CardButtonTextForPersonalReport,
				Url:  link,
			},
		},
	})
	cardMsg.ActionModules = append(cardMsg.ActionModules, &pushPb.CardActionModule{
		Type: pushPb.CardActionModuleType_Hr,
	})
	m := map[string]interface{}{
		consts.FsCardValueCardType: consts.FsCardTypeNoticeSubscribe,
		consts.FsCardValueAction:   consts.FsCardUnsubscribePersonalReport,
		consts.FsCardValueOrgId:    orgId,
	}
	values, err := structpb.NewStruct(m)
	if err != nil {
		log.Errorf("[GetDailyIssueReportCard] err:%v", err)
		return nil
	}
	cardMsg.ActionModules = append(cardMsg.ActionModules, &pushPb.CardActionModule{
		Type: pushPb.CardActionModuleType_Action,
		Actions: []*pushPb.CardAction{
			&pushPb.CardAction{
				Level:  pushPb.CardElementLevel_Default,
				Type:   pushPb.CardActionType_Link,
				Text:   consts.CardUnSubscribe,
				Values: values,
			},
		},
	})

	return cardMsg
}

// 项目日报
func GetDailyProjectCard(msg *bo.DailyProjectReportMsgBo, link string) *pushPb.TemplateCard {
	cardMsg := &pushPb.TemplateCard{}
	cardMsg.Title = &pushPb.CardTitle{
		Level: pushPb.CardElementLevel_Default,
		Title: consts.DailyProjectReportTitle,
	}

	cardMsg.Divs = append(cardMsg.Divs, &pushPb.CardTextDiv{
		Fields: []*pushPb.CardTextField{
			&pushPb.CardTextField{
				Key:   consts.CardElementProject,
				Value: msg.ProjectName,
			},
			&pushPb.CardTextField{
				Key:   consts.DailyFinishCountTitle,
				Value: cast.ToString(msg.DailyFinishCount),
			},
			&pushPb.CardTextField{
				Key:   consts.DailyRemainingCountTitle,
				Value: cast.ToString(msg.RemainingCount),
			},
			&pushPb.CardTextField{
				Key:   consts.DailyOverdueCountTitle,
				Value: cast.ToString(msg.OverdueCount),
			},
		},
	})

	cardMsg.ActionModules = append(cardMsg.ActionModules, &pushPb.CardActionModule{
		Type: pushPb.CardActionModuleType_Action,
		Actions: []*pushPb.CardAction{
			&pushPb.CardAction{
				Type: pushPb.CardActionType_ButtonJumpUrl,
				Text: consts.CardButtonTextForViewDetail,
				Url:  link,
			},
		},
	})
	cardMsg.ActionModules = append(cardMsg.ActionModules, &pushPb.CardActionModule{
		Type: pushPb.CardActionModuleType_Hr,
	})
	m := map[string]interface{}{
		consts.FsCardValueCardType: consts.FsCardTypeNoticeSubscribe,
		consts.FsCardValueAction:   consts.FsCardUnsubscribeProjectReport,
		consts.FsCardValueOrgId:    msg.OrgId,
	}
	values, err := structpb.NewStruct(m)
	if err != nil {
		log.Errorf("[GetDailyProjectCard] err:%v", err)
		return nil
	}
	cardMsg.ActionModules = append(cardMsg.ActionModules, &pushPb.CardActionModule{
		Type: pushPb.CardActionModuleType_Action,
		Actions: []*pushPb.CardAction{
			&pushPb.CardAction{
				Level:  pushPb.CardElementLevel_Default,
				Type:   pushPb.CardActionType_Link,
				Text:   consts.CardUnSubscribe,
				Values: values,
			},
		},
	})

	return cardMsg
}
