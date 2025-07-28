package trendsvo

import (
	"github.com/star-table/startable-server/common/model/bo"
	"github.com/star-table/startable-server/common/model/vo"
)

type CreateCommentReqVo struct {
	CommentBo bo.CommentBo `json:"commentBo"`
}

type CreateCommentRespVo struct {
	CommentId int64 `json:"data"`

	vo.Err
}
