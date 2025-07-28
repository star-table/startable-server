package util

import (
	"strings"

	"github.com/star-table/startable-server/common/core/util/file"
)

func GetCurrentPath() string {
	return file.GetCurrentPath()
}

func GetMobile(mobile string) string {
	strs := strings.Split(mobile, "-")
	if len(strs) == 1 {
		return strs[0]
	}
	return strs[1]
}
