package api

import (
	"net/http"
	"github.com/gin-gonic/gin"
	msgsvcService "github.com/star-table/startable-server/app/service/msgsvc/service"
	"github.com/star-table/startable-server/common/core/errs"
	"github.com/star-table/startable-server/common/model/vo"
	"github.com/star-table/startable-server/common/model/vo/msgvo"
)

// FixAddOrderForFeiShu 重跑消费订单异常的订单消息
func FixAddOrderForFeiShu(c *gin.Context) {
	var msg msgvo.FixAddOrderForFeiShuReqVo
	if err := c.ShouldBindJSON(&msg); err != nil {
		response := vo.VoidErr{Err: vo.NewErr(errs.BuildSystemErrorInfo(errs.ParamError, err))}
		c.JSON(http.StatusOK, response)
		return
	}
	
	err := msgsvcService.FixAddOrderForFeiShu(msg.Input)
	if err != nil {
		response := vo.VoidErr{Err: vo.NewErr(errs.BuildSystemErrorInfo(errs.ServerError, err))}
		c.JSON(http.StatusOK, response)
		return
	}
	response := vo.VoidErr{Err: vo.NewErr(nil)}
	c.JSON(http.StatusOK, response)
}
