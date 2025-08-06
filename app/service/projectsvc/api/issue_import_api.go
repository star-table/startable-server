package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/star-table/startable-server/app/service/projectsvc/service"
	"github.com/star-table/startable-server/common/core/errs"
	"github.com/star-table/startable-server/common/model/vo"
	"github.com/star-table/startable-server/common/model/vo/projectvo"
)

func ImportIssues(c *gin.Context) {
	var reqVo projectvo.ImportIssuesReqVo
	if err := c.ShouldBindJSON(&reqVo); err != nil {
		// Replaced: c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		c.JSON(http.StatusBadRequest, projectvo.ImportIssuesRespVo{Err: vo.NewErr(errs.ReqParamsValidateError)})
		return
	}

	data, err := service.ImportIssues(reqVo.OrgId, reqVo.UserId, reqVo.Input)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, projectvo.ImportIssuesRespVo{Err: vo.NewErr(err), Data: data})
}

func ExportIssueTemplate(c *gin.Context) {
	var reqVo projectvo.ExportIssueTemplateReqVo
	if err := c.ShouldBindJSON(&reqVo); err != nil {
		// Replaced: c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		c.JSON(http.StatusBadRequest, projectvo.ExportIssueTemplateRespVo{Err: vo.NewErr(errs.ReqParamsValidateError)})
		return
	}

	url, err := service.ExportIssueTemplate(reqVo.OrgId, reqVo.ProjectId, reqVo.ProjectObjectTypeId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, projectvo.ExportIssueTemplateRespVo{Data: &vo.ExportIssueTemplateResp{URL: url}, Err: vo.NewErr(err)})
}

func ExportData(c *gin.Context) {
	var reqVo projectvo.ExportIssueReqVo
	if err := c.ShouldBindJSON(&reqVo); err != nil {
		// Replaced: c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		c.JSON(http.StatusBadRequest, projectvo.ExportIssueTemplateRespVo{Err: vo.NewErr(errs.ReqParamsValidateError)})
		return
	}

	url, err := service.ExportData(reqVo.OrgId, reqVo)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, projectvo.ExportIssueTemplateRespVo{Data: &vo.ExportIssueTemplateResp{URL: url}, Err: vo.NewErr(err)})
}

// ExportUserOrDeptSameNameList 导出同名部门和用户列表
func ExportUserOrDeptSameNameList(c *gin.Context) {
	var reqVo projectvo.ExportUserOrDeptSameNameListReqVo
	if err := c.ShouldBindJSON(&reqVo); err != nil {
		// Replaced: c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		c.JSON(http.StatusBadRequest, projectvo.ExportIssueTemplateRespVo{Err: vo.NewErr(errs.ReqParamsValidateError)})
		return
	}

	url, err := service.ExportUserOrDeptSameNameList(reqVo.OrgId, reqVo.ProjectId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, projectvo.ExportIssueTemplateRespVo{Data: &vo.ExportIssueTemplateResp{URL: url}, Err: vo.NewErr(err)})
}

// GetUserOrDeptSameNameList 获取同名部门或用户列表
func GetUserOrDeptSameNameList(c *gin.Context) {
	var reqVo projectvo.GetUserOrDeptSameNameListReq
	if err := c.ShouldBindJSON(&reqVo); err != nil {
		// Replaced: c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		c.JSON(http.StatusBadRequest, projectvo.GetUserOrDeptSameNameListRespVo{Err: vo.NewErr(errs.ReqParamsValidateError)})
		return
	}

	res, err := service.GetUserOrDeptSameNameList(reqVo.OrgId, reqVo.DataType)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, projectvo.GetUserOrDeptSameNameListRespVo{Data: res, Err: vo.NewErr(err)})
}

// ImportIssues2022 2022 实现的任务导入
// 无需解析 excel，直接根据入参进行任务新增/更新
//func (PostGreeter) ImportIssues2022(reqVo projectvo.ImportIssues2022ReqVo) projectvo.ImportIssuesRespVo {
//	data, err := service.ImportIssues2022(reqVo.OrgId, reqVo.UserId, &reqVo.Input)
//	return projectvo.ImportIssuesRespVo{Err: vo.NewErr(err), Data: data}
//}
