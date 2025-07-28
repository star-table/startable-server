package callsvc

import "github.com/go-laoji/wecom-go-sdk/pkg/svr/logic"

type PayHandle struct {
	logic.BizEvent
	ServiceCorpId string `xml:"ServiceCorpId"`
	AuthCorpId    string `xml:"AuthCorpId"`
	OrderId       string `xml:"OrderId"`
	BuyerUserId   string `xml:"BuyerUserId"`
	OrderStatus   string `xml:"OrderStatus"`
}

type OpenOrder struct {
	logic.BizEvent
	PaidCorpId string `xml:"PaidCorpId"`
	OrderId    string `xml:"OrderId"`
	OperatorId string `xml:"OperatorId"`
}

type ChangeOrder struct {
	logic.BizEvent
	PaidCorpId string `xml:"PaidCorpId"`
	OldOrderId string `xml:"OldOrderId"`
	NewOrderId string `xml:"NewOrderId"`
}

type PayForAppSuccess struct {
	logic.BizEvent
	PaidCorpId string `xml:"PaidCorpId"`
	OrderId    string `xml:"OrderId"`
}
