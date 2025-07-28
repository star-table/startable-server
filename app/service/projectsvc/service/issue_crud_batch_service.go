package service

import (
	"github.com/star-table/startable-server/app/service/projectsvc/domain"
	"github.com/star-table/startable-server/common/core/errs"
	"github.com/star-table/startable-server/common/model/vo/projectvo"
	"github.com/star-table/startable-server/common/model/vo/uservo"
)

// SyncBatchCreateIssue 每 BATCH_SIZE 条为一个批次，执行同步批量插入
func SyncBatchCreateIssue(req *projectvo.BatchCreateIssueReqVo, withoutAuth bool, triggerBy *projectvo.TriggerBy) (
	[]map[string]interface{}, map[string]*uservo.MemberDept, map[string]map[string]interface{}, errs.SystemErrorInfo) {
	return domain.SyncBatchCreateIssue(req, withoutAuth, triggerBy)
}

// AsyncBatchCreateIssue 每 BATCH_SIZE 条为一个批次，执行异步批量插入
func AsyncBatchCreateIssue(req *projectvo.BatchCreateIssueReqVo, withoutAuth bool, triggerBy *projectvo.TriggerBy, asyncTaskId string) {
	domain.AsyncBatchCreateIssue(req, withoutAuth, triggerBy, asyncTaskId)
}

func BatchCreateIssue(req *projectvo.BatchCreateIssueReqVo, withoutAuth bool, triggerBy *projectvo.TriggerBy,
	asyncTaskId string, createCount int64) (
	[]map[string]interface{}, map[string]*uservo.MemberDept, map[string]map[string]interface{}, errs.SystemErrorInfo) {
	return domain.BatchCreateIssue(req, withoutAuth, triggerBy, asyncTaskId, createCount)
}

// SyncBatchUpdateIssue 每BATCH_SIZE条为一个批次，执行同步批量更新
func SyncBatchUpdateIssue(req *projectvo.BatchUpdateIssueReqInnerVo, withoutAuth bool, triggerBy *projectvo.TriggerBy) errs.SystemErrorInfo {
	return domain.SyncBatchUpdateIssue(req, withoutAuth, triggerBy)
}

func BatchUpdateIssue(reqVo *projectvo.BatchUpdateIssueReqInnerVo, withoutAuth bool, triggerBy *projectvo.TriggerBy) errs.SystemErrorInfo {
	return domain.BatchUpdateIssue(reqVo, withoutAuth, triggerBy)
}

func BatchAuditIssue(reqVo *projectvo.BatchAuditIssueReqVo) ([]int64, errs.SystemErrorInfo) {
	return domain.BatchAuditIssue(reqVo)
}

func BatchUrgeIssue(reqVo *projectvo.BatchUrgeIssueReqVo) ([]int64, errs.SystemErrorInfo) {
	return domain.BatchUrgeIssue(reqVo)
}
