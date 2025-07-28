package openapi

import (
	"strconv"

	"github.com/star-table/startable-server/app/facade/projectfacade"
	"github.com/star-table/startable-server/common/core/errs"
	"github.com/star-table/startable-server/common/model/vo"
	"github.com/star-table/startable-server/common/model/vo/projectvo"
	"github.com/gin-gonic/gin"
)

func CreateProject(c *gin.Context) {
	authData, err := ParseOpenAuthInfo(c)
	if err != nil {
		Fail(c, err)
		return
	}
	req := projectvo.OpenCreateProjectReq{}
	err1 := c.BindJSON(&req)
	if err1 != nil {
		Fail(c, errs.BuildSystemErrorInfo(errs.ReqParamsValidateError, err1))
		return
	}
	if req.OperatorId == int64(0) {
		Fail(c, errs.OperatorInvalid)
		return
	}
	respVo := projectfacade.OpenCreateProject(projectvo.CreateProjectReqVo{
		Input:  req.CreateProjectReq,
		OrgId:  authData.OrgID,
		UserId: req.OperatorId,
	})

	if respVo.Failure() {
		Fail(c, respVo.Error())
	} else {
		Suc(c, respVo.Project)
	}
}

func Projects(c *gin.Context) {
	authData, err := ParseOpenAuthInfo(c)
	if err != nil {
		Fail(c, err)
		return
	}
	req := projectvo.OpenProjectsReq{}
	err1 := c.BindJSON(&req)
	if err1 != nil {
		Fail(c, errs.BuildSystemErrorInfo(errs.ReqParamsValidateError, err1))
		return
	}
	page := 1
	size := 10
	orderBy := make([]*string, 0)
	if c.Query("page") != "" {
		page = ParseInt(c.Query("page"))
	}
	if c.Query("size") != "" {
		size = ParseInt(c.Query("size"))
	}
	if c.Query("order") != "" {
		order := ParseInt(c.Query("order"))
		orderStr := ""
		switch order {
		case 1:
			orderStr = "create_time asc"
		case 2:
			orderStr = "create_time desc"
		}
		if orderStr != "" {
			orderBy = append(orderBy, &orderStr)
		}
	}
	respVo := projectfacade.OpenProjects(projectvo.ProjectsRepVo{
		Page: page,
		Size: size,
		ProjectExtraBody: projectvo.ProjectExtraBody{
			Params: nil,
			Order:  orderBy,
			Input:  &req.ProjectsReq,
		},
		OrgId:  authData.OrgID,
		UserId: req.OperatorId,
	})
	if respVo.Failure() {
		Fail(c, respVo.Error())
	} else {
		Suc(c, respVo.ProjectList)
	}
}

func ProjectInfo(c *gin.Context) {
	authData, err := ParseOpenAuthInfo(c)
	if err != nil {
		Fail(c, err)
		return
	}

	projectId, err1 := strconv.ParseInt(c.Param("projectId"), 10, 64)
	if err1 != nil {
		Fail(c, errs.BuildSystemErrorInfo(errs.ReqParamsValidateError, err1))
		return
	}
	respVo := projectfacade.OpenProjectInfo(projectvo.ProjectInfoReqVo{
		Input:  vo.ProjectInfoReq{ProjectID: projectId},
		OrgId:  authData.OrgID,
		UserId: 0,
	})
	if respVo.Failure() {
		Fail(c, respVo.Error())
	} else {
		Suc(c, respVo.ProjectInfo)
	}
}

func UpdateProject(c *gin.Context) {
	authData, err := ParseOpenAuthInfo(c)
	if err != nil {
		Fail(c, err)
		return
	}
	req := projectvo.OpenUpdateProjectReq{}
	err1 := c.BindJSON(&req)
	if err1 != nil {
		Fail(c, errs.BuildSystemErrorInfo(errs.ReqParamsValidateError, err1))
		return
	}
	if req.OperatorId == int64(0) {
		Fail(c, errs.OperatorInvalid)
		return
	}
	projectId, err1 := strconv.ParseInt(c.Param("projectId"), 10, 64)
	if err1 != nil {
		Fail(c, errs.BuildSystemErrorInfo(errs.ReqParamsValidateError, err1))
		return
	}

	req.UpdateProjectReq.ID = projectId

	respVo := projectfacade.OpenUpdateProject(projectvo.UpdateProjectReqVo{
		Input:  req.UpdateProjectReq,
		OrgId:  authData.OrgID,
		UserId: req.OperatorId,
	})
	if respVo.Failure() {
		Fail(c, respVo.Error())
	} else {
		Suc(c, vo.Void{ID: projectId})
	}
}

func DeleteProject(c *gin.Context) {
	authData, err := ParseOpenAuthInfo(c)
	if err != nil {
		Fail(c, err)
		return
	}
	req := projectvo.OpenOperatorReq{}
	err1 := c.BindJSON(&req)
	if err1 != nil {
		Fail(c, errs.BuildSystemErrorInfo(errs.ReqParamsValidateError, err1))
		return
	}
	if req.OperatorId == int64(0) {
		Fail(c, errs.OperatorInvalid)
		return
	}

	projectId, err1 := strconv.ParseInt(c.Param("projectId"), 10, 64)
	if err1 != nil {
		Fail(c, errs.BuildSystemErrorInfo(errs.ReqParamsValidateError, err1))
		return
	}

	respVo := projectfacade.OpenDeleteProject(projectvo.ProjectIdReqVo{
		ProjectId: projectId,
		OrgId:     authData.OrgID,
		UserId:    req.OperatorId,
	})
	if respVo.Failure() {
		Fail(c, respVo.Error())
	} else {
		Suc(c, respVo.Void)
	}
}
