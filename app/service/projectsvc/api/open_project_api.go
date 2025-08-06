package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/star-table/startable-server/app/facade/orgfacade"
	"github.com/star-table/startable-server/app/service/projectsvc/service"
	"github.com/star-table/startable-server/common/core/consts"
	"github.com/star-table/startable-server/common/core/errs"
	"github.com/star-table/startable-server/common/model/vo"
	"github.com/star-table/startable-server/common/model/vo/projectvo"
)

func OpenCreateProject(c *gin.Context) {
	var reqVo projectvo.CreateProjectReqVo
	if err := c.ShouldBindJSON(&reqVo); err != nil {
		// Replaced: c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		c.JSON(http.StatusBadRequest, projectvo.ProjectRespVo{Err: vo.NewErr(errs.ReqParamsValidateError), Project: nil})
		return
	}

	// 校验操作人是否存在
	_, err := orgfacade.GetBaseUserInfoRelaxed(reqVo.OrgId, reqVo.UserId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// // common/core/consts/business_consts.go:83
	// consts.ProjectTypeCommon2022V47
	// 空白普通项目
	projectTypeID := int64(consts.ProjectTypeCommon2022V47)
	// 私有项目,公开项目暂时不支持
	projectStatus := int(consts.PrivateProject)

	//
	reqVo.Input.ProjectTypeID = &projectTypeID
	reqVo.Input.PublicStatus = projectStatus

	res, err := service.CreateProjectWithoutAuth(reqVo)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, projectvo.ProjectRespVo{Err: vo.NewErr(err), Project: res})
}

func OpenProjects(c *gin.Context) {
	var reqVo projectvo.ProjectsRepVo
	if err := c.ShouldBindJSON(&reqVo); err != nil {
		// Replaced: c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		c.JSON(http.StatusBadRequest, projectvo.ProjectsRespVo{Err: vo.NewErr(errs.ReqParamsValidateError), ProjectList: nil})
		return
	}

	res, err := service.Projects(reqVo)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, projectvo.ProjectsRespVo{
		Err:         vo.NewErr(err),
		ProjectList: res,
	})
}

func OpenProjectInfo(c *gin.Context) {
	var reqVo projectvo.ProjectInfoReqVo
	if err := c.ShouldBindJSON(&reqVo); err != nil {
		// Replaced: c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		c.JSON(http.StatusBadRequest, projectvo.ProjectInfoRespVo{Err: vo.NewErr(errs.ReqParamsValidateError), ProjectInfo: nil})
		return
	}

	res, err := service.ProjectInfo(reqVo.OrgId, reqVo.UserId, reqVo.Input)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, projectvo.ProjectInfoRespVo{
		Err:         vo.NewErr(err),
		ProjectInfo: res,
	})
}

func OpenUpdateProject(c *gin.Context) {
	var reqVo projectvo.UpdateProjectReqVo
	if err := c.ShouldBindJSON(&reqVo); err != nil {
		// Replaced: c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		c.JSON(http.StatusBadRequest, projectvo.ProjectRespVo{Err: vo.NewErr(errs.ReqParamsValidateError), Project: nil})
		return
	}

	// 校验操作人是否存在
	_, err := orgfacade.GetBaseUserInfoRelaxed(reqVo.OrgId, reqVo.UserId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	//兼容updateFields
	input := reqVo.Input
	updateFields := make([]string, 0)
	if input.Name != nil {
		updateFields = append(updateFields, "name")
	}
	if input.Owner != nil {
		updateFields = append(updateFields, "owner")
	}
	if input.PriorityID != nil {
		updateFields = append(updateFields, "priorityId")
	}
	if input.PlanStartTime != nil {
		updateFields = append(updateFields, "planStartTime")
	}
	if input.PlanEndTime != nil {
		updateFields = append(updateFields, "planEndTime")
	}
	if input.PublicStatus != nil {
		updateFields = append(updateFields, "publicStatus")
	}
	if input.ResourceID != nil {
		updateFields = append(updateFields, "resourceId")
	}
	if input.ResourcePath != nil {
		updateFields = append(updateFields, "resourcePath")
	}
	if input.Remark != nil {
		updateFields = append(updateFields, "remark")
	}
	if input.Status != nil {
		updateFields = append(updateFields, "status")
	}
	if input.MemberIds != nil {
		updateFields = append(updateFields, "memberIds")
	}
	if input.MemberForDepartmentID != nil {
		updateFields = append(updateFields, "memberForDepartmentId")
	}
	if input.IsAllMember != nil {
		//这里也是用memberIds
		updateFields = append(updateFields, "memberIds")
	}
	if input.FollowerIds != nil {
		updateFields = append(updateFields, "followerIds")
	}
	if input.IsSyncOutCalendar != nil {
		updateFields = append(updateFields, "isSyncOutCalendar")
	}
	if input.SyncCalendarStatusList != nil {
		updateFields = append(updateFields, "syncCalendarStatusList")
	}
	if input.IsCreateFsChat != nil {
		updateFields = append(updateFields, "isCreateFsChat")
	}
	reqVo.Input.UpdateFields = updateFields

	judgeErr := service.JudgeProjectFiling(reqVo.OrgId, reqVo.Input.ID)
	if judgeErr != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": judgeErr.Error()})
		return
	}

	res, err := service.UpdateProjectWithoutAuth(reqVo)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, projectvo.ProjectRespVo{Err: vo.NewErr(err), Project: res})
}

func OpenDeleteProject(c *gin.Context) {
	var reqVo projectvo.ProjectIdReqVo
	if err := c.ShouldBindJSON(&reqVo); err != nil {
		// Replaced: c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		c.JSON(http.StatusBadRequest, vo.CommonRespVo{Err: vo.NewErr(errs.ReqParamsValidateError), Void: nil})
		return
	}

	// 校验操作人是否存在
	_, err := orgfacade.GetBaseUserInfoRelaxed(reqVo.OrgId, reqVo.UserId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	judgeErr := service.JudgeProjectFiling(reqVo.OrgId, reqVo.ProjectId)
	if judgeErr != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": judgeErr.Error()})
		return
	}

	res, err := service.DeleteProjectWithoutAuth(reqVo.OrgId, reqVo.UserId, reqVo.ProjectId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, vo.CommonRespVo{
		Err:  vo.NewErr(err),
		Void: res,
	})
}

// 优先级 废弃
func OpenGetPriorityList(c *gin.Context) {
	var reqVo projectvo.OpenPriorityListReqVo
	if err := c.ShouldBindJSON(&reqVo); err != nil {
		// Replaced: c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		c.JSON(http.StatusBadRequest, projectvo.OpenSomeAttrListRespVo{Err: vo.NewErr(errs.ReqParamsValidateError), Data: nil})
		return
	}

	resp, err := service.OpenPriorityList(reqVo.OrgId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, projectvo.OpenSomeAttrListRespVo{
		Err:  vo.NewErr(err),
		Data: resp,
	})
}

// ProjectObjectType 废弃
func OpenGetProjectObjectTypeList(c *gin.Context) {
	var reqVo projectvo.OpenPriorityListReqVo
	if err := c.ShouldBindJSON(&reqVo); err != nil {
		// Replaced: c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		c.JSON(http.StatusBadRequest, projectvo.OpenSomeAttrListRespVo{Err: vo.NewErr(errs.ReqParamsValidateError), Data: nil})
		return
	}

	resp, err := service.OpenGetProjectObjectTypeList(reqVo)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, projectvo.OpenSomeAttrListRespVo{
		Err:  vo.NewErr(err),
		Data: resp,
	})
}

func OpenGetIterationList(c *gin.Context) {
	var reqVo projectvo.OpenGetIterationListReqVo
	if err := c.ShouldBindJSON(&reqVo); err != nil {
		// Replaced: c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		c.JSON(http.StatusBadRequest, projectvo.OpenGetIterationListRespVo{Err: vo.NewErr(errs.ReqParamsValidateError), Data: nil})
		return
	}

	resp, err := service.OpenGetIterationList(reqVo)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, projectvo.OpenGetIterationListRespVo{
		Err:  vo.NewErr(err),
		Data: resp,
	})
}

// 需求来源 废弃
func OpenGetIssueSourceList(c *gin.Context) {
	var reqVo projectvo.OpenGetDemandSourceListReqVo
	if err := c.ShouldBindJSON(&reqVo); err != nil {
		// Replaced: c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		c.JSON(http.StatusBadRequest, projectvo.OpenSomeAttrListRespVo{Err: vo.NewErr(errs.ReqParamsValidateError), Data: nil})
		return
	}

	resp, err := service.OpenGetDemandSourceList(reqVo)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, projectvo.OpenSomeAttrListRespVo{
		Err:  vo.NewErr(err),
		Data: resp,
	})
}

func OpenGetPropertyList(c *gin.Context) {
	var reqVo projectvo.OpenGetPropertyListReqVo
	if err := c.ShouldBindJSON(&reqVo); err != nil {
		// Replaced: c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		c.JSON(http.StatusBadRequest, projectvo.OpenSomeAttrListRespVo{Err: vo.NewErr(errs.ReqParamsValidateError), Data: nil})
		return
	}

	resp, err := service.OpenGetPropertyList(reqVo)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, projectvo.OpenSomeAttrListRespVo{
		Err:  vo.NewErr(err),
		Data: resp,
	})
}
