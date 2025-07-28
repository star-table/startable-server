package msgfacade

import (
	"fmt"

	"github.com/star-table/startable-server/app/facade"
	"github.com/star-table/startable-server/common/core/config"
	"github.com/star-table/startable-server/common/core/consts"
	"github.com/star-table/startable-server/common/model/vo"
	"github.com/star-table/startable-server/common/model/vo/msgvo"
)

func InsertMqConsumeFailMsg(req msgvo.InsertMqConsumeFailMsgReqVo) vo.VoidErr {
	respVo := &vo.VoidErr{}
	reqUrl := fmt.Sprintf("%s/api/msgsvc/insertMqConsumeFailMsg", config.GetPreUrl("msgsvc"))
	queryParams := map[string]interface{}{}
	queryParams["msgType"] = req.MsgType
	queryParams["orgId"] = req.OrgId
	requestBody := &req.Msg
	err := facade.Request(consts.HttpMethodPost, reqUrl, queryParams, nil, requestBody, respVo)
	if err.Failure() {
		respVo.Err = err
	}
	return *respVo
}

func PushMsgToMq(req msgvo.PushMsgToMqReqVo) vo.VoidErr {
	respVo := &vo.VoidErr{}
	reqUrl := fmt.Sprintf("%s/api/msgsvc/pushMsgToMq", config.GetPreUrl("msgsvc"))
	queryParams := map[string]interface{}{}
	queryParams["msgType"] = req.MsgType
	queryParams["orgId"] = req.OrgId
	requestBody := &req.Msg
	err := facade.Request(consts.HttpMethodPost, reqUrl, queryParams, nil, requestBody, respVo)
	if err.Failure() {
		respVo.Err = err
	}
	return *respVo
}

func SendLoginSMS(req msgvo.SendLoginSMSReqVo) vo.VoidErr {
	respVo := &vo.VoidErr{}
	reqUrl := fmt.Sprintf("%s/api/msgsvc/sendLoginSMS", config.GetPreUrl("msgsvc"))
	queryParams := map[string]interface{}{}
	queryParams["phoneNumber"] = req.PhoneNumber
	queryParams["code"] = req.Code
	err := facade.Request(consts.HttpMethodPost, reqUrl, queryParams, nil, nil, respVo)
	if err.Failure() {
		respVo.Err = err
	}
	return *respVo
}

func SendMail(req msgvo.SendMailReqVo) vo.VoidErr {
	respVo := &vo.VoidErr{}
	reqUrl := fmt.Sprintf("%s/api/msgsvc/sendMail", config.GetPreUrl("msgsvc"))
	requestBody := &req.Input
	err := facade.Request(consts.HttpMethodPost, reqUrl, nil, nil, requestBody, respVo)
	if err.Failure() {
		respVo.Err = err
	}
	return *respVo
}

func SendSMS(req msgvo.SendSMSReqVo) vo.VoidErr {
	respVo := &vo.VoidErr{}
	reqUrl := fmt.Sprintf("%s/api/msgsvc/sendSMS", config.GetPreUrl("msgsvc"))
	requestBody := &req.Input
	err := facade.Request(consts.HttpMethodPost, reqUrl, nil, nil, requestBody, respVo)
	if err.Failure() {
		respVo.Err = err
	}
	return *respVo
}

func UpdateMsgStatus(req msgvo.UpdateMsgStatusReqVo) vo.VoidErr {
	respVo := &vo.VoidErr{}
	reqUrl := fmt.Sprintf("%s/api/msgsvc/updateMsgStatus", config.GetPreUrl("msgsvc"))
	queryParams := map[string]interface{}{}
	queryParams["msgId"] = req.MsgId
	queryParams["newStatus"] = req.NewStatus
	err := facade.Request(consts.HttpMethodPost, reqUrl, queryParams, nil, nil, respVo)
	if err.Failure() {
		respVo.Err = err
	}
	return *respVo
}

func WriteSomeFailedMsg(req msgvo.WriteSomeFailedMsgReqVo) vo.VoidErr {
	respVo := &vo.VoidErr{}
	reqUrl := fmt.Sprintf("%s/api/msgsvc/writeSomeFailedMsg", config.GetPreUrl("msgsvc"))
	queryParams := map[string]interface{}{}
	queryParams["orgId"] = req.OrgId
	requestBody := &req.Input
	err := facade.Request(consts.HttpMethodPost, reqUrl, queryParams, nil, requestBody, respVo)
	if err.Failure() {
		respVo.Err = err
	}
	return *respVo
}
