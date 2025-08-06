package api

import (
	"net/http"
	"github.com/gin-gonic/gin"
	orgsvcService "github.com/star-table/startable-server/app/service/orgsvc/service"
	"github.com/star-table/startable-server/common/model/vo"
	"github.com/star-table/startable-server/common/model/vo/orgvo"
	"github.com/star-table/startable-server/common/core/errs"
	"github.com/star-table/startable-server/common/core/logger"
)

func GetWeiXinJsAPISign(c *gin.Context) {
	var input orgvo.JsAPISignReq
	if err := c.ShouldBindJSON(&input); err != nil {
		logger.Error("GetWeiXinJsAPISign bind request failed", err)
		response := orgvo.GetJsAPISignRespVo{Err: vo.NewErr(errs.BuildSystemErrorInfo(errs.ReqParamsValidateError, err))}
		c.JSON(http.StatusOK, response)
		return
	}

	res, err := orgsvcService.GetWeiXinJsAPISign(input)
	response := orgvo.GetJsAPISignRespVo{Err: vo.NewErr(err), GetJsAPISign: res}
	c.JSON(http.StatusOK, response)
}

func GetWeiXinRegisterUrl(c *gin.Context) {
	var input orgvo.GetRegisterUrlReq
	if err := c.ShouldBindJSON(&input); err != nil {
		logger.Error("GetWeiXinRegisterUrl bind request failed", err)
		response := orgvo.GetRegisterUrlResp{Err: vo.NewErr(errs.BuildSystemErrorInfo(errs.ReqParamsValidateError, err))}
		c.JSON(http.StatusOK, response)
		return
	}

	res, err := orgsvcService.GetWeiXinRegisterUrl()
	response := orgvo.GetRegisterUrlResp{Err: vo.NewErr(err), Data: res}
	c.JSON(http.StatusOK, response)
}

func GetWeiXinInstallUrl(c *gin.Context) {
	var input orgvo.GetRegisterUrlReq
	if err := c.ShouldBindJSON(&input); err != nil {
		logger.Error("GetWeiXinInstallUrl bind request failed", err)
		response := orgvo.GetRegisterUrlResp{Err: vo.NewErr(errs.BuildSystemErrorInfo(errs.ReqParamsValidateError, err))}
		c.JSON(http.StatusOK, response)
		return
	}

	res, err := orgsvcService.GetWeiXinInstallUrl()
	response := orgvo.GetRegisterUrlResp{Err: vo.NewErr(err), Data: res}
	c.JSON(http.StatusOK, response)
}

func AuthWeiXinCode(c *gin.Context) {
	var input orgvo.AuthWeiXinCodeReqVo
	if err := c.ShouldBindJSON(&input); err != nil {
		logger.Error("AuthWeiXinCode bind request failed", err)
		response := orgvo.AuthDingCodeRespVo{Err: vo.NewErr(errs.BuildSystemErrorInfo(errs.ReqParamsValidateError, err))}
		c.JSON(http.StatusOK, response)
		return
	}

	res, err := orgsvcService.AuthWeiXinCode(input)
	response := orgvo.AuthDingCodeRespVo{
		Err:  vo.NewErr(err),
		Data: res,
	}
	c.JSON(http.StatusOK, response)
}

// PersonWeiXinLogin 个人微信登陆
func PersonWeiXinLogin(c *gin.Context) {
	var input orgvo.PersonWeiXinLoginReqVo
	if err := c.ShouldBindJSON(&input); err != nil {
		logger.Error("PersonWeiXinLogin bind request failed", err)
		response := orgvo.PersonWeiXinLoginRespVo{Err: vo.NewErr(errs.BuildSystemErrorInfo(errs.ReqParamsValidateError, err))}
		c.JSON(http.StatusOK, response)
		return
	}

	res, err := orgsvcService.PersonWeiXinLogin(&input)
	response := orgvo.PersonWeiXinLoginRespVo{
		Err:  vo.NewErr(err),
		Data: res,
	}
	c.JSON(http.StatusOK, response)
}

func PersonWeiXinBindExistAccount(c *gin.Context) {
	var input orgvo.PersonWeiXinBindExistAccountReq
	if err := c.ShouldBindJSON(&input); err != nil {
		logger.Error("PersonWeiXinBindExistAccount bind request failed", err)
		response := vo.CommonRespVo{Err: vo.NewErr(errs.BuildSystemErrorInfo(errs.ReqParamsValidateError, err))}
		c.JSON(http.StatusOK, response)
		return
	}

	err := orgsvcService.PersonWeiXinBindExistAccount(&input)
	response := vo.CommonRespVo{Err: vo.NewErr(err)}
	c.JSON(http.StatusOK, response)
}

func PersonWeiXinBind(c *gin.Context) {
	var input orgvo.PersonWeiXinBindReq
	if err := c.ShouldBindJSON(&input); err != nil {
		logger.Error("PersonWeiXinBind bind request failed", err)
		response := orgvo.PersonWeiXinBindResp{Err: vo.NewErr(errs.BuildSystemErrorInfo(errs.ReqParamsValidateError, err))}
		c.JSON(http.StatusOK, response)
		return
	}

	res, err := orgsvcService.PersonWeiXinBind(input.Data)
	response := orgvo.PersonWeiXinBindResp{Err: vo.NewErr(err), Data: res}
	c.JSON(http.StatusOK, response)
}

func PersonWeiXinQrCode(c *gin.Context) {
	var input orgvo.PersonWeiXinQrCodeReqVo
	if err := c.ShouldBindJSON(&input); err != nil {
		logger.Error("PersonWeiXinQrCode bind request failed", err)
		response := orgvo.PersonWeiXinQrCodeRespVo{Err: vo.NewErr(errs.BuildSystemErrorInfo(errs.ReqParamsValidateError, err))}
		c.JSON(http.StatusOK, response)
		return
	}

	res, err := orgsvcService.PersonWeiXinQrCode(&input)
	response := orgvo.PersonWeiXinQrCodeRespVo{Err: vo.NewErr(err), Data: res}
	c.JSON(http.StatusOK, response)
}

func PersonWeiXinQrCodeScan(c *gin.Context) {
	var input orgvo.QrCodeScanReq
	if err := c.ShouldBindJSON(&input); err != nil {
		logger.Error("PersonWeiXinQrCodeScan bind request failed", err)
		response := vo.CommonRespVo{Err: vo.NewErr(errs.BuildSystemErrorInfo(errs.ReqParamsValidateError, err))}
		c.JSON(http.StatusOK, response)
		return
	}

	err := orgsvcService.PersonWeiXinQrCodeScan(&input)
	response := vo.CommonRespVo{Err: vo.NewErr(err)}
	c.JSON(http.StatusOK, response)
}

func CheckPersonWeiXinQrCode(c *gin.Context) {
	var input orgvo.CheckQrCodeScanReq
	if err := c.ShouldBindJSON(&input); err != nil {
		logger.Error("CheckPersonWeiXinQrCode bind request failed", err)
		response := orgvo.CheckQrCodeScanResp{Err: vo.NewErr(errs.BuildSystemErrorInfo(errs.ReqParamsValidateError, err))}
		c.JSON(http.StatusOK, response)
		return
	}

	res, err := orgsvcService.CheckPersonWeiXinQrCode(&input)
	response := orgvo.CheckQrCodeScanResp{Err: vo.NewErr(err), Data: res}
	c.JSON(http.StatusOK, response)
}

func CheckMobileHasBind(c *gin.Context) {
	var input orgvo.CheckMobileHasBindReq
	if err := c.ShouldBindJSON(&input); err != nil {
		logger.Error("CheckMobileHasBind bind request failed", err)
		response := orgvo.CheckMobileHasBindResp{Err: vo.NewErr(errs.BuildSystemErrorInfo(errs.ReqParamsValidateError, err))}
		c.JSON(http.StatusOK, response)
		return
	}

	res, err := orgsvcService.CheckMobileHasBind(input.Input)
	response := orgvo.CheckMobileHasBindResp{Err: vo.NewErr(err), Data: res}
	c.JSON(http.StatusOK, response)
}

func UnbindThirdAccount(c *gin.Context) {
	var input orgvo.UnbindAccountReq
	if err := c.ShouldBindJSON(&input); err != nil {
		logger.Error("UnbindThirdAccount bind request failed", err)
		response := vo.CommonRespVo{Err: vo.NewErr(errs.BuildSystemErrorInfo(errs.ReqParamsValidateError, err))}
		c.JSON(http.StatusOK, response)
		return
	}

	err := orgsvcService.UnbindThirdAccount(input)
	response := vo.CommonRespVo{Err: vo.NewErr(err)}
	c.JSON(http.StatusOK, response)
}
