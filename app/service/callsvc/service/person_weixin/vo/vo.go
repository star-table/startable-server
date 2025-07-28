package callsvc

type Verify struct {
	MsgSign   string `form:"signature" binding:"required"`
	Timestamp string `form:"timestamp" binding:"required"`
	Nonce     string `form:"nonce" binding:"required"`
	EchoStr   string `form:"echostr"`
}

type Message struct {
	Event        string `xml:"Event"`
	ToUserName   string `xml:"ToUserName"`
	FromUserName string `xml:"FromUserName"`
	MsgType      string `xml:"MsgType"`
	EventKey     string `xml:"EventKey"`
	Ticket       string `xml:"Ticket"`
}
