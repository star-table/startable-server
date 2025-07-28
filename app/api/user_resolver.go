package api

import (
	"context"

	"github.com/star-table/startable-server/app/facade/orgfacade"
	"github.com/star-table/startable-server/common/core/errs"
	"github.com/star-table/startable-server/common/model/vo"
	"github.com/star-table/startable-server/common/model/vo/orgvo"
)

func (r *queryResolver) PersonalInfo(ctx context.Context) (*vo.PersonalInfo, error) {
	cacheUserInfo, err := orgfacade.GetCurrentUserRelaxed(ctx)

	if err != nil {
		log.Error(err)
		return nil, errs.BuildSystemErrorInfo(errs.TokenAuthError, err)
	}

	orgId := cacheUserInfo.OrgId
	userId := cacheUserInfo.UserId

	req := orgvo.PersonalInfoReqVo{
		OrgId:         orgId,
		UserId:        userId,
		SourceChannel: cacheUserInfo.SourceChannel,
	}

	resp := orgfacade.PersonalInfo(req)

	return resp.PersonalInfo, resp.Error()
}

func (r *queryResolver) UserConfigInfo(ctx context.Context) (*vo.UserConfig, error) {
	cacheUserInfo := orgfacade.GetCurrentUserWithoutPayVerify(ctx)

	if cacheUserInfo.Failure() {
		log.Error(cacheUserInfo.Err.Error())
		return nil, errs.BuildSystemErrorInfo(errs.TokenAuthError, cacheUserInfo.Err.Error())
	}

	req := orgvo.UserConfigInfoReqVo{
		UserId: cacheUserInfo.CacheInfo.UserId,
		OrgId:  cacheUserInfo.CacheInfo.OrgId,
	}

	resp := orgfacade.UserConfigInfo(req)
	if resp.Failure() {
		return nil, resp.Error()
	}

	return resp.UserConfigInfo, nil
}

func (r *mutationResolver) UpdateUserConfig(ctx context.Context, input vo.UpdateUserConfigReq) (*vo.UpdateUserConfigResp, error) {
	cacheUserInfo := orgfacade.GetCurrentUserWithoutPayVerify(ctx)

	if cacheUserInfo.Failure() {
		log.Error(cacheUserInfo.Err.Error())
		return nil, errs.BuildSystemErrorInfo(errs.TokenAuthError, cacheUserInfo.Err.Error())
	}

	resp := orgfacade.UpdateUserConfig(orgvo.UpdateUserConfigReqVo{
		UpdateUserConfigReq: input,
		UserId:              cacheUserInfo.CacheInfo.UserId,
		OrgId:               cacheUserInfo.CacheInfo.OrgId,
	})
	if resp.Failure() {
		return nil, resp.Error()
	}

	return resp.UpdateUserConfig, nil
}

func (r *mutationResolver) UpdateUserPcConfig(ctx context.Context, input vo.UpdateUserPcConfigReq) (*vo.UpdateUserConfigResp, error) {
	cacheUserInfo := orgfacade.GetCurrentUserWithoutPayVerify(ctx)

	if cacheUserInfo.Failure() {
		log.Error(cacheUserInfo.Err.Error())
		return nil, errs.BuildSystemErrorInfo(errs.TokenAuthError, cacheUserInfo.Err.Error())
	}

	resp := orgfacade.UpdateUserPcConfig(orgvo.UpdateUserPcConfigReqVo{
		UpdateUserPcConfigReq: input,
		UserId:                cacheUserInfo.CacheInfo.UserId,
		OrgId:                 cacheUserInfo.CacheInfo.OrgId,
	})
	if resp.Failure() {
		return nil, resp.Error()
	}
	return resp.UpdateUserConfig, nil
}

func (r *mutationResolver) UpdateUserInfo(ctx context.Context, input vo.UpdateUserInfoReq) (*vo.Void, error) {
	cacheUserInfo, err := orgfacade.GetCurrentUserRelaxed(ctx)

	if err != nil {
		log.Error(err)
		return nil, errs.BuildSystemErrorInfo(errs.TokenAuthError, err)
	}

	resp := orgfacade.UpdateUserInfo(orgvo.UpdateUserInfoReqVo{
		UpdateUserInfoReq: input,
		OrgId:             cacheUserInfo.OrgId,
		UserId:            cacheUserInfo.UserId,
	})

	if resp.Failure() {
		return nil, resp.Error()
	}

	return resp.Void, nil

}

func (r *queryResolver) GetPayRemind(ctx context.Context) (*vo.GetPayRemindResp, error) {
	cacheUserInfo, err := orgfacade.GetCurrentUserRelaxed(ctx)
	if err != nil {
		log.Error(err)
		return nil, errs.BuildSystemErrorInfo(errs.TokenAuthError, err)
	}
	resp := orgfacade.GetPayRemind(orgvo.GetPayRemindReq{
		OrgId:  cacheUserInfo.OrgId,
		UserId: cacheUserInfo.UserId,
	})

	if resp.Failure() {
		return nil, resp.Error()
	}
	return resp.Data, nil
}
