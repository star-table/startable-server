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

func OpenAPIAuth(c *gin.Context) {
	var req orgvo.OpenAPIAuthReq
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Error("OpenAPIAuth bind request failed", err)
		response := orgvo.OpenAPIAuthResp{Err: vo.NewErr(errs.BuildSystemErrorInfo(errs.ReqParamsValidateError, err))}
		c.JSON(http.StatusOK, response)
		return
	}

	data, err := orgsvcService.OpenAPIAuth(req)
	response := orgvo.OpenAPIAuthResp{Data: data, Err: vo.NewErr(err)}
	c.JSON(http.StatusOK, response)
}

func GetAppTicket(c *gin.Context) {
	var req orgvo.GetAppTicketReq
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Error("GetAppTicket bind request failed", err)
		response := orgvo.GetAppTicketResp{Err: vo.NewErr(errs.BuildSystemErrorInfo(errs.ReqParamsValidateError, err))}
		c.JSON(http.StatusOK, response)
		return
	}

	data, err := orgsvcService.GetAppTicket(req)
	response := orgvo.GetAppTicketResp{
		Err:  vo.NewErr(err),
		Data: data,
	}
	c.JSON(http.StatusOK, response)
}
