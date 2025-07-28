package ordersvc

import (
	"github.com/star-table/startable-server/app/service"
	"github.com/star-table/startable-server/common/model/vo"
	"github.com/star-table/startable-server/common/model/vo/ordervo"
)

func (PostGreeter) AddFsOrder(req ordervo.AddFsOrderReq) vo.CommonRespVo {
	res, err := service.AddFsOrder(req.Data)
	return vo.CommonRespVo{
		Err:  vo.NewErr(err),
		Void: &vo.Void{ID: res},
	}
}

func (PostGreeter) AddDingOrder(req ordervo.AddDingOrderReq) vo.CommonRespVo {
	res, err := service.AddDingOrder(req.Data)
	return vo.CommonRespVo{
		Err:  vo.NewErr(err),
		Void: &vo.Void{ID: res},
	}
}

func (PostGreeter) AddWeiXinOrder(req ordervo.AddWeiXinOrderReq) vo.CommonRespVo {
	res, err := service.AddWeiXinOrder(req.Data)
	return vo.CommonRespVo{
		Err:  vo.NewErr(err),
		Void: &vo.Void{ID: res},
	}
}

func (PostGreeter) UpdateDingOrder(req ordervo.UpdateDingOrderReq) ordervo.UpdateDingOrderDataResp {
	res, err := service.UpdateOrderInfo(req.Input.OrgId, req.Input.OutOrgId)
	return ordervo.UpdateDingOrderDataResp{
		Err:  vo.NewErr(err),
		Data: res,
	}
}

func (PostGreeter) DeleteDingOrder(req ordervo.DeleteDingOrderReqVo) vo.CommonRespVo {
	res, err := service.DeleteDingOrderByOrderId(req.Data.OutOrgId, req.Data.OrderId)
	return vo.CommonRespVo{
		Err:  vo.NewErr(err),
		Void: &vo.Void{ID: res},
	}
}

// CheckFsOrgIsInTrial 检查一个 fs 组织目前是否是处于试用期间
// 如果不是 fs 组织，查询不到试用订单，则视为不在试用时间内。
func (PostGreeter) CheckFsOrgIsInTrial(req ordervo.CheckFsOrgIsInTrialReq) vo.BoolRespVo {
	isTrial, err := service.CheckFsOrgIsInTrial(req.OrgId)
	return vo.BoolRespVo{
		Err:    vo.NewErr(err),
		IsTrue: isTrial,
	}
}

func (PostGreeter) CreateWeiXinLicenceOrder(req ordervo.CreateWeiXinLicenceOrderReq) vo.CommonRespVo {
	err := service.CreateLicenceOrder(&req)
	return vo.CommonRespVo{
		Err:  vo.NewErr(err),
		Void: nil,
	}
}

func (PostGreeter) GetOrderPayInfo(req ordervo.GetOrderPayInfoReq) ordervo.GetOrderPayInfoResp {
	data, err := service.GetOrderPayInfo(req.OrgId)
	return ordervo.GetOrderPayInfoResp{
		Err:  vo.NewErr(err),
		Data: data,
	}
}
