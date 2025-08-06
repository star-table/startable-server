package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/spf13/cast"
	"github.com/star-table/startable-server/app/service/projectsvc/domain"
	"github.com/star-table/startable-server/app/service/projectsvc/service"
	"github.com/star-table/startable-server/common/core/errs"

	"github.com/star-table/startable-server/common/core/util/json"
	"github.com/star-table/startable-server/common/model/vo"
	"github.com/star-table/startable-server/common/model/vo/projectvo"
)

func InnerIssueFilter(c *gin.Context) {
	req := projectvo.InnerIssueFilterReq{}
	if err := c.ShouldBindJSON(&req); err != nil {
		// Replaced: c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		c.String(http.StatusOK, json.ToJsonIgnoreError(&projectvo.LcHomeIssuesRespVo{Err: vo.NewErr(errs.BuildSystemErrorInfo(errs.ReqParamsValidateError, err)), Data: "{}"}))
		return
	}

	appId := cast.ToString(req.Input.AppId)
	tableId := cast.ToString(req.Input.TableId)
	projectId, err := domain.GetProjectIdByAppId(req.OrgId, req.Input.AppId)
	if err != nil {
		// Replaced: c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()}); return
		c.String(http.StatusOK, json.ToJsonIgnoreError(&projectvo.LcHomeIssuesRespVo{Err: vo.NewErr(err), Data: "{}"}))
		return
	}
	res, err := service.LcHomeIssuesForProject(req.OrgId, req.UserId, req.Page, req.Size, &projectvo.HomeIssueInfoReq{
		MenuAppID:     &appId,
		ProjectID:     &projectId,
		TableID:       &tableId,
		FilterColumns: req.Input.Columns,
		LessConds:     req.Input.Condition,
		LessOrder:     req.Input.Orders,
	}, true)
	if err != nil {
		// Replaced: c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()}); return
		c.String(http.StatusOK, json.ToJsonIgnoreError(&projectvo.LcHomeIssuesRespVo{Err: vo.NewErr(err), Data: "{}"}))
		return
	}
	c.String(http.StatusOK, res.Data)
}

func InnerIssueCreate(c *gin.Context) {
	req := projectvo.InnerIssueCreateReq{}
	if err := c.ShouldBindJSON(&req); err != nil {
		// Replaced: c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		c.JSON(http.StatusOK, &projectvo.LcDataListRespVo{Err: vo.NewErr(errs.BuildSystemErrorInfo(errs.ReqParamsValidateError, err)), Data: nil})
		return
	}

	projectId, err := domain.GetProjectIdByAppId(req.OrgId, req.Input.AppId)
	if err != nil {
		// Replaced: c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()}); return
		c.JSON(http.StatusOK, &projectvo.LcDataListRespVo{Err: vo.NewErr(err), Data: nil})
		return
	}
	reqVo := &projectvo.BatchCreateIssueReqVo{
		OrgId:  req.OrgId,
		UserId: req.UserId,
		Input: &projectvo.BatchCreateIssueInput{
			AppId:     req.Input.AppId,
			ProjectId: projectId,
			TableId:   req.Input.TableId,
			Data:      req.Input.Data,
		},
	}
	res, userDept, relateData, err := service.SyncBatchCreateIssue(reqVo, true, req.Input.TriggerBy)
	if err != nil {
		// Replaced: c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()}); return
	}
	c.JSON(http.StatusOK, &projectvo.LcDataListRespVo{Err: vo.NewErr(err), Data: res, UserDept: userDept, RelateData: relateData})
}

func InnerIssueCreateByCopy(c *gin.Context) {
	req := projectvo.InnerIssueCreateByCopyReq{}
	if err := c.ShouldBindJSON(&req); err != nil {
		// Replaced: c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		c.JSON(http.StatusOK, &projectvo.LcDataListRespVo{Err: vo.NewErr(errs.BuildSystemErrorInfo(errs.ReqParamsValidateError, err)), Data: nil})
		return
	}

	projectId, err := domain.GetProjectIdByAppId(req.OrgId, req.Input.AppId)
	if err != nil {
		// Replaced: c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()}); return
		c.JSON(http.StatusOK, &projectvo.LcDataListRespVo{Err: vo.NewErr(err), Data: nil})
		return
	}
	res, userDept, relateData, err := service.CopyIssueBatchWithData(req.OrgId, req.UserId, projectId, cast.ToInt64(req.Input.TableId),
		req.Input.IssueIds, req.Input.Data, req.Input.TriggerBy,
		true, req.Input.IsStaticCopy, req.Input.IsCreateMissingSelectOptions, true)
	if err != nil {
		// Replaced: c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()}); return
	}
	c.JSON(http.StatusOK, &projectvo.LcDataListRespVo{Err: vo.NewErr(err), Data: res, UserDept: userDept, RelateData: relateData})
}

func InnerIssueUpdate(c *gin.Context) {
	req := projectvo.InnerIssueUpdateReq{}
	if err := c.ShouldBindJSON(&req); err != nil {
		// Replaced: c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		c.JSON(http.StatusOK, &vo.VoidErr{Err: vo.NewErr(errs.BuildSystemErrorInfo(errs.ReqParamsValidateError, err))})
		return
	}

	projectId, err := domain.GetProjectIdByAppId(req.OrgId, req.Input.AppId)
	if err != nil {
		// Replaced: c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()}); return
		c.JSON(http.StatusOK, &vo.VoidErr{Err: vo.NewErr(err)})
		return
	}
	reqVo := &projectvo.BatchUpdateIssueReqInnerVo{
		OrgId:     req.OrgId,
		UserId:    req.UserId,
		AppId:     req.Input.AppId,
		ProjectId: projectId,
		TableId:   req.Input.TableId,
		Data:      req.Input.Data,
	}
	err = service.SyncBatchUpdateIssue(reqVo, true, req.Input.TriggerBy)
	if err != nil {
		// Replaced: c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()}); return
	}
	c.JSON(http.StatusOK, &vo.VoidErr{Err: vo.NewErr(err)})
}
