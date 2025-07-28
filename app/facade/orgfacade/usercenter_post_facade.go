package orgfacade

import (
	"fmt"

	"github.com/star-table/startable-server/app/facade"
	"github.com/star-table/startable-server/common/core/config"
	"github.com/star-table/startable-server/common/core/consts"
	"github.com/star-table/startable-server/common/model/vo/orgvo"
)

func GetUserCountByDeptIds(req *orgvo.GetUserCountByDeptIdsReq) *orgvo.GetUserCountByDeptIdsResp {
	respVo := &orgvo.GetUserCountByDeptIdsResp{}
	reqUrl := fmt.Sprintf("%s/usercenter/inner/api/v1/dept/getUserCountByDeptIds", config.GetPreUrl("usercentersvc"))
	err := facade.Request(consts.HttpMethodPost, reqUrl, nil, nil, req, respVo)
	if err.Failure() {
		respVo.Err = err
	}
	return respVo
}

func GetUserDeptIds(req *orgvo.GetUserDeptIdsReq) *orgvo.GetUserDeptIdsResp {
	respVo := &orgvo.GetUserDeptIdsResp{}
	reqUrl := fmt.Sprintf("%s/usercenter/inner/api/v1/dept/getUserDeptIds", config.GetPreUrl("usercentersvc"))
	err := facade.Request(consts.HttpMethodPost, reqUrl, nil, nil, &req, respVo)
	if err.Failure() {
		respVo.Err = err
	}
	return respVo
}

func GetUserDeptIdsBatch(req *orgvo.GetUserDeptIdsBatchReq) *orgvo.GetUserDeptIdsBatchResp {
	respVo := &orgvo.GetUserDeptIdsBatchResp{}
	reqUrl := fmt.Sprintf("%s/usercenter/inner/api/v1/dept/getUserDeptIdsBatch", config.GetPreUrl("usercentersvc"))
	err := facade.Request(consts.HttpMethodPost, reqUrl, nil, nil, &req, respVo)
	if err.Failure() {
		respVo.Err = err
	}
	return respVo
}

func GetUserIdsByDeptIds(req *orgvo.GetUserIdsByDeptIdsReq) *orgvo.GetUserIdsByDeptIdsResp {
	respVo := &orgvo.GetUserIdsByDeptIdsResp{}
	reqUrl := fmt.Sprintf("%s/usercenter/inner/api/v1/dept/getUserIdsByDeptIds", config.GetPreUrl("usercentersvc"))
	err := facade.Request(consts.HttpMethodPost, reqUrl, nil, nil, &req, respVo)
	if err.Failure() {
		respVo.Err = err
	}
	return respVo
}

func GetUserDeptIdsWithParentId(req orgvo.GetUserDeptIdsWithParentIdReq) *orgvo.GetUserDeptIdsWithParentIdResp {
	respVo := &orgvo.GetUserDeptIdsWithParentIdResp{}
	reqUrl := fmt.Sprintf("%s/usercenter/inner/api/v1/dept/getUserDeptIdsWithParentId", config.GetPreUrl("usercentersvc"))
	err := facade.Request(consts.HttpMethodPost, reqUrl, nil, nil, &req, respVo)
	if err.Failure() {
		respVo.Err = err
	}
	return respVo
}
