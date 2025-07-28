package msgsvc

import (
	"github.com/star-table/startable-server/app/service"
	"github.com/star-table/startable-server/common/model/vo"
	"github.com/star-table/startable-server/common/model/vo/msgvo"
)

func (PostGreeter) SendLoginSMS(req msgvo.SendLoginSMSReqVo) vo.VoidErr {
	err := service.SendLoginSMS(req)
	return vo.VoidErr{Err: vo.NewErr(err)}
}

func (PostGreeter) SendMail(req msgvo.SendMailReqVo) vo.VoidErr {
	err := service.SendMail(req)
	return vo.VoidErr{Err: vo.NewErr(err)}
}

func (PostGreeter) SendSMS(req msgvo.SendSMSReqVo) vo.VoidErr {
	err := service.SendSMS(req)
	return vo.VoidErr{Err: vo.NewErr(err)}
}
