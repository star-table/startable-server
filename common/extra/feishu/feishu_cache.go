package feishu

import (
	"gitea.bjx.cloud/allstar/feishu-sdk-golang/sdk"
	"github.com/star-table/startable-server/common/core/config"
	"github.com/star-table/startable-server/common/core/consts"
	"github.com/star-table/startable-server/common/core/errs"
	"github.com/star-table/startable-server/common/core/util/json"
	"github.com/star-table/startable-server/common/library/cache"
)

// 60s * 70 -> 70分钟
const CacheInvalidSeconds = 60 * 70

func ResendAppTicket() {
	fsConfig := config.GetConfig().FeiShu
	if fsConfig == nil {
		log.Error("飞书配置为空")
		return
	}
	resp, err := sdk.AppTicketResend(fsConfig.AppId, fsConfig.AppSecret)
	if err != nil {
		log.Error(err)
		return
	}
	log.Infof("app_ticket 重新发送请求响应 %s", json.ToJsonIgnoreError(resp))
}

func ClearFsScopeCache(tenantKey string) errs.SystemErrorInfo {
	key := consts.CacheFeiShuScope + tenantKey
	_, err := cache.Del(key)
	if err != nil {
		log.Error(err)
		return errs.CacheProxyError
	}
	return nil
}

func GetScopeOpenIdsFromCache(tenantKey string) ([]string, errs.SystemErrorInfo) {
	key := consts.CacheFeiShuScope + tenantKey

	value, cacheErr := cache.Get(key)
	if cacheErr != nil {
		log.Error(cacheErr)
		return nil, errs.CacheProxyError
	}
	if value != "" {
		result := &[]string{}
		_ = json.FromJson(value, result)
		return *result, nil
	} else {
		scopeOpenIds, err := GetScopeOpenIds(tenantKey)
		if err != nil {
			log.Error(err)
			return nil, err
		}

		value := json.ToJsonIgnoreError(scopeOpenIds)
		cacheErr = cache.SetEx(key, value, consts.GetCacheBaseExpire())
		if cacheErr != nil {
			log.Error(cacheErr)
			return nil, errs.CacheProxyError
		}

		return scopeOpenIds, nil
	}
}

func getAppAccessTokenCacheKey(appId, appSecret string) string {
	return consts.CacheFeiShuAppAccessToken + appId + ":" + appSecret
}

func getTenantAccessTokenCacheKey(tenant string) string {
	return consts.CacheFeiShuTenantAccessToken + tenant
}
