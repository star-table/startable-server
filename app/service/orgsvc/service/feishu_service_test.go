package orgsvc

import (
	"context"
	"fmt"
	"net/url"
	"testing"
	"time"

	"gitea.bjx.cloud/allstar/polaris-backend/service/platform/orgsvc/po"
	"github.com/star-table/startable-server/app/service/projectsvc/test"
	"github.com/star-table/startable-server/common/extra/third_platform_sdk"

	"github.com/star-table/startable-server/common/model/vo/orgvo"

	sdk_const "gitea.bjx.cloud/allstar/platform-sdk/consts"

	"github.com/star-table/startable-server/common/core/consts"
	"github.com/star-table/startable-server/common/core/util"
	"github.com/star-table/startable-server/common/core/util/json"
	"github.com/star-table/startable-server/common/library/db/mysql"
	"github.com/star-table/startable-server/common/model/vo"
	"github.com/smartystreets/goconvey/convey"
	"upper.io/db.v3"
)

func TestGetFsAccessToken(t *testing.T) {
	convey.Convey("Test", t, test.StartUpWithUserInfo(func(userId, orgId int64) {
		t.Log(GetFsAccessToken(2373, 29611))
	}))
}

func TestAssemblyFeiShuUserOrgRelationInfo(t *testing.T) {
	end, _ := time.Parse(consts.AppTimeFormat, "2020-10-16 12:10:10")

	t.Log(time.Now().Format(consts.AppTimeFormat))
	t.Log(time.Now().After(end.Local()))
}

func TestFeiShuAuth(t *testing.T) {
	convey.Convey("TestUpdateOrgMemberStatus", t, test.StartUpWithUserInfo(func(userId, orgId int64) {
		deptInfo := &po.PpmOrgDepartmentOutInfo{}
		err := mysql.SelectOneByCond(consts.TableDepartmentOutInfo, db.Cond{
			consts.TcIsDelete:           consts.AppIsNoDelete,
			consts.TcOrgId:              1111,
			consts.TcOutOrgDepartmentId: "aa",
		}, deptInfo)
		t.Log(deptInfo, err)

	}))

}

func TestFsDeptChange(t *testing.T) {
	convey.Convey("TestUpdateOrgMemberStatus", t, test.StartUpWithUserInfo(func(userId, orgId int64) {
		t.Log(PlatformDeptChange(1417, "dept_add", sdk_const.SourceChannelFeishu, "5e5292a64cb77c4c"))
	}))
}

func TestChangeDeptScope(t *testing.T) {
	convey.Convey("TestUpdateOrgMemberStatus", t, test.StartUpWithUserInfo(func(userId, orgId int64) {
		t.Log(ChangeDeptScope(1417))
	}))
}

func TestFeiShuAuth1(t *testing.T) {
	convey.Convey("Test", t, test.StartUp(func(ctx context.Context) {
		codeType := 2
		resp, err := FeiShuAuth(vo.FeiShuAuthReq{
			Code:     "FJD8pKdxSnTQnSnTnxdPug",
			CodeType: &codeType,
		})
		if err != nil {
			t.Error(err)
			return
		}
		t.Log(json.ToJsonIgnoreError(resp))
	}))
}

func TestGetJsAPITicket(t *testing.T) {
	convey.Convey("Test", t, test.StartUp(func(ctx context.Context) {
		ticket, err := GetJsAPITicket(2376)
		t.Log(ticket)
		t.Log(err)
	}))
}

func TestGetDifIds1(t *testing.T) {
	oldIds := []int64{1, 12, 1, 2, 3}
	newIds := []int64{2, 3, 2, 10}
	delIds, addIds := util.GetDifMemberIds(oldIds, newIds)
	fmt.Printf("del: %s, add: %s", json.ToJsonIgnoreError(delIds), json.ToJsonIgnoreError(addIds))
}

func TestClearOutDept1(t *testing.T) {
	convey.Convey("Test", t, test.StartUp(func(ctx context.Context) {
		orgId := int64(2466)
		if err := domain.CheckAndClearExtraOutDept(orgId); err != nil {
			t.Error(err)
			return
		}
		if err := domain.ClearDeptForSyncDept(orgId); err != nil {
			t.Error(err)
		}
		t.Log("--end--")
	}))
}

func TestSyncDept1(t *testing.T) {
	convey.Convey("Test", t, test.StartUp(func(ctx context.Context) {
		orgId := int64(2466)
		if err := ChangeDeptScope(orgId); err != nil {
			t.Error(err)
			return
		}
		t.Log("--end--")
	}))
}

func TestEncode(t *testing.T) {
	convey.Convey("Test", t, test.StartUp(func(ctx context.Context) {
		fmt.Println(url.QueryEscape("https://www.baidu.com"))
	}))
}

func TestBindProject(t *testing.T) {
	convey.Convey("Test", t, test.StartUp(func(ctx context.Context) {
		third_platform_sdk.InitPlatFormSdk(nil)
		fmt.Println(BindCoolApp(orgvo.BindCoolAppReq{
			OrgId: 2373,
			Input: &orgvo.BindCoolAppData{
				OpenConversationId: "cidPKzlav39ckmNXTWpcKP9yQ==",
				ProjectId:          62640,
				AppId:              1625841128089681922,
			},
		}))
	}))
}

func TestGetSpaceList(t *testing.T) {
	convey.Convey("Test", t, test.StartUp(func(ctx context.Context) {
		third_platform_sdk.InitPlatFormSdk(nil)
		fmt.Println(GetSpaceList(orgvo.GetSpaceListReq{
			OrgId:  2588,
			UserId: 43637,
		}))
	}))
}

func TestGetDingJsAPISign(t *testing.T) {
	convey.Convey("Test", t, test.StartUp(func(ctx context.Context) {
		third_platform_sdk.InitPlatFormSdk(nil)
		fmt.Println(GetDingJsAPISign(&orgvo.JsAPISignReq{
			Type:   "",
			URL:    "http://localhost:9201/workbench",
			CorpID: "dingbd385c1cf7fd62a724f2f5cc6abecb85",
		}))
	}))
}
