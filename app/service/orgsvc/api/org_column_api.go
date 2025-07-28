package orgsvc

import (
	"github.com/star-table/startable-server/app/service"
	"github.com/star-table/startable-server/common/model/vo"
	"github.com/star-table/startable-server/common/model/vo/orgvo"
)

func (PostGreeter) CreateOrgColumn(req orgvo.CreateOrgColumnReq) orgvo.CreateOrgColumnRespVo {
	resp, err := service.CreateOrgColumn(req)
	return orgvo.CreateOrgColumnRespVo{
		Err:  vo.NewErr(err),
		Data: resp,
	}
}

func (PostGreeter) GetOrgColumns(req orgvo.GetOrgColumnsReq) orgvo.GetOrgColumnsRespVo {
	resp, err := service.GetOrgColumns(req)
	return orgvo.GetOrgColumnsRespVo{
		Err:  vo.NewErr(err),
		Data: resp,
	}
}

func (PostGreeter) DeleteOrgColumn(req orgvo.DeleteOrgColumnReq) orgvo.DeleteOrgColumnRespVo {
	resp, err := service.DeleteOrgColumn(req)
	return orgvo.DeleteOrgColumnRespVo{
		Err:  vo.NewErr(err),
		Data: resp,
	}
}
