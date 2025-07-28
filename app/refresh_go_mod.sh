#!/bin/bash
echo $(PWD)
go mod edit -replace=gitea.bjx.cloud/LessCode/go-common=gitea.startable.cn/LessCode/go-common@master
go mod tidy
go mod edit -replace=gitea.bjx.cloud/LessCode/interface=gitea.startable.cn/LessCode/interface@master
go mod tidy
go mod edit -replace=gitea.bjx.cloud/allstar/common=gitea.startable.cn/allstar/common@master
go mod tidy
go mod edit -replace=gitea.bjx.cloud/allstar/dingtalk-sdk-golang=gitea.startable.cn/allstar/dingtalk-sdk-golang@master
go mod tidy
go mod edit -replace=gitea.bjx.cloud/allstar/emitter-go-client=gitea.startable.cn/allstar/emitter-go-client@master
go mod tidy
go mod edit -replace=gitea.bjx.cloud/allstar/feishu-sdk-golang=gitea.startable.cn/allstar/feishu-sdk-golang@master
go mod tidy
go mod edit -replace=gitea.bjx.cloud/allstar/platform-sdk=gitea.bjx.cloud/allstar/platform-sdk@master
go mod tidy
go mod edit -replace=github.com/go-laoji/wecom-go-sdk=gitea.startable.cn/LessCode/wecom-go-sdk@main
go mod tidy
