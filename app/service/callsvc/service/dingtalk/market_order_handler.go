package callsvc

import (
	"fmt"
	"time"

	"github.com/star-table/startable-server/common/core/util/json"
	"github.com/star-table/startable-server/common/model/bo"
	"github.com/star-table/startable-server/common/model/vo/ordervo"
)

type MarketOrderHandler struct{}

type MarketOrderCallbackMsg struct {
	OrderId            int64   `json:"orderId"`  // 订单id
	SuiteId            int64   `json:"suiteId"`  //  用户购买第三方企业应用的suiteId
	SuiteKey           string  `json:"suiteKey"` // 用户购买第三方企业应用的suiteKey
	CorpId             string  `json:"corpId"`
	GoodsCode          string  `json:"goodsCode"`          //商品码
	GoodsName          string  `json:"goodsName"`          // 商品名称
	ItemCode           string  `json:"itemCode"`           //规格码
	ItemName           string  `json:"itemName"`           // 规格名称
	SubQuantity        int     `json:"subQuantity"`        // 订购数量
	MaxOfPeople        int     `json:"maxOfPeople"`        // 规格支持最大使用人数
	MinOfPeople        int     `json:"minOfPeople"`        // 规格支持最小使用人数
	PaidTime           int64   `json:"paidtime"`           // 支付时间
	PayFee             int64   `json:"payFee"`             // 实际支付价格，单位：分
	ServiceStartTime   int64   `json:"serviceStartTime"`   // 服务开始时间
	ServiceStopTime    int64   `json:"serviceStopTime"`    // 该订单的服务结束时间  单位毫秒
	OrderCreateSource  string  `json:"orderCreateSource"`  // 订单来源
	DiscountFee        int64   `json:"discountFee"`        // 折扣减免费用（单位：分）
	DisCount           float64 `json:"disCount"`           // 折扣，现值为1.00
	OrderType          string  `json:"orderType"`          // 订单类型, BUY：新购,RENEW：续费,UPGRADE：升级,RENEW_UPGRADE：续费升配,RENEW_DEGRADE：续费降配
	OrderChargeType    string  `json:"orderChargeType"`    // 订单收费类型，FREE：免费开通，TRYOUT：试用开通
	SaleModelType      string  `json:"saleModelType"`      // 售卖模式，取值：CYC_UPGRADE_MEMBER：按周期 + 数量（人数）售卖;CYC_UPGRADE：按周期售卖;QUANTITY：按数量（人数）售卖
	OpenId             string  `json:"openId"`             // 用户在当前开放应用内的唯一标识
	UnionId            string  `json:"unionId"`            // 用户在当前钉钉开放平台账号范围内的唯一标识
	OutTradeNo         string  `json:"outTradeNo"`         // 外部订单号
	ArticleType        string  `json:"articleType"`        // 商品类型，normal：普通商品，image：OXM镜像商品
	AppId              int64   `json:"appId"`              // 应用Id
	PurchaserType      int     `json:"purchaserType"`      // 购买类型, 1：组织购买, 2：个人购买
	AutoChangeFreeItem bool    `json:"autoChangeFreeItem"` // 自动转免费规格
}

// 用户的内购订单支付成功后触发
func (MarketOrderHandler) Handle(data req_vo.DingEventBizData) error {
	msg := &MarketOrderCallbackMsg{}
	errJson := json.FromJson(data.BizData, msg)
	if errJson != nil {
		log.Errorf("[MarketOrderHandler] handle json err:%v", errJson)
		return errJson
	}
	log.Infof("dingTalk order callback -> %v", data)

	//client, err := platform_sdk.GetClient(sdk_const.SourceChannelDingTalk, msg.CorpId)
	//if err != nil {
	//	return errs.BuildSystemErrorInfo(errs.PlatFormOpenApiCallError, err)
	//}
	//dingClient := client.GetOriginClient().(*dingtalk.DingTalk)

	//// 通知钉钉完成订单处理
	//finshedResp, err := dingClient.OrderFinished(msg.OrderId)
	//if err != nil {
	//	log.Errorf("[MarketOrderHandler] handle OrderFinished err:%v", err)
	//	return errs.BuildSystemErrorInfo(errs.DingTalkOpenApiCallError, err)
	//}
	//if finshedResp.Code != 0 {
	//	log.Errorf("[MarketOrderHandler] handle 通知钉钉完成订单处理失败, code:%v, msg:%v", finshedResp.Code, finshedResp.Msg)
	//	return errs.BuildSystemErrorInfo(errs.DingTalkOpenApiCallError, err)
	//}

	// 具体订单的真实支付状态 需要查询用户订单信息接口来判断
	//orderInfo, err := dingClient.GetAppMarketOrderInfo(msg.OrderId)
	//if err != nil {
	//	log.Errorf("[MarketOrderHandler] handle GetOrderInfo err:%v", err)
	//	return errs.BuildSystemErrorInfo(errs.DingTalkOpenApiCallError, err)
	//}
	//log.Infof("[MarketOrderHandler] handle GetAppMarketOrderInfo data:%v", json.ToJsonIgnoreError(orderInfo))
	//if orderInfo.Code != 0 {
	//	log.Errorf("[MarketOrderHandler] handle 查询订单消息失败, code:%v, msg:%v", orderInfo.Code, orderInfo.Msg)
	//	return errs.BuildSystemErrorInfo(errs.DingTalkOpenApiCallError, err)
	//}

	//插入订单
	//orderCreateTime := time.Unix(orderInfo.CreateTimestamp/1000, 0)
	orderBo := bo.OrderDingBo{
		OutOrgId:      msg.CorpId,
		OrderId:       fmt.Sprintf("%d", msg.OrderId),
		GoodsName:     msg.GoodsName,
		GoodsCode:     msg.GoodsCode,
		ItemName:      msg.ItemName,
		ItemCode:      msg.ItemCode,
		Quantity:      msg.SubQuantity,
		OrderPayPrice: msg.PayFee,
		OrderType:     msg.OrderType,
		//Status:          orderInfo.Status,
		OrderSource:     msg.OrderCreateSource,
		SaleModelType:   msg.SaleModelType,
		OrderChargeType: msg.OrderChargeType,
		//OrderCreateTime: orderCreateTime,
	}
	if msg.AutoChangeFreeItem {
		orderBo.AutoChangeFreeItem = 1
	} else {
		orderBo.AutoChangeFreeItem = 2
	}
	if msg.PaidTime != 0 {
		orderBo.PaidTime = time.Unix(msg.PaidTime/1000, 0)
		orderBo.OrderCreateTime = time.Unix(msg.PaidTime/1000, 0)
	}
	if msg.ServiceStopTime != 0 {
		orderBo.EndTime = time.Unix(msg.ServiceStopTime/1000, 0)
	}
	if msg.ServiceStartTime != 0 {
		orderBo.StartTime = time.Unix(msg.ServiceStartTime/1000, 0)
	}

	// 试用订单也会有市场订单事件
	resp := orderfacade.AddDingOrder(ordervo.AddDingOrderReq{Data: orderBo})
	if resp.Failure() {
		log.Error(resp.Error())
		return resp.Error()
	}

	errCache := base.DelPayRangeInfoCache(0, msg.CorpId)
	if errCache != nil {
		log.Errorf("[MarketOrderHandler] DelPayRangeInfoCache err:%v, corpId:%v", errCache, msg.CorpId)
	}

	return nil
}
