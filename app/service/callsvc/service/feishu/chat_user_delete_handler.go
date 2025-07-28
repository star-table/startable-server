package callsvc

import (
	"github.com/star-table/startable-server/common/core/util/json"
	"github.com/star-table/startable-server/common/core/util/slice"
	sdk_const "gitea.bjx.cloud/allstar/platform-sdk/consts"
	"github.com/star-table/startable-server/common/core/errs"
	"github.com/star-table/startable-server/common/extra/card"
	"github.com/star-table/startable-server/common/model/bo"
	"github.com/star-table/startable-server/common/model/vo/callvo"
	"github.com/star-table/startable-server/common/model/vo/orgvo"
	"github.com/star-table/startable-server/app/facade/orgfacade"
	"github.com/star-table/startable-server/app/facade/tablefacade"
	"github.com/star-table/startable-server/app/service/callsvc/domain"
	"github.com/star-table/startable-server/common/model/vo/callvo"
)
)

// 飞书事件 - 用户被移出群
// https://open.feishu.cn/document/uAjLw4CM/ukTMukTMukTM/reference/im-v1/chat-member-user/events/deleted

type ChatUserDeleteHandler struct{}

type ChatUserDeleteReq struct {
	req_vo.FeiShuBaseEventV2
	Event ChatUserDeleteReqEvent `json:"event"`
}

type ChatUserDeleteReqEvent struct {
	req_vo.FeiShuEventBaseInfo
	Users []req_vo.ChatUserItem `json:"users"`
}

// Handle 处理用户出群（退群或被移出）
// 针对任务群聊，出群的如果是协作人，则向操作人发送群卡片（仅操作人可见）
func (ChatUserDeleteHandler) Handle(data string) (string, errs.SystemErrorInfo) {
	var oriErr error
	log.Infof("[ChatUserDeleteHandler.Handle] request body: %s", data)
	req := ChatUserDeleteReq{}
	if oriErr = json.FromJson(data, &req); oriErr != nil {
		log.Errorf("[ChatUserDeleteHandler.Handle] err: %v", oriErr)
	}

	// 移出群后，触发一些调用
	// 如果是任务群聊，则向操作人（拉人的人）**单独**发送一条提示信息：你已将 xxx 移除出群聊，但他/她仍是该条任务的协作人...
	infoBox, err := domain.GetGroupChatHandleInfos(req.Header.TenantKey, &callvo.MessageReqData{
		TenantKey:  req.Header.TenantKey,
		OpenChatId: req.Event.ChatId,
		OpenId:     req.Event.OperatorId.OpenId,
	}, sdk_const.SourceChannelFeishu)
	if err != nil {
		log.Errorf("[ChatUserDeleteHandler.Handle] GetGroupChatHandleInfos err: %v, chatId: %s", err, req.Event.ChatId)
		return "", err
	}
	if !infoBox.IsIssueChat {
		log.Infof("[ChatUserDeleteHandler.Handle] this is not issue chat, chatId: %s", req.Event.ChatId)
		return "", nil
	}
	if domain.CheckIssueIsDelete(&infoBox.IssueInfo) {
		log.Infof("[ChatUserDeleteHandler.Handle] issue not exist or deleted, chatId: %s", req.Event.ChatId)
		return "", nil
	}
	// 检查被移出群的人是否是任务的协作人
	collaborateUidArr, err := tablefacade.GetDataCollaborateUserIds(infoBox.OrgInfo.OrgId, 0, infoBox.IssueInfo.DataId)
	if err != nil {
		log.Errorf("[ChatUserDeleteHandler.Handle] GetDataCollaborateUserIds err: %v, chatId: %s", err, req.Event.ChatId)
		return "", err
	}
	delUserOpenIds := make([]string, 0, len(req.Event.Users))
	for _, user := range req.Event.Users {
		delUserOpenIds = append(delUserOpenIds, user.UserId.OpenId)
	}
	usersResp := orgfacade.GetBaseUserInfoByEmpIdBatch(orgvo.GetBaseUserInfoByEmpIdBatchReqVo{
		OrgId: infoBox.OrgInfo.OrgId,
		Input: orgvo.GetBaseUserInfoByEmpIdBatchReqVoInput{
			OpenIds: delUserOpenIds,
		},
	})
	if usersResp.Failure() {
		log.Errorf("[ChatUserDeleteHandler.Handle] GetBaseUserInfoByEmpIdBatch err: %v, openIds: %s",
			usersResp.Error(), json.ToJsonIgnoreError(delUserOpenIds))
		return "", errs.BuildSystemErrorInfo(errs.UserNotFoundError, usersResp.Error())
	}
	delUsersMap := make(map[string]bo.BaseUserInfoBo, len(delUserOpenIds))
	for _, user := range usersResp.Data {
		delUsersMap[user.OutUserId] = user
	}
	userNamesIsCollaborate := make([]string, 0, len(req.Event.Users))
	for _, user := range req.Event.Users {
		if tmpUser, ok := delUsersMap[user.UserId.OpenId]; ok {
			isCollaborateUser, _ := slice.Contain(collaborateUidArr, tmpUser.UserId)
			// 是协作人，才收集其名字
			if isCollaborateUser {
				userNamesIsCollaborate = append(userNamesIsCollaborate, user.Name)
			}
		}
	}

	if len(userNamesIsCollaborate) > 0 {
		if err := card.SendCardFeiShuIssueChatUserDeleteToOpUser(req.Event.ChatId, infoBox, userNamesIsCollaborate); err != nil {
			log.Errorf("[ChatUserDeleteHandler.Handle] SendCardFeiShuIssueChatUserJoinToOpUser err: %v, chatId: %s", err,
				req.Event.ChatId)
			return "", err
		}
	}

	return "", nil
}
