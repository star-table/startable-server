package callsvc

import (
	sdk_const "gitea.bjx.cloud/allstar/platform-sdk/consts"
	"github.com/star-table/startable-server/common/core/errs"
	"github.com/star-table/startable-server/common/core/util/json"
	"github.com/star-table/startable-server/common/model/vo/orgvo"
)

type ChangeAuthHandler struct {
	base.CallBackBase
}

type ChangeAuthCallbackMsg struct {
	SyncAction   string       `json:"syncAction"`
	AuthCorpInfo AuthCorpInfo `json:"auth_corp_info"`
}

type AuthCorpInfo struct {
	CorpId   string `json:"corpid"`
	CorpName string `json:"corp_name"`
}

func (c ChangeAuthHandler) Handle(data req_vo.DingEventBizData) error {
	log.Infof("[ChangeAuthHandler] handle data:%v", data)
	msg := &ChangeAuthCallbackMsg{}
	_ = json.FromJson(data.BizData, msg)

	outOrgResp := orgfacade.GetOrgOutInfoByOutOrgId(orgvo.GetOutOrgInfoByOutOrgIdReqVo{
		OutOrgId: data.CorpId,
	})
	if outOrgResp.Failure() {
		if outOrgResp.Code == errs.OrgOutInfoNotExist.Code() {
			orgName := msg.AuthCorpInfo.CorpName
			if orgName == "" {
				return nil
			}
			corpId := data.CorpId
			return InitDingOrg(corpId, orgName)
		} else {
			log.Errorf("[ChangeAuthHandler] GetOrgOutInfoByOutOrgId err:%v", outOrgResp.Error())
			return outOrgResp.Error()
		}
	}

	return c.ChangeAuthScope(data.CorpId, sdk_const.SourceChannelDingTalk)
}
