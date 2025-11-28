#!/bin/bash

# 星穹通行证 - Swagger API文档生成脚本

echo "=== 生成Swagger API文档 ==="

# 检查swag是否安装
if ! command -v swag &> /dev/null
then
    echo "swag未安装，正在安装..."
    go install github.com/swaggo/swag/cmd/swag@latest
fi

# 生成文档
echo "正在生成API文档..."
swag init -g main.go -o ./docs

if [ $? -eq 0 ]; then
    echo "✅ API文档生成成功！"
    echo "文档位置: ./docs/swagger.json"
    echo "访问地址: http://localhost:8080/swagger/index.html"
else
    echo "❌ API文档生成失败"
    exit 1
fi
