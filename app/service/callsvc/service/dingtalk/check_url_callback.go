package callsvc

import (
	"fmt"

	"github.com/star-table/startable-server/common/core/util/json"
)

type CheckUrlCallbackHandler struct{}

type CheckUrlCallbackMsg struct {
	Random       string
	EventType    string
	TestSuiteKey string
}

func (CheckUrlCallbackHandler) Handle(data req_vo.DingEventBizData) error {
	msg := &CheckUrlCallbackMsg{}
	json.FromJson(data.BizData, msg)
	fmt.Println("callback -> ", data)
	return nil
}
