package api

import (
	"net/http"
	"github.com/gin-gonic/gin"
	orgsvcService "github.com/star-table/startable-server/app/service/orgsvc/service"
	"github.com/star-table/startable-server/app/service/orgsvc/domain"
	"github.com/star-table/startable-server/common/core/consts"
	"github.com/star-table/startable-server/common/core/util/slice"
	"github.com/star-table/startable-server/common/model/vo"
	"github.com/star-table/startable-server/common/model/vo/orgvo"
	"github.com/star-table/startable-server/common/core/errs"
	"github.com/star-table/startable-server/common/core/logger"
)

func SendSMSLoginCode(c *gin.Context) {
	var req orgvo.SendSMSLoginCodeReqVo
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Error("SendSMSLoginCode bind request failed", err)
		response := vo.VoidErr{Err: vo.NewErr(errs.BuildSystemErrorInfo(errs.ReqParamsValidateError, err))}
		c.JSON(http.StatusOK, response)
		return
	}

	phoneNumber := req.Input.PhoneNumber
	//phone format check
	verifyErr := orgsvcService.VerifyCaptcha(req.Input.CaptchaID, req.Input.CaptchaPassword, req.Input.PhoneNumber, req.Input.YidunValidate)
	if verifyErr != nil {
		response := vo.VoidErr{Err: vo.NewErr(verifyErr)}
		c.JSON(http.StatusOK, response)
		return
	}

	err := orgsvcService.SendSMSLoginCode(phoneNumber)
	response := vo.VoidErr{Err: vo.NewErr(err)}
	c.JSON(http.StatusOK, response)
}

func SendAuthCode(c *gin.Context) {
	var req orgvo.SendAuthCodeReqVo
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Error("SendAuthCode bind request failed", err)
		response := vo.VoidErr{Err: vo.NewErr(errs.BuildSystemErrorInfo(errs.ReqParamsValidateError, err))}
		c.JSON(http.StatusOK, response)
		return
	}

	if ok, _ := slice.Contain([]int{consts.AuthCodeTypeBind, consts.AuthCodeTypeUnBind}, req.Input.AuthType); !ok {
		verifyErr := orgsvcService.VerifyCaptcha(req.Input.CaptchaID, req.Input.CaptchaPassword, req.Input.Address, req.Input.YidunValidate)
		if verifyErr != nil {
			response := vo.VoidErr{Err: vo.NewErr(verifyErr)}
			c.JSON(http.StatusOK, response)
			return
		}
	}

	err := orgsvcService.SendAuthCode(req)
	response := vo.VoidErr{Err: vo.NewErr(err)}
	c.JSON(http.StatusOK, response)
}

func GetPwdLoginCode(c *gin.Context) {
	var req orgvo.GetPwdLoginCodeReqVo
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Error("GetPwdLoginCode bind request failed", err)
		response := orgvo.GetPwdLoginCodeRespVo{Err: vo.NewErr(errs.BuildSystemErrorInfo(errs.ReqParamsValidateError, err))}
		c.JSON(http.StatusOK, response)
		return
	}

	res, err := domain.GetPwdLoginCode(req.CaptchaId)
	response := orgvo.GetPwdLoginCodeRespVo{Err: vo.NewErr(err), CaptchaPassword: res}
	c.JSON(http.StatusOK, response)
}

func SetPwdLoginCode(c *gin.Context) {
	var req orgvo.SetPwdLoginCodeReqVo
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Error("SetPwdLoginCode bind request failed", err)
		response := vo.VoidErr{Err: vo.NewErr(errs.BuildSystemErrorInfo(errs.ReqParamsValidateError, err))}
		c.JSON(http.StatusOK, response)
		return
	}

	response := vo.VoidErr{Err: vo.NewErr(domain.SetPwdLoginCode(req.CaptchaId, req.CaptchaPassword))}
	c.JSON(http.StatusOK, response)
}
