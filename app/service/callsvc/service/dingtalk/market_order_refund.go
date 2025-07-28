package callsvc

import (
	"github.com/star-table/startable-server/common/core/util/json"
	"github.com/star-table/startable-server/common/model/vo/ordervo"
)

type MarketOrderRefundHandler struct{}

type MarketOrderRefundMsg struct {
	SyncAction string `json:"syncAction"`
	CorpId     string `json:"corpId"`
	OrderId    string `json:"orderId"`
	RefundFee  int64  `json:"refundFee"`
}

func (MarketOrderRefundHandler) Handle(data req_vo.DingEventBizData) error {
	msg := &MarketOrderRefundMsg{}
	errJson := json.FromJson(data.BizData, msg)
	if errJson != nil {
		log.Errorf("[MarketOrderRefundHandler] handle json err:%v", errJson)
		return errJson
	}
	log.Infof("dingTalk MarketOrderRefundHandler callback -> %v", data)

	// 订单退款
	// 删除订单，并重置付费等级
	orderResp := orderfacade.DeleteDingOrder(ordervo.DeleteDingOrderReqVo{
		Data: ordervo.DeleteDingOrderReq{
			OutOrgId: msg.CorpId,
			OrderId:  msg.OrderId,
		},
	})
	if orderResp.Failure() {
		log.Errorf("[MarketOrderRefundHandler] DeleteDingOrder err:%v, outOrgId:%v, orderId:%v",
			orderResp.Error(), msg.CorpId, msg.OrderId)
		return orderResp.Error()
	}

	return nil
}
