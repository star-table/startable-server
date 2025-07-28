package api

import (
	"github.com/star-table/startable-server/app/service/projectsvc/service"
	"github.com/star-table/startable-server/common/model/vo"
	"github.com/star-table/startable-server/common/model/vo/projectvo"
)

func (PostGreeter) IterationStats(reqVo projectvo.IterationStatsReqVo) projectvo.IterationStatsRespVo {
	res, err := service.IterationStats(reqVo.OrgId, reqVo.Page, reqVo.Size, reqVo.Input)
	return projectvo.IterationStatsRespVo{Err: vo.NewErr(err), IterationStats: res}
}
