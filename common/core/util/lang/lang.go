package lang

import (
	"github.com/star-table/startable-server/common/core/consts"
	"github.com/star-table/startable-server/common/core/threadlocal"
)

func GetLang() string {
	return threadlocal.GetValue(consts.AppHeaderLanguage)
}

func IsEnglish() bool {
	return GetLang() == consts.LangEnglish
}
