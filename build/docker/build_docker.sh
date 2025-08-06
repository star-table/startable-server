#!/bin/bash

set -x

#取当前目录
BASE_PATH=`cd "$(dirname "$0")"; pwd`

cd $BASE_PATH

mkdir -p /data/logs/startable-server



# 停止服务，并删除镜像
# docker-compose down --rmi all

# 构建镜像
docker-compose build

# 推送镜像到仓库
docker-compose push

# 启动服务
#docker-compose up -d


