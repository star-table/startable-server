package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/star-table/startable-server/app/service/projectsvc/domain"
	"github.com/star-table/startable-server/app/service/projectsvc/service"
	"github.com/star-table/startable-server/common/core/consts"
	"github.com/star-table/startable-server/common/model/vo"
	"github.com/star-table/startable-server/common/model/vo/projectvo"
)

// GetProjectTemplateInner 获取项目模板内部接口
func GetProjectTemplateInner(c *gin.Context) {
	req := projectvo.GetProjectTemplateReq{}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	res, err := service.GetProjectTemplateInner(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, projectvo.GetProjectTemplateResp{Err: vo.NewErr(err), Data: res})
}

// ApplyProjectTemplateInner 应用项目模板
func ApplyProjectTemplateInner(c *gin.Context) {
	req := projectvo.ApplyProjectTemplateReq{}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	data, err := service.ApplyProjectTemplateInner(&req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, projectvo.ApplyProjectTemplateResp{
		Err:  vo.NewErr(err),
		Data: data,
	})
}

// 更新项目群聊成员
func ChangeProjectChatMember(c *gin.Context) {
	req := projectvo.ChangeProjectMemberReq{}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := service.ChangeProjectChatMember(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, vo.CommonRespVo{Err: vo.NewErr(err)})
}

// 批量删除模板项目
func DeleteProjectBatchInner(c *gin.Context) {
	req := projectvo.DeleteProjectInnerReq{}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := service.DeleteProjectBatchInner(req.OrgId, req.UserId, req.ProjectIds)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, vo.CommonRespVo{Err: vo.NewErr(err)})
}

// 是否允许创建项目
func AuthCreateProject(c *gin.Context) {
	req := projectvo.AuthCreateProjectReq{}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := domain.AuthPayProjectNum(req.OrgId, consts.FunctionProjectCreate)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, vo.CommonRespVo{Err: vo.NewErr(err)})
}
