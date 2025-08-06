package api

import (
	"net/http"
	"github.com/gin-gonic/gin"
	"github.com/star-table/startable-server/app/service/projectsvc/service"
	"github.com/star-table/startable-server/common/model/vo"
	"github.com/star-table/startable-server/common/model/vo/projectvo"
)

func AddRecycleRecord(c *gin.Context) {
	var req projectvo.AddRecycleRecordReqVo
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	res, err := service.AddRecycleRecord(req.OrgId, req.UserId, req.Input.ProjectId, req.Input.RelationIds, req.Input.RelationType)
	response := projectvo.AddRecycleRecordRespVo{Data: res, Err: vo.NewErr(err)}
	c.JSON(http.StatusOK, response)
}

func GetRecycleList(c *gin.Context) {
	var req projectvo.GetRecycleListReqVo
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	res, err := service.GetRecycleList(req.OrgId, req.UserId, req.Input.ProjectID, req.Input.RelationType, req.Page, req.Size)
	response := projectvo.GetRecycleListRespVo{
		Err:  vo.NewErr(err),
		Data: res,
	}
	c.JSON(http.StatusOK, response)
}

func CompleteDelete(c *gin.Context) {
	var req projectvo.CompleteDeleteReqVo
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	res, err := service.CompleteDelete(req.OrgId, req.UserId, req.Input.ProjectID, req.Input.RecycleID, req.Input.RelationID, req.Input.RelationType)
	response := vo.CommonRespVo{
		Err:  vo.NewErr(err),
		Void: res,
	}
	c.JSON(http.StatusOK, response)
}

func RecoverRecycleBin(c *gin.Context) {
	var req projectvo.CompleteDeleteReqVo
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	res, err := service.RecoverRecycleBin(req.OrgId, req.UserId, req.Input.ProjectID, req.Input.RecycleID, req.Input.RelationID, req.Input.RelationType, req.SourceChannel)
	response := vo.CommonRespVo{
		Err:  vo.NewErr(err),
		Void: res,
	}
	c.JSON(http.StatusOK, response)
}
