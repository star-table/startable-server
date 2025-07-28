package domain

import (
	"fmt"
	"strconv"

	"github.com/star-table/startable-server/app/facade/tablefacade"

	tablePb "gitea.bjx.cloud/LessCode/interface/golang/table/v1"
	fsvo "gitea.bjx.cloud/allstar/feishu-sdk-golang/core/model/vo"
	"gitea.bjx.cloud/allstar/feishu-sdk-golang/sdk"
	"github.com/star-table/startable-server/app/facade/formfacade"
	"github.com/star-table/startable-server/app/facade/orgfacade"
	"github.com/star-table/startable-server/app/facade/userfacade"
	"github.com/star-table/startable-server/common/core/businees"
	"github.com/star-table/startable-server/common/core/consts"
	"github.com/star-table/startable-server/common/core/errs"
	"github.com/star-table/startable-server/common/core/util"
	"github.com/star-table/startable-server/common/core/util/asyn"
	"github.com/star-table/startable-server/common/core/util/copyer"
	"github.com/star-table/startable-server/common/core/util/json"
	"github.com/star-table/startable-server/common/core/util/slice"
	"github.com/star-table/startable-server/common/extra/card"
	"github.com/star-table/startable-server/common/extra/feishu"
	"github.com/star-table/startable-server/common/library/cache"
	"github.com/star-table/startable-server/common/model/bo"
	"github.com/star-table/startable-server/common/model/vo/commonvo"
	"github.com/star-table/startable-server/common/model/vo/formvo"
	"github.com/star-table/startable-server/common/model/vo/lc_table"
	"github.com/star-table/startable-server/common/model/vo/orgvo"
	"github.com/star-table/startable-server/common/model/vo/projectvo"
	"github.com/spf13/cast"
)

// IssueChat 任务群聊
// 目前只有飞书的实现，后续如果需要接入钉钉，则只需实现该接口。
type IssueChat interface {
	// StartIssueChat 发起讨论
	StartIssueChat() (string, errs.SystemErrorInfo)
	NewIssueChat() (string, errs.SystemErrorInfo)
	// DeleteMembers 删除讨论群的成员
	DeleteMembers(delMemberIds []int64) errs.SystemErrorInfo
	// AddMembers 新增讨论群的成员
	AddMembers(addMemberIds []int64) errs.SystemErrorInfo
	// SendTopicCard 向讨论群发送讨论主题卡片
	SendTopicCard() errs.SystemErrorInfo
	// GenIssueTrend 生成发起讨论的动态
	GenIssueTrend() errs.SystemErrorInfo
	// AuthCreateChat 检查是否有权限创建讨论群
	AuthCreateChat() errs.SystemErrorInfo
}

// FeiShuIssueChat 任务群聊-飞书实现
type FeiShuIssueChat struct {
	Org                  *bo.BaseOrgInfoBo
	Project              *bo.ProjectBo
	IssueId              int64 `json:"issueId"`
	Issue                *bo.IssueBo
	ChatId               string `json:"chatId"`
	IsNewChat            bool   `json:"isNewChat"` // 是否新创建的群
	OpUserId             int64  `json:"opUserId"`  // 操作人的 userId
	MemberIds            []int64
	Topic                string
	TenantSdk            *sdk.Tenant
	MemberUserIds        []int64  // 包含任务协作人的 id，群聊发起人 id
	InitialMemberOpenIds []string // 包含任务协作人的 openId 和讨论群发起人的 openId
	UserOpenIdMap        map[int64]string
	ProjectAdminIds      []int64 // 项目管理员。不包括组织管理员（虽然组织管理员一定是项目管理员，但这里仅保存关联项目管理员角色的成员）
	OrgAdminIds          []int64 // 组织项目管理员
	ChatIsValid          bool    `json:"chatIsValid"`
	HasProBelong         bool    `json:"hasProBelong"` // 当前任务是否有项目归属
}

// NewFsIssueChat 查询实例化任务群聊所需的一些信息
func NewFsIssueChat(orgId, opUserId, issueId int64, topic string) (*FeiShuIssueChat, errs.SystemErrorInfo) {
	result := &FeiShuIssueChat{
		IssueId:  issueId,
		Topic:    topic,
		OpUserId: opUserId,
	}
	// 查询任务详情
	issues, err := GetIssueInfosLc(orgId, opUserId, []int64{issueId})
	if err != nil {
		log.Errorf("[NewFsIssueChat] GetIssueInfosLc err: %v", err)
		return result, err
	}
	if len(issues) < 1 {
		err := errs.IssueNotExist
		log.Errorf("[NewFsIssueChat] err: %v", err)
		return result, err
	}
	result.Issue = issues[0]

	issue := issues[0]
	if issue.OutChatId != "" {
		result.ChatId = issue.OutChatId
	}
	result.HasProBelong = CheckIssueHasProBelong(issue)

	return result, nil
}

// CollectBaseInfo 收集一些基础信息，用于创建群聊
func (ic *FeiShuIssueChat) CollectBaseInfo() errs.SystemErrorInfo {
	result := ic
	orgId := ic.Issue.OrgId
	// 查询组织详情
	baseOrgInfo, err := orgfacade.GetBaseOrgInfoRelaxed(orgId)
	if err != nil {
		log.Errorf("[CollectBaseInfo] GetBaseOrgInfoRelaxed err: %v, orgId: %d", err, orgId)
		return err
	}
	if baseOrgInfo.OutOrgId == "" {
		return errs.CannotBindChat
	}
	result.Org = baseOrgInfo

	tenant, err := feishu.GetTenant(baseOrgInfo.OutOrgId)
	if err != nil {
		log.Errorf("[CollectBaseInfo] GetTenant err: %v, outOrgId: %s", err, baseOrgInfo.OutOrgId)
		return err
	}
	result.TenantSdk = tenant

	initialChatMemberIds := make([]int64, 0)
	// 查询任务的协作人信息
	if ic.HasProBelong {
		initialChatMemberIds, err = tablefacade.GetDataCollaborateUserIds(orgId, ic.OpUserId, ic.Issue.DataId)
		if err != nil {
			log.Errorf("[NewFsIssueChat] GetDataCollaborateUserIds err: %v, DataId: %d", err, ic.Issue.DataId)
			return err
		}
		initialChatMemberIds = append(initialChatMemberIds, ic.OpUserId)
	} else {
		// 查询任务的负责人、关注人、确认人，默认视其为协作人 —— 伪协作人
		initialChatMemberIds = append(initialChatMemberIds, ic.Issue.OwnerIdI64...)
		initialChatMemberIds = append(initialChatMemberIds, ic.Issue.FollowerIdsI64...)
		initialChatMemberIds = append(initialChatMemberIds, ic.Issue.AuditorIdsI64...)
		initialChatMemberIds = append(initialChatMemberIds, ic.OpUserId)
	}
	initialChatMemberIds = slice.SliceUniqueInt64(initialChatMemberIds)
	result.UserOpenIdMap, err = GetOpenIdMapByUserIds(orgId, result.Org.SourceChannel, initialChatMemberIds)
	if err != nil {
		log.Errorf("[NewFsIssueChat] GetOpenIdMapByUserIds err:%v", err)
		return err
	}
	allUserOpenIds := make([]string, 0, len(result.UserOpenIdMap))
	for _, openId := range result.UserOpenIdMap {
		allUserOpenIds = append(allUserOpenIds, openId)
	}
	result.MemberUserIds = initialChatMemberIds
	result.InitialMemberOpenIds = allUserOpenIds

	// 查询项目管理员
	if result.Issue.ProjectId > 0 {
		project, err := GetProject(orgId, result.Issue.ProjectId)
		if err != nil {
			log.Errorf("[NewFsIssueChat] GetProject err:%v", err)
			return err
		}
		result.Project = project
		result.ProjectAdminIds = project.OwnerIds
	}
	// 组织管理员
	orgAdminResp := userfacade.GetUsersCouldManage(orgId, result.Issue.AppId)
	if orgAdminResp.Failure() {
		log.Errorf("[NewFsIssueChat] GetUsersCouldManage err:%v, orgId: %d", orgAdminResp.Error(), orgId)
		return orgAdminResp.Error()
	}
	orgAdminIds := make([]int64, 0)
	for _, user := range orgAdminResp.Data.List {
		orgAdminIds = append(orgAdminIds, user.Id)
	}
	result.OrgAdminIds = orgAdminIds

	return nil
}

// StartIssueChat 创建一个新的任务讨论群
func (ic *FeiShuIssueChat) StartIssueChat() (string, errs.SystemErrorInfo) {
	var err errs.SystemErrorInfo
	if err := ic.CollectBaseInfo(); err != nil {
		log.Errorf("[FeiShuIssueChat.NewIssueChat] err: %v, issueId: %d", err, ic.IssueId)
		return "", err
	}
	if ic.Project != nil && ic.Project.Id > 0 && ic.Project.IsFiling == 1 {
		err := errs.ProjectIsFilingYet
		log.Errorf("[FeiShuIssueChat.NewIssueChat] err: %v, issueId: %d", err, ic.IssueId)
		return "", err
	}
	if val, ok := ic.Issue.LessData[consts.BasicFieldRecycleFlag]; ok {
		isInRecycle := int(val.(float64)) == consts.AppIsDeleted
		if isInRecycle {
			log.Errorf("[FeiShuIssueChat.NewIssueChat] err: %v, issueId: %d, recycleFlag: %v", err, ic.IssueId, val)
			return "", err
		}
	}
	// 权限验证
	if err = ic.AuthCreateChat(); err != nil {
		log.Errorf("[FeiShuIssueChat.StartIssueChat] AuthCreateChat: %v, issueId: %d", err, ic.IssueId)
		return "", err
	}
	if ic.ChatId == "" {
		ic.ChatId, err = ic.NewIssueChat()
		if err != nil {
			log.Errorf("[FeiShuIssueChat.StartIssueChat] err: %v, issueId: %d", err, ic.IssueId)
			return "", err
		}
		ic.IsNewChat = true
		// 群聊 id 保存到任务中
		if err := ic.SaveIssueToLc(); err != nil {
			log.Errorf("[FeiShuIssueChat.StartIssueChat] SaveIssueToLc err: %v, issueId: %d", err, ic.IssueId)
			return "", err
		}
	} else {
		// 机器人不在群中/群被解散 isValid 为 false
		isValid, err := CheckChatIsValid(ic.TenantSdk, ic.ChatId)
		if err != nil {
			log.Errorf("[FeiShuIssueChat.StartIssueChat] CheckChatIsValid err: %v, issueId: %d, chatId: %s", err,
				ic.IssueId, ic.ChatId)
			return "", err
		}
		ic.ChatIsValid = isValid
		if isValid {
			// 检查发起人是否再群聊中，如果在，则无需处理；如不在，则将其拉入任务群聊中。
			// 无论不存在，拉入即可；如果存在，调用“拉入”也无妨。
			if opUserOpenId, ok := ic.UserOpenIdMap[ic.OpUserId]; ok {
				resp, oriErr := ic.TenantSdk.AddChatUser(fsvo.UpdateChatMemberReqVo{
					OpenIds: []string{opUserOpenId},
					ChatId:  ic.ChatId,
				})
				if oriErr != nil {
					log.Errorf("[FeiShuIssueChat.StartIssueChat] AddChatUser err: %v, issueId:%d", err, ic.IssueId)
					return "", errs.BuildSystemErrorInfo(errs.FeiShuOpenApiCallError, err)
				}
				if resp.Code != 0 {
					log.Errorf("[FeiShuIssueChat.StartIssueChat] AddChatUser code err resp: %s, issueId: %d, chatId: %s",
						json.ToJsonIgnoreError(resp), ic.IssueId, ic.ChatId)
					return "", errs.BuildSystemErrorInfoWithMessage(errs.FeiShuOpenApiCallError, resp.Msg)
				}
			}
		}
	}
	asyn.Execute(func() {
		if err = ic.GenIssueTrend(); err != nil {
			// 动态异常，不影响主流程
			log.Errorf("[FeiShuIssueChat.StartIssueChat] GenIssueTrend err: %v, issueId: %d", err, ic.IssueId)
		}
		if err = ic.SendTopicCard(); err != nil {
			log.Errorf("[FeiShuIssueChat.StartIssueChat] SendTopicCard err: %v, issueId: %d", err, ic.IssueId)
		}
	})

	return ic.ChatId, nil
}

// NewIssueChat 创建一个新的任务群聊
func (ic *FeiShuIssueChat) NewIssueChat() (string, errs.SystemErrorInfo) {
	// 任务创建群聊时加锁，防止并发
	issueIdStr := strconv.FormatInt(ic.IssueId, 10)
	lockKey := consts.CreateIssueChatLock + issueIdStr
	exist, err := cache.Exist(lockKey)
	if err != nil {
		log.Errorf("[NewIssueChat] check lock err: %v, key: %s", err, lockKey)
		return "", errs.BuildSystemErrorInfo(errs.CacheProxyError, err)
	}
	if exist {
		// 未获取到锁，直接响应错误信息
		return "", errs.CreateIssueChatDuplicate
	} else {
		err := cache.SetEx(lockKey, "1", 2)
		if err != nil {
			log.Errorf("[NewIssueChat] 获取%s锁时异常 %v", lockKey, err)
			return "", errs.BuildSystemErrorInfo(errs.CacheProxyError, err)
		}
	}

	// 创建飞书群聊
	trulyRemark := ic.Issue.Title
	// 查看有没有重复的名称（已删除的，因为创建的时候飞书的判断逻辑是所有信息和某个群一致，那么就返回旧有的群id，这里就判断名称相同，则在描述里面加点标识）
	//isExist, isExistErr := mysql.IsExistByCond(consts.TableIssue, db.Cond{
	//	consts.TcOrgId: ic.Org.OrgId,
	//	consts.TcTitle: ic.Issue.Title,
	//	consts.TcId:    db.NotEq(ic.Issue.Id),
	//})
	//if isExistErr != nil {
	//	log.Errorf("[FeiShuIssueChat.NewIssueChat] err: %v", isExistErr)
	//	return "", errs.BuildSystemErrorInfo(errs.MysqlOperateError, isExistErr)
	//}
	//if isExist {
	//	trulyRemark += fmt.Sprintf("（任务ID：%d）", ic.Issue.Id)
	//}

	reply, errSys := GetRawRows(ic.Org.OrgId, 0, &tablePb.ListRawRequest{
		FilterColumns: []string{consts.BasicFieldId},
		Condition: &tablePb.Condition{Type: tablePb.ConditionType_and, Conditions: []*tablePb.Condition{
			GetRowsCondition(consts.TcTitle, tablePb.ConditionType_equal, ic.Issue.Title, nil),
			GetRowsCondition(consts.BasicFieldIssueId, tablePb.ConditionType_un_equal, ic.Issue.Id, nil),
		}},
	})
	if err != nil {
		log.Errorf("[FeiShuIssueChat.NewIssueChat] GetRawRows err: %v, issueId:%v", err, ic.Issue.Id)
		return "", errSys
	}
	if len(reply.Data) != 0 {
		trulyRemark += fmt.Sprintf("（任务ID：%d）", ic.Issue.Id)
	}

	batchNum := 200
	firstOpenIds := make([]string, 0, batchNum)
	secondOpenIds := make([]string, 0, batchNum)
	if len(ic.InitialMemberOpenIds) <= batchNum {
		firstOpenIds = ic.InitialMemberOpenIds
	} else {
		firstOpenIds = ic.InitialMemberOpenIds[:batchNum]
		secondOpenIds = ic.InitialMemberOpenIds[batchNum:]
	}
	opUserOpenId, ok := ic.UserOpenIdMap[ic.OpUserId]
	if !ok {
		// 如果未找到，则机器人为群主
		log.Infof("[FeiShuIssueChat.NewIssueChat] CreateChat 时未找到操作人的 openId, issueId: %d", ic.IssueId)
	}
	resp, err := ic.TenantSdk.CreateChat(fsvo.CreateChatReqVo{
		Name:        consts.IssueChatTitle(ic.Issue.Title),
		Description: trulyRemark,
		OpenIds:     firstOpenIds,
		OwnerOpenId: opUserOpenId,
	})
	if err != nil {
		log.Errorf("[FeiShuIssueChat.NewIssueChat] CreateChat err: %v", err)
		return "", errs.BuildSystemErrorInfo(errs.FeiShuOpenApiCallError, err)
	}
	if resp.Code != 0 {
		log.Errorf("[FeiShuIssueChat.NewIssueChat] CreateChat code err resp: %s", json.ToJsonIgnoreError(resp))
		return "", errs.BuildSystemErrorInfoWithMessage(errs.FeiShuOpenApiCallError, resp.Msg)
	}
	ic.ChatId = resp.Data.ChatId
	ic.ChatIsValid = true

	// 异步：继续增加其余的群成员（每次只能传200个）
	asyn.Execute(func() {
		count := len(secondOpenIds)
		for i := 0; i < count; i += batchNum {
			max := i + batchNum
			if max > count {
				max = count
			}
			if err := ic.AddMembers(secondOpenIds[i:max]); err != nil {
				log.Errorf("[FeiShuIssueChat.NewIssueChat] err: %v, issueId: %d", err, ic.IssueId)
			}
		}
	})

	return resp.Data.ChatId, nil
}

func (ic *FeiShuIssueChat) SaveIssueToLc() errs.SystemErrorInfo {
	var err errs.SystemErrorInfo
	if ic.ChatId == "" {
		log.Infof("[SaveIssueToLc] do not need update for issue chat. issueId: %d", ic.IssueId)
		return nil
	}
	updForm := make([]map[string]interface{}, 0, 1)
	updForm = append(updForm, map[string]interface{}{
		consts.BasicFieldIssueId:   ic.IssueId,
		consts.BasicFieldOutChatId: ic.ChatId,
	})
	if ic.Issue.AppId <= 0 {
		ic.Issue.AppId, err = GetOrgSummaryAppId(ic.Org.OrgId)
		if err != nil {
			log.Errorf("[SaveIssueToLc] GetOrgSummaryAppId err: %v, issueId: %d", err, ic.IssueId)
			return err
		}
	}
	lcReq := formvo.LessUpdateIssueReq{
		OrgId:   ic.Issue.OrgId,
		AppId:   ic.Issue.AppId,
		TableId: ic.Issue.TableId,
		UserId:  ic.OpUserId,
		Form:    updForm,
	}
	// 更新任务信息到无码
	resp := formfacade.LessUpdateIssue(lcReq)
	if resp.Failure() {
		log.Errorf("[SaveIssueToLc] LessUpdateIssue: %v, issueId: %d", resp.Error(), ic.IssueId)
		return resp.Error()
	}

	return nil
}

// GenIssueTrend 发起任务群聊时，生成动态
func (ic *FeiShuIssueChat) GenIssueTrend() errs.SystemErrorInfo {
	trendInfo := projectvo.IssueChatTrendObj{
		Topic:     ic.Topic,
		IsNewChat: ic.IsNewChat,
	}
	issueTrendsBo := &bo.IssueTrendsBo{
		PushType:   consts.PushTypeIssueStartChat,
		OrgId:      ic.Issue.OrgId,
		ProjectId:  ic.Issue.ProjectId,
		TableId:    ic.Issue.TableId,
		DataId:     ic.Issue.DataId,
		IssueId:    ic.IssueId,
		IssueTitle: ic.Issue.Title,
		OperatorId: ic.OpUserId,
		NewValue:   json.ToJsonIgnoreError(trendInfo),
	}
	PushIssueTrends(issueTrendsBo)

	return nil
}

// DeleteMembers 从群聊中移除部分成员 todo
func (ic FeiShuIssueChat) DeleteMembers(delMemberIds []int64) errs.SystemErrorInfo {
	return nil
}

// AddMembers 向群聊中增加部分成员
func (ic FeiShuIssueChat) AddMembers(addMemberOpenIds []string) errs.SystemErrorInfo {
	resp, err := ic.TenantSdk.AddChatUser(fsvo.UpdateChatMemberReqVo{
		OpenIds: addMemberOpenIds,
		ChatId:  ic.ChatId,
	})
	if err != nil {
		log.Errorf("[FeiShuIssueChat.AddMembers] err: %v, chatId: %s", err, ic.ChatId)
		return errs.BuildSystemErrorInfo(errs.FeiShuOpenApiCallError, err)
	}
	if resp.Code != 0 {
		log.Errorf("[FeiShuIssueChat.AddMembers] AddMembers resp: %s, issueId: %d, chatId: %s",
			json.ToJsonIgnoreError(resp),
			ic.IssueId, ic.ChatId)
		return errs.BuildSystemErrorInfoWithMessage(errs.FeiShuOpenApiCallError, resp.Msg)
	}

	return nil
}

func (ic FeiShuIssueChat) SendTopicCard() errs.SystemErrorInfo {
	if ic.Topic == "" || !ic.ChatIsValid {
		return nil
	}
	infoBox, err := CollectCardInfo(ic.Org.OrgId, ic.IssueId, ic.OpUserId, map[string]interface{}{
		"issue":   ic.Issue,
		"project": ic.Project,
	})
	if err != nil {
		log.Errorf("[FeiShuIssueChat.SendTopicCard] CollectCardInfo err: %v, issueId: %d", err, ic.IssueId)
		return errs.BuildSystemErrorInfo(errs.FeiShuOpenApiCallError, err)
	}
	cardMsg := card.GetFsCardStartIssueChat(ic.Org.OrgId, ic.Topic, infoBox)
	errSys := card.PushCard(ic.Org.OrgId, &commonvo.PushCard{
		OrgId:         ic.Org.OrgId,
		OutOrgId:      ic.Org.OutOrgId,
		SourceChannel: ic.Org.SourceChannel,
		ChatIds:       []string{ic.ChatId},
		CardMsg:       cardMsg,
	})
	if errSys != nil {
		log.Errorf("[SendTopicCard] err:%v", errSys)
		return errSys
	}
	//cardInfo := card.GetFsCardStartIssueChat(ic.Topic, infoBox)
	//resp, oriErr := ic.TenantSdk.SendMessage(fsvo.MsgVo{
	//	ChatId:  ic.ChatId,
	//	MsgType: "interactive",
	//	Card:    cardInfo,
	//})
	//if oriErr != nil {
	//	log.Errorf("[FeiShuIssueChat.SendTopicCard] SendTopicCard err: %v, issueId: %d", oriErr, ic.IssueId)
	//	return errs.BuildSystemErrorInfo(errs.FeiShuOpenApiCallError, oriErr)
	//}
	//log.Infof("[FeiShuIssueChat.SendTopicCard] resp: %s", json.ToJsonIgnoreError(resp))

	return nil
}

func (ic *FeiShuIssueChat) AuthCreateChat() errs.SystemErrorInfo {
	if exist, _ := slice.Contain(ic.ProjectAdminIds, ic.OpUserId); exist {
		return nil
	}
	if exist, _ := slice.Contain(ic.OrgAdminIds, ic.OpUserId); exist {
		return nil
	}
	if exist, _ := slice.Contain(ic.MemberUserIds, ic.OpUserId); exist {
		return nil
	}

	return errs.DenyStartIssueChat
}

// GetOpenIdsByUserIds 获取协作人和一些额外成员的 openId 列表
func GetOpenIdsByUserIds(orgId int64, sourceChannel string,
	userIds []int64,
) ([]string, errs.SystemErrorInfo) {
	allOpenIds := make([]string, 0, len(userIds))
	if len(userIds) < 1 {
		return allOpenIds, nil
	}
	users, err := orgfacade.GetBaseUserInfoBatchRelaxed(orgId, userIds)
	if err != nil {
		log.Errorf("[GetCollaborateOpenIdsByCollaboratorInfo] GetBaseUserInfoBatchRelaxed err:%v", err)
		return nil, err
	}
	for _, user := range users {
		allOpenIds = append(allOpenIds, user.OutUserId)
	}

	return allOpenIds, nil
}

// GetOpenIdMapByUserIds 获取协作人和一些额外成员的 openId map: {uid: openId}
func GetOpenIdMapByUserIds(orgId int64, sourceChannel string, userIds []int64) (map[int64]string, errs.SystemErrorInfo) {
	userOpenIdMap := make(map[int64]string, len(userIds))
	if len(userIds) < 1 {
		return userOpenIdMap, nil
	}
	users, err := orgfacade.GetBaseUserInfoBatchRelaxed(orgId, userIds)
	if err != nil {
		log.Errorf("[GetCollaborateOpenIdsByCollaboratorInfo] GetBaseUserInfoBatchRelaxed err:%v", err)
		return nil, err
	}
	for _, user := range users {
		userOpenIdMap[user.UserId] = user.OutUserId
	}

	return userOpenIdMap, nil
}

// GetUserIdsByCollaboratorInfo 获取协作人和一些额外成员的 id 列表
func GetUserIdsByCollaboratorInfo(collaborateInfo []*projectvo.DataCollaborators, extraUserIds []int64,
) ([]int64, errs.SystemErrorInfo) {
	allUserIds := make([]int64, 0, 20)
	for _, issueCollaborate := range collaborateInfo {
		for _, item := range issueCollaborate.ColumnCollaborators {
			allUserIds = append(allUserIds, item.UserIds...)
		}
	}
	if len(extraUserIds) > 0 {
		allUserIds = append(allUserIds, extraUserIds...)
	}
	allUserIds = slice.SliceUniqueInt64(allUserIds)

	return allUserIds, nil
}

// 同步任务群名称
func SyncIssueChatTitle(orgId, issueId int64, chatGroupName string, sourceChannel string) errs.SystemErrorInfo {
	// 查询是否有任务群聊
	issueHasChat, chatId, errorInfo := CheckIssueHasChat(orgId, issueId)
	if errorInfo != nil {
		log.Errorf("[SyncIssueChatTitle] CheckIssueHasChat err:%v, orgId:%d, issueId:%d",
			errorInfo, orgId, issueId)
		return errorInfo
	}
	if !issueHasChat {
		// 没有任务群聊直接返回，不处理
		return nil
	}
	baseOrgInfoBo, err := orgfacade.GetBaseOrgInfoRelaxed(orgId)
	if err != nil {
		log.Errorf("[SyncIssueChatTitle] orgfacade.GetBaseOrgInfoRelaxe err:%v, orgId:%d", err, orgId)
		return err
	}
	tenant, tenantErr := feishu.GetTenant(baseOrgInfoBo.OutOrgId)
	if tenantErr != nil {
		log.Errorf("[SyncIssueChatTitle] feishu.GetTenant err:%v, orgId:%d, issueId:%d", tenantErr, orgId, issueId)
		return tenantErr
	}
	chatResp, errChat := tenant.UpdateChat(fsvo.UpdateChatReqVo{
		ChatId: chatId,
		Name:   consts.IssueChatTitle(chatGroupName),
	})
	if errChat != nil {
		log.Errorf("[SyncIssueChatTitle] UpdateChat 修改群聊信息失败, err:%v, orgId:%d, issueId:%d", errChat, orgId, issueId)
		return errs.FeiShuOpenApiCallError
	}
	if chatResp.Code != 0 {
		log.Errorf("[SyncIssueChatTitle] UpdateChat 修改群聊信息失败, err:%v, respCode:%v, orgId:%d, issueId:%d",
			errChat, chatResp.Code, orgId, issueId)
		return errs.FeiShuOpenApiCallError
	}
	return nil
}

// 查询是否有任务群聊
func CheckIssueHasChat(orgId int64, issueId int64) (bool, string, errs.SystemErrorInfo) {
	issueBos, err := GetIssueInfosLc(orgId, 0, []int64{issueId})
	if err != nil {
		log.Errorf("[CheckIssueHasChat] GetIssueInfosLc err:%v, orgId:%d, issueId:%d", err, orgId, issueId)
		return false, "", err
	}
	if len(issueBos) > 0 {
		issueBo := issueBos[0]
		if issueBo.OutChatId != "" {
			return true, issueBo.OutChatId, nil
		}
	}
	return false, "", nil
}

// CollectCardInfo 有些地方调用该方法可能已经有数据了，因此提供传入 sourceMap
// 请确保传入 sourceMap 的数据源符合 projectvo.CardInfoBoxForStartIssueChat 中对应的类型，如果不确定，请不要传
func CollectCardInfo(orgId, issueId, opUserId int64, sourceMap map[string]interface{}) (*projectvo.CardInfoBoxForStartIssueChat, errs.SystemErrorInfo) {
	infoBox := &projectvo.CardInfoBoxForStartIssueChat{}
	if val, ok := sourceMap["issue"]; ok {
		infoBox.IssueInfo = val.(*bo.IssueBo)
	} else {
		issues, err := GetIssueInfosLc(orgId, opUserId, []int64{issueId})
		if err != nil {
			log.Errorf("[CollectCardInfo] err: %v, issueId: %v", err, issueId)
			return nil, err
		}
		if len(issues) > 0 {
			infoBox.IssueInfo = issues[0]
		} else {
			err := errs.IssueNotExist
			log.Errorf("[CollectCardInfo] err: %v, issueId: %d", err, issueId)
			return nil, err
		}
	}
	if infoBox.IssueInfo.ParentId > 0 {
		if val, ok := sourceMap["parentIssue"]; ok {
			infoBox.ParentIssue = val.(*bo.IssueBo)
		} else {
			parentIssues, err := GetIssueInfosLc(orgId, opUserId, []int64{infoBox.IssueInfo.ParentId})
			if err != nil {
				log.Errorf("[CollectCardInfo] err: %v, issueId: %v", err, infoBox.IssueInfo.ParentId)
				return nil, err
			}
			if len(parentIssues) > 0 {
				infoBox.ParentIssue = parentIssues[0]
			} else {
				err := errs.IssueNotExist
				log.Errorf("[CollectCardInfo] err: %v, issueId: %d", err, infoBox.IssueInfo.ParentId)
				return nil, err
			}
		}
	}

	if val, ok := sourceMap["project"]; ok {
		infoBox.ProjectBo = val.(*bo.ProjectBo)
	} else {
		if infoBox.IssueInfo.ProjectId > 0 {
			project, err := GetProject(orgId, infoBox.IssueInfo.ProjectId)
			if err != nil {
				log.Errorf("[CollectCardInfo] GetProject err:%v", err)
				return nil, err
			}
			infoBox.ProjectBo = project
		}
	}

	//if val, ok := sourceMap["table"]; ok {
	//	infoBox.IssueTableInfo = val.(*projectvo.TableMetaData)
	//} else {
	{
		if infoBox.IssueInfo.TableId > 0 {
			tableInfo, err := GetTableByTableId(orgId, opUserId, infoBox.IssueInfo.TableId)
			if err != nil {
				log.Errorf("[CollectCardInfo] GetTableByTableId err:%v, tableId:%d", err, infoBox.IssueInfo.TableId)
				return nil, err
			}
			infoBox.IssueTableInfo = tableInfo
		} else {
			// 默认一个任务表
			infoBox.IssueTableInfo = &projectvo.TableMetaData{
				Name: consts.DefaultTableName,
			}
			// err := errs.InvalidTableId
			// log.Errorf("[CollectCardInfo] err: %v, issueId: %d", err, issueId)
			// return infoBox, err
		}
	}

	//if val, ok := sourceMap["tableHeader"]; ok {
	//	infoBox.IssueTableInfo = val.(*projectvo.TableMetaData)
	//} else {
	{
		headers := make(map[string]lc_table.LcCommonField, 0)
		if infoBox.IssueInfo.TableId > 0 {
			tableColumns, err := GetTableColumnsMap(orgId, infoBox.IssueInfo.TableId, nil)
			if err != nil {
				log.Errorf("[CollectCardInfo] 获取表头失败 org:%d, table:%d, user:%d, err: %v",
					orgId, infoBox.IssueInfo.TableId, opUserId, err)
				return nil, err
			}
			infoBox.ProjectTableColumnMap = tableColumns
			copyer.Copy(tableColumns, &headers)
		}
		infoBox.TableColumnMap = headers
	}

	userIds := make([]int64, 0, len(infoBox.IssueInfo.OwnerIdI64)+1)
	userIds = append(userIds, opUserId)
	userIds = append(userIds, infoBox.IssueInfo.OwnerIdI64...)
	userIds = slice.SliceUniqueInt64(userIds)
	userInfoArr, err := orgfacade.GetBaseUserInfoBatchRelaxed(orgId, userIds)
	if err != nil {
		log.Errorf("[CollectCardInfo] 查询组织 %d 用户信息出现异常 %v", orgId, err)
		return nil, err
	}
	userMap := make(map[int64]bo.BaseUserInfoBo, len(userInfoArr))
	for _, user := range userInfoArr {
		userMap[user.UserId] = user
	}
	if opUser, ok := userMap[opUserId]; ok {
		infoBox.OperateUser = &opUser
	}
	ownerInfos := []*bo.BaseUserInfoBo{}
	for _, uid := range infoBox.IssueInfo.OwnerIdI64 {
		if user, ok := userMap[uid]; ok {
			ownerInfos = append(ownerInfos, &user)
		}
	}
	infoBox.IssueOwners = ownerInfos

	links := GetIssueLinks("", orgId, issueId)

	if val, ok := sourceMap["issueInfoUrl"]; ok {
		infoBox.IssueInfoUrl = val.(string)
	} else {
		infoBox.IssueInfoUrl = links.SideBarLink
	}

	if val, ok := sourceMap["issuePcUrl"]; ok {
		infoBox.IssuePcUrl = val.(string)
	} else {
		infoBox.IssuePcUrl = links.Link
	}

	return infoBox, nil
}

// GetProjectAdminIds 获取项目管理员ids，包括组织管理员ids
func GetProjectAdminIds(orgId int64, appId int64) ([]int64, errs.SystemErrorInfo) {
	adminIds := []int64{}
	projectAdminIds, errSys := GetProAdminUserIdsForLessCodeApp(orgId, 0, appId)
	if errSys != nil {
		log.Errorf("[GetProjectAdminIds]GetProAdminUserIdsForLessCodeApp err:%v, orgId:%d, appId:%d",
			errSys, orgId, appId)
		return nil, errSys
	}
	adminIds = append(adminIds, projectAdminIds...)

	userIdsOfOrg, errSys := GetAdminUserIdsOfOrg(orgId, appId)
	if errSys != nil {
		log.Errorf("[GetProjectAdminIds]GetAdminUserIdsOfOrg err:%v, orgId:%d, appId:%d", errSys, orgId, appId)
		return nil, errSys
	}
	adminIds = append(adminIds, userIdsOfOrg...)
	adminIds = slice.SliceUniqueInt64(adminIds)
	return adminIds, nil
}

// 找出更新的协作人
func GetDiffCollaborators(updateIssueVo projectvo.UpdateChatIssueVo, tableColumns map[string]*projectvo.TableColumnData) ([]int64, []int64, errs.SystemErrorInfo) {
	originCollaboratorIds := []int64{}
	curCollaboratorIds := []int64{}
	for k, v := range updateIssueVo.LcData {
		// 成员
		if header, ok := tableColumns[k]; ok {
			if header.Field.Type == consts.LcColumnFieldTypeMember {
				props := lc_table.LcProps{}
				err := copyer.Copy(header.Field.Props, &props)
				if err != nil {
					log.Errorf("[GetDiffCollaborators] copy err:%v", err)
					return nil, nil, errs.ObjectCopyError
				}
				// 只找协作人的
				if CheckColumnIsCollaborate(k, props) {
					newMember := cast.ToStringSlice(v)
					newMemberIds := businees.LcMemberToUserIds(newMember)
					oldMember := cast.ToStringSlice(updateIssueVo.OldLcData[k])
					oldMemberIds := businees.LcMemberToUserIds(oldMember)
					originCollaboratorIds = append(originCollaboratorIds, oldMemberIds...)
					curCollaboratorIds = append(curCollaboratorIds, newMemberIds...)
				}
			}
		}
	}

	delCollaboratorIds, addCollaboratorIds := util.GetDifMemberIds(originCollaboratorIds, curCollaboratorIds)

	return delCollaboratorIds, addCollaboratorIds, nil
}

// CheckColumnIsCollaborate 检查列是否是开启协作人
func CheckColumnIsCollaborate(columnId string, columnProps lc_table.LcProps) bool {
	// 确认人一定是协作人
	isCollaborate := columnProps.Member.IsCollaborators || columnId == consts.BasicFieldAuditorIds ||
		len(columnProps.CollaboratorRoles) > 0

	return isCollaborate
}

// 处理协作人变更 被拉进群还是被踢出群
func DealChatWithUpdateCollaborators(orgId, appId int64, updateIssueVo projectvo.UpdateChatIssueVo, tableColumns map[string]*projectvo.TableColumnData) errs.SystemErrorInfo {
	issueId := updateIssueVo.IssueId
	issueInfosLc, err := GetIssueInfosLc(orgId, 0, []int64{issueId})
	if err != nil {
		log.Errorf("[DealChatWithUpdateCollaborators] GetIssueInfosLc err:%v, orgId:%v, issueId:%v", err, orgId, issueId)
		return err
	}
	outChatId := ""
	if len(issueInfosLc) > 0 {
		outChatId = issueInfosLc[0].OutChatId
	}
	// 没有群聊id，直接返回
	if outChatId == "" {
		return nil
	}

	delCollaboratorIds, addCollaboratorIds, err := GetDiffCollaborators(updateIssueVo, tableColumns)
	if err != nil {
		log.Errorf("[DealChatWithUpdateCollaborators]GetDiffCollaborators err:%v", err)
		return err
	}

	// 获取管理员ids
	adminIds, errSys := GetProjectAdminIds(orgId, appId)
	if errSys != nil {
		log.Errorf("[DealChatWithUpdateCollaborators]GetProjectAdminIds err:%v", errSys)
		return errSys
	}
	// 删除的协作人如果是管理员，需要过滤出来
	delIds := businees.DifferenceInt64Set(delCollaboratorIds, adminIds)

	baseOrgInfo, err := orgfacade.GetBaseOrgInfoRelaxed(orgId)
	if err != nil {
		log.Errorf("[DealChatWithUpdateCollaborators]GetBaseOrgInfoRelaxed err:%v, orgId:%d, appId:%d",
			err, orgId, appId)
		return err
	}

	userIds := []int64{}
	userIds = append(userIds, delIds...)
	userIds = append(userIds, addCollaboratorIds...)
	userIds = slice.SliceUniqueInt64(userIds)

	userInfoBatch := orgfacade.GetBaseUserInfoBatch(orgvo.GetBaseUserInfoBatchReqVo{
		OrgId:   orgId,
		UserIds: userIds,
	})
	if userInfoBatch.Failure() {
		log.Errorf("[DealChatWithUpdateCollaborators] GetBaseUserInfoBatch error:%v, orgId:%d",
			userInfoBatch.Error(), orgId)
		return userInfoBatch.Error()
	}

	userInfosMap := map[int64]bo.BaseUserInfoBo{}

	for _, user := range userInfoBatch.BaseUserInfos {
		userInfosMap[user.UserId] = user
	}

	if len(delIds) > 0 {
		// 直接踢出群聊, todo 暂不处理
		//openIds := []string{}
		//for _, id := range delIds {
		//	openIds = append(openIds, userInfosMap[id].OutUserId)
		//}
		//errSys := RemoveCollaboratorsFromChat(baseOrgInfo.OutOrgId, outChatId, openIds)
		//if errSys != nil {
		//	log.Error(errSys)
		//	return errSys
		//}

	}

	if len(addCollaboratorIds) > 0 {
		// 拉人进群
		openIds := []string{}
		for _, id := range addCollaboratorIds {
			openIds = append(openIds, userInfosMap[id].OutUserId)
		}
		errSys := AddCollaboratorsToChat(baseOrgInfo.OutOrgId, outChatId, openIds)
		if errSys != nil {
			log.Error(errSys)
			return errSys
		}
	}
	return nil
}

// 移除协作人 群聊
func RemoveCollaboratorsFromChat(outOrgId string, chatId string, openIds []string) errs.SystemErrorInfo {
	tenant, tenantErr := feishu.GetTenant(outOrgId)
	if tenantErr != nil {
		log.Errorf("[RemoveCollaboratorsFromChat]GetTenant err:%v, outOrgId:%v, openIds:%v", tenantErr, outOrgId, openIds)
		return errs.FeiShuClientTenantError
	}

	// 限制一次操作50人
	batchNum := 50
	count := len(openIds)
	for i := 0; i < count; i += batchNum {
		max := i + batchNum
		if max > count {
			max = count
		}
		resp, err := tenant.RemoveChatUser(fsvo.UpdateChatMemberReqVo{
			ChatId:  chatId,
			OpenIds: openIds[i:max],
		})
		if err != nil {
			log.Errorf("[RemoveCollaboratorsFromChat] RemoveChatUser err:%v, outOrgId:%v, openIds:%v",
				err, outOrgId, openIds)
			return errs.FeiShuOpenApiCallError
		}
		if resp.Code != 0 {
			log.Errorf("[RemoveCollaboratorsFromChat] 移除用户出群异常:%v, code:%v", resp.Msg, resp.Code)
			return errs.FeiShuOpenApiCallError
		}
	}

	log.Infof("[RemoveCollaboratorsFromChat] 移除用户出群成功！chatId:%v, outOrgId:%s, openIds:%v",
		chatId, outOrgId, openIds)

	return nil
}

// 拉人进群
func AddCollaboratorsToChat(outOrgId string, chatId string, openIds []string) errs.SystemErrorInfo {
	tenant, tenantErr := feishu.GetTenant(outOrgId)
	if tenantErr != nil {
		log.Errorf("[AddCollaboratorsToChat]GetTenant err:%v, outOrgId:%v, openIds:%v", tenantErr, outOrgId, openIds)
		return errs.FeiShuClientTenantError
	}

	// 限制一次操作50人
	batchNum := 50
	count := len(openIds)
	for i := 0; i < count; i += batchNum {
		max := i + batchNum
		if max > count {
			max = count
		}
		resp, err := tenant.AddChatUser(fsvo.UpdateChatMemberReqVo{
			ChatId:  chatId,
			OpenIds: openIds[i:max],
		})
		if err != nil {
			log.Errorf("[AddCollaboratorsToChat] AddChatUser err:%v, outOrgId:%v, openIds:%v",
				err, outOrgId, openIds)
			return errs.FeiShuOpenApiCallError
		}
		if resp.Code != 0 {
			log.Errorf("[AddCollaboratorsToChat] 拉用户进群异常:%v, code:%v", resp.Msg, resp.Code)
			return errs.FeiShuOpenApiCallError
		}
	}

	log.Infof("[AddCollaboratorsToChat] 拉用户进群成功！chatId:%v, outOrgId:%s, openIds:%v",
		chatId, outOrgId, openIds)
	return nil
}

func dealDeleteCollaboratorsColumnToIssueChat() errs.SystemErrorInfo {
	return nil
}

func dealUpdateCollaboratorsColumnToIssueChat(orgId, tableId int64, outOrgId string, newField lc_table.LcCommonField,
	userIds []int64, chatUserMap map[string][]int64) errs.SystemErrorInfo {
	oldIsCollaborator := false
	newIsCollaborator := false
	columns, errSys := GetTableColumns(orgId, 0, tableId, false)
	if errSys != nil {
		log.Errorf("[DealUpdateCollaboratorsColumnToIssueChat] GetTableColumns err:%v, orgId:%d, tableId:%d", errSys, orgId, tableId)
		return errSys
	}
	if len(columns) > 0 {
		oldColumnData := columns[0]
		oldLcProps := lc_table.LcProps{}
		err := copyer.Copy(oldColumnData.Field.Props, &oldLcProps)
		if err != nil {
			log.Errorf("[DealCollaboratorsColumnsToIssueChat] copy err:%v", err)
			return errs.ObjectCopyError
		}
		oldIsCollaborator = oldLcProps.Member.IsCollaborators
	}

	newIsCollaborator = newField.Field.Props.Member.IsCollaborators

	if newIsCollaborator == oldIsCollaborator {
		// 协作人开关没变，直接返回
		return nil
	}

	if oldIsCollaborator && !newIsCollaborator {
		// 协作人开关关闭，暂不处理 todo
	}

	if !oldIsCollaborator && newIsCollaborator {
		// 协作人打开开关
		err := OpenCollaboratorColumnWithIssueChat(orgId, outOrgId, userIds, chatUserMap)
		if err != nil {
			log.Errorf("[DealUpdateCollaboratorsColumnToIssueChat] err:%v, orgId:%d", err, orgId)
			return err
		}
	}

	return nil
}

func OpenCollaboratorColumnWithIssueChat(orgId int64, outOrgId string, userIds []int64, chatUsersMap map[string][]int64) errs.SystemErrorInfo {
	// 查询成员外部信息
	userIds = slice.SliceUniqueInt64(userIds)
	userInfoBatch := orgfacade.GetBaseUserInfoBatch(orgvo.GetBaseUserInfoBatchReqVo{
		OrgId:   orgId,
		UserIds: userIds,
	})
	if userInfoBatch.Failure() {
		log.Errorf("[OpenCollaboratorColumnWithIssueChat]GetBaseUserInfoBatch err:%v, orgId:%d", userInfoBatch.Error(), orgId)
		return userInfoBatch.Error()
	}
	userInfosMap := map[int64]bo.BaseUserInfoBo{}
	for _, user := range userInfoBatch.BaseUserInfos {
		userInfosMap[user.UserId] = user
	}

	for chaId, us := range chatUsersMap {
		openIds := make([]string, 0, len(us))
		for _, userId := range us {
			user, ok := userInfosMap[userId]
			if !ok {
				continue
			}
			openIds = append(openIds, user.OutUserId)
		}

		// 拉人进群
		errorChat := AddCollaboratorsToChat(outOrgId, chaId, openIds)
		if errorChat != nil {
			log.Errorf("[OpenCollaboratorColumnWithIssueChat]AddCollaboratorsToChat err:%v", errorChat)
			return errorChat
		}
	}

	return nil
}

// 创建/更新/删除列，处理协作人字段开启/关闭，涉及到的成员同步到群
func DealCollaboratorsColumnsToIssueChat(orgId int64, tableId int64, newTableColumn *projectvo.TableColumnData, columnAction int) errs.SystemErrorInfo {
	baseOrgInfoBo, errSys := orgfacade.GetBaseOrgInfoRelaxed(orgId)
	if errSys != nil {
		log.Errorf("[DealCollaboratorsColumnsToIssueChat]GetBaseOrgInfoRelaxed err:%v, orgId:%d, tableId:%v",
			errSys, orgId, tableId)
		return errSys
	}

	outOrgId := baseOrgInfoBo.OutOrgId

	// 获取该表下所有 有任务群聊的任务
	issueBosMap, errSys := GetIssueListByTableIdWithChaId(orgId, tableId)
	if errSys != nil {
		log.Errorf("[DealCollaboratorsColumnsToIssueChat] GetIssueListByTableId err:%v", errSys)
		return errSys
	}

	// 过滤掉协作人字段不为空的任务
	chatUsersMap := map[string][]int64{} // chatId: [userId1,userId2] 存群聊对应的协作人ids
	userIds := []int64{}
	//issueMaps := make([]map[string]interface{}, 0, len(issueBosMap))
	for _, issueBoMap := range issueBosMap {
		collaboratorIdsStrSlice := []string{}
		if collaborators, ok := issueBoMap[newTableColumn.Name]; ok {
			if cols, ok2 := collaborators.([]interface{}); ok2 {
				if len(cols) == 0 {
					continue
				}
				for _, col := range cols {
					ss := col.(string)
					collaboratorIdsStrSlice = append(collaboratorIdsStrSlice, ss)
				}
			}
		}
		collaboratorIds := businees.LcMemberToUserIds(collaboratorIdsStrSlice)
		if chatIdI, ok2 := issueBoMap[consts.BasicFieldOutChatId]; ok2 {
			chatId := chatIdI.(string)
			if chatId != "" {
				chatUsersMap[chatId] = collaboratorIds
				userIds = append(userIds, collaboratorIds...)
			}
		}
	}

	userIds = slice.SliceUniqueInt64(userIds)

	newFields := lc_table.LcCommonField{}
	err := copyer.Copy(newTableColumn, &newFields)
	if err != nil {
		log.Errorf("[DealCollaboratorsColumnsToIssueChat] copy err:%v", err)
		return errs.ObjectCopyError
	}

	switch columnAction {
	case consts.UpdateColumn:
		chatErr := dealUpdateCollaboratorsColumnToIssueChat(orgId, tableId, outOrgId, newFields, userIds, chatUsersMap)
		if chatErr != nil {
			log.Errorf("[DealCollaboratorsColumnsToIssueChat] 处理更新列失败:%v", chatErr)
			return chatErr
		}
	case consts.DeleteColumn:
		// 删除列 todo
		//dealDeleteCollaboratorsColumnToIssueChat()
	}

	return nil
}

// CheckChatIsValid 根据 chatId 检查群聊是否有效
func CheckChatIsValid(clientSdk *sdk.Tenant, chatId string) (bool, errs.SystemErrorInfo) {
	chatInfoResp, err := clientSdk.ChatInfo(chatId)
	if err != nil {
		log.Errorf("[CheckChatIsValid] ChatInfo err: %v, chatId: %s", err, chatId)
		return false, errs.BuildSystemErrorInfo(errs.FeiShuOpenApiCallError, err)
	}
	if chatInfoResp.Code != 0 {
		log.Infof("[CheckChatIsValid] ChatInfo code err resp: %s, chatId: %s", json.ToJsonIgnoreError(chatInfoResp), chatId)
		// 群异常或者群不存在，则返回 false
		return false, nil
	}
	if chatInfoResp.Data.ChatId != "" {
		return true, nil
	}

	return false, nil
}

//// IssueChatInviteUsersWhenMoveIssue 跨表移动时，拉人到对应的任务群聊
//func IssueChatInviteUsersWhenMoveIssue(orgId int64, issueIds []int64, targetTableId int64) errs.SystemErrorInfo {
//	// 查询任务的协作人信息
//	collaborateResp := authfacade.GetCollaboratorIdsByIssueIds(projectvo.GetCollaboratorIdsByIssueIdsRequest{
//		OrgId: orgId,
//		Input: &permissionV1.GetDataCollaboratorsRequest{
//			IssueIds: issueIds,
//		},
//	})
//	if collaborateResp.Failure() {
//		log.Errorf("[IssueChatInviteUsersWhenMoveIssue] GetCollaboratorIdsByIssueIds err: %v, IssueIds: %v",
//			collaborateResp.Failure(), issueIds)
//		return collaborateResp.Error()
//	}
//	collaborateDataArr := make([]*projectvo.DataCollaborators, len(collaborateResp.Data.DataCollaborators))
//	errCopy := copyer.Copy(collaborateResp.Data.DataCollaborators, &collaborateDataArr)
//	if errCopy != nil {
//		log.Errorf("[IssueChatInviteUsersWhenMoveIssue] copy err: %v", errCopy)
//		return errs.BuildSystemErrorInfo(errs.TypeConvertError, errCopy)
//	}
//	collaboratorIdsMap := map[int64][]int64{} // key:issueId, value:[userId]
//	userIds := []int64{}
//	for _, collaborates := range collaborateDataArr {
//		for _, c := range collaborates.ColumnCollaborators {
//			userIds = append(userIds, c.UserIds...)
//			collaboratorIdsMap[collaborates.IssueId] = c.UserIds
//		}
//	}
//
//	// 根据issueIds 找出任务群聊id
//	issueList, errSys := GetIssueInfosLc(orgId, 0, issueIds)
//	if errSys != nil {
//		log.Errorf("[IssueChatInviteUsersWhenMoveIssue]GetIssueInfosLc err:%v, orgId:%d, issueIds:%v",
//			errSys, orgId, issueIds)
//		return errSys
//	}
//	issueChatIdMap := map[int64]string{} // key:issueId  value:chatId
//	for _, issue := range issueList {
//		if issue.OutChatId == "" {
//			continue
//		}
//		issueChatIdMap[issue.Id] = issue.OutChatId
//	}
//
//	chatUsersMap := map[string][]int64{} // chatId: [userId1,userId2] 存群聊对应的协作人ids
//	for k, v := range issueChatIdMap {
//		chatUsersMap[v] = collaboratorIdsMap[k]
//	}
//
//	// 添加群聊
//	baseOrgInfoBo, err := orgfacade.GetBaseOrgInfoRelaxed(orgId)
//	if err != nil {
//		log.Errorf("[IssueChatInviteUsersWhenMoveIssue]GetBaseOrgInfoRelaxed err:%v, orgId:%d, issueIds:%v",
//			errSys, orgId, issueIds)
//		return err
//	}
//
//	outOrgId := baseOrgInfoBo.OutOrgId
//
//	userIds = slice.SliceUniqueInt64(userIds)
//	userInfoBatch := orgfacade.GetBaseUserInfoBatch(orgvo.GetBaseUserInfoBatchReqVo{
//		OrgId:   orgId,
//		UserIds: userIds,
//	})
//	if userInfoBatch.Failure() {
//		log.Errorf("[IssueChatInviteUsersWhenMoveIssue]GetBaseUserInfoBatch err:%v, orgId:%d, issueIds:%v",
//			errSys, orgId, issueIds)
//		return userInfoBatch.Error()
//	}
//	userInfosMap := map[int64]bo.BaseUserInfoBo{}
//
//	for _, user := range userInfoBatch.BaseUserInfos {
//		userInfosMap[user.UserId] = user
//	}
//
//	for chaId, us := range chatUsersMap {
//		openIds := make([]string, 0, len(us))
//		for _, userId := range us {
//			user, ok := userInfosMap[userId]
//			if !ok {
//				continue
//			}
//			openIds = append(openIds, user.OutUserId)
//		}
//
//		// 拉人进群  注：未判断该成员是否已经在群里了
//		if len(openIds) > 0 {
//			chatErr := AddCollaboratorsToChat(outOrgId, chaId, openIds)
//			if chatErr != nil {
//				log.Errorf("[IssueChatInviteUsersWhenMoveIssue] AddCollaboratorsToChat err: %v, chatId: %s", chatErr, chaId)
//				return chatErr
//			}
//		}
//	}
//
//	return nil
//}

// CheckIssueHasProBelong 检查任务是否有项目归属
func CheckIssueHasProBelong(issue *bo.IssueBo) bool {
	if issue == nil {
		return false
	}
	if issue.AppId >= 0 && issue.ProjectId >= 0 {
		return true
	}

	return false
}
