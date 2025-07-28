package service

/// 飞书群聊-用户指令 回复卡片

import (
	pushPb "gitea.bjx.cloud/LessCode/interface/golang/push/v1"
	sdk_const "gitea.bjx.cloud/allstar/platform-sdk/consts"
	"github.com/star-table/startable-server/common/core/consts"
	"github.com/star-table/startable-server/common/core/errs"
	"github.com/star-table/startable-server/common/extra/card"
	"github.com/star-table/startable-server/common/model/vo"
	"github.com/star-table/startable-server/common/model/vo/commonvo"
)

///// 发送卡片：催办任务
//func GetCardForUrgeIssue(orgId, issueId, opUserId int64, opUserName string, projectId int64, projectName string, issue *bo.IssueBo, urgeText string) (*fsSdkVo.Card, errs.SystemErrorInfo) {
//	issueLinks := domain.GetIssueLinks("", orgId, issueId)
//	elements := []interface{}{
//		fsSdkVo.CardElementContentModule{
//			Tag: "div",
//			Fields: []fsSdkVo.CardElementField{
//				{
//					Text: fsSdkVo.CardElementText{
//						Tag:     "lark_md",
//						Content: issue.Title,
//					},
//				},
//			},
//		},
//		fsSdkVo.CardElementContentModule{
//			Tag: "div",
//			Fields: []fsSdkVo.CardElementField{
//				{
//					Text: fsSdkVo.CardElementText{
//						Tag:     "lark_md",
//						Content: "催办：" + urgeText,
//					},
//				},
//			},
//		},
//	}
//	if issue.ParentId != 0 {
//		parentInfo, err := domain.GetIssueBo(orgId, issue.ParentId)
//		if err != nil {
//			log.Error(err)
//			return nil, err
//		}
//		elements = append(elements, fsSdkVo.CardElementContentModule{
//			Tag: "div",
//			Fields: []fsSdkVo.CardElementField{
//				{
//					Text: fsSdkVo.CardElementText{
//						Tag:     "lark_md",
//						Content: "父记录：" + parentInfo.Title,
//					},
//				},
//			},
//		})
//	}
//	elements = append(elements, fsSdkVo.CardElementContentModule{
//		Tag: "div",
//		Fields: []fsSdkVo.CardElementField{
//			{
//				Text: fsSdkVo.CardElementText{
//					Tag:     "lark_md",
//					Content: "所属项目：" + projectName,
//				},
//			},
//		},
//	},
//		fsSdkVo.CardElementActionModule{
//			Tag: "action",
//			Actions: []interface{}{
//				fsSdkVo.ActionButton{
//					Tag: "button",
//					Text: fsSdkVo.CardElementText{
//						Tag:     "plain_text",
//						Content: "任务详情",
//					},
//					Url:  issueLinks.SideBarLink,
//					Type: "default",
//				},
//				fsSdkVo.ActionSelectMenu{
//					Tag: "select_static",
//					Placeholder: &fsSdkVo.CardElementText{
//						Tag:     "plain_text",
//						Content: "现在忙，晚点处理",
//					},
//					Options: []fsSdkVo.CardElementOption{
//						{
//							Text: &fsSdkVo.CardElementText{
//								Tag:     "lark_md",
//								Content: "现在忙，晚点处理",
//							},
//							Value: "现在忙，晚点处理",
//						}, {
//							Text: &fsSdkVo.CardElementText{
//								Tag:     "lark_md",
//								Content: "抱歉，现在不方便处理",
//							},
//							Value: "抱歉，现在不方便处理",
//						}, {
//							Text: &fsSdkVo.CardElementText{
//								Tag:     "lark_md",
//								Content: "暂不处理",
//							},
//							Value: "暂不处理",
//						}, {
//							Text: &fsSdkVo.CardElementText{
//								Tag:     "lark_md",
//								Content: "忽略",
//							},
//							Value: "忽略",
//						},
//					},
//					Value: map[string]interface{}{
//						consts.FsCardValueIssueId:  issueId,
//						consts.FsCardValueCardType: consts.FsCardTypeIssueUrge,
//						consts.FsCardValueOrgId:    orgId,
//						consts.FsCardValueUserId:   opUserId,
//						//consts.FsCardValueUrgeIssueBeUrgedUserId: beUrgedUserId,
//						consts.FsCardValueUrgeIssueUrgeText: urgeText,
//					},
//				},
//			},
//		})
//	cardMsg := &fsSdkVo.Card{
//		Header: &fsSdkVo.CardHeader{
//			Title: &fsSdkVo.CardHeaderTitle{
//				Tag:     "plain_text",
//				Content: opUserName + " 催办了任务",
//			},
//		},
//		//Config: &fsSdkVo.CardConfig{UpdateMulti: true},
//		Elements: elements,
//	}
//
//	return cardMsg, nil
//}

//func GetCardForTest(orgId, issueId, opUserId int64, opUserName string, project *bo.Project, issue *bo.IssueBo, urgeText string) *fsSdkVo.Card {
//	issueLinks := domain.GetIssueLinks("", orgId, issueId)
//	cardMsg := &fsSdkVo.Card{
//		Header: &fsSdkVo.CardHeader{
//			Title: &fsSdkVo.CardHeaderTitle{
//				Tag:     "plain_text",
//				Content: opUserName + " 催办了任务",
//			},
//		},
//		Elements: []interface{}{
//			fsSdkVo.CardElementContentModule{
//				Tag: "div",
//				Fields: []fsSdkVo.CardElementField{
//					{
//						Text: fsSdkVo.CardElementText{
//							Tag:     "lark_md",
//							Content: issue.Title,
//						},
//					},
//				},
//			},
//			fsSdkVo.CardElementContentModule{
//				Tag: "div",
//				Fields: []fsSdkVo.CardElementField{
//					{
//						Text: fsSdkVo.CardElementText{
//							Tag:     "lark_md",
//							Content: "催办：" + urgeText,
//						},
//					},
//				},
//			},
//			fsSdkVo.CardElementContentModule{
//				Tag: "div",
//				Fields: []fsSdkVo.CardElementField{
//					{
//						Text: fsSdkVo.CardElementText{
//							Tag:     "lark_md",
//							Content: "所属项目：" + project.Name,
//						},
//					},
//				},
//			},
//			fsSdkVo.CardElementActionModule{
//				Tag: "action",
//				Actions: []interface{}{
//					fsSdkVo.ActionButton{
//						Tag: "button",
//						Text: fsSdkVo.CardElementText{
//							Tag:     "plain_text",
//							Content: "任务详情",
//						},
//						Url:  issueLinks.SideBarLink,
//						Type: "default",
//					},
//					//fsSdkVo.ActionSelectMenu{
//					//	Tag: "select_static",
//					//	Placeholder: &fsSdkVo.CardElementText{
//					//		Tag:     "plain_text",
//					//		Content: "现在忙，晚点处理",
//					//	},
//					//	Value: map[string]interface{}{
//					//		consts.FsCardValueIssueId:  issueId,
//					//		consts.FsCardValueCardType: consts.FsCardTypeIssueUrge,
//					//		consts.FsCardValueOrgId:    orgId,
//					//		consts.FsCardValueUserId:   opUserId,
//					//	},
//					//	Options: []fsSdkVo.CardElementOption{
//					//		{
//					//			Text: &fsSdkVo.CardElementText{
//					//				Tag:     "lark_md",
//					//				Content: "现在忙，晚点处理",
//					//			},
//					//			Value: "现在忙，晚点处理",
//					//		}, {
//					//			Text: &fsSdkVo.CardElementText{
//					//				Tag:     "lark_md",
//					//				Content: "抱歉，现在不方便处理",
//					//			},
//					//			Value: "抱歉，现在不方便处理",
//					//		}, {
//					//			Text: &fsSdkVo.CardElementText{
//					//				Tag:     "lark_md",
//					//				Content: "暂不处理",
//					//			},
//					//			Value: "暂不处理",
//					//		}, {
//					//			Text: &fsSdkVo.CardElementText{
//					//				Tag:     "lark_md",
//					//				Content: "忽略",
//					//			},
//					//			Value: "忽略",
//					//		},
//					//	},
//					//},
//				},
//			},
//		},
//	}
//
//	return cardMsg
//}

// 机器人回复消息卡片
func GroupChatReplyToUserInsProIssue(sourceChannel string, outOrgId string, chatId string, orgId, projectId, projectTypeId int64, proStatResp *vo.IssueStatusTypeStatResp) errs.SystemErrorInfo {
	switch sourceChannel {
	case sdk_const.SourceChannelFeishu:
		//tenant, err := feishu.GetTenant(outOrgId)
		//if err != nil {
		//	log.Error(err)
		//	return err
		//}
		var msgCard *pushPb.TemplateCard
		if projectTypeId == consts.ProjectTypeAgileId {
			msgCard = card.GetFsCardProIssueCardForAgile(projectId, projectTypeId, "项目任务", proStatResp)
		} else {
			msgCard = card.GetFsCardProIssueCardForCommon(projectId, projectTypeId, "项目任务", proStatResp)
		}
		errSys := card.PushCard(orgId, &commonvo.PushCard{
			OrgId:         orgId,
			OutOrgId:      outOrgId,
			SourceChannel: sourceChannel,
			ChatIds:       []string{chatId},
			CardMsg:       msgCard,
		})
		if errSys != nil {
			log.Errorf("[GroupChatReplyToUserInsProIssue] err:%v", errSys)
			return errSys
		}
		//fsResp, oriErr := tenant.SendMessage(fsSdkVo.MsgVo{
		//	ChatId:  chatId,
		//	MsgType: "interactive",
		//	Card:    msgCard,
		//})
		//if oriErr != nil {
		//	log.Error(oriErr)
		//	return err
		//}
		//if fsResp.Code != 0 {
		//	log.Error("fs发送消息异常")
		//}
		//log.Infof("飞书群聊-用户指令处理-HandleGroupChatUserInsProIssue-响应用户调用结果 code: %v", fsResp.Code)
	case sdk_const.SourceChannelDingTalk:
		// 钉钉 群助手发送的互动卡片
	}

	return nil
}

func GroupChatReplyToUserInProgress(sourceChannel string, outOrgId string, chatId string, orgId, projectId int64, proStatusName string) errs.SystemErrorInfo {
	// 机器人将信息回复给用户
	switch sourceChannel {
	case sdk_const.SourceChannelFeishu:
		//tenant, err := feishu.GetTenant(outOrgId)
		//if err != nil {
		//	log.Error(err)
		//	return err
		//}
		msgCard := card.GetFsCardProProgress(projectId, "项目进展："+proStatusName)
		errSys := card.PushCard(orgId, &commonvo.PushCard{
			OrgId:         orgId,
			OutOrgId:      outOrgId,
			SourceChannel: sourceChannel,
			ChatIds:       []string{chatId},
			CardMsg:       msgCard,
		})
		if errSys != nil {
			log.Errorf("[GroupChatReplyToUserInProgress] err:%v", errSys)
			return errSys
		}
		//fsResp, oriErr := tenant.SendMessage(fsSdkVo.MsgVo{
		//	ChatId:  chatId,
		//	MsgType: "interactive",
		//	Card:    msgCard,
		//})
		//if oriErr != nil {
		//	log.Error(oriErr)
		//	return err
		//}
		//if fsResp.Code != 0 {
		//	log.Error("发送消息异常")
		//}
		//log.Infof("飞书群聊-用户指令处理-GroupChatReplyToUserInProgress-响应用户调用结果 code: %v", fsResp.Code)
	case sdk_const.SourceChannelDingTalk:

	}

	return nil
}

func GroupChatReplyToUserInsProjectSettings(sourceChannel string, outOrgId string, chatId string, orgId int64) errs.SystemErrorInfo {
	switch sourceChannel {
	case sdk_const.SourceChannelFeishu:
		// 机器人将信息回复给用户
		//tenant, err := feishu.GetTenant(outOrgId)
		//if err != nil {
		//	log.Error(err)
		//	return err
		//}
		msgCard := card.GetFsCardProjectSettings(orgId, chatId)
		errSys := card.PushCard(orgId, &commonvo.PushCard{
			OrgId:         orgId,
			OutOrgId:      outOrgId,
			SourceChannel: sourceChannel,
			ChatIds:       []string{chatId},
			CardMsg:       msgCard,
		})
		if errSys != nil {
			log.Errorf("[GroupChatReplyToUserInsProjectSettings] err:%v", errSys)
			return errSys
		}
		//fsResp, oriErr := tenant.SendMessage(fsSdkVo.MsgVo{
		//	ChatId:  chatId,
		//	MsgType: "interactive",
		//	Card:    msgCard,
		//})
		//if oriErr != nil {
		//	log.Error(oriErr)
		//	return err
		//}
		//if fsResp.Code != 0 {
		//	log.Error("发送消息异常")
		//}
		//log.Infof("飞书群聊-用户指令处理-HandleGroupChatUserInsProjectSettings-响应用户调用结果 code: %v", fsResp.Code)
		return nil
	case sdk_const.SourceChannelDingTalk:
		// 钉钉发送项目动态设置卡片
	}

	return nil
}
