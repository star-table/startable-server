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

func GetBaseOrgInfo(c *gin.Context) {
	var reqVo orgvo.GetBaseOrgInfoReqVo
	if err := c.ShouldBindJSON(&reqVo); err != nil {
		c.JSON(http.StatusOK, orgvo.GetBaseOrgInfoRespVo{Err: vo.NewErr(errs.BuildSystemErrorInfo(errs.ReqParamsValidateError, err))})
		return
	}

	res, err := orgsvcService.GetBaseOrgInfo(reqVo.OrgId)
	if err != nil {
		logger.Error("GetBaseOrgInfo error", logger.ErrorField(err))
	}
	c.JSON(http.StatusOK, orgvo.GetBaseOrgInfoRespVo{Err: vo.NewErr(err), BaseOrgInfo: res})
}

func GetBaseOrgInfoByOutOrgId(c *gin.Context) {
	var reqVo orgvo.GetBaseOrgInfoByOutOrgIdReqVo
	if err := c.ShouldBindJSON(&reqVo); err != nil {
		c.JSON(http.StatusOK, orgvo.GetBaseOrgInfoByOutOrgIdRespVo{Err: vo.NewErr(errs.BuildSystemErrorInfo(errs.ReqParamsValidateError, err))})
		return
	}

	res, err := orgsvcService.GetBaseOrgInfoByOutOrgId(reqVo.OutOrgId)
	if err != nil {
		logger.Error("GetBaseOrgInfoByOutOrgId error", logger.ErrorField(err))
	}
	c.JSON(http.StatusOK, orgvo.GetBaseOrgInfoByOutOrgIdRespVo{Err: vo.NewErr(err), BaseOrgInfo: res})
}

func GetOrgOutInfoByOutOrgIdBatch(c *gin.Context) {
	var reqVo orgvo.GetOrgOutInfoByOutOrgIdBatchReqVo
	if err := c.ShouldBindJSON(&reqVo); err != nil {
		c.JSON(http.StatusOK, orgvo.GetOrgOutInfoByOutOrgIdBatchRespVo{Err: vo.NewErr(errs.BuildSystemErrorInfo(errs.ReqParamsValidateError, err))})
		return
	}

	orgBoList, err := orgsvcService.GetOrgOutInfoByOutOrgIdBatch(reqVo.Input)
	if err != nil {
		logger.Error("GetOrgOutInfoByOutOrgIdBatch error", logger.ErrorField(err))
	}
	c.JSON(http.StatusOK, orgvo.GetOrgOutInfoByOutOrgIdBatchRespVo{Err: vo.NewErr(err), Data: orgBoList})
}

func GetBaseUserInfoByEmpId(c *gin.Context) {
	var reqVo orgvo.GetBaseUserInfoByEmpIdReqVo
	if err := c.ShouldBindJSON(&reqVo); err != nil {
		c.JSON(http.StatusOK, orgvo.GetBaseUserInfoByEmpIdRespVo{Err: vo.NewErr(errs.BuildSystemErrorInfo(errs.ReqParamsValidateError, err))})
		return
	}

	res, err := orgsvcService.GetBaseUserInfoByEmpId(reqVo.OrgId, reqVo.EmpId)
	if err != nil {
		logger.Error("GetBaseUserInfoByEmpId error", logger.ErrorField(err))
	}
	c.JSON(http.StatusOK, orgvo.GetBaseUserInfoByEmpIdRespVo{Err: vo.NewErr(err), BaseUserInfo: res})
}

func GetBaseUserInfoByEmpIdBatch(c *gin.Context) {
	var reqVo orgvo.GetBaseUserInfoByEmpIdBatchReqVo
	if err := c.ShouldBindJSON(&reqVo); err != nil {
		c.JSON(http.StatusOK, orgvo.GetBaseUserInfoByEmpIdBatchRespVo{Err: vo.NewErr(errs.BuildSystemErrorInfo(errs.ReqParamsValidateError, err))})
		return
	}

	res, err := orgsvcService.GetBaseUserInfoByEmpIdBatch(reqVo.OrgId, reqVo.Input)
	if err != nil {
		logger.Error("GetBaseUserInfoByEmpIdBatch error", logger.ErrorField(err))
	}
	c.JSON(http.StatusOK, orgvo.GetBaseUserInfoByEmpIdBatchRespVo{Err: vo.NewErr(err), Data: res})
}

func GetUserConfigInfo(c *gin.Context) {
	var reqVo orgvo.GetUserConfigInfoReqVo
	if err := c.ShouldBindJSON(&reqVo); err != nil {
		c.JSON(http.StatusOK, orgvo.GetUserConfigInfoRespVo{Err: vo.NewErr(errs.BuildSystemErrorInfo(errs.ReqParamsValidateError, err))})
		return
	}

	res, err := orgsvcService.GetUserConfigInfo(reqVo.OrgId, reqVo.UserId)
	if err != nil {
		logger.Error("GetUserConfigInfo error", logger.ErrorField(err))
	}
	c.JSON(http.StatusOK, orgvo.GetUserConfigInfoRespVo{Err: vo.NewErr(err), UserConfigInfo: res})
}

func GetUserConfigInfoBatch(c *gin.Context) {
	var reqVo orgvo.GetUserConfigInfoBatchReqVo
	if err := c.ShouldBindJSON(&reqVo); err != nil {
		c.JSON(http.StatusOK, orgvo.GetUserConfigInfoBatchRespVo{Err: vo.NewErr(errs.BuildSystemErrorInfo(errs.ReqParamsValidateError, err))})
		return
	}

	boList, err := orgsvcService.GetUserConfigInfoBatch(reqVo.OrgId, &reqVo.Input)
	if err != nil {
		logger.Error("GetUserConfigInfoBatch error", logger.ErrorField(err))
	}
	c.JSON(http.StatusOK, orgvo.GetUserConfigInfoBatchRespVo{Err: vo.NewErr(err), Data: boList})
}

func GetBaseUserInfo(c *gin.Context) {
	var reqVo orgvo.GetBaseUserInfoReqVo
	if err := c.ShouldBindJSON(&reqVo); err != nil {
		c.JSON(http.StatusOK, orgvo.GetBaseUserInfoRespVo{Err: vo.NewErr(errs.BuildSystemErrorInfo(errs.ReqParamsValidateError, err))})
		return
	}

	res, err := orgsvcService.GetBaseUserInfo(reqVo.OrgId, reqVo.UserId)
	if err != nil {
		logger.Error("GetBaseUserInfo error", logger.ErrorField(err))
	}
	c.JSON(http.StatusOK, orgvo.GetBaseUserInfoRespVo{Err: vo.NewErr(err), BaseUserInfo: res})
}

func GetDingTalkBaseUserInfo(c *gin.Context) {
	var reqVo orgvo.GetDingTalkBaseUserInfoReqVo
	if err := c.ShouldBindJSON(&reqVo); err != nil {
		c.JSON(http.StatusOK, orgvo.GetBaseUserInfoRespVo{Err: vo.NewErr(errs.BuildSystemErrorInfo(errs.ReqParamsValidateError, err))})
		return
	}

	res, err := orgsvcService.GetDingTalkBaseUserInfo(reqVo.OrgId, reqVo.UserId)
	if err != nil {
		logger.Error("GetDingTalkBaseUserInfo error", logger.ErrorField(err))
	}
	c.JSON(http.StatusOK, orgvo.GetBaseUserInfoRespVo{Err: vo.NewErr(err), BaseUserInfo: res})
}

func GetBaseUserInfoBatch(c *gin.Context) {
	var reqVo orgvo.GetBaseUserInfoBatchReqVo
	if err := c.ShouldBindJSON(&reqVo); err != nil {
		c.JSON(http.StatusOK, orgvo.GetBaseUserInfoBatchRespVo{Err: vo.NewErr(errs.BuildSystemErrorInfo(errs.ReqParamsValidateError, err))})
		return
	}

	res, err := orgsvcService.GetBaseUserInfoBatch(reqVo.OrgId, reqVo.UserIds)
	if err != nil {
		logger.Error("GetBaseUserInfoBatch error", logger.ErrorField(err))
	}
	c.JSON(http.StatusOK, orgvo.GetBaseUserInfoBatchRespVo{Err: vo.NewErr(err), BaseUserInfos: res})
}

func GetShareUrl(c *gin.Context) {
	var req orgvo.GetShareUrlReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusOK, orgvo.GetShareUrlResp{Err: vo.NewErr(errs.BuildSystemErrorInfo(errs.ReqParamsValidateError, err))})
		return
	}

	res, err := orgsvcService.GetShareUrl(req.Key)
	if err != nil {
		logger.Error("GetShareUrl error", logger.ErrorField(err))
	}
	c.JSON(http.StatusOK, orgvo.GetShareUrlResp{
		Err: vo.NewErr(err),
		Url: res,
	})
}

func ClearOrgUsersPayCache(c *gin.Context) {
	var reqVo orgvo.GetBaseOrgInfoReqVo
	if err := c.ShouldBindJSON(&reqVo); err != nil {
		c.JSON(http.StatusOK, vo.CommonRespVo{Err: vo.NewErr(errs.BuildSystemErrorInfo(errs.ReqParamsValidateError, err))})
		return
	}

	err := orgsvcService.ClearOrgUsersPayCache(reqVo.OrgId)
	if err != nil {
		logger.Error("ClearOrgUsersPayCache error", logger.ErrorField(err))
	}
	c.JSON(http.StatusOK, vo.CommonRespVo{
		Err:  vo.NewErr(err),
		Void: &vo.Void{ID: 0},
	})
}

func SetShareUrl(c *gin.Context) {
	var req orgvo.SetShareUrlReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusOK, vo.CommonRespVo{Err: vo.NewErr(errs.BuildSystemErrorInfo(errs.ReqParamsValidateError, err))})
		return
	}

	err := orgsvcService.SetShareUrl(req.Data.Key, req.Data.Value)
	if err != nil {
		logger.Error("SetShareUrl error", logger.ErrorField(err))
	}
	c.JSON(http.StatusOK, vo.CommonRespVo{
		Err:  vo.NewErr(err),
		Void: &vo.Void{ID: 0},
	})
}
