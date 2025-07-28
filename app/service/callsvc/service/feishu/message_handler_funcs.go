package callsvc

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/star-table/startable-server/app/facade/orgfacade"
	"github.com/star-table/startable-server/app/facade/projectfacade"
	"github.com/star-table/startable-server/common/model/vo/commonvo"

	v1 "gitea.bjx.cloud/LessCode/interface/golang/table/v1"
	"gitea.bjx.cloud/allstar/feishu-sdk-golang/core/model/vo"
	"gitea.bjx.cloud/allstar/feishu-sdk-golang/sdk"
	sdk_const "gitea.bjx.cloud/allstar/platform-sdk/consts"
	"github.com/star-table/startable-server/common/core/consts"
	"github.com/star-table/startable-server/common/core/errs"
	"github.com/star-table/startable-server/common/core/util/json"
	"github.com/star-table/startable-server/common/extra/card"
	"github.com/star-table/startable-server/common/extra/feishu"
	"github.com/star-table/startable-server/common/model/bo"
	vo1 "github.com/star-table/startable-server/common/model/vo"
	"github.com/star-table/startable-server/common/model/vo/callvo"
	"github.com/star-table/startable-server/common/model/vo/orgvo"
	"github.com/star-table/startable-server/common/model/vo/projectvo"
	"github.com/spf13/cast"
)

// 用户指令：@负责人姓名 任务标题1
func instructAtUserNameWithIssueTitleFunc(tenantKey string, event callvo.MessageReqData, msg string, extraInfo map[string]interface{}) error {
	dealInfo, ok := extraInfo["dealInfo"].(*callvo.GroupChatHandleInfo)
	if !ok {
		log.Errorf("[instructAtUserNameWithIssueTitleFunc] get dealInfo err. tenantKey: %s", tenantKey)
		return nil
	}
	hasBindMultiPro, ok := extraInfo["hasBindMultiPro"].(bool)
	if !ok || hasBindMultiPro {
		// 发送卡片，提示不支持该群聊指令
		PushBotNoticeForNotSupportChatInstruction(tenantKey, &event, dealInfo)
		return nil
	}

	atUserOpenId := ""
	issueTitle := ""
	if tmpVal, ok := extraInfo["atUserOpenId"]; !ok {
		return errors.New("指令请求中[用户名]不合法。")
	} else {
		atUserOpenId = tmpVal.(string)
	}
	if tmpVal, ok := extraInfo["issueTitle"]; !ok {
		return errors.New("指令请求中[任务标题]不合法。")
	} else {
		issueTitle = tmpVal.(string)
	}

	dealInfo, err := domain.GetGroupChatHandleInfos(tenantKey, &event, sdk_const.SourceChannelFeishu)
	if err != nil {
		log.Errorf("[instructAtUserNameWithIssueTitleFunc] err: %v, tenantKey: %s", err, tenantKey)
		return err
	}

	if len(dealInfo.BindProjectIds) < 1 {
		err := sendSimpleNotice(tenantKey, event, map[string]interface{}{
			"isSimpleHelperInfo": false,
		})
		if err != nil {
			log.Error(err)
			return err
		}
	} else if len(dealInfo.BindProjectIds) > 1 {
		if err := PushBotNoticeForNotSupportChatInstruction(tenantKey, &event, dealInfo); err != nil {
			log.Errorf("[instructAtUserNameWithIssueTitleFunc] PushBotNoticeForNotSupportChatInstruction err: %v, tenantKey: %s", err, tenantKey)
			return err
		}
		return nil
	}
	// 只绑定**一个**项目的群聊，才能正常执行群聊指令
	// 调用接口处理指令
	resp := projectfacade.HandleGroupChatUserInsAtUserNameWithIssueTitle(projectvo.HandleGroupChatUserInsAtUserNameWithIssueTitleReqVo{
		Input: &projectvo.HandleGroupChatUserInsAtUserNameWithIssueTitleReq{
			OpenChatId:    event.OpenChatId,
			OpUserOpenId:  event.OpenId,
			AtUserOpenId:  atUserOpenId,
			IssueTitleStr: issueTitle,
		},
	})
	if resp.Failure() {
		log.Error(resp.Error())
		return resp.Error()
	}
	log.Infof("飞书群聊-处理用户指令-instructAtUserNameWithIssueTitleFunc-code: %v", resp.Code)
	return nil
}

// 用户指令：@负责人姓名
func instructAtUserNameFunc(tenantKey string, event callvo.MessageReqData, msg string, extraInfo map[string]interface{}) error {
	dealInfo, ok := extraInfo["dealInfo"].(*callvo.GroupChatHandleInfo)
	if !ok {
		log.Errorf("[instructAtUserNameFunc] get dealInfo err. tenantKey: %s", tenantKey)
		return nil
	}
	hasBindMultiPro, ok := extraInfo["hasBindMultiPro"].(bool)
	if !ok || hasBindMultiPro {
		// 发送卡片，提示不支持该群聊指令
		PushBotNoticeForNotSupportChatInstruction(tenantKey, &event, dealInfo)
		return nil
	}

	log.Infof("[instructAtUserNameFunc] tenantKey: %s, event: %s, extraInfo: %s", tenantKey, json.ToJsonIgnoreError(event),
		json.ToJsonIgnoreError(extraInfo))
	atUserOpenId := ""
	if tmpVal, ok := extraInfo["atUserOpenId"]; !ok {
		return errors.New("指令请求中[用户名]不合法。")
	} else {
		atUserOpenId = tmpVal.(string)
	}
	resp := projectfacade.HandleGroupChatUserInsAtUserName(projectvo.HandleGroupChatUserInsAtUserNameReqVo{
		Input: &projectvo.HandleGroupChatUserInsAtUserNameReq{
			OpenChatId:   event.OpenChatId,
			AtUserOpenId: atUserOpenId,
			OpUserOpenId: event.OpenId,
		},
	})
	if resp.Failure() {
		log.Error(resp.Error())
		return resp.Error()
	}
	log.Infof("飞书群聊-处理用户指令-instructAtUserNameFunc-code: %v", resp.Code)
	return nil
}

// 用户指令：项目任务
func instructUserInsProIssueFunc(tenantKey string, event callvo.MessageReqData, msg string, extraInfo map[string]interface{}) error {
	dealInfo, err := domain.GetGroupChatHandleInfos(tenantKey, &event, sdk_const.SourceChannelFeishu)
	if err != nil {
		log.Errorf("[instructUserInsProIssueFunc] err: %v, tenantKey: %s", err, tenantKey)
		return err
	}
	hasBindMultiPro, ok := extraInfo["hasBindMultiPro"].(bool)
	if !ok || hasBindMultiPro {
		// 发送卡片，提示不支持该群聊指令
		PushBotNoticeForNotSupportChatInstruction(tenantKey, &event, dealInfo)
		return nil
	}
	resp := projectfacade.HandleGroupChatUserInsProIssue(projectvo.HandleGroupChatUserInsProIssueReqVo{
		Input: &projectvo.HandleGroupChatUserInsProIssueReq{
			OpenChatId:    event.OpenChatId,
			OpenId:        event.OpenId,
			SourceChannel: sdk_const.SourceChannelFeishu,
		},
	})
	if resp.Failure() {
		log.Errorf("[instructUserInsProIssueFunc] err: %v, chatId: %s", resp.Error(), event.OpenChatId)
		return resp.Error()
	}
	log.Infof("[instructUserInsProIssueFunc] 飞书群聊-处理用户指令 resp code: %v", resp.Code)
	return nil
}

// 用户指令：项目进展
func instructUserInsProProgressFunc(tenantKey string, event callvo.MessageReqData, msg string, extraInfo map[string]interface{}) error {
	dealInfo, ok := extraInfo["dealInfo"].(*callvo.GroupChatHandleInfo)
	if !ok {
		log.Errorf("[instructUserInsProProgressFunc] get dealInfo err. tenantKey: %s", tenantKey)
		return nil
	}
	hasBindMultiPro, ok := extraInfo["hasBindMultiPro"].(bool)
	if !ok || hasBindMultiPro {
		// 发送卡片，提示不支持该群聊指令
		PushBotNoticeForNotSupportChatInstruction(tenantKey, &event, dealInfo)
		return nil
	}
	resp := projectfacade.HandleGroupChatUserInsProProgress(projectvo.HandleGroupChatUserInsProProgressReqVo{
		Input: &projectvo.HandleGroupChatUserInsProProgressReq{
			OpenChatId:    event.OpenChatId,
			OpenId:        event.OpenId,
			SourceChannel: sdk_const.SourceChannelFeishu,
		},
	})
	if resp.Failure() {
		log.Error(resp.Error())
		return resp.Error()
	}
	log.Infof("[instructUserInsProProgressFunc] 飞书群聊-处理用户指令-code: %v", resp.Code)
	return nil
}

// 群聊用户指令：项目设置
func instructUserInsProjectSettingsFunc(tenantKey string, event callvo.MessageReqData, msg string, extraInfo map[string]interface{}) error {
	//dealInfo, ok := extraInfo["dealInfo"].(*callvo.GroupChatHandleInfo)
	//if !ok {
	//	log.Errorf("[instructUserInsProjectSettingsFunc] get dealInfo err. tenantKey: %s", tenantKey)
	//	return nil
	//}
	resp := projectfacade.HandleGroupChatUserInsProjectSettings(projectvo.HandleGroupChatUserInsProjectSettingsReq{
		OpenChatId:    event.OpenChatId,
		SourceChannel: sdk_const.SourceChannelFeishu,
	})
	if resp.Failure() {
		log.Error(resp.Error())
		return resp.Error()
	}
	log.Infof("[instructUserInsProjectSettingsFunc] 飞书群聊-处理用户指令-code: %v", resp.Code)
	return nil
}

// 设置
func instructSettingsFunc(tenantKey string, event callvo.MessageReqData, msg string, extraInfo map[string]interface{}) error {
	tenantClient, err := feishu.GetTenant(tenantKey)
	if err != nil {
		log.Error(err)
		return err
	}
	PushBotSettingsNotice(*tenantClient, "", event.OpenChatId)

	return nil
}

// 创建任务
func instructCreateIssueFunc(tenantKey string, event callvo.MessageReqData, msg string, extraInfo map[string]interface{}) error {
	chatId := event.OpenChatId

	msg = strings.Trim(msg, " ")
	instructs := strings.SplitN(msg, MessageSeparator, 2)
	if len(instructs) < 2 {
		errMsg := "创建任务指令缺少任务标题"
		log.Error(errMsg)
		PushErrorNotice(0, tenantKey, "", chatId, errMsg)
		return nil
	}
	issueTitle := instructs[1]

	openId := event.OpenId

	orgInfoRespVo := orgfacade.GetBaseOrgInfoByOutOrgId(orgvo.GetBaseOrgInfoByOutOrgIdReqVo{OutOrgId: tenantKey})
	if orgInfoRespVo.Failure() {
		log.Error(orgInfoRespVo.Message)
		PushErrorNotice(0, tenantKey, "", chatId, "组织认证异常")
		return orgInfoRespVo.Error()
	}
	orgId := orgInfoRespVo.BaseOrgInfo.OrgId

	userInfoRespVo := orgfacade.GetBaseUserInfoByEmpId(orgvo.GetBaseUserInfoByEmpIdReqVo{
		OrgId: orgId,
		EmpId: openId,
	})
	if userInfoRespVo.Failure() {
		log.Error(userInfoRespVo.Message)
		PushErrorNotice(orgId, tenantKey, "", chatId, "用户认证异常")
		return userInfoRespVo.Error()
	}
	userId := userInfoRespVo.BaseUserInfo.UserId
	projectBo := &bo.ProjectAuthBo{
		Name:        consts.CardDefaultIssueProjectName,
		ProjectType: consts.ProjectTypeCommon2022V47,
	}

	lessCodeData := make(map[string]interface{}, 0)
	lessCodeData[consts.BasicFieldTitle] = issueTitle
	lessCodeData[consts.BasicFieldOwnerId] = []string{fmt.Sprintf("%s%d", consts.LcCustomFieldUserType, userId)}

	createIssueResp := projectfacade.CreateIssue(&projectvo.CreateIssueReqVo{
		UserId:        userId,
		OrgId:         orgId,
		SourceChannel: sdk_const.SourceChannelFeishu,
		//InputAppId:    projectBo.AppId,
		CreateIssue: vo1.CreateIssueReq{
			//ProjectID: userConfigBo.DefaultProjectId,
			//Title:              issueTitle,
			//OwnerID:            []int64{userId},
			//PriorityID:         defaultPriority.ID,
			//TableID:            tableIdStr,
			LessCreateIssueReq: lessCodeData,
		},
	})
	if createIssueResp.Failure() {
		log.Errorf("[instructCreateIssueFunc] err: %v, orgId: %d, userId: %d", createIssueResp.Message, orgId, userId)
		tipsMsg := createIssueResp.Message
		if createIssueResp.Code == errs.CommonUserCreateTaskLimit.Code() {
			tipsMsg = "标准版用户可创建记录上限为1000条（包括已删除记录），为不影响使用，可点击下方配置按钮前往升级版本。"
		}
		PushErrorNotice(orgId, tenantKey, "", chatId, tipsMsg)
		return createIssueResp.Error()
	} else {
		newIssue := &vo1.Issue{}
		newIssues := createIssueResp.Data
		if len(newIssues) > 0 {
			m := newIssues[0]
			newIssue.Title = cast.ToString(m[consts.BasicFieldTitle])
			newIssue.OrgID = cast.ToInt64(m[consts.BasicFieldOrgId])
			newIssue.ID = cast.ToInt64(m[consts.BasicFieldIssueId])
			newIssue.ProjectID = cast.ToInt64(m[consts.BasicFieldProjectId])
			newIssue.TableID = cast.ToString(m[consts.BasicFieldTableId])
		}
		//发送创建成功通知
		PushCreateIssueSuccessNotice(tenantKey, "", chatId, newIssue, userInfoRespVo.BaseUserInfo, projectBo)
	}
	return nil
}

func PushErrorNotice(orgId int64, outOrgId, openId string, chatId string, errMsg string) {
	cardMsg := card.ChatCreateIssueFailNotice(errMsg)

	errSys := card.PushCard(orgId, &commonvo.PushCard{
		OrgId:         orgId,
		OutOrgId:      outOrgId,
		SourceChannel: sdk_const.SourceChannelFeishu,
		ChatIds:       []string{chatId},
		OpenIds:       []string{openId},
		CardMsg:       cardMsg,
	})
	if errSys != nil {
		log.Errorf("[SendCardUrgeIssue] err:%v", errSys)
		return
	}
	log.Infof("发送机器人报错通知成功")
}

// 指令创建的不需要加退订
func PushCreateIssueSuccessNotice(outOrgId, openId string, chatId string, issueVo *vo1.Issue, userInfo *bo.BaseUserInfoBo, projectVo *bo.ProjectAuthBo) {
	title := issueVo.Title
	tableId, tableIdErr := strconv.ParseInt(issueVo.TableID, 10, 64)
	if tableIdErr != nil {
		log.Errorf("[GetCardReplyUrgeIssue] ParseInt failed:%v", tableIdErr)
		return
	}
	tableName := consts.DefaultTableName
	if tableId != 0 {
		tableInfoResp := projectfacade.GetTable(projectvo.GetTableInfoReq{
			OrgId:  issueVo.OrgID,
			UserId: issueVo.Creator,
			Input:  &v1.ReadTableRequest{TableId: tableId},
		})
		if tableInfoResp.Failure() {
			log.Error(tableInfoResp.Error())
			return
		}
		tableName = tableInfoResp.Data.Table.Name
	}

	issueLinks := projectfacade.GetIssueLinks(projectvo.GetIssueLinksReqVo{
		SourceChannel: sdk_const.SourceChannelFeishu,
		OrgId:         issueVo.OrgID,
		IssueId:       issueVo.ID,
	}).Data

	cardMsg := card.GetFsCardCreateNewIssue(title, projectVo.Name, tableName, []*bo.BaseUserInfoBo{userInfo}, issueLinks, *issueVo)

	errSys := card.PushCard(issueVo.OrgID, &commonvo.PushCard{
		OrgId:         issueVo.OrgID,
		OutOrgId:      outOrgId,
		SourceChannel: sdk_const.SourceChannelFeishu,
		OpenIds:       []string{openId},
		ChatIds:       []string{chatId},
		CardMsg:       cardMsg,
	})
	if errSys != nil {
		log.Errorf("[PushCreateIssueSuccessNotice] err:%v", errSys)
		return
	}
	log.Infof("发送机器人报错通知成功")
}

func PushBotSettingsNotice(tenantClient sdk.Tenant, openId string, chatId string) {
	resp, err1 := tenantClient.SendMessage(vo.MsgVo{
		OpenId:  openId,
		ChatId:  chatId,
		MsgType: "interactive",
		Card: &vo.Card{
			Config: &vo.CardConfig{
				WideScreenMode: true,
			},
			Header: &vo.CardHeader{
				Title: &vo.CardHeaderTitle{
					Tag:     "plain_text",
					Content: "您还没有完成对我的配置哦",
				},
			},
			Elements: []interface{}{
				vo.CardElementContentModule{
					Tag: "div",
					Fields: []vo.CardElementField{
						{
							Text: vo.CardElementText{
								Tag:     "lark_md",
								Content: "请先在桌面端完成配置，即可使用快速创建任务等功能 ^ ^ ！",
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
								Tag:     "lark_md",
								Content: "打开配置项",
							},
							Type: "default",
							Url:  feishu.GetDefaultProjectConfigureAppLink(),
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
		log.Errorf("发送机器人settings，错误信息 %s", resp.Msg)
		return
	}
	log.Infof("发送机器人settings通知成功")

}

// 个人机器人-设置
func instructNotification(tenantKey string, event callvo.MessageReqData, msg string, extraInfo map[string]interface{}) error {
	orgResp := orgfacade.GetBaseOrgInfoByOutOrgId(orgvo.GetBaseOrgInfoByOutOrgIdReqVo{OutOrgId: tenantKey})
	if orgResp.Failure() {
		return orgResp.Error()
	}
	orgId := orgResp.BaseOrgInfo.OrgId
	PushBotNotification(orgId, orgResp.BaseOrgInfo.OutOrgId, "", event.OpenChatId)

	return nil
}

// PushBotNotification 个人机器人-卡片（设置）推送
func PushBotNotification(orgId int64, outOrgId string, openId string, chatId string) {

	cardMsg := card.GetFsBotNotification(orgId)
	errSys := card.PushCard(orgId, &commonvo.PushCard{
		OrgId:         orgId,
		OutOrgId:      outOrgId,
		SourceChannel: sdk_const.SourceChannelFeishu,
		//OpenIds:       []string{openId},
		ChatIds: []string{chatId},
		CardMsg: cardMsg,
	})
	if errSys != nil {
		log.Errorf("[PushBotNotification] err:%v", errSys)
		return
	}
	log.Infof("[PushBotNotification] 发送机器人通知设置通知成功")

}

// PushBotNoticeForNotSupportChatInstruction 推送卡片，表示不支持该群聊指令
func PushBotNoticeForNotSupportChatInstruction(tenantKey string, event *callvo.MessageReqData,
	dealInfo *callvo.GroupChatHandleInfo) errs.SystemErrorInfo {
	//tenantClient := dealInfo.TenantClient

	orgInfoResp := orgfacade.GetBaseOrgInfoByOutOrgId(orgvo.GetBaseOrgInfoByOutOrgIdReqVo{OutOrgId: tenantKey})
	if orgInfoResp.Failure() {
		log.Errorf("[PushBotNoticeForNotSupportChatInstruction] GetBaseOrgInfoByOutOrgId failed:%v, tenantKey:%s",
			orgInfoResp.Error(), tenantKey)
		return orgInfoResp.Error()
	}
	// 通过操作人的 openId 获取其姓名
	opUserResp := orgfacade.GetBaseUserInfoByEmpId(orgvo.GetBaseUserInfoByEmpIdReqVo{
		OrgId: orgInfoResp.BaseOrgInfo.OrgId,
		EmpId: event.OpenId,
	})
	if opUserResp.Failure() {
		log.Errorf("[PushBotNoticeForNotSupportChatInstruction] orgfacade.GetBaseUserInfoByEmpId failed:%v, tenantKey:%s",
			opUserResp.Error(), tenantKey)
		return opUserResp.Error()
	}
	bindProText := fmt.Sprintf(" %d 个", len(dealInfo.BindProjectIds))
	atSomeoneUserStr := RenderAtSomeoneStr(event.OpenId, opUserResp.BaseUserInfo.Name)

	content := atSomeoneUserStr + "抱歉~关联" + bindProText + "项目时，暂无法进行快捷操作。\n" +
		"您是否需要重新进行推送设置？"
	cardMsg := card.GetFsCardNotSupportChatInstruction(dealInfo.OrgInfo.OrgId, event.OpenChatId, content)

	errSys := card.PushCard(dealInfo.OrgInfo.OrgId, &commonvo.PushCard{
		OrgId:         dealInfo.OrgInfo.OrgId,
		OutOrgId:      tenantKey,
		SourceChannel: orgInfoResp.BaseOrgInfo.SourceChannel,
		ChatIds:       []string{event.OpenChatId},
		CardMsg:       cardMsg,
	})
	if errSys != nil {
		log.Errorf("[HandleGroupChatAtUserNameWithIssueTitle] err:%v", errSys)
		return errSys
	}
	log.Infof("[PushBotNoticeForNotSupportChatInstruction] 发送机器人通知设置通知成功")

	return nil
}
