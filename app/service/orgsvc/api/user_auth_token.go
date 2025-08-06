package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/star-table/startable-server/common/core/logger"
	"github.com/star-table/startable-server/common/model/vo"
	"github.com/star-table/startable-server/common/model/vo/orgvo"
	orgsvcService "github.com/star-table/startable-server/app/service/orgsvc/service"
)

func GetCurrentUser(c *gin.Context) {
	info, err := orgsvcService.GetCurrentUser(c.Request.Context())
	res := orgvo.CacheUserInfoVo{Err: vo.NewErr(err)}
	if info != nil {
		res.CacheInfo = *info
	}
	if err != nil {
		logger.Error("GetCurrentUser error", logger.ErrorField(err))
	}
	c.JSON(http.StatusOK, res)
}

func GetCurrentUserWithoutOrgVerify(c *gin.Context) {
	info, err := orgsvcService.GetCurrentUserWithoutOrgVerify(c.Request.Context())
	res := orgvo.CacheUserInfoVo{Err: vo.NewErr(err)}
	if info != nil {
		res.CacheInfo = *info
	}
	if err != nil {
		logger.Error("GetCurrentUserWithoutOrgVerify error", logger.ErrorField(err))
	}
	c.JSON(http.StatusOK, res)
}

func GetCurrentUserWithoutPayVerify(c *gin.Context) {
	info, err := orgsvcService.GetCurrentUserWithCond(c.Request.Context(), true, false)
	res := orgvo.CacheUserInfoVo{Err: vo.NewErr(err)}
	if info != nil {
		res.CacheInfo = *info
	}
	if err != nil {
		logger.Error("GetCurrentUserWithoutPayVerify error", logger.ErrorField(err))
	}
	c.JSON(http.StatusOK, res)
}
