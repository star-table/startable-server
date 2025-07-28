package automationfacade

import (
	"fmt"

	automationPb "gitea.bjx.cloud/LessCode/interface/golang/automation/v1"
	"github.com/star-table/startable-server/app/facade"
	"github.com/star-table/startable-server/common/core/config"
	"github.com/star-table/startable-server/common/core/consts"
	"github.com/star-table/startable-server/common/core/logger"
	"github.com/star-table/startable-server/common/model/vo/commonvo"
)

var log = logger.GetDefaultLogger()

var (
	ApiV1PrefixInner = "/inner/v1"
	ApiV1Prefix      = "/v1"
)

func GlobalSwitchOff(orgId, userId int64) *commonvo.GlobalSwitchOffResp {
	respVo := &commonvo.GlobalSwitchOffResp{Data: &automationPb.GlobalSwitchOffReply{}}
	reqUrl := fmt.Sprintf("%s%s/global/switch/off", config.GetPreUrl(consts.ServiceAutomation), ApiV1PrefixInner)
	err := facade.RequestWithCommonHeaderGrpc(orgId, userId, consts.HttpMethodPost, reqUrl, nil, nil, &automationPb.GlobalSwitchOffReq{}, respVo.Data)
	if err.Failure() {
		respVo.Err = err
	}
	return respVo
}

func ApplyTemplate(req *commonvo.ApplyTemplateReq) *commonvo.ApplyTemplateResp {
	respVo := &commonvo.ApplyTemplateResp{Data: &automationPb.ApplyTemplateReply{}}
	reqUrl := fmt.Sprintf("%s%s/workflows/applyTemplate", config.GetPreUrl(consts.ServiceAutomation), ApiV1PrefixInner)
	err := facade.RequestWithCommonHeaderGrpc(req.OrgId, req.UserId, consts.HttpMethodPost, reqUrl, nil, nil, req.Input, respVo.Data)
	if err.Failure() {
		respVo.Err = err
	}
	return respVo
}

func GetTodo(orgId, userId, todoId int64) *commonvo.GetTodoResp {
	respVo := &commonvo.GetTodoResp{Data: &automationPb.GetTodoReply{}}
	reqUrl := fmt.Sprintf("%s%s/todo/%d", config.GetPreUrl(consts.ServiceAutomation), ApiV1Prefix, todoId)
	err := facade.RequestWithCommonHeaderGrpc(orgId, userId, consts.HttpMethodGet, reqUrl, nil, nil, nil, respVo.Data)
	if err.Failure() {
		respVo.Err = err
	}
	return respVo
}

func UpdateTodo(req *commonvo.UpdateTodoReq) *commonvo.UpdateTodoResp {
	respVo := &commonvo.UpdateTodoResp{Data: &automationPb.UpdateTodoReply{}}
	reqUrl := fmt.Sprintf("%s%s/todo", config.GetPreUrl(consts.ServiceAutomation), ApiV1PrefixInner)
	err := facade.RequestWithCommonHeaderGrpc(req.OrgId, req.UserId, consts.HttpMethodPut, reqUrl, nil, nil, req.Input, respVo.Data)
	if err.Failure() {
		respVo.Err = err
	}
	return respVo
}
