package orgsvc

import (
	"gitea.bjx.cloud/allstar/platform-sdk/sdk_interface"
	lru "github.com/hashicorp/golang-lru"
	"github.com/spf13/cast"
)

var lruCache, _ = lru.New(6000)

func SetCacheCorpInfo(corpInfo *sdk_interface.CorpInfo) {
	lruCache.Add(corpInfo.CorpId, corpInfo)
	lruCache.Add(corpInfo.OrgId, corpInfo)
}

func GetCorpInfoFromDB(orgId int64, corpId string) (*sdk_interface.CorpInfo, error) {
	var (
		value  interface{}
		isFind bool
	)
	if orgId > 0 {
		value, isFind = lruCache.Get(orgId)
	} else {
		value, isFind = lruCache.Get(corpId)
	}
	if isFind {
		return value.(*sdk_interface.CorpInfo), nil
	}

	outInfo, err := GetOrgOutInfoByOutOrgId(orgId, corpId)
	if err != nil {
		return nil, err
	}
	return &sdk_interface.CorpInfo{
		OrgId:         outInfo.OrgId,
		AgentId:       cast.ToInt64(outInfo.TenantCode),
		CorpId:        outInfo.OutOrgId,
		PermanentCode: outInfo.AuthTicket,
	}, nil
}
