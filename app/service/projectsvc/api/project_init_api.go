package api

import (
	"github.com/star-table/startable-server/app/service/projectsvc/service"
	"github.com/star-table/startable-server/common/model/vo"
	"github.com/star-table/startable-server/common/model/vo/projectvo"
)

func (PostGreeter) ProjectInit(req projectvo.ProjectInitReqVo) projectvo.ProjectInitRespVo {
	respVo := projectvo.ProjectInitRespVo{}

	//_ = mysql.TransX(func(tx sqlbuilder.Tx) error {
	//	contextMap, err := service.ProjectInit(req.OrgId, tx)
	//
	//	respVo.ContextMap = contextMap
	//	respVo.Err = vo.NewErr(err)
	//
	//	return err
	//})
	return respVo
}

// 初始化组织的时候创建左侧目录应用、视图
func (PostGreeter) CreateOrgDirectoryAppsAndViews(reqVo projectvo.CreateOrgDirectoryAppsReq) vo.CommonRespVo {
	err := service.CreateOrgDirectoryAppsAndViews(reqVo)
	return vo.CommonRespVo{
		Err:  vo.NewErr(err),
		Void: &vo.Void{ID: 1},
	}
}
