package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	commonService "github.com/star-table/startable-server/app/service/commonsvc/service"
	"github.com/star-table/startable-server/common/model/vo"
	"github.com/star-table/startable-server/common/model/vo/commonvo"
)

// UploadOssByFsImageKey 通过飞书图片key上传到OSS
func UploadOssByFsImageKey(c *gin.Context) {
	var req commonvo.UploadOssByFsImageKeyReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response := commonvo.UploadOssByFsImageKeyResp{Err: vo.NewErr(err)}
		c.JSON(http.StatusOK, response)
		return
	}

	res, err := commonService.UploadOssByFsImageKey(req.OrgId, req.ImageKey, req.IsApp)
	response := commonvo.UploadOssByFsImageKeyResp{
		Err: vo.NewErr(err),
		Url: res,
	}
	c.JSON(http.StatusOK, response)
}
