package orgsvc

import (
	"testing"

	"github.com/star-table/startable-server/common/core/config"
	"github.com/star-table/startable-server/common/core/consts"
	"github.com/star-table/startable-server/common/core/util/strs"
	"github.com/star-table/startable-server/common/library/cache"
	"github.com/star-table/startable-server/common/library/db/mysql"
)

func TestOrgInit(t *testing.T) {
	config.LoadConfig("F:\\workspace-golang-polaris\\polaris-backend\\polaris-server\\configs", "application")

	cache.Set(consts.CacheDingTalkSuiteTicket, "abc")

	conn, _ := mysql.GetConnect()

	tx, _ := conn.NewTx(nil)
	defer func() {
		if tx != nil {
			if err := tx.Close(); err != nil {
				log.Info(strs.ObjectToString(err))
			}
		}
	}()

	_, err := OrgInit("ding8ac2bab2b708b3cc35c2f4657eb6378f", "efg", tx)
	if err != nil {
		err := tx.Rollback()
		if err != nil {
			log.Info(strs.ObjectToString(err))
		}
	} else {
		err := tx.Commit()
		if err != nil {
			log.Info(strs.ObjectToString(err))
		}
	}

}
