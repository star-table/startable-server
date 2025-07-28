package trendssvc

import (
	"github.com/star-table/startable-server/app/facade/idfacade"
	"github.com/star-table/startable-server/common/core/consts"
	"github.com/star-table/startable-server/common/core/errs"
	"github.com/star-table/startable-server/common/core/util/copyer"
	"github.com/star-table/startable-server/common/model/bo"
)

func CreateComment(commentBo bo.CommentBo) (int64, errs.SystemErrorInfo) {
	commentPo := &po.PpmTreComment{}
	_ = copyer.Copy(commentBo, commentPo)

	commentId, err1 := idfacade.ApplyPrimaryIdRelaxed(consts.TableComment)
	if err1 != nil {
		log.Error(err1)
		return 0, errs.BuildSystemErrorInfo(errs.ApplyIdError)
	}

	commentPo.Id = commentId

	err := dao.InsertComment(*commentPo)
	if err != nil {
		log.Error(err)
		return 0, errs.BuildSystemErrorInfo(errs.MysqlOperateError, err)
	}

	return commentId, nil
}
