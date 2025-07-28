package callsvc

import (
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/star-table/startable-server/app/facade/orgfacade"
	"github.com/star-table/startable-server/common/core/consts"
	"github.com/star-table/startable-server/common/core/model"
	"github.com/star-table/startable-server/common/core/threadlocal"
	"github.com/star-table/startable-server/common/core/util/json"
	"github.com/star-table/startable-server/common/core/util/net"
	"github.com/star-table/startable-server/common/model/vo/orgvo"
	"github.com/gin-gonic/gin"
	"github.com/jtolds/gls"
)

func PAReport(c *gin.Context) {
	httpContext := c.Request.Context().Value(consts.HttpContextKey).(model.HttpContext)
	threadlocal.Mgr.SetValues(gls.Values{consts.HttpContextKey: httpContext, consts.TraceIdKey: httpContext.TraceId}, func() {
		PAReportHandler(c.Writer, c.Request)
	})
}

func PAReportHandler(w http.ResponseWriter, r *http.Request) {
	reqBody, _ := ioutil.ReadAll(r.Body)
	paReportMsgReqVo := orgvo.PAReportMsgReqVo{}
	_ = json.FromJson(string(reqBody), &paReportMsgReqVo)
	paReportMsgReqVo.Body.OutSideIp, _ = net.GetIP(r)
	paReportResp := orgfacade.PAReport(orgvo.PAReportMsgReqVo{})
	_, _ = fmt.Fprintln(w, paReportResp)
}
