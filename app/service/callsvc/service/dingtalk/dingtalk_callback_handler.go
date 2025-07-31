package callsvc

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"gitea.bjx.cloud/allstar/dingtalk-sdk-golang"
	platform_sdk "gitea.bjx.cloud/allstar/platform-sdk"
	sdk_const "gitea.bjx.cloud/allstar/platform-sdk/consts"
	"github.com/gin-gonic/gin"
	"github.com/jtolds/gls"
	"github.com/star-table/startable-server/app/facade/msgfacade"
	"github.com/star-table/startable-server/common/core/config"
	"github.com/star-table/startable-server/common/core/consts"
	consts2 "github.com/star-table/startable-server/common/core/consts"
	"github.com/star-table/startable-server/common/core/errs"
	"github.com/star-table/startable-server/common/core/model"
	"github.com/star-table/startable-server/common/core/threadlocal"
	"github.com/star-table/startable-server/common/core/util/json"
	"github.com/star-table/startable-server/common/core/util/stack"
	"github.com/star-table/startable-server/common/core/util/strs"
	"github.com/star-table/startable-server/common/core/util/uuid"
)

type DingTalkReq struct {
	Encrypt string
}

type BaseCallbackMsg struct {
	EventType string                    `json:"EventType"`
	BizData   []req_vo.DingEventBizData `json:"bizData"`
}

type SyncActionData struct {
	SyncAction string `json:"syncAction"`
}

type AbstractHandler interface {
	Handle(data req_vo.DingEventBizData) error
}

var DingTalkHandlers = map[string]AbstractHandler{
	EventSuitTicket:     SuiteTicketHandler{},
	EventCheckServerUrl: CheckUrlCallbackHandler{},

	EventAuthApp:        TempAuthCodeHandler{},
	EventSuitRelieve:    SuiteRelieveHandler{},
	EventChangeAuth:     ChangeAuthHandler{},
	EventAppScopeUpdate: ChangeAuthHandler{},

	EventMarketOrder:       MarketOrderHandler{},
	EventMarketOrderRefund: MarketOrderRefundHandler{},
	EventOrderTrail:        MarketOrderTrailHandler{},

	EventDeptCreate: DeptChangeHandler{EventType: consts2.EventDeptAdd},
	EventDeptModify: DeptChangeHandler{EventType: consts2.EventDeptUpdate},
	EventDeptRemove: DeptChangeHandler{EventType: consts2.EventDeptDel},

	EventUserAdd:    UserAddHandler{},
	EventUserLeave:  UserLeaveHandler{},
	EventUserModify: UserUpdateHandler{},
	EventUserActive: UserActiveHandler{},

	EventChatDisband:        ChatDisbandHandler{},
	EventChatTemplateChange: ChatTemplateChangeHandler{},

	EventCoolAppInstall:   CoolAppInstallHandler{},
	EventCoolAppUninstall: CoolAppUninstallHandler{},
}

func DingTalkCallBackHandlerFunc(c *gin.Context) {
	log.Infof("ding:%v", *config.GetDingTalkSdkConfig())
	httpContext := c.Request.Context().Value(consts.HttpContextKey).(model.HttpContext)
	threadlocal.Mgr.SetValues(gls.Values{consts.HttpContextKey: httpContext, consts.TraceIdKey: httpContext.TraceId}, func() {
		DingTalkCallbackHandler(c.Writer, c.Request)
	})
}

func DingTalkCallbackHandler(w http.ResponseWriter, r *http.Request) {
	defer func() {
		if r := recover(); r != nil {
			log.Error(errs.BuildSystemErrorInfoWithPanicRecover(r, stack.GetStack()))
		}
	}()

	s, _ := ioutil.ReadAll(r.Body)

	dingTalkEncrypt := &DingTalkReq{}
	err := json.FromJson(string(s), dingTalkEncrypt)
	if err != nil {
		log.Error("call back happend err: " + strs.ObjectToString(err))
		return
	}

	reqEncrypt := dingTalkEncrypt.Encrypt
	reqSignature := r.FormValue("signature")
	reqTimestamp := r.FormValue("timestamp")
	reqNonce := r.FormValue("nonce")

	dingConfig := config.GetConfig().DingTalk
	log.Infof("dingConfig:%s", json.ToJsonIgnoreError(dingConfig))
	client, err := platform_sdk.GetClient(sdk_const.SourceChannelDingTalk, "")
	if err != nil {
		log.Errorf("[DingTalkCallbackHandler] err:%v", err)
		return
	}
	dingTalk := client.GetOriginClient().(*dingtalk.DingTalk)
	crypto, err := dingTalk.GetDingTalkCrypto(dingConfig.Token, dingConfig.AesKey)
	if err != nil {
		log.Errorf("dingTalk.GetDingTalkCrypto err:%v", err)
		return
	}
	//crypto.SuiteKey = dingConfig.SuiteKey
	jsonMsg, err := crypto.Decrypt(reqEncrypt, reqSignature, reqTimestamp, reqNonce)

	//jsonMsg, err := dingtalk.GetCrypto().DecryptMsg(reqSignature, reqTimestamp, reqNonce, reqEncrypt)
	if err != nil {
		log.Errorf("call back happend reqEncrypt:%v, reqSignature:%v, reqTimestamp:%v, reqNonce:%v err: %v", reqEncrypt, reqSignature, reqTimestamp, reqNonce, strs.ObjectToString(err))
		return
	}

	log.Infof("dingTalk回调事件 data:%v", jsonMsg)

	baseMsg := &BaseCallbackMsg{}
	_ = json.FromJson(jsonMsg, baseMsg)

	for _, bizData := range baseMsg.BizData {
		syncAction := &SyncActionData{}
		json.FromJson(bizData.BizData, syncAction)

		eventType := syncAction.SyncAction
		if bizData.BizType == EventDingAppTrialBizType {
			eventType = EventOrderTrail
		}

		// 分为立即执行和队列执行
		switch eventType {
		case EventCheckServerUrl, EventSuitTicket:
			handler := DingTalkHandlers[eventType]
			if handler == nil {
				fmt.Println("un know event ", eventType)
			} else {
				err = handler.Handle(bizData)
				if err != nil {
					log.Errorf("[DingTalkCallbackHandler] data:%v, err:%v", jsonMsg, err)
				}
			}
		default:
			pushDingTalkCallBack(req_vo.DingTalkCallBackMqVo{
				Type:    eventType,
				Data:    bizData,
				MsgTime: time.Now(),
			})
		}
	}
	//response success
	//timestamp := time.Now().UnixNano() / 1e6
	//timestampStr := strconv.FormatInt(timestamp, 10)
	//nonce := random.RandomString(8)
	//
	//encrypt, sign, err := dingtalk.GetCrypto().EncryptMsg("success", timestampStr, nonce)
	//if err != nil {
	//	HandleError(err)
	//}
	//m := map[string]interface{}{
	//	"msg_signature": sign,
	//	"timeStamp":     timestampStr,
	//	"nonce":         nonce,
	//	"encrypt":       encrypt,
	//}
	//
	//data, _ := json.ToJson(m)
	//fmt.Fprintln(w, data)

	msgResp, err := crypto.Encrypt("success")
	if err != nil {
		log.Errorf("crypto.Encrypt err:%v", err)
		return
	}
	fmt.Fprintln(w, msgResp.String())
}

func HandleError(err error) {
	if err != nil {
		log.Error(strs.ObjectToString(err))
		panic(err)
	}
}

func pushDingTalkCallBack(dingCallBack req_vo.DingTalkCallBackMqVo) {
	message, err := json.ToJson(dingCallBack)
	if err != nil {
		log.Error(err)
	}
	mqMessage := &model.MqMessage{
		Topic:          config.GetDingTalkCallBackTopic().Topic,
		Keys:           uuid.NewUuid(),
		Body:           message,
		DelayTimeLevel: 3,
	}
	msgErr := msgfacade.PushMsgToMqRelaxed(*mqMessage, 0, 0)
	if msgErr != nil {
		log.Errorf("[pushDingTalkCallBack] mq消息推送失败，入表失败，消息内容：%s, 失败信息: %v", json.ToJsonIgnoreError(mqMessage), msgErr)
	}
}
