package api

import (
	msgPb "gitea.bjx.cloud/LessCode/interface/golang/msg/v1"
	"github.com/star-table/startable-server/app/facade/common/report"
	"github.com/star-table/startable-server/common/model/vo"
	"github.com/star-table/startable-server/common/model/vo/commonvo"
)

func (PostGreeter) ReportAppEvent(req *commonvo.ReportAppEventReq) *vo.CommonRespVo {
	report.ReportAppEvent(msgPb.EventType(req.EventType), req.TraceId, req.AppEvent)
	return &vo.CommonRespVo{}
}

func (PostGreeter) ReportTableEvent(req *commonvo.ReportTableEventReq) *vo.CommonRespVo {
	report.ReportTableEvent(msgPb.EventType(req.EventType), req.TraceId, req.TableEvent)
	return &vo.CommonRespVo{}
}
