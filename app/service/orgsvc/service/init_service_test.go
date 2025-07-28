package orgsvc

import (
	"context"
	"fmt"
	"testing"

	"gitea.bjx.cloud/allstar/polaris-backend/service/platform/orgsvc/consts"
	"github.com/star-table/startable-server/common/extra/third_platform_sdk"
	"github.com/star-table/startable-server/common/test"

	consts2 "github.com/star-table/startable-server/common/core/consts"
	"github.com/star-table/startable-server/common/core/util/json"
	"github.com/star-table/startable-server/common/library/cache"
	"github.com/star-table/startable-server/common/library/db/mysql"
	"github.com/star-table/startable-server/common/model/bo"
	"github.com/google/martian/log"
	"github.com/smartystreets/goconvey/convey"
	"upper.io/db.v3/lib/sqlbuilder"
)

func TestInitOrg(t *testing.T) {
	convey.Convey("Test login", t, test.StartUp(func(ctx context.Context) {
		third_platform_sdk.InitPlatFormSdk(orgfacade.GetCorpInfoFromDB)

		fmt.Println(InitOrg(bo.InitOrgBo{
			OutOrgId:      "1075dede29da1740",
			OutOrgOwnerId: "ou_e9a41cd346239bb5ef87232348649da8",
			OrgName:       "飞书平台组织",
			SourceChannel: "fs",
		}))
	}))
}

func TestLarkInit(t *testing.T) {
	convey.Convey("Test login", t, test.StartUp(func(ctx context.Context) {

		cacheUserInfo := bo.CacheUserInfoBo{OutUserId: "213213", SourceChannel: "lark_test", UserId: int64(1001), CorpId: "1", OrgId: 222}

		cacheUserInfoJson, _ := json.ToJson(cacheUserInfo)

		log.Info("缓存用户信息" + cacheUserInfoJson)

		cache.Set(consts.CacheUserToken+"abc", cacheUserInfoJson)

		mysql.TransX(func(tx sqlbuilder.Tx) error {
			err := NewbieGuideInit(cacheUserInfo.OrgId, cacheUserInfo.UserId, "fs")
			return err
		})
	}))
}

func TestUpdateOrg(t *testing.T) {
	convey.Convey("Test", t, test.StartUp(func(ctx context.Context) {
		err := mysql.TransX(func(tx sqlbuilder.Tx) error {
			if err := domain.UpdateOrgWithTx(bo.UpdateOrganizationBo{
				Bo: bo.OrganizationBo{
					Id: 1574,
				},
				OrganizationUpdateCond: mysql.Upd{
					consts2.TcVersion: 301001,
				},
			}, tx); err != nil {
				t.Error(err)
				return nil
			}
			return nil
		})
		if err != nil {
			t.Error(err)
			return
		}
	}))
}
