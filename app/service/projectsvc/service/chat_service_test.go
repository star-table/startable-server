package service

import (
	"context"
	"testing"

	sdk_const "gitea.bjx.cloud/allstar/platform-sdk/consts"
	"github.com/star-table/startable-server/app/facade/orgfacade"
	"github.com/star-table/startable-server/app/service/projectsvc/domain"
	"github.com/star-table/startable-server/app/service/projectsvc/test"
	"github.com/star-table/startable-server/common/core/util/json"
	"github.com/star-table/startable-server/common/core/util/str"
	"github.com/star-table/startable-server/common/model/vo"
	"github.com/star-table/startable-server/common/model/vo/orgvo"
	"github.com/star-table/startable-server/common/model/vo/projectvo"
	"github.com/smartystreets/goconvey/convey"
	"gotest.tools/assert"
)

func TestAddFsChat(t *testing.T) {
	convey.Convey("Test login", t, test.StartUp(func(ctx context.Context) {
		des := "aa"
		t.Log(AddFsChat(1242, 1883, []int64{}, 10125, "aabb", &des, []int64{1883}, []int64{0}))
	}))
}

func TestAllChatList(t *testing.T) {
	convey.Convey("Test login", t, test.StartUp(func(ctx context.Context) {
		//name := "刘"
		t.Log(AllChatList(1284, 1111, nil))
	}))
}

func TestChatList(t *testing.T) {
	convey.Convey("Test login", t, test.StartUp(func(ctx context.Context) {
		res, err := ChatList(10, 1007, 10125, 1003, nil, nil)
		t.Log(json.ToJsonIgnoreError(res))
		assert.Equal(t, err, nil)
	}))
}

func TestDisbandChat(t *testing.T) {
	convey.Convey("Test login", t, test.StartUp(func(ctx context.Context) {
		t.Log(DisbandChat([]string{"oc_c09294289150324643033f7cdd74f0db", "bbbb"}, 1003, 1007, 10125))
	}))
}

func TestAddChat(t *testing.T) {
	convey.Convey("Test login", t, test.StartUp(func(ctx context.Context) {
		//syncMember := false
		t.Log(AddChat([]string{"oc_c09294289150324643033f7cdd74f0db", "aaa"}, 1003, 1007, 10125))
	}))
}

func TestGetFsProjectChatPushSettings(t *testing.T) {
	convey.Convey("Test", t, test.StartUp(func(ctx context.Context) {
		settings, err := GetFsProjectChatPushSettings(2373, &projectvo.GetFsProjectChatPushSettingsReq{
			OrgId:         2373,
			ChatId:        "", // oc_1ca75cc5ace116072aae6e2369f4003f
			ProjectId:     13908,
			SourceChannel: sdk_const.SourceChannelFeishu,
		})
		if err != nil {
			t.Error(err)
			return
		}
		t.Log(json.ToJsonIgnoreError(settings))
	}))
}

func TestGetProjectMainChatId(t *testing.T) {
	convey.Convey("Test", t, test.StartUp(func(ctx context.Context) {
		result, err := GetProjectMainChatId(2373, 29611, 60927, "fs")
		if err != nil {
			t.Error(err)
			return
		}
		t.Log(result)
	}))
}

func TestGetFsProjectChatPushSettings1(t *testing.T) {
	convey.Convey("Test", t, test.StartUp(func(ctx context.Context) {
		settings, err := GetFsProjectChatPushSettings(2373, &projectvo.GetFsProjectChatPushSettingsReq{
			OrgId:         2373,
			ChatId:        "", // oc_55ee4bf34666e23072335c190e65f675
			ProjectId:     13908,
			SourceChannel: sdk_const.SourceChannelFeishu,
		})
		if err != nil {
			t.Error(err)
			return
		}
		t.Log(json.ToJsonIgnoreError(settings))
	}))
}

func TestUpdateFsProjectChatPushSettings1(t *testing.T) {
	convey.Convey("Test", t, test.StartUp(func(ctx context.Context) {
		tables := []*vo.UpdateFsProjectChatPushSettingsOfTableParam{
			// {
			// 	ProjectID: 13777,
			// 	TableID:   "1530081219977416704",
			// }, {
			// 	ProjectID: 13777,
			// 	TableID:   "1533729135686324224",
			// }, {
			// 	ProjectID: 13825,
			// 	TableID:   "1533748534338129921",
			// },
		}
		err := UpdateFsProjectChatPushSettings(2373, 29611, sdk_const.SourceChannelFeishu,
			vo.UpdateFsProjectChatPushSettingsReq{
				ChatID:              "oc_55ee4bf34666e23072335c190e65f675",
				Tables:              tables,
				CreateIssue:         1,
				CreateIssueComment:  2,
				UpdateIssueCase:     2,
				ModifyColumnsOfSend: []*string{},
			})
		if err != nil {
			t.Error(err)
			return
		}
		t.Log("--end--")
	}))
}

func TestCutStr1(t *testing.T) {
	originStr := "我爱罗创建了新记录0"
	resStr := originStr
	if str.RuneLen(originStr) > 10 {
		resStr = str.Substr(originStr, 0, 10) + "......"
	}
	t.Log(resStr)
	originStr = "0123456789Hello world"
	resStr = originStr
	if str.RuneLen(originStr) > 10 {
		resStr = str.Substr(originStr, 0, 10) + "......"
	}
	t.Log(resStr)
}

func TestGetProTableChatSettingBoArrByTableId1(t *testing.T) {
	convey.Convey("Test", t, test.StartUp(func(ctx context.Context) {
		settings, err := domain.GetProTableChatSettingBoArrByTableId(2373, 13777, 1530081219977416704)
		if err != nil {
			t.Error(err)
			return
		}
		t.Log(json.ToJsonIgnoreError(settings))
	}))
}

func TestCheckIsAppCollaborator1(t *testing.T) {
	convey.Convey("Test", t, test.StartUp(func(ctx context.Context) {
		isCollaborate := domain.CheckIsAppCollaborator(2373, 1561911130903904257, 29480)
		t.Log(isCollaborate)
	}))
}

func TestCheckIsInChat1(t *testing.T) {
	convey.Convey("Test", t, test.StartUp(func(ctx context.Context) {
		orgId := int64(2373)
		currentUserId := int64(29480)
		baseOrgInfo, err := orgfacade.GetBaseOrgInfoRelaxed(orgId)
		if err != nil {
			log.Errorf("[GetProjectMainChatId] GetBaseOrgInfoRelaxed err: %v", err)
			return
		}
		curUserResp := orgfacade.GetBaseUserInfo(orgvo.GetBaseUserInfoReqVo{
			OrgId:  orgId,
			UserId: currentUserId,
		})
		isIn, err := domain.CheckUserIsInFsChat(baseOrgInfo, curUserResp.BaseUserInfo, "oc_7bfc74753cff56186e77e0ae9e300c8a")
		if err != nil {
			t.Error(err)
			return
		}
		t.Log(isIn)
	}))
}
