package service

import (
	"github.com/star-table/startable-server/app/facade/automationfacade"
	"github.com/star-table/startable-server/app/facade/orgfacade"
	"github.com/star-table/startable-server/app/service/projectsvc/domain"
	"github.com/star-table/startable-server/common/core/errs"
	"github.com/star-table/startable-server/common/extra/card"
	"github.com/star-table/startable-server/common/model/bo"
	"github.com/star-table/startable-server/common/model/vo/commonvo"
	"github.com/star-table/startable-server/common/model/vo/orgvo"
	"github.com/star-table/startable-server/common/model/vo/projectvo"
)

func CreateTodoHook(req *projectvo.InnerCreateTodoHookReq) errs.SystemErrorInfo {
	log.Infof("[CreateTodoHook] org: %v, user: %v, todo: %v", req.OrgId, req.UserId, req.Input.TodoId)

	// todo
	resp := automationfacade.GetTodo(req.OrgId, req.UserId, req.Input.TodoId)
	if resp.Failure() {
		log.Errorf("[CreateTodoHook] get todo failed, error: %v", resp.Error())
		return resp.Error()
	}
	todo := resp.Data.Todo

	// org
	orgBaseInfo, errSys := orgfacade.GetBaseOrgInfoRelaxed(req.OrgId)
	if errSys != nil {
		log.Errorf("[CreateTodoHook] GetBaseOrgInfoRelaxed failed, orgId:%v, err: %v", req.OrgId, errSys)
		return errSys
	}

	// card
	cd, err := domain.GetTodoCard(req.OrgId, req.UserId, orgBaseInfo.SourceChannel, todo)
	if err != nil {
		log.Errorf("[CreateTodoHook] get todo card failed, error: %v", err)
		return err
	}

	var userIds []int64
	for id, _ := range todo.Operators {
		userIds = append(userIds, id)
	}

	if len(userIds) == 0 {
		log.Error("[CreateTodoHook] no operator need to notice")
		return nil
	}

	// 查询用户信息
	userResp := orgfacade.GetBaseUserInfoBatch(orgvo.GetBaseUserInfoBatchReqVo{
		OrgId:   req.OrgId,
		UserIds: userIds,
	})
	if userResp.Failure() || len(userResp.BaseUserInfos) == 0 {
		log.Errorf("[CreateTodoHook] GetBaseUserInfoBatch failed, err: %v", userResp.Error())
		return userResp.Error()
	}
	users := map[int64]*bo.BaseUserInfoBo{}
	for i, info := range userResp.BaseUserInfos {
		users[info.UserId] = &userResp.BaseUserInfos[i]
	}

	// 机器人将信息回复给用户
	var allOutUserIds []string
	for _, id := range userIds {
		if info, ok := users[id]; ok {
			if info.OutUserId != "" {
				allOutUserIds = append(allOutUserIds, info.OutUserId)
			}
		}
	}
	if len(allOutUserIds) == 0 {
		log.Error("[CreateTodoHook] no out users")
		return nil
	}

	cardMsg := &commonvo.PushCard{
		OrgId:         req.OrgId,
		OutOrgId:      orgBaseInfo.OutOrgId,
		SourceChannel: orgBaseInfo.SourceChannel,
		OpenIds:       allOutUserIds,
		CardMsg:       cd,
	}
	errSys = card.PushCard(req.OrgId, cardMsg)
	if errSys != nil {
		log.Errorf("[CreateTodoHook] PushCard err:%v", errSys)
		return errSys
	}

	return nil
}

func TodoUrge(req *projectvo.TodoUrgeReq) errs.SystemErrorInfo {
	return domain.TodoUrge(req.OrgId, req.UserId, req.Input.TodoId, req.Input.Msg)
}
