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

// 飞书事件 - 用户进群
// https://open.feishu.cn/document/uAjLw4CM/ukTMukTMukTM/reference/im-v1/chat-member-user/events/added

type ChatUserJoinHandler struct{}

type ChatUserJoinReq struct {
	req_vo.FeiShuBaseEventV2
	Event ChatUserJoinReqEvent `json:"event"`
}

type ChatUserJoinReqEvent struct {
	req_vo.FeiShuEventBaseInfo
	Users []req_vo.ChatUserItem `json:"users"`
}

// Handle 处理用户进群
// 针对任务群聊，进群的如果不是协作人，则向邀请人发送群卡片（仅操作人可见）
func (ChatUserJoinHandler) Handle(data string) (string, errs.SystemErrorInfo) {
	var oriErr error
	log.Infof("[ChatUserJoinHandler.Handle] request body: %s", data)
	req := ChatUserJoinReq{}
	if oriErr = json.FromJson(data, &req); oriErr != nil {
		log.Errorf("[ChatUserJoinHandler.Handle] err: %v", oriErr)
	}

	// 进群后，触发一些调用
	// 如果是任务群聊，则向操作人（拉人的人）**单独**发送一条提示信息
	infoBox, err := domain.GetGroupChatHandleInfos(req.Header.TenantKey, &callvo.MessageReqData{
		TenantKey:  req.Header.TenantKey,
		OpenChatId: req.Event.ChatId,
		OpenId:     req.Event.OperatorId.OpenId,
	}, sdk_const.SourceChannelFeishu)
	if err != nil {
		log.Errorf("[ChatUserJoinHandler.Handle] GetGroupChatHandleInfos err: %v, chatId: %s", err, req.Event.ChatId)
		return "", err
	}
	if !infoBox.IsIssueChat {
		log.Infof("[ChatUserJoinHandler.Handle] this is not issue chat, chatId: %s", req.Event.ChatId)
		return "", nil
	}
	if domain.CheckIssueIsDelete(&infoBox.IssueInfo) {
		log.Infof("[ChatUserJoinHandler.Handle] issue not exist or deleted, chatId: %s", req.Event.ChatId)
		return "", nil
	}
	// 检查被拉的人是否是任务的协作人
	collaborateUidArr, err := tablefacade.GetDataCollaborateUserIds(infoBox.OrgInfo.OrgId, 0, infoBox.IssueInfo.DataId)
	if err != nil {
		log.Errorf("[ChatUserJoinHandler.Handle] GetDataCollaborateUserIds err: %v, chatId: %s", err, req.Event.ChatId)
		return "", err
	}
	joinUserOpenIds := make([]string, 0, len(req.Event.Users))
	for _, user := range req.Event.Users {
		joinUserOpenIds = append(joinUserOpenIds, user.UserId.OpenId)
	}
	joinUsersResp := orgfacade.GetBaseUserInfoByEmpIdBatch(orgvo.GetBaseUserInfoByEmpIdBatchReqVo{
		OrgId: infoBox.OrgInfo.OrgId,
		Input: orgvo.GetBaseUserInfoByEmpIdBatchReqVoInput{
			OpenIds: joinUserOpenIds,
		},
	})
	if joinUsersResp.Failure() {
		log.Errorf("[ChatUserJoinHandler.Handle] GetBaseUserInfoByEmpIdBatch err: %v, openIds: %s",
			joinUsersResp.Error(), json.ToJsonIgnoreError(joinUserOpenIds))
		return "", errs.BuildSystemErrorInfo(errs.UserNotFoundError, joinUsersResp.Error())
	}
	joinUsersMap := make(map[string]bo.BaseUserInfoBo, len(joinUserOpenIds))
	for _, user := range joinUsersResp.Data {
		joinUsersMap[user.OutUserId] = user
	}
	userNamesIsNotCollaborate := make([]string, 0, len(req.Event.Users))
	for _, user := range req.Event.Users {
		if tmpUser, ok := joinUsersMap[user.UserId.OpenId]; ok {
			isCollaborateUser, _ := slice.Contain(collaborateUidArr, tmpUser.UserId)
			// 不是协作人，才收集其名字
			if !isCollaborateUser {
				userNamesIsNotCollaborate = append(userNamesIsNotCollaborate, user.Name)
			}
		}
	}
	if len(userNamesIsNotCollaborate) > 0 {
		if err := card.SendCardFeiShuIssueChatUserJoinToOpUser(req.Event.ChatId, infoBox, userNamesIsNotCollaborate); err != nil {
			log.Errorf("[ChatUserJoinHandler.Handle] SendCardFeiShuIssueChatUserJoinToOpUser err: %v, chatId: %s", err,
				req.Event.ChatId)
			return "", err
		}
	}

	return "", nil
}
