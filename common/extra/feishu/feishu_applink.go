package feishu

import (
	"encoding/base64"
	"fmt"
	url2 "net/url"
	"strings"

	"github.com/star-table/startable-server/common/core/config"
	"github.com/star-table/startable-server/common/core/util/temp"
)

const (
	AppLinkParamNameAppId               = "AppId"
	AppLinkParamNameProAppId            = "ProAppId"
	AppLinkParamNameIssueId             = "IssueId"
	AppLinkParamNameParentId            = "ParentId"
	AppLinkParamNameProjectId           = "ProjectId"
	AppLinkParamNameChatId              = "ChatId"
	AppLinkParamNameProjectTypeId       = "ProjectTypeId"
	AppLinkParamNameProjectObjectTypeId = "TableId"
	AppLinkParamNameProjectName         = "ProjectName"
	AppLinkParamNameOrgId               = "OrgId"

	// 小程序打开 mode
	AppModeSidebar   = "sidebar-semi"
	AppModeAppCenter = "appCenter"

	AppLinkDefaultMini = "https://applink.feishu.cn/client/mini_program/open?appId={{.AppId}}&mode=sidebar-semi"
	AppLinkDefault     = "https://applink.feishu.cn/client/web_app/open?appId={{.AppId}}&path=feishu%2Fauth&mode=sidebar-semi"
	//任务详情 小程序 sidebar
	AppLinkIssueInfoMini = "https://applink.feishu.cn/client/mini_program/open?appId={{.AppId}}&mode=sidebar-semi&path=components/widget/AddTask/index%3fid%3d{{.IssueId}}%26projectId%3d{{.ProjectId}}"
	AppLinkIssueInfo     = "https://applink.feishu.cn/client/web_app/open?appId={{.AppId}}&mode=sidebar-semi&path=/taskDetail%3FissueId%3D{{.IssueId}}%26platform%3DFEISHU%26fsmode%3Dwindow%26isRedirect%3Dtrue%26orgId%3D{{.OrgId}}"
	// 任务详情 小程序 appCenter
	AppLinkIssueInfoForAppCenter = "https://applink.feishu.cn/client/mini_program/open?appId={{.AppId}}&mode=appCenter&path=pages%2fPC%2fAddTask%2findex%3fid%3d{{.IssueId}}%26parentId%3d{{.ParentId}}%26projectId%3d{{.ProjectId}}"

	AppLinkMobileOpenQRCode        = "https://applink.feishu.cn/client/mini_program/open?appId={{.AppId}}&mode=sidebar-semi&path=pages/PC/Home/index"
	AppLinkDefaultProjectConfigure = "https://applink.feishu.cn/client/mini_program/open?appId={{.AppId}}&mode=sidebar-semi&path=pages/PC/Configure/index"
	//通知设置中心
	AppLinkDefaultNotification = "https://applink.feishu.cn/client/mini_program/open?appId={{.AppId}}&mode=sidebar-semi&path=pages/Notification/index"

	//项目统计
	AppLinkProjectStatistic = "https://applink.feishu.cn/client/mini_program/open?appId={{.AppId}}&mode=sidebar-semi&path=pages%2fPC%2fProjectStatistics%2findex%3fprojectId%3d{{.ProjectId}}"
	//小程序pc地址
	AppLinkPcWelcomeUrlMini = "https://applink.feishu.cn/client/mini_program/open?appId={{.AppId}}&mode=appCenter"
	AppLinkPcWelcomeUrl     = "https://applink.feishu.cn/client/web_app/open?appId={{.AppId}}&path=feishu/auth"
	// 项目页地址 项目详情，任务列表
	AppLinkProjectUrl = "https://applink.feishu.cn/client/mini_program/open?appId={{.AppId}}&mode=appCenter&path=pages%2fPC%2fIndex%2findex%3fredirectUrl%3d%2fproject%2f{{.ProjectId}}%2ftask%2f{{.TableId}}%3fprojectTypeId%3D{{.ProjectTypeId}}"
	// 项目概览地址
	AppLinkProjectOverviewUrl = "https://applink.feishu.cn/client/mini_program/open?appId={{.AppId}}&mode=appCenter&path=pages%2fPC%2fIndex%2findex%3fredirectUrl%3d%2fproject%2f{{.ProjectId}}%2foverview"
	//项目详情页 手机端
	AppLinkMobileProjectUrl = "https://applink.feishu.cn/client/mini_program/open?appId={{.AppId}}&path=pages%2fProjectDetail%2findex%3fprojectId%3d{{.ProjectId}}"
	// 任务模块页
	AppLinkMobileTaskModuleUrl = "https://applink.feishu.cn/client/mini_program/open?appId={{.AppId}}&path="

	AppGuide           = "https://polaris.feishu.cn/docs/doccn5PXJnUDDhLee2Uksn337mb"
	AppCustomerService = "https://applink.feishu.cn/client/helpdesk/open?id=7031382079124373507&extra=%7B%22channel%22%3A1%2C%22created_at%22%3A1637122101%7D"
	AppConfigDoc       = "https://polaris.feishu.cn/docs/doccn56SrnQKluE4VnLnL4MTMKe"
	AppGuideToEntrance = "https://polaris.feishu.cn/docs/doccn5PXJnUDDhLee2Uksn337mb#KA1MPN"

	//网页版app
	AppLinkWebDefault = "https://applink.feishu.cn/client/web_app/open?appId={{.AppId}}&mode=appCenter"
	// 飞书群聊：项目动态设置
	AppLinkGroupChatProjectTrendPushConfig = "https://applink.feishu.cn/client/mini_program/open?appId={{." +
		"AppId}}&mode=sidebar-semi&path=pages%2fPC%2fChatMessageSetting%2findex%3fprojectId%3d{{." +
		"ProjectId}}%26chatId%3d{{.ChatId}}"
	// 飞书管理后台-配置极星应用
	FsLinkForSettingPolarisApp = "https://feishu.cn/admin/appCenter/manage/%s"
)

func IsGrayOrg(orgId int64) bool {
	// 使用新版的 appLink 所以这里直接返回 true
	return true
	// key := consts.CacheGrayFeishuOrg
	// value, err := cache.Get(key)
	// if err != nil {
	// 	log.Error(err)
	// 	return false
	// }
	//
	// var grayOrgList []int64
	// if value != "" {
	// 	_ = json.FromJson(value, &grayOrgList)
	// } else {
	// 	return true
	// }
	//
	// ok, _ := slice.Contain(grayOrgList, orgId)
	//
	// return ok
}

func GetDefaultAppLink(orgId int64) string {
	fsConfig := config.GetConfig().FeiShu
	if fsConfig == nil {
		log.Error("飞书配置为空")
		return ""
	}
	params := map[string]interface{}{
		AppLinkParamNameAppId: fsConfig.AppId,
	}
	template := AppLinkDefaultMini
	if IsGrayOrg(orgId) {
		template = AppLinkDefault
	}
	link, err := temp.Render(template, params)
	if err != nil {
		log.Error(err)
		return ""
	}
	return link
}

func GetMobileOpenQRCodeAppLink(orgId int64) string {
	fsConfig := config.GetConfig().FeiShu
	if fsConfig == nil {
		log.Error("飞书配置为空")
		return ""
	}
	params := map[string]interface{}{
		AppLinkParamNameAppId: fsConfig.AppId,
	}
	template := AppLinkMobileOpenQRCode
	if IsGrayOrg(orgId) {
		template = AppLinkDefault
	}
	link, err := temp.Render(template, params)
	if err != nil {
		log.Error(err)
		return ""
	}
	return link
}

func GetDefaultProjectConfigureAppLink() string {
	fsConfig := config.GetConfig().FeiShu
	if fsConfig == nil {
		log.Error("飞书配置为空")
		return ""
	}
	params := map[string]interface{}{
		AppLinkParamNameAppId: fsConfig.AppId,
	}
	link, err := temp.Render(AppLinkDefaultProjectConfigure, params)
	if err != nil {
		log.Error(err)
		return ""
	}
	return link
}

//func GetIssueSideBarLinkFeiShu(projectId int64, issueId int64) string {
//	fsConfig := config.GetConfig().FeiShu
//	if fsConfig == nil {
//		log.Error("飞书配置为空")
//		return ""
//	}
//	// 旧的条件：IsGrayOrg(orgId)，做个备份
//	if true {
//		redirectUri := fmt.Sprintf("/taskDetail?issueId=%d", issueId)
//		encodeUri := EncodeUriForFeiShu(redirectUri)
//		return fmt.Sprintf("https://applink.feishu.cn/client/web_app/open?appId=%s&mode=sidebar-semi&path=feishu/auth&callback=%s", fsConfig.AppId, encodeUri)
//	} else {
//		params := map[string]interface{}{
//			AppLinkParamNameAppId:     fsConfig.AppId,
//			AppLinkParamNameIssueId:   issueId,
//			AppLinkParamNameProjectId: projectId,
//		}
//
//		link, err := temp.Render(AppLinkIssueInfoMini, params)
//		if err != nil {
//			log.Error(err)
//			return ""
//		}
//		return link
//	}
//}

// FsGetIssueInfoAppLinkForDeleteIssue 回收站查看任务详情
// @param openType `sidebar-semi` 侧边；`appCenter`工作台打开；
// openType 更多枚举值参考 https://open.feishu.cn/document/uAjLw4CM/uYjL24iN/applink-protocol/supported-protocol/open-a-gadget
func FsGetIssueInfoAppLinkForDeleteIssue(issueId int64, proAppId int64, openType string) string {
	fsConfig := config.GetConfig().FeiShu
	if fsConfig == nil {
		log.Error("飞书配置为空")
		return ""
	}
	redirectUri := fmt.Sprintf("/trash/app%d?issueId=%d", proAppId, issueId)
	encodeUri := EncodeUriForFeiShu(redirectUri)

	return fmt.Sprintf("https://applink.feishu.cn/client/web_app/open?appId=%s&mode=%s&path=feishu/auth&callback=%s",
		fsConfig.AppId, openType, encodeUri)
}

func EncodeUriForFeiShu(uri string) string {
	encodeUri := base64.StdEncoding.EncodeToString([]byte(uri))
	//替换特殊字符
	encodeUri = strings.ReplaceAll(encodeUri, "+", "-")
	encodeUri = strings.ReplaceAll(encodeUri, "/", "_")
	encodeUri = strings.ReplaceAll(encodeUri, "=", "$")

	return encodeUri
}

// GetTaskModuleUrl 获取任务模块的地址
// isPc 为 true 表示从工作台打开
func GetTaskModuleUrl(isPc bool) string {
	fsConfig := config.GetConfig().FeiShu
	if fsConfig == nil {
		log.Error("飞书配置为空")
		return ""
	}

	redirectUri := fmt.Sprintf("/task")
	encodeUri := EncodeUriForFeiShu(redirectUri)
	if isPc {
		// pc 端，从工作台打开
		// https://open.feishu.cn/document/uAjLw4CM/uYjL24iN/applink-protocol/supported-protocol/open-an-h5-app
		link := fmt.Sprintf("https://applink.feishu.cn/client/web_app/open?appId=%s&path=feishu/auth?callback=%s", fsConfig.AppId, encodeUri)

		return link
	} else {
		return fmt.Sprintf("https://applink.feishu.cn/client/web_app/open?appId=%s&mode=sidebar-semi&path=feishu/auth&callback=%s", fsConfig.AppId, encodeUri)
	}
}

func GetProjectWelcomeUrl() string {
	fsConfig := config.GetConfig().FeiShu
	if fsConfig == nil {
		log.Error("飞书配置为空")
		return ""
	}

	params := map[string]interface{}{
		AppLinkParamNameAppId: fsConfig.AppId,
	}

	link, err := temp.Render(fsConfig.CardJumpLink.ProjectWelcomeUrl, params)
	if err != nil {
		log.Error(err)
		return ""
	}
	log.Infof("飞书pc首页跳转 %s", link)

	begin := strings.Index(link, "http%3a%2f%2f")
	url, _ := url2.QueryUnescape(link[begin:])

	return url
}

func GetAppLinkPcWelcomeUrl(orgId int64) string {
	fsConfig := config.GetConfig().FeiShu
	if fsConfig == nil {
		log.Error("飞书配置为空")
		return ""
	}
	params := map[string]interface{}{
		AppLinkParamNameAppId: fsConfig.AppId,
	}
	template := AppLinkPcWelcomeUrlMini
	if IsGrayOrg(orgId) {
		template = AppLinkPcWelcomeUrl
	}
	link, err := temp.Render(template, params)
	if err != nil {
		log.Error(err)
		return ""
	}
	return link
}

func GetAppLinkWebPcWelcomeUrl() string {
	fsConfig := config.GetConfig().FeiShu
	if fsConfig == nil {
		log.Error("飞书配置为空")
		return ""
	}
	params := map[string]interface{}{
		AppLinkParamNameAppId: fsConfig.AppId,
	}
	link, err := temp.Render(AppLinkWebDefault, params)
	if err != nil {
		log.Error(err)
		return ""
	}
	return link
}

// 飞书群聊：项目动态设置配置
func GetGroupChatProjectTrendPushConfigAppLink(projectId int64, orgId int64, chatId string) string {
	fsConfig := config.GetConfig().FeiShu
	if fsConfig == nil {
		log.Error("飞书配置为空")
		return ""
	}
	if IsGrayOrg(orgId) {
		encodeUri := base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("/projectPushSetting?projectId=%d&chatId=%s",
			projectId, chatId)))
		//替换特殊字符
		encodeUri = strings.ReplaceAll(encodeUri, "+", "-")
		encodeUri = strings.ReplaceAll(encodeUri, "/", "_")
		encodeUri = strings.ReplaceAll(encodeUri, "=", "$")

		return fmt.Sprintf("https://applink.feishu.cn/client/web_app/open?appId=%s&path=feishu/auth&mode=sidebar-semi&callback=%s", fsConfig.AppId, encodeUri)
	} else {
		params := map[string]interface{}{
			AppLinkParamNameAppId:     fsConfig.AppId,
			AppLinkParamNameProjectId: projectId,
			AppLinkParamNameChatId:    chatId,
		}

		link, err := temp.Render(AppLinkGroupChatProjectTrendPushConfig, params)
		if err != nil {
			log.Error(err)
			return ""
		}
		return link
	}
}

// GetGroupChatSettingBtnAppLink 飞书群聊：项目动态设置配置
func GetGroupChatSettingBtnAppLink(orgId int64, chatId string) string {
	fsConfig := config.GetConfig().FeiShu
	if fsConfig == nil {
		log.Errorf("[GetGroupChatSettingBtnAppLink] 飞书配置为空 orgId: %d", orgId)
		return ""
	}
	if IsGrayOrg(orgId) {
		encodeUri := base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("/projectPushSetting?projectId=%d&chatId=%s",
			0, chatId)))
		//替换特殊字符
		encodeUri = strings.ReplaceAll(encodeUri, "+", "-")
		encodeUri = strings.ReplaceAll(encodeUri, "/", "_")
		encodeUri = strings.ReplaceAll(encodeUri, "=", "$")

		return fmt.Sprintf("https://applink.feishu."+
			"cn/client/web_app/open?appId=%s&path=/feishu/auth&mode=sidebar-semi&callback=%s", fsConfig.AppId, encodeUri)
	} else {
		params := map[string]interface{}{
			AppLinkParamNameAppId:     fsConfig.AppId,
			AppLinkParamNameProjectId: 0, // 无用参数
			AppLinkParamNameChatId:    chatId,
		}

		link, err := temp.Render(AppLinkGroupChatProjectTrendPushConfig, params)
		if err != nil {
			log.Errorf("[GetGroupChatSettingBtnAppLink] err: %v", err)
			return ""
		}
		return link
	}
}

func GetAppLinkProjectUrl(projectId, projectTypeId, projectObjectTypeId int64) string {
	params := map[string]interface{}{
		AppLinkParamNameProjectId:           projectId,
		AppLinkParamNameProjectTypeId:       projectTypeId,
		AppLinkParamNameProjectObjectTypeId: projectObjectTypeId,
	}
	return RenderOneAppLink(AppLinkProjectUrl, params)
}

// 通过链接模板+参数，渲染出链接
func RenderOneAppLink(linkTmp string, params map[string]interface{}) string {
	fsConfig := config.GetConfig().FeiShu
	if fsConfig == nil {
		log.Error("飞书配置为空")
		return ""
	}
	// 追加 appId 参数
	params[AppLinkParamNameAppId] = fsConfig.AppId
	link, err := temp.Render(linkTmp, params)
	if err != nil {
		log.Error(err)
		return ""
	}
	return link
}

// 项目概览链接
func GetAppLinkProjectOverviewUrl(projectId int64) string {
	params := map[string]interface{}{
		AppLinkParamNameProjectId: projectId,
	}
	return RenderOneAppLink(AppLinkProjectOverviewUrl, params)
}

func GetMobileProjectUrl(projectId int64) string {
	params := map[string]interface{}{
		AppLinkParamNameProjectId: projectId,
	}
	return RenderOneAppLink(AppLinkMobileProjectUrl, params)
}

// 通知设置页面
func GetDefaultNotificationAppLink(orgId int64) string {
	fsConfig := config.GetConfig().FeiShu
	if fsConfig == nil {
		log.Error("飞书配置为空")
		return ""
	}

	if IsGrayOrg(orgId) {
		encodeUri := base64.StdEncoding.EncodeToString([]byte("/notificationSettings"))
		//替换特殊字符
		encodeUri = strings.ReplaceAll(encodeUri, "+", "-")
		encodeUri = strings.ReplaceAll(encodeUri, "/", "_")
		encodeUri = strings.ReplaceAll(encodeUri, "=", "$")

		return fmt.Sprintf("https://applink.feishu.cn/client/web_app/open?appId=%s&path=/feishu/auth&mode=sidebar-semi&callback=%s",
			fsConfig.AppId, encodeUri)
	} else {
		params := map[string]interface{}{
			AppLinkParamNameAppId: fsConfig.AppId,
		}
		link, err := temp.Render(AppLinkDefaultNotification, params)
		if err != nil {
			log.Error(err)
			return ""
		}
		return link
	}
}
