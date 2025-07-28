package ordersvc

import (
	sdk_const "gitea.bjx.cloud/allstar/platform-sdk/consts"
	"github.com/star-table/startable-server/app/facade/orgfacade"
	"github.com/star-table/startable-server/common/core/consts"
	"github.com/star-table/startable-server/common/core/errs"
	"github.com/star-table/startable-server/common/library/db/mysql"
	"github.com/star-table/startable-server/common/model/vo/ordervo"
	"github.com/star-table/startable-server/common/model/vo/orgvo"
	"upper.io/db.v3"
)

func GetOrderListByCond(orgIds []int64, cond db.Cond, orderBy interface{}, page, size uint) ([]po.PpmOrdOrderFs, errs.SystemErrorInfo) {
	fsOrder := make([]po.PpmOrdOrderFs, 0)
	conn, dbErr := mysql.GetConnect()
	if dbErr != nil {
		log.Errorf("[GetOrderListByCond] err: %v", dbErr)
		return fsOrder, errs.BuildSystemErrorInfo(errs.MysqlOperateError, dbErr)
	}
	if len(orgIds) > 0 {
		cond[consts.TcOrgId] = db.In(orgIds)
	}
	mid := conn.Select("*").From(consts.TableOrderFs).Where(cond)
	defaultPage, defaultSize := uint(1), uint(20)
	if size > 0 && page > 0 {
		defaultPage = page
		defaultSize = size
	}
	if orderBy != nil {
		mid = mid.OrderBy(orderBy)
	} else {
		mid = mid.OrderBy(db.Raw("id asc"))
	}
	if err := mid.Paginate(defaultSize).Page(defaultPage).All(&fsOrder); err != nil {
		log.Errorf("[GetOrderListByCond] find err: %v", err)
		return fsOrder, errs.BuildSystemErrorInfo(errs.MysqlOperateError, err)
	}

	return fsOrder, nil
}

func GetOrderPayInfo(orgId int64) (*ordervo.GetOrderPayInfo, errs.SystemErrorInfo) {
	baseOrgInfo := orgfacade.GetBaseOrgInfo(orgvo.GetBaseOrgInfoReqVo{OrgId: orgId})
	if baseOrgInfo.Failure() {
		log.Errorf("[GetOrderPayInfo] GetBaseOrgInfo err:%v, orgId:%v", baseOrgInfo.Error(), orgId)
		return nil, baseOrgInfo.Error()
	}
	sourceChannel := baseOrgInfo.BaseOrgInfo.SourceChannel
	switch sourceChannel {
	case sdk_const.SourceChannelFeishu:
		return GetFsOrderInfoByOrgId(orgId)
	case sdk_const.SourceChannelDingTalk:
		return GetDingOrderInfoByOrgId(orgId)
	case sdk_const.SourceChannelWeixin:
		return GetWeiXinOrderInfoByOrgId(orgId)
	default:
		return nil, errs.SourceNotExist
	}
}

func GetFsOrderInfoByOrgId(orgId int64) (*ordervo.GetOrderPayInfo, errs.SystemErrorInfo) {
	notInPriceType := []string{consts.FsTrial, consts.FsPermanent, consts.FsActiveDay}
	orderFsPos := []*po.PpmOrdOrderFs{}
	_, err := mysql.SelectAllByCondWithPageAndOrder(consts.TableOrderFs, db.Cond{
		consts.TcOrgId:         orgId,
		consts.TcIsDelete:      consts.AppIsNoDelete,
		consts.TcPricePlanType: db.NotIn(notInPriceType),
	}, nil, 1, 1, "id desc", &orderFsPos)
	if err != nil {
		log.Errorf("[GetFsOrderInfoByOrgId] err:%v, orgId:%v", err, orgId)
		return nil, errs.MysqlOperateError
	}
	payInfo := &ordervo.GetOrderPayInfo{PayNum: 0}
	if len(orderFsPos) > 0 {
		payInfo.PayNum = orderFsPos[0].Seats
	}
	return payInfo, nil
}

func GetDingOrderInfoByOrgId(orgId int64) (*ordervo.GetOrderPayInfo, errs.SystemErrorInfo) {
	orderDing := []*po.PpmOrdOrderDing{}
	_, err := mysql.SelectAllByCondWithPageAndOrder(consts.TableOrderDing, db.Cond{
		consts.TcOrgId:    orgId,
		consts.TcIsDelete: consts.AppIsNoDelete,
		consts.TcItemCode: db.NotIn([]string{consts.DingFreeCode, consts.DingTrailCode}),
	}, nil, 1, 1, "id desc", &orderDing)
	if err != nil {
		log.Errorf("[GetDingOrderInfoByOrgId] err:%v, orgId:%v", err, orgId)
		return nil, errs.MysqlOperateError
	}
	payInfo := &ordervo.GetOrderPayInfo{PayNum: 0}
	if len(orderDing) > 0 {
		payInfo.PayNum = orderDing[0].Quantity
	}
	return payInfo, nil
}

func GetWeiXinOrderInfoByOrgId(orgId int64) (*ordervo.GetOrderPayInfo, errs.SystemErrorInfo) {
	orderWeiXin := []*po.PpmOrdOrderWeiXin{}
	_, err := mysql.SelectAllByCondWithPageAndOrder(consts.TableOrderWeiXin, db.Cond{
		consts.TcOrgId:     orgId,
		consts.TcIsDelete:  consts.AppIsNoDelete,
		consts.TcEditionId: db.NotEq(consts.WeiXinStandardId),
	}, nil, 1, 1, "id desc", &orderWeiXin)
	if err != nil {
		log.Errorf("[GetDingOrderInfoByOrgId] err:%v, orgId:%v", err, orgId)
		return nil, errs.MysqlOperateError
	}
	payInfo := &ordervo.GetOrderPayInfo{PayNum: 0}
	if len(orderWeiXin) > 0 {
		payInfo.PayNum = orderWeiXin[0].UserCount
	}
	return payInfo, nil
}

func CheckDingFreeOrder(orgId int64) bool {
	dingOrder := []*po.PpmOrdOrderDing{}
	err := mysql.SelectAllByCond(consts.TableOrderDing, db.Cond{
		consts.TcOrgId:    orgId,
		consts.TcIsDelete: consts.AppIsNoDelete,
	}, &dingOrder)
	if err != nil {
		return false
	}
	for _, order := range dingOrder {
		if order.ItemCode == consts.DingFreeCode {
			return true
		}
	}
	return false
}
