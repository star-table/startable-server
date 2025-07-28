package orgsvc

import (
	"context"
	"testing"

	sdk_const "gitea.bjx.cloud/allstar/platform-sdk/consts"
	"github.com/star-table/startable-server/app/service/projectsvc/test"
	"github.com/star-table/startable-server/common/core/util/json"
	"github.com/star-table/startable-server/common/core/util/maps"
	"github.com/star-table/startable-server/common/model/vo"
	"github.com/star-table/startable-server/common/model/vo/orgvo"
	"github.com/magiconair/properties/assert"
	"github.com/smartystreets/goconvey/convey"
)

func TestUpdateOrgMemberStatus(t *testing.T) {
	convey.Convey("TestUpdateOrgMemberStatus", t, test.StartUpWithUserInfo(func(userId, orgId int64) {
		_, err := UpdateOrgMemberStatus(orgvo.UpdateOrgMemberStatusReq{
			UserId: 1004,
			OrgId:  orgId,
			Input: vo.UpdateOrgMemberStatusReq{
				MemberIds: []int64{1019},
				Status:    2,
			},
		})
		assert.Equal(t, err, nil)
	}))
}

func TestUpdateOrgMemberCheckStatus(t *testing.T) {
	convey.Convey("TestUpdateOrgMemberStatus", t, test.StartUpWithUserInfo(func(userId, orgId int64) {
		_, err := UpdateOrgMemberCheckStatus(orgvo.UpdateOrgMemberCheckStatusReq{
			UserId: 1004,
			OrgId:  orgId,
			Input: vo.UpdateOrgMemberCheckStatusReq{
				MemberIds:   []int64{1019},
				CheckStatus: 2,
			},
		})
		assert.Equal(t, err, nil)
	}))
}

func TestOrgUserList(t *testing.T) {
	type TestStruct struct {
		Id   int64
		Name string
	}
	list := &[]TestStruct{
		{
			Id:   1,
			Name: "hello",
		},
		{
			Id:   1,
			Name: "world",
		},
		{
			Id:   3,
			Name: "haha",
		},
	}

	cacheMap := maps.NewMap("Id", list)
	t.Log(json.ToJsonIgnoreError(cacheMap))
	t.Log(cacheMap)
	t.Log(cacheMap[int64(1)])
}

func TestGetOrgUserInfoListBySourceChannel(t *testing.T) {
	convey.Convey("TestGetOrgUserInfoListBySourceChannel", t, test.StartUpWithUserInfo(func(userId, orgId int64) {

		resp, err := GetOrgUserInfoListBySourceChannel(orgvo.GetOrgUserInfoListBySourceChannelReq{
			SourceChannel: sdk_const.SourceChannelFeishu,
			Page:          -1,
			Size:          -1,
		})
		t.Log(err)
		assert.Equal(t, err, nil)
		t.Log(json.ToJsonIgnoreError(resp))

		for _, userInfo := range resp.List {
			t.Log(json.ToJsonIgnoreError(userInfo))
		}

		t.Log("===================================")

		resp, err = GetOrgUserInfoListBySourceChannel(orgvo.GetOrgUserInfoListBySourceChannelReq{
			SourceChannel: sdk_const.SourceChannelFeishu,
			Page:          2,
			Size:          10,
		})
		t.Log(err)
		assert.Equal(t, err, nil)
		t.Log(json.ToJsonIgnoreError(resp))

		for _, userInfo := range resp.List {
			t.Log(json.ToJsonIgnoreError(userInfo))
		}

		t.Log("===================================")

		resp, err = GetOrgUserInfoListBySourceChannel(orgvo.GetOrgUserInfoListBySourceChannelReq{
			SourceChannel: sdk_const.SourceChannelFeishu,
			Page:          -1,
			Size:          -1,
		})
		t.Log(err)
		assert.Equal(t, err, nil)
		t.Log(json.ToJsonIgnoreError(resp))

		for _, userInfo := range resp.List {
			t.Log(json.ToJsonIgnoreError(userInfo))
		}
	}))
}

func TestOrgUserList2(t *testing.T) {
	convey.Convey("TestUpdateOrgMemberStatus", t, test.StartUpWithUserInfo(func(userId, orgId int64) {
		info, err := OrgUserList(1111, 1269, 1, 10, nil)
		t.Log(json.ToJsonIgnoreError(info))
		t.Log(err)
	}))
}

func TestGetFunctionConfig(t *testing.T) {
	convey.Convey("TestUpdateOrgMemberStatus", t, test.StartUpWithUserInfo(func(userId, orgId int64) {
		res, err := GetFunctionKeysByOrg(1205)
		t.Log(res)
		t.Log(err)
	}))

}

func TestUserList(t *testing.T) {
	convey.Convey("TestUpdateOrgMemberStatus", t, test.StartUpWithUserInfo(func(userId, orgId int64) {
		res, err := UserList(1004, orgvo.UserListReq{
			SearchCode:   nil,
			IsAllocate:   nil,
			Status:       nil,
			RoleId:       nil,
			DepartmentId: nil,
			Page:         1,
			Size:         2,
		})
		t.Log(json.ToJsonIgnoreError(res))
		t.Log(err)
	}))
}

func TestGetUserListWithCreateTimeRange(t *testing.T) {
	convey.Convey("Test", t, test.StartUp(func(ctx context.Context) {
		resp, err := GetUserListWithCreateTimeRange(&orgvo.GetUserListWithCreateTimeRangeReq{
			SourceChannel: "fs",
			CreateTime1:   "2020-12-22 00:00:00",
			CreateTime2:   "2020-12-24 23:59:59",
			Page:          1,
			Size:          10,
		})
		if err != nil {
			t.Error(err)
		}
		t.Log(json.ToJsonIgnoreError(resp))
	}))
}
