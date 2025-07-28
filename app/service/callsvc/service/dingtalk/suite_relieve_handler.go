package callsvc

import (
	"fmt"

	"github.com/star-table/startable-server/common/core/util/json"
)

type SuiteRelieveHandler struct{}

type SuiteRelieveCallbackMsg struct {
	SuiteKey   string
	EventType  string
	TimeStamp  string
	AuthCorpId string
}

func (SuiteRelieveHandler) Handle(data req_vo.DingEventBizData) error {
	msg := &SuiteRelieveCallbackMsg{}
	json.FromJson(data.BizData, msg)
	fmt.Println("callback -> ", data)
	return nil
}
