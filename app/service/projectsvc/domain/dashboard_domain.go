package domain

import (
	"github.com/star-table/startable-server/app/facade/dashboardfacade"
	"github.com/star-table/startable-server/common/core/errs"
	"github.com/star-table/startable-server/common/model/vo/projectvo"
)

func GetAppDashboards(orgId, appId int64) ([]*projectvo.DashboardInfo, errs.SystemErrorInfo) {
	resp := dashboardfacade.GetDashboardList(orgId, 0, []int64{appId})
	if resp.Failure() {
		log.Error(resp.Error())
		return nil, resp.Error()
	}
	return resp.Data, nil
}
