package card

import (
	"fmt"
	"strings"

	"github.com/star-table/startable-server/common/core/consts"
	"github.com/star-table/startable-server/common/model/vo/commonvo"

	"github.com/star-table/startable-server/common/core/util/json"
	"github.com/star-table/startable-server/common/extra/weixin"
	wework "github.com/go-laoji/wecom-go-sdk"
)

func SendCardMetaWeiXin(corpId string, meta *commonvo.CardMeta, openIds []string) {
	if len(openIds) == 0 || meta == nil {
		return
	}
	wx, sysErr := weixin.GetWeWork(corpId)
	if sysErr != nil {
		log.Error(sysErr)
		return
	}
	SendCardMetaWeiXinByWeWork(wx, meta, openIds)
}

func SendCardMetaWeiXinByWeWork(wx *commonvo.WeWork, meta *commonvo.CardMeta, openIds []string) {
	if len(openIds) == 0 || meta == nil {
		return
	}

	card := GenerateWeiXinCard(meta)
	toUsers := strings.Join(openIds, "|")
	message := wework.MarkDownMessage{
		Message: wework.Message{
			ToUser: toUsers,
		},
		MarkDown: wework.Text{
			Content: card.Content,
		},
	}
	log.Infof("[SendCardMetaWeiXin] MessageSend to: %v, markdown: [%v]", message.Message.ToUser, message.MarkDown.Content)
	resp := wx.Sdk.MessageSend(wx.CorpIdUint, message)
	if resp.ErrCode != 0 {
		log.Errorf("[SendCardMetaWeiXin] MessageSend code: %d, errMsg: %s, card: %s", resp.ErrCode, resp.ErrorMsg, json.ToJsonIgnoreError(message))
		return
	}
}

func SendCardWeiXin(corpId string, card *commonvo.Card, openIds []string) {
	if len(openIds) == 0 || card == nil || len(card.MarkdownElements) == 0 {
		return
	}
	wx, sysErr := weixin.GetWeWork(corpId)
	if sysErr != nil {
		log.Error(sysErr)
		return
	}
	SendCardWeiXinByWeWork(wx, card, openIds)
}

func SendCardWeiXinByWeWork(wx *commonvo.WeWork, card *commonvo.Card, openIds []string) {
	if len(openIds) == 0 || card == nil || len(card.MarkdownElements) == 0 {
		return
	}

	markdownElements := append([]string{fmt.Sprintf(consts.MarkdownTitle, card.Title)}, card.MarkdownElements...)

	toUsers := strings.Join(openIds, "|")
	message := wework.MarkDownMessage{
		Message: wework.Message{
			ToUser: toUsers,
		},
		MarkDown: wework.Text{
			Content: strings.Join(markdownElements, consts.MarkdownBr),
		},
	}
	log.Infof("[SendCardWeiXin] MessageSend to: %v, markdown: [%v]", message.Message.ToUser, message.MarkDown.Content)
	resp := wx.Sdk.MessageSend(wx.CorpIdUint, message)
	if resp.ErrCode != 0 {
		log.Errorf("[SendCardWeiXin] MessageSend code: %d, errMsg: %s, card: %s", resp.ErrCode, resp.ErrorMsg, json.ToJsonIgnoreError(message))
		return
	}
}
