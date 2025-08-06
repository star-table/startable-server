package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/star-table/startable-server/app/service/projectsvc/service"
	"github.com/star-table/startable-server/common/core/errs"

	"github.com/star-table/startable-server/common/model/vo"
	"github.com/star-table/startable-server/common/model/vo/projectvo"
)

func InnerCreateTodoHook(c *gin.Context) {
	req := projectvo.InnerCreateTodoHookReq{}
	if err := c.ShouldBindJSON(&req); err != nil {
		// Replaced: c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		c.JSON(http.StatusOK, &vo.VoidErr{Err: vo.NewErr(errs.BuildSystemErrorInfo(errs.ReqParamsValidateError, err))})
		return
	}

	err := service.CreateTodoHook(&req)
	if err != nil {
		// Replaced: c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()}); return
	}
	c.JSON(http.StatusOK, &vo.VoidErr{Err: vo.NewErr(err)})
}
