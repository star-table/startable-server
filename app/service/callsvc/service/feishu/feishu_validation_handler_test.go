package callsvc

import (
	"context"
	"testing"

	"github.com/star-table/startable-server/common/core/util/json"
	"github.com/star-table/startable-server/common/test"
	"github.com/smartystreets/goconvey/convey"
)

func TestGetIssues(t *testing.T) {
	convey.Convey("Test", t, test.StartUp(func(ctx context.Context) {
		issues, err := domain.GetIssueByChatId(2373, "oc_8e618d2b4177b361a566147a351de52b", false)
		if err != nil {
			t.Error(err)
			return
		}
		t.Log(json.ToJsonIgnoreError(issues))
	}))
}
