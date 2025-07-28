package ordersvc

import (
	"fmt"
	"time"

	"github.com/star-table/startable-server/app/facade/idfacade"
	"github.com/star-table/startable-server/app/facade/orgfacade"
	"github.com/star-table/startable-server/common/core/util/asyn"

	sdk_const "gitea.bjx.cloud/allstar/platform-sdk/consts"
	"github.com/star-table/startable-server/common/core/consts"
	"github.com/star-table/startable-server/common/core/errs"
	"github.com/star-table/startable-server/common/core/logger"
	"github.com/star-table/startable-server/common/core/util/copyer"
	"github.com/star-table/startable-server/common/core/util/json"
	"github.com/star-table/startable-server/common/core/util/slice"
	"github.com/star-table/startable-server/common/library/db/mysql"
	"github.com/star-table/startable-server/common/model/bo"
	"github.com/star-table/startable-server/common/model/vo/ordervo"
	"github.com/star-table/startable-server/common/model/vo/orgvo"
	"upper.io/db.v3"
	"upper.io/db.v3/lib/sqlbuilder"
)

var log = *logger.GetDefaultLogger()

func AddFsOrder(data bo.OrderFsBo) (int64, errs.SystemErrorInfo) {
	//查询是否已经插入
	isExist, existErr := mysql.IsExistByCond(consts.TableOrderFs, db.Cond{
		consts.TcIsDelete: consts.AppIsNoDelete,
		consts.TcOrderId:  data.OrderId,
	})
	if existErr != nil {
		log.Error(existErr)
		return 0, errs.MysqlOperateError
	}
	if isExist {
		// 如果已经存在，则**更新**。主要针对 orgsvc 下的 SetFsOrgPayLevel 函数中的调用。
		log.Infof("数据已插入系统,%s", json.ToJsonIgnoreError(data))
		return UpdateFsOrder(data)
	}
	orderFsId, err := idfacade.ApplyPrimaryIdRelaxed(consts.TableOrderFs)
	if err != nil {
		log.Error(err)
		return 0, err
	}

	orderId, err := idfacade.ApplyPrimaryIdRelaxed(consts.TableOrder)
	if err != nil {
		log.Error(err)
		return 0, err
	}
	//获取组织id
	orgInfoResp := orgfacade.GetBaseOrgInfoByOutOrgId(orgvo.GetBaseOrgInfoByOutOrgIdReqVo{OutOrgId: data.TenantKey})
	orgId := int64(0)
	if orgInfoResp.Successful() {
		orgId = orgInfoResp.BaseOrgInfo.OrgId
		asyn.Execute(func() {
			resp := orgfacade.ClearOrgUsersPayCache(orgvo.GetBaseOrgInfoReqVo{
				OrgId: orgId,
			})
			if resp.Failure() {
				log.Errorf("[ClearOrgUsersPayCache] err:%v", resp)
			}
		})
	} else {
		log.Infof("[AddFsOrder] 新用户新增订单时，还没有创建组织，因此 orgId 为 0，属正常情况。info: %s", json.ToJsonIgnoreError(orgInfoResp))
	}

	//拼装数据
	data.Id = orderFsId
	data.OrgId = orgId
	orderFsPo := &po.PpmOrdOrderFs{}
	copyer.Copy(data, orderFsPo)
	transErr := mysql.TransX(func(tx sqlbuilder.Tx) error {
		//飞书订单表
		err := mysql.TransInsert(tx, orderFsPo)
		if err != nil {
			log.Error(err)
			return err
		}

		//订单表
		orderPo := &po.PpmOrdOrder{
			Id:              orderId,
			OrgId:           data.OrgId,
			OutOrderNo:      orderFsId,
			Status:          0,
			OrderCreateTime: orderFsPo.OrderCreateTime,
			PaidTime:        orderFsPo.PaidTime,
			EffectiveTime:   0,
			Seats:           orderFsPo.Seats,
			TotalPrice:      orderFsPo.OrderPayPrice,
			OrderPayPrice:   orderFsPo.OrderPayPrice,
			SourceChannel:   sdk_const.SourceChannelFeishu,
			BuyType:         orderFsPo.BuyType,
			BuyCount:        orderFsPo.BuyCount,
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
		return 0, errs.MysqlOperateError
	}

	// 再更新 org pay config
	if err := UpdateFsOrgPayConfig(orgId, orderFsPo); err != nil {
		log.Errorf("[AddFsOrder] orgId: %d, UpdateOrgPayConfig err: %v", orgId, err)
		return 0, err
	}

	return orderId, nil
}

// UpdateFsOrder 飞书订单的更新，主要更新订单的 org_id
func UpdateFsOrder(data bo.OrderFsBo) (int64, errs.SystemErrorInfo) {
	// 根据订单号，查询订单
	targetFsOrder := po.PpmOrdOrderFs{}
	oriErr := mysql.SelectOneByCond(consts.TableOrderFs, db.Cond{
		consts.TcIsDelete: consts.AppIsNoDelete,
		consts.TcOrderId:  data.OrderId,
	}, &targetFsOrder)
	if oriErr != nil {
		log.Errorf("[UpdateFsOrder] err: %v", oriErr)
		return 0, errs.BuildSystemErrorInfo(errs.MysqlOperateError, oriErr)
	}
	// 如果 org_id 不为 0，则无需更新
	if targetFsOrder.OrgId > 0 {
		log.Infof("[UpdateFsOrder] 订单信息是最新的，无需更新。")
		return 0, nil
	}
	// 查询组织
	orgId := int64(0)
	orgInfoResp := orgfacade.GetBaseOrgInfoByOutOrgId(orgvo.GetBaseOrgInfoByOutOrgIdReqVo{OutOrgId: data.TenantKey})
	if orgInfoResp.Successful() {
		orgId = orgInfoResp.BaseOrgInfo.OrgId
		asyn.Execute(func() {
			resp := orgfacade.ClearOrgUsersPayCache(orgvo.GetBaseOrgInfoReqVo{
				OrgId: orgId,
			})
			if resp.Failure() {
				log.Errorf("[ClearOrgUsersPayCache] err:%v", resp)
			}
		})
	} else {
		// 如果没查到，这打印一下日志。
		log.Infof("[UpdateFsOrder] 新用户新增订单时，还没有创建组织，因此 orgId 为 0，属正常情况。info: %s", json.ToJsonIgnoreError(orgInfoResp))
		return 0, nil
	}
	transErr := mysql.TransX(func(tx sqlbuilder.Tx) error {
		oriErr = mysql.TransUpdateSmart(tx, consts.TableOrderFs, targetFsOrder.Id, mysql.Upd{
			consts.TcOrgId: orgId,
		})
		if oriErr != nil {
			log.Errorf("[UpdateFsOrder] update TableOrderFs err: %v", oriErr)
			return oriErr
		}

		_, oriErr := mysql.TransUpdateSmartWithCond(tx, consts.TableOrder, db.Cond{
			consts.TcOutOrderNo: targetFsOrder.Id,
		}, mysql.Upd{
			consts.TcOrgId: orgId,
		})
		if oriErr != nil {
			log.Errorf("[UpdateFsOrder] update TableOrder err: %v", oriErr)
			return oriErr
		}
		return nil
	})
	if transErr != nil {
		log.Errorf("[UpdateFsOrder] trans err: %v", transErr)
		return 0, errs.BuildSystemErrorInfo(errs.MysqlOperateError, transErr)
	}
	targetFsOrder.OrgId = orgId

	// 再更新 org pay config
	if err := UpdateFsOrgPayConfig(orgId, &targetFsOrder); err != nil {
		log.Errorf("[UpdateFsOrder] UpdateFsOrgPayConfig err: %v", err)
		return 0, err
	}

	return targetFsOrder.Id, nil
}

// UpdateFsOrgPayConfig 更细用户的组织配置，使之使用对应的版本（标准版、企业版）
func UpdateFsOrgPayConfig(orgId int64, fsOrder *po.PpmOrdOrderFs) errs.SystemErrorInfo {
	//获取订单方案对应的等级
	level := consts.GetFsOrderLevel(fsOrder.PricePlanId)
	//判断有效期（目前只有处理: 按年,按月,试用） 2021-5-13支持定向方案的按截止日期、按有效天数
	//订单类型有3种（正常buy，续费renew，增购upgrade）
	if ok, _ := slice.Contain([]string{consts.FsPerSeatPerYear, consts.FsPerSeatPerMonth, consts.FsTrial, consts.FsActiveDay, consts.FsActiveEndDate},
		fsOrder.PricePlanType); ok {
		//设置组织等级以及设置功能项
		resp := orgfacade.UpdateOrgFunctionConfig(orgvo.UpdateOrgFunctionConfigReq{
			OrgId:         orgId,
			SourceChannel: sdk_const.SourceChannelFeishu,
			Input: orgvo.UpdateFunctionConfigData{
				Level:         level,
				BuyType:       fsOrder.BuyType,
				PricePlanType: fsOrder.PricePlanType,
				PayTime:       fsOrder.PaidTime,
				//Seats:         0,
				//ExpireDays:    0,
				//EndDate:       time.Time{},
				//TrailDays:     0,
			},
		})
		if resp.Failure() {
			log.Errorf("[UpdateOrgPayConfig] orgId: %d, outOrderId: %v, err: %v", orgId, fsOrder.OrderId, resp.Error())
			return resp.Error()
		}
	}

	return nil
}

func AddDingOrder(data bo.OrderDingBo) (int64, errs.SystemErrorInfo) {
	//查询是否已经插入
	isExist, existErr := mysql.IsExistByCond(consts.TableOrderDing, db.Cond{
		consts.TcIsDelete: consts.AppIsNoDelete,
		consts.TcOrderId:  data.OrderId,
	})
	if existErr != nil {
		log.Errorf("[AddDingOrder]select err:%v, orderId:%v", existErr, data.OrderId)
		return 0, errs.MysqlOperateError
	}
	if isExist {
		log.Infof("[AddDingOrder] 数据已插入,%s", json.ToJsonIgnoreError(data))
		return updateDingOrder(data)
	}

	orderDingPo, errSys := domain.CreateDingOrder(data)
	if errSys != nil {
		log.Errorf("[AddDingOrder] domain.CreateDingOrder err:%v, orderDingPo:%v", errSys, orderDingPo)
		return 0, errSys
	}

	orgId := orderDingPo.OrgId

	err := updateDingOrgPayConfig(orgId, orderDingPo)
	if err != nil {
		log.Errorf("[AddDingOrder]updateDingOrgPayConfig err:%v, orgId:%v, orderId:%v", err, orgId, orderDingPo.OrderId)
		return 0, err
	}
	return orderDingPo.Id, nil
}

func updateDingOrder(data bo.OrderDingBo) (int64, errs.SystemErrorInfo) {
	dingOrder := po.PpmOrdOrderDing{}
	err := mysql.SelectOneByCond(consts.TableOrderDing, db.Cond{
		consts.TcIsDelete: consts.AppIsNoDelete,
		consts.TcOrderId:  data.OrderId,
	}, &dingOrder)
	if err != nil {
		log.Errorf("[updateDingOrder] err:%v", err)
		return 0, errs.BuildSystemErrorInfo(errs.MysqlOperateError, err)
	}
	if dingOrder.OrgId > 0 {
		log.Infof("[updateDingOrder] 订单信息是最新的，无需更新。")
		return 0, nil
	}
	// 查询组织
	orgId := int64(0)
	orgInfoResp := orgfacade.GetBaseOrgInfoByOutOrgId(orgvo.GetBaseOrgInfoByOutOrgIdReqVo{OutOrgId: data.OutOrgId})
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
		// 如果没查到，这打印一下日志。
		log.Infof("[updateDingOrder] 新用户新增订单时，还没有创建组织，因此 orgId 为 0，属正常情况。info: %s", json.ToJsonIgnoreError(orgInfoResp))
		return 0, nil
	}

	transErr := mysql.TransX(func(tx sqlbuilder.Tx) error {
		err := mysql.TransUpdateSmart(tx, consts.TableOrderDing, dingOrder.Id, mysql.Upd{
			consts.TcOrgId: orgId,
		})
		if err != nil {
			log.Errorf("[updateDingOrder] update TableOrderDing err: %v", err)
			return err
		}
		_, err = mysql.TransUpdateSmartWithCond(tx, consts.TableOrder, db.Cond{
			consts.TcOutOrderNo: dingOrder.Id,
		}, mysql.Upd{
			consts.TcOrgId: orgId,
		})
		if err != nil {
			log.Errorf("[updateDingOrder] update TableOrder err: %v", err)
			return err
		}
		return nil
	})

	if transErr != nil {
		log.Errorf("[updateDingOrder] trans err: %v", transErr)
		return 0, errs.BuildSystemErrorInfo(errs.MysqlOperateError, transErr)
	}

	// 更新 org pay config
	errSys := updateDingOrgPayConfig(orgId, &dingOrder)
	if errSys != nil {
		log.Errorf("[updateDingOrder] updateDingOrgPayConfig err: %v", err)
		return 0, errSys
	}

	return dingOrder.Id, nil
}

// 更新钉钉用户的组织配置，标准版、企业版、旗舰版
func updateDingOrgPayConfig(orgId int64, dingOrder *po.PpmOrdOrderDing) errs.SystemErrorInfo {
	log.Infof("[updateDingOrgPayConfig] orgId:%v, dingOrder:%v", orgId, json.ToJsonIgnoreError(dingOrder))
	//level, err := domain.GetDingOrderLevel(dingOrder.GoodsCode, dingOrder.ItemCode)
	//if err != nil {
	//	log.Error(err)
	//	return err
	//}
	level := consts.GetDingOrderLevel(dingOrder.ItemCode)

	// 查询该组织的订单是不是已经有免费版本的订单了，如果有 就跳过不更新payLevel
	if domain.CheckDingFreeOrder(orgId) {
		return nil
	}

	//设置组织等级以及设置功能项
	resp := orgfacade.UpdateOrgFunctionConfig(orgvo.UpdateOrgFunctionConfigReq{
		OrgId:         orgId,
		SourceChannel: sdk_const.SourceChannelDingTalk,
		Input: orgvo.UpdateFunctionConfigData{
			Level:         level,
			PayTime:       dingOrder.PaidTime,
			EndDate:       dingOrder.EndTime,
			BuyType:       dingOrder.OrderType,
			PricePlanType: dingOrder.OrderChargeType,
		},
	})
	if resp.Failure() {
		log.Errorf("[updateDingOrgPayConfig] orgId: %d, outOrderId: %v, err: %v", orgId, dingOrder.OrderId, resp.Error())
		return resp.Error()
	}

	return nil
}

// CheckFsOrgIsInTrial 检查一个 fs 组织目前是否是处于试用期间
func CheckFsOrgIsInTrial(orgId int64) (bool, errs.SystemErrorInfo) {
	fsOrders, err := domain.GetOrderListByCond(nil, db.Cond{
		consts.TcOrgId:         orgId,
		consts.TcPricePlanType: consts.FsTrial,
	}, db.Raw("id desc"), 1, 1)
	if err != nil {
		log.Errorf("[CheckFsOrgIsInTrial] err: %v, orgId: %d", err, orgId)
		return false, err
	}
	if len(fsOrders) < 1 {
		return false, nil
	}
	// 试用订单创建时间 + 15days 后意味着试用过期
	fsTrialOrder := fsOrders[0]
	endDate := fsTrialOrder.CreateTime.AddDate(0, 0, consts.FsOrderTrailDays).Format(consts.AppDateFormat)
	nowDate := time.Now().Format(consts.AppDateFormat)
	if nowDate <= endDate {
		return true, nil
	}

	return false, nil
}

func UpdateOrderInfo(orgId int64, outOrgId string) (*ordervo.UpdateDingOrderRespData, errs.SystemErrorInfo) {
	return domain.UpdateOrderInfo(orgId, outOrgId)
}

func DeleteDingOrderByOrderId(outOrgId string, orderId string) (int64, errs.SystemErrorInfo) {
	orgInfoResp := orgfacade.GetBaseOrgInfoByOutOrgId(orgvo.GetBaseOrgInfoByOutOrgIdReqVo{OutOrgId: outOrgId})
	if orgInfoResp.Failure() {
		log.Errorf("[DeleteDingOrderByOrderId]GetBaseOrgInfoByOutOrgId err:%v", orgInfoResp.Error())
		return 0, orgInfoResp.Error()
	}
	orgId := orgInfoResp.BaseOrgInfo.OrgId
	// 删除订单
	// 服务版本降级
	return domain.DeleteDingOrder(orgId, orderId)
}

func GetOrderPayInfo(orgId int64) (*ordervo.GetOrderPayInfo, errs.SystemErrorInfo) {
	return domain.GetOrderPayInfo(orgId)
}

func AddWeiXinOrder(data bo.OrderWeiXinBo) (int64, errs.SystemErrorInfo) {
	isExist, existErr := mysql.IsExistByCond(consts.TableOrderWeiXin, db.Cond{
		consts.TcIsDelete: consts.AppIsNoDelete,
		consts.TcOrderId:  data.OrderId,
	})
	if existErr != nil {
		log.Errorf("[AddWeiXinOrder]select err:%v, orderId:%v", existErr, data.OrderId)
		return 0, errs.MysqlOperateError
	}
	if isExist {
		log.Infof("[AddWeiXinOrder] 数据已插入,%s", json.ToJsonIgnoreError(data))
		// update
		return updateWeiXinOrder(data)
	}

	order, errSys := domain.CreateWeiXinOrder(data)
	if errSys != nil {
		log.Errorf("[AddWeiXinOrder] CreateWeiXinOrder err:%v", errSys)
		return 0, nil
	}
	orgId := order.OrgId
	// 再更新 org pay config
	errSys = updateWeiXinOrgPayConfig(orgId, order)
	if errSys != nil {
		log.Errorf("[AddWeiXinOrder]updateWeiXinOrgPayConfig err:%v, orgId:%v, orderId:%v", errSys, orgId, order.OrderId)
		return 0, errSys
	}

	return order.Id, nil
}

func updateWeiXinOrgPayConfig(orgId int64, weixinOrder *po.PpmOrdOrderWeiXin) errs.SystemErrorInfo {
	log.Infof("[updateWeiXinOrgPayConfig] orgId:%v, weixinOrder:%v", orgId, json.ToJsonIgnoreError(weixinOrder))
	level := consts.GetWeiXinOrderLevel(weixinOrder.EditionId)
	//设置组织等级以及设置功能项
	resp := orgfacade.UpdateOrgFunctionConfig(orgvo.UpdateOrgFunctionConfigReq{
		OrgId:         orgId,
		SourceChannel: sdk_const.SourceChannelWeixin,
		Input: orgvo.UpdateFunctionConfigData{
			Level:   level,
			PayTime: weixinOrder.PaidTime,
			EndDate: weixinOrder.EndTime,
			BuyType: fmt.Sprintf("%d", weixinOrder.OrderType),
			//PricePlanType: fmt.Sprintf("%d", weixinOrder.OrderType),
		},
	})
	if resp.Failure() {
		log.Errorf("[updateWeiXinOrgPayConfig] orgId: %d, outOrderId: %v, err: %v", orgId, weixinOrder.OrderId, resp.Error())
		return resp.Error()
	}
	return nil
}

func updateWeiXinOrder(data bo.OrderWeiXinBo) (int64, errs.SystemErrorInfo) {
	weixinOrder := po.PpmOrdOrderWeiXin{}
	err := mysql.SelectOneByCond(consts.TableOrderWeiXin, db.Cond{
		consts.TcIsDelete: consts.AppIsNoDelete,
		consts.TcOrderId:  data.OrderId,
	}, &weixinOrder)
	if err != nil {
		log.Errorf("[updateWeiXinOrder] err:%v", err)
		return 0, errs.BuildSystemErrorInfo(errs.MysqlOperateError, err)
	}
	if weixinOrder.OrgId > 0 {
		log.Infof("[updateWeiXinOrder] 订单信息是最新的，无需更新。")
		return 0, nil
	}
	// 查询组织
	orgId := int64(0)
	orgInfoResp := orgfacade.GetBaseOrgInfoByOutOrgId(orgvo.GetBaseOrgInfoByOutOrgIdReqVo{OutOrgId: data.OutOrgId})
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
		// 如果没查到，这打印一下日志。
		log.Infof("[updateWeiXinOrder] 新用户新增订单时，还没有创建组织，因此 orgId 为 0，属正常情况。info: %s", json.ToJsonIgnoreError(orgInfoResp))
		return 0, nil
	}

	transErr := mysql.TransX(func(tx sqlbuilder.Tx) error {
		err := mysql.TransUpdateSmart(tx, consts.TableOrderWeiXin, weixinOrder.Id, mysql.Upd{
			consts.TcOrgId: orgId,
		})
		if err != nil {
			log.Errorf("[updateWeiXinOrder] update TableOrderWeiXin err: %v", err)
			return err
		}
		_, err = mysql.TransUpdateSmartWithCond(tx, consts.TableOrder, db.Cond{
			consts.TcOutOrderNo: weixinOrder.Id,
		}, mysql.Upd{
			consts.TcOrgId: orgId,
		})
		if err != nil {
			log.Errorf("[updateWeiXinOrder] update TableOrder err: %v", err)
			return err
		}
		return nil
	})

	if transErr != nil {
		log.Errorf("[updateWeiXinOrder] trans err: %v", transErr)
		return 0, errs.BuildSystemErrorInfo(errs.MysqlOperateError, transErr)
	}

	// 更新 org pay config
	errSys := updateWeiXinOrgPayConfig(orgId, &weixinOrder)
	if errSys != nil {
		log.Errorf("[updateWeiXinOrder] updateWeiXinOrgPayConfig err: %v", err)
		return 0, errSys
	}
	return weixinOrder.Id, nil
}
