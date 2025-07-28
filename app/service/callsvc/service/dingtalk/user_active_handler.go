package callsvc

import (
	"fmt"

	"github.com/star-table/startable-server/common/core/util/json"
)

type UserActiveHandler struct {
	base.CallBackBase
}

type userActiveEvent struct {
	SyncAction string     `json:"syncAction"`
	UserId     string     `json:"userid"`
	Department []string   `json:"department"`
	Role       []userRole `json:"role"`
}

// Handle 好像不知道有啥用。。
func (u UserActiveHandler) Handle(data req_vo.DingEventBizData) error {
	msg := &userActiveEvent{}
	json.FromJson(data.BizData, msg)

	fmt.Println("callback -> ", data)

	return nil
}
