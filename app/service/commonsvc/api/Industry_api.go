package api

import (
	"net/http"
	"github.com/gin-gonic/gin"
	commonService "github.com/star-table/startable-server/app/service/commonsvc/service"
	"github.com/star-table/startable-server/common/model/vo"
	"github.com/star-table/startable-server/common/model/vo/commonvo"
)

// IndustryList 获取行业列表
func IndustryList(c *gin.Context) {
	res, err := commonService.IndustryList()
	response := commonvo.IndustryListRespVo{Err: vo.NewErr(err), IndustryList: res}
	c.JSON(http.StatusOK, response)
}
