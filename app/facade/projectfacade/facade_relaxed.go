package projectfacade

import (
	"gitea.bjx.cloud/allstar/platform-sdk/sdk_interface"
	"github.com/star-table/startable-server/app/facade/orgfacade"
	"github.com/star-table/startable-server/common/core/errs"
	"github.com/star-table/startable-server/common/model/bo"
	"github.com/star-table/startable-server/common/model/vo/orgvo"
	"github.com/spf13/cast"
)

func GetBaseUserInfoBatchRelaxed(orgId int64, userIds []int64) ([]bo.BaseUserInfoBo, errs.SystemErrorInfo) {
	respVo := orgfacade.GetBaseUserInfoBatch(orgvo.GetBaseUserInfoBatchReqVo{
		OrgId:   orgId,
		UserIds: userIds,
	})
	return respVo.BaseUserInfos, respVo.Error()
}

func GetCorpInfoFromDB(orgId int64, corpId string) (*sdk_interface.CorpInfo, error) {
	respVo := orgfacade.GetOrgOutInfoByOutOrgId(orgvo.GetOutOrgInfoByOutOrgIdReqVo{
		OrgId:    orgId,
		OutOrgId: corpId,
	})
	if respVo.Failure() {
		return nil, respVo.Error()
	}

	return &sdk_interface.CorpInfo{
		OrgId:         respVo.Data.OrgId,
		AgentId:       cast.ToInt64(respVo.Data.TenantCode),
		CorpId:        respVo.Data.OutOrgId,
		PermanentCode: respVo.Data.AuthTicket,
	}, respVo.Error()
}
