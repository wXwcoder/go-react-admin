package main

import (
	"fmt"
	"log"

	"go-react-admin/global"
	"go-react-admin/initialize"
	"go-react-admin/router"

	"github.com/gin-gonic/gin"
	"github.com/gin-contrib/cors"
)

func main() {
	// 加载配置
	initialize.LoadConfig()

	// 初始化数据库
	initialize.InitDB()

	// 数据库迁移
	initialize.Migrate()

	// 初始化Redis
	initialize.InitRedis()

	// 初始化Casbin权限管理
	initialize.InitCasbin()

	// 初始化默认权限策略
	initialize.InitDefaultPolicies()

	// 初始化管理员用户
	initialize.InitAdminUser()

	// 初始化菜单数据
	initialize.InitMenuData()

	// 初始化API数据
	initialize.InitApiData()

	// 创建Gin路由器
	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()

	// 添加CORS中间件
	config := cors.DefaultConfig()
	config.AllowAllOrigins = true
	config.AllowMethods = []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}
	config.AllowHeaders = []string{"Origin", "Content-Type", "Accept", "Authorization"}
	r.Use(cors.New(config))

	// 初始化API路由
	router.InitApiRoutes(r)

	// 启动服务器
	port := ":" + global.GlobalConfig.Server.Port
	fmt.Printf("Server is running on port %s\n", port)
	log.Fatal(r.Run(port))
}
