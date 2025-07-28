package consts

import (
	"math/rand"
	"time"

	sdk_const "gitea.bjx.cloud/allstar/platform-sdk/consts"
)

const CacheKeyPrefix = "polaris:"
const CacheKeyOfSys = "sys:"
const CacheKeyOfOrg = "org_{{." + CacheKeyOrgIdConstName + "}}:"
const CacheKeyOfOutOrg = "outorg_{{." + CacheKeyOutOrgIdConstName + "}}:"
const CacheKeyOfUser = "user_{{." + CacheKeyUserIdConstName + "}}:"
const CacheKeyOfOutUser = "outuser_{{." + CacheKeyOutUserIdConstName + "}}:"
const CacheKeyOfProject = "project_{{." + CacheKeyProjectIdConstName + "}}:"
const CacheKeyOfProcess = "process_{{." + CacheKeyProcessIdConstName + "}}:"
const CacheKeyOfRole = "role_{{." + CacheKeyRoleIdConstName + "}}:"
const CacheKeyOfSourceChannel = "source_channel_{{." + CacheKeySourceChannelConstName + "}}:"
const CacheKeyOfPhone = "phone_{{." + CacheKeyPhoneConstName + "}}:"
const CacheKeyOfAuthType = "authType_{{." + CacheKeyAuthTypeConstName + "}}:"
const CacheKeyOfAddressType = "addressType_{{." + CacheKeyAddressTypeConstName + "}}:"
const CacheKeyOfLoginName = "login_name_{{." + CacheKeyLoginNameConstName + "}}:"
const CacheKeyOfRoleGroup = "group_{{." + CacheKeyRoleGroupConstName + "}}"
const CacheKeyOfDepartment = "user_{{." + CacheKeyDepartmentIdConstName + "}}:"
const CacheKeyOfIssue = "user_{{." + CacheKeyIssueIdConstName + "}}:"
const CacheKeyOfSchedule = CacheKeyPrefix + "schedule:"
const CacheKeyOfApp = "app_{{." + CacheKeyAppIdConstName + "}}:"
const CacheKeyOfAsyncTask = "asyncTask:"
const CacheKeyOfAsyncTaskId = "id_{{." + CacheKeyAsyncTaskIdConstName + "}}"
const CacheKeyOfJsapiTicket = "config_{{." + CacheKeyJsapiTicketConfigType + "}}:"

//服务名
const (
	AppsvcApplicationName      = "appsvc:"
	IdsvcApplicationName       = "idsvc:"
	MsgsvcApplicationName      = "msgsvc:"
	CallsvcApplicationName     = "callsvc:"
	OrgsvcApplicationName      = "orgsvc:"
	ProcesssvcApplicationName  = "processsvc:"
	ProjectsvcApplicationName  = "projectsvc:"
	ResourcesvcApplicationName = "resourcesvc:"
	RolesvcApplicationName     = "rolesvc:"
	RrendssvcApplicationName   = "trendssvc:"
	SchedulesvcApplicationName = "scheduletsvc:"
	CommonsvcApplicationName   = "commonsvc"
	OrderApplicationName       = "ordersvc:"
)

//失效时间
const (
	//用户Token失效时间: 15天
	CacheUserTokenExpire = 60 * 60 * 24 * 15
	//通用失效时间: 3小时
	CacheBaseExpire = int64(60 * 60 * 3)
	//用户信息缓存: 1小时
	CacheBaseUserInfoExpire = int64(60 * 60 * 1)
	CacheExpire1Day         = int64(60 * 60 * 24)
)

func GetCacheBaseExpire() int64 {
	rand.Seed(time.Now().Unix())
	return CacheBaseExpire + rand.Int63n(30)
}
func GetCacheBaseUserInfoExpire() int64 {
	rand.Seed(time.Now().Unix())
	return CacheBaseUserInfoExpire + rand.Int63n(30)
}

const (
	CacheKeyOrgIdConstName         = "orgId"
	CacheKeyUserIdConstName        = "userId"
	CacheKeyAppIdConstName         = "appId"
	CacheKeyOutOrgIdConstName      = "outOrgId"
	CacheKeyOutUserIdConstName     = "outUserId"
	CacheKeyProjectIdConstName     = "projectId"
	CacheKeyIssueIdConstName       = "issueId"
	CacheKeyObjectCodeConstName    = "objectCode"
	CacheKeyProcessIdConstName     = "processId"
	CacheKeyRoleIdConstName        = "roleId"
	CacheKeySourceChannelConstName = "sourceChannel"
	CacheKeyYearConstName          = "year"
	CacheKeyMonthConstName         = "month"
	CacheKeyDayConstName           = "day"
	CacheKeyPhoneConstName         = "phone"
	CacheKeyAuthTypeConstName      = "authType"
	CacheKeyAddressTypeConstName   = "addressType"
	CacheKeyLoginNameConstName     = "loginName"
	CacheKeyRoleGroupConstName     = "roleGroup"
	CacheKeyDepartmentIdConstName  = "departmentId"
	CacheKeyAsyncTaskIdConstName   = "asyncTaskId"
	CacheKeyJsapiTicketConfigType  = "configType"
)

//系统缓存
const (
	//DingTalk Suite Ticket
	CacheDingTalkSuiteTicket     = CacheKeyPrefix + CacheKeyOfSys + sdk_const.SourceChannelDingTalk + ":suite_ticket"
	CacheDingTalkCorpAccessToken = CacheKeyPrefix + CacheKeyOfSys + sdk_const.SourceChannelDingTalk + ":corp_access_token:"
	CacheDingTalkAppAccessToken  = CacheKeyPrefix + CacheKeyOfSys + sdk_const.SourceChannelDingTalk + ":app_access_token:"

	//weixin
	CacheWeixinSuiteTicket     = CacheKeyPrefix + CacheKeyOfSys + sdk_const.SourceChannelWeixin + ":suite_ticket"
	CacheWeixinCorpAccessToken = CacheKeyPrefix + CacheKeyOfSys + sdk_const.SourceChannelWeixin + ":corp_access_token:"
	CacheWeixinAppAccessToken  = CacheKeyPrefix + CacheKeyOfSys + sdk_const.SourceChannelWeixin + ":app_access_token:"

	//飞书 AppTicket
	CacheFeiShuAppTicket = CacheKeyPrefix + CacheKeyOfSys + sdk_const.SourceChannelFeishu + ":app_ticket"
	//飞书 AppAccessToken
	CacheFeiShuAppAccessToken = CacheKeyPrefix + CacheKeyOfSys + sdk_const.SourceChannelFeishu + ":app_access_token:"
	//飞书 TenantAccessToken
	CacheFeiShuTenantAccessToken = CacheKeyPrefix + CacheKeyOfSys + sdk_const.SourceChannelFeishu + ":tenant_access_token:"
	//飞书 授权范围
	CacheFeiShuScope = CacheKeyPrefix + CacheKeyOfSys + sdk_const.SourceChannelFeishu + ":scope:"

	//fs用户refresh_token和user_access_token
	CacheFsUserAccessToken = CacheKeyPrefix + OrgsvcApplicationName + CacheKeyOfOrg + CacheKeyOfUser + "fstoken"

	//mqtt root key
	CacheMQTTRootKey = CacheKeyPrefix + CacheKeyOfSys + "mqtt:root_key"

	//飞书 卡片通知回调消息refresh-token, 网络抖动等极端情况下，会出现卡片点击失败但是业务方已经处理过 action 的现象，所以业务方接口存在被重复调用的风险。X-Refresh-Token 只有在卡片点击事件成功被响应并通知到客户端的时候才会刷新，如果业务方的接口非幂等，可以通过缓存并验证该字段防止接口被重复调用。
	CacheFeiShuCardCallBackRefreshToken = CacheKeyPrefix + CacheKeyOfSys + sdk_const.SourceChannelFeishu + ":card_call_back:refresh_token:"

	//用户token
	//CacheUserToken = CacheKeyPrefix + CacheKeyOfSys + "user:token:"
	////对象id缓存key前缀
	//CacheObjectIdPreKey = CacheKeyPrefix + CacheKeyOfSys + "object_id:"
	// 角色操作列表
	//CacheRoleOperationList = CacheKeyPrefix + CacheKeyOfSys + "role_operation_list"
	////部门对应关系
	//CacheDeptRelation = CacheKeyPrefix + CacheKeyOfSys + "dept_relation_list"

	//屏蔽部分飞书组织同步用户信息
	CacheFeishuNotSyncUser = CacheKeyPrefix + CacheKeyOfSys + sdk_const.SourceChannelFeishu + ":not_sync_user_tenant"

	//灰度企业名单
	CacheGrayFeishuOrg = CacheKeyPrefix + CacheKeyOfSys + sdk_const.SourceChannelFeishu + ":gray_org_list"
	// 异步任务进度信息
	CacheKeyOfAsyncTaskInfo = CacheKeyPrefix + OrgsvcApplicationName + CacheKeyOfOrg + CacheKeyOfAsyncTask + CacheKeyOfAsyncTaskId

	// 极星付费范围
	CachePayRangeInfo = CacheKeyPrefix + CacheKeyOfSys + "pay_range_info"

	CacheUrgeIssueAudit   = "urge_issue_audit"
	CacheUrgeIssue        = "urge_issue"
	CacheDeletedIssueNums = "deleted_issue_nums"
	CacheProjectTypeList  = "project_type_list"
	CacheFsPushSettings   = "fs_push_settings"
)
