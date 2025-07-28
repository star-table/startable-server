package handler

import (
	"github.com/star-table/startable-server/app/facade/projectfacade"
	"github.com/star-table/startable-server/common/core/errs"
	"github.com/star-table/startable-server/common/model/vo"
	"github.com/star-table/startable-server/common/model/vo/projectvo"
	"github.com/gin-gonic/gin"
)

type todoHandler struct{}

var TodoHandler todoHandler

func (todoHandler) Urge(c *gin.Context) {
	cacheUserInfo, sysErr := GetCacheUserInfo(c)
	if sysErr != nil {
		Fail(c, sysErr)
		return
	}

	req := vo.TodoUrgeReq{}
	if err := c.BindJSON(&req); err != nil {
		Fail(c, errs.ReqParamsValidateError)
		return
	}

	respVo := projectfacade.TodoUrge(&projectvo.TodoUrgeReq{
		OrgId:  cacheUserInfo.OrgId,
		UserId: cacheUserInfo.UserId,
		Input: &projectvo.TodoUrgeInput{
			TodoId: req.TodoId,
			Msg:    req.Msg,
		},
	})
	if respVo.Failure() {
		Fail(c, respVo.Error())
	} else {
		Success(c, respVo.Void)
	}
}
