package commonsvc

import (
	"github.com/star-table/startable-server/common/core/consts"
)

var (
	//对象id缓存key前缀
	CacheObjectIdPreKey = consts.CacheKeyPrefix + consts.IdsvcApplicationName + consts.CacheKeyOfSys + "object_id:"
)
