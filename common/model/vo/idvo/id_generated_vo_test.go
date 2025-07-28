package idvo

import (
	"testing"

	"github.com/star-table/startable-server/common/core/errs"
	"github.com/star-table/startable-server/common/core/util/json"
	"github.com/star-table/startable-server/common/model/vo"
)

func Test_Convert(t *testing.T) {
	obj := ApplyPrimaryIdRespVo{
		Id:  123,
		Err: vo.NewErr(errs.SystemError),
	}
	t.Log(json.ToJsonIgnoreError(obj))

	obj1 := vo.VoidErr{
		Err: vo.NewErr(errs.SystemError),
	}
	t.Log(json.ToJsonIgnoreError(obj1))
}
