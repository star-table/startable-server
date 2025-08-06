package domain

import (
	"github.com/star-table/startable-server/app/service/projectsvc/po"
	"github.com/star-table/startable-server/common/core/consts"
	"github.com/star-table/startable-server/common/core/errs"
	"github.com/star-table/startable-server/common/core/logger"

	"github.com/star-table/startable-server/common/core/util/slice"
	"github.com/star-table/startable-server/common/library/db/mysql"

	"strconv"

	"upper.io/db.v3"
)

var log = logger.GetDefaultLogger()

func CheckCodeRepetition(orgId int64, code string) (bool, errs.SystemErrorInfo) {
	count, err := mysql.SelectCountByCond(consts.TableProject, db.Cond{
		consts.TcOrgId:    orgId,
		consts.TcPreCode:  code,
		consts.TcIsDelete: consts.AppIsNoDelete,
	})
	if err != nil {
		log.Error(err)
		return false, errs.BuildSystemErrorInfo(errs.MysqlOperateError, err)
	}
	if count > 0 {
		return true, nil
	}
	return false, nil
}

func GetUniquelyCode(orgId int64, code string) (string, errs.SystemErrorInfo) {
	projects := &[]po.PpmProProject{}
	err := mysql.SelectAllByCond(consts.TableProject, db.Cond{
		consts.TcOrgId:    orgId,
		consts.TcPreCode:  db.Like(code + "%"),
		consts.TcIsDelete: consts.AppIsNoDelete,
	}, projects)
	if err != nil {
		log.Error(err)
		return "", errs.BuildSystemErrorInfo(errs.MysqlOperateError, err)
	}
	alreadyExistedCodes := make([]string, len(*projects))
	for i, pro := range *projects {
		alreadyExistedCodes[i] = pro.PreCode
	}
	for i := 0; i < len(alreadyExistedCodes)+1; i++ {
		c := code + strconv.Itoa(i)
		con, err := slice.Contain(alreadyExistedCodes, c)
		if err != nil {
			log.Error(err)
			return "", errs.BuildSystemErrorInfo(errs.ObjectTypeError, err)
		}
		if con {
			continue
		}
		return c, nil
	}
	return "", errs.BuildSystemErrorInfo(errs.ProjectDomainError)
}
