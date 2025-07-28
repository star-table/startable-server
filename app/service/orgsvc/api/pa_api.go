package orgsvc

import (
	"github.com/star-table/startable-server/app/service"
	"github.com/star-table/startable-server/common/model/vo"
	"github.com/star-table/startable-server/common/model/vo/orgvo"
)

func (PostGreeter) PAReport(req orgvo.PAReportMsgReqVo) vo.VoidErr {
	err := service.PAReport(req.Body)
	return vo.VoidErr{Err: vo.NewErr(err)}
}
