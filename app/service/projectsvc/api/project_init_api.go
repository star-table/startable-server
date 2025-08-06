package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/star-table/startable-server/app/service/projectsvc/service"
	"github.com/star-table/startable-server/common/model/vo"
	"github.com/star-table/startable-server/common/model/vo/projectvo"
)

// ProjectInit 项目初始化
func ProjectInit(c *gin.Context) {
	var req projectvo.ProjectInitReqVo
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	respVo := projectvo.ProjectInitRespVo{}

	//_ = mysql.TransX(func(tx sqlbuilder.Tx) error {
	//	contextMap, err := service.ProjectInit(req.OrgId, tx)
	//
	//	respVo.ContextMap = contextMap
	//	respVo.Err = vo.NewErr(err)
	//
	//	return err
	//})

	c.JSON(http.StatusOK, respVo)
}

// CreateOrgDirectoryAppsAndViews 初始化组织的时候创建左侧目录应用、视图
func CreateOrgDirectoryAppsAndViews(c *gin.Context) {
	var reqVo projectvo.CreateOrgDirectoryAppsReq
	if err := c.ShouldBindJSON(&reqVo); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := service.CreateOrgDirectoryAppsAndViews(reqVo)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": &vo.Void{ID: 1},
	})
}
