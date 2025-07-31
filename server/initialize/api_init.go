package initialize

import (
	"go-react-admin/global"
	"go-react-admin/model"
	"log"
)

// InitApiData 初始化API数据
func InitApiData() {
	// 检查是否已经初始化过API数据
	var count int64
	global.DB.Model(&model.Api{}).Count(&count)
	if count > 0 {
		log.Println("API数据已存在，跳过初始化")
		return
	}

	// 创建默认API权限数据
	apis := []model.Api{
		// 用户相关API
		{ID: 1, Path: "/api/v1/user/list", Method: "GET", Description: "获取用户列表", Category: "用户管理"},
		{ID: 2, Path: "/api/v1/user/create", Method: "POST", Description: "创建用户", Category: "用户管理"},
		{ID: 3, Path: "/api/v1/user/update/:id", Method: "PUT", Description: "更新用户", Category: "用户管理"},
		{ID: 4, Path: "/api/v1/user/delete/:id", Method: "DELETE", Description: "删除用户", Category: "用户管理"},
		{ID: 5, Path: "/api/v1/user/info", Method: "GET", Description: "获取用户信息", Category: "用户管理"},
		{ID: 6, Path: "/api/v1/user/logout", Method: "POST", Description: "用户登出", Category: "用户管理"},

		// 角色相关API
		{ID: 7, Path: "/api/v1/role/list", Method: "GET", Description: "获取角色列表", Category: "角色管理"},
		{ID: 8, Path: "/api/v1/role/create", Method: "POST", Description: "创建角色", Category: "角色管理"},
		{ID: 9, Path: "/api/v1/role/update/:id", Method: "PUT", Description: "更新角色", Category: "角色管理"},
		{ID: 10, Path: "/api/v1/role/delete/:id", Method: "DELETE", Description: "删除角色", Category: "角色管理"},

		// 菜单相关API
		{ID: 11, Path: "/api/v1/menu/list", Method: "GET", Description: "获取菜单列表", Category: "菜单管理"},
		{ID: 12, Path: "/api/v1/menus/user", Method: "GET", Description: "获取用户菜单", Category: "菜单管理"},
		{ID: 13, Path: "/api/v1/menu/create", Method: "POST", Description: "创建菜单", Category: "菜单管理"},
		{ID: 14, Path: "/api/v1/menu/update/:id", Method: "PUT", Description: "更新菜单", Category: "菜单管理"},
		{ID: 15, Path: "/api/v1/menu/delete/:id", Method: "DELETE", Description: "删除菜单", Category: "菜单管理"},

		// API管理相关API
		{ID: 16, Path: "/api/v1/api/list", Method: "GET", Description: "获取API列表", Category: "API管理"},
		{ID: 17, Path: "/api/v1/api/create", Method: "POST", Description: "创建API", Category: "API管理"},
		{ID: 18, Path: "/api/v1/api/update/:id", Method: "PUT", Description: "更新API", Category: "API管理"},
		{ID: 19, Path: "/api/v1/api/delete/:id", Method: "DELETE", Description: "删除API", Category: "API管理"},

		// 权限管理相关API
		{ID: 20, Path: "/api/v1/permissions/role", Method: "POST", Description: "分配角色权限", Category: "权限管理"},
		{ID: 21, Path: "/api/v1/permissions/role/:id", Method: "GET", Description: "获取角色权限", Category: "权限管理"},
		{ID: 22, Path: "/api/v1/permissions/user", Method: "POST", Description: "分配用户角色", Category: "权限管理"},
		{ID: 23, Path: "/api/v1/permissions/user/:id/roles", Method: "GET", Description: "获取用户角色", Category: "权限管理"},
		{ID: 24, Path: "/api/v1/permissions/user/:id/permissions", Method: "GET", Description: "获取用户权限", Category: "权限管理"},
		{ID: 25, Path: "/api/v1/permissions/check", Method: "GET", Description: "检查权限", Category: "权限管理"},

		// 日志管理相关API
		{ID: 26, Path: "/api/v1/log/list", Method: "GET", Description: "获取日志列表", Category: "日志管理"},

		// 租户相关API
		{ID: 27, Path: "/api/v1/tenant/list", Method: "GET", Description: "获取租户列表", Category: "租户管理"},

		// 认证相关API
		{ID: 28, Path: "/api/v1/login", Method: "POST", Description: "用户登录", Category: "认证管理"},
		{ID: 29, Path: "/api/v1/register", Method: "POST", Description: "用户注册", Category: "认证管理"},
	}

	// 批量创建API数据
	for _, api := range apis {
		if err := global.DB.Create(&api).Error; err != nil {
			log.Printf("创建API数据失败: %v", err)
		} else {
			log.Printf("创建API: %s %s", api.Method, api.Path)
		}
	}

	log.Println("API数据初始化完成")
}