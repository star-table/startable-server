package commonvo

import (
	"github.com/star-table/startable-server/common/model/vo"
)

type IndustryListRespVo struct {
	vo.Err
	IndustryList *vo.IndustryListResp `json:"data"`
}
