package service

import (
	"github.com/star-table/startable-server/app/service/projectsvc/domain"
	"github.com/star-table/startable-server/common/core/errs"
)

func AuthProjectPermission(orgId, userId, projectId int64, path string, operation string, authFiling bool) errs.SystemErrorInfo {
	return domain.AuthProjectWithCond(orgId, userId, projectId, path, operation, authFiling, false, 0)
}
