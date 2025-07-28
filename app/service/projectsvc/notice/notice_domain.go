package notice

import (
	"github.com/star-table/startable-server/common/core/logger"
	"github.com/star-table/startable-server/common/model/bo"
)

var log = logger.GetDefaultLogger()

const (
	//推送范围
	PushNoticeRangeTypeOwner       = 1
	PushNoticeRangeTypeParticipant = 2
	PushNoticeRangeTypeAttention   = 3

	//推送目标
	PushNoticeTargetTypeRemindMsg    = 1
	PushNoticeTargetTypeCommentAtMsg = 2
	PushNoticeTargetTypeCreateMsg    = 3
	PushNoticeTargetTypeModifyMsg    = 4
	PushNoticeTargetTypeRelationMsg  = 5

	//无限制类型
	PushNoticeNotLimit = 0
)

// userPushRangeConfigContinueFlag 推送范围
// rangeType 成员类型 0 负责人, 1 参与人，2 关注人  3协作人
// 返回 true 表示不推送
// 不满足条件，返回 true
func userPushRangeConfigContinueFlag(rangeType int, userConfig *bo.UserConfigBo) bool {
	if rangeType == 0 && userConfig.OwnerRangeStatus != 1 {
		return true
	} else if rangeType == 1 && userConfig.ParticipantRangeStatus != 1 {
		return true
	} else if rangeType == 2 && userConfig.AttentionRangeStatus != 1 {
		return true
	} else if rangeType == 3 && userConfig.CollaborateMessageStatus != 1 {
		return true
	}
	return false
}

// 推送目标, 1 任务提醒，2任务评论/详情被 at，3，创建任务，4修改任务，5关联信息
func userPushTargetConfigContinueFlag(target int, userConfig *bo.UserConfigBo) bool {
	if target == 1 && userConfig.RemindMessageStatus != 1 {
		return true
	} else if target == 2 && userConfig.CommentAtMessageStatus != 1 {
		return true
	} else if (target == 3 || target == 4) && userConfig.CreateRangeStatus != 1 && userConfig.ModifyMessageStatus != 1 {
		return true
	} else if target == 5 && userConfig.RelationMessageStatus != 1 {
		return true
	}
	return false
}
