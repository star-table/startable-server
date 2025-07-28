package callsvc

import (
	"encoding/xml"
	"io/ioutil"
	"net/http"

	"github.com/spf13/cast"

	consts4 "github.com/star-table/startable-server/common/core/consts"

	consts3 "gitea.bjx.cloud/allstar/platform-sdk/consts"

	consts2 "github.com/star-table/startable-server/app/service"

	"github.com/star-table/startable-server/common/core/config"
	"github.com/star-table/startable-server/common/core/consts"
	"github.com/star-table/startable-server/common/core/logger"
	"github.com/star-table/startable-server/common/core/model"
	"github.com/star-table/startable-server/common/core/threadlocal"
	"github.com/gin-gonic/gin"
	"github.com/go-laoji/wecom-go-sdk/pkg/svr/logic"
	"github.com/go-laoji/wxbizmsgcrypt"
	"github.com/jtolds/gls"
)

var log = logger.GetDefaultLogger()

type customCallBackHandler struct {
	base.CallBackBase
	verifyCallBackHandler
}

var CustomCallBackHandler = &customCallBackHandler{}

func (a *customCallBackHandler) Handler(c *gin.Context) {
	httpContext := c.Request.Context().Value(consts.HttpContextKey).(model.HttpContext)
	threadlocal.Mgr.SetValues(gls.Values{consts.HttpContextKey: httpContext, consts.TraceIdKey: httpContext.TraceId}, func() {
		if c.Request.Method == "GET" {
			a.verify(c)
		} else {
			a.handlerPost(c)
		}
	})
}

func (a *customCallBackHandler) handlerPost(c *gin.Context) {
	conf := config.GetConfig().WeCom
	var params logic.EventPushQueryBinding
	if ok := c.ShouldBindQuery(&params); ok == nil {
		body, err := ioutil.ReadAll(c.Request.Body)
		if err != nil {
			log.Errorf("[wecom] custom handlerPost url:%v, err:%v", c.Request.RequestURI, err)
			c.JSON(http.StatusOK, gin.H{"errno": 500, "errmsg": err.Error()})
			return
		} else {
			var bizData logic.BizData
			xml.Unmarshal(body, &bizData)
			wxcpt := wxbizmsgcrypt.NewWXBizMsgCrypt(conf.SuiteToken, conf.SuiteAesKey,
				bizData.ToUserName, wxbizmsgcrypt.XmlType)
			if msg, err := wxcpt.DecryptMsg(params.MsgSign, params.Timestamp, params.Nonce, body); err != nil {
				log.Errorf("[wecom] custom handlerPost DecryptMsg url:%v, err:%v", c.Request.RequestURI, err)
				c.JSON(http.StatusOK, gin.H{"errno": 500, "errmsg": err.ErrMsg})
				return
			} else {
				log.Infof("[wecom] custom callback get msg:%v", string(msg))
				event := &vo.CustomEvent{}
				xml.Unmarshal(msg, event)
				switch event.Event {
				case consts2.EventTypeContact:
					switch event.ChangeType {
					case consts2.EventChangeTypeCreateUser, consts2.EventChangeTypeDeleteUser, consts2.EventChangeTypeUpdateUser:
						contactEvent := &vo.CustomUserEvent{}
						xml.Unmarshal(msg, contactEvent)
						a.handlerUser(contactEvent)
					case consts2.EventChangeTypeCreateDept, consts2.EventChangeTypeDeleteDept, consts2.EventChangeTypeUpdateDept:
						deptEvent := &vo.CustomDeptEvent{}
						xml.Unmarshal(msg, deptEvent)
						a.handlerDept(deptEvent)
					}
				}

				c.Writer.WriteString("success")
			}
		}
	} else {
		log.Errorf("[wecom] params error, url:%v", c.Request.RequestURI)
		c.JSON(http.StatusOK, gin.H{"errno": 400, "errmsg": ok.Error()})
	}
}

func (a *customCallBackHandler) handlerUser(event *vo.CustomUserEvent) {
	var err error
	switch event.ChangeType {
	case consts2.EventChangeTypeCreateUser:
		err = a.UserAdd(consts3.SourceChannelWeixin, event.ToUserName, event.UserID, cast.ToStringSlice(event.Department)...)
	case consts2.EventChangeTypeUpdateUser:
		err = a.UserUpdate(consts3.SourceChannelWeixin, event.ToUserName, event.UserID, event.NewUserID, cast.ToStringSlice(event.Department)...)
	case consts2.EventChangeTypeDeleteUser:
		err = a.UserLeave(consts3.SourceChannelWeixin, event.ToUserName, event.UserID)
	}
	if err != nil {
		log.Errorf("[handlerUser] event:%v, error:%v", event, err)
	}
}

func (a *customCallBackHandler) handlerDept(event *vo.CustomDeptEvent) {
	eventType := ""
	switch event.ChangeType {
	case consts2.EventChangeTypeCreateDept:
		eventType = consts4.EventDeptAdd
	case consts2.EventChangeTypeUpdateDept:
		eventType = consts4.EventDeptUpdate
	case consts2.EventChangeTypeDeleteDept:
		eventType = consts4.EventDeptDel
	}
	err := a.HandleDeptChange(consts3.SourceChannelWeixin, event.ToUserName, eventType, cast.ToString(event.Id))
	if err != nil {
		log.Errorf("[handlerDept] event:%v, error:%v", event, err)
	}
}
