package msgfacade

import (
	"fmt"

	"github.com/star-table/startable-server/app/facade"
	"github.com/star-table/startable-server/common/core/config"
	"github.com/star-table/startable-server/common/core/consts"
	"github.com/star-table/startable-server/common/model/vo"
	"github.com/star-table/startable-server/common/model/vo/msgvo"
)

func FixAddOrderForFeiShu(req msgvo.FixAddOrderForFeiShuReqVo) vo.VoidErr {
	respVo := &vo.VoidErr{}
	reqUrl := fmt.Sprintf("%s/api/msgsvc/fixAddOrderForFeiShu", config.GetPreUrl("msgsvc"))
	requestBody := &req.Input
	err := facade.Request(consts.HttpMethodGet, reqUrl, nil, nil, requestBody, respVo)
	if err.Failure() {
		respVo.Err = err
	}
	return *respVo
}

func GetFailMsgList(req msgvo.GetFailMsgListReqVo) msgvo.GetFailMsgListRespVo {
	respVo := &msgvo.GetFailMsgListRespVo{}
	reqUrl := fmt.Sprintf("%s/api/msgsvc/getFailMsgList", config.GetPreUrl("msgsvc"))
	queryParams := map[string]interface{}{}
	queryParams["orgId"] = req.OrgId
	queryParams["msgType"] = req.MsgType
	queryParams["page"] = req.Page
	queryParams["size"] = req.Size
	err := facade.Request(consts.HttpMethodGet, reqUrl, queryParams, nil, nil, respVo)
	if err.Failure() {
		respVo.Err = err
	}
	return *respVo
}
