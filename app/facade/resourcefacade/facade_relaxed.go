package resourcefacade

import (
	"github.com/star-table/startable-server/common/core/errs"
	"github.com/star-table/startable-server/common/model/bo"
	"github.com/star-table/startable-server/common/model/vo/resourcevo"
)

func GetResourceBoListRelaxed(page uint, size uint, cond resourcevo.GetResourceBoListCond) (*[]bo.ResourceBo, int64, errs.SystemErrorInfo) {
	respVo := GetResourceBoList(resourcevo.GetResourceBoListReqVo{
		Page:  page,
		Size:  size,
		Input: cond,
	})

	if respVo.Failure() {
		return nil, 0, respVo.Error()
	}

	return respVo.ResourceBos, respVo.Total, nil
}
