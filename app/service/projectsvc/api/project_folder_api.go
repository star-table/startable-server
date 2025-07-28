package api

import (
	"github.com/star-table/startable-server/app/service/projectsvc/service"
	"github.com/star-table/startable-server/common/model/vo"
	"github.com/star-table/startable-server/common/model/vo/projectvo"
)

func (PostGreeter) CreateProjectFolder(reqVo projectvo.CreateProjectFolderReqVo) vo.CommonRespVo {
	res, err := service.CreateProjectFolder(reqVo.OrgId, reqVo.UserId, reqVo.Input)
	return vo.CommonRespVo{Err: vo.NewErr(err), Void: res}
}

func (PostGreeter) UpdateProjectFolder(reqVo projectvo.UpdateProjectFolderReqVo) projectvo.UpdateProjectFolderRespVo {
	res, err := service.UpdateProjectFolder(reqVo.OrgId, reqVo.UserId, reqVo.Input)
	return projectvo.UpdateProjectFolderRespVo{Err: vo.NewErr(err), Output: res}
}

func (PostGreeter) DeleteProjectFolder(reqVo projectvo.DeleteProjectFolerReqVo) projectvo.DeleteProjectFolerRespVo {
	res, err := service.DeleteProjectFolder(reqVo.OrgId, reqVo.UserId, reqVo.Input)
	return projectvo.DeleteProjectFolerRespVo{Err: vo.NewErr(err), Output: res}
}

func (PostGreeter) GetProjectFolder(reqVo projectvo.GetProjectFolderReqVo) projectvo.GetProjectFolderRespVo {
	res, err := service.GetProjectFolder(reqVo.OrgId, reqVo.UserId, reqVo.Page, reqVo.Size, reqVo.Input)
	return projectvo.GetProjectFolderRespVo{Err: vo.NewErr(err), Output: res}
}
