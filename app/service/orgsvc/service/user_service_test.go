package orgsvc

import (
	"context"
	"fmt"
	"testing"
	"time"

	platform_sdk "gitea.bjx.cloud/allstar/platform-sdk"
	consts3 "gitea.bjx.cloud/allstar/platform-sdk/consts"
	wework "github.com/go-laoji/wecom-go-sdk"
	"github.com/google/martian/log"

	"github.com/star-table/startable-server/app/service/projectsvc/test"
	"github.com/star-table/startable-server/common/core/config"
	consts1 "github.com/star-table/startable-server/common/core/consts"
	"github.com/star-table/startable-server/common/core/threadlocal"
	"github.com/star-table/startable-server/common/core/types"
	"github.com/star-table/startable-server/common/core/util/json"
	"github.com/star-table/startable-server/common/extra/third_platform_sdk"
	"github.com/star-table/startable-server/common/model/vo"
	"github.com/star-table/startable-server/common/model/vo/orgvo"
	"github.com/getsentry/sentry-go"
	"github.com/smartystreets/goconvey/convey"
)

func TestUserConfigInfo(t *testing.T) {
	convey.Convey("Test sql", t, test.StartUp(func(ctx context.Context) {
		config, err := UserConfigInfo(1574, 24462)

		convey.ShouldBeNil(err)
		convey.ShouldNotBeNil(config)
		t.Log(json.ToJsonIgnoreError(config))
	}))
}

func TestWeixinScan(t *testing.T) {
	convey.Convey("Test", t, test.StartUp(func(ctx context.Context) {
		fmt.Println(CheckPersonWeiXinQrCode(&orgvo.CheckQrCodeScanReq{SceneKey: "19c2ea2addbf4889940889fd78f5bef7"}))
	}))
}

func TestWeixinRegisterUrl(t *testing.T) {
	convey.Convey("Test", t, test.StartUp(func(ctx context.Context) {
		third_platform_sdk.InitPlatFormSdk(orgfacade.GetCorpInfoFromDB)
		fmt.Println(GetWeiXinRegisterUrl())
	}))
}

func TestWeixinInstallUrl(t *testing.T) {
	convey.Convey("Test", t, test.StartUp(func(ctx context.Context) {
		third_platform_sdk.InitPlatFormSdk(orgfacade.GetCorpInfoFromDB)
		fmt.Println(GetWeiXinInstallUrl())
	}))
}

func TestWeixinContactSuccess(t *testing.T) {
	convey.Convey("Test", t, test.StartUp(func(ctx context.Context) {
		third_platform_sdk.InitPlatFormSdk(orgfacade.GetCorpInfoFromDB)
		sdk, err := platform_sdk.GetClient(consts3.SourceChannelWeixin, "")
		if err != nil {
			fmt.Printf("[handlerRegisterCorp] GetClient error::%v", err)
			return
		}
		weWork := sdk.GetOriginClient().(wework.IWeWork)
		resp := weWork.ContactSyncSuccess("kPCE4TvFh_p1kNpQuE9zNNjx-9UT_eaq-GAnTkPmobaSx4erkAZ-vTa1uUW8hF5PD4TR5T6GsrKl5mo-WbTDOQDAyw03vPCdVpBrDbS2huec-JEgq9OMKtorKBffYOKAMWkUI2g8cdNM6ebUGV8-ObM514dR_v755Ehi5XmLeUHEqCAXxwDzT6x-p7D_PRlxNWlmGwwfN0BqCHAOv4JULg")
		if resp.ErrCode != 0 {
			fmt.Printf("[handlerRegisterCorp] ContactSyncSuccess error::%v, token:%v", err, "")
		}
		fmt.Println(resp)
	}))
}

func TestDingTalkAuthCode(t *testing.T) {
	convey.Convey("Test sql", t, test.StartUp(func(ctx context.Context) {
		third_platform_sdk.InitPlatFormSdk(orgfacade.GetCorpInfoFromDB)
		fmt.Println(AuthWeiXinCode(orgvo.AuthWeiXinCodeReqVo{
			Code:     "tcZPj1REFOKnYCFC3RRMsQbgcMza-l7KXEEDIFuw7Zg",
			CodeType: 0,
			CorpId:   "wwf36b5e6ef0b569ac",
		}))
	}))
}

func TestPersonalInfo(t *testing.T) {
	convey.Convey("Test", t, test.StartUp(func(ctx context.Context) {
		curUser, err := orgfacade.GetCurrentUserRelaxed(ctx)
		if err != nil {
			log.Error(err)
			return
		}
		t.Logf("user: %s", json.ToJsonIgnoreError(curUser))
		info, err := PersonalInfo(curUser.OrgId, curUser.UserId, curUser.SourceChannel)
		if err != nil {
			t.Error(err)
			return
		}
		t.Log(json.ToJsonIgnoreError(info))
	}))
}

func TestGetQrCode(t *testing.T) {
	convey.Convey("Test", t, test.StartUp(func(ctx context.Context) {
		fmt.Println(PersonWeiXinQrCode(&orgvo.PersonWeiXinQrCodeReqVo{}))
	}))
}

func TestQrCodeScan(t *testing.T) {
	convey.Convey("Test", t, test.StartUp(func(ctx context.Context) {
		fmt.Println(PersonWeiXinQrCodeScan(&orgvo.QrCodeScanReq{
			OpenId:   "oKf_M6uwXaORyyShO02gNNRP57hk",
			SceneKey: "8c48a4972a544d23888d3413ea6acf0c",
		}))
	}))
}

func TestCheckQrCodeScan(t *testing.T) {
	convey.Convey("Test", t, test.StartUp(func(ctx context.Context) {
		fmt.Println(CheckPersonWeiXinQrCode(&orgvo.CheckQrCodeScanReq{SceneKey: "8c48a4972a544d23888d3413ea6acf0c"}))
	}))
}

func TestPersonWeiXinBind(t *testing.T) {
	convey.Convey("Test", t, test.StartUp(func(ctx context.Context) {
	}))
}

func TestUpdateUserInfo(t *testing.T) {

	convey.Convey("Test UpdateUserInfo", t, test.StartUpWithUserInfo(func(userId, orgId int64) {

		updateFields := []string{"name", "avatar", "sex", "birthday"}

		avatar := "https://timgsa.baidu.com/timg?image&quality=80&size=b9999_10000&sec=1576325272960&di=64d46c5b9578915bbc91ddf4a282d917&imgtype=0&src=http%3A%2F%2Fimg1.xiazaizhijia.com%2Fwalls%2F20150910%2Fmiddle_eea71d559bf1341.jpg"
		name := "更改后的名字"
		sex := 1
		nowTime := types.NowTime()

		input := vo.UpdateUserInfoReq{
			Name:         &name,
			Avatar:       &avatar,
			UpdateFields: updateFields,
			Sex:          &sex,
			Birthday:     &nowTime,
		}

		resp, err := UpdateUserInfo(0, userId, input)

		convey.ShouldBeNil(err)
		convey.ShouldNotBeNil(resp)
	}))

}

func TestVerifyOrgUsers(t *testing.T) {
	convey.Convey("Test UpdateUserInfo", t, test.StartUpWithUserInfo(func(userId, orgId int64) {
		t.Log(VerifyOrgUsers(1001, []int64{1, 2}))
	}))

}

func TestJudgeUserIsAdmin(t *testing.T) {
	convey.Convey("Test UpdateUserInfo", t, test.StartUpWithUserInfo(func(userId, orgId int64) {
		t.Log(JudgeUserIsAdmin("2eb972ddc1cf1652", "ou_06d3425ac73fe283541a2045b10512e4", "fs"))
	}))
}

// 测试同步成员信息和组织架构
func TestSyncUserInfoFromFeiShu(t *testing.T) {
	convey.Convey("TestSyncUserInfoFromFeiShu", t, test.StartUpWithUserInfo(func(userId, orgId int64) {
		orgId = 1213
		userId = 1624
		err := SyncUserInfoFromFeiShu(orgId, userId, vo.SyncUserInfoFromFeiShuReq{
			NeedSyncName:       true,
			NeedSyncAvatar:     true,
			NeedSyncDepartment: false,
			SourceChannel:      "fs",
		})
		if err != nil {
			t.Fatalf("%v\n", err)
		}
		convey.So(err, convey.ShouldEqual, nil)
	}))
}

// 同步部门信息
//func TestSyncDepartments(t *testing.T) {
//	convey.Convey("TestSyncUserInfoFromFeiShu", t, test.StartUpWithUserInfo(func(userId, orgId int64) {
//		err := mysql.TransX(func(tx sqlbuilder.Tx) error {
//			err := syncDepartments(orgId, )
//			return err
//		})
//		if err != nil {
//			t.Fatalf("%v\n", err)
//		}
//		convey.So(err, convey.ShouldEqual, nil)
//	}))
//}

func TestFilterResignedUserIdsName(t *testing.T) {
	convey.Convey("TestSyncUserInfoFromFeiShu", t, test.StartUp(func(ctx context.Context) {
		orgId := int64(1004)
		userIds := []int64{1010, 1011, 1012, 1013, 1014}
		validUserIds, err := FilterResignedUserIds(orgId, userIds)
		if err != nil {
			t.Error(err)
		}
		t.Log(json.ToJsonIgnoreError(validUserIds))
	}))
}

func TestGetUserInfoByUserIds(t *testing.T) {
	convey.Convey("TestSyncUserInfoFromFeiShu", t, test.StartUp(func(ctx context.Context) {
		resp, err := GetUserInfoByUserIds(orgvo.GetUserInfoByUserIdsReqVo{
			UserIds: []int64{1624},
			OrgId:   1213,
		})
		if err != nil {
			t.Error(err)
		}
		t.Log(json.ToJsonIgnoreError(resp))
	}))
}

func TestGetPayRemind(t *testing.T) {
	convey.Convey("TestGetPayRemind", t, test.StartUp(func(ctx context.Context) {
		t.Log(GetPayRemind(2373, 29480))
	}))
}

func TestSentryError(t *testing.T) {
	// errs.UserInitError
	convey.Convey("TestSyncUserInfoFromFeiShu", t, test.StartUp(func(ctx context.Context) {

	}))
}

func TestDoSentry(t *testing.T) {
	convey.Convey("Test TestExportData", t, test.StartUp(func(ctx context.Context) {
		sentryConfig := config.GetSentryConfig()
		sentryDsn := ""
		if sentryConfig != nil {
			sentryDsn = sentryConfig.Dsn
		}
		dsn := sentryDsn
		applicationName := config.GetApplication().Name
		env := "dev"
		connectSentryErr := sentry.Init(sentry.ClientOptions{
			Dsn:         dsn,
			ServerName:  applicationName,
			Environment: env,
		})
		errMsg := "thisIsTestErrorForSentry"
		stackInfo := "【北极星警告】【告警日志】test error stack info. 004"
		event := sentry.NewEvent()
		event.EventID = sentry.EventID(errMsg)
		event.Message = errMsg + "\n" + stackInfo
		event.Environment = env
		event.ServerName = applicationName
		event.Tags[consts1.TraceIdKey] = threadlocal.GetTraceId()
		if connectSentryErr == nil {
			sentry.CaptureEvent(event)
			sentry.Flush(time.Second * 5)
		} else {
			fmt.Println(connectSentryErr)
		}
		fmt.Println(event.Message)
	}))
}

func TestSetVisitUserGuideStatus(t *testing.T) {
	convey.Convey("Test", t, test.StartUp(func(ctx context.Context) {
		info, err := SetVisitUserGuideStatus(orgvo.SetVisitUserGuideStatusReq{
			OrgId:  1574,
			UserId: 24462,
			Flag:   "newUserGuideStatus",
		})
		if err != nil {
			t.Error(err)
			return
		}
		t.Log(json.ToJsonIgnoreError(info))
	}))
}

func TestCreateOrgUser(t *testing.T) {
	convey.Convey("test", t, test.StartUp(func(ctx context.Context) {
		err := CreateOrgUser(2797, 34236, orgvo.CreateOrgMemberReq{
			PhoneRegion: "+86",
			PhoneNumber: "13812111210",
			Name:        "123333",
		})
		if err != nil {
			t.Error(err)
		}
	}))
}
