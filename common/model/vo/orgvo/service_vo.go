package orgvo

import (
	"github.com/star-table/startable-server/common/model/bo"
	"github.com/star-table/startable-server/common/model/vo"
)

type CacheUserInfoVo struct {
	vo.Err

	CacheInfo bo.CacheUserInfoBo `json:"data"`
}
