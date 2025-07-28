package service

import (
	"context"
	"testing"

	"github.com/star-table/startable-server/app/service/projectsvc/test"
	"github.com/star-table/startable-server/common/core/util/json"
	"github.com/smartystreets/goconvey/convey"
)

func TestPayLimitNum1(t *testing.T) {
	convey.Convey("Test GetProjectRelation", t, test.StartUp(func(ctx context.Context) {
		res, err := PayLimitNum(2627)
		if err != nil {
			t.Error(err)
			return
		}
		t.Log(json.ToJsonIgnoreError(res))
	}))
}
