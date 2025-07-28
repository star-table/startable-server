package callsvc

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/star-table/startable-server/app/service/callsvc/domain"
	"github.com/star-table/startable-server/common/core/consts"
	"github.com/star-table/startable-server/common/core/util/json"
	"github.com/star-table/startable-server/common/model/bo"
	"github.com/star-table/startable-server/common/test"
	"github.com/smartystreets/goconvey/convey"
)

func TestFsCallbackHandler(t *testing.T) {

	mqBo := bo.FeiShuCallBackMqBo{
		MsgTime: time.Now(),
	}

	data := json.ToJsonIgnoreError(mqBo)

	newBo := &bo.FeiShuCallBackMqBo{}
	_ = json.FromJson(data, newBo)

	t.Log(json.ToJsonIgnoreError(data))

	msgTime := newBo.MsgTime

	fmt.Println(msgTime)
	fmt.Println(time.Now())
}

func TestUrgeCallback(t *testing.T) {
	convey.Convey("Test", t, test.StartUp(func(ctx context.Context) {
		handle := FsIssueUrgeHandler{}
		res, err := handle.Handle(CardReq{
			OpenId:        "ou_83776f77689ab1a5f399d5aafa27a0a4",
			OpenMessageId: "om_6c7bb6006d25f4902b7ff689dc95241a",
			TenantKey:     "2c0c7abea54f9758",
			Token:         "c-2d61c91b1376d695018b91ed894e5754da877fa7",
			Action: CardReqAction{
				Tag:      "select_static",
				Option:   "抱歉，现在不方便处理",
				Timezone: "",
				Value: &CardReqActionValue{
					IssueId:                           50541,
					OrgId:                             1213,
					UserId:                            1624,
					StatusType:                        0,
					CardType:                          "FsCardTypeIssueUrge",
					Action:                            "",
					FsCardValueUrgeIssueBeUrgedUserId: 1624,
					FsCardValueUrgeIssueUrgeText:      "看看进度啥样了。",
					FsCardValueIsGroupChat:            false,
					FsCardValueGroupChatId:            "",
				},
			},
		})
		if err != nil {
			t.Error(err)
		}
		t.Log(res)
	}))
}

func TestHandleError(t *testing.T) {
	convey.Convey("Test", t, test.StartUp(func(ctx context.Context) {
		//addFsUser("2e99b3ab0b0f1654", "xxxxxxxxxx")
	}))
}

func TestOrderPaidHandler_Handle(t *testing.T) {
	convey.Convey("Test", t, test.StartUp(func(ctx context.Context) {
		handler := OrderPaidHandler{}
		data := OrderPaidReq{
			Ts:    "",
			Uuid:  "",
			Token: "",
			Type:  "",
			Event: OrderPaidReqData{
				Type:          "order_paid",
				AppId:         "cli_9d40f5bf08f95108",
				OrderId:       "price_xxxxxxxxxxxxx",
				PricePlanId:   "price_xxxxxxxxxxxx",
				PricePlanType: consts.FsActiveEndDate,
				Seats:         0,
				BuyCount:      0,
				CreateTime:    "",
				PayTime:       "1620962188",
				BuyType:       "buy",
				SrcOrderId:    "",
				OrderPayPrice: 0,
				TenantKey:     "13383006d6cfd747",
			},
		}
		res, err := handler.Handle(json.ToJsonIgnoreError(data))
		t.Log(res)
		t.Log(err)
	}))
}

func TestAtTextMatch1(t *testing.T) {
	openId, userName := domain.MatchUserInsUserInfo("<at open_id=\"ou_b927e11615d382b034588908aedc988f\">苏汉宇</at>")
	t.Log(openId, userName)
}

// 拉人进群
func TestMockCallback1(t *testing.T) {
	inputText := "{\"schema\":\"2.0\",\"header\":{\"event_id\":\"1197b3b256e0101f9ae12b022c1e4b3e\",\"token\":\"QBUiTJGYZ90ynrLJVmu6Dh0loATMhUQt\",\"create_time\":\"1657246329997\",\"event_type\":\"im.chat.member.user.added_v1\",\"tenant_key\":\"12a813fb180f1758\",\"app_id\":\"cli_a02f641485b9d00b\"},\"event\":{\"chat_id\":\"oc_5c6acb9fd896eff99b24922c8178e753\",\"external\":false,\"operator_id\":{\"open_id\":\"ou_3ab7fe596cf91692218f744558ae157f\",\"union_id\":\"on_97cbc4023f1475c700f2a9f3d3919c59\",\"user_id\":null},\"operator_tenant_key\":\"12a813fb180f1758\",\"users\":[{\"name\":\"冉一萍\",\"tenant_key\":\"12a813fb180f1758\",\"user_id\":{\"open_id\":\"ou_5d296a1582a335dd5d9ddafa694a9e2d\",\"union_id\":\"on_695d8d6a88fd9170372be5d7751dafa8\",\"user_id\":null}}]}}"
	convey.Convey("Test", t, test.StartUp(func(ctx context.Context) {
		handler := ChatUserJoinHandler{}
		handler.Handle(inputText)
	}))
}

// 任务群聊-移除成员发送卡片
func TestMockCallback2(t *testing.T) {
	inputText := "{\"schema\":\"2.0\",\"header\":{\"event_id\":\"0468a06da4a341bb171ded63ba0d3517\",\"token\":\"QBUiTJGYZ90ynrLJVmu6Dh0loATMhUQt\",\"create_time\":\"1657090992748\",\"event_type\":\"im.chat.member.user.deleted_v1\",\"tenant_key\":\"12a813fb180f1758\",\"app_id\":\"cli_a02f641485b9d00b\"},\"event\":{\"chat_id\":\"oc_1f299b02fb39c783f4a5514b6b0343e0\",\"external\":false,\"operator_id\":{\"open_id\":\"ou_3ab7fe596cf91692218f744558ae157f\",\"union_id\":\"on_97cbc4023f1475c700f2a9f3d3919c59\",\"user_id\":null},\"operator_tenant_key\":\"12a813fb180f1758\",\"users\":[{\"name\":\"成卫忠\",\"tenant_key\":\"12a813fb180f1758\",\"user_id\":{\"open_id\":\"ou_2bbc27672c6ce4f0c2c73856f73ccbdb\",\"union_id\":\"on_4527d5e64525613a36461fcffd30df34\",\"user_id\":null}}]}}"
	convey.Convey("Test", t, test.StartUp(func(ctx context.Context) {
		handler := ChatUserDeleteHandler{}
		handler.Handle(inputText)
	}))
}

// 机器人进群回调测试
func TestMockBotEntry(t *testing.T) {
	inputText := "{\"schema\":\"2.0\",\"header\":{\"event_id\":\"fbcfdcdfd5ecb47afd28a209c50adbd0\",\"token\":\"QBUiTJGYZ90ynrLJVmu6Dh0loATMhUQt\",\"create_time\":\"1657098774993\",\"event_type\":\"im.chat.member.bot.added_v1\",\"tenant_key\":\"12a813fb180f1758\",\"app_id\":\"cli_a02f641485b9d00b\"},\"event\":{\"chat_id\":\"oc_1f299b02fb39c783f4a5514b6b0343e0\",\"external\":false,\"operator_id\":{\"open_id\":\"ou_3ab7fe596cf91692218f744558ae157f\",\"union_id\":\"on_97cbc4023f1475c700f2a9f3d3919c59\",\"user_id\":null},\"operator_tenant_key\":\"12a813fb180f1758\"}}"
	convey.Convey("Test", t, test.StartUp(func(ctx context.Context) {
		handler := ChatBotEntryHandler{}
		handler.Handle(inputText)
	}))
}

// 群聊解散时，清除群 id
func TestMockChatDisabledV1(t *testing.T) {
	inputText := "{\"schema\":\"2.0\",\"header\":{\"event_id\":\"66a5928a0eef47355b9ce2b4b87a4565\",\"token\":\"QBUiTJGYZ90ynrLJVmu6Dh0loATMhUQt\",\"create_time\":\"1657171394368\",\"event_type\":\"im.chat.disbanded_v1\",\"tenant_key\":\"12a813fb180f1758\",\"app_id\":\"cli_a02f641485b9d00b\"},\"event\":{\"chat_id\":\"oc_d5671919c5d4403b2d1a8f725d2a5917\",\"external\":false,\"operator_id\":{\"open_id\":\"ou_5d296a1582a335dd5d9ddafa694a9e2d\",\"union_id\":\"on_695d8d6a88fd9170372be5d7751dafa8\",\"user_id\":null},\"operator_tenant_key\":\"12a813fb180f1758\"}}"
	convey.Convey("Test", t, test.StartUp(func(ctx context.Context) {
		handler := ImChatDisbandV1Handler{}
		handler.Handle(inputText)
	}))
}

func TestParseFeishuEventData(t *testing.T) {
	data := `{"encrypt":"5GlQQKLRc+gg3DX9NhbpbCLvvBPLtCdj7/m+EfNO3fNi8uHtmOdA154Rmi4N0sfqb7cwLgpObfBkAE5ZDkbPxcXzc0Bx6CSbn+jM+bYEctZFQX55VEbpnHDRticBwqqhRqJ/aNfVdJ203mWM2efJBUTCHmiyra8Zbo52q8Z/nkyut6ybg91HbXtUVHSTR7Mo"}`
	//data := `{"encrypt":"klJfS9YFUH6KJ9lgybNdcuxH2YIhuFDRQmlkS+3SEXd7s/zi8LgkMHF/m5yjIXvzRc3Z5Y+DlQuS2wcu+johxRHa1TQMoQj29QKzwfczzwlJco6QVLj/Ql/+aswKqoCB1P4zviB3c+G5Gb95hj1JfrsSrOvAgRwAp4fKZbwSTuXNMXuiuUnzLpX3oO9WhMddurkHf5JnBcBhKYODkTwE/6hPbehzxyEN1kLyVMS+HC/gCEAcv6CUdPPTWNlW86/ytxzEzODA9T/PjSZxOjEMiCCm7SvhqyuwWgJ0BDccLafwVVosd3lLMEIHMHxxpXX0tpwoeQb+XiKiPSkju3BWxmwWq4cyJ6y92F3w4Bzh8ho="}`
	//data2 := `{"encrypt":"unU/BjjXBEy2WzE0D/tHCZVgPsnv8EVRsQPlF6W4ICZHGk8mbqnVne6q7YAYExErdwARivNuhBsR6J0GrEUpfVsae4rJzn2B72wBML8VB6kmKrTYq8pbHBsb8yz+nCJhQvDq/D5rWEgkyLWMg13feBtAgnqH7I3tXDeWpmioinxoLulVPzBExcAKIhKI7wv3YmkqyPK5oBdH/YBu+6yg/NinfkxHap5S53TqdMM+XnhAegA5mNZKv4AEibhRPqX++DhXBwPoRlDNItdggCrY7F/GV/HaYiZ1dxxv8jb8LyjFqNbmQdO9NE2jBxTAqoaCf0+NuqCuWZxw3rdA5dsvU98/ARqSipAZ6TGtCTaFNvOGfsLDt2NfUCA2NpK96ftpUKVabrSspyZs8yZ9ONR+lFW8osSXzvBAHD6XnAJSSuCOjImgJE795NulJrGp9JeDLNGjLiDw9vHKc8GfC3l4r6nxvtFCLROXWirLE/+mRhM="}`
	//data3 := `{"encrypt":"VfMCMMbx47yrqYp+tnppTFh71oX1l9+e95wDZOi+NJQhS8E/kGW0IeTHE2wXBcRFGNxA96mHiH1Y/IxFkXq3ovqJ6phne8lMyVZ2Sra0kcvGAKShmHmfrZ4o1ROwN4slIVbUrkCJAGW9WXa4zK7T7igD/dn3Y4oBRMIvxmnV0lVNt08zSoOzHEE1BAGhNhNSjVn17qjTeSf3YY4F+WpbKfYLBIQZvewaBCy3m8AbdfplFb0ZyClQCqcrRG3q7MLWnjeLasM7yiO1Gf0vyxGFoZ9eEe6UChkqjg2EcYkHRYBQHgMIOTnm/wC+2qZdcDXlPGl2wtoU5ygGzoLo7YwLcuGul0Adq9fD0dn4lTgIcXvf9M4jWYlF3H4lRAa8Jxuo3LY4qLC2pGHVlsMXSbyX+qErJa1/MfDDgtPhhdHs4Gpya5XJK5Et7ARE0f2w6t9DRtkh305/BpYHz3Y89frAeUiAoItghNU1S3tY/2y7Awo="}`
	//data4 := `{"encrypt":"4PxCdtEjviCtvTT0F+ItddyXCNK2IOg1V9/EueNqFyR4zOnl59ZbGsoE29ZQMy+qK5LCV1pYPt4p21uFzyw2Lz+fiMu3K5Lhrgl3OMgzXjBRClsBvQBdSuap3FLgYeguH3SPxanXzReAHSx+mzWtycZ7bIJbFA4wAWgi6DTL1ekK9sV6YTFt0W7PhLF5FQJcV6/w1ROftfMxHcvvD4kpfMv8X1BnYYTkbHTbxReQvPRG9flSmIeI0lfFB0VBwp+0RpzNUkMOXqfhepiRtudNGa3t3nX+NGdKJxqSCRyRIxxQriX61yhlvG+n5ujfojUJ0xJLaGZbnCO6+iYrL9Vip5vOTubJOuTav1sKkaffTe5rc/HS1vczmn66uDhsATVXNZw+tvaSK+QuAB4ooQIZHipEVtRoEhGStUjeROfz9W4lEa5NEKoFj5yw6OTbdf7+b/DvZAslRkK9kgpHo09tNE5/4Pj1zIcmhoLf0dC6RkXzAEGBVIo3JD/xo4J448xA1ldKC6rgDV9vQPbcPgKICpZ3EpsCxHMm0zGT3H0mvlcy7raZ6p0xOHc7Wen6eIM62VyMOZ0xI/X48EkV+AevqKknvuhASnjhM37n7erfWkhavL/AMMpJu+apru6Ymxnq5nYC6d/qW/EsvflotB6FNuHSDz52eCb+UE+ecpfA4NMaHHxh7OLAL9vTZnuI+2lOdYQr9knUrHBoxq1nqnGgugJ3lJg72Bcgmld16KJxuWzqVBBjJT7wfJsIBldGTUmnKyNW5ZhDibtD/2LBABWYpxR6eNwOxKTqub42/3pNBNMpvobkY1EqEOfangVhE743UkfhhZv04DCCN+VTjH4N0S4vPUMFN5YG/I8NIavXFSxhsj+VTkSYDy+9M4F3CIF4Ay/ptGoLbYtuhVyYG8ejOCX5V1GNrFL7JQo8gDWJsIg="}`
	convey.Convey("test", t, test.StartUp(func(ctx context.Context) {
		eventType, eventStr, err := ParseFeiShuEventData(data)
		if err != nil {
			t.Error(err)
		}
		t.Log(eventType)
		t.Log(eventStr)
	}))
}

func TestFsIssueUrgeHandler_Handle(t *testing.T) {
	data := `{
    "open_id":"ou_2bbc27672c6ce4f0c2c73856f73ccbdb",
    "open_message_id":"om_a2033d568c3d05682a1182164c450241",
    "tenant_key":"12a813fb180f1758",
    "token":"c-451c6b80901cdc04cf1e221aa663345b5edd645e",
    "action":{
        "tag":"select_static",
        "option":"忽略",
        "timezone":"",
        "value":{
            "issueId":10084963,
            "appId":"",
            "projectId":0,
            "tableId":"",
            "orgId":2373,
            "userId":29611,
            "statusType":0,
            "cardType":"FsCardTypeIssueUrge",
            "action":"",
            "fsCardValueUrgeIssueBeUrgedUserId":0,
            "fsCardValueUrgeIssueUrgeText":"请问这个工作进展的怎么样了？",
            "fsCardValueAuditIssueResult":0,
            "fsCardValueAuditIssueStarterId":0,
            "fsCardValueIsGroupChat":false,
            "fsCardValueGroupChatId":"",
            "cardTitle":""
        }
    }
}`
	convey.Convey("test", t, test.StartUp(func(ctx context.Context) {
		cardReq := CardReq{}
		json.FromJson(data, &cardReq)
		h := FsIssueUrgeHandler{}
		h.Handle(cardReq)
	}))
}

func TestFsIssueAuditHandler_Handle(t *testing.T) {
	data := `{
    "open_id": "ou_2bbc27672c6ce4f0c2c73856f73ccbdb",
    "open_message_id": "om_80b720ece8e243ac1e2f47b511f9550a",
    "tenant_key": "12a813fb180f1758",
    "token": "c-6dbc76b1064ca7d62985361debd1b03bf326a154",
    "action": {
        "tag": "button",
        "option": "",
        "timezone": "",
        "value": {
            "issueId": 10093389,
            "appId": "1537674994417500161",
            "projectId": 0,
            "tableId": "1537674995105271808",
            "orgId": 2373,
            "userId": 29610,
            "statusType": 0,
            "cardType": "FsCardTypeIssueAudit",
            "action": "",
            "fsCardValueUrgeIssueBeUrgedUserId": 0,
            "fsCardValueUrgeIssueUrgeText": "",
            "fsCardValueAuditIssueResult": 3,
            "fsCardValueAuditIssueStarterId": 29610,
            "fsCardValueIsGroupChat": false,
            "fsCardValueGroupChatId": "",
            "cardTitle": ""
        }
    }
}`
	convey.Convey("test", t, test.StartUp(func(ctx context.Context) {
		cardReq := CardReq{}
		json.FromJson(data, &cardReq)
		h := FsIssueAuditHandler{}
		h.Handle(cardReq)
	}))
}
