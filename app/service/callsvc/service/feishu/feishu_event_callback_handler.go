package callsvc

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/star-table/startable-server/app/facade/msgfacade"
	"github.com/star-table/startable-server/common/core/config"
	"github.com/star-table/startable-server/common/core/consts"
	"github.com/star-table/startable-server/common/core/errs"
	"github.com/star-table/startable-server/common/core/logger"
	"github.com/star-table/startable-server/common/core/model"
	"github.com/star-table/startable-server/common/core/threadlocal"
	"github.com/star-table/startable-server/common/core/util/asyn"
	"github.com/star-table/startable-server/common/core/util/encrypt"
	"github.com/star-table/startable-server/common/core/util/json"
	"github.com/star-table/startable-server/common/core/util/strs"
	"github.com/star-table/startable-server/common/model/bo"
	"github.com/star-table/startable-server/common/model/vo/msgvo"
	"github.com/gin-gonic/gin"
	"github.com/jtolds/gls"
)

// event
const (
	// 飞书新版事件回调
	EventChatBotEntry          = "im.chat.member.bot.added_v1"      // 机器人进群
	EventAddUserToChat         = "im.chat.member.user.added_v1"     //用户进群
	EventRemoveUserFromChat    = "im.chat.member.user.deleted_v1"   //用户被移除群
	EventRevokeAddUserFromChat = "im.chat.member.user.withdrawn_v1" //撤销拉用户进群
	EventImChatDisabledV1      = "im.chat.disbanded_v1"             // 群解散

	// v3 版本的通讯录授权返回发生变更：https://open.feishu.cn/document/uAjLw4CM/ukTMukTMukTM/reference/contact-v3/scope/events/updated#52f59572
	EventContactScopeChangeV3 = "contact.scope.updated_v3"

	// 飞书旧版事件
	//EventAddBot = "add_bot" // 现在应该不用这个了
	// 授权范围变更  https://open.feishu.cn/document/ukTMukTMukTM/uETNz4SM1MjLxUzM//event/scope-change
	EventContactScopeChange = "contact_scope_change"

	// 员工变更 https://open.feishu.cn/document/ukTMukTMukTM/uETNz4SM1MjLxUzM//event/employee-change
	EventUserLeave  = "user_leave"
	EventUserAdd    = "user_add"
	EventUserUpdate = "user_update"

	// 用户状态变更 https://open.feishu.cn/document/ukTMukTMukTM/uETNz4SM1MjLxUzM//event/user-status-changed
	EventUserStatusChange = "user_status_change" // 比如用户从正常变为暂停、从暂停恢复正常等

	// 解散群 https://open.feishu.cn/document/ukTMukTMukTM/uQDOwUjL0gDM14CN4ATN/event/group-closed
	EventChatDisband = "chat_disband" //解散群

	// 部门变更  https://open.feishu.cn/document/ukTMukTMukTM/uETNz4SM1MjLxUzM//event/department-update
	EventDeptAdd    = "dept_add"    //新建部门
	EventDeptDelete = "dept_delete" //删除部门
	EventDeptUpdate = "dept_update" //修改部门

)

const (
	EventValidation = "url_verification" // 请求网址校验 https://open.feishu.cn/document/ukTMukTMukTM/uYDNxYjL2QTM24iN0EjN/event-subscription-configure-/request-url-configuration-case

	EventAppTicket = "app_ticket" // app_ticket 事件

	// 首次启用应用触发此事件 https://open.feishu.cn/document/uAjLw4CM/ukTMukTMukTM/application-v6/event/app-first-enabled
	EventAppOpen = "app_open"

	// 用户和机器人的会话首次被创建 https://open.feishu.cn/document/ukTMukTMukTM/uYDNxYjL2QTM24iN0EjN/bot-events
	EventP2pChatCreate = "p2p_chat_create"

	EventMessage = "message" // 飞书机器人群聊会话事件

	// 应用停启用 https://open.feishu.cn/document/uAjLw4CM/ukTMukTMukTM/application-v6/event/app-enabled-or-disabled
	EventAppStatusChange = "app_status_change"

	// 应用商店应用购买  https://open.feishu.cn/document/uAjLw4CM/ukTMukTMukTM/application-v6/event/public-app-purchase
	EventOrderPaid = "order_paid"
)

var log = logger.GetDefaultLogger()

// 这个结构体 会被event事件和 卡片回调使用
// event 关注encrypt
// 卡片回调 type 要验证 url_verification
type ValidationReq struct {
	Challenge string `json:"challenge"`
	Token     string `json:"token"`
	Type      string `json:"type"`
	Encrypt   string `json:"encrypt"`
}

type EventReq struct {
	Token string       `json:"token"`
	Event EventReqData `json:"event"`
}

type EventReqData struct {
	Type string `json:"type"`
}

type AbstractHandler interface {
	Handle(data string) (string, errs.SystemErrorInfo)
}

var FsCallBackHandlers = map[string]AbstractHandler{
	EventValidation:           FsValidationHandler{},
	EventAppTicket:            AppTicketHandler{},
	EventAppOpen:              AppOpenHandler{},
	EventP2pChatCreate:        P2pChatCreateHandler{},
	EventMessage:              MessageHandler{}, // 机器人、群聊、消息相关的
	EventContactScopeChange:   ContactScopeChangeHandler{},
	EventContactScopeChangeV3: ContactScopeChangeV3Handler{},
	EventUserLeave:            UserLeaveHandler{},
	EventUserAdd:              UserAddHandler{},
	// 用户修改信息事件
	EventUserUpdate:            UserUpdateHandler{},
	EventUserStatusChange:      UserStatusChangeHandler{},
	EventAppStatusChange:       AppStatusChangeHandler{},
	EventOrderPaid:             OrderPaidHandler{},
	EventChatDisband:           ChatDisbandHandler{},
	EventImChatDisabledV1:      ImChatDisbandV1Handler{},
	EventAddUserToChat:         ChatUserJoinHandler{},
	EventRemoveUserFromChat:    ChatUserDeleteHandler{},
	EventRevokeAddUserFromChat: ChatUserChangeHandler{},
	EventDeptAdd:               DeptChangeHandler{},
	EventDeptDelete:            DeptChangeHandler{},
	EventDeptUpdate:            DeptChangeHandler{},
	EventChatBotEntry:          ChatBotEntryHandler{},
}

func FsCallbackHandlerFunc(c *gin.Context) {
	httpContext := c.Request.Context().Value(consts.HttpContextKey).(model.HttpContext)
	threadlocal.Mgr.SetValues(gls.Values{consts.HttpContextKey: httpContext, consts.TraceIdKey: httpContext.TraceId}, func() {
		FsCallbackHandler(c.Writer, c.Request)
	})
}

// 解析事件
// 1.0版本事件和2.0版本事件结构体区别 https://open.feishu.cn/document/ukTMukTMukTM/uUTNz4SN1MjL1UzM
func ParseFeiShuEventData(data string) (string, string, errs.SystemErrorInfo) {
	fsConfig := config.GetConfig().FeiShu
	if fsConfig == nil {
		log.Error("[ParseFeiShuEventData] 飞书配置信息不存在")
		return "", "", errs.FeiShuConfigNotExistError
	}
	encryptKey := fsConfig.EventEncryptKey
	verificationToken := fsConfig.EventVerifyToken

	validationReq := &ValidationReq{}
	errSys := json.FromJson(data, validationReq)
	if errSys != nil {
		log.Errorf("[ParseFeiShuEventData] json err:%v, 原始数据:%s", errSys, data)
	}

	// 如果加密，先解密, 把飞书回调的encrypt字段解密
	if validationReq.Encrypt != "" {
		encryptStr := validationReq.Encrypt
		afterDecrypt, err := encrypt.AesDecrypt(encryptKey, encryptStr)
		if err != nil {
			log.Errorf("[ParseFeiShuEventData] err: %v", err)
			return "", "", errs.BuildSystemErrorInfo(errs.DecryptError, err)
		}
		data = afterDecrypt
		log.Infof("[ParseFeiShuEventData] 解密后的数据：%s", afterDecrypt)
		err = json.FromJson(data, validationReq)
		if err != nil {
			log.Errorf("[ParseFeiShuEventData]err:%v", err)
		}
	}

	eventV1Req := &req_vo.FeiShuEventV1{}
	eventV2Req := &req_vo.FeiShuBaseEventV2{}
	err := json.FromJson(data, eventV1Req)
	if err != nil {
		// 这里如果解析失败，暂时不return，现阶段有两个版本的事件，一次解析失败，再去解析
		log.Errorf("[ParseFeiShuEventData] json err:%v", err)
	}
	if eventV1Req.Token == "" {
		// 说明是2.0版本的事件
		err = json.FromJson(data, eventV2Req)
		if err != nil {
			log.Errorf("[ParseFeiShuEventData] json err:%v, 原始数据:%s", err, data)
		}
		// 转换为旧版格式
		eventV1Req.Token = eventV2Req.Header.Token
		eventV1Req.Event = req_vo.FeiShuEventV1Base{
			Type: eventV2Req.Header.EventType,
		}
	}

	// 验证token
	if eventV1Req.Token != verificationToken {
		log.Error("[ParseFeiShuEventData] 无效的token")
		return "", "", errs.FeiShuEventCheckTokenFailed
	}

	eventType := ""
	if validationReq.Type == EventValidation {
		eventType = EventValidation
	} else {
		eventType = eventV1Req.Event.Type
	}

	// 返回事件type和解密后的数据
	return eventType, data, nil

}

func FsCallbackHandler(w http.ResponseWriter, r *http.Request) {
	fsConfig := config.GetConfig().FeiShu
	if fsConfig == nil {
		log.Error("飞书配置信息不存在")
		return
	}

	reqBody, _ := ioutil.ReadAll(r.Body)
	reqObj := string(reqBody)
	log.Infof("[FsCallbackHandler] 飞书事件请求原始数据 %s", reqObj)

	// 解析事件
	// 1.0版本事件和2.0版本事件结构体区别 https://open.feishu.cn/document/ukTMukTMukTM/uUTNz4SN1MjL1UzM
	//eventReq := &EventReq{}
	//_ = json.FromJson(reqObj, eventReq)
	//if eventReq.Token == "" {
	//	eventV3Req := &req_vo.EventV3Req{}
	//	_ = json.FromJson(reqObj, eventV3Req)
	//	// 转换为旧版格式
	//	eventReq.Token = eventV3Req.Header.Token
	//	eventReq.CustomEvent = EventReqData{
	//		Type: eventV3Req.Header.EventType,
	//	}
	//}

	eventType, parseData, errSys := ParseFeiShuEventData(reqObj)
	if errSys != nil {
		log.Errorf("[FsCallbackHandler]ParseFeiShuEventData, err:%v", errSys)
		return
	}

	respBody := ""
	// 将处理封装到一个方法中，如果处理异常，则将事件记录到数据库中
	// 飞书回调的事件都在 handleWithData 中处理了
	if resp, err := handleWithData(eventType, parseData); err != nil {
		respBody = resp
		// 记录到数据库
		resp := msgfacade.WriteSomeFailedMsg(msgvo.WriteSomeFailedMsgReqVo{
			OrgId: 0,
			Input: &msgvo.WriteSomeFailedMsgReqVoData{
				MsgType: 0,
				MsgBody: reqObj,
			},
		})
		if resp.Failure() {
			log.Errorf("handleWithData error:%v, 请求数据%s, 响应数据%s", resp.Error(), reqObj, respBody)
		}
	} else {
		respBody = resp
	}
	log.Infof("请求数据%s, 响应数据%s", reqObj, respBody)
	_, ioErr := fmt.Fprintln(w, respBody)
	if ioErr != nil {
		log.Error(ioErr)
	}
	return
}

func handleWithData(eventType, reqObj string) (string, errs.SystemErrorInfo) {
	//fsConfig := config.GetConfig().FeiShu
	//if fsConfig == nil {
	//	log.Error("[handleWithData] 飞书配置信息不存在")
	//	return "", errs.FeiShuConfigNotExistError
	//}
	//encryptKey := fsConfig.EventEncryptKey
	//verificationToken := fsConfig.EventVerifyToken
	//
	//validationReq := &ValidationReq{}
	//_ = json.FromJson(reqObj, validationReq)
	//
	//// 如果加密，先解密, 把飞书回调的encrypt字段解密
	//if validationReq.Encrypt != "" {
	//	encryptStr := validationReq.Encrypt
	//	afterDecrypt, err := encrypt.AesDecrypt(encryptKey, encryptStr)
	//	if err != nil {
	//		log.Errorf("[handleWithData] err: %v", err)
	//		return "", errs.BuildSystemErrorInfo(errs.DecryptError, err)
	//	}
	//	reqObj = afterDecrypt
	//	validationReq = &ValidationReq{}
	//	_ = json.FromJson(reqObj, validationReq)
	//}
	//
	//// 针对通讯录授权范围发生变更为了兼容新旧版本：[旧版](https://open.feishu.cn/document/ukTMukTMukTM/uITNxYjLyUTM24iM1EjN)可使用 `EventReq` 结构体
	//// 而[新版](https://open.feishu.cn/document/uAjLw4CM/ukTMukTMukTM/reference/contact-v3/scope/events/updated)则使用新的结构体 `req_vo.EventV3Req`
	//// 解析事件 1.0版本和2.0版本事件 https://open.feishu.cn/document/ukTMukTMukTM/uUTNz4SN1MjL1UzM
	//eventReq := &EventReq{}
	//eventV3Req := &req_vo.EventV3Req{}
	//_ = json.FromJson(reqObj, eventReq)
	//if eventReq.Token == "" {
	//	_ = json.FromJson(reqObj, eventV3Req)
	//	log.Infof("[handleWithData] eventV3Req: %s", reqObj)
	//	// 转换为旧版格式
	//	eventReq.Token = eventV3Req.Header.Token
	//	eventReq.CustomEvent = EventReqData{
	//		Type: eventV3Req.Header.EventType,
	//	}
	//}
	//
	//// 验证token
	//if eventReq.Token != verificationToken {
	//	log.Error("[handleWithData] 无效的token")
	//	return "", errs.FeiShuEventCheckTokenFailed
	//}
	//
	//reqEvent := ""
	//if validationReq.Type == EventValidation {
	//	reqEvent = EventValidation
	//} else {
	//	reqEvent = eventReq.CustomEvent.Type
	//}

	log.Infof("[handleWithData] 处理事件 %s，event body: %s", eventType, reqObj)

	// 处理事件
	handler, ok := FsCallBackHandlers[eventType]
	if !ok {
		log.Infof("[handleWithData] %s事件不在正常处理范围内", eventType)
		return "", errs.FeiShuEventNotSupport
	}

	respBody := ""
	reqBo := bo.FeiShuCallBackMqBo{
		Type:    eventType,
		Data:    reqObj,
		MsgTime: time.Now(),
	}

	switch eventType {
	case EventValidation, EventAppTicket, EventChatDisband, EventAddUserToChat, EventImChatDisabledV1,
		EventRemoveUserFromChat, EventRevokeAddUserFromChat, EventDeptUpdate, EventDeptDelete, EventDeptAdd, EventChatBotEntry:
		respBody, _ = handler.Handle(reqObj)
	//case EventAddBot:
	//	msgPush(reqBo)
	case EventMessage:
		msgPush(reqBo)
	case EventOrderPaid:
		orderPush(reqBo)
	case EventUserAdd, EventUserLeave, EventUserUpdate, EventUserStatusChange:
		userAddReq := &UserAddReq{}
		_ = json.FromJson(reqObj, userAddReq)
		//判断是否在白名单里面，如果在的话就不处理
		baseCallBack := base.CallBackBase{}
		isWhiteList, err := baseCallBack.CheckIsInWhiteList(userAddReq.Event.TenantKey)
		if err != nil {
			return "", err
		}
		if isWhiteList {
			return "", nil
		}
		userChangePush(reqBo)
	case EventContactScopeChange, EventContactScopeChangeV3:
		//通讯录变更回调
		contactScopeChangePush(reqBo)
	default:
		commonPush(reqBo)
	}

	return respBody, nil
}

func orderPush(req bo.FeiShuCallBackMqBo) {
	asyn.Execute(func() {
		domain.PushFeiShuCallBackOrder(req)
	})
}

func msgPush(req bo.FeiShuCallBackMqBo) {
	asyn.Execute(func() {
		domain.PushFeiShuCallBackMsg(req)
	})
}

func userChangePush(req bo.FeiShuCallBackMqBo) {
	asyn.Execute(func() {
		domain.PushFeiShuCallBackUserChange(req)
	})
}

func contactScopeChangePush(req bo.FeiShuCallBackMqBo) {
	asyn.Execute(func() {
		domain.PushFeiShuCallBackContactScopeChange(req)
	})
}

func commonPush(req bo.FeiShuCallBackMqBo) {
	asyn.Execute(func() {
		domain.PushFeiShuCallBack(req)
	})
}

func HandleError(err error) {
	if err != nil {
		log.Error(strs.ObjectToString(err))
		panic(err)
	}
}
