package callsvc

import (
	"github.com/star-table/startable-server/common/core/errs"
	"github.com/star-table/startable-server/common/core/util/json"
)

type IssueCreateHandler struct{}

func (i IssueCreateHandler) Handle(cardReq CardReq) (string, errs.SystemErrorInfo) {
	log.Infof("[IssueCreateHandler] handle info:%s", json.ToJsonIgnoreError(cardReq))
	//value := cardReq.Action.Value
	//
	//// 操作人的 openId
	//opUserOpenId := cardReq.OpenId
	//issueId := value.IssueId
	//orgId := value.OrgId
	//optionsDate := cardReq.Action.Option
	//action := value.Action
	//
	//// 查询操作人信息
	//resp := orgfacade.GetBaseUserInfoByEmpId(orgvo.GetBaseUserInfoByEmpIdReqVo{
	//	OrgId: orgId,
	//	EmpId: opUserOpenId,
	//})
	//if resp.Failure() {
	//	return "", errs.BuildSystemErrorInfo(errs.UserNotFoundError, resp.Error())
	//}
	//opUserInfo := resp.BaseUserInfo
	//
	//issueLinks := projectfacade.GetIssueLinks(projectvo.GetIssueLinksReqVo{
	//	SourceChannel: SourceChannel,
	//	OrgId:         orgId,
	//	IssueId:       issueId,
	//}).NewData
	//
	//if (action == consts.FsCardActionUpdatePlanEndTime || action == consts.FsCardActionUpdatePlanStartTime) && optionsDate != "" {
	//	targetTime, parseErr := time.Parse(consts.AppTimeFormatYYYYMMDDHHmmTimezone, optionsDate)
	//	if parseErr != nil {
	//		log.Error(parseErr)
	//		return "", errs.BuildSystemErrorInfo(errs.DateParseError)
	//	}
	//	targetTypesTime := types.Time(targetTime)
	//	form := make([]map[string]interface{}, 0)
	//	data := make(map[string]interface{})
	//	data[consts.BasicFieldId] = issueId
	//
	//	if action == consts.FsCardActionUpdatePlanStartTime {
	//		data[consts.BasicFieldPlanStartTime] = targetTypesTime
	//	}
	//	if action == consts.FsCardActionUpdatePlanEndTime {
	//		data[consts.BasicFieldPlanEndTime] = targetTypesTime
	//	}
	//
	//	form = append(form, data)
	//
	//	updateRespVo := projectfacade.BatchUpdateIssue(projectvo.BatchUpdateIssueReqVo{
	//		UserId:    opUserInfo.UserId,
	//		OrgId:     orgId,
	//		AppId:     -1, // 这里拿不到，只能传-1
	//		ProjectId: -1, // 这里拿不到，只能传-1
	//		TableId:   -1, // 这里拿不到，只能传-1
	//		NewData:      form,
	//	})
	//	// 查询issue和project 相关信息
	//	// 查询出任务数据
	//	issuesResp := projectfacade.GetIssueInfoList(projectvo.IssueInfoListReqVo{
	//		UserId:   opUserInfo.UserId,
	//		OrgId:    orgId,
	//		IssueIds: []int64{issueId},
	//	})
	//	if issuesResp.Failure() {
	//		log.Errorf("[IssueCreateHandler] err:%v, orgId:%d, userId:%d", issuesResp.Error(), orgId, opUserInfo.UserId)
	//		return "", issuesResp.Error()
	//	}
	//
	//	issueInfo := issuesResp.IssueInfos[0]
	//	// 项目/表
	//	projectContent := fmt.Sprintf("**%s**%s%s", consts.CardElementProjectTable, consts.FsCard3Tab, consts.CardDefaultIssueProjectName)
	//
	//	// 负责人
	//	ownerDisplayName := opUserInfo.Name
	//	if ownerDisplayName == "" {
	//		ownerDisplayName = consts.CardDefaultOwnerNameForUpdateIssue
	//	}
	//
	//	if updateRespVo.Failure() {
	//		log.Errorf("[IssueCreateHandler] BatchUpdateIssue err:%v, orgId:%d, userId:%d",
	//			updateRespVo.Error(), orgId, opUserInfo.UserId)
	//		cardResp := card.GetFsCardCreateNewIssue(issueInfo.Title, projectContent, ownerDisplayName, issueLinks.SideBarLink, issueLinks.Link, issueInfo, updateRespVo.Message)
	//		return json.ToJsonIgnoreError(cardResp), nil
	//	}
	//
	//	cardResp := card.GetFsCardCreateNewIssue(issueInfo.Title, projectContent, ownerDisplayName, issueLinks.SideBarLink, issueLinks.Link, issueInfo, "")
	//	return json.ToJsonIgnoreError(cardResp), nil
	//}

	return "", nil

}
