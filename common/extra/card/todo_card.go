package card

import (
	"fmt"

	automationPb "gitea.bjx.cloud/LessCode/interface/golang/automation/v1"

	pushPb "gitea.bjx.cloud/LessCode/interface/golang/push/v1"
	"github.com/star-table/startable-server/common/core/consts"
	"github.com/star-table/startable-server/common/core/errs"
	"github.com/star-table/startable-server/common/core/util/json"
	"github.com/star-table/startable-server/common/model/vo/commonvo"
	"github.com/star-table/startable-server/common/model/vo/projectvo"
)

// SendCardTodoUrge 催办待办事项
func SendCardTodoUrge(sourceChannel string, card *projectvo.UrgeIssueCard) errs.SystemErrorInfo {
	cd, err := GetCardTodoUrge(card)
	if err != nil {
		log.Errorf("[SendCardTodoUrge] GetCardTodoUrge err: %v", err)
		return err
	}
	pushCard := &commonvo.PushCard{
		OrgId:         card.OrgId,
		OutOrgId:      card.OutOrgId,
		SourceChannel: sourceChannel,
		OpenIds:       card.OpenIds,
		CardMsg:       cd,
	}

	errSys := PushCard(card.OrgId, pushCard)
	if errSys != nil {
		log.Errorf("[SendCardUrgeIssue] err:%v", errSys)
		return errSys
	}
	return nil
}

// GetCardTodoUrge 催办待办事项-卡片组装
func GetCardTodoUrge(card *projectvo.UrgeIssueCard) (*pushPb.TemplateCard, errs.SystemErrorInfo) {
	log.Infof("[GetCardUrgeIssueNew] UrgeIssueCard: %s", json.ToJsonIgnoreError(card.IssueBo))
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
		SpecialAction:     nil,
		TableColumnInfo:   card.TableColumn,
	})

	return cardMsg, nil
}

// GetCardTodo 待办事项-卡片组装
func GetCardTodo(card *projectvo.TodoCard) (*pushPb.TemplateCard, errs.SystemErrorInfo) {
	var title string
	var button string
	if card.TodoType == automationPb.TodoType_Audit {
		title = consts.CardTodoAuditTitle
		button = consts.CardButtonTextAudit
	} else if card.TodoType == automationPb.TodoType_FillIn {
		title = consts.CardTodoFillInTitle
		button = consts.CardButtonTextFillIn
	}
	cardMsg := &pushPb.TemplateCard{}
	cardMsg.Title = &pushPb.CardTitle{
		Level: pushPb.CardElementLevel_Default,
		Title: fmt.Sprintf(title, card.Trigger.Name),
	}

	contentStr := ""
	if card.Issue.ProjectId == 0 || card.Issue.TableId == 0 {
		contentStr = consts.CardDefaultIssueProjectName
	} else {
		contentStr = fmt.Sprintf(consts.CardTablePro, card.Project.Name, card.Table.Name)
	}

	columnTitle, _ := GetCardTitleAndOwnerColumn(card.TableColumns)

	cardMsg.Divs = append(cardMsg.Divs, &pushPb.CardTextDiv{
		Fields: []*pushPb.CardTextField{
			&pushPb.CardTextField{
				Key:   columnTitle,
				Value: card.Issue.Title,
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

	cardActions := []*pushPb.CardAction{}

	// 去审批/去填写
	cardActions = append(cardActions, &pushPb.CardAction{
		Level: pushPb.CardElementLevel_Default,
		Type:  pushPb.CardActionType_ButtonJumpUrl,
		Text:  button,
		Url:   card.TodoSideBarUrl,
	})
	// 应用内处理
	cardActions = append(cardActions, &pushPb.CardAction{
		Level: pushPb.CardElementLevel_NotImportant,
		Type:  pushPb.CardActionType_ButtonJumpUrl,
		Text:  consts.CardButtonTextForDoInsideApp,
		Url:   card.TodoUrl,
	})

	cardMsg.ActionModules = append(cardMsg.ActionModules, &pushPb.CardActionModule{
		Type:    pushPb.CardActionModuleType_Action,
		Actions: cardActions,
	})

	return cardMsg, nil
}
