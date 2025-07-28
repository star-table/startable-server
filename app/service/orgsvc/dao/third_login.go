package orgsvc

import (
	"fmt"

	platform_sdk "gitea.bjx.cloud/allstar/platform-sdk"
	sdk_const "gitea.bjx.cloud/allstar/platform-sdk/consts"
	sdkVo "gitea.bjx.cloud/allstar/platform-sdk/vo"
	sconsts "github.com/star-table/startable-server/app/service"
	"github.com/star-table/startable-server/common/core/consts"
	"github.com/star-table/startable-server/common/core/errs"
	"github.com/star-table/startable-server/common/core/util/json"
	"github.com/star-table/startable-server/common/library/cache"
	"github.com/star-table/startable-server/common/model/bo"
	"github.com/star-table/startable-server/common/model/vo/orgvo"
	wework "github.com/go-laoji/wecom-go-sdk"
	"github.com/google/martian/log"
)

func SetThirdLoginCache(token string, info *bo.CacheTempThirdUserInfo) errs.SystemErrorInfo {
	key := fmt.Sprintf(sconsts.CacheThirdUserInfo, token)
	s, _ := json.ToJson(info)
	if err := cache.SetEx(key, s, consts.GetCacheBaseExpire()); err != nil {
		return errs.BuildSystemErrorInfo(errs.RedisOperateError, err)
	}

	return nil
}

func GetThirdLoginCache(token string) (*bo.CacheTempThirdUserInfo, errs.SystemErrorInfo) {
	info := &bo.CacheTempThirdUserInfo{}
	key := fmt.Sprintf(sconsts.CacheThirdUserInfo, token)
	s, err := cache.Get(key)
	if err != nil {
		return nil, errs.BuildSystemErrorInfo(errs.RedisOperateError, err)
	}

	_ = json.FromJson(s, info)

	return info, nil
}

func GetOfficialWeiXinAccessToken() (string, errs.SystemErrorInfo) {
	token, err := cache.Get(sconsts.CacheOfficialAccessToken)
	if err != nil {
		return "", errs.BuildSystemErrorInfo(errs.RedisOperateError, err)
	}

	return token, nil
}

func GetWeiXinLoginUserInfo(input orgvo.AuthWeiXinCodeReqVo) (*bo.AuthCodeInfoBo, errs.SystemErrorInfo) {
	client, sdkErr := platform_sdk.GetClient(sdk_const.SourceChannelWeixin, "")
	if sdkErr != nil {
		log.Errorf("[GetWeiXinLoginUserInfo] platform_sdk.GetClient err: %v", sdkErr)
		return nil, errs.BuildSystemErrorInfo(errs.PlatFormOpenApiCallError, sdkErr)
	}
	weWork := client.GetOriginClient().(wework.IWeWork)
	loginInfoResp := weWork.GetLoginInfo(input.Code)
	if loginInfoResp.ErrCode != 0 {
		log.Errorf("[GetWeiXinLoginUserInfo] GetLoginInfo err: %v, code:%v", loginInfoResp.ErrorMsg, loginInfoResp.ErrCode)
		return nil, errs.PlatFormOpenApiCallError
	}

	log.Infof("[GetWeiXinLoginUserInfo] info:%v", json.ToJsonIgnoreError(loginInfoResp))

	res := &bo.AuthCodeInfoBo{
		CorpId:    loginInfoResp.CorpInfo.CorpId,
		OutUserId: loginInfoResp.UserInfo.UserId,
		IsAdmin:   false,
		Binding:   true,
		Name:      loginInfoResp.UserInfo.Name,
	}
	orgInfoResp, _ := client.GetOrgInfo(&sdkVo.OrgInfoReq{CorpId: loginInfoResp.CorpInfo.CorpId})
	if orgInfoResp != nil {
		res.OutOrgName = orgInfoResp.Name
	}
	return res, nil
}
