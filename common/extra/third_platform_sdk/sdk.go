package third_platform_sdk

import (
	"fmt"

	sdk "gitea.bjx.cloud/allstar/platform-sdk"
	sdk_const "gitea.bjx.cloud/allstar/platform-sdk/consts"
	"gitea.bjx.cloud/allstar/platform-sdk/sdk_interface"
	"gitea.bjx.cloud/allstar/platform-sdk/vo/register_info"
	"github.com/star-table/startable-server/common/core/config"
	"github.com/spf13/cast"
)

func InitPlatFormSdk(f sdk_interface.GetCorpInfoFromDB) {
	redis := config.GetConfig().Redis
	redisAddr := redis.Host + ":" + cast.ToString(redis.Port)

	fsConfig := config.GetConfig().FeiShu
	if fsConfig == nil {
		panic("feishu config not found")
	}
	sdk.RegisterPlatformInfo(&register_info.PlatformInfo{
		AppKey:        fsConfig.AppId,
		AppSecret:     fsConfig.AppSecret,
		RunType:       sdk_const.RunTypeCorp,
		SourceChannel: sdk_const.SourceChannelFeishu,
	}, redisAddr, redis.Pwd, redis.Database, redis.IsSentinel, f)

	dingConfig := config.GetConfig().DingTalk
	if dingConfig == nil {
		panic("dingTalk config not found")
	}
	sdk.RegisterPlatformInfo(&register_info.PlatformInfo{
		AppId:         fmt.Sprintf("%d", dingConfig.AppId),
		AppKey:        dingConfig.SuiteKey,
		AppSecret:     dingConfig.SuiteSecret,
		RunType:       sdk_const.RunTypeCorp,
		SourceChannel: sdk_const.SourceChannelDingTalk,
	}, redisAddr, redis.Pwd, redis.Database, redis.IsSentinel, f)

	WeComConfig := config.GetConfig().WeCom
	if WeComConfig == nil {
		panic("WeCom config not found")
	}
	customInfos := make([]*register_info.CustomInfo, 0, 10)
	for _, info := range WeComConfig.CustomInfos {
		customInfos = append(customInfos, &register_info.CustomInfo{
			CorpId:  info.CorpId,
			AgentId: info.AgentId,
			Secret:  info.Secret,
		})
	}
	sdk.RegisterPlatformInfo(&register_info.PlatformInfo{
		AppId:          WeComConfig.CorpId,
		ProviderSecret: WeComConfig.ProviderSecret,
		AppKey:         WeComConfig.SuiteId,
		AppSecret:      WeComConfig.Secret,
		RunType:        sdk_const.RunTypeCorp,
		SourceChannel:  sdk_const.SourceChannelWeixin,
		CustomInfos:    customInfos,
	}, redisAddr, redis.Pwd, redis.Database, redis.IsSentinel, f)
}
