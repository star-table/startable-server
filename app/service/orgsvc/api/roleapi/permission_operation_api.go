package orgsvc

import (
	service "github.com/star-table/startable-server/app/service"
	"github.com/star-table/startable-server/common/model/vo"
	"github.com/star-table/startable-server/common/model/vo/rolevo"
)

// TODO 替换为lesscode-permission
func (GetGreeter) GetPersonalPermissionInfo(req rolevo.GetPersonalPermissionInfoReqVo) rolevo.GetPersonalPermissionInfoRespVo {
	res, err := service.GetPersonalPermissionInfoForFuse(req.OrgId, req.UserId, req.ProjectId, req.IssueId, req.SourceChannel)
	return rolevo.GetPersonalPermissionInfoRespVo{Err: vo.NewErr(err), Data: res}
}
