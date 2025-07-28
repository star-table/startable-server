package orgsvc

import (
	"github.com/star-table/startable-server/common/core/consts"
	"github.com/star-table/startable-server/common/core/errs"
	"github.com/star-table/startable-server/common/core/util/asyn"
	"github.com/star-table/startable-server/common/library/cache"
	"github.com/star-table/startable-server/common/library/db/mysql"
	"github.com/star-table/startable-server/common/model/bo"
	"upper.io/db.v3/lib/sqlbuilder"
)

func PlatformAuth(sourceChannel, corpId, outUserId string, accessToken string, deptIds ...string) (*bo.BaseUserInfoBo, errs.SystemErrorInfo) {
	//获取组织信息
	orgInfo, err := GetOrgInfoByOutOrgId(corpId, sourceChannel)
	if err != nil {
		log.Error(err)
		if err.Code() == errs.OrgOutInfoNotExist.Code() {
			return nil, errs.PlatFormAppUnauthorizedError
		}
		return nil, errs.BuildSystemErrorInfo(errs.OrgNotInitError)
	}
	//获取用户信息
	baseUserInfo, err := GetPlatformBaseUserInfoByEmpId(orgInfo.OrgId, sourceChannel, outUserId, accessToken, corpId, true)
	if err != nil {
		//这里做用户初始化的兜底
		lockKey := consts.InitUserLock + sourceChannel + outUserId
		suc, err := cache.TryGetDistributedLock(lockKey, outUserId)
		log.Infof("准备获取分布式锁 %v", suc)
		if err != nil {
			log.Error(err)
			return nil, errs.BuildSystemErrorInfo(errs.RedisOperateError)
		}
		// 如果是初始化用户，则对此标记置为 true。用于判断是否是新加的用户。
		isNewFlag := false
		if suc {
			log.Infof("获取分布式锁成功 %v", suc)
			defer func() {
				if _, lockErr := cache.ReleaseDistributedLock(lockKey, outUserId); lockErr != nil {
					log.Error(lockErr)
				}
			}()
			//double check
			baseUserInfo, err = GetPlatformBaseUserInfoByEmpId(orgInfo.OrgId, sourceChannel, outUserId, accessToken, corpId, true)
			if err != nil {
				err1 := mysql.TransX(func(tx sqlbuilder.Tx) error {
					_, err := InitPlatformUser(sourceChannel, orgInfo.OrgId, corpId, outUserId, tx, accessToken, deptIds...)
					if err != nil {
						log.Error(err)
						return err
					}
					isNewFlag = true
					return nil
				})
				if err1 != nil {
					log.Error(err1)
					if sysErr, ok := err1.(errs.SystemErrorInfo); ok {
						return nil, sysErr
					}
					return nil, errs.BuildSystemErrorInfo(errs.FeiShuUserNotInAppUseScopeOfAuthority, err1)
				}
			}
		}
		if baseUserInfo == nil {
			//上面已经初始化了，不需要再从飞书获取新的信息更新了，所以传false
			baseUserInfo, err = GetPlatformBaseUserInfoByEmpId(orgInfo.OrgId, sourceChannel, outUserId, accessToken, corpId, false)
			if err != nil {
				log.Error(err)
				return nil, errs.UserInitError
			}
			// 新增用户时，给前端推送事件
			if isNewFlag && baseUserInfo != nil {
				userDeptIdMap, oriErr := GetUserDepartmentIdMap(baseUserInfo.OrgId, []int64{baseUserInfo.UserId})
				if oriErr != nil {
					log.Error(oriErr)
					return nil, errs.BuildSystemErrorInfo(errs.MysqlOperateError, oriErr)
				}
				currentUserDeptId := int64(0)
				if tmpDeptId, ok := userDeptIdMap[baseUserInfo.UserId]; ok {
					currentUserDeptId = tmpDeptId
				}
				asyn.Execute(func() {
					PushAddOrgMemberNotice(baseUserInfo.OrgId, currentUserDeptId, []int64{baseUserInfo.UserId}, 0)
				})
			}
		}
	}
	return baseUserInfo, nil
}
