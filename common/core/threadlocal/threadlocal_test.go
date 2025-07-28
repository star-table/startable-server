package threadlocal

import (
	"testing"

	"github.com/star-table/startable-server/common/core/consts"
	"github.com/star-table/startable-server/common/core/util/uuid"
	"github.com/jtolds/gls"
)

func TestSetTraceId(t *testing.T) {

	SetTraceId()
	t.Log(GetTraceId())

	Mgr.SetValues(gls.Values{consts.TraceIdKey: uuid.NewUuid()}, func() {

		t.Log("in ", GetTraceId())
	})
}
