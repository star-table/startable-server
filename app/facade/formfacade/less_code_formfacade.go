package formfacade

import (
	"fmt"

	"github.com/star-table/startable-server/app/facade"

	"github.com/star-table/startable-server/common/core/config"
	"github.com/star-table/startable-server/common/core/consts"
	"github.com/star-table/startable-server/common/model/vo/appvo"
	"github.com/star-table/startable-server/common/model/vo/formvo"
)

// FormCreateIssue 新增任务数据
func FormCreateIssue(req *appvo.FormCreateIssueReq) *appvo.FormCreateOneResp {
	respVo := &appvo.FormCreateOneResp{}
	reqUrl := fmt.Sprintf("%s/form/inner/api/v1/apps/%d/values", config.GetPreUrl(consts.ServiceNameLcForm), req.AppId)
	queryParams := map[string]interface{}{}
	err := facade.Request(consts.HttpMethodPost, reqUrl, queryParams, nil, req, respVo)
	if err.Failure() {
		respVo.Err = err
	}
	return respVo
}

func FormCreatePriority(req *appvo.FormCreatePriorityReq) *appvo.FormCreateOneResp {
	respVo := &appvo.FormCreateOneResp{}
	reqUrl := fmt.Sprintf("%s/form/inner/api/v1/apps/%d/values", config.GetPreUrl(consts.ServiceNameLcForm), req.AppId)
	queryParams := map[string]interface{}{}
	err := facade.Request(consts.HttpMethodPost, reqUrl, queryParams, nil, req, respVo)
	if err.Failure() {
		respVo.Err = err
	}
	return respVo
}

// LessCreateIssue 新增任务数据
func LessCreateIssue(req formvo.LessCreateIssueReq) *formvo.LessUpdateIssueResp {
	respVo := &formvo.LessUpdateIssueResp{}
	reqUrl := fmt.Sprintf("%s/form/inner/api/v1/apps/%d/values", config.GetPreUrl(consts.ServiceNameLcForm), req.AppId)
	queryParams := map[string]interface{}{}
	err := facade.Request(consts.HttpMethodPost, reqUrl, queryParams, nil, req, respVo)
	if err.Failure() {
		respVo.Err = err
	}
	return respVo
}

// LessUpdateIssue 更新任务数据
func LessUpdateIssue(req formvo.LessUpdateIssueReq) *formvo.LessUpdateIssueResp {
	respVo := &formvo.LessUpdateIssueResp{}
	reqUrl := fmt.Sprintf("%s/form/inner/api/v1/apps/%d/values", config.GetPreUrl(consts.ServiceNameLcForm), req.AppId)
	queryParams := map[string]interface{}{}
	err := facade.Request(consts.HttpMethodPut, reqUrl, queryParams, nil, req, respVo)
	if err.Failure() {
		respVo.Err = err
	}
	return respVo
}

// LessUpdateIssueBatchRaw 批量更新任务数据
func LessUpdateIssueBatchRaw(req *formvo.LessUpdateIssueBatchReq) *formvo.LessCommonIssueResp {
	respVo := &formvo.LessCommonIssueResp{}
	reqUrl := fmt.Sprintf("%s/form/inner/api/v1/apps/%d/values/batch", config.GetPreUrl(consts.ServiceNameLcForm), req.AppId)
	queryParams := map[string]interface{}{}
	err := facade.Request(consts.HttpMethodPut, reqUrl, queryParams, nil, req, respVo)
	if err.Failure() {
		respVo.Err = err
	}
	return respVo
}

// LessRecycleIssue 回收任务数据
func LessRecycleIssue(req formvo.LessRecycleIssueReq) *formvo.LessRecycleIssueResp {
	respVo := &formvo.LessRecycleIssueResp{}
	reqUrl := fmt.Sprintf("%s/form/inner/api/v1/apps/%d/values/recycle", config.GetPreUrl(consts.ServiceNameLcForm), req.AppId)
	queryParams := map[string]interface{}{}
	err := facade.Request(consts.HttpMethodPost, reqUrl, queryParams, nil, req, respVo)
	if err.Failure() {
		respVo.Err = err
	}
	return respVo
}

// LessRecoverIssue 恢复任务数据
func LessRecoverIssue(req formvo.LessRecoverIssueReq) *formvo.LessRecoverIssueResp {
	respVo := &formvo.LessRecoverIssueResp{}
	reqUrl := fmt.Sprintf("%s/form/inner/api/v1/apps/%d/values/recover", config.GetPreUrl(consts.ServiceNameLcForm), req.AppId)
	queryParams := map[string]interface{}{}
	err := facade.Request(consts.HttpMethodPost, reqUrl, queryParams, nil, req, respVo)
	if err.Failure() {
		respVo.Err = err
	}
	return respVo
}
