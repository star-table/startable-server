package callsvc

import (
	"crypto/sha1"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/star-table/startable-server/common/core/config"
	"github.com/star-table/startable-server/common/core/consts"
	consts2 "github.com/star-table/startable-server/common/core/consts"
	"github.com/star-table/startable-server/common/core/errs"
	"github.com/star-table/startable-server/common/core/model"
	"github.com/star-table/startable-server/common/core/threadlocal"
	"github.com/star-table/startable-server/common/core/util/json"
	"github.com/star-table/startable-server/common/core/util/stack"
	"github.com/star-table/startable-server/common/library/cache"
	"github.com/gin-gonic/gin"
	"github.com/jtolds/gls"
)

type CardReq struct {
	OpenId        string        `json:"open_id"`
	OpenMessageId string        `json:"open_message_id"`
	TenantKey     string        `json:"tenant_key"`
	Token         string        `json:"token"`
	Action        CardReqAction `json:"action"`
}

type CardReqAction struct {
	Tag      string              `json:"tag"`
	Option   string              `json:"option"`
	Timezone string              `json:"timezone"`
	Value    *CardReqActionValue `json:"value"`
}

type CardReqActionValue struct {
	IssueId                           int64  `json:"issueId"`
	AppId                             string `json:"appId"`
	ProjectId                         int64  `json:"projectId"`
	TableId                           string `json:"tableId"`
	OrgId                             int64  `json:"orgId"`
	UserId                            int64  `json:"userId"`
	StatusType                        int    `json:"statusType"`
	CardType                          string `json:"cardType"`
	Action                            string `json:"action"`
	FsCardValueUrgeIssueBeUrgedUserId int64  `json:"fsCardValueUrgeIssueBeUrgedUserId"`
	FsCardValueUrgeIssueUrgeText      string `json:"fsCardValueUrgeIssueUrgeText"`
	FsCardValueAuditIssueResult       int    `json:"fsCardValueAuditIssueResult"`    // 任务审批审批结果
	FsCardValueAuditIssueStarterId    int64  `json:"fsCardValueAuditIssueStarterId"` // 审批审批的发起人
	FsCardValueIsGroupChat            bool   `json:"fsCardValueIsGroupChat"`
	FsCardValueGroupChatId            string `json:"fsCardValueGroupChatId"`
	CardTitle                         string `json:"cardTitle"`
}

type CardAbstractHandler interface {
	Handle(cardReq CardReq) (string, errs.SystemErrorInfo)
}

// 卡片回调事件处理
var cardHandlers = map[string]CardAbstractHandler{
	consts2.FsCardTypeIssueRemind:     IssueRemindHandler{}, // 事件：任务即将到达截止事件提醒处理
	consts2.FsCardTypeIssueUrge:       FsIssueUrgeHandler{},
	consts2.FsCardTypeIssueAudit:      FsIssueAuditHandler{}, // 任务审核的卡片
	consts2.FsCardTypeNoticeSubscribe: NoticeSubscribe{},
	consts2.FsCardTypeIssueCreate:     IssueCreateHandler{},   // 通过机器人指令创建记录 并和卡片交互的事件
	consts2.FsCardTypeDealTime:        IssueDealTimeHandler{}, // 卡片交互选择时间
}

func FsCardCallbackHandlerFunc(c *gin.Context) {
	httpContext := c.Request.Context().Value(consts.HttpContextKey).(model.HttpContext)
	threadlocal.Mgr.SetValues(gls.Values{consts.HttpContextKey: httpContext, consts.TraceIdKey: httpContext.TraceId}, func() {
		FsCardCallbackHandler(c.Writer, c.Request)
	})
}

func FsCardCallbackHandler(w http.ResponseWriter, r *http.Request) {
	defer func() {
		if r := recover(); r != nil {
			log.Error(errs.BuildSystemErrorInfoWithPanicRecover(r, stack.GetStack()))
		}
	}()

	fsConfig := config.GetConfig().FeiShu
	if fsConfig == nil {
		log.Error("飞书配置信息不存在")
		return
	}
	refreshToken := r.Header.Get("X-Refresh-Token")
	verificationToken := fsConfig.EventVerifyToken
	reqBody, _ := ioutil.ReadAll(r.Body)
	reqObj := string(reqBody)

	validationReq := &ValidationReq{}
	_ = json.FromJson(reqObj, validationReq)

	//url校验
	if validationReq.Type == EventValidation {
		respBody, _ := FsCallBackHandlers[EventValidation].Handle(reqObj)
		_, _ = w.Write([]byte(respBody))
		return
	}

	log.Infof("[FsCardCallbackHandler] 飞书卡片回调数据：%s", reqObj)
	cardReq, checkErr := checkCardReqValidity(r, reqObj, verificationToken)
	if checkErr != nil {
		log.Error(checkErr)
		return
	}

	log.Infof("[FsCardCallbackHandler] 飞书卡片回调 %s", json.ToJsonIgnoreError(cardReq))
	action := cardReq.Action
	actionValue := action.Value

	respBody := ""
	var err errs.SystemErrorInfo = nil
	if actionValue != nil {
		if handler, ok := cardHandlers[actionValue.CardType]; ok {
			respBody, err = handler.Handle(*cardReq)
		}
		log.Infof("[FsCardCallbackHandler] 飞书卡片回调响应：%s", respBody)
	}

	if err != nil {
		w.WriteHeader(500)
	}
	respCode, ioErr := w.Write([]byte(respBody))
	log.Infof("[FsCardCallbackHandler] 卡片回调消息响应结果 %d %v", respCode, ioErr)
	if err == nil && ioErr == nil {
		//三十分钟就够了
		cacheErr := cache.SetEx(consts2.CacheFeiShuCardCallBackRefreshToken+refreshToken, consts.BlankString, 60*30)
		if cacheErr != nil {
			log.Errorf("refresh_token 添加失败 %s %v", refreshToken, cacheErr)
		}
	}
}

// 校验请求
func checkCardReqValidity(r *http.Request, requestBody string, verificationToken string) (*CardReq, errs.SystemErrorInfo) {
	cardReq := &CardReq{}
	_ = json.FromJson(requestBody, cardReq)

	timestamp := r.Header.Get("X-Lark-Request-Timestamp")
	nonce := r.Header.Get("X-Lark-Request-Nonce")
	signature := r.Header.Get("X-Lark-Signature")
	refreshToken := r.Header.Get("X-Refresh-Token")

	repeatErr := checkCallMsgRepeatHandle(refreshToken)
	if repeatErr != nil {
		log.Error(repeatErr)
		return nil, repeatErr
	}

	log.Infof("飞书卡片回调请求头： timestamp %s, nonce %s, sign %s, refresh_token %s", timestamp, nonce, signature, refreshToken)
	if cardReq.OpenId == "" ||
		cardReq.TenantKey == "" ||
		cardReq.Token == "" ||
		timestamp == "" ||
		nonce == "" ||
		signature == "" {
		return nil, errs.BuildSystemErrorInfo(errs.ReqParamsValidateError)
	}
	var b strings.Builder
	b.WriteString(timestamp)
	b.WriteString(nonce)
	b.WriteString(verificationToken)
	b.WriteString(requestBody)

	bs := []byte(b.String())
	h := sha1.New()
	h.Write(bs)
	bs = h.Sum(nil)
	sig := fmt.Sprintf("%x", bs)

	if sig != signature {
		log.Errorf("签名校验失败 %s != %s", sig, signature)
		return nil, errs.BuildSystemErrorInfo(errs.FeiShuCardCallSignVerifyError)
	}
	return cardReq, nil
}

func checkCallMsgRepeatHandle(refreshToken string) errs.SystemErrorInfo {
	exist, err := cache.Exist(consts2.CacheFeiShuCardCallBackRefreshToken + refreshToken)
	if err != nil {
		log.Error(err)
		return errs.BuildSystemErrorInfo(errs.RedisOperateError)
	}
	if exist {
		log.Errorf("消息重复推送 %s", refreshToken)
		return errs.BuildSystemErrorInfo(errs.FeiShuCardCallMsgRepetError)
	}
	return nil
}
