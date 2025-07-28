package asyn

import (
	"github.com/star-table/startable-server/common/core/consts"
	sconsts "github.com/star-table/startable-server/common/core/consts"
	"github.com/star-table/startable-server/common/core/errs"
	"github.com/star-table/startable-server/common/core/logger"
	"github.com/star-table/startable-server/common/core/threadlocal"
	"github.com/star-table/startable-server/common/core/util/stack"
	"github.com/jtolds/gls"
)

var log = logger.GetDefaultLogger()

func Execute(fn func()) {
	httpContextKey, _ := threadlocal.Mgr.GetValue(consts.HttpContextKey)
	traceId, _ := threadlocal.Mgr.GetValue(consts.TraceIdKey)
	openTraceId, _ := threadlocal.Mgr.GetValue(sconsts.JaegerContextTraceKey)
	pmLang, _ := threadlocal.Mgr.GetValue(sconsts.AppHeaderLanguage)
	go func() {
		threadlocal.Mgr.SetValues(gls.Values{
			consts.HttpContextKey:         httpContextKey,
			consts.TraceIdKey:             traceId,
			sconsts.JaegerContextTraceKey: openTraceId,
			sconsts.AppHeaderLanguage:     pmLang,
		}, func() {
			defer func() {
				if r := recover(); r != nil {
					log.Error(errs.BuildSystemErrorInfoWithPanicRecover(r, stack.GetStack()))
				}
			}()
			fn()
		})
	}()
}
