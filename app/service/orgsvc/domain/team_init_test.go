package orgsvc

import (
	"fmt"
	"testing"

	"gitea.bjx.cloud/allstar/polaris-backend/service/platform/orgsvc/po"
	"github.com/star-table/startable-server/common/core/config"
	"github.com/star-table/startable-server/common/core/consts"
	"github.com/star-table/startable-server/common/library/db/mysql"
	"github.com/smartystreets/goconvey/convey"
	"upper.io/db.v3"
)

func TestTeamInit(t *testing.T) {

	config.LoadEnvConfig("F:\\polaris-backend-clone\\config", "application.common", "local")

	conn, _ := mysql.GetConnect()
	tx, _ := conn.NewTx(nil)

	TeamInit(2, tx)
	tx.Commit()

}

func TestGetCorpAuthInfo(t *testing.T) {
	convey.Convey("测试加载env2", t, func() {
		config.LoadEnvConfig("config", "application", "dev")
		info := &[]po.PpmOrgDepartmentOutInfo{}
		mysql.SelectAllByCond(consts.TableDepartmentOutInfo, db.Cond{}, info)
		fmt.Println(info)
	})
}
