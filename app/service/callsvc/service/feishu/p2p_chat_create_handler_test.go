package callsvc

import (
	"context"
	"testing"
	"time"

	"github.com/star-table/startable-server/common/test"
	"github.com/smartystreets/goconvey/convey"
)

func TestPushBotHelpNotice(t *testing.T) {
	convey.Convey("Test TestPushBotHelpNotice", t, test.StartUp(func(ctx context.Context) {
		//PushBotHelpNotice("2ed263bf32cf1651", "oc_c6c692ff9b21486009f64222ede7003f", "ou_433f2b1abf3e0ec316fd9e60d4cda654")
		PushBotHelpNotice("2ed263bf32cf1651", "", "ou_87f1b2210acad10a90cc3690802626d7")
	}))
}

func TestPushBotInstallSuccessNotice(t *testing.T) {
	convey.Convey("Test TestPushBotInstallSuccessNotice", t, test.StartUp(func(ctx context.Context) {
		PushBotInstallSuccessNotice("2ed263bf32cf1651", "oc_c6c692ff9b21486009f64222ede7003f", "")
	}))
}

func TestTime(t *testing.T) {
	t.Log(time.Unix(1594290392, 0))
}
