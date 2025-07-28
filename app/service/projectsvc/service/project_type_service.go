package service

import (
	"github.com/star-table/startable-server/app/service/projectsvc/domain"
	"github.com/star-table/startable-server/common/core/errs"
	"github.com/star-table/startable-server/common/core/util/copyer"
	"github.com/star-table/startable-server/common/core/util/lang"
	"github.com/star-table/startable-server/common/language/english"
	"github.com/star-table/startable-server/common/model/vo"
)

func ProjectTypes(orgId int64) ([]*vo.ProjectType, errs.SystemErrorInfo) {
	projectTypeBo, err := domain.GetProjectTypeList(orgId)
	if err != nil {
		return nil, err
	}

	projectTypeVo := &[]*vo.ProjectType{}
	copyErr := copyer.Copy(projectTypeBo, projectTypeVo)
	if copyErr != nil {
		return nil, errs.BuildSystemErrorInfo(errs.JSONConvertError, copyErr)
	}
	if lang.IsEnglish() {
		for i, projectType := range *projectTypeVo {
			if projectType.LangCode == "" {
				continue
			}
			if enName, ok := english.ProjectTypeLang[projectType.LangCode]; ok {
				(*projectTypeVo)[i].Name = enName
			}
		}
	}

	return *projectTypeVo, nil
}
