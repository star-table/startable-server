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

// GetFunctionConfig 获取当前组织支持的功能 key 列表
func GetFunctionConfig(c *gin.Context) {
	var req orgvo.GetOrgFunctionConfigReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusOK, orgvo.GetOrgFunctionConfigResp{Err: vo.NewErr(errs.BuildSystemErrorInfo(errs.ReqParamsValidateError, err))})
		return
	}

	res, err := orgsvcService.GetFunctionKeysByOrg(req.OrgId)
	if err != nil {
		logger.Error("GetFunctionConfig error", logger.ErrorField(err))
	}
	c.JSON(http.StatusOK, orgvo.GetOrgFunctionConfigResp{
		Data: res,
		Err:  vo.NewErr(err),
	})
}

// GetFunctionObjArrByOrg 获取当前组织支持的功能列表信息
func GetFunctionObjArrByOrg(c *gin.Context) {
	var req orgvo.GetOrgFunctionConfigReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusOK, orgvo.GetFunctionArrByOrgResp{Err: vo.NewErr(errs.BuildSystemErrorInfo(errs.ReqParamsValidateError, err))})
		return
	}

	res, err := orgsvcService.GetFunctionsByOrg(req.OrgId)
	if err != nil {
		logger.Error("GetFunctionObjArrByOrg error", logger.ErrorField(err))
	}
	c.JSON(http.StatusOK, orgvo.GetFunctionArrByOrgResp{
		Data: orgvo.GetFunctionArrByOrgRespData{
			Functions: res,
		},
		Err: vo.NewErr(err),
	})
}

func UpdateOrgFunctionConfig(c *gin.Context) {
	var req orgvo.UpdateOrgFunctionConfigReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusOK, vo.CommonRespVo{Err: vo.NewErr(errs.BuildSystemErrorInfo(errs.ReqParamsValidateError, err))})
		return
	}

	err := orgsvcService.UpdateOrgFunctionConfig(req.OrgId, req.SourceChannel, req.Input.Level, req.Input.BuyType, req.Input.PricePlanType, req.Input.PayTime, req.Input.Seats, req.Input.ExpireDays, req.Input.EndDate, req.Input.TrailDays)
	if err != nil {
		logger.Error("UpdateOrgFunctionConfig error", logger.ErrorField(err))
	}
	c.JSON(http.StatusOK, vo.CommonRespVo{
		Err:  vo.NewErr(err),
		Void: &vo.Void{ID: req.OrgId},
	})
}
