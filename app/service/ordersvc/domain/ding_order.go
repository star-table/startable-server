package ordersvc

import (
	sdk_const "gitea.bjx.cloud/allstar/platform-sdk/consts"
	"github.com/star-table/startable-server/app/facade/idfacade"
	"github.com/star-table/startable-server/app/facade/orgfacade"
	"github.com/star-table/startable-server/common/core/consts"
	"github.com/star-table/startable-server/common/core/errs"
	"github.com/star-table/startable-server/common/core/util/asyn"
	"github.com/star-table/startable-server/common/core/util/copyer"
	"github.com/star-table/startable-server/common/core/util/json"
	"github.com/star-table/startable-server/common/library/db/mysql"
	"github.com/star-table/startable-server/common/model/bo"
	"github.com/star-table/startable-server/common/model/vo/ordervo"
	"github.com/star-table/startable-server/common/model/vo/orgvo"
	"upper.io/db.v3"
	"upper.io/db.v3/lib/sqlbuilder"
)

func CreateDingOrder(data bo.OrderDingBo) (*po.PpmOrdOrderDing, errs.SystemErrorInfo) {
	orderDingId, err := idfacade.ApplyPrimaryIdRelaxed(consts.TableOrderDing)
	if err != nil {
		log.Errorf("[CreateDingOrder]TableOrderDing id gen error:%v, orderId:%v", err, data.OrderId)
		return nil, err
	}

	orderId, err := idfacade.ApplyPrimaryIdRelaxed(consts.TableOrder)
	if err != nil {
		log.Errorf("[CreateDingOrder]TableOrder id gen error:%v, orderId:%v", err, data.OrderId)
		return nil, err
	}

	//获取组织id
	outOrgId := data.OutOrgId
	orgInfoResp := orgfacade.GetBaseOrgInfoByOutOrgId(orgvo.GetBaseOrgInfoByOutOrgIdReqVo{OutOrgId: outOrgId})

	orgId := int64(0)
	if orgInfoResp.Successful() {
		if orgInfoResp.BaseOrgInfo != nil {
			orgId = orgInfoResp.BaseOrgInfo.OrgId
		}
		asyn.Execute(func() {
			resp := orgfacade.ClearOrgUsersPayCache(orgvo.GetBaseOrgInfoReqVo{
				OrgId: orgId,
			})
			if resp.Failure() {
				log.Errorf("[ClearOrgUsersPayCache] err:%v", resp)
			}
		})
	} else {
		log.Infof("[AddDingOrder] 新用户新增订单时，还没有创建组织，因此 orgId 为 0，属正常情况。info: %s", json.ToJsonIgnoreError(orgInfoResp))
	}

	//拼装数据
	data.Id = orderDingId
	data.OrgId = orgId
	orderDingPo := &po.PpmOrdOrderDing{}
	copyer.Copy(data, orderDingPo)
	transErr := mysql.TransX(func(tx sqlbuilder.Tx) error {
		//钉钉订单表
		err := mysql.TransInsert(tx, orderDingPo)
		if err != nil {
			log.Error(err)
			return err
		}

		//订单表
		orderPo := &po.PpmOrdOrder{
			Id:              orderId,
			OrgId:           orgId,
			OutOrderNo:      orderDingId,
			Status:          0,
			OrderCreateTime: orderDingPo.OrderCreateTime,
			PaidTime:        orderDingPo.PaidTime,
			EffectiveTime:   0,
			//Seats:           orderDingPo.Seats,
			TotalPrice:    orderDingPo.OrderPayPrice,
			OrderPayPrice: orderDingPo.OrderPayPrice,
			SourceChannel: sdk_const.SourceChannelDingTalk,
			BuyType:       orderDingPo.OrderType,
			BuyCount:      orderDingPo.Quantity,
		}

		err1 := mysql.TransInsert(tx, orderPo)
		if err1 != nil {
			log.Error(err1)
			return err1
		}

		return nil
	})
	if transErr != nil {
		log.Error(transErr)
		return nil, errs.MysqlOperateError
	}

	return orderDingPo, nil
}

func UpdateOrderInfo(orgId int64, outOrgId string) (*ordervo.UpdateDingOrderRespData, errs.SystemErrorInfo) {
	log.Infof("[UpdateOrderInfo] orgId:%v, outOrgId:%v", orgId, outOrgId)
	orderDing := &po.PpmOrdOrderDing{}
	cond := db.Cond{
		consts.TcIsDelete: consts.AppIsNoDelete,
		consts.TcOutOrgId: outOrgId,
		consts.TcOrgId:    0,
	}
	errSys := mysql.SelectOneByCond(consts.TableOrderDing, cond, orderDing)
	if errSys != nil {
		log.Errorf("[UpdateOrderInfo] err:%v", errSys)
		if errSys == db.ErrNoMoreRows {
			return nil, nil
		}
		return nil, errs.MysqlOperateError
	}

	//level, err := GetDingOrderLevel(orderDing.GoodsCode, orderDing.ItemCode)
	//if err != nil {
	//	log.Errorf("[UpdateOrderInfo] GetDingOrderLevel err:%v", err)
	//	return nil, err
	//}
	level := consts.GetDingOrderLevel(orderDing.ItemCode)

	tranErr := mysql.TransX(func(tx sqlbuilder.Tx) error {
		upd := mysql.Upd{consts.TcOrgId: orgId}
		_, err := mysql.TransUpdateSmartWithCond(tx, consts.TableOrderDing, db.Cond{
			consts.TcIsDelete: consts.AppIsNoDelete,
			consts.TcOutOrgId: outOrgId,
		}, upd)
		if err != nil {
			log.Errorf("[UpdateOrderInfo] err:%v", err)
			return err
		}

		_, err = mysql.TransUpdateSmartWithCond(tx, consts.TableOrder, db.Cond{
			consts.TcOutOrderNo:    orderDing.Id,
			consts.TcSourceChannel: sdk_const.SourceChannelDingTalk,
		}, upd)
		if err != nil {
			log.Errorf("[UpdateOrderInfo] err:%v", err)
			return err
		}
		return nil
	})

	if tranErr != nil {
		log.Errorf("[UpdateOrderInfo] err:%v, orgId:%v, outOrgId:%v", tranErr, orgId, outOrgId)
		return nil, errs.BuildSystemErrorInfo(errs.MysqlOperateError, tranErr)
	}

	orderDingBo := bo.OrderDingBo{}
	errCopy := copyer.Copy(orderDing, &orderDingBo)
	if errCopy != nil {
		log.Errorf("[UpdateOrderInfo] err:%v", errCopy)
		return nil, errs.ObjectCopyError
	}

	return &ordervo.UpdateDingOrderRespData{
		Level:    level,
		DingData: orderDingBo,
	}, nil
}

func DeleteDingOrder(orgId int64, orderId string) (int64, errs.SystemErrorInfo) {
	dingOrderPo := po.PpmOrdOrderDing{}
	err := mysql.SelectOneByCond(consts.TableOrderDing, db.Cond{
		consts.TcIsDelete: consts.AppIsNoDelete,
		consts.TcOrgId:    orgId,
		consts.TcOrderId:  orderId,
	}, &dingOrderPo)
	if err != nil {
		log.Errorf("[DeleteDingOrder] err:%v, orgId:%v, orderId:%v", orgId, orderId)
		return 0, errs.MysqlOperateError
	}

	errTrans := mysql.TransX(func(tx sqlbuilder.Tx) error {
		_, err = mysql.TransUpdateSmartWithCond(tx, consts.TableOrderDing, db.Cond{
			consts.TcIsDelete: consts.AppIsNoDelete,
			consts.TcOrgId:    orgId,
			consts.TcOrderId:  orderId,
		}, mysql.Upd{
			consts.TcIsDelete: consts.AppIsDeleted,
		})
		if err != nil {
			log.Errorf("[DeleteDingOrder] err:%v, orgId:%v, orderId:%v", orgId, orderId)
			return err
		}
		_, err = mysql.TransUpdateSmartWithCond(tx, consts.TableOrder, db.Cond{
			consts.TcOutOrderNo: dingOrderPo.Id,
		}, mysql.Upd{
			consts.TcIsDelete: consts.AppIsDeleted,
		})
		if err != nil {
			log.Errorf("[DeleteDingOrder] err:%v, orgId:%v, orderId:%v", orgId, orderId)
			return err
		}

		return nil
	})
	if errTrans != nil {
		log.Errorf("[DeleteDingOrder] err:%v, orgId:%v, orderId:%v", orgId, orderId)
		return 0, errs.MysqlOperateError
	}

	orgConfigResp := orgfacade.UpdateOrgFunctionConfig(orgvo.UpdateOrgFunctionConfigReq{
		OrgId:         orgId,
		UserId:        0,
		SourceChannel: sdk_const.SourceChannelDingTalk,
		Input: orgvo.UpdateFunctionConfigData{
			Level:         consts.PayLevelStandard,
			BuyType:       dingOrderPo.OrderType,
			PricePlanType: dingOrderPo.OrderChargeType,
			PayTime:       dingOrderPo.PaidTime,
			EndDate:       dingOrderPo.EndTime,
		},
	})
	if orgConfigResp.Failure() {
		log.Errorf("[DeleteDingOrder] UpdateOrgFunctionConfig err:%v, orgId:%v, orderId:%v", orgId, orderId)
		return 0, orgConfigResp.Error()
	}

	return orgId, nil
}

func CreateWeiXinOrder(data bo.OrderWeiXinBo) (*po.PpmOrdOrderWeiXin, errs.SystemErrorInfo) {
	orderWeixinId, err := idfacade.ApplyPrimaryIdRelaxed(consts.TableOrderWeiXin)
	if err != nil {
		log.Errorf("[CreateWeiXinOrder]TableOrderWeiXin id gen error:%v, orderId:%v", err, data.OrderId)
		return nil, err
	}

	orderId, err := idfacade.ApplyPrimaryIdRelaxed(consts.TableOrder)
	if err != nil {
		log.Errorf("[CreateDingOrder]TableOrder id gen error:%v, orderId:%v", err, data.OrderId)
		return nil, err
	}

	//获取组织id
	outOrgId := data.OutOrgId
	orgInfoResp := orgfacade.GetBaseOrgInfoByOutOrgId(orgvo.GetBaseOrgInfoByOutOrgIdReqVo{OutOrgId: outOrgId})

	orgId := int64(0)
	if orgInfoResp.Successful() {
		if orgInfoResp.BaseOrgInfo != nil {
			orgId = orgInfoResp.BaseOrgInfo.OrgId
		}
		asyn.Execute(func() {
			resp := orgfacade.ClearOrgUsersPayCache(orgvo.GetBaseOrgInfoReqVo{
				OrgId: orgId,
			})
			if resp.Failure() {
				log.Errorf("[ClearOrgUsersPayCache] err:%v", resp)
			}
		})
	} else {
		log.Infof("[CreateWeiXinOrder] 新用户新增订单时，还没有创建组织，因此 orgId 为 0，属正常情况。info: %s", json.ToJsonIgnoreError(orgInfoResp))
	}
	//拼装数据
	data.Id = orderWeixinId
	data.OrgId = orgId
	orderWeixinPo := &po.PpmOrdOrderWeiXin{}
	copyer.Copy(data, orderWeixinPo)
	transErr := mysql.TransX(func(tx sqlbuilder.Tx) error {
		// 微信订单表
		err := mysql.TransInsert(tx, orderWeixinPo)
		if err != nil {
			log.Error(err)
			return err
		}
		//订单表
		orderPo := &po.PpmOrdOrder{
			Id:            orderId,
			OrgId:         orgId,
			OutOrderNo:    orderWeixinId,
			PaidTime:      orderWeixinPo.PaidTime,
			TotalPrice:    orderWeixinPo.OrderPayPrice,
			OrderPayPrice: orderWeixinPo.OrderPayPrice,
			SourceChannel: sdk_const.SourceChannelWeixin,
			Seats:         orderWeixinPo.UserCount,
		}

		err1 := mysql.TransInsert(tx, orderPo)
		if err1 != nil {
			log.Error(err1)
			return err1
		}

		return nil
	})
	if transErr != nil {
		log.Error(transErr)
		return nil, errs.MysqlOperateError
	}

	return orderWeixinPo, nil
}
