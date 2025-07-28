package ordersvc

import (
	"context"
	"testing"

	"github.com/star-table/startable-server/common/test"
	"github.com/smartystreets/goconvey/convey"
)

func TestFixOrderOfZeroOrgId1(t *testing.T) {
	convey.Convey("Test", t, test.StartUp(func(ctx context.Context) {
		err := FixOrderOfZeroOrgId()
		if err != nil {
			t.Error(err)
			return
		}
		t.Log("--end--")
	}))
}
