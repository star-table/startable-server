package service

import (
	"context"
	"testing"

	"github.com/star-table/startable-server/app/service/projectsvc/test"
	"github.com/star-table/startable-server/common/core/util/json"
	"github.com/star-table/startable-server/common/model/vo"
	"github.com/smartystreets/goconvey/convey"
)

func TestProjectIssueRelatedStatus(t *testing.T) {
	convey.Convey("Test GetProjectRelation", t, test.StartUp(func(ctx context.Context) {
		res, err := ProjectIssueRelatedStatus(1013, vo.ProjectIssueRelatedStatusReq{ProjectID: 1319, TableID: "1479"})
		t.Log(json.ToJsonIgnoreError(res))
		t.Log(err)
	}))
}
