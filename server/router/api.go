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
		public.POST("/login", middleware.LoginLogger(), api.Login)
		public.POST("/register", api.Register)
	}

	// 受保护的路由
	protected := r.Group("/api/v1")
	protected.Use(middleware.Logger()) // 添加操作日志记录
	protected.Use(middleware.JWTAuth())
	{
		// 用户相关路由
		protected.GET("/user/info", api.GetUserInfo)
		protected.POST("/user/logout", middleware.LogoutLogger(), api.Logout)
		protected.GET("/user/list", api.GetUserList)
		protected.POST("/user/create", api.CreateUser)
		protected.PUT("/user/update/:id", api.UpdateUser)
		protected.DELETE("/user/delete/:id", api.DeleteUser)
		protected.POST("/user/upload-avatar", api.UploadAvatar)

		// 角色相关路由
		protected.GET("/role/list", api.GetRoleList)
		protected.POST("/role/create", api.CreateRole)
		protected.PUT("/role/update/:id", api.UpdateRole)
		protected.DELETE("/role/delete/:id", api.DeleteRole)

		// 菜单相关路由
		protected.GET("/menu/list", api.GetMenuList)
		protected.GET("/menus/user", api.GetUserMenus)
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

		// 第三方客户管理路由
		protected.GET("/admin/customers", api.CustomerApi.GetCustomerList)
		protected.POST("/admin/customers", api.CustomerApi.CreateCustomer)
		protected.GET("/admin/customers/:id", api.CustomerApi.GetCustomer)
		protected.PUT("/admin/customers/:id", api.CustomerApi.UpdateCustomer)
		protected.DELETE("/admin/customers/:id", api.CustomerApi.DeleteCustomer)
		protected.PUT("/admin/customers/:id/status", api.CustomerApi.UpdateCustomerStatus)
		protected.POST("/admin/customers/:id/ban", api.CustomerApi.BanCustomer)
		protected.POST("/admin/customers/:id/unban", api.CustomerApi.UnbanCustomer)

		// 管理员消息管理路由
		adminMessageAPI := &api.AdminMessageAPI{}
		protected.GET("/admin/messages", adminMessageAPI.GetMessages)
		protected.GET("/admin/messages/stats", adminMessageAPI.GetMessageStats)
		protected.GET("/admin/messages/:id", adminMessageAPI.GetMessageDetail)
		protected.POST("/admin/messages", adminMessageAPI.CreateMessage)
		protected.PUT("/admin/messages/:id", adminMessageAPI.UpdateMessage)
		protected.DELETE("/admin/messages/:id", adminMessageAPI.DeleteMessage)
		protected.POST("/admin/messages/:id/send", adminMessageAPI.SendMessage)
		protected.POST("/admin/messages/:id/cancel", adminMessageAPI.CancelMessage)
		
		// 管理员公告管理路由
		protected.GET("/admin/announcements", adminMessageAPI.GetAnnouncementList)
		protected.GET("/admin/announcements/stats", adminMessageAPI.GetAnnouncementStats)
		protected.GET("/admin/announcements/:id", adminMessageAPI.GetAnnouncementDetail)
		protected.POST("/admin/announcements", adminMessageAPI.CreateAnnouncement)
		protected.PUT("/admin/announcements/:id", adminMessageAPI.UpdateAnnouncement)
		protected.DELETE("/admin/announcements/:id", adminMessageAPI.DeleteAnnouncement)
		protected.POST("/admin/announcements/:id/publish", adminMessageAPI.PublishAnnouncement)
		protected.POST("/admin/announcements/:id/revoke", adminMessageAPI.RevokeAnnouncement)

		// 第三方客户消息管理路由
		protected.POST("/admin/customer-messages", api.CustomerMessageApi.AdminCreateMessage)
		protected.POST("/admin/customer-messages/batch", api.CustomerMessageApi.AdminCreateMessagesBatch)

		// 动态数据管理路由
		InitDynamicRoutes(protected)
	}

	// 第三方客户认证路由
	customerAuth := r.Group("/api/v1/customer")
	{
		customerAuth.POST("/register", api.CustomerAuthApi.Register)
		customerAuth.POST("/login", api.CustomerAuthApi.Login)
		customerAuth.POST("/reset-password", api.CustomerAuthApi.ResetPassword)
	}

	// 第三方客户受保护路由
	customerProtected := r.Group("/api/v1/customer")
	customerProtected.Use(middleware.CustomerJWTAuth())
	{
		customerProtected.GET("/profile", api.CustomerAuthApi.GetProfile)
		customerProtected.PUT("/profile", api.CustomerAuthApi.UpdateProfile)
		customerProtected.PUT("/password", api.CustomerAuthApi.UpdatePassword)
		customerProtected.POST("/logout", api.CustomerAuthApi.Logout)
		customerProtected.POST("/refresh-token", api.CustomerAuthApi.RefreshToken)

		// 第三方客户消息路由
		customerProtected.GET("/messages", api.CustomerMessageApi.GetMessages)
		customerProtected.GET("/messages/:id", api.CustomerMessageApi.GetMessageDetail)
		customerProtected.PUT("/messages/:id/read", api.CustomerMessageApi.MarkMessageAsRead)
		customerProtected.POST("/messages/:id/read", api.CustomerMessageApi.MarkMessageAsRead)
		customerProtected.PUT("/messages/batch-read", api.CustomerMessageApi.MarkMessagesAsReadBatch)
		customerProtected.GET("/messages/unread-count", api.CustomerMessageApi.GetUnreadCount)
		customerProtected.DELETE("/messages/:id", api.CustomerMessageApi.DeleteMessage)

		// 第三方客户公告路由
		customerProtected.GET("/announcements", api.CustomerMessageApi.GetCustomerAnnouncements)
		customerProtected.GET("/announcements/:id", api.CustomerMessageApi.GetCustomerAnnouncementDetail)
		customerProtected.POST("/announcements/:id/read", api.CustomerMessageApi.MarkAnnouncementRead)
		customerProtected.POST("/announcements/batch-read", api.CustomerMessageApi.MarkAnnouncementsBatchRead)
		customerProtected.GET("/announcements/unread-count", api.CustomerMessageApi.GetUnreadAnnouncementCount)
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
