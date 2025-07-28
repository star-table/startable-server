package orgsvc

import (
	"context"
	"fmt"
	"testing"
	"time"

	sdk_const "gitea.bjx.cloud/allstar/platform-sdk/consts"
	"github.com/star-table/startable-server/common/core/consts"
	"github.com/star-table/startable-server/common/core/errs"
	"github.com/star-table/startable-server/common/core/util/json"
	"github.com/star-table/startable-server/common/core/util/strs"
	"github.com/star-table/startable-server/common/model/vo/orgvo"
	"github.com/star-table/startable-server/common/test"
	"github.com/google/martian/log"
	"github.com/magiconair/properties/assert"
	"github.com/smartystreets/goconvey/convey"
)

func TestCreateOrg(t *testing.T) {

	orgName := "0123456789123456"
	orgNameLen := strs.Len(orgName)

	var err error = nil
	if orgNameLen == 0 || orgNameLen > 256 {
		log.Error("组织名称长度错误")
		err = errs.BuildSystemErrorInfo(errs.OrgNameLenError)
	}

	assert.Equal(t, err, nil)

}

//func TestGetOutOrgInfoByOrgIdBatch(t *testing.T) {
//	convey.Convey("Test", t, test.StartUp(func(ctx context.Context) {
//		resp, err := GetOutOrgInfoByOrgIdBatch([]int64{1004})
//		if err != nil {
//			t.Error(err)
//		}
//		log.Info(json.ToJsonIgnoreError(resp))
//	}))
//}

// 暂不支持创建文件夹，需要界面上创建。测试时，开发环境就以 1380077674080432129 作为文件夹 id。
func TestCreateLcFolder(t *testing.T) {
	convey.Convey("Test", t, test.StartUp(func(ctx context.Context) {
		orgId := int64(202505060582035456)
		opUserId := int64(3)
		resp, err := domain.CreateLcFolder(orgId, opUserId, "dev测试应用001")
		if err != nil {
			t.Error(err)
		}
		log.Info(json.ToJsonIgnoreError(resp))
	}))
}

//func TestInitExistOrg(t *testing.T) {
//	convey.Convey("Test", t, test.StartUp(func(ctx context.Context) {
//		orgId := int64(202505060582035456)
//		opUserId := int64(3)
//		isQueryTokenInfo := 1
//		isNeedSetPaid := 0
//		resp, err := InitExistOrg(orgId, opUserId, vo.InitExistOrgReq{
//			NeedOcrConfig:      0,
//			NeedPriority:       0,
//			NeedSummaryTable:   0,
//			NeedSetToPaid:      &isNeedSetPaid,
//			QueryAuthTokenInfo: &isQueryTokenInfo,
//		})
//		if err != nil {
//			t.Error(err)
//		}
//		log.Info(json.ToJsonIgnoreError(resp))
//	}))
//}

func TestTimeAddDate1(t *testing.T) {
	endTime := time.Now().AddDate(100, 0, 0)
	t.Log(endTime.Format(consts.AppTimeFormatYYYYMMDDHHmm))
}

func Test_CheckSpecificScope(t *testing.T) {
	convey.Convey("Test", t, test.StartUp(func(ctx context.Context) {
		res, err := CheckSpecificScope(1213, 1624, "calendar:calendar:access")
		if err != nil {
			t.Error(err)
			return
		}
		t.Log(json.ToJsonIgnoreError(res))
	}))
}

func Test_ApplyScopes(t *testing.T) {
	convey.Convey("Test", t, test.StartUp(func(ctx context.Context) {
		res, err := ApplyScopes(1213, 1624)
		if err != nil {
			t.Error(err)
			return
		}
		t.Log(json.ToJsonIgnoreError(res))
	}))
}

func TestGetOrgIdListByPage1(t *testing.T) {
	convey.Convey("Test", t, test.StartUp(func(ctx context.Context) {
		res, err := GetOrgIdListByPage(orgvo.GetOrgIdListByPageReqVoData{}, 1, 100)
		if err != nil {
			t.Error(err)
			return
		}
		t.Log(json.ToJsonIgnoreError(res))
	}))
}

func TestGetOrgIdListByPage2(t *testing.T) {
	convey.Convey("Test", t, test.StartUp(func(ctx context.Context) {
		page := 1
		size := 50
		count1 := 0
		for {
			orgIdsResp := orgfacade.GetOrgIdListByPage(orgvo.GetOrgIdListByPageReqVo{
				Page: page,
				Size: size,
			})
			if orgIdsResp.Failure() {
				log.Error(orgIdsResp.Message)
				return
			}
			count1 += len(orgIdsResp.Data.OrgIds)
			t.Log(json.ToJsonIgnoreError(orgIdsResp))
			if len(orgIdsResp.Data.OrgIds) < 50 {
				break
			}
			page++
		}
		fmt.Printf("count1 val is: %v", count1)
	}))
}

func TestGetOrgBoListByPage(t *testing.T) {
	convey.Convey("Test", t, test.StartUp(func(ctx context.Context) {
		resp, err := GetOrgBoListByPage(orgvo.GetOrgIdListByPageReqVoData{}, 100, 100)
		if err != nil {
			t.Error(err)
			return
		}
		t.Log(json.ToJsonIgnoreError(resp))
	}))
}

func TestUpdateOrgRemarkSetting(t *testing.T) {
	convey.Convey("Test", t, test.StartUp(func(ctx context.Context) {
		isOk, err := UpdateOrgRemarkSetting(orgvo.UpdateOrgRemarkSettingReqVo{
			OrgId:  1574,
			UserId: 0,
			Input: orgvo.SaveOrgSummaryTableAppIdReqVoData{
				ProjectFolderAppId: 1474355741010747394,
				IssueFolderAppId:   1474355740213829634,
			},
		})
		if err != nil {
			t.Error(err)
			return
		}
		t.Log(isOk)
	}))
}

func TestUserOrgList(t *testing.T) {
	convey.Convey("Test", t, test.StartUp(func(ctx context.Context) {
		fmt.Println(UserOrganizationList(1014))
	}))
}

func TestSwitchUserOrganization(t *testing.T) {
	convey.Convey("Test", t, test.StartUp(func(ctx context.Context) {
		fmt.Println(SwitchUserOrganization(2373, 33991, "o2373u33990ta19cbb3542c740d0b9d6c2a95b369feb"))
	}))
}

func TestRemainDays1(t *testing.T) {
	payEndTime := time.Date(2038, 8, 29, 18, 32, 0, 0, time.Local)
	diffDays := domain.GetRemainDays(payEndTime)
	t.Log(diffDays)
}

func TestSendCardToAdminForUpgrade1(t *testing.T) {
	convey.Convey("Test", t, test.StartUp(func(ctx context.Context) {
		err := SendCardToAdminForUpgrade(2373, 29611, sdk_const.SourceChannelFeishu)
		if err != nil {
			t.Error(err)
			return
		}
		t.Log("--end--")
	}))
}
