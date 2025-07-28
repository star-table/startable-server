package callsvc

import (
	"github.com/star-table/startable-server/common/core/errs"
	"github.com/star-table/startable-server/common/core/util/json"
)

type FsValidationHandler struct{}

type FsValidationHandlerResp struct {
	Challenge string `json:"challenge"`
}

func (FsValidationHandler) Handle(data string) (string, errs.SystemErrorInfo) {
	validationReq := &ValidationReq{}
	_ = json.FromJson(data, validationReq)

	result := validationReq.Challenge

	return json.ToJsonIgnoreError(FsValidationHandlerResp{
		Challenge: result,
	}), nil
}
