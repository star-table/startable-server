package callsvc

import (
	"github.com/star-table/startable-server/common/core/util/json"
)

type MarketOrderTrailHandler struct{}

type MarketOrderTrail struct {
	BuyerUnionId string `json:"buyerUnionId"`
	SuiteId      int64  `json:"suiteId"`
	CorpId       string `json:"corpId"`
	FromDate     string `json:"fromDate"`
	EndDate      string `json:"endDate"`
	EventType    string `json:"EventType"`
	AppId        int64  `json:"appId"`
	TryoutType   string `json:"tryoutType"`
	GoodsCode    string `json:"goodsCode"`
	ItemType     string `json:"itemType"`
	TimeStamp    string `json:"TimeStamp"`
}

func (MarketOrderTrailHandler) Handle(data req_vo.DingEventBizData) error {
	log.Infof("ding [MarketOrderTrailHandler] handle:%v", data)
	msg := MarketOrderTrail{}
	err := json.FromJson(data.BizData, &msg)
	if err != nil {
		log.Errorf("[MarketOrderTrailHandler] handle json err:%v", err)
		return err
	}

	return nil
}
