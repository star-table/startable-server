package api

import (
	"net/http"
	"github.com/gin-gonic/gin"
	msgsvcService "github.com/star-table/startable-server/app/service/msgsvc/service"
	"github.com/star-table/startable-server/common/core/errs"
	"github.com/star-table/startable-server/common/model/vo"
	"github.com/star-table/startable-server/common/model/vo/msgvo"
)

// SendLoginSMS 发送登录短信
func SendLoginSMS(c *gin.Context) {
	var req msgvo.SendLoginSMSReqVo
	if err := c.ShouldBindJSON(&req); err != nil {
		response := vo.VoidErr{Err: vo.NewErr(errs.BuildSystemErrorInfo(errs.ParamError, err))}
		c.JSON(http.StatusOK, response)
		return
	}
	
	err := msgsvcService.SendLoginSMS(req)
	if err != nil {
		response := vo.VoidErr{Err: vo.NewErr(errs.BuildSystemErrorInfo(errs.ServerError, err))}
		c.JSON(http.StatusOK, response)
		return
	}
	response := vo.VoidErr{Err: vo.NewErr(nil)}
	c.JSON(http.StatusOK, response)
}

// SendMail 发送邮件
func SendMail(c *gin.Context) {
	var req msgvo.SendMailReqVo
	if err := c.ShouldBindJSON(&req); err != nil {
		response := vo.VoidErr{Err: vo.NewErr(errs.BuildSystemErrorInfo(errs.ParamError, err))}
		c.JSON(http.StatusOK, response)
		return
	}
	
	err := msgsvcService.SendMail(req)
	if err != nil {
		response := vo.VoidErr{Err: vo.NewErr(errs.BuildSystemErrorInfo(errs.ServerError, err))}
		c.JSON(http.StatusOK, response)
		return
	}
	response := vo.VoidErr{Err: vo.NewErr(nil)}
	c.JSON(http.StatusOK, response)
}

// SendSMS 发送短信
func SendSMS(c *gin.Context) {
	var req msgvo.SendSMSReqVo
	if err := c.ShouldBindJSON(&req); err != nil {
		response := vo.VoidErr{Err: vo.NewErr(errs.BuildSystemErrorInfo(errs.ParamError, err))}
		c.JSON(http.StatusOK, response)
		return
	}
	
	err := msgsvcService.SendSMS(req)
	if err != nil {
		response := vo.VoidErr{Err: vo.NewErr(errs.BuildSystemErrorInfo(errs.ServerError, err))}
		c.JSON(http.StatusOK, response)
		return
	}
	response := vo.VoidErr{Err: vo.NewErr(nil)}
	c.JSON(http.StatusOK, response)
}
