package api

import (
	"github.com/star-table/startable-server/app/service/projectsvc/service"
	"github.com/star-table/startable-server/common/model/vo"
	"github.com/star-table/startable-server/common/model/vo/projectvo"
)

func (PostGreeter) ConvertCode(reqVo projectvo.ConvertCodeReqVo) projectvo.ConvertCodeRespVo {
	res, err := service.ConvertCode(reqVo.OrgId, reqVo.Input)
	return projectvo.ConvertCodeRespVo{Err: vo.NewErr(err), ConvertCode: res}
}
