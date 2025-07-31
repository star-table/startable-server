package orgsvc

import (
	"time"

	"github.com/star-table/startable-server/app/facade/trendsfacade"
	"github.com/star-table/startable-server/common/model/bo"
	"github.com/star-table/startable-server/common/model/vo/trendsvo"
)

func PushOrgTrends(orgTrendsBo bo.OrgTrendsBo) {
	orgTrendsBo.OperateTime = time.Now()
	//动态改成同步的
	resp := trendsfacade.AddOrgTrends(trendsvo.AddOrgTrendsReqVo{OrgTrendsBo: orgTrendsBo})
	if resp.Failure() {
		log.Error(resp.Message)
	}
}
