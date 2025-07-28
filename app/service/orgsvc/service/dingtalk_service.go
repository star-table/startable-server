package orgsvc

import (
	"fmt"
	"net/url"
	"strconv"
	"strings"
	"time"

	"gitea.bjx.cloud/allstar/polaris-backend/facade/projectfacade"
	"github.com/star-table/startable-server/common/core/util/json"

	"upper.io/db.v3"

	"github.com/star-table/startable-server/common/core/util/str"

	"github.com/spf13/cast"

	"github.com/star-table/startable-server/common/model/vo"

	"github.com/star-table/startable-server/common/model/vo/projectvo"

	"gitea.bjx.cloud/allstar/platform-sdk/sdk_interface"

	"gitea.bjx.cloud/allstar/dingtalk-sdk-golang/request"

	"gitea.bjx.cloud/allstar/dingtalk-sdk-golang"

	"github.com/google/martian/log"
	"github.com/google/uuid"

	platform_sdk "gitea.bjx.cloud/allstar/platform-sdk"
	sdk_const "gitea.bjx.cloud/allstar/platform-sdk/consts"
	vo2 "gitea.bjx.cloud/allstar/platform-sdk/vo"
	"github.com/star-table/startable-server/common/core/config"
	"github.com/star-table/startable-server/common/core/errs"
	"github.com/star-table/startable-server/common/core/util/encrypt"
	"github.com/star-table/startable-server/common/core/util/random"
	"github.com/star-table/startable-server/common/model/vo/orgvo"
)

func GetDingJsAPISign(input *orgvo.JsAPISignReq) (*orgvo.JsAPISignResp, errs.SystemErrorInfo) {
	resp := &orgvo.JsAPISignResp{}

	client, err := platform_sdk.GetClient(sdk_const.SourceChannelDingTalk, input.CorpID)
	if err != nil {
		return nil, errs.BuildSystemErrorInfo(errs.DingTalkOpenApiCallError, err)
	}

	outOrgInfo, err := domain.GetOrgOutInfoByTenantKey(input.CorpID)
	if err != nil {
		return nil, errs.BuildSystemErrorInfo(errs.DingTalkAvoidCodeInvalidError, err)
	}
	ticket, err := GetJsAPITicket(outOrgInfo.OrgId)
	if err != nil {
		return nil, errs.BuildSystemErrorInfo(errs.DingTalkAvoidCodeInvalidError, err)
	}
	resp.NoceStr = random.RandomString(5)
	resp.AgentID = config.GetDingTalkSdkConfig().AgentId

	// 如果是isv的情况下，每个corpId下面的agentId会不一样，需要拿一下
	orgInfo, err := client.GetOrgInfo(&vo2.OrgInfoReq{CorpId: input.CorpID})
	if err != nil {
		return nil, errs.BuildSystemErrorInfo(errs.DingTalkAvoidCodeInvalidError, err)
	}
	if orgInfo.AgentId != 0 {
		resp.AgentID = orgInfo.AgentId
	}

	url := input.URL
	timestamp := time.Now().UnixNano() / 1e6
	resp.TimeStamp = strconv.FormatInt(timestamp, 10)
	plain := "jsapi_ticket=" + ticket + "&noncestr=" + resp.NoceStr + "&timestamp=" + resp.TimeStamp + "&url=" + url
	resp.Signature = encrypt.SHA1(plain)
	//resp.Signature = sdk.CalculateJsApiSign(ticket, resp.NoceStr, timestamp, url)
	return resp, nil
}

func AuthDingCode(input orgvo.AuthDingCodeReqVo) (*orgvo.PlatformAuthCodeData, errs.SystemErrorInfo) {
	return thirdAuthCode(sdk_const.SourceChannelDingTalk, input.CorpId, input.Code, input.CodeType)
}

func GetSpaceList(input orgvo.GetSpaceListReq) (*orgvo.SpaceList, errs.SystemErrorInfo) {
	resp := &orgvo.SpaceList{}
	outInfo, err := domain.GetOrgOutInfo(input.OrgId)
	if err == nil && outInfo.SourceChannel == sdk_const.SourceChannelDingTalk {
		userOutInfo, err := domain.GetUserOutInfoByUserIdAndOrgId(input.UserId, input.OrgId, sdk_const.SourceChannelDingTalk)
		if err == nil {
			client, err := platform_sdk.GetClient(sdk_const.SourceChannelDingTalk, outInfo.OutOrgId)
			if err != nil {
				return nil, errs.BuildSystemErrorInfo(errs.DingTalkOpenApiCallError, err)
			}
			originClient := client.GetOriginClient().(*dingtalk.DingTalk)
			nextToken := ""
			for {
				spaceList, err := originClient.GetDriveSpaces(userOutInfo.OutOrgUserId, "org", nextToken, 50)
				if err != nil {
					return nil, errs.BuildSystemErrorInfo(errs.DingTalkOpenApiCallError, err)
				}

				for _, space := range spaceList.Spaces {
					resp.List = append(resp.List, &orgvo.SpaceInfo{
						SpaceId:   space.SpaceId,
						SpaceName: space.Name,
					})
				}

				if spaceList.Token == "" {
					break
				}
				nextToken = spaceList.Token
			}
		}
	}

	return resp, nil
}

func CreateCoolApp(input *orgvo.CreateCoolAppData) errs.SystemErrorInfo {
	client, apiErr := platform_sdk.GetClient(sdk_const.SourceChannelDingTalk, input.CorpId)
	if apiErr != nil {
		return errs.BuildSystemErrorInfo(errs.DingTalkOpenApiCallError, apiErr)
	}
	err := registerRefreshCardCallBack(client)
	if err != nil {
		log.Errorf("[CreateCoolApp] registerRefreshCardCallBack err:%v", err)
	}

	traceId := uuid.New().String()
	err = sendFirstCard(client.GetOriginClient().(*dingtalk.DingTalk), input.OpenConversationId, input.CorpId, traceId, input.RobotCode)
	if err != nil {
		log.Errorf("[CreateCoolApp] sendOrUpdateFirstCard input:%v,  err:%v", input, err)
	}

	return dao.GetDingCoolApp().CreateOrUpdate(&po.PpmOrgDingCoolApp{
		CorpId:         input.CorpId,
		ConversationId: input.OpenConversationId,
		RobotCode:      input.RobotCode,
		FirstTraceId:   traceId,
	})
}

func GetCoolAppInfo(conversationId string) (*orgvo.CoolAppInfo, errs.SystemErrorInfo) {
	coolApp, err := dao.GetDingCoolApp().Get(conversationId)
	if err != nil {
		if err == db.ErrNoMoreRows {
			return &orgvo.CoolAppInfo{}, nil
		}
		return nil, err
	}

	return &orgvo.CoolAppInfo{
		ProjectId: coolApp.ProjectId,
		AppId:     coolApp.AppId,
	}, nil
}

func DeleteCoolApp(input orgvo.DeleteCoolAppReq) errs.SystemErrorInfo {
	return dao.GetDingCoolApp().Delete(input.Input.OpenConversationId)
}

func DeleteCoolAppByProject(orgId, projectId int64) errs.SystemErrorInfo {
	coolApps, err := dao.GetDingCoolApp().GetByProjectId(orgId, projectId)
	if err != nil {
		if err == db.ErrNoMoreRows {
			return nil
		}
		return err
	}
	if len(coolApps) == 0 {
		return nil
	}

	client, sdkErr := platform_sdk.GetClient(sdk_const.SourceChannelDingTalk, coolApps[0].CorpId)
	if sdkErr != nil {
		log.Errorf("[sendTopCard] GetClient err:%v", sdkErr)
		return errs.BuildSystemErrorInfo(errs.DingTalkOpenApiCallError, sdkErr)
	}
	originClient := client.GetOriginClient().(*dingtalk.DingTalk)

	for _, app := range coolApps {
		err := updateFirstCard(originClient, app.ConversationId, app.CorpId, app.FirstTraceId, consts.BindTitle)
		if err != nil {
			log.Errorf("[sendTopCard] updateFirstCard, coolApp:%v, err:%v", app, err)
		}

		_, sdkErr = originClient.TopCardClose(&request.TopCardCloseReq{
			OutTrackId:         app.TopTraceId,
			ConversationType:   consts.ConversationTypeGroup,
			OpenConversationId: app.ConversationId,
			RobotCode:          app.RobotCode,
			CoolAppCode:        config.GetDingTalkSdkConfig().CoolAppCode,
		})
		if sdkErr != nil {
			log.Errorf("[sendTopCard] TopCardClose, coolApp:%v, err:%v", app, sdkErr)
		}
	}

	return dao.GetDingCoolApp().DeleteByProjectId(orgId, projectId)
}

func BindCoolApp(input orgvo.BindCoolAppReq) errs.SystemErrorInfo {
	topTraceId := uuid.New().String()
	input.Input.OpenConversationId = strings.ReplaceAll(input.Input.OpenConversationId, " ", "+")
	coolApp, err := dao.GetDingCoolApp().Get(input.Input.OpenConversationId)
	if err != nil {
		log.Errorf("[BindCoolAppReq] Get input:%v, err:%v", input, err)
		return err
	}
	if coolApp.TopTraceId != "" {
		topTraceId = ""
	}

	updateModel := &po.PpmOrgDingCoolApp{
		OrgId:          input.OrgId,
		AppId:          input.Input.AppId,
		ProjectId:      input.Input.ProjectId,
		ConversationId: input.Input.OpenConversationId,
		TopTraceId:     topTraceId,
	}

	err = dao.GetDingCoolApp().CreateOrUpdate(updateModel)
	if err != nil {
		return err
	}

	client, sdkErr := platform_sdk.GetClient(sdk_const.SourceChannelDingTalk, coolApp.CorpId)
	if sdkErr != nil {
		return errs.BuildSystemErrorInfo(errs.DingTalkOpenApiCallError, sdkErr)
	}

	if topTraceId != "" {
		err = sendTopCard(client.GetOriginClient().(*dingtalk.DingTalk), input.Input.OpenConversationId, coolApp.CorpId, coolApp.RobotCode, topTraceId, input.OrgId, input.Input.ProjectId)
		if err != nil {
			log.Errorf("[BindCoolAppReq] sendTopCard input:%v, err:%v", input, err)
		}
	}

	err = updateFirstCard(client.GetOriginClient().(*dingtalk.DingTalk), input.Input.OpenConversationId, coolApp.CorpId, coolApp.FirstTraceId, consts.CreateIssueTitle)
	if err != nil {
		log.Errorf("[CreateCoolApp] sendOrUpdateFirstCard coolApp:%v,  err:%v", coolApp, err)
	}

	return nil
}

func UpdateTopCard(orgId, projectId int64) errs.SystemErrorInfo {
	coolApps, err := dao.GetDingCoolApp().GetByProjectId(orgId, projectId)
	if err != nil {
		if err == db.ErrNoMoreRows {
			return nil
		}
		return err
	}

	if len(coolApps) == 0 {
		return nil
	}

	cardData, err := getTopCardData(orgId, projectId, coolApps[0].CorpId, "")
	if err != nil {
		log.Errorf("[UpdateTopCard] getTopCardData err:%v", err)
		return err
	}

	client, sdkErr := platform_sdk.GetClient(sdk_const.SourceChannelDingTalk, coolApps[0].CorpId)
	if sdkErr != nil {
		log.Errorf("[sendTopCard] GetClient err:%v", sdkErr)
		return errs.BuildSystemErrorInfo(errs.DingTalkOpenApiCallError, sdkErr)
	}
	for _, app := range coolApps {
		originClient := client.GetOriginClient().(*dingtalk.DingTalk)
		cardData.CardParamMap["openConversationId"] = app.ConversationId
		_, sdkErr = originClient.CardUpdate(&request.CardUpdateReq{
			OutTrackId: app.TopTraceId,
			CardData:   cardData,
		})
		if sdkErr != nil {
			log.Errorf("[UpdateTopCard] CardUpdate err:%v", sdkErr)
		}
	}

	return nil
}

func GetTopCardData(conversationId string) (string, errs.SystemErrorInfo) {
	coolApp, err := dao.GetDingCoolApp().Get(conversationId)
	if err != nil {
		return "", err
	}

	cardData, err := getTopCardData(coolApp.OrgId, coolApp.ProjectId, coolApp.CorpId, conversationId)
	if err != nil {
		log.Errorf("[UpdateTopCard] getTopCardData err:%v", err)
		return "", err
	}

	return json.ToJsonIgnoreError(map[string]interface{}{"cardData": cardData}), nil
}

func thirdAuthCode(sourceChannel, corpId, code string, codeType int) (*orgvo.PlatformAuthCodeData, errs.SystemErrorInfo) {
	authCodeInfo, err := getPlatformUserInfoByAuthCode(sourceChannel, corpId, code, codeType)
	if err != nil {
		log.Error(err)
		return nil, err
	}

	if authCodeInfo.Binding {
		err := platformLoginByAuthCodeInfo(sourceChannel, authCodeInfo)
		if err != nil {
			return nil, err
		}
	}

	res := &orgvo.PlatformAuthCodeData{
		CorpId:       authCodeInfo.CorpId,
		OutUserId:    authCodeInfo.OutUserId,
		IsAdmin:      authCodeInfo.IsAdmin,
		Binding:      authCodeInfo.Binding,
		RefreshToken: authCodeInfo.RefreshToken,
		AccessToken:  authCodeInfo.AccessToken,
		CodeToken:    random.Token(),
		OrgID:        authCodeInfo.OrgID,
		OrgName:      authCodeInfo.OrgName,
		OutOrgName:   authCodeInfo.OutOrgName,
		OrgCode:      authCodeInfo.OrgCode,
		UserID:       authCodeInfo.UserID,
		Name:         authCodeInfo.Name,
		Token:        authCodeInfo.Token,
	}

	return res, nil
}

// registerRefreshCardCallBack 注册刷新回调
func registerRefreshCardCallBack(client sdk_interface.Sdk) errs.SystemErrorInfo {
	topCardRefreshUrl := config.GetDingTalkSdkConfig().CallBackUrl + consts.RefreshUrlPath
	originClient := client.GetOriginClient().(*dingtalk.DingTalk)

	_, apiErr := originClient.CardRegisterCallback(&request.CardRegisterCallbackReq{
		CallbackRouteKey: consts.CallBackKeyRefresh,
		CallbackUrl:      topCardRefreshUrl,
		ForceUpdate:      true,
	})
	if apiErr != nil {
		log.Errorf("[sendFirstCard] CardCreateAndSend err:%v", apiErr)
		return errs.BuildSystemErrorInfo(errs.DingTalkOpenApiCallError, apiErr)
	}

	return nil
}

// 发送或者更新第一个卡片，用于绑定项目
func sendFirstCard(originClient *dingtalk.DingTalk, conversationId, corpId, traceId, robotCode string) errs.SystemErrorInfo {
	cardData := getFirstCardData(conversationId, corpId, consts.BindTitle, consts.BindProjectPath)
	_, apiErr := originClient.CardCreateAndSend(&request.CardCreateAndSendReq{
		OpenSpaceId:             "dtv1.card//im_group." + conversationId,
		CardTemplateId:          consts.FirstCardTemplateId,
		OutTrackId:              traceId,
		CardData:                cardData,
		ImGroupOpenDeliverModel: &request.ImGroupOpenDeliverModel{RobotCode: robotCode},
		ImGroupOpenSpaceModel:   &request.ImGroupOpenSpaceModel{SupportForward: false, LastMessageI18n: map[string]string{"ZH_CN": "卡片"}},
	})
	if apiErr != nil {
		log.Errorf("[sendFirstCard] CardCreateAndSend err:%v", apiErr)
		return errs.BuildSystemErrorInfo(errs.DingTalkOpenApiCallError, apiErr)
	}

	return nil
}

func updateFirstCard(originClient *dingtalk.DingTalk, conversationId, corpId, traceId, title string) errs.SystemErrorInfo {
	cardData := getFirstCardData(conversationId, corpId, title, consts.CreateIssuePath)
	_, apiErr := originClient.CardUpdate(&request.CardUpdateReq{
		OutTrackId: traceId,
		CardData:   cardData,
	})
	if apiErr != nil {
		log.Errorf("[sendFirstCard] CardCreateAndSend err:%v", apiErr)
		return errs.BuildSystemErrorInfo(errs.DingTalkOpenApiCallError, apiErr)
	}

	return nil
}

func getFirstCardData(conversationId, corpId, title, bindPath string) *request.GroupSendInteractiveCardReqCardData {
	bindUrl := url.QueryEscape(fmt.Sprintf(bindPath, conversationId, corpId))
	loginUrl := url.QueryEscape(fmt.Sprintf(config.GetDingTalkSdkConfig().FrontUrl+consts.LoginUrl, bindUrl))
	bindPcUrl := url.QueryEscape(fmt.Sprintf(consts.DingPcSide, loginUrl))
	bindMobileUrl := url.QueryEscape(fmt.Sprintf(consts.DingMobileSide, loginUrl))
	bindPlatformUrl := fmt.Sprintf(consts.DingPlatform, bindPcUrl, bindMobileUrl)

	helpUrl := url.QueryEscape(consts.HelpPath)
	helpPcUrl := url.QueryEscape(fmt.Sprintf(consts.DingPcSide, helpUrl))
	helpMobileUrl := url.QueryEscape(fmt.Sprintf(consts.DingMobileSide, helpUrl))
	helpPlatformUrl := fmt.Sprintf(consts.DingPlatform, helpPcUrl, helpMobileUrl)

	return &request.GroupSendInteractiveCardReqCardData{CardParamMap: map[string]string{
		"openConversationId": conversationId,
		"bindTitle":          title,
		"bindUrl":            bindPlatformUrl,
		"helpUrl":            helpPlatformUrl,
	}}
}

// sendTopCard 发送吊顶卡片
func sendTopCard(originClient *dingtalk.DingTalk, conversationId, corpId, robotCode, traceId string, orgId, projectId int64) errs.SystemErrorInfo {
	cardData, err := getTopCardData(orgId, projectId, corpId, conversationId)
	if err != nil {
		return err
	}

	_, apiErr := originClient.TopCardSend(&request.TopCardSendReq{
		ConversationType:   consts.ConversationTypeGroup,
		OpenConversationId: conversationId,
		CardTemplateId:     consts.TopCardTemplateId,
		OutTrackId:         traceId,
		RobotCode:          robotCode,
		CoolAppCode:        config.GetDingTalkSdkConfig().CoolAppCode,
		CardData:           cardData,
		CallbackRouteKey:   consts.CallBackKeyRefresh,
	})
	if apiErr != nil {
		log.Errorf("[sendTopCard] TopCardSend err:%v", apiErr)
		return errs.BuildSystemErrorInfo(errs.DingTalkOpenApiCallError, apiErr)
	}

	return nil
}

func getTopCardData(orgId, projectId int64, corpId, conversationId string) (*request.GroupSendInteractiveCardReqCardData, errs.SystemErrorInfo) {
	statResp := projectfacade.IssueStatusTypeStat(projectvo.IssueStatusTypeStatReqVo{
		Input: &vo.IssueStatusTypeStatReq{
			ProjectID: &projectId,
		},
		OrgId: orgId,
	})
	if statResp.Failure() {
		log.Errorf("[sendTopCard] IssueStatusTypeStat projectId:%v, err:%v", projectId, statResp.Error())
		return nil, statResp.Error()
	}

	projectInfoResp := projectfacade.ProjectInfo(projectvo.ProjectInfoReqVo{
		Input: vo.ProjectInfoReq{ProjectID: projectId},
		OrgId: orgId,
	})
	if projectInfoResp.Failure() {
		log.Errorf("[sendTopCard] ProjectInfo projectId:%v, err:%v", projectId, statResp.Error())
		return nil, statResp.Error()
	}
	projectUrl := url.QueryEscape(fmt.Sprintf(config.GetDingTalkSdkConfig().FrontUrl+consts.ProjectDetailPath, projectId, projectInfoResp.ProjectInfo.ProjectTypeID, projectInfoResp.ProjectInfo.AppID))
	percent := 0
	if statResp.IssueStatusTypeStat.Total > 0 {
		percent = int(float64(100*statResp.IssueStatusTypeStat.CompletedTotal) / float64(statResp.IssueStatusTypeStat.Total))
	}
	return &request.GroupSendInteractiveCardReqCardData{CardParamMap: map[string]string{
		"openConversationId": conversationId,
		"projectUrl":         fmt.Sprintf(consts.DingConsole, corpId, config.GetDingTalkSdkConfig().AppId, projectUrl),
		"projectName":        str.TruncateName(projectInfoResp.ProjectInfo.Name, 4),
		"overdueToday":       cast.ToString(statResp.IssueStatusTypeStat.OverdueTodayTotal),
		"processing":         cast.ToString(statResp.IssueStatusTypeStat.ProcessingTotal + statResp.IssueStatusTypeStat.NotStartTotal),
		"overdue":            cast.ToString(statResp.IssueStatusTypeStat.OverdueTotal),
		"completed":          cast.ToString(statResp.IssueStatusTypeStat.CompletedTotal),
		"percent":            cast.ToString(percent) + "%",
		"updateTime":         time.Now().Format("01/02 15:04"),
	}}, nil
}
