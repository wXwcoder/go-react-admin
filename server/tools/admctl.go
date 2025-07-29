package main

import (
	"fmt"
	"go-react-admin/global"
	"go-react-admin/initialize"
	"go-react-admin/model"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"log"
	"time"
)

func main() {
	// 加载配置
	initialize.LoadConfig()
	
	// 初始化数据库和Redis
	initialize.InitDB()
	initialize.InitRedis()
	
	// 初始化Casbin
	initialize.InitCasbin()
	
	// 开始测试：为admin超级管理员账号授权
	setupAdminTestData()
	
	log.Println("Admin超级管理员账号授权测试完成")
}

func setupAdminTestData() {
	// 1. 创建超级管理员用户
	adminUser := model.User{
		Username:  "admin",
		Password:  hashPassword("admin123"),
		Email:     "admin@example.com",
		Status:    1,
		TenantID:  1, // 设置默认租户ID
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	
	// 检查用户是否已存在
	var existingUser model.User
	if err := global.DB.Where("username = ?", "admin").First(&existingUser).Error; err == nil {
		log.Printf("管理员用户已存在，ID: %d", existingUser.ID)
		adminUser.ID = existingUser.ID
	} else if err == gorm.ErrRecordNotFound {
		// 创建新用户
		if err := global.DB.Create(&adminUser).Error; err != nil {
			log.Fatalf("创建管理员用户失败: %v", err)
		}
		log.Printf("创建管理员用户成功，ID: %d", adminUser.ID)
	} else {
		log.Fatalf("查询管理员用户失败: %v", err)
	}
	
	// 2. 创建超级管理员角色
	superAdminRole := model.Role{
		Name:        "超级管理员",
		Description: "系统超级管理员角色，拥有所有权限",
		Status:      1,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
	
	var existingRole model.Role
	if err := global.DB.Where("name = ?", "超级管理员").First(&existingRole).Error; err == nil {
		log.Printf("超级管理员角色已存在，ID: %d", existingRole.ID)
		superAdminRole.ID = existingRole.ID
	} else if err == gorm.ErrRecordNotFound {
		if err := global.DB.Create(&superAdminRole).Error; err != nil {
			log.Fatalf("创建超级管理员角色失败: %v", err)
		}
		log.Printf("创建超级管理员角色成功，ID: %d", superAdminRole.ID)
	} else {
		log.Fatalf("查询超级管理员角色失败: %v", err)
	}
	
	// 3. 创建测试菜单和API权限
	createTestMenusAndAPIs()
	
	// 4. 获取所有菜单和API
	var allMenus []model.Menu
	var allApis []model.Api
	global.DB.Find(&allMenus)
	global.DB.Find(&allApis)
	
	log.Printf("共找到 %d 个菜单和 %d 个API权限", len(allMenus), len(allApis))
	
	// 5. 分配所有权限给超级管理员角色
	assignAllPermissionsToRole(superAdminRole.ID, allMenus, allApis)
	
	// 6. 建立用户与角色的关联
	userRole := model.UserRole{
		UserID:    adminUser.ID,
		RoleID:    superAdminRole.ID,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	
	var existingUserRole model.UserRole
	if err := global.DB.Where("user_id = ? AND role_id = ?", adminUser.ID, superAdminRole.ID).First(&existingUserRole).Error; err == nil {
		log.Printf("用户角色关联已存在")
	} else if err == gorm.ErrRecordNotFound {
		if err := global.DB.Create(&userRole).Error; err != nil {
			log.Fatalf("创建用户角色关联失败: %v", err)
		}
		log.Printf("用户角色关联创建成功")
	}
	
	// 7. 更新Casbin策略
	roleName := "admin" // 使用固定的角色名称
	updateCasbinPolicies(adminUser.ID, roleName, allMenus, allApis)

	log.Printf("管理员授权完成：用户ID=%d, 角色=%s", adminUser.ID, roleName)
}

func createTestMenusAndAPIs() {
	// 创建测试菜单
	menus := []model.Menu{
		{Name: "系统管理", Path: "/system", Component: "System", Icon: "setting", Sort: 1, Status: 1},
		{Name: "用户管理", Path: "/user", Component: "UserManagement", Icon: "user", Sort: 1, Status: 1},
		{Name: "角色管理", Path: "/role", Component: "RoleManagement", Icon: "team", Sort: 2, Status: 1},
		{Name: "菜单管理", Path: "/menu", Component: "MenuManagement", Icon: "menu", Sort: 3, Status: 1},
		{Name: "API管理", Path: "/api", Component: "ApiManagement", Icon: "api", Sort: 4, Status: 1},
		{Name: "权限管理", Path: "/permission", Component: "PermissionManagement", Icon: "lock", Sort: 5, Status: 1},
	}
	
	for _, menu := range menus {
		var existing model.Menu
		if err := global.DB.Where("name = ?", menu.Name).First(&existing).Error; err == gorm.ErrRecordNotFound {
			menu.CreatedAt = time.Now()
			menu.UpdatedAt = time.Now()
			if err := global.DB.Create(&menu).Error; err != nil {
				log.Printf("创建菜单失败: %v", err)
			}
		}
	}
	
	// 创建测试API
	apis := []model.Api{
		{Path: "/api/v1/users", Method: "GET", Description: "获取用户列表", Category: "用户管理"},
		{Path: "/api/v1/users", Method: "POST", Description: "创建用户", Category: "用户管理"},
		{Path: "/api/v1/users/:id", Method: "PUT", Description: "更新用户", Category: "用户管理"},
		{Path: "/api/v1/users/:id", Method: "DELETE", Description: "删除用户", Category: "用户管理"},
		{Path: "/api/v1/roles", Method: "GET", Description: "获取角色列表", Category: "角色管理"},
		{Path: "/api/v1/roles", Method: "POST", Description: "创建角色", Category: "角色管理"},
		{Path: "/api/v1/roles/:id", Method: "PUT", Description: "更新角色", Category: "角色管理"},
		{Path: "/api/v1/roles/:id", Method: "DELETE", Description: "删除角色", Category: "角色管理"},
		{Path: "/api/v1/menus", Method: "GET", Description: "获取菜单列表", Category: "菜单管理"},
		{Path: "/api/v1/menus", Method: "POST", Description: "创建菜单", Category: "菜单管理"},
		{Path: "/api/v1/menus/:id", Method: "PUT", Description: "更新菜单", Category: "菜单管理"},
		{Path: "/api/v1/menus/:id", Method: "DELETE", Description: "删除菜单", Category: "菜单管理"},
		{Path: "/api/v1/apis", Method: "GET", Description: "获取API列表", Category: "API管理"},
		{Path: "/api/v1/apis", Method: "POST", Description: "创建API", Category: "API管理"},
		{Path: "/api/v1/apis/:id", Method: "PUT", Description: "更新API", Category: "API管理"},
		{Path: "/api/v1/apis/:id", Method: "DELETE", Description: "删除API", Category: "API管理"},
		{Path: "/api/v1/permissions/user/:id/permissions", Method: "GET", Description: "获取用户权限", Category: "权限管理"},
		{Path: "/api/v1/permissions/role/:id/permissions", Method: "GET", Description: "获取角色权限", Category: "权限管理"},
		{Path: "/api/v1/permissions/role/:id/permissions", Method: "POST", Description: "分配角色权限", Category: "权限管理"},
	}
	
	for _, api := range apis {
		var existing model.Api
		if err := global.DB.Where("path = ? AND method = ?", api.Path, api.Method).First(&existing).Error; err == gorm.ErrRecordNotFound {
			api.CreatedAt = time.Now()
			api.UpdatedAt = time.Now()
			if err := global.DB.Create(&api).Error; err != nil {
				log.Printf("创建API失败: %v", err)
			}
		}
	}
}

func assignAllPermissionsToRole(roleID uint, menus []model.Menu, apis []model.Api) {
	// 分配所有菜单权限
	for _, menu := range menus {
		roleMenu := model.RoleMenu{
				RoleID:    roleID,
				MenuID:    menu.ID,
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			}
		
		var existing model.RoleMenu
		if err := global.DB.Where("role_id = ? AND menu_id = ?", roleID, menu.ID).First(&existing).Error; err == gorm.ErrRecordNotFound {
			if err := global.DB.Create(&roleMenu).Error; err != nil {
				log.Printf("分配菜单权限失败: %v", err)
			}
		}
	}
	
	// 分配所有API权限
	for _, api := range apis {
		roleApi := model.RoleApi{
				RoleID:    roleID,
				ApiID:     api.ID,
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			}
		
		var existing model.RoleApi
		if err := global.DB.Where("role_id = ? AND api_id = ?", roleID, api.ID).First(&existing).Error; err == gorm.ErrRecordNotFound {
			if err := global.DB.Create(&roleApi).Error; err != nil {
				log.Printf("分配API权限失败: %v", err)
			}
		}
	}
	
	log.Printf("已为角色ID=%d分配所有权限", roleID)
}

func updateCasbinPolicies(userID uint, roleCode string, menus []model.Menu, apis []model.Api) {
	if global.Enforcer == nil {
		log.Println("Casbin enforcer未初始化")
		return
	}
	
	// 清除现有的策略
	global.Enforcer.RemoveFilteredPolicy(0, roleCode)
	
	// 添加所有API权限到Casbin策略
	for _, api := range apis {
		policy := []string{roleCode, api.Path, api.Method, "1"}
		global.Enforcer.AddPolicy(policy)
	}
	
	// 保存策略
	global.Enforcer.SavePolicy()
	log.Printf("Casbin策略更新完成，角色=%s", roleCode)
	
	// 添加用户角色关系到Casbin
	userIDStr := fmt.Sprintf("%d", userID)
	global.Enforcer.AddGroupingPolicy(userIDStr, "admin", "1")
	log.Printf("用户角色关系已添加到Casbin: 用户ID=%s, 角色=%s, 租户ID=1", userIDStr, "admin")
}

func hashPassword(password string) string {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Fatalf("密码加密失败: %v", err)
	}
	return string(hashedPassword)
}