package ordersvc

import (
	"strconv"

	"github.com/star-table/startable-server/common/model/vo/ordervo"

	sconsts "github.com/star-table/startable-server/app/service"
	"github.com/star-table/startable-server/common/core/consts"
	"github.com/star-table/startable-server/common/core/errs"
	"github.com/star-table/startable-server/common/core/logger"
	"github.com/star-table/startable-server/common/core/util/json"
	"github.com/star-table/startable-server/common/library/cache"
	"github.com/star-table/startable-server/common/library/db/mysql"
	"upper.io/db.v3"
)

var log = logger.GetDefaultLogger()

func GetFunctionByLevel(level int64) ([]ordervo.FunctionLimitObj, errs.SystemErrorInfo) {
	result := make([]ordervo.FunctionLimitObj, 0)
	key := sconsts.CacheFunctionByLevel

	// 双11旗舰版大促 和 旗舰版功能一致
	if level == consts.PayLevelDouble11Activity {
		level = consts.PayLevelFlagship
	}

	infoJson, err := cache.HGet(key, strconv.FormatInt(level, 10))
	if err != nil {
		log.Errorf("[GetFunctionByLevel] HGet err: %v", err)
		return nil, errs.BuildSystemErrorInfo(errs.RedisOperateError)
	}

	if infoJson != "" {
		err = json.FromJson(infoJson, &result)
		if err != nil {
			log.Errorf("[GetFunctionByLevel] FromJson err: %v", err)
			return nil, errs.BuildSystemErrorInfo(errs.JSONConvertError)
		}
		return result, nil
	} else {
		conn, err := mysql.GetConnect()
		//defer func() {
		//	if err := conn.Close(); err != nil {
		//		log.Error(err)
		//	}
		//}()
		if err != nil {
			return nil, errs.MysqlOperateError
		}

		midPo := &[]po.FunctionLevelPo{}
		err = conn.Select(db.Raw("l.function_id, f.code, f.name, l.limit_info")).From("ppm_ord_function_level as l").
			Join("ppm_ord_function as f").On("l.function_id = f.id").Where(db.Cond{
			"l.is_delete": consts.AppIsNoDelete,
			"f.is_delete": consts.AppIsNoDelete,
			"l.level":     level,
		}).All(midPo)
		if err != nil {
			log.Errorf("[GetFunctionByLevel] select err: %v", err)
			return nil, errs.MysqlOperateError
		}

		for _, levelPo := range *midPo {
			fnLimitObj := ordervo.FunctionLimitObj{
				Limit: make([]ordervo.FunctionLimitItem, 0),
			}
			json.FromJson(levelPo.LimitInfo, &fnLimitObj)
			for i, _ := range fnLimitObj.Limit {
				fnLimitObj.Limit[i].Typ = "default"
			}
			fnLimitObj.Key = levelPo.Code
			fnLimitObj.Name = levelPo.Name
			result = append(result, fnLimitObj)
		}

		resJson, err := json.ToJson(result)
		if err != nil {
			log.Errorf("[GetFunctionByLevel] ToJson err: %v", err)
			return nil, errs.BuildSystemErrorInfo(errs.JSONConvertError)
		}
		err = cache.HSet(key, strconv.FormatInt(level, 10), resJson)
		if err != nil {
			log.Errorf("[GetFunctionByLevel] HSet err: %v", err)
			return nil, errs.BuildSystemErrorInfo(errs.RedisOperateError)
		}

		return result, nil
	}
}
