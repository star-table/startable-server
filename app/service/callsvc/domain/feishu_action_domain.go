package callsvc

import (
	"regexp"
	"strings"

	"github.com/star-table/startable-server/common/core/consts"
	"github.com/star-table/startable-server/common/model/bo"

	"github.com/star-table/startable-server/app/facade/orgfacade"
	"github.com/star-table/startable-server/app/facade/projectfacade"
	"github.com/star-table/startable-server/common/core/errs"
	"github.com/star-table/startable-server/common/core/util/json"
	commonFeiShu "github.com/star-table/startable-server/common/extra/feishu"
	"github.com/star-table/startable-server/common/model/vo/callvo"
	"github.com/star-table/startable-server/common/model/vo/orgvo"
	"github.com/star-table/startable-server/common/model/vo/projectvo"
)

// GetGroupChatHandleInfos 处理群聊指令时用到的一些数据在这里组装
func GetGroupChatHandleInfos(tenantKey string, event *callvo.MessageReqData, sourceChannel string) (*callvo.GroupChatHandleInfo, errs.SystemErrorInfo) {
	info := &callvo.GroupChatHandleInfo{SourceChannel: sourceChannel}
	orgResp := orgfacade.GetBaseOrgInfoByOutOrgId(orgvo.GetBaseOrgInfoByOutOrgIdReqVo{OutOrgId: tenantKey})
	if orgResp.Failure() {
		log.Errorf("[GetGroupChatHandleInfos] GetBaseOrgInfoByOutOrgId err: %v, tenantKey: %v", orgResp.Error(), tenantKey)
		return info, orgResp.Error()
	}
	info.OrgInfo = *orgResp.BaseOrgInfo

	tenantClient, err := commonFeiShu.GetTenant(tenantKey)
	if err != nil {
		log.Errorf("[GetGroupChatHandleInfos] GetTenant err: %v, tenantKey: %s", err, tenantKey)
		return info, err
	}
	info.TenantClient = tenantClient

	// 查询绑定的项目ids
	bindProResp := projectfacade.GetProjectIdsByChatId(projectvo.GetProjectIdsByChatIdReqVo{
		OrgId:      info.OrgInfo.OrgId,
		OpenChatId: event.OpenChatId,
	})
	if bindProResp.Failure() {
		log.Errorf("[GetGroupChatHandleInfos] GetProjectIdsByChatId err: %v, tenantKey: %v", bindProResp.Error(), tenantKey)
		return info, bindProResp.Error()
	}
	info.BindProjectIds = bindProResp.Data.ProjectIds

	// 查询操作人的信息
	if event.OpenId != "" {
		opUserResp := orgfacade.GetBaseUserInfoByEmpId(orgvo.GetBaseUserInfoByEmpIdReqVo{
			OrgId: info.OrgInfo.OrgId,
			EmpId: event.OpenId,
		})
		if opUserResp.Failure() {
			log.Errorf("[GetGroupChatHandleInfos] err: %v, operate user openId: %s", opUserResp.Error(), event.OpenId)
			return info, errs.BuildSystemErrorInfo(errs.UserNotFoundError, opUserResp.Error())
		}
		info.OperateUser = *opUserResp.BaseUserInfo
	}
	info.OperatorOpenId = event.OpenId

	// 检查是否是任务群聊
	if len(info.BindProjectIds) > 0 {
		info.IsIssueChat = false
	} else {
		// 查询任务群聊对应的任务，如果没查到，则不是任务群聊
		issues, err := GetIssueByChatId(info.OrgInfo.OrgId, event.OpenChatId, true)
		if err != nil {
			if err == errs.IssueNotExist {
				info.IsIssueChat = false
			} else {
				log.Errorf("[GetGroupChatHandleInfos] err: %v, tenantKey: %s", err, tenantKey)
				return info, err
			}
		} else {
			if len(issues) > 0 {
				info.IsIssueChat = true
				info.IssueInfo = issues[0]

				issueLinks := projectfacade.GetIssueLinks(projectvo.GetIssueLinksReqVo{
					SourceChannel: sourceChannel,
					OrgId:         info.OrgInfo.OrgId,
					IssueId:       info.IssueInfo.Id,
				}).Data
				info.IssueInfoUrl = issueLinks.SideBarLink
				info.IssuePcUrl = issueLinks.Link
			}
		}
	}

	return info, nil
}

// GetOrgSummaryAppId 获取某个组织的汇总表 appId
func GetOrgSummaryAppId(orgId int64) (int64, errs.SystemErrorInfo) {
	orgResp := orgfacade.GetOrgBoListByPage(orgvo.GetOrgIdListByPageReqVo{
		Page: 1,
		Size: 10000, // 不要超过1w
		Input: orgvo.GetOrgIdListByPageReqVoData{
			OrgIds: []int64{orgId},
		},
	})
	if orgResp.Failure() {
		log.Errorf("[GetOrgSummaryAppId] err: %v", orgResp.Error())
		return 0, orgResp.Error()
	}
	if len(orgResp.Data.List) < 1 {
		return 0, nil
	}
	org := orgResp.Data.List[0]
	orgRemarkObj := &orgvo.OrgRemarkConfigType{}
	if len(org.Remark) > 0 {
		oriErr := json.FromJson(org.Remark, orgRemarkObj)
		if oriErr != nil {
			log.Errorf("[GetOrgSummaryAppId] 组织 remark 反序列化异常，组织id:%d,原因:%v", org.Id, oriErr)
			return 0, errs.BuildSystemErrorInfo(errs.JSONConvertError, oriErr)
		}
	}

	return orgRemarkObj.OrgSummaryTableAppId, nil
}

func CheckIssueIsDelete(issue *bo.IssueBo) bool {
	if issue.Id <= 0 {
		return true
	}
	if issue.IsDelete == consts.AppIsDeleted {
		return true
	}
	flag, ok := issue.LessData[consts.BasicFieldRecycleFlag]
	if !ok {
		return false
	}
	flagInt := int(flag.(float64))
	if flagInt == 1 {
		return true
	}
	flagInt = int(flag.(float64))
	if flagInt == 1 {
		return true
	}

	return false
}

func CheckIsUserName(text string) bool {
	reg1 := regexp.MustCompile(`^(<at open_id="\w+">[^<]+</at>)$`)
	matchRes := reg1.FindString(text)
	if len(matchRes) < 1 {
		return false
	}
	return true
}

func CheckIsUserNameWithIssueTitle(text string) bool {
	_, _, title := MatchUserInsUserNameWithIssueTitle(text)
	if len(title) < 1 {
		return false
	}
	return true
}

// 根据用户输入的指令，匹配出`at`的目标用户信息。
func MatchUserInsUserInfo(text string) (openId, userName string) {
	reg1 := regexp.MustCompile(`<at open_id="(\w+)">([^<]+)</at>`)
	matchList := reg1.FindStringSubmatch(text)
	if len(matchList) >= 2 {
		openId = matchList[1]
		userName = matchList[2]
	}
	return
}

// 根据用户输入的指令，匹配出`at`的目标用户信息，以及任务标题
func MatchUserInsUserNameWithIssueTitle(text string) (atUserName, atUserOpenId, issueTitle string) {
	reg1 := regexp.MustCompile(`([^\s]+)*\s*<at open_id="(\w+)">([^<]+)</at>\s*([^\n]+)?`)
	matchList := reg1.FindStringSubmatch(text)
	if len(matchList) >= 3 {
		titlePart1 := matchList[1]
		atUserOpenId = matchList[2]
		atUserName = matchList[3]
		titlePart2 := matchList[4]
		// 去除标签字符串。如：`<at user_id="xx">@Name1</at>`，则只保留：`@Name1`
		issueTitle = RemoveAtHtmlTag(titlePart1 + titlePart2)
		issueTitle = RemoveHtmlTag(issueTitle)
		issueTitle = RemoveHtmlEntity(issueTitle)
		issueTitle = strings.Trim(issueTitle, " ")
	}
	return
}

// CheckIsIssueChat 检查是否是任务的群聊
func CheckIsIssueChat(orgId int64, chatId string) (bool, errs.SystemErrorInfo) {
	isIssueChat := false
	issues, err := GetIssueByChatId(orgId, chatId, false)
	if err != nil {
		if err == errs.IssueNotExist {
			isIssueChat = false
		} else {
			log.Errorf("[CheckIsIssueChat] err: %v, chatId: %s", err, chatId)
			return false, err
		}
	}
	if len(issues) > 0 {
		return true, nil
	}

	return isIssueChat, nil
}

// 去除 <at>xxx</at> 标签
func RemoveAtHtmlTag(s string) string {
	re, _ := regexp.Compile(`<\/?at[^>]*>`)
	s = re.ReplaceAllString(s, "")
	return s
}

// 去除 html 标签
func RemoveHtmlTag(s string) string {
	re, _ := regexp.Compile(`<\/?[a-zA-Z]*[\s\w\'\"=]*>`)
	s = re.ReplaceAllString(s, "")
	// 对单个尖括号的标签处理：如果匹配上了 `<[a-zA-Z]{1,}`，则将 < 过滤掉
	re, _ = regexp.Compile(`<([a-zA-Z]{1,})`)
	matchRes := re.FindString(s)
	if len(matchRes) > 0 {
		s = re.ReplaceAllString(s, `$1`)
	}
	return s
}

// 去除 html 实体字符
func RemoveHtmlEntity(s string) string {
	re, _ := regexp.Compile(`\&(&#32|&#34|nbsp|amp|&#39|&#43|&lt|&gt);`)
	s = re.ReplaceAllString(s, "")
	return s
}
