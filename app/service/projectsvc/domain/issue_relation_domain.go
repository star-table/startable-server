package domain

import (
	"strconv"
	"time"

	"github.com/spf13/cast"

	"upper.io/db.v3"
	"upper.io/db.v3/lib/sqlbuilder"

	tablePb "gitea.bjx.cloud/LessCode/interface/golang/table/v1"
	"github.com/star-table/startable-server/app/facade/idfacade"
	"github.com/star-table/startable-server/app/facade/resourcefacade"
	"github.com/star-table/startable-server/app/service/projectsvc/dao"
	"github.com/star-table/startable-server/app/service/projectsvc/po"
	"github.com/star-table/startable-server/common/core/consts"
	"github.com/star-table/startable-server/common/core/errs"
	"github.com/star-table/startable-server/common/core/types"
	"github.com/star-table/startable-server/common/core/util/asyn"
	"github.com/star-table/startable-server/common/core/util/copyer"
	"github.com/star-table/startable-server/common/core/util/json"
	"github.com/star-table/startable-server/common/core/util/slice"
	"github.com/star-table/startable-server/common/core/util/uuid"
	"github.com/star-table/startable-server/common/extra/lc_helper"
	"github.com/star-table/startable-server/common/library/cache"
	"github.com/star-table/startable-server/common/library/db/mysql"
	"github.com/star-table/startable-server/common/model/bo"
	"github.com/star-table/startable-server/common/model/vo/resourcevo"
)

func DeleteAllIssueRelation(tx sqlbuilder.Tx, operatorId int64, orgId int64, issueIds []int64, recycleVersionId int64) errs.SystemErrorInfo {
	//删除之前的关联
	_, err := mysql.TransUpdateSmartWithCond(tx, consts.TableIssueRelation, db.Cond{
		consts.TcOrgId:    orgId,
		consts.TcIssueId:  db.In(issueIds),
		consts.TcIsDelete: consts.AppIsNoDelete,
	}, mysql.Upd{
		consts.TcIsDelete: consts.AppIsDeleted,
		consts.TcVersion:  recycleVersionId,
		consts.TcUpdator:  operatorId,
	})
	if err != nil {
		log.Error(err)
		return errs.BuildSystemErrorInfo(errs.MysqlOperateError)
	}
	return nil
}

func DeleteAllIssueRelationByIds(orgId, operatorId int64, relationIds []int64, recycleVersionId int64, tx ...sqlbuilder.Tx) errs.SystemErrorInfo {
	//删除之前的关联
	_, err := mysql.TransUpdateSmartWithCond(tx[0], consts.TableIssueRelation, db.Cond{
		consts.TcOrgId:    orgId,
		consts.TcId:       db.In(relationIds),
		consts.TcIsDelete: consts.AppIsNoDelete,
	}, mysql.Upd{
		consts.TcIsDelete:   consts.AppIsDeleted,
		consts.TcUpdator:    operatorId,
		consts.TcUpdateTime: time.Now(),
		consts.TcVersion:    recycleVersionId,
	})
	if err != nil {
		log.Error(err)
		return errs.BuildSystemErrorInfo(errs.MysqlOperateError)
	}
	return nil
}

func DeleteIssueRelation(operatorId int64, issueBo bo.IssueBo, relationType int) errs.SystemErrorInfo {
	orgId := issueBo.OrgId
	issueId := issueBo.Id
	//删除之前的关联
	_, err := mysql.UpdateSmartWithCond(consts.TableIssueRelation, db.Cond{
		consts.TcOrgId:        orgId,
		consts.TcIssueId:      issueId,
		consts.TcRelationType: relationType,
		consts.TcIsDelete:     consts.AppIsNoDelete,
	}, mysql.Upd{
		consts.TcIsDelete: consts.AppIsDeleted,
		consts.TcUpdator:  operatorId,
	})
	if err != nil {
		log.Error(err)
		return errs.BuildSystemErrorInfo(errs.MysqlOperateError)
	}
	return nil
}

func DeleteIssueRelationByIds(operatorId int64, issueBo bo.IssueBo, relationType int, relationIds []int64) errs.SystemErrorInfo {
	if relationIds == nil || len(relationIds) == 0 {
		return nil
	}

	orgId := issueBo.OrgId
	issueId := issueBo.Id

	err := mysql.TransX(func(tx sqlbuilder.Tx) error {
		//删除之前的关联
		_, err := mysql.TransUpdateSmartWithCond(tx, consts.TableIssueRelation, db.Cond{
			consts.TcOrgId:        orgId,
			consts.TcIssueId:      issueId,
			consts.TcRelationId:   db.In(relationIds),
			consts.TcRelationType: relationType,
			consts.TcIsDelete:     consts.AppIsNoDelete,
		}, mysql.Upd{
			consts.TcIsDelete: consts.AppIsDeleted,
			consts.TcUpdator:  operatorId,
		})
		if err != nil {
			log.Error(err)
			return err
		}

		if relationType == consts.IssueRelationTypeResource {
			//删除文件
			resp := resourcefacade.DeleteResource(resourcevo.DeleteResourceReqVo{
				Input: bo.DeleteResourceBo{
					ResourceIds: relationIds,
					UserId:      operatorId,
					OrgId:       orgId,
					ProjectId:   issueBo.ProjectId,
					IssueId:     issueId,
				},
			})
			if resp.Failure() {
				log.Error(resp.Message)
				return resp.Error()
			}
		}

		return nil
	})

	if err != nil {
		log.Error(err)
		return errs.BuildSystemErrorInfo(errs.MysqlOperateError, err)
	}

	return nil
}

func UpdateIssueRelationSingle(operatorId int64, issueBo *bo.IssueBo, relationType int, newUserIds []int64) (*bo.IssueRelationBo, errs.SystemErrorInfo) {
	bos, err := UpdateIssueRelation(operatorId, issueBo, relationType, newUserIds, "")
	if err != nil {
		return nil, err
	}
	if len(bos) == 0 {
		return nil, errs.BuildSystemErrorInfo(errs.IssueDomainError)
	}
	return &bos[0], nil
}

func UpdateIssueRelation(operatorId int64, issueBo *bo.IssueBo, relationType int, newUserIds []int64, relationCode string) ([]bo.IssueRelationBo, errs.SystemErrorInfo) {
	orgId := issueBo.OrgId
	issueId := issueBo.Id

	//防止项目成员重复插入
	uid := uuid.NewUuid()
	issueIdStr := strconv.FormatInt(issueId, 10)
	relationTypeStr := strconv.Itoa(relationType)
	lockKey := consts.AddIssueRelationLock + issueIdStr + "#" + relationTypeStr
	suc, err := cache.TryGetDistributedLock(lockKey, uid)
	if err != nil {
		log.Errorf("获取%s锁时异常 %v", lockKey, err)
		return nil, errs.TryDistributedLockError
	}
	if suc {
		defer func() {
			if _, err := cache.ReleaseDistributedLock(lockKey, uid); err != nil {
				log.Error(err)
			}
		}()
	} else {
		return nil, errs.BuildSystemErrorInfo(errs.GetDistributedLockError)
	}

	//预先查询已有的关联
	issueRelations := &[]po.PpmPriIssueRelation{}
	err5 := mysql.SelectAllByCond(consts.TableIssueRelation, db.Cond{
		consts.TcIssueId:      issueBo.Id,
		consts.TcRelationId:   db.In(newUserIds),
		consts.TcRelationType: relationType,
		consts.TcIsDelete:     consts.AppIsNoDelete,
	}, issueRelations)
	if err5 != nil {
		log.Error(err5)
		return nil, errs.BuildSystemErrorInfo(errs.MysqlOperateError)
	}

	//check，去掉已有的关联
	if len(*issueRelations) > 0 {
		alreadyExistRelationIdMap := map[int64]bool{}
		for _, issueRelation := range *issueRelations {
			alreadyExistRelationIdMap[issueRelation.RelationId] = true
		}
		notRelationUserIds := make([]int64, 0)
		for _, newUserId := range newUserIds {
			if _, ok := alreadyExistRelationIdMap[newUserId]; !ok {
				notRelationUserIds = append(notRelationUserIds, newUserId)
			}
		}
		newUserIds = notRelationUserIds
	}
	newUserIds = slice.SliceUniqueInt64(newUserIds)

	issueRelationBos := make([]bo.IssueRelationBo, len(newUserIds))
	if len(newUserIds) == 0 {
		return issueRelationBos, nil
	}

	ids, err := idfacade.ApplyMultiplePrimaryIdRelaxed(consts.TableIssueRelation, len(newUserIds))
	if err != nil {
		log.Errorf("id generate: %q\n", err)
		return nil, errs.BuildSystemErrorInfo(errs.ApplyIdError, err)
	}

	issueRelationPos := make([]po.PpmPriIssueRelation, len(newUserIds))
	for i, newUserId := range newUserIds {
		id := ids.Ids[i].Id
		issueRelation := &po.PpmPriIssueRelation{}
		issueRelation.Id = id
		issueRelation.OrgId = orgId
		issueRelation.ProjectId = issueBo.ProjectId
		issueRelation.IssueId = issueBo.Id
		issueRelation.RelationId = newUserId
		issueRelation.RelationType = relationType
		issueRelation.Creator = operatorId
		issueRelation.Updator = operatorId
		issueRelation.IsDelete = consts.AppIsNoDelete
		issueRelation.RelationCode = relationCode
		issueRelationPos[i] = *issueRelation

		issueRelationBos[i] = bo.IssueRelationBo{
			Id:           id,
			OrgId:        issueBo.OrgId,
			IssueId:      issueBo.Id,
			RelationId:   newUserId,
			RelationType: relationType,
			Creator:      operatorId,
			CreateTime:   types.NowTime(),
			Updator:      operatorId,
			UpdateTime:   types.NowTime(),
			Version:      1,
		}
	}

	err2 := PaginationInsert(slice.ToSlice(issueRelationPos), &po.PpmPriIssueRelation{})
	if err2 != nil {
		log.Errorf("mysql.BatchInsert(): %q\n", err2)
		return nil, errs.BuildSystemErrorInfo(errs.MysqlOperateError, err2)
	}
	return issueRelationBos, nil
}

func GetIssueRelationIdsByRelateType(orgId int64, issueId int64, relationType int) (*[]int64, errs.SystemErrorInfo) {
	issueParticipantRelations, _, err := dao.SelectIssueRelationByPage(db.Cond{
		consts.TcOrgId:        orgId,
		consts.TcIssueId:      issueId,
		consts.TcRelationType: relationType,
		consts.TcIsDelete:     consts.AppIsNoDelete,
	}, bo.PageBo{
		Order: consts.TcCreateTime + " asc",
	})
	if err != nil {
		log.Error(err)
		return nil, errs.BuildSystemErrorInfo(errs.MysqlOperateError)
	}
	relationIds := make([]int64, len(*issueParticipantRelations))
	for i, participantRelation := range *issueParticipantRelations {
		relationIds[i] = participantRelation.RelationId
	}
	return &relationIds, nil
}

func GetIssueRelationByRelateTypeList(orgId int64, issueId int64, relationTypes []int) ([]bo.IssueRelationBo, errs.SystemErrorInfo) {
	issueParticipantRelations, _, err := dao.SelectIssueRelationByPage(db.Cond{
		consts.TcOrgId:        orgId,
		consts.TcIssueId:      issueId,
		consts.TcRelationType: db.In(relationTypes),
		consts.TcIsDelete:     consts.AppIsNoDelete,
	}, bo.PageBo{
		Order: consts.TcCreateTime + " desc",
	})
	if err != nil {
		log.Error(err)
		return nil, errs.BuildSystemErrorInfo(errs.MysqlOperateError)
	}
	relationBos := &[]bo.IssueRelationBo{}
	_ = copyer.Copy(issueParticipantRelations, relationBos)
	return *relationBos, nil
}

// 创建者也可以删除资源
func GetIssueResourceIdsByCreator(orgId int64, issueId int64, ids []int64, creatorId int64) (*[]int64, errs.SystemErrorInfo) {
	issueRelations, _, err := dao.SelectIssueRelationByPage(db.Cond{
		consts.TcId:           db.In(ids),
		consts.TcOrgId:        orgId,
		consts.TcIssueId:      issueId,
		consts.TcRelationType: consts.IssueRelationTypeResource,
		consts.TcCreator:      creatorId,
		consts.TcIsDelete:     consts.AppIsNoDelete,
	}, bo.PageBo{
		Order: consts.TcCreateTime + " desc",
	})
	if err != nil {
		log.Error(err)
		return nil, errs.BuildSystemErrorInfo(errs.MysqlOperateError)
	}
	relationIds := make([]int64, len(*issueRelations))
	for i, participantRelation := range *issueRelations {
		relationIds[i] = participantRelation.Id
	}
	return &relationIds, nil
}

func GetIssueRelationByResource(orgId int64, projectId int64, resourceIds []int64) (*[]po.PpmPriIssueRelation, errs.SystemErrorInfo) {
	issueParticipantRelations, _, err := dao.SelectIssueRelationByPage(db.Cond{
		consts.TcOrgId:        orgId,
		consts.TcProjectId:    projectId,
		consts.TcRelationType: consts.IssueRelationTypeResource,
		consts.TcIsDelete:     consts.AppIsNoDelete,
		consts.TcRelationId:   db.In(resourceIds),
	}, bo.PageBo{
		Order: consts.TcId + " desc",
	})
	if err != nil {
		log.Error(err)
		return nil, errs.BuildSystemErrorInfo(errs.MysqlOperateError)
	}
	return issueParticipantRelations, nil
}

func GetTotalResourceByRelationCond(cond db.Cond) (*[]po.PpmPriIssueRelation, errs.SystemErrorInfo) {
	pos := &[]po.PpmPriIssueRelation{}
	err := mysql.SelectAllByCond(consts.TableIssueRelation, cond, pos)
	if err != nil {
		log.Error(err)
		return nil, errs.BuildSystemErrorInfo(errs.MysqlOperateError, err)
	}
	resourceIds := []int64{}
	for _, value := range *pos {
		if isContain, _ := slice.Contain(resourceIds, value.RelationId); !isContain {
			resourceIds = append(resourceIds, value.RelationId)
		}
	}
	return pos, nil
}

func DeleteProjectAttachment(orgId, operatorId, projectId int64, resourceIds []int64) errs.SystemErrorInfo {
	issueRelationPos, err := GetIssueRelationByResource(orgId, projectId, resourceIds)
	if err != nil {
		log.Error(err)
		return err
	}
	realResourceIds := make([]int64, 0)
	relationIds := make([]int64, len(*issueRelationPos))
	realResourceMap := make(map[int64]bool)
	for index, value := range *issueRelationPos {
		realResourceMap[value.RelationId] = true
		relationIds[index] = value.Id
	}
	for key, _ := range realResourceMap {
		realResourceIds = append(realResourceIds, key)
	}
	if len(realResourceIds) != len(resourceIds) {
		return errs.InvalidResourceIdsError
	}
	_ = mysql.TransX(func(tx sqlbuilder.Tx) error {
		recycleVersionId, versionErr := AddRecycleRecord(orgId, operatorId, projectId, resourceIds, consts.RecycleTypeAttachment, tx)
		if versionErr != nil {
			log.Error(versionErr)
			return versionErr
		}
		err := DeleteAllIssueRelationByIds(orgId, operatorId, relationIds, recycleVersionId, tx)
		if err != nil {
			log.Error(err)
			return nil
		}

		projectinfo, err := GetProjectSimple(orgId, projectId)
		if err != nil {
			return err
		}
		deleteInput := bo.DeleteResourceBo{
			ResourceIds:      resourceIds,
			UserId:           operatorId,
			OrgId:            orgId,
			ProjectId:        projectId,
			RecycleVersionId: recycleVersionId,
			AppId:            projectinfo.AppId,
		}
		resp := resourcefacade.DeleteResource(resourcevo.DeleteResourceReqVo{Input: deleteInput})
		if resp.Failure() {
			log.Error(resp.Error())
			return resp.Error()
		}
		return nil
	})

	asyn.Execute(func() {
		reqVo := resourcevo.GetResourceByIdReqBody{
			ResourceIds: resourceIds,
		}
		req := resourcevo.GetResourceByIdReqVo{GetResourceByIdReqBody: reqVo}
		resp := resourcefacade.GetResourceById(req)
		resourceBos := resp.ResourceBos
		resourceNames := make([]string, len(resourceBos))
		resourceTrend := []bo.ResourceInfoBo{}
		for index, value := range resourceBos {
			resourceNames[index] = value.Name
			resourceTrend = append(resourceTrend, bo.ResourceInfoBo{
				Name:   value.Name,
				Url:    value.Host + value.Path,
				Size:   value.Size,
				Suffix: value.Suffix,
			})
		}

		trendBo := bo.ProjectTrendsBo{
			PushType:   consts.PushTypeDeleteResource,
			OrgId:      orgId,
			ProjectId:  projectId,
			OperatorId: operatorId,
			NewValue:   json.ToJsonIgnoreError(resourceNames),
			Ext: bo.TrendExtensionBo{
				ResourceInfo: resourceTrend,
			},
		}

		asyn.Execute(func() {
			PushProjectTrends(trendBo)
		})
		//asyn.Execute(func() {
		//	PushProjectThirdPlatformNotice(trendBo)
		//})
	})

	return nil
}

func GetRelationInfoByIssueIds(issueIds []int64, relationTypes []int) ([]bo.IssueRelationBo, errs.SystemErrorInfo) {
	relationInfos := &[]po.PpmPriIssueRelation{}
	cond := db.Cond{
		consts.TcIsDelete: consts.AppIsNoDelete,
		consts.TcIssueId:  db.In(issueIds),
	}
	if len(relationTypes) != 0 {
		cond[consts.TcRelationType] = db.In(relationTypes)
	}
	err := mysql.SelectAllByCond(consts.TableIssueRelation, cond, relationInfos)
	if err != nil {
		log.Error(err)
		return nil, errs.MysqlOperateError
	}

	bos := &[]bo.IssueRelationBo{}
	copyErr := copyer.Copy(relationInfos, bos)
	if copyErr != nil {
		log.Error(copyErr)
		return nil, errs.ObjectCopyError
	}

	return *bos, nil
}

func GetIssueRelationResource(page, size int) ([]bo.IssueRelationBo, errs.SystemErrorInfo) {
	pos := &[]po.PpmPriIssueRelation{}
	_, err := mysql.SelectAllByCondWithPageAndOrder(consts.TableIssueRelation, db.Cond{
		consts.TcRelationType: consts.IssueRelationTypeResource,
		consts.TcIsDelete:     consts.AppIsNoDelete,
	}, nil, page, size, nil, pos)

	if err != nil {
		log.Error(err)
		return nil, errs.MysqlOperateError
	}

	res := &[]bo.IssueRelationBo{}
	_ = copyer.Copy(pos, res)
	return *res, nil
}

// GetIssueIdsByProIds 通过项目id，获取项目下所有的任务 id
func GetIssueIdsByProIds(orgId int64, projectIds []int64) ([]int64, errs.SystemErrorInfo) {
	issueIds := make([]int64, 0)
	conds := []*tablePb.Condition{
		GetRowsCondition(consts.BasicFieldProjectId, tablePb.ConditionType_in, nil, projectIds),
	}
	issueInfoDatas, err := GetIssueInfosMapLc(orgId, 0, &tablePb.Condition{
		Type:       tablePb.ConditionType_and,
		Conditions: conds,
	}, []string{lc_helper.ConvertToFilterColumn(consts.BasicFieldIssueId)}, 0, 0)
	if err != nil {
		log.Errorf("[GetIssueIdsByProIds]GetIssueInfosMapLc err:%v, orgId:%v, projectIds:%v", err, orgId, projectIds)
		return nil, err
	}

	for _, data := range issueInfoDatas {
		issueId := cast.ToInt64(data[consts.BasicFieldIssueId])
		issueIds = append(issueIds, issueId)
	}

	return issueIds, nil
}

func AddResourceRelation(orgId, userId, projectId, issueId int64, resourceIds []int64, sourceType int, columnId string) errs.SystemErrorInfo {
	resourceRelationResp := resourcefacade.AddResourceRelationWithType(resourcevo.AddResourceRelationReq{
		OrgId:  orgId,
		UserId: userId,
		Input: &resourcevo.AddResourceRelationData{
			ProjectId:   projectId,
			IssueId:     issueId,
			ResourceIds: resourceIds,
			SourceType:  sourceType,
			ColumnId:    columnId,
		},
	})
	if resourceRelationResp.Failure() {
		log.Errorf("[AddResourceRelation] CreateResourceRelation err:%v, orgId:%v, userId:%v, projectId:%v, issueId:%v, resourceIds:%v",
			resourceRelationResp.Error(), orgId, userId, projectId, issueId, resourceIds)
		return resourceRelationResp.Error()
	}
	return nil
}

// GetResourceTypeByResourceIds 获取资源id和sourceType的对应关系
func GetResourceTypeByResourceIds(orgId, projectId int64, resourceIds []int64) (map[int64]int, errs.SystemErrorInfo) {
	relationList := resourcefacade.GetResourceRelationsByProjectId(resourcevo.GetResourceRelationsByProjectIdReqVo{
		UserId: 0,
		OrgId:  orgId,
		Input: resourcevo.GetResourceRelationsByProjectIdData{
			ProjectId:   projectId,
			SourceTypes: []int32{consts.OssPolicyTypeProjectResource, consts.OssPolicyTypeLesscodeResource},
		},
	})
	if relationList.Failure() {
		log.Errorf("[GetResourceTypeByResourceIds] err:%v, orgId:%d, projectId:%d, resourceIds:%v",
			relationList.Error(), orgId, projectId, resourceIds)
		return nil, relationList.Error()
	}
	relationMap := make(map[int64]resourcevo.ResourceRelationVo, len(relationList.ResourceRelations))
	for _, r := range relationList.ResourceRelations {
		relationMap[r.ResourceId] = r
	}

	// map[resourceId][sourceType]
	resourceIdType := map[int64]int{}
	for _, resourceId := range resourceIds {
		if _, ok := relationMap[resourceId]; ok {
			resourceIdType[resourceId] = relationMap[resourceId].SourceType
		}
	}

	return resourceIdType, nil
}
