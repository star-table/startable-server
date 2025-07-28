package callsvc

import (
	"encoding/xml"
	"io/ioutil"
	"net/http"
	"time"

	platform_sdk "gitea.bjx.cloud/allstar/platform-sdk"
	sdk_const "gitea.bjx.cloud/allstar/platform-sdk/consts"
	"github.com/star-table/startable-server/app/facade/orderfacade"
	"github.com/star-table/startable-server/app/facade/orgfacade"
	consts2 "github.com/star-table/startable-server/app/service"
	"github.com/star-table/startable-server/common/core/config"
	"github.com/star-table/startable-server/common/core/consts"
	consts3 "github.com/star-table/startable-server/common/core/consts"
	"github.com/star-table/startable-server/common/core/model"
	"github.com/star-table/startable-server/common/core/threadlocal"
	"github.com/star-table/startable-server/common/core/util/asyn"
	"github.com/star-table/startable-server/common/core/util/json"
	"github.com/star-table/startable-server/common/core/util/uuid"
	"github.com/star-table/startable-server/common/library/cache"
	"github.com/star-table/startable-server/common/model/bo"
	"github.com/star-table/startable-server/common/model/vo/ordervo"
	"github.com/star-table/startable-server/common/model/vo/orgvo"
	"github.com/gin-gonic/gin"
	wework "github.com/go-laoji/wecom-go-sdk"
	"github.com/go-laoji/wecom-go-sdk/pkg/svr/logic"
	"github.com/go-laoji/wxbizmsgcrypt"
	"github.com/jtolds/gls"
	"github.com/spf13/cast"
)

type cmdCallBackHandler struct {
	base.CallBackBase
	verifyCallBackHandler
}

var CmdCallBackHandler = &cmdCallBackHandler{}

func (a *cmdCallBackHandler) Handler(c *gin.Context) {
	httpContext := c.Request.Context().Value(consts.HttpContextKey).(model.HttpContext)
	threadlocal.Mgr.SetValues(gls.Values{consts.HttpContextKey: httpContext, consts.TraceIdKey: httpContext.TraceId}, func() {
		if c.Request.Method == "GET" {
			a.verify(c)
		} else {
			a.handlerPost(c)
		}
	})
}

func (a *cmdCallBackHandler) handlerPost(c *gin.Context) {
	conf := config.GetConfig().WeCom
	var params logic.EventPushQueryBinding
	if ok := c.ShouldBindQuery(&params); ok == nil {
		body, err := ioutil.ReadAll(c.Request.Body)
		log.Infof("[wecom] cmd handlerPost data: %v", string(body))
		if err != nil {
			log.Errorf("[wecom] cmd handlerPost url:%v, err:%v", c.Request.RequestURI, err)
			c.JSON(http.StatusOK, gin.H{"errno": 500, "errmsg": err.Error()})
			return
		} else {
			wxcpt := wxbizmsgcrypt.NewWXBizMsgCrypt(conf.SuiteToken, conf.SuiteAesKey,
				conf.SuiteId, wxbizmsgcrypt.XmlType)
			if msg, err := wxcpt.DecryptMsg(params.MsgSign, params.Timestamp, params.Nonce, body); err != nil {
				log.Errorf("[wecom] cmd handlerPost DecryptMsg url:%v, err:%v", c.Request.RequestURI, err)
				c.JSON(http.StatusOK, gin.H{"errno": 500, "errmsg": err.ErrMsg})
				return
			} else {
				log.Infof("[wecom] cmd callback get msg:%v", string(msg))

				var bizEvent vo.DataEvent
				if e := xml.Unmarshal(msg, &bizEvent); e != nil {
					c.JSON(http.StatusOK, gin.H{"errno": 500, "errmsg": e})
					return
				}
				switch bizEvent.InfoType {
				case logic.SuiteTicket:
					a.setTicket(msg)
				case logic.CreateAuth:
					// 服务商的响应必须在1000ms内完成，以保证用户安装应用的体验。
					// 建议在接收到此事件时，先记录下AuthCode，并立即回应企业微信，之后再做相关业务的处理。
					a.createAuth(msg)
				case logic.ChangeAuth, logic.CorpArchAuth:
					a.changeAuth(msg)
				case logic.CancelAuth:
				case logic.ResetPermanentCode:

				case consts2.EventTypeContact:
					switch bizEvent.ChangeType {
					case consts2.EventChangeTypeCreateUser, consts2.EventChangeTypeDeleteUser, consts2.EventChangeTypeUpdateUser:
						contactEvent := &vo.DataUserEvent{}
						xml.Unmarshal(msg, contactEvent)
						a.handlerUser(contactEvent)
					case consts2.EventChangeTypeCreateDept, consts2.EventChangeTypeDeleteDept, consts2.EventChangeTypeUpdateDept:
						deptEvent := &vo.DataDeptEvent{}
						xml.Unmarshal(msg, deptEvent)
						a.handlerDept(deptEvent)
					}
				case consts2.EventOpenOrder:
				case consts2.EventChangeOrder:
				case consts2.EventPayForAppSuccess:
					a.handlerOrderPayForSuccess(msg)

				case consts2.EventChangeEdition:
				case consts2.EventRefund:

				}
				c.Writer.WriteString("success")
			}
		}
	} else {
		log.Errorf("[wecom] params error, url:%v", c.Request.RequestURI)
		c.JSON(http.StatusOK, gin.H{"errno": 400, "errmsg": ok.Error()})
	}
}

func (a *cmdCallBackHandler) setTicket(data []byte) {
	var suiteEvent logic.SuiteTicketEvent
	if err := xml.Unmarshal(data, &suiteEvent); err != nil {
		log.Errorf("[wecom] setTicket error::%v", err)
	} else {
		err = platform_sdk.GetPlatformInfo(sdk_const.SourceChannelWeixin).Cache.SetTicket(suiteEvent.SuiteTicket)
		if err != nil {
			log.Errorf("[wecom] setTicket error::%v", err)
		}
	}
}

func (a *cmdCallBackHandler) createAuth(data []byte) {
	var event logic.CreateAuthEvent
	if err := xml.Unmarshal(data, &event); err != nil {
		log.Errorf("[wecom] createAuth error::%v", err)
		return
	}
	sdk, err := platform_sdk.GetClient(sdk_const.SourceChannelWeixin, "")
	if err != nil {
		log.Errorf("[wecom] GetClient error::%v", err)
		return
	}
	client := sdk.GetOriginClient().(wework.IWeWork)
	resp := client.GetPermanentCode(event.AuthCode)
	if resp.ErrCode != 0 {
		log.Errorf("[wecom] GetPermanentCode code:%v, msg:%v", resp.ErrCode, resp.ErrorMsg)
		return
	}
	log.Infof("[wecom]createAuth GetPermanentCode 获取永久授权码信息:%v, 永久授权码:%s", resp, resp.PermanentCode)

	corpId := resp.AuthCorpInfo.CorpId
	uid := uuid.NewUuid()
	lockKey := consts3.WeiXinTalkCorpInitKey + corpId
	suc, err := cache.TryGetDistributedLock(lockKey, uid)
	if err != nil {
		log.Errorf("企业%s初始化时，获取分布式锁异常:%v", corpId, err)
		return
	}
	if suc {
		defer func() {
			_, e := cache.ReleaseDistributedLock(lockKey, uid)
			if e != nil {
				log.Error(e)
			}
		}()
	} else {
		log.Infof("企业%s正在初始化中，当前请求无效", corpId)
		return
	}

	asyn.Execute(func() {
		//开始初始化组织
		orgInitResp := orgfacade.InitOrg(orgvo.InitOrgReqVo{
			InitOrg: bo.InitOrgBo{
				OutOrgId:      resp.AuthCorpInfo.CorpId,
				OrgName:       resp.AuthCorpInfo.CorpName,
				SourceChannel: sdk_const.SourceChannelWeixin,
				PermanentCode: resp.PermanentCode,
				TenantCode:    cast.ToString(resp.AuthInfo.Agent[0].AgentId),
			},
		})
		if orgInitResp.Failure() {
			log.Error(orgInitResp.Message)
		}
	})

	// 规定要1秒内返回，所以只停留500ms，免得超时
	time.Sleep(500 * time.Millisecond)
}

func (a *cmdCallBackHandler) changeAuth(data []byte) {
	var event vo.ChangeAuth
	if err := xml.Unmarshal(data, &event); err != nil {
		log.Errorf("[wecom] changeAuth unmarshal error::%v", err)
		return
	}

	err := a.ChangeAuthScope(event.AuthCorpId, sdk_const.SourceChannelWeixin)
	if err != nil {
		log.Errorf("[wecom] ChangeAuthScope error::%v, corpId:%v", err, event.AuthCorpId)
	}
}

func (a *cmdCallBackHandler) handlerUser(event *vo.DataUserEvent) {
	var err error
	switch event.ChangeType {
	case consts2.EventChangeTypeCreateUser:
		err = a.UserAdd(sdk_const.SourceChannelWeixin, event.AuthCorpId, event.UserID, cast.ToStringSlice(event.Department)...)
	case consts2.EventChangeTypeUpdateUser:
		err = a.UserUpdate(sdk_const.SourceChannelWeixin, event.AuthCorpId, event.UserID, event.NewUserID, cast.ToStringSlice(event.Department)...)
	case consts2.EventChangeTypeDeleteUser:
		err = a.UserLeave(sdk_const.SourceChannelWeixin, event.AuthCorpId, event.UserID)
	}
	if err != nil {
		log.Errorf("[handlerUser] event:%v, error:%v", event, err)
	}
}

func (a *cmdCallBackHandler) handlerDept(event *vo.DataDeptEvent) {
	eventType := ""
	switch event.ChangeType {
	case consts2.EventChangeTypeCreateDept:
		eventType = consts3.EventDeptAdd
	case consts2.EventChangeTypeUpdateDept:
		eventType = consts3.EventDeptUpdate
	case consts2.EventChangeTypeDeleteDept:
		eventType = consts3.EventDeptDel
	}
	err := a.HandleDeptChange(sdk_const.SourceChannelWeixin, event.AuthCorpId, eventType, cast.ToString(event.Id))
	if err != nil {
		log.Errorf("[handlerDept] event:%v, error:%v", event, err)
	}
}

// 支订单付成功回调
func (a *cmdCallBackHandler) handlerOrderPayForSuccess(data []byte) {
	var event vo.PayForAppSuccess
	if err := xml.Unmarshal(data, &event); err != nil {
		log.Errorf("[wecom] handlerOrderPayForSuccess unmarshal error::%v", err)
		return
	}

	log.Infof("[wecom]handlerOrderPayForSuccess info:%v", string(data))

	// 根据orderId查询订单详情，插入订单表
	client, err := platform_sdk.GetClient(sdk_const.SourceChannelWeixin, event.PaidCorpId)
	if err != nil {
		log.Errorf("[handlerOrderPayForSuccess] GetClient corpId:%v, err:%v", event.PaidCorpId, err)
		return
	}

	weWork := client.GetOriginClient().(wework.IWeWork)
	orderInfo := weWork.GetOrderInfo(event.OrderId)
	if orderInfo.ErrCode != 0 {
		log.Errorf("[handlerOrderPayForSuccess] GetOrderInfo err: %v, code:%v", orderInfo.ErrorMsg, orderInfo.ErrCode)
		return
	}
	log.Infof("[wecom] orderInfo:%v", json.ToJsonIgnoreError(orderInfo))
	if orderInfo.OrderStatus != consts3.WeiXinOrderStatusPaySuccess {
		return
	}

	// 插入微信订单表
	orderBo := bo.OrderWeiXinBo{
		OutOrgId:      event.PaidCorpId,
		OrderId:       event.OrderId,
		EditionId:     orderInfo.EditionId,
		EditionName:   orderInfo.EditionName,
		OrderType:     orderInfo.OrderType,
		OrderStatus:   orderInfo.OrderStatus,
		UserCount:     orderInfo.UserCount,
		OrderPeriod:   orderInfo.OrderPeriod,
		OrderPayPrice: orderInfo.Price,
		PaidTime:      time.Unix(orderInfo.PaidTime, 0),
		BeginTime:     time.Unix(orderInfo.BeginTime, 0),
		EndTime:       time.Unix(orderInfo.EndTime, 0),
	}
	resp := orderfacade.AddWeiXinOrder(ordervo.AddWeiXinOrderReq{Data: orderBo})
	if resp.Failure() {
		log.Error(resp.Error())
		return
	}

	// 如果是企业版和旗舰版 下单购买接口调用许可
	level := consts3.GetWeiXinOrderLevel(orderInfo.EditionId)
	if level == consts3.PayLevelFlagship || level == consts3.PayLevelEnterprise {
		weComConfig := config.GetConfig().WeCom
		buyerUserId := weComConfig.BuyerUserId
		reply := weWork.CreateNewOrder(wework.CreateOrderRequest{
			CorpId:      event.PaidCorpId,
			BuyerUserid: buyerUserId,
			AccountCount: wework.AccountCount{
				BaseCount: orderInfo.UserCount,
			},
			AccountDuration: wework.AccountDuration{
				Days: orderInfo.OrderPeriod,
			},
		})
		if reply.ErrCode != 0 {
			log.Errorf("[handlerOrderPayForSuccess] CreateNewOrder err:%v, errCode:%v", reply.ErrorMsg, reply.ErrCode)
		} else {
			//licenceOrder := orderfacade.CreateWeiXinLicenceOrder(ordervo.CreateWeiXinLicenceOrderReq{
			//	SuitId:  weComConfig.SuiteId,
			//	CorpId:  event.PaidCorpId,
			//	OrderId: reply.OrderId,
			//})
			//if licenceOrder.Failure() {
			//	log.Errorf("[handlerOrderPayForSuccess] CreateWeiXinLicenceOrder licenceOrderId:%v, err:%v, errCode:%v",
			//		reply.OrderId, reply.ErrorMsg, reply.ErrCode)
			//}
		}
	}
}
