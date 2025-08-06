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

// RegisterRoutes 注册commonsvc的所有路由
// 替代原来的mvc自动路由注册机制
// 注意：需要确保所有处理函数都已在对应的文件中定义
func RegisterRoutes(r *gin.Engine, applicationName, version string) {
	// 构建API路径前缀
	apiPrefix := buildAPIPrefix(applicationName, version)

	// GET路由 - 使用小写首字母的函数名（符合mvc.BuildApi中str.LcFirst的逻辑）

	r.GET(apiPrefix+"/health", Health)
	r.GET(apiPrefix+"/industryList", IndustryList)

	// POST路由 - 使用小写首字母的函数名
	r.POST(apiPrefix+"/areaLinkageList", AreaLinkageList)
	r.POST(apiPrefix+"/areaInfo", AreaInfo)
	r.POST(apiPrefix+"/uploadOssByFsImageKey", UploadOssByFsImageKey)
}

// buildAPIPrefix 构建API路径前缀
// 模拟mvc.BuildApi的逻辑: api/{applicationName}/{version}/{funcName}
func buildAPIPrefix(applicationName, version string) string {
	if version == "" {
		return "/api/" + applicationName
	}
	return "/api/" + applicationName + "/" + version
}
