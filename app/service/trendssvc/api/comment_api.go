package trendssvc

import (
	"github.com/star-table/startable-server/app/service"
	"github.com/star-table/startable-server/common/model/vo"
	"github.com/star-table/startable-server/common/model/vo/trendsvo"
)

func (PostGreeter) CreateComment(req trendsvo.CreateCommentReqVo) trendsvo.CreateCommentRespVo {
	id, err := service.CreateComment(req.CommentBo)
	return trendsvo.CreateCommentRespVo{CommentId: id, Err: vo.NewErr(err)}
}
