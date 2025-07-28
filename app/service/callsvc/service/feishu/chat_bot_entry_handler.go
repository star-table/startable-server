package callsvc

import (
	"github.com/star-table/startable-server/app/facade/orgfacade"
	"github.com/star-table/startable-server/common/core/errs"
	"github.com/star-table/startable-server/common/core/util/json"
	"github.com/star-table/startable-server/common/model/vo/orgvo"
)

// 机器人进群
type ChatBotEntryHandler struct{}

type ChatBotEntryReq struct {
	req_vo.FeiShuBaseEventV2
	Event req_vo.FeiShuEventBaseInfo `json:"event"`
}

// Handle 处理机器人进群
// https://open.feishu.cn/document/uAjLw4CM/ukTMukTMukTM/reference/im-v1/chat-member-bot/events/added
func (ChatBotEntryHandler) Handle(data string) (string, errs.SystemErrorInfo) {
	req := &ChatBotEntryReq{}
	_ = json.FromJson(data, req)
	log.Infof("[ChatBotEntryHandler.Handle] 机器人进群事件 %s", data)
	tenantKey := req.Header.TenantKey
	chatId := req.Event.ChatId

	orgResp := orgfacade.GetBaseOrgInfoByOutOrgId(orgvo.GetBaseOrgInfoByOutOrgIdReqVo{OutOrgId: tenantKey})
	if orgResp.Failure() {
		log.Errorf("[ChatBotEntryHandler.Handle] GetBaseOrgInfoByOutOrgId err: %v, tenantKey: %v", orgResp.Error(),
			tenantKey)
		return "", orgResp.Error()
	}

	isIssueChat, err := domain.CheckIsIssueChat(orgResp.BaseOrgInfo.OrgId, chatId)
	if err != nil {
		log.Errorf("[ChatBotEntryHandler.Handle] GetBaseOrgInfoByOutOrgId err: %v, tenantKey: %v", err,
			tenantKey)
		return "", err
	}
	if isIssueChat {
		// do nothing
	} else {
		err := PushBotNoticeNotBindProject(tenantKey, chatId)
		if err != nil {
			log.Error(err)
			return "error", err
		}
	}

	return "ok", nil
}
