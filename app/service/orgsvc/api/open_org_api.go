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

func OpenOrgUserList(c *gin.Context) {
	var reqVo orgvo.OpenOrgUserListReqVo
	if err := c.ShouldBindJSON(&reqVo); err != nil {
		logger.Error("OpenOrgUserList bind request failed", err)
		response := orgvo.OpenOrgUserListRespVo{Err: vo.NewErr(errs.BuildSystemErrorInfo(errs.ReqParamsValidateError, err))}
		c.JSON(http.StatusOK, response)
		return
	}

	res, err := orgsvcService.OrgUserList(reqVo.OrgId, reqVo.UserId, reqVo.Page, reqVo.Size, reqVo.Input)
	response := orgvo.OpenOrgUserListRespVo{Err: vo.NewErr(err), Data: res}
	c.JSON(http.StatusOK, response)
}
