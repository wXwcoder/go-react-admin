package initialize

import (
	"go-react-admin/global"
	"go-react-admin/model"
	"log"
)

// InitMenuData 初始化菜单数据
func InitMenuData() {
	// 检查是否已经初始化过菜单
	var count int64
	global.DB.Model(&model.Menu{}).Count(&count)
	if count > 0 {
		log.Println("菜单数据已存在，跳过初始化")
		return
	}

	// 创建默认菜单结构
	menus := []model.Menu{
		// 首页
		{
			ID:       1,
			Name:     "dashboard",
			Path:     "/dashboard",
			Component: "Dashboard",
			Icon:     "fas fa-tachometer-alt",
			Title:    "仪表板",
			ParentID: 0,
			Sort:     1,
			Level:    1,
			Type:     "menu",
			Status:   1,
		},
		// 系统管理（一级菜单组）
		{
			ID:       2,
			Name:     "system",
			Path:     "/system",
			Component: "",
			Icon:     "fas fa-cogs",
			Title:    "系统管理",
			ParentID: 0,
			Sort:     2,
			Level:    1,
			Type:     "group",
			Status:   1,
		},
		// 用户管理
		{
			ID:       3,
			Name:     "users",
			Path:     "/system/users",
			Component: "UserManagement",
			Icon:     "fas fa-users",
			Title:    "用户管理",
			ParentID: 2,
			Sort:     1,
			Level:    2,
			Type:     "menu",
			Status:   1,
		},
		// 角色管理
		{
			ID:       4,
			Name:     "roles",
			Path:     "/system/roles",
			Component: "RoleManagement",
			Icon:     "fas fa-user-tag",
			Title:    "角色管理",
			ParentID: 2,
			Sort:     2,
			Level:    2,
			Type:     "menu",
			Status:   1,
		},
		// 菜单管理
		{
			ID:       5,
			Name:     "menus",
			Path:     "/system/menus",
			Component: "MenuManagement",
			Icon:     "fas fa-list",
			Title:    "菜单管理",
			ParentID: 2,
			Sort:     3,
			Level:    2,
			Type:     "menu",
			Status:   1,
		},
		// API管理
		{
			ID:       6,
			Name:     "apis",
			Path:     "/system/apis",
			Component: "ApiManagement",
			Icon:     "fas fa-code",
			Title:    "API管理",
			ParentID: 2,
			Sort:     4,
			Level:    2,
			Type:     "menu",
			Status:   1,
		},
		// 权限管理
		{
			ID:       7,
			Name:     "permissions",
			Path:     "/system/permissions",
			Component: "PermissionManagement",
			Icon:     "fas fa-shield-alt",
			Title:    "权限管理",
			ParentID: 2,
			Sort:     5,
			Level:    2,
			Type:     "menu",
			Status:   1,
		},
		// 日志管理
		{
			ID:       8,
			Name:     "logs",
			Path:     "/system/logs",
			Component: "LogManagement",
			Icon:     "fas fa-history",
			Title:    "日志管理",
			ParentID: 2,
			Sort:     6,
			Level:    2,
			Type:     "menu",
			Status:   1,
		},
		// 动态数据管理（一级菜单组）
		{
			ID:       9,
			Name:     "dynamic",
			Path:     "/dynamic",
			Component: "",
			Icon:     "fas fa-database",
			Title:    "动态数据管理",
			ParentID: 0,
			Sort:     3,
			Level:    1,
			Type:     "group",
			Status:   1,
		},
		// 动态表管理
		{
			ID:       10,
			Name:     "dynamic-tables",
			Path:     "/dynamic/tables",
			Component: "DynamicTableManagement",
			Icon:     "fas fa-table",
			Title:    "动态表管理",
			ParentID: 9,
			Sort:     1,
			Level:    2,
			Type:     "menu",
			Status:   1,
		},
	}

	// 批量创建菜单
	for _, menu := range menus {
		if err := global.DB.Create(&menu).Error; err != nil {
			log.Printf("创建菜单失败: %v", err)
		}
	}

	log.Println("菜单数据初始化完成")
}