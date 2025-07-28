package api

import (
	"github.com/star-table/startable-server/app/service/projectsvc/service"
	"github.com/star-table/startable-server/common/model/vo"
	"github.com/star-table/startable-server/common/model/vo/projectvo"
)

func (PostGreeter) InnerCreateTodoHook(reqVo *projectvo.InnerCreateTodoHookReq) *vo.VoidErr {
	return &vo.VoidErr{Err: vo.NewErr(service.CreateTodoHook(reqVo))}
}
