package commonsvc

import (
	"github.com/star-table/startable-server/app/service"
	"github.com/star-table/startable-server/common/model/vo"
	"github.com/star-table/startable-server/common/model/vo/commonvo"
)

func (PostGreeter) AreaLinkageList(req commonvo.AreaLinkageListReqVo) commonvo.AreaLinkageListRespVo {
	res, err := service.AreaLinkageList(req.Input)
	return commonvo.AreaLinkageListRespVo{Err: vo.NewErr(err), AreaLinkageListResp: res}
}

func (PostGreeter) AreaInfo(req commonvo.AreaInfoReqVo) commonvo.AreaInfoRespVo {

	res, err := service.OrgAreaInfo(req)

	return commonvo.AreaInfoRespVo{Err: vo.NewErr(err), AreaInfoResp: res}

}
