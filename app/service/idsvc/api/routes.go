package api

import (
	"github.com/gin-gonic/gin"
)

// RegisterRoutes 注册idsvc的所有路由
// 替代原来的mvc自动路由注册机制
func RegisterRoutes(r *gin.Engine, applicationName, version string) {
	// 构建API路径前缀
	apiPrefix := buildAPIPrefix(applicationName, version)
	
	// GET路由
	r.GET(apiPrefix+"/health", Health)
	
	// POST路由 - 所有ID生成相关接口都是POST请求
	r.POST(apiPrefix+"/applyPrimaryId", ApplyPrimaryId)
	r.POST(apiPrefix+"/applyCode", ApplyCode)
	r.POST(apiPrefix+"/applyMultipleId", ApplyMultipleId)
	r.POST(apiPrefix+"/applyMultiplePrimaryId", ApplyMultiplePrimaryId)
	r.POST(apiPrefix+"/applyMultiplePrimaryIdByCodes", ApplyMultiplePrimaryIdByCodes)
}

// buildAPIPrefix 构建API路径前缀
// 模拟mvc.BuildApi的逻辑: api/{applicationName}/{version}/{funcName}
func buildAPIPrefix(applicationName, version string) string {
	if version == "" {
		return "/api/" + applicationName
	}
	return "/api/" + applicationName + "/" + version
}