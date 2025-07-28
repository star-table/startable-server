package handler

import (
	pushPb "gitea.bjx.cloud/LessCode/interface/golang/push/v1"
	"github.com/star-table/startable-server/app/facade/pushfacade"
	"github.com/star-table/startable-server/common/core/errs"
	"github.com/gin-gonic/gin"
	"github.com/spf13/cast"
)

type appHandlers struct{}

var AppHandler appHandlers

// GetMqttChannelKey
// @Security ApiKeyAuth
// @Summary 获取mqtt的key
// @Description 获取mqtt的key
// @Tags 应用
// @accept application/json
// @Produce application/json
// @Success 200 {object} MqttKeyResp
// @Failure 400
// @Router /api/rest/app/getMqttChannelKey [post]
func (appHandlers) GetMqttChannelKey(c *gin.Context) {

	cacheUserInfo, errSys := GetCacheUserInfo(c)
	if errSys != nil {
		Fail(c, errSys)
		return
	}

	type GetMqttChannelKeyReq struct {
		ProjectId   int `json:"projectId"`
		ChannelType int `json:"channelType"`
	}
	req := GetMqttChannelKeyReq{}
	err := c.BindJSON(&req)
	if err != nil {
		Fail(c, errs.ReqParamsValidateError)
		return
	}
	projectId, err1 := cast.ToInt64E(req.ProjectId)
	channelType, err2 := cast.ToInt64E(req.ChannelType)
	if err1 != nil || err2 != nil {
		Fail(c, errs.ReqParamsValidateError)
		return
	}
	resp := pushfacade.GenerateMqttKey(&pushPb.GenerateMqttKeyReq{
		ChannelType: int32(channelType),
		OrgId:       cacheUserInfo.OrgId,
		ProjectId:   projectId,
	})
	if resp.Failure() {
		Fail(c, resp.Error())
	} else {
		type MqttKeyResp struct {
			Address string
			Host    string
			Port    int
			Channel string
			Key     string
		}
		Success(c, &MqttKeyResp{
			Address: resp.Data.Address,
			Host:    resp.Data.Host,
			Port:    int(resp.Data.Port),
			Channel: resp.Data.Channel,
			Key:     resp.Data.Key,
		})
	}
}
