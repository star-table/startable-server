package commonsvc

import (
	"github.com/star-table/startable-server/common/core/logger"
	"github.com/star-table/startable-server/common/extra/gin/mvc"
)

var log = logger.GetDefaultLogger()

type PostGreeter struct {
	mvc.Greeter
}

type GetGreeter struct {
	mvc.Greeter
}

func (GetGreeter) Health() string {
	return "ok"
}
