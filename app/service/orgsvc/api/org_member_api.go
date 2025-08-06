package api

import (
	"net/http"
	"github.com/gin-gonic/gin"
	orgsvcService "github.com/star-table/startable-server/app/service/orgsvc/service"
	"github.com/star-table/startable-server/common/core/errs"
	"github.com/star-table/startable-server/common/model/vo"
	"github.com/star-table/startable-server/common/model/vo/orgvo"
	"github.com/star-table/startable-server/common/core/logger"
)

func UpdateOrgMemberStatus(c *gin.Context) {
	var reqVo orgvo.UpdateOrgMemberStatusReq
	if err := c.ShouldBindJSON(&reqVo); err != nil {
		logger.Error("UpdateOrgMemberStatus bind request failed", err)
		response := vo.CommonRespVo{Err: vo.NewErr(errs.BuildSystemErrorInfo(errs.ReqParamsValidateError, err))}
		c.JSON(http.StatusOK, response)
		return
	}

	status := reqVo.Input.Status
	if status != 1 && status != 2 {
		logger.Errorf("修改组织%d成员状态，状态不在正常范围内 %d", reqVo.OrgId, status)
		response := vo.CommonRespVo{Err: vo.NewErr(errs.BuildSystemErrorInfo(errs.ReqParamsValidateError)), Void: nil}
		c.JSON(http.StatusOK, response)
		return
	}
	res, err := orgsvcService.UpdateOrgMemberStatus(reqVo)
	response := vo.CommonRespVo{Err: vo.NewErr(err), Void: res}
	c.JSON(http.StatusOK, response)
}

func UpdateOrgMemberCheckStatus(c *gin.Context) {
	var reqVo orgvo.UpdateOrgMemberCheckStatusReq
	if err := c.ShouldBindJSON(&reqVo); err != nil {
		logger.Error("UpdateOrgMemberCheckStatus bind request failed", err)
		response := vo.CommonRespVo{Err: vo.NewErr(errs.BuildSystemErrorInfo(errs.ReqParamsValidateError, err))}
		c.JSON(http.StatusOK, response)
		return
	}

	status := reqVo.Input.CheckStatus
	if status != 1 && status != 2 && status != 3 {
		logger.Errorf("修改组织%d成员审核状态，状态不在正常范围内 %d", reqVo.OrgId, status)
		response := vo.CommonRespVo{Err: vo.NewErr(errs.BuildSystemErrorInfo(errs.ReqParamsValidateError)), Void: nil}
		c.JSON(http.StatusOK, response)
		return
	}
	res, err := orgsvcService.UpdateOrgMemberCheckStatus(reqVo)
	response := vo.CommonRespVo{Err: vo.NewErr(err), Void: res}
	c.JSON(http.StatusOK, response)
}

//func RemoveOrgMember(c *gin.Context) {
//	var reqVo orgvo.RemoveOrgMemberReq
//	if err := c.ShouldBindJSON(&reqVo); err != nil {
//		logger.Error("RemoveOrgMember bind request failed", err)
//		response := vo.CommonRespVo{Err: vo.NewErr(errs.BuildSystemErrorInfo(errs.ReqParamsValidateError, err))}
//		c.JSON(http.StatusOK, response)
//		return
//	}
//	res, err := orgsvcService.RemoveOrgMember(reqVo)
//	response := vo.CommonRespVo{Err: vo.NewErr(err), Void: res}
//	c.JSON(http.StatusOK, response)
//}

func OrgUserList(c *gin.Context) {
	var reqVo orgvo.OrgUserListReq
	if err := c.ShouldBindJSON(&reqVo); err != nil {
		logger.Error("OrgUserList bind request failed", err)
		response := orgvo.OrgUserListResp{Err: vo.NewErr(errs.BuildSystemErrorInfo(errs.ReqParamsValidateError, err))}
		c.JSON(http.StatusOK, response)
		return
	}

	res, err := orgsvcService.OrgUserList(reqVo.OrgId, reqVo.UserId, reqVo.Page, reqVo.Size, reqVo.Input)
	response := orgvo.OrgUserListResp{Err: vo.NewErr(err), Data: res}
	c.JSON(http.StatusOK, response)
}

func GetOrgUserInfoListBySourceChannel(c *gin.Context) {
	var reqVo orgvo.GetOrgUserInfoListBySourceChannelReq
	if err := c.ShouldBindJSON(&reqVo); err != nil {
		logger.Error("GetOrgUserInfoListBySourceChannel bind request failed", err)
		response := orgvo.GetOrgUserInfoListBySourceChannelResp{Err: vo.NewErr(errs.BuildSystemErrorInfo(errs.ReqParamsValidateError, err))}
		c.JSON(http.StatusOK, response)
		return
	}

	res, err := orgsvcService.GetOrgUserInfoListBySourceChannel(reqVo)
	response := orgvo.GetOrgUserInfoListBySourceChannelResp{Data: res, Err: vo.NewErr(err)}
	c.JSON(http.StatusOK, response)
}
