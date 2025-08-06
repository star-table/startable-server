/*
*
2 * @Author: Nico
3 * @Date: 2020/1/31 11:17
4
*/
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

func UserRegister(c *gin.Context) {
	var req orgvo.UserRegisterReqVo
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Error("UserRegister bind request failed", err)
		response := orgvo.UserRegisterRespVo{Err: vo.NewErr(errs.BuildSystemErrorInfo(errs.ReqParamsValidateError, err))}
		c.JSON(http.StatusOK, response)
		return
	}

	res, err := orgsvcService.UserRegister(req)
	response := orgvo.UserRegisterRespVo{Data: res, Err: vo.NewErr(err)}
	c.JSON(http.StatusOK, response)
}
