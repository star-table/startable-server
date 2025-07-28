package callsvc

import (
	"fmt"

	"gitea.bjx.cloud/allstar/dingtalk-sdk-golang"
	"gitea.bjx.cloud/allstar/dingtalk-sdk-golang/request"
	"gitea.bjx.cloud/allstar/platform-sdk/consts"

	platform_sdk "gitea.bjx.cloud/allstar/platform-sdk"
	"github.com/star-table/startable-server/common/core/util/json"
)

type ChatTemplateChangeHandler struct{}

type chatTemplateChangeEvent struct {
	EventType          string `json:"EventType"`
	CorpId             string `json:"CorpId"`
	Status             string `json:"Status"`
	OpenConversationId string `json:"OpenConversationId"`
	TemplateId         string `json:"TemplateId"`
	Operator           string `json:"Operator"`
}

func (c ChatTemplateChangeHandler) Handle(data req_vo.DingEventBizData) error {
	msg := &chatTemplateChangeEvent{}
	json.FromJson(data.BizData, msg)

	fmt.Println("callback->", data)

	if msg != nil {
		// 模版启用通知
		client, err := platform_sdk.GetClient(consts.SourceChannelDingTalk, msg.CorpId)
		if err != nil {
			log.Errorf("[ChatTemplateChangeHandler] 获取钉钉客户端失败:%v, corpId:%s, templateId", err, msg.CorpId)
			return err
		}
		dingTalk := client.GetOriginClient().(*dingtalk.DingTalk)
		resp, err := dingTalk.SceneGroupTemplateApply(&request.SceneGroupTemplate{
			OwnerUserId:        msg.Operator,
			TemplateId:         msg.TemplateId,
			OpenConversationId: msg.OpenConversationId,
		})
		if err != nil {
			log.Errorf("[ChatTemplateChangeHandler]")
			return err
		}
		if resp.Code != 0 {
			log.Errorf("[ChatTemplateChangeHandler] SceneGroupTemplateApply failed")
		}
	}

	return nil
}
