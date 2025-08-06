package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/star-table/startable-server/app/service/projectsvc/service/issue_view"
	"github.com/star-table/startable-server/common/core/errs"
	"github.com/star-table/startable-server/common/model/vo"
	"github.com/star-table/startable-server/common/model/vo/projectvo"
)

// / 任务视图相关的操作
// 创建视图
func CreateTaskView(c *gin.Context) {
	var reqVo projectvo.CreateTaskViewReqVo
	if err := c.ShouldBindJSON(&reqVo); err != nil {
		c.JSON(http.StatusBadRequest, projectvo.CreateTaskViewRespVo{Err: vo.NewErr(errs.ReqParamsValidateError), Data: nil})
		return
	}

	res, err := issue_view.CreateTaskView(reqVo.OrgId, reqVo.UserId, &reqVo.Input)
	c.JSON(http.StatusOK, projectvo.CreateTaskViewRespVo{
		Err:  vo.NewErr(err),
		Data: res,
	})
}

func GetTaskViewList(c *gin.Context) {
	var reqVo projectvo.GetTaskViewListReqVo
	if err := c.ShouldBindJSON(&reqVo); err != nil {
		c.JSON(http.StatusBadRequest, projectvo.GetTaskViewListRespVo{Err: vo.NewErr(errs.ReqParamsValidateError), Data: nil})
		return
	}

	res, err := issue_view.GetTaskViewList(reqVo.OrgId, reqVo.UserId, &reqVo.Input)
	c.JSON(http.StatusOK, projectvo.GetTaskViewListRespVo{
		Err:  vo.NewErr(err),
		Data: res,
	})
}

func UpdateTaskView(c *gin.Context) {
	var reqVo projectvo.UpdateTaskViewReqVo
	if err := c.ShouldBindJSON(&reqVo); err != nil {
		c.JSON(http.StatusBadRequest, vo.CommonRespVo{Err: vo.NewErr(errs.ReqParamsValidateError), Void: nil})
		return
	}

	res, err := issue_view.UpdateTaskView(reqVo.OrgId, reqVo.UserId, &reqVo.Input)
	c.JSON(http.StatusOK, vo.CommonRespVo{
		Err:  vo.NewErr(err),
		Void: res,
	})
}

func DeleteTaskView(c *gin.Context) {
	var reqVo projectvo.DeleteTaskViewReqVo
	if err := c.ShouldBindJSON(&reqVo); err != nil {
		c.JSON(http.StatusBadRequest, vo.CommonRespVo{Err: vo.NewErr(errs.ReqParamsValidateError), Void: nil})
		return
	}

	res, err := issue_view.DeleteTaskView(reqVo.OrgId, reqVo.UserId, reqVo.Input.ID)
	c.JSON(http.StatusOK, vo.CommonRespVo{
		Err:  vo.NewErr(err),
		Void: res,
	})
}

// 筛选时，获取值列表
//func (PostGreeter) GetOptionList(reqVo *projectvo.GetOptionListReqVo) projectvo.GetOptionListRespVo {
//	return projectvo.GetOptionListRespVo{}
//
//	//res, err := issue_view.GetOptionList(reqVo.OrgId, reqVo.Input)
//	//return projectvo.GetOptionListRespVo{
//	//	Err:  vo.NewErr(err),
//	//	NewData: res,
//	//}
//}
