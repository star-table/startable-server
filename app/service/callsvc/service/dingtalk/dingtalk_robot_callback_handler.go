package callsvc

import (
	"io/ioutil"

	"github.com/star-table/startable-server/common/core/util/json"
	"github.com/gin-gonic/gin"
)

type DingRobotMsg struct {
	MsgType           string      `json:"msgtype"` // 目前只支持text
	TEXT              TextContent `json:"text"`
	MsgId             string      `json:"msgId"`             // 加密的消息id
	CreateAt          string      `json:"createAt"`          // 消息的时间戳，单位ms
	ConversationType  string      `json:"conversationType"`  // 1单聊  2群聊
	ConversationId    string      `json:"conversationId"`    // 加密的会话id
	ConversationTitle string      `json:"conversationTitle"` // 群聊时才会有的会话标题
	SenderId          string      `json:"senderId"`          // 加密的发送者id
	SenderNick        string      `json:"senderNick"`        // 发送者昵称
	SenderCorpId      string      `json:"senderCorpId"`      // 企业内部群有的发送者当前群的企业corpId
	SenderStaffId     string      `json:"senderStaffId"`     // 企业内部群有的发送者在企业内的userid
	ChatbotUserId     string      `json:"chatbotUserId"`     // 加密的机器人ID
	AtUsers           []AtUsers   `json:"atUsers"`           // 被@人的信息dingtalkId：加密的发送者ID staffId：企业内部群有的发送者在企业内的userid
}

type AtUsers struct {
	DingtalkId string `json:"dingtalkId"`
	StaffId    string `json:"staffId"`
}

type TextContent struct {
	Content string `json:"content"`
}

// 当用户@群助手时，接收用户指令信息
func DingTalkRobotCallBackHandlerFunc(c *gin.Context) {
	// 需要计算接收的header
	//timestampFromDing := c.Request.Header.Get("timestamp")
	//signFromDing := c.Request.Header.Get("sign")
	//
	//timeSFromD, err := strconv.ParseInt(timestampFromDing, 10, 64)
	//if err != nil {
	//	log.Errorf("[DingTalkRobotCallBackHandlerFunc] 类型转化错误 err:%v", err)
	//	return
	//}
	//
	//timestampNow := time.Now().Unix() * 1000
	//// timestamp与系统当前时间戳如果相差1小时以上，则认为是非法的请求。
	//if timestampNow > (timeSFromD+3600*1000) || timestampNow < (timeSFromD-3600*1000) {
	//	log.Errorf("[DingTalkRobotCallBackHandlerFunc] timestamp与当前时间不一致")
	//	return
	//}
	//
	//// header中的timestamp + "\n" + 机器人的appSecret当做签名字符串，
	//// 使用HmacSHA256算法计算签名，然后进行Base64 encode，得到最终的签名值。
	//secret := timestampFromDing + "\n" + "appsecret"
	//message := ""
	//hashData := hmac.New(sha256.New, []byte(secret))
	//hashData.Write([]byte(message))
	//signNow := base64.StdEncoding.EncodeToString([]byte(message))
	//
	//if signNow != signFromDing {
	//	log.Errorf("sign不合法, oldSign:%s, newSign:%s", signFromDing, signNow)
	//	return
	//}

	robotMsg := DingRobotMsg{}
	defer c.Request.Body.Close()
	body, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		log.Error(err)
		return
	}
	json.FromJson(string(body), &robotMsg)

	log.Infof("[DingTalkRobotCallBackHandlerFunc]ding 群助手消息: %s", body)

	// 处理群助手指令
	HandleDingRobotMsg(robotMsg)
}
