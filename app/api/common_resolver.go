package api

import (
	"context"

	"github.com/star-table/startable-server/app/facade/commonfacade"
	"github.com/star-table/startable-server/app/facade/orgfacade"
	"github.com/star-table/startable-server/common/core/errs"
	"github.com/star-table/startable-server/common/model/vo"
	"github.com/star-table/startable-server/common/model/vo/commonvo"
)

func (r *queryResolver) AreaLinkageList(ctx context.Context, input vo.AreaLinkageListReq) (*vo.AreaLinkageListResp, error) {
	_, err := orgfacade.GetCurrentUserRelaxed(ctx)

	if err != nil {
		log.Error(err)
		return nil, errs.BuildSystemErrorInfo(errs.TokenAuthError, err)
	}

	respVo := commonfacade.AreaLinkageList(commonvo.AreaLinkageListReqVo{
		Input: input,
	})

	if respVo.Failure() {
		log.Error(respVo.Message)
		return nil, respVo.Error()
	}
	return respVo.AreaLinkageListResp, nil

}

func (r *queryResolver) IndustryList(ctx context.Context) (*vo.IndustryListResp, error) {

	//_, err := orgfacade.GetCurrentUserRelaxed(ctx)
	//
	//if err != nil {
	//	log.Error(err)
	//	return nil, errs.BuildSystemErrorInfo(errs.TokenAuthError, err)
	//}

	respVo := commonfacade.IndustryList()

	if respVo.Failure() {
		log.Error(respVo.Message)
		return nil, respVo.Error()
	}
	return respVo.IndustryList, nil
}
