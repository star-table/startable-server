package api

import (
	"net/http"
	"github.com/gin-gonic/gin"
	"github.com/star-table/startable-server/common/core/logger"
)

var log = logger.GetDefaultLogger()

// Health 健康检查接口
func Health(c *gin.Context) {
	c.JSON(http.StatusOK, "ok")
}