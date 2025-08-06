package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/spf13/cast"
	"github.com/star-table/startable-server/app/service/projectsvc/service"
	"github.com/star-table/startable-server/common/core/consts"
	"github.com/star-table/startable-server/common/core/errs"
	"github.com/star-table/startable-server/common/model/vo"
	"github.com/star-table/startable-server/common/model/vo/projectvo"
)

func CreateIssue(c *gin.Context) {
	reqVo := &projectvo.CreateIssueReqVo{}
	if err := c.ShouldBindJSON(reqVo); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var beforeDataId, afterDataId int64
	if reqVo.CreateIssue.BeforeDataID != nil {
		beforeDataId = cast.ToInt64(*reqVo.CreateIssue.BeforeDataID)
	}
	if reqVo.CreateIssue.AfterDataID != nil {
		afterDataId = cast.ToInt64(*reqVo.CreateIssue.AfterDataID)
	}
	req := &projectvo.BatchCreateIssueReqVo{
		OrgId:  reqVo.OrgId,
		UserId: reqVo.UserId,
		Input: &projectvo.BatchCreateIssueInput{
			AppId:        reqVo.InputAppId,
			ProjectId:    reqVo.CreateIssue.ProjectID,
			TableId:      cast.ToInt64(reqVo.CreateIssue.TableID),
			BeforeDataId: beforeDataId,
			AfterDataId:  afterDataId,
			Data:         []map[string]interface{}{reqVo.CreateIssue.LessCreateIssueReq},
		},
	}
	data, userDepts, relateData, err := service.BatchCreateIssue(req, false, &projectvo.TriggerBy{
		TriggerBy: consts.TriggerByNormal,
	}, "", 0)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, &projectvo.IssueRespVo{Err: vo.NewErr(err), Data: data, UserDepts: userDepts, RelateData: relateData})
}

func DeleteIssue(c *gin.Context) {
	reqVo := projectvo.DeleteIssueReqVo{}
	if err := c.ShouldBindJSON(&reqVo); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	res, err := service.DeleteIssue(reqVo)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, projectvo.IssueRespVo{Err: vo.NewErr(err), Issue: res})
}

//func (PostGreeter) MoveIssue(reqVo projectvo.MoveIssueReqVo) vo.CommonRespVo {
//	res, err := service.MoveIssue(reqVo.OrgId, reqVo.UserId, reqVo.InputAppId, reqVo.Input)
//	return vo.CommonRespVo{Err: vo.NewErr(err), Void: res}
//}
//
//func (PostGreeter) MoveIssueBatch(req projectvo.MoveIssueBatchReqVo) projectvo.MoveIssueBatchRespVo {
//	res, err := service.MoveIssueBatch(req.OrgId, req.UserId, req.Input)
//	return projectvo.MoveIssueBatchRespVo{
//		Err:  vo.NewErr(err),
//		Data: res,
//	}
//}

func GetIssueInfo(c *gin.Context) {
	reqVo := projectvo.GetIssueInfoReqVo{}
	if err := c.ShouldBindJSON(&reqVo); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	res, err := service.GetIssueInfo(reqVo.OrgId, reqVo.UserId, reqVo.IssueId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, projectvo.GetIssueInfoRespVo{
		Err:       vo.NewErr(err),
		IssueInfo: res,
	})
}

func GetIssueInfoList(c *gin.Context) {
	reqVo := projectvo.IssueInfoListReqVo{}
	if err := c.ShouldBindJSON(&reqVo); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	res, err := service.GetIssueInfoList(reqVo.OrgId, reqVo.UserId, reqVo.IssueIds)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, projectvo.IssueInfoListRespVo{
		Err:        vo.NewErr(err),
		IssueInfos: res,
	})
}

func GetIssueInfoByDataIdsList(c *gin.Context) {
	reqVo := projectvo.IssueInfoListByDataIdsReqVo{}
	if err := c.ShouldBindJSON(&reqVo); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	res, err := service.GetIssueInfoListByDataIds(reqVo.OrgId, reqVo.UserId, reqVo.DataIds)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, projectvo.IssueInfoListByDataIdsRespVo{
		Err:        vo.NewErr(err),
		IssueInfos: res,
	})
}

//func (PostGreeter) CopyIssue(req projectvo.CopyIssueReqVo) *projectvo.LcDataListRespVo {
//	res, userDept, relateData, err := service.CopyIssue(req.OrgId, req.UserId, req.Input)
//	return &projectvo.LcDataListRespVo{Err: vo.NewErr(err), Data: res, UserDept: userDept, RelateData: relateData}
//}

func BatchCopyIssue(c *gin.Context) {
	req := &projectvo.LcCopyIssuesReq{}
	if err := c.ShouldBindJSON(req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	res, userDept, relateData, err := service.CopyIssueBatch(req.OrgId, req.UserId, req.Input)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, &projectvo.LcDataListRespVo{Err: vo.NewErr(err), Data: res, UserDept: userDept, RelateData: relateData})
}

func BatchMoveIssue(c *gin.Context) {
	req := &projectvo.LcMoveIssuesReq{}
	if err := c.ShouldBindJSON(req); err != nil {
		// Replaced: c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		c.JSON(http.StatusOK, &projectvo.LcMoveIssuesResp{Err: vo.NewErr(errs.BuildSystemErrorInfo(errs.ReqParamsValidateError, err))})
		return
	}

	issueIds, err := service.MoveIssueBatch(req.OrgId, req.UserId, req.Input)
	if err != nil {
		// Replaced: c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()}); return
	}
	c.JSON(http.StatusOK, &projectvo.LcMoveIssuesResp{Err: vo.NewErr(err), IssueIds: issueIds})
}

func GetFieldMapping(c *gin.Context) {
	req := &projectvo.LcGetFieldMappingReq{}
	if err := c.ShouldBindJSON(req); err != nil {
		// Replaced: c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		c.JSON(http.StatusOK, &projectvo.LcGetFieldMappingResp{Err: vo.NewErr(errs.BuildSystemErrorInfo(errs.ReqParamsValidateError, err))})
		return
	}

	fieldMapping := service.GetFieldMapping(req.OrgId, req.Input.FromTableId, req.Input.ToTableId)
	c.JSON(http.StatusOK, &projectvo.LcGetFieldMappingResp{Err: vo.NewErr(nil), FieldMapping: fieldMapping})
}

//func (PostGreeter) CopyIssueBatch(req projectvo.CopyIssueBatchReqVo) projectvo.CopyIssueBatchRespVo {
//	res, _, _, err := service.CopyIssueBatch(req.OrgId, req.UserId, &req.Input, &projectvo.TriggerBy{TriggerBy: consts.TriggerByCopy})
//	return projectvo.CopyIssueBatchRespVo{Data: int64(len(res)), Err: vo.NewErr(err)}
//}

func DeleteIssueBatch(c *gin.Context) {
	req := projectvo.DeleteIssueBatchReqVo{}
	if err := c.ShouldBindJSON(&req); err != nil {
		// Replaced: c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		c.JSON(http.StatusOK, projectvo.DeleteIssueBatchRespVo{Err: vo.NewErr(errs.BuildSystemErrorInfo(errs.ReqParamsValidateError, err))})
		return
	}

	res, err := service.DeleteIssueBatch(req)
	if err != nil {
		// Replaced: c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()}); return
	}
	c.JSON(http.StatusOK, projectvo.DeleteIssueBatchRespVo{Err: vo.NewErr(err), Data: res})
}

func ConvertIssueToParent(c *gin.Context) {
	req := projectvo.ConvertIssueToParentReq{}
	if err := c.ShouldBindJSON(&req); err != nil {
		// Replaced: c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		c.JSON(http.StatusOK, vo.CommonRespVo{Err: vo.NewErr(errs.BuildSystemErrorInfo(errs.ReqParamsValidateError, err))})
		return
	}

	res, err := service.ConvertIssueToParent(req.OrgId, req.UserId, req.Input)
	if err != nil {
		// Replaced: c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()}); return
	}
	c.JSON(http.StatusOK, vo.CommonRespVo{Err: vo.NewErr(err), Void: res})
}

func ChangeParentIssue(c *gin.Context) {
	req := projectvo.ChangeParentIssueReq{}
	if err := c.ShouldBindJSON(&req); err != nil {
		// Replaced: c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		c.JSON(http.StatusOK, vo.CommonRespVo{Err: vo.NewErr(errs.BuildSystemErrorInfo(errs.ReqParamsValidateError, err))})
		return
	}

	res, err := service.ChangeParentIssue(req.OrgId, req.UserId, req.Input)
	if err != nil {
		// Replaced: c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()}); return
	}
	c.JSON(http.StatusOK, vo.CommonRespVo{Err: vo.NewErr(err), Void: res})
}

func AuditIssue(c *gin.Context) {
	req := projectvo.AuditIssueReq{}
	if err := c.ShouldBindJSON(&req); err != nil {
		// Replaced: c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		c.JSON(http.StatusOK, vo.CommonRespVo{Err: vo.NewErr(errs.BuildSystemErrorInfo(errs.ReqParamsValidateError, err))})
		return
	}

	message := ""
	if req.Params.Comment != nil {
		message = *req.Params.Comment
	}
	issueIds, err := service.BatchAuditIssue(&projectvo.BatchAuditIssueReqVo{
		OrgId:  req.OrgId,
		UserId: req.UserId,
		Input: &vo.LessBatchAuditIssueReq{
			IssueIds:    []int64{req.Params.IssueID},
			AuditStatus: req.Params.Status,
			Message:     message,
			Attachments: req.Params.Attachments,
		},
	})
	var issueId int64
	if len(issueIds) > 0 {
		issueId = issueIds[0]
	}
	if err != nil {
		// Replaced: c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()}); return
	}
	c.JSON(http.StatusOK, vo.CommonRespVo{
		Err: vo.NewErr(err),
		Void: &vo.Void{
			ID: issueId,
		},
	})
}

func UrgeIssue(c *gin.Context) {
	req := projectvo.UrgeIssueReqVo{}
	if err := c.ShouldBindJSON(&req); err != nil {
		// Replaced: c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		c.JSON(http.StatusOK, vo.BoolRespVo{Err: vo.NewErr(errs.BuildSystemErrorInfo(errs.ReqParamsValidateError, err))})
		return
	}

	message := ""
	if req.Input.UrgeText != nil {
		message = *req.Input.UrgeText
	}
	issueIds, err := service.BatchUrgeIssue(&projectvo.BatchUrgeIssueReqVo{
		OrgId:  req.OrgId,
		UserId: req.UserId,
		Input: &vo.LessBatchUrgeIssueReq{
			IssueIds: []int64{req.Input.IssueID},
			Message:  message,
		},
	})
	if err != nil {
		// Replaced: c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()}); return
	}
	c.JSON(http.StatusOK, vo.BoolRespVo{
		Err:    vo.NewErr(err),
		IsTrue: len(issueIds) > 0,
	})
}

func UrgeAuditIssue(c *gin.Context) {
	req := projectvo.UrgeAuditIssueReq{}
	if err := c.ShouldBindJSON(&req); err != nil {
		// Replaced: c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		c.JSON(http.StatusOK, vo.BoolRespVo{Err: vo.NewErr(errs.BuildSystemErrorInfo(errs.ReqParamsValidateError, err))})
		return
	}

	issueIds, err := service.BatchUrgeIssue(&projectvo.BatchUrgeIssueReqVo{
		OrgId:  req.OrgId,
		UserId: req.UserId,
		Input: &vo.LessBatchUrgeIssueReq{
			IssueIds: []int64{req.Input.IssueID},
			Message:  req.Input.UrgeText,
		},
	})
	if err != nil {
		// Replaced: c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()}); return
	}
	c.JSON(http.StatusOK, vo.BoolRespVo{
		Err:    vo.NewErr(err),
		IsTrue: len(issueIds) > 0,
	})
}

func ViewAuditIssue(c *gin.Context) {
	req := projectvo.ViewAuditIssueReq{}
	if err := c.ShouldBindJSON(&req); err != nil {
		// Replaced: c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		c.JSON(http.StatusOK, vo.CommonRespVo{Err: vo.NewErr(errs.BuildSystemErrorInfo(errs.ReqParamsValidateError, err))})
		return
	}

	err := service.ViewAuditIssue(req.OrgId, req.UserId, req.IssueId)
	if err != nil {
		// Replaced: c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()}); return
	}
	c.JSON(http.StatusOK, vo.CommonRespVo{
		Err:  vo.NewErr(err),
		Void: &vo.Void{ID: req.IssueId},
	})
}

func WithdrawIssue(c *gin.Context) {
	req := projectvo.WithdrawIssueReq{}
	if err := c.ShouldBindJSON(&req); err != nil {
		// Replaced: c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		c.JSON(http.StatusOK, vo.CommonRespVo{Err: vo.NewErr(errs.BuildSystemErrorInfo(errs.ReqParamsValidateError, err))})
		return
	}

	err := service.WithdrawIssue(req.OrgId, req.UserId, req.IssueId)
	if err != nil {
		// Replaced: c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()}); return
	}
	c.JSON(http.StatusOK, vo.CommonRespVo{
		Err:  vo.NewErr(err),
		Void: &vo.Void{ID: req.IssueId},
	})
}

//func (PostGreeter) GetIssueIdsByOrgId(req projectvo.GetIssueIdsByOrgIdReqVo) projectvo.GetIssueIdsByOrgIdRespVo {
//	resp, err := service.GetIssueIdsByOrgId(req.OrgId, req.UserId, req.Input)
//	return projectvo.GetIssueIdsByOrgIdRespVo{
//		Err:  vo.NewErr(err),
//		Data: resp,
//	}
//}

// InsertIssueProRelation 插入任务归属于项目的关联关系
//func (PostGreeter) InsertIssueProRelation(req projectvo.InsertIssueProRelationReqVo) projectvo.BoolRespVo {
//	err := service.InsertIssueProRelation(req.OrgId, req.UserId, req.Input)
//	return projectvo.BoolRespVo{
//		Err: vo.NewErr(err),
//		Data: &vo.BoolResp{
//			IsTrue: err == nil,
//		},
//	}
//}

// GetTableStatusByTableId 查询任务状态 供其他服务调用
func GetTableStatusByTableId(c *gin.Context) {
	req := projectvo.GetTableStatusReq{}
	if err := c.ShouldBindJSON(&req); err != nil {
		// Replaced: c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		c.JSON(http.StatusOK, projectvo.GetTableStatusResp{Err: vo.NewErr(errs.BuildSystemErrorInfo(errs.ReqParamsValidateError, err))})
		return
	}

	resp, err := service.GetTableStatus(req.OrgId, req.Input.TableId)
	if err != nil {
		// Replaced: c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()}); return
	}
	c.JSON(http.StatusOK, projectvo.GetTableStatusResp{
		Err:  vo.NewErr(err),
		Data: resp,
	})
}

// IssueStartChat 发起任务的群聊讨论
func IssueStartChat(c *gin.Context) {
	req := projectvo.IssueStartChatReqVo{}
	if err := c.ShouldBindJSON(&req); err != nil {
		// Replaced: c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		c.JSON(http.StatusOK, projectvo.IssueStartChatRespVo{Err: vo.NewErr(errs.BuildSystemErrorInfo(errs.ReqParamsValidateError, err))})
		return
	}

	res, err := service.IssueStartChat(req.OrgId, req.UserId, req.SourceChannel, req.Input)
	if err != nil {
		// Replaced: c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()}); return
	}
	c.JSON(http.StatusOK, projectvo.IssueStartChatRespVo{
		Err:  vo.NewErr(err),
		Data: res,
	})
}

func GetIssueRowList(c *gin.Context) {
	req := projectvo.IssueRowListReq{}
	if err := c.ShouldBindJSON(&req); err != nil {
		// Replaced: c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		c.JSON(http.StatusOK, projectvo.IssueRowListResp{Err: vo.NewErr(errs.BuildSystemErrorInfo(errs.ReqParamsValidateError, err))})
		return
	}

	resp, err := service.GetIssueRowList(req.OrgId, req.UserId, req.Input)
	if err != nil {
		// Replaced: c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()}); return
	}
	c.JSON(http.StatusOK, projectvo.IssueRowListResp{
		Err:  vo.NewErr(err),
		Data: resp,
	})
}
