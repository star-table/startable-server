package api

import (
	"github.com/star-table/startable-server/app/service/projectsvc/service"
	"github.com/star-table/startable-server/common/core/util/copyer"
	"github.com/star-table/startable-server/common/model/vo"
	"github.com/star-table/startable-server/common/model/vo/projectvo"
)

func (PostGreeter) PayLimitNum(req projectvo.PayLimitNumReq) projectvo.PayLimitNumResp {
	res, err := service.PayLimitNum(req.OrgId)
	respData := vo.PayLimitNumResp{}
	copyer.Copy(res, &respData)

	return projectvo.PayLimitNumResp{
		Err:  vo.NewErr(err),
		Data: &respData,
	}
}

func (PostGreeter) PayLimitNumForRest(req projectvo.PayLimitNumReq) projectvo.PayLimitNumForRestResp {
	res, err := service.PayLimitNum(req.OrgId)

	return projectvo.PayLimitNumForRestResp{
		Err:  vo.NewErr(err),
		Data: res,
	}
}

func (PostGreeter) GetProjectStatistics(req projectvo.GetProjectStatisticsReqVo) projectvo.GetProjectStatisticsResp {
	res, err := service.GetProjectStatistics(req.OrgId, req.UserId, req.Input)
	return projectvo.GetProjectStatisticsResp{
		Err:  vo.NewErr(err),
		Data: res,
	}
}
