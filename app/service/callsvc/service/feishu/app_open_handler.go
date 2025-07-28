package callsvc

import (
	"time"

	sdk_const "gitea.bjx.cloud/allstar/platform-sdk/consts"

	"gitea.bjx.cloud/allstar/feishu-sdk-golang/core/model/vo"
	"github.com/star-table/startable-server/app/facade/orgfacade"
	"github.com/star-table/startable-server/common/core/consts"
	"github.com/star-table/startable-server/common/core/errs"
	"github.com/star-table/startable-server/common/core/util/asyn"
	"github.com/star-table/startable-server/common/core/util/json"
	"github.com/star-table/startable-server/common/core/util/uuid"
	"github.com/star-table/startable-server/common/extra/feishu"
	"github.com/star-table/startable-server/common/library/cache"
	"github.com/star-table/startable-server/common/model/bo"
	"github.com/star-table/startable-server/common/model/vo/orgvo"
)

type AppOpenHandler struct{}

type AppOpenReq struct {
	Ts    string `json:"ts"`
	Uuid  string `json:"uuid"`
	Token string `json:"token"`
	Type  string `json:"type"`

	Event AppOpenReqData `json:"event"`
}

type AppOpenReqData struct {
	AppId      string               `json:"app_id"`
	TenantKey  string               `json:"tenant_key"`
	Type       string               `json:"type"`
	Applicants []AppOpenApplication `json:"applicants"`
	// 当应用被管理员安装时，返回此字段。如果是自动安装或由普通成员获取时，没有此字段
	Installer AppOpenInstaller `json:"installer"`
	// `app_open` 事件，结构体新加的字段。用于**普通成员安装**
	// https://open.feishu.cn/document/uAjLw4CM/ukTMukTMukTM/application-v6/event/app-first-enabled
	InstallerEmployee AppOpenInstaller `json:"installer_employee"`
}

type AppOpenApplication struct {
	OpenId string `json:"open_id"`
	UserId string `json:"user_id"`
}

type AppOpenInstaller struct {
	OpenId string `json:"open_id"`
	UserId string `json:"user_id"`
}

// InstallInfo 调用 fsInit 时传入的参数
type InstallInfo struct {
	// 安装者的 openId
	InstallerOpenId string `json:"installerOpenId"`
	// 普通用户安装时为 false；管理员安装为 true
	IsAdmin bool `json:"isAdmin"`
}

// 飞书处理
func (AppOpenHandler) Handle(data string) (string, errs.SystemErrorInfo) {
	appOpenReq := &AppOpenReq{}
	_ = json.FromJson(data, appOpenReq)

	log.Infof("飞书开通应用通知 %s", data)

	tenantKey := appOpenReq.Event.TenantKey
	installerOpenId := appOpenReq.Event.Installer.OpenId
	isAdmin := true
	// CustomEvent.Installer.OpenId 不存在时，则取 CustomEvent.InstallerEmployee.OpenId（普通成员安装方式）
	if installerOpenId == "" {
		installerOpenId = appOpenReq.Event.InstallerEmployee.OpenId
		isAdmin = false
	}
	err := FsInit(tenantKey, InstallInfo{
		InstallerOpenId: installerOpenId,
		IsAdmin:         isAdmin,
	})
	if err != nil {
		log.Error(err)
		return "", err
	}

	return "ok", nil
}

func PushBotToInstallerRemindMsg(tenantKey string, openIds []string) errs.SystemErrorInfo {
	if len(openIds) == 0 {
		log.Infof("openId为空，不用推送，tenant:%s", tenantKey)
		return nil
	}
	tenantClient, err := feishu.GetTenant(tenantKey)
	if err != nil {
		log.Error(err)
		return err
	}

	orgResp := orgfacade.GetBaseOrgInfoByOutOrgId(orgvo.GetBaseOrgInfoByOutOrgIdReqVo{OutOrgId: tenantKey})
	if orgResp.Failure() {
		return orgResp.Error()
	}

	cardInfo := GetCardInfoForFirst(orgResp.BaseOrgInfo.OrgId)

	size := 500 //每批次500条（飞书：消息请求体最大不能超过30k）
	offset := 0
	length := len(openIds)
	for {
		limit := offset + size
		if length < limit {
			limit = length
		}
		curOpenIds := openIds[offset:limit]

		cardInfo.OpenIds = curOpenIds
		resp, err1 := tenantClient.SendMessageBatch(cardInfo)
		if err1 != nil {
			log.Error(err1)
			return errs.FeiShuOpenApiCallError
		}
		if resp.Code != 0 {
			log.Errorf("发送应用安装成功通知失败，错误信息 %s，错误码 %d", resp.Msg, resp.Code)
			return errs.BuildSystemErrorInfo(errs.FeiShuOpenApiCallError)
		}

		if length <= limit {
			break
		}
		offset += size
	}

	log.Infof("发送应用安装成功通知成功")

	return nil
}

func FsInit(tenantKey string, installParam InstallInfo) errs.SystemErrorInfo {
	installerOpenId := installParam.InstallerOpenId
	//防止并发重复初始化
	uuid := uuid.NewUuid()
	suc, err := cache.TryGetDistributedLock(consts.FeiShuCorpInitKey+tenantKey, uuid)
	if err != nil {
		log.Error(err)
		return errs.BuildSystemErrorInfo(errs.TryDistributedLockError)
	}
	if !suc {
		log.Errorf("飞书组织初始化异常，外部组织key: %s, 原因:其他组织正在初始化", tenantKey)
		return errs.BuildSystemErrorInfo(errs.OrgNotInitError)
	}
	defer func() {
		if _, lockErr := cache.ReleaseDistributedLock(consts.FeiShuCorpInitKey+tenantKey, uuid); lockErr != nil {
			log.Error(lockErr)
		}
	}()

	orgName := "飞书平台组织"
	log.Infof("开始初始化组织, orgName: %s, tenantKey: %s, initParam: %s", orgName, tenantKey, json.ToJsonIgnoreError(installParam))

	//开始初始化组织
	orgInitResp := orgfacade.InitOrg(orgvo.InitOrgReqVo{
		InitOrg: bo.InitOrgBo{
			OutOrgId:      tenantKey,
			OrgName:       orgName,
			SourceChannel: sdk_const.SourceChannelFeishu,
			OutOrgOwnerId: installerOpenId,
		},
	})
	if orgInitResp.Failure() {
		log.Error(orgInitResp.Message)
		//组织已初始化不用推送消息
		if orgInitResp.Code != errs.OrgNotNeedInitError.Code() {
			PushOrgInitErrorNotice(tenantKey, installerOpenId)
		}

		return orgInitResp.Error()
	} else {
		asyn.Execute(func() {
			time.Sleep(3 * time.Second)
			//surplusOpenIds, err := feishu.GetScopeOpenIdsLimit(tenantKey)
			//if err != nil {
			//	log.Error(err)
			//	return
			//}

			surplusOpenIds := make([]string, 0)
			// 管理员初始化组织后，只需向安装者本人推送欢迎消息
			if installerOpenId != "" {
				surplusOpenIds = append(surplusOpenIds, installerOpenId)
			}
			if installParam.IsAdmin {
				sendHelpErr := PushBotToInstallerRemindMsg(tenantKey, surplusOpenIds)
				if sendHelpErr != nil {
					log.Errorf("[FsInit] tenantKey: %s, surplusOpenIds: %s, PushBotToInstallerRemindMsg err: %v", tenantKey, json.ToJsonIgnoreError(surplusOpenIds), sendHelpErr)
				}
			}
		})
	}

	return nil
}

func PushOrgInitErrorNotice(tenantKey string, openId string) {
	if openId == consts.BlankString {
		log.Infof("openId为空，不用推送，tenant:%s", tenantKey)
		return
	}
	tenantClient, err := feishu.GetTenant(tenantKey)
	if err != nil {
		log.Error(err)
		return
	}

	resp, err1 := tenantClient.SendMessage(vo.MsgVo{
		OpenId:  openId,
		MsgType: "interactive",
		Card: &vo.Card{
			Config: &vo.CardConfig{
				WideScreenMode: true,
			},
			Header: &vo.CardHeader{
				Title: &vo.CardHeaderTitle{
					Tag:     "plain_text",
					Content: "应用安装失败",
				},
			},
			Elements: []interface{}{
				vo.CardElementContentModule{
					Tag: "div",
					Fields: []vo.CardElementField{
						{
							Text: vo.CardElementText{
								Tag:     "lark_md",
								Content: "您可以通过企业后台管理关闭并重新开启应用来完成企业安装！",
							},
						},
					},
				},
			},
		},
	})
	if err1 != nil {
		log.Error(err1)
		return
	}
	if resp.Code != 0 {
		log.Errorf("发送企业初始化失败通知失败，错误信息 %s，错误码 %d", resp.Msg, resp.Code)
		return
	}
	log.Infof("发送企业初始化失败通知成功")

}

// GetCardInfoForFirst 管理员安装时向其推送；或者普通成员首次点开应用时，机器人推送帮助信息；
func GetCardInfoForFirst(orgId int64) vo.BatchMsgVo {
	return vo.BatchMsgVo{
		MsgType: "interactive",
		Card: &vo.Card{
			Header: &vo.CardHeader{
				Title: &vo.CardHeaderTitle{
					Tag:     "plain_text",
					Content: "Hi～欢迎使用✨ 极星协作✨",
				},
				Template: "blue",
			},
			Elements: []interface{}{
				vo.CardElementImageModule{
					Tag:    "img",
					ImgKey: "img_v2_c8d16837-ae0f-4712-8323-05a2f02a4f3g",
					Alt: vo.CardElementText{
						Tag:     "lark_md",
						Content: "极星协作",
					},
					CustomWidth: 580,
				},
				vo.CardElementContentModule{
					Tag: "div",
					Fields: []vo.CardElementField{
						{
							Text: vo.CardElementText{
								Tag: "lark_md",
								Content: "**极星协作是一款灵活安全的项目管理工具，可以帮您：**\n" +
									"📌  任务分配，进度跟踪\n" +
									"📃  文件管理，工时统计\n" +
									"🔐  权限设置，灵活安全\n" +
									"✍️  团队协作，提升效率",
							},
						},
					},
				},
				vo.CardElementActionModule{
					Tag: "action",
					Actions: []interface{}{
						vo.ActionButton{
							Tag: "button",
							Text: vo.CardElementText{
								Tag:     "plain_text",
								Content: "开始使用",
							},
							MultiUrl: &vo.CardElementUrl{
								PcUrl:      feishu.GetAppLinkWebPcWelcomeUrl(),
								IosUrl:     feishu.GetDefaultAppLink(orgId),
								AndroidUrl: feishu.GetDefaultAppLink(orgId),
							},
							Type: "primary",
						},
						vo.ActionButton{
							Tag: "button",
							Text: vo.CardElementText{
								Tag:     "plain_text",
								Content: "入门视频",
							},
							Url:  "https://startable.feishu.cn/docs/doccnPxGpclfzFA4Ea6Id42cjYd",
							Type: "primary",
						},
						vo.ActionButton{
							Tag: "button",
							Text: vo.CardElementText{
								Tag:     "plain_text",
								Content: "联系客服",
							},
							Url:  feishu.AppCustomerService,
							Type: "danger",
						},
					},
				},
			},
		},
	}
}

func PushBotInstallSuccessNotice(tenantKey, chatId string, openId string) errs.SystemErrorInfo {
	if openId == consts.BlankString {
		log.Infof("openId为空，不用推送，tenant:%s", tenantKey)
		return nil
	}
	tenantClient, err := feishu.GetTenant(tenantKey)
	if err != nil {
		log.Error(err)
		return err
	}
	orgResp := orgfacade.GetBaseOrgInfoByOutOrgId(orgvo.GetBaseOrgInfoByOutOrgIdReqVo{OutOrgId: tenantKey})
	if orgResp.Failure() {
		return orgResp.Error()
	}
	orgId := orgResp.BaseOrgInfo.OrgId

	resp, err1 := tenantClient.SendMessage(vo.MsgVo{
		ChatId:  chatId,
		OpenId:  openId,
		MsgType: "interactive",
		Card: &vo.Card{
			Header: &vo.CardHeader{
				Title: &vo.CardHeaderTitle{
					Tag:     "plain_text",
					Content: "极星协作应用安装成功",
				},
			},
			Elements: []interface{}{
				vo.CardElementContentModule{
					Tag: "div",
					Fields: []vo.CardElementField{
						{
							Text: vo.CardElementText{
								Tag:     "lark_md",
								Content: "极星是一个高效而稳定的项目协作平台",
							},
						},
					},
				},
				vo.CardElementContentModule{
					Tag: "div",
					Fields: []vo.CardElementField{
						{
							Text: vo.CardElementText{
								Tag:     "lark_md",
								Content: "移动端：[点击进入]($urlVal)",
								Href: map[string]vo.CardElementUrl{
									"urlVal": vo.CardElementUrl{
										Url:        feishu.GetDefaultAppLink(orgId),
										AndroidUrl: feishu.GetDefaultAppLink(orgId),
										IosUrl:     feishu.GetDefaultAppLink(orgId),
										PcUrl:      feishu.GetMobileOpenQRCodeAppLink(orgId),
									},
								},
							},
						},
					},
				},
				vo.CardElementContentModule{
					Tag: "div",
					Fields: []vo.CardElementField{
						{
							Text: vo.CardElementText{
								Tag:     "lark_md",
								Content: "PC端：[点击进入](" + feishu.GetAppLinkPcWelcomeUrl(orgId) + ")",
								//Href: map[string]vo.CardElementUrl{
								//	"urlVal":vo.CardElementUrl{
								//		Url:feishu.GetDefaultAppLink(),
								//		AndroidUrl:feishu.GetDefaultAppLink(),
								//		IosUrl:feishu.GetDefaultAppLink(),
								//		PcUrl:feishu.GetMobileOpenQRCodeAppLink(),
								//	},
								//},
							},
						},
					},
				},
				vo.CardElementContentModule{
					Tag: "div",
					Fields: []vo.CardElementField{
						{
							Text: vo.CardElementText{
								Tag:     "lark_md",
								Content: "您可以通过@极星协作机器人进行快捷操作",
							},
						},
					},
				},
				vo.CardElementContentModule{
					Tag: "div",
					Fields: []vo.CardElementField{
						{
							Text: vo.CardElementText{
								Tag:     "lark_md",
								Content: "发送** create 名称 **可以快速创建任务",
							},
						},
					},
				},
				//vo.CardElementContentModule{
				//	Tag: "div",
				//	Fields: []vo.CardElementField{
				//		{
				//			Text: vo.CardElementText{
				//				Tag:     "lark_md",
				//				Content: "发送** settings **可以进入配置修改页面",
				//			},
				//		},
				//	},
				//},
				vo.CardElementContentModule{
					Tag: "div",
					Fields: []vo.CardElementField{
						{
							Text: vo.CardElementText{
								Tag:     "lark_md",
								Content: "发送** notification **可以进入通知设置修改页面",
							},
						},
					},
				},
				vo.CardElementContentModule{
					Tag: "div",
					Fields: []vo.CardElementField{
						{
							Text: vo.CardElementText{
								Tag:     "lark_md",
								Content: "在使用中有任何问题，均可以直接联系我们。【联系电话：18001735738】",
							},
						},
					},
				},
				vo.CardElementContentModule{
					Tag: "div",
					Fields: []vo.CardElementField{
						{
							Text: vo.CardElementText{
								Tag:     "lark_md",
								Content: "详细使用说明 **[帮助手册](" + feishu.AppGuide + ")**",
							},
						},
					},
				},
				vo.CardElementContentModule{
					Tag: "div",
					Fields: []vo.CardElementField{
						{
							Text: vo.CardElementText{
								Tag:     "lark_md",
								Content: "可以随时联系 **[在线客服](" + feishu.AppCustomerService + ")**",
							},
						},
					},
				},
			},
		},
	})
	if err1 != nil {
		log.Error(err1)
		return errs.FeiShuOpenApiCallError
	}
	if resp.Code != 0 {
		log.Errorf("发送应用安装成功通知失败，错误信息 %s，错误码 %d, openId %s", resp.Msg, resp.Code, openId)
		return errs.BuildSystemErrorInfo(errs.FeiShuOpenApiCallError)
	}
	log.Infof("发送应用安装成功通知成功")

	return nil
}
