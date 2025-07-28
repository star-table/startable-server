package notice

import (
	"container/list"

	"github.com/star-table/startable-server/app/facade/tablefacade"

	"github.com/star-table/startable-server/common/model/vo/commonvo"

	pushV1 "gitea.bjx.cloud/LessCode/interface/golang/push/v1"
	"github.com/star-table/startable-server/app/facade/orgfacade"
	"github.com/star-table/startable-server/app/service/projectsvc/domain"
	"github.com/star-table/startable-server/common/core/businees"
	"github.com/star-table/startable-server/common/core/consts"
	"github.com/star-table/startable-server/common/core/errs"
	"github.com/star-table/startable-server/common/core/util"
	"github.com/star-table/startable-server/common/core/util/asyn"
	"github.com/star-table/startable-server/common/core/util/copyer"
	"github.com/star-table/startable-server/common/core/util/json"
	"github.com/star-table/startable-server/common/core/util/maps"
	"github.com/star-table/startable-server/common/core/util/slice"
	int642 "github.com/star-table/startable-server/common/core/util/slice/int64"
	"github.com/star-table/startable-server/common/core/util/str"
	"github.com/star-table/startable-server/common/extra/card"
	"github.com/star-table/startable-server/common/model/bo"
	"github.com/star-table/startable-server/common/model/vo/orgvo"
	"github.com/star-table/startable-server/common/model/vo/projectvo"
)

func PushIssue(issueTrendsBo bo.IssueTrendsBo) {
	asyn.Execute(func() {
		baseOrgInfo := orgfacade.GetBaseOrgInfo(orgvo.GetBaseOrgInfoReqVo{OrgId: issueTrendsBo.OrgId})
		if !baseOrgInfo.Failure() {
			PushIssueByChannel(issueTrendsBo, baseOrgInfo.BaseOrgInfo.SourceChannel)
		}
	})
}

func PushIssueByChannel(issueTrendsBo bo.IssueTrendsBo, sourceChannel string) {
	log.Infof("[PushIssueByChannel]sourceChannel:%v, issueTrendsBo:%v", sourceChannel, json.ToJsonIgnoreError(issueTrendsBo))
	pushType := issueTrendsBo.PushType
	switch pushType {
	case consts.PushTypeUpdateIssue: // 任务更新
		HandlePushCardForUpdateIssue(issueTrendsBo, sourceChannel)
	case consts.PushTypeCreateIssue: // , consts.PushTypeRecoverIssue, consts.PushTypeDeleteIssue , consts.PushTypeUpdateRelationIssue
		HandlePushCardForDefaultEvent(issueTrendsBo, sourceChannel)
	default:
	}
}

// PushIssueComment 评论时被 at，或者任务详情中被 at，此外，非当事人也能收到推送，注释卡片标题有所不同
func PushIssueComment(issueTrendsBo bo.IssueTrendsBo, content string, mentionedUserIds []int64, pushType consts.IssueNoticePushType) {
	orgId := issueTrendsBo.OrgId
	operatorId := issueTrendsBo.OperatorId
	infoBox, err := GetBaseInfoForCardPush(&issueTrendsBo)
	if err != nil {
		log.Errorf("[PushIssueComment] err: %v, issueId: %d", err, issueTrendsBo.IssueId)
		return
	}
	issueNoticeBo := bo.IssueNoticeBo{}
	errSys := copyer.Copy(issueTrendsBo, &issueNoticeBo)
	if errSys != nil {
		return
	}

	baseOrgInfo := orgfacade.GetBaseOrgInfo(orgvo.GetBaseOrgInfoReqVo{OrgId: orgId})
	if baseOrgInfo.Failure() {
		log.Errorf("[PushIssueComment] GetBaseOrgInfo err:%v, orgId:%v", baseOrgInfo.Error(), orgId)
		return
	}

	// 查询协作人信息（忽略负责人、确认人列）
	collaborateUserIds := infoBox.CollaboratorIds
	ownerIds := infoBox.IssueInfo.OwnerIdI64
	bePushedUserIds := make([]int64, 0, len(collaborateUserIds)+len(ownerIds)+len(mentionedUserIds))
	bePushedUserIds = append(bePushedUserIds, collaborateUserIds...)
	// 如果评论或任务描述中被 at（@），则强制推送（不受个人推送配置影响）
	beAtUserIds := util.GetCommentAtUserIds(content)
	bePushedUserIds = append(bePushedUserIds, ownerIds...)
	// 即使不是任务协作人，也需要被推送卡片
	bePushedUserIds = append(bePushedUserIds, mentionedUserIds...)
	bePushedUserIds = slice.SliceUniqueInt64(bePushedUserIds)

	log.Infof("[PushIssueComment] 要推送的用户ids:%s, issueId: %d", json.ToJsonIgnoreError(bePushedUserIds), issueTrendsBo.IssueId)
	if len(bePushedUserIds) == 0 {
		return
	}
	userNoticeInfos := make([]bo.UserNoticeInfoBo, 0)
	baseUserInfos, err := orgfacade.GetBaseUserInfoBatchRelaxed(orgId, bePushedUserIds)
	if err != nil {
		log.Errorf("[PushIssueComment] err: %v, issueId: %d", err, issueTrendsBo.IssueId)
		return
	}
	// 查询推送开关配置信息
	uConfigResp := orgfacade.GetUserConfigInfoBatch(orgvo.GetUserConfigInfoBatchReqVo{
		OrgId: issueTrendsBo.OrgId,
		Input: orgvo.GetUserConfigInfoBatchReqVoInput{
			UserIds: bePushedUserIds,
		},
	})
	if uConfigResp.Failure() {
		log.Errorf("[PushIssueComment] err: %v, issueId: %d", uConfigResp.Error(), issueTrendsBo.IssueId)
		return
	}
	uConfigMap := make(map[int64]bo.UserConfigBo, len(uConfigResp.Data))
	for _, item := range uConfigResp.Data {
		uConfigMap[item.UserId] = item
	}

	userMap := maps.NewMap("UserId", baseUserInfos)
	for _, userId := range bePushedUserIds {
		if operatorId == userId {
			continue
		}
		if isBeAt, _ := slice.Contain(beAtUserIds, userId); !isBeAt {
			// 如果是负责人
			if exist, _ := slice.Contain(ownerIds, userId); exist {
				if tmpConfig, ok := uConfigMap[userId]; ok {
					if !CheckUserPushConfigIsTurnOn(tmpConfig, consts.PersonalPushConfigForOwnerRange) ||
						!CheckUserPushConfigIsTurnOn(tmpConfig, consts.PersonalPushConfigForIssueBeComment) {
						continue
					}
				}
			} else {
				if tmpConfig, ok := uConfigMap[userId]; ok {
					if !CheckUserPushConfigIsTurnOn(tmpConfig, consts.PersonalPushConfigForCollaborateRange) ||
						!CheckUserPushConfigIsTurnOn(tmpConfig, consts.PersonalPushConfigForIssueBeComment) {
						continue
					}
				}
			}
		} else {
			// 被 at，一定会收到通知
		}
		if baseUserInfoInterface, ok := userMap[userId]; ok {
			if baseUserInfo, ok := baseUserInfoInterface.(bo.BaseUserInfoBo); ok {
				userNoticeInfos = append(userNoticeInfos, bo.UserNoticeInfoBo{
					UserId:    baseUserInfo.UserId,
					OutUserId: baseUserInfo.OutUserId,
					Name:      baseUserInfo.Name,
				})
			}
		}
	}

	if userNoticeInfos != nil && len(userNoticeInfos) > 0 {
		if pushType == consts.PushTypeIssueComment {
			cardForSelf, cardForOther, err := GetIssueCardComment(issueNoticeBo, infoBox, content)
			if err != nil {
				log.Errorf("[PushIssueComment] GetIssueCardComment err: %v", err)
				return
			}
			IssueNoticePush(baseOrgInfo.BaseOrgInfo.SourceChannel, issueTrendsBo.OrgId, userNoticeInfos, cardForSelf, cardForOther, mentionedUserIds, nil)
		} else if pushType == consts.PushTypeIssueRemarkRemind {
			cardForSelf, cardForOther, err := GetIssueCardRemarkAt(issueNoticeBo, content)
			if err != nil {
				log.Errorf("[PushIssueComment] GetIssueCardRemarkAt err: %v", err)
				return
			}
			IssueNoticePush(baseOrgInfo.BaseOrgInfo.SourceChannel, issueTrendsBo.OrgId, userNoticeInfos, cardForSelf, cardForOther, mentionedUserIds, nil)
		}
	}
}

// HandlePushCardForUpdateIssue 机器人个人推送-处理任务更新时，机器人个人推送卡片。
// 有成员更新时，向被提及的成员发送第一人称视角标题的卡片；向其他非提及的成员发送非第一人称视角的卡片。
func HandlePushCardForUpdateIssue(issueTrendsBo bo.IssueTrendsBo, sourceChannel string) {
	infoObj, err := GetBaseInfoForCardPush(&issueTrendsBo)
	if err != nil {
		log.Errorf("[HandlePushCardForUpdateIssue] err: %v, issueId: %d", err, issueTrendsBo.IssueId)
		return
	}

	if len(infoObj.TableColumnMap) < 1 {
		//log.Errorf("[HandlePushCardForUpdateIssue] 没有需要推送的字段, issueId:%v", issueTrendsBo.IssueId)
		return
	}

	// 查询被推送的成员：负责人、协作人
	ownerIds := infoObj.IssueInfo.OwnerIdI64
	collaboratorIds := infoObj.CollaboratorIds
	// 根据配置，过滤负责人、协作人
	filteredOwnerIdArr, filteredCollaborateIdArr, err := filterBePushedUserIdsForUpdateIssue(&issueTrendsBo, ownerIds,
		collaboratorIds)
	if err != nil {
		log.Errorf("[HandlePushCardForUpdateIssue] filterBePushedUserIdsForUpdateIssue err: %v, issueId: %d", err, issueTrendsBo.IssueId)
		return
	}
	allUserId := make([]int64, 0, len(filteredOwnerIdArr)+len(filteredCollaborateIdArr))
	allUserId = append(allUserId, filteredOwnerIdArr...)
	allUserId = append(allUserId, filteredCollaborateIdArr...)
	allUserId = slice.SliceUniqueInt64(allUserId)
	log.Infof("[HandlePushCardForUpdateIssue] collaborate user count: %d", len(allUserId))

	cardMeta := &pushV1.TemplateCard{}

	// 如果只有新增成员的变动，则只发送“被添加到记录”的推送
	isOnlyForAddMember := false
	if len(allUserId) > 0 {
		filteredUsers, err := getUserListByIdsForPushCard(issueTrendsBo.OrgId, allUserId)
		if err != nil {
			log.Errorf("[HandlePushCardForUpdateIssue] getUserListByIdsForPushCard err: %v, issueId: %d", err, issueTrendsBo.IssueId)
			return
		}
		addMemberColumns, notAddMemberColumns := domain.AssemblyUpdateContentForUpdateIssue(infoObj.TableColumnMap, &issueTrendsBo, cardMeta)

		if len(cardMeta.Divs) < 1 {
			//log.Infof("[HandlePushCardForUpdateIssue] 没有需要推送的消息卡片")
			return
		}

		log.Infof("[HandlePushCardForUpdateIssue] issueId: %d, addMemberColumns: %s, notAddMemberColumns: %s",
			issueTrendsBo.IssueId,
			json.ToJsonIgnoreError(addMemberColumns),
			json.ToJsonIgnoreError(notAddMemberColumns))
		if len(notAddMemberColumns) < 1 && len(addMemberColumns) == 1 {
			// 因为需要向被添加的人专门推送“被添加到记录”的卡片，因此不予推送“任务被更新”的卡片（去重）
			// 去重条件是：一次更新请求中，有且只有一个成员类型的列被更新（新增成员）
			isOnlyForAddMember = true
		}

		ignoreMemberIds := make([]int64, 0)
		if isOnlyForAddMember {
			mentionedUserGroup, _ := GetRelateCustomFiledMemberIdsForUpdateIssue(&issueTrendsBo)
			for _, uidArr := range mentionedUserGroup {
				if len(uidArr) > 0 {
					ignoreMemberIds = append(ignoreMemberIds, uidArr...)
					break
				}
			}
		}
		// 无需向编辑者自己推送卡片
		ignoreMemberIds = append(ignoreMemberIds, issueTrendsBo.OperatorId)

		cardForSelf, cardForOther, err := GetIssueCardUpdateIssue(&issueTrendsBo, infoObj, cardMeta)
		if err != nil {
			if err == errs.CardTitleEmpty {
				return
			}
			log.Errorf("[HandlePushCardForUpdateIssue] GetCollaboratorIdsForOneIssue err: %v, issueId: %d", err, issueTrendsBo.IssueId)
			return
		}

		// 进行卡片推送（任务被更新推送）
		IssueNoticePush(sourceChannel, issueTrendsBo.OrgId, filteredUsers, cardForSelf, cardForOther, nil, ignoreMemberIds)
	}

	// 成员字段变更，向变更的目标人员发送卡片提醒（“被添加到记录”推送）
	if err := SendCardToMentionedUserIdsForUpdateIssue(&issueTrendsBo, infoObj, sourceChannel); err != nil {
		log.Errorf("[HandlePushCardForUpdateIssue] SendCardToMentionedUserIdsForUpdateIssue err: %v, issueId: %d", err, issueTrendsBo.IssueId)
		return
	}
}

// getUserListByIdsForPushCard 根据 userIds 查询对应的成员转成 bo.UserNoticeInfoBo 对象
func getUserListByIdsForPushCard(orgId int64, userIds []int64) ([]bo.UserNoticeInfoBo, errs.SystemErrorInfo) {
	users := make([]bo.UserNoticeInfoBo, 0)
	if len(userIds) < 1 {
		return users, nil
	}
	baseUserInfos, err := orgfacade.GetBaseUserInfoBatchRelaxed(orgId, userIds)
	if err != nil {
		log.Errorf("[GetUserInfoMapByIds] GetBaseUserInfoBatchRelaxed err: %v, orgId: %d", err, orgId)
		return users, err
	}
	for _, user := range baseUserInfos {
		users = append(users, bo.UserNoticeInfoBo{
			UserId:       user.UserId,
			OutUserId:    user.OutUserId,
			Name:         user.Name,
			OutOrgUserId: user.OutOrgUserId,
		})
	}

	return users, nil
}

// 更新任务时，推送卡片，根据用户配置的推送开关进行过滤
func filterBePushedUserIdsForUpdateIssue(issueTrendsBo *bo.IssueTrendsBo, ownerIds, collaboratorIds []int64) ([]int64, []int64, errs.SystemErrorInfo) {
	filteredOwnerUidArr := make([]int64, 0, len(ownerIds))
	filteredCollaborateUidArr := make([]int64, 0, len(collaboratorIds))
	allUserIds := append(collaboratorIds, ownerIds...)
	if len(allUserIds) < 1 {
		return filteredOwnerUidArr, filteredCollaborateUidArr, nil
	}
	// 查询 user config
	uConfigResp := orgfacade.GetUserConfigInfoBatch(orgvo.GetUserConfigInfoBatchReqVo{
		OrgId: issueTrendsBo.OrgId,
		Input: orgvo.GetUserConfigInfoBatchReqVoInput{
			UserIds: allUserIds,
		},
	})
	if uConfigResp.Failure() {
		log.Errorf("[filterBePushedUserIdsForUpdateIssue] err: %v, issueId: %d", uConfigResp.Error(), issueTrendsBo.IssueId)
		return filteredOwnerUidArr, filteredCollaborateUidArr, uConfigResp.Error()
	}
	uConfigMap := make(map[int64]bo.UserConfigBo, len(uConfigResp.Data))
	for _, item := range uConfigResp.Data {
		uConfigMap[item.UserId] = item
	}

	uidGroup := map[string][]int64{
		consts.BasicFieldOwnerId:     ownerIds,
		consts.UserTypeOfCollaborate: collaboratorIds,
	}

	for uTypeFlag, userIds := range uidGroup {
		switch uTypeFlag {
		case consts.BasicFieldOwnerId:
			for _, uid := range userIds {
				// 去除操作人，即不向自己推送卡片
				if uid == issueTrendsBo.OperatorId {
					continue
				}
				tmpConfig, ok := uConfigMap[uid]
				if !ok {
					continue
				}
				if !CheckUserPushConfigIsTurnOn(tmpConfig, consts.PersonalPushConfigForOwnerRange) ||
					!CheckUserPushConfigIsTurnOn(tmpConfig, consts.PersonalPushConfigForIssueUpdate) {
					continue
				}
				filteredOwnerUidArr = append(filteredOwnerUidArr, uid)
			}
		case consts.UserTypeOfCollaborate:
			for _, uid := range userIds {
				if uid == issueTrendsBo.OperatorId {
					continue
				}
				tmpConfig, ok := uConfigMap[uid]
				if !ok {
					continue
				}
				if !CheckUserPushConfigIsTurnOn(tmpConfig, consts.PersonalPushConfigForCollaborateRange) ||
					!CheckUserPushConfigIsTurnOn(tmpConfig, consts.PersonalPushConfigForIssueUpdate) {
					continue
				}
				filteredCollaborateUidArr = append(filteredCollaborateUidArr, uid)
			}
		}
	}

	return filteredOwnerUidArr, filteredCollaborateUidArr, nil
}

// CheckUserPushConfigIsTurnOn 检查用户的个人卡片推送配置项对应的开关是否开启
// switchFlag 开关标识
func CheckUserPushConfigIsTurnOn(config bo.UserConfigBo, switchFlag string) bool {
	var configVal int
	switch switchFlag {
	case consts.PersonalPushConfigForMyDailyReport:
		configVal = config.DailyReportMessageStatus
	case consts.PersonalPushConfigForProDailyReport:
		configVal = config.DailyProjectReportMessageStatus
	case consts.PersonalPushConfigForOwnerRange:
		configVal = config.OwnerRangeStatus
	case consts.PersonalPushConfigForCollaborateRange:
		configVal = config.CollaborateMessageStatus
	case consts.PersonalPushConfigForIssueUpdate:
		configVal = config.ModifyMessageStatus
	case consts.PersonalPushConfigForIssueBeComment:
		configVal = config.CommentAtMessageStatus
	case consts.PersonalPushConfigForRelateIssueTrends:
		configVal = config.RelationMessageStatus
	case consts.PersonalPushConfigForRemindMessage: // 被添加到记录
		configVal = config.RemindMessageStatus
	default:
	}

	return configVal == 1
}

// SendCardToMentionedUserIdsForUpdateIssue 被添加到记录的相关推送。修改了任务成员时，取出对应修改的目标成员，并向这些目标成员推送卡片
// 例如：新增了负责人 1001，则拿出 1001，告知他（1001）（当事人）被添加为负责人
func SendCardToMentionedUserIdsForUpdateIssue(issueTrendsBo *bo.IssueTrendsBo, infoObj *projectvo.BaseInfoBoxForIssueCard, sourceChannel string) errs.SystemErrorInfo {
	mentionedUserGroup, _ := GetRelateCustomFiledMemberIdsForUpdateIssue(issueTrendsBo)
	allUserIds := make([]int64, 0)
	for _, addUserIds := range mentionedUserGroup {
		allUserIds = append(allUserIds, addUserIds...)
	}
	userMap := make(map[int64]bo.UserNoticeInfoBo, 0)
	userInfoArr, err := getUserListByIdsForPushCard(issueTrendsBo.OrgId, allUserIds)
	if err != nil {
		log.Errorf("[SendCardToMentionedUserIdsForUpdateIssue] err: %v, issueId: %d", err, issueTrendsBo.IssueId)
		return err
	}
	// 检查是否勾选了“被添加到记录”配置
	// 查询 user config
	uConfigResp := orgfacade.GetUserConfigInfoBatch(orgvo.GetUserConfigInfoBatchReqVo{
		OrgId: issueTrendsBo.OrgId,
		Input: orgvo.GetUserConfigInfoBatchReqVoInput{
			UserIds: allUserIds,
		},
	})
	if uConfigResp.Failure() {
		log.Errorf("[SendCardToMentionedUserIdsForUpdateIssue] err: %v, issueId: %d", uConfigResp.Error(), issueTrendsBo.IssueId)
		return uConfigResp.Error()
	}
	uConfigMap := make(map[int64]bo.UserConfigBo, len(uConfigResp.Data))
	for _, item := range uConfigResp.Data {
		uConfigMap[item.UserId] = item
	}

	for _, item := range userInfoArr {
		userMap[item.UserId] = item
	}
	for columnKey, addUserIds := range mentionedUserGroup {
		if len(addUserIds) < 1 {
			continue
		}
		// 检查是否是协作人列，如果不是，则不推送
		columnIsCollaborate := CheckColumnIsCollaborate(columnKey, infoObj)
		tmpUsers := make([]bo.UserNoticeInfoBo, 0, len(addUserIds))
		for _, uid := range addUserIds {
			if uid < 1 || uid == issueTrendsBo.OperatorId {
				continue
			}
			tmpUserConfig, ok := uConfigMap[uid]
			if !ok {
				continue
			}
			// 未开启“被添加到记录”，则不发送卡片
			if columnKey == consts.BasicFieldOwnerId {
				if !CheckUserPushConfigIsTurnOn(tmpUserConfig, consts.PersonalPushConfigForOwnerRange) ||
					!CheckUserPushConfigIsTurnOn(tmpUserConfig, consts.PersonalPushConfigForRemindMessage) {
					continue
				}
				if user, ok := userMap[uid]; ok {
					tmpUsers = append(tmpUsers, user)
				}
			} else {
				if !columnIsCollaborate {
					continue
				}
				if !CheckUserPushConfigIsTurnOn(tmpUserConfig, consts.PersonalPushConfigForCollaborateRange) ||
					!CheckUserPushConfigIsTurnOn(tmpUserConfig, consts.PersonalPushConfigForRemindMessage) {
					continue
				}
				if user, ok := userMap[uid]; ok {
					tmpUsers = append(tmpUsers, user)
				}
			}
		}
		if len(tmpUsers) > 0 {
			cardMeta, err := GetIssueAddMemberMsg(issueTrendsBo, columnKey, infoObj)
			if err != nil {
				log.Errorf("[SendCardToMentionedUserIdsForUpdateIssue] err: %v, issueId: %d, card: %s", err,
					issueTrendsBo.IssueId, json.ToJsonIgnoreError(cardMeta))
				continue
			}
			IssueNoticePush(sourceChannel, issueTrendsBo.OrgId, tmpUsers, cardMeta, cardMeta, nil, nil)
		}
		log.Infof("[SendCardToMentionedUserIdsForUpdateIssue] 发送卡片，将用户添加为记录成员。columnKey: %s, userCount: %d",
			columnKey, len(tmpUsers))
	}

	return nil
}

// CheckColumnIsCollaborate 检查列是否是协作人列
func CheckColumnIsCollaborate(columnKey string, infoObj *projectvo.BaseInfoBoxForIssueCard) bool {
	// 确认人一定是协作人（强制逻辑）
	if columnKey == consts.BasicFieldAuditorIds {
		return true
	}
	curColumn, ok := infoObj.TableColumnMap[columnKey]
	if !ok {
		return false
	}

	return curColumn.Field.Props.Member.IsCollaborators
}

// GetIssueAddMemberMsg 添加为记录的成员时的卡片组装
func GetIssueAddMemberMsg(issueTrendsBo *bo.IssueTrendsBo, columnKey string, infoObj *projectvo.BaseInfoBoxForIssueCard) (*pushV1.TemplateCard, errs.SystemErrorInfo) {
	curColumn, ok := infoObj.TableColumnMap[columnKey]
	if !ok {
		err := errs.CustomFieldNotExist
		log.Errorf("[GetIssueAddMemberMsg] err: %v, columnKey: %s", err, columnKey)
		return nil, err
	}
	if infoObj.IssueInfo.Title == "" {
		return nil, errs.CardTitleEmpty
	}
	columnDisplayName := str.TruncateColumnName(curColumn)

	msg := card.GetCardIssueAddMember(issueTrendsBo, infoObj, columnDisplayName)
	return msg, nil
}

// GetRelateCustomFiledMemberIdsForUpdateIssue 取出成员字段变更的目标成员，用于向这些成员额外推送卡片信息
// 只要有新增的成员，就会向这批目标成员额外推送卡片；删除成员暂时不额外推送
// 负责人：只有开启了我负责的，且开启了“添加到记录”，才会返回
// 其他成员列：只有开启了我协作的，且开启了“添加到记录”，才会返回
func GetRelateCustomFiledMemberIdsForUpdateIssue(issueTrendsBo *bo.IssueTrendsBo) (map[string][]int64, errs.SystemErrorInfo) {
	resultMap := make(map[string][]int64, 0)
	// 负责人
	if issueTrendsBo.UpdateOwner {
		if _, add, _ := int642.CompareSliceAddDelInt64(issueTrendsBo.AfterOwner, issueTrendsBo.BeforeOwner); len(add) > 0 {
			resultMap[consts.BasicFieldOwnerId] = add
		}
	}
	// 确认人
	if issueTrendsBo.UpdateAuditor {
		if _, add, _ := int642.CompareSliceAddDelInt64(issueTrendsBo.AfterChangeAuditors, issueTrendsBo.BeforeChangeAuditors); len(add) > 0 {
			resultMap[consts.BasicFieldAuditorIds] = add
		}
	}
	// 关注人
	if issueTrendsBo.UpdateFollower {
		if _, add, _ := int642.CompareSliceAddDelInt64(issueTrendsBo.AfterChangeFollowers, issueTrendsBo.BeforeChangeFollowers); len(add) > 0 {
			resultMap[consts.BasicFieldFollowerIds] = add
		}
	}
	// 其他成员类型字段
	for _, changeItem := range issueTrendsBo.Ext.ChangeList {
		if changeItem.FieldType != consts.LcColumnFieldTypeMember {
			continue
		}
		if _, add, _ := int642.CompareSliceAddDelString(changeItem.NewUserIdsOrDeptIdsValue, changeItem.OldUserIdsOrDeptIdsValue); len(add) > 0 {
			addInt64Ids, err := businees.LcMemberToUserIdsWithError(add)
			if err != nil {
				log.Errorf("[GetRelateCustomFiledMemberIdsForUpdateIssue] LcMemberToUserIdsWithError err: %v, "+
					"data: %s", err, json.ToJsonIgnoreError(add))
			}
			resultMap[changeItem.Field] = addInt64Ids
		}
	}

	return resultMap, nil
}

// HandlePushCardForDefaultEvent 机器人个人推送-处理第三方卡片推送服务，如钉钉卡片、飞书卡片（目前只支持飞书卡片）。推送触发的事件如任务创建等
// 目前，该逻辑触发的事件有 consts.PushTypeCreateIssue, consts.PushTypeDeleteIssue, consts.PushTypeRecoverIssue, consts.PushTypeUpdateRelationIssue
func HandlePushCardForDefaultEvent(issueTrendsBo bo.IssueTrendsBo, sourceChannel string) {
	orgId := issueTrendsBo.OrgId
	pushType := issueTrendsBo.PushType
	issueNoticeBo := bo.IssueNoticeBo{}
	err := copyer.Copy(issueTrendsBo, &issueNoticeBo)
	if err != nil {
		log.Errorf("[HandlePushCardForDefaultEvent] copy err: %v, issueId: %d", err, issueTrendsBo.IssueId)
		return
	}

	bePushedUserInfos := GetNormalUserIds(issueNoticeBo, sourceChannel, PushNoticeTargetTypeModifyMsg)
	log.Infof("[HandlePushCardForDefaultEvent] 任务推送，需要推送的用户信息列表为 %s", json.ToJsonIgnoreError(bePushedUserInfos))

	if len(bePushedUserInfos) == 0 {
		log.Infof("[HandlePushCardForDefaultEvent] 需要推送的人员数为0，不需要推送，消息%s", json.ToJsonIgnoreError(issueNoticeBo))
	}

	// 创建新任务，其实就是分配一个任务的负责人
	if pushType == consts.PushTypeCreateIssue {
		// 创建事件，推送负责人专属消息
		ownerIds := issueNoticeBo.AfterOwner
		// 如果变动后的负责人和操作人不是同一个人，对该负责人推送负责人消息
		// 处理负责人信息
		err := dealOwnerInfo(issueTrendsBo, []int64{}, issueNoticeBo.AfterOwner, orgId, sourceChannel, []int64{})
		if err != nil {
			log.Error(err)
			return
		}
		// 去掉负责人
		bePushedNormalUerInfos := make([]bo.UserNoticeInfoBo, 0)
		for _, bePushedUserInfo := range bePushedUserInfos {
			if ok, _ := slice.Contain(ownerIds, bePushedUserInfo.UserId); !ok {
				bePushedNormalUerInfos = append(bePushedNormalUerInfos, bePushedUserInfo)
			}
		}
		bePushedUserInfos = bePushedNormalUerInfos
	}

	if len(bePushedUserInfos) > 0 {
		bePushedMsg, err := GetIssueCardNormal(issueTrendsBo)
		log.Infof("[HandlePushCardForDefaultEvent] 飞书任务推送，需要推送的消息为 %s", json.ToJsonIgnoreError(bePushedMsg))
		if err != nil {
			log.Error(err)
			return
		}

		IssueNoticePush(sourceChannel, orgId, bePushedUserInfos, bePushedMsg, bePushedMsg, nil, nil)
	}
	////删除任务附带的子任务通知
	//if pushType == consts.PushTypeDeleteIssue && len(issueNoticeBo.IssueChildren) > 0 {
	//	PushMsgForDeleteRelateChildrenIssue(issueNoticeBo)
	//}
}

// dealOwnerInfo 添加负责人时的卡片推送
func dealOwnerInfo(issueTrendsBo bo.IssueTrendsBo, beforeOwnerId []int64, afterOwnerId []int64, orgId int64, sourceChannel string, collaborateUserIds []int64) errs.SystemErrorInfo {
	issueNoticeBo := bo.IssueNoticeBo{}
	oriErr := copyer.Copy(issueTrendsBo, &issueNoticeBo)
	if oriErr != nil {
		log.Errorf("[dealOwnerInfo] copy err: %v, issueId: %d", oriErr, issueTrendsBo.IssueId)
		return errs.BuildSystemErrorInfo(errs.ObjectCopyError, oriErr)
	}

	added := []int64{}
	for _, i2 := range afterOwnerId {
		if i2 == issueNoticeBo.OperatorId {
			continue
		}
		if ok, _ := slice.Contain(beforeOwnerId, i2); ok {
			continue
		}
		added = append(added, i2)
	}
	//给后续增加的人推送新记录提醒
	if len(added) == 0 {
		return nil
	}

	issueBos, issuesInfoErr := domain.GetIssueInfosLc(orgId, 0, []int64{issueNoticeBo.IssueId})
	if issuesInfoErr != nil {
		log.Error(issuesInfoErr)
		return issuesInfoErr
	}
	if len(issueBos) < 1 {
		log.Errorf("[dealOwnerInfo] not found issue issueId:%v", issueNoticeBo.IssueId)
		return errs.IssueNotExist
	}
	issueInfo := issueBos[0]
	//任务标题为空，则不需要推送
	if issueInfo.Title == "" {
		return errs.CardTitleEmpty
	}
	afterOwnerInfos, err := orgfacade.GetBaseUserInfoBatchRelaxed(orgId, added)
	if err != nil {
		log.Error(err)
		return err
	}

	// 如果卡片接收人是新增的负责人，则对应的标题显示会有所不同，即当事人和非当事人的卡片标题有区别
	msgForSelf, msgForOther, err := GetIssueCardOwnerChange(&issueTrendsBo, issueNoticeBo, "owner")
	if err != nil {
		log.Errorf("[dealOwnerInfo] err: %v", err)
		return err
	}

	ownerIds := make([]int64, 0, len(afterOwnerInfos))
	for _, afterOwnerInfo := range afterOwnerInfos {
		ownerIds = append(ownerIds, afterOwnerInfo.UserId)
	}
	userConfigArrResp := orgfacade.GetUserConfigInfoBatch(orgvo.GetUserConfigInfoBatchReqVo{
		OrgId: orgId,
		Input: orgvo.GetUserConfigInfoBatchReqVoInput{
			UserIds: ownerIds,
		},
	})
	if userConfigArrResp.Failure() {
		log.Errorf("[dealOwnerInfo] err: %v, ownerIds: %s", err, json.ToJsonIgnoreError(ownerIds))
		return err
	}
	userConfigMap := make(map[int64]bo.UserConfigBo, len(ownerIds))
	for _, tmpConfig := range userConfigArrResp.Data {
		userConfigMap[tmpConfig.UserId] = tmpConfig
	}
	var noticeUsers []bo.UserNoticeInfoBo
	for _, afterOwnerInfo := range afterOwnerInfos {
		userConfig, ok := userConfigMap[afterOwnerInfo.UserId]
		if !ok {
			continue
		}
		if userConfig.OwnerRangeStatus != 1 || userConfig.ModifyMessageStatus != 1 {
			continue
		}
		noticeUsers = append(noticeUsers, bo.UserNoticeInfoBo{
			UserId:       afterOwnerInfo.UserId,
			OutUserId:    afterOwnerInfo.OutUserId,
			Name:         afterOwnerInfo.Name,
			OutOrgUserId: afterOwnerInfo.OutOrgUserId,
		})
	}
	IssueNoticePush(sourceChannel, orgId, noticeUsers, msgForSelf, msgForOther, added, nil)

	return nil
}

// GetNormalUserIdsWithFilter filter: 需要推送的人，限定推送范围
// collaborateUidArr 协作人 uid 列表
func GetNormalUserIdsWithFilter(orgId int64, operatorId int64, ownerIds []int64, participants []int64, followers []int64, collaborateUidArr []int64, sourceChannel string, targetType int, filter map[int64]bool) []bo.UserNoticeInfoBo {
	ownerArray := ownerIds
	userIdsList := [][]int64{ownerArray, participants, followers, collaborateUidArr}

	//返回链表
	noticeUserList := dealIssueNoticeUserIdsList(userIdsList, operatorId, orgId, targetType)

	//转换array并去重
	bePushedUserIds := make([]int64, noticeUserList.Len())
	i := 0
	for e := noticeUserList.Front(); e != nil; e = e.Next() {
		bePushedUserIds[i] = e.Value.(int64)
		i++
	}
	bePushedUserIds = slice.SliceUniqueInt64(bePushedUserIds)

	if filter != nil && len(filter) > 0 {
		userIds := make([]int64, 0)
		for _, bePushedUserId := range bePushedUserIds {
			if _, ok := filter[bePushedUserId]; ok {
				userIds = append(userIds, bePushedUserId)
			}
		}
		bePushedUserIds = userIds
	}

	userNoticeInfos := make([]bo.UserNoticeInfoBo, 0)
	if len(bePushedUserIds) > 0 {
		baseUserInfos, err := orgfacade.GetBaseUserInfoBatchRelaxed(orgId, bePushedUserIds)
		if err != nil {
			log.Error(err)
			return userNoticeInfos
		}
		userMap := maps.NewMap("UserId", baseUserInfos)
		for _, userId := range bePushedUserIds {
			if baseUserInfoInterface, ok := userMap[userId]; ok {
				if baseUserInfo, ok := baseUserInfoInterface.(bo.BaseUserInfoBo); ok {
					userNoticeInfos = append(userNoticeInfos, bo.UserNoticeInfoBo{
						UserId:       baseUserInfo.UserId,
						OutUserId:    baseUserInfo.OutUserId,
						Name:         baseUserInfo.Name,
						OutOrgUserId: baseUserInfo.OutOrgUserId,
					})
				}
			}
		}
	}
	return userNoticeInfos
}

func GetNormalUserIds(issueNoticeBo bo.IssueNoticeBo, sourceChannel string, targetType int) []bo.UserNoticeInfoBo {
	collaboratorIds, err := tablefacade.GetDataCollaborateUserIds(issueNoticeBo.OrgId, 0, issueNoticeBo.DataId)
	if err != nil {
		log.Errorf("[GetNormalUserIds] GetDataCollaborateUserIds err: %v, DataId: %d", err, issueNoticeBo.DataId)
		return nil
	}
	return GetNormalUserIdsWithFilter(issueNoticeBo.OrgId, issueNoticeBo.OperatorId, issueNoticeBo.BeforeOwner,
		issueNoticeBo.BeforeChangeParticipants, issueNoticeBo.BeforeChangeFollowers, collaboratorIds, sourceChannel, targetType, nil)
}

func dealIssueNoticeUserIdsList(userIdsList [][]int64, operatorId int64, orgId int64, targetType int) *list.List {
	noticeUserList := list.New()

	for i, userIds := range userIdsList {
		if userIds != nil {
			dealIssueNoticeUserIds(i, userIds, operatorId, orgId, targetType, noticeUserList)
		}
	}
	return noticeUserList
}

//func delNoticeUserIdsList(userIds []int64, operatorId int64, orgId int64, targetType int) []int64 {
//	bePushedUserIds := make([]int64, 0)
//	for _, userId := range userIds {
//		if userId == operatorId {
//			continue
//		}
//		//去重
//		if ok, _ := slice.Contain(bePushedUserIds, userId); ok {
//			continue
//		}
//		userConfig, err := orgfacade.GetUserConfigInfoRelaxed(orgId, userId)
//		if err != nil {
//			log.Errorf("获取%d用户配置失败, %v", userId, err)
//			continue
//		}
//		if userPushTargetConfigContinueFlag(targetType, userConfig) {
//			continue
//		}
//		bePushedUserIds = append(bePushedUserIds, userId)
//	}
//
//	return bePushedUserIds
//}

func dealIssueNoticeUserIds(i int, userIds []int64, operatorId int64, orgId int64, targetType int, noticeUserList *list.List) {
	for _, userId := range userIds {
		// 获取用户信息时 id 为 0 的可以跳过。
		if userId == operatorId || userId == 0 {
			continue
		}
		userConfig, err := orgfacade.GetUserConfigInfoRelaxed(orgId, userId)
		if err != nil {
			log.Errorf("获取%d用户配置失败, %v", userId, err)
			continue
		}
		if userPushRangeConfigContinueFlag(i, userConfig) {
			continue
		}
		if userPushTargetConfigContinueFlag(targetType, userConfig) {
			continue
		}
		noticeUserList.PushBack(userId)
	}
}

// IssueNoticePush 任务通知卡片推送
// userInfos 被通知的人员
// msgForSelf 被 at 的卡片信息（当事人）
// msgForOther 非当事人的卡片信息。如果不存在，可以传与发送给当事人一样的卡片。
// mentionedUserIds 当事人的 id 列表。这些人接受到的卡片标题时第一人称（我）的描述
// ignoreMemberIds 忽略推送的成员 id 列表
func IssueNoticePush(sourceChannel string, orgId int64, userInfos []bo.UserNoticeInfoBo, cardForSelf, cardForOther *pushV1.TemplateCard, mentionedUserIds, ignoreMemberIds []int64) {
	if mentionedUserIds == nil {
		mentionedUserIds = make([]int64, 0)
	}
	if ignoreMemberIds == nil {
		ignoreMemberIds = make([]int64, 0)
	}
	openIds := make([]string, 0, len(userInfos))
	openIdMap := make(map[int64]string, len(userInfos))
	for _, userInfo := range userInfos {
		if exist, _ := slice.Contain(ignoreMemberIds, userInfo.UserId); exist {
			continue
		}
		if userInfo.OutUserId != "" {
			openIds = append(openIds, userInfo.OutUserId)
			openIdMap[userInfo.UserId] = userInfo.OutUserId
		}
	}
	// 区分出当事人用户和非当事人用户
	beAtOpenIds := make([]string, 0, len(userInfos))
	otherOpenIds := make([]string, 0, len(userInfos))
	for uid, openId := range openIdMap {
		if exist, _ := slice.Contain(mentionedUserIds, uid); exist {
			beAtOpenIds = append(beAtOpenIds, openId)
		} else {
			otherOpenIds = append(otherOpenIds, openId)
		}
	}
	orgBaseInfo, err := orgfacade.GetBaseOrgInfoRelaxed(orgId)
	if err != nil {
		log.Errorf("[IssueNoticeFeiShuToSelfOther] 组织外部信息不存在 %v, orgId: %d", err, orgId)
		return
	}
	log.Infof("[IssueNoticePush] beAtOpenIds:%v, otherOpenIds:%v, cardForSelf:%v, cardForOther:%v",
		beAtOpenIds, otherOpenIds, json.ToJsonIgnoreError(cardForSelf), json.ToJsonIgnoreError(cardForOther))

	cardSelf := &commonvo.PushCard{
		OrgId:         orgId,
		OutOrgId:      orgBaseInfo.OutOrgId,
		SourceChannel: sourceChannel,
		OpenIds:       beAtOpenIds,
		CardMsg:       cardForSelf,
	}
	cardOther := &commonvo.PushCard{
		OrgId:         orgId,
		OutOrgId:      orgBaseInfo.OutOrgId,
		SourceChannel: sourceChannel,
		OpenIds:       otherOpenIds,
		CardMsg:       cardForOther,
	}
	errSys := card.PushCard(orgId, cardSelf)
	if errSys != nil {
		log.Errorf("[IssueNoticePush] err:%v", errSys)
	}
	errSys = card.PushCard(orgId, cardOther)
	if errSys != nil {
		log.Errorf("[IssueNoticePush] err:%v", errSys)
	}
	//card.SendCardMeta(sourceChannel, orgBaseInfo.OutOrgId, cardForSelf, beAtOpenIds)
	//card.SendCardMeta(sourceChannel, orgBaseInfo.OutOrgId, cardForOther, otherOpenIds)
}
