package callsvc

import (
	"context"
	"testing"

	"gitea.bjx.cloud/allstar/feishu-sdk-golang/sdk"
	"github.com/star-table/startable-server/common/extra/feishu"
	"github.com/star-table/startable-server/common/test"
	"github.com/smartystreets/goconvey/convey"
)

func GetTenant() *sdk.Tenant {
	// 做大做强：12a813fb180f1758
	sdkClient, err := feishu.GetTenant("12a813fb180f1758")
	if err != nil {
		log.Error(err)
		return nil
	}
	return sdkClient
}

func TestPushErrorNotice1(t *testing.T) {
	convey.Convey("Test", t, test.StartUp(func(ctx context.Context) {
		//sdkClient := GetTenant()
		//PushErrorNotice(*sdkClient, "ou_3ab7fe596cf91692218f744558ae157f", "oc_7bfc74753cff56186e77e0ae9e300c8a",
		//	errs.CommonUserCreateTaskLimit.Message())
	}))
}
