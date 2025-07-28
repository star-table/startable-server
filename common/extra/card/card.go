package card

import (
	"fmt"
	"strings"
	"unicode"

	fsSdkVo "gitea.bjx.cloud/allstar/feishu-sdk-golang/core/model/vo"
	sdk_const "gitea.bjx.cloud/allstar/platform-sdk/consts"
	"github.com/star-table/startable-server/common/core/consts"
	"github.com/star-table/startable-server/common/core/util/str"
	"github.com/star-table/startable-server/common/model/vo/commonvo"
)

func GenerateFeiShuCard(meta *commonvo.CardMeta) *fsSdkVo.Card {
	cd := &fsSdkVo.Card{}

	if meta.IsWide {
		cd.Config = &fsSdkVo.CardConfig{
			WideScreenMode: true,
		}
	} else {
		cd.Config = &fsSdkVo.CardConfig{
			WideScreenMode: false,
		}
	}

	// 标题
	cd.Header = &fsSdkVo.CardHeader{
		Title: &fsSdkVo.CardHeaderTitle{
			Tag:     "plain_text",
			Content: meta.Title,
		},
	}
	switch meta.Level {
	case consts.CardLevelInfo:
		cd.Header.Template = consts.FsCardTitleColorBlue
	case consts.CardLevelWarning:
		cd.Header.Template = consts.FsCardTitleColorOrange
	}

	// 预处理
	maxCh := 0
	maxEn := 0
	for _, div := range meta.Divs {
		for _, field := range div.Fields {
			ch, en := calcFieldLength(field.Key)
			if ch > maxCh {
				maxCh = ch
			}
			if en > maxEn {
				maxEn = en
			}
		}
	}

	// 内容
	for _, div := range meta.Divs {
		d := &fsSdkVo.CardElementContentModule{
			Tag:    "div",
			Fields: []fsSdkVo.CardElementField{},
		}
		cd.Elements = append(cd.Elements, d)

		for _, field := range div.Fields {
			d.Fields = append(d.Fields, fsSdkVo.CardElementField{
				Text: fsSdkVo.CardElementText{
					Tag:     "lark_md",
					Content: generateLine(sdk_const.SourceChannelFeishu, field, maxCh, maxEn),
				},
			})
		}
	}

	// Actions
	cd.Elements = append(cd.Elements, meta.FsActionElements...)

	return cd
}

func GenerateWeiXinCard(meta *commonvo.CardMeta) *commonvo.WxCard {
	var markdowns []string

	// 标题
	switch meta.Level {
	case consts.CardLevelInfo:
		markdowns = append(markdowns, fmt.Sprintf(consts.MarkdownBold, meta.Title))
	case consts.CardLevelWarning:
		markdowns = append(markdowns, fmt.Sprintf(consts.MarkdownColorTitle, consts.WxCardTitleColorWarning, meta.Title))
	}

	// 预处理
	maxCh := 0
	maxEn := 0
	for _, div := range meta.Divs {
		for _, field := range div.Fields {
			ch, en := calcFieldLength(field.Key)
			if ch > maxCh {
				maxCh = ch
			}
			if en > maxEn {
				maxEn = en
			}
		}
	}

	// 内容
	for _, div := range meta.Divs {
		for _, field := range div.Fields {
			markdowns = append(markdowns, generateLine(sdk_const.SourceChannelWeixin, field, maxCh, maxEn))
		}
	}

	// Actions
	markdowns = append(markdowns, meta.ActionMarkdowns...)

	return &commonvo.WxCard{
		Content: strings.Join(markdowns, consts.MarkdownBr),
	}
}

func GenerateDingTalkCard(meta *commonvo.CardMeta) *commonvo.DtCard {
	cd := &commonvo.DtCard{}
	var markdowns []string

	// 标题
	cd.Title = meta.Title

	// 预处理
	maxCh := 0
	maxEn := 0
	for _, div := range meta.Divs {
		for _, field := range div.Fields {
			ch, en := calcFieldLength(field.Key)
			if ch > maxCh {
				maxCh = ch
			}
			if en > maxEn {
				maxEn = en
			}
		}
	}

	// 内容
	for _, div := range meta.Divs {
		for _, field := range div.Fields {
			markdowns = append(markdowns, generateLine(sdk_const.SourceChannelDingTalk, field, maxCh, maxEn))
		}
	}

	// Actions
	markdowns = append(markdowns, meta.ActionMarkdowns...)

	cd.Content = strings.Join(markdowns, consts.MarkdownBr)
	return cd
}

func generateLine(sourceChannel string, field *commonvo.CardField, maxCh, maxEn int) string {
	if field.Value == "" {
		// 整行展示key的时候不存在对齐问题
		return field.Key
	} else {
		var builder strings.Builder

		// key
		if sourceChannel == sdk_const.SourceChannelFeishu {
			builder.WriteString(fmt.Sprintf(consts.MarkdownBold, field.Key))
		} else {
			builder.WriteString(field.Key)
		}

		// padding
		ch, en := calcFieldLength(field.Key)
		for i := 0; i < maxCh+2-ch; i++ {
			builder.WriteString(consts.CardPaddingChBlank)
		}
		for i := 0; i < maxEn-en; i++ {
			builder.WriteString(consts.CardPaddingEnBlank)
		}

		// value
		if sourceChannel == sdk_const.SourceChannelWeixin {
			if strings.Contains(field.Value, consts.MarkDownDel) {
				field.Value = strings.Replace(field.Value, consts.MarkDownDel, "", -1)
			}
			// 去除评论@标签，企微不支持评论@
			field.Value = str.TrimComment(field.Value)
			switch field.Level {
			case consts.CardLevelInfo:
				builder.WriteString(fmt.Sprintf(consts.MarkdownColor, consts.WxCardTitleColorComment, field.Value))
			case consts.CardLevelWarning:
				builder.WriteString(fmt.Sprintf(consts.MarkdownColor, consts.WxCardTitleColorWarning, field.Value))
			}
		} else if sourceChannel == sdk_const.SourceChannelDingTalk {
			// 钉钉的  markdown删除~~用不了，特殊处理一下
			if strings.Contains(field.Value, consts.MarkDownDel) {
				field.Value = strings.Replace(field.Value, consts.MarkDownDel, "", -1)
			}
			builder.WriteString(field.Value)
		} else {
			builder.WriteString(field.Value)
		}

		return builder.String()
	}
}

func calcFieldLength(key string) (int, int) {
	rs := []rune(key)
	ch := 0
	en := 0
	for _, r := range rs {
		if unicode.Is(unicode.Han, r) {
			ch += 1
		} else {
			en += 1
		}
	}
	return ch, en
}
