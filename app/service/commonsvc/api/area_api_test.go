package api

import (
	"context"
	"fmt"
	"testing"

	"github.com/smartystreets/goconvey/convey"
	"github.com/star-table/startable-server/app/facade/commonfacade"
	"github.com/star-table/startable-server/app/facade/orgfacade"
	"github.com/star-table/startable-server/common/core/config"
	"github.com/star-table/startable-server/common/core/util/json"
	"github.com/star-table/startable-server/common/model/vo"
	"github.com/star-table/startable-server/common/model/vo/commonvo"
	"github.com/star-table/startable-server/common/test"
)

func TestAreaLinkageList(t *testing.T) {
	convey.Convey("Test 行业列表", t, test.StartUp(func(ctx context.Context) {

		mysqlJson, _ := json.ToJson(config.GetMysqlConfig())
		redisJson, _ := json.ToJson(config.GetRedisConfig())

		fmt.Println("数据库配置:", mysqlJson)
		fmt.Println("redis配置", redisJson)

		orgfacade.GetCurrentUser(ctx)

		isRoot := true

		resp := commonfacade.AreaLinkageList(commonvo.AreaLinkageListReqVo{
			Input: vo.AreaLinkageListReq{
				IsRoot: &isRoot,
			},
		})

		fmt.Printf("地区列表的返回%+v", json.ToJsonIgnoreError(resp))

		convey.ShouldBeFalse(resp.Failure())
	}))
}
