package domain

import (
	"time"

	"github.com/star-table/startable-server/common/core/consts"
	"github.com/star-table/startable-server/common/model/vo"

	"github.com/star-table/startable-server/app/facade/idfacade"
	"github.com/star-table/startable-server/app/facade/orgfacade"
	"github.com/star-table/startable-server/app/service/projectsvc/po"
	"github.com/star-table/startable-server/common/core/errs"
	"github.com/star-table/startable-server/common/core/util/copyer"
	"github.com/star-table/startable-server/common/core/util/json"
	"github.com/star-table/startable-server/common/core/util/slice"
	"github.com/star-table/startable-server/common/library/db/mysql"
	"github.com/star-table/startable-server/common/model/bo"
	"github.com/star-table/startable-server/common/model/vo/orgvo"
	"upper.io/db.v3"
	"upper.io/db.v3/lib/sqlbuilder"
)

func HandleProjectMember(orgId int64, currentUserId int64, owner int64, projectId int64, memberIds []int64, followerIds []int64, isAllMember *bool, departmentIds []int64, ownerIds []int64) ([]interface{}, []int64, errs.SystemErrorInfo) {
	//插入项目成员
	memberEntities := []interface{}{}
	addedMemberIds := []int64{}

	//1.负责人
	allOwner := []int64{}
	if owner != int64(0) {
		allOwner = append(allOwner, owner)
	}
	if ownerIds != nil && len(ownerIds) > 0 {
		allOwner = append(allOwner, ownerIds...)
	}
	allOwner = slice.SliceUniqueInt64(allOwner)
	if len(allOwner) != 0 {
		verifyOrgUserFlag := orgfacade.VerifyOrgUsersRelaxed(orgId, allOwner)
		if !verifyOrgUserFlag {
			log.Error("存在用户组织校验失败")
			return memberEntities, nil, errs.VerifyOrgError
		}
		memberRelationIds, idErr := idfacade.ApplyMultiplePrimaryIdRelaxed(consts.TableProjectRelation, len(allOwner))
		if idErr != nil {
			log.Error(idErr)
			return memberEntities, nil, idErr
		}
		for i, v := range allOwner {
			memberEntities = append(memberEntities, po.PpmProProjectRelation{
				Id:           memberRelationIds.Ids[i].Id,
				OrgId:        orgId,
				ProjectId:    projectId,
				RelationId:   v,
				RelationType: consts.ProjectRelationTypeOwner,
				Creator:      currentUserId,
				CreateTime:   time.Now(),
				IsDelete:     consts.AppIsNoDelete,
				Status:       consts.ProjectMemberEffective,
				Updator:      currentUserId,
				UpdateTime:   time.Now(),
				Version:      1,
			})
			addedMemberIds = append(addedMemberIds, v)
		}
	}

	//2.项目成员
	if isAllMember != nil && *isAllMember == true {
		//全选，获取所有成员
		tempMemberIds := []int64{}
		resp := orgfacade.GetOrgUserIds(orgvo.GetOrgUserIdsReq{
			OrgId: orgId,
		})

		if resp.Failure() {
			log.Error(resp.Error())
			return nil, nil, resp.Error()
		}

		for _, info := range resp.Data {
			tempMemberIds = append(tempMemberIds, info)
		}
		memberIds = tempMemberIds
	} else {
		//默认创建者也是项目成员
		if owner != currentUserId {
			if bool, _ := slice.Contain(memberIds, currentUserId); !bool {
				memberIds = append(memberIds, currentUserId)
			}
		}
	}
	memberIds = slice.SliceUniqueInt64(memberIds)

	if len(memberIds) != 0 {
		verifyOrgUserFlag := orgfacade.VerifyOrgUsersRelaxed(orgId, memberIds)
		if !verifyOrgUserFlag {
			log.Error("存在用户组织校验失败")
			return memberEntities, nil, errs.VerifyOrgError
		}
		memberRelationIds, idErr := idfacade.ApplyMultiplePrimaryIdRelaxed(consts.TableProjectRelation, len(memberIds))
		if idErr != nil {
			log.Error(idErr)
			return memberEntities, nil, idErr
		}
		for i, v := range memberIds {
			memberEntities = append(memberEntities, po.PpmProProjectRelation{
				Id:           memberRelationIds.Ids[i].Id,
				OrgId:        orgId,
				ProjectId:    projectId,
				RelationId:   v,
				RelationType: consts.ProjectRelationTypeParticipant,
				Creator:      currentUserId,
				CreateTime:   time.Now(),
				IsDelete:     consts.AppIsNoDelete,
				Status:       consts.ProjectMemberEffective,
				Updator:      currentUserId,
				UpdateTime:   time.Now(),
				Version:      1,
			})
			addedMemberIds = append(addedMemberIds, v)
		}
	}

	//成员部门处理
	if len(departmentIds) > 0 {
		departmentIds = slice.SliceUniqueInt64(departmentIds)
		if ok, _ := slice.Contain(departmentIds, int64(0)); ok {
			departmentIds = []int64{0}
		} else {
			verifyDepartment := orgfacade.VerifyDepartments(orgvo.VerifyDepartmentsReq{DepartmentIds: departmentIds, OrgId: orgId})
			if !verifyDepartment.IsTrue {
				log.Errorf("存在无效部门, 组织id:【%d】,部门：【%s】", orgId, json.ToJsonIgnoreError(departmentIds))
				return memberEntities, nil, errs.DepartmentNotExist
			}
		}

		memberRelationIds, idErr := idfacade.ApplyMultiplePrimaryIdRelaxed(consts.TableProjectRelation, len(departmentIds))
		if idErr != nil {
			log.Error(idErr)
			return memberEntities, nil, idErr
		}
		for i, v := range departmentIds {
			memberEntities = append(memberEntities, po.PpmProProjectRelation{
				Id:           memberRelationIds.Ids[i].Id,
				OrgId:        orgId,
				ProjectId:    projectId,
				RelationId:   v,
				RelationType: consts.ProjectRelationTypeDepartmentParticipant,
				Creator:      currentUserId,
				Updator:      currentUserId,
				UpdateTime:   time.Now(),
			})
		}

		//todo 成员（为了日历和群聊服务）
	}

	return memberEntities, addedMemberIds, nil
}

// 我参与的
func GetParticipantMembers(orgId, currentUserId int64) ([]int64, errs.SystemErrorInfo) {
	projectIdsNeed := []int64{}
	memberEntities := &[]*po.PpmProProjectRelation{}
	err := mysql.SelectAllByCond((&po.PpmProProjectRelation{}).TableName(), db.Cond{
		consts.TcIsDelete:     db.Eq(consts.AppIsNoDelete),
		consts.TcRelationType: db.Eq(consts.ProjectRelationTypeParticipant),
		consts.TcRelationId:   db.Eq(currentUserId),
		consts.TcOrgId:        orgId,
	}, memberEntities)
	if err != nil {
		return projectIdsNeed, errs.BuildSystemErrorInfo(errs.MysqlOperateError, err)
	}
	for _, v := range *memberEntities {
		projectIdsNeed = append(projectIdsNeed, v.ProjectId)
	}

	return projectIdsNeed, nil
}

func GetProjectMemberInfo(projectIds []int64, orgId int64, creatorIds []int64) (map[int64][]bo.UserIDInfoBo, map[int64][]bo.UserIDInfoBo, map[int64][]bo.UserIDInfoBo, map[int64]bo.UserIDInfoBo, errs.SystemErrorInfo) {
	conn, err := mysql.GetConnect()
	if err != nil {
		return nil, nil, nil, nil, errs.BuildSystemErrorInfo(errs.MysqlOperateError, err)
	}

	relatedInfo := &[]bo.RelationInfoTypeBo{}
	err1 := conn.Select("relation_id", "relation_type", "project_id").From("ppm_pro_project_relation").
		Where(db.Cond{
			consts.TcIsDelete:     consts.AppIsNoDelete,
			consts.TcProjectId:    db.In(projectIds),
			consts.TcStatus:       1,
			consts.TcOrgId:        orgId,
			consts.TcRelationType: db.In([]int64{consts.ProjectRelationTypeOwner, consts.ProjectRelationTypeParticipant}),
		}).All(relatedInfo)
	if err1 != nil {
		log.Error(err1)
		return nil, nil, nil, nil, errs.BuildSystemErrorInfo(errs.MysqlOperateError)
	}
	creatorInfo := map[int64]bo.UserIDInfoBo{}

	ownerInfo := map[int64][]bo.UserIDInfoBo{}
	participantInfo := map[int64][]bo.UserIDInfoBo{}
	followerInfo := map[int64][]bo.UserIDInfoBo{}
	allRelationIds := []int64{}
	for _, v := range *relatedInfo {
		allRelationIds = append(allRelationIds, v.RelationId)
	}

	creatorIds = slice.SliceUniqueInt64(creatorIds)

	allRelationIds = append(allRelationIds, creatorIds...)
	allRelationIds = slice.SliceUniqueInt64(allRelationIds)
	allUserInfo, err := orgfacade.GetBaseUserInfoBatchRelaxed(orgId, allRelationIds)
	if err != nil {
		return nil, nil, nil, nil, errs.BuildSystemErrorInfo(errs.ProjectDomainError, err)
	}
	userInfoById := map[int64]bo.BaseUserInfoBo{}
	for _, v := range allUserInfo {
		userInfoById[v.UserId] = v
	}
	for _, v := range creatorIds {
		userInfo, ok := userInfoById[v]
		if !ok {
			continue
		}
		temp := bo.UserIDInfoBo{}
		temp.Id = userInfo.UserId
		temp.Name = userInfo.Name
		temp.NamePy = userInfo.NamePy
		temp.Avatar = userInfo.Avatar
		temp.UserID = userInfo.UserId
		temp.EmplID = userInfo.OutUserId
		temp.IsDeleted = userInfo.OrgUserIsDelete == consts.AppIsDeleted
		temp.IsDisabled = userInfo.OrgUserStatus == consts.AppStatusDisabled
		creatorInfo[v] = temp
	}
	for _, v := range *relatedInfo {
		userInfo, ok := userInfoById[v.RelationId]
		if !ok {
			continue
		}
		temp := bo.UserIDInfoBo{}
		if userInfo.OrgUserIsDelete == consts.AppIsDeleted {
			continue
		}

		temp.Id = userInfo.UserId
		temp.Name = userInfo.Name
		temp.NamePy = userInfo.NamePy
		temp.Avatar = userInfo.Avatar
		temp.UserID = userInfo.UserId
		temp.EmplID = userInfo.OutUserId
		temp.IsDeleted = userInfo.OrgUserIsDelete == consts.AppIsDeleted
		temp.IsDisabled = userInfo.OrgUserStatus == consts.AppStatusDisabled
		if v.RelationType == consts.ProjectRelationTypeOwner {
			ownerInfo[v.ProjectId] = append(ownerInfo[v.ProjectId], temp)
		} else if v.RelationType == consts.ProjectRelationTypeParticipant {
			participantInfo[v.ProjectId] = append(participantInfo[v.ProjectId], temp)
			//} else if v.RelationType == consts.IssueRelationTypeFollower {
			//	followerInfo[v.ProjectId] = append(followerInfo[v.ProjectId], temp)
		}
	}

	//项目成员展示去重
	uniqueParticipantInfo := map[int64][]bo.UserIDInfoBo{}
	for i, bos := range participantInfo {
		relationIdsForProject := []int64{}
		for _, infoBo := range bos {
			if ok, _ := slice.Contain(relationIdsForProject, infoBo.UserID); !ok {
				uniqueParticipantInfo[i] = append(uniqueParticipantInfo[i], infoBo)
				relationIdsForProject = append(relationIdsForProject, infoBo.UserID)
			}
		}
	}
	return ownerInfo, uniqueParticipantInfo, followerInfo, creatorInfo, nil
}

// 项目的负责人、关注人 id。
func GetProjectParticipantIds(orgId, projectId int64) ([]int64, error) {
	uids := make([]int64, 0)
	memberEntities := &[]po.PpmProProjectRelation{}
	err := mysql.SelectAllByCond(consts.TableProjectRelation, db.Cond{
		consts.TcOrgId:        orgId,
		consts.TcProjectId:    projectId,
		consts.TcRelationType: db.In([]int{consts.ProjectRelationTypeParticipant, consts.ProjectRelationTypeOwner}),
		consts.TcIsDelete:     consts.AppIsNoDelete,
	}, memberEntities)
	if err != nil {
		return uids, err
	}
	for _, item := range *memberEntities {
		uids = append(uids, item.RelationId)
	}
	return uids, nil
}

// 分组插入新数据
func PaginationInsert(list []interface{}, domainObj mysql.Domain, tx ...sqlbuilder.Tx) errs.SystemErrorInfo {
	isTx := tx != nil && len(tx) > 0
	totalSize := len(list)
	batch := 1000
	offset := 0
	for {
		limit := offset + batch
		if totalSize < limit {
			limit = totalSize
		}
		oneBatch := list[offset:limit]
		if isTx {
			batchInsert := mysql.TransBatchInsert(tx[0], domainObj, oneBatch)
			if batchInsert != nil {
				log.Error(batchInsert)
				return errs.BuildSystemErrorInfo(errs.MysqlOperateError, batchInsert)
			}
		} else {
			batchInsert := mysql.BatchInsert(domainObj, oneBatch)
			if batchInsert != nil {
				log.Error(batchInsert)
				return errs.BuildSystemErrorInfo(errs.MysqlOperateError, batchInsert)
			}
		}

		if totalSize <= limit {
			break
		}
		offset += batch
	}
	return nil
}

func GetProjectAllMember(orgId, projectId int64, page, size int) (int64, []bo.ProjectRelationBo, errs.SystemErrorInfo) {
	conn, err := mysql.GetConnect()
	if err != nil {
		log.Error(err)
		return 0, nil, errs.MysqlOperateError
	}

	cond := db.Cond{
		consts.TcOrgId:     orgId,
		consts.TcIsDelete:  consts.AppIsNoDelete,
		consts.TcProjectId: projectId,
		//consts.TcRelationId:db.NotEq(0),
		consts.TcRelationType: db.In([]int64{consts.ProjectRelationTypeOwner, consts.ProjectRelationTypeParticipant, consts.ProjectRelationTypeDepartmentParticipant}),
	}

	pos := &[]po.PpmProProjectRelation{}
	//获取所有成员（最小的relation_id代表最高的用户角色（负责人），最早的创建时间表示加入时间，随机挑选一名操作人）
	//项目部门放在前面
	mid := conn.Select(db.Raw("relation_id, relation_type, create_time, creator")).
		From(consts.TableProjectRelation).
		Where(cond).OrderBy("relation_type asc")

	selectErr := mid.All(pos)
	if selectErr != nil {
		log.Error(selectErr)
		return 0, nil, errs.MysqlOperateError
	}

	bos := &[]bo.ProjectRelationBo{}
	copyErr := copyer.Copy(pos, bos)
	if copyErr != nil {
		log.Error(copyErr)
		return 0, nil, errs.BuildSystemErrorInfo(errs.ObjectCopyError, err)
	}

	if len(*bos) == 0 {
		return 0, *bos, nil
	}
	userIds := []int64{}
	deptIds := []int64{}
	//排列顺序 重新组合（部门->负责人->成员）
	userMap := map[int64]bo.ProjectRelationBo{}
	deptMap := map[int64]bo.ProjectRelationBo{}
	ownerId := int64(0)
	for _, relationBo := range *bos {
		if relationBo.RelationType != consts.ProjectRelationTypeDepartmentParticipant {
			if _, ok := userMap[relationBo.RelationId]; !ok {
				userIds = append(userIds, relationBo.RelationId)
				//一个是为了去重，一个是为了把负责人放到第一个
				userMap[relationBo.RelationId] = relationBo
			}
			if relationBo.RelationType == consts.ProjectRelationTypeOwner {
				ownerId = relationBo.RelationId
			}
		} else {
			deptIds = append(deptIds, relationBo.RelationId)
			deptMap[relationBo.RelationId] = relationBo
		}
	}

	res := []bo.ProjectRelationBo{}
	if len(deptIds) > 0 {
		departmentsInfo := orgfacade.Departments(orgvo.DepartmentsReqVo{
			Page: nil,
			Size: nil,
			Params: &vo.DepartmentListReq{
				DepartmentIds: deptIds,
			},
			OrgId: orgId,
		})
		if departmentsInfo.Failure() {
			log.Error(departmentsInfo.Error())
			return 0, nil, departmentsInfo.Error()
		}
		if ok, _ := slice.Contain(deptIds, int64(0)); ok {
			if relationInfo, ok := deptMap[int64(0)]; ok {
				res = append(res, relationInfo)
			}
		}
		for _, department := range departmentsInfo.DepartmentList.List {
			if relationInfo, ok := deptMap[department.ID]; ok {
				res = append(res, relationInfo)
			}
		}
	}
	if len(userIds) > 0 {
		deletedUserIds := []int64{}
		userInfos := orgfacade.GetUserInfoByUserIds(orgvo.GetUserInfoByUserIdsReqVo{
			UserIds: userIds,
			OrgId:   orgId,
		})
		if userInfos.Failure() {
			log.Error(userInfos.Error())
			return 0, nil, userInfos.Error()
		}
		if userInfos.GetUserInfoByUserIdsRespVo != nil {
			for _, respVo := range *userInfos.GetUserInfoByUserIdsRespVo {
				if respVo.OrgUserIsDelete == 1 {
					deletedUserIds = append(deletedUserIds, respVo.UserId)
				}
			}
		}
		//负责人放到上面
		if ok, _ := slice.Contain(deletedUserIds, ownerId); !ok {
			res = append(res, userMap[ownerId])
		}
		for i, relationBo := range userMap {
			if i == ownerId {
				continue
			}
			if ok, _ := slice.Contain(deletedUserIds, i); !ok {
				res = append(res, relationBo)
			}
		}
	}

	count := len(res)
	if page > 0 && size > 0 {
		offset := (page - 1) * size
		end := offset + size
		if offset > count {
			res = []bo.ProjectRelationBo{}
		} else {
			if end > count {
				end = count
			}
			res = res[offset:end]
		}
	}

	return int64(count), res, nil
}

func GetProjectMembers(orgId, projectId int64, relationTypes []int64) ([]bo.ProjectRelationBo, errs.SystemErrorInfo) {
	pos := &[]po.PpmProProjectRelation{}
	cond := db.Cond{
		consts.TcIsDelete:     consts.AppIsNoDelete,
		consts.TcOrgId:        orgId,
		consts.TcRelationType: db.In(relationTypes),
	}
	if projectId > 0 {
		cond[consts.TcProjectId] = projectId
	}
	err := mysql.SelectAllByCond(consts.TableProjectRelation, cond, pos)
	if err != nil {
		log.Error(err)
		return nil, errs.MysqlOperateError
	}

	bos := &[]bo.ProjectRelationBo{}
	copyErr := copyer.Copy(pos, bos)
	if copyErr != nil {
		log.Error(copyErr)
		return nil, errs.BuildSystemErrorInfo(errs.ObjectCopyError, err)
	}

	return *bos, nil
}

func GetProjectMemberDepartmentsInfo(orgId, projectId int64) ([]bo.DepartmentSimpleInfoBo, errs.SystemErrorInfo) {
	pos := &[]po.PpmProProjectRelation{}
	err := mysql.SelectAllByCond(consts.TableProjectRelation, db.Cond{
		consts.TcIsDelete:     consts.AppIsNoDelete,
		consts.TcOrgId:        orgId,
		consts.TcProjectId:    projectId,
		consts.TcRelationType: consts.ProjectRelationTypeDepartmentParticipant,
	}, pos)
	if err != nil {
		log.Error(err)
		return nil, errs.MysqlOperateError
	}

	departmentIds := []int64{}
	for _, relation := range *pos {
		departmentIds = append(departmentIds, relation.RelationId)
	}
	res := []bo.DepartmentSimpleInfoBo{}
	if len(departmentIds) == 0 {
		return res, nil
	}
	departmentIds = slice.SliceUniqueInt64(departmentIds)

	departmentsInfo := orgfacade.Departments(orgvo.DepartmentsReqVo{
		Page: nil,
		Size: nil,
		Params: &vo.DepartmentListReq{
			DepartmentIds: departmentIds,
		},
		OrgId: orgId,
	})
	if departmentsInfo.Failure() {
		log.Error(departmentsInfo.Error())
		return nil, departmentsInfo.Error()
	}

	if departmentsInfo.DepartmentList != nil {
		for _, department := range departmentsInfo.DepartmentList.List {
			res = append(res, bo.DepartmentSimpleInfoBo{
				ID:   department.ID,
				Name: department.Name,
			})
		}
	}
	deptUserCountResp := orgfacade.GetUserCountByDeptIds(&orgvo.GetUserCountByDeptIdsReq{
		OrgId:   orgId,
		DeptIds: departmentIds,
	})
	if deptUserCountResp.Failure() {
		log.Error(deptUserCountResp.Error())
		return nil, deptUserCountResp.Error()
	}
	if ok, _ := slice.Contain(departmentIds, int64(0)); ok {
		res = append(res, bo.DepartmentSimpleInfoBo{
			ID:   0,
			Name: "全部",
		})
	}
	if deptUserCountResp.Data != nil {
		userCount := deptUserCountResp.Data.UserCount
		for i, re := range res {
			if count, ok := userCount[re.ID]; ok {
				res[i].UserCount = int64(count)
			}
		}
	}

	return res, nil
}
