package initialize

import (
	"fmt"
	"log"
	"path/filepath"

	"github.com/casbin/casbin/v2"
	gormadapter "github.com/casbin/gorm-adapter/v3"
	"go-react-admin/global"
)

// InitCasbin 初始化Casbin
func InitCasbin() {
	// 获取模型文件路径
	modelPath := filepath.Join("config", "rbac_model.conf")
	
	// 创建Gorm适配器
	adapter, err := gormadapter.NewAdapterByDB(global.DB)
	if err != nil {
		log.Fatalf("Failed to create casbin adapter: %v", err)
	}

	// 创建Enforcer
	enforcer, err := casbin.NewEnforcer(modelPath, adapter)
	if err != nil {
		log.Fatalf("Failed to create casbin enforcer: %v", err)
	}

	// 启用自动保存
	enforcer.EnableAutoSave(true)

	// 加载策略
	err = enforcer.LoadPolicy()
	if err != nil {
		log.Fatalf("Failed to load casbin policy: %v", err)
	}

	// 设置全局Enforcer
	global.Enforcer = enforcer

	log.Println("Casbin initialized successfully")
}

// InitDefaultPolicies 初始化默认权限策略
func InitDefaultPolicies() {
	if global.Enforcer == nil {
		log.Println("Casbin enforcer not initialized")
		return
	}

	// 清除现有策略
	global.Enforcer.ClearPolicy()

	// 添加默认角色权限策略
	// 超级管理员角色拥有所有权限
	policies := [][]string{
		// 用户管理权限
		{"admin", "/api/users", "GET", "1"},
		{"admin", "/api/users", "POST", "1"},
		{"admin", "/api/users/*", "PUT", "1"},
		{"admin", "/api/users/*", "DELETE", "1"},
		
		// 角色管理权限
		{"admin", "/api/roles", "GET", "1"},
		{"admin", "/api/roles", "POST", "1"},
		{"admin", "/api/roles/*", "PUT", "1"},
		{"admin", "/api/roles/*", "DELETE", "1"},
		{"admin", "/api/roles/*/permissions", "GET", "1"},
		{"admin", "/api/roles/*/permissions", "PUT", "1"},
		
		// 菜单管理权限
		{"admin", "/api/menus", "GET", "1"},
		{"admin", "/api/menus", "POST", "1"},
		{"admin", "/api/menus/*", "PUT", "1"},
		{"admin", "/api/menus/*", "DELETE", "1"},
		
		// API管理权限
		{"admin", "/api/apis", "GET", "1"},
		{"admin", "/api/apis", "POST", "1"},
		{"admin", "/api/apis/*", "PUT", "1"},
		{"admin", "/api/apis/*", "DELETE", "1"},
		
		// 普通用户权限
		{"user", "/api/profile", "GET", "1"},
		{"user", "/api/profile", "PUT", "1"},
		{"user", "/api/menus/user", "GET", "1"},
	}

	// 批量添加策略
	for _, policy := range policies {
		_, err := global.Enforcer.AddPolicy(policy)
		if err != nil {
			log.Printf("Failed to add policy %v: %v", policy, err)
		}
	}

	// 保存策略
	err := global.Enforcer.SavePolicy()
	if err != nil {
		log.Printf("Failed to save policies: %v", err)
	}

	log.Println("Default policies initialized successfully")
}

// AddUserRole 为用户添加角色
func AddUserRole(userID uint, roleID uint, tenantID uint) error {
	if global.Enforcer == nil {
		return fmt.Errorf("casbin enforcer not initialized")
	}

	userIDStr := fmt.Sprintf("%d", userID)
	roleIDStr := fmt.Sprintf("role_%d", roleID)
	tenantIDStr := fmt.Sprintf("%d", tenantID)

	_, err := global.Enforcer.AddGroupingPolicy(userIDStr, roleIDStr, tenantIDStr)
	return err
}

// RemoveUserRole 移除用户角色
func RemoveUserRole(userID uint, roleID uint, tenantID uint) error {
	if global.Enforcer == nil {
		return fmt.Errorf("casbin enforcer not initialized")
	}

	userIDStr := fmt.Sprintf("%d", userID)
	roleIDStr := fmt.Sprintf("role_%d", roleID)
	tenantIDStr := fmt.Sprintf("%d", tenantID)

	_, err := global.Enforcer.RemoveGroupingPolicy(userIDStr, roleIDStr, tenantIDStr)
	return err
}

// AddRolePermission 为角色添加权限
func AddRolePermission(roleID uint, resource, action string, tenantID uint) error {
	if global.Enforcer == nil {
		return fmt.Errorf("casbin enforcer not initialized")
	}

	roleIDStr := fmt.Sprintf("role_%d", roleID)
	tenantIDStr := fmt.Sprintf("%d", tenantID)

	_, err := global.Enforcer.AddPolicy(roleIDStr, resource, action, tenantIDStr)
	return err
}

// RemoveRolePermission 移除角色权限
func RemoveRolePermission(roleID uint, resource, action string, tenantID uint) error {
	if global.Enforcer == nil {
		return fmt.Errorf("casbin enforcer not initialized")
	}

	roleIDStr := fmt.Sprintf("role_%d", roleID)
	tenantIDStr := fmt.Sprintf("%d", tenantID)

	_, err := global.Enforcer.RemovePolicy(roleIDStr, resource, action, tenantIDStr)
	return err
}