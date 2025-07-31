package callsvc

import (
	"gitea.bjx.cloud/allstar/dingtalk-sdk-golang"
	"gitea.bjx.cloud/allstar/dingtalk-sdk-golang/response"
	platform_sdk "gitea.bjx.cloud/allstar/platform-sdk"
	sdk_const "gitea.bjx.cloud/allstar/platform-sdk/consts"
	"github.com/star-table/startable-server/app/facade/orgfacade"
	"github.com/star-table/startable-server/common/core/config"
	"github.com/star-table/startable-server/common/core/consts"
	"github.com/star-table/startable-server/common/core/errs"
	"github.com/star-table/startable-server/common/core/logger"
	"github.com/star-table/startable-server/common/core/util/json"
	"github.com/star-table/startable-server/common/core/util/strs"
	"github.com/star-table/startable-server/common/core/util/uuid"
	"github.com/star-table/startable-server/common/library/cache"
	"github.com/star-table/startable-server/common/model/bo"
	"github.com/star-table/startable-server/common/model/vo/orgvo"
)

type TempAuthCodeHandler struct{}

type OrgSuitAuthCallbackMsg struct {
	SyncAction   string       `json:"syncAction"`
	AuthCorpInfo AuthCorpInfo `json:"auth_corp_info"`
}

var log = logger.GetDefaultLogger()

func (t TempAuthCodeHandler) Handle(data req_vo.DingEventBizData) error {
	log.Infof("ding TempAuthCodeHandler data:%v", data)

	msg := &OrgSuitAuthCallbackMsg{}
	_ = json.FromJson(data.BizData, msg)

	orgName := msg.AuthCorpInfo.CorpName
	corpId := data.CorpId
	return InitDingOrg(corpId, orgName)
}

func (t TempAuthCodeHandler) activeApp(corpId, authCode string) (*response.CorpPermanentCode, error) {
	ticket, err := platform_sdk.GetPlatformInfo(sdk_const.SourceChannelDingTalk).Cache.GetTicket()
	if err != nil {
		log.Errorf("[activeApp] GetTicket err:%v", err)
		return nil, err
	}

	dingConfig := config.GetConfig().DingTalk
	dingClient, err := dingtalk.NewClient(dingConfig.SuiteKey, dingConfig.SuiteSecret,
		dingtalk.WithTicket(ticket), dingtalk.WithCorpId(corpId))
	if err != nil {
		log.Errorf("[GetClient] err:%v", err)
		return nil, errs.DingTalkClientError
	}
	//获取永久授权码
	permanentCodeResp, err := dingClient.GetCorpPermanentCode(authCode)
	if err != nil {
		log.Errorf("[GetCorpPermanentCode] err:%v", err)
		return nil, err
	}

	permanentCode := permanentCodeResp.PermanentCode
	log.Infof("permanentCode : %v", permanentCodeResp)

	//激活应用
	activateSuiteResp, err := dingClient.ActivateSuite(corpId, permanentCode)
	if err != nil {
		log.Error("[activeApp] ActivateSuite:" + strs.ObjectToString(err))
		return nil, err
	}
	log.Info("激活状态：" + json.ToJsonIgnoreError(activateSuiteResp))

	//注册回调
	//resp, err := dingClient.RegisterEvent(&request.RegisterEvent{
	//	Tags: []string{"user_add_org", "user_modify_org", "user_leave_org", "user_active_org",
	//		"org_dept_create", "org_dept_modify", "org_dept_remove", "org_remove"},
	//	Token:  dingConfig.Token,
	//	Url:    "https://apifuse.bjx.cloud/api/callsvc/callback/dingtalk",
	//	Secret: dingConfig.AesKey,
	//})
	//if err != nil {
	//	log.Error("[activeAndRegisterCallback] RegisterEvent:" + strs.ObjectToString(err))
	//	return nil, err
	//}
	//log.Info("注册回调：" + json.ToJsonIgnoreError(resp))

	return &permanentCodeResp, nil
}

func InitDingOrg(corpId, corpName string) error {
	uid := uuid.NewUuid()
	lockKey := consts.DingTalkCorpInitKey + corpId
	suc, err := cache.TryGetDistributedLock(lockKey, uid)
	if err != nil {
		log.Errorf("企业%s初始化时，获取分布式锁异常:%v", corpId, err)
		return err
	}
	if suc {
		defer func() {
			_, e := cache.ReleaseDistributedLock(lockKey, uid)
			if e != nil {
				log.Error(e)
			}
		}()
	} else {
		log.Infof("企业%s正在初始化中，当前请求无效", corpId)
		return nil
	}
	//开始初始化组织
	orgInitResp := orgfacade.InitOrg(orgvo.InitOrgReqVo{
		InitOrg: bo.InitOrgBo{
			OutOrgId:      corpId,
			OrgName:       corpName,
			SourceChannel: sdk_const.SourceChannelDingTalk,
			//PermanentCode: permanentCodeResp.PermanentCode,
		},
	})
	if orgInitResp.Failure() {
		log.Error(orgInitResp.Message)
		return orgInitResp.Error()
	}
	return nil
}
