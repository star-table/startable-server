package api

import (
	"fmt"
	"net/http"
	"github.com/gin-gonic/gin"
	msgsvcService "github.com/star-table/startable-server/app/service/msgsvc/service"
	"github.com/star-table/startable-server/common/core/errs"
	"github.com/star-table/startable-server/common/model/vo"
	"github.com/star-table/startable-server/common/model/vo/msgvo"
)

// PushMsgToMq 推送消息到MQ
func PushMsgToMq(c *gin.Context) {
	var msg msgvo.PushMsgToMqReqVo
	if err := c.ShouldBindJSON(&msg); err != nil {
		response := vo.VoidErr{Err: vo.NewErr(errs.BuildSystemErrorInfo(errs.ParamError, err))}
		c.JSON(http.StatusOK, response)
		return
	}
	
	err := msgsvcService.PushMsgToMq(msg)
	if err != nil {
		response := vo.VoidErr{Err: vo.NewErr(errs.BuildSystemErrorInfo(errs.ServerError, err))}
		c.JSON(http.StatusOK, response)
		return
	}
	response := vo.VoidErr{Err: vo.NewErr(nil)}
	c.JSON(http.StatusOK, response)
}

// InsertMqConsumeFailMsg 插入MQ消费失败消息
func InsertMqConsumeFailMsg(c *gin.Context) {
	var msg msgvo.InsertMqConsumeFailMsgReqVo
	if err := c.ShouldBindJSON(&msg); err != nil {
		response := vo.VoidErr{Err: vo.NewErr(errs.BuildSystemErrorInfo(errs.ParamError, err))}
		c.JSON(http.StatusOK, response)
		return
	}
	
	err := msgsvcService.InsertMqConsumeFailMsg(msg)
	if err != nil {
		response := vo.VoidErr{Err: vo.NewErr(errs.BuildSystemErrorInfo(errs.ServerError, err))}
		c.JSON(http.StatusOK, response)
		return
	}
	response := vo.VoidErr{Err: vo.NewErr(nil)}
	c.JSON(http.StatusOK, response)
}

// GetFailMsgList 获取失败消息列表
func GetFailMsgList(c *gin.Context) {
	var req msgvo.GetFailMsgListReqVo
	if err := c.ShouldBindJSON(&req); err != nil {
		response := msgvo.GetFailMsgListRespVo{Err: vo.NewErr(errs.BuildSystemErrorInfo(errs.ParamError, err))}
		c.JSON(http.StatusOK, response)
		return
	}
	
	res, err := msgsvcService.GetFailMsgList(req)
	if err != nil {
		response := msgvo.GetFailMsgListRespVo{Err: vo.NewErr(errs.BuildSystemErrorInfo(errs.ServerError, err))}
		c.JSON(http.StatusOK, response)
		return
	}
	response := msgvo.GetFailMsgListRespVo{Data: msgvo.GetFailMsgListRespData{MsgList: res}, Err: vo.NewErr(nil)}
	c.JSON(http.StatusOK, response)
}

// UpdateMsgStatus 更新消息状态
func UpdateMsgStatus(c *gin.Context) {
	var req msgvo.UpdateMsgStatusReqVo
	if err := c.ShouldBindJSON(&req); err != nil {
		response := vo.VoidErr{Err: vo.NewErr(errs.BuildSystemErrorInfo(errs.ParamError, err))}
		c.JSON(http.StatusOK, response)
		return
	}
	
	err := msgsvcService.UpdateMsgStatus(req)
	if err != nil {
		response := vo.VoidErr{Err: vo.NewErr(errs.BuildSystemErrorInfo(errs.ServerError, err))}
		c.JSON(http.StatusOK, response)
		return
	}
	response := vo.VoidErr{Err: vo.NewErr(nil)}
	c.JSON(http.StatusOK, response)
}

// WriteSomeFailedMsg 将一些场景下的任务处理失败的消息记录到表中，使用 msgType 进行区分查询。
func WriteSomeFailedMsg(c *gin.Context) {
	var msg msgvo.WriteSomeFailedMsgReqVo
	if err := c.ShouldBindJSON(&msg); err != nil {
		response := vo.VoidErr{Err: vo.NewErr(errs.BuildSystemErrorInfo(errs.ParamError, err))}
		c.JSON(http.StatusOK, response)
		return
	}
	
	err := msgsvcService.WriteSomeFailedMsg(msg)
	if err != nil {
		// 处理错误，打印错误栈，并返回错误。
		errMsg := fmt.Sprintf("%+v", err)
		log.Error(errMsg)
		busiErr := errs.BuildSystemErrorInfo(errs.ServerError, err)
		response := vo.VoidErr{Err: vo.NewErr(busiErr)}
		c.JSON(http.StatusOK, response)
		return
	}
	response := vo.VoidErr{Err: vo.NewErr(nil)}
	c.JSON(http.StatusOK, response)
}
