package trendssvc

import (
	"github.com/star-table/startable-server/app/service"
	"github.com/star-table/startable-server/common/model/vo"
	"github.com/star-table/startable-server/common/model/vo/trendsvo"
)

func (PostGreeter) UnreadNoticeCount(req trendsvo.UnreadNoticeCountReqVo) trendsvo.UnreadNoticeCountRespVo {
	res, err := service.UnreadNoticeCount(req.OrgId, req.UserId)
	return trendsvo.UnreadNoticeCountRespVo{Err: vo.NewErr(err), Count: int64(res)}
}

func (PostGreeter) CallMeCount(req trendsvo.CallMeCountReqVo) trendsvo.UnreadNoticeCountRespVo {
	res, err := service.GetCallMeCount(req)
	return trendsvo.UnreadNoticeCountRespVo{Err: vo.NewErr(err), Count: int64(res)}
}

func (PostGreeter) NoticeList(req trendsvo.NoticeListReqVo) trendsvo.NoticeListRespVo {
	res, err := service.NoticeList(req.OrgId, req.UserId, req.Page, req.Size, req.Input)
	return trendsvo.NoticeListRespVo{Err: vo.NewErr(err), Data: res}
}
