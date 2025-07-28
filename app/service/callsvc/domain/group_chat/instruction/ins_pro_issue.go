package callsvc

/// 飞书群聊，用户指令处理

import (
	"github.com/star-table/startable-server/common/core/errs"
	"github.com/star-table/startable-server/common/core/logger"
	"github.com/star-table/startable-server/common/model/vo/projectvo"
)

var log = logger.GetDefaultLogger()

// / 定义一个用户指令 handle：项目任务
type UserInsProIssue struct {
	InsName string
}

func NewUserInsProIssue() *UserInsProIssue {
	return &UserInsProIssue{
		InsName: "UserInsProIssue",
	}
}

// 获取指令名称
func (h *UserInsProIssue) GetInsName() string {
	return h.InsName
}

// 指令的业务处理
func (h *UserInsProIssue) Handler(reqEvent *req_vo.UserAtBotReqEvent) errs.SystemErrorInfo {
	resp := projectfacade.HandleGroupChatUserInsProIssue(projectvo.HandleGroupChatUserInsProIssueReqVo{
		Input: &projectvo.HandleGroupChatUserInsProIssueReq{
			OpenChatId: reqEvent.OpenChatID,
			OpenId:     reqEvent.OpenID,
		},
	})
	if resp.Failure() {
		log.Error(resp.Error())
	}
	log.Infof("飞书群聊-处理用户指令-UserInsProIssue-code: %v", resp.Code)
	return nil
}
