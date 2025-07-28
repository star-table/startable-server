package trendssvc

import (
	"context"
	"testing"

	"github.com/star-table/startable-server/common/core/consts"
	"github.com/star-table/startable-server/common/model/bo"
	"github.com/star-table/startable-server/common/test"
	"github.com/smartystreets/goconvey/convey"
)

func TestAddProjectTrends(t *testing.T) {
	convey.Convey("Test ArchiveProject", t, test.StartUp(func(ctx context.Context) {
		ext := bo.TrendExtensionBo{}
		ext.ObjName = "二号测试应用aa"
		AddProjectTrends(bo.ProjectTrendsBo{
			PushType:              consts.PushTypeUpdateProjectMembers,
			OrgId:                 10101,
			ProjectId:             10113,
			OperatorId:            10201,
			BeforeChangeMembers:   []int64{},
			AfterChangeMembers:    []int64{},
			BeforeChangeFollowers: []int64{},
			AfterChangeFollowers:  []int64{10202},
			Ext:                   ext,
			SourceChannel:         "",
		})
	}))
}
