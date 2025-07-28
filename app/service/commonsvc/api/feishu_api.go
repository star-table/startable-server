package commonsvc

import (
	"github.com/star-table/startable-server/app/service"
	"github.com/star-table/startable-server/common/model/vo"
	"github.com/star-table/startable-server/common/model/vo/commonvo"
)

func (PostGreeter) UploadOssByFsImageKey(req commonvo.UploadOssByFsImageKeyReq) commonvo.UploadOssByFsImageKeyResp {
	res, err := service.UploadOssByFsImageKey(req.OrgId, req.ImageKey, req.IsApp)
	return commonvo.UploadOssByFsImageKeyResp{
		Err: vo.NewErr(err),
		Url: res,
	}
}
