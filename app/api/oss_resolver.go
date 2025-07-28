package api

import (
	"context"

	"github.com/star-table/startable-server/app/facade/orgfacade"
	"github.com/star-table/startable-server/app/facade/resourcefacade"
	"github.com/star-table/startable-server/common/core/errs"
	"github.com/star-table/startable-server/common/core/util/json"
	"github.com/star-table/startable-server/common/model/vo"
	"github.com/star-table/startable-server/common/model/vo/resourcevo"
)

func (r *queryResolver) GetOssPostPolicy(ctx context.Context, input vo.OssPostPolicyReq) (*vo.OssPostPolicyResp, error) {

	cacheUserInfo, err := orgfacade.GetCurrentUserRelaxed(ctx)

	if err != nil {
		log.Error(err)
		return nil, errs.BuildSystemErrorInfo(errs.TokenAuthError, err)
	}
	req := resourcevo.GetOssPostPolicyReqVo{
		Input:  input,
		OrgId:  cacheUserInfo.OrgId,
		UserId: cacheUserInfo.UserId,
	}

	respVo := resourcefacade.GetOssPostPolicy(req)
	log.Infof("[GetOssPostPolicy] resp: %v", json.ToJsonIgnoreError(respVo))
	return respVo.GetOssPostPolicy, respVo.Error()
}
