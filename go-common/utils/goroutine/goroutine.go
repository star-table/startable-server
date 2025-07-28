package goroutine

import (
	"gitea.bjx.cloud/LessCode/go-common/utils/stack"
	"github.com/go-kratos/kratos/v2/log"
)

func SafeRun(fn func(), logger *log.Helper) {
	go func() {
		defer func() {
			if r := recover(); r != nil {
				logger.Errorf("err:%v, stack:%v", r, stack.GetStack())
			}
		}()
		fn()
	}()
}
