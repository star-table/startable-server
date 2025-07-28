package api

import (
	"github.com/star-table/startable-server/app/service/projectsvc/service"
	"github.com/star-table/startable-server/common/model/vo"
	"github.com/star-table/startable-server/common/model/vo/projectvo"
)

func (PostGreeter) ProjectChatList(req projectvo.ProjectChatListReqVo) projectvo.ProjectChatListRespVo {
	var pageSize int64 = 10
	if req.Input.PageSize != nil {
		pageSize = *req.Input.PageSize
	}

	pageSize = 500

	res, err := service.ChatList(pageSize, req.UserId, req.Input.ProjectID, req.OrgId, req.Input.LastRelationID, req.Input.Name)
	return projectvo.ProjectChatListRespVo{
		Err:  vo.NewErr(err),
		List: res,
	}
}

func (PostGreeter) UnrelatedChatList(req projectvo.UnrelatedChatListReqVo) projectvo.ProjectChatListRespVo {
	var pageSize int64 = 10
	if req.Input.PageSize != nil {
		pageSize = *req.Input.PageSize
	}

	pageSize = 500

	res, err := service.UnrelatedChatList(pageSize, req.UserId, req.Input.ProjectID, req.OrgId, req.Input.LastOutChatID, req.Input.Name)
	return projectvo.ProjectChatListRespVo{
		Err:  vo.NewErr(err),
		List: res,
	}
}

func (PostGreeter) AddProjectChat(req projectvo.UpdateRelateChatReqVo) vo.CommonRespVo {
	res, err := service.AddChat(req.Input.OutChatIds, req.OrgId, req.UserId, req.Input.ProjectID)
	return vo.CommonRespVo{
		Err:  vo.NewErr(err),
		Void: res,
	}
}

func (PostGreeter) DisbandProjectChat(req projectvo.UpdateRelateChatReqVo) vo.CommonRespVo {
	res, err := service.DisbandChat(req.Input.OutChatIds, req.OrgId, req.UserId, req.Input.ProjectID)
	return vo.CommonRespVo{
		Err:  vo.NewErr(err),
		Void: res,
	}
}

func (PostGreeter) FsChatDisbandCallback(req projectvo.FsChatDisbandCallbackReq) vo.CommonRespVo {
	err := service.FsChatDisbandCallback(req.OrgId, req.ChatId)
	return vo.CommonRespVo{
		Err:  vo.NewErr(err),
		Void: nil,
	}
}

func (PostGreeter) GetProjectMainChatId(req projectvo.GetProjectMainChatIdReq) projectvo.GetProjectMainChatIdResp {
	res, err := service.GetProjectMainChatId(req.OrgId, req.UserId, req.ProjectId, req.SourceChannel)
	return projectvo.GetProjectMainChatIdResp{
		Err:    vo.NewErr(err),
		ChatId: res,
	}
}

func (PostGreeter) CheckIsShowProChatIcon(req projectvo.CheckIsShowProChatIconReq) projectvo.CheckIsShowProChatIconResp {
	isShow, err := service.CheckIsShowProChatIcon(req.OrgId, req.UserId, req.SourceChannel, req.Input)
	return projectvo.CheckIsShowProChatIconResp{
		Err: vo.NewErr(err),
		Data: projectvo.CheckShowProChatIconRespData{
			IsShow: isShow,
		},
	}
}

func (PostGreeter) UpdateFsProjectChatPushSettings(req projectvo.UpdateFsProjectChatPushSettingsReq) vo.CommonRespVo {
	err := service.UpdateFsProjectChatPushSettings(req.OrgId, req.UserId, req.SourceChannel, req.Input)
	return vo.CommonRespVo{
		Err:  vo.NewErr(err),
		Void: &vo.Void{ID: 1},
	}
}

func (PostGreeter) GetFsProjectChatPushSettings(req projectvo.GetFsProjectChatPushSettingsReq) projectvo.GetFsProjectChatPushSettingsResp {
	res, err := service.GetFsProjectChatPushSettings(req.OrgId, &req)
	return projectvo.GetFsProjectChatPushSettingsResp{
		Err:  vo.NewErr(err),
		Data: res,
	}
}

func (PostGreeter) DeleteChatCallback(req projectvo.DeleteChatReq) vo.CommonRespVo {
	err := service.DeleteChatCallback(req.OutOrgId, req.ChatId, req.SourceChannel)
	return vo.CommonRespVo{
		Err:  vo.NewErr(err),
		Void: nil,
	}
}
