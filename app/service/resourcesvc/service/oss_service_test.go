package resourcesvc

import (
	"context"
	"testing"

	"github.com/star-table/startable-server/app/service/projectsvc/domain"
	"github.com/star-table/startable-server/common/core/consts"
	"github.com/star-table/startable-server/common/core/util/json"
	"github.com/star-table/startable-server/common/test"
	"github.com/smartystreets/goconvey/convey"
)

func TestGetPolicy1(t *testing.T) {
	convey.Convey("Test login", t, test.StartUp(func(ctx context.Context) {
		orgId := int64(2739)
		res, err := domain.GetPayFunctionLimitResourceNum(orgId, consts.FunctionAttachSizeLimit)
		if err != nil {
			t.Error(err)
			return
		}
		t.Log(json.ToJsonIgnoreError(res))
	}))
}
