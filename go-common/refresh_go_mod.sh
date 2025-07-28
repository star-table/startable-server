#!/bin/bash
go mod edit -replace=gitea.bjx.cloud/allstar/emitter-go-client=gitea.startable.cn/allstar/emitter-go-client@master
go mod tidy
