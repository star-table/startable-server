package orgsvc

import (
	"github.com/star-table/startable-server/app/service"
	"github.com/star-table/startable-server/common/model/vo"
	"github.com/star-table/startable-server/common/model/vo/orgvo"
)

func (GetGreeter) GetWeiXinJsAPISign(input orgvo.JsAPISignReq) orgvo.GetJsAPISignRespVo {
	res, err := service.GetWeiXinJsAPISign(input)
	return orgvo.GetJsAPISignRespVo{Err: vo.NewErr(err), GetJsAPISign: res}
}

func (GetGreeter) GetWeiXinRegisterUrl(input orgvo.GetRegisterUrlReq) orgvo.GetRegisterUrlResp {
	res, err := service.GetWeiXinRegisterUrl()
	return orgvo.GetRegisterUrlResp{Err: vo.NewErr(err), Data: res}
}

func (GetGreeter) GetWeiXinInstallUrl(input orgvo.GetRegisterUrlReq) orgvo.GetRegisterUrlResp {
	res, err := service.GetWeiXinInstallUrl()
	return orgvo.GetRegisterUrlResp{Err: vo.NewErr(err), Data: res}
}

func (PostGreeter) AuthWeiXinCode(input orgvo.AuthWeiXinCodeReqVo) orgvo.AuthDingCodeRespVo {
	res, err := service.AuthWeiXinCode(input)
	return orgvo.AuthDingCodeRespVo{
		Err:  vo.NewErr(err),
		Data: res,
	}
}

// PersonWeiXinLogin 个人微信登陆
func (PostGreeter) PersonWeiXinLogin(input orgvo.PersonWeiXinLoginReqVo) orgvo.PersonWeiXinLoginRespVo {
	res, err := service.PersonWeiXinLogin(&input)
	return orgvo.PersonWeiXinLoginRespVo{
		Err:  vo.NewErr(err),
		Data: res,
	}
}

func (PostGreeter) PersonWeiXinBindExistAccount(input orgvo.PersonWeiXinBindExistAccountReq) vo.CommonRespVo {
	err := service.PersonWeiXinBindExistAccount(&input)
	return vo.CommonRespVo{Err: vo.NewErr(err)}
}

func (PostGreeter) PersonWeiXinBind(input orgvo.PersonWeiXinBindReq) orgvo.PersonWeiXinBindResp {
	res, err := service.PersonWeiXinBind(input.Data)
	return orgvo.PersonWeiXinBindResp{Err: vo.NewErr(err), Data: res}
}

func (GetGreeter) PersonWeiXinQrCode(input orgvo.PersonWeiXinQrCodeReqVo) orgvo.PersonWeiXinQrCodeRespVo {
	res, err := service.PersonWeiXinQrCode(&input)
	return orgvo.PersonWeiXinQrCodeRespVo{Err: vo.NewErr(err), Data: res}
}

func (GetGreeter) PersonWeiXinQrCodeScan(input orgvo.QrCodeScanReq) vo.CommonRespVo {
	err := service.PersonWeiXinQrCodeScan(&input)
	return vo.CommonRespVo{Err: vo.NewErr(err)}
}

func (GetGreeter) CheckPersonWeiXinQrCode(input orgvo.CheckQrCodeScanReq) orgvo.CheckQrCodeScanResp {
	res, err := service.CheckPersonWeiXinQrCode(&input)
	return orgvo.CheckQrCodeScanResp{Err: vo.NewErr(err), Data: res}
}

func (PostGreeter) CheckMobileHasBind(input orgvo.CheckMobileHasBindReq) orgvo.CheckMobileHasBindResp {
	res, err := service.CheckMobileHasBind(input.Input)
	return orgvo.CheckMobileHasBindResp{Err: vo.NewErr(err), Data: res}
}

func (PostGreeter) UnbindThirdAccount(input orgvo.UnbindAccountReq) vo.CommonRespVo {
	err := service.UnbindThirdAccount(input)
	return vo.CommonRespVo{Err: vo.NewErr(err)}
}
