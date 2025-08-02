// @title Go-React-Admin API
// @version 1.0
// @description 基于Gin + React的后台管理系统API文档
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:8080
// @BasePath /
// @schemes http https

package main

import (
	"fmt"
	"log"

	"go-react-admin/global"
	"go-react-admin/initialize"
	"go-react-admin/router"

	"github.com/gin-gonic/gin"
	"github.com/gin-contrib/cors"
	_ "go-react-admin/docs" // 引入生成的docs包
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
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
	config.AllowHeaders = []string{"Origin", "Content-Type", "Accept", "Authorization", "Cache-Control", "Pragma", "Expires"}
	r.Use(cors.New(config))

	// 初始化API路由
	router.InitApiRoutes(r)

	// 添加Swagger路由
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// 启动服务器
	port := ":" + global.GlobalConfig.Server.Port
	fmt.Printf("Server is running on port %s\n", port)
	fmt.Printf("Swagger文档地址: http://localhost%s/swagger/index.html\n", port)
	log.Fatal(r.Run(port))
}
