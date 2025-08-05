@echo off
echo 正在停止Go-React-Admin开发环境...
echo.
echo 使用命令: docker-compose -f docker-compose.dev.yml down
echo.
docker-compose -p gar -f docker-compose.dev.yml down
echo.
echo 环境已停止！
pause