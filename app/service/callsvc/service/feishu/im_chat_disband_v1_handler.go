package callsvc

import (
	sdk_const "gitea.bjx.cloud/allstar/platform-sdk/consts"
	"github.com/star-table/startable-server/common/core/consts"
	"github.com/star-table/startable-server/common/core/errs"
	"github.com/star-table/startable-server/common/core/util/json"
	"github.com/star-table/startable-server/common/model/vo/callvo"
)

// 飞书事件 - 群解散
// https://open.feishu.cn/document/uAjLw4CM/ukTMukTMukTM/reference/im-v1/chat/events/disbanded

type ImChatDisbandV1Handler struct{}

type ImChatDisabledV1Req struct {
	req_vo.FeiShuBaseEventV2
	Event ChatUserJoinReqEvent `json:"event"`
}

type ImChatDisabledV1ReqEvent struct {
	req_vo.FeiShuEventBaseInfo
}

// Handle 群解散
// 针对任务群聊，解散后，更新任务数据，将 outChatId 改为 `""`
func (ImChatDisbandV1Handler) Handle(data string) (string, errs.SystemErrorInfo) {
	var oriErr error
	log.Infof("[ImChatDisbandV1Handler.Handle] request body: %s", data)
	req := ChatUserJoinReq{}
	if oriErr = json.FromJson(data, &req); oriErr != nil {
		log.Errorf("[ImChatDisbandV1Handler.Handle] err: %v", oriErr)
	}

	infoBox, err := domain.GetGroupChatHandleInfos(req.Header.TenantKey, &callvo.MessageReqData{
		TenantKey:  req.Header.TenantKey,
		OpenChatId: req.Event.ChatId,
		OpenId:     req.Event.OperatorId.OpenId,
	}, sdk_const.SourceChannelFeishu)
	if err != nil {
		log.Errorf("[ImChatDisbandV1Handler.Handle] GetGroupChatHandleInfos err: %v, chatId: %s", err, req.Event.ChatId)
		return "", err
	}
	if !infoBox.IsIssueChat {
		log.Infof("[ImChatDisbandV1Handler.Handle] this is not issue chat, chatId: %s", req.Event.ChatId)
		return "", nil
	}
	if val, ok := infoBox.IssueInfo.LessData[consts.BasicFieldOutChatId]; ok {
		chatId := val.(string)
		if chatId == req.Event.ChatId {
			if err := domain.SaveIssueForClearChatId(infoBox.OrgInfo.OrgId, &infoBox.IssueInfo); err != nil {
				log.Infof("[ImChatDisbandV1Handler.Handle] SaveIssueForClearChatId err: %v, chatId: %s", err, req.Event.ChatId)
				return "", err
			}
		}
	} else {
		log.Infof("[ImChatDisbandV1Handler.Handle] 没有找到群聊id issueId: %d", infoBox.IssueInfo.Id)
		return "", nil
	}

	return "", nil
}
