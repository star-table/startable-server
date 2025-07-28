package domain

import (
	"context"
	"fmt"
	"os"

	"github.com/star-table/startable-server/common/core/config"
	"github.com/star-table/startable-server/common/core/consts"
	"github.com/gin-gonic/gin"
)

var env = ""

const BaseConfigPath = "./../../../../config"
const SelfConfigPath = "./../config"

func init() {
	env = os.Getenv("POL_ENV")
	if "" == env {
		env = "unittest"
	}
	//配置文件
	err := config.LoadEnvConfig(BaseConfigPath, "application.common", env)

	if err != nil {
		fmt.Printf("err:%s\n", err)
	}

	err = config.LoadEnvConfig(SelfConfigPath, "application", "")

	if err != nil {
		fmt.Printf("err:%s\n", err)
	}
	//添加token操作
	ginCtx := gin.Context{}
	ginCtx.Set(consts.AppHeaderTokenName, "abc")
	//获得一个顶级上下文
	ctx := context.Background()
	//返回父上下文  在父上下文中设置 key/value  这边是 GinContextKey:ginCtx
	ctx = context.WithValue(ctx, "GinContextKey", &ginCtx)
}
