package model

import (
	"encoding/json"
	"time"

	"gorm.io/gorm"
)

// TablePermission 表权限控制
type TablePermission struct {
	ID          uint           `gorm:"primaryKey" json:"id"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
	
	TableID     uint           `gorm:"index" json:"table_id" validate:"required"`
	RoleID      uint           `gorm:"index" json:"role_id" validate:"required"`
	CanView     bool           `gorm:"default:false" json:"can_view"`
	CanCreate   bool           `gorm:"default:false" json:"can_create"`
	CanUpdate   bool           `gorm:"default:false" json:"can_update"`
	CanDelete   bool           `gorm:"default:false" json:"can_delete"`
	CanExport   bool           `gorm:"default:false" json:"can_export"`
	FieldPermissions json.RawMessage `gorm:"type:json" json:"field_permissions"`
}

// TableName 自定义表名
func (TablePermission) TableName() string {
	return "table_permissions"
}

// FieldPermission 字段级权限
type FieldPermission struct {
	FieldName string `json:"field_name"`
	CanView   bool   `json:"can_view"`
	CanEdit   bool   `json:"can_edit"`
}

// PermissionRequest 权限设置请求
type PermissionRequest struct {
	TableID          uint                        `json:"table_id" validate:"required"`
	RoleID           uint                        `json:"role_id" validate:"required"`
	CanView          bool                        `json:"can_view"`
	CanCreate        bool                        `json:"can_create"`
	CanUpdate        bool                        `json:"can_update"`
	CanDelete        bool                        `json:"can_delete"`
	CanExport        bool                        `json:"can_export"`
	FieldPermissions map[string]FieldPermission `json:"field_permissions"`
}

// PermissionCheck 权限检查请求
type PermissionCheck struct {
	TableID uint   `json:"table_id" validate:"required"`
	Action  string `json:"action" validate:"required,oneof=view create update delete export"`
	Field   string `json:"field,omitempty"`
}

// UserPermission 用户权限汇总
type UserPermission struct {
	TableID          uint                        `json:"table_id"`
	CanView          bool                        `json:"can_view"`
	CanCreate        bool                        `json:"can_create"`
	CanUpdate        bool                        `json:"can_update"`
	CanDelete        bool                        `json:"can_delete"`
	CanExport        bool                        `json:"can_export"`
	FieldPermissions map[string]FieldPermission `json:"field_permissions"`
}

// GetFieldPermissions 获取字段权限
func (tp *TablePermission) GetFieldPermissions() (map[string]FieldPermission, error) {
	if tp.FieldPermissions == nil {
		return make(map[string]FieldPermission), nil
	}
	
	var permissions map[string]FieldPermission
	err := json.Unmarshal(tp.FieldPermissions, &permissions)
	if err != nil {
		return nil, err
	}
	
	return permissions, nil
}

// SetFieldPermissions 设置字段权限
func (tp *TablePermission) SetFieldPermissions(permissions map[string]FieldPermission) error {
	data, err := json.Marshal(permissions)
	if err != nil {
		return err
	}
	tp.FieldPermissions = data
	return nil
}

// HasPermission 检查是否有指定权限
func (tp *TablePermission) HasPermission(action string) bool {
	switch action {
	case "view":
		return tp.CanView
	case "create":
		return tp.CanCreate
	case "update":
		return tp.CanUpdate
	case "delete":
		return tp.CanDelete
	case "export":
		return tp.CanExport
	default:
		return false
	}
}

// HasFieldPermission 检查字段权限
func (tp *TablePermission) HasFieldPermission(fieldName, action string) bool {
	permissions, err := tp.GetFieldPermissions()
	if err != nil {
		return false
	}
	
	fieldPerm, exists := permissions[fieldName]
	if !exists {
		// 如果没有设置字段权限，默认继承表权限
		switch action {
		case "view":
			return tp.CanView
		case "edit":
			return tp.CanUpdate
		default:
			return false
		}
	}
	
	switch action {
	case "view":
		return fieldPerm.CanView
	case "edit":
		return fieldPerm.CanEdit
	default:
		return false
	}
}

// MergePermissions 合并多个角色的权限（取最大权限）
func MergePermissions(permissions []TablePermission) *UserPermission {
	if len(permissions) == 0 {
		return &UserPermission{
			FieldPermissions: make(map[string]FieldPermission),
		}
	}
	
	result := &UserPermission{
		TableID:          permissions[0].TableID,
		FieldPermissions: make(map[string]FieldPermission),
	}
	
	// 合并表级权限（任一角色有权限即可）
	for _, perm := range permissions {
		if perm.CanView {
			result.CanView = true
		}
		if perm.CanCreate {
			result.CanCreate = true
		}
		if perm.CanUpdate {
			result.CanUpdate = true
		}
		if perm.CanDelete {
			result.CanDelete = true
		}
		if perm.CanExport {
			result.CanExport = true
		}
		
		// 合并字段权限
		fieldPerms, err := perm.GetFieldPermissions()
		if err != nil {
			continue
		}
		
		for fieldName, fieldPerm := range fieldPerms {
			existing, exists := result.FieldPermissions[fieldName]
			if !exists {
				result.FieldPermissions[fieldName] = fieldPerm
			} else {
				// 取最大权限
				result.FieldPermissions[fieldName] = FieldPermission{
					FieldName: fieldName,
					CanView:   existing.CanView || fieldPerm.CanView,
					CanEdit:   existing.CanEdit || fieldPerm.CanEdit,
				}
			}
		}
	}
	
	return result
}

// DefaultPermissions 获取默认权限配置
func DefaultPermissions() *PermissionRequest {
	return &PermissionRequest{
		CanView:          true,
		CanCreate:        false,
		CanUpdate:        false,
		CanDelete:        false,
		CanExport:        false,
		FieldPermissions: make(map[string]FieldPermission),
	}
}

// AdminPermissions 获取管理员权限配置
func AdminPermissions() *PermissionRequest {
	return &PermissionRequest{
		CanView:          true,
		CanCreate:        true,
		CanUpdate:        true,
		CanDelete:        true,
		CanExport:        true,
		FieldPermissions: make(map[string]FieldPermission),
	}
}

// ReadOnlyPermissions 获取只读权限配置
func ReadOnlyPermissions() *PermissionRequest {
	return &PermissionRequest{
		CanView:          true,
		CanCreate:        false,
		CanUpdate:        false,
		CanDelete:        false,
		CanExport:        true,
		FieldPermissions: make(map[string]FieldPermission),
	}
}