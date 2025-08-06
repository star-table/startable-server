package api

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/star-table/startable-server/common/core/util/copyer"

	"github.com/star-table/startable-server/app/service/projectsvc/service"
	"github.com/star-table/startable-server/common/core/errs"
	"github.com/star-table/startable-server/common/core/util/format"
	"github.com/star-table/startable-server/common/model/vo"
	"github.com/star-table/startable-server/common/model/vo/projectvo"
)

func CreateIssueComment(c *gin.Context) {
	var reqVo projectvo.CreateIssueCommentReqVo
	if err := c.ShouldBindJSON(&reqVo); err != nil {
		c.JSON(http.StatusBadRequest, vo.CommonRespVo{Err: vo.NewErr(errs.ReqParamsValidateError), Void: nil})
		return
	}

	reqVo.Input.Comment = strings.TrimSpace(reqVo.Input.Comment)
	attachmentIds := reqVo.Input.AttachmentIds
	isCommentRight := format.VerifyIssueCommenFormat(reqVo.Input.Comment)
	if !isCommentRight && (attachmentIds == nil || len(attachmentIds) == 0) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "评论格式验证失败"})
		return
	}

	res, err := service.CreateIssueComment(reqVo.OrgId, reqVo.UserId, reqVo.Input)
	c.JSON(http.StatusOK, vo.CommonRespVo{Err: vo.NewErr(err), Void: res})
}

func CreateIssueResource(c *gin.Context) {
	var reqVo projectvo.CreateIssueResourceReqVo
	if err := c.ShouldBindJSON(&reqVo); err != nil {
		c.JSON(http.StatusBadRequest, vo.CommonRespVo{Err: vo.NewErr(errs.ReqParamsValidateError), Void: nil})
		return
	}

	res, err := service.CreateIssueResource(reqVo.OrgId, reqVo.UserId, &reqVo.Input)
	c.JSON(http.StatusOK, vo.CommonRespVo{Err: vo.NewErr(err), Void: res})
}

func DeleteIssueResource(c *gin.Context) {
	var reqVo projectvo.DeleteIssueResourceReqVo
	if err := c.ShouldBindJSON(&reqVo); err != nil {
		c.JSON(http.StatusBadRequest, vo.CommonRespVo{Err: vo.NewErr(errs.ReqParamsValidateError), Void: nil})
		return
	}

	res, err := service.DeleteIssueResource(reqVo.OrgId, reqVo.UserId, reqVo.Input)
	c.JSON(http.StatusOK, vo.CommonRespVo{Err: vo.NewErr(err), Void: res})
}

func IssueResources(c *gin.Context) {
	var reqVo projectvo.IssueResourcesReqVo
	if err := c.ShouldBindJSON(&reqVo); err != nil {
		c.JSON(http.StatusBadRequest, projectvo.IssueResourcesRespVo{Err: vo.NewErr(errs.ReqParamsValidateError), IssueResources: nil})
		return
	}

	res, err := service.IssueResources(reqVo.OrgId, reqVo.Page, reqVo.Size, reqVo.Input)
	c.JSON(http.StatusOK, projectvo.IssueResourcesRespVo{Err: vo.NewErr(err), IssueResources: res})
}

func GetIssueRelationResource(c *gin.Context) {
	var req projectvo.GetIssueRelationResourceReqVo
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, projectvo.GetIssueRelationResourceRespVo{Err: vo.NewErr(errs.ReqParamsValidateError), Data: nil})
		return
	}

	res, err := service.GetIssueRelationResource(req.Page, req.Size)
	c.JSON(http.StatusOK, projectvo.GetIssueRelationResourceRespVo{
		Err:  vo.NewErr(err),
		Data: res,
	})
}

func AddIssueAttachmentFs(c *gin.Context) {
	var req projectvo.AddIssueAttachmentFsReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, projectvo.AddIssueAttachmentFsRespVo{Err: vo.NewErr(errs.ReqParamsValidateError), Data: nil})
		return
	}

	res, err := service.AddIssueAttachmentFs(req.OrgId, req.UserId, req.Input)
	var resources []*vo.Resource
	copyer.Copy(&res, &resources)
	c.JSON(http.StatusOK, projectvo.AddIssueAttachmentFsRespVo{
		Err: vo.NewErr(err),
		Data: &vo.AddIssueAttachmentFsResp{
			Resources: resources,
		},
	})
}

func AddIssueAttachment(c *gin.Context) {
	var req projectvo.AddIssueAttachmentReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, vo.CommonRespVo{Err: vo.NewErr(errs.ReqParamsValidateError), Void: nil})
		return
	}

	res, err := service.AddIssueAttachment(req.OrgId, req.UserId, req.Input)
	c.JSON(http.StatusOK, vo.CommonRespVo{
		Err:  vo.NewErr(err),
		Void: res,
	})
}

//func (PostGreeter) UpdateIssueBeforeAfterIssues(req projectvo.UpdateIssueBeforeAfterIssuesReq) vo.CommonRespVo {
//	res, err := service.UpdateIssueBeforeAfterIssues(req.OrgId, req.UserId, req.Input)
//	return vo.CommonRespVo{
//		Err:  vo.NewErr(err),
//		Void: res,
//	}
//}

//func (PostGreeter) BeforeAfterIssueList(req projectvo.BeforeAfterIssueListReq) projectvo.BeforeAfterIssueListResp {
//	res, err := service.BeforeAfterIssueList(req.OrgId, req.IssueId, req.SourceChannel)
//	return projectvo.BeforeAfterIssueListResp{
//		Err:  vo.NewErr(err),
//		NewData: res,
//	}
//}

// 分享任务卡片消息
func IssueShareCard(c *gin.Context) {
	var req projectvo.IssueCardShareReqVo
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	resp, err := service.ShareIssueCard(req.OrgId, req.UserId, req.Input)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, projectvo.IssueCardShareResp{
		Err:  vo.NewErr(err),
		Data: resp,
	})
}
