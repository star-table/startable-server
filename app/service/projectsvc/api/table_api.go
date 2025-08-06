package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/star-table/startable-server/app/service/projectsvc/service"
	"github.com/star-table/startable-server/common/core/errs"

	"github.com/star-table/startable-server/common/model/vo"
	"github.com/star-table/startable-server/common/model/vo/projectvo"
)

func GetOneTableColumns(c *gin.Context) {
	req := projectvo.GetTableColumnReq{}
	if err := c.ShouldBindJSON(&req); err != nil {
		// Replaced: c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		c.JSON(http.StatusOK, projectvo.TableColumnsResp{Err: vo.NewErr(errs.BuildSystemErrorInfo(errs.ReqParamsValidateError, err))})
		return
	}

	resp, err := service.GetOneTableColumns(req)
	if err != nil {
		// Replaced: c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()}); return
	}
	c.JSON(http.StatusOK, projectvo.TableColumnsResp{
		Err:  vo.NewErr(err),
		Data: resp,
	})
}

func GetTablesColumns(c *gin.Context) {
	req := projectvo.GetTablesColumnsReq{}
	if err := c.ShouldBindJSON(&req); err != nil {
		// Replaced: c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		c.JSON(http.StatusOK, projectvo.TablesColumnsResp{Err: vo.NewErr(errs.BuildSystemErrorInfo(errs.ReqParamsValidateError, err))})
		return
	}

	resp, err := service.GetTablesColumns(req)
	if err != nil {
		// Replaced: c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()}); return
	}
	c.JSON(http.StatusOK, projectvo.TablesColumnsResp{
		Err:  vo.NewErr(err),
		Data: resp,
	})
}

func CreateTable(c *gin.Context) {
	req := projectvo.CreateTableReq{}
	if err := c.ShouldBindJSON(&req); err != nil {
		// Replaced: c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		c.JSON(http.StatusOK, projectvo.CreateTableRespVo{Err: vo.NewErr(errs.BuildSystemErrorInfo(errs.ReqParamsValidateError, err))})
		return
	}

	res, err := service.CreateTable(req)
	if err != nil {
		// Replaced: c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()}); return
	}
	c.JSON(http.StatusOK, projectvo.CreateTableRespVo{
		Err:  vo.NewErr(err),
		Data: res,
	})
}

func RenameTable(c *gin.Context) {
	req := projectvo.RenameTableReq{}
	if err := c.ShouldBindJSON(&req); err != nil {
		// Replaced: c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		c.JSON(http.StatusOK, projectvo.RenameTableResp{Err: vo.NewErr(errs.BuildSystemErrorInfo(errs.ReqParamsValidateError, err))})
		return
	}

	res, err := service.RenameTable(req)
	if err != nil {
		// Replaced: c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()}); return
	}
	c.JSON(http.StatusOK, projectvo.RenameTableResp{
		Err:  vo.NewErr(err),
		Data: res,
	})
}

func DeleteTable(c *gin.Context) {
	req := projectvo.DeleteTableReq{}
	if err := c.ShouldBindJSON(&req); err != nil {
		// Replaced: c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		c.JSON(http.StatusOK, projectvo.DeleteTableResp{Err: vo.NewErr(errs.BuildSystemErrorInfo(errs.ReqParamsValidateError, err))})
		return
	}

	res, err := service.DeleteTable(req)
	if err != nil {
		// Replaced: c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()}); return
	}
	c.JSON(http.StatusOK, projectvo.DeleteTableResp{
		Err:  vo.NewErr(err),
		Data: res,
	})
}

func SetAutoSchedule(c *gin.Context) {
	req := projectvo.SetAutoScheduleReq{}
	if err := c.ShouldBindJSON(&req); err != nil {
		// Replaced: c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		c.JSON(http.StatusOK, projectvo.SetAutoScheduleResp{Err: vo.NewErr(errs.BuildSystemErrorInfo(errs.ReqParamsValidateError, err))})
		return
	}

	res, err := service.SetAutoSchedule(req)
	if err != nil {
		// Replaced: c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()}); return
	}
	c.JSON(http.StatusOK, projectvo.SetAutoScheduleResp{
		Err:  vo.NewErr(err),
		Data: res,
	})
}

func GetTable(c *gin.Context) {
	req := projectvo.GetTableInfoReq{}
	if err := c.ShouldBindJSON(&req); err != nil {
		// Replaced: c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		c.JSON(http.StatusOK, projectvo.GetTableInfoResp{Err: vo.NewErr(errs.BuildSystemErrorInfo(errs.ReqParamsValidateError, err))})
		return
	}

	res, err := service.GetTable(req)
	if err != nil {
		// Replaced: c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()}); return
	}
	c.JSON(http.StatusOK, projectvo.GetTableInfoResp{
		Err:  vo.NewErr(err),
		Data: res,
	})
}

func GetTables(c *gin.Context) {
	req := projectvo.GetTablesReqVo{}
	if err := c.ShouldBindJSON(&req); err != nil {
		// Replaced: c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		c.JSON(http.StatusOK, projectvo.GetTablesDataResp{Err: vo.NewErr(errs.BuildSystemErrorInfo(errs.ReqParamsValidateError, err))})
		return
	}

	res, err := service.GetTables(req.OrgId, req.UserId, req.Input.AppId)
	if err != nil {
		// Replaced: c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()}); return
	}
	c.JSON(http.StatusOK, projectvo.GetTablesDataResp{
		Err:  vo.NewErr(err),
		Data: res,
	})
}

func GetTablesByApps(c *gin.Context) {
	req := projectvo.ReadTablesByAppsReqVo{}
	if err := c.ShouldBindJSON(&req); err != nil {
		// Replaced: c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		c.JSON(http.StatusOK, projectvo.ReadTablesByAppsRespVo{Err: vo.NewErr(errs.BuildSystemErrorInfo(errs.ReqParamsValidateError, err))})
		return
	}

	res, err := service.GetTablesByApps(req.OrgId, req.UserId, req.Input.AppIds)
	if err != nil {
		// Replaced: c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()}); return
	}
	c.JSON(http.StatusOK, projectvo.ReadTablesByAppsRespVo{
		Err:  vo.NewErr(err),
		Data: res,
	})
}

func GetTablesByOrg(c *gin.Context) {
	req := projectvo.GetTablesByOrgReq{}
	if err := c.ShouldBindJSON(&req); err != nil {
		// Replaced: c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		c.JSON(http.StatusOK, projectvo.GetTablesByOrgRespVo{Err: vo.NewErr(errs.BuildSystemErrorInfo(errs.ReqParamsValidateError, err))})
		return
	}

	res, err := service.GetTablesByOrg(req.OrgId, req.UserId)
	if err != nil {
		// Replaced: c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()}); return
	}
	c.JSON(http.StatusOK, projectvo.GetTablesByOrgRespVo{
		Err:  vo.NewErr(err),
		Data: res,
	})
}

func GetBigTableModeConfig(c *gin.Context) {
	req := projectvo.GetBigTableModeConfigReqVo{}
	if err := c.ShouldBindJSON(&req); err != nil {
		// Replaced: c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		c.JSON(http.StatusOK, projectvo.GetBigTableModeConfigResp{Err: vo.NewErr(errs.BuildSystemErrorInfo(errs.ReqParamsValidateError, err))})
		return
	}

	res, err := service.GetBigTableModeConfig(req)
	if err != nil {
		// Replaced: c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()}); return
	}
	c.JSON(http.StatusOK, projectvo.GetBigTableModeConfigResp{
		Err:  vo.NewErr(err),
		Data: res,
	})
}

func SwitchBigTableMode(c *gin.Context) {
	req := projectvo.SwitchBigTableModeReqVo{}
	if err := c.ShouldBindJSON(&req); err != nil {
		// Replaced: c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		c.JSON(http.StatusOK, vo.CommonRespVo{Err: vo.NewErr(errs.BuildSystemErrorInfo(errs.ReqParamsValidateError, err))})
		return
	}

	err := service.SwitchBigTableMode(req)
	if err != nil {
		// Replaced: c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()}); return
	}
	c.JSON(http.StatusOK, vo.CommonRespVo{
		Err: vo.NewErr(err),
	})
}
