package resourcesvc

import (
	resourcePo "github.com/star-table/startable-server/app/service/resourcesvc/po"
	"github.com/star-table/startable-server/common/core/consts"
	"github.com/star-table/startable-server/common/library/db/mysql"
	"github.com/star-table/startable-server/common/model/bo"
	"upper.io/db.v3"
	"upper.io/db.v3/lib/sqlbuilder"
)

func InsertResource(po resourcePo.PpmResResource, tx ...sqlbuilder.Tx) error {
	var err error = nil
	if tx != nil && len(tx) > 0 {
		err = mysql.TransInsert(tx[0], &po)
	} else {
		err = mysql.Insert(&po)
	}
	if err != nil {
		log.Errorf("Resource dao Insert err %v", err)
	}
	return nil
}

func InsertResourceBatch(pos []resourcePo.PpmResResource, tx ...sqlbuilder.Tx) error {
	var err error = nil

	isTx := tx != nil && len(tx) > 0

	var batch *sqlbuilder.BatchInserter

	if !isTx {
		//没有事务
		conn, err := mysql.GetConnect()
		if err != nil {
			return err
		}

		batch = conn.InsertInto(consts.TableResource).Batch(len(pos))
	}

	if batch == nil {
		batch = tx[0].InsertInto(consts.TableResource).Batch(len(pos))
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

func UpdateResource(po resourcePo.PpmResResource, tx ...sqlbuilder.Tx) error {
	var err error = nil
	if tx != nil && len(tx) > 0 {
		err = mysql.TransUpdate(tx[0], &po)
	} else {
		err = mysql.Update(&po)
	}
	if err != nil {
		log.Errorf("Resource dao Update err %v", err)
	}
	return err
}

func UpdateResourceById(id int64, upd mysql.Upd, tx ...sqlbuilder.Tx) (int64, error) {
	return UpdateResourceByCond(db.Cond{
		consts.TcId:       id,
		consts.TcIsDelete: consts.AppIsNoDelete,
	}, upd, tx...)
}

func UpdateResourceByOrg(id int64, orgId int64, upd mysql.Upd, tx ...sqlbuilder.Tx) (int64, error) {
	return UpdateResourceByCond(db.Cond{
		consts.TcId:       id,
		consts.TcOrgId:    orgId,
		consts.TcIsDelete: consts.AppIsNoDelete,
	}, upd, tx...)
}

func UpdateResourceByCond(cond db.Cond, upd mysql.Upd, tx ...sqlbuilder.Tx) (int64, error) {
	var mod int64 = 0
	var err error = nil
	if tx != nil && len(tx) > 0 {
		mod, err = mysql.TransUpdateSmartWithCond(tx[0], consts.TableResource, cond, upd)
	} else {
		mod, err = mysql.UpdateSmartWithCond(consts.TableResource, cond, upd)
	}
	if err != nil {
		log.Errorf("Resource dao Update err %v", err)
	}
	return mod, err
}

func DeleteResourceById(id int64, operatorId int64, tx ...sqlbuilder.Tx) (int64, error) {
	upd := mysql.Upd{
		consts.TcIsDelete: consts.AppIsDeleted,
	}
	if operatorId > 0 {
		upd[consts.TcUpdator] = operatorId
	}
	return UpdateResourceByCond(db.Cond{
		consts.TcId:       id,
		consts.TcIsDelete: consts.AppIsNoDelete,
	}, upd, tx...)
}

func DeleteResourceByOrg(id int64, orgId int64, operatorId int64, tx ...sqlbuilder.Tx) (int64, error) {
	upd := mysql.Upd{
		consts.TcIsDelete: consts.AppIsDeleted,
	}
	if operatorId > 0 {
		upd[consts.TcUpdator] = operatorId
	}
	return UpdateResourceByCond(db.Cond{
		consts.TcId:       id,
		consts.TcOrgId:    orgId,
		consts.TcIsDelete: consts.AppIsNoDelete,
	}, upd, tx...)
}

func SelectResourceById(id int64) (*resourcePo.PpmResResource, error) {
	po := &resourcePo.PpmResResource{}
	err := mysql.SelectById(consts.TableResource, id, po)
	if err != nil {
		log.Errorf("Resource dao SelectById err %v", err)
	}
	return po, err
}

func SelectResourceByIdAndOrg(id int64, orgId int64) (*resourcePo.PpmResResource, error) {
	po := &resourcePo.PpmResResource{}
	err := mysql.SelectOneByCond(consts.TableResource, db.Cond{
		consts.TcId:       id,
		consts.TcOrgId:    orgId,
		consts.TcIsDelete: consts.AppIsNoDelete,
	}, po)
	if err != nil {
		log.Errorf("Resource dao Select err %v", err)
	}
	return po, err
}

func SelectResource(cond db.Cond) (*[]resourcePo.PpmResResource, error) {
	pos := &[]resourcePo.PpmResResource{}
	err := mysql.SelectAllByCond(consts.TableResource, cond, pos)
	if err != nil {
		log.Errorf("Resource dao SelectList err %v", err)
	}
	return pos, err
}

func SelectOneResource(cond db.Cond) (*resourcePo.PpmResResource, error) {
	po := &resourcePo.PpmResResource{}
	err := mysql.SelectOneByCond(consts.TableResource, cond, po)
	if err != nil {
		log.Errorf("Resource dao Select err %v", err)
	}
	return po, err
}

func SelectResourceByPage(cond db.Cond, pageBo bo.PageBo) (*[]resourcePo.PpmResResource, uint64, error) {
	pos := &[]resourcePo.PpmResResource{}
	total, err := mysql.SelectAllByCondWithPageAndOrder(consts.TableResource, cond, nil, pageBo.Page, pageBo.Size, pageBo.Order, pos)
	if err != nil {
		log.Errorf("Resource dao SelectPage err %v", err)
	}
	return pos, total, err
}
