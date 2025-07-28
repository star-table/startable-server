package api

import (
	"github.com/star-table/startable-server/app/service/projectsvc/service"
	"github.com/star-table/startable-server/common/model/vo"
	"github.com/star-table/startable-server/common/model/vo/projectvo"
)

func (PostGreeter) DeleteProjectAttachment(reqVo projectvo.DeleteProjectAttachmentReqVo) projectvo.DeleteProjectAttachmentRespVo {
	res, err := service.DeleteProjectAttachment(reqVo.OrgId, reqVo.UserId, reqVo.Input)
	return projectvo.DeleteProjectAttachmentRespVo{Err: vo.NewErr(err), Output: res}
}

func (PostGreeter) GetProjectAttachment(reqVo projectvo.GetProjectAttachmentReqVo) projectvo.GetProjectAttachmentRespVo {
	if reqVo.Size == 0 {
		reqVo.Size = 20
	}
	res, err := service.GetProjectAttachment(reqVo.OrgId, reqVo.UserId, reqVo.Page, reqVo.Size, reqVo.Input)
	return projectvo.GetProjectAttachmentRespVo{Err: vo.NewErr(err), Output: res}
}

func (PostGreeter) GetProjectAttachmentInfo(reqVo projectvo.GetProjectAttachmentInfoReqVo) projectvo.GetProjectAttachmentInfoRespVo {
	res, err := service.GetProjectAttachmentInfo(reqVo.OrgId, reqVo.UserId, reqVo.Input)
	return projectvo.GetProjectAttachmentInfoRespVo{Err: vo.NewErr(err), Output: res}
}
