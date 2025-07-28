package orgsvc

import (
	"fmt"
	"strings"
	"time"

	"gitea.bjx.cloud/allstar/feishu-sdk-golang/sdk"
	platform_sdk "gitea.bjx.cloud/allstar/platform-sdk"
	sdk_const "gitea.bjx.cloud/allstar/platform-sdk/consts"
	"gitea.bjx.cloud/allstar/platform-sdk/sdk_interface"
	sdkVo "gitea.bjx.cloud/allstar/platform-sdk/vo"
	"gitea.bjx.cloud/allstar/polaris-backend/facade/idfacade"
	sconsts "github.com/star-table/startable-server/app/service"
	"github.com/star-table/startable-server/common/core/config"
	"github.com/star-table/startable-server/common/core/consts"
	"github.com/star-table/startable-server/common/core/errs"
	"github.com/star-table/startable-server/common/core/util"
	"github.com/star-table/startable-server/common/core/util/json"
	"github.com/star-table/startable-server/common/core/util/random"
	"github.com/star-table/startable-server/common/core/util/stack"
	"github.com/star-table/startable-server/common/core/util/uuid"
	"github.com/star-table/startable-server/common/extra/feishu"
	"github.com/star-table/startable-server/common/library/cache"
	"github.com/star-table/startable-server/common/library/db/mysql"
	"github.com/star-table/startable-server/common/model/bo"
	"github.com/star-table/startable-server/common/model/vo"
	"github.com/google/martian/log"
	"upper.io/db.v3"
	"upper.io/db.v3/lib/sqlbuilder"
)

func FeiShuAuth(input vo.FeiShuAuthReq) (*vo.FeiShuAuthResp, errs.SystemErrorInfo) {
	codeType := 1
	if input.CodeType != nil {
		codeType = *input.CodeType
	}
	log.Info(json.ToJsonIgnoreError(input))
	log.Infof("codeType %d", codeType)

	codeAuthInfo, err := getPlatformUserInfoByAuthCode(sdk_const.SourceChannelFeishu, "", input.Code, codeType)
	if err != nil {
		log.Error(err)
		return nil, err
	}

	openId := codeAuthInfo.OutUserId
	tenantKey := codeAuthInfo.CorpId
	refreshToken := codeAuthInfo.RefreshToken
	accessToken := codeAuthInfo.AccessToken
	isAdmin := codeAuthInfo.IsAdmin

	log.Info("开始初始化当前用户")
	baseUserInfo, err1 := domain.PlatformAuth(sdk_const.SourceChannelFeishu, tenantKey, openId, accessToken)
	if err1 != nil {
		log.Error(err1)
		return &vo.FeiShuAuthResp{
			TenantKey: tenantKey,
			OpenID:    openId,
		}, err1
	}

	res := &vo.FeiShuAuthResp{
		Token:     random.GenTokenWithOrgId(baseUserInfo.OrgId, baseUserInfo.UserId),
		UserID:    baseUserInfo.UserId,
		Name:      baseUserInfo.Name,
		IsAdmin:   isAdmin,
		TenantKey: tenantKey,
		OpenID:    openId,
	}
	//封装组织信息
	if baseUserInfo.OrgId != 0 {
		organizationBo, err := domain.GetOrgBoById(baseUserInfo.OrgId)
		if err != nil {
			log.Error(err)
		} else {
			res.OrgCode = organizationBo.Code
			res.OrgID = organizationBo.Id
			res.OrgName = organizationBo.Name
			res.Token = random.GenTokenWithOrgId(res.OrgID, res.UserID)
		}
	}
	cacheUserInfo := &bo.CacheUserInfoBo{
		OutUserId:     baseUserInfo.OutUserId,
		UserId:        baseUserInfo.UserId,
		CorpId:        tenantKey,
		SourceChannel: sdk_const.SourceChannelFeishu,
		OrgId:         baseUserInfo.OrgId,
	}
	cacheUserInfoJson, err2 := json.ToJson(cacheUserInfo)
	if err2 != nil {
		return nil, errs.BuildSystemErrorInfo(errs.JSONConvertError, err2)
	}
	_ = cache.SetEx(sconsts.CacheUserToken+res.Token, cacheUserInfoJson, consts.CacheUserTokenExpire)

	//缓存refresh_token
	platform_sdk.GetPlatformInfo(sdk_const.SourceChannelFeishu).Cache.SetRefreshToken(baseUserInfo.OrgId, baseUserInfo.UserId, refreshToken, accessToken)

	//更新上次登录时间
	go func() {
		defer func() {
			if r := recover(); r != nil {
				log.Error(errs.BuildSystemErrorInfoWithPanicRecover(r, stack.GetStack()))
			}
		}()
		if err := domain.UserLoginHook(cacheUserInfo.UserId, cacheUserInfo.OrgId, sdk_const.SourceChannelFeishu); err != nil {
			log.Error(err)
		}
	}()

	return res, nil
}

func FeiShuAuthCode(code string, codeType int) (*vo.FeiShuAuthCodeResp, errs.SystemErrorInfo) {
	authCodeInfo, err := getPlatformUserInfoByAuthCode(sdk_const.SourceChannelFeishu, "", code, codeType)
	if err != nil {
		log.Error(err)
		return nil, err
	}

	if authCodeInfo.Binding {
		err := platformLoginByAuthCodeInfo(sdk_const.SourceChannelFeishu, authCodeInfo)
		if err != nil {
			return nil, err
		}
	}

	authCodeToken := random.Token()
	res := &vo.FeiShuAuthCodeResp{
		TenantKey:    authCodeInfo.CorpId,
		OpenID:       authCodeInfo.OutUserId,
		IsAdmin:      authCodeInfo.IsAdmin,
		Binding:      authCodeInfo.Binding,
		RefreshToken: authCodeInfo.RefreshToken,
		AccessToken:  authCodeInfo.AccessToken,
		CodeToken:    authCodeToken,
		OrgID:        authCodeInfo.OrgID,
		OrgName:      authCodeInfo.OrgName,
		OutOrgName:   authCodeInfo.OutOrgName,
		OrgCode:      authCodeInfo.OrgCode,
		UserID:       authCodeInfo.UserID,
		Name:         authCodeInfo.Name,
		Token:        authCodeInfo.Token,
	}

	_ = cache.SetEx(sconsts.CacheFsAuthCodeToken+authCodeToken, json.ToJsonIgnoreError(res), consts.CacheBaseExpire)

	return res, nil
}

func platformLoginByAuthCodeInfo(sourceChannel string, infoBo *bo.AuthCodeInfoBo) errs.SystemErrorInfo {
	log.Infof("[platformLoginByAuthCodeInfo]开始登录当前用户（没有则初始化）  authCodeInfo:%v", infoBo)

	// 封装用户信息和token
	baseUserInfo, err1 := domain.PlatformAuth(sourceChannel, infoBo.CorpId, infoBo.OutUserId, infoBo.AccessToken)
	if err1 != nil {
		log.Errorf("[platformLoginByAuthCodeInfo] PlatformAuth failed:%v, tenantKey:%s, openId:%s, accessToken:%s", err1, infoBo.CorpId, infoBo.OutUserId, infoBo.AccessToken)
		return err1
	}
	infoBo.Token = random.GenTokenWithOrgId(baseUserInfo.OrgId, baseUserInfo.UserId)
	infoBo.UserID = baseUserInfo.UserId
	infoBo.Name = baseUserInfo.Name

	//封装组织信息
	if baseUserInfo.OrgId != 0 {
		organizationBo, err := domain.GetOrgBoById(baseUserInfo.OrgId)
		if err != nil {
			log.Error(err)
		} else {
			infoBo.OrgCode = organizationBo.Code
			infoBo.OrgID = organizationBo.Id
			infoBo.OrgName = organizationBo.Name
			err := domain.UpdateGlobalUserLastLoginInfo(baseUserInfo.UserId, baseUserInfo.OrgId)
			if err != nil {
				log.Errorf("[platformLoginByAuthCodeInfo] UpdateGlobalUserLastLoginInfo err:%v", err)
			}
		}
	}

	// 设置登陆token
	cacheUserInfo := &bo.CacheUserInfoBo{
		OutUserId:     baseUserInfo.OutUserId,
		UserId:        baseUserInfo.UserId,
		CorpId:        infoBo.CorpId,
		SourceChannel: sourceChannel,
		OrgId:         baseUserInfo.OrgId,
	}
	cacheUserInfoJson, err2 := json.ToJson(cacheUserInfo)
	if err2 != nil {
		return errs.BuildSystemErrorInfo(errs.JSONConvertError, err2)
	}
	expireTime := consts.CacheUserTokenExpire
	if sourceChannel == sdk_const.SourceChannelDingTalk {
		expireTime = 3600 * 24 * 7
	}
	_ = cache.SetEx(sconsts.CacheUserToken+infoBo.Token, cacheUserInfoJson, int64(expireTime))

	//缓存refresh_token
	platform_sdk.GetPlatformInfo(sdk_const.SourceChannelFeishu).Cache.SetRefreshToken(baseUserInfo.OrgId, baseUserInfo.UserId, infoBo.RefreshToken, infoBo.AccessToken)

	//更新上次登录时间
	go func() {
		defer func() {
			if r := recover(); r != nil {
				log.Error(errs.BuildSystemErrorInfoWithPanicRecover(r, stack.GetStack()))
			}
		}()
		if err := domain.UserLoginHook(cacheUserInfo.UserId, cacheUserInfo.OrgId, sourceChannel); err != nil {
			log.Error(err)
		}
		// 授权登陆的时候获取到了头像，这个时候更新下
		if baseUserInfo.Avatar == "" && infoBo.Avatar != "" {
			if err := domain.UpdateUserInfo(baseUserInfo.UserId, mysql.Upd{
				consts.TcAvatar: infoBo.Avatar,
			}); err != nil {
				log.Errorf("[platformLoginByAuthCodeInfo] UpdateUserInfo userId:%v, avatar:%v, err:%v", baseUserInfo.UserId, infoBo.Avatar, err)
			}
			errCache := domain.ClearBaseUserInfo(baseUserInfo.OrgId, baseUserInfo.UserId)
			if errCache != nil {
				log.Errorf("[platformLoginByAuthCodeInfo] ClearBaseUserInfo userId:%v, avatar:%v, err:%v", baseUserInfo.UserId, infoBo.Avatar, errCache)
			}
		}
	}()

	return nil
}

func getPlatformUserInfoByAuthCode(sourceChannel, corpId, code string, codeType int) (*bo.AuthCodeInfoBo, errs.SystemErrorInfo) {
	client, sdkErr := platform_sdk.GetClient(sourceChannel, corpId)
	if sdkErr != nil {
		log.Errorf("[getPlatformUserInfoByAuthCode] platform_sdk.GetClient err: %v", sdkErr)
		return nil, errs.BuildSystemErrorInfo(errs.PlatFormOpenApiCallError, sdkErr)
	}

	// 根据authCode 获取相关信息
	authCodeInfo, err := getAuthCodeInfo(client, code, codeType)
	if err != nil {
		return nil, err
	}
	log.Infof("[getPlatformUserInfoByAuthCode] authCodeInfo corpId:%v", authCodeInfo.CorpId)
	//2022-2-19应最新需求，飞书创建组织改为默认安装时创建
	isInit, err := checkUserIsBindOrg(sourceChannel, authCodeInfo.CorpId, authCodeInfo.OutUserId)
	if err != nil {
		return nil, err
	}
	//组织未初始化，引导用户去创建组织（历史用户可能安装了应用，但没有初始化组织）
	authCodeInfo.Binding = isInit

	orgInfoResp, _ := client.GetOrgInfo(&sdkVo.OrgInfoReq{CorpId: authCodeInfo.CorpId})
	if orgInfoResp != nil {
		authCodeInfo.OutOrgName = orgInfoResp.Name
	}

	return authCodeInfo, nil
}

// 根据authCode获取信息
func getAuthCodeInfo(client sdk_interface.Sdk, code string, codeType int) (*bo.AuthCodeInfoBo, errs.SystemErrorInfo) {
	userAuthInfo, err := client.GetUserInfoByAuthCode(&sdkVo.GetUserInfoByAuthCodeReq{Code: code, CodeType: codeType})
	if err != nil {
		if strings.Contains(err.Message(), consts.DingAuthError) {
			return nil, errs.PlatFormAppUnauthorizedError
		}
		return nil, errs.BuildSystemErrorInfo(errs.PlatFormOpenApiCallError, err)
	}

	authCodeInfo := &bo.AuthCodeInfoBo{
		Binding:      true, //默认已经初始化组织，下面逻辑判断是否未初始化
		OutUserId:    userAuthInfo.OutUserId,
		CorpId:       userAuthInfo.CorpId,
		AccessToken:  userAuthInfo.AccessToken,
		RefreshToken: userAuthInfo.RefreshToken,
	}
	if userAuthInfo.UserInfo != nil {
		authCodeInfo.Name = userAuthInfo.UserInfo.Name
		authCodeInfo.IsAdmin = userAuthInfo.UserInfo.IsAdmin
		authCodeInfo.Avatar = userAuthInfo.UserInfo.Avatar
	}

	return authCodeInfo, nil
}

func checkUserIsBindOrg(sourceChannel, corpId, outUserId string) (bool, errs.SystemErrorInfo) {
	orgInfo, orgErr := domain.GetOrgByOutOrgId(corpId)
	if orgErr != nil {
		if orgErr == errs.OrgOutInfoNotExist {
			cacheInfo, _ := cache.Get(fmt.Sprintf("%s%s", sconsts.CacheFsOrgInit, corpId))
			if cacheInfo == corpId {
				//判断组织是否正在初始化
				return false, errs.OrgInitDoing
			} else {
				//组织未初始化，引导用户去创建组织（历史用户可能安装了应用，但没有初始化组织）
				return false, nil
			}
		} else {
			log.Errorf("[checkUserIsBindOrg] GetOrgByOutOrgId failed:%v, tenantKey:%s", orgErr, corpId)
			return false, orgErr
		}
	}

	_, userOutInfoErr := GetUserOutInfoByOpenID(orgInfo.Id, sourceChannel, outUserId)
	if userOutInfoErr != nil {
		if userOutInfoErr == errs.UserOutInfoNotExist {
			_, orgErr := domain.GetOrgByOutOrgId(corpId)
			if orgErr != nil {
				if orgErr == errs.OrgOutInfoNotExist {
					cacheInfo, _ := cache.Get(fmt.Sprintf("%s%s", sconsts.CacheFsOrgInit, corpId))
					if cacheInfo == corpId {
						//判断组织是否正在初始化
						return false, errs.OrgInitDoing
					} else {
						//组织未初始化，引导用户去创建组织（历史用户可能安装了应用，但没有初始化组织）
						return false, nil
					}
				} else {
					log.Errorf("[checkUserIsBindOrg] GetOrgByOutOrgId failed:%v, tenantKey:%s", orgErr, corpId)
					return false, orgErr
				}
			}
		} else {
			log.Errorf("[checkUserIsBindOrg] GetUserOutInfoByOpenID failed:%v, openId:%s", userOutInfoErr, outUserId)
			return false, userOutInfoErr
		}
	}

	return true, nil
}

func GetFsAccessToken(orgId, userId int64) (string, errs.SystemErrorInfo) {
	key, err := util.ParseCacheKey(sconsts.CacheFsUserAccessToken, map[string]interface{}{
		consts.CacheKeyOrgIdConstName:  orgId,
		consts.CacheKeyUserIdConstName: userId,
	})
	if err != nil {
		log.Error(err)
		return "", err
	}

	baseInfoJson, cacheErr := platform_sdk.GetPlatformInfo(sdk_const.SourceChannelFeishu).Cache.GetRefreshToken(orgId, userId)
	if cacheErr != nil {
		return "", errs.BuildSystemErrorInfo(errs.RedisOperateError, cacheErr)
	}
	baseInfo := &bo.CacheFsUserTokenBo{}
	if baseInfoJson != "" {
		err := json.FromJson(baseInfoJson, baseInfo)
		if err != nil {
			return "", errs.BuildSystemErrorInfo(errs.JSONConvertError)
		}
		//1小时更新一次access_token
		if time.Now().Unix()-baseInfo.FetchTime > 3600 {
			fsConfig := config.GetConfig().FeiShu
			if fsConfig == nil {
				log.Error("飞书配置为空")
				return "", errs.BuildSystemErrorInfo(errs.FeiShuConfigNotExistError)
			}
			appTicket, err := platform_sdk.GetPlatformInfo(sdk_const.SourceChannelFeishu).Cache.GetTicket()
			if err != nil {
				log.Errorf("[GetFsAccessToken] GetFeiShuCache err:%v", err)
				return "", errs.BuildSystemErrorInfo(errs.FeiShuAppTicketNotExistError)
			}
			resp, err1 := sdk.RefreshUserAccessToken(fsConfig.AppId, fsConfig.AppSecret, appTicket, baseInfo.RefreshToken)
			if err1 != nil {
				log.Error(err1)
				return "", errs.FeiShuOpenApiCallError
			}
			if resp.Code != 0 {
				if resp.Code == 20007 {
					// 这种情况下，需要用户重新扫码登陆
					// https://open.feishu.cn/document/ukTMukTMukTM/ugjM14COyUjL4ITN
					log.Errorf("[GetFsAccessToken] 需要用户重新扫码登陆。err: %v, accessToken: %s, refreshToken: %s",
						err1,
						baseInfo.AccessToken, baseInfo.RefreshToken)
					// 删除缓存
					cache.Del(key)
					// return "", errs.TokenExpires
					return "", errs.TokenAuthError
				}
				log.Error(resp.Msg)
				return "", errs.BuildSystemErrorInfo(errs.FeiShuOpenApiCallError)
			}
			platform_sdk.GetPlatformInfo(sdk_const.SourceChannelFeishu).Cache.SetRefreshToken(orgId, userId, resp.Data.RefreshToken, resp.Data.AccessToken)

			return resp.Data.AccessToken, nil
		} else {
			return baseInfo.AccessToken, nil
		}
	} else {
		//如果没有的话走重新登录逻辑
		return "", errs.TokenAuthError
	}
}

func PlatformDeptChange(orgId int64, eventType string, sourceChannel, deptOpenId string) (*vo.Void, errs.SystemErrorInfo) {
	orgInfo, orgInfoErr := domain.GetBaseOrgInfo(orgId)
	if orgInfoErr != nil {
		log.Error(orgInfoErr)
		return nil, orgInfoErr
	}

	id := int64(0)
	var err errs.SystemErrorInfo
	switch eventType {
	case consts.EventDeptAdd:
		id, err = domain.DeptAdd(orgId, deptOpenId, sourceChannel, orgInfo.OutOrgId)
	case consts.EventDeptDel:
		id, err = domain.DeptDelete(orgId, deptOpenId)
	case consts.EventDeptUpdate:
		id, err = domain.DeptUpdate(orgId, deptOpenId, sourceChannel, orgInfo.OutOrgId)
	}

	return &vo.Void{ID: id}, err
}

func ChangeDeptScope(orgId int64) errs.SystemErrorInfo {
	orgInfo, orgInfoErr := domain.GetBaseOrgInfo(orgId)
	if orgInfoErr != nil {
		log.Error(orgInfoErr)
		return orgInfoErr
	}

	return domain.ChangeDeptScope(orgId, orgInfo.SourceChannel, orgInfo.OutOrgId)
}

// BoundFeiShu 绑定飞书
// 判断是否无超管
// 判断团队是否已绑定外部平台
// 判断外部平台是否已被绑定
// 重置通讯录
func BoundFeiShu(userId, orgId int64, codeToken string, token string) (*vo.FeiShuAuthCodeResp, errs.SystemErrorInfo) {
	codeInfoJson, redisErr := cache.Get(sconsts.CacheFsAuthCodeToken + codeToken)
	if redisErr != nil {
		log.Error(redisErr)
		return nil, errs.RedisOperateError
	}
	if codeInfoJson == "" {
		return nil, errs.CodeTokenInvalid
	}
	codeInfo := vo.FeiShuAuthCodeResp{}
	_ = json.FromJson(codeInfoJson, &codeInfo)
	tenantKey := codeInfo.TenantKey

	//// 判断是否为超管
	//orgSysAdminIds, err := domain.GetOrgSysAdmin([]int64{orgId})
	//if err != nil {
	//	log.Error(err)
	//	return nil, err
	//}
	//suc, _ := slice.Contain(orgSysAdminIds[orgId], userId)
	//if !suc {
	//	log.Infof("绑定飞书组织失败，%d 非组织 %d 超管", userId, orgId)
	//	return nil, errs.BindFeiShuNoPermissionErr
	//}
	// 判断当前团队是否被其它团队绑定
	_, err := domain.GetOrgOutInfo(orgId)
	if err == nil {
		log.Infof("组织 %d 已绑定过外部平台", orgId)
		return nil, errs.OrgAlreadyBindPlatform
	}
	// 绑定飞书是不是被其它团队绑定过
	_, err = domain.GetOrgOutInfoByTenantKey(tenantKey)
	if err == nil {
		log.Infof("外部团队 %s 已被其它企业绑定", tenantKey)
		return nil, errs.TenantKeyAlreadyBindOrg
	}

	initOrgBo := bo.InitOrgBo{
		OutOrgId:      tenantKey,
		OutOrgOwnerId: codeInfo.OpenID,
		SourceChannel: sdk_const.SourceChannelFeishu,
	}

	openIds, err := domain.FilterBoundingOpenIds(orgId, sdk_const.SourceChannelFeishu, []string{codeInfo.OpenID}, nil)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	if len(openIds) > 0 {
		// 先为管理员绑定飞书外部信息
		err = domain.BoundFsOutInfo(userId, orgId, codeInfo.Name, codeInfo.OpenID)
		if err != nil {
			log.Error(err)
			return nil, err
		}
	}

	// 开始绑定企业
	// 部门架构重置，人员通过邮箱或手机号匹配，未匹配用户设置为未加入状态
	mysqlErr := mysql.TransX(func(tx sqlbuilder.Tx) error {
		// 清理当前企业之前的用户外部信息
		_, mysqlErr := mysql.TransUpdateSmartWithCond(tx, consts.TableUserOutInfo, db.Cond{
			consts.TcOrgId:     orgId,
			consts.TcOutUserId: db.NotEq(codeInfo.OpenID),
		}, mysql.Upd{
			consts.TcIsDelete: consts.AppIsDeleted,
		})
		if mysqlErr != nil {
			log.Error(mysqlErr)
			return errs.MysqlOperateError
		}
		//尝试获取企业信息，如果没有权限之类也不用报错
		tenantInfo, tenantInfoErr := feishu.GetTenantInfo(tenantKey)
		if tenantInfoErr != nil {
			log.Error(tenantInfoErr)
		} else {
			initOrgBo.OrgName = tenantInfo.Tenant.Name
			initOrgBo.TenantCode = tenantInfo.Tenant.DisplayId
		}
		_, err = domain.OrgOutInfoInit(initOrgBo, orgId, tx)
		if err != nil {
			log.Error(err)
			return err
		}
		deptMapping, err := domain.BoundFsDepartment(orgId, tenantKey, tx)
		if err != nil {
			log.Error(err)
			return err
		}
		return domain.BoundFsUser(orgId, userId, tenantKey, deptMapping, tx)
	})
	if mysqlErr != nil {
		log.Error(mysqlErr)
		return nil, errs.MysqlOperateError
	}
	//_ = SendFeishuMemberHelpMsg(codeInfo.TenantKey, orgId, codeInfo.OpenID)
	UpdateCacheUserFsInfo(token, codeInfo.TenantKey, codeInfo.OpenID, sdk_const.SourceChannelFeishu)
	return &codeInfo, nil
}

// BoundFeiShuAccount 绑定飞书账号
// 组织已绑定过飞书：当前绑定账号必须为该飞书组织下成员
// 组织未绑定过飞书：仅绑定外部信息，可扫码和免密登录
// 判断当前用户是否已绑定过外部账号
// 判断当前外部账号是否已被其他用户绑定过
func BoundFeiShuAccount(userId, orgId int64, codeToken string) (*vo.FeiShuAuthCodeResp, errs.SystemErrorInfo) {
	codeInfoJson, redisErr := cache.Get(sconsts.CacheFsAuthCodeToken + codeToken)
	if redisErr != nil {
		log.Error(redisErr)
		return nil, errs.RedisOperateError
	}
	if codeInfoJson == "" {
		return nil, errs.CodeTokenInvalid
	}
	codeInfo := vo.FeiShuAuthCodeResp{}
	_ = json.FromJson(codeInfoJson, &codeInfo)
	tenantKey := codeInfo.TenantKey
	openId := codeInfo.OpenID

	// 判断当前组织是否绑定的该团队
	orgOutInfo, err := domain.GetOrgOutInfo(orgId)
	if err != nil && err != errs.OrgOutInfoNotExist {
		log.Error(err)
		return nil, err
	}
	if orgOutInfo != nil && orgOutInfo.OutOrgId != tenantKey {
		return nil, errs.OrgAlreadyBindPlatform
	}
	// 区分三方登录和组织绑定两个场景
	platformOrgId := int64(0)
	if orgOutInfo != nil {
		platformOrgId = orgOutInfo.OrgId
	}
	// 判断用户是否绑定过外部信息
	userOutInfos, err := GetUserOutInfoByPlatformAndOrgId(sdk_const.SourceChannelFeishu, userId, platformOrgId)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	if len(userOutInfos) > 0 {
		return nil, errs.UserAlreadyBindPlatform
	}
	// 判断外部信息是不是被其他人绑定过
	userOutInfo, _ := GetUserOutInfoByOpenID(orgId, sdk_const.SourceChannelFeishu, openId)
	if userOutInfo != nil {
		return nil, errs.OutUserAlreadyBind
	}
	if orgOutInfo != nil {
		// 绑定过组织
		return &codeInfo, domain.BoundFsAccount(userId, platformOrgId, tenantKey, openId)
	}
	return &codeInfo, domain.BoundFsOutInfo(userId, platformOrgId, codeInfo.Name, openId)
}

// InitFeiShuAccount 初始化飞书账户
func InitFeiShuAccount(codeToken string) (*vo.FeiShuAuthCodeResp, errs.SystemErrorInfo) {
	codeInfoJson, redisErr := cache.Get(sconsts.CacheFsAuthCodeToken + codeToken)
	if redisErr != nil {
		log.Error(redisErr)
		return nil, errs.RedisOperateError
	}
	if codeInfoJson == "" {
		return nil, errs.CodeTokenInvalid
	}
	codeInfo := vo.FeiShuAuthCodeResp{}
	_ = json.FromJson(codeInfoJson, &codeInfo)

	openId := codeInfo.OpenID
	userOutInfo, _ := GetUserOutInfoByOpenID(0, sdk_const.SourceChannelFeishu, openId)
	if userOutInfo != nil {
		return nil, errs.OutUserAlreadyBind
	}
	// 如果飞书已被绑定过，使用该orgId
	orgOutInfo, _ := domain.GetOrgOutInfoByTenantKey(codeInfo.TenantKey)
	orgId := int64(0)
	if orgOutInfo != nil {
		orgId = orgOutInfo.OrgId
	}
	// orgId 为 0 是正常的流程，前端先调用该接口创建用户，再进行创建组织并绑定。

	// 新建用户
	userId, err := idfacade.ApplyPrimaryIdRelaxed(consts.TableUser)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	userPo := domain.AssemblyFeiShuUserInfo(orgId, userId, &sdkVo.ScopeUser{Name: codeInfo.Name})
	mysqlErr := mysql.Insert(&userPo)
	if mysqlErr != nil {
		log.Error(mysqlErr)
		return nil, errs.MysqlOperateError
	}

	// 新建用户外部信息
	err = domain.BoundFsOutInfo(userId, orgId, codeInfo.Name, openId)
	if err != nil {
		log.Error(err)
		return nil, err
	}

	// 新建token
	codeInfo.Token = random.GenTokenWithOrgId(orgId, userId)
	cacheUserInfo := &bo.CacheUserInfoBo{
		OutUserId:     openId,
		UserId:        userId,
		CorpId:        codeInfo.TenantKey,
		SourceChannel: sdk_const.SourceChannelFeishu,
		OrgId:         0,
	}
	cacheUserInfoJson, err2 := json.ToJson(cacheUserInfo)
	if err2 != nil {
		return nil, errs.BuildSystemErrorInfo(errs.JSONConvertError, err2)
	}
	_ = cache.SetEx(sconsts.CacheUserToken+codeInfo.Token, cacheUserInfoJson, consts.CacheUserTokenExpire)
	return &codeInfo, nil
}

func GetJsAPITicket(orgId int64, configType ...string) (string, errs.SystemErrorInfo) {
	params := map[string]interface{}{
		consts.CacheKeyOrgIdConstName:        orgId,
		consts.CacheKeyJsapiTicketConfigType: orgId,
	}
	if len(configType) > 0 && configType[0] != "" {
		params[consts.CacheKeyJsapiTicketConfigType] = configType[0]
	}
	key, err := util.ParseCacheKey(sconsts.CacheOrgJSTicket, params)
	if err != nil {
		log.Error(err)
		return "", err
	}
	ticket, ticketErr := GetTicketFromCache(orgId, key)
	if ticketErr != nil || ticket == "" {
		//防止缓存击穿
		uuid := uuid.NewUuid()
		suc, err := cache.TryGetDistributedLock(consts.FeiShuGetJSApiTicketLockKey, uuid)
		if err != nil {
			log.Error("获取锁异常")
			return "", errs.BuildSystemErrorInfo(errs.TryDistributedLockError)
		}
		if suc {
			//释放锁
			defer func() {
				if _, e := cache.ReleaseDistributedLock(consts.FeiShuGetJSApiTicketLockKey, uuid); e != nil {
					log.Error(e)
				}
			}()

			ticket, ticketErr := GetTicketFromCache(orgId, key)
			if ticketErr != nil || ticket == "" {
				orgOutInfo, err := domain.GetOrgOutInfo(orgId)
				if err != nil {
					// 如果不是三方来源组织，GetJsAPITicket 返回空字符串
					if err.Code() == errs.OrgOutInfoNotExist.Code() {
						return "", nil
					}
					log.Errorf("[GetJsAPITicket]GetOrgOutInfo failed:%v, orgId:%d", err, orgId)
					return "", err
				}
				client, sdkErr := platform_sdk.GetClient(orgOutInfo.SourceChannel, orgOutInfo.OutOrgId)
				if sdkErr != nil {
					log.Errorf("[GetJsAPITicket] GetClient failed:%v, corpId:%s", sdkErr, orgOutInfo.OutOrgId)
					return "", errs.PlatFormOpenApiCallError
				}
				jsConfigType := ""
				if len(configType) > 0 {
					jsConfigType = configType[0]
				}
				resp, sdkErr := client.GetJsTicket(jsConfigType)
				if sdkErr != nil {
					log.Errorf("[GetJsAPITicket]GetJsTicket failed:%v", sdkErr)
					return "", errs.PlatFormOpenApiCallError
				}

				//本身过期时间是2小时，提前10分钟重新获取新的
				cacheExpire := resp.ExpireIn
				if cacheExpire > 600 {
					cacheExpire = cacheExpire - 600
				}
				cacheErr := cache.SetEx(key, resp.Ticket, int64(cacheExpire))
				if cacheErr != nil {
					log.Error(cacheErr)
				}
				return resp.Ticket, nil
			} else {
				return ticket, nil
			}
		} else {
			return "", errs.SystemBusy
		}
	}

	return ticket, nil
}

func GetTicketFromCache(orgId int64, key string) (string, errs.SystemErrorInfo) {
	ticket, ticketErr := cache.Get(key)
	if ticketErr != nil {
		log.Errorf("[GetTicketFromCache]CacheGet failed:%v", ticketErr)
		return "", errs.BuildSystemErrorInfo(errs.RedisOperateError, ticketErr)
	}

	return ticket, nil
}
