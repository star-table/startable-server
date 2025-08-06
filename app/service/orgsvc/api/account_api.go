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

func SetPassword(c *gin.Context) {
	var req orgvo.SetPasswordReqVo
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Error("SetPassword bind request failed", err)
		response := vo.VoidErr{Err: vo.NewErr(errs.BuildSystemErrorInfo(errs.ReqParamsValidateError, err))}
		c.JSON(http.StatusOK, response)
		return
	}

	err := orgsvcService.SetPassword(req)
	response := vo.VoidErr{Err: vo.NewErr(err)}
	c.JSON(http.StatusOK, response)
}

func ResetPassword(c *gin.Context) {
	var req orgvo.ResetPasswordReqVo
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Error("ResetPassword bind request failed", err)
		response := vo.VoidErr{Err: vo.NewErr(errs.BuildSystemErrorInfo(errs.ReqParamsValidateError, err))}
		c.JSON(http.StatusOK, response)
		return
	}

	err := orgsvcService.ResetPassword(req)
	response := vo.VoidErr{Err: vo.NewErr(err)}
	c.JSON(http.StatusOK, response)
}

func RetrievePassword(c *gin.Context) {
	var req orgvo.RetrievePasswordReqVo
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Error("RetrievePassword bind request failed", err)
		response := vo.VoidErr{Err: vo.NewErr(errs.BuildSystemErrorInfo(errs.ReqParamsValidateError, err))}
		c.JSON(http.StatusOK, response)
		return
	}

	err := orgsvcService.RetrievePassword(req)
	response := vo.VoidErr{Err: vo.NewErr(err)}
	c.JSON(http.StatusOK, response)
}

// UnbindLoginName 目前情况先不允许解绑，要不会导致本地账户丢失，再也无法使用
func UnbindLoginName(c *gin.Context) {
	var req orgvo.UnbindLoginNameReqVo
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Error("UnbindLoginName bind request failed", err)
		response := vo.VoidErr{Err: vo.NewErr(errs.BuildSystemErrorInfo(errs.ReqParamsValidateError, err))}
		c.JSON(http.StatusOK, response)
		return
	}

	err := orgsvcService.UnbindLoginName(req)
	response := vo.VoidErr{Err: vo.NewErr(err)}
	c.JSON(http.StatusOK, response)
}

func BindLoginName(c *gin.Context) {
	var req orgvo.BindLoginNameReqVo
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Error("BindLoginName bind request failed", err)
		response := vo.VoidErr{Err: vo.NewErr(errs.BuildSystemErrorInfo(errs.ReqParamsValidateError, err))}
		c.JSON(http.StatusOK, response)
		return
	}

	err := orgsvcService.BindLoginName(req)
	response := vo.VoidErr{Err: vo.NewErr(err)}
	c.JSON(http.StatusOK, response)
}

func CheckLoginName(c *gin.Context) {
	var req orgvo.CheckLoginNameReqVo
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Error("CheckLoginName bind request failed", err)
		response := vo.VoidErr{Err: vo.NewErr(errs.BuildSystemErrorInfo(errs.ReqParamsValidateError, err))}
		c.JSON(http.StatusOK, response)
		return
	}

	err := orgsvcService.CheckLoginName(req)
	response := vo.VoidErr{Err: vo.NewErr(err)}
	c.JSON(http.StatusOK, response)
}

func VerifyOldName(c *gin.Context) {
	var req orgvo.UnbindLoginNameReqVo
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Error("VerifyOldName bind request failed", err)
		response := vo.VoidErr{Err: vo.NewErr(errs.BuildSystemErrorInfo(errs.ReqParamsValidateError, err))}
		c.JSON(http.StatusOK, response)
		return
	}

	err := orgsvcService.VerifyOldName(req)
	response := vo.VoidErr{Err: vo.NewErr(err)}
	c.JSON(http.StatusOK, response)
}

func ChangeLoginName(c *gin.Context) {
	var req orgvo.BindLoginNameReqVo
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Error("ChangeLoginName bind request failed", err)
		response := vo.VoidErr{Err: vo.NewErr(errs.BuildSystemErrorInfo(errs.ReqParamsValidateError, err))}
		c.JSON(http.StatusOK, response)
		return
	}

	err := orgsvcService.ChangeLoginName(req)
	response := vo.VoidErr{Err: vo.NewErr(err)}
	c.JSON(http.StatusOK, response)
}

func DisbandThirdAccount(c *gin.Context) {
	var req vo.CommonReqVo
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Error("DisbandThirdAccount bind request failed", err)
		response := vo.VoidErr{Err: vo.NewErr(errs.BuildSystemErrorInfo(errs.ReqParamsValidateError, err))}
		c.JSON(http.StatusOK, response)
		return
	}

	err := orgsvcService.DisbandThirdAccount(req.OrgId, req.UserId, req.SourceChannel)
	response := vo.VoidErr{Err: vo.NewErr(err)}
	c.JSON(http.StatusOK, response)
}

func ThirdAccountBindList(c *gin.Context) {
	var req vo.CommonReqVo
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Error("ThirdAccountBindList bind request failed", err)
		response := orgvo.ThirdAccountListResp{Err: vo.NewErr(errs.BuildSystemErrorInfo(errs.ReqParamsValidateError, err))}
		c.JSON(http.StatusOK, response)
		return
	}

	res, err := orgsvcService.ThirdAccountBindList(req.OrgId, req.UserId)
	response := orgvo.ThirdAccountListResp{
		Err:  vo.NewErr(err),
		Data: res,
	}
	c.JSON(http.StatusOK, response)
}
