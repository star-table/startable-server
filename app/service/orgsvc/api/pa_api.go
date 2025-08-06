package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/star-table/startable-server/common/core/errs"
	"github.com/star-table/startable-server/common/core/logger"
	"github.com/star-table/startable-server/common/model/vo"
	"github.com/star-table/startable-server/common/model/vo/orgvo"
	orgsvcService "github.com/star-table/startable-server/app/service/orgsvc/service"
)

func PAReport(c *gin.Context) {
	var req orgvo.PAReportMsgReqVo
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusOK, vo.VoidErr{Err: vo.NewErr(errs.BuildSystemErrorInfo(errs.ReqParamsValidateError, err))})
		return
	}

	err := orgsvcService.PAReport(req.Body)
	if err != nil {
		logger.Error("PAReport error", logger.ErrorField(err))
	}
	c.JSON(http.StatusOK, vo.VoidErr{Err: vo.NewErr(err)})
}
