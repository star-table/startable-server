package api

// 注释掉的MQTT相关接口，如需要可以取消注释并重构
//func GetMQTTChannelKey(c *gin.Context) {
//	var req msgvo.GetMQTTChannelKeyReqVo
//	if err := c.ShouldBindJSON(&req); err != nil {
//		response := msgvo.GetMQTTChannelKeyRespVo{Err: vo.NewErr(err)}
//		c.JSON(http.StatusOK, response)
//		return
//	}
//	
//	res, err := msgsvcService.GetMQTTChannelKey(req.OrgId, req.UserId, req.Input)
//	response := msgvo.GetMQTTChannelKeyRespVo{Err: vo.NewErr(err), Data: res}
//	c.JSON(http.StatusOK, response)
//}

//func MqttPublish(c *gin.Context) {
//	var req msgvo.MqttPublishReqVo
//	if err := c.ShouldBindJSON(&req); err != nil {
//		response := vo.VoidErr{Err: vo.NewErr(err)}
//		c.JSON(http.StatusOK, response)
//		return
//	}
//	
//	err := msgsvcService.MqttPublish(&req)
//	response := vo.VoidErr{Err: vo.NewErr(err)}
//	c.JSON(http.StatusOK, response)
//}
