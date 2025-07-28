package feishu

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/star-table/startable-server/common/core/util/json"
	"github.com/star-table/startable-server/common/test"
	"github.com/smartystreets/goconvey/convey"
)

func TestGetDeptUserInfosByDeptIds(t *testing.T) {

	convey.Convey("TestGetDeptUserInfosByDeptIds", t, test.StartUp(func(ctx context.Context) {

		list, err := GetDeptUserInfosByDeptIds("2e99b3ab0b0f1654", []string{"0", "od-f91d6d3208ee337d7a3aea8095dea658"}, true)
		t.Log(err)
		for _, u := range list {
			t.Log(json.ToJsonIgnoreError(u))
		}
	}))

}

func TestGetScopeOpenIds2(t *testing.T) {
	start := time.Now().UnixNano()
	GetScopeOpenIds("2d079cf8960e175d") //2ed263bf32cf1651
	end := time.Now().UnixNano()
	fmt.Println((end - start) / 1e6)
}

func TestGetDeptList(t *testing.T) {
	convey.Convey("TestGetDeptList", t, test.StartUp(func(ctx context.Context) {

		list, err := GetDeptList("2ed263bf32cf1651")
		t.Log(err)
		for _, u := range list {
			t.Log(json.ToJsonIgnoreError(u))
		}
	}))
}

func TestGetScopeDeps(t *testing.T) {
	convey.Convey("TestGetScopeDeps", t, test.StartUp(func(ctx context.Context) {

		list, err := GetScopeDeps("2eb972ddc1cf1652")
		t.Log(err)
		for _, u := range list {
			t.Log(json.ToJsonIgnoreError(u))
		}
	}))
}
