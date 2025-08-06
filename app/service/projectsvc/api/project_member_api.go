package api

import (
	"net/http"
	"github.com/gin-gonic/gin"
	"github.com/star-table/startable-server/app/service/projectsvc/service"
	"github.com/star-table/startable-server/common/model/vo"
	"github.com/star-table/startable-server/common/model/vo/projectvo"
)

func GetProjectRelationUserIds(c *gin.Context) {
	var req projectvo.GetProjectRelationUserIdsReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	res, err := service.GetProjectRelationUserIds(req.ProjectId, req.RelationType)
	response := projectvo.GetProjectRelationUserIdsResp{
		Err:     vo.NewErr(err),
		UserIds: res,
	}
	c.JSON(http.StatusOK, response)
}

func ProjectMemberIdList(c *gin.Context) {
	var req projectvo.ProjectMemberIdListReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	res, err := service.ProjectMemberIdList(req.OrgId, req.ProjectId, req.Data)
	response := projectvo.ProjectMemberIdListResp{
		Err:  vo.NewErr(err),
		Data: res,
	}
	c.JSON(http.StatusOK, response)
}

func GetProjectMemberIds(c *gin.Context) {
	var req projectvo.GetProjectMemberIdsReqVo
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	res, err := service.GetProjectMemberIds(req.OrgId, req.Input.ProjectId, req.Input.IncludeAdmin)
	response := projectvo.GetProjectMemberIdsResp{
		Err:  vo.NewErr(err),
		Data: res,
	}
	c.JSON(http.StatusOK, response)
}

func GetTrendsMembers(c *gin.Context) {
	var req projectvo.GetTrendListMembersReqVo
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	res, err := service.GetTrendsMembers(req.OrgId, req.UserId, req.Input.ProjectId)
	response := projectvo.GetTrendListMembersResp{
		Err:  vo.NewErr(err),
		Data: res,
	}
	c.JSON(http.StatusOK, response)
}
