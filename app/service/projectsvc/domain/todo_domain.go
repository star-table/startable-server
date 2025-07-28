package domain

import (
	automationPb "gitea.bjx.cloud/LessCode/interface/golang/automation/v1"
	"github.com/star-table/startable-server/app/facade/automationfacade"
	"github.com/star-table/startable-server/app/facade/orgfacade"
	"github.com/star-table/startable-server/common/core/consts"
	"github.com/star-table/startable-server/common/core/errs"
	"github.com/star-table/startable-server/common/core/util/slice"
	"github.com/star-table/startable-server/common/extra/card"
	"github.com/star-table/startable-server/common/model/bo"
	"github.com/star-table/startable-server/common/model/vo/orgvo"
	"github.com/star-table/startable-server/common/model/vo/projectvo"
)

func TodoUrge(orgId, userId, todoId int64, msg string) errs.SystemErrorInfo {
	log.Infof("[TodoUrge] org:%v, user:%v, todo:%v", orgId, userId, todoId)

	// 获取待办信息
	todoResp := automationfacade.GetTodo(orgId, userId, todoId)
	if todoResp.Failure() || todoResp.Data.Todo == nil {
		log.Errorf("[TodoUrge] GetTodo, org:%d user:%d todo:%d, err: %v", orgId, userId, todoId, todoResp.Error())
		return todoResp.Error()
	}
	todo := todoResp.Data.Todo

	// check
	if todo.Status != automationPb.TodoStatus_SUnFinished {
		log.Infof("[TodoUrge] todo cannot be urged by status: %v, ignored. org:%d user:%d todo:%d.",
			todo.Status, orgId, userId, todoId)
		return nil
	}
	if todo.TriggerUserId != userId {
		log.Errorf("[TodoUrge] user is not trigger:%v. org:%d user:%d todo:%d.",
			todo.TriggerUserId, orgId, userId, todoId)
		return errs.TodoInvalidOperator
	}
	if !todo.AllowUrgeByTrigger {
		log.Errorf("[TodoUrge] urge by trigger is not allowed. org:%d user:%d todo:%d.", orgId, userId, todoId)
		return errs.TodoInvalidOp
	}

	// 获取 org base info
	orgResp := orgfacade.GetBaseOrgInfo(orgvo.GetBaseOrgInfoReqVo{OrgId: orgId})
	if orgResp.Failure() {
		log.Errorf("[TodoUrge] GetBaseOrgInfo failed. org:%d user:%d, err: %v", orgId, userId, orgResp.Error())
		return orgResp.Error()
	}
	org := orgResp.BaseOrgInfo

	if org.OutOrgId == "" {
		log.Infof("[TodoUrge] local org, ignored. org:%d user:%d", orgId, userId)
		return nil
	}

	// 获取任务信息
	issue, errSys := GetIssueInfoLc(orgId, userId, todo.IssueId)
	if errSys != nil {
		log.Errorf("[TodoUrge] GetIssueInfoLc failed. org:%d user:%d issue:%d, err: %v", orgId, userId, todo.IssueId, errSys)
		return errSys
	}

	// 获取项目信息
	project, errSys := GetProject(orgId, issue.ProjectId)
	if errSys != nil {
		log.Errorf("[TodoUrge] GetProject failed. org:%d project:%d, err: %v", orgId, issue.ProjectId, errSys)
		return errSys
	}

	// 获取表格信息
	tableMeta, errSys := GetTableByTableId(orgId, userId, issue.TableId)
	if errSys != nil {
		log.Errorf("[TodoUrge] GetTableByTableId failed. org:%d project:%d, err: %v", orgId, userId, issue.TableId, errSys)
		return errSys
	}

	var userIds []int64
	var urgeUserIds []int64
	userIds = append(userIds, userId)
	for id, op := range todo.Operators {
		if op.Op == automationPb.TodoOp_OpInit {
			userIds = append(userIds, id)
			urgeUserIds = append(urgeUserIds, id)
		}
	}
	for _, id := range issue.OwnerIdI64 {
		userIds = append(userIds, id)
	}
	userIds = slice.SliceUniqueInt64(userIds)
	userResp := orgfacade.GetBaseUserInfoBatch(orgvo.GetBaseUserInfoBatchReqVo{
		OrgId:   orgId,
		UserIds: userIds,
	})
	if userResp.Failure() {
		log.Errorf("[TodoUrge] GetBaseUserInfoBatch failed. org:%d users:%v, err: %v", orgId, userIds, errSys)
		return orgResp.Error()
	}
	users := make(map[int64]*bo.BaseUserInfoBo)
	for i, info := range userResp.BaseUserInfos {
		users[info.UserId] = &userResp.BaseUserInfos[i]
	}

	// selfUser
	selfUser, ok := users[userId]
	if !ok {
		log.Errorf("[TodoUrge] user not exists. org:%d user:%v", orgId, userId)
		return errs.UserNotExist
	}

	// toOpenIds
	var toOpenIds []string
	for _, id := range urgeUserIds {
		if id != userId { // 自己不urge
			if u, ok := users[id]; ok {
				toOpenIds = append(toOpenIds, u.OutUserId)
			}
		}
	}

	// owners
	var owners []*bo.BaseUserInfoBo
	for _, id := range issue.OwnerIdI64 {
		if u, ok := users[id]; ok {
			owners = append(owners, u)
		}
	}

	// table columns
	tableColumns, errSys := GetTableColumnsMap(orgId, issue.TableId, []string{consts.BasicFieldTitle, consts.BasicFieldOwnerId})
	if errSys != nil {
		log.Errorf("[TodoUrge] GetTableColumnConfig failed. org:%v, table:%v, err: %v",
			orgId, issue.TableId, errSys)
		return errSys
	}

	urgeIssueCard := &projectvo.UrgeIssueCard{
		OrgId:           orgId,
		OutOrgId:        org.OutOrgId,
		OperateUserId:   userId,
		OperateUserName: selfUser.Name,
		OpenIds:         toOpenIds,
		ProjectName:     project.Name,
		TableName:       tableMeta.Name,
		ProjectId:       issue.ProjectId,
		ProjectTypeId:   project.ProjectTypeId,
		IssueId:         issue.Id,
		IssueBo:         issue,
		OwnerInfos:      owners,
		UrgeText:        msg,
		IssueLinks:      GetIssueLinks(org.SourceChannel, orgId, issue.Id),
		TableColumn:     tableColumns,
	}
	errSys = card.SendCardTodoUrge(org.SourceChannel, urgeIssueCard)
	if errSys != nil {
		log.Errorf("[TodoUrge] SendCardTodoUrge failed. org:%d, user:%d, issue:%d", orgId, userId, issue.Id, errSys)
		return errSys
	}
	return nil
}
