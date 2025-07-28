package commonsvc

import (
	"github.com/star-table/startable-server/app/service"
	"github.com/star-table/startable-server/common/model/vo"
	"github.com/star-table/startable-server/common/model/vo/commonvo"
)

func (GetGreeter) IndustryList() commonvo.IndustryListRespVo {
	res, err := service.IndustryList()
	return commonvo.IndustryListRespVo{Err: vo.NewErr(err), IndustryList: res}
}
