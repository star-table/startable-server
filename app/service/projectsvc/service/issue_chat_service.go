package service

import (
	"strings"

	sdk_const "gitea.bjx.cloud/allstar/platform-sdk/consts"
	"github.com/star-table/startable-server/app/service/projectsvc/domain"
	"github.com/star-table/startable-server/common/core/errs"
	"github.com/star-table/startable-server/common/model/vo/projectvo"
)

// IssueStartChat 发起任务群聊讨论
func IssueStartChat(orgId, userId int64, sourceChannel string, input projectvo.IssueStartChatReq) (projectvo.IssueStartChatRespData, errs.SystemErrorInfo) {
	result := projectvo.IssueStartChatRespData{}
	// 目前只支持飞书
	if sourceChannel != sdk_const.SourceChannelFeishu {
		return result, nil
	}

	issueChatObj, err := domain.NewFsIssueChat(orgId, userId, input.IssueId, strings.Trim(input.Topic, " \n\t\r"))
	if err != nil {
		log.Errorf("[IssueStartChat] NewFsIssueChat err: %v, issueId: %d", err, input.IssueId)
		return result, err
	}
	outChatId, err := issueChatObj.StartIssueChat()
	if err != nil {
		log.Errorf("[IssueStartChat] StartIssueChat err: %v, issueId: %d", err, input.IssueId)
		return result, err
	}
	result.ChatId = outChatId
	return result, nil
}
