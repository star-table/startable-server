package util

import (
	"github.com/star-table/startable-server/common/core/errs"
	"github.com/star-table/startable-server/common/core/util/temp"
)

func ParseCacheKey(key string, params map[string]interface{}) (string, errs.SystemErrorInfo) {

	target, err := temp.Render(key, params)
	if err != nil {
		log.Error(err)
		return "", errs.BuildSystemErrorInfo(errs.TemplateRenderError, err)
	}
	return target, nil
}
