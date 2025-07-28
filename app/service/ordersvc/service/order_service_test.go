package ordersvc

import (
	"context"
	"testing"
	"time"

	"github.com/star-table/startable-server/common/core/consts"
	"github.com/star-table/startable-server/common/test"

	"github.com/star-table/startable-server/common/model/bo"
	"github.com/smartystreets/goconvey/convey"
)

func TestAddFsOrder1(t *testing.T) {
	convey.Convey("Test", t, test.StartUp(func(ctx context.Context) {
		payTime := time.Now()
		createTime := payTime.Add(-(3600 * time.Second))
		orderId, err := AddFsOrder(bo.OrderFsBo{
			Id:              1001001,
			OrgId:           1574000,
			OrderId:         "7083115533758955522",
			PricePlanId:     "price_a15167c3baaf5013",
			PricePlanType:   "trial",
			Seats:           0,
			BuyCount:        1,
			PaidTime:        payTime,
			OrderCreateTime: createTime,
			Status:          "normal",
			BuyType:         "buy",
			OrderPayPrice:   1200,
			TenantKey:       "2c15f533e1cf5758",
			CreateTime:      createTime,
			UpdateTime:      createTime,
			Version:         1,
			IsDelete:        2,
		})
		if err != nil {
			t.Error(err)
			return
		}
		t.Log("--end--", orderId)
	}))
}

func TestDateCmp1(t *testing.T) {
	d1 := time.Date(2022, 9, 31, 11, 59, 10, 0, time.Local)
	endDateStr := d1.Format(consts.AppDateFormat)
	nowDateStr := time.Now().Format(consts.AppDateFormat)
	t.Log(nowDateStr <= endDateStr)
}

func TestCheckIsTrial1(t *testing.T) {
	convey.Convey("Test", t, test.StartUp(func(ctx context.Context) {
		isTrial, err := CheckFsOrgIsInTrial(2373)
		if err != nil {
			t.Error(err)
			return
		}
		t.Log(isTrial)
	}))
}

//func TestGetOutOrgId(t *testing.T)  {
//	convey.Convey("test", t ,test.StartUp(func(ctx context.Context) {
//		outOrgId := ""
//		orgInfoResp := orgfacade.GetBaseOrgInfoByOutOrgId(orgvo.GetBaseOrgInfoByOutOrgIdReqVo{OutOrgId: outOrgId})
//		if orgInfoResp.Successful() {
//			fmt.Println(1111111111)
//		}
//		if orgInfoResp.Failure() {
//			t.Error(orgInfoResp.Error())
//		}
//		t.Log(orgInfoResp.BaseOrgInfo)
//	}))
//}

func TestAddDingOrder(t *testing.T) {
	convey.Convey("test", t, test.StartUp(func(ctx context.Context) {
		data := bo.OrderDingBo{OrderId: ""}
		order, err := AddDingOrder(data)
		if err != nil {
			t.Error(err)
		}
		t.Log(order)
	}))
}

func TestAddWeiXinOrder(t *testing.T) {
	convey.Convey("test", t, test.StartUp(func(ctx context.Context) {
		paidTime := "2023-02-22 18:12:12"
		endTime := "2023-10-22 18:12:12"
		paid, _ := time.Parse("2006-01-02 15:04:05", paidTime)
		end, _ := time.Parse("2006-01-02 15:04:05", endTime)
		data := bo.OrderWeiXinBo{
			OrderId:   "OI000007A8809163F58BFC66D45D2N",
			OutOrgId:  "wpkYCoBwAAmwNhPL3EElsSbNRvySUKKA",
			EditionId: consts.WeiXinFlagshipId,
			PaidTime:  paid,
			EndTime:   end,
		}
		order, err := AddWeiXinOrder(data)
		if err != nil {
			t.Error(err)
		}
		t.Log(order)
	}))
}

func TestGetOrderPayInfo(t *testing.T) {
	convey.Convey("test", t, test.StartUp(func(ctx context.Context) {
		orderPayInfo, err := GetOrderPayInfo(3075)
		if err != nil {
			t.Error(err)
		}
		t.Log(orderPayInfo)
	}))
}
