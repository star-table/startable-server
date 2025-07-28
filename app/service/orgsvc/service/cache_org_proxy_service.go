package orgsvc

import (
	"github.com/star-table/startable-server/common/core/errs"
	"github.com/star-table/startable-server/common/model/bo"
	"github.com/star-table/startable-server/common/model/vo/orgvo"
)

func GetBaseOrgInfo(orgId int64) (*bo.BaseOrgInfoBo, errs.SystemErrorInfo) {
	return domain.GetBaseOrgInfo(orgId)
}

func GetBaseUserInfoByEmpId(orgId int64, empId string) (*bo.BaseUserInfoBo, errs.SystemErrorInfo) {
	return domain.GetBaseUserInfoByEmpId(orgId, empId)
}

func GetBaseUserInfoByEmpIdBatch(orgId int64, input orgvo.GetBaseUserInfoByEmpIdBatchReqVoInput) ([]bo.BaseUserInfoBo, errs.SystemErrorInfo) {
	return domain.GetBaseUserInfoByOpenIdBatch(orgId, input.OpenIds)
}

func GetUserConfigInfo(orgId int64, userId int64) (*bo.UserConfigBo, errs.SystemErrorInfo) {
	return domain.GetUserConfigInfo(orgId, userId)
}

func GetUserConfigInfoBatch(orgId int64, input *orgvo.GetUserConfigInfoBatchReqVoInput) ([]bo.UserConfigBo, errs.SystemErrorInfo) {
	return domain.GetUserConfigInfoBatch(orgId, input.UserIds)
}

func GetBaseUserInfo(orgId int64, userId int64) (*bo.BaseUserInfoBo, errs.SystemErrorInfo) {
	return domain.GetBaseUserInfo(orgId, userId)
}

func GetDingTalkBaseUserInfo(orgId int64, userId int64) (*bo.BaseUserInfoBo, errs.SystemErrorInfo) {
	return domain.GetDingTalkBaseUserInfo(orgId, userId)
}

func GetBaseUserInfoBatch(orgId int64, userIds []int64) ([]bo.BaseUserInfoBo, errs.SystemErrorInfo) {
	return domain.GetBaseUserInfoBatch(orgId, userIds)
}

func SetShareUrl(key, url string) errs.SystemErrorInfo {
	return domain.SetShareUrl(key, url)
}

func GetShareUrl(key string) (string, errs.SystemErrorInfo) {
	return domain.GetShareUrl(key)
}

func ClearOrgUsersPayCache(orgId int64) errs.SystemErrorInfo {
	return domain.ClearAllOrgUserPayCache(orgId)
}
