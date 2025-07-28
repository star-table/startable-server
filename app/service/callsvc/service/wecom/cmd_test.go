package callsvc

import (
	"context"
	"encoding/xml"
	"fmt"
	"testing"

	"github.com/go-laoji/wecom-go-sdk/pkg/svr/logic"
	"github.com/go-laoji/wxbizmsgcrypt"

	platform_sdk "gitea.bjx.cloud/allstar/platform-sdk"
	sdk_const "gitea.bjx.cloud/allstar/platform-sdk/consts"
	wework "github.com/go-laoji/wecom-go-sdk"

	"github.com/star-table/startable-server/app/facade/orgfacade"
	"github.com/star-table/startable-server/common/extra/third_platform_sdk"

	"github.com/star-table/startable-server/common/test"
	"github.com/smartystreets/goconvey/convey"
)

func TestCreateAuthHandle(t *testing.T) {
	convey.Convey("TestCreateAuthHandle", t, test.StartUp(func(ctx context.Context) {
		third_platform_sdk.InitPlatFormSdk(orgfacade.GetCorpInfoFromDB)

		sdk, err := platform_sdk.GetClient(sdk_const.SourceChannelWeixin, "wwf36b5e6ef0b569ac")
		if err != nil {
			log.Errorf("[wecom] GetClient error::%v", err)
			return
		}
		client := sdk.GetOriginClient().(wework.IWeWork)
		fmt.Println(client.CreateNewOrder(wework.CreateOrderRequest{
			CorpId:      "wwf36b5e6ef0b569ac",
			BuyerUserid: "JiaoWoLaoJiangHaoLe",
			AccountCount: struct {
				BaseCount            int `json:"base_count" validate:"required_without=ExternalContactCount,max=1000000"`
				ExternalContactCount int `json:"external_contact_count" validate:"required_without=BaseCount,max=1000000"`
			}{BaseCount: 10},
			AccountDuration: wework.AccountDuration{
				Months: 1,
			},
		}))
	}))
}

func TestSystem(t *testing.T) {
	body := []byte("<xml><ToUserName><![CDATA[ww9b85ae8ff033ee89]]></ToUserName><Encrypt><![CDATA[r8luKJZsCnOGxAUBnS+qvXyHl3V+AzNzqHiAMTYmz1MvwByZ9jrLRjMSF4J+iyfp8owjVKK4HdFl6ZlCtBI2293QH2MvCMFJNT5Z6WLuyig8z//ziT/CZ5ZfkxnTIGg3t+fH1FKGkN+z1GZOTM1+Oh5y4thSxyYlqPYiXyaC0ZfCOE6ckvIGds0sCaOU1Pa/FCQpb4JlDwI2BuRj3ghbtKlrIvVLOt5spPeyyR0I9dSaV7vRT0uXomkjsgyZords28oZD+6SqGEXFHZHI1PSTICGIuAc62h+96IxhmTixt0HJBMTFuQI95t9ZVWNBKE/+rlTjYWM1PZTZCfoFaGPO+rcXXW8Akj6j1O2DNTWhx5lf6JUo2h//SJ367ParPK1]]></Encrypt><AgentID><![CDATA[]]></AgentID></xml>")
	var params = logic.EventPushQueryBinding{
		MsgSign:   "495d6718b26e018f4b7a2bc32a6011e462336cd0",
		Timestamp: "1684227827",
		Nonce:     "1684326459",
	}
	var bizData logic.BizData
	xml.Unmarshal(body, &bizData)
	wxcpt := wxbizmsgcrypt.NewWXBizMsgCrypt("LV5N2mHIhRu78KqgQDPODB", "lqPZVA4mtwhRMAXR6li5aHqkFc0C57DeHnPLBQacRiP",
		bizData.ToUserName, wxbizmsgcrypt.XmlType)
	fmt.Println(wxcpt.DecryptMsg(params.MsgSign, params.Timestamp, params.Nonce, body))
}
