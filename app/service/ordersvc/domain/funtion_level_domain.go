package ordersvc

import (
	"github.com/star-table/startable-server/app/facade/orgfacade"
	"github.com/star-table/startable-server/common/core/errs"
	"github.com/star-table/startable-server/common/model/vo/orgvo"
)

func SetOrgFunction(orgId int64, level int64, sourceChannel string) errs.SystemErrorInfo {
	resp := orgfacade.UpdateOrgFunctionConfig(orgvo.UpdateOrgFunctionConfigReq{
		OrgId:         orgId,
		UserId:        0,
		SourceChannel: sourceChannel,
		Input: orgvo.UpdateFunctionConfigData{
			Level: level,
		},
	})
	if resp.Failure() {
		log.Error(resp.Error())
		return resp.Error()
	}

	return nil
}
