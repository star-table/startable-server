package callsvc

import (
	"encoding/xml"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/star-table/startable-server/common/model/vo/orgvo"

	"github.com/star-table/startable-server/app/facade/orgfacade"
	consts2 "github.com/star-table/startable-server/app/service"

	"github.com/star-table/startable-server/common/core/config"

	"github.com/go-laoji/wecom-go-sdk/pkg/svr/logic"
	"github.com/go-laoji/wxbizmsgcrypt"

	"github.com/star-table/startable-server/common/core/consts"
	"github.com/star-table/startable-server/common/core/model"
	"github.com/star-table/startable-server/common/core/threadlocal"
	"github.com/gin-gonic/gin"
	"github.com/jtolds/gls"
)

type msgCallBackHandler struct {
	base.CallBackBase
	verifyCallBackHandler
}

var MsgCallBackHandler = &msgCallBackHandler{}

func (a *msgCallBackHandler) Handler(c *gin.Context) {
	httpContext := c.Request.Context().Value(consts.HttpContextKey).(model.HttpContext)
	threadlocal.Mgr.SetValues(gls.Values{consts.HttpContextKey: httpContext, consts.TraceIdKey: httpContext.TraceId}, func() {
		if c.Request.Method == "GET" {
			a.verify(c)
		} else {
			a.handlePost(c)
		}
	})
}

func (a *msgCallBackHandler) handlePost(c *gin.Context) {
	var params logic.EventPushQueryBinding
	if ok := c.ShouldBindQuery(&params); ok == nil {
		body, err := ioutil.ReadAll(c.Request.Body)
		log.Infof("[person_weixin] body:%v, params:%v, err:%v", string(body), params, err)
		if err != nil {
			return
		}
		var bizData logic.BizData
		xml.Unmarshal(body, &bizData)
		conf := config.GetConfig().PersonWeiXin
		wxcpt := wxbizmsgcrypt.NewWXBizMsgCrypt(conf.Token, conf.AesKey,
			conf.AppId, wxbizmsgcrypt.XmlType)
		if msg, err := wxcpt.DecryptMsg(params.MsgSign, params.Timestamp, params.Nonce, body); err != nil {
			log.Errorf("[person_weixin] handlerPost DecryptMsg url:%v, err:%v", c.Request.RequestURI, err)
			c.JSON(http.StatusOK, gin.H{"errno": 500, "errmsg": err.ErrMsg})
			return
		} else {
			event := &vo.Message{}
			xml.Unmarshal(msg, event)
			log.Infof("[person_weixin] plain:%v, object:%v", string(msg), event)

			switch strings.ToLower(event.Event) {
			case consts2.MsgTypeScan, consts2.MsgTypeSubscribe:
				eventKey := strings.Replace(event.EventKey, "qrscene_", "", 1)
				resp := orgfacade.PersonWeiXinQrCodeScan(orgvo.QrCodeScanReq{OpenId: event.FromUserName, SceneKey: eventKey})
				if resp.Failure() {
					log.Errorf("[person_weixin] PersonWeiXinQrCodeScan openId:%v, sceneKey:%v, err:%v", event.FromUserName, eventKey, resp.Error())
				}
			}
		}
	}
}
