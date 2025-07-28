package callsvc

import (
	"gitea.bjx.cloud/allstar/feishu-sdk-golang/core/model/vo"
	sdk_const "gitea.bjx.cloud/allstar/platform-sdk/consts"
	"github.com/star-table/startable-server/app/facade/orgfacade"
	"github.com/star-table/startable-server/common/core/consts"
	"github.com/star-table/startable-server/common/core/errs"
	"github.com/star-table/startable-server/common/core/util/json"
	"github.com/star-table/startable-server/common/core/util/temp"
	"github.com/star-table/startable-server/common/extra/card"
	"github.com/star-table/startable-server/common/extra/feishu"
	"github.com/star-table/startable-server/common/model/vo/callvo"
	"github.com/star-table/startable-server/common/model/vo/commonvo"
	"github.com/star-table/startable-server/common/model/vo/orgvo"
)

type P2pChatCreateHandler struct{}

type P2pChatCreateReq struct {
	Ts    string `json:"ts"`
	Uuid  string `json:"uuid"`
	Token string `json:"token"`
	Type  string `json:"type"`

	Event P2pChatCreateReqData `json:"event"`
}

type P2pChatCreateReqData struct {
	AppId     string `json:"app_id"`
	ChatId    string `json:"chat_id"`
	TenantKey string `json:"tenant_key"`
	Type      string `json:"type"`

	Operator P2pChatCreateOperator `json:"operator"`
	User     P2pChatCreateUser     `json:"user"`
}

type P2pChatCreateOperator struct {
	OpenId string `json:"open_id"`
	UserId string `json:"user_id"`
}

type P2pChatCreateUser struct {
	OpenId string `json:"open_id"`
	UserId string `json:"user_id"`
	Name   string `json:"name"`
}

// 飞书处理
func (P2pChatCreateHandler) Handle(data string) (string, errs.SystemErrorInfo) {
	p2pChatCreateReq := &P2pChatCreateReq{}
	_ = json.FromJson(data, p2pChatCreateReq)

	log.Infof("[Handle] 用户和机器人的会话首次被创建 data: %s", json.ToJsonIgnoreError(data))

	tenantKey := p2pChatCreateReq.Event.TenantKey
	user := p2pChatCreateReq.Event.User
	chatId := p2pChatCreateReq.Event.ChatId

	log.Infof("[Handle] tenantKey %s, user %s, chatId %s", tenantKey, json.ToJsonIgnoreError(user), chatId)

	resp := orgfacade.JudgeUserIsAdmin(orgvo.JudgeUserIsAdminReq{
		SourceChannel: sdk_const.SourceChannelFeishu,
		OutOrgId:      tenantKey,
		OutUserId:     user.OpenId,
	})
	if resp.Failure() {
		log.Errorf("[Handle] orgfacade.JudgeUserIsAdmin err: %v", resp.Error())
	} else if !resp.IsTrue {
		isOrgCreator, err := CheckIsOrgCreator(user.OpenId, tenantKey)
		if err != nil {
			log.Errorf("[Handle] tenantKey: %s, chatId: %s, CheckIsOrgCreator err: %v", tenantKey, chatId, err)
			return "error", err
		}

		// err := PushBotHelpNotice(tenantKey, chatId, "")
		// 首次打开机器人的成员，如果不是组织创建者，才能发送帮助卡片
		if !isOrgCreator {
			err := PushWelcomeMsgToOrgMember(tenantKey, chatId)
			if err != nil {
				log.Errorf("[Handle] tenantKey: %s, chatId: %s, PushWelcomeMsgToOrgMember err: %v", tenantKey, chatId, err)
				return "error", err
			}
		}
	}

	return "ok", nil
}

// CheckIsOrgCreator 检查 curOpenId 对应的是否是组织创建人
func CheckIsOrgCreator(curOpenId, tenantKey string) (bool, errs.SystemErrorInfo) {
	isCreator := false
	orgResp := orgfacade.GetBaseOrgInfoByOutOrgId(orgvo.GetBaseOrgInfoByOutOrgIdReqVo{OutOrgId: tenantKey})
	if orgResp.Failure() {
		log.Errorf("[Handle] CheckIsOrgCreator err: %v", orgResp.Error())
		return isCreator, orgResp.Error()
	}
	opUserResp := orgfacade.GetUserId(orgvo.GetUserIdReqVo{
		SourceChannel: sdk_const.SourceChannelFeishu,
		EmpId:         curOpenId,
		CorpId:        tenantKey,
		OrgId:         orgResp.BaseOrgInfo.OrgId,
	})
	if opUserResp.Failure() {
		log.Errorf("[Handle] orgfacade.GetUserId err: %v", orgResp.Error())
		return isCreator, opUserResp.Error()
	}
	isCreator = orgResp.BaseOrgInfo.Creator == opUserResp.GetUserId.UserID

	return isCreator, nil
}

func PushWelcomeMsgToOrgMember(tenantKey string, chatId string, openIds ...string) errs.SystemErrorInfo {
	tenantClient, err := feishu.GetTenant(tenantKey)
	if err != nil {
		log.Error(err)
		return err
	}
	orgResp := orgfacade.GetBaseOrgInfoByOutOrgId(orgvo.GetBaseOrgInfoByOutOrgIdReqVo{OutOrgId: tenantKey})
	if orgResp.Failure() {
		return orgResp.Error()
	}
	cardInfo := GetCardInfoForFirst(orgResp.BaseOrgInfo.OrgId)
	card := cardInfo.Card

	var resp *vo.MsgResp
	var err1 error

	if chatId != "" {
		resp, err1 = tenantClient.SendMessage(vo.MsgVo{
			ChatId:  chatId,
			MsgType: "interactive",
			Card:    card,
		})
	} else {
		resp, err1 = tenantClient.SendMessageBatch(vo.BatchMsgVo{
			OpenIds: openIds,
			MsgType: "interactive",
			Card:    card,
		})
	}
	if err1 != nil {
		log.Errorf("[PushWelcomeMsgToOrgMember] chatId: %s, openIds: %s, SendMessage/SendMessageBatch err: %v", chatId, json.ToJsonIgnoreError(openIds), err1)
		return errs.FeiShuOpenApiCallError
	}
	if resp.Code != 0 {
		log.Errorf("发送机器人欢迎通知失败，错误信息: %s, chatId: %s, openIds: %s", resp.Msg, chatId, json.ToJsonIgnoreError(openIds))
		return errs.BuildSystemErrorInfo(errs.FeiShuOpenApiCallError)
	}
	log.Infof("发送机器人欢迎通知成功")

	return nil
}

// PushBotHelpNotice 给**个人用户**响应帮助卡片
func PushBotHelpNotice(tenantKey, chatId string, openIds ...string) errs.SystemErrorInfo {
	log.Infof("[PushBotHelpNotice] tenantKey:%s, chatId:%s, openIds:%v", tenantKey, chatId, openIds)
	tenantClient, err := feishu.GetTenant(tenantKey)
	if err != nil {
		log.Error(err)
		return err
	}
	orgResp := orgfacade.GetBaseOrgInfoByOutOrgId(orgvo.GetBaseOrgInfoByOutOrgIdReqVo{OutOrgId: tenantKey})
	if orgResp.Failure() {
		return orgResp.Error()
	}
	orgId := orgResp.BaseOrgInfo.OrgId
	card := &vo.Card{
		Header: &vo.CardHeader{
			Title: &vo.CardHeaderTitle{
				Tag:     "plain_text",
				Content: "欢迎使用 极星 协作应用",
			},
		},
		Elements: []interface{}{
			vo.CardElementContentModule{
				Tag: "div",
				Fields: []vo.CardElementField{
					{
						Text: vo.CardElementText{
							Tag:     "lark_md",
							Content: "极星是一个高效而稳定的项目协作平台",
						},
					},
				},
			},
			vo.CardElementContentModule{
				Tag: "div",
				Fields: []vo.CardElementField{
					{
						Text: vo.CardElementText{
							Tag:     "lark_md",
							Content: "移动端：[点击进入]($urlVal)",
							Href: map[string]vo.CardElementUrl{
								"urlVal": vo.CardElementUrl{
									Url:        feishu.GetDefaultAppLink(orgId),
									AndroidUrl: feishu.GetDefaultAppLink(orgId),
									IosUrl:     feishu.GetDefaultAppLink(orgId),
									PcUrl:      feishu.GetMobileOpenQRCodeAppLink(orgId),
								},
							},
						},
					},
				},
			},
			vo.CardElementContentModule{
				Tag: "div",
				Fields: []vo.CardElementField{
					{
						Text: vo.CardElementText{
							Tag:     "lark_md",
							Content: "PC端：[点击进入](" + feishu.GetAppLinkPcWelcomeUrl(orgId) + ")",
							//Href: map[string]vo.CardElementUrl{
							//	"urlVal":vo.CardElementUrl{
							//		Url:feishu.GetDefaultAppLink(),
							//		AndroidUrl:feishu.GetDefaultAppLink(),
							//		IosUrl:feishu.GetDefaultAppLink(),
							//		PcUrl:feishu.GetMobileOpenQRCodeAppLink(),
							//	},
							//},
						},
					},
				},
			},
			vo.CardElementContentModule{
				Tag: "div",
				Fields: []vo.CardElementField{
					{
						Text: vo.CardElementText{
							Tag:     "lark_md",
							Content: "您可以通过@极星协作机器人进行快捷操作",
						},
					},
				},
			},
			vo.CardElementContentModule{
				Tag: "div",
				Fields: []vo.CardElementField{
					{
						Text: vo.CardElementText{
							Tag:     "lark_md",
							Content: "发送** create 名称 **可以快速创建任务",
						},
					},
				},
			},
			//vo.CardElementContentModule{
			//	Tag: "div",
			//	Fields: []vo.CardElementField{
			//		{
			//			Text: vo.CardElementText{
			//				Tag:     "lark_md",
			//				Content: "发送** settings **可以进入配置修改页面",
			//			},
			//		},
			//	},
			//},
			vo.CardElementContentModule{
				Tag: "div",
				Fields: []vo.CardElementField{
					{
						Text: vo.CardElementText{
							Tag:     "lark_md",
							Content: "发送** notification **可以进入通知设置修改页面",
						},
					},
				},
			},
			vo.CardElementContentModule{
				Tag: "div",
				Fields: []vo.CardElementField{
					{
						Text: vo.CardElementText{
							Tag:     "lark_md",
							Content: "在使用中有任何问题，均可以直接联系我们。【联系电话：18001735738】",
						},
					},
				},
			},
			vo.CardElementContentModule{
				Tag: "div",
				Fields: []vo.CardElementField{
					{
						Text: vo.CardElementText{
							Tag:     "lark_md",
							Content: "详细使用说明 **[帮助手册](" + feishu.AppGuide + ")**",
						},
					},
				},
			},
			vo.CardElementContentModule{
				Tag: "div",
				Fields: []vo.CardElementField{
					{
						Text: vo.CardElementText{
							Tag:     "lark_md",
							Content: "可以随时联系 **[在线客服](" + feishu.AppCustomerService + ")**",
						},
					},
				},
			},
		},
	}
	var resp *vo.MsgResp
	var err1 error

	if chatId != "" {
		resp, err1 = tenantClient.SendMessage(vo.MsgVo{
			ChatId:  chatId,
			MsgType: "interactive",
			Card:    card,
		})
	} else {
		resp, err1 = tenantClient.SendMessageBatch(vo.BatchMsgVo{
			OpenIds: openIds,
			MsgType: "interactive",
			Card:    card,
		})
	}

	if err1 != nil {
		log.Error(err1)
		return errs.FeiShuOpenApiCallError
	}
	if resp.Code != 0 {
		log.Errorf("发送机器人欢迎通知失败，错误信息 %s", resp.Msg)
		return errs.BuildSystemErrorInfo(errs.FeiShuOpenApiCallError)
	}
	log.Infof("发送机器人欢迎通知成功")

	return nil
}

// 飞书群聊：chatType 为 group 的群聊，并且群聊未关联项目的帮助信息卡片。
func PushBotHelpNoticeForNoProGroup(tenantKey, chatId string, openIds ...string) errs.SystemErrorInfo {
	tenantClient, err := feishu.GetTenant(tenantKey)
	if err != nil {
		log.Error(err)
		return err
	}
	card := &vo.Card{
		Header: &vo.CardHeader{
			Title: &vo.CardHeaderTitle{
				Tag:     "plain_text",
				Content: "欢迎使用 极星 协作应用",
			},
		},
		Elements: []interface{}{
			vo.CardElementContentModule{
				Tag: "div",
				Fields: []vo.CardElementField{
					{
						Text: vo.CardElementText{
							Tag:     "lark_md",
							Content: "您可以通过**@极星协作**进行快捷操作",
						},
					},
				},
			},
			vo.CardElementContentModule{
				Tag: "div",
				Fields: []vo.CardElementField{
					{
						Text: vo.CardElementText{
							Tag:     "lark_md",
							Content: "发送** create 名称 **可以快速创建任务",
						},
					},
				},
			},
			//vo.CardElementContentModule{
			//	Tag: "div",
			//	Fields: []vo.CardElementField{
			//		{
			//			Text: vo.CardElementText{
			//				Tag:     "lark_md",
			//				Content: "发送** settings **可以进入配置修改页面",
			//			},
			//		},
			//	},
			//},
			vo.CardElementContentModule{
				Tag: "div",
				Fields: []vo.CardElementField{
					{
						Text: vo.CardElementText{
							Tag:     "lark_md",
							Content: "发送** notification **可以进入通知设置修改页面",
						},
					},
				},
			},
		},
	}
	var resp *vo.MsgResp
	var err1 error

	if chatId != "" {
		resp, err1 = tenantClient.SendMessage(vo.MsgVo{
			ChatId:  chatId,
			MsgType: "interactive",
			Card:    card,
		})
	} else {
		resp, err1 = tenantClient.SendMessageBatch(vo.BatchMsgVo{
			OpenIds: openIds,
			MsgType: "interactive",
			Card:    card,
		})
	}

	if err1 != nil {
		log.Error(err1)
		return errs.FeiShuOpenApiCallError
	}
	if resp.Code != 0 {
		log.Errorf("发送机器人欢迎通知失败，错误信息 %s", resp.Msg)
		return errs.BuildSystemErrorInfo(errs.FeiShuOpenApiCallError)
	}
	log.Infof("发送机器人欢迎通知成功")

	return nil
}

// 根据参数渲染出一个 at 某个人的字符串
func RenderAtSomeoneStr(openId, userName string) string {
	tmp := `<at id="{{.OpenId}}">@{{.UserName}}</at>`
	params := map[string]interface{}{
		"OpenId":   openId,
		"UserName": userName,
	}
	result, err := temp.Render(tmp, params)
	if err != nil {
		log.Error(err)
		return ""
	}
	return result
}

// PushBotHelpNoticeForGroup 飞书群聊 - chatType 为 group 的群聊帮助信息卡片。
func PushBotHelpNoticeForGroup(orgId int64, tenantKey string, event callvo.MessageReqData, bindProNum int,
	openIds ...string) errs.SystemErrorInfo {
	hasBindMultiPro := bindProNum > 1
	chatId := event.OpenChatId
	//tenantClient, err := feishu.GetTenant(tenantKey)
	//if err != nil {
	//	log.Errorf("[PushBotHelpNoticeForGroup] err: %v, tenantKey: %s", err, tenantKey)
	//	return err
	//}
	// 通过操作人的 openId 获取其姓名
	opUserResp := orgfacade.GetBaseUserInfoByEmpId(orgvo.GetBaseUserInfoByEmpIdReqVo{
		OrgId: orgId,
		EmpId: event.OpenId,
	})
	if opUserResp.Failure() {
		log.Error(opUserResp.Error())
		return opUserResp.Error()
	}
	opUser := opUserResp.BaseUserInfo
	atSomeoneUserStr := RenderAtSomeoneStr(event.OpenId, opUser.Name)
	cardMsg := card.GetFsCardMsgForBindPro(atSomeoneUserStr + consts.FsCardNotUnderStand + consts.FsCardHelpInfoForInstructionHelp)
	if hasBindMultiPro {
		// 如果群聊绑定了多个项目，则推送另一种帮助卡片信息
		cardMsg = card.GetFsCardMsgForBindPro(atSomeoneUserStr + consts.FsCardBindMultiPro)
	}

	errSys := card.PushCard(orgId, &commonvo.PushCard{
		OrgId:         orgId,
		OutOrgId:      tenantKey,
		SourceChannel: sdk_const.SourceChannelFeishu,
		OpenIds:       openIds,
		ChatIds:       []string{chatId},
		CardMsg:       cardMsg,
	})
	if errSys != nil {
		log.Errorf("[HandleGroupChatAtUserName] err:%v", errSys)
		return errSys
	}

	//var resp *vo.MsgResp
	//var err1 error
	//
	//if chatId != "" {
	//	resp, err1 = tenantClient.SendMessage(vo.MsgVo{
	//		ChatId:  chatId,
	//		MsgType: "interactive",
	//		Card:    card,
	//	})
	//} else {
	//	resp, err1 = tenantClient.SendMessageBatch(vo.BatchMsgVo{
	//		OpenIds: openIds,
	//		MsgType: "interactive",
	//		Card:    card,
	//	})
	//}
	//
	//if err1 != nil {
	//	log.Error(err1)
	//	return errs.FeiShuOpenApiCallError
	//}
	//if resp.Code != 0 {
	//	log.Errorf("发送机器人欢迎通知失败，错误信息 %s", resp.Msg)
	//	return errs.BuildSystemErrorInfo(errs.FeiShuOpenApiCallError)
	//}
	log.Infof("发送机器人欢迎通知成功")

	return nil
}

//func GetFsCardMsgForBindOnePro(atSomeoneUserStr string) *vo.Card {
//	card := &vo.Card{
//		Elements: []interface{}{
//			//vo.CardElementContentModule{
//			//	Tag: "div",
//			//	Fields: []vo.CardElementField{
//			//		{
//			//			Text: vo.CardElementText{
//			//				Tag:     "lark_md",
//			//				Content: atSomeoneUserStr + " 我没有理解你的意思，你可以通过以下方式与项目交互：",
//			//			},
//			//		},
//			//	},
//			//},
//			vo.CardElementContentModule{
//				Tag: "div",
//				Fields: []vo.CardElementField{
//					{
//						Text: vo.CardElementText{
//							Tag:     "lark_md",
//							Content: atSomeoneUserStr + " 我没有理解你的意思，你可以通过以下方式与项目交互：\n" + consts.FsCardHelpInfoForInstructionHelp,
//						},
//					},
//				},
//			},
//		},
//	}
//
//	return card
//}

//func GetFsCardMsgForBindMultiPro(atSomeoneUserStr string) *vo.Card {
//	card := &vo.Card{
//		Elements: []interface{}{
//			vo.CardElementContentModule{
//				Tag: "div",
//				Fields: []vo.CardElementField{
//					{
//						Text: vo.CardElementText{
//							Tag: "lark_md",
//							Content: atSomeoneUserStr + "**已将2个项目关联到本群聊**\n" +
//								" 您关联的所有项目都将推送动态至本群聊。\n" +
//								"您还可以发送 @极星协作 设置 ，重新进行推送设置。",
//						},
//					},
//				},
//			},
//		},
//	}
//
//	return card
//}

// 在已有群聊手动添加机器人，推送相关的消息卡片(没有绑定项目的)
func PushBotNoticeNotBindProject(tenantKey string, chatId string) errs.SystemErrorInfo {
	orgInfoResp := orgfacade.GetBaseOrgInfoByOutOrgId(orgvo.GetBaseOrgInfoByOutOrgIdReqVo{OutOrgId: tenantKey})
	if orgInfoResp.Failure() {
		log.Errorf("[PushBotNoticeNoBindProject] orgfacade.GetBaseOrgInfoByOutOrgId err:%v, tenantKey:%s, chatId:%s",
			orgInfoResp.Message, tenantKey, chatId)
		return orgInfoResp.Error()
	}

	orgId := orgInfoResp.BaseOrgInfo.OrgId

	//tenantClient, err := feishu.GetTenant(tenantKey)
	//if err != nil {
	//	log.Errorf("[PushBotNoticeNoBindProject] err: %v, tenantKey: %s", err, tenantKey)
	//	return errs.FeiShuClientTenantError
	//}

	cardMsg := card.GetFsCardForNotBindProjectMsg(orgId, chatId)
	errSys := card.PushCard(orgId, &commonvo.PushCard{
		OrgId:         orgId,
		OutOrgId:      tenantKey,
		SourceChannel: orgInfoResp.BaseOrgInfo.SourceChannel,
		ChatIds:       []string{chatId},
		CardMsg:       cardMsg,
	})
	if errSys != nil {
		log.Errorf("[PushBotNoticeNotBindProject] err:%v", errSys)
		return errSys
	}
	//resp, err1 := tenantClient.SendMessage(vo.MsgVo{
	//	ChatId:  chatId,
	//	MsgType: "interactive",
	//	Card:    card,
	//})
	//
	//if err1 != nil {
	//	log.Errorf("[PushBotNoticeNoBindProject] chatId: %s, SendMessage err: %v", chatId, err1)
	//	return errs.FeiShuOpenApiCallError
	//}
	//if resp.Code != 0 {
	//	log.Errorf("手动添加机器人 发送欢迎通知失败，错误信息: %s, chatId: %s", resp.Msg, chatId)
	//	return errs.BuildSystemErrorInfo(errs.FeiShuOpenApiCallError)
	//}
	log.Infof("手动添加机器人 发送欢迎通知成功")

	return nil
}

//func GetFsCardForNotBindProjectMsg(orgId int64, chatId string) *vo.Card {
//	card := &vo.Card{
//		Elements: []interface{}{
//			vo.CardElementContentModule{
//				Tag: "div",
//				Fields: []vo.CardElementField{
//					{
//						Text: vo.CardElementText{
//							Tag: "lark_md",
//							Content: "**Hi~您好，我是极星协作的机器人。**\n" +
//								"**我可以为您推送关注的项目动态。快来对我设置吧^ ^ ~**" + "",
//						},
//					},
//				},
//			},
//			vo.CardElementActionModule{
//				Tag: "action",
//				Actions: []interface{}{
//					vo.ActionButton{
//						Tag: "button",
//						Text: vo.CardElementText{
//							Tag:     "plain_text",
//							Content: "项目动态设置",
//						},
//						Url:  feishu.GetGroupChatSettingBtnAppLink(orgId, chatId),
//						Type: "default",
//					},
//				},
//			},
//		},
//	}
//	return card
//}

//// SendCardMsgBindOneProject 手动添加机器人 关联一个项目时发送的卡片信息
//func SendCardMsgBindOneProject(tenantKey string, chatId string) errs.SystemErrorInfo {
//	card := &vo.Card{
//		Elements: []interface{}{
//			vo.CardElementContentModule{
//				Tag: "div",
//				Fields: []vo.CardElementField{
//					{
//						Text: vo.CardElementText{
//							Tag: "lark_md",
//							Content: "**已将1个项目关联到本群聊**\n" +
//								"您关联项目的消息将推送到本群聊，您还可以通过以下方式与项目交互，例如：\n" + consts.FsCardHelpInfoForInstructionHelp,
//						},
//					},
//				},
//			},
//		},
//	}
//	tenantClient, err := feishu.GetTenant(tenantKey)
//	if err != nil {
//		log.Errorf("[SendCardMsgBindOneProject] err: %v, tenantKey: %s", err, tenantKey)
//		return errs.FeiShuClientTenantError
//	}
//	resp, err1 := tenantClient.SendMessage(vo.MsgVo{
//		ChatId:  chatId,
//		MsgType: "interactive",
//		Card:    card,
//	})
//	if err1 != nil {
//		log.Errorf("[SendCardMsgBindOneProject] chatId: %s, SendMessage err: %v", chatId, err1)
//		return errs.FeiShuOpenApiCallError
//	}
//	if resp.Code != 0 {
//		log.Errorf("发送消息失败，错误信息: %s, chatId: %s", resp.Msg, chatId)
//		return errs.BuildSystemErrorInfo(errs.FeiShuOpenApiCallError)
//	}
//
//	return nil
//}
