#!/bin/bash

# 数据库初始化脚本

echo "正在初始化数据库..."

# 检查Go是否安装
if ! command -v go &> /dev/null; then
    echo "错误: 未找到Go，请先安装Go 1.23+"
    exit 1
fi

# 检查配置文件
if [ ! -f "config.yaml" ]; then
    echo "错误: 未找到config.yaml文件，请先配置数据库连接信息"
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

# 运行数据库初始化脚本
echo "正在创建数据库表..."
go run scripts/init_db.go

echo "数据库初始化完成！"
echo "默认用户信息："
echo "  用户名: admin"
echo "  密码: admin123"
echo "  邮箱: admin@example.com" 