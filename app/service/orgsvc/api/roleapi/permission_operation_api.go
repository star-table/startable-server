package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/star-table/startable-server/common/core/errs"
	"github.com/star-table/startable-server/common/core/logger"
	"github.com/star-table/startable-server/common/model/vo"
	"github.com/star-table/startable-server/common/model/vo/rolevo"
	orgsvcService "github.com/star-table/startable-server/app/service/orgsvc/service"
)

// TODO 替换为lesscode-permission
func GetPersonalPermissionInfo(c *gin.Context) {
	var req rolevo.GetPersonalPermissionInfoReqVo
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusOK, rolevo.GetPersonalPermissionInfoRespVo{Err: vo.NewErr(errs.BuildSystemErrorInfo(errs.ReqParamsValidateError, err))})
		return
	}

	res, err := orgsvcService.GetPersonalPermissionInfoForFuse(req.OrgId, req.UserId, req.ProjectId, req.IssueId, req.SourceChannel)
	if err != nil {
		logger.Error("GetPersonalPermissionInfo error", logger.ErrorField(err))
	}
	c.JSON(http.StatusOK, rolevo.GetPersonalPermissionInfoRespVo{Err: vo.NewErr(err), Data: res})
}
