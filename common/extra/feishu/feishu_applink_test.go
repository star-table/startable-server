package feishu

import (
	"context"
	"testing"

	"github.com/star-table/startable-server/common/test"
	"github.com/smartystreets/goconvey/convey"
)

func TestGetProjectWelcomeUrl(t *testing.T) {
	convey.Convey("TestGetProjectWelcomeUrl", t, test.StartUp(func(ctx context.Context) {
		t.Log(GetAppLinkProjectUrl(123, 1, 504))
		t.Log(GetProjectWelcomeUrl())
	}))
}
