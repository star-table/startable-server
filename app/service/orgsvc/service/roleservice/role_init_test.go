package orgsvc

import (
	"context"
	"fmt"
	"testing"

	domain "github.com/star-table/startable-server/app/service"
	"github.com/star-table/startable-server/common/core/config"
	"github.com/star-table/startable-server/common/library/db/mysql"
	"github.com/star-table/startable-server/common/test"
	"github.com/smartystreets/goconvey/convey"
)

func TestRoleInit(t *testing.T) {
	config.LoadConfig("F:\\polaris-backend\\polaris-server\\configs", "application")

	convey.Convey("Test GetRoleOperationList", t, test.StartUp(func(ctx context.Context) {
		convey.Convey("权限init", func() {
			conn, _ := mysql.GetConnect()
			tx, _ := conn.NewTx(nil)
			_, err := domain.RoleInit(1, tx)
			if err != nil {
				tx.Rollback()
				fmt.Println(err)
			}
			tx.Commit()
		})
	}))
}

func TestChangeDefaultRole(t *testing.T) {
	convey.Convey("Test GetRoleOperationList", t, test.StartUp(func(ctx context.Context) {
		t.Log(getDataIdByIssueId(2373, 0, 10236478))
	}))
}
