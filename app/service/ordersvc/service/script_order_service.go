package ordersvc

import (
	"fmt"
	"strings"
	"time"

	"github.com/star-table/startable-server/app/facade/orgfacade"
	"github.com/star-table/startable-server/common/core/consts"
	"github.com/star-table/startable-server/common/core/errs"
	"github.com/star-table/startable-server/common/core/util/str"
	"github.com/star-table/startable-server/common/library/db/mysql"
	"github.com/star-table/startable-server/common/model/vo/orgvo"
	"upper.io/db.v3"
	"upper.io/db.v3/lib/sqlbuilder"
)

// FixOrderOfZeroOrgId 修复组织 id 为 0 的订单，更新为正确的组织 id
func FixOrderOfZeroOrgId() errs.SystemErrorInfo {
	size := uint(300)
	// 查询所有的 orgId 为 0 的订单
	// 取出订单中的 tenantKey，通过 tenantKey 查询组织信息
	// 将组织信息中的 orgId 更新到订单中
	// 更新完后，数据被修复，所以查询 orgId 为 0 的无需再使用 page 进行 offset
	cmpId := int64(0)
	for {
		fsOrders, err := domain.GetOrderListByCond(nil, db.Cond{
			consts.TcOrgId: 0,
			consts.TcId:    db.Gt(cmpId), // 因为存在一些无法更新的数据，所以设定一个 offset，保证查询跳过无效的数据。
		}, nil, 1, size)
		if err != nil {
			log.Errorf("[FixOrderOfZeroOrgId] err: %v", err)
			return err
		}
		if len(fsOrders) < 1 {
			break
		}
		outOrgIds := make([]string, 0, len(fsOrders))
		for _, order := range fsOrders {
			outOrgIds = append(outOrgIds, order.TenantKey)
		}
		// 更新用于 offset 作用的主键值
		cmpId = fsOrders[len(fsOrders)-1].Id
		outOrgRes := orgfacade.GetOrgOutInfoByOutOrgIdBatch(orgvo.GetOrgOutInfoByOutOrgIdBatchReqVo{
			Input: orgvo.GetOrgOutInfoByOutOrgIdBatchReqInput{
				OutOrgIds: outOrgIds,
			},
		})
		if outOrgRes.Failure() {
			log.Errorf("[FixOrderOfZeroOrgId] GetOrgOutInfoByOutOrgIdBatch err: %v", outOrgRes.Error())
			return outOrgRes.Error()
		}
		if len(outOrgRes.Data) < 1 {
			continue
		}
		orgOutIdMap := make(map[string]int64, len(outOrgRes.Data))
		for _, item := range outOrgRes.Data {
			orgOutIdMap[item.OutOrgId] = item.OrgId
		}
		if len(orgOutIdMap) < 1 {
			log.Info("没有需要更新的组织订单。")
			return nil
		}
		// 更新订单
		dbErr := mysql.TransX(func(tx sqlbuilder.Tx) error {
			if err := UpdateFsOrderForFix(orgOutIdMap, fsOrders, tx); err != nil {
				log.Errorf("[FixOrderOfZeroOrgId] UpdateFsOrderForFix err: %v", err)
				return err
			}
			if err := UpdateOrderForFix(orgOutIdMap, fsOrders, tx); err != nil {
				log.Errorf("[FixOrderOfZeroOrgId] UpdateOrderForFix err: %v", err)
				return err
			}
			return nil
		})
		if dbErr != nil {
			log.Errorf("[FixOrderOfZeroOrgId] transX err: %v", dbErr)
			return errs.BuildSystemErrorInfo(errs.MysqlOperateError, dbErr)
		}
		time.Sleep(100 * time.Millisecond)
		log.Infof("[FixOrderOfZeroOrgId] 处理完一批 curCmpId: %d ", cmpId)
	}
	log.Infof("[FixOrderOfZeroOrgId] 脚本执行完毕！")

	return nil
}

// UpdateFsOrderForFix 更新 ppm_ord_order_fs 表
func UpdateFsOrderForFix(outOrgIdMap map[string]int64, fsOrders []po.PpmOrdOrderFs, tx sqlbuilder.Tx) errs.SystemErrorInfo {
	if len(fsOrders) < 1 {
		return nil
	}
	tx.SetLogging(true)
	// 不同订单，更新的 orgId 不同
	sqlPartStr := strings.Builder{}
	limitCount := 0
	fsOrderIds := make([]int64, 0, len(fsOrders))
	for _, fsOrder := range fsOrders {
		if tmpOrgId, ok := outOrgIdMap[fsOrder.TenantKey]; ok {
			sqlPartStr.WriteString(fmt.Sprintf(" WHEN %s THEN %d ", fsOrder.OrderId, tmpOrgId))
			limitCount += 1
			fsOrderIds = append(fsOrderIds, fsOrder.Id)
		}
	}
	// update 语句不能带有 limit，所以换一种方式。
	_, err := tx.Exec("UPDATE `" + consts.TableOrderFs + "` " +
		" SET " + " org_id=  " +
		" CASE order_id " +
		sqlPartStr.String() +
		"  END " +
		fmt.Sprintf(" WHERE id IN (%s)", str.Int64Implode(fsOrderIds, ",")) +
		fmt.Sprintf(" LIMIT %d", limitCount),
	)
	if err != nil {
		log.Errorf("[UpdateFsOrderForFix] err: %v", err)
		return errs.BuildSystemErrorInfo(errs.MysqlOperateError, err)
	}

	return nil
}

// UpdateOrderForFix 更新 ppm_ord_order 表
func UpdateOrderForFix(outOrgIdMap map[string]int64, fsOrders []po.PpmOrdOrderFs, tx sqlbuilder.Tx) errs.SystemErrorInfo {
	if len(fsOrders) < 1 {
		return nil
	}
	tx.SetLogging(true)
	// 不同订单，更新的 orgId 不同
	sqlPartStr := strings.Builder{}
	limitCount := 0
	fsOrderIds := make([]int64, 0, len(fsOrders))
	for _, fsOrder := range fsOrders {
		if tmpOrgId, ok := outOrgIdMap[fsOrder.TenantKey]; ok {
			sqlPartStr.WriteString(fmt.Sprintf(" WHEN %d THEN %d ", fsOrder.Id, tmpOrgId))
			limitCount += 1
			fsOrderIds = append(fsOrderIds, fsOrder.Id)
		}
	}
	// update 语句不能带有 limit，所以换一种方式。
	_, err := tx.Exec("UPDATE `" + consts.TableOrder + "` " +
		" SET " + "org_id=  " +
		" CASE out_order_no " +
		sqlPartStr.String() +
		"  END " +
		fmt.Sprintf(" WHERE out_order_no IN (%s)", str.Int64Implode(fsOrderIds, ",")) +
		fmt.Sprintf(" LIMIT %d", len(fsOrderIds)),
	)
	if err != nil {
		log.Errorf("[UpdateOrderForFix] err: %v", err)
		return errs.BuildSystemErrorInfo(errs.MysqlOperateError, err)
	}

	return nil
}
