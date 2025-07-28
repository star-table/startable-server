package api

import (
	"github.com/star-table/startable-server/app/service/projectsvc/service"
	"github.com/star-table/startable-server/common/model/vo"
	"github.com/star-table/startable-server/common/model/vo/projectvo"
)

func (PostGreeter) TodoUrge(reqVo *projectvo.TodoUrgeReq) *vo.CommonRespVo {
	err := service.TodoUrge(reqVo)
	return &vo.CommonRespVo{Err: vo.NewErr(err)}
}
