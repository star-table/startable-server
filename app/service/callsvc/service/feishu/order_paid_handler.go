package callsvc

import (
	"strconv"
	"time"

	"github.com/star-table/startable-server/app/facade/orderfacade"
	"github.com/star-table/startable-server/common/core/consts"
	"github.com/star-table/startable-server/common/core/errs"
	"github.com/star-table/startable-server/common/core/util/json"
	"github.com/star-table/startable-server/common/model/bo"
	"github.com/star-table/startable-server/common/model/vo/ordervo"
)

type OrderPaidHandler struct{}

type OrderPaidReq struct {
	Ts    string `json:"ts"`
	Uuid  string `json:"uuid"`
	Token string `json:"token"`
	Type  string `json:"type"`

	Event OrderPaidReqData `json:"event"`
}

type OrderPaidReqData struct {
	Type          string `json:"type"`
	AppId         string `json:"app_id"`
	OrderId       string `json:"order_id"`
	PricePlanId   string `json:"price_plan_id"`
	PricePlanType string `json:"price_plan_type"`
	Seats         int    `json:"seats"`
	BuyCount      int    `json:"buy_count"`
	CreateTime    string `json:"create_time"`
	PayTime       string `json:"pay_time"`
	BuyType       string `json:"buy_type"`
	SrcOrderId    string `json:"src_order_id"`
	OrderPayPrice int64  `json:"order_pay_price"`
	TenantKey     string `json:"tenant_key"`
}

// 飞书处理
func (OrderPaidHandler) Handle(data string) (string, errs.SystemErrorInfo) {
	req := &OrderPaidReq{}
	_ = json.FromJson(data, req)

	log.Infof("飞书应用购买通知 %s", data)
	eventData := req.Event

	orderFsBo := bo.OrderFsBo{
		OrderId:       eventData.OrderId,
		PricePlanId:   eventData.PricePlanId,
		PricePlanType: eventData.PricePlanType,
		Seats:         eventData.Seats,
		BuyCount:      eventData.BuyCount,
		Status:        consts.FsOrderStatusNormal,
		BuyType:       eventData.BuyType,
		SrcOrderId:    eventData.SrcOrderId,
		DstOrderId:    "",
		OrderPayPrice: eventData.OrderPayPrice,
		TenantKey:     eventData.TenantKey,
	}
	if eventData.PayTime != "" {
		payTime, err := strconv.ParseInt(eventData.PayTime, 10, 64)
		if err != nil {
			log.Error(err)
		}
		orderFsBo.PaidTime = time.Unix(payTime, 0)
	}
	if eventData.CreateTime != "" {
		createTime, err := strconv.ParseInt(eventData.CreateTime, 10, 64)
		if err != nil {
			log.Error(err)
		}
		orderFsBo.CreateTime = time.Unix(createTime, 0)
	}

	resp := orderfacade.AddFsOrder(ordervo.AddFsOrderReq{Data: orderFsBo})
	if resp.Failure() {
		log.Error(resp.Error())
		return "error", resp.Error()
	}

	errCache := base.DelPayRangeInfoCache(0, eventData.TenantKey)
	if errCache != nil {
		log.Errorf("[OrderPaidHandler] DelPayRangeInfoCache err:%v, corpId:%v", errCache, eventData.TenantKey)
	}

	return "ok", nil
}
