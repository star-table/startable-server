package domain

import (
	"time"

	"github.com/star-table/startable-server/common/core/errs"
	"github.com/star-table/startable-server/common/core/types"
	"github.com/star-table/startable-server/common/model/bo"
)

func SelectIssueAssignCount(orgId int64, projectId int64, rankTop int) ([]bo.IssueAssignCountBo, errs.SystemErrorInfo) {
	return nil, nil
}

func SelectIssueAssignTotalCount(orgId int64, projectId int64) (int64, errs.SystemErrorInfo) {
	return 0, nil

}

func GetIssueCountByOwnerId(orgId, ownerId int64) (int64, errs.SystemErrorInfo) {

	return 0, nil
}

// IssueDailyPersonalWorkCompletionStatForLc 任务完成统计图，根据endTime做groupBy，统计当前负责人在指定时间段之内的每天任务完成数量
// 通过无码查询统计
func IssueDailyPersonalWorkCompletionStatForLc(orgId, ownerId int64, startDate *types.Time, endDate *types.Time) ([]bo.IssueDailyPersonalWorkCompletionStatBo, errs.SystemErrorInfo) {
	return nil, nil

}

// 计算日期范围条件
// startDate: 开始时间
// endDate: 结束时间
// cond: 条件
// defaultDay: 默认查询多少天前
func CalDateRangeCond(startDate *types.Time, endDate *types.Time, defaultDay int) (*time.Time, *time.Time, errs.SystemErrorInfo) {
	if startDate == nil && endDate == nil {
		currentTime := time.Now()
		sd := types.Time(currentTime.AddDate(0, 0, -(defaultDay - 1)))
		ed := types.Time(currentTime)
		startDate = &sd
		endDate = &ed
	}
	if startDate != nil && endDate != nil {
	} else {
		log.Error("没有提供明确的时间范围")
		return nil, nil, errs.BuildSystemErrorInfo(errs.DateRangeError)
	}
	startTime := time.Time(*startDate)
	endTime := time.Time(*endDate)

	return &startTime, &endTime, nil
}
