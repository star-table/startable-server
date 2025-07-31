package orgsvc

import (
	"gitea.bjx.cloud/allstar/feishu-sdk-golang/sdk"
	sdk_const "gitea.bjx.cloud/allstar/platform-sdk/consts"
	"github.com/pkg/errors"
	"github.com/star-table/startable-server/app/facade/idfacade"
	sconsts "github.com/star-table/startable-server/app/service"
	"github.com/star-table/startable-server/app/service/orgsvc/po"
	"github.com/star-table/startable-server/common/core/consts"
	"github.com/star-table/startable-server/common/core/errs"
	"github.com/star-table/startable-server/common/core/util"
	"github.com/star-table/startable-server/common/core/util/json"
	"github.com/star-table/startable-server/common/core/util/pinyin"
	"github.com/star-table/startable-server/common/core/util/slice"
	"github.com/star-table/startable-server/common/extra/feishu"
	"github.com/star-table/startable-server/common/library/cache"
	"github.com/star-table/startable-server/common/library/db/mysql"
	"github.com/star-table/startable-server/common/model/bo"
	"upper.io/db.v3"
	"upper.io/db.v3/lib/sqlbuilder"
)

func GetUserIdBatchByEmpId(sourceChannel string, orgId int64, empIds []string) ([]int64, errs.SystemErrorInfo) {
	keys := make([]interface{}, len(empIds))
	for i, empId := range empIds {
		key, _ := util.ParseCacheKey(sconsts.CacheOutUserIdRelationId, map[string]interface{}{
			consts.CacheKeyOrgIdConstName:         orgId,
			consts.CacheKeySourceChannelConstName: sourceChannel,
			consts.CacheKeyOutUserIdConstName:     empId,
		})
		keys[i] = key
	}
	resultList := make([]string, 0)
	if len(keys) > 0 {
		list, err := cache.MGet(keys...)
		if err != nil {
			log.Error(err)
			return nil, errs.BuildSystemErrorInfo(errs.RedisOperateError)
		}
		resultList = list
	}
	userIds := make([]int64, 0)
	validEmpIds := make([]string, 0)
	for _, empInfoJson := range resultList {
		empIdInfo := &bo.UserEmpIdInfo{}
		err := json.FromJson(empInfoJson, empIdInfo)
		if err != nil {
			log.Error(err)
			return nil, errs.BuildSystemErrorInfo(errs.JSONConvertError)
		}
		userIds = append(userIds, empIdInfo.UserId)
		validEmpIds = append(validEmpIds, empIdInfo.EmpId)
	}
	//找不存在的
	if len(empIds) != len(validEmpIds) {
		for _, empId := range empIds {
			exist, _ := slice.Contain(validEmpIds, empId)
			if !exist {
				userId, err := GetUserIdByEmpId(orgId, empId)
				if err != nil {
					log.Error(err)
					continue
				}
				userIds = append(userIds, userId)
			}
		}
	}
	return userIds, nil
}

func GetUserIdByEmpId(orgId int64, empId string) (int64, errs.SystemErrorInfo) {
	key, err5 := util.ParseCacheKey(sconsts.CacheOutUserIdRelationId, map[string]interface{}{
		consts.CacheKeyOrgIdConstName:     orgId,
		consts.CacheKeyOutUserIdConstName: empId,
	})
	if err5 != nil {
		log.Error(err5)
		return 0, err5
	}

	empInfoJson, err := cache.Get(key)
	if err != nil {
		log.Error(err)
		return 0, errs.BuildSystemErrorInfo(errs.RedisOperateError, err)
	}
	if empInfoJson != "" {
		empIdInfo := &bo.UserEmpIdInfo{}
		err := json.FromJson(empInfoJson, empIdInfo)
		if err != nil {
			log.Error(err)
			return 0, errs.BuildSystemErrorInfo(errs.JSONConvertError)
		}
		return empIdInfo.UserId, nil
	} else {
		userOutInfo := &po.PpmOrgUserOutInfo{}
		cond := db.Cond{
			consts.TcOrgId:     orgId,
			consts.TcOutUserId: empId,
			consts.TcIsDelete:  consts.AppIsNoDelete,
			consts.TcStatus:    consts.AppStatusEnable,
		}
		err := mysql.SelectOneByCond(userOutInfo.TableName(), cond, userOutInfo)
		if err != nil {
			log.Error(err)
			return 0, errs.BuildSystemErrorInfo(errs.UserNotExist, errors.New(" empId:"+empId))
		}
		empIdInfo := bo.UserEmpIdInfo{
			EmpId:  empId,
			UserId: userOutInfo.UserId,
		}
		err = cache.Set(key, json.ToJsonIgnoreError(empIdInfo))
		if err != nil {
			log.Error(err)
			return 0, errs.BuildSystemErrorInfo(errs.RedisOperateError, err)
		}
		return userOutInfo.UserId, nil
	}
}

func GetDingTalkBaseUserInfoByEmpId(orgId int64, empId string) (*bo.BaseUserInfoBo, errs.SystemErrorInfo) {
	userId, err := GetUserIdByEmpId(orgId, empId)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	return GetDingTalkBaseUserInfo(orgId, userId)
}

func GetBaseUserInfoByEmpId(orgId int64, empId string) (*bo.BaseUserInfoBo, errs.SystemErrorInfo) {
	userId, err := GetUserIdByEmpId(orgId, empId)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	return GetBaseUserInfo(orgId, userId)
}

func GetPlatformBaseUserInfoByEmpId(orgId int64, sourceChannel, empId, accessToken, corpId string, isNeedUpdate bool) (*bo.BaseUserInfoBo, errs.SystemErrorInfo) {
	userId, err := GetUserIdByEmpId(orgId, empId)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	info, infoErr := GetPlatformBaseUserInfo(orgId, userId)
	if infoErr != nil {
		log.Error(infoErr)
		return nil, infoErr
	}
	if !isNeedUpdate {
		return info, nil
	}

	// 飞书做了很多操作。。不知道整的啥，先抽一个函数出去，目前钉钉应该用不到，而且感觉有bug，可能报错
	if sourceChannel == sdk_const.SourceChannelFeishu {
		info, err = completionFeiShuInfo(info, corpId, accessToken)
	}

	return info, err
}

func completionFeiShuInfo(baseInfo *bo.BaseUserInfoBo, corpId, accessToken string) (*bo.BaseUserInfoBo, errs.SystemErrorInfo) {
	//判断信息是否不全，不全的话补全信息
	if baseInfo.Name == "未激活" && baseInfo.Avatar == consts.AvatarForUnallocated {
		isInScope, isInScopeErr := feishu.JudgeIsInAppScopes(corpId)
		if isInScopeErr != nil {
			log.Error(isInScopeErr)
			return nil, isInScopeErr
		}
		if isInScope {
			//有通讯录权限，就获取通讯录权限，补全部门信息。同时更新用户基本信息
			tenant, err1 := feishu.GetTenant(corpId)
			if err1 != nil {
				log.Error(err1)
				return nil, err1
			}
			userBatchGetResp, err := tenant.GetUserBatchGetV2(nil, []string{baseInfo.OutOrgId})
			if err != nil {
				log.Error(err)
				return nil, errs.BuildSystemErrorInfo(errs.FeiShuOpenApiCallError, err)
			}
			if userBatchGetResp.Code != 0 {
				log.Errorf("err %s", json.ToJsonIgnoreError(userBatchGetResp))
				return nil, errs.BuildSystemErrorInfoWithMessage(errs.FeiShuOpenApiCallError, userBatchGetResp.Msg)
			}
			userInfos := userBatchGetResp.Data.Users
			if len(userInfos) == 0 {
				log.Errorf("user不存在")
				return nil, errs.FeiShuUserNotInAppUseScopeOfAuthority
			}
			fsUser := userInfos[0]

			insertUserDeptPos := []po.PpmOrgUserDepartment{}
			if fsUser.Departments != nil && len(fsUser.Departments) > 0 {
				//查找部门
				depts := &[]po.PpmOrgDepartmentOutInfo{}
				conn, err := mysql.GetConnect()
				if err != nil {
					log.Error(err)
					return nil, errs.BuildSystemErrorInfo(errs.MysqlOperateError)
				}
				selectErr := conn.Select("o.department_id").From("ppm_org_department_out_info o", "ppm_org_department d").Where(db.Cond{
					"o." + consts.TcIsDelete:     consts.AppIsNoDelete,
					"d." + consts.TcIsDelete:     consts.AppIsNoDelete,
					"o." + consts.TcOrgId:        baseInfo.OrgId,
					"d." + consts.TcOrgId:        baseInfo.OrgId,
					"o." + consts.TcStatus:       consts.AppStatusEnable,
					"d." + consts.TcStatus:       consts.AppStatusEnable,
					"o." + consts.TcDepartmentId: db.Raw("d." + consts.TcId),
				}).All(depts)
				if selectErr != nil {
					log.Error(err1)
					return nil, errs.MysqlOperateError
				}

				if len(*depts) > 0 {
					userDepIds, idErr := idfacade.ApplyMultiplePrimaryIdRelaxed(consts.TableUserDepartment, len(*depts))
					if idErr != nil {
						log.Error(idErr)
						return nil, idErr
					}
					for i, info := range *depts {
						insertUserDeptPos = append(insertUserDeptPos, po.PpmOrgUserDepartment{
							Id:           userDepIds.Ids[i].Id,
							OrgId:        info.OrgId,
							UserId:       baseInfo.UserId,
							DepartmentId: info.DepartmentId,
						})
					}
				}
			}

			//更新用户信息
			namePy := pinyin.ConvertToPinyin(fsUser.Name)
			transErr := mysql.TransX(func(tx sqlbuilder.Tx) error {
				if fsUser.Name != "" {
					_, err := mysql.TransUpdateSmartWithCond(tx, consts.TableUserOutInfo, db.Cond{
						consts.TcIsDelete: consts.AppIsNoDelete,
						consts.TcUserId:   baseInfo.UserId,
						consts.TcOrgId:    baseInfo.OrgId,
					}, mysql.Upd{
						consts.TcName:   fsUser.Name,
						consts.TcAvatar: fsUser.Avatar.AvatarOrigin,
					})
					if err != nil {
						log.Error(err)
						return err
					}

					_, err1 := mysql.TransUpdateSmartWithCond(tx, consts.TableUser, db.Cond{
						consts.TcIsDelete: consts.AppIsNoDelete,
						consts.TcId:       baseInfo.UserId,
					}, mysql.Upd{
						consts.TcName:       fsUser.Name,
						consts.TcAvatar:     fsUser.Avatar.AvatarOrigin,
						consts.TcNamePinyin: namePy,
					})
					if err1 != nil {
						log.Error(err1)
						return err1
					}
				}

				if len(insertUserDeptPos) > 0 {
					dbErr := mysql.TransBatchInsert(tx, &po.PpmOrgUserDepartment{}, slice.ToSlice(insertUserDeptPos))
					if dbErr != nil {
						log.Error(dbErr)
						return dbErr
					}
				}

				return nil
			})
			if transErr != nil {
				log.Error(transErr)
				return nil, errs.MysqlOperateError
			}

			baseInfo.Name = fsUser.Name
			baseInfo.NamePy = namePy
			baseInfo.Avatar = fsUser.Avatar.AvatarOrigin
		} else {
			//没有通讯录权限，就补全用户基本信息
			userInfo, err := sdk.GetOAuth2UserInfo(accessToken)
			if err != nil {
				log.Error(err)
				return nil, errs.FeiShuOpenApiCallError
			}

			if userInfo.Name != "" {
				//更新用户信息
				namePy := pinyin.ConvertToPinyin(userInfo.Name)
				transErr := mysql.TransX(func(tx sqlbuilder.Tx) error {
					_, err := mysql.TransUpdateSmartWithCond(tx, consts.TableUserOutInfo, db.Cond{
						consts.TcIsDelete: consts.AppIsNoDelete,
						consts.TcUserId:   baseInfo.UserId,
						consts.TcOrgId:    baseInfo.OrgId,
					}, mysql.Upd{
						consts.TcName:   userInfo.Name,
						consts.TcAvatar: userInfo.AvatarUrl,
					})
					if err != nil {
						log.Error(err)
						return err
					}

					_, err1 := mysql.TransUpdateSmartWithCond(tx, consts.TableUser, db.Cond{
						consts.TcIsDelete: consts.AppIsNoDelete,
						consts.TcId:       baseInfo.UserId,
					}, mysql.Upd{
						consts.TcName:       userInfo.Name,
						consts.TcAvatar:     userInfo.AvatarUrl,
						consts.TcNamePinyin: namePy,
					})
					if err1 != nil {
						log.Error(err1)
						return err1
					}

					return nil
				})
				if transErr != nil {
					log.Error(transErr)
					return nil, errs.MysqlOperateError
				}

				baseInfo.Name = userInfo.Name
				baseInfo.NamePy = namePy
				baseInfo.Avatar = userInfo.AvatarUrl
			}
		}

		clearErr := ClearBaseUserInfo(baseInfo.OrgId, baseInfo.UserId)
		if clearErr != nil {
			log.Error(clearErr)
			return nil, clearErr
		}
	}

	return baseInfo, nil
}

func GetDingTalkBaseUserInfo(orgId int64, userId int64) (*bo.BaseUserInfoBo, errs.SystemErrorInfo) {
	return GetBaseUserInfo(orgId, userId)
}

func GetPlatformBaseUserInfo(orgId int64, userId int64) (*bo.BaseUserInfoBo, errs.SystemErrorInfo) {
	return GetBaseUserInfo(orgId, userId)
}
