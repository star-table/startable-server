package callsvc

import (
	"strings"

	"github.com/star-table/startable-server/common/core/util/json"
	sdk_const "gitea.bjx.cloud/allstar/platform-sdk/consts"
	"github.com/star-table/startable-server/common/core/errs"
	"github.com/star-table/startable-server/common/extra/card"
	"github.com/star-table/startable-server/common/model/vo/callvo"
	"github.com/star-table/startable-server/common/model/vo/orgvo"
	"github.com/star-table/startable-server/common/model/vo/projectvo"
	"github.com/star-table/startable-server/app/facade/orgfacade"
	"github.com/star-table/startable-server/app/facade/projectfacade"
	callconsts "github.com/star-table/startable-server/app/service/callsvc/consts"
	"github.com/star-table/startable-server/app/service/callsvc/domain"
)
)

//消息分隔符
const (
	MessageSeparator = " "
)

//消息类型
const (
	MessageTypeText = "text"
)

type MessageHandler struct{}

type MessageReq struct {
	Ts    string `json:"ts"`
	Uuid  string `json:"uuid"`
	Token string `json:"token"`
	Type  string `json:"type"`

	Event callvo.MessageReqData `json:"event"`
}

// 飞书群聊，机器人消息和会话事件处理配置，用于 private bot 聊天或手动拉群群聊。
var instructHandlerMap = map[string]func(tenantKey string, event callvo.MessageReqData, text string, extraInfo map[string]interface{}) error{
	callconsts.InstructCreateIssue: instructCreateIssueFunc,
	//InstructSettings:                 instructSettingsFunc,
	callconsts.InstructNotification: instructNotification,
}

// 飞书群聊，机器人消息和会话事件处理配置，用于**项目群聊**。
var instructHandlerMapForGroup = map[string]func(tenantKey string, event callvo.MessageReqData, text string, extraInfo map[string]interface{}) error{
	callconsts.InstructUserInsProIssue:          instructUserInsProIssueFunc,
	callconsts.InstructUserInsProProgress:       instructUserInsProProgressFunc,
	callconsts.InstructAtUserNameWithIssueTitle: instructAtUserNameWithIssueTitleFunc,
	callconsts.InstructAtUserName:               instructAtUserNameFunc,
	callconsts.InstructUserInsProDynSetting:     instructUserInsProjectSettingsFunc,
}

// Handle 对用户主动发给应用或者群内@机器人的消息做对应的处理，目前只处理text消息
// doc:https://open.feishu.cn/document/uQjL04CN/uUTNz4SN1MjL1UzM#%E6%8E%A5%E6%94%B6%E6%B6%88%E6%81%AF
func (MessageHandler) Handle(data string) (string, errs.SystemErrorInfo) {
	messageReq := &MessageReq{}
	_ = json.FromJson(data, messageReq)

	log.Infof("用户主动发消息到机器人 %s", data)

	event := messageReq.Event
	tenantKey := event.TenantKey
	openId := event.OpenId

	log.Infof("tenantKey %s, msgType %s, openId %s", tenantKey, event.MsgType, openId)
	extraParam := make(map[string]interface{}, 0)
	if event.MsgType == MessageTypeText {
		if event.ChatType == "private" {
			text := event.TextWithoutAtBot
			handleErr := handleMessageInstruct(tenantKey, event, text, extraParam)
			if handleErr != nil {
				log.Error(handleErr)
				return "error", nil
			}
		} else if event.ChatType == "group" {
			if event.IsMention {
				text := event.TextWithoutAtBot
				handleErr := handleMessageInstructForGroup(tenantKey, event, text, extraParam)
				if handleErr != nil {
					log.Error(handleErr)
					return "error", nil
				}
			}
		}
	}

	return "ok", nil
}

// 处理 chatType 为 private 的 bot 群聊消息，或者**未关联项目**的群聊消息。
func handleMessageInstruct(tenantKey string, event callvo.MessageReqData, msg string, extraInfo map[string]interface{}) error {
	msg = strings.Trim(msg, " ")
	instructs := strings.Split(msg, MessageSeparator)

	//nco.2019-11-28: 目前只处理创建任务指令和主动设置指令
	// eg: create 这是一个任务
	// eg: settings
	if handler, ok := instructHandlerMap[strings.ToLower(instructs[0])]; ok {
		err := handler(tenantKey, event, msg, extraInfo)
		if err != nil {
			log.Error(err)
			return err
		}
	} else {
		//欢迎使用极星协作
		//您可以通过@极星协作机器人发送 create指令快速创建任务
		//例：
		//@极星协作机器人  create 这是一条新任务
		err := sendSimpleNotice(tenantKey, event, extraInfo)
		if err != nil {
			log.Error(err)
			return err
		}
	}

	return nil
}

// 处理 chatType 为 group 的 bot 群聊消息，并且群聊已经关联了项目。
func handleMessageInstructForGroup(tenantKey string, event callvo.MessageReqData, msg string, extraInfo map[string]interface{}) error {
	msg = strings.Trim(msg, " ")
	// 通过群聊 id 查询该群关联了多少个项目。
	// 关联一个和关联多个的交互会有所不同
	// 通过 chatId 查询是否关联了项目，如无法匹配到对应的项目时，则回复用户“简化版”通用帮助信息卡片。
	dealInfo, err := domain.GetGroupChatHandleInfos(tenantKey, &event, sdk_const.SourceChannelFeishu)
	if err != nil {
		log.Errorf("[handleMessageInstructForGroup] err: %v, tenantKey: %s", err, tenantKey)
		return err
	}
	if dealInfo.IsIssueChat {
		// 任务群聊，则不支持群指令，返回卡片提示
		if err := card.SendCardFeiShuReplyIssueChatIns(event.OpenChatId, dealInfo); err != nil {
			log.Errorf("[handleMessageInstructForGroup] SendCardFeiShuReplyIssueChatIns: %v", err)
			return err
		}
		return nil
	}
	hasBindMultiPro := len(dealInfo.BindProjectIds) > 1
	extraInfo["hasBindMultiPro"] = hasBindMultiPro
	// 如果是群聊，并且没有关联项目，则执行 `handleMessageInstruct`
	if len(dealInfo.BindProjectIds) == 0 {
		// 没有关联项目的群聊，需将 isSimpleHelperInfo 置为 true。表示群聊中得指令帮助信息（简化版）
		extraInfo["isSimpleHelperInfo"] = true
		errSys := PushBotNoticeForNotSupportChatInstruction(tenantKey, &event, dealInfo)
		if errSys != nil {
			log.Errorf("[handleMessageInstructForGroup] err:%v, tenantKey: %s", errSys, tenantKey)
			return errSys
		}
		return nil
	}
	extraInfo["dealInfo"] = dealInfo

	// 匹配指令：项目任务|项目进展
	if handler, ok := instructHandlerMapForGroup[msg]; ok {
		err := handler(tenantKey, event, msg, extraInfo)
		if err != nil {
			log.Errorf("[handleMessageInstructForGroup] err: %v", err)
			return err
		}
	} else {
		// 尝试解析以下指令：
		// 1.@负责人姓名
		// 2.@负责人姓名 任务标题
		insText := ParseInsForUserName(msg, extraInfo)
		if handler, ok := instructHandlerMapForGroup[insText]; ok {
			err := handler(tenantKey, event, msg, extraInfo)
			if err != nil {
				log.Error(err)
				return err
			}
		} else {
			//若没有匹配到指令，则回复用户**帮助信息**卡片
			if len(dealInfo.BindProjectIds) == 0 {
				// 如果绑定了0个项目
				errSys := PushBotNoticeForNotSupportChatInstruction(tenantKey, &event, dealInfo)
				if errSys != nil {
					log.Errorf("[handleMessageInstructForGroup] err:%v", errSys)
					return errSys
				}
				return nil
			} else if len(dealInfo.BindProjectIds) == 1 {
				// 关联一个项目
				errSys := PushBotHelpNoticeForGroup(dealInfo.OrgInfo.OrgId, tenantKey, event, len(dealInfo.BindProjectIds))
				if errSys != nil {
					log.Errorf("[handleMessageInstructForGroup] err:%v", errSys)
					return errSys
				}
				return nil
			} else if len(dealInfo.BindProjectIds) >= 2 {
				errSys := PushBotNoticeForNotSupportChatInstruction(tenantKey, &event, dealInfo)
				if errSys != nil {
					log.Errorf("[handleMessageInstructForGroup] err:%v", errSys)
					return errSys
				}
				return nil
			}
			err := sendSimpleNoticeForGroup(tenantKey, event, len(dealInfo.BindProjectIds))
			if err != nil {
				log.Errorf("[handleMessageInstructForGroup] err: %v", err)
				return err
			}
		}
	}

	return nil
}

// 推送“详细版”通用帮助信息卡片。将 chatType 为 private 的群聊帮助信息，响应给用户。
func sendSimpleNotice(tenantKey string, event callvo.MessageReqData, extraInfo map[string]interface{}) errs.SystemErrorInfo {
	log.Infof("[sendSimpleNotice] eventData:%s", json.ToJsonIgnoreError(event))
	//isSimpleHelperInfo := false
	//if tmpVal, ok := extraInfo["isSimpleHelperInfo"]; ok {
	//	isSimpleHelperInfo = tmpVal.(bool)
	//}
	//var err1 errs.SystemErrorInfo
	//if isSimpleHelperInfo {
	//	err1 = PushBotHelpNoticeForNoProGroup(tenantKey, event.OpenChatId, "")
	//} else {
	//	err1 = PushBotHelpNotice(tenantKey, event.OpenChatId, "")
	//}
	//if err1 != nil {
	//	log.Error(err1)
	//	return errs.FeiShuOpenApiCallError
	//}
	orgResp := orgfacade.GetBaseOrgInfoByOutOrgId(orgvo.GetBaseOrgInfoByOutOrgIdReqVo{OutOrgId: tenantKey})
	if orgResp.Failure() {
		return orgResp.Error()
	}
	orgId := orgResp.BaseOrgInfo.OrgId
	err := card.SendCardFeiShuBotNotKnown(tenantKey, orgId, event.OpenChatId, event.OpenId)
	if err != nil {
		log.Errorf("[sendSimpleNotice] err:%v, tenantKey:%s, chatId:%s", err, tenantKey, event.OpenChatId)
		return err
	}
	return nil
}

// 推送“简化版”通用帮助信息卡片。将 chatType 为 group 的群聊，但群聊未关联上项目帮助信息，响应给用户。
func sendSimpleNoticeForNoProGroup(tenantKey string, event callvo.MessageReqData) errs.SystemErrorInfo {
	err1 := PushBotHelpNoticeForNoProGroup(tenantKey, event.OpenChatId, "")
	if err1 != nil {
		log.Error(err1)
		return errs.FeiShuOpenApiCallError
	}
	return nil
}

// 将 chatType 为 group 的群聊帮助信息，响应给用户。
func sendSimpleNoticeForGroup(tenantKey string, event callvo.MessageReqData, bindProNum int) errs.SystemErrorInfo {
	// 通过 event.OpenChatId 获取对应项目的 id
	resp := projectfacade.GetProjectIdByChatId(projectvo.GetProjectIdByChatIdReqVo{
		OpenChatId: event.OpenChatId,
	})
	if resp.Failure() {
		log.Errorf("[sendSimpleNoticeForGroup] err: %v, tenantKey: %s", resp.Error(), tenantKey)
		return resp.Error()
	}
	err1 := PushBotHelpNoticeForGroup(resp.Data.OrgId, tenantKey, event, bindProNum)
	if err1 != nil {
		log.Error(err1)
		return errs.FeiShuOpenApiCallError
	}

	return nil
}

// 解析 group 群聊的特殊指令：用户名相关
func ParseInsForUserName(text string, extraInfo map[string]interface{}) (insText string) {
	// 指令名
	insText = ""
	text = strings.Trim(text, " ")
	log.Infof("ParseInsForUserName text: %v", text)
	// 1.@负责人姓名：<at open_id="yyy">@张三</at>
	if domain.CheckIsUserName(text) {
		insText = callconsts.InstructAtUserName
		atUserOpenId, atUserName := domain.MatchUserInsUserInfo(text)
		extraInfo["atUserName"] = atUserName
		extraInfo["atUserOpenId"] = atUserOpenId
	} else if domain.CheckIsUserNameWithIssueTitle(text) {
		// 2.@负责人姓名 任务标题1：<at open_id="yyy">@张三</at> 任务标题001
		insText = callconsts.InstructAtUserNameWithIssueTitle
		atUserName, atUserOpenId, issueTitle := domain.MatchUserInsUserNameWithIssueTitle(text)
		extraInfo["atUserName"] = atUserName
		extraInfo["atUserOpenId"] = atUserOpenId
		extraInfo["issueTitle"] = issueTitle
	} else {
		// 3.不支持的指令
		insText = "Unknown"
	}
	return insText
}
