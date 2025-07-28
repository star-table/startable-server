package callsvc

import (
	service "github.com/star-table/startable-server/app/service"
	"github.com/star-table/startable-server/common/core/errs"
)

func AuthOrgRole(orgId, userId int64, path string, operation string) errs.SystemErrorInfo {
	return service.Authenticate(orgId, userId, nil, nil, path, operation, nil)
}
