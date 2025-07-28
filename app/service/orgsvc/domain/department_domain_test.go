package orgsvc

import (
	"context"
	"fmt"
	"testing"
	"time"

	wework "github.com/go-laoji/wecom-go-sdk"

	"gitea.bjx.cloud/allstar/platform-sdk/sdk_interface"

	platform_sdk "gitea.bjx.cloud/allstar/platform-sdk"
	sdkConsts "gitea.bjx.cloud/allstar/platform-sdk/consts"
	"github.com/star-table/startable-server/app/service/projectsvc/test"
	"github.com/star-table/startable-server/common/extra/third_platform_sdk"

	"github.com/star-table/startable-server/common/core/util/slice"
	"github.com/star-table/startable-server/common/library/db/mysql"
	"github.com/smartystreets/goconvey/convey"
	"upper.io/db.v3/lib/sqlbuilder"
)

func TestGetDepartmentIdsByOutDepartmentIds(t *testing.T) {
	convey.Convey("TestGetDepartmentIdsByOutDepartmentIds", t, test.StartUpWithUserInfo(func(userId, orgId int64) {
		queryParam := []string{"od-e262795e10f1d4777293b9be65d5703e"}
		res, err := GetDepartmentIdsByOutDepartmentIds(orgId, queryParam)
		if err != nil {
			t.Logf("%s\n", err)
		}
		log.Infof("orgId: %v, dept ids: %v\n", orgId, res)
		convey.So(res, convey.ShouldNotEqual, nil)
	}))
}

func TestUpper(t *testing.T) {
	convey.Convey("Test", t, test.StartUpWithUserInfo(func(userId, orgId int64) {
		configList := []po.PpmOrcFunctionConfig{
			{
				Id:           5,
				OrgId:        0,
				FunctionCode: "suhanyu-test5",
				Remark:       "",
				Creator:      0,
				CreateTime:   time.Time{},
				Updator:      0,
				UpdateTime:   time.Time{},
				Version:      0,
				IsDelete:     2,
			},
			{
				Id:           6,
				OrgId:        0,
				FunctionCode: "suhanyu-test6",
				Remark:       "",
				Creator:      0,
				CreateTime:   time.Time{},
				Updator:      0,
				UpdateTime:   time.Time{},
				Version:      0,
				IsDelete:     2,
			},
		}
		oneObj := &po.PpmOrcFunctionConfig{}

		mysql.SelectById(oneObj.TableName(), 5, oneObj)
		err := mysql.TransX(func(tx sqlbuilder.Tx) error {
			return mysql.TransBatchInsert(tx, &po.PpmOrcFunctionConfig{}, slice.ToSlice(configList))
		})
		convey.ShouldEqual(err, nil)
	}))
}

func TestGetTopDepartmentInfoList(t *testing.T) {
	convey.Convey("TestGetDepartmentIdsByOutDepartmentIds", t, test.StartUpWithUserInfo(func(userId, orgId int64) {
		//name := "F"
		t.Log(GetDepartmentMembersPaginate(1004, "", nil, 0, 0, nil, false, 0, 0))

	}))

}

func TestWeixin(t *testing.T) {
	convey.Convey("TestGetDepartmentIdsByOutDepartmentIds", t, test.StartUp(func(ctx context.Context) {
		third_platform_sdk.InitPlatFormSdk(GetCorpInfoFromDB)
		SetCacheCorpInfo(&sdk_interface.CorpInfo{
			OrgId:         1,
			AgentId:       1000010,
			PermanentCode: "xGXP4GU8Iat0HILJMtj15Ro9PiY1LEljXLZ5n4TriiQ",
			CorpId:        "wpkYCoBwAAlfm2nZzlQpTqrq5P4jKMXg",
		})
		client, sdkErr := platform_sdk.GetClient(sdkConsts.SourceChannelWeixin, "wpkYCoBwAAlfm2nZzlQpTqrq5P4jKMXg")
		if sdkErr != nil {
			log.Errorf("[InitOrg] platform_sdk.GetClient err: %v", sdkErr)
		}
		//fmt.Println(client.CheckIsAdminAuth())
		//fmt.Println(client.CheckIsAdminAuth())
		weixinClient := client.GetOriginClient().(wework.IWeWork)
		fmt.Println(weixinClient.UserId2OpenId(1, "wokYCoBwAAVITjs-L935_l2J_dcBIwLg"))
		//fmt.Println(weixinClient.UserSimpleList(2820, 1, 0))
		//fmt.Println(weixinClient.GetAuthInfo("wpkYCoBwAAlfm2nZzlQpTqrq5P4jKMXg", "8qWagmoJ68FlhuKE3OIT0y8JcpioYr19JLOSfWfDxkg"))
		//fmt.Println(weixinClient.UserGet(2820, "wokYCoBwAAVITjs-L935_l2J_dcBIwLg"))
	}))
}
