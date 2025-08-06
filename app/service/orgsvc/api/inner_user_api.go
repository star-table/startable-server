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

func InnerUserInfo(c *gin.Context) {
	var req orgvo.InnerUserInfosReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusOK, vo.DataRespVo{Err: vo.NewErr(errs.BuildSystemErrorInfo(errs.ReqParamsValidateError, err))})
		return
	}

	userInfos, err := orgsvcService.InnerGetUserInfos(req.OrgId, req.Input.Ids)
	if err != nil {
		logger.Error("InnerUserInfo error", logger.ErrorField(err))
	}
	c.JSON(http.StatusOK, vo.DataRespVo{Err: vo.NewErr(err), Data: userInfos})
}
