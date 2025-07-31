package orgsvc

import (
	"strconv"
	"time"

	"github.com/spf13/cast"

	"gitea.bjx.cloud/allstar/platform-sdk/sdk_interface"

	sdk_const "gitea.bjx.cloud/allstar/platform-sdk/consts"

	platform_sdk "gitea.bjx.cloud/allstar/platform-sdk"
	platformVo "gitea.bjx.cloud/allstar/platform-sdk/vo"
	"github.com/star-table/startable-server/app/facade/idfacade"
	"github.com/star-table/startable-server/app/service/orgsvc/po"
	"github.com/star-table/startable-server/common/core/consts"
	"github.com/star-table/startable-server/common/core/errs"
	"github.com/star-table/startable-server/common/core/util/json"
	"github.com/star-table/startable-server/common/core/util/strs"
	"github.com/star-table/startable-server/common/library/db/mysql"
	"github.com/star-table/startable-server/common/model/bo"
	"github.com/star-table/startable-server/common/model/vo/ordervo"
	"upper.io/db.v3/lib/sqlbuilder"
)

// 留着做对比方便的注释
func InitOrg(initOrgBo bo.InitOrgBo, tx sqlbuilder.Tx) (*po.PpmOrgOrganization, errs.SystemErrorInfo) {
	_, err := GetOrgInfoByOutOrgId(initOrgBo.OutOrgId, initOrgBo.SourceChannel)
	if err == nil {
		log.Errorf("组织已经存在，不需要初始化，初始化信息为：%s", json.ToJsonIgnoreError(initOrgBo))
		return nil, errs.BuildSystemErrorInfo(errs.OrgNotNeedInitError)
	}

	orgId, err := idfacade.ApplyPrimaryIdRelaxed(consts.TableOrganization)
	if err != nil {
		return nil, err
	}
	// 企微有点特别。。需要设置一个缓存，因为数据库里面还没有outOrgInfo
	if initOrgBo.SourceChannel == sdk_const.SourceChannelWeixin {
		SetCacheCorpInfo(&sdk_interface.CorpInfo{
			OrgId:         orgId,
			AgentId:       cast.ToInt64(initOrgBo.TenantCode),
			PermanentCode: initOrgBo.PermanentCode,
			CorpId:        initOrgBo.OutOrgId,
		})
	}
	//尝试获取企业信息，如果没有权限之类也不用报错
	client, sdkErr := platform_sdk.GetClient(initOrgBo.SourceChannel, initOrgBo.OutOrgId)
	if sdkErr != nil {
		log.Errorf("[InitOrg] platform_sdk.GetClient err: %v", sdkErr)
		return nil, errs.BuildSystemErrorInfo(errs.PlatFormOpenApiCallError, sdkErr)
	}
	thirdOrgInfo, sdkErr := client.GetOrgInfo(&platformVo.OrgInfoReq{CorpId: initOrgBo.OutOrgId})
	if sdkErr != nil {
		log.Error(sdkErr)
	} else {
		if thirdOrgInfo.Name != "" {
			initOrgBo.OrgName = thirdOrgInfo.Name
		}
		if thirdOrgInfo.DisplayId != "" {
			initOrgBo.TenantCode = thirdOrgInfo.DisplayId
		}
	}

	orgInfo, err := OrgInfoInit(orgId, initOrgBo, tx)
	if err != nil {
		log.Error(err)
		return nil, err
	}

	_, err = OrgOutInfoInit(initOrgBo, orgInfo.Id, tx)
	if err != nil {
		log.Error(err)
		return nil, err
	}

	_, err = OrgConfigInfoInit(orgInfo.Id, tx)
	if err != nil {
		log.Error(err)
		return nil, err
	}

	//权限、角色初始化
	//roleInitResp, err := domain.RoleInit(orgId, tx)
	//if err != nil {
	//	log.Error(err)
	//	return 0, err
	//}
	//log.Info("权限、角色初始化成功")

	//管理组初始化
	initResp := userfacade.InitDefaultManageGroup(orgInfo.Id, 0)
	if initResp.Failure() {
		log.Error(initResp.Message)
		return nil, initResp.Error()
	}

	////优先级，任务类型，任务来源初始化
	//priorityInfo := projectfacade.ProjectInit(projectvo.ProjectInitReqVo{OrgId: orgInfo.Id})
	//if priorityInfo.Failure() {
	//	log.Error(priorityInfo.Message)
	//	return nil, priorityInfo.Error()
	//}
	//log.Info("优先级初始化成功")

	err = InitDepartment(
		client,
		orgInfo.Id,
		initOrgBo.OutOrgId,
		initOrgBo.SourceChannel,
		initOrgBo.OutOrgOwnerId,
		tx)
	if err != nil {
		log.Error(err)
		return nil, err
	}

	return orgInfo, nil
}

func SetFsOrgPayLevel(outOrgId string) errs.SystemErrorInfo {
	effectiveStatus := "normal"
	client, err := platform_sdk.GetClient(sdk_const.SourceChannelFeishu, outOrgId)
	if err != nil {
		log.Error(err)
		return errs.FeiShuClientTenantError
	}
	orderList, sdkError := client.GetOrderList(&platformVo.GetOrderListReq{
		Status:   effectiveStatus,
		Size:     10,
		OutOrgId: outOrgId,
	})
	if sdkError != nil {
		log.Error(sdkError)
		return errs.FeiShuOpenApiCallError
	}

	if len(orderList.List) < 1 {
		log.Errorf("[SetFsOrgPayLevel] 没有飞书订单, outOrgId:%s", outOrgId)
		return nil
	}

	for _, order := range orderList.List {
		orderFsBo := bo.OrderFsBo{
			OrderId:       order.OrderId,
			PricePlanId:   order.PricePlanId,
			PricePlanType: order.PricePlanType,
			Seats:         order.Seats,
			BuyCount:      order.BuyCount,
			Status:        consts.FsOrderStatusNormal,
			BuyType:       order.BuyType,
			SrcOrderId:    order.SrcOrderId,
			DstOrderId:    order.DstOrderId,
			OrderPayPrice: order.OrderPayPrice,
			TenantKey:     order.CorpId,
		}
		if order.PayTime != "" {
			payTime, err := strconv.ParseInt(order.PayTime, 10, 64)
			if err != nil {
				log.Error(err)
			}
			orderFsBo.PaidTime = time.Unix(payTime, 0)
		}
		if order.CreateTime != "" {
			createTime, err := strconv.ParseInt(order.CreateTime, 10, 64)
			if err != nil {
				log.Error(err)
			}
			orderFsBo.CreateTime = time.Unix(createTime, 0)
		}

		fsOrder := orderfacade.AddFsOrder(ordervo.AddFsOrderReq{Data: orderFsBo})
		if fsOrder.Failure() {
			log.Errorf("[SetFsOrgPayLevel] orderfacade.AddFsOrder failed, outOrgId:%s, err:%v", outOrgId, fsOrder.Error())
			return fsOrder.Error()
		}
	}

	return nil
}

func SetWeiXinOrgPayLevel(outOrgId string) errs.SystemErrorInfo {
	client, err := platform_sdk.GetClient(sdk_const.SourceChannelWeixin, outOrgId)
	if err != nil {
		return errs.PlatFormOpenApiCallError
	}
	d, _ := time.ParseDuration("-24h")
	orderReq := &platformVo.GetOrderListReq{
		StartTime: time.Now().Add(d).Unix(),
		EndTime:   time.Now().Add(time.Hour * 24 * 6).Unix(),
		OutOrgId:  outOrgId,
	}
	orderList, sdkError := client.GetOrderList(orderReq)
	if sdkError != nil {
		return errs.PlatFormOpenApiCallError
	}
	if len(orderList.List) < 1 {
		log.Errorf("没有未处理的微信 订单, outOrgId:%s", outOrgId)
		return nil
	}

	log.Infof("[SetWeiXinOrgPayLevel] weixinOrderList:%v", json.ToJsonIgnoreError(orderList.List))

	for _, order := range orderList.List {
		if order.OrderStatus != consts.WeiXinOrderStatusPaySuccess {
			continue
		}
		paidTime, _ := time.Parse(consts.AppTimeFormat, order.PayTime)
		beginTime, _ := time.Parse(consts.AppTimeFormat, order.BeginTime)
		endTime, _ := time.Parse(consts.AppTimeFormat, order.EndTime)
		orderBo := ordervo.AddWeiXinOrderReq{Data: bo.OrderWeiXinBo{
			OutOrgId:      order.CorpId,
			OrderId:       order.OrderId,
			EditionId:     order.EditionId,
			EditionName:   order.EditionName,
			OrderType:     order.OrderType,
			OrderStatus:   order.OrderStatus,
			UserCount:     order.Seats,
			OrderPeriod:   order.OrderPeriod,
			OrderPayPrice: order.OrderPayPrice,
			PaidTime:      paidTime,
			BeginTime:     beginTime,
			EndTime:       endTime,
		}}
		orderReply := orderfacade.AddWeiXinOrder(orderBo)
		if orderReply.Failure() {
			log.Errorf("[SetWeiXinOrgPayLevel]orderfacade.AddWeiXinOrder failed, outOrgId:%s, err:%v", outOrgId, orderReply.Error())
			return orderReply.Error()
		}
	}

	return nil
}

func SetDingOrgPayLevel(outOrgId string) errs.SystemErrorInfo {
	client, err := platform_sdk.GetClient(sdk_const.SourceChannelDingTalk, outOrgId)
	if err != nil {
		return nil
	}
	orderReq := &platformVo.GetOrderListReq{
		Page:     0,
		Size:     100,
		OutOrgId: outOrgId,
	}
	orderList, sdkError := client.GetOrderList(orderReq)
	if sdkError != nil {
		return errs.DingTalkClientError
	}
	if len(orderList.List) < 1 {
		log.Errorf("没有未处理的钉钉 订单, outOrgId:%s", outOrgId)
		return nil
	}

	//orderBos := make([]bo.OrderDingBo, len(orderList.List))

	for _, order := range orderList.List {
		paidTime, _ := time.Parse(consts.AppTimeFormat, order.PayTime)
		orderCreateTime, _ := time.Parse(consts.AppTimeFormat, order.CreateTime)
		//orderBos = append(orderBos, bo.OrderDingBo{
		//	CorpId:          order.CorpId,
		//	OrderId:         order.OrderId,
		//	GoodsName:       "",
		//	GoodsCode:       order.GoodsCode,
		//	ItemName:        "",
		//	ItemCode:        order.ItemCode,
		//	Quantity:        int(order.Quantity),
		//	OrderPayPrice:   order.OrderPayPrice,
		//	PaidTime:        paidTime,
		//	OrderCreateTime: orderCreateTime,
		//	Status:          0,
		//})
		dingOrder := orderfacade.AddDingOrder(ordervo.AddDingOrderReq{Data: bo.OrderDingBo{
			OutOrgId:        order.CorpId,
			OrderId:         order.OrderId,
			GoodsName:       "",
			GoodsCode:       order.GoodsCode,
			ItemName:        "",
			ItemCode:        order.ItemCode,
			Quantity:        int(order.Quantity),
			OrderPayPrice:   order.OrderPayPrice,
			PaidTime:        paidTime,
			OrderCreateTime: orderCreateTime,
			Status:          0,
		}})
		if dingOrder.Failure() {
			log.Errorf("[SetDingOrgPayLevel]orderfacade.AddDingOrder failed, outOrgId:%s, err:%v", outOrgId, dingOrder.Error())
			return dingOrder.Error()
		}
	}

	return nil
}

// 组织信息初始化
func OrgInfoInit(orgId int64, initOrgBo bo.InitOrgBo, tx sqlbuilder.Tx) (*po.PpmOrgOrganization, errs.SystemErrorInfo) {
	isAuth := 0
	if initOrgBo.IsAuthenticated {
		isAuth = 1
	}

	org := &po.PpmOrgOrganization{}
	org.Id = orgId
	org.Status = consts.AppStatusEnable
	org.IsDelete = consts.AppIsNoDelete
	org.SourceChannel = initOrgBo.SourceChannel
	org.Name = initOrgBo.OrgName
	org.LogoUrl = initOrgBo.OrgLogo
	org.Address = initOrgBo.OrgProvince + initOrgBo.OrgCity
	org.IsAuthenticated = isAuth
	err2 := mysql.TransInsert(tx, org)
	if err2 != nil {
		log.Error("组织初始化，添加组织时异常:" + strs.ObjectToString(err2))
		return nil, errs.BuildSystemErrorInfo(errs.MysqlOperateError, err2)
	}
	return org, nil
}

// 组织外部信息初始化
func OrgOutInfoInit(initOrgBo bo.InitOrgBo, orgId int64, tx sqlbuilder.Tx) (int64, errs.SystemErrorInfo) {
	isAuth := 0
	if initOrgBo.IsAuthenticated {
		isAuth = 1
	}

	orgOutInfoId, err := idfacade.ApplyPrimaryIdRelaxed(consts.TableOrganizationOutInfo)
	if err != nil {
		return 0, err
	}
	orgOutInfo := &po.PpmOrgOrganizationOutInfo{}
	orgOutInfo.Id = orgOutInfoId
	orgOutInfo.OrgId = orgId
	orgOutInfo.IsDelete = consts.AppIsNoDelete
	orgOutInfo.Status = consts.AppStatusEnable
	orgOutInfo.SourceChannel = initOrgBo.SourceChannel
	orgOutInfo.Name = initOrgBo.OrgName
	orgOutInfo.OutOrgId = initOrgBo.OutOrgId
	orgOutInfo.Industry = initOrgBo.Industry
	orgOutInfo.IsAuthenticated = isAuth
	orgOutInfo.AuthLevel = strconv.Itoa(initOrgBo.AuthLevel)
	orgOutInfo.TenantCode = initOrgBo.TenantCode
	orgOutInfo.AuthTicket = initOrgBo.PermanentCode

	err2 := mysql.TransInsert(tx, orgOutInfo)
	if err2 != nil {
		log.Error("组织初始化，添加外部组织信息时异常: " + strs.ObjectToString(err2))
		return 0, errs.BuildSystemErrorInfo(errs.MysqlOperateError, err2)
	}
	return orgOutInfoId, nil
}

// 组织配置初始化
func OrgConfigInfoInit(orgId int64, tx sqlbuilder.Tx) (int64, errs.SystemErrorInfo) {
	sysConfig := &po.PpmOrcConfig{}

	//payLevel := &po.PpmBasPayLevel{}
	//err := mysql.SelectById(payLevel.TableName(), 1, payLevel)
	//if err != nil {
	//	log.Error(err)
	//	return 0, errs.BuildSystemErrorInfo(errs.MysqlOperateError, err)
	//}

	orgConfigId, err1 := idfacade.ApplyPrimaryIdRelaxed(consts.TableOrgConfig)
	if err1 != nil {
		log.Error(err1)
		return 0, err1
	}

	sysConfig.Id = orgConfigId
	sysConfig.OrgId = orgId
	sysConfig.TimeZone = "Asia/Shanghai"
	sysConfig.TimeDifference = "+08:00"
	sysConfig.PayLevel = 1
	sysConfig.PayStartTime = time.Now()
	sysConfig.PayEndTime = time.Now().Add(time.Duration(0) * time.Second)
	sysConfig.Language = "zh-CN"
	sysConfig.RemindSendTime = "09:00:00"
	sysConfig.DatetimeFormat = "yyyy-MM-dd HH:mm:ss"
	sysConfig.PasswordLength = 6
	sysConfig.PasswordRule = 1
	sysConfig.MaxLoginFailCount = 0
	sysConfig.Status = consts.AppStatusEnable

	if CheckIsPrivateDeploy() {
		sysConfig.PayLevel = consts.PayLevelPrivateDeploy
		sysConfig.PayEndTime = time.Now().AddDate(100, 0, 0)
	}
	err2 := mysql.TransInsert(tx, sysConfig)
	if err2 != nil {
		log.Error("组织初始化，添加组织配置信息时异常: " + strs.ObjectToString(err2))
		return 0, errs.BuildSystemErrorInfo(errs.MysqlOperateError, err2)
	}

	return orgConfigId, nil
}
