#!/bin/bash
echo "正在启动Go-React-Admin开发环境..."
echo
echo "使用命令: docker-compose -f docker-compose.dev.yml up --build"
echo
docker-compose -p GAR -f docker-compose.dev.yml up --build
read -p "按任意键继续..." -n1 -s
echo