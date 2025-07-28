package api

import (
	"github.com/star-table/startable-server/app/service/projectsvc/service"
	"github.com/star-table/startable-server/common/model/vo"
	"github.com/star-table/startable-server/common/model/vo/projectvo"
)

func (PostGreeter) MirrorsStat(reqVo *projectvo.MirrorsStatReq) projectvo.MirrorsStatResp {
	res, err := service.MirrorStat(reqVo.OrgId, reqVo.UserId, reqVo.Input.AppIds)
	return projectvo.MirrorsStatResp{
		Err: vo.NewErr(err),
		Data: vo.MirrorsStatResp{
			DataStat: res,
		},
	}
}
