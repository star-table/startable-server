package domain

import (
	"context"
	"fmt"
	"testing"

	"github.com/star-table/startable-server/app/service/projectsvc/test"
	"github.com/smartystreets/goconvey/convey"
)

func TestAddFsChatMembers(t *testing.T) {
	convey.Convey("Test UpdateIssue", t, test.StartUp(func(ctx context.Context) {
		t.Log(AddFsChatMembers(1242, 7984, []int64{}, []int64{}, []int64{0}, []int64{}))
	}))
}

func TestGetProTableChatSettingBoArrByTableId(t *testing.T) {
	convey.Convey("Test UpdateIssue", t, test.StartUp(func(ctx context.Context) {
		list, err := GetProTableChatSettingBoArrByTableId(2373, 13882, 1535105783056830464)
		fmt.Println(list, err)
	}))
}
