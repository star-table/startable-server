package service

import (
	"github.com/star-table/startable-server/app/service/projectsvc/domain"
	"github.com/star-table/startable-server/common/core/errs"
	"github.com/star-table/startable-server/common/core/util/json"
	"github.com/star-table/startable-server/common/model/vo/projectvo"
)

func GetMenu(orgId int64, appId int64) (projectvo.MenuData, errs.SystemErrorInfo) {
	menuInfo, err := domain.GetMenu(orgId, appId)
	if err != nil {
		return projectvo.MenuData{}, err
	}

	tables, _ := GetTables(orgId, 0, appId)
	//if err != nil {
	//	return projectvo.MenuData{}, err
	//}

	dashboards, err := domain.GetAppDashboards(orgId, appId)
	if err != nil {
		return projectvo.MenuData{}, err
	}

	result := projectvo.MenuData{
		AppId:      menuInfo.AppId,
		Config:     menuInfo.Config,
		Dashboards: dashboards,
	}
	if tables != nil {
		result.Tables = tables.Tables
	}

	return result, nil
}

func SaveMenu(orgId int64, input projectvo.SaveMenuData) (projectvo.SaveMenuResp, errs.SystemErrorInfo) {
	configStr := json.ToJsonIgnoreError(input.Config)
	appId, err := domain.SaveMenu(orgId, input.AppId, configStr)
	if err != nil {
		log.Errorf("[SaveMenu] 错误, err: %v", err)
		return projectvo.SaveMenuResp{}, err
	}
	return projectvo.SaveMenuResp{AppId: appId}, nil
}
