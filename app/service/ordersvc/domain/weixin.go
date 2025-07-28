package ordersvc

import (
	"time"

	"github.com/star-table/startable-server/common/core/errs"

	"github.com/star-table/startable-server/common/library/db/mysql"

	"github.com/star-table/startable-server/common/model/vo/ordervo"
)

func CreateLicenceOrder(req *ordervo.CreateWeiXinLicenceOrderReq) errs.SystemErrorInfo {
	m := &po.PpmOrdWeiXinLicenceOrder{
		SuitId:     req.SuitId,
		CorpId:     req.CorpId,
		OrderId:    req.OrderId,
		CreateTime: time.Now(),
	}

	err := mysql.Insert(m)
	if err != nil {
		log.Errorf("[CreateLicenceOrder] mysql.Insert model:%v, err:%v", m, err)
		return errs.BuildSystemErrorInfo(errs.MysqlOperateError, err)
	}

	return nil
}
