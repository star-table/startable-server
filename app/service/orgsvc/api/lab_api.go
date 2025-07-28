package orgsvc

import (
	"github.com/star-table/startable-server/app/service"
	"github.com/star-table/startable-server/common/model/vo"
	"github.com/star-table/startable-server/common/model/vo/orgvo"
)

func (PostGreeter) SetLabConfig(reqVo orgvo.SetLabReqVo) orgvo.SetLabRespVo {
	resp, err := service.SetLabConfig(reqVo.OrgId, reqVo.UserId, reqVo.Input)
	return orgvo.SetLabRespVo{
		Err:  vo.NewErr(err),
		Data: resp,
	}
}

func (GetGreeter) GetLabConfig(reqVo orgvo.GetLabReqVo) orgvo.GetLabRespVo {
	resp, err := service.GetLabConfig(reqVo)
	return orgvo.GetLabRespVo{
		Err:  vo.NewErr(err),
		Data: resp,
	}
}
