package callsvc

import "github.com/star-table/startable-server/common/core/consts"

var (
	// 飞书用户信息更改的事件推送处理
	CacheFsUserUpdateWhiteList = consts.CacheKeyPrefix + consts.OrgsvcApplicationName + consts.CacheKeyOfSys + "fs_user_update_white_list"
	//飞书 终止授权范围回调消费
	CacheFeiShuStopSyncScope = consts.CacheKeyPrefix + consts.CacheKeyOfSys + "stop_sync_scope"
)

// 消息指令
const (
	//创建任务指令
	InstructCreateIssue = "create"
	//设置指令
	InstructSettings = "settings"
	// 通知设置修改指令 notification
	InstructNotification = consts.BotPrivateInsNameOfSetting
	// 项目任务
	InstructUserInsProIssue = "项目任务"
	// 项目进展
	InstructUserInsProProgress = "项目进展"
	//  项目动态设置
	InstructUserInsProDynSetting = "设置"
	// 用户指令：@负责人姓名
	InstructAtUserName = "InstructAtUserName"
	// 用户指令：@负责人姓名 任务标题
	InstructAtUserNameWithIssueTitle = "InstructAtUserNameWithIssueTitle"
)

const (
	EventV1Type = "event_callback"
)

const (
	EventSuitTicketBizType   = 2
	EventDingAppTrialBizType = 63 // 应用试用开通记录的回调信息
	EventDingOrderTrail      = "ding_order_trail"
)
