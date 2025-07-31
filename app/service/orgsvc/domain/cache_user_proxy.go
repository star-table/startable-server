package orgsvc

import (
	"fmt"
	"math/rand"
	"strconv"
	"time"

	fsvo "gitea.bjx.cloud/allstar/feishu-sdk-golang/core/model/vo"
	"gitea.bjx.cloud/allstar/feishu-sdk-golang/sdk"
	platform_sdk "gitea.bjx.cloud/allstar/platform-sdk"
	sdk_const "gitea.bjx.cloud/allstar/platform-sdk/consts"
	sdkVo "gitea.bjx.cloud/allstar/platform-sdk/vo"
	sconsts "github.com/star-table/startable-server/app/service"
	"github.com/star-table/startable-server/app/service/orgsvc/po"
	"github.com/star-table/startable-server/common/core/consts"
	"github.com/star-table/startable-server/common/core/errs"
	"github.com/star-table/startable-server/common/core/util"
	"github.com/star-table/startable-server/common/core/util/copyer"
	"github.com/star-table/startable-server/common/core/util/json"
	"github.com/star-table/startable-server/common/core/util/maps"
	"github.com/star-table/startable-server/common/core/util/slice"
	"github.com/star-table/startable-server/common/core/util/stack"
	"github.com/star-table/startable-server/common/core/util/uuid"
	"github.com/star-table/startable-server/common/extra/feishu"
	"github.com/star-table/startable-server/common/language/english"
	"github.com/star-table/startable-server/common/library/cache"
	"github.com/star-table/startable-server/common/library/db/mysql"
	"github.com/star-table/startable-server/common/model/bo"
	"github.com/star-table/startable-server/common/model/vo/ordervo"
	"upper.io/db.v3"
)

func GetBaseUserOutInfoBatch(orgId int64, userIds []int64) ([]bo.BaseUserOutInfoBo, errs.SystemErrorInfo) {
	keys := make([]interface{}, len(userIds))
	for i, userId := range userIds {
		key, _ := util.ParseCacheKey(sconsts.CacheBaseUserOutInfo, map[string]interface{}{
			consts.CacheKeyOrgIdConstName:  orgId,
			consts.CacheKeyUserIdConstName: userId,
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
	baseUserOutInfoList := make([]bo.BaseUserOutInfoBo, 0)
	validUserIds := map[int64]bool{}
	for _, userInfoJson := range resultList {
		userOutInfoBo := &bo.BaseUserOutInfoBo{}
		err := json.FromJson(userInfoJson, userOutInfoBo)
		if err != nil {
			return nil, errs.BuildSystemErrorInfo(errs.JSONConvertError)
		}
		baseUserOutInfoList = append(baseUserOutInfoList, *userOutInfoBo)
		validUserIds[userOutInfoBo.UserId] = true
	}

	missUserIds := make([]int64, 0)
	//找不存在的
	if len(userIds) != len(validUserIds) {
		for _, userId := range userIds {
			if _, ok := validUserIds[userId]; !ok {
				missUserIds = append(missUserIds, userId)
			}
		}
	}

	//批量查外部信息
	outInfos, userErr := GetBaseUserOutInfoByUserIds(orgId, missUserIds)
	if userErr != nil {
		log.Error(userErr)
		return nil, userErr
	}

	if len(outInfos) > 0 {
		baseUserOutInfoList = append(baseUserOutInfoList, outInfos...)
	}

	return baseUserOutInfoList, nil
}

func GetBaseUserInfoBatch(orgId int64, originUserIds []int64) ([]bo.BaseUserInfoBo, errs.SystemErrorInfo) {
	//去重
	userIds := slice.SliceUniqueInt64(originUserIds)

	keys := make([]interface{}, len(userIds))
	for i, userId := range userIds {
		key, _ := util.ParseCacheKey(sconsts.CacheBaseUserInfo, map[string]interface{}{
			consts.CacheKeyOrgIdConstName:  orgId,
			consts.CacheKeyUserIdConstName: userId,
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
	baseUserInfoList := make([]bo.BaseUserInfoBo, 0)
	validUserIds := map[int64]bool{}
	for _, userInfoJson := range resultList {
		userInfoBo := &bo.BaseUserInfoBo{}
		err := json.FromJson(userInfoJson, userInfoBo)
		if err != nil {
			return nil, errs.BuildSystemErrorInfo(errs.JSONConvertError)
		}
		baseUserInfoList = append(baseUserInfoList, *userInfoBo)
		validUserIds[userInfoBo.UserId] = true
	}

	log.Infof("from cache %s", json.ToJsonIgnoreError(baseUserInfoList))
	missUserIds := make([]int64, 0)
	//找不存在的
	if len(userIds) != len(validUserIds) {
		for _, userId := range userIds {
			if _, ok := validUserIds[userId]; !ok {
				missUserIds = append(missUserIds, userId)
			}
		}
	}

	missUserInfos, userErr := getLocalBaseUserInfoBatch(orgId, missUserIds)
	if userErr != nil {
		log.Error(userErr)
		return nil, userErr
	}
	if len(missUserInfos) > 0 {
		baseUserInfoList = append(baseUserInfoList, missUserInfos...)
	}
	if ok, _ := slice.Contain(originUserIds, int64(0)); ok {
		baseUserInfoList = append(baseUserInfoList, bo.BaseUserInfoBo{
			UserId:             0,
			OutUserId:          "",
			OrgId:              orgId,
			OutOrgId:           "",
			Name:               english.WordTransLate("未分配"),
			NamePy:             "",
			Avatar:             consts.AvatarForUnallocated,
			HasOutInfo:         false,
			HasOrgOutInfo:      false,
			OutOrgUserId:       "",
			OrgUserIsDelete:    0,
			OrgUserStatus:      0,
			OrgUserCheckStatus: 0,
		})
	}

	//获取用户外部信息
	baseUserOutInfos, err := GetBaseUserOutInfoBatch(orgId, userIds)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	outInfoMap := maps.NewMap("UserId", baseUserOutInfos)
	for i, _ := range baseUserInfoList {
		userInfo := baseUserInfoList[i]
		if outInfoInterface, ok := outInfoMap[userInfo.UserId]; ok {
			outInfo := outInfoInterface.(bo.BaseUserOutInfoBo)

			userInfo.OutUserId = outInfo.OutUserId
			userInfo.OutOrgUserId = outInfo.OutOrgUserId
			userInfo.OutOrgId = outInfo.OutOrgId
			userInfo.HasOutInfo = outInfo.OutUserId != ""
			userInfo.HasOrgOutInfo = outInfo.OutOrgId != ""
		}
		baseUserInfoList[i] = userInfo
	}

	//按照原始顺序排序
	baseUserInfoMap := map[int64]bo.BaseUserInfoBo{}
	for _, infoBo := range baseUserInfoList {
		baseUserInfoMap[infoBo.UserId] = infoBo
	}

	resBo := []bo.BaseUserInfoBo{}
	for _, id := range originUserIds {
		if info, ok := baseUserInfoMap[id]; ok {
			resBo = append(resBo, info)
		}
	}
	log.Infof("user info: %s", json.ToJsonIgnoreError(resBo))

	return resBo, nil
}

func ClearBaseUserInfo(orgId, userId int64) errs.SystemErrorInfo {
	key, err5 := util.ParseCacheKey(sconsts.CacheBaseUserInfo, map[string]interface{}{
		consts.CacheKeyOrgIdConstName:  orgId,
		consts.CacheKeyUserIdConstName: userId,
	})
	if err5 != nil {
		log.Error(err5)
		return err5
	}
	_, err := cache.Del(key)
	if err != nil {
		log.Error(err)
		return errs.BuildSystemErrorInfo(errs.RedisOperateError)
	}
	return nil
}

// 批量清楚用户缓存信息
func ClearBaseUserInfoBatch(orgId int64, userIds []int64) errs.SystemErrorInfo {
	keys := make([]interface{}, 0)
	for _, userId := range userIds {
		key, err5 := util.ParseCacheKey(sconsts.CacheBaseUserInfo, map[string]interface{}{
			consts.CacheKeyOrgIdConstName:  orgId,
			consts.CacheKeyUserIdConstName: userId,
		})
		if err5 != nil {
			log.Error(err5)
			return err5
		}
		keys = append(keys, key)
	}
	_, err := cache.Del(keys...)
	if err != nil {
		log.Error(err)
		return errs.BuildSystemErrorInfo(errs.RedisOperateError)
	}
	return nil
}

// sourceChannel可以为空
func GetBaseUserInfo(orgId int64, userId int64) (*bo.BaseUserInfoBo, errs.SystemErrorInfo) {
	if userId == 0 {
		//系统创建
		return &bo.BaseUserInfoBo{
			OrgId:  orgId,
			Name:   english.WordTransLate("未分配"),
			Avatar: consts.AvatarForUnallocated,
		}, nil
	}

	key, err5 := util.ParseCacheKey(sconsts.CacheBaseUserInfo, map[string]interface{}{
		consts.CacheKeyOrgIdConstName:  orgId,
		consts.CacheKeyUserIdConstName: userId,
	})
	if err5 != nil {
		log.Error(err5)
		return nil, err5
	}

	baseUserInfoJson, err := cache.Get(key)
	if err != nil {
		return nil, errs.BuildSystemErrorInfo(errs.RedisOperateError, err)
	}
	baseUserInfo := &bo.BaseUserInfoBo{}
	if baseUserInfoJson != "" {
		err := json.FromJson(baseUserInfoJson, baseUserInfo)
		if err != nil {
			return nil, errs.BuildSystemErrorInfo(errs.JSONConvertError, err)
		}
	} else {
		userInfo, errorInfo := getLocalBaseUserInfo(orgId, userId, key)
		if errorInfo != nil {
			log.Error(errorInfo)
			return nil, errorInfo
		}
		baseUserInfo = userInfo
	}

	//这里不存缓存，动态获取
	baseUserOutInfo, sysErr := GetBaseUserOutInfo(orgId, userId)
	if sysErr != nil {
		log.Error(sysErr)
		return nil, sysErr
	}
	baseUserInfo.OutUserId = baseUserOutInfo.OutUserId
	baseUserInfo.OutOrgId = baseUserOutInfo.OutOrgId
	baseUserInfo.HasOutInfo = baseUserInfo.OutUserId != ""
	baseUserInfo.HasOrgOutInfo = baseUserInfo.OutOrgId != ""
	baseUserInfo.OutOrgUserId = baseUserOutInfo.OutOrgUserId

	return baseUserInfo, nil
}

func GetBaseUserOutInfo(orgId int64, userId int64) (*bo.BaseUserOutInfoBo, errs.SystemErrorInfo) {
	if userId == 0 {
		//系统创建
		return &bo.BaseUserOutInfoBo{
			OrgId: orgId,
		}, nil
	}

	key, err5 := util.ParseCacheKey(sconsts.CacheBaseUserOutInfo, map[string]interface{}{
		consts.CacheKeyOrgIdConstName:  orgId,
		consts.CacheKeyUserIdConstName: userId,
	})
	if err5 != nil {
		log.Error(err5)
		return nil, err5
	}
	baseUserOutInfoJson, err := cache.Get(key)
	if err != nil {
		return nil, errs.BuildSystemErrorInfo(errs.RedisOperateError, err)
	}
	if baseUserOutInfoJson != "" {
		baseUserOutInfo := &bo.BaseUserOutInfoBo{}
		err := json.FromJson(baseUserOutInfoJson, baseUserOutInfo)
		if err != nil {
			return nil, errs.BuildSystemErrorInfo(errs.JSONConvertError, err)
		}
		return baseUserOutInfo, nil
	} else {
		//用户外部信息
		userOutInfo := &po.PpmOrgUserOutInfo{}
		_ = mysql.SelectOneByCond(consts.TableUserOutInfo, db.Cond{
			consts.TcIsDelete: consts.AppIsNoDelete,
			consts.TcOrgId:    orgId,
			consts.TcUserId:   userId,
		}, userOutInfo)

		outInfo := bo.BaseUserOutInfoBo{
			UserId:       userId,
			OrgId:        orgId,
			OutUserId:    userOutInfo.OutUserId,
			OutOrgUserId: userOutInfo.OutOrgUserId,
		}
		//组织外部信息
		orgOutInfo := &po.PpmOrgOrganizationOutInfo{}
		err = mysql.SelectOneByCond(consts.TableOrganizationOutInfo, db.Cond{
			consts.TcIsDelete: consts.AppIsNoDelete,
			consts.TcOrgId:    orgId,
		}, orgOutInfo)
		if err != nil {
			if err == db.ErrNoMoreRows {

			} else {
				log.Error(err)
				return nil, errs.MysqlOperateError
			}
		} else {
			outInfo.OutOrgId = orgOutInfo.OutOrgId
		}

		baseUserOutInfoJson, err := json.ToJson(outInfo)
		if err != nil {
			log.Error(err)
			return nil, errs.BuildSystemErrorInfo(errs.JSONConvertError, err)
		}
		err = cache.SetEx(key, baseUserOutInfoJson, consts.GetCacheBaseExpire())
		if err != nil {
			log.Error(err)
			return nil, errs.BuildSystemErrorInfo(errs.RedisOperateError, err)
		}
		return &outInfo, nil
	}
}

func GetBaseUserOutInfoByUserIds(orgId int64, userIds []int64) ([]bo.BaseUserOutInfoBo, errs.SystemErrorInfo) {
	log.Infof("批量获取用户外部信息 %d, %s", orgId, json.ToJsonIgnoreError(userIds))

	resultList := make([]bo.BaseUserOutInfoBo, 0)

	if userIds == nil || len(userIds) == 0 {
		return resultList, nil
	}

	//用户外部信息
	userOutInfos := &[]po.PpmOrgUserOutInfo{}
	err := mysql.SelectAllByCond(consts.TableUserOutInfo, db.Cond{
		consts.TcIsDelete: consts.AppIsNoDelete,
		consts.TcOrgId:    orgId,
		consts.TcUserId:   db.In(userIds),
	}, userOutInfos)
	if err != nil {
		log.Error(err)
		return nil, errs.MysqlOperateError
	}

	//组织外部信息
	orgOutInfo := &po.PpmOrgOrganizationOutInfo{}
	err = mysql.SelectOneByCond(consts.TableOrganizationOutInfo, db.Cond{
		consts.TcIsDelete: consts.AppIsNoDelete,
		consts.TcOrgId:    orgId,
	}, orgOutInfo)

	msetArgs := map[string]string{}
	keys := make([]string, 0)
	for _, userOutInfo := range *userOutInfos {
		key, err5 := util.ParseCacheKey(sconsts.CacheBaseUserOutInfo, map[string]interface{}{
			consts.CacheKeyOrgIdConstName:  orgId,
			consts.CacheKeyUserIdConstName: userOutInfo.UserId,
		})
		if err5 != nil {
			log.Error(err5)
			return nil, err5
		}
		keys = append(keys, key)

		outInfo := bo.BaseUserOutInfoBo{
			UserId:       userOutInfo.UserId,
			OrgId:        orgId,
			OutUserId:    userOutInfo.OutUserId,
			OutOrgId:     orgOutInfo.OutOrgId,
			OutOrgUserId: userOutInfo.OutOrgUserId,
		}

		resultList = append(resultList, outInfo)
		msetArgs[key] = json.ToJsonIgnoreError(outInfo)
	}

	if len(msetArgs) > 0 {
		err = cache.MSet(msetArgs)
		if err != nil {
			log.Error(err)
			return nil, errs.BuildSystemErrorInfo(errs.RedisOperateError, err)
		}
	}

	go func() {
		defer func() {
			if r := recover(); r != nil {
				log.Error(errs.BuildSystemErrorInfoWithPanicRecover(r, stack.GetStack()))
			}
		}()

		for _, key := range keys {
			_, _ = cache.Expire(key, consts.GetCacheBaseUserInfoExpire())
		}
	}()

	return resultList, nil
}

// sourceChannel可以为空
func getLocalBaseUserInfo(orgId, userId int64, key string) (*bo.BaseUserInfoBo, errs.SystemErrorInfo) {
	user, err := GetUserInfoByUserId(userId)
	if err != nil {
		return nil, err
	}

	baseUserInfo := &bo.BaseUserInfoBo{
		UserId:             user.Id,
		Name:               user.Name,
		NamePy:             user.NamePinyin,
		Avatar:             user.Avatar,
		OrgId:              orgId,
		OrgUserIsDelete:    2,
		OrgUserStatus:      1,
		OrgUserCheckStatus: 1,
	}

	if orgId > 0 {
		newestUserOrganization, err1 := GetUserOrganizationNewestRelation(orgId, userId)
		if err1 != nil {
			log.Error(err1)
			return nil, err1
		}
		baseUserInfo.OrgUserIsDelete = newestUserOrganization.IsDelete
		baseUserInfo.OrgUserStatus = newestUserOrganization.Status
		baseUserInfo.OrgUserCheckStatus = newestUserOrganization.CheckStatus
	}

	baseUserInfoJson, err2 := json.ToJson(baseUserInfo)
	if err2 != nil {
		log.Error(err2)
		return nil, errs.BuildSystemErrorInfo(errs.JSONConvertError, err2)
	}
	err2 = cache.SetEx(key, baseUserInfoJson, consts.GetCacheBaseUserInfoExpire())
	if err2 != nil {
		log.Error(err2)
		return nil, errs.BuildSystemErrorInfo(errs.RedisOperateError, err2)
	}
	return baseUserInfo, nil
}

func GetUserInfoByUserId(userId int64) (*po.PpmOrgUser, errs.SystemErrorInfo) {
	user := &po.PpmOrgUser{}
	err := mysql.SelectById(user.TableName(), userId, user)
	if err != nil {
		log.Errorf("[GetUserInfoByUserId] userId: %d, err: %v", userId, err)
		return nil, errs.BuildSystemErrorInfo(errs.MysqlOperateError, err)
	}

	return user, nil
}

func getLocalBaseUserInfoBatch(orgId int64, userIds []int64) ([]bo.BaseUserInfoBo, errs.SystemErrorInfo) {
	log.Infof("批量获取用户信息 %d, %s", orgId, json.ToJsonIgnoreError(userIds))

	baseUserInfos := make([]bo.BaseUserInfoBo, 0)

	if userIds == nil || len(userIds) == 0 {
		return baseUserInfos, nil
	}

	users := &[]po.PpmOrgUser{}
	err := mysql.SelectAllByCond(consts.TableUser, db.Cond{
		consts.TcId: db.In(userIds),
	}, users)
	if err != nil {
		log.Error(err)
		return nil, errs.BuildSystemErrorInfo(errs.MysqlOperateError, err)
	}

	//获取关联列表，要做去重
	userOrganizationPos := &[]po.PpmOrgUserOrganization{}
	_, selectErr := mysql.SelectAllByCondWithPageAndOrder(consts.TableUserOrganization, db.Cond{
		consts.TcOrgId:  orgId,
		consts.TcUserId: db.In(userIds),
	}, nil, 0, -1, "id asc", userOrganizationPos)
	if selectErr != nil {
		log.Error(selectErr)
		return nil, errs.MysqlOperateError
	}

	//id升序，保留最新: key: userId, value: po
	userOrgMap := map[int64]po.PpmOrgUserOrganization{}
	for _, userOrg := range *userOrganizationPos {
		userOrgMap[userOrg.UserId] = userOrg
	}

	for _, user := range *users {
		baseUserInfo := bo.BaseUserInfoBo{
			UserId: user.Id,
			Name:   user.Name,
			NamePy: user.NamePinyin,
			Avatar: user.Avatar,
			OrgId:  orgId,
		}

		if userOrg, ok := userOrgMap[user.Id]; ok {
			baseUserInfo.OrgUserIsDelete = userOrg.IsDelete
			baseUserInfo.OrgUserStatus = userOrg.Status
			baseUserInfo.OrgUserCheckStatus = userOrg.CheckStatus
		}

		baseUserInfos = append(baseUserInfos, baseUserInfo)
	}

	msetArgs := map[string]string{}
	keys := make([]string, 0)
	for _, baseUserInfo := range baseUserInfos {
		key, err5 := util.ParseCacheKey(sconsts.CacheBaseUserInfo, map[string]interface{}{
			consts.CacheKeyOrgIdConstName:  orgId,
			consts.CacheKeyUserIdConstName: baseUserInfo.UserId,
		})
		if err5 != nil {
			log.Error(err5)
			return nil, err5
		}
		msetArgs[key] = json.ToJsonIgnoreError(baseUserInfo)
		keys = append(keys, key)
	}

	if len(msetArgs) > 0 {
		err = cache.MSet(msetArgs)
		if err != nil {
			log.Error(err)
			return nil, errs.BuildSystemErrorInfo(errs.RedisOperateError, err)
		}
	}

	go func() {
		defer func() {
			if r := recover(); r != nil {
				log.Error(errs.BuildSystemErrorInfoWithPanicRecover(r, stack.GetStack()))
			}
		}()

		for _, key := range keys {
			_, _ = cache.Expire(key, consts.GetCacheBaseUserInfoExpire())
		}
	}()
	return baseUserInfos, nil
}

func GetUserConfigInfo(orgId int64, userId int64) (*bo.UserConfigBo, errs.SystemErrorInfo) {
	userConfig, err := getUserConfigInfo(orgId, userId)
	if err != nil {
		userConfig = &bo.UserConfigBo{}
		userConfigBo, err2 := InsertUserConfig(orgId, userId)
		if err2 != nil {
			log.Error(err2)
			return nil, errs.BuildSystemErrorInfo(errs.UserConfigUpdateError, err2)
		}
		err3 := copyer.Copy(userConfigBo, userConfig)
		if err3 != nil {
			log.Error(err3)
			return nil, errs.BuildSystemErrorInfo(errs.ObjectCopyError, err3)
		}
	}
	return userConfig, nil
}

func getUserConfigInfo(orgId int64, userId int64) (*bo.UserConfigBo, errs.SystemErrorInfo) {
	key, err5 := util.ParseCacheKey(sconsts.CacheUserConfig, map[string]interface{}{
		consts.CacheKeyOrgIdConstName:  orgId,
		consts.CacheKeyUserIdConstName: userId,
	})
	if err5 != nil {
		log.Error(err5)
		return nil, err5
	}

	userConfigJson, err := cache.Get(key)
	if err != nil {
		log.Error(err)
		return nil, errs.BuildSystemErrorInfo(errs.RedisOperateError, err)
	}
	userConfigBo := &bo.UserConfigBo{}
	if userConfigJson != "" {
		err := json.FromJson(userConfigJson, userConfigBo)
		if err != nil {
			log.Error(err)
			return nil, errs.BuildSystemErrorInfo(errs.JSONConvertError, err)
		}
		return userConfigBo, nil
	} else {
		userConfig := &po.PpmOrgUserConfig{}
		err = mysql.SelectOneByCond(userConfig.TableName(), db.Cond{
			consts.TcOrgId:    orgId,
			consts.TcUserId:   userId,
			consts.TcIsDelete: consts.AppIsNoDelete,
		}, userConfig)
		if err != nil {
			log.Error(err)
			return nil, errs.BuildSystemErrorInfo(errs.MysqlOperateError, err)
		}
		_ = copyer.Copy(userConfig, userConfigBo)
		userConfigJson, err = json.ToJson(userConfigBo)
		if err != nil {
			return nil, errs.BuildSystemErrorInfo(errs.JSONConvertError, err)
		}
		err = cache.Set(key, userConfigJson)
		if err != nil {
			log.Error(err)
			return nil, errs.BuildSystemErrorInfo(errs.RedisOperateError, err)
		}
		return userConfigBo, nil
	}
}

func DeleteUserConfigInfo(orgId int64, userId int64) errs.SystemErrorInfo {
	key, err5 := util.ParseCacheKey(sconsts.CacheUserConfig, map[string]interface{}{
		consts.CacheKeyOrgIdConstName:  orgId,
		consts.CacheKeyUserIdConstName: userId,
	})
	if err5 != nil {
		log.Error(err5)
		return err5
	}
	_, err := cache.Del(key)
	if err != nil {
		log.Error(err)
		return errs.BuildSystemErrorInfo(errs.RedisOperateError, err)
	}
	return nil
}

func ClearUserCacheInfo(token string) errs.SystemErrorInfo {
	userCacheKey := sconsts.CacheUserToken + token
	_, err := cache.Del(userCacheKey)
	if err != nil {
		log.Error(err)
		return errs.BuildSystemErrorInfo(errs.RedisOperateError)
	}
	return nil
}

func SetDingCodeCache(outUserId string) errs.SystemErrorInfo {
	key, err5 := util.ParseCacheKey(sconsts.LoginByDingCode, map[string]interface{}{
		consts.CacheKeySourceChannelConstName: sdk_const.SourceChannelDingTalk,
		consts.CacheKeyOutUserIdConstName:     outUserId,
	})
	if err5 != nil {
		log.Error(err5)
		return err5
	}

	err := cache.SetEx(key, outUserId, int64(60*5))
	if err != nil {
		log.Error(err)
		return errs.RedisOperateError
	}

	return nil
}

func GetDingCodeCache(outUserId string) errs.SystemErrorInfo {
	key, err5 := util.ParseCacheKey(sconsts.LoginByDingCode, map[string]interface{}{
		consts.CacheKeySourceChannelConstName: sdk_const.SourceChannelDingTalk,
		consts.CacheKeyOutUserIdConstName:     outUserId,
	})
	if err5 != nil {
		log.Error(err5)
		return err5
	}

	res, err := cache.Get(key)
	if err != nil {
		log.Error(err)
		return errs.RedisOperateError
	}

	if res == consts.BlankString {
		return errs.DingCodeCacheInvalid
	}

	//查到之后清掉
	_, clearErr := cache.Del(key)
	if clearErr != nil {
		log.Error(clearErr)
		return errs.RedisOperateError
	}

	return nil
}

// 用户是否是付费用户
func CheckPaidUser(orgId, userId int64, outOrgId, outUserId, sourceChannel string, orgConfig *bo.OrgConfigBo, errInfoMsg string) errs.SystemErrorInfo {
	payInfoJson, errCache := cache.HGet(consts.CachePayRangeInfo, fmt.Sprintf("%d", orgId))
	if errCache != nil {
		log.Errorf("[CheckPaidUser] cache err:%v", errCache)
		return errs.RedisOperateError
	}

	if payInfoJson == "" {
		// 没有付费的相关数据，这时候需要重新添加
		return SetPayUserInfo(orgId, userId, outOrgId, outUserId, sourceChannel, orgConfig, errInfoMsg)
	}

	payRangeData := bo.PayRangeData{}
	errJson := json.FromJson(payInfoJson, &payRangeData)
	if errJson != nil {
		log.Error(errJson)
		return errs.JSONConvertError
	}

	// 如果飞书设置了付费范围的flag，就只走原来的checkPay逻辑了
	if payRangeData.IsFsSetPayRange == consts.SetFsPayRange {
		return CheckUserPay(orgId, userId, outUserId, errInfoMsg)
	}

	if sourceChannel != sdk_const.SourceChannelWeixin && payRangeData.PayNum == 0 {
		// 理论上在飞书上付费了payNum不会为0的，但 如果是飞书付费版本的试用 回调的付费人数为0
		// 钉钉 如果是试用的话 payNum为 999999，但是在获取付费订单信息的时候过滤掉了，即 payNum还是为0
		return nil
	}

	//if cast.ToTime(payRangeData.EndTime).Before(time.Now()) {
	//	// 过了付费版本有效期
	//	return errs.OrgUserInvalid
	//}
	// 如果席位数大于可见范围人数时，直接进入极星
	if payRangeData.PayNum >= payRangeData.ScopeNum {
		return nil
	} else {
		// 席位数小于可见范围人数
		if ok, errSlice := slice.Contain(payRangeData.UserIds, userId); errSlice == nil && ok {
			return nil
		}
		// 可用范围内前n个人可以进入极星
		// 判断登录的人数是否已经超过了付费人数，超过了就直接报错禁止进入极星
		payRangeData.UserIds = append(payRangeData.UserIds, userId)
		payRangeData.UserIds = slice.SliceUniqueInt64(payRangeData.UserIds)
		if len(payRangeData.UserIds) > payRangeData.PayNum {
			if sourceChannel == sdk_const.SourceChannelFeishu {
				if CheckFsUserPay(outOrgId, outUserId) {
					// 如果当前用户在付费范围内，就 走原来 checkPay的老逻辑 并删除pay_range缓存
					payRangeData = bo.PayRangeData{IsFsSetPayRange: consts.SetFsPayRange}
					errCache = cache.HSet(consts.CachePayRangeInfo, fmt.Sprintf("%d", orgId), json.ToJsonIgnoreError(payRangeData))
					if errCache != nil {
						log.Errorf("[CheckPaidUser] cache err:%v", errCache)
						return errs.RedisOperateError
					}
					return CheckUserPay(orgId, userId, outUserId, errInfoMsg)
				}
			}
			return errs.BuildSystemErrorInfoWithMessage(errs.OrgUserInvalid, errInfoMsg)
		}
		errCache = cache.HSet(consts.CachePayRangeInfo, fmt.Sprintf("%d", orgId), json.ToJsonIgnoreError(payRangeData))
		if errCache != nil {
			log.Errorf("[CheckPaidUser] cache err:%v", errCache)
			return errs.RedisOperateError
		}
	}

	return nil
}

func CheckFsUserPay(outOrgId, outUserId string) bool {
	client, sdkError := platform_sdk.GetClient(sdk_const.SourceChannelFeishu, outOrgId)
	if sdkError != nil {
		log.Errorf("[CheckPaidUser] sdkErr:%v, corpId:%v", sdkError, outOrgId)
		return false
	}
	tenant := client.GetOriginClient().(*sdk.Tenant)
	checkReply, sdkError := tenant.CheckUser(fsvo.CheckUserReq{OpenId: outUserId})
	if sdkError != nil {
		log.Errorf("[CheckPaidUser] sdkErr:%v, corpId:%v", sdkError, outOrgId)
		return false
	}
	if checkReply.Code != 0 {
		log.Errorf("[CheckPaidUser] corpId:%v, code:%v, message:%v", outOrgId, checkReply.Code, checkReply.Msg)
		return false
	}
	if checkReply.Data.Status == consts.FsUserStatusValid {
		return true
	}
	return false
}

func SetPayUserInfo(orgId, userId int64, outOrgId, outUserId, sourceChannel string, orgConfig *bo.OrgConfigBo, errInfoMsg string) errs.SystemErrorInfo {
	// 查询订单的付费信息
	orderPayInfo := orderfacade.GetOrderPayInfo(ordervo.GetOrderPayInfoReq{OrgId: orgId})
	if orderPayInfo.Failure() {
		log.Errorf("[SetPayUserInfo] GetOrderPayInfo err:%v, orgId:%v, userId:%v", orderPayInfo.Error(), orgId, userId)
		return orderPayInfo.Error()
	}
	payNum := orderPayInfo.Data.PayNum

	client, sdkError := platform_sdk.GetClient(sourceChannel, outOrgId)
	if sdkError != nil {
		log.Errorf("[SetPayUserInfo] sdkErr:%v, corpId:%v", sdkError, outOrgId)
		return errs.PlatFormOpenApiCallError
	}
	depts := []*sdkVo.DepartmentInfo{}
	if sourceChannel == sdk_const.SourceChannelDingTalk {
		depResp, sdkError := client.GetScopeDeps()
		if sdkError != nil {
			log.Errorf("[SetPayUserInfo] sdkErr:%v, corpId:%v", sdkError, outOrgId)
			return errs.PlatFormOpenApiCallError
		}
		depts = depResp.Depts
	}
	usersReply, sdkError := client.GetScopeUsers(&sdkVo.GetScopeUsersReq{Depts: depts})
	if sdkError != nil {
		log.Errorf("[SetPayUserInfo] sdkErr:%v, corpId:%v", sdkError, outOrgId)
		return errs.PlatFormOpenApiCallError
	}
	scopeOpenIds := make([]string, 0, len(usersReply.Users))
	for _, u := range usersReply.Users {
		if sourceChannel == sdk_const.SourceChannelDingTalk {
			scopeOpenIds = append(scopeOpenIds, u.UserId)
		} else {
			scopeOpenIds = append(scopeOpenIds, u.OpenId)
		}
	}
	userIds := []int64{}
	// 当前用户在可见范围内
	isInScope := false
	if ok, errSlice := slice.Contain(scopeOpenIds, outUserId); errSlice == nil && ok {
		userIds = append(userIds, userId)
		isInScope = true
	}
	// 不在可见范围内直接提示
	if !isInScope {
		return errs.BuildSystemErrorInfoWithMessage(errs.OrgUserInvalid, errInfoMsg)
	}

	//if sourceChannel == sdk_const.SourceChannelFeishu {
	//	// 如果飞书平台 该用户是第一次进入极星  不在付费范围内直接报错
	//	if !CheckFsUserPay(outOrgId, outUserId) {
	//		log.Infof("[SetPayUserInfo] 不在付费范围内, orgId:%v, userId:%v", orgId, userId)
	//		userIds = businees.DifferenceInt64Set(userIds, []int64{userId})
	//	}
	//}

	//if payNum == 0 {
	//	// 有些很早的组织订单表中可能没有付费人数的数据，但payLevel是2、3，直接放开吧
	//	return nil
	//}

	if len(userIds) == 0 {
		return errs.BuildSystemErrorInfoWithMessage(errs.OrgUserInvalid, errInfoMsg)
	}

	endTime := orgConfig.PayEndTime
	payData := bo.PayRangeData{
		OrgId:         orgId,
		OutOrgId:      outOrgId,
		SourceChannel: sourceChannel,
		ScopeNum:      len(usersReply.Users),
		PayNum:        payNum,
		UserIds:       userIds,
		EndTime:       endTime.Format(consts.AppTimeFormat),
	}

	cacheErr := cache.HSet(consts.CachePayRangeInfo, fmt.Sprintf("%d", orgId), json.ToJsonIgnoreError(payData))
	if cacheErr != nil {
		log.Errorf("[SetPayUserInfo] cache err:%v, orgId:%v, userId:%v", cacheErr, orgId, userId)
		return errs.RedisOperateError
	}
	return nil
}

func ResetOrgPayNum(orgId int64) error {
	// 查询订单的付费信息
	orderPayInfo := orderfacade.GetOrderPayInfo(ordervo.GetOrderPayInfoReq{OrgId: orgId})
	if orderPayInfo.Failure() {
		log.Errorf("[SetPayUserInfo] GetOrderPayInfo err:%v, orgId:%v", orderPayInfo.Error(), orgId)
		return orderPayInfo.Error()
	}
	payNum := orderPayInfo.Data.PayNum

	payRangeInfoJson, errCache := cache.HGet(consts.CachePayRangeInfo, fmt.Sprintf("%d", orgId))
	if errCache != nil {
		log.Error(errCache)
		return errCache
	}
	if payRangeInfoJson != "" {
		rangeData := bo.PayRangeData{}
		_ = json.FromJson(payRangeInfoJson, &rangeData)
		rangeData.PayNum = payNum
		errCache = cache.HSet(consts.CachePayRangeInfo, fmt.Sprintf("%d", orgId), json.ToJsonIgnoreError(rangeData))
		if errCache != nil {
			log.Error(errCache)
			return errCache
		}
	}

	return nil
}

func CheckUserPay(orgId, userId int64, outUserId, errInfoMsg string) errs.SystemErrorInfo {
	key, err1 := util.ParseCacheKey(sconsts.CacheUserCheckPay, map[string]interface{}{
		consts.CacheKeyOrgIdConstName:  orgId,
		consts.CacheKeyUserIdConstName: userId,
	})
	if err1 != nil {
		log.Error(err1)
		return err1
	}

	res, err := cache.Get(key)
	if err != nil {
		log.Error(err)
		return errs.RedisOperateError
	}
	nowTime := time.Now()
	info := &bo.CacheUserCheckPayBo{}
	if res != "" {
		_ = json.FromJson(res, info)
		// 只需判断是否在付费范围。如果在付费范围并且时间已过，则最多多使用几天时间而已，不 care。
		if info.Status == consts.FsUserStatusValid {
			return nil
		} else {
			return errs.BuildSystemErrorInfoWithMessage(errs.OrgUserInvalid, errInfoMsg)
		}
	} else {
		// 查询飞书接口，用户是否在付费范围内
		lockKey := consts.CheckFsUserLock + strconv.FormatInt(orgId, 10) + strconv.FormatInt(userId, 10)
		lockUuid := uuid.NewUuid()
		suc, err := cache.TryGetDistributedLock(lockKey, lockUuid)
		log.Infof("准备获取分布式锁 %v", suc)
		if err != nil {
			log.Error(err)
			return errs.BuildSystemErrorInfo(errs.RedisOperateError)
		}

		if suc {
			log.Infof("获取分布式锁成功 %v", suc)
			defer func() {
				if _, lockErr := cache.ReleaseDistributedLock(lockKey, lockUuid); lockErr != nil {
					log.Error(lockErr)
				}
			}()
			res, err := cache.Get(key)
			if err != nil {
				log.Error(err)
				return errs.RedisOperateError
			}
			if res != "" {
				_ = json.FromJson(res, info)
				if info.Status == consts.FsUserStatusValid {
					return nil
				}
			}

			baseOrgInfo, businessErr := GetBaseOrgInfo(orgId)
			if businessErr != nil {
				log.Errorf("[CheckUserPay] err: %v", businessErr)
				return businessErr
			}
			if baseOrgInfo.OutOrgId == "" {
				err := errs.OrgOutInfoNotExist
				log.Errorf("[CheckUserPay] err: %v", err)
				return err
			}
			tenant, err1 := feishu.GetTenant(baseOrgInfo.OutOrgId)
			if err1 != nil {
				log.Errorf("[CheckUserPay] outOrgId: %s, GetTenant err: %v", baseOrgInfo.OutOrgId, err1)
				return err1
			}
			resp, err2 := tenant.CheckUser(fsvo.CheckUserReq{
				OpenId: outUserId,
			})
			if err2 != nil {
				log.Errorf("[CheckUserPay] outUserId: %s, GetTenant err: %v", outUserId, err1)
				return errs.FeiShuOpenApiCallError
			}
			if resp.Code != 0 {
				log.Errorf("[CheckUserPay] resp: %s", json.ToJsonIgnoreError(resp))
				return errs.FeiShuOpenApiCallError
			}

			info = &bo.CacheUserCheckPayBo{
				UserId: userId,
				Status: resp.Data.Status,
			}
			timeUnix, _ := strconv.ParseInt(resp.Data.ServiceStopTime, 10, 64)
			info.ServiceStopTime = time.Unix(timeUnix, 0)
			// 如果不在付费范围内，则缓存时间为 5s，否则为 7 天（并且失效时间一定是在某一天的中午时间段）
			cacheDurationSec := int64(5)
			if info.Status == consts.FsUserStatusValid && info.ServiceStopTime.After(nowTime) {
				cacheDurationSec = GetDuration2Noon() + int64(7*24*3600)
			}
			// 缓存。
			cacheErr := cache.SetEx(key, json.ToJsonIgnoreError(info), cacheDurationSec)
			if cacheErr != nil {
				log.Errorf("[CheckUserPay] SetEx err: %v", cacheErr)
				return errs.BuildSystemErrorInfo(errs.RedisOperateError, cacheErr)
			}

			if info.Status != consts.FsUserStatusValid {
				return errs.BuildSystemErrorInfoWithMessage(errs.OrgUserInvalid, errInfoMsg)
			}
		}
	}

	return nil
}

// ClearAllOrgUserPayCache 清除org的用户的付费缓存，解决付费后延迟问题
func ClearAllOrgUserPayCache(orgId int64) errs.SystemErrorInfo {
	users, err := GetUserOrgInfos(orgId)
	if err != nil {
		log.Errorf("[ClearAllOrgUserPayCache] GetUserOrgInfos orgId:%v, err:%v", orgId, err)
		return err
	}
	keys := make([]interface{}, 0, len(users))
	for _, user := range users {
		key, _ := util.ParseCacheKey(sconsts.CacheUserCheckPay, map[string]interface{}{
			consts.CacheKeyOrgIdConstName:  orgId,
			consts.CacheKeyUserIdConstName: user.UserId,
		})
		keys = append(keys, key)
	}
	if len(keys) > 0 {
		_, err := cache.Del(keys)
		if err != nil {
			log.Errorf("[ClearAllOrgUserPayCache] Del keys:%v, err:%v", keys, err)
		}
	}

	return nil
}

// GetDuration2Noon 计算当前时间到中午或第二天中午的秒数
func GetDuration2Noon() int64 {
	now := time.Now()
	hour := now.Hour()
	durationSec := int64(0)
	if hour <= 12 {
		// 当天中午
		durationSec = int64((12 - hour) * 3600)
	} else {
		// 当第二天中午
		durationSec = int64((hour - 12) * 3600)
	}
	randSec := rand.Intn(3600)

	return int64(randSec) + durationSec
}

func SetShareUrl(key, url string) errs.SystemErrorInfo {
	if key == "" {
		return nil
	}
	cacheErr := cache.SetEx(sconsts.CacheShareUrl+key, url, 60*60*24*3)
	if cacheErr != nil {
		log.Error(cacheErr)
		return errs.CacheProxyError
	}

	return nil
}

func GetShareUrl(key string) (string, errs.SystemErrorInfo) {
	url, cacheErr := cache.Get(sconsts.CacheShareUrl + key)
	if cacheErr != nil {
		log.Error(cacheErr)
		return "", errs.CacheProxyError
	}

	return url, nil
}

func NeedRemindPayExpire(orgId, userId int64, payEndTime time.Time) (bool, errs.SystemErrorInfo) {
	expireTime := payEndTime.Unix() - time.Now().Unix()
	if expireTime <= 0 {
		//超出时间不需要提醒
		return false, nil
	}

	key, err5 := util.ParseCacheKey(sconsts.CachePayExpireRemind, map[string]interface{}{
		consts.CacheKeyOrgIdConstName: orgId,
	})
	if err5 != nil {
		log.Error(err5)
		return false, err5
	}
	isExist, isExistErr := cache.Exist(key)
	if isExistErr != nil {
		log.Error(isExistErr)
		return false, errs.CacheProxyError
	}
	if isExist {
		isRemind, err := cache.HGet(key, strconv.FormatInt(userId, 10))
		if err != nil {
			log.Error(err)
			return false, errs.CacheProxyError
		}
		if isRemind == "" {
			err := cache.HSet(key, strconv.FormatInt(userId, 10), "true")
			if err != nil {
				log.Error(err)
				return false, errs.CacheProxyError
			}
		} else {
			return false, nil
		}
	} else {
		err := cache.HSet(key, strconv.FormatInt(userId, 10), "true")
		if err != nil {
			log.Error(err)
			return false, errs.CacheProxyError
		}
		_, expireErr := cache.Expire(key, expireTime)
		if expireErr != nil {
			log.Error(expireErr)
			return false, errs.CacheProxyError
		}
		return true, nil
	}

	return false, nil
}

func NeedRemindPayOverdue(orgId, userId int64, payEndTime time.Time) (bool, errs.SystemErrorInfo) {
	expireTime := payEndTime.AddDate(0, 0, 1).Unix() - time.Now().Unix()
	if expireTime <= 0 {
		//超出时间不需要提醒
		return false, nil
	}

	key, err5 := util.ParseCacheKey(sconsts.CachePayOverdueRemind, map[string]interface{}{
		consts.CacheKeyOrgIdConstName: orgId,
	})
	if err5 != nil {
		log.Error(err5)
		return false, err5
	}
	isExist, isExistErr := cache.Exist(key)
	if isExistErr != nil {
		log.Error(isExistErr)
		return false, errs.CacheProxyError
	}
	if isExist {
		isRemind, err := cache.HGet(key, strconv.FormatInt(userId, 10))
		if err != nil {
			log.Error(err)
			return false, errs.CacheProxyError
		}
		if isRemind == "" {
			err := cache.HSet(key, strconv.FormatInt(userId, 10), "true")
			if err != nil {
				log.Error(err)
				return false, errs.CacheProxyError
			}
			return true, nil
		} else {
			return false, nil
		}
	} else {
		err := cache.HSet(key, strconv.FormatInt(userId, 10), "true")
		if err != nil {
			log.Error(err)
			return false, errs.CacheProxyError
		}
		_, expireErr := cache.Expire(key, expireTime)
		if expireErr != nil {
			log.Error(expireErr)
			return false, errs.CacheProxyError
		}
		return true, nil
	}

	return false, nil
}
