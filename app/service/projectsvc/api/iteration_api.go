package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/star-table/startable-server/app/service/projectsvc/service"
	"github.com/star-table/startable-server/common/core/errs"

	"github.com/star-table/startable-server/common/model/vo"
	"github.com/star-table/startable-server/common/model/vo/projectvo"
)

func IterationList(c *gin.Context) {
	req := projectvo.IterationListReqVo{}
	if err := c.ShouldBindJSON(&req); err != nil {
		// Replaced: c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		c.JSON(http.StatusOK, projectvo.IterationListRespVo{Err: vo.NewErr(errs.BuildSystemErrorInfo(errs.ReqParamsValidateError, err))})
		return
	}

	res, err := service.IterationList(req.OrgId, req.Page, req.Size, req.Input)
	if err != nil {
		// Replaced: c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()}); return
	}
	c.JSON(http.StatusOK, projectvo.IterationListRespVo{Err: vo.NewErr(err), IterationList: res})
}

func CreateIteration(c *gin.Context) {
	req := projectvo.CreateIterationReqVo{}
	if err := c.ShouldBindJSON(&req); err != nil {
		// Replaced: c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		c.JSON(http.StatusOK, vo.CommonRespVo{Err: vo.NewErr(errs.BuildSystemErrorInfo(errs.ReqParamsValidateError, err))})
		return
	}

	res, err := service.CreateIteration(req.OrgId, req.UserId, req.Input)
	if err != nil {
		// Replaced: c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()}); return
	}
	c.JSON(http.StatusOK, vo.CommonRespVo{Err: vo.NewErr(err), Void: res})
}

func UpdateIteration(c *gin.Context) {
	req := projectvo.UpdateIterationReqVo{}
	if err := c.ShouldBindJSON(&req); err != nil {
		// Replaced: c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		c.JSON(http.StatusOK, vo.CommonRespVo{Err: vo.NewErr(errs.BuildSystemErrorInfo(errs.ReqParamsValidateError, err))})
		return
	}

	res, err := service.UpdateIteration(req.OrgId, req.UserId, req.Input)
	if err != nil {
		// Replaced: c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()}); return
	}
	c.JSON(http.StatusOK, vo.CommonRespVo{Err: vo.NewErr(err), Void: res})
}

func DeleteIteration(c *gin.Context) {
	req := projectvo.DeleteIterationReqVo{}
	if err := c.ShouldBindJSON(&req); err != nil {
		// Replaced: c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		c.JSON(http.StatusOK, vo.CommonRespVo{Err: vo.NewErr(errs.BuildSystemErrorInfo(errs.ReqParamsValidateError, err))})
		return
	}

	res, err := service.DeleteIteration(req.OrgId, req.UserId, req.Input)
	if err != nil {
		// Replaced: c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()}); return
	}
	c.JSON(http.StatusOK, vo.CommonRespVo{Err: vo.NewErr(err), Void: res})
}

func IterationStatusTypeStat(c *gin.Context) {
	req := projectvo.IterationStatusTypeStatReqVo{}
	if err := c.ShouldBindJSON(&req); err != nil {
		// Replaced: c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		c.JSON(http.StatusOK, projectvo.IterationStatusTypeStatRespVo{Err: vo.NewErr(errs.BuildSystemErrorInfo(errs.ReqParamsValidateError, err))})
		return
	}

	res, err := service.IterationStatusTypeStat(req.OrgId, req.Input)
	if err != nil {
		// Replaced: c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()}); return
	}
	c.JSON(http.StatusOK, projectvo.IterationStatusTypeStatRespVo{Err: vo.NewErr(err), IterationStatusTypeStat: res})
}

//func (PostGreeter) IterationIssueRelate(reqVo projectvo.IterationIssueRelateReqVo) vo.CommonRespVo {
//	res, err := service.IterationIssueRelate(reqVo.OrgId, reqVo.UserId, reqVo.Input)
//	return vo.CommonRespVo{Err: vo.NewErr(err), Void: res}
//}

func UpdateIterationStatus(c *gin.Context) {
	req := projectvo.UpdateIterationStatusReqVo{}
	if err := c.ShouldBindJSON(&req); err != nil {
		// Replaced: c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		c.JSON(http.StatusOK, vo.CommonRespVo{Err: vo.NewErr(errs.BuildSystemErrorInfo(errs.ReqParamsValidateError, err))})
		return
	}

	res, err := service.UpdateIterationStatus(req.OrgId, req.UserId, req.Input)
	if err != nil {
		// Replaced: c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()}); return
	}
	c.JSON(http.StatusOK, vo.CommonRespVo{Err: vo.NewErr(err), Void: res})
}

func IterationInfo(c *gin.Context) {
	req := projectvo.IterationInfoReqVo{}
	if err := c.ShouldBindJSON(&req); err != nil {
		// Replaced: c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		c.JSON(http.StatusOK, projectvo.IterationInfoRespVo{Err: vo.NewErr(errs.BuildSystemErrorInfo(errs.ReqParamsValidateError, err))})
		return
	}

	res, err := service.IterationInfo(req.OrgId, req.Input)
	if err != nil {
		// Replaced: c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()}); return
	}
	c.JSON(http.StatusOK, projectvo.IterationInfoRespVo{Err: vo.NewErr(err), IterationInfo: res})
}

// 获取未完成的的迭代列表
func GetNotCompletedIterationBoList(c *gin.Context) {
	req := projectvo.GetNotCompletedIterationBoListReqVo{}
	if err := c.ShouldBindJSON(&req); err != nil {
		// Replaced: c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		c.JSON(http.StatusOK, projectvo.GetNotCompletedIterationBoListRespVo{Err: vo.NewErr(errs.BuildSystemErrorInfo(errs.ReqParamsValidateError, err))})
		return
	}

	res, err := service.GetNotCompletedIterationBoList(req.OrgId, req.ProjectId)
	if err != nil {
		// Replaced: c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()}); return
	}
	c.JSON(http.StatusOK, projectvo.GetNotCompletedIterationBoListRespVo{Err: vo.NewErr(err), IterationBoList: res})
}

func UpdateIterationStatusTime(c *gin.Context) {
	req := projectvo.UpdateIterationStatusTimeReqVo{}
	if err := c.ShouldBindJSON(&req); err != nil {
		// Replaced: c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		c.JSON(http.StatusOK, vo.CommonRespVo{Err: vo.NewErr(errs.BuildSystemErrorInfo(errs.ReqParamsValidateError, err))})
		return
	}

	res, err := service.UpdateIterationStatusTime(req.OrgId, req.UserId, req.Input)
	if err != nil {
		// Replaced: c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()}); return
	}
	c.JSON(http.StatusOK, vo.CommonRespVo{Err: vo.NewErr(err), Void: res})
}

func GetSimpleIterationInfo(c *gin.Context) {
	req := projectvo.GetSimpleIterationInfoReqVo{}
	if err := c.ShouldBindJSON(&req); err != nil {
		// Replaced: c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		c.JSON(http.StatusOK, projectvo.GetSimpleIterationInfoRespVo{Err: vo.NewErr(errs.BuildSystemErrorInfo(errs.ReqParamsValidateError, err))})
		return
	}

	res, err := service.GetSimpleIterationInfo(req.OrgId, req.IterationIds)
	if err != nil {
		// Replaced: c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()}); return
	}
	c.JSON(http.StatusOK, projectvo.GetSimpleIterationInfoRespVo{
		Err:  vo.NewErr(err),
		Data: res,
	})
}

func UpdateIterationSort(c *gin.Context) {
	req := projectvo.UpdateIterationSortReqVo{}
	if err := c.ShouldBindJSON(&req); err != nil {
		// Replaced: c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		c.JSON(http.StatusOK, vo.CommonRespVo{Err: vo.NewErr(errs.BuildSystemErrorInfo(errs.ReqParamsValidateError, err))})
		return
	}

	res, err := service.UpdateIterationSort(req.OrgId, req.UserId, req.Params)
	if err != nil {
		// Replaced: c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()}); return
	}
	c.JSON(http.StatusOK, vo.CommonRespVo{
		Err:  vo.NewErr(err),
		Void: res,
	})
}
