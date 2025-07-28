package service

import (
	"context"
	"testing"

	"github.com/star-table/startable-server/app/service/projectsvc/test"
	"github.com/star-table/startable-server/common/core/util/json"
	"github.com/star-table/startable-server/common/model/vo/projectvo"
	"github.com/smartystreets/goconvey/convey"
	"gopkg.in/go-playground/assert.v1"
)

func TestIssueAndProjectCountStat(t *testing.T) {

	convey.Convey("TestIssueAndProjectCountStat", t, test.StartUpWithUserInfo(func(userId, orgId int64) {
		resp, err := IssueAndProjectCountStat(projectvo.IssueAndProjectCountStatReqVo{
			OrgId:  orgId,
			UserId: userId,
		})

		log.Error(err)
		log.Error(json.ToJsonIgnoreError(resp))
		assert.Equal(t, err, nil)
		assert.NotEqual(t, resp, nil)
		assert.Equal(t, resp.ProjectNotCompletedCount, int64(1))
		assert.Equal(t, resp.IssueNotCompletedCount, int64(1))
	}))

}

func TestIssueDailyPersonalWorkCompletionStat(t *testing.T) {
	convey.Convey("TestIssueDailyPersonalWorkCompletionStat", t, test.StartUpWithUserInfo(func(userId, orgId int64) {
		resp, err := IssueDailyPersonalWorkCompletionStat(projectvo.IssueDailyPersonalWorkCompletionStatReqVo{
			OrgId:  orgId,
			UserId: userId,
		})

		log.Error(err)
		log.Error(json.ToJsonIgnoreError(resp))
		assert.Equal(t, err, nil)
		assert.NotEqual(t, resp, nil)
	}))
}

func TestIssueAssignRank1(t *testing.T) {
	convey.Convey("TestIssueDailyPersonalWorkCompletionStat", t, test.StartUp(func(ctx context.Context) {
		res, err := IssueAssignRank(8055, 42839, 10)
		if err != nil {
			t.Error(err)
			return
		}
		t.Log(json.ToJsonIgnoreError(res))
	}))
}
