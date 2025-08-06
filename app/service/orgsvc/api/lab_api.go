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

func SetLabConfig(c *gin.Context) {
	var reqVo orgvo.SetLabReqVo
	if err := c.ShouldBindJSON(&reqVo); err != nil {
		logger.Error("SetLabConfig bind request failed", err)
		response := orgvo.SetLabRespVo{Err: vo.NewErr(errs.BuildSystemErrorInfo(errs.ReqParamsValidateError, err))}
		c.JSON(http.StatusOK, response)
		return
	}

	resp, err := orgsvcService.SetLabConfig(reqVo.OrgId, reqVo.UserId, reqVo.Input)
	response := orgvo.SetLabRespVo{
		Err:  vo.NewErr(err),
		Data: resp,
	}
	c.JSON(http.StatusOK, response)
}

func GetLabConfig(c *gin.Context) {
	var reqVo orgvo.GetLabReqVo
	if err := c.ShouldBindJSON(&reqVo); err != nil {
		logger.Error("GetLabConfig bind request failed", err)
		response := orgvo.GetLabRespVo{Err: vo.NewErr(errs.BuildSystemErrorInfo(errs.ReqParamsValidateError, err))}
		c.JSON(http.StatusOK, response)
		return
	}

	resp, err := orgsvcService.GetLabConfig(reqVo)
	response := orgvo.GetLabRespVo{
		Err:  vo.NewErr(err),
		Data: resp,
	}
	c.JSON(http.StatusOK, response)
}
