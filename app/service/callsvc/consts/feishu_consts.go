package callsvc

import "gitea.bjx.cloud/allstar/feishu-sdk-golang/core/model/vo"

var FsHelpPostMsg = &vo.MsgPost{
	ZhCn: &vo.MsgPostValue{
		Title: "欢迎使用极星协作",
		Content: []interface{}{
			[]interface{}{
				vo.MsgPostContentText{
					Tag:      "text",
					UnEscape: true,
					Text:     "您可以通过@极星协作机器人发送 create 指令快速创建任务",
				},
			},
			[]interface{}{
				vo.MsgPostContentText{
					Tag:      "text",
					UnEscape: true,
					Text:     "同时可以@极星协作机器人发送 settings 指令变更默认项目",
				},
			},
			[]interface{}{
				vo.MsgPostContentText{
					Tag:      "text",
					UnEscape: true,
					Text:     "在使用过程中遇到任务和问题，您可以查看帮助手册",
				},
			},
			[]interface{}{
				vo.MsgPostContentText{
					Tag:      "text",
					UnEscape: true,
					Text:     "配置文档：",
				},
				vo.MsgPostContentA{
					Tag:  "a",
					Text: "点我查看",
					Href: "https://polaris.feishu.cn/docs/doccn56SrnQKluE4VnLnL4MTMKe",
				},
			},
			[]interface{}{
				vo.MsgPostContentText{
					Tag:      "text",
					UnEscape: true,
					Text:     "操作手册：",
				},
				vo.MsgPostContentA{
					Tag:  "a",
					Text: "点我查看",
					Href: "https://polaris.feishu.cn/docs/doccn6Vp6ysqoYQJAWcj6nCYyzf",
				},
			},
		},
	},
}
