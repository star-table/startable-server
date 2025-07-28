package orgsvc

import (
	"github.com/star-table/startable-server/app/service"
	"github.com/star-table/startable-server/common/model/vo"
	"github.com/star-table/startable-server/common/model/vo/orgvo"
)

func (GetGreeter) OpenAPIAuth(req orgvo.OpenAPIAuthReq) orgvo.OpenAPIAuthResp {
	data, err := service.OpenAPIAuth(req)
	return orgvo.OpenAPIAuthResp{Data: data, Err: vo.NewErr(err)}
}

func (PostGreeter) GetAppTicket(req orgvo.GetAppTicketReq) orgvo.GetAppTicketResp {
	data, err := service.GetAppTicket(req)
	return orgvo.GetAppTicketResp{
		Err:  vo.NewErr(err),
		Data: data,
	}
}
