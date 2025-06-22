#!/bin/bash

# AI智能体后台管理系统启动脚本

echo "正在启动AI智能体后台管理系统..."

# 检查Go是否安装
if ! command -v go &> /dev/null; then
    echo "错误: 未找到Go，请先安装Go 1.23+"
    exit 1
fi

# 检查Go版本
GO_VERSION=$(go version | awk '{print $3}' | sed 's/go//')
REQUIRED_VERSION="1.23"

if [ "$(printf '%s\n' "$REQUIRED_VERSION" "$GO_VERSION" | sort -V | head -n1)" != "$REQUIRED_VERSION" ]; then
    echo "错误: Go版本过低，需要Go 1.23+，当前版本: $GO_VERSION"
    exit 1
fi

echo "Go版本检查通过: $GO_VERSION"

# 检查配置文件
if [ ! -f "config.yaml" ]; then
    echo "错误: 未找到config.yaml文件，请先配置系统参数"
    exit 1
fi

# 检查Redis是否运行
if ! command -v redis-cli &> /dev/null; then
    echo "警告: 未找到Redis客户端，请确保Redis服务正在运行"
else
    if ! redis-cli ping &> /dev/null; then
        echo "警告: Redis服务未运行，请先启动Redis服务"
    else
        echo "Redis服务检查通过"
    fi
fi

# 安装依赖
echo "正在安装依赖..."
go mod tidy

# 创建上传目录
mkdir -p uploads

# 启动服务器
echo "正在启动服务器..."
go run main.go 