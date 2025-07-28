package callsvc

/// 接受用户消息的结构体
type UserAtBotReq struct {
	Ts    string `json:"ts"`
	Uuid  string `json:"uuid"`
	Token string `json:"token"`
	Type  string `json:"type"`

	Event UserAtBotReqEvent `json:"event"`
}

// event 结构体
type UserAtBotReqEvent struct {
	AppId     string `json:"app_id"`
	TenantKey string `json:"tenant_key"`
	Type      string `json:"type"`

	RootID           string `json:"root_id"`
	ParentID         string `json:"parent_id"`
	OpenChatID       string `json:"open_chat_id"`
	ChatType         string `json:"chat_type"`
	MsgType          string `json:"msg_type"`
	OpenID           string `json:"open_id"`
	EmployeeID       string `json:"employee_id"`
	UnionID          string `json:"union_id"`
	OpenMessageID    string `json:"open_message_id"`
	IsMention        bool   `json:"is_mention"`
	Text             string `json:"text"`
	TextWithoutAtBot string `json:"text_without_at_bot"`
}
