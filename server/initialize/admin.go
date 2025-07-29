package initialize

import (
	"fmt"
	"log"

	"go-react-admin/global"
	"go-react-admin/model"
)

// InitAdminUser 初始化管理员用户和超级管理员角色
func InitAdminUser() {
	// 默认租户ID为1
	defaultTenantID := uint(1)

	// 检查是否已存在管理员用户
	var count int64
	global.DB.Model(&model.User{}).Where("username = ? AND tenant_id = ?", "admin", defaultTenantID).Count(&count)

	var adminUser model.User
	if count == 0 {
		// 创建管理员用户
		adminUser = model.User{
			Username: "admin",
			Password: "123456", // 实际项目中应该加密密码
			Nickname: "管理员",
			Email:    "admin@example.com",
			Phone:    "13800138000",
			Status:   1,
			TenantID: defaultTenantID,
		}

		if err := global.DB.Create(&adminUser).Error; err != nil {
			log.Printf("创建管理员用户失败: %v", err)
			return
		} else {
			log.Println("管理员用户创建成功")
		}
	} else {
		// 获取已存在的管理员用户
		global.DB.Where("username = ? AND tenant_id = ?", "admin", defaultTenantID).First(&adminUser)
		log.Println("管理员用户已存在")
	}

	// 创建或获取超级管理员角色
	var superAdminRole model.Role
	err := global.DB.Where("name = ? AND tenant_id = ?", "超级管理员", defaultTenantID).First(&superAdminRole).Error
	if err != nil {
		// 创建超级管理员角色
		superAdminRole = model.Role{
			Name:        "超级管理员",
			Description: "拥有系统所有权限的超级管理员角色",
			Status:      1,
			TenantID:    defaultTenantID,
		}
		if err := global.DB.Create(&superAdminRole).Error; err != nil {
			log.Printf("创建超级管理员角色失败: %v", err)
			return
		}
		log.Println("超级管理员角色创建成功")
	} else {
		log.Println("超级管理员角色已存在")
	}

	// 获取所有菜单ID
	var allMenus []model.Menu
	if err := global.DB.Where("tenant_id = ?", defaultTenantID).Find(&allMenus).Error; err != nil {
		log.Printf("获取菜单列表失败: %v", err)
		return
	}

	// 获取所有API ID
	var allApis []model.Api
	if err := global.DB.Where("tenant_id = ?", defaultTenantID).Find(&allApis).Error; err != nil {
		log.Printf("获取API列表失败: %v", err)
		return
	}

	// 清除现有的角色权限关联
	global.DB.Where("role_id = ? AND tenant_id = ?", superAdminRole.ID, defaultTenantID).Delete(&model.RoleMenu{})
	global.DB.Where("role_id = ? AND tenant_id = ?", superAdminRole.ID, defaultTenantID).Delete(&model.RoleApi{})

	// 给超级管理员角色分配所有菜单权限
	for _, menu := range allMenus {
		roleMenu := model.RoleMenu{
			RoleID:   superAdminRole.ID,
			MenuID:   menu.ID,
			TenantID: defaultTenantID,
		}
		global.DB.Create(&roleMenu)
	}

	// 给超级管理员角色分配所有API权限
	for _, api := range allApis {
		roleApi := model.RoleApi{
			RoleID:   superAdminRole.ID,
			ApiID:    api.ID,
			TenantID: defaultTenantID,
		}
		global.DB.Create(&roleApi)
	}

	// 清除现有的用户角色关联
	global.DB.Where("user_id = ? AND tenant_id = ?", adminUser.ID, defaultTenantID).Delete(&model.UserRole{})

	// 将管理员用户关联到超级管理员角色
	userRole := model.UserRole{
		UserID:   adminUser.ID,
		RoleID:   superAdminRole.ID,
		TenantID: defaultTenantID,
	}
	if err := global.DB.Create(&userRole).Error; err != nil {
		log.Printf("关联用户角色失败: %v", err)
	} else {
		log.Println("管理员用户关联超级管理员角色成功")
	}

	// 更新Casbin策略，给超级管理员角色添加所有权限
	if global.Enforcer != nil {
		roleIDStr := fmt.Sprintf("role_%d", superAdminRole.ID)
		tenantIDStr := fmt.Sprintf("%d", defaultTenantID)
		
		// 清除现有的角色策略
		global.Enforcer.RemoveFilteredPolicy(0, roleIDStr, "", "", tenantIDStr)
		
		// 添加所有API权限到超级管理员角色
		for _, api := range allApis {
			policy := []string{roleIDStr, api.Path, api.Method, tenantIDStr}
			global.Enforcer.AddPolicy(policy)
		}
		
		// 添加用户角色关系到Casbin
		userIDStr := fmt.Sprintf("%d", adminUser.ID)
		global.Enforcer.AddGroupingPolicy(userIDStr, roleIDStr, tenantIDStr)
		
		// 保存策略
		global.Enforcer.SavePolicy()
		log.Println("超级管理员Casbin策略更新成功")
	}

	log.Printf("管理员用户权限初始化完成")
}
