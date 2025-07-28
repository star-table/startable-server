package api

import (
	"context"

	sdk_const "gitea.bjx.cloud/allstar/platform-sdk/consts"
	"github.com/star-table/startable-server/app/facade/orgfacade"
	"github.com/star-table/startable-server/app/facade/projectfacade"
	"github.com/star-table/startable-server/common/core/errs"
	"github.com/star-table/startable-server/common/model/vo"
	"github.com/star-table/startable-server/common/model/vo/projectvo"
)

func (r *queryResolver) GetProjectMainChatID(ctx context.Context, params vo.GetProjectMainChatIDReq) (*vo.GetProjectMainChatIDResp, error) {
	cacheUserInfo, err := orgfacade.GetCurrentUserRelaxed(ctx)
	if err != nil {
		log.Error(err)
		return nil, errs.BuildSystemErrorInfo(errs.TokenAuthError, err)
	}
	respVo := projectfacade.GetProjectMainChatId(projectvo.GetProjectMainChatIdReq{
		OrgId:         cacheUserInfo.OrgId,
		UserId:        cacheUserInfo.UserId,
		SourceChannel: cacheUserInfo.SourceChannel,
		ProjectId:     params.ProjectID,
	})

	return &vo.GetProjectMainChatIDResp{ChatID: respVo.ChatId}, respVo.Error()
}

func (r *queryResolver) GetFsProjectChatPushSettings(ctx context.Context, params vo.GetFsProjectChatPushSettingsReq) (*vo.GetFsProjectChatPushSettingsResp, error) {
	cacheUserInfo, err := orgfacade.GetCurrentUserRelaxed(ctx)
	if err != nil {
		log.Error(err)
		return nil, errs.BuildSystemErrorInfo(errs.TokenAuthError, err)
	}
	if cacheUserInfo.SourceChannel != sdk_const.SourceChannelFeishu {
		return nil, errs.CannotBindChat
	}
	respVo := projectfacade.GetFsProjectChatPushSettings(projectvo.GetFsProjectChatPushSettingsReq{
		OrgId:         cacheUserInfo.OrgId,
		ChatId:        params.ChatID,
		ProjectId:     params.ProjectID,
		SourceChannel: cacheUserInfo.SourceChannel,
	})

	return respVo.Data, respVo.Error()
}

func (r *mutationResolver) UpdateFsProjectChatPushSettings(ctx context.Context, params vo.UpdateFsProjectChatPushSettingsReq) (*vo.Void, error) {
	cacheUserInfo, err := orgfacade.GetCurrentUserRelaxed(ctx)
	if err != nil {
		log.Error(err)
		return nil, errs.BuildSystemErrorInfo(errs.TokenAuthError, err)
	}
	if cacheUserInfo.SourceChannel != sdk_const.SourceChannelFeishu {
		return nil, errs.CannotBindChat
	}
	respVo := projectfacade.UpdateFsProjectChatPushSettings(projectvo.UpdateFsProjectChatPushSettingsReq{
		OrgId:         cacheUserInfo.OrgId,
		UserId:        cacheUserInfo.UserId,
		SourceChannel: cacheUserInfo.SourceChannel,
		Input:         params,
	})

	return respVo.Void, respVo.Error()
}
