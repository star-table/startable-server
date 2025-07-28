package ordersvc

import (
	"github.com/star-table/startable-server/app/service"
	"github.com/star-table/startable-server/common/model/vo"
	"github.com/star-table/startable-server/common/model/vo/ordervo"
)

func (PostGreeter) GetFunctionByLevel(req ordervo.FunctionReq) ordervo.FunctionResp {
	res, err := service.GetFunctionByLevel(req.Level)
	return ordervo.FunctionResp{
		Err:  vo.NewErr(err),
		Data: res,
	}
}
