package commonfacade

import (
	"fmt"

	"github.com/star-table/startable-server/app/facade"
	"github.com/star-table/startable-server/common/core/config"
	"github.com/star-table/startable-server/common/core/consts"
	"github.com/star-table/startable-server/common/model/vo/commonvo"
)

func IndustryList() commonvo.IndustryListRespVo {
	respVo := &commonvo.IndustryListRespVo{}
	reqUrl := fmt.Sprintf("%s/api/commonsvc/industryList", config.GetPreUrl("commonsvc"))
	err := facade.Request(consts.HttpMethodGet, reqUrl, nil, nil, nil, respVo)
	if err.Failure() {
		respVo.Err = err
	}
	return *respVo
}
