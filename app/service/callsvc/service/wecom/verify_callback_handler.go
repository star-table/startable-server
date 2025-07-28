package callsvc

import (
	"net/http"

	"github.com/star-table/startable-server/common/core/config"
	"github.com/star-table/startable-server/common/core/consts"
	"github.com/star-table/startable-server/common/core/model"
	"github.com/star-table/startable-server/common/core/threadlocal"
	"github.com/gin-gonic/gin"
	"github.com/go-laoji/wecom-go-sdk/pkg/svr/logic"
	"github.com/go-laoji/wxbizmsgcrypt"
	"github.com/jtolds/gls"
)

type verifyCallBackHandler struct {
	base.CallBackBase
}

func (a *verifyCallBackHandler) verify(c *gin.Context) {
	httpContext := c.Request.Context().Value(consts.HttpContextKey).(model.HttpContext)
	threadlocal.Mgr.SetValues(gls.Values{consts.HttpContextKey: httpContext, consts.TraceIdKey: httpContext.TraceId}, func() {
		conf := config.GetConfig().WeCom
		wxcpt := wxbizmsgcrypt.NewWXBizMsgCrypt(conf.SuiteToken, conf.SuiteAesKey,
			conf.CorpId, wxbizmsgcrypt.XmlType)
		var params logic.EventPushQueryBinding
		if ok := c.ShouldBindQuery(&params); ok == nil {
			echoStr, cryptErr := wxcpt.VerifyURL(params.MsgSign, params.Timestamp, params.Nonce, params.EchoStr)
			if nil != cryptErr {
				log.Errorf("[wecom] handlerGet url:%v, err:%v", c.Request.RequestURI, cryptErr)
				c.JSON(http.StatusLocked, gin.H{"err": cryptErr, "echoStr": echoStr})
			} else {
				c.Writer.Write(echoStr)
			}
		} else {
			log.Errorf("[wecom] params error, url:%v", c.Request.RequestURI)
			c.JSON(http.StatusInternalServerError, gin.H{"errno": 500, "errmsg": "no echostr"})
		}
	})
}
