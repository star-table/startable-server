package ordersvc

import (
	"github.com/star-table/startable-server/common/core/errs"
	"github.com/star-table/startable-server/common/model/vo/ordervo"
)

func GetFunctionByLevel(level int64) ([]ordervo.FunctionLimitObj, errs.SystemErrorInfo) {
	return domain.GetFunctionByLevel(level)
}
