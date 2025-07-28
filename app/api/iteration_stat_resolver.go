package api

import (
	"context"

	"github.com/star-table/startable-server/app/facade/orgfacade"
	"github.com/star-table/startable-server/app/facade/projectfacade"
	"github.com/star-table/startable-server/common/core/errs"
	"github.com/star-table/startable-server/common/model/vo"
	"github.com/star-table/startable-server/common/model/vo/projectvo"
)

func (r *queryResolver) IterationStats(ctx context.Context, page *int, size *int, params vo.IterationStatReq) (*vo.IterationStatList, error) {
	pageA := uint(0)
	sizeA := uint(0)
	if page != nil && size != nil && *page > 0 && *size > 0 {
		pageA = uint(*page)
		sizeA = uint(*size)
	}

	cacheUserInfo, err := orgfacade.GetCurrentUserRelaxed(ctx)
	if err != nil {
		log.Error(err)
		return nil, errs.BuildSystemErrorInfo(errs.TokenAuthError, err)
	}

	respVo := projectfacade.IterationStats(projectvo.IterationStatsReqVo{
		Page:  pageA,
		Size:  sizeA,
		Input: params,
		OrgId: cacheUserInfo.OrgId,
	})
	return respVo.IterationStats, respVo.Error()
}

func (r *mutationResolver) UpdateIterationSort(ctx context.Context, input vo.UpdateIterationSortReq) (*vo.Void, error) {
	cacheUserInfo, err := orgfacade.GetCurrentUserRelaxed(ctx)
	if err != nil {
		log.Error(err)
		return nil, errs.BuildSystemErrorInfo(errs.TokenAuthError, err)
	}

	respVo := projectfacade.UpdateIterationSort(projectvo.UpdateIterationSortReqVo{
		Params: input,
		OrgId:  cacheUserInfo.OrgId,
		UserId: cacheUserInfo.UserId,
	})
	return respVo.Void, respVo.Error()
}
