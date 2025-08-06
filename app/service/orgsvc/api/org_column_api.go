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

func CreateOrgColumn(c *gin.Context) {
	var req orgvo.CreateOrgColumnReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusOK, orgvo.CreateOrgColumnRespVo{Err: vo.NewErr(errs.BuildSystemErrorInfo(errs.ReqParamsValidateError, err))})
		return
	}

	resp, err := orgsvcService.CreateOrgColumn(req)
	if err != nil {
		logger.Error("CreateOrgColumn error", logger.ErrorField(err))
	}
	c.JSON(http.StatusOK, orgvo.CreateOrgColumnRespVo{
		Err:  vo.NewErr(err),
		Data: resp,
	})
}

func GetOrgColumns(c *gin.Context) {
	var req orgvo.GetOrgColumnsReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusOK, orgvo.GetOrgColumnsRespVo{Err: vo.NewErr(errs.BuildSystemErrorInfo(errs.ReqParamsValidateError, err))})
		return
	}

	resp, err := orgsvcService.GetOrgColumns(req)
	if err != nil {
		logger.Error("GetOrgColumns error", logger.ErrorField(err))
	}
	c.JSON(http.StatusOK, orgvo.GetOrgColumnsRespVo{
		Err:  vo.NewErr(err),
		Data: resp,
	})
}

func DeleteOrgColumn(c *gin.Context) {
	var req orgvo.DeleteOrgColumnReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusOK, orgvo.DeleteOrgColumnRespVo{Err: vo.NewErr(errs.BuildSystemErrorInfo(errs.ReqParamsValidateError, err))})
		return
	}

	resp, err := orgsvcService.DeleteOrgColumn(req)
	if err != nil {
		logger.Error("DeleteOrgColumn error", logger.ErrorField(err))
	}
	c.JSON(http.StatusOK, orgvo.DeleteOrgColumnRespVo{
		Err:  vo.NewErr(err),
		Data: resp,
	})
}
