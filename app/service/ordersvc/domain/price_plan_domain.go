package ordersvc

import (
	"time"

	"github.com/star-table/startable-server/common/core/consts"
	"github.com/star-table/startable-server/common/core/errs"
	"github.com/star-table/startable-server/common/core/util/copyer"
	"github.com/star-table/startable-server/common/library/db/mysql"
	"github.com/star-table/startable-server/common/model/bo"
	"upper.io/db.v3"
)

func GetFsOrderLevel(outPlanId string) (int64, errs.SystemErrorInfo) {
	info := &po.PpmOrdPricePlanFs{}
	err := mysql.SelectOneByCond(consts.TablePricePlanFs, db.Cond{
		consts.TcIsDelete:  consts.AppIsNoDelete,
		consts.TcOutPlanId: outPlanId,
	}, info)
	if err != nil {
		if err == db.ErrNoMoreRows {
			log.Errorf("该方案尚未与系统关联，飞书方案id:%s", outPlanId)
			//return 0, errs.FsPricePlanNotExist
			//没有关联默认为付费版本（目前只有两个级别）
			return consts.PayLevelStandard, nil
		}
		log.Error(err)
		return 0, errs.MysqlOperateError
	}

	return info.Level, nil
}

func GetFsOrderPlanBo(outPlanId string) (*bo.OrdPricePlanFs, errs.SystemErrorInfo) {
	info := &po.PpmOrdPricePlanFs{}
	err := mysql.SelectOneByCond(consts.TablePricePlanFs, db.Cond{
		consts.TcIsDelete:  consts.AppIsNoDelete,
		consts.TcOutPlanId: outPlanId,
	}, info)
	res := &bo.OrdPricePlanFs{
		Id:         0,
		OutPlanId:  "",
		Name:       "",
		Seats:      0,
		TrialDays:  0,
		ExpireDays: 0,
		EndDate:    time.Time{},
		MonthPrice: 0,
		YearPrice:  0,
		Level:      consts.PayLevelStandard,
		Creator:    0,
		CreateTime: time.Time{},
		Updator:    0,
		UpdateTime: time.Time{},
		Version:    0,
		IsDelete:   0,
	}
	if err != nil {
		if err == db.ErrNoMoreRows {
			log.Errorf("该方案尚未与系统关联，飞书方案id:%s", outPlanId)
			//return 0, errs.FsPricePlanNotExist
			//没有关联默认为付费版本（目前只有两个级别）
			return res, nil
		}
		log.Error(err)
		return nil, errs.MysqlOperateError
	}

	_ = copyer.Copy(info, res)
	return res, nil
}

func GetDingOrderLevel(goodsCode, itemCode string) (int64, errs.SystemErrorInfo) {
	info := &po.PpmOrdPricePlanDing{}
	err := mysql.SelectOneByCond(consts.TablePricePlanDing, db.Cond{
		consts.TcIsDelete:  consts.AppIsNoDelete,
		consts.TcGoodsCode: goodsCode,
		consts.TcItemCode:  itemCode,
	}, info)
	if err != nil {
		if err == db.ErrNoMoreRows {
			log.Errorf("该方案尚未与系统关联，钉钉商品码:%s, 规格码:%s", goodsCode, itemCode)
			//return 0, errs.FsPricePlanNotExist
			return consts.PayLevelStandard, nil
		}
		log.Error(err)
		return 0, errs.MysqlOperateError
	}

	// 试用期限范围内，等级相当于旗舰版本
	if info.Level == consts.PayLevelDingFreeTrial {
		info.Level = consts.PayLevelFlagship
	}

	return info.Level, nil
}
