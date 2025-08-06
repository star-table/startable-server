package api

import (
	"fmt"
	"github.com/gin-gonic/gin"
)

// RegisterRoutes 注册msgsvc的所有路由
// serviceName: 服务名称，如 "msgsvc"
// version: API版本，如 "v1"
func RegisterRoutes(r *gin.Engine, serviceName, version string) {
	// 构建API路径前缀，模拟mvc.BuildApi的逻辑
	apiPrefix := fmt.Sprintf("/api/%s/%s", serviceName, version)
	
	// 创建路由组
	apiGroup := r.Group(apiPrefix)
	
	// 注册GET路由
	apiGroup.GET("/health", Health)
	apiGroup.GET("/fixAddOrderForFeiShu", FixAddOrderForFeiShu)
	
	// 注册POST路由
	apiGroup.POST("/pushMsgToMq", PushMsgToMq)
	apiGroup.POST("/insertMqConsumeFailMsg", InsertMqConsumeFailMsg)
	apiGroup.POST("/getFailMsgList", GetFailMsgList)
	apiGroup.POST("/updateMsgStatus", UpdateMsgStatus)
	apiGroup.POST("/writeSomeFailedMsg", WriteSomeFailedMsg)
	apiGroup.POST("/sendLoginSMS", SendLoginSMS)
	apiGroup.POST("/sendMail", SendMail)
	apiGroup.POST("/sendSMS", SendSMS)
	
	// MQTT相关路由（已注释）
	// apiGroup.POST("/getMQTTChannelKey", GetMQTTChannelKey)
	// apiGroup.POST("/mqttPublish", MqttPublish)
}