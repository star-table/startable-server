package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Health 健康检查接口
func Health(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}
