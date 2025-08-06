package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/star-table/startable-server/app/service/projectsvc/service/work_hour"
	"github.com/star-table/startable-server/common/core/errs"
	"github.com/star-table/startable-server/common/model/vo"
	"github.com/star-table/startable-server/common/model/vo/projectvo"
)

// GetIssueWorkHoursList 获取工时列表
func GetIssueWorkHoursList(c *gin.Context) {
	var reqVo projectvo.GetIssueWorkHoursListReqVo
	if err := c.ShouldBindJSON(&reqVo); err != nil {
		c.JSON(http.StatusOK, projectvo.GetIssueWorkHoursListRespVo{
			Err: vo.NewErr(errs.BuildSystemErrorInfo(errs.ReqParamsValidateError, err)),
		})
		return
	}

	resp, err := work_hour.GetIssueWorkHoursList(reqVo.OrgId, reqVo.UserId, reqVo.Input)
	c.JSON(http.StatusOK, projectvo.GetIssueWorkHoursListRespVo{
		Err:  vo.NewErr(err),
		Data: resp,
	})
}

// CreateIssueWorkHours 创建工时记录
func CreateIssueWorkHours(c *gin.Context) {
	var reqVo projectvo.CreateIssueWorkHoursReqVo
	if err := c.ShouldBindJSON(&reqVo); err != nil {
		c.JSON(http.StatusOK, projectvo.BoolRespVo{
			Err: vo.NewErr(errs.BuildSystemErrorInfo(errs.ReqParamsValidateError, err)),
		})
		return
	}

	resp, err := work_hour.CreateIssueWorkHours(reqVo.OrgId, reqVo.UserId, reqVo.Input, true)
	c.JSON(http.StatusOK, projectvo.BoolRespVo{
		Err:  vo.NewErr(err),
		Data: resp,
	})
}

// CreateMultiIssueWorkHours 创建详细工时记录，含多条预估记录
func CreateMultiIssueWorkHours(c *gin.Context) {
	var reqVo projectvo.CreateMultiIssueWorkHoursReqVo
	if err := c.ShouldBindJSON(&reqVo); err != nil {
		c.JSON(http.StatusOK, projectvo.BoolRespVo{
			Err: vo.NewErr(errs.BuildSystemErrorInfo(errs.ReqParamsValidateError, err)),
		})
		return
	}

	resp, err := work_hour.CreateMultiIssueWorkHours(reqVo.OrgId, reqVo.UserId, reqVo.Input)
	c.JSON(http.StatusOK, projectvo.BoolRespVo{
		Err:  vo.NewErr(err),
		Data: resp,
	})
}

// UpdateIssueWorkHours 编辑工时
func UpdateIssueWorkHours(c *gin.Context) {
	var reqVo projectvo.UpdateIssueWorkHoursReqVo
	if err := c.ShouldBindJSON(&reqVo); err != nil {
		c.JSON(http.StatusOK, projectvo.BoolRespVo{
			Err: vo.NewErr(errs.BuildSystemErrorInfo(errs.ReqParamsValidateError, err)),
		})
		return
	}

	resp, err := work_hour.UpdateIssueWorkHours(reqVo.OrgId, reqVo.UserId, reqVo.Input, true, true, true, true)
	c.JSON(http.StatusOK, projectvo.BoolRespVo{
		Err:  vo.NewErr(err),
		Data: resp,
	})
}

// UpdateMultiIssueWorkHours 编辑详细预估工时记录
func UpdateMultiIssueWorkHours(c *gin.Context) {
	var reqVo projectvo.UpdateMultiIssueWorkHoursReqVo
	if err := c.ShouldBindJSON(&reqVo); err != nil {
		c.JSON(http.StatusOK, projectvo.BoolRespVo{
			Err: vo.NewErr(errs.BuildSystemErrorInfo(errs.ReqParamsValidateError, err)),
		})
		return
	}

	resp, err := work_hour.UpdateMultiIssueWorkHourWithDelete(reqVo.OrgId, reqVo.UserId, reqVo.Input)
	c.JSON(http.StatusOK, projectvo.BoolRespVo{
		Err:  vo.NewErr(err),
		Data: resp,
	})
}

// DeleteIssueWorkHours 删除工时列表
func DeleteIssueWorkHours(c *gin.Context) {
	var reqVo projectvo.DeleteIssueWorkHoursReqVo
	if err := c.ShouldBindJSON(&reqVo); err != nil {
		c.JSON(http.StatusOK, projectvo.BoolRespVo{
			Err: vo.NewErr(errs.BuildSystemErrorInfo(errs.ReqParamsValidateError, err)),
		})
		return
	}

	resp, err := work_hour.DeleteIssueWorkHours(reqVo.OrgId, reqVo.UserId, reqVo.Input, true, true)
	c.JSON(http.StatusOK, projectvo.BoolRespVo{
		Err:  vo.NewErr(err),
		Data: resp,
	})
}

// DisOrEnableIssueWorkHours 启用、关闭一个项目的工时功能
func DisOrEnableIssueWorkHours(c *gin.Context) {
	var reqVo projectvo.DisOrEnableIssueWorkHoursReqVo
	if err := c.ShouldBindJSON(&reqVo); err != nil {
		c.JSON(http.StatusOK, projectvo.BoolRespVo{
			Err: vo.NewErr(errs.BuildSystemErrorInfo(errs.ReqParamsValidateError, err)),
		})
		return
	}

	resp, err := work_hour.DisOrEnableIssueWorkHours(reqVo.OrgId, reqVo.UserId, reqVo.Input)
	c.JSON(http.StatusOK, projectvo.BoolRespVo{
		Err:  vo.NewErr(err),
		Data: resp,
	})
}

// GetIssueWorkHoursInfo 获取一个任务的工时信息
func GetIssueWorkHoursInfo(c *gin.Context) {
	var reqVo projectvo.GetIssueWorkHoursInfoReqVo
	if err := c.ShouldBindJSON(&reqVo); err != nil {
		c.JSON(http.StatusOK, projectvo.GetIssueWorkHoursInfoRespVo{
			Err: vo.NewErr(errs.BuildSystemErrorInfo(errs.ReqParamsValidateError, err)),
		})
		return
	}

	resp, err := work_hour.GetIssueWorkHoursInfo(reqVo.OrgId, reqVo.UserId, reqVo.Input)
	c.JSON(http.StatusOK, projectvo.GetIssueWorkHoursInfoRespVo{
		Err:  vo.NewErr(err),
		Data: resp,
	})
}

// GetWorkHourStatistic 获取工时统计数据
func GetWorkHourStatistic(c *gin.Context) {
	var reqVo projectvo.GetWorkHourStatisticReqVo
	if err := c.ShouldBindJSON(&reqVo); err != nil {
		c.JSON(http.StatusOK, projectvo.GetWorkHourStatisticRespVo{
			Err: vo.NewErr(errs.BuildSystemErrorInfo(errs.ReqParamsValidateError, err)),
		})
		return
	}

	resp, err := work_hour.GetWorkHourStatistic(reqVo.OrgId, reqVo.Input)
	c.JSON(http.StatusOK, projectvo.GetWorkHourStatisticRespVo{
		Err:  vo.NewErr(err),
		Data: resp,
	})
}

// CheckIsIssueMember 检查用户是否是任务的关注人
func CheckIsIssueMember(c *gin.Context) {
	var reqVo projectvo.CheckIsIssueMemberReqVo
	if err := c.ShouldBindJSON(&reqVo); err != nil {
		c.JSON(http.StatusOK, projectvo.BoolRespVo{
			Err: vo.NewErr(errs.BuildSystemErrorInfo(errs.ReqParamsValidateError, err)),
		})
		return
	}

	resp, err := work_hour.CheckIsIssueMember(reqVo.OrgId, reqVo.Input)
	c.JSON(http.StatusOK, projectvo.BoolRespVo{
		Err:  vo.NewErr(err),
		Data: resp,
	})
}

// SetUserJoinIssue 将用户加入到任务成员中
func SetUserJoinIssue(c *gin.Context) {
	var reqVo projectvo.SetUserJoinIssueReqVo
	if err := c.ShouldBindJSON(&reqVo); err != nil {
		c.JSON(http.StatusOK, projectvo.BoolRespVo{
			Err: vo.NewErr(errs.BuildSystemErrorInfo(errs.ReqParamsValidateError, err)),
		})
		return
	}

	resp, err := work_hour.SetUserJoinIssue(reqVo.OrgId, reqVo.UserId, reqVo.Input)
	c.JSON(http.StatusOK, projectvo.BoolRespVo{
		Err:  vo.NewErr(err),
		Data: resp,
	})
}

// CheckIsEnableWorkHour 检查项目的工时是否开启
func CheckIsEnableWorkHour(c *gin.Context) {
	var reqVo projectvo.CheckIsEnableWorkHourReqVo
	if err := c.ShouldBindJSON(&reqVo); err != nil {
		c.JSON(http.StatusOK, projectvo.CheckIsEnableWorkHourRespVo{
			Err: vo.NewErr(errs.BuildSystemErrorInfo(errs.ReqParamsValidateError, err)),
		})
		return
	}

	resp, err := work_hour.CheckIsEnableWorkHour(reqVo.OrgId, reqVo.UserId, reqVo.Input)
	c.JSON(http.StatusOK, projectvo.CheckIsEnableWorkHourRespVo{
		Err:  vo.NewErr(err),
		Data: resp,
	})
}

// ExportWorkHourStatistic 工时统计数据的导出
func ExportWorkHourStatistic(c *gin.Context) {
	var reqVo projectvo.GetWorkHourStatisticReqVo
	if err := c.ShouldBindJSON(&reqVo); err != nil {
		c.JSON(http.StatusOK, projectvo.ExportWorkHourStatisticRespVo{
			Err: vo.NewErr(errs.BuildSystemErrorInfo(errs.ReqParamsValidateError, err)),
		})
		return
	}

	resp, err := work_hour.ExportWorkHourStatistic(reqVo.OrgId, reqVo.Input)
	c.JSON(http.StatusOK, projectvo.ExportWorkHourStatisticRespVo{
		Err:  vo.NewErr(err),
		Data: resp,
	})
}

// CheckIsIssueRelatedPeople 检查是否是任务相关人员（创建人，关注人，负责人）
func CheckIsIssueRelatedPeople(c *gin.Context) {
	var reqVo projectvo.CheckIsIssueMemberReqVo
	if err := c.ShouldBindJSON(&reqVo); err != nil {
		c.JSON(http.StatusOK, projectvo.BoolRespVo{
			Err: vo.NewErr(errs.BuildSystemErrorInfo(errs.ReqParamsValidateError, err)),
		})
		return
	}

	resp, err := work_hour.CheckIsIssueRelatedPeople(reqVo.OrgId, reqVo.Input)
	c.JSON(http.StatusOK, projectvo.BoolRespVo{
		Err:  vo.NewErr(err),
		Data: resp,
	})
}
