package orgsvc

import (
	"github.com/star-table/startable-server/app/service"
	"github.com/star-table/startable-server/common/model/vo"
	"github.com/star-table/startable-server/common/model/vo/orgvo"
)

func (PostGreeter) InnerUserInfo(req *orgvo.InnerUserInfosReq) *vo.DataRespVo {
	userInfos, err := service.InnerGetUserInfos(req.OrgId, req.Input.Ids)
	return &vo.DataRespVo{Err: vo.NewErr(err), Data: userInfos}
}
