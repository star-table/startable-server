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

func GetDingJsAPISign(c *gin.Context) {
	var input orgvo.GetDingApiSignReq
	if err := c.ShouldBindJSON(&input); err != nil {
		logger.Error("GetDingJsAPISign bind request failed", err)
		response := orgvo.GetJsAPISignRespVo{Err: vo.NewErr(errs.BuildSystemErrorInfo(errs.ReqParamsValidateError, err))}
		c.JSON(http.StatusOK, response)
		return
	}

	res, err := orgsvcService.GetDingJsAPISign(input.Input)
	response := orgvo.GetJsAPISignRespVo{Err: vo.NewErr(err), GetJsAPISign: res}
	c.JSON(http.StatusOK, response)
}

func AuthDingCode(c *gin.Context) {
	var input orgvo.AuthDingCodeReqVo
	if err := c.ShouldBindJSON(&input); err != nil {
		logger.Error("AuthDingCode bind request failed", err)
		response := orgvo.AuthDingCodeRespVo{Err: vo.NewErr(errs.BuildSystemErrorInfo(errs.ReqParamsValidateError, err))}
		c.JSON(http.StatusOK, response)
		return
	}

	res, err := orgsvcService.AuthDingCode(input)
	response := orgvo.AuthDingCodeRespVo{
		Err:  vo.NewErr(err),
		Data: res,
	}
	c.JSON(http.StatusOK, response)
}

func CreateCoolApp(c *gin.Context) {
	var input orgvo.CreateCoolAppReq
	if err := c.ShouldBindJSON(&input); err != nil {
		logger.Error("CreateCoolApp bind request failed", err)
		response := vo.CommonRespVo{Err: vo.NewErr(errs.BuildSystemErrorInfo(errs.ReqParamsValidateError, err))}
		c.JSON(http.StatusOK, response)
		return
	}

	err := orgsvcService.CreateCoolApp(input.Input)
	response := vo.CommonRespVo{
		Err: vo.NewErr(err),
	}
	c.JSON(http.StatusOK, response)
}

func GetCoolAppInfo(c *gin.Context) {
	var input orgvo.GetCoolAppInfoReq
	if err := c.ShouldBindJSON(&input); err != nil {
		logger.Error("GetCoolAppInfo bind request failed", err)
		response := orgvo.GetCoolAppInfoResp{Err: vo.NewErr(errs.BuildSystemErrorInfo(errs.ReqParamsValidateError, err))}
		c.JSON(http.StatusOK, response)
		return
	}

	info, err := orgsvcService.GetCoolAppInfo(input.Input.OpenConversationId)
	response := orgvo.GetCoolAppInfoResp{
		Err:  vo.NewErr(err),
		Data: info,
	}
	c.JSON(http.StatusOK, response)
}

func DeleteCoolApp(c *gin.Context) {
	var input orgvo.DeleteCoolAppReq
	if err := c.ShouldBindJSON(&input); err != nil {
		logger.Error("DeleteCoolApp bind request failed", err)
		response := vo.CommonRespVo{Err: vo.NewErr(errs.BuildSystemErrorInfo(errs.ReqParamsValidateError, err))}
		c.JSON(http.StatusOK, response)
		return
	}

	err := orgsvcService.DeleteCoolApp(input)
	response := vo.CommonRespVo{
		Err: vo.NewErr(err),
	}
	c.JSON(http.StatusOK, response)
}

func DeleteCoolAppByProject(c *gin.Context) {
	var input orgvo.DeleteCoolAppByProjectReq
	if err := c.ShouldBindJSON(&input); err != nil {
		logger.Error("DeleteCoolAppByProject bind request failed", err)
		response := vo.CommonRespVo{Err: vo.NewErr(errs.BuildSystemErrorInfo(errs.ReqParamsValidateError, err))}
		c.JSON(http.StatusOK, response)
		return
	}

	err := orgsvcService.DeleteCoolAppByProject(input.OrgId, input.ProjectId)
	response := vo.CommonRespVo{
		Err: vo.NewErr(err),
	}
	c.JSON(http.StatusOK, response)
}

func BindCoolApp(c *gin.Context) {
	var input orgvo.BindCoolAppReq
	if err := c.ShouldBindJSON(&input); err != nil {
		logger.Error("BindCoolApp bind request failed", err)
		response := vo.CommonRespVo{Err: vo.NewErr(errs.BuildSystemErrorInfo(errs.ReqParamsValidateError, err))}
		c.JSON(http.StatusOK, response)
		return
	}

	err := orgsvcService.BindCoolApp(input)
	response := vo.CommonRespVo{
		Err: vo.NewErr(err),
	}
	c.JSON(http.StatusOK, response)
}

func UpdateCoolAppTopCard(c *gin.Context) {
	var input orgvo.UpdateCoolAppTopCardReq
	if err := c.ShouldBindJSON(&input); err != nil {
		logger.Error("UpdateCoolAppTopCard bind request failed", err)
		response := vo.CommonRespVo{Err: vo.NewErr(errs.BuildSystemErrorInfo(errs.ReqParamsValidateError, err))}
		c.JSON(http.StatusOK, response)
		return
	}

	err := orgsvcService.UpdateTopCard(input.OrgId, input.ProjectId)
	response := vo.CommonRespVo{
		Err: vo.NewErr(err),
	}
	c.JSON(http.StatusOK, response)
}

func GetUpdateCoolAppTopCardData(c *gin.Context) {
	var input orgvo.GetCoolAppTopCardDataReq
	if err := c.ShouldBindJSON(&input); err != nil {
		logger.Error("GetUpdateCoolAppTopCardData bind request failed", err)
		response := orgvo.GetCoolAppTopCardDataResp{Err: vo.NewErr(errs.BuildSystemErrorInfo(errs.ReqParamsValidateError, err))}
		c.JSON(http.StatusOK, response)
		return
	}

	data, err := orgsvcService.GetTopCardData(input.Input.OpenConversationId)
	response := orgvo.GetCoolAppTopCardDataResp{
		Err:  vo.NewErr(err),
		Data: data,
	}
	c.JSON(http.StatusOK, response)
}

func GetSpaceList(c *gin.Context) {
	var input orgvo.GetSpaceListReq
	if err := c.ShouldBindJSON(&input); err != nil {
		logger.Error("GetSpaceList bind request failed", err)
		response := orgvo.GetSpaceListResp{Err: vo.NewErr(errs.BuildSystemErrorInfo(errs.ReqParamsValidateError, err))}
		c.JSON(http.StatusOK, response)
		return
	}

	data, err := orgsvcService.GetSpaceList(input)
	response := orgvo.GetSpaceListResp{
		Err:  vo.NewErr(err),
		Data: data,
	}
	c.JSON(http.StatusOK, response)
}
