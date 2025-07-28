package notice

import (
	"github.com/star-table/startable-server/app/facade/orgfacade"
	"github.com/star-table/startable-server/app/facade/tablefacade"
	"github.com/star-table/startable-server/app/service/projectsvc/domain"
	"github.com/star-table/startable-server/common/core/consts"
	"github.com/star-table/startable-server/common/core/errs"
	"github.com/star-table/startable-server/common/core/util/copyer"
	"github.com/star-table/startable-server/common/model/bo"
	"github.com/star-table/startable-server/common/model/vo/lc_table"
	"github.com/star-table/startable-server/common/model/vo/projectvo"
)

// GetBaseInfoForCardPush 卡片推送时查询一些基础信息
func GetBaseInfoForCardPush(issueTrendsBo *bo.IssueTrendsBo) (*projectvo.BaseInfoBoxForIssueCard, errs.SystemErrorInfo) {
	var err errs.SystemErrorInfo
	infoObj := &projectvo.BaseInfoBoxForIssueCard{}

	// orgInfo, err := orgfacade.GetBaseOrgInfoRelaxed(issueTrendsBo.SourceChannel, issueTrendsBo.OrgId)
	// if err != nil {
	// 	log.Errorf("[GetBaseInfoForCardPush] GetBaseOrgInfoRelaxed err: %v, orgId: %d", err, issueTrendsBo.OrgId)
	// 	return nil, err
	// }
	// infoObj.OrgInfo = *orgInfo

	issues, err := domain.GetIssueInfosLc(issueTrendsBo.OrgId, issueTrendsBo.OperatorId, []int64{issueTrendsBo.IssueId})
	if err != nil {
		log.Errorf("[GetBaseInfoForCardPush] err: %v, issueId: %v", err, issueTrendsBo.IssueId)
		return nil, err
	}
	if len(issues) > 0 {
		infoObj.IssueInfo = *issues[0]
	} else {
		err := errs.IssueNotExist
		log.Errorf("[GetBaseInfoForCardPush] err: %v, issueId: %d", err, issueTrendsBo.IssueId)
		return nil, err
	}

	// 查询任务负责人
	userIds := make([]int64, 0, len(infoObj.IssueInfo.OwnerIdI64)+1)
	userIds = append(userIds, issueTrendsBo.OperatorId)
	userIds = append(userIds, infoObj.IssueInfo.OwnerIdI64...)
	userInfoArr, err := orgfacade.GetBaseUserInfoBatchRelaxed(issueTrendsBo.OrgId, userIds)
	if err != nil {
		log.Errorf("[GetBaseInfoForCardPush] 查询组织 %d 用户信息出现异常 %v", issueTrendsBo.OrgId, err)
		return nil, err
	}
	userMap := make(map[int64]bo.BaseUserInfoBo, len(userInfoArr))
	for _, user := range userInfoArr {
		userMap[user.UserId] = user
	}
	if opUser, ok := userMap[issueTrendsBo.OperatorId]; ok {
		infoObj.OperateUser = opUser
	}
	//ownerUserNameArr := make([]string, 0)
	ownerInfos := make([]*bo.BaseUserInfoBo, 0)
	for _, uid := range infoObj.IssueInfo.OwnerIdI64 {
		if user, ok := userMap[uid]; ok {
			//ownerUserNameArr = append(ownerUserNameArr, user.Name)
			ownerInfos = append(ownerInfos, &user)
		}
	}
	infoObj.OwnerInfos = ownerInfos
	//infoObj.IssueOwnerNameArr = ownerUserNameArr

	projectAuthInfo, err := domain.LoadProjectAuthBo(issueTrendsBo.OrgId, issueTrendsBo.ProjectId)
	if err != nil {
		if err == errs.ProjectNotExist {
			// projectId为0，需要查询汇总表id
			summaryAppId, errSys := domain.GetOrgSummaryAppId(issueTrendsBo.OrgId)
			if errSys != nil {
				log.Errorf("[GetBaseInfoForCardPush] GetOrgSummaryAppId error:%v, orgId:%d, issueId:%d",
					errSys, issueTrendsBo.OrgId, issueTrendsBo.IssueId)
			}

			projectAuthInfo = &bo.ProjectAuthBo{
				Id:            0,
				AppId:         summaryAppId,
				Name:          consts.CardDefaultIssueProjectName,
				ProjectTypeId: consts.ProjectTypeCommon2022V47,
				Status:        0,
				PublicStatus:  consts.PublicProject,
			}
		} else {
			log.Errorf("[GetBaseInfoForCardPush] err: %v, projectId: %d", err, issueTrendsBo.ProjectId)
			return nil, err
		}
	}
	infoObj.ProjectAuthBo = *projectAuthInfo

	if issueTrendsBo.TableId > 0 {
		headers := make(map[string]lc_table.LcCommonField, 0)
		tableColumns, err := domain.GetTableColumnsMap(issueTrendsBo.OrgId, issueTrendsBo.TableId, nil)
		if err != nil {
			log.Errorf("[GetBaseInfoForCardPush] 获取表头失败 org:%d proj:%d table:%d user:%d, err: %v",
				issueTrendsBo.OrgId, issueTrendsBo.ProjectId, issueTrendsBo.TableId, issueTrendsBo.OperatorId, err)
			return nil, err
		}
		infoObj.ProjectTableColumn = tableColumns
		copyer.Copy(tableColumns, &headers)
		tableColumnMeta := make(map[string]lc_table.LcCommonField, 0)
		for columnId, c := range headers {
			if c.Field.Props.PushMsg {
				tableColumnMeta[columnId] = c
			}
		}
		infoObj.TableColumnMap = tableColumnMeta
	}

	// 查询任务所在表名
	if issueTrendsBo.TableId > 0 {
		tableInfo, err := domain.GetTableByTableId(issueTrendsBo.OrgId, issueTrendsBo.OperatorId, issueTrendsBo.TableId)
		if err != nil {
			log.Errorf("[GetBaseInfoForCardPush] GetTableByTableId err:%v, userId:%d, tableId:%d", err, issueTrendsBo.OperatorId, issueTrendsBo.TableId)
			return nil, err
		}
		infoObj.IssueTableInfo = *tableInfo
	} else {
		// err := errs.InvalidTableId
		// log.Errorf("[GetBaseInfoForCardPush] err: %v, issueId: %d", err, issueTrendsBo.IssueId)
		// return infoObj, err
		// 默认一个任务表
		infoObj.IssueTableInfo = projectvo.TableMetaData{
			Name: consts.DefaultTableName,
		}
	}

	// 查询父任务信息
	if issueTrendsBo.ParentId > 0 {
		parents, err := domain.GetIssueInfosLc(issueTrendsBo.OrgId, 0, []int64{issueTrendsBo.ParentId})
		if err != nil {
			log.Errorf("[GetBaseInfoForCardPush] GetIssueBo err: %v, parentId: %d", err, issueTrendsBo.ParentId)
			return nil, err
		}
		if len(parents) < 1 {
			log.Errorf("[GetBaseInfoForCardPush] not found issue parentId:%v", issueTrendsBo.ParentId)
			return nil, errs.IssueNotExist
		}
		parent := parents[0]
		infoObj.ParentIssue = *parent
	}

	// 查看详情、应用内查看等按钮的 url 链接
	links := domain.GetIssueLinks("", issueTrendsBo.OrgId, issueTrendsBo.IssueId)
	infoObj.IssueInfoUrl = links.SideBarLink
	infoObj.IssuePcUrl = links.Link

	// 查询任务的协作人信息
	collaboratorIds, err := tablefacade.GetDataCollaborateUserIds(issueTrendsBo.OrgId, 0, issueTrendsBo.DataId)
	if err != nil {
		log.Errorf("[GetBaseInfoForCardPush] GetDataCollaborateUserIds err: %v, DataId: %d", err, issueTrendsBo.DataId)
		return nil, err
	}
	infoObj.CollaboratorIds = collaboratorIds
	return infoObj, nil
}
