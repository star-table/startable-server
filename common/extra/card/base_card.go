package card

import (
	"fmt"
	"strings"

	pushV1 "gitea.bjx.cloud/LessCode/interface/golang/push/v1"
	"github.com/star-table/startable-server/common/core/consts"
	"github.com/star-table/startable-server/common/core/errs"
	"github.com/star-table/startable-server/common/model/vo/commonvo"
	"github.com/star-table/startable-server/common/model/vo/projectvo"
	"google.golang.org/protobuf/types/known/structpb"
)

// 催办回复卡片发送给用户
func SendCardUrgeIssue(sourceChannel string, card *projectvo.UrgeIssueCard) errs.SystemErrorInfo {
	cd, err := GetCardUrgeIssue(sourceChannel, card)
	if err != nil {
		log.Errorf("[SendCardUrgeIssue] GetCardUrgeIssue err: %v", err)
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

	//SendCardMeta(sourceChannel, card.OutOrgId, cd, card.OpenIds)
	return nil
}

// 飞书 信息回复给催办人，被催着卡片状态信息变更
// 被催者点击卡片交互后，给催办人回复卡片
// 被催者点击卡片交互后，卡片样式变动
func SendCardBeUrgeIssue(sourceChannel string, card projectvo.UrgeReplyCard) errs.SystemErrorInfo {
	cd, err := GetCardReplyUrgeIssue(card)
	if err != nil {
		log.Errorf("[SendCardBeUrgeIssue] GetCardReplyUrgeIssue err: %v", err)
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
		log.Errorf("[SendCardBeUrgeIssue] err:%v", errSys)
		return errSys
	}
	//SendCardMeta(sourceChannel, card.OutOrgId, cd, card.OpenIds)
	return nil
}

// 生成通用的更新记录的card
func GenerateCard(card *commonvo.GenerateCard) *pushV1.TemplateCard {
	cardMsg := &pushV1.TemplateCard{}
	cardMsg.Title = &pushV1.CardTitle{
		Level: card.CardLevel,
		Title: card.CardTitle,
	}

	contentStr := ""
	if card.ProjectId == 0 {
		contentStr = consts.CardDefaultIssueProjectName
	} else {
		contentStr = fmt.Sprintf(consts.CardTablePro, card.ProjectName, card.TableName)
	}

	// 父记录
	parentTitle := consts.CardDefaultRelationIssueTitle
	if card.ParentId != 0 {
		if card.ParentTitle != "" {
			parentTitle = card.ParentTitle
		}
	}
	columnTitle, columnOwnerId := GetCardTitleAndOwnerColumn(card.TableColumnInfo)
	// 负责人
	ownerSlice := []string{}
	for _, ownerInfo := range card.OwnerInfos {
		ownerSlice = append(ownerSlice, ownerInfo.Name)
	}
	ownerDisplayName := strings.Join(ownerSlice, "，")
	if ownerDisplayName == "" {
		ownerDisplayName = consts.CardDefaultOwnerNameForUpdateIssue
	}

	cardMsg.Divs = append(cardMsg.Divs, &pushV1.CardTextDiv{
		Fields: []*pushV1.CardTextField{
			&pushV1.CardTextField{
				Key:   columnTitle,
				Value: card.IssueTitle,
			},
		},
	})

	projectTableName := consts.CardElementProjectTable
	if card.ProjectTypeId == consts.ProjectTypeEmpty {
		projectTableName = consts.CardElementAppTable
	}
	cardMsg.Divs = append(cardMsg.Divs, &pushV1.CardTextDiv{
		Fields: []*pushV1.CardTextField{
			&pushV1.CardTextField{
				Key:   projectTableName,
				Value: contentStr,
			},
		},
	})
	if card.ParentId != 0 {
		cardMsg.Divs = append(cardMsg.Divs, &pushV1.CardTextDiv{
			Fields: []*pushV1.CardTextField{
				&pushV1.CardTextField{
					Key:   consts.CardElementParent,
					Value: parentTitle,
				},
			},
		})
	}

	if card.ProjectTypeId != consts.ProjectTypeEmpty {
		cardMsg.Divs = append(cardMsg.Divs, &pushV1.CardTextDiv{
			Fields: []*pushV1.CardTextField{
				&pushV1.CardTextField{
					Key:   columnOwnerId,
					Value: ownerDisplayName,
				},
			},
		})
	}

	if card.SpecialDiv != nil && len(card.SpecialDiv) > 0 {
		cardMsg.Divs = append(cardMsg.Divs, card.SpecialDiv...)
	}

	if card.SpecialAction != nil && len(card.SpecialAction) > 0 {
		cardMsg.ActionModules = append(cardMsg.ActionModules, card.SpecialAction...)
	}

	// 查看详情等按钮
	cardMsg.ActionModules = append(cardMsg.ActionModules, &pushV1.CardActionModule{
		Type: pushV1.CardActionModuleType_Action,
		Actions: []*pushV1.CardAction{
			&pushV1.CardAction{
				Level: pushV1.CardElementLevel_Default,
				Type:  pushV1.CardActionType_ButtonJumpUrl,
				Text:  consts.CardButtonTextForViewDetail,
				Url:   card.Links.SideBarLink,
			},
			&pushV1.CardAction{
				Level: pushV1.CardElementLevel_NotImportant,
				Type:  pushV1.CardActionType_ButtonJumpUrl,
				Text:  consts.CardButtonTextForViewInsideApp,
				Url:   card.Links.Link,
			},
		},
	})
	if card.IsNeedUnsubscribe {
		m := map[string]interface{}{
			consts.FsCardValueCardType: consts.FsCardTypeNoticeSubscribe,
			consts.FsCardValueAction:   consts.FsCardUnsubscribeIssueRemind,
			consts.FsCardValueOrgId:    card.OrgId,
		}
		values, err := structpb.NewStruct(m)
		if err != nil {
			log.Errorf("[GenerateUpdateCard] err:%v", err)
			return nil
		}
		cardMsg.ActionModules = append(cardMsg.ActionModules,
			&pushV1.CardActionModule{
				Type: pushV1.CardActionModuleType_Hr,
			},
			&pushV1.CardActionModule{
				Type: pushV1.CardActionModuleType_Action,
				Actions: []*pushV1.CardAction{
					&pushV1.CardAction{
						Level:  pushV1.CardElementLevel_Default,
						Type:   pushV1.CardActionType_Link,
						Text:   consts.CardUnSubscribe,
						Values: values,
					},
				},
			},
		)
	}

	return cardMsg
}

func GetCardTitleAndOwnerColumn(tableColumn map[string]*projectvo.TableColumnData) (string, string) {
	columnTitle := consts.CardElementTitle
	columnOwnerId := consts.CardElementOwner

	for _, column := range tableColumn {
		if column.Name == consts.BasicFieldTitle {
			if column.AliasTitle != "" {
				columnTitle = column.AliasTitle
			}
		}
		if column.Name == consts.BasicFieldOwnerId {
			if column.AliasTitle != "" {
				columnOwnerId = column.AliasTitle
			}
		}
	}
	return columnTitle, columnOwnerId
}
