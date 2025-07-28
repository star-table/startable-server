package api

import (
	"context"

	"github.com/star-table/startable-server/app/facade/orgfacade"
	"github.com/star-table/startable-server/common/core/errs"
	"github.com/star-table/startable-server/common/model/vo"
	"github.com/star-table/startable-server/common/model/vo/orgvo"
)

func (r *mutationResolver) JoinOrgByInviteCode(ctx context.Context, params vo.JoinOrgByInviteCodeReq) (*vo.Void, error) {
	cacheUserInfo, err := orgfacade.GetCurrentUserRelaxed(ctx)

	if err != nil {
		log.Error(err)
		return nil, errs.BuildSystemErrorInfo(errs.TokenAuthError, err)
	}

	respVo := orgfacade.JoinOrgByInviteCode(orgvo.JoinOrgByInviteCodeReq{
		UserId:     cacheUserInfo.UserId,
		OrgId:      cacheUserInfo.OrgId,
		InviteCode: params.InviteCode,
	})

	if respVo.Failure() {
		return nil, respVo.Error()
	}
	return &vo.Void{
		ID: cacheUserInfo.UserId,
	}, nil
}
