package api

import (
	"net/http"
	"strings"
	"github.com/gin-gonic/gin"
	orgsvcService "github.com/star-table/startable-server/app/service/orgsvc/service"
	"github.com/star-table/startable-server/common/core/errs"
	"github.com/star-table/startable-server/common/core/util/format"
	"github.com/star-table/startable-server/common/model/vo"
	"github.com/star-table/startable-server/common/model/vo/orgvo"
)

// UserLogin 用户登录
func UserLogin(c *gin.Context) {
	var req orgvo.UserLoginReqVo
	if err := c.ShouldBindJSON(&req); err != nil {
		response := orgvo.UserSMSLoginRespVo{Err: vo.NewErr(errs.BuildSystemErrorInfo(errs.ParamError, err))}
		c.JSON(http.StatusOK, response)
		return
	}
	
	if req.UserLoginReq.Name != nil && *req.UserLoginReq.Name != "" {
		creatorName := strings.TrimSpace(*req.UserLoginReq.Name)
		isNameRight := format.VerifyUserNameFormat(creatorName)
		if !isNameRight {
			response := orgvo.UserSMSLoginRespVo{Err: vo.NewErr(errs.UserNameLenError), Data: nil}
			c.JSON(http.StatusOK, response)
			return
		}
		req.UserLoginReq.Name = &creatorName
	}

	res, err := orgsvcService.UserLogin(req.UserLoginReq)
	if err != nil {
		response := orgvo.UserSMSLoginRespVo{Err: vo.NewErr(errs.BuildSystemErrorInfo(errs.ServerError, err))}
		c.JSON(http.StatusOK, response)
		return
	}
	response := orgvo.UserSMSLoginRespVo{Err: vo.NewErr(nil), Data: res}
	c.JSON(http.StatusOK, response)
}

// UserQuit 用户退出
func UserQuit(c *gin.Context) {
	var req orgvo.UserQuitReqVo
	if err := c.ShouldBindJSON(&req); err != nil {
		response := vo.VoidErr{Err: vo.NewErr(errs.BuildSystemErrorInfo(errs.ParamError, err))}
		c.JSON(http.StatusOK, response)
		return
	}
	
	err := orgsvcService.UserQuit(req)
	if err != nil {
		response := vo.VoidErr{Err: vo.NewErr(errs.BuildSystemErrorInfo(errs.ServerError, err))}
		c.JSON(http.StatusOK, response)
		return
	}
	response := vo.VoidErr{Err: vo.NewErr(nil)}
	c.JSON(http.StatusOK, response)
}

// JoinOrgByInviteCode 通过邀请码加入组织
func JoinOrgByInviteCode(c *gin.Context) {
	var req orgvo.JoinOrgByInviteCodeReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response := vo.VoidErr{Err: vo.NewErr(errs.BuildSystemErrorInfo(errs.ParamError, err))}
		c.JSON(http.StatusOK, response)
		return
	}
	
	err := orgsvcService.JoinOrgByInviteCode(req)
	if err != nil {
		response := vo.VoidErr{Err: vo.NewErr(errs.BuildSystemErrorInfo(errs.ServerError, err))}
		c.JSON(http.StatusOK, response)
		return
	}
	response := vo.VoidErr{Err: vo.NewErr(nil)}
	c.JSON(http.StatusOK, response)
}

// ExchangeShareToken 交换分享令牌
func ExchangeShareToken(c *gin.Context) {
	var req orgvo.ExchangeShareTokenReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response := orgvo.UserSMSLoginRespVo{Err: vo.NewErr(errs.BuildSystemErrorInfo(errs.ParamError, err))}
		c.JSON(http.StatusOK, response)
		return
	}
	
	data, err := orgsvcService.ExchangeShareToken(req)
	if err != nil {
		response := orgvo.UserSMSLoginRespVo{Err: vo.NewErr(errs.BuildSystemErrorInfo(errs.ServerError, err))}
		c.JSON(http.StatusOK, response)
		return
	}
	response := orgvo.UserSMSLoginRespVo{
		Err:  vo.NewErr(nil),
		Data: data,
	}
	c.JSON(http.StatusOK, response)
}
