package orgsvc

import (
	"github.com/star-table/startable-server/app/service"
	"github.com/star-table/startable-server/common/model/vo"
	"github.com/star-table/startable-server/common/model/vo/orgvo"
)

func (PostGreeter) OpenOrgUserList(reqVo orgvo.OpenOrgUserListReqVo) orgvo.OpenOrgUserListRespVo {
	res, err := service.OrgUserList(reqVo.OrgId, reqVo.UserId, reqVo.Page, reqVo.Size, reqVo.Input)
	return orgvo.OpenOrgUserListRespVo{Err: vo.NewErr(err), Data: res}
}
