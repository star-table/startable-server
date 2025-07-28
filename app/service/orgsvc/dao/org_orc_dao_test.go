package orgsvc

import (
	"context"
	"testing"

	"github.com/star-table/startable-server/common/test"
	"github.com/smartystreets/goconvey/convey"
)

func TestPostGreeter_Departments(t *testing.T) {
	convey.Convey("Test sql", t, test.StartUp(func(ctx context.Context) {

		_, _, _ = OrcConfigPageList(1, 100)

	}))
}
