package callsvc

import "time"

type EventV3Req struct {
	Schema string           `json:"schema"`
	Header EventV3ReqHeader `json:"header"`
	Event  EventV3ReqEvent  `json:"event"` // 这里应该使用 interface{} 装载，否则无法透传 event 中的结构
}

type EventV3ReqHeader struct {
	EventId    string `json:"event_id"`
	Token      string `json:"token"`
	CreateTime string `json:"create_time"`
	EventType  string `json:"event_type"`
	TenantKey  string `json:"tenant_key"`
	AppId      string `json:"app_id"`
}

type EventV3ReqEvent struct {
	Added   EventV3ReqEventAddOrRemoved `json:"added"`
	Removed EventV3ReqEventAddOrRemoved `json:"removed"`

	FeiShuEventBaseInfo
	Users []ChatUserItem `json:"users"`
}

type EventV3ReqEventAddOrRemoved struct {
	Departments []*EventV3ReqEventDeptItem `json:"departments"`
	Users       []*EventV3ReqEventUserItem `json:"users"`
	UserGroups  interface{}                `json:"user_groups"` // 暂未接入
}

type EventV3ReqEventDeptItem struct {
	Name               string                        `json:"name"`
	I18NName           I18NName                      `json:"i18n_name"`
	ParentDepartmentID string                        `json:"parent_department_id"`
	DepartmentID       string                        `json:"department_id"`
	OpenDepartmentID   string                        `json:"open_department_id"`
	LeaderUserID       string                        `json:"leader_user_id"`
	ChatID             string                        `json:"chat_id"`
	Order              string                        `json:"order"`
	UnitIds            []string                      `json:"unit_ids"`
	MemberCount        int                           `json:"member_count"`
	Status             EventV3ReqEventDeptItemStatus `json:"status"`
}

type EventV3ReqEventUserItem struct {
	UnionID  string                        `json:"union_id"`
	UserID   string                        `json:"user_id"`
	OpenID   string                        `json:"open_id"`
	Name     string                        `json:"name"`
	EnName   string                        `json:"en_name"`
	Email    string                        `json:"email"`
	Mobile   string                        `json:"mobile"`
	Gender   int                           `json:"gender"`
	Avatar   EventV3ReqEventUserItemAvatar `json:"avatar"`
	Status   EventV3ReqEventUserItemStatus `json:"status"`
	IsFrozen bool                          `json:"is_frozen"`
}

type EventV3ReqEventUserItemAvatar struct {
	Avatar72     string `json:"avatar_72"`
	Avatar240    string `json:"avatar_240"`
	Avatar640    string `json:"avatar_640"`
	AvatarOrigin string `json:"avatar_origin"`
}

type EventV3ReqEventUserItemStatus struct {
	IsFrozen    bool `json:"is_frozen"`
	IsResigned  bool `json:"is_resigned"`
	IsActivated bool `json:"is_activated"`
}

type EventV3ReqEventDeptItemStatus struct {
	IsDeleted bool `json:"is_deleted"`
}

type I18NName struct {
	ZhCn string `json:"zh_cn"`
	JaJp string `json:"ja_jp"`
	EnUs string `json:"en_us"`
}

// 飞书1.0版本事件回调结构体
type FeiShuEventV1 struct {
	UUID  string            `json:"uuid"`
	Event FeiShuEventV1Base `json:"event"`
	Token string            `json:"token"`
	Ts    string            `json:"ts"`
	Type  string            `json:"type"`
}

type FeiShuEventV1Base struct {
	AppId          string   `json:"app_id"`
	OpenChatId     string   `json:"open_chat_id"`
	OpenId         string   `json:"open_id"`
	OpenMessageIds []string `json:"open_message_ids"`
	TenantKey      string   `json:"tenant_key"`
	Type           string   `json:"type"`
}

// 飞书2.0版本事件回调基本结构体
type FeiShuBaseEventV2 struct {
	Schema string              `json:"schema"` // 事件模式
	Header FeiShuEventHeaderV2 `json:"header"` // 事件头
}

type FeiShuEventHeaderV2 struct {
	EventId    string `json:"event_id"`    // 事件id
	EventType  string `json:"event_type"`  // 事件类型
	CreateTime string `json:"create_time"` // 事件创建时间戳（单位：毫秒）
	Token      string `json:"token"`       // 事件Token
	AppId      string `json:"app_id"`      // 应用ID
	TenantKey  string `json:"tenant_key"`  // 租户key
}

type FeiShuUserIds struct {
	UnionId string `json:"union_id"`
	UserId  string `json:"user_id"`
	OpenId  string `json:"open_id"`
}

type FeiShuEventBaseInfo struct {
	ChatId            string        `json:"chat_id"`
	OperatorId        FeiShuUserIds `json:"operator_id"`
	External          bool          `json:"external"`
	OperatorTenantKey string        `json:"operator_tenant_key"`
}

type ChatUserItem struct {
	Name      string        `json:"name"`
	TenantKey string        `json:"tenant_key"`
	UserId    FeiShuUserIds `json:"user_id"`
}

type DingEventBizData struct {
	GmtCreate   int64  `json:"gmt_create"`
	BizType     int    `json:"biz_type"`
	OpenCursor  int    `json:"open_cursor"`
	SubscribeId string `json:"subscribe_id"`
	Id          int    `json:"id"`
	GmtModified int64  `json:"gmt_modified"`
	BizId       string `json:"biz_id"`
	BizData     string `json:"biz_data"`
	CorpId      string `json:"corp_id"`
	Status      int    `json:"status"`
}

type DingTalkCallBackMqVo struct {
	//类型
	Type string `json:"type"`
	//数据
	Data DingEventBizData `json:"data"`
	//消息时间
	MsgTime time.Time `json:"msgTime"`
}
