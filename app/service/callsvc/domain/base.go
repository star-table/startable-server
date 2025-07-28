package callsvc

import (
	"github.com/star-table/startable-server/app/facade/orgfacade"
	"github.com/star-table/startable-server/app/facade/projectfacade"
	"github.com/star-table/startable-server/common/core/errs"
	"github.com/star-table/startable-server/common/model/vo/callvo"
	"github.com/star-table/startable-server/common/model/vo/orgvo"
	"github.com/star-table/startable-server/common/model/vo/projectvo"
)

func GetDingGroupChatHandleInfos(outOrgId string, chatId string, sourceChannel string) (*callvo.GroupChatHandleInfo, errs.SystemErrorInfo) {
	info := &callvo.GroupChatHandleInfo{
		SourceChannel: sourceChannel,
	}
	orgResp := orgfacade.GetBaseOrgInfoByOutOrgId(orgvo.GetBaseOrgInfoByOutOrgIdReqVo{OutOrgId: outOrgId})
	if orgResp.Failure() {
		log.Errorf("[GetGroupChatHandleInfos] GetBaseOrgInfoByOutOrgId err: %v, outOrgId: %v", orgResp.Error(),
			outOrgId)
		return info, orgResp.Error()
	}
	info.OrgInfo = *orgResp.BaseOrgInfo
	// 查询绑定的项目ids
	bindProResp := projectfacade.GetProjectIdsByChatId(projectvo.GetProjectIdsByChatIdReqVo{
		OrgId:      info.OrgInfo.OrgId,
		OpenChatId: chatId,
	})
	if bindProResp.Failure() {
		log.Errorf("[GetGroupChatHandleInfos] GetProjectIdsByChatId err: %v, outOrgId: %v", orgResp.Error(), outOrgId)
		return info, orgResp.Error()
	}
	info.BindProjectIds = bindProResp.Data.ProjectIds

	return info, nil
}
