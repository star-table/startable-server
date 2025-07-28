package domain

import (
	"errors"
	"strconv"
	"time"

	tablePb "gitea.bjx.cloud/LessCode/interface/golang/table/v1"
	vo2 "gitea.bjx.cloud/allstar/feishu-sdk-golang/core/model/vo"
	platform_sdk "gitea.bjx.cloud/allstar/platform-sdk"
	sdk_const "gitea.bjx.cloud/allstar/platform-sdk/consts"
	"gitea.bjx.cloud/allstar/platform-sdk/sdk_interface"
	vo3 "gitea.bjx.cloud/allstar/platform-sdk/vo"
	"github.com/star-table/startable-server/app/facade/idfacade"
	"github.com/star-table/startable-server/app/facade/orgfacade"
	"github.com/star-table/startable-server/app/service/projectsvc/dao"
	"github.com/star-table/startable-server/app/service/projectsvc/po"
	"github.com/star-table/startable-server/common/core/businees"
	"github.com/star-table/startable-server/common/core/consts"
	"github.com/star-table/startable-server/common/core/errs"
	"github.com/star-table/startable-server/common/core/util"
	"github.com/star-table/startable-server/common/core/util/json"
	"github.com/star-table/startable-server/common/core/util/slice"
	int64Util "github.com/star-table/startable-server/common/core/util/slice/int64"
	"github.com/star-table/startable-server/common/core/util/uuid"
	"github.com/star-table/startable-server/common/extra/feishu"
	"github.com/star-table/startable-server/common/library/cache"
	"github.com/star-table/startable-server/common/library/db/mysql"
	"github.com/star-table/startable-server/common/model/bo"
	"github.com/star-table/startable-server/common/model/vo"
	"upper.io/db.v3"
)

var VersionV3 = "v3"
var VersionV4 = "v4"

// 更新日历订阅者
func UpdateCalendarAttendees(orgId int64, calendarId string, creatorUserId int64, addMembers []int64, delMembers []int64, projectId int64) {
	if calendarId == "" {
		projectCalendarInfo, err := GetProjectCalendarInfo(orgId, projectId)
		if err != nil {
			log.Error(err)
			return
		}
		if projectCalendarInfo.IsSyncOutCalendar != consts.IsSyncOutCalendar || projectCalendarInfo.CalendarId == consts.BlankString {
			log.Error("无对应项目日历或未设置导出日历")
			return
		}
		calendarId = projectCalendarInfo.CalendarId
	}
	orgBaseInfo, err := orgfacade.GetBaseOrgInfoRelaxed(orgId)
	if err != nil {
		log.Errorf("组织外部信息不存在 %v", err)
		return
	}

	if orgBaseInfo.OutOrgId != "" {
		allUserIds := make([]int64, 0, len(addMembers)+len(delMembers)+1)
		allUserIds = append(allUserIds, creatorUserId)
		allUserIds = append(allUserIds, addMembers...)
		allUserIds = append(allUserIds, delMembers...)
		allUserIds = slice.SliceUniqueInt64(allUserIds)
		baseInfoMap, err := getBaseUserInfoMap(orgId, allUserIds)
		if err != nil {
			log.Errorf("[UpdateCalendarAttendees] getBaseUserInfoMap, orgId:%v, userIds:%v, err:%v", orgId, allUserIds, err)
			return
		}
		client, platformUserId, err := getPlatformClientAndCreatorIdWithInfoMap(orgId, orgBaseInfo.OutOrgId, orgBaseInfo.SourceChannel, creatorUserId, baseInfoMap)
		if err != nil {
			log.Errorf("[UpdateCalendarAttendees] GetClient, sourceChannel:%v, outOrgId:%v, err:%v", orgBaseInfo.SourceChannel, orgBaseInfo.OutOrgId, err)
			return
		}
		_, err2 := client.UpdateCalendarAttend(&vo3.UpdateCalendarAttendReq{
			UserId:        platformUserId,
			CalendarId:    calendarId,
			AddUserIds:    getPlatformUserIdsWithInfoMap(orgBaseInfo.SourceChannel, addMembers, baseInfoMap),
			DeleteUserIds: getPlatformUserIdsWithInfoMap(orgBaseInfo.SourceChannel, delMembers, baseInfoMap),
		})
		if err2 != nil {
			log.Errorf("[UpdateCalendarAttendees] UpdateCalendarAttend, orgId:%v, userIds:%v, err:%v", orgId, allUserIds, err2)
			return
		}
	}
}

// 创建日历
func CreateCalendar(isSyncOutCalendar *int, orgId int64, projectId int64, userId int64, addMemberIds []int64) {
	// 2、0 表示不推送。
	if isSyncOutCalendar != nil && *isSyncOutCalendar > 0 && *isSyncOutCalendar != consts.IsNotSyncOutCalendar {
		orgBaseInfo, err := orgfacade.GetBaseOrgInfoRelaxed(orgId)
		if err != nil {
			log.Errorf("组织外部信息不存在 %v", err)
		}
		if orgBaseInfo.OutOrgId == "" {
			return
		}

		//防止重复插入
		uid := uuid.NewUuid()
		projectIdStr := strconv.FormatInt(projectId, 10)
		lockKey := consts.CreateCalendarLock + projectIdStr
		suc, err2 := cache.TryGetDistributedLock(lockKey, uid)
		if err2 != nil {
			log.Errorf("获取%s锁时异常 %v", lockKey, err2)
			return
		}
		if suc {
			defer func() {
				if _, err := cache.ReleaseDistributedLock(lockKey, uid); err != nil {
					log.Error(err)
				}
			}()
		}
		projectCalendarInfo, err := GetProjectCalendarInfo(orgId, projectId)
		if err != nil {
			log.Error(err)
			return
		}
		if projectCalendarInfo.CalendarId != "" {
			log.Info("日历已插入" + projectIdStr)
			return
		}
		//把缓存删掉
		delErr := DeleteProjectCalendarInfo(orgId, projectId)
		if delErr != nil {
			log.Error("删除缓存失败")
			return
		}

		projectInfo, err := GetProjectSimple(orgId, projectId)
		if err != nil {
			log.Error(err)
			return
		}

		client, platformUserId, err := getPlatformClientAndCreatorId(orgId, orgBaseInfo.OutOrgId, orgBaseInfo.SourceChannel, userId)
		if err != nil {
			log.Errorf("[CreateCalendar] GetClient, sourceChannel:%v, outOrgId:%v, err:%v", orgBaseInfo.SourceChannel, orgBaseInfo.OutOrgId, err)
			return
		}

		reply, err2 := client.CreateCalendar(&vo3.CreateCalendarReq{
			UserId:      platformUserId,
			Name:        projectInfo.Name,
			Description: projectInfo.Remark,
		})
		if err2 != nil {
			log.Errorf("[CreateCalendar] GetClient, sourceChannel:%v, outOrgId:%v, err:%v", orgBaseInfo.SourceChannel, orgBaseInfo.OutOrgId, err2)
			return
		}

		//创建关联关系
		memberId, err := idfacade.ApplyPrimaryIdRelaxed(consts.TableProjectRelation)
		if err != nil {
			log.Error(err)
			return
		}
		err1 := dao.InsertProjectRelation(po.PpmProProjectRelation{
			Id:           memberId,
			OrgId:        orgId,
			ProjectId:    projectId,
			RelationType: consts.IssueRelationTypeCalendar,
			RelationCode: reply.CalendarId,
			Creator:      userId,
			CreateTime:   time.Now(),
			IsDelete:     consts.AppIsNoDelete,
			Status:       consts.ProjectMemberEffective,
			Updator:      userId,
			UpdateTime:   time.Now(),
			Version:      1,
		})
		if err1 != nil {
			log.Error(err1)
		}
		_, err2 = GetProjectCalendarInfo(orgId, projectId)
		if err2 != nil {
			log.Error("缓存信息失败" + err2.Error())
		}
		//访问控制
		// 如果未勾选“订阅日历”，则不增加访问控制
		if CheckCalendarIsSyncToSubCalendar(*isSyncOutCalendar) {
			UpdateCalendarAttendees(orgId, reply.CalendarId, userId, addMemberIds, []int64{}, projectId)
		}

		//创建日程（可能创建日历之前就有任务）
		createCalendarEventsBefore(orgId, projectId, userId, false)
	}
}

// 创建日程（可能创建日历之前就有任务）
func createCalendarEventsBefore(orgId, projectId, userId int64, dealAlreadyExist bool) {
	conds := []*tablePb.Condition{
		GetRowsCondition(consts.BasicFieldProjectId, tablePb.ConditionType_equal, projectId, nil),
		GetRowsCondition(consts.BasicFieldRecycleFlag, tablePb.ConditionType_equal, consts.DeleteFlagNotDel, nil),
	}
	issueInfos, err := GetIssueInfosMapLc(orgId, userId, &tablePb.Condition{
		Type:       tablePb.ConditionType_and,
		Conditions: conds,
	}, nil, -1, -1)
	if err != nil {
		log.Errorf("[createCalendarEventsBefore] GetIssueInfosMapLc err:%v, orgId:%v, projectId:%v", err, orgId, projectId)
		return
	}
	if len(issueInfos) == 0 {
		return
	}

	issueBos := make([]*bo.IssueBo, 0, len(issueInfos))
	issueIds := make([]int64, 0, len(issueInfos))
	for _, data := range issueInfos {
		issue, errSys := ConvertIssueDataToIssueBo(data)
		if errSys != nil {
			log.Errorf("[createCalendarEventsBefore]ConvertIssueDataToIssueBo err:%v, orgId:%v, issueId:%v", errSys, orgId, issue.Id)
			return
		}
		issueIds = append(issueIds, issue.Id)
		issueBos = append(issueBos, issue)
	}

	//issuePos := []po.PpmPriIssue{}
	//err := mysql.SelectAllByCond(consts.TableIssue, db.Cond{
	//	consts.TcProjectId: projectId,
	//	consts.TcIsDelete:  consts.AppIsNoDelete,
	//	consts.TcOrgId:     orgId,
	//}, &issuePos)
	//if err != nil {
	//	log.Error(err)
	//	return
	//}
	//if len(issuePos) == 0 {
	//	return
	//}
	//
	//issueIds := []int64{}
	//for _, issueBo := range issuePos {
	//	issueIds = append(issueIds, issueBo.Id)
	//}

	//issuesInfo, issueInfoErr := GetIssueInfosLc(orgId, userId, issueIds)
	//if issueInfoErr != nil {
	//	log.Errorf("[createCalendarEventsBefore] failed:%v", issueInfoErr)
	//	return
	//}

	//获取已经创建过日程的任务
	relationInfo, err1 := GetRelationInfoByIssueIds(issueIds, []int{consts.IssueRelationTypeCalendar})
	if err1 != nil {
		log.Error(err1)
		return
	}

	alreadyCreateIds := []int64{}
	for _, relationBo := range relationInfo {
		alreadyCreateIds = append(alreadyCreateIds, relationBo.IssueId)
	}

	needCreateIds := []int64{}
	for _, id := range issueIds {
		if ok, _ := slice.Contain(alreadyCreateIds, id); !ok {
			needCreateIds = append(needCreateIds, id)
		}
	}
	// 如果需要处理已存在的日程，则更新对应的日程订阅者。
	if dealAlreadyExist {
		if len(needCreateIds) == 0 && len(alreadyCreateIds) == 0 {
			return
		}
	} else {
		if len(needCreateIds) == 0 {
			return
		}
	}

	for _, issueBo := range issueBos {
		if ok, _ := slice.Contain(needCreateIds, issueBo.Id); ok {
			CreateCalendarEvent(issueBo, userId, issueBo.FollowerIdsI64)
		} else {
			// 如果需要对已存在的日程更新其订阅者，则进行处理。
			if dealAlreadyExist {
				if ok1, _ := slice.Contain(alreadyCreateIds, issueBo.Id); ok1 {
					CreateCalendarEvent(issueBo, userId, issueBo.FollowerIdsI64)
				}
			}
		}
	}
}

// 更新日历
func UpdateCalendar(input vo.UpdateProjectReq, orgId int64, userId int64, oldCalendarInfo *bo.CacheProjectCalendarInfoBo, oldProjectInfo bo.ProjectBo) {
	projectCalendarInfo, err := GetProjectCalendarInfo(orgId, input.ID)
	if err != nil {
		log.Error(err)
		return
	}
	var deleted, added []int64
	var oldMembers, newMembers []int64
	newSyncCalendarFlag := TransferSyncOutCalendarStatusIntoOne(input.SyncCalendarStatusList)
	if projectCalendarInfo.CalendarId != consts.BlankString {
		//创建成员订阅日历
		//deleted, added := util.GetDifMemberIds(oldMembers, newMembers)
		//if !util.FieldInUpdate(input.UpdateFields, "memberIds") {
		// 如果没有更新 memberIds，则需主动查询 member
		participantUids, err := GetProjectParticipantIds(orgId, input.ID)
		if err != nil {
			log.Errorf("UpdateCalendar 查询项目成员异常 err: %v", err)
			return
		}
		oldMembers = participantUids
		newMembers = participantUids
		//}
		// 计算之前是否勾选“订阅日历”。
		oldHasSubCalender := CheckCalendarIsSyncToSubCalendar(oldCalendarInfo.IsSyncOutCalendar)
		newHasSubCalender := CheckCalendarIsSyncToSubCalendar(newSyncCalendarFlag)
		// 更新日历时，如果状态是不同步给订阅日历。则 newMembers 置为空。
		// 如果“订阅日历”从"勾选"切换到“未勾选”，则给项目负责人、关注人去除访问控制。
		if oldHasSubCalender != newHasSubCalender && !CheckCalendarIsSyncToSubCalendar(newSyncCalendarFlag) {
			deleted = oldMembers
		}
		// 如果“订阅日历”从“未勾选”切换到“勾选”，则给项目负责人、关注人增加访问控制。
		if oldHasSubCalender != newHasSubCalender && CheckCalendarIsSyncToSubCalendar(newSyncCalendarFlag) {
			added = newMembers
			if len(added) < 1 {
				added = oldMembers
			}
		}
		log.Infof("UpdateCalendar debug. debug info: %s. ", json.ToJsonIgnoreError(map[string]interface{}{
			"added":           added,
			"deleted":         deleted,
			"oldMembers":      oldMembers,
			"newMembers":      newMembers,
			"oldCalendarInfo": oldCalendarInfo,
		}))
		UpdateCalendarAttendees(orgId, projectCalendarInfo.CalendarId, projectCalendarInfo.Creator, added, deleted, input.ID)
	}
	// 大于 0 的值才是有效的值。
	if projectCalendarInfo.IsSyncOutCalendar < 1 {
		log.Error("未设置导出日历")
		return
	}
	if projectCalendarInfo.CalendarId == consts.BlankString {
		log.Info("无对应项目日历,重新生成日历")
		CreateCalendar(&projectCalendarInfo.IsSyncOutCalendar, orgId, input.ID, userId, newMembers)
		return
	}

	orgBaseInfo, err := orgfacade.GetBaseOrgInfoRelaxed(orgId)
	if err != nil {
		log.Errorf("组织外部信息不存在 %v", err)
		return
	}

	updateReq := &vo3.UpdateCalendarReq{
		CalendarId:  projectCalendarInfo.CalendarId,
		Name:        oldProjectInfo.Name,
		Description: oldProjectInfo.Remark,
	}
	needUpdate := 0
	for _, v := range input.UpdateFields {
		if v == "name" && input.Name != nil {
			updateReq.Name = *input.Name
			needUpdate = 1
		}
		if v == "remark" && input.Remark != nil {
			updateReq.Description = *input.Remark
			needUpdate = 1
		}
	}

	if needUpdate != 0 {
		client, platformUserId, err := getPlatformClientAndCreatorId(orgId, orgBaseInfo.OutOrgId, orgBaseInfo.SourceChannel, userId)
		if err != nil {
			log.Errorf("[UpdateCalendar] GetClient, sourceChannel:%v, outOrgId:%v, err:%v", orgBaseInfo.SourceChannel, orgBaseInfo.OutOrgId, err)
			return
		}

		updateReq.UserId = platformUserId
		_, err2 := client.UpdateCalendar(updateReq)
		if err2 != nil {
			log.Errorf("[CreateCalendar] UpdateCalendar, updateReq:%v, err:%v", updateReq, err2)
			return
		}
	}

	// 检查是否更改了 syncCalendarStatusList 中的 关注人、负责人配置，如更改了，则需要对已存在的日程更新订阅者
	if CheckIsUpdateCalendarOwnerOrFollower(oldCalendarInfo.IsSyncOutCalendar, newSyncCalendarFlag) {
		createCalendarEventsBefore(orgId, input.ID, userId, true)
	} else {
		createCalendarEventsBefore(orgId, input.ID, userId, false)
	}
}

// 对比新旧的值，看是否更改了 syncCalendarStatusList 中的 关注人、负责人配置
func CheckIsUpdateCalendarOwnerOrFollower(old, new int) bool {
	result := false
	oldVal := CheckCalendarIsSyncToOwner(old)
	newVal := CheckCalendarIsSyncToOwner(new)
	if oldVal != newVal {
		result = true
		return result
	}
	oldVal = CheckCalendarIsSyncToFollower(old)
	newVal = CheckCalendarIsSyncToFollower(new)
	if oldVal != newVal {
		result = true
		return result
	}
	return result
}

// 检查同步到日历配置是否需要同步给负责人
func CheckCalendarIsSyncToOwner(syncFlag int) bool {
	isOk := false
	if syncFlag == consts.IsSyncOutCalendar {
		isOk = true
		return isOk
	}
	if syncFlag&consts.IsSyncOutCalendarForOwner == consts.IsSyncOutCalendarForOwner {
		isOk = true
		return isOk
	}
	if syncFlag == consts.IsSyncOutCalendarForOwnerAndFollower {
		isOk = true
		return isOk
	}
	return isOk
}

// 检查同步到日历配置是否需要同步给关注人
func CheckCalendarIsSyncToFollower(syncFlag int) bool {
	isOk := false
	if syncFlag == consts.IsSyncOutCalendar {
		isOk = true
		return isOk
	}
	if syncFlag&consts.IsSyncOutCalendarForFollower == consts.IsSyncOutCalendarForFollower {
		isOk = true
		return isOk
	}
	if syncFlag == consts.IsSyncOutCalendarForOwnerAndFollower {
		isOk = true
		return isOk
	}
	return isOk
}

// 检查是否同步到"订阅日历"中
func CheckCalendarIsSyncToSubCalendar(syncFlag int) bool {
	isOk := false
	if syncFlag == consts.IsSyncOutCalendar {
		isOk = true
		return isOk
	}
	if syncFlag&consts.IsSyncOutCalendarForSubCalendar == consts.IsSyncOutCalendarForSubCalendar {
		isOk = true
		return isOk
	}
	return isOk
}

// 创建日程
func CreateCalendarEvent(issueBo *bo.IssueBo, userId int64, followers []int64) {
	//防止重复插入
	uid := uuid.NewUuid()
	issueIdStr := strconv.FormatInt(issueBo.Id, 10)
	lockKey := consts.CreateCalendarEventLock + issueIdStr
	suc, err := cache.TryGetDistributedLock(lockKey, uid)
	if err != nil {
		log.Errorf("获取%s锁时异常 %v", lockKey, err)
		return
	}
	if suc {
		defer func() {
			if _, err := cache.ReleaseDistributedLock(lockKey, uid); err != nil {
				log.Error(err)
			}
		}()
	}

	ok, calendarInfo := GetCalendarInfo(issueBo.OrgId, issueBo.ProjectId)
	if !ok {
		return
	}

	//如果原来没有则新建
	calendarEventRelation, err := dao.SelectOneIssueRelation(db.Cond{
		consts.TcIsDelete:     consts.AppIsNoDelete,
		consts.TcRelationType: consts.IssueRelationTypeCalendar,
		consts.TcIssueId:      issueBo.Id,
		consts.TcOrgId:        issueBo.OrgId,
	})
	if err != nil {
		if err != db.ErrNoMoreRows {
			log.Error(err)
			return
		}
	} else {
		log.Info("日程已创建，只对订阅人进行更新。issueIdStr: " + issueIdStr)
		// 对于已经创建的日程，更新一下对应的订阅人群

		attendeesId := []int64{}
		addedParticipantIds := []int64{}
		deletedParticipantIds := []int64{}
		if CheckCalendarIsSyncToOwner(calendarInfo.SyncCalendarFlag) {
			addedParticipantIds = append(addedParticipantIds, businees.LcMemberToUserIds(issueBo.OwnerId)...)
		} else {
			deletedParticipantIds = append(deletedParticipantIds, businees.LcMemberToUserIds(issueBo.OwnerId)...)
		}
		if CheckCalendarIsSyncToFollower(calendarInfo.SyncCalendarFlag) {
			addedParticipantIds = append(addedParticipantIds, followers...)
		} else {
			deletedParticipantIds = append(deletedParticipantIds, followers...)
		}
		// 聚合，查询对应的用户信息。
		attendeesId = append(addedParticipantIds, deletedParticipantIds...)
		attendeesId = slice.SliceUniqueInt64(attendeesId)
		if len(attendeesId) == 0 {
			log.Info("无需更新日程的订阅人。")
			return
		}
		infosMap, err := getBaseUserInfoMap(calendarInfo.OrgId, attendeesId)
		if err != nil {
			log.Error(err)
			return
		}

		//更新日程订阅者
		updateCalendarEventAttendees(calendarInfo, calendarEventRelation.RelationCode, calendarEventRelation.Creator,
			getPlatformAttendUserWithInfoMap(calendarInfo.SourceChannel, addedParticipantIds, infosMap),
			getPlatformAttendUserWithInfoMap(calendarInfo.SourceChannel, deletedParticipantIds, infosMap),
		)

		return
	}

	if time.Time(issueBo.PlanStartTime).Unix() <= 0 || time.Time(issueBo.PlanEndTime).Unix() <= 0 {
		return
	}

	attendeesId := []int64{}
	if CheckCalendarIsSyncToOwner(calendarInfo.SyncCalendarFlag) {
		attendeesId = append(attendeesId, businees.LcMemberToUserIds(issueBo.OwnerId)...)
	}
	if CheckCalendarIsSyncToFollower(calendarInfo.SyncCalendarFlag) {
		attendeesId = append(attendeesId, followers...)
	}
	attendeesId = slice.SliceUniqueInt64(attendeesId)
	var attendees []*vo3.AttendUser
	if len(attendeesId) == 0 {
		log.Info("日程的订阅人为空，但是还需创建日程，因为日历的订阅人需要看到。")
		// return
	} else {
		attendees, err = getPlatformAttendUser(calendarInfo.OrgId, calendarInfo.SourceChannel, attendeesId)
		if err != nil {
			log.Error(err)
			return
		}
	}

	//创建日程
	log.Infof("[CreateCalendarEvent] getPlatformClientAndCreatorId orgId:%v, outOrgId:%v, userId:%v, start:%v, end:%v",
		calendarInfo.OrgId, calendarInfo.OutOrgId, userId, time.Time(issueBo.PlanStartTime).Unix(), time.Time(issueBo.PlanEndTime).Unix())

	client, platformUserId, err := getPlatformClientAndCreatorId(calendarInfo.OrgId, calendarInfo.OutOrgId, calendarInfo.SourceChannel, userId)
	if err != nil {
		log.Errorf("[CreateCalendarEvent] GetClient, sourceChannel:%v, outOrgId:%v, err:%v", calendarInfo.SourceChannel, calendarInfo.OutOrgId, err)
		return
	}

	log.Infof("[CreateCalendarEvent] sourceChannel:%v, userId:%v, platformUserId:%v, attendees:%v", calendarInfo.SourceChannel, userId, platformUserId, json.ToJsonIgnoreError(attendees))

	reply, err := client.CreateCalendarEvent(&vo3.CreateCalendarEventReq{
		UserId:      platformUserId,
		CalendarId:  calendarInfo.CalendarId,
		Summary:     issueBo.Title,
		Description: getCalendarEventDescription(issueBo.OrgId, issueBo.Id, issueBo.ParentId),
		Start:       time.Time(issueBo.PlanStartTime),
		End:         time.Time(issueBo.PlanEndTime),
		AttendUsers: attendees,
	})
	if err != nil {
		log.Errorf("[CreateCalendarEvent] err:%v", err)
		return
	}

	//创建关联关系
	relationId, err := idfacade.ApplyPrimaryIdRelaxed(consts.TableIssueRelation)
	if err != nil {
		log.Error(err)
		return
	}
	err1 := dao.InsertIssueRelation(po.PpmPriIssueRelation{
		Id:           relationId,
		OrgId:        issueBo.OrgId,
		IssueId:      issueBo.Id,
		RelationType: consts.IssueRelationTypeCalendar,
		RelationCode: reply.EventId,
		Creator:      userId,
		CreateTime:   time.Now(),
		Updator:      userId,
		UpdateTime:   time.Now(),
		IsDelete:     consts.AppIsNoDelete,
	})
	if err1 != nil {
		log.Error(err1)
		return
	}
}

func updateCalendarEventAttendees(calendarInfo bo.CalendarInfo, eventId string, eventCreator int64, addUsers, deleteUsers []*vo3.AttendUser) {
	client, platformUserId, err := getPlatformClientAndCreatorId(calendarInfo.OrgId, calendarInfo.OutOrgId, calendarInfo.SourceChannel, eventCreator)
	if err != nil {
		log.Errorf("[updateCalendarEventAttendees] GetClient, sourceChannel:%v, outOrgId:%v, err:%v", calendarInfo.SourceChannel, calendarInfo.OutOrgId, err)
		return
	}

	req := &vo3.UpdateCalendarEventAttendReq{
		UserId:      platformUserId,
		CalendarId:  calendarInfo.CalendarId,
		EventId:     eventId,
		AddUsers:    addUsers,
		DeleteUsers: deleteUsers,
	}
	_, err2 := client.UpdateCalendarEventAttend(req)
	if err2 != nil {
		log.Errorf("[updateCalendarEventAttendees] UpdateCalendarEventAttend 更新日程订阅者失败, req:%v, err:%v", req, err2)
	}
}

func getCalendarEventDescription(orgId, issueId, parentId int64) string {
	description := ""
	if parentId != 0 {
		//parentInfo, err := GetIssueBo(orgId, parentId)
		//if err != nil {
		//	log.Error(err)
		//} else {
		//	description += "父任务：" + parentInfo.Title + "\n"
		//}
		parentInfo, err := GetIssueInfoLc(orgId, 0, parentId)
		if err != nil {
			log.Errorf("[getCalendarEventDescription] GetIssueInfoLc err:%v, orgId:%v, parentId:%v", err, orgId, parentId)
			return ""
		}
		description += "父任务：" + parentInfo.Title + "\n"
	}

	issueLink := GetIssueLinks("", orgId, issueId).Link
	description += "任务链接：" + issueLink + "\n"
	//获取任务详情
	issueInfo, err := GetIssueInfoLc(orgId, 0, issueId)
	if err != nil {
		log.Errorf("[getCalendarEventDescription] GetIssueInfoLc err:%v, orgId:%v, issueId:%v", err, orgId, issueId)
		return ""
	}
	//issueDetailPo := &po.PpmPriIssueDetail{}
	//detailErr := mysql.SelectOneByCond(consts.TableIssueDetail, db.Cond{
	//	consts.TcIsDelete: consts.AppIsNoDelete,
	//	consts.TcIssueId:  issueId,
	//}, issueDetailPo)
	//if detailErr != nil {
	//	if detailErr != db.ErrNoMoreRows {
	//		log.Error(detailErr)
	//		return description
	//	}
	//} else {
	//	if issueDetailPo.Remark != nil && *issueDetailPo.Remark != consts.BlankString {
	//		description += "任务详情：" + *issueDetailPo.Remark
	//	}
	//}
	if issueInfo.Remark != consts.BlankString {
		description += "任务详情：" + issueInfo.Remark
	}
	return description
}

// UpdateCalendarEvent 更新日程
// 标题、起止时间、负责人、关注人的变更会触发日程的变更
func UpdateCalendarEvent(orgId int64, operatorId int64, issueId int64, oldIssueBo *bo.IssueBo, newIssueBo *bo.IssueBo, beforeFollowers, afterFollowers []int64) {
	ok, calendarInfo := GetCalendarInfo(orgId, oldIssueBo.ProjectId)
	if !ok {
		log.Info("[UpdateCalendarEvent] 更新日历日程时，查询不到日历信息。")
		return
	}

	//如果原来没有则新建
	relation, err := dao.SelectOneIssueRelation(db.Cond{
		consts.TcIsDelete:     consts.AppIsNoDelete,
		consts.TcRelationType: consts.IssueRelationTypeCalendar,
		consts.TcIssueId:      issueId,
		consts.TcOrgId:        orgId,
	})
	if err != nil {
		log.Infof("[UpdateCalendarEvent] 更新时生成日程 issueId: %d", issueId)
		CreateCalendarEvent(newIssueBo, operatorId, afterFollowers)
		return
	}
	//如果有则更新（基本信息，订阅者）
	//更新日程
	if time.Time(newIssueBo.PlanStartTime).Unix() <= 0 || time.Time(newIssueBo.PlanEndTime).Unix() <= 0 {
		log.Infof("[UpdateCalendarEvent] 起止时间为空时间，则删除对应的日程 issueId: %d, planStartTime: %s, planEndTime: %s", issueId,
			newIssueBo.PlanStartTime, newIssueBo.PlanEndTime)
		if err := DeleteCalendarEventBatch(orgId, newIssueBo.ProjectId, []int64{issueId}); err != nil {
			log.Errorf("[UpdateCalendarEvent] 删除日程异常：%v, issueId: %d", err, issueId)
			return
		}
		return
	}

	isSyncToOwner := CheckCalendarIsSyncToOwner(calendarInfo.SyncCalendarFlag)
	isSyncToFollower := CheckCalendarIsSyncToFollower(calendarInfo.SyncCalendarFlag)

	// 人员变动影响订阅者
	deletedParticipantIds := make([]int64, 0)
	addedParticipantIds := make([]int64, 0)
	if isSyncToFollower {
		deletedParticipantIds, addedParticipantIds = util.GetDifMemberIds(beforeFollowers, afterFollowers)
	}

	oldOwnerIds := businees.LcMemberToUserIds(oldIssueBo.OwnerId)
	newOwnerIds := businees.LcMemberToUserIds(newIssueBo.OwnerId)

	if isSyncToOwner {
		del, _ := util.GetDifMemberIds(oldOwnerIds, newOwnerIds)
		deletedParticipantIds = append(deletedParticipantIds, del...)
		addedParticipantIds = append(addedParticipantIds, newOwnerIds...)
	} else {
		// 如果不同步负责人，则将其视为要被删除的 attendance
		if len(oldOwnerIds) > 0 {
			deletedParticipantIds = append(deletedParticipantIds, oldOwnerIds...)
		}
	}

	if isSyncToFollower {
		addedParticipantIds = append(addedParticipantIds, afterFollowers...)
	} else {
		deletedParticipantIds = append(deletedParticipantIds, afterFollowers...)
	}
	deletedParticipantIds = slice.SliceUniqueInt64(deletedParticipantIds)
	addedParticipantIds = slice.SliceUniqueInt64(addedParticipantIds)
	// 针对在 deletedParticipantIds 中，又在 addedParticipantIds 中的元素，我们视为，对应的用户还需订阅。因此需要将其从 deletedParticipantIds 中去除。
	// 如：删除关注人时，恰好关注人和负责人是同一人，此时如果还需同步给负责人。则需要将关注人从 deletedParticipantIds 中删除。
	intersectUid := int64Util.Int64Intersect(deletedParticipantIds, addedParticipantIds)
	deletedParticipantIds = int64Util.Int64RemoveSomeVal(deletedParticipantIds, intersectUid)
	attendeesId := slice.SliceUniqueInt64(append(deletedParticipantIds, addedParticipantIds...))
	// 推送目标人员的增减标识。
	editParticipantFlag := true
	if len(attendeesId) == 0 {
		editParticipantFlag = false
		log.Info("[UpdateCalendarEvent] 无相关人员")
		// return
	} else {
		// 如果无需增减人员，还需向已经关联日历的人员进行推送“更新的日程”
		attendeesId = append(attendeesId, afterFollowers...)
	}
	log.Infof("[UpdateCalendarEvent] debug info: %s. ", json.ToJsonIgnoreError(map[string]interface{}{
		"beforeFollowers": beforeFollowers,
		"afterFollowers":  afterFollowers,
		"attendeesId":     attendeesId,
		"NewIssueBo":      newIssueBo,
	}))
	infosMap, err := getBaseUserInfoMap(orgId, attendeesId)
	if err != nil {
		log.Errorf("[UpdateCalendarEvent] orgId:%v, userIds:%v err:%v", orgId, attendeesId, err)
		return
	}

	// 更新日历 start。
	if editParticipantFlag {
		log.Info("[UpdateCalendarEvent] step4-1")
		//更新日程参与者
		updateCalendarEventAttendees(calendarInfo, relation.RelationCode, relation.Creator,
			getPlatformAttendUserWithInfoMap(calendarInfo.SourceChannel, addedParticipantIds, infosMap),
			getPlatformAttendUserWithInfoMap(calendarInfo.SourceChannel, deletedParticipantIds, infosMap))
	}

	client, platformUserId, err := getPlatformClientAndCreatorId(calendarInfo.OrgId, calendarInfo.OutOrgId, calendarInfo.SourceChannel, relation.Creator)
	if err != nil {
		log.Errorf("[UpdateCalendarEvent] GetClient, sourceChannel:%v, outOrgId:%v, err:%v", calendarInfo.SourceChannel, calendarInfo.OutOrgId, err)
		return
	}

	req := &vo3.UpdateCalendarEventReq{
		UserId:      platformUserId,
		CalendarId:  calendarInfo.CalendarId,
		EventId:     relation.RelationCode,
		Summary:     newIssueBo.Title,
		Description: getCalendarEventDescription(orgId, issueId, oldIssueBo.ParentId),
		Start:       time.Time(newIssueBo.PlanStartTime),
		End:         time.Time(newIssueBo.PlanEndTime),
	}
	_, err = client.UpdateCalendarEvent(req)
	if err != nil {
		log.Errorf("[UpdateCalendarEvent] UpdateCalendarEvent, req:%v, err:%v", req, err)
	}
}

// 中途同步日历
func SyncCalendarConfirm(orgId, userId, projectId int64) {
	projectCalendarInfo, err := GetProjectCalendarInfo(orgId, projectId)
	if err != nil {
		log.Error(err)
		return
	}
	if projectCalendarInfo.CalendarId != "" {
		log.Info("项目已同步日历")
		createCalendarEventsBefore(orgId, projectId, userId, false)
		return
	}
	info, err := GetProjectRelationByType(projectId, []int64{consts.ProjectRelationTypeOwner, consts.ProjectRelationTypeParticipant})
	if err != nil {
		log.Error(err)
		return
	}
	addIds := []int64{}
	for _, v := range *info {
		addIds = append(addIds, v.RelationId)
	}
	CreateCalendar(&projectCalendarInfo.IsSyncOutCalendar, orgId, projectId, userId, addIds)
}

func SwitchCalendar(orgId, oldProjectId int64, issueIds []int64, operatorId int64, newProjectId int64) errs.SystemErrorInfo {
	errDelete := DeleteCalendarEventBatch(orgId, oldProjectId, issueIds)
	if errDelete != nil {
		log.Errorf("[SwitchCalendar] DeleteCalendarEventBatch orgId:%v, projectId:%v, issueIds:%v", orgId, oldProjectId, issueIds)
		return errs.BuildSystemErrorInfo(errs.SystemBusy, errDelete)
	}

	issuesInfo, issueErr := GetIssueInfosLc(orgId, operatorId, issueIds)
	if issueErr != nil {
		log.Error(issueErr)
		return issueErr
	}

	for _, issueBo := range issuesInfo {
		CreateCalendarEvent(issueBo, operatorId, issueBo.FollowerIdsI64)
	}

	return nil
}

// 删除一个日历，以及日历对应的日程
func DeleteOneCalendar(orgId, projectId int64) error {
	ok, calendarInfo := GetCalendarInfo(orgId, projectId)
	if !ok {
		return nil
	}

	client, platformUserId, err := getPlatformClientAndCreatorId(calendarInfo.OrgId, calendarInfo.OutOrgId, calendarInfo.SourceChannel, calendarInfo.Creator)
	if err != nil {
		log.Errorf("[UpdateCalendarEvent] GetClient, sourceChannel:%v, outOrgId:%v, err:%v", calendarInfo.SourceChannel, calendarInfo.OutOrgId, err)
		return err
	}
	_, err2 := client.DeleteCalendar(&vo3.DeleteCalendarReq{
		UserId:     platformUserId,
		CalendarId: calendarInfo.CalendarId,
	})
	if err2 != nil {
		log.Errorf("[UpdateCalendarEvent] DeleteCalendar calendarId:%v, err:%v", calendarInfo.CalendarId, err)
		return err
	}

	return DeleteCalendarEventByProjectId(orgId, projectId)
}

// 删除一个项目下所有的日程
func DeleteCalendarEventByProjectId(orgId, projectId int64) error {
	// 因“任务多项目”需求之故，对查询条件做一下转变，查询项目下的所有任务，需要通过 issue_relation 表，查询 relation_type 为 `consts.IssueRelationTypeBelongManyPro` 的关联数据。
	allowIssueIds, err := GetIssueIdsByProIds(orgId, []int64{projectId})
	if err != nil {
		log.Error(err)
		return err
	}
	// 如果没有关联的任务id，但又有项目的限定，此时用 0 表示查到的任务 id。
	if allowIssueIds == nil || len(allowIssueIds) < 1 {
		allowIssueIds = []int64{0}
	}

	log.Infof("DeleteCalendarEventByProjectId 删除日程 ok")
	return DeleteCalendarEventBatch(orgId, projectId, allowIssueIds)
}

// 批量删除日程。
// 遍历，1.调用删除日程接口，2删除日程与任务的关联
// 因为删除任务的操作已执行，因此需要将删除状态的 Calendar issue relation 查询出来，删除对应的日程。
func DeleteCalendarEventBatch(orgId, projectId int64, issueIds []int64) error {
	if len(issueIds) < 1 {
		return nil
	}
	ok, calendarInfo := GetCalendarInfo(orgId, projectId)
	if !ok {
		return nil
	}

	client, platformUserId, err := getPlatformClientAndCreatorId(calendarInfo.OrgId, calendarInfo.OutOrgId, calendarInfo.SourceChannel, calendarInfo.Creator)
	if err != nil {
		log.Errorf("[DeleteCalendarEventBatch] GetClient, sourceChannel:%v, outOrgId:%v, err:%v", calendarInfo.SourceChannel, calendarInfo.OutOrgId, err)
		return err
	}

	// 查询任务何日程的关联
	// 这里是查询删除状态下的 relation。
	relationInfos, err1 := GetRelationInfoByIssueIds(issueIds, []int{consts.IssueRelationTypeCalendar})
	if err1 != nil {
		log.Error(err1)
		return err1
	}
	alreadyExistIssueIdMap := map[int64]string{}
	for _, relationBo := range relationInfos {
		alreadyExistIssueIdMap[relationBo.IssueId] = relationBo.RelationCode
	}
	for _, eventCode := range alreadyExistIssueIdMap {
		_, err := client.DeleteCalendarEvent(&vo3.DeleteCalendarEventReq{
			UserId:     platformUserId,
			CalendarId: calendarInfo.CalendarId,
			EventId:    eventCode,
		})
		if err != nil {
			log.Errorf("[DeleteCalendarEventBatch] calendarId:%v, eventId:%v, err:%v", calendarInfo.CalendarId, eventCode, err)
			return err
		}
	}

	// 删除日程与日历的关联
	condition := db.Cond{
		consts.TcIsDelete:     consts.AppIsNoDelete,
		consts.TcIssueId:      db.In(issueIds),
		consts.TcRelationType: consts.IssueRelationTypeCalendar,
	}
	updateData := mysql.Upd{
		consts.TcIsDelete: consts.AppIsDeleted,
	}
	_, busiErr := mysql.UpdateSmartWithCond(consts.TableIssueRelation, condition, updateData)
	if busiErr != nil {
		return errors.New("DeleteCalendarEventBatch 删除日程与日历的关联异常：" + busiErr.Error())
	}
	log.Infof("批量删除日程 ok")
	return nil
}

// 恢复若干个日程。任务从回收站恢复时触发。
// 恢复子任务/父任务，1创建日程，2创建日程与任务的关联。
func RecoveryCalendarEventBatch(orgId, projectId int64, issueIds []int64, opUserId int64) error {
	// 查询所有任务
	issueBos, busiErr := GetIssueInfosLc(orgId, opUserId, issueIds)
	if busiErr != nil {
		log.Error(busiErr)
		return busiErr
	}
	for _, issueBo := range issueBos {
		CreateCalendarEvent(issueBo, opUserId, issueBo.FollowerIdsI64)
	}
	return nil
}

// 获取应用（极星）的授权状态，检查是否有需要的权限。
func GetAppScopes(outOrgId string) (*vo2.GetScopesResp, error) {
	tenant, err := feishu.GetTenant(outOrgId)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	return tenant.GetScopes()
}

// 检查极星应用是否有访问日历的权限。
func CheckHasCalendarPower(outOrgId string) (bool, error) {
	var powerFlag = "calendar:calendar:access"
	resp, err := GetAppScopes(outOrgId)
	if err != nil {
		return false, err
	}
	for _, item := range resp.Data.Scopes {
		if item.ScopeName == powerFlag && item.GrantStatus == 1 {
			return true, nil
		}
	}
	return false, nil
}

func JudgeCalendarSdkVersion(outOrgId string) (string, errs.SystemErrorInfo) {
	resp, err := GetAppScopes(outOrgId)
	if err != nil {
		return "", errs.FeiShuOpenApiCallError
	}

	for _, item := range resp.Data.Scopes {
		if item.ScopeName == consts.ScopeNameCalendarAccess && item.GrantStatus == 1 {
			return VersionV3, nil
		} else if item.ScopeName == consts.ScopeNameCalendarCalendar && item.GrantStatus == 1 {
			return VersionV4, nil
		}
	}

	//暂时没有v4版本权限就日历先不可用
	return "", errs.FeiShuNotScopeInCalendar
	//如果没有明确展示的话默认用v3吧
	//return VersionV3, nil
}

func getBaseUserInfoMap(orgId int64, userIds []int64) (map[int64]*bo.BaseUserInfoBo, errs.SystemErrorInfo) {
	m := make(map[int64]*bo.BaseUserInfoBo)
	resp, err := orgfacade.GetBaseUserInfoBatchRelaxed(orgId, userIds)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	for i := range resp {
		m[resp[i].UserId] = &resp[i]
	}

	return m, nil
}

func getPlatformUserIdsWithInfoMap(sourceChannel string, userIds []int64, baseInfoMap map[int64]*bo.BaseUserInfoBo) []string {
	otherUserIds := make([]string, 0, len(userIds))
	for _, id := range userIds {
		if info, ok := baseInfoMap[id]; ok {
			// 钉钉日历使用unionId
			if sourceChannel == sdk_const.SourceChannelDingTalk {
				otherUserIds = append(otherUserIds, info.OutOrgUserId)
			} else {
				otherUserIds = append(otherUserIds, info.OutUserId)
			}
		}
	}

	return otherUserIds
}

func getPlatformUserIds(orgId int64, sourceChannel string, userIds []int64) ([]string, errs.SystemErrorInfo) {
	baseInfoMap, err := getBaseUserInfoMap(orgId, userIds)
	if err != nil {
		log.Errorf("[getPlatformUserIds] getBaseUserInfoMap orgId:%v, userIds:%v, err:%v", orgId, userIds, err)
		return nil, err
	}
	otherUserIds := make([]string, 0, len(userIds))
	for _, id := range userIds {
		if info, ok := baseInfoMap[id]; ok {
			// 钉钉日历使用unionId
			if sourceChannel == sdk_const.SourceChannelDingTalk {
				otherUserIds = append(otherUserIds, info.OutOrgUserId)
			} else {
				otherUserIds = append(otherUserIds, info.OutUserId)
			}
		}
	}

	return otherUserIds, nil
}

func getPlatformAttendUser(orgId int64, sourceChannel string, userIds []int64) ([]*vo3.AttendUser, errs.SystemErrorInfo) {
	baseInfoMap, err := getBaseUserInfoMap(orgId, userIds)
	if err != nil {
		log.Errorf("[getPlatformAttendUser] getBaseUserInfoMap orgId:%v, userIds:%v, err:%v", orgId, userIds, err)
		return nil, err
	}

	return getPlatformAttendUserWithInfoMap(sourceChannel, userIds, baseInfoMap), nil
}

func getPlatformAttendUserWithInfoMap(sourceChannel string, userIds []int64, baseInfoMap map[int64]*bo.BaseUserInfoBo) []*vo3.AttendUser {
	attendUsers := make([]*vo3.AttendUser, 0, len(userIds))
	for _, id := range userIds {
		if info, ok := baseInfoMap[id]; ok {
			// 钉钉日历使用unionId
			if sourceChannel == sdk_const.SourceChannelDingTalk {
				attendUsers = append(attendUsers, &vo3.AttendUser{
					UserId: info.OutOrgUserId,
					Name:   info.Name,
				})
			} else {
				attendUsers = append(attendUsers, &vo3.AttendUser{
					UserId: info.OutUserId,
					Name:   info.Name,
				})
			}
		}
	}

	return attendUsers
}

// getPlatformClientAndCreatorIdWithInfoMap 钉钉日程相关的都要创建人的id，巨坑
func getPlatformClientAndCreatorIdWithInfoMap(orgId int64, outOrgId string, sourceChannel string, creator int64, baseInfoMap map[int64]*bo.BaseUserInfoBo) (sdk_interface.Sdk, string, errs.SystemErrorInfo) {
	client, err2 := platform_sdk.GetClient(sourceChannel, outOrgId)
	if err2 != nil {
		log.Errorf("[UpdateCalendarAttendees] GetClient, sourceChannel:%v, outOrgId:%v, err:%v", sourceChannel, outOrgId, err2)
		return nil, "", errs.BuildSystemErrorInfo(errs.PlatFormOpenApiCallError, err2)
	}

	platformUserId := ""

	var (
		platformUserIds []string
		err             errs.SystemErrorInfo
	)
	if baseInfoMap == nil {
		platformUserIds, err = getPlatformUserIds(orgId, sourceChannel, []int64{creator})
		if err != nil {
			return nil, "", err
		}
	} else {
		platformUserIds = getPlatformUserIdsWithInfoMap(sourceChannel, []int64{creator}, baseInfoMap)
	}
	if len(platformUserIds) > 0 {
		platformUserId = platformUserIds[0]
	}

	return client, platformUserId, nil
}

func getPlatformClientAndCreatorId(orgId int64, outOrgId string, sourceChannel string, creator int64) (sdk_interface.Sdk, string, errs.SystemErrorInfo) {
	baseInfoMap, err := getBaseUserInfoMap(orgId, []int64{creator})
	if err != nil {
		log.Errorf("[getPlatformUserIds] getBaseUserInfoMap orgId:%v, userIds:%v, err:%v", orgId, creator, err)
		return nil, "", err
	}

	return getPlatformClientAndCreatorIdWithInfoMap(orgId, outOrgId, sourceChannel, creator, baseInfoMap)
}
