package commonsvc

import (
	commonPo "github.com/star-table/startable-server/app/service/commonsvc/po"
	"github.com/star-table/startable-server/common/core/consts"
	"github.com/star-table/startable-server/common/core/logger"
	"github.com/star-table/startable-server/common/library/db/mysql"
	"github.com/star-table/startable-server/common/model/bo"
	"upper.io/db.v3"
	"upper.io/db.v3/lib/sqlbuilder"
)

var log = logger.GetDefaultLogger()

func InsertArea(commonPo commonPo.PpmCmmArea, tx ...sqlbuilder.Tx) error {
	var err error = nil
	if tx != nil && len(tx) > 0 {
		err = mysql.TransInsert(tx[0], &commonPo)
	} else {
		err = mysql.Insert(&commonPo)
	}
	if err != nil {
		log.Errorf("Area dao Insert err %v", err)
	}
	return nil
}

func InsertAreaBatch(commonPos []commonPo.PpmCmmArea, tx ...sqlbuilder.Tx) error {
	var err error = nil
	isTx := tx != nil && len(tx) > 0
	var batch *sqlbuilder.BatchInserter
	if !isTx {
		conn, err := mysql.GetConnect()
		if err != nil {
			return err
		}
		batch = conn.InsertInto(consts.TableArea).Batch(len(commonPos))
	}
	if batch == nil {
		batch = tx[0].InsertInto(consts.TableArea).Batch(len(commonPos))
	}
	go func() {
		defer batch.Done()
		for i := range commonPos {
			batch.Values(commonPos[i])
		}
	}()
	err = batch.Wait()
	if err != nil {
		log.Errorf("Iteration dao InsertBatch err %v", err)
		return err
	}
	return nil
}

func UpdateArea(commonPo commonPo.PpmCmmArea, tx ...sqlbuilder.Tx) error {
	var err error = nil
	if tx != nil && len(tx) > 0 {
		err = mysql.TransUpdate(tx[0], &commonPo)
	} else {
		err = mysql.Update(&commonPo)
	}
	if err != nil {
		log.Errorf("Area dao Update err %v", err)
	}
	return err
}

func UpdateAreaById(id int64, upd mysql.Upd, tx ...sqlbuilder.Tx) (int64, error) {
	return UpdateAreaByCond(db.Cond{
		consts.TcId:       id,
		consts.TcIsDelete: consts.AppIsNoDelete,
	}, upd, tx...)
}

func UpdateAreaByOrg(id int64, orgId int64, upd mysql.Upd, tx ...sqlbuilder.Tx) (int64, error) {
	return UpdateAreaByCond(db.Cond{
		consts.TcId:       id,
		consts.TcOrgId:    orgId,
		consts.TcIsDelete: consts.AppIsNoDelete,
	}, upd, tx...)
}

func UpdateAreaByCond(cond db.Cond, upd mysql.Upd, tx ...sqlbuilder.Tx) (int64, error) {
	var mod int64 = 0
	var err error = nil
	if tx != nil && len(tx) > 0 {
		mod, err = mysql.TransUpdateSmartWithCond(tx[0], consts.TableArea, cond, upd)
	} else {
		mod, err = mysql.UpdateSmartWithCond(consts.TableArea, cond, upd)
	}
	if err != nil {
		log.Errorf("Area dao Update err %v", err)
	}
	return mod, err
}

func DeleteAreaById(id int64, operatorId int64, tx ...sqlbuilder.Tx) (int64, error) {
	upd := mysql.Upd{
		consts.TcIsDelete: consts.AppIsDeleted,
	}
	if operatorId > 0 {
		upd[consts.TcUpdator] = operatorId
	}
	return UpdateAreaByCond(db.Cond{
		consts.TcId:       id,
		consts.TcIsDelete: consts.AppIsNoDelete,
	}, upd, tx...)
}

func DeleteAreaByOrg(id int64, orgId int64, operatorId int64, tx ...sqlbuilder.Tx) (int64, error) {
	upd := mysql.Upd{
		consts.TcIsDelete: consts.AppIsDeleted,
	}
	if operatorId > 0 {
		upd[consts.TcUpdator] = operatorId
	}
	return UpdateAreaByCond(db.Cond{
		consts.TcId:       id,
		consts.TcOrgId:    orgId,
		consts.TcIsDelete: consts.AppIsNoDelete,
	}, upd, tx...)
}

func SelectAreaById(id int64) (*commonPo.PpmCmmArea, error) {
	commonPo := &commonPo.PpmCmmArea{}
	err := mysql.SelectById(consts.TableArea, id, commonPo)
	if err != nil {
		log.Errorf("Area dao SelectById err %v", err)
	}
	return commonPo, err
}

func SelectAreaByIdAndOrg(id int64, orgId int64) (*commonPo.PpmCmmArea, error) {
	commonPo := &commonPo.PpmCmmArea{}
	err := mysql.SelectOneByCond(consts.TableArea, db.Cond{
		consts.TcId:       id,
		consts.TcOrgId:    orgId,
		consts.TcIsDelete: consts.AppIsNoDelete,
	}, commonPo)
	if err != nil {
		log.Errorf("Area dao Select err %v", err)
	}
	return commonPo, err
}

func SelectArea(cond db.Cond) (*[]commonPo.PpmCmmArea, error) {
	commonPos := &[]commonPo.PpmCmmArea{}
	err := mysql.SelectAllByCond(consts.TableArea, cond, commonPos)
	if err != nil {
		log.Errorf("Area dao SelectList err %v", err)
	}
	return commonPos, err
}

func SelectOneArea(cond db.Cond) (*commonPo.PpmCmmArea, error) {
	commonPo := &commonPo.PpmCmmArea{}
	err := mysql.SelectOneByCond(consts.TableArea, cond, commonPo)
	if err != nil {
		log.Errorf("Area dao Select err %v", err)
	}
	return commonPo, err
}

func SelectAreaByPage(cond db.Cond, pageBo bo.PageBo) (*[]commonPo.PpmCmmArea, uint64, error) {
	commonPos := &[]commonPo.PpmCmmArea{}
	total, err := mysql.SelectAllByCondWithPageAndOrder(consts.TableArea, cond, nil, pageBo.Page, pageBo.Size, pageBo.Order, commonPos)
	if err != nil {
		log.Errorf("Area dao SelectPage err %v", err)
	}
	return commonPos, total, err
}
