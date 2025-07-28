package domain

import (
	tablePb "gitea.bjx.cloud/LessCode/interface/golang/table/v1"
	"github.com/star-table/startable-server/app/facade/resourcefacade"
	"github.com/star-table/startable-server/app/service/projectsvc/po"
	"github.com/star-table/startable-server/common/core/consts"
	"github.com/star-table/startable-server/common/core/errs"
	"github.com/star-table/startable-server/common/library/db/mysql"
	"github.com/star-table/startable-server/common/model/vo/projectvo"
	"github.com/star-table/startable-server/common/model/vo/resourcevo"
	"upper.io/db.v3"
)

func GetProjectCountByOwnerId(orgId, ownerId int64) (int64, errs.SystemErrorInfo) {
	finishedIds, _ := consts.GetStatusIdsByCategory(orgId, consts.ProcessStatusCategoryProject, consts.StatusTypeComplete)
	if len(finishedIds) == 0 {
		log.Errorf("[GetProjectCountByOwnerId] GetStatusIdsByCategory err: not finished status id")
		return 0, errs.BuildSystemErrorInfo(errs.CacheProxyError, errs.ProcessStatusNotExist)
	}

	total, err1 := mysql.SelectCountByCond(consts.TableProject, db.Cond{
		consts.TcOrgId:    orgId,
		consts.TcOwner:    ownerId,
		consts.TcStatus:   db.NotIn(finishedIds),
		consts.TcIsDelete: consts.AppIsNoDelete,
		consts.TcIsFiling: consts.AppIsNotFilling,
	})
	if err1 != nil {
		log.Error(err1)
		return 0, errs.BuildSystemErrorInfo(errs.MysqlOperateError)
	}
	return int64(total), nil
}

func PayLimitNum(orgId int64) (*projectvo.PayLimitNumRespData, errs.SystemErrorInfo) {
	projectNum, err := mysql.SelectCountByCond(consts.TableProject, db.Cond{
		consts.TcIsDelete:     consts.AppIsNoDelete,
		consts.TcOrgId:        orgId,
		consts.TcTemplateFlag: consts.TemplateFalse,
	})
	if err != nil {
		return nil, errs.MysqlOperateError
	}

	projectPos := []*po.PpmProProject{}
	err = mysql.SelectAllByCond(consts.TableProject, db.Cond{
		consts.TcOrgId:        orgId,
		consts.TcIsDelete:     consts.AppIsNoDelete,
		consts.TcTemplateFlag: consts.TemplateTrue,
	}, &projectPos)
	if err != nil {
		log.Errorf("[PayLimitNum] err:%v, orgId:%v", err, orgId)
		return nil, errs.MysqlOperateError
	}
	projectIds := make([]int64, 0, len(projectPos))
	for _, pro := range projectPos {
		projectIds = append(projectIds, pro.Id)
	}

	//issueNum, err := mysql.SelectCountByCond(consts.TableIssue, db.Cond{
	//	consts.TcOrgId:     orgId,
	//	consts.TcProjectId: db.NotIn(db.Raw("select id from ppm_pro_project where org_id = ? and is_delete = 2 and template_flag = 1", orgId)),
	//})
	//if err != nil {
	//	log.Error(err)
	//	return nil, errs.MysqlOperateError
	//}

	var condition *tablePb.Condition
	if len(projectIds) > 0 {
		condition = &tablePb.Condition{}
		condition.Type = tablePb.ConditionType_and
		condition.Conditions = append(condition.Conditions, GetRowsCondition(consts.BasicFieldProjectId, tablePb.ConditionType_not_in, nil, projectIds))
	}
	issueNum, errSys := getRowsCountByCondition(orgId, 0, condition)
	if errSys != nil {
		log.Errorf("[PayLimitNum] getRowsCountByCondition err:%v, orgId:%v", err, orgId)
		return nil, errSys
	}

	// 加上被删除的任务
	deletedIssueNum, errCache := GetDeletedIssueNum(orgId)
	if errCache != nil {
		log.Errorf("[PayLimitNum] GetDeletedIssueNum err:%v, orgId:%v", errCache, orgId)
	}

	log.Infof("[PayLimitNum] issueNum:%v, deletedIssueNum:%v, orgId:%v", issueNum, deletedIssueNum, orgId)

	resourceResp := resourcefacade.CacheResourceSize(resourcevo.CacheResourceSizeReq{OrgId: orgId})
	if resourceResp.Failure() {
		log.Error(resourceResp.Error())
		return nil, resourceResp.Error()
	}

	return &projectvo.PayLimitNumRespData{
		ProjectNum: int64(projectNum),
		IssueNum:   issueNum + deletedIssueNum,
		FileSize:   resourceResp.Size,
		DashNum:    int64(projectNum), // 现有仪表盘数量等于项目数。后续如果规则发生变化，这里的逻辑也会发生变化
	}, nil
}
