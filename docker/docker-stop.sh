#!/bin/bash
echo "正在停止Go-React-Admin开发环境..."
echo
echo "使用命令: docker-compose -f docker-compose.dev.yml down"
echo
docker-compose -p GAR -f docker-compose.dev.yml down
echo
echo "环境已停止！"
read -p "按任意键继续..." -n1 -s
echo