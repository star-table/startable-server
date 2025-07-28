package dingtalk

import (
	"gitea.bjx.cloud/allstar/dingtalk-sdk-golang"
	platform_sdk "gitea.bjx.cloud/allstar/platform-sdk"
	sdk_const "gitea.bjx.cloud/allstar/platform-sdk/consts"
	"github.com/star-table/startable-server/common/core/errs"
	"github.com/star-table/startable-server/common/core/logger"
)

var log = logger.GetDefaultLogger()

func GetDingTalk(corpId string) (*dingtalk.DingTalk, errs.SystemErrorInfo) {
	client, err := platform_sdk.GetClient(sdk_const.SourceChannelDingTalk, corpId)
	if err != nil {
		log.Errorf("[GetDingTalk] GetClient corpId:%v, err:%v", corpId, err)
		return nil, errs.BuildSystemErrorInfo(errs.PlatFormOpenApiCallError, err)
	}

	return client.GetOriginClient().(*dingtalk.DingTalk), nil
}
