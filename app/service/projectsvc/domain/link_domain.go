package domain

import (
	platform_sdk "gitea.bjx.cloud/allstar/platform-sdk"
	sdk_const "gitea.bjx.cloud/allstar/platform-sdk/consts"
	sdk_vo "gitea.bjx.cloud/allstar/platform-sdk/vo"
	"github.com/star-table/startable-server/app/facade/orgfacade"
	"github.com/star-table/startable-server/common/core/config"
	"github.com/star-table/startable-server/common/core/util/slice"
	"github.com/star-table/startable-server/common/model/vo/orgvo"
	"github.com/star-table/startable-server/common/model/vo/projectvo"
)

// todo 上下文传递项目信息 表信息 记录信息
func GetIssueLinks(sourceChannel string, orgId int64, issueId int64) *projectvo.IssueLinks {
	issueLink := &projectvo.IssueLinks{}

	// 支持不传sourceChannel，通过orgId去拿
	orgInfoResp := orgfacade.GetBaseOrgInfo(orgvo.GetBaseOrgInfoReqVo{OrgId: orgId})
	if orgInfoResp.Failure() {
		log.Errorf("[GetIssueLinks] GetBaseOrgInfo failed: %v", orgInfoResp.Error())
		return issueLink
	}
	if sourceChannel == "" {
		sourceChannel = orgInfoResp.BaseOrgInfo.SourceChannel
	}

	if ok, _ := slice.Contain([]string{sdk_const.SourceChannelFeishu, sdk_const.SourceChannelWeixin,
		sdk_const.SourceChannelDingTalk}, sourceChannel); !ok {
		return issueLink
	}

	issueInfo, issueErr := GetIssueInfoLc(orgId, 0, issueId)
	if issueErr != nil {
		log.Errorf("[GetIssueLinks] GetIssueBo failed: %v", issueErr)
		return issueLink
	}

	appId, appIdErr := GetAppIdFromProjectId(orgId, issueInfo.ProjectId)
	if appIdErr != nil {
		log.Errorf("[GetIssueLinks] GetAppIdFromProjectId failed: %v", appIdErr)
		return issueLink
	}
	var projectTypeId int64
	if issueInfo.ProjectId > 0 {
		projectInfo, projectInfoErr := GetProjectSimple(orgId, issueInfo.ProjectId)
		if projectInfoErr != nil {
			log.Errorf("[GetIssueLinks] GetProject failed: %v", projectInfoErr)
			return issueLink
		}
		projectTypeId = projectInfo.ProjectTypeId
	}

	outOrgId := orgInfoResp.BaseOrgInfo.OutOrgId

	sdk, sdkErr := platform_sdk.GetClient(sourceChannel, "")
	if sdkErr != nil {
		log.Errorf("[GetIssueLinks] platform_sdk.GetClient failed: %v", sdkErr.Error())
		return issueLink
	}
	host := ""
	serverCommon := config.GetConfig().ServerCommon
	if serverCommon != nil {
		if sourceChannel == sdk_const.SourceChannelWeixin {
			host = serverCommon.WeiXinHost
		} else {
			host = serverCommon.Host
		}
	}
	// 拿link
	link, sdkErr := sdk.GetIssueLink(&sdk_vo.GetIssueLinkReq{
		AppId:         appId,
		ProjectId:     issueInfo.ProjectId,
		ProjectTypeId: projectTypeId,
		TableId:       issueInfo.TableId,
		IssueId:       issueId,
		Host:          host,
		CorpId:        outOrgId,
	})
	if sdkErr != nil {
		log.Errorf("[GetIssueLinks] sdk.GetIssueLink failed: %v", sdkErr.Error())
		return issueLink
	}
	issueLink.Link = link.Url

	// 拿sidebar link
	sideBarLink, sdkErr := sdk.GetIssueSideBarLink(&sdk_vo.GetIssueSideBarLinkReq{
		IssueId: issueId,
		Host:    host,
		CorpId:  outOrgId,
	})
	if sdkErr != nil {
		log.Errorf("[GetIssueLinks] sdk.GetIssueSideBarLink failed: %v", sdkErr.Error())
		return issueLink
	}
	issueLink.SideBarLink = sideBarLink.Url

	return issueLink
}

func GetTodoLinks(sourceChannel string, orgId int64, todoId int64) *projectvo.TodoLinks {
	todoLink := &projectvo.TodoLinks{}

	// 支持不传sourceChannel，通过orgId去拿
	orgInfoResp := orgfacade.GetBaseOrgInfo(orgvo.GetBaseOrgInfoReqVo{OrgId: orgId})
	if orgInfoResp.Failure() {
		log.Errorf("[GetTodoLinks] GetBaseOrgInfo failed: %v", orgInfoResp.Error())
		return todoLink
	}
	if sourceChannel == "" {
		sourceChannel = orgInfoResp.BaseOrgInfo.SourceChannel
	}

	if ok, _ := slice.Contain([]string{sdk_const.SourceChannelFeishu, sdk_const.SourceChannelWeixin,
		sdk_const.SourceChannelDingTalk}, sourceChannel); !ok {
		return todoLink
	}

	outOrgId := orgInfoResp.BaseOrgInfo.OutOrgId

	sdk, sdkErr := platform_sdk.GetClient(sourceChannel, "")
	if sdkErr != nil {
		log.Errorf("[GetTodoLinks] platform_sdk.GetClient failed: %v", sdkErr.Error())
		return todoLink
	}
	host := ""
	serverCommon := config.GetConfig().ServerCommon
	if serverCommon != nil {
		if sourceChannel == sdk_const.SourceChannelWeixin {
			host = serverCommon.WeiXinHost
		} else {
			host = serverCommon.Host
		}
	}

	// 拿link
	link, sdkErr := sdk.GetTodoLink(&sdk_vo.GetTodoLinkReq{
		TodoId: todoId,
		Host:   host,
		CorpId: outOrgId,
	})
	if sdkErr != nil {
		log.Errorf("[GetTodoLinks] sdk.GetTodoLink failed: %v", sdkErr.Error())
		return todoLink
	}
	todoLink.Link = link.Url

	// 拿sidebar link
	sideBarLink, sdkErr := sdk.GetTodoSideBarLink(&sdk_vo.GetTodoLinkReq{
		TodoId: todoId,
		Host:   host,
		CorpId: outOrgId,
	})
	if sdkErr != nil {
		log.Errorf("[GetTodoLinks] sdk.GetTodoSideBarLink failed: %v", sdkErr.Error())
		return todoLink
	}
	todoLink.SideBarLink = sideBarLink.Url

	return todoLink
}
