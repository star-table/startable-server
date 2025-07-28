package msgsvc

import (
	"github.com/star-table/startable-server/app/service"
	"github.com/star-table/startable-server/common/model/vo"
	"github.com/star-table/startable-server/common/model/vo/msgvo"
)

// FixAddOrderForFeiShu 重跑消费订单异常的订单消息
func (GetGreeter) FixAddOrderForFeiShu(msg msgvo.FixAddOrderForFeiShuReqVo) vo.VoidErr {
	err := service.FixAddOrderForFeiShu(msg.Input)
	return vo.VoidErr{Err: vo.NewErr(err)}
}
