package orgsvc

import (
	"fmt"
	"net/url"
	"strconv"
	"time"

	"github.com/star-table/startable-server/common/core/config"

	wework "github.com/go-laoji/wecom-go-sdk"
	"github.com/google/martian/log"

	"github.com/star-table/startable-server/common/library/db/mysql"
	"upper.io/db.v3/lib/sqlbuilder"

	"gitea.bjx.cloud/allstar/feishu-sdk-golang/core/util/encrypt"
	platform_sdk "gitea.bjx.cloud/allstar/platform-sdk"
	sdk_const "gitea.bjx.cloud/allstar/platform-sdk/consts"
	"github.com/star-table/startable-server/common/core/consts"
	"github.com/star-table/startable-server/common/core/errs"
	"github.com/star-table/startable-server/common/core/util/random"
	"github.com/star-table/startable-server/common/model/bo"
	"github.com/star-table/startable-server/common/model/vo"
	"github.com/star-table/startable-server/common/model/vo/orgvo"
	"upper.io/db.v3"
)

const (
	registerUrl = "https://open.work.weixin.qq.com/3rdservice/wework/register?register_code="
	installUrl  = "https://open.work.weixin.qq.com/3rdapp/install?suite_id=%v&pre_auth_code=%v&redirect_uri=%v&state=STATE&from=&auth_from_open=&scene=&open_report_key=#pageType=Install"
	redirectUrl = "https://app.startable.cn"
)

func GetWeiXinJsAPISign(input orgvo.JsAPISignReq) (*orgvo.JsAPISignResp, errs.SystemErrorInfo) {
	resp := &orgvo.JsAPISignResp{}
	outOrgInfo, err := domain.GetOrgOutInfoByTenantKey(input.CorpID)
	if err != nil {
		return nil, errs.BuildSystemErrorInfo(errs.DingTalkAvoidCodeInvalidError, err)
	}
	if outOrgInfo.OrgId == consts.PreviewTplOrgId {
		return resp, nil
	}
	ticket, err := GetJsAPITicket(outOrgInfo.OrgId, input.Type)
	if err != nil {
		return nil, errs.BuildSystemErrorInfo(errs.DingTalkAvoidCodeInvalidError, err)
	}
	resp.NoceStr = random.RandomString(5)
	resp.TimeStamp = strconv.FormatInt(time.Now().Unix(), 10)
	plain := "jsapi_ticket=" + ticket + "&noncestr=" + resp.NoceStr + "&timestamp=" + resp.TimeStamp + "&url=" + input.URL
	resp.Signature = encrypt.SHA1(plain)
	corpInfo, err2 := platform_sdk.GetPlatformInfo(sdk_const.SourceChannelWeixin).Cache.GetCorpInfo(input.CorpID, 0)
	if err2 != nil {
		return nil, errs.BuildSystemErrorInfo(errs.DingTalkAvoidCodeInvalidError, err2)
	}
	resp.AgentID = corpInfo.AgentId
	return resp, nil
}

func GetWeiXinRegisterUrl() (*orgvo.GetRegisterUrlData, errs.SystemErrorInfo) {
	client, sdkErr := platform_sdk.GetClient(sdk_const.SourceChannelWeixin, "")
	if sdkErr != nil {
		log.Errorf("[GetWeiXinLoginUserInfo] platform_sdk.GetClient err: %v", sdkErr)
		return nil, errs.BuildSystemErrorInfo(errs.PlatFormOpenApiCallError, sdkErr)
	}
	weWork := client.GetOriginClient().(wework.IWeWork)

	resp := weWork.GetRegisterCode(&wework.GetRegisterCodeReq{
		TemplateId: "tpl584e0d30ffa9b8a2",
		FollowUser: config.GetConfig().WeCom.BuyerUserId,
	})
	if resp.ErrCode != 0 {
		log.Errorf("[GetWeiXinRegisterUrl] GetRegisterCode err: %v, code:%v", resp.ErrorMsg, resp.ErrCode)
		return nil, errs.PlatFormOpenApiCallError
	}

	return &orgvo.GetRegisterUrlData{Url: registerUrl + resp.RegisterCode}, nil
}

func GetWeiXinInstallUrl() (*orgvo.GetRegisterUrlData, errs.SystemErrorInfo) {
	client, sdkErr := platform_sdk.GetClient(sdk_const.SourceChannelWeixin, "")
	if sdkErr != nil {
		log.Errorf("[GetWeiXinInstallUrl] platform_sdk.GetClient err: %v", sdkErr)
		return nil, errs.BuildSystemErrorInfo(errs.PlatFormOpenApiCallError, sdkErr)
	}
	weWork := client.GetOriginClient().(wework.IWeWork)

	resp := weWork.GetPreAuthCode()
	if resp.ErrCode != 0 {
		log.Errorf("[GetWeiXinInstallUrl] GetRegisterCode err: %v, code:%v", resp.ErrorMsg, resp.ErrCode)
		return nil, errs.PlatFormOpenApiCallError
	}

	return &orgvo.GetRegisterUrlData{Url: fmt.Sprintf(installUrl, config.GetConfig().WeCom.SuiteId, resp.PreAuthCode, url.QueryEscape(redirectUrl))}, nil
}

func AuthWeiXinCode(input orgvo.AuthWeiXinCodeReqVo) (*orgvo.PlatformAuthCodeData, errs.SystemErrorInfo) {
	if input.CodeType == consts.WeiXinPlatformCodeType {
		return thirdAuthCode(sdk_const.SourceChannelWeixin, input.CorpId, input.Code, input.CodeType)
	}
	authCodeInfo, err := domain.GetWeiXinLoginUserInfo(input)
	if err != nil {
		log.Errorf("[AuthWeiXinCode] GetWeiXinLoginUserInfo err:%v", err)
		return nil, err
	}
	err = platformLoginByAuthCodeInfo(sdk_const.SourceChannelWeixin, authCodeInfo)
	if err != nil {
		log.Errorf("[AuthWeiXinCode] platformLoginByAuthCodeInfo err:%v", err)
		return nil, err
	}
	res := &orgvo.PlatformAuthCodeData{
		CorpId:       authCodeInfo.CorpId,
		OutUserId:    authCodeInfo.OutUserId,
		IsAdmin:      authCodeInfo.IsAdmin,
		Binding:      authCodeInfo.Binding,
		RefreshToken: authCodeInfo.RefreshToken,
		AccessToken:  authCodeInfo.AccessToken,
		CodeToken:    random.Token(),
		OrgID:        authCodeInfo.OrgID,
		OrgName:      authCodeInfo.OrgName,
		OutOrgName:   authCodeInfo.OutOrgName,
		OrgCode:      authCodeInfo.OrgCode,
		UserID:       authCodeInfo.UserID,
		Name:         authCodeInfo.Name,
		Token:        authCodeInfo.Token,
	}
	return res, nil
}

func PersonWeiXinQrCode(input *orgvo.PersonWeiXinQrCodeReqVo) (*orgvo.PersonWeiXinRqCodeData, errs.SystemErrorInfo) {
	token, err := domain.GetOfficialWeiXinAccessToken()
	if err != nil {
		log.Errorf("[PersonWeiXinQrCode] err:%v", err)
		return nil, err
	}

	sceneKey := random.Token()
	qrCode, err2 := weixinfacade.GetQrCode(token, sceneKey)
	if err2 != nil {
		log.Errorf("[PersonWeiXinQrCode] GetQrCode token:%v", token)
		return nil, errs.PlatFormOpenApiCallError
	}

	return &orgvo.PersonWeiXinRqCodeData{
		Ticket:   qrCode.Ticket,
		SceneKey: sceneKey,
	}, nil
}

func PersonWeiXinLogin(input *orgvo.PersonWeiXinLoginReqVo) (*orgvo.PersonWeiXinLoginData, errs.SystemErrorInfo) {
	accessInfo, err := weixinfacade.GetAccessTokenByCode(input.Code)
	if err != nil {
		return nil, errs.OpenAccessTokenInvalid
	}

	thirdLoginUser, err := dao.GetGlobalUser().GetGlobalUserByWeiXinId(accessInfo.Openid)
	if err != nil {
		if err == db.ErrNoMoreRows {
			err2 := domain.SetThirdLoginCache(input.Code, &bo.CacheTempThirdUserInfo{
				AccessToken:  accessInfo.AccessToken,
				RefreshToken: accessInfo.RefreshToken,
				OpenId:       accessInfo.Openid,
				UnionId:      accessInfo.UnionId,
			})
			return &orgvo.PersonWeiXinLoginData{
				Binding: false,
			}, err2
		}
		log.Errorf("[PersonWeiXinLogin] error:%v", err)
		return nil, errs.MysqlOperateError
	}

	userBo, err2 := domain.GetUserInfoByGlobalUser(thirdLoginUser)
	if err2 != nil {
		return nil, err2
	}

	res, err2 := userLogin(userBo, "", consts.TokenTypeNormal)
	if err2 != nil {
		return nil, err2
	}

	return &orgvo.PersonWeiXinLoginData{
		Binding: true,
		Token:   res.Token,
		OrgID:   res.OrgID,
		OrgName: res.OrgName,
		UserID:  res.UserID,
		Name:    res.Name,
	}, nil
}

func PersonWeiXinQrCodeScan(input *orgvo.QrCodeScanReq) errs.SystemErrorInfo {
	return domain.SetThirdLoginCache(input.SceneKey, &bo.CacheTempThirdUserInfo{
		OpenId: input.OpenId,
	})
}

func CheckPersonWeiXinQrCode(input *orgvo.CheckQrCodeScanReq) (*orgvo.CheckQrCodeScanData, errs.SystemErrorInfo) {
	info, err := domain.GetThirdLoginCache(input.SceneKey)
	if err != nil {
		return nil, err
	}
	if info.OpenId == "" {
		return &orgvo.CheckQrCodeScanData{}, nil
	}

	globalUser, err2 := dao.GetGlobalUser().GetGlobalUserByWeiXinId(info.OpenId)
	if err2 != nil {
		if err2 == db.ErrNoMoreRows {
			return &orgvo.CheckQrCodeScanData{IsScan: true}, nil
		}
		log.Errorf("[CheckPersonWeiXinQrCode] error:%v", err2)
		return nil, errs.MysqlOperateError
	}

	userBo, err := domain.GetUserInfoByGlobalUser(globalUser)
	if err != nil {
		return nil, err
	}

	res, err := userLogin(userBo, "", consts.TokenTypeNormal)
	if err != nil {
		return nil, err
	}

	return &orgvo.CheckQrCodeScanData{
		IsScan:  true,
		Binding: true,
		Token:   res.Token,
		OrgID:   res.OrgID,
		OrgName: res.OrgName,
		UserID:  res.UserID,
		Name:    res.Name,
	}, nil
}

func CheckMobileHasBind(input *orgvo.CheckMobileHasBindData) (*orgvo.CheckMobileHasBindResult, errs.SystemErrorInfo) {
	if input.SourcePlatform == "" {
		input.SourcePlatform = consts.AppSourcePlatformPersonWeixin
	}

	globalUser, err := dao.GetGlobalUser().GetGlobalUserByMobile(input.Mobile)
	if err != nil {
		if err == db.ErrNoMoreRows {
			return &orgvo.CheckMobileHasBindResult{IsBind: false}, nil
		}
		return nil, errs.MysqlOperateError
	}
	if globalUser.WeiXinOpenId != "" {
		return &orgvo.CheckMobileHasBindResult{IsBind: true}, nil
	}

	return &orgvo.CheckMobileHasBindResult{IsBind: false}, nil
}

// PersonWeiXinBindExistAccount 绑定已经存在的账号
func PersonWeiXinBindExistAccount(req *orgvo.PersonWeiXinBindExistAccountReq) errs.SystemErrorInfo {
	info, err := domain.GetThirdLoginCache(req.Input.SceneKey)
	if err != nil {
		return err
	}
	_, err2 := dao.GetGlobalUser().GetGlobalUserByWeiXinId(info.OpenId)
	if err2 == nil {
		return errs.WeiXinAlreadyBindError
	} else if err2 != db.ErrNoMoreRows {
		log.Errorf("[PersonWeiXinBindExistAccount] error:%v", err2)
		return errs.MysqlOperateError
	}

	globalId, errDao := dao.GetGlobalUserRelation().GetGlobalUserIdByUserId(req.UserId)
	if errDao != nil {
		return errs.MobileNotBindError
	}

	globalUser, errDao := dao.GetGlobalUser().GetGlobalUserById(globalId)
	if errDao != nil {
		return errs.MobileNotBindError
	}
	if globalUser.WeiXinOpenId != "" {
		return errs.AccountHadBindWeiXin
	}

	err2 = mysql.TransX(func(tx sqlbuilder.Tx) error {
		err := dao.GetGlobalUser().UpdateWeiXinOpenId(tx, globalUser.Mobile, info.OpenId)
		if err != nil {
			return err
		}
		return dao.GetThirdLogin().AddThirdLoginInfo(tx, &po.PpmOrgThirdLoginUser{
			UserId:         req.UserId,
			ThirdUserId:    info.OpenId,
			SourcePlatform: consts.AppSourcePlatformPersonWeixin,
		})
	})
	if err2 != nil {
		log.Errorf("[PersonWeiXinBindExistAccount] AddThirdLoginInfo error:%v", err2)
		return errs.MysqlOperateError
	}

	return nil
}

func PersonWeiXinBind(input *orgvo.PersonWeiXinBindData) (*orgvo.PersonWeiXinLoginData, errs.SystemErrorInfo) {
	info, err := domain.GetThirdLoginCache(input.SceneKey)
	if err != nil {
		return nil, err
	}

	_, err2 := dao.GetGlobalUser().GetGlobalUserByWeiXinId(info.OpenId)
	if err2 == nil {
		return nil, errs.WeiXinAlreadyBindError
	} else if err2 != db.ErrNoMoreRows {
		log.Errorf("[PersonWeiXinBind] error:%v", err2)
		return nil, errs.MysqlOperateError
	}

	userBo, err := UserAuthCodeLogin(vo.UserLoginReq{
		LoginType: 1,
		LoginName: input.Mobile,
		Password:  nil,
		AuthCode:  &input.AuthCode,
	})
	if err != nil {
		log.Errorf("[PersonWeiXinBind] UserAuthCodeLogin mobile:%v, code:%v, err:%v", input.Mobile, input.AuthCode, err)
		return nil, err
	}

	err2 = mysql.TransX(func(tx sqlbuilder.Tx) error {
		err := dao.GetGlobalUser().UpdateWeiXinOpenId(tx, input.Mobile, info.OpenId)
		if err != nil {
			return err
		}
		return dao.GetThirdLogin().AddThirdLoginInfo(tx, &po.PpmOrgThirdLoginUser{
			UserId:         userBo.ID,
			ThirdUserId:    info.OpenId,
			SourcePlatform: consts.AppSourcePlatformPersonWeixin,
		})
	})

	if err2 != nil {
		log.Errorf("[PersonWeiXinBind] UpdateWeiXinOpenId error:%v", err2)
		return nil, errs.MysqlOperateError
	}

	res, err := userLogin(userBo, "", consts.TokenTypeNormal)
	if err != nil {
		return nil, err
	}
	return &orgvo.PersonWeiXinLoginData{
		Binding: true,
		Token:   res.Token,
		OrgID:   res.OrgID,
		OrgName: res.OrgName,
		UserID:  res.UserID,
		Name:    res.Name,
	}, nil
}

func UnbindThirdAccount(input orgvo.UnbindAccountReq) errs.SystemErrorInfo {
	if input.Input.SourcePlatform == "" {
		input.Input.SourcePlatform = consts.AppSourcePlatformPersonWeixin
	}

	globalId, errDao := dao.GetGlobalUserRelation().GetGlobalUserIdByUserId(input.UserId)
	if errDao != nil {
		return errs.MysqlOperateError
	}

	globalUser, errDao := dao.GetGlobalUser().GetGlobalUserById(globalId)
	if errDao != nil {
		return errs.MysqlOperateError
	}
	if globalUser.WeiXinOpenId == "" {
		return errs.AccountNotBindWeiXin
	}

	err := mysql.TransX(func(tx sqlbuilder.Tx) error {
		err := dao.GetGlobalUser().UnbindWeiXinOpenId(tx, globalId)
		if err != nil {
			return err
		}

		return dao.GetThirdLogin().UnBindThirdLogin(tx, globalUser.WeiXinOpenId)
	})

	if err != nil {
		return errs.MysqlOperateError
	}

	return nil
}
