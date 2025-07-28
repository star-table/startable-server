package card

import (
	"fmt"

	"strconv"
	"strings"

	pushV1 "gitea.bjx.cloud/LessCode/interface/golang/push/v1"
	fsSdkVo "gitea.bjx.cloud/allstar/feishu-sdk-golang/core/model/vo"
	"gitea.bjx.cloud/allstar/feishu-sdk-golang/sdk"
	sdk_const "gitea.bjx.cloud/allstar/platform-sdk/consts"
	"github.com/star-table/startable-server/common/core/consts"
	"github.com/star-table/startable-server/common/core/errs"
	"github.com/star-table/startable-server/common/core/logger"
	"github.com/star-table/startable-server/common/core/util/json"
	"github.com/star-table/startable-server/common/core/util/str"
	"github.com/star-table/startable-server/common/extra/feishu"
	"github.com/star-table/startable-server/common/model/bo"
	"github.com/star-table/startable-server/common/model/vo"
	"github.com/star-table/startable-server/common/model/vo/callvo"
	"github.com/star-table/startable-server/common/model/vo/commonvo"
	"github.com/star-table/startable-server/common/model/vo/projectvo"
)

var log = logger.GetDefaultLogger()

func SendCardMetaFeiShu(corpId string, meta *commonvo.CardMeta, openIds []string) {
	if len(openIds) == 0 || meta == nil {
		return
	}
	tenant, sysErr := feishu.GetTenant(corpId)
	if sysErr != nil {
		log.Errorf("[SendCardMetaFeiShu] GetTenant err: %v", sysErr.Error())
		return
	}
	SendCardMetaFeiShuByTenant(tenant, meta, openIds)
}

func SendCardMetaFeiShuByTenant(tenant *sdk.Tenant, meta *commonvo.CardMeta, openIds []string) {
	if len(openIds) == 0 || meta == nil {
		return
	}

	card := GenerateFeiShuCard(meta)
	log.Infof("[SendCardMetaFeiShu] SendMessageBatch to: %v, card: %v", openIds, json.ToJsonIgnoreError(card))
	resp, err := tenant.SendMessageBatch(fsSdkVo.BatchMsgVo{
		OpenIds: openIds,
		MsgType: "interactive",
		Card:    card,
	})
	if err != nil {
		log.Errorf("[SendCardMetaFeiShu] SendMessageBatch err: %v", err)
		return
	}
	if resp.Code != 0 {
		log.Errorf("[SendCardMetaFeiShu] SendMessageBatch code: %d, errMsg: %s, card: %s", resp.Code, resp.Msg, json.ToJsonIgnoreError(card))
		return
	}
	//str, _ := json.ToJson(resp)
	//log.Info("[SendCardFeiShu] openIds 发送成功: " + str)
}

func SendCardFeiShu(corpId string, card *commonvo.Card, openIds []string) {
	if len(openIds) == 0 || card == nil || card.FsCard == nil {
		return
	}
	tenant, sysErr := feishu.GetTenant(corpId)
	if sysErr != nil {
		log.Errorf("[SendCardFeiShu] GetTenant err: %v", sysErr.Error())
		return
	}
	SendCardFeiShuByTenant(tenant, card, openIds)
}

func SendCardFeiShuByTenant(tenant *sdk.Tenant, card *commonvo.Card, openIds []string) {
	if len(openIds) == 0 || card == nil || card.FsCard == nil {
		return
	}

	log.Infof("[SendCardFeiShu] SendMessageBatch to: %v, card: %v", openIds, json.ToJsonIgnoreError(card.FsCard))
	resp, err := tenant.SendMessageBatch(fsSdkVo.BatchMsgVo{
		OpenIds: openIds,
		MsgType: "interactive",
		Card:    card.FsCard,
	})
	if err != nil {
		log.Errorf("[SendCardFeiShu] SendMessageBatch err: %v", err)
		return
	}
	if resp.Code != 0 {
		log.Errorf("[SendCardFeiShu] SendMessageBatch code: %d, errMsg: %s, card: %s", resp.Code, resp.Msg, json.ToJsonIgnoreError(card))
		return
	}
	//str, _ := json.ToJson(resp)
	//log.Info("[SendCardFeiShu] openIds 发送成功: " + str)
}

// 飞书发送卡片消息：当用户在个人机器人中做未知的操作指令
func SendCardFeiShuBotNotKnown(tenantKey string, orgId int64, chatId string, openIds ...string) errs.SystemErrorInfo {
	//tenantClient, err := feishu.GetTenant(tenantKey)
	//if err != nil {
	//	log.Error(err)
	//	return err
	//}

	cardInfo := GetFsCardBotSimpleInfo(orgId)
	pushCardReq := &commonvo.PushCard{
		OrgId:         orgId,
		OutOrgId:      tenantKey,
		SourceChannel: sdk_const.SourceChannelFeishu,
		CardMsg:       cardInfo,
	}
	if chatId != "" {
		pushCardReq.ChatIds = []string{chatId}
	} else {
		pushCardReq.OpenIds = openIds
	}
	errSys := PushCard(orgId, pushCardReq)
	if errSys != nil {
		log.Errorf("[SendCardFeiShuBotNotKnown] err:%v", errSys)
		return errSys
	}

	log.Infof("发送机器人欢迎通知成功")

	return nil
}

// SendCardFeiShuReplyIssueChatIns 飞书群聊 - 响应用户在任务群聊中输入的指令。
func SendCardFeiShuReplyIssueChatIns(chatId string, infoBox *callvo.GroupChatHandleInfo) errs.SystemErrorInfo {
	cardMsg := &pushV1.TemplateCard{
		Divs: []*pushV1.CardTextDiv{
			&pushV1.CardTextDiv{
				Fields: []*pushV1.CardTextField{
					&pushV1.CardTextField{
						Key:   consts.FsCardReplyIssueChatIns,
						Value: "",
					},
				},
			},
		},
		ActionModules: []*pushV1.CardActionModule{
			&pushV1.CardActionModule{
				Type: pushV1.CardActionModuleType_Action,
				Actions: []*pushV1.CardAction{
					&pushV1.CardAction{
						Level: pushV1.CardElementLevel_Default,
						Type:  pushV1.CardActionType_ButtonJumpUrl,
						Text:  consts.CardButtonTextForViewDetail,
						Url:   infoBox.IssueInfoUrl,
					},
					&pushV1.CardAction{
						Level: pushV1.CardElementLevel_NotImportant,
						Type:  pushV1.CardActionType_ButtonJumpUrl,
						Text:  consts.CardButtonTextForViewInsideApp,
						Url:   infoBox.IssuePcUrl,
					},
				},
			},
		},
	}
	errSys := PushCard(infoBox.OrgInfo.OrgId, &commonvo.PushCard{
		OrgId:         infoBox.OrgInfo.OrgId,
		OutOrgId:      infoBox.OrgInfo.OutOrgId,
		SourceChannel: sdk_const.SourceChannelFeishu,
		ChatIds:       []string{chatId},
		CardMsg:       cardMsg,
	})
	if errSys != nil {
		log.Errorf("[SendCardFeiShuReplyIssueChatIns] pushCard err:%v", errSys)
		return errSys
	}

	return nil
}

// SendCardFeiShuIssueChatUserJoinToOpUser 给拉人的操作发送提示卡片（仅你可见的临时消息）
func SendCardFeiShuIssueChatUserJoinToOpUser(chatId string, infoBox *callvo.GroupChatHandleInfo, joinUserNames []string) errs.SystemErrorInfo {
	elements := make([]interface{}, 0, 3)
	elements = append(elements, fsSdkVo.CardElementContentModule{
		Tag: "div",
		Fields: []fsSdkVo.CardElementField{
			{
				Text: fsSdkVo.CardElementText{
					Tag:     "lark_md",
					Content: fmt.Sprintf("**%s不是任务协作者，只是临时参与讨论哦~**", strings.Join(joinUserNames, "、")),
				},
			},
		},
	})
	cardMsg := &fsSdkVo.Card{
		Elements: elements,
	}
	resp, err := infoBox.TenantClient.SendEphemeralMessage(fsSdkVo.EphemeralMsg{
		ChatId:  chatId,
		OpenId:  infoBox.OperatorOpenId,
		MsgType: "interactive",
		Card:    cardMsg,
	})
	if err != nil {
		log.Errorf("[SendCardFeiShuIssueChatUserJoinToOpUser] SendMessage err: %v", err)
		return errs.BuildSystemErrorInfoWithMessage(errs.FeiShuOpenApiCallError, err.Error())
	}
	if resp.Code != 0 {
		log.Errorf("[SendCardFeiShuIssueChatUserJoinToOpUser] resp: %s", json.ToJsonIgnoreError(resp))
		return errs.BuildSystemErrorInfoWithMessage(errs.FeiShuOpenApiCallError, resp.Msg)
	}

	return nil
}

// SendCardFeiShuIssueChatUserDeleteToOpUser 用户被移出任务群聊时，给操作者发送提示卡片（仅你可见的临时消息）
func SendCardFeiShuIssueChatUserDeleteToOpUser(chatId string, infoBox *callvo.GroupChatHandleInfo, joinUserNames []string) errs.SystemErrorInfo {
	elements := make([]interface{}, 0, 3)
	elements = append(elements, fsSdkVo.CardElementContentModule{
		Tag: "div",
		Fields: []fsSdkVo.CardElementField{
			{
				Text: fsSdkVo.CardElementText{
					Tag:     "lark_md",
					Content: fmt.Sprintf("**你已将%s移除出群聊，但他/她仍是该条任务的协作人~\n若要移除他/她的协作人身份，请去任务内操作。**", strings.Join(joinUserNames, "、")),
				},
			},
		},
	})
	elements = append(elements, fsSdkVo.CardElementActionModule{
		Tag: "action",
		Actions: []interface{}{
			fsSdkVo.ActionButton{
				Tag: "button",
				Text: fsSdkVo.CardElementText{
					Tag:     "plain_text",
					Content: consts.CardButtonTextForViewDetail,
				},
				Url:  infoBox.IssueInfoUrl,
				Type: consts.FsCardButtonColorPrimary,
			},
			fsSdkVo.ActionButton{
				Tag: "button",
				Text: fsSdkVo.CardElementText{
					Tag:     "plain_text",
					Content: consts.CardButtonTextForViewInsideApp,
				},
				Url:  infoBox.IssuePcUrl,
				Type: consts.FsCardButtonColorDefault,
			},
		},
	})
	cardMsg := &fsSdkVo.Card{
		Elements: elements,
	}
	resp, err := infoBox.TenantClient.SendEphemeralMessage(fsSdkVo.EphemeralMsg{
		ChatId:  chatId,
		OpenId:  infoBox.OperatorOpenId,
		MsgType: "interactive",
		Card:    cardMsg,
	})
	if err != nil {
		log.Errorf("[SendCardFeiShuIssueChatUserDeleteToOpUser] SendMessage err: %v", err)
		return errs.BuildSystemErrorInfoWithMessage(errs.FeiShuOpenApiCallError, err.Error())
	}
	if resp.Code != 0 {
		log.Errorf("[SendCardFeiShuIssueChatUserDeleteToOpUser] resp: %s", json.ToJsonIgnoreError(resp))
		return errs.BuildSystemErrorInfoWithMessage(errs.FeiShuOpenApiCallError, resp.Msg)
	}

	return nil
}

// 被催者回复后 卡片变为
func GetFsCardAfterReplyUrgeIssue(card projectvo.UrgeReplyCard) *pushV1.TemplateCard {
	issue := card.Issue
	// 催办
	cardDivs := []*pushV1.CardTextDiv{
		&pushV1.CardTextDiv{
			Fields: []*pushV1.CardTextField{
				&pushV1.CardTextField{
					Key:   fmt.Sprintf(consts.CardUrgePeople, card.OperateUserName),
					Value: fmt.Sprintf(consts.CardDoubleQuotationMarks, card.UrgeText),
				},
			},
		},
	}

	// 回复
	cardDivs = append(cardDivs, &pushV1.CardTextDiv{
		Fields: []*pushV1.CardTextField{
			&pushV1.CardTextField{
				Key:   consts.CardElementReply,
				Value: fmt.Sprintf(consts.CardDoubleQuotationMarks, card.ReplyMsg),
			},
		},
	})

	cardMsg := GenerateCard(&commonvo.GenerateCard{
		OrgId:             card.OrgId,
		CardTitle:         consts.CardReplyUrgeSelf,
		CardLevel:         pushV1.CardElementLevel_Default,
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

	return cardMsg
}

// 当用户在个人机器人中做未知的操作指令时推送的卡片信息
func GetFsCardBotSimpleInfo(orgId int64) *pushV1.TemplateCard {
	cardMsg := &pushV1.TemplateCard{
		Divs: []*pushV1.CardTextDiv{
			&pushV1.CardTextDiv{
				Fields: []*pushV1.CardTextField{
					&pushV1.CardTextField{
						Key:   "抱歉，我没有理解你的意思。你是否要做如下操作：",
						Value: "",
					},
				},
			},
			&pushV1.CardTextDiv{
				Fields: []*pushV1.CardTextField{
					&pushV1.CardTextField{
						Key: "发送 **create 任务名称**，快速创建任务；\n" +
							"发送 **设置**，进入极星机器人通知设置页面。",
						Value: "",
					},
				},
			},
			&pushV1.CardTextDiv{
				Fields: []*pushV1.CardTextField{
					&pushV1.CardTextField{
						Key:   "你也可以点击 [极星协作](" + feishu.GetAppLinkPcWelcomeUrl(orgId) + ") 直接进入应用。",
						Value: "",
					},
				},
			},
		},
	}

	return cardMsg
}

// 即将逾期的卡片，点击完成任务的交互后的卡片信息
func GetFsCardBeOverdueComplete(title string, projectId int64, projectName, tableName string,
	ownerInfos []*bo.BaseUserInfoBo, issueLinks *projectvo.IssueLinks, tableColumnMap map[string]*projectvo.TableColumnData) *pushV1.TemplateCard {

	cardMsg := GenerateCard(&commonvo.GenerateCard{
		CardTitle:       consts.CardCompleteIssue,
		CardLevel:       pushV1.CardElementLevel_Default,
		ProjectName:     projectName,
		TableName:       tableName,
		IssueTitle:      title,
		ProjectId:       projectId,
		OwnerInfos:      ownerInfos,
		Links:           issueLinks,
		TableColumnInfo: tableColumnMap,
	})

	return cardMsg
}

// 用户指令创建任务
func GetFsCardCreateNewIssue(title, projectName, tableName string, ownerInfos []*bo.BaseUserInfoBo, issueLink *projectvo.IssueLinks, issueVo vo.Issue) *pushV1.TemplateCard {
	cardMsg := GenerateCard(&commonvo.GenerateCard{
		CardTitle:   "新记录",
		CardLevel:   pushV1.CardElementLevel_Default,
		IssueTitle:  title,
		ProjectName: projectName,
		TableName:   tableName,
		ProjectId:   issueVo.ProjectID,
		OwnerInfos:  ownerInfos,
		Links:       issueLink,
	})

	return cardMsg
}

// GetFsCardStartIssueChat 飞书卡片-发起任务群聊
func GetFsCardStartIssueChat(orgId int64, topic string, infoBox *projectvo.CardInfoBoxForStartIssueChat) *pushV1.TemplateCard {
	issue := infoBox.IssueInfo
	if issue.Title == "" {
		issue.Title = consts.IssueChatTitleNoNameTitle
	}

	// 父记录
	parentTitle := consts.CardDefaultRelationIssueTitle
	if issue.ParentId > 0 {
		if infoBox.ParentIssue.Title != "" {
			parentTitle = infoBox.ParentIssue.Title
		}
	}

	cardDivs := []*pushV1.CardTextDiv{
		&pushV1.CardTextDiv{
			Fields: []*pushV1.CardTextField{
				&pushV1.CardTextField{
					Key:   consts.CardElementChatTopic,
					Value: topic,
				},
			},
		},
	}

	projectName := ""
	projectTypeId := int64(0)
	if infoBox.ProjectBo != nil {
		projectName = infoBox.ProjectBo.Name
		projectTypeId = infoBox.ProjectBo.ProjectTypeId
	}

	cardMsg := GenerateCard(&commonvo.GenerateCard{
		OrgId:         orgId,
		CardTitle:     fmt.Sprintf("%s 发起了任务讨论", infoBox.OperateUser.Name),
		CardLevel:     pushV1.CardElementLevel_Default,
		ParentTitle:   parentTitle,
		ProjectName:   projectName,
		ProjectTypeId: projectTypeId,
		TableName:     infoBox.IssueTableInfo.Name,
		IssueTitle:    issue.Title,
		ProjectId:     issue.ProjectId,
		ParentId:      issue.ParentId,
		OwnerInfos:    infoBox.IssueOwners,
		Links: &projectvo.IssueLinks{
			Link:        infoBox.IssuePcUrl,
			SideBarLink: infoBox.IssueInfoUrl,
		},
		IsNeedUnsubscribe: false,
		SpecialDiv:        cardDivs,
		TableColumnInfo:   infoBox.ProjectTableColumnMap,
	})

	return cardMsg
}

// / 回复提示信息卡片：抱歉，您指定的用户不是项目成员。
func GetFsCardForSomeTip(tipText string, extraParam map[string]interface{}) *pushV1.TemplateCard {
	// 实现 at 某个成员
	openId := ""
	userName := ""
	if tmpVal, ok := extraParam["atOpenId"]; ok {
		openId = tmpVal.(string)
	}
	if tmpVal, ok := extraParam["atUserName"]; ok {
		userName = tmpVal.(string)
	}
	atSomeoneStr := str.RenderAtSomeoneStr(openId, userName)
	if len(openId) > 0 && len(atSomeoneStr) > 0 {
		tipText = atSomeoneStr + tipText
	}

	cardMsg := &pushV1.TemplateCard{
		Title: &pushV1.CardTitle{
			Title: "提示信息",
		},
		Divs: []*pushV1.CardTextDiv{
			&pushV1.CardTextDiv{
				Fields: []*pushV1.CardTextField{
					&pushV1.CardTextField{
						Key:   tipText,
						Value: "",
					},
				},
			},
		},
	}

	return cardMsg
}

// GetFsCardHelpOfBind0ProForGroupChat 组装绑定 0 个项目时的帮助卡片
func GetFsCardHelpOfBind0ProForGroupChat(orgId int64, chatId string) *pushV1.TemplateCard {
	return GetFsCardSettingMsg(orgId, chatId, consts.FsCardHelpOfBind0ProForGroupChat)
}

// GetFsCardHelpOfBind1ProForGroupChat 组装绑定 1 个项目时的帮助卡片
func GetFsCardHelpOfBind1ProForGroupChat() *pushV1.TemplateCard {
	cardMsg := &pushV1.TemplateCard{
		Divs: []*pushV1.CardTextDiv{
			&pushV1.CardTextDiv{
				Fields: []*pushV1.CardTextField{
					&pushV1.CardTextField{
						Key: "已将1个项目关联到本群聊\n" +
							"您关联项目的消息将推送到本群聊，您还可以通过以下方式与项目交互，例如：\n" + consts.FsCardHelpInfoForInstructionHelp,
					},
				},
			},
		},
	}

	return cardMsg
}

// GetFsCardHelpOfBind2ProForGroupChat 组装绑定 2 个及以上项目时的帮助卡片
func GetFsCardHelpOfBind2ProForGroupChat(proNum int) *pushV1.TemplateCard {

	cardMsg := &pushV1.TemplateCard{
		Divs: []*pushV1.CardTextDiv{
			&pushV1.CardTextDiv{
				Fields: []*pushV1.CardTextField{
					&pushV1.CardTextField{
						Key: "已将" + strconv.Itoa(proNum) + "个项目关联到本群聊\n" +
							"您关联项目的消息将推送到本群聊。" +
							"您还可以发送 @极星协作 设置，重新进行推送设置。",
					},
				},
			},
		},
	}

	return cardMsg
}

// / 响应指令：项目进展
func GetFsCardProProgress(projectId int64, title string) *pushV1.TemplateCard {
	cardMsg := &pushV1.TemplateCard{
		Title: &pushV1.CardTitle{
			Title: title,
			Level: pushV1.CardElementLevel_Default,
		},
	}
	return cardMsg
}

// / 响应指令：项目设置
func GetFsCardProjectSettings(orgId int64, chatId string) *pushV1.TemplateCard {
	return GetFsCardSettingMsg(orgId, chatId, consts.FsCardProjectSettings)
}

// / 响应指令：项目任务 敏捷项目
func GetFsCardProIssueCardForAgile(projectId, projectTypeId int64, title string, proIssueStat *vo.IssueStatusTypeStatResp) *pushV1.TemplateCard {

	cardMsg := &pushV1.TemplateCard{
		Title: &pushV1.CardTitle{
			Title: title,
		},
		Divs: []*pushV1.CardTextDiv{
			&pushV1.CardTextDiv{
				Fields: []*pushV1.CardTextField{
					&pushV1.CardTextField{
						Key:   "待处理：" + strconv.FormatInt(proIssueStat.NotStartTotal, 10),
						Value: "",
					},
					&pushV1.CardTextField{
						Key:   "正在处理：" + strconv.FormatInt(proIssueStat.ProcessingTotal, 10),
						Value: "",
					},
					&pushV1.CardTextField{
						Key:   "即将逾期：" + strconv.FormatInt(proIssueStat.OverdueTomorrowTotal, 10),
						Value: "",
					},
					&pushV1.CardTextField{
						Key:   "已逾期：" + strconv.FormatInt(proIssueStat.OverdueTotal, 10),
						Value: "",
					},
					&pushV1.CardTextField{
						Key:   "已完成：" + strconv.FormatInt(proIssueStat.CompletedTotal, 10),
						Value: "",
					},
					&pushV1.CardTextField{
						Key:   "任务总数：" + strconv.FormatInt(proIssueStat.Total, 10),
						Value: "",
					},
				},
			},
		},
	}

	return cardMsg
}

// / 响应指令：项目任务 通用项目
func GetFsCardProIssueCardForCommon(projectId, projectTypeId int64, title string, proIssueStat *vo.IssueStatusTypeStatResp) *pushV1.TemplateCard {
	cardMsg := &pushV1.TemplateCard{
		Title: &pushV1.CardTitle{
			Title: title,
		},
		Divs: []*pushV1.CardTextDiv{
			&pushV1.CardTextDiv{
				Fields: []*pushV1.CardTextField{
					&pushV1.CardTextField{
						Key:   "未完成：" + strconv.FormatInt(proIssueStat.NotStartTotal+proIssueStat.ProcessingTotal, 10),
						Value: "",
					},
					&pushV1.CardTextField{
						Key:   "即将逾期：" + strconv.FormatInt(proIssueStat.BeAboutToOverdueSum, 10),
						Value: "",
					},
					&pushV1.CardTextField{
						Key:   "已逾期：" + strconv.FormatInt(proIssueStat.OverdueTotal, 10),
						Value: "",
					},
					&pushV1.CardTextField{
						Key:   "已完成：" + strconv.FormatInt(proIssueStat.CompletedTotal, 10),
						Value: "",
					},
					&pushV1.CardTextField{
						Key:   "任务总数：" + strconv.FormatInt(proIssueStat.Total, 10),
						Value: "",
					},
				},
			},
		},
	}

	return cardMsg
}

// 回复卡片：帮助信息
func GetFsCardGroupChatHelp(orgId int64, projectId int64, chaId string) *pushV1.TemplateCard {

	card := &pushV1.TemplateCard{
		Divs: []*pushV1.CardTextDiv{
			&pushV1.CardTextDiv{
				Fields: []*pushV1.CardTextField{
					&pushV1.CardTextField{
						Key:   "你好，我是你的项目管理小助手。\n项目相关的动态会同步到这个群里",
						Value: "",
					},
				},
			},
			&pushV1.CardTextDiv{
				Fields: []*pushV1.CardTextField{
					&pushV1.CardTextField{
						Key:   "您可以通过以下方式与项目交互，例如：\n" + consts.FsCardHelpInfoForInstructionHelp,
						Value: "",
					},
				},
			},
			&pushV1.CardTextDiv{
				Fields: []*pushV1.CardTextField{
					&pushV1.CardTextField{
						Key:   "您还可以点击下面的按钮调整动态通知范围：",
						Value: "",
					},
				},
			},
		},
		ActionModules: []*pushV1.CardActionModule{
			&pushV1.CardActionModule{
				Type: pushV1.CardActionModuleType_Action,
				Actions: []*pushV1.CardAction{
					&pushV1.CardAction{
						Level: pushV1.CardElementLevel_NotImportant,
						Type:  pushV1.CardActionType_ButtonJumpUrl,
						Text:  consts.CardProTrendsSettings,
						Url:   feishu.GetGroupChatProjectTrendPushConfigAppLink(projectId, orgId, chaId),
					},
				},
			},
		},
	}
	return card
}

func GetFsCardMsgForBindPro(content string) *pushV1.TemplateCard {
	cardMsg := &pushV1.TemplateCard{
		Divs: []*pushV1.CardTextDiv{
			&pushV1.CardTextDiv{
				Fields: []*pushV1.CardTextField{
					&pushV1.CardTextField{
						Key:   content,
						Value: "",
					},
				},
			},
		},
	}

	return cardMsg
}

func GetFsCardForNotBindProjectMsg(orgId int64, chatId string) *pushV1.TemplateCard {
	return GetFsCardSettingMsg(orgId, chatId, consts.FsCardNoBindPro)

}

func GetFsCardSettingMsg(orgId int64, chatId, content string) *pushV1.TemplateCard {
	cardMsg := &pushV1.TemplateCard{
		Divs: []*pushV1.CardTextDiv{
			&pushV1.CardTextDiv{
				Fields: []*pushV1.CardTextField{
					&pushV1.CardTextField{
						Key:   content,
						Value: "",
					},
				},
			},
		},
		ActionModules: []*pushV1.CardActionModule{
			&pushV1.CardActionModule{
				Type: pushV1.CardActionModuleType_Action,
				Actions: []*pushV1.CardAction{
					&pushV1.CardAction{
						Level: pushV1.CardElementLevel_NotImportant,
						Type:  pushV1.CardActionType_ButtonJumpUrl,
						Text:  consts.CardProTrendsSettings,
						Url:   feishu.GetGroupChatSettingBtnAppLink(orgId, chatId),
					},
				},
			},
		},
	}
	return cardMsg
}

func GetFsCardNotSupportChatInstruction(orgId int64, chatId, content string) *pushV1.TemplateCard {
	return GetFsCardSettingMsg(orgId, chatId, content)
}

func GetFsBotNotification(orgId int64) *pushV1.TemplateCard {
	cardMsg := &pushV1.TemplateCard{
		Divs: []*pushV1.CardTextDiv{
			&pushV1.CardTextDiv{
				Fields: []*pushV1.CardTextField{
					&pushV1.CardTextField{
						Key:   consts.FsCardBotNotificationText,
						Value: "",
					},
				},
			},
		},
		ActionModules: []*pushV1.CardActionModule{
			&pushV1.CardActionModule{
				Type: pushV1.CardActionModuleType_Action,
				Actions: []*pushV1.CardAction{
					&pushV1.CardAction{
						Level: pushV1.CardElementLevel_NotImportant,
						Type:  pushV1.CardActionType_ButtonJumpUrl,
						Text:  consts.FsCardBotNotificationButton,
						Url:   feishu.GetDefaultNotificationAppLink(orgId),
					},
				},
			},
		},
	}
	return cardMsg
}

func GetFsCardSubscribeNotice(content string) *pushV1.TemplateCard {
	cardMsg := &pushV1.TemplateCard{
		Divs: []*pushV1.CardTextDiv{
			&pushV1.CardTextDiv{
				Fields: []*pushV1.CardTextField{
					&pushV1.CardTextField{
						Key:   content,
						Value: "",
					},
				},
			},
		},
	}
	return cardMsg
}
