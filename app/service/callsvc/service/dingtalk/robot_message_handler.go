package callsvc

import (
	"strings"

	sdk_const "gitea.bjx.cloud/allstar/platform-sdk/consts"
	"github.com/star-table/startable-server/common/model/vo/projectvo"
)

// 钉钉群聊助手指令
func HandleDingRobotMsg(msg DingRobotMsg) error {
	// 查询是否关联多个项目
	dealInfo, errSys := domain.GetDingGroupChatHandleInfos(msg.SenderCorpId, msg.ConversationId, sdk_const.SourceChannelDingTalk)
	if errSys != nil {
		log.Errorf("[HandleDingRobotMsg] GetDingGroupChatHandleInfos err:%v", errSys)
		return errSys
	}
	bindProjectIds := dealInfo.BindProjectIds

	if len(bindProjectIds) <= 0 || len(bindProjectIds) > 1 {
		// 发送卡片，提示不支持该群聊指令
		return nil
	}

	// 需要解析用户指令
	content := strings.Trim(msg.TEXT.Content, " ")
	switch content {
	case consts.InstructUserInsProIssue:
		// 项目任务
		resp := projectfacade.HandleGroupChatUserInsProIssue(projectvo.HandleGroupChatUserInsProIssueReqVo{
			Input: &projectvo.HandleGroupChatUserInsProIssueReq{
				OpenChatId:    msg.ConversationId,
				OpenId:        msg.SenderStaffId,
				SourceChannel: sdk_const.SourceChannelDingTalk,
			},
		})
		if resp.Failure() {
			log.Errorf("[HandleDingRobotMsg] failed:%v, chatId:%s", resp.Error(), msg.ConversationId)
			return resp.Error()
		}
		if resp.Data.IsTrue {
			log.Infof("[HandleDingRobotMsg] 钉钉 项目任务消息发送成功！")
		}
		return nil
	case consts.InstructUserInsProProgress:
		// 项目进展
		resp := projectfacade.HandleGroupChatUserInsProProgress(projectvo.HandleGroupChatUserInsProProgressReqVo{
			Input: &projectvo.HandleGroupChatUserInsProProgressReq{
				OpenChatId: msg.MsgId,
				OpenId:     msg.ConversationId,
			},
		})
		if resp.Failure() {
			log.Error(resp.Error())
			return resp.Error()
		}
		if resp.Data.IsTrue {
			log.Infof("[HandleDingRobotMsg] 钉钉 项目进展消息发送成功！")
		}
		return nil
	case consts.InstructUserInsProDynSetting:
		// 项目动态设置
		resp := projectfacade.HandleGroupChatUserInsProjectSettings(projectvo.HandleGroupChatUserInsProjectSettingsReq{
			OpenChatId:    msg.ConversationId,
			SourceChannel: sdk_const.SourceChannelDingTalk,
		})
		if resp.Failure() {
			log.Error(resp.Error())
			return resp.Error()
		}
		if resp.Data.IsTrue {
			log.Infof("[HandleDingRobotMsg] 钉钉 项目动态设置消息发送成功！")
		}
		return nil
	default:
		// 解析@负责人姓名的指令，如果指令匹配不上，就提示听不懂的信息
		outUserIds := []string{}
		for _, v := range msg.AtUsers {
			if v.StaffId == "" {
				continue
			}
			outUserIds = append(outUserIds, v.StaffId)
		}
		// 没有outUserIds，说明没有@群成员
		if len(outUserIds) <= 0 || len(outUserIds) > 1 {
			// 错误的指令
			return nil
		}
		if content != "" {
			// @用户 创建记录

		} else {
			// 查看该用户的任务完成情况
		}
	}
	return nil
}
