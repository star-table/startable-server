package api

import (
	"net/http"
	"github.com/gin-gonic/gin"
	orgsvcService "github.com/star-table/startable-server/app/service/orgsvc/service"
	"github.com/star-table/startable-server/common/core/errs"
	"github.com/star-table/startable-server/common/model/vo"
	"github.com/star-table/startable-server/common/model/vo/orgvo"
)

// GetInviteCode 获取邀请码
func GetInviteCode(c *gin.Context) {
	var req orgvo.GetInviteCodeReqVo
	if err := c.ShouldBindJSON(&req); err != nil {
		response := orgvo.GetInviteCodeRespVo{Err: vo.NewErr(errs.BuildSystemErrorInfo(errs.ParamError, err))}
		c.JSON(http.StatusOK, response)
		return
	}
	
	res, err := orgsvcService.GetInviteCode(req.CurrentUserId, req.OrgId, req.SourcePlatform)
	if err != nil {
		response := orgvo.GetInviteCodeRespVo{Err: vo.NewErr(errs.BuildSystemErrorInfo(errs.ServerError, err))}
		c.JSON(http.StatusOK, response)
		return
	}
	response := orgvo.GetInviteCodeRespVo{Err: vo.NewErr(nil), Data: res}
	c.JSON(http.StatusOK, response)
}

// GetInviteInfo 获取邀请信息
func GetInviteInfo(c *gin.Context) {
	var req orgvo.GetInviteInfoReqVo
	if err := c.ShouldBindJSON(&req); err != nil {
		response := orgvo.GetInviteInfoRespVo{Err: vo.NewErr(errs.BuildSystemErrorInfo(errs.ParamError, err))}
		c.JSON(http.StatusOK, response)
		return
	}
	
	res, err := orgsvcService.GetInviteInfo(req.InviteCode)
	if err != nil {
		response := orgvo.GetInviteInfoRespVo{Err: vo.NewErr(errs.BuildSystemErrorInfo(errs.ServerError, err))}
		c.JSON(http.StatusOK, response)
		return
	}
	response := orgvo.GetInviteInfoRespVo{Err: vo.NewErr(nil), Data: res}
	c.JSON(http.StatusOK, response)
}

// InviteUserByPhones 通过多个手机号邀请
func InviteUserByPhones(c *gin.Context) {
	var req orgvo.InviteUserByPhonesReqVo
	if err := c.ShouldBindJSON(&req); err != nil {
		response := orgvo.InviteUserByPhonesRespVo{Err: vo.NewErr(errs.BuildSystemErrorInfo(errs.ParamError, err))}
		c.JSON(http.StatusOK, response)
		return
	}
	
	res, err := orgsvcService.InviteUserByPhones(req.OrgId, req.UserId, req.Input)
	if err != nil {
		response := orgvo.InviteUserByPhonesRespVo{Err: vo.NewErr(errs.BuildSystemErrorInfo(errs.ServerError, err))}
		c.JSON(http.StatusOK, response)
		return
	}
	response := orgvo.InviteUserByPhonesRespVo{Err: vo.NewErr(nil), Data: res}
	c.JSON(http.StatusOK, response)
}

// ExportInviteTemplate 下载导入模板
func ExportInviteTemplate(c *gin.Context) {
	var req orgvo.ExportInviteTemplateReqVo
	if err := c.ShouldBindJSON(&req); err != nil {
		response := orgvo.ExportInviteTemplateRespVo{Err: vo.NewErr(errs.BuildSystemErrorInfo(errs.ParamError, err))}
		c.JSON(http.StatusOK, response)
		return
	}
	
	res, err := orgsvcService.ExportInviteTemplate(req.OrgId, req.UserId)
	if err != nil {
		response := orgvo.ExportInviteTemplateRespVo{Err: vo.NewErr(errs.BuildSystemErrorInfo(errs.ServerError, err))}
		c.JSON(http.StatusOK, response)
		return
	}
	response := orgvo.ExportInviteTemplateRespVo{Err: vo.NewErr(nil), Data: res}
	c.JSON(http.StatusOK, response)
}

// ImportMembers 批量导入成员
func ImportMembers(c *gin.Context) {
	var req orgvo.ImportMembersReqVo
	if err := c.ShouldBindJSON(&req); err != nil {
		response := orgvo.ImportMembersRespVo{Err: vo.NewErr(errs.BuildSystemErrorInfo(errs.ParamError, err))}
		c.JSON(http.StatusOK, response)
		return
	}
	
	res, err := orgsvcService.ImportMembers(req.OrgId, req.UserId, req.Input)
	if err != nil {
		response := orgvo.ImportMembersRespVo{Err: vo.NewErr(errs.BuildSystemErrorInfo(errs.ServerError, err))}
		c.JSON(http.StatusOK, response)
		return
	}
	response := orgvo.ImportMembersRespVo{Err: vo.NewErr(nil), Data: res}
	c.JSON(http.StatusOK, response)
}

// AddUser 直接添加用户
func AddUser(c *gin.Context) {
	var req orgvo.AddUserReqVo
	if err := c.ShouldBindJSON(&req); err != nil {
		response := vo.CommonRespVo{Err: vo.NewErr(errs.BuildSystemErrorInfo(errs.ParamError, err))}
		c.JSON(http.StatusOK, response)
		return
	}
	
	err := orgsvcService.AddUser(&req)
	if err != nil {
		response := vo.CommonRespVo{Err: vo.NewErr(errs.BuildSystemErrorInfo(errs.ServerError, err))}
		c.JSON(http.StatusOK, response)
		return
	}
	response := vo.CommonRespVo{Err: vo.NewErr(nil)}
	c.JSON(http.StatusOK, response)
}
