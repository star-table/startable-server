package orgsvc

import (
	"github.com/star-table/startable-server/app/service"
	"github.com/star-table/startable-server/common/model/vo"
	"github.com/star-table/startable-server/common/model/vo/orgvo"
)

func (PostGreeter) GetDingJsAPISign(input orgvo.GetDingApiSignReq) orgvo.GetJsAPISignRespVo {
	res, err := service.GetDingJsAPISign(input.Input)
	return orgvo.GetJsAPISignRespVo{Err: vo.NewErr(err), GetJsAPISign: res}
}

func (PostGreeter) AuthDingCode(input orgvo.AuthDingCodeReqVo) orgvo.AuthDingCodeRespVo {
	res, err := service.AuthDingCode(input)
	return orgvo.AuthDingCodeRespVo{
		Err:  vo.NewErr(err),
		Data: res,
	}
}

func (PostGreeter) CreateCoolApp(input orgvo.CreateCoolAppReq) vo.CommonRespVo {
	err := service.CreateCoolApp(input.Input)
	return vo.CommonRespVo{
		Err: vo.NewErr(err),
	}
}

func (PostGreeter) GetCoolAppInfo(input orgvo.GetCoolAppInfoReq) orgvo.GetCoolAppInfoResp {
	info, err := service.GetCoolAppInfo(input.Input.OpenConversationId)
	return orgvo.GetCoolAppInfoResp{
		Err:  vo.NewErr(err),
		Data: info,
	}
}

func (PostGreeter) DeleteCoolApp(input orgvo.DeleteCoolAppReq) vo.CommonRespVo {
	err := service.DeleteCoolApp(input)
	return vo.CommonRespVo{
		Err: vo.NewErr(err),
	}
}

func (PostGreeter) DeleteCoolAppByProject(input orgvo.DeleteCoolAppByProjectReq) vo.CommonRespVo {
	err := service.DeleteCoolAppByProject(input.OrgId, input.ProjectId)
	return vo.CommonRespVo{
		Err: vo.NewErr(err),
	}
}

func (PostGreeter) BindCoolApp(input orgvo.BindCoolAppReq) vo.CommonRespVo {
	err := service.BindCoolApp(input)
	return vo.CommonRespVo{
		Err: vo.NewErr(err),
	}
}

func (PostGreeter) UpdateCoolAppTopCard(input orgvo.UpdateCoolAppTopCardReq) vo.CommonRespVo {
	err := service.UpdateTopCard(input.OrgId, input.ProjectId)
	return vo.CommonRespVo{
		Err: vo.NewErr(err),
	}
}

func (PostGreeter) GetUpdateCoolAppTopCardData(input orgvo.GetCoolAppTopCardDataReq) orgvo.GetCoolAppTopCardDataResp {
	data, err := service.GetTopCardData(input.Input.OpenConversationId)
	return orgvo.GetCoolAppTopCardDataResp{
		Err:  vo.NewErr(err),
		Data: data,
	}
}

func (PostGreeter) GetSpaceList(input orgvo.GetSpaceListReq) orgvo.GetSpaceListResp {
	data, err := service.GetSpaceList(input)
	return orgvo.GetSpaceListResp{
		Err:  vo.NewErr(err),
		Data: data,
	}
}
