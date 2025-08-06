package api

import (
	"net/http"
	"github.com/gin-gonic/gin"
	"github.com/star-table/startable-server/app/service/projectsvc/service"
	"github.com/star-table/startable-server/common/model/vo"
	"github.com/star-table/startable-server/common/model/vo/projectvo"
)

func ProjectChatList(c *gin.Context) {
	var req projectvo.ProjectChatListReqVo
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var pageSize int64 = 10
	if req.Input.PageSize != nil {
		pageSize = *req.Input.PageSize
	}

	pageSize = 500

	res, err := service.ChatList(pageSize, req.UserId, req.Input.ProjectID, req.OrgId, req.Input.LastRelationID, req.Input.Name)
	response := projectvo.ProjectChatListRespVo{
		Err:  vo.NewErr(err),
		List: res,
	}
	c.JSON(http.StatusOK, response)
}

func UnrelatedChatList(c *gin.Context) {
	var req projectvo.UnrelatedChatListReqVo
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var pageSize int64 = 10
	if req.Input.PageSize != nil {
		pageSize = *req.Input.PageSize
	}

	pageSize = 500

	res, err := service.UnrelatedChatList(pageSize, req.UserId, req.Input.ProjectID, req.OrgId, req.Input.LastOutChatID, req.Input.Name)
	response := projectvo.ProjectChatListRespVo{
		Err:  vo.NewErr(err),
		List: res,
	}
	c.JSON(http.StatusOK, response)
}

func AddProjectChat(c *gin.Context) {
	var req projectvo.UpdateRelateChatReqVo
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	res, err := service.AddChat(req.Input.OutChatIds, req.OrgId, req.UserId, req.Input.ProjectID)
	response := vo.CommonRespVo{
		Err:  vo.NewErr(err),
		Void: res,
	}
	c.JSON(http.StatusOK, response)
}

func DisbandProjectChat(c *gin.Context) {
	var req projectvo.UpdateRelateChatReqVo
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	res, err := service.DisbandChat(req.Input.OutChatIds, req.OrgId, req.UserId, req.Input.ProjectID)
	response := vo.CommonRespVo{
		Err:  vo.NewErr(err),
		Void: res,
	}
	c.JSON(http.StatusOK, response)
}

func FsChatDisbandCallback(c *gin.Context) {
	var req projectvo.FsChatDisbandCallbackReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := service.FsChatDisbandCallback(req.OrgId, req.ChatId)
	response := vo.CommonRespVo{
		Err:  vo.NewErr(err),
		Void: nil,
	}
	c.JSON(http.StatusOK, response)
}

func GetProjectMainChatId(c *gin.Context) {
	var req projectvo.GetProjectMainChatIdReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	res, err := service.GetProjectMainChatId(req.OrgId, req.UserId, req.ProjectId, req.SourceChannel)
	response := projectvo.GetProjectMainChatIdResp{
		Err:    vo.NewErr(err),
		ChatId: res,
	}
	c.JSON(http.StatusOK, response)
}

func CheckIsShowProChatIcon(c *gin.Context) {
	var req projectvo.CheckIsShowProChatIconReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	isShow, err := service.CheckIsShowProChatIcon(req.OrgId, req.UserId, req.SourceChannel, req.Input)
	response := projectvo.CheckIsShowProChatIconResp{
		Err: vo.NewErr(err),
		Data: projectvo.CheckShowProChatIconRespData{
			IsShow: isShow,
		},
	}
	c.JSON(http.StatusOK, response)
}

func UpdateFsProjectChatPushSettings(c *gin.Context) {
	var req projectvo.UpdateFsProjectChatPushSettingsReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := service.UpdateFsProjectChatPushSettings(req.OrgId, req.UserId, req.SourceChannel, req.Input)
	response := vo.CommonRespVo{
		Err:  vo.NewErr(err),
		Void: &vo.Void{ID: 1},
	}
	c.JSON(http.StatusOK, response)
}

func GetFsProjectChatPushSettings(c *gin.Context) {
	var req projectvo.GetFsProjectChatPushSettingsReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	res, err := service.GetFsProjectChatPushSettings(req.OrgId, &req)
	response := projectvo.GetFsProjectChatPushSettingsResp{
		Err:  vo.NewErr(err),
		Data: res,
	}
	c.JSON(http.StatusOK, response)
}

func DeleteChatCallback(c *gin.Context) {
	var req projectvo.DeleteChatReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := service.DeleteChatCallback(req.OutOrgId, req.ChatId, req.SourceChannel)
	response := vo.CommonRespVo{
		Err:  vo.NewErr(err),
		Void: nil,
	}
	c.JSON(http.StatusOK, response)
}
