package dao

import (
	"github.com/star-table/startable-server/app/service/projectsvc/po"
	"github.com/star-table/startable-server/common/core/consts"

	"github.com/star-table/startable-server/common/library/db/mysql"
	"upper.io/db.v3/lib/sqlbuilder"
)

func InsertChatSettingBatch(pos []po.PpmProProjectChat, tx ...sqlbuilder.Tx) error {
	var err error = nil
	isTx := tx != nil && len(tx) > 0
	var batch *sqlbuilder.BatchInserter
	if !isTx {
		conn, err := mysql.GetConnect()
		if err != nil {
			return err
		}
		batch = conn.InsertInto(consts.TableProjectChat).Batch(len(pos))
	}
	if batch == nil {
		batch = tx[0].InsertInto(consts.TableProjectChat).Batch(len(pos))
	}
	go func() {
		defer batch.Done()
		for i := range pos {
			batch.Values(pos[i])
		}
	}()
	err = batch.Wait()
	if err != nil {
		log.Errorf("Project setting dao InsertBatch err %v", err)
		return err
	}
	return nil
}
