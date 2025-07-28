package domain

import (
	"github.com/star-table/startable-server/app/service/projectsvc/po"
	"github.com/star-table/startable-server/common/core/consts"
	"github.com/star-table/startable-server/common/core/util/copyer"
	"github.com/star-table/startable-server/common/library/db/mysql"
	"github.com/star-table/startable-server/common/model/bo"
	"upper.io/db.v3"
)

func GetIssueReport(id string) (bo.ShareBo, error) {
	shareInfo := &po.PpmShaShare{}
	err := mysql.SelectOneByCond(consts.TableShare, db.Cond{
		consts.TcIsDelete: consts.AppIsNoDelete,
		consts.TcStatus:   1,
		consts.TcId:       id,
	}, shareInfo)
	shareBo := &bo.ShareBo{}
	if err != nil {
		return *shareBo, err
	}
	_ = copyer.Copy(shareInfo, shareBo)

	return *shareBo, nil
}
