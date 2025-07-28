package orgsvc

import (
	sconsts "github.com/star-table/startable-server/app/service"
	"github.com/star-table/startable-server/common/core/consts"
	"github.com/star-table/startable-server/common/core/errs"
	"github.com/star-table/startable-server/common/core/util/copyer"
	"github.com/star-table/startable-server/common/core/util/json"
	"github.com/star-table/startable-server/common/library/cache"
	"github.com/star-table/startable-server/common/library/db/mysql"
	bo2 "github.com/star-table/startable-server/common/model/bo"
	"upper.io/db.v3"
)

func GetPermissionOperationList() (*[]bo2.PermissionOperationBo, errs.SystemErrorInfo) {
	key := sconsts.CachePermissionOperationList
	listJson, err := cache.Get(key)
	if err != nil {
		log.Error(err)
		return nil, errs.BuildSystemErrorInfo(errs.RedisOperateError, err)
	}

	bo := &[]bo2.PermissionOperationBo{}
	if listJson != "" {
		err := json.FromJson(listJson, bo)
		if err != nil {
			return nil, errs.BuildSystemErrorInfo(errs.JSONConvertError)
		}
	} else {
		// TODO 切换lesscode-permission
		po := &[]po.PpmRolPermissionOperation{}
		selectErr := mysql.SelectAllByCond(consts.TablePermissionOperation, db.Cond{
			consts.TcIsDelete: consts.AppIsNoDelete,
			consts.TcStatus:   consts.AppStatusEnable,
		}, po)
		if selectErr != nil {
			log.Error(selectErr)
			return nil, errs.BuildSystemErrorInfo(errs.MysqlOperateError, selectErr)
		}
		_ = copyer.Copy(po, bo)
		listJson, err = json.ToJson(bo)
		if err != nil {
			log.Error(err)
			return nil, errs.BuildSystemErrorInfo(errs.JSONConvertError)
		}
		err = cache.SetEx(key, listJson, consts.GetCacheBaseExpire())
		if err != nil {
			return nil, errs.BuildSystemErrorInfo(errs.RedisOperateError, err)
		}
	}

	return bo, nil
}

func GetPermissionOperationListByPermissionId(permissionId int64) ([]bo2.PermissionOperationBo, errs.SystemErrorInfo) {
	res := []bo2.PermissionOperationBo{}

	list, err := GetPermissionOperationList()
	if err != nil {
		log.Error(err)
		return nil, err
	}

	for _, v := range *list {
		//过滤所有读的操作项，默认拥有
		if v.IsShow == consts.AppShowEnable && v.PermissionId == permissionId && v.OperationCodes != consts.RoleOperationView {
			res = append(res, v)
		}
	}

	return res, nil
}
