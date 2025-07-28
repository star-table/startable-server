package api

import (
	"github.com/star-table/startable-server/app/service/projectsvc/service"
	"github.com/star-table/startable-server/common/model/vo"
	"github.com/star-table/startable-server/common/model/vo/projectvo"
)

func (PostGreeter) AuthProjectPermission(reqVo projectvo.AuthProjectPermissionReqVo) vo.CommonRespVo {
	input := reqVo.Input
	err := service.AuthProjectPermission(input.OrgId, input.UserId, input.ProjectId, input.Path, input.Operation, input.AuthFiling)
	return vo.CommonRespVo{Err: vo.NewErr(err), Void: nil}
}
