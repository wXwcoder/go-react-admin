package router

import (
	"go-react-admin/api"
	"go-react-admin/middleware"

	"github.com/gin-gonic/gin"
)

// InitApiRoutes 初始化API路由
func InitApiRoutes(r *gin.Engine) {
	// 初始化API实例
	permissionApi := api.NewPermissionApi()

	// 公开路由
	public := r.Group("/api/v1")
	{
		// 用户相关路由
		public.POST("/login", api.Login)
		public.POST("/register", api.Register)
	}

	// 受保护的路由
	protected := r.Group("/api/v1")
	protected.Use(middleware.JWTAuth())
	{
		// 用户相关路由
		protected.GET("/user/info", api.GetUserInfo)
		protected.POST("/user/logout", api.Logout)
		protected.GET("/user/list", api.GetUserList)
		protected.POST("/user/create", api.CreateUser)
		protected.PUT("/user/update/:id", api.UpdateUser)
		protected.DELETE("/user/delete/:id", api.DeleteUser)

		// 角色相关路由
		protected.GET("/role/list", api.GetRoleList)
		protected.POST("/role/create", api.CreateRole)
		protected.PUT("/role/update/:id", api.UpdateRole)
		protected.DELETE("/role/delete/:id", api.DeleteRole)

		// 菜单相关路由
		protected.GET("/menu/list", api.GetMenuList)
		protected.POST("/menu/create", api.CreateMenu)
		protected.PUT("/menu/update/:id", api.UpdateMenu)
		protected.DELETE("/menu/delete/:id", api.DeleteMenu)

		// API相关路由
		protected.GET("/api/list", api.GetApiList)
		protected.POST("/api/create", api.CreateApi)
		protected.PUT("/api/update/:id", api.UpdateApi)
		protected.DELETE("/api/delete/:id", api.DeleteApi)

		// 权限管理路由
		permissions := protected.Group("/permissions")
		{
			// 角色权限管理
			permissions.POST("/role", permissionApi.AssignRolePermissions)
			permissions.GET("/role/:id", permissionApi.GetRolePermissions)
			
			// 用户角色管理
			permissions.POST("/user", permissionApi.AssignUserRoles)
			permissions.GET("/user/:id/roles", permissionApi.GetUserRoles)
			permissions.GET("/user/:id/permissions", permissionApi.GetUserPermissions)
			
			// 权限检查
			permissions.GET("/check", permissionApi.CheckPermission)
		}

		// 日志相关路由
		protected.GET("/log/list", api.GetLogList)

		// 租户相关路由
		protected.GET("/tenant/list", api.GetTenantList)
	}

	// 需要权限验证的路由
	authorized := r.Group("/api/v1")
	authorized.Use(middleware.JWTAuth())
	authorized.Use(middleware.CasbinMiddleware())
	{
		// 这里可以添加需要特定权限的路由
		// 例如：只有管理员才能访问的路由
		authorized.GET("/admin/users", api.GetUserList)
		authorized.GET("/admin/roles", api.GetRoleList)
		authorized.GET("/admin/menus", api.GetMenuList)
		authorized.GET("/admin/apis", api.GetApiList)
	}
}
