package dao

import (
	"github.com/star-table/startable-server/app/service/projectsvc/po"
	"github.com/star-table/startable-server/common/core/consts"

	"github.com/star-table/startable-server/common/library/db/mysql"
	"github.com/star-table/startable-server/common/model/bo"
	"upper.io/db.v3"
	"upper.io/db.v3/lib/sqlbuilder"
)

func InsertIssueSource(po po.PpmPrsIssueSource, tx ...sqlbuilder.Tx) error {
	var err error = nil
	if tx != nil && len(tx) > 0 {
		err = mysql.TransInsert(tx[0], &po)
	} else {
		err = mysql.Insert(&po)
	}
	if err != nil {
		log.Errorf("IssueSource dao Insert err %v", err)
	}
	return nil
}

func InsertIssueSourceBatch(pos []po.PpmPrsIssueSource, tx ...sqlbuilder.Tx) error {
	var err error = nil
	isTx := tx != nil && len(tx) > 0
	var batch *sqlbuilder.BatchInserter
	if !isTx {
		conn, err := mysql.GetConnect()
		if err != nil {
			return err
		}
		batch = conn.InsertInto(consts.TableIssueSource).Batch(len(pos))
	}
	if batch == nil {
		batch = tx[0].InsertInto(consts.TableIssueSource).Batch(len(pos))
	}
	go func() {
		defer batch.Done()
		for i := range pos {
			batch.Values(pos[i])
		}
	}()
	err = batch.Wait()
	if err != nil {
		log.Errorf("Iteration dao InsertBatch err %v", err)
		return err
	}
	return nil
}

func UpdateIssueSource(po po.PpmPrsIssueSource, tx ...sqlbuilder.Tx) error {
	var err error = nil
	if tx != nil && len(tx) > 0 {
		err = mysql.TransUpdate(tx[0], &po)
	} else {
		err = mysql.Update(&po)
	}
	if err != nil {
		log.Errorf("IssueSource dao Update err %v", err)
	}
	return err
}

func UpdateIssueSourceById(id int64, upd mysql.Upd, tx ...sqlbuilder.Tx) (int64, error) {
	return UpdateIssueSourceByCond(db.Cond{
		consts.TcId:       id,
		consts.TcIsDelete: consts.AppIsNoDelete,
	}, upd, tx...)
}

func UpdateIssueSourceByOrg(id int64, orgId int64, upd mysql.Upd, tx ...sqlbuilder.Tx) (int64, error) {
	return UpdateIssueSourceByCond(db.Cond{
		consts.TcId:       id,
		consts.TcOrgId:    orgId,
		consts.TcIsDelete: consts.AppIsNoDelete,
	}, upd, tx...)
}

func UpdateIssueSourceByCond(cond db.Cond, upd mysql.Upd, tx ...sqlbuilder.Tx) (int64, error) {
	var mod int64 = 0
	var err error = nil
	if tx != nil && len(tx) > 0 {
		mod, err = mysql.TransUpdateSmartWithCond(tx[0], consts.TableIssueSource, cond, upd)
	} else {
		mod, err = mysql.UpdateSmartWithCond(consts.TableIssueSource, cond, upd)
	}
	if err != nil {
		log.Errorf("IssueSource dao Update err %v", err)
	}
	return mod, err
}

func DeleteIssueSourceById(id int64, operatorId int64, tx ...sqlbuilder.Tx) (int64, error) {
	upd := mysql.Upd{
		consts.TcIsDelete: consts.AppIsDeleted,
	}
	if operatorId > 0 {
		upd[consts.TcUpdator] = operatorId
	}
	return UpdateIssueSourceByCond(db.Cond{
		consts.TcId:       id,
		consts.TcIsDelete: consts.AppIsNoDelete,
	}, upd, tx...)
}

func DeleteIssueSourceByOrg(id int64, orgId int64, operatorId int64, tx ...sqlbuilder.Tx) (int64, error) {
	upd := mysql.Upd{
		consts.TcIsDelete: consts.AppIsDeleted,
	}
	if operatorId > 0 {
		upd[consts.TcUpdator] = operatorId
	}
	return UpdateIssueSourceByCond(db.Cond{
		consts.TcId:       id,
		consts.TcOrgId:    orgId,
		consts.TcIsDelete: consts.AppIsNoDelete,
	}, upd, tx...)
}

func SelectIssueSourceById(id int64) (*po.PpmPrsIssueSource, error) {
	po := &po.PpmPrsIssueSource{}
	err := mysql.SelectById(consts.TableIssueSource, id, po)
	if err != nil {
		log.Errorf("IssueSource dao SelectById err %v", err)
	}
	return po, err
}

func SelectIssueSourceByIdAndOrg(id int64, orgId int64) (*po.PpmPrsIssueSource, error) {
	po := &po.PpmPrsIssueSource{}
	err := mysql.SelectOneByCond(consts.TableIssueSource, db.Cond{
		consts.TcId:       id,
		consts.TcOrgId:    orgId,
		consts.TcIsDelete: consts.AppIsNoDelete,
	}, po)
	if err != nil {
		log.Errorf("IssueSource dao Select err %v", err)
	}
	return po, err
}

func SelectIssueSource(cond db.Cond) (*[]po.PpmPrsIssueSource, error) {
	pos := &[]po.PpmPrsIssueSource{}
	err := mysql.SelectAllByCond(consts.TableIssueSource, cond, pos)
	if err != nil {
		log.Errorf("IssueSource dao SelectList err %v", err)
	}
	return pos, err
}

func SelectOneIssueSource(cond db.Cond) (*po.PpmPrsIssueSource, error) {
	po := &po.PpmPrsIssueSource{}
	err := mysql.SelectOneByCond(consts.TableIssueSource, cond, po)
	if err != nil {
		log.Errorf("IssueSource dao Select err %v", err)
	}
	return po, err
}

func SelectIssueSourceByPage(cond db.Cond, pageBo bo.PageBo) (*[]po.PpmPrsIssueSource, uint64, error) {
	pos := &[]po.PpmPrsIssueSource{}
	total, err := mysql.SelectAllByCondWithPageAndOrder(consts.TableIssueSource, cond, nil, pageBo.Page, pageBo.Size, pageBo.Order, pos)
	if err != nil {
		log.Errorf("IssueSource dao SelectPage err %v", err)
	}
	return pos, total, err
}
