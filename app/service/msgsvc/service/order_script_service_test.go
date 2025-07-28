package msgsvc

import (
	"context"
	"testing"

	"github.com/star-table/startable-server/common/model/vo/msgvo"
	"github.com/star-table/startable-server/common/test"

	"github.com/smartystreets/goconvey/convey"
)

func TestFixAddOrderForFeiShu(t *testing.T) {
	convey.Convey("Test", t, test.StartUp(func(ctx context.Context) {
		if err := FixAddOrderForFeiShu(msgvo.FixAddOrderForFeiShuReqData{
			StartTime: nil,
			EndTime:   nil,
		}); err != nil {
			t.Error(err)
			return
		}
		t.Log("...end...")
	}))
}
