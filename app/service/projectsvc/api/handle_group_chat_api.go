package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/star-table/startable-server/app/service/projectsvc/service"
	"github.com/star-table/startable-server/common/core/errs"

	"github.com/star-table/startable-server/common/model/vo"
	"github.com/star-table/startable-server/common/model/vo/projectvo"
)

// 项目群聊，响应用户指令，项目任务
func HandleGroupChatUserInsProIssue(c *gin.Context) {
	req := projectvo.HandleGroupChatUserInsProIssueReqVo{}
	if err := c.ShouldBindJSON(&req); err != nil {
		// Replaced: c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		c.JSON(http.StatusOK, projectvo.BoolRespVo{Err: vo.NewErr(errs.BuildSystemErrorInfo(errs.ReqParamsValidateError, err))})
		return
	}

	isOk, err := service.HandleGroupChatUserInsProIssue(req.Input)
	if err != nil {
		// Replaced: c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()}); return
	}
	c.JSON(http.StatusOK, projectvo.BoolRespVo{
		Err: vo.NewErr(err),
		Data: &vo.BoolResp{
			IsTrue: isOk,
		},
	})
}

// 项目群聊，响应用户指令，项目进展
func HandleGroupChatUserInsProProgress(c *gin.Context) {
	req := projectvo.HandleGroupChatUserInsProProgressReqVo{}
	if err := c.ShouldBindJSON(&req); err != nil {
		// Replaced: c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		c.JSON(http.StatusOK, projectvo.BoolRespVo{Err: vo.NewErr(errs.BuildSystemErrorInfo(errs.ReqParamsValidateError, err))})
		return
	}

	isOk, err := service.HandleGroupChatUserInsProProgress(req.Input)
	if err != nil {
		// Replaced: c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()}); return
	}
	c.JSON(http.StatusOK, projectvo.BoolRespVo{
		Err: vo.NewErr(err),
		Data: &vo.BoolResp{
			IsTrue: isOk,
		},
	})
}

// 项目群聊，响应用户指令，项目设置
func HandleGroupChatUserInsProjectSettings(c *gin.Context) {
	req := projectvo.HandleGroupChatUserInsProjectSettingsReq{}
	if err := c.ShouldBindJSON(&req); err != nil {
		// Replaced: c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		c.JSON(http.StatusOK, projectvo.BoolRespVo{Err: vo.NewErr(errs.BuildSystemErrorInfo(errs.ReqParamsValidateError, err))})
		return
	}

	isOk, err := service.HandleGroupChatUserInsProjectSettings(req.OpenChatId, req.SourceChannel)
	if err != nil {
		// Replaced: c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()}); return
	}
	c.JSON(http.StatusOK, projectvo.BoolRespVo{
		Err: vo.NewErr(err),
		Data: &vo.BoolResp{
			IsTrue: isOk,
		},
	})
}

// 项目群聊，响应用户指令：@用户姓名
func HandleGroupChatUserInsAtUserName(c *gin.Context) {
	req := projectvo.HandleGroupChatUserInsAtUserNameReqVo{}
	if err := c.ShouldBindJSON(&req); err != nil {
		// Replaced: c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		c.JSON(http.StatusOK, projectvo.BoolRespVo{Err: vo.NewErr(errs.BuildSystemErrorInfo(errs.ReqParamsValidateError, err))})
		return
	}

	isOk, err := service.HandleGroupChatAtUserName(req.Input)
	if err != nil {
		// Replaced: c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()}); return
	}
	c.JSON(http.StatusOK, projectvo.BoolRespVo{
		Err: vo.NewErr(err),
		Data: &vo.BoolResp{
			IsTrue: isOk,
		},
	})
}

// 项目群聊，响应用户指令:@用户姓名 任务标题1
func HandleGroupChatUserInsAtUserNameWithIssueTitle(c *gin.Context) {
	req := projectvo.HandleGroupChatUserInsAtUserNameWithIssueTitleReqVo{}
	if err := c.ShouldBindJSON(&req); err != nil {
		// Replaced: c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		c.JSON(http.StatusOK, projectvo.BoolRespVo{Err: vo.NewErr(errs.BuildSystemErrorInfo(errs.ReqParamsValidateError, err))})
		return
	}

	isOk, err := service.HandleGroupChatAtUserNameWithIssueTitle(req.Input)
	if err != nil {
		// Replaced: c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()}); return
	}
	c.JSON(http.StatusOK, projectvo.BoolRespVo{
		Err: vo.NewErr(err),
		Data: &vo.BoolResp{
			IsTrue: isOk,
		},
	})
}

// GetProjectIdByChatId 通过 chatId 获取群聊对应的项目id
func GetProjectIdByChatId(c *gin.Context) {
	req := projectvo.GetProjectIdByChatIdReqVo{}
	if err := c.ShouldBindJSON(&req); err != nil {
		// Replaced: c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		c.JSON(http.StatusOK, projectvo.GetProjectIdByChatIdRespVo{Err: vo.NewErr(errs.BuildSystemErrorInfo(errs.ReqParamsValidateError, err))})
		return
	}

	resp, err := service.GetProjectIdByChatId(req)
	if err != nil {
		// Replaced: c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()}); return
	}
	c.JSON(http.StatusOK, projectvo.GetProjectIdByChatIdRespVo{
		Err:  vo.NewErr(err),
		Data: resp,
	})
}

// GetProjectIdsByChatId 通过 chatId 获取群聊对应的项目。因为现在一个群聊可以关联多个项目 id todo
func GetProjectIdsByChatId(c *gin.Context) {
	req := projectvo.GetProjectIdsByChatIdReqVo{}
	if err := c.ShouldBindJSON(&req); err != nil {
		// Replaced: c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		c.JSON(http.StatusOK, projectvo.GetProjectIdsByChatIdRespVo{Err: vo.NewErr(errs.BuildSystemErrorInfo(errs.ReqParamsValidateError, err))})
		return
	}

	resp, err := service.GetProjectIdsByChatId(req.OrgId, req)
	if err != nil {
		// Replaced: c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()}); return
	}
	c.JSON(http.StatusOK, projectvo.GetProjectIdsByChatIdRespVo{
		Err:  vo.NewErr(err),
		Data: resp,
	})
}
