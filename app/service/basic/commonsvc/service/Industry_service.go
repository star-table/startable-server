package commonsvc

import (
	"github.com/star-table/startable-server/common/core/consts"
	"github.com/star-table/startable-server/common/core/errs"
	"github.com/star-table/startable-server/common/core/logger"
	"github.com/star-table/startable-server/common/core/util/copyer"
	"github.com/star-table/startable-server/common/model/vo"
	"upper.io/db.v3"
)

var log = logger.GetDefaultLogger()

func IndustryList() (*vo.IndustryListResp, errs.SystemErrorInfo) {

	cond := db.Cond{
		consts.TcIsShow: consts.AppShowEnable,
	}

	bos, err := domain.GetIndustryBoAllList(cond)

	if err != nil {
		return nil, err
	}

	resultList := &[]*vo.IndustryResp{}

	copyErr := copyer.Copy(bos, resultList)
	if copyErr != nil {
		log.Errorf("对象copy异常: %v", copyErr)
		return nil, errs.BuildSystemErrorInfo(errs.ObjectCopyError, copyErr)
	}

	return &vo.IndustryListResp{
		List: *resultList,
	}, nil

}
