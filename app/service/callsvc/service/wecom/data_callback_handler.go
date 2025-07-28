package callsvc

import (
	"encoding/xml"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/star-table/startable-server/common/core/util/asyn"

	wework "github.com/go-laoji/wecom-go-sdk"

	platform_sdk "gitea.bjx.cloud/allstar/platform-sdk"
	consts3 "gitea.bjx.cloud/allstar/platform-sdk/consts"
	sdkVo "gitea.bjx.cloud/allstar/platform-sdk/vo"
	"github.com/star-table/startable-server/app/facade/orderfacade"
	"github.com/star-table/startable-server/app/facade/orgfacade"
	consts2 "github.com/star-table/startable-server/app/service"
	"github.com/star-table/startable-server/common/core/config"
	"github.com/star-table/startable-server/common/core/consts"
	consts4 "github.com/star-table/startable-server/common/core/consts"
	"github.com/star-table/startable-server/common/core/model"
	"github.com/star-table/startable-server/common/core/threadlocal"
	"github.com/star-table/startable-server/common/core/util/slice"
	int642 "github.com/star-table/startable-server/common/core/util/slice/int64"
	"github.com/star-table/startable-server/common/model/vo/ordervo"
	"github.com/star-table/startable-server/common/model/vo/orgvo"
	"github.com/gin-gonic/gin"
	"github.com/go-laoji/wecom-go-sdk/pkg/svr/logic"
	"github.com/go-laoji/wxbizmsgcrypt"
	"github.com/jtolds/gls"
	"github.com/spf13/cast"
)

type dataCallBackHandler struct {
	base.CallBackBase
	verifyCallBackHandler
}

var DataCallBackHandler = &dataCallBackHandler{}

func (a *dataCallBackHandler) Handler(c *gin.Context) {
	httpContext := c.Request.Context().Value(consts.HttpContextKey).(model.HttpContext)
	threadlocal.Mgr.SetValues(gls.Values{consts.HttpContextKey: httpContext, consts.TraceIdKey: httpContext.TraceId}, func() {
		if c.Request.Method == "GET" {
			a.verify(c)
		} else {
			a.handlerPost(c)
		}
	})
}

func (a *dataCallBackHandler) handlerPost(c *gin.Context) {
	conf := config.GetConfig().WeCom
	var params logic.EventPushQueryBinding
	if ok := c.ShouldBindQuery(&params); ok == nil {
		body, err := ioutil.ReadAll(c.Request.Body)
		if err != nil {
			log.Errorf("[wecom] data handlerPost url:%v, err:%v", c.Request.RequestURI, err)
			c.JSON(http.StatusOK, gin.H{"errno": 500, "errmsg": err.Error()})
			return
		} else {
			var bizData logic.BizData
			xml.Unmarshal(body, &bizData)
			wxcpt := wxbizmsgcrypt.NewWXBizMsgCrypt(conf.SuiteToken, conf.SuiteAesKey,
				bizData.ToUserName, wxbizmsgcrypt.XmlType)
			if msg, err := wxcpt.DecryptMsg(params.MsgSign, params.Timestamp, params.Nonce, body); err != nil {
				log.Errorf("[wecom] data handlerPost DecryptMsg url:%v, err:%v", c.Request.RequestURI, err)
				c.JSON(http.StatusOK, gin.H{"errno": 500, "errmsg": err.ErrMsg})
				return
			} else {
				log.Infof("[wecom] data callback get msg:%v", string(msg))
				event := &vo.DataEvent{}
				xml.Unmarshal(msg, event)
				switch event.InfoType {
				//case consts2.EventTypeContact:
				//	switch event.ChangeType {
				//	case consts2.EventChangeTypeCreateUser, consts2.EventChangeTypeDeleteUser, consts2.EventChangeTypeUpdateUser:
				//		contactEvent := &vo.DataUserEvent{}
				//		xml.Unmarshal(msg, contactEvent)
				//		a.handlerUser(contactEvent)
				//	case consts2.EventChangeTypeCreateDept, consts2.EventChangeTypeDeleteDept, consts2.EventChangeTypeUpdateDept:
				//		deptEvent := &vo.DataDeptEvent{}
				//		xml.Unmarshal(msg, deptEvent)
				//		a.handlerDept(deptEvent)
				//	}
				case consts2.InfoTypePaySuccess, consts2.InfoTypePayRefund:
					a.handlerPay(msg)
				case consts2.InfoTypeRegisterCorp:
					a.handlerRegisterCorp(msg)
				}

				data := &vo.ChangeAppAdmin{}
				xml.Unmarshal(msg, data)
				if data.Event == consts2.EventChangeAppAdmin {
					a.handlerChangeAdmin(data)
				}

				c.Writer.WriteString("success")
			}
		}
	} else {
		log.Errorf("[wecom] params error, url:%v", c.Request.RequestURI)
		c.JSON(http.StatusOK, gin.H{"errno": 400, "errmsg": ok.Error()})
	}
}

//func (a *dataCallBackHandler) handlerUser(event *vo.DataUserEvent) {
//	var err error
//	switch event.ChangeType {
//	case consts2.EventChangeTypeCreateUser:
//		err = a.UserAdd(consts3.SourceChannelWeixin, event.AuthCorpId, event.UserID, cast.ToStringSlice(event.Department)...)
//	case consts2.EventChangeTypeUpdateUser:
//		err = a.UserUpdate(consts3.SourceChannelWeixin, event.AuthCorpId, event.UserID, event.NewUserID, cast.ToStringSlice(event.Department)...)
//	case consts2.EventChangeTypeDeleteUser:
//		err = a.UserLeave(consts3.SourceChannelWeixin, event.AuthCorpId, event.UserID)
//	}
//	if err != nil {
//		log.Errorf("[handlerUser] event:%v, error:%v", event, err)
//	}
//}
//
//func (a *dataCallBackHandler) handlerDept(event *vo.DataDeptEvent) {
//	eventType := ""
//	switch event.ChangeType {
//	case consts2.EventChangeTypeCreateDept:
//		eventType = consts4.EventDeptAdd
//	case consts2.EventChangeTypeUpdateDept:
//		eventType = consts4.EventDeptUpdate
//	case consts2.EventChangeTypeDeleteDept:
//		eventType = consts4.EventDeptDel
//	}
//	err := a.HandleDeptChange(consts3.SourceChannelWeixin, event.AuthCorpId, eventType, cast.ToString(event.Id))
//	if err != nil {
//		log.Errorf("[handlerDept] event:%v, error:%v", event, err)
//	}
//}

func (a *dataCallBackHandler) handlerPay(data []byte) {
	var event vo.PayHandle
	if err := xml.Unmarshal(data, &event); err != nil {
		log.Errorf("[wecom] handlerPay unmarshal error::%v", err)
		return
	}
	switch event.InfoType {
	case consts2.InfoTypePaySuccess:
		resp := orderfacade.CreateWeiXinLicenceOrder(ordervo.CreateWeiXinLicenceOrderReq{
			SuitId:  event.ServiceCorpId,
			CorpId:  event.AuthCorpId,
			OrderId: event.OrderId,
		})
		if resp.Failure() {
			log.Errorf("[CreateWeiXinLicenceOrder] params:%v, err:%v", event, resp.Err)
		}
	case consts2.InfoTypePayRefund:
	}
}

func (a *dataCallBackHandler) handlerRegisterCorp(data []byte) {
	asyn.Execute(func() {
		time.Sleep(10 * time.Second)

		var event vo.RegisterFinish
		if err := xml.Unmarshal(data, &event); err != nil {
			log.Errorf("[wecom] handlerRegisterCorp unmarshal error::%v", err)
			return
		}

		sdk, err := platform_sdk.GetClient(consts3.SourceChannelWeixin, "")
		if err != nil {
			log.Errorf("[handlerRegisterCorp] GetClient error::%v", err)
			return
		}
		weWork := sdk.GetOriginClient().(wework.IWeWork)
		resp := weWork.ContactSyncSuccess(event.ContactSync.AccessToken)
		if resp.ErrCode != 0 {
			log.Errorf("[handlerRegisterCorp] ContactSyncSuccess error::%v, token:%v", err, event.ContactSync.AccessToken)
		}
	})
}

func (a *dataCallBackHandler) handlerChangeAdmin(data *vo.ChangeAppAdmin) {
	sdk, err := platform_sdk.GetClient(consts3.SourceChannelWeixin, data.ToUserName)
	if err != nil {
		log.Errorf("[handlerChangeAdmin] GetClient error::%v, corpId:%v", err, data.ToUserName)
		return
	}
	// 查询应用管理员列表 与现在的超管比对
	adminList, sdkError := sdk.GetAdminList(&sdkVo.GetAdminListReq{
		CorpId:  data.ToUserName,
		AgentId: cast.ToInt(data.AgentID),
	})
	if sdkError != nil {
		log.Errorf("[handlerChangeAdmin] GetAdminList err:%v, corpId:%v", sdkError, data.ToUserName)
		return
	}
	adminOpenIds := adminList.OpenIds
	baseOrgInfo := orgfacade.GetBaseOrgInfoByOutOrgId(orgvo.GetBaseOrgInfoByOutOrgIdReqVo{OutOrgId: data.ToUserName})
	if baseOrgInfo.Failure() {
		log.Errorf("[handlerChangeAdmin] GetBaseOrgInfoByOutOrgId err:%v, corpId:%v", baseOrgInfo.Message, data.ToUserName)
		return
	}
	orgId := baseOrgInfo.BaseOrgInfo.OrgId
	// 获取超管的信息
	superAdminInfo := orgfacade.GetOrgSuperAdminInfo(orgvo.GetOrgSuperAdminInfoReq{OrgId: orgId})
	if superAdminInfo.Failure() {
		log.Errorf("[handlerChangeAdmin] GetOrgSuperAdminInfo err:%v, corpId:%v", superAdminInfo.Message, data.ToUserName)
		return
	}
	currAdminIdMap := map[string]int64{}
	currAdminUserIds := []int64{}
	currAdminOpenIds := []string{}
	for _, u := range superAdminInfo.Data {
		currAdminIdMap[u.OpenId] = u.UserId
		currAdminUserIds = append(currAdminUserIds, u.UserId)
		currAdminOpenIds = append(currAdminOpenIds, u.OpenId)
	}
	allOpenIds := []string{}
	allOpenIds = append(allOpenIds, adminOpenIds...)
	allOpenIds = append(allOpenIds, currAdminOpenIds...)
	allOpenIds = slice.SliceUniqueString(allOpenIds)
	outUserResp := orgfacade.GetOrgUserIdsByEmIds(orgvo.GetOrgUserIdsByEmIdsReq{
		OrgId:         orgId,
		SourceChannel: consts3.SourceChannelWeixin,
		EmpIds:        allOpenIds,
	})
	if outUserResp.Failure() {
		log.Errorf("[handlerChangeAdmin] GetOrgUserIdsByEmIds err:%v, corpId:%v", outUserResp.Message, data.ToUserName)
		return
	}
	allAdminUserIdsMap := outUserResp.Data
	// 对比更改前后的管理员
	_, add, del := int642.CompareSliceAddDelString(adminOpenIds, currAdminOpenIds)

	updateAdminIds := []int64{}
	updateType := 0
	if len(add) > 0 && len(del) == 0 {
		// 添加超管
		//addAdminIds := []int64{}
		for _, outUserId := range add {
			if addId, ok := allAdminUserIdsMap[outUserId]; ok {
				updateAdminIds = append(updateAdminIds, addId)
			}
		}
		updateType = consts4.AddType
	}

	if len(del) > 0 && len(add) == 0 {
		// 删除超管
		//delAdminIds := []int64{}
		for _, outUserId := range del {
			if delId, ok := allAdminUserIdsMap[outUserId]; ok {
				updateAdminIds = append(updateAdminIds, delId)
			}
		}
		updateType = consts4.DelType
	}
	if len(updateAdminIds) > 0 {
		resp := orgfacade.UpdateUserToSysManageGroup(orgvo.UpdateUserToSysManageGroupReq{
			OrgId: orgId,
			Input: orgvo.UpdateUserToSysManageGroupData{
				UserIds:    updateAdminIds,
				UpdateType: updateType,
			},
		})
		if resp.Failure() {
			log.Errorf("[handlerChangeAdmin] UpdateUserToSysManageGroup err:%v, corpId:%v", resp.Message, data.ToUserName)
			return
		}
	}

}
