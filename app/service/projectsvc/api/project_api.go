package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/star-table/startable-server/app/service/projectsvc/service"
	"github.com/star-table/startable-server/common/core/errs"

	"github.com/star-table/startable-server/common/model/vo"
	"github.com/star-table/startable-server/common/model/vo/projectvo"
)

func Projects(c *gin.Context) {
	reqVo := projectvo.ProjectsRepVo{}
	if err := c.ShouldBindJSON(&reqVo); err != nil {
		// Replaced: c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		c.JSON(http.StatusOK, projectvo.ProjectsRespVo{Err: vo.NewErr(errs.BuildSystemErrorInfo(errs.ReqParamsValidateError, err))})
		return
	}

	res, err := service.Projects(reqVo)
	if err != nil {
		// Replaced: c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()}); return
	}
	c.JSON(http.StatusOK, projectvo.ProjectsRespVo{Err: vo.NewErr(err), ProjectList: res})
}

func CreateProject(c *gin.Context) {
	reqVo := projectvo.CreateProjectReqVo{}
	if err := c.ShouldBindJSON(&reqVo); err != nil {
		// Replaced: c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		c.JSON(http.StatusOK, projectvo.ProjectRespVo{Err: vo.NewErr(errs.BuildSystemErrorInfo(errs.ReqParamsValidateError, err))})
		return
	}

	res, err := service.CreateProject(reqVo)
	if err != nil {
		// Replaced: c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()}); return
	}
	c.JSON(http.StatusOK, projectvo.ProjectRespVo{Err: vo.NewErr(err), Project: res})
}

func UpdateProject(c *gin.Context) {
	reqVo := projectvo.UpdateProjectReqVo{}
	if err := c.ShouldBindJSON(&reqVo); err != nil {
		// Replaced: c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		c.JSON(http.StatusOK, projectvo.ProjectRespVo{Err: vo.NewErr(errs.BuildSystemErrorInfo(errs.ReqParamsValidateError, err))})
		return
	}

	res, err := service.UpdateProject(reqVo)
	if err != nil {
		// Replaced: c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()}); return
	}
	c.JSON(http.StatusOK, projectvo.ProjectRespVo{Err: vo.NewErr(err), Project: res})
}

func UpdateProjectStatus(c *gin.Context) {
	reqVo := projectvo.UpdateProjectStatusReqVo{}
	if err := c.ShouldBindJSON(&reqVo); err != nil {
		// Replaced: c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		c.JSON(http.StatusOK, vo.CommonRespVo{Err: vo.NewErr(errs.BuildSystemErrorInfo(errs.ReqParamsValidateError, err))})
		return
	}

	res, err := service.UpdateProjectStatus(reqVo)
	if err != nil {
		// Replaced: c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()}); return
	}
	c.JSON(http.StatusOK, vo.CommonRespVo{Err: vo.NewErr(err), Void: res})
}

func ProjectInfo(c *gin.Context) {
	reqVo := projectvo.ProjectInfoReqVo{}
	if err := c.ShouldBindJSON(&reqVo); err != nil {
		// Replaced: c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		c.JSON(http.StatusOK, projectvo.ProjectInfoRespVo{Err: vo.NewErr(errs.BuildSystemErrorInfo(errs.ReqParamsValidateError, err))})
		return
	}

	res, err := service.ProjectInfo(reqVo.OrgId, reqVo.UserId, reqVo.Input)
	if err != nil {
		// Replaced: c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()}); return
	}
	c.JSON(http.StatusOK, projectvo.ProjectInfoRespVo{Err: vo.NewErr(err), ProjectInfo: res})
}

// GetProjectBoListByProjectTypeLangCode 通过项目类型langCode获取项目列表
func GetProjectBoListByProjectTypeLangCode(c *gin.Context) {
	req := projectvo.GetProjectBoListByProjectTypeLangCodeReqVo{}
	if err := c.ShouldBindJSON(&req); err != nil {
		// Replaced: c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		c.JSON(http.StatusOK, projectvo.GetProjectBoListByProjectTypeLangCodeRespVo{Err: vo.NewErr(errs.BuildSystemErrorInfo(errs.ReqParamsValidateError, err))})
		return
	}

	res, err := service.GetProjectBoListByProjectTypeLangCode(req.OrgId, req.ProjectTypeLangCode)
	if err != nil {
		// Replaced: c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()}); return
	}
	c.JSON(http.StatusOK, projectvo.GetProjectBoListByProjectTypeLangCodeRespVo{ProjectBoList: res, Err: vo.NewErr(err)})
}

func GetSimpleProjectInfo(c *gin.Context) {
	req := projectvo.GetSimpleProjectInfoReqVo{}
	if err := c.ShouldBindJSON(&req); err != nil {
		// Replaced: c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		c.JSON(http.StatusOK, projectvo.GetSimpleProjectInfoRespVo{Err: vo.NewErr(errs.BuildSystemErrorInfo(errs.ReqParamsValidateError, err))})
		return
	}

	res, err := service.GetSimpleProjectInfo(req.OrgId, req.Ids)
	if err != nil {
		// Replaced: c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()}); return
	}
	c.JSON(http.StatusOK, projectvo.GetSimpleProjectInfoRespVo{Err: vo.NewErr(err), Data: res})
}

func GetProjectDetails(c *gin.Context) {
	req := projectvo.GetSimpleProjectInfoReqVo{}
	if err := c.ShouldBindJSON(&req); err != nil {
		// Replaced: c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		c.JSON(http.StatusOK, projectvo.GetProjectDetailsRespVo{Err: vo.NewErr(errs.BuildSystemErrorInfo(errs.ReqParamsValidateError, err))})
		return
	}

	res, err := service.GetProjectDetails(req.OrgId, req.Ids)
	if err != nil {
		// Replaced: c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()}); return
	}
	c.JSON(http.StatusOK, projectvo.GetProjectDetailsRespVo{Err: vo.NewErr(err), Data: res})
}

func GetProjectRelation(c *gin.Context) {
	req := projectvo.GetProjectRelationReqVo{}
	if err := c.ShouldBindJSON(&req); err != nil {
		// Replaced: c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		c.JSON(http.StatusOK, projectvo.GetProjectRelationRespVo{Err: vo.NewErr(errs.BuildSystemErrorInfo(errs.ReqParamsValidateError, err))})
		return
	}

	res, err := service.GetProjectRelation(req.ProjectId, req.RelationType)
	if err != nil {
		// Replaced: c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()}); return
	}
	c.JSON(http.StatusOK, projectvo.GetProjectRelationRespVo{Err: vo.NewErr(err), Data: res})
}

func GetProjectRelationBatch(c *gin.Context) {
	req := projectvo.GetProjectRelationBatchReqVo{}
	if err := c.ShouldBindJSON(&req); err != nil {
		// Replaced: c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		c.JSON(http.StatusOK, projectvo.GetProjectRelationBatchRespVo{Err: vo.NewErr(errs.BuildSystemErrorInfo(errs.ReqParamsValidateError, err))})
		return
	}

	res, err := service.GetProjectRelationBatch(req.OrgId, req.Data)
	if err != nil {
		// Replaced: c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()}); return
	}
	c.JSON(http.StatusOK, projectvo.GetProjectRelationBatchRespVo{Err: vo.NewErr(err), Data: res})
}

func ArchiveProject(c *gin.Context) {
	reqVo := projectvo.ProjectIdReqVo{}
	if err := c.ShouldBindJSON(&reqVo); err != nil {
		// Replaced: c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		c.JSON(http.StatusOK, vo.CommonRespVo{Err: vo.NewErr(errs.BuildSystemErrorInfo(errs.ReqParamsValidateError, err))})
		return
	}

	res, err := service.ArchiveProject(reqVo.OrgId, reqVo.UserId, reqVo.ProjectId)
	if err != nil {
		// Replaced: c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()}); return
	}
	c.JSON(http.StatusOK, vo.CommonRespVo{Err: vo.NewErr(err), Void: res})
}

func CancelArchivedProject(c *gin.Context) {
	reqVo := projectvo.ProjectIdReqVo{}
	if err := c.ShouldBindJSON(&reqVo); err != nil {
		// Replaced: c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		c.JSON(http.StatusOK, vo.CommonRespVo{Err: vo.NewErr(errs.BuildSystemErrorInfo(errs.ReqParamsValidateError, err))})
		return
	}

	res, err := service.CancelArchivedProject(reqVo.OrgId, reqVo.UserId, reqVo.ProjectId)
	if err != nil {
		// Replaced: c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()}); return
	}
	c.JSON(http.StatusOK, vo.CommonRespVo{Err: vo.NewErr(err), Void: res})
}

func GetCacheProjectInfo(c *gin.Context) {
	reqVo := projectvo.GetCacheProjectInfoReqVo{}
	if err := c.ShouldBindJSON(&reqVo); err != nil {
		// Replaced: c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		c.JSON(http.StatusOK, projectvo.GetCacheProjectInfoRespVo{Err: vo.NewErr(errs.BuildSystemErrorInfo(errs.ReqParamsValidateError, err))})
		return
	}

	res, err := service.GetCacheProjectInfo(reqVo)
	if err != nil {
		// Replaced: c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()}); return
	}
	c.JSON(http.StatusOK, projectvo.GetCacheProjectInfoRespVo{Err: vo.NewErr(err), ProjectCacheBo: res})
}

// GetProjectInfoByOrgIds 通过组织id集合获取未删除 未归档的项目
func GetProjectInfoByOrgIds(c *gin.Context) {
	req := projectvo.GetProjectInfoListByOrgIdsReqVo{}
	if err := c.ShouldBindJSON(&req); err != nil {
		// Replaced: c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		c.JSON(http.StatusOK, projectvo.GetProjectInfoListByOrgIdsListRespVo{Err: vo.NewErr(errs.BuildSystemErrorInfo(errs.ReqParamsValidateError, err))})
		return
	}

	res, err := service.GetProjectInfoByOrgIds(req.OrgIds)
	if err != nil {
		// Replaced: c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()}); return
	}
	c.JSON(http.StatusOK, projectvo.GetProjectInfoListByOrgIdsListRespVo{ProjectInfoListByOrgIdsRespVo: res, Err: vo.NewErr(err)})
}

func OrgProjectMember(c *gin.Context) {
	reqVo := projectvo.OrgProjectMemberReqVo{}
	if err := c.ShouldBindJSON(&reqVo); err != nil {
		// Replaced: c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		c.JSON(http.StatusOK, projectvo.OrgProjectMemberListRespVo{Err: vo.NewErr(errs.BuildSystemErrorInfo(errs.ReqParamsValidateError, err))})
		return
	}

	res, err := service.OrgProjectMembers(reqVo)
	if err != nil {
		// Replaced: c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()}); return
	}
	c.JSON(http.StatusOK, projectvo.OrgProjectMemberListRespVo{Err: vo.NewErr(err), OrgProjectMemberRespVo: res})
}

func DeleteProject(c *gin.Context) {
	reqVo := projectvo.ProjectIdReqVo{}
	if err := c.ShouldBindJSON(&reqVo); err != nil {
		// Replaced: c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		c.JSON(http.StatusOK, vo.CommonRespVo{Err: vo.NewErr(errs.BuildSystemErrorInfo(errs.ReqParamsValidateError, err))})
		return
	}

	res, err := service.DeleteProject(reqVo.OrgId, reqVo.UserId, reqVo.ProjectId)
	if err != nil {
		// Replaced: c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()}); return
	}
	c.JSON(http.StatusOK, vo.CommonRespVo{Err: vo.NewErr(err), Void: res})
}

func GetSimpleProjectsByOrgId(c *gin.Context) {
	req := projectvo.GetSimpleProjectsByOrgIdReq{}
	if err := c.ShouldBindJSON(&req); err != nil {
		// Replaced: c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		c.JSON(http.StatusOK, projectvo.GetSimpleProjectsByOrgIdResp{Err: vo.NewErr(errs.BuildSystemErrorInfo(errs.ReqParamsValidateError, err))})
		return
	}

	res, err := service.GetSimpleProjectsByOrgId(req.OrgId)
	if err != nil {
		// Replaced: c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()}); return
	}
	c.JSON(http.StatusOK, projectvo.GetSimpleProjectsByOrgIdResp{
		Err:  vo.NewErr(err),
		Data: res,
	})
}

func GetOrgIssueAndProjectCount(c *gin.Context) {
	req := projectvo.GetOrgIssueAndProjectCountReq{}
	if err := c.ShouldBindJSON(&req); err != nil {
		// Replaced: c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		c.JSON(http.StatusOK, projectvo.GetOrgIssueAndProjectCountResp{Err: vo.NewErr(errs.BuildSystemErrorInfo(errs.ReqParamsValidateError, err))})
		return
	}

	res, err := service.GetOrgIssueAndProjectCount(req.OrgId)
	if err != nil {
		// Replaced: c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()}); return
	}
	c.JSON(http.StatusOK, projectvo.GetOrgIssueAndProjectCountResp{
		Err:  vo.NewErr(err),
		Data: res,
	})
}

// QueryProcessForAsyncTask 异步任务查询进度条
func QueryProcessForAsyncTask(c *gin.Context) {
	req := projectvo.QueryProcessForAsyncTaskReqVo{}
	if err := c.ShouldBindJSON(&req); err != nil {
		// Replaced: c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		c.JSON(http.StatusOK, projectvo.QueryProcessForAsyncTaskRespVo{Err: vo.NewErr(errs.BuildSystemErrorInfo(errs.ReqParamsValidateError, err))})
		return
	}

	res, err := service.QueryProcessForAsyncTask(req.OrgId, &req.Input)
	if err != nil {
		// Replaced: c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()}); return
	}
	c.JSON(http.StatusOK, projectvo.QueryProcessForAsyncTaskRespVo{
		Err:  vo.NewErr(err),
		Data: res,
	})
}
