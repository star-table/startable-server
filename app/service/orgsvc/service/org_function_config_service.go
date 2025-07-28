package orgsvc

import (
	"time"

	"github.com/google/martian/log"
	"github.com/spf13/cast"

	msgPb "gitea.bjx.cloud/LessCode/interface/golang/msg/v1"
	"github.com/star-table/startable-server/common/core/threadlocal"
	"github.com/star-table/startable-server/common/model/vo/commonvo"

	"github.com/star-table/startable-server/common/core/util/asyn"

	"github.com/star-table/startable-server/common/core/util/copyer"
	"github.com/star-table/startable-server/common/model/vo/orgvo"

	sdk_const "gitea.bjx.cloud/allstar/platform-sdk/consts"
	"github.com/star-table/startable-server/common/core/consts"
	"github.com/star-table/startable-server/common/core/errs"
	"github.com/star-table/startable-server/common/core/util/date"
	"github.com/star-table/startable-server/common/library/db/mysql"
	"github.com/star-table/startable-server/common/model/vo"
	"upper.io/db.v3"
)

func GetFunctionKeysByOrg(orgId int64) (*vo.FunctionConfigResp, errs.SystemErrorInfo) {
	bos, err := domain.GetOrgPayFunction(orgId)
	if err != nil {
		log.Error(err)
		return nil, err
	}

	funcKeys := domain.GetFunctionKeyListByFunctions(bos)
	result := &vo.FunctionConfigResp{FunctionCodes: funcKeys}

	return result, nil
}

func GetFunctionsByOrg(orgId int64) ([]orgvo.FunctionLimitObj, errs.SystemErrorInfo) {
	bos, err := domain.GetOrgPayFunction(orgId)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	newFuncArr := make([]orgvo.FunctionLimitObj, 0, len(bos))
	copyer.Copy(bos, &newFuncArr)

	return newFuncArr, nil
}

func UpdateOrgFunctionConfig(orgId int64, sourceChannel string, level int64, buyType string, pricePlanType string, payTime time.Time, seats int, expireDays int, endDate time.Time, trailDays int) errs.SystemErrorInfo {
	log.Infof("[UpdateOrgFunctionConfig] sourceChannel:%v, orgId:%v, buyType:%v, pricePlanType:%v, payTime:%v, endDate:%v, payLeve:%v",
		sourceChannel, orgId, buyType, pricePlanType, payTime, endDate, level)
	//查看当前等级是否存在

	errSet := domain.ResetOrgPayNum(orgId)
	if errSet != nil {
		log.Errorf("[ResetOrgPayNum] orgId:%v, err:%v", orgId, errSet)
	}

	levelIsExist, existErr := domain.PayLevelIsExist(level)
	if existErr != nil {
		log.Error(existErr)
		return existErr
	}
	if !levelIsExist {
		return errs.PayLevelNotExist
	}
	//获取当前组织等级
	orgConfig, err := domain.GetOrgConfig(orgId)
	if err != nil {
		log.Error(err)
		return err
	}
	upd := mysql.Upd{}
	if int64(orgConfig.PayLevel) != level {
		upd[consts.TcPayLevel] = level
	}
	if seats != 0 {
		upd[consts.TcSeats] = seats //目前只用于定向方案
	}
	//支付当天的零点
	payStartTime := date.GetZeroTime(payTime)
	if sourceChannel == sdk_const.SourceChannelFeishu {
		//有效期设置(增购的话不需要修改有效期，人数限制从新的获取就好了)
		if pricePlanType == consts.FsPerSeatPerMonth {
			if buyType == consts.FsOrderBuy {
				//普通支付
				if orgConfig.PayEndTime.After(payStartTime) {
					upd[consts.TcPayEndTime] = orgConfig.PayEndTime.AddDate(0, 1, 0).Add(-1 * time.Second)
				} else {
					upd[consts.TcPayStartTime] = payStartTime
					upd[consts.TcPayEndTime] = payStartTime.AddDate(0, 1, 0).Add(-1 * time.Second)
				}
			} else if buyType == consts.FsOrderRenew {
				if orgConfig.PayEndTime.After(payStartTime) {
					upd[consts.TcPayEndTime] = orgConfig.PayEndTime.AddDate(0, 1, 0).Add(-1 * time.Second)
				} else {
					upd[consts.TcPayStartTime] = payStartTime
					upd[consts.TcPayEndTime] = payStartTime.AddDate(0, 1, 0).Add(-1 * time.Second)
				}
			}
		} else if pricePlanType == consts.FsPerSeatPerYear {
			if buyType == consts.FsOrderBuy {
				//普通支付
				if orgConfig.PayEndTime.After(payStartTime) {
					upd[consts.TcPayEndTime] = orgConfig.PayEndTime.AddDate(1, 0, 0).Add(-1 * time.Second)
				} else {
					upd[consts.TcPayStartTime] = payStartTime
					upd[consts.TcPayEndTime] = payStartTime.AddDate(1, 0, 0).Add(-1 * time.Second)
				}
			} else if buyType == consts.FsOrderRenew {
				if orgConfig.PayEndTime.After(payStartTime) {
					upd[consts.TcPayEndTime] = orgConfig.PayEndTime.AddDate(1, 0, 0).Add(-1 * time.Second)
				} else {
					upd[consts.TcPayStartTime] = payStartTime
					upd[consts.TcPayEndTime] = payStartTime.AddDate(1, 0, 0).Add(-1 * time.Second)
				}
			}
		} else if pricePlanType == consts.FsTrial {
			//目前试用是15天
			if trailDays <= 0 {
				trailDays = consts.FsOrderTrailDays
			}
			if orgConfig.PayEndTime.After(payStartTime) {
				upd[consts.TcPayEndTime] = orgConfig.PayEndTime.AddDate(0, 0, trailDays).Add(-1 * time.Second)
			} else {
				upd[consts.TcPayStartTime] = payStartTime
				upd[consts.TcPayEndTime] = payStartTime.AddDate(0, 0, trailDays).Add(-1 * time.Second)
			}
		} else if pricePlanType == consts.FsActiveDay {
			if expireDays > 0 {
				if buyType == consts.FsOrderBuy {
					//普通支付
					if orgConfig.PayEndTime.After(payStartTime) {
						upd[consts.TcPayEndTime] = orgConfig.PayEndTime.AddDate(0, 0, expireDays).Add(-1 * time.Second)
					} else {
						upd[consts.TcPayStartTime] = payStartTime
						upd[consts.TcPayEndTime] = payStartTime.AddDate(0, 0, expireDays).Add(-1 * time.Second)
					}
				} else if buyType == consts.FsOrderRenew {
					if orgConfig.PayEndTime.After(payStartTime) {
						upd[consts.TcPayEndTime] = orgConfig.PayEndTime.AddDate(0, 0, expireDays).Add(-1 * time.Second)
					} else {
						upd[consts.TcPayStartTime] = payStartTime
						upd[consts.TcPayEndTime] = payStartTime.AddDate(0, 0, expireDays).Add(-1 * time.Second)
					}
				}
			}
		} else if pricePlanType == consts.FsActiveEndDate {
			if endDate.After(time.Now()) {
				upd[consts.TcPayStartTime] = payStartTime
				upd[consts.TcPayEndTime] = endDate
			}
		}
	}

	if sourceChannel == sdk_const.SourceChannelDingTalk {
		if pricePlanType == consts.DingOrderChargeTryout {
			// 试用
			upd[consts.TcPayStartTime] = payStartTime
			upd[consts.TcPayEndTime] = endDate
		}
		if buyType == consts.DingOrderBy {
			// 新购
			upd[consts.TcPayStartTime] = payStartTime
			upd[consts.TcPayEndTime] = endDate
		} else if buyType == consts.DingOrderRenew || (buyType == consts.DingOrderBy && pricePlanType == consts.DingOrderChargeTryout) {
			// 续费, 试用之后付费 订单类型是BUY
			if orgConfig.PayEndTime.After(payStartTime) {
				upd[consts.TcPayEndTime] = orgConfig.PayEndTime.AddDate(1, 0, 0).Add(-1 * time.Second)
			} else {
				upd[consts.TcPayStartTime] = payStartTime
				upd[consts.TcPayEndTime] = payStartTime.AddDate(1, 0, 0).Add(-1 * time.Second)
			}
		} else if buyType == consts.DingOrderRenewUpgrade {
			// 续费升配
			upd[consts.TcPayLevel] = level
		} else if buyType == consts.DingOrderRenewDegrade {
			// 续费降配
			upd[consts.TcPayLevel] = level
		} else if buyType == consts.DingOrderUpgrade {
			// 升级
			upd[consts.TcPayLevel] = level
		}
	}

	if sourceChannel == sdk_const.SourceChannelWeixin {
		upd[consts.TcPayStartTime] = payStartTime
		upd[consts.TcPayEndTime] = endDate
	}

	if len(upd) > 0 {
		_, updateErr := mysql.UpdateSmartWithCond(consts.TableOrgConfig, db.Cond{
			consts.TcOrgId: orgId,
		}, upd)
		if updateErr != nil {
			log.Error(updateErr)
			return errs.MysqlOperateError
		}
	}

	clearErr := domain.ClearOrgConfig(orgId)
	if clearErr != nil {
		log.Error(clearErr)
		return clearErr
	}

	// 上报事件
	asyn.Execute(func() {
		orgConfig, err := domain.GetOrgConfigRich(orgId)
		if err != nil {
			return
		}
		orgEvent := &commonvo.OrgEvent{}
		orgEvent.OrgId = orgId
		orgEvent.New = orgConfig

		openTraceId, _ := threadlocal.Mgr.GetValue(consts.JaegerContextTraceKey)
		openTraceIdStr := cast.ToString(openTraceId)

		report.ReportOrgEvent(msgPb.EventType_OrgConfigUpdated, openTraceIdStr, orgEvent, true)
	})

	return nil
}
