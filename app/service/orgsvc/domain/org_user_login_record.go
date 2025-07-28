package orgsvc

import (
	"time"

	"github.com/star-table/startable-server/app/facade/idfacade"
	"github.com/star-table/startable-server/common/core/consts"
	"github.com/star-table/startable-server/common/core/errs"
	"github.com/star-table/startable-server/common/library/db/mysql"
)

// 添加登录记录
func AddLoginRecord(orgId, userId int64, sourceChannel string) errs.SystemErrorInfo {
	id, err := idfacade.ApplyPrimaryIdRelaxed(consts.TableOrgUserLoginRecord)
	if err != nil {
		log.Error(err)
		return err
	}
	userLoginRecord := po.PpmOrgUserLoginRecord{
		Id:            id,
		OrgId:         orgId,
		UserId:        userId,
		LoginTime:     time.Now(),
		SourceChannel: sourceChannel,
		Creator:       userId,
		Updator:       userId,
	}

	insertErr := mysql.Insert(&userLoginRecord)
	if insertErr != nil {
		log.Error(insertErr)
		return errs.MysqlOperateError
	}

	return nil
}
