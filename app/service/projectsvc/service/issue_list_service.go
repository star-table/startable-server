package service

import (
	"fmt"

	pushPb "gitea.bjx.cloud/LessCode/interface/golang/push/v1"
	sdk_const "gitea.bjx.cloud/allstar/platform-sdk/consts"
	"github.com/star-table/startable-server/app/facade/orgfacade"
	"github.com/star-table/startable-server/app/facade/trendsfacade"
	"github.com/star-table/startable-server/app/service/projectsvc/domain"
	"github.com/star-table/startable-server/common/core/consts"
	"github.com/star-table/startable-server/common/core/errs"
	"github.com/star-table/startable-server/common/core/logger"
	lang2 "github.com/star-table/startable-server/common/core/util/lang"
	"github.com/star-table/startable-server/common/extra/card"
	"github.com/star-table/startable-server/common/model/bo"
	"github.com/star-table/startable-server/common/model/bo/status"
	"github.com/star-table/startable-server/common/model/vo"
	"github.com/star-table/startable-server/common/model/vo/commonvo"
	"github.com/star-table/startable-server/common/model/vo/orgvo"
	"github.com/star-table/startable-server/common/model/vo/projectvo"
	"github.com/star-table/startable-server/common/model/vo/trendsvo"
)

var log = *logger.GetDefaultLogger()

// HandleGroupChatAtUserName 项目群聊-查询用户在项目内的任务统计信息。响应用户指令：@用户名a
func HandleGroupChatAtUserName(input *projectvo.HandleGroupChatUserInsAtUserNameReq) (bool, errs.SystemErrorInfo) {
	orgId, projectId, err := domain.GetProjectIdByOpenChatId(input.OpenChatId)
	if err != nil {
		log.Error(err)
		return false, errs.BuildSystemErrorInfo(errs.ParamError, err)
	}
	// 查询操作人信息
	resp1 := orgfacade.GetBaseUserInfoByEmpId(orgvo.GetBaseUserInfoByEmpIdReqVo{
		OrgId: orgId,
		EmpId: input.OpUserOpenId,
	})
	if resp1.Failure() {
		return false, resp1.Error()
	}
	opUser := resp1.BaseUserInfo
	// 查询用户信息
	resp := orgfacade.GetBaseUserInfoByEmpId(orgvo.GetBaseUserInfoByEmpIdReqVo{
		OrgId: orgId,
		EmpId: input.AtUserOpenId,
	})
	if resp.Failure() {
		log.Errorf("[HandleGroupChatAtUserName] err: %v, orgId: %d, ourOrgId: %s", err, orgId,
			resp.BaseUserInfo.OutOrgId)
		return false, resp.Error()
	}
	atUser := resp.BaseUserInfo
	// 机器人将信息回复给用户
	//tenant, err := feishu.GetTenant(resp.BaseUserInfo.OutOrgId)
	//if err != nil {
	//	log.Errorf("[HandleGroupChatAtUserName] err: %v, ourOrgId: %s", err, resp.BaseUserInfo.OutOrgId)
	//	return false, err
	//}
	// 查询负责人是被“at的用户”的任务统计信息
	relationType := consts.IssueRelationTypeOwner
	proStatResp, err := IssueStatusTypeStat(orgId, atUser.UserId, &vo.IssueStatusTypeStatReq{
		ProjectID:    &projectId,
		IterationID:  nil,
		RelationType: &relationType,
	})
	if err != nil {
		return false, err
	}
	// 检查用户是否是项目成员
	isIn, err := domain.IsProjectParticipant(orgId, atUser.UserId, projectId)
	if err != nil {
		return false, err
	}

	pushCard := commonvo.PushCard{
		OrgId:         orgId,
		OutOrgId:      resp.BaseUserInfo.OutOrgId,
		SourceChannel: sdk_const.SourceChannelFeishu,
		ChatIds:       []string{input.OpenChatId},
	}

	// 如果不是项目成员，则返回提示信息
	if !isIn {
		msgCard := card.GetFsCardForSomeTip("抱歉，您指定的用户不是项目成员。", map[string]interface{}{
			"atOpenId":   opUser.OutUserId,
			"atUserName": opUser.Name,
		})
		pushCard.CardMsg = msgCard
		errSys := card.PushCard(orgId, &pushCard)
		if errSys != nil {
			log.Errorf("[HandleGroupChatAtUserName] err:%v", errSys)
			return false, errSys
		}
		//fsResp, oriErr := tenant.SendMessage(fsvo.MsgVo{
		//	ChatId:  input.OpenChatId,
		//	MsgType: "interactive",
		//	Card:    msgCard,
		//})
		//if oriErr != nil {
		//	log.Error(oriErr)
		//	return false, err
		//}
		//if fsResp.Code != 0 {
		//	log.Error("HandleGroupChatAtUserName isIn 发送消息异常")
		//}
		return true, nil
	}
	project, err := domain.GetProjectSimple(orgId, projectId)
	if err != nil {
		log.Error(err)
		return false, err
	}
	var msgCard *pushPb.TemplateCard
	cardTitle := atUser.Name + " 负责的任务完成情况"
	if project.ProjectTypeId == consts.ProjectTypeAgileId {
		msgCard = card.GetFsCardProIssueCardForAgile(projectId, project.ProjectTypeId, cardTitle, proStatResp)
	} else {
		msgCard = card.GetFsCardProIssueCardForCommon(projectId, project.ProjectTypeId, cardTitle, proStatResp)
	}
	pushCard.CardMsg = msgCard
	errSys := card.PushCard(orgId, &pushCard)
	if errSys != nil {
		log.Errorf("[HandleGroupChatAtUserName] err:%v", errSys)
		return false, errSys
	}
	//fsResp, oriErr := tenant.SendMessage(fsvo.MsgVo{
	//	ChatId:  input.OpenChatId,
	//	MsgType: "interactive",
	//	Card:    msgCard,
	//})
	//if oriErr != nil {
	//	log.Error(oriErr)
	//	return false, err
	//}
	//if fsResp.Code != 0 {
	//	log.Error("HandleGroupChatAtUserName 发送消息异常")
	//}
	//log.Infof("飞书群聊-用户指令处理-HandleGroupChatAtUserName-响应用户调用结果 code: %v", fsResp.Code)

	return true, nil
}

// 项目群聊，响应用户指令：@用户名a 任务标题1。表示创建负责人为 `用户名a`，标题为 `任务标题1` 的任务。
func HandleGroupChatAtUserNameWithIssueTitle(input *projectvo.HandleGroupChatUserInsAtUserNameWithIssueTitleReq) (bool, errs.SystemErrorInfo) {
	orgId, projectId, err := domain.GetProjectIdByOpenChatId(input.OpenChatId)
	if err != nil {
		log.Error(err)
		return false, errs.BuildSystemErrorInfo(errs.ParamError, err)
	}
	resp := orgfacade.GetBaseUserInfoByEmpId(orgvo.GetBaseUserInfoByEmpIdReqVo{
		OrgId: orgId,
		EmpId: input.OpUserOpenId,
	})
	if resp.Failure() {
		return false, errs.BuildSystemErrorInfo(errs.UserNotFoundError, resp.Error())
	}
	opUserInfo := resp.BaseUserInfo
	resp = orgfacade.GetBaseUserInfoByEmpId(orgvo.GetBaseUserInfoByEmpIdReqVo{
		OrgId: orgId,
		EmpId: input.AtUserOpenId,
	})
	if resp.Failure() {
		return false, errs.BuildSystemErrorInfo(errs.UserNotFoundError, resp.Error())
	}
	atUserInfo := resp.BaseUserInfo
	// 机器人将信息回复给用户
	//tenant, err := feishu.GetTenant(opUserInfo.OutOrgId)
	//if err != nil {
	//	log.Error(err)
	//	return false, err
	//}
	// 检查创建者用户（操作人）是否是项目成员
	opUserIsIn, err := domain.IsProjectParticipant(orgId, opUserInfo.UserId, projectId)
	if err != nil {
		return false, err
	}

	pushCard := commonvo.PushCard{
		OrgId:         orgId,
		OutOrgId:      resp.BaseUserInfo.OutOrgId,
		SourceChannel: sdk_const.SourceChannelFeishu,
		ChatIds:       []string{input.OpenChatId},
	}

	if !opUserIsIn {
		msgCard := card.GetFsCardForSomeTip("抱歉，您不是项目成员，无法创建任务。", map[string]interface{}{
			"atOpenId":   opUserInfo.OutUserId,
			"atUserName": opUserInfo.Name,
		})
		pushCard.CardMsg = msgCard
		errSys := card.PushCard(orgId, &pushCard)
		if errSys != nil {
			log.Errorf("[HandleGroupChatAtUserNameWithIssueTitle] err:%v", errSys)
			return false, errSys
		}
		//fsResp, oriErr := tenant.SendMessage(fsvo.MsgVo{
		//	ChatId:  input.OpenChatId,
		//	MsgType: "interactive",
		//	Card:    msgCard,
		//})
		//if oriErr != nil {
		//	log.Error(oriErr)
		//	return false, err
		//}
		//if fsResp.Code != 0 {
		//	log.Error("HandleGroupChatAtUserName opUserIsIn 发送消息异常")
		//}
		return true, nil
	}
	// 检查负责人用户是否是项目成员
	atUserIsIn, err := domain.IsProjectParticipant(orgId, atUserInfo.UserId, projectId)
	if err != nil {
		return false, err
	}
	// 如果不是项目成员，则返回提示信息
	if !atUserIsIn {
		msgCard := card.GetFsCardForSomeTip("抱歉，您指定的用户不是项目成员。", map[string]interface{}{
			"atOpenId":   opUserInfo.OutUserId,
			"atUserName": opUserInfo.Name,
		})
		pushCard.CardMsg = msgCard
		errSys := card.PushCard(orgId, &pushCard)
		if errSys != nil {
			log.Errorf("[HandleGroupChatAtUserNameWithIssueTitle] err:%v", errSys)
			return false, errSys
		}
		//fsResp, oriErr := tenant.SendMessage(fsvo.MsgVo{
		//	ChatId:  input.OpenChatId,
		//	MsgType: "interactive",
		//	Card:    msgCard,
		//})
		//if oriErr != nil {
		//	log.Error(oriErr)
		//	return false, err
		//}
		//if fsResp.Code != 0 {
		//	log.Error("HandleGroupChatAtUserName isIn 发送消息异常")
		//}
		return true, nil
	}
	// 创建任务 start
	//priorityList, err := PriorityList(orgId, 1, 10000, db.Cond{
	//	consts.TcType: consts.PriorityTypeIssue,
	//})
	//if err != nil {
	//	log.Error(err)
	//	return false, err
	//}
	//var defaultPriority *vo.Priority = nil
	//list := priorityList.List
	//if len(list) > 0 {
	//	for _, priority := range list {
	//		if priority.IsDefault == 1 {
	//			defaultPriority = priority
	//			break
	//		}
	//	}
	//	if defaultPriority == nil {
	//		defaultPriority = list[len(list)-1]
	//	}
	//}
	tableId, err := getProjectFirstTableId(orgId, projectId)
	if err != nil {
		log.Errorf("[HandleGroupChatAtUserNameWithIssueTitle] getProjectFirstTableId, orgId:%v, projectId:%v, err:%v", orgId, projectId, err)
		return false, err
	}

	project, err := domain.GetProjectSimple(orgId, projectId)
	if err != nil {
		log.Errorf("[HandleGroupChatAtUserNameWithIssueTitle] GetProjectSimple, orgId:%v, projectId:%v, err:%v", orgId, projectId, err)
		return false, err
	}

	// TODO: 默认优先级？
	d := map[string]interface{}{
		consts.BasicFieldTitle:   input.IssueTitleStr,
		consts.BasicFieldOwnerId: []string{fmt.Sprintf("U_%d", atUserInfo.UserId)},
	}

	req := &projectvo.BatchCreateIssueReqVo{
		OrgId:  orgId,
		UserId: opUserInfo.UserId,
		Input: &projectvo.BatchCreateIssueInput{
			AppId:     project.AppId,
			ProjectId: projectId,
			TableId:   tableId,
			Data:      []map[string]interface{}{d},
		},
	}
	_, _, _, err = BatchCreateIssue(req, false, &projectvo.TriggerBy{
		TriggerBy: consts.TriggerByNormal,
	}, "", 0)
	//_, err = CreateIssue(projectvo.CreateIssueReqVo{
	//	UserId:        opUserInfo.UserId,
	//	OrgId:         orgId,
	//	SourceChannel: sdk_const.SourceChannelFeishu,
	//	CreateIssue: vo.CreateIssueReq{
	//		ProjectID:  projectId,
	//		Title:      input.IssueTitleStr,
	//		OwnerID:    []int64{atUserInfo.UserId},
	//		PriorityID: defaultPriority.ID,
	//		TableID:    cast.ToString(tableId),
	//	},
	//})
	// 创建任务异常时，回复提示信息卡片到群聊中。
	if err != nil {
		msgCard := card.GetFsCardForSomeTip(err.Message(), map[string]interface{}{
			"atOpenId":   opUserInfo.OutUserId,
			"atUserName": opUserInfo.Name,
		})
		pushCard.CardMsg = msgCard
		errSys := card.PushCard(orgId, &pushCard)
		if errSys != nil {
			log.Errorf("[HandleGroupChatAtUserNameWithIssueTitle] err:%v", errSys)
			return false, errSys
		}
		//fsResp, oriErr := tenant.SendMessage(fsvo.MsgVo{
		//	ChatId:  input.OpenChatId,
		//	MsgType: "interactive",
		//	Card:    msgCard,
		//})
		//if oriErr != nil {
		//	log.Error(oriErr)
		//	return false, err
		//}
		//if fsResp.Code != 0 {
		//	log.Error("HandleGroupChatAtUserName createIssueFailed 发送消息异常")
		//}
		//log.Error(err)
		return false, err
	}
	// 创建任务 end
	log.Infof("飞书群聊-用户指令处理-HandleGroupChatAtUserNameWithIssueTitle-执行完成。")

	return true, nil
}

func getProjectFirstTableId(orgId, projectId int64) (int64, errs.SystemErrorInfo) {
	appId, err := domain.GetAppIdFromProjectId(orgId, projectId)
	if err != nil {
		return 0, err
	}
	tables, err := domain.GetAppTableList(orgId, appId)
	if err != nil {
		return 0, err
	}
	if len(tables) > 0 {
		return tables[0].TableId, nil
	}

	return 0, nil
}

// 项目群聊，响应用户指令，项目进展
func HandleGroupChatUserInsProProgress(input *projectvo.HandleGroupChatUserInsProProgressReq) (bool, errs.SystemErrorInfo) {
	orgId, projectId, err := domain.GetProjectIdByOpenChatId(input.OpenChatId)
	if err != nil {
		log.Error(err)
		return false, errs.BuildSystemErrorInfo(errs.ParamError, err)
	}
	resp := orgfacade.GetBaseUserInfoByEmpId(orgvo.GetBaseUserInfoByEmpIdReqVo{
		OrgId: orgId,
		EmpId: input.OpenId,
	})
	if resp.Failure() {
		return false, resp.Error()
	}
	// 获取项目进展
	proStatusList, err := domain.GetProjectStatus(orgId, projectId)
	if err != nil {
		return false, errs.BuildSystemErrorInfo(errs.CacheProxyError, err)
	}
	statusMap := make(map[int64]bo.CacheProcessStatusBo, 0)
	for _, item := range proStatusList {
		statusMap[item.StatusId] = item
	}
	projectBo, err1 := domain.GetProjectSimple(orgId, projectId)
	if err1 != nil {
		log.Error(err1)
		return false, errs.BuildSystemErrorInfo(errs.ProjectNotExist)
	}
	proStatusName := ""
	if tmp, ok := statusMap[projectBo.Status]; ok {
		proStatusName = tmp.Name
	}
	// 如果是归档项目，则显示归档状态
	if projectBo.IsFiling == consts.ProjectIsFiling {
		proStatusName = consts.ProjectIsFilingDesc
	}
	// 机器人将信息回复给用户
	err = GroupChatReplyToUserInProgress(input.SourceChannel, resp.BaseUserInfo.OutOrgId, input.OpenChatId, projectBo.OrgId, projectId, proStatusName)
	if err != nil {
		log.Errorf("[HandleGroupChatUserInsProProgress] err:%v", err)
		return false, err
	}

	return true, nil
}

// HandleGroupChatUserInsProjectSettings 项目群聊指令，响应用户指令，群聊推送设置
func HandleGroupChatUserInsProjectSettings(openChatId string, sourceChannel string) (bool, errs.SystemErrorInfo) {
	orgId, _, err := domain.GetProjectIdByOpenChatId(openChatId)
	if err != nil {
		log.Error(err)
		return false, errs.BuildSystemErrorInfo(errs.ParamError, err)
	}
	resp := orgfacade.GetBaseOrgInfo(orgvo.GetBaseOrgInfoReqVo{OrgId: orgId})
	if resp.Failure() {
		return false, resp.Error()
	}

	err = GroupChatReplyToUserInsProjectSettings(sourceChannel, resp.BaseOrgInfo.OutOrgId, openChatId, orgId)
	if err != nil {
		log.Errorf("[HandleGroupChatUserInsProjectSettings] err: %v", err)
		return false, nil
	}

	return true, nil
}

// 项目群聊，响应用户指令，项目任务，通过 chat id 等参数，获取对应项目下任务的统计信息
func HandleGroupChatUserInsProIssue(input *projectvo.HandleGroupChatUserInsProIssueReq) (bool, errs.SystemErrorInfo) {
	orgId, projectId, err := domain.GetProjectIdByOpenChatId(input.OpenChatId)
	if err != nil {
		log.Error(err)
		return false, errs.BuildSystemErrorInfo(errs.ParamError, err)
	}
	resp := orgfacade.GetBaseUserInfoByEmpId(orgvo.GetBaseUserInfoByEmpIdReqVo{
		OrgId: orgId,
		EmpId: input.OpenId,
	})
	if resp.Failure() {
		return false, resp.Error()
	}
	// user := resp.BaseUserInfo
	proStatResp, err := IssueStatusTypeStat(orgId, 0, &vo.IssueStatusTypeStatReq{
		ProjectID:    &projectId,
		IterationID:  nil,
		RelationType: nil,
	})
	if err != nil {
		return false, err
	}
	projectInfo, err := domain.GetProjectSimple(orgId, projectId)
	if err != nil {
		log.Error(err)
		return false, err
	}
	// 机器人将信息回复给用户
	errSys := GroupChatReplyToUserInsProIssue(input.SourceChannel, resp.BaseUserInfo.OutOrgId, input.OpenChatId, projectInfo.OrgId, projectId, projectInfo.ProjectTypeId, proStatResp)
	if errSys != nil {
		log.Errorf("[HandleGroupChatUserInsProIssue] failed:%v", errSys)
		return false, errSys
	}

	return true, nil
}

func IssueStatusTypeStat(orgId, currentUserId int64, input *vo.IssueStatusTypeStatReq) (*vo.IssueStatusTypeStatResp, errs.SystemErrorInfo) {
	if input.ProjectID != nil && !domain.JudgeProjectIsExist(orgId, *input.ProjectID) {
		return nil, errs.BuildSystemErrorInfo(errs.ProjectNotExist)
	}
	if input.IterationID != nil {
		iterationBo, err := domain.GetIterationBoByOrgId(*input.IterationID, orgId)
		if err != nil {
			log.Error(err)
			return nil, errs.BuildSystemErrorInfo(errs.IterationNotExist)
		}
		input.ProjectID = &iterationBo.ProjectId
	}
	// 兼容一下查询项目下的各个状态任务的数量，而非只查询某项目下某个成员的任务统计。
	issueStatusStatBos, err1 := domain.GetIssueStatusStatWithLc(bo.IssueStatusStatCondBo{
		OrgId:        orgId,
		ProjectId:    input.ProjectID,
		IterationId:  input.IterationID,
		RelationType: input.RelationType,
		UserId:       currentUserId,
	})
	if err1 != nil {
		log.Error(err1)
		return nil, errs.BuildSystemErrorInfo(errs.IssueDomainError, err1)
	}
	result := &vo.IssueStatusTypeStatResp{}

	for _, issueStatusStatBo := range issueStatusStatBos {
		result.NotStartTotal += int64(issueStatusStatBo.IssueWaitCount)
		result.ProcessingTotal += int64(issueStatusStatBo.IssueRunningCount)
		result.CompletedTotal += int64(issueStatusStatBo.IssueEndCount)
		result.CompletedTodayTotal += int64(issueStatusStatBo.IssueEndTodayCount)
		result.Total += int64(issueStatusStatBo.IssueCount)
		result.OverdueCompletedTotal += int64(issueStatusStatBo.IssueOverdueEndCount)
		result.OverdueTodayTotal += int64(issueStatusStatBo.IssueOverdueTodayCount)
		result.OverdueTotal += int64(issueStatusStatBo.IssueOverdueCount)
		result.OverdueTomorrowTotal += int64(issueStatusStatBo.IssueOverdueTomorrowCount)
		result.TodayCount += int64(issueStatusStatBo.TodayCount)
		result.TodayCreateCount += int64(issueStatusStatBo.TodayCreateCount)
		result.WaitConfirmedTotal += int64(issueStatusStatBo.WaitConfirmedCount)
	}

	//即将到期 今天到期的主子任务数+明日逾期的主子任务数
	result.BeAboutToOverdueSum = result.OverdueTomorrowTotal + result.OverdueTodayTotal

	//@我的数量
	if currentUserId > 0 {
		callMeCountResp := trendsfacade.CallMeCount(trendsvo.CallMeCountReqVo{
			ProjectId: input.ProjectID,
			UserId:    currentUserId,
			OrgId:     orgId,
		})
		if callMeCountResp.Failure() {
			log.Error(callMeCountResp.Error())
			return nil, callMeCountResp.Error()
		}
		result.CallMeTotal = callMeCountResp.Count
	}

	result.List = append(result.List, &vo.StatCommon{
		Name:  "已逾期",
		Count: result.OverdueTotal,
	})
	result.List = append(result.List, &vo.StatCommon{
		Name:  "进行中",
		Count: result.ProcessingTotal,
	})
	result.List = append(result.List, &vo.StatCommon{
		Name:  "未完成",
		Count: result.NotStartTotal + result.ProcessingTotal,
	})
	result.List = append(result.List, &vo.StatCommon{
		Name:  "已完成",
		Count: result.CompletedTotal,
	})
	lang := lang2.GetLang()
	isOtherLang := lang2.IsEnglish()
	if isOtherLang {
		for index, item := range result.List {
			if tmpMap, ok1 := consts.LANG_ISSUE_STAT_DESC_MAP[lang]; ok1 {
				if tmpVal, ok2 := tmpMap[item.Name]; ok2 {
					(*result.List[index]).Name = tmpVal
				}
			}
		}
	}

	return result, nil
}

func IssueStatusTypeStatDetail(orgId, currentUserId int64, input *vo.IssueStatusTypeStatReq) (*vo.IssueStatusTypeStatDetailResp, errs.SystemErrorInfo) {
	projectId := input.ProjectID

	if projectId != nil && !domain.JudgeProjectIsExist(orgId, *projectId) {
		return nil, errs.BuildSystemErrorInfo(errs.ProjectNotExist)
	}

	issueStatusStatBos, err1 := domain.GetIssueStatusStatWithLc(bo.IssueStatusStatCondBo{
		OrgId:        orgId,
		UserId:       currentUserId,
		ProjectId:    input.ProjectID,
		IterationId:  input.IterationID,
		RelationType: input.RelationType,
	})
	if err1 != nil {
		log.Error(err1)
		return nil, errs.BuildSystemErrorInfo(errs.IssueDomainError, err1)
	}

	result := &vo.IssueStatusTypeStatDetailResp{
		NotStart:   []*vo.IssueStatByObjectType{},
		Processing: []*vo.IssueStatByObjectType{},
		Completed:  []*vo.IssueStatByObjectType{},
	}

	for _, issueStatusStatBo := range issueStatusStatBos {
		temp := issueStatusStatBo
		if issueStatusStatBo.IssueWaitCount > 0 {
			result.NotStart = append(result.NotStart, &vo.IssueStatByObjectType{
				ProjectObjectTypeID:   &temp.ProjectTypeId,
				ProjectObjectTypeName: &temp.ProjectTypeName,
				Total:                 int64(temp.IssueWaitCount),
			})
		}
		if issueStatusStatBo.IssueRunningCount > 0 {
			result.Processing = append(result.Processing, &vo.IssueStatByObjectType{
				ProjectObjectTypeID:   &temp.ProjectTypeId,
				ProjectObjectTypeName: &temp.ProjectTypeName,
				Total:                 int64(temp.IssueRunningCount),
			})
		}
		if issueStatusStatBo.IssueEndCount > 0 {
			result.Completed = append(result.Completed, &vo.IssueStatByObjectType{
				ProjectObjectTypeID:   &temp.ProjectTypeId,
				ProjectObjectTypeName: &temp.ProjectTypeName,
				Total:                 int64(temp.IssueEndCount),
			})
		}
	}

	return result, nil
}

//func GetSimpleIssueInfoBatch(orgId int64, ids []int64) (*[]vo.Issue, errs.SystemErrorInfo) {
//	list, _, err := domain.SelectList(db.Cond{
//		consts.TcOrgId: orgId,
//		consts.TcId:    db.In(ids),
//		//consts.TcIsDelete: consts.AppIsNoDelete,
//	}, nil, 0, 0, nil, false)
//	if err != nil {
//		log.Error(err)
//		return nil, err
//	}
//	issueVo := &[]vo.Issue{}
//	copyErr := copyer.Copy(list, issueVo)
//	if copyErr != nil {
//		log.Error(copyErr)
//		return nil, errs.BuildSystemErrorInfo(errs.ObjectCopyError, copyErr)
//	}
//
//	return issueVo, nil
//}

func GetLcIssueInfoBatch(orgId int64, issueIds []int64) ([]*bo.IssueBo, errs.SystemErrorInfo) {
	issueBos, err := domain.GetIssueInfosLc(orgId, 0, issueIds)
	if err != nil {
		log.Errorf("[GetSimpleIssueInfoBatch] domain.GetIssueInfosLc err:%v, orgId:%v, issueIds:%v", err, orgId, issueIds)
		return nil, err
	}
	return issueBos, nil
}

//func GetIssueRemindInfoList(reqVo projectvo.GetIssueRemindInfoListReqVo) (*projectvo.GetIssueRemindInfoListRespData, errs.SystemErrorInfo) {
//	if reqVo.Page < 0 {
//		return nil, errs.BuildSystemErrorInfo(errs.PageInvalidError)
//	}
//	if reqVo.Size < 0 || reqVo.Size > 100 {
//		return nil, errs.BuildSystemErrorInfo(errs.PageSizeInvalidError)
//	}
//
//	selectIssueIdsCondBo := bo.SelectIssueIdsCondBo{}
//
//	input := reqVo.Input
//	//计划结束时间条件
//	selectIssueIdsCondBo.BeforePlanEndTime = input.BeforePlanEndTime
//	selectIssueIdsCondBo.AfterPlanEndTime = input.AfterPlanEndTime
//	selectIssueIdsCondBo.BeforePlanStartTime = input.BeforePlanStartTime
//	selectIssueIdsCondBo.AfterPlanStartTime = input.AfterPlanStartTime
//
//	issueRemindInfos, total, err := domain.SelectIssueRemindInfoList(selectIssueIdsCondBo, reqVo.Page, reqVo.Size)
//	if err != nil {
//		log.Error(err)
//		return nil, err
//	}
//
//	return &projectvo.GetIssueRemindInfoListRespData{
//		Total: total,
//		List:  issueRemindInfos,
//	}, nil
//}

func IssueListStat(orgId, userId, projectId int64) (*vo.IssueListStatResp, errs.SystemErrorInfo) {
	projectInfo, projectErr := domain.GetProject(orgId, projectId)
	if projectErr != nil {
		log.Error(projectErr)
		return nil, errs.ProjectNotExist
	}
	projectTypeList, err := domain.GetAppTableList(orgId, projectInfo.AppId)
	if err != nil {
		log.Error(err)
		return nil, err
	}

	data, err := domain.GetIssueStatusStatWithLc(bo.IssueStatusStatCondBo{
		OrgId:     orgId,
		UserId:    userId,
		ProjectId: &projectId,
	})
	if err != nil {
		log.Error(err)
		return nil, err
	}
	midStat := map[int64]vo.IssueListStatData{}
	for _, datum := range data {
		midStat[datum.ProjectTypeId] = vo.IssueListStatData{
			Total:                 int64(datum.IssueCount),
			FinishedCount:         int64(datum.IssueEndCount),
			OverdueCount:          int64(datum.IssueOverdueCount),
			ProjectObjectTypeID:   datum.ProjectTypeId,
			ProjectObjectTypeName: datum.ProjectTypeName,
		}
	}

	resStat := []*vo.IssueListStatData{}

	for _, objectType := range projectTypeList {
		if _, ok := midStat[objectType.TableId]; ok {
			mid := midStat[objectType.TableId]
			resStat = append(resStat, &mid)
		} else {
			resStat = append(resStat, &vo.IssueListStatData{
				ProjectObjectTypeID:   objectType.TableId,
				ProjectObjectTypeName: objectType.Name,
				Total:                 0,
				FinishedCount:         0,
				OverdueCount:          0,
			})
		}
	}

	return &vo.IssueListStatResp{List: resStat}, nil
}

//func HomeIssuesGroup(orgId, currentUserId int64, page int, size int, input *vo.HomeIssueInfoReq) (*vo.HomeIssueInfoGroupResp, errs.SystemErrorInfo) {
//	if input.ProjectID == nil || *input.ProjectID == 0 {
//		return nil, errs.ParamError
//	}
//	projectId := *input.ProjectID
//	projectInfo, err := domain.GetProject(orgId, projectId)
//	if err != nil {
//		log.Error(err)
//		return nil, err
//	}
//
//	list, err := HomeIssuesForTableMode(orgId, currentUserId, page, size, input, false)
//	if err != nil {
//		log.Error(err)
//		return nil, err
//	}
//	timeSpan := int64(0)
//	if len(list.List) > 0 {
//		var timeRange []time.Time
//		for _, info := range list.List {
//			start := time.Time(info.Issue.PlanStartTime)
//			end := time.Time(info.Issue.PlanEndTime)
//
//			if start.After(consts.BlankTimeObject) && end.After(consts.BlankTimeObject) {
//				timeRange = append(timeRange, start)
//				timeRange = append(timeRange, end)
//			}
//		}
//
//		if len(timeRange) > 0 {
//			minTime := timeRange[0]
//			maxTime := timeRange[0]
//			for _, t := range timeRange {
//				if t.Before(minTime) {
//					minTime = t
//				}
//				if t.After(maxTime) {
//					maxTime = t
//				}
//			}
//
//			timeSpan = maxTime.Unix() - minTime.Unix()
//		}
//	}
//
//	commonRes := &vo.HomeIssueInfoGroupResp{
//		Total:       list.Total,
//		ActualTotal: list.ActualTotal,
//		TimeSpan:    timeSpan,
//		Group: []*vo.HomeIssueGroup{
//			&vo.HomeIssueGroup{
//				ID:        0,
//				Name:      "",
//				Avatar:    "",
//				BgStyle:   "",
//				FontStyle: "",
//				TimeSpan:  0,
//				List:      list.List,
//			},
//		},
//	}
//	if input.GroupType == nil {
//		return commonRes, nil
//	}
//	lang := lang2.GetLang()
//	isOtherLang := lang2.IsEnglish()
//	//任务栏
//	//projectTypeList, err := ProjectObjectTypesWithProject(orgId, projectId)
//	//if err != nil {
//	//	log.Error(err)
//	//	return nil, err
//	//}
//	trulyProjectObjectTypeList := []vo.ProjectObjectType{}
//	//projectObjectTypeIds := []int64{}
//	//for _, objectType := range projectTypeList.List {
//	//	if projectInfo.ProjectTypeId == consts.ProjectTypeAgileId && objectType.Name == "迭代" {
//	//		continue
//	//	}
//	//	projectObjectTypeIds = append(projectObjectTypeIds, objectType.ID)
//	//	temp := objectType
//	//	trulyProjectObjectTypeList = append(trulyProjectObjectTypeList, *temp)
//	//}
//	//if input.ProjectObjectTypeID != nil && *input.ProjectObjectTypeID != 0 {
//	//	projectObjectTypeIds = []int64{*input.ProjectObjectTypeID}
//	//}
//	//获取项目所有的状态
//	tableList, tableListErr := domain.GetAppTableList(orgId, projectInfo.AppId)
//	if tableListErr != nil {
//		log.Errorf("[HomeIssuesGroup] GetAppTableList failed:%d, orgId:%d, appId:%d", tableListErr, orgId, projectInfo.AppId)
//		return nil, tableListErr
//	}
//	tableIds := []int64{}
//	for _, typeBo := range tableList {
//		tableIds = append(tableIds, typeBo.Id)
//	}
//	//allStatus, err := domain.GetIssueAllStatusNew(orgId, []int64{projectId}, tableIds) //这里每个表中的statusId可能会一样，下面的分组需要注意
//	//if err != nil {
//	//	log.Error(err)
//	//	return nil, err
//	//}
//	//statusMap := map[int][]int64{}
//	//for _, bos := range allStatus {
//	//	for _, infoBo := range bos {
//	//		if temp, ok := statusMap[infoBo.Type]; ok {
//	//			if ok1, _ := slice.Contain(temp, infoBo.ID); !ok1 {
//	//				statusMap[infoBo.Type] = append(statusMap[infoBo.Type], infoBo.ID)
//	//			}
//	//		} else {
//	//			statusMap[infoBo.Type] = append(statusMap[infoBo.Type], infoBo.ID)
//	//		}
//	//	}
//	//}
//	group := []*vo.HomeIssueGroup{}
//	switch *input.GroupType {
//	case 1:
//		//获取所有任务的负责人
//		//userIds := []int64{}
//		//for _, info := range list.List {
//		//	userIds = append(userIds, info.Issue.Owner)
//		//}
//		//userIds = slice.SliceUniqueInt64(userIds)
//		//ownerInfos, err := orgfacade.GetBaseUserInfoBatchRelaxed("", orgId, userIds)
//		//if err != nil {
//		//	log.Error(err)
//		//	return nil, err
//		//}
//		//userMap := maps.NewMap("UserId", ownerInfos)
//		//for _, id := range userIds {
//		//	if _, ok := userMap[id]; ok {
//		//		temp := userMap[id].(bo.BaseUserInfoBo)
//		//		group = append(group, &vo.HomeIssueGroup{
//		//			ID:       temp.UserId,
//		//			Name:     temp.Name,
//		//			Avatar:   temp.Avatar,
//		//			TimeSpan: 0,
//		//			List:     []*vo.HomeIssueInfo{},
//		//		})
//		//	}
//		//}
//	case 2:
//		//状态
//		//if len(statusMap) == 3 {
//		//	group = append(group, &vo.HomeIssueGroup{
//		//		ID:        1,
//		//		Name:      "未开始",
//		//		Avatar:    "#FFFFFF",
//		//		BgStyle:   "#DBDBDB",
//		//		FontStyle: "",
//		//		TimeSpan:  0,
//		//		List:      []*vo.HomeIssueInfo{},
//		//	})
//		//	//敏捷包含进行中
//		//	group = append(group, &vo.HomeIssueGroup{
//		//		ID:        2,
//		//		Name:      "进行中",
//		//		Avatar:    "",
//		//		BgStyle:   "#FFCD1C",
//		//		FontStyle: "#FFFFFF",
//		//		TimeSpan:  0,
//		//		List:      []*vo.HomeIssueInfo{},
//		//	})
//		//} else {
//		//	group = append(group, &vo.HomeIssueGroup{
//		//		ID:        1,
//		//		Name:      "未完成",
//		//		Avatar:    "#FFFFFF",
//		//		BgStyle:   "#DBDBDB",
//		//		FontStyle: "",
//		//		TimeSpan:  0,
//		//		List:      []*vo.HomeIssueInfo{},
//		//	})
//		//}
//		group = append(group, &vo.HomeIssueGroup{
//			ID:        3,
//			Name:      "已完成",
//			Avatar:    "",
//			BgStyle:   "#69A922",
//			FontStyle: "#FFFFFF",
//			TimeSpan:  0,
//			List:      []*vo.HomeIssueInfo{},
//		})
//		if isOtherLang {
//			otherLanguageMap := make(map[string]string, 0)
//			if tmpMap, ok1 := consts.LANG_ISSUE_STAT_DESC_MAP[lang]; ok1 {
//				otherLanguageMap = tmpMap
//			}
//			for index, item := range group {
//				if tmpVal, ok2 := otherLanguageMap[item.Name]; ok2 {
//					group[index].Name = tmpVal
//				}
//			}
//		}
//	case 3:
//		//优先级
//		//priorities, err := domain.GetPriorityListByType(orgId, consts.PriorityTypeIssue)
//		//if err != nil {
//		//	log.Error(err)
//		//	return nil, err
//		//}
//		//bo.SortPriorityBo(*priorities)
//		//for _, priorityBo := range *priorities {
//		//	group = append(group, &vo.HomeIssueGroup{
//		//		ID:        priorityBo.Id,
//		//		Name:      priorityBo.Name,
//		//		Avatar:    "",
//		//		BgStyle:   priorityBo.BgStyle,
//		//		FontStyle: priorityBo.FontStyle,
//		//		TimeSpan:  0,
//		//		List:      []*vo.HomeIssueInfo{},
//		//	})
//		//}
//	case 4:
//		//任务栏
//		for _, objectType := range trulyProjectObjectTypeList {
//			group = append(group, &vo.HomeIssueGroup{
//				ID:        objectType.ID,
//				Name:      objectType.Name,
//				Avatar:    "",
//				BgStyle:   objectType.BgStyle,
//				FontStyle: objectType.FontStyle,
//				TimeSpan:  0,
//				List:      []*vo.HomeIssueInfo{},
//			})
//		}
//	case 5:
//		//迭代
//		iterationList, _, err := domain.GetIterationBoList(0, 0, db.Cond{
//			consts.TcIsDelete:  consts.AppIsNoDelete,
//			consts.TcProjectId: projectId,
//			consts.TcOrgId:     orgId,
//		}, nil)
//		if err != nil {
//			log.Error(err)
//			return nil, err
//		}
//		for _, iterationBo := range *iterationList {
//			group = append(group, &vo.HomeIssueGroup{
//				ID:        iterationBo.Id,
//				Name:      iterationBo.Name,
//				Avatar:    "",
//				BgStyle:   "",
//				FontStyle: "",
//				TimeSpan:  0,
//				List:      []*vo.HomeIssueInfo{},
//			})
//		}
//		// 迭代的"未规划"多语言
//		notPlanName := "未规划"
//		if isOtherLang {
//			otherLanguageMap := make(map[string]string, 0)
//			if tmpMap, ok1 := consts.LANG_ROLE_NAME_MAP[lang]; ok1 {
//				otherLanguageMap = tmpMap
//			}
//			if tmpVal, ok2 := otherLanguageMap[notPlanName]; ok2 {
//				notPlanName = tmpVal
//			}
//		}
//		group = append(group, &vo.HomeIssueGroup{
//			ID:        0,
//			Name:      notPlanName,
//			Avatar:    "",
//			BgStyle:   "",
//			FontStyle: "",
//			TimeSpan:  0,
//			List:      []*vo.HomeIssueInfo{},
//		})
//	case 6:
//		//具体状态
//		//statusArr := []int64{}
//		//for _, bos := range allStatus {
//		//	for _, infoBo := range bos {
//		//		//去重
//		//		if ok, _ := slice.Contain(statusArr, infoBo.ID); ok {
//		//			continue
//		//		}
//		//		group = append(group, &vo.HomeIssueGroup{
//		//			ID:        infoBo.ID,
//		//			Name:      infoBo.Name,
//		//			Avatar:    "",
//		//			BgStyle:   infoBo.BgStyle,
//		//			FontStyle: infoBo.FontStyle,
//		//			TimeSpan:  0,
//		//			List:      []*vo.HomeIssueInfo{},
//		//		})
//		//		statusArr = append(statusArr, infoBo.ID)
//		//	}
//		//}
//	}
//	if len(group) == 0 {
//		return commonRes, nil
//	}
//
//	for i, issueGroup := range group {
//		for _, info := range list.List {
//			switch *input.GroupType {
//			case 1:
//				//负责人
//				//if info.Issue.Owner == issueGroup.ID {
//				//	group[i].List = append(group[i].List, info)
//				//}
//			case 2:
//				//状态
//				//if status, ok := statusMap[int(issueGroup.ID)]; ok {
//				//	if ok1, _ := slice.Contain(status, info.Issue.Status); ok1 {
//				//		group[i].List = append(group[i].List, info)
//				//	}
//				//}
//			case 3:
//				//优先级
//				if info.Issue.PriorityID == issueGroup.ID {
//					group[i].List = append(group[i].List, info)
//				}
//			case 4:
//				//任务栏
//				if info.Issue.ProjectObjectTypeID == issueGroup.ID {
//					group[i].List = append(group[i].List, info)
//				}
//			case 5:
//				//迭代
//				if info.Issue.IterationID == issueGroup.ID {
//					group[i].List = append(group[i].List, info)
//				}
//			case 6:
//				//具体状态
//				if info.Issue.Status == issueGroup.ID {
//					group[i].List = append(group[i].List, info)
//				}
//			}
//		}
//	}
//
//	for i, issueGroup := range group {
//		var timeRange []time.Time
//		for _, info := range issueGroup.List {
//			start := time.Time(info.Issue.PlanStartTime)
//			end := time.Time(info.Issue.PlanEndTime)
//
//			if start.After(consts.BlankTimeObject) && end.After(consts.BlankTimeObject) {
//				timeRange = append(timeRange, start)
//				timeRange = append(timeRange, end)
//				group[i].FitTotal++
//			}
//		}
//
//		if len(timeRange) == 0 {
//			continue
//		}
//		minTime := timeRange[0]
//		maxTime := timeRange[0]
//		for _, t := range timeRange {
//			if t.Before(minTime) {
//				minTime = t
//			}
//			if t.After(maxTime) {
//				maxTime = t
//			}
//		}
//		group[i].TimeSpan = maxTime.Unix() - minTime.Unix()
//	}
//
//	return &vo.HomeIssueInfoGroupResp{
//		Total:       list.Total,
//		ActualTotal: list.ActualTotal,
//		TimeSpan:    timeSpan,
//		Group:       group,
//	}, nil
//}

// GetIssueIdsByOrgId 根据一个组织的一批任务 id 列表
//func GetIssueIdsByOrgId(orgId, userId int64, input *projectvo.GetIssueIdsByOrgIdReq) (*projectvo.GetIssueIdsByOrgIdResp, errs.SystemErrorInfo) {
//	result := &projectvo.GetIssueIdsByOrgIdResp{
//		List: make([]int64, 0),
//	}
//	ids, total, err := domain.GetIssueIdsByOrgId(orgId, input.Page, input.Size)
//	if err != nil {
//		log.Error(err)
//		return result, err
//	}
//	result.List = ids
//	result.Total = total
//
//	return result, nil
//}

// InsertIssueProRelation 根据任务 id，将其与项目的关联，新增一条关联数据到 issue_relation 中。
//func InsertIssueProRelation(orgId, userId int64, input *projectvo.InsertIssueProRelationReq) errs.SystemErrorInfo {
//	issueIds := input.IssueIds
//	// 检查该任务是否已存在关联关系，已存在则跳过，不存在，则插入关联关系
//	existRelations, err := domain.GetIssueRelationsByCond(issueIds, []int{consts.IssueRelationTypeBelongManyPro}, db.Cond{
//		consts.TcOrgId: orgId,
//	})
//	if err != nil {
//		log.Errorf("[InsertIssueProRelation] err: %v", err)
//		return err
//	}
//	existIssueIds := make([]int64, 0)
//	for _, item := range existRelations {
//		existIssueIds = append(existIssueIds, item.IssueId)
//	}
//	notExistIssueIds := int642.ArrayDiff(issueIds, existIssueIds)
//	if len(notExistIssueIds) < 1 {
//		// 没有要处理的任务
//		return nil
//	}
//	// 查询任务信息
//	cond1 := db.Cond{
//		consts.TcOrgId: orgId,
//		consts.TcId:    notExistIssueIds,
//	}
//	bos, _, err := domain.SelectList(cond1, nil, 1, 10_0000, "id desc", false)
//	if err != nil {
//		log.Error(err)
//		return err
//	}
//	waitRelationBos := make([]*po.PpmPriIssueRelation, 0)
//	// 组装 issue relation
//	oids, err := idfacade.ApplyMultiplePrimaryIdRelaxed(consts.TableIssueRelation, len(*bos))
//	if err != nil {
//		log.Error(err)
//		return err
//	}
//	for index, issueBo := range *bos {
//		oid := oids.Ids[index].Id
//		waitRelationBos = append(waitRelationBos, &po.PpmPriIssueRelation{
//			Id:           oid,
//			OrgId:        issueBo.OrgId,
//			ProjectId:    issueBo.ProjectId,
//			IssueId:      issueBo.Id,
//			RelationId:   issueBo.ProjectObjectTypeId, // 这里的值是 TableId，也就是任务栏。
//			RelationCode: consts.IssueRelationBelongManyProCode,
//			RelationType: consts.IssueRelationTypeBelongManyPro,
//			Creator:      issueBo.Creator,
//			CreateTime:   time.Time(issueBo.CreateTime),
//			Updator:      issueBo.Updator,
//			UpdateTime:   time.Time(issueBo.UpdateTime),
//			Version:      issueBo.Version,
//			IsDelete:     issueBo.IsDelete,
//		})
//	}
//	insertErr := mysql.BatchInsert(&po.PpmPriIssueRelation{}, slice.ToSlice(waitRelationBos))
//	if insertErr != nil {
//		log.Error(insertErr)
//		return errs.BuildSystemErrorInfo(errs.MysqlOperateError, insertErr)
//	}
//
//	return nil
//}

func GetTableStatus(orgId int64, tableId int64) ([]status.StatusInfoBo, errs.SystemErrorInfo) {
	tableStatus, err := domain.GetTableStatus(orgId, tableId)
	if err != nil {
		log.Errorf("[GetTableStatus]错误, orgId:%v, tableId:%v, err:%v", orgId, tableId, err)
		return nil, err
	}
	return tableStatus, nil
}
