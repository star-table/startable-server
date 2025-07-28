package weixin

import (
	platform_sdk "gitea.bjx.cloud/allstar/platform-sdk"
	sdk_const "gitea.bjx.cloud/allstar/platform-sdk/consts"
	"github.com/star-table/startable-server/common/core/errs"
	"github.com/star-table/startable-server/common/core/logger"
	"github.com/star-table/startable-server/common/model/vo/commonvo"
	wework "github.com/go-laoji/wecom-go-sdk"
)

var log = logger.GetDefaultLogger()

func GetWeWork(corpId string) (*commonvo.WeWork, errs.SystemErrorInfo) {
	client, err := platform_sdk.GetClient(sdk_const.SourceChannelWeixin, corpId)
	if err != nil {
		log.Errorf("[GetWeWork] GetClient corpId:%v, err:%v", corpId, err)
		return nil, errs.BuildSystemErrorInfo(errs.PlatFormOpenApiCallError, err)
	}

	return &commonvo.WeWork{
		Sdk:        client.GetOriginClient().(wework.IWeWork),
		CorpIdUint: client.GetCorpId().(uint),
	}, nil
}
