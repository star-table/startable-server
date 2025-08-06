package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/star-table/startable-server/app/service/projectsvc/service"
	"github.com/star-table/startable-server/common/model/vo"
	"github.com/star-table/startable-server/common/model/vo/projectvo"
)

func ProjectTypes(c *gin.Context) {
	reqVo := &projectvo.ProjectTypesReqVo{}
	if err := c.ShouldBindJSON(reqVo); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	resp, err := service.ProjectTypes(reqVo.OrgId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, &projectvo.ProjectTypesRespVo{Err: vo.NewErr(err), ProjectTypes: resp})
}
