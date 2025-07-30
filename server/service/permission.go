package service

import (
	"errors"
	"fmt"

	"go-react-admin/global"
	"go-react-admin/initialize"
	"go-react-admin/model"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type PermissionService struct{}

// PermissionRequest 权限请求结构
type PermissionRequest struct {
	RoleID   uint   `json:"role_id" binding:"required"`
	MenuIDs  []uint `json:"menu_ids"`
	ApiIDs   []uint `json:"api_ids"`
	TenantID uint   `json:"tenant_id"`
}

// UserRoleRequest 用户角色请求结构
type UserRoleRequest struct {
	UserID   uint   `json:"user_id" binding:"required"`
	RoleIDs  []uint `json:"role_ids"`
	TenantID uint   `json:"tenant_id"`
}

// RolePermissionResponse 角色权限响应结构
type RolePermissionResponse struct {
	RoleID uint         `json:"role_id"`
	Role   *model.Role  `json:"role"`
	Menus  []model.Menu `json:"menus"`
	Apis   []model.Api  `json:"apis"`
}

// UserRoleResponse 用户角色响应结构
type UserRoleResponse struct {
	UserID uint         `json:"user_id"`
	User   *model.User  `json:"user"`
	Roles  []model.Role `json:"roles"`
}

// AssignRolePermissions 分配角色权限
func (s *PermissionService) AssignRolePermissions(req *PermissionRequest) error {
	return global.DB.Transaction(func(tx *gorm.DB) error {
		// 验证角色是否存在
		var role model.Role
		if err := tx.Where("id = ? AND tenant_id = ?", req.RoleID, req.TenantID).First(&role).Error; err != nil {
			return errors.New("角色不存在")
		}

		// 删除现有的角色菜单关联
		if err := tx.Where("role_id = ? AND tenant_id = ?", req.RoleID, req.TenantID).Delete(&model.RoleMenu{}).Error; err != nil {
			return err
		}

		// 删除现有的角色API关联
		if err := tx.Where("role_id = ? AND tenant_id = ?", req.RoleID, req.TenantID).Delete(&model.RoleApi{}).Error; err != nil {
			return err
		}

		// 添加新的菜单权限
		for _, menuID := range req.MenuIDs {
			roleMenu := model.RoleMenu{
				RoleID:   req.RoleID,
				MenuID:   menuID,
				TenantID: req.TenantID,
			}
			if err := tx.Clauses(clause.OnConflict{
				Columns:   []clause.Column{{Name: "role_id"}, {Name: "menu_id"}, {Name: "tenant_id"}},
				UpdateAll: true,
			}).Create(&roleMenu).Error; err != nil {
				return err
			}
		}

		// 添加新的API权限
		for _, apiID := range req.ApiIDs {
			roleApi := model.RoleApi{
				RoleID:   req.RoleID,
				ApiID:    apiID,
				TenantID: req.TenantID,
			}
			if err := tx.Clauses(clause.OnConflict{
				Columns:   []clause.Column{{Name: "role_id"}, {Name: "api_id"}, {Name: "tenant_id"}},
				UpdateAll: true,
			}).Create(&roleApi).Error; err != nil {
				return err
			}
		}

		// 更新Casbin策略
		if err := s.updateCasbinPoliciesForRole(req.RoleID, req.ApiIDs, req.TenantID); err != nil {
			return err
		}

		// 保存Casbin策略
		if global.Enforcer != nil {
			if err := global.Enforcer.SavePolicy(); err != nil {
				return err
			}
		}

		return nil
	})
}

// AssignUserRoles 分配用户角色
func (s *PermissionService) AssignUserRoles(req *UserRoleRequest) error {
	return global.DB.Transaction(func(tx *gorm.DB) error {
		// 验证用户是否存在
		var user model.User
		if err := tx.Where("id = ? AND tenant_id = ?", req.UserID, req.TenantID).First(&user).Error; err != nil {
			return errors.New("用户不存在")
		}

		// 删除现有的用户角色关联
		if err := tx.Where("user_id = ? AND tenant_id = ?", req.UserID, req.TenantID).Delete(&model.UserRole{}).Error; err != nil {
			return err
		}

		// 删除用户在Casbin中的所有角色关联
		if err := s.removeAllUserRoles(req.UserID, req.TenantID); err != nil {
			return err
		}

		// 添加新的用户角色关联
		for _, roleID := range req.RoleIDs {
			userRole := model.UserRole{
				UserID:   req.UserID,
				RoleID:   roleID,
				TenantID: req.TenantID,
			}
			// 使用OnConflict实现存在则更新，不存在则插入
			if err := tx.Clauses(clause.OnConflict{
				Columns:   []clause.Column{{Name: "user_id"}, {Name: "role_id"}, {Name: "tenant_id"}},
				UpdateAll: true,
			}).Create(&userRole).Error; err != nil {
				return err
			}

			// 更新Casbin用户角色关联
			if err := initialize.AddUserRole(req.UserID, roleID, req.TenantID); err != nil {
				return err
			}
		}

		// 保存Casbin策略
		if global.Enforcer != nil {
			if err := global.Enforcer.SavePolicy(); err != nil {
				return err
			}
		}

		return nil
	})
}

// GetRolePermissions 获取角色权限
func (s *PermissionService) GetRolePermissions(roleID, tenantID uint) (*RolePermissionResponse, error) {
	var role model.Role
	if err := global.DB.Where("id = ? AND tenant_id = ?", roleID, tenantID).First(&role).Error; err != nil {
		return nil, errors.New("角色不存在")
	}

	// 获取角色关联的菜单
	var menus []model.Menu
	if err := global.DB.Table("menus").
		Joins("JOIN role_menus ON menus.id = role_menus.menu_id").
		Where("role_menus.role_id = ? AND role_menus.tenant_id = ?", roleID, tenantID).
		Find(&menus).Error; err != nil {
		return nil, err
	}

	// 获取角色关联的API
	var apis []model.Api
	if err := global.DB.Table("apis").
		Joins("JOIN role_apis ON apis.id = role_apis.api_id").
		Where("role_apis.role_id = ? AND role_apis.tenant_id = ?", roleID, tenantID).
		Find(&apis).Error; err != nil {
		return nil, err
	}

	return &RolePermissionResponse{
		RoleID: roleID,
		Role:   &role,
		Menus:  menus,
		Apis:   apis,
	}, nil
}

// GetUserRoles 获取用户角色
func (s *PermissionService) GetUserRoles(userID, tenantID uint) (*UserRoleResponse, error) {
	var user model.User
	if err := global.DB.Where("id = ? AND tenant_id = ?", userID, tenantID).First(&user).Error; err != nil {
		return nil, errors.New("用户不存在")
	}

	// 获取用户关联的角色
	var roles []model.Role
	if err := global.DB.Table("roles").
		Joins("JOIN user_roles ON roles.id = user_roles.role_id").
		Where("user_roles.user_id = ? AND user_roles.tenant_id = ?", userID, tenantID).
		Find(&roles).Error; err != nil {
		return nil, err
	}

	return &UserRoleResponse{
		UserID: userID,
		User:   &user,
		Roles:  roles,
	}, nil
}

// GetUserPermissions 获取用户权限（通过角色）
func (s *PermissionService) GetUserPermissions(userID, tenantID uint) ([]model.Menu, []model.Api, error) {
	// 获取用户角色
	var roleIDs []uint
	if err := global.DB.Model(&model.UserRole{}).
		Where("user_id = ? AND tenant_id = ?", userID, tenantID).
		Pluck("role_id", &roleIDs).Error; err != nil {
		return nil, nil, err
	}

	if len(roleIDs) == 0 {
		return []model.Menu{}, []model.Api{}, nil
	}

	// 获取角色关联的菜单
	var menus []model.Menu
	if err := global.DB.Table("menus").
		Joins("JOIN role_menus ON menus.id = role_menus.menu_id").
		Where("role_menus.role_id IN ? AND role_menus.tenant_id = ?", roleIDs, tenantID).
		Group("menus.id").
		Find(&menus).Error; err != nil {
		return nil, nil, err
	}

	// 获取角色关联的API
	var apis []model.Api
	if err := global.DB.Table("apis").
		Joins("JOIN role_apis ON apis.id = role_apis.api_id").
		Where("role_apis.role_id IN ? AND role_apis.tenant_id = ?", roleIDs, tenantID).
		Group("apis.id").
		Find(&apis).Error; err != nil {
		return nil, nil, err
	}

	return menus, apis, nil
}

// updateCasbinPoliciesForRole 更新角色的Casbin策略
func (s *PermissionService) updateCasbinPoliciesForRole(roleID uint, apiIDs []uint, tenantID uint) error {
	if global.Enforcer == nil {
		return errors.New("casbin enforcer not initialized")
	}

	roleIDStr := fmt.Sprintf("role_%d", roleID)
	tenantIDStr := fmt.Sprintf("%d", tenantID)

	// 删除角色的所有现有策略
	_, err := global.Enforcer.RemoveFilteredPolicy(0, roleIDStr, "", "", tenantIDStr)
	if err != nil {
		return err
	}

	// 根据API权限添加新策略
	var apis []model.Api
	if len(apiIDs) > 0 {
		if err := global.DB.Where("id IN ? AND tenant_id = ?", apiIDs, tenantID).Find(&apis).Error; err != nil {
			return err
		}

		for _, api := range apis {
			_, err := global.Enforcer.AddPolicy(roleIDStr, api.Path, api.Method, tenantIDStr)
			if err != nil {
				return err
			}
		}
	}

	// 不在这里保存策略，由调用方统一保存
	return nil
}

// CheckUserPermission 检查用户是否有特定权限
func (s *PermissionService) CheckUserPermission(userID uint, resource, action string, tenantID uint) (bool, error) {
	if global.Enforcer == nil {
		return false, errors.New("casbin enforcer not initialized")
	}

	// 特殊处理：如果是admin用户（ID为1），直接返回true
	if userID == 1 {
		return true, nil
	}

	// 资源名称到API路径的映射
	resourceToPath := map[string]string{
		"dashboard":  "/api/v1/dashboard",
		"user":       "/api/v1/users",
		"role":       "/api/v1/roles",
		"menu":       "/api/v1/menus",
		"api":        "/api/v1/apis",
		"permission": "/api/v1/permissions",
	}

	// 动作到HTTP方法的映射
	actionToMethod := map[string]string{
		"read":   "GET",
		"create": "POST",
		"update": "PUT",
		"delete": "DELETE",
	}

	// 获取对应的API路径和HTTP方法
	apiPath, exists := resourceToPath[resource]
	if !exists {
		// 如果没有映射，直接使用资源名称
		apiPath = resource
	}

	httpMethod, exists := actionToMethod[action]
	if !exists {
		// 如果没有映射，直接使用动作名称
		httpMethod = action
	}

	userIDStr := fmt.Sprintf("%d", userID)
	tenantIDStr := fmt.Sprintf("%d", tenantID)

	// 使用Casbin的RBAC模型检查用户权限
	// 在RBAC模型中，权限是赋予角色的，用户通过分组策略关联到角色
	allowed, err := global.Enforcer.Enforce(userIDStr, apiPath, httpMethod, tenantIDStr)
	if err != nil {
		fmt.Printf("权限检查错误: userID=%d, path=%s, method=%s, tenantID=%d, error: %v", userID, apiPath, httpMethod, tenantID, err)
		return false, err
	}

	// 如果映射路径检查失败，尝试直接使用原始资源名称检查
	if !allowed {
		fmt.Printf("权限检查失败，尝试原始资源: userID=%d, resource=%s, action=%s, tenantID=%d", userID, resource, action, tenantID)
		allowed, err = global.Enforcer.Enforce(userIDStr, resource, action, tenantIDStr)
		if err != nil {
			fmt.Printf("原始资源权限检查错误: userID=%d, resource=%s, action=%s, tenantID=%d, error: %v", userID, resource, action, tenantID, err)
			return false, err
		}
	}

	fmt.Printf("权限检查结果: userID=%d, resource=%s, action=%s, tenantID=%d, allowed=%v", userID, resource, action, tenantID, allowed)
	return allowed, nil
}

// removeAllUserRoles 删除用户的所有角色关联
func (s *PermissionService) removeAllUserRoles(userID, tenantID uint) error {
	if global.Enforcer == nil {
		return errors.New("casbin enforcer not initialized")
	}

	userIDStr := fmt.Sprintf("%d", userID)
	tenantIDStr := fmt.Sprintf("%d", tenantID)

	// 删除用户的所有分组策略（角色关联）
	_, err := global.Enforcer.RemoveFilteredGroupingPolicy(0, userIDStr, "", tenantIDStr)
	return err
}
