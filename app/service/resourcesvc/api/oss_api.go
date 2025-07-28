package resourcesvc

import (
	"github.com/star-table/startable-server/app/service"
	"github.com/star-table/startable-server/common/model/vo"
	"github.com/star-table/startable-server/common/model/vo/resourcevo"
)

func (PostGreeter) GetOssSignURL(req resourcevo.OssApplySignURLReqVo) resourcevo.GetOssSignURLRespVo {
	res, err := service.GetOssSignURL(req)
	return resourcevo.GetOssSignURLRespVo{Err: vo.NewErr(err), GetOssSignURL: res}
}

func (PostGreeter) GetOssPostPolicy(req resourcevo.GetOssPostPolicyReqVo) resourcevo.GetOssPostPolicyRespVo {
	res, err := service.GetOssPostPolicy(req)
	return resourcevo.GetOssPostPolicyRespVo{Err: vo.NewErr(err), GetOssPostPolicy: res}
}
