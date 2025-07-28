package api

import (
	"github.com/star-table/startable-server/app/service/projectsvc/service"
	"github.com/star-table/startable-server/common/model/vo"
	"github.com/star-table/startable-server/common/model/vo/projectvo"
)

func (PostGreeter) SaveMenu(req projectvo.SaveMenuReqVo) projectvo.SaveMenuRespVo {
	res, err := service.SaveMenu(req.OrgId, req.Input)
	return projectvo.SaveMenuRespVo{
		Err:  vo.NewErr(err),
		Data: res,
	}
}

func (PostGreeter) GetMenu(req projectvo.GetMenuReqVo) projectvo.GetMenuRespVo {
	res, err := service.GetMenu(req.OrgId, req.AppId)
	return projectvo.GetMenuRespVo{
		Err:  vo.NewErr(err),
		Data: res,
	}
}
