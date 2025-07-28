package api

import (
	"github.com/star-table/startable-server/app/service/projectsvc/service"
	"github.com/star-table/startable-server/common/model/vo"
	"github.com/star-table/startable-server/common/model/vo/projectvo"
)

func (GetGreeter) ProjectTypes(req projectvo.ProjectTypesReqVo) projectvo.ProjectTypesRespVo {
	resp, err := service.ProjectTypes(req.OrgId)
	return projectvo.ProjectTypesRespVo{Err: vo.NewErr(err), ProjectTypes: resp}
}
