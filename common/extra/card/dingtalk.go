package card

import (
	"strings"

	"gitea.bjx.cloud/allstar/dingtalk-sdk-golang"
	dingRequest "gitea.bjx.cloud/allstar/dingtalk-sdk-golang/request"
	platform_sdk "gitea.bjx.cloud/allstar/platform-sdk"
	sdk_const "gitea.bjx.cloud/allstar/platform-sdk/consts"
	"gitea.bjx.cloud/allstar/platform-sdk/sdk_interface"
	sdk_vo "gitea.bjx.cloud/allstar/platform-sdk/vo"
	"github.com/star-table/startable-server/common/core/config"
	"github.com/star-table/startable-server/common/core/consts"
	"github.com/star-table/startable-server/common/core/util/json"
	ding "github.com/star-table/startable-server/common/extra/dingtalk"
	"github.com/star-table/startable-server/common/model/vo/commonvo"
)

func SendCardMetaDingTalk(corpId string, meta *commonvo.CardMeta, openIds []string) {
	if len(openIds) == 0 || meta == nil {
		return
	}
	client, errSys := platform_sdk.GetClient(sdk_const.SourceChannelDingTalk, corpId)
	if errSys != nil {
		log.Errorf("[SendCardMetaDingTalk] GetClient err:%v", errSys)
		return
	}

	SendCardMetaDingTalkByDingTalk(client, meta, corpId, openIds)
}

func SendCardMetaDingTalkByDingTalk(sdk sdk_interface.Sdk, meta *commonvo.CardMeta, corpId string, openIds []string) {
	if len(openIds) == 0 || meta == nil {
		return
	}

	card := GenerateDingTalkCard(meta)
	dingConfig := config.GetConfig().DingTalk
	if dingConfig == nil {
		log.Error("[SendCardDingTalkByDingTalk] GetConfig DingTalk not found")
		return
	}

	orgInfo, sdkError := sdk.GetOrgInfo(&sdk_vo.OrgInfoReq{CorpId: corpId})
	if sdkError != nil {
		log.Errorf("[SendCardMetaDingTalkByDingTalk] GetOrgInfo err:%v", sdkError)
		return
	}

	dingTalk := sdk.GetOriginClient().(*dingtalk.DingTalk)

	log.Infof("dingConfig:%v", json.ToJsonIgnoreError(dingConfig))

	toUsers := strings.Join(openIds, ",")
	msgResp, sdkErr := dingTalk.SendTemplateMessage(&dingRequest.SendTemplateMessage{
		AgentId:    orgInfo.AgentId, // isv应用 每个corpId下面的agentId会不一样
		TemplateId: dingConfig.TemplateId,
		UserIdList: toUsers,
		Data:       json.ToJsonIgnoreError(card),
	})
	if sdkErr != nil {
		log.Errorf("[SendCardDingTalkByDingTalk] err: %v", sdkErr)
		return
	}
	if msgResp.Code != 0 {
		log.Errorf("[SendCardDingTalkByDingTalk] SendTemplateMessage, code: %d, errMsg: %s", msgResp.Code, msgResp.Msg)
		return
	}
}

func SendCardDingTalk(corpId string, card *commonvo.Card, openIds []string) {
	if len(openIds) == 0 || card == nil || len(card.MarkdownElements) == 0 {
		return
	}
	dingClient, errSys := ding.GetDingTalk(corpId)
	if errSys != nil {
		log.Errorf("[SendCardDingTalk] GetDingTalk err:%v", errSys)
		return
	}

	SendCardDingTalkByDingTalk(dingClient, card, openIds)
}

func SendCardDingTalkByDingTalk(ding *dingtalk.DingTalk, card *commonvo.Card, openIds []string) {
	if len(openIds) == 0 || card == nil || len(card.MarkdownElements) == 0 {
		return
	}

	dingConfig := config.GetConfig().DingTalk
	if dingConfig == nil {
		log.Error("[SendCardDingTalkByDingTalk] GetConfig DingTalk not found")
		return
	}

	log.Infof("dingConfig:%v", json.ToJsonIgnoreError(dingConfig))

	dingMsg := commonvo.DtCard{
		Title:   card.Title,
		Content: strings.Join(card.MarkdownElements, consts.MarkdownBr),
	}
	toUsers := strings.Join(openIds, ",")
	msgResp, sdkErr := ding.SendTemplateMessage(&dingRequest.SendTemplateMessage{
		//AgentId:    dingConfig.AgentId, // isv应用 每个corpId下面的agentId会不一样
		TemplateId: dingConfig.TemplateId,
		UserIdList: toUsers,
		Data:       json.ToJsonIgnoreError(dingMsg),
	})
	if sdkErr != nil {
		log.Errorf("[SendCardDingTalkByDingTalk] err: %v", sdkErr)
		return
	}
	if msgResp.Code != 0 {
		log.Errorf("[SendCardDingTalkByDingTalk] SendTemplateMessage, code: %d, errMsg: %s", msgResp.Code, msgResp.Msg)
		return
	}
}
