package callsvc

import (
	"bytes"
	"crypto/sha1"
	"fmt"
	"net/http"
	"sort"

	"github.com/star-table/startable-server/common/core/logger"

	"github.com/star-table/startable-server/common/core/config"
	"github.com/star-table/startable-server/common/core/consts"
	"github.com/star-table/startable-server/common/core/model"
	"github.com/star-table/startable-server/common/core/threadlocal"
	"github.com/gin-gonic/gin"
	"github.com/jtolds/gls"
)

var log = logger.GetDefaultLogger()

type verifyCallBackHandler struct {
	base.CallBackBase
}

func (a *verifyCallBackHandler) verify(c *gin.Context) {
	httpContext := c.Request.Context().Value(consts.HttpContextKey).(model.HttpContext)
	threadlocal.Mgr.SetValues(gls.Values{consts.HttpContextKey: httpContext, consts.TraceIdKey: httpContext.TraceId}, func() {
		var params vo.Verify
		if ok := c.ShouldBindQuery(&params); ok == nil {
			sign := a.calSignature(params.Timestamp, params.Nonce)
			if sign == params.MsgSign {
				c.Writer.Write([]byte(params.EchoStr))
			} else {
				c.JSON(http.StatusLocked, gin.H{"err": "sign no equal", "echoStr": params.EchoStr})
			}
		} else {
			log.Errorf("[person_weixin] params error, url:%v", c.Request.RequestURI)
			c.JSON(http.StatusInternalServerError, gin.H{"errno": 500, "errmsg": "no echostr"})
		}
	})
}

func (a *verifyCallBackHandler) calSignature(timestamp, nonce string) string {
	conf := config.GetConfig().PersonWeiXin
	sortArr := []string{conf.Token, timestamp, nonce}
	sort.Strings(sortArr)
	var buffer bytes.Buffer
	for _, value := range sortArr {
		buffer.WriteString(value)
	}

	sha := sha1.New()
	sha.Write(buffer.Bytes())
	signature := fmt.Sprintf("%x", sha.Sum(nil))
	return signature
}
