package n8nfacade

import (
	"fmt"

	"github.com/star-table/startable-server/common/core/errs"

	"github.com/star-table/startable-server/app/facade"
	"github.com/star-table/startable-server/common/core/config"
	"github.com/star-table/startable-server/common/core/consts"
	"github.com/star-table/startable-server/common/core/logger"
)

var log = logger.GetDefaultLogger()

func CallWebhookWaiting(executionId int64, body interface{}) errs.SystemErrorInfo {
	reqUrl := fmt.Sprintf("%s/webhook-waiting/%d", config.GetPreUrl("n8nwebhook"), executionId)
	err := facade.Request(consts.HttpMethodPost, reqUrl, nil, nil, body, nil)
	if err.Failure() {
		return err.Error()
	}
	return nil
}
