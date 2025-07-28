package commonsvc

import (
	"context"
	"fmt"
	"testing"

	"gitea.bjx.cloud/allstar/polaris-backend/facade/commonfacade"
	"gitea.bjx.cloud/allstar/polaris-backend/facade/orgfacade"
	"github.com/star-table/startable-server/app/service/projectsvc/test"
	"github.com/star-table/startable-server/common/core/config"
	"github.com/star-table/startable-server/common/core/util/json"
	"github.com/smartystreets/goconvey/convey"
)

func TestIndustryList(t *testing.T) {
	convey.Convey("Test 行业列表", t, test.StartUp(func(ctx context.Context) {

		mysqlJson, _ := json.ToJson(config.GetMysqlConfig())
		redisJson, _ := json.ToJson(config.GetRedisConfig())

		fmt.Println("数据库配置:", mysqlJson)
		fmt.Println("redis配置", redisJson)

		orgfacade.GetCurrentUser(ctx)

		resp := commonfacade.IndustryList()

		convey.ShouldBeFalse(resp.Failure())
	}))
}

func TestAA(t *testing.T) {
	a := "a啊"
	t.Log(len(a))
	t.Log(len(a[:5]))
	t.Log(a[:2])
}
