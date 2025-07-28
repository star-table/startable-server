package orgsvc

import (
	"context"
	"testing"

	"github.com/star-table/startable-server/common/model/vo"
	"github.com/star-table/startable-server/common/test"
	"github.com/smartystreets/goconvey/convey"
)

func TestCreateRole(t *testing.T) {
	convey.Convey("Test GetRoleOperationList", t, test.StartUp(func(ctx context.Context) {
		t.Log(CreateRole(10113, 10201, vo.CreateRoleReq{
			RoleGroupType: 1,
			Name:          "测试角色",
		}))
	}))
}
