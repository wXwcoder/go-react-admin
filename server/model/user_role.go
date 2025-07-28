package model

import (
	"time"

	"gorm.io/gorm"
)

// UserRole 用户角色关联模型
type UserRole struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at"`
	UserID    uint           `gorm:"index" json:"user_id"`
	RoleID    uint           `gorm:"index" json:"role_id"`
	TenantID  uint           `gorm:"index" json:"tenant_id"` // 租户ID
}

// TableName 自定义表名
func (UserRole) TableName() string {
	return "user_roles"
}

// BeforeCreate 创建前钩子
func (ur *UserRole) BeforeCreate(tx *gorm.DB) error {
	// 可以在这里添加创建前的逻辑
	return nil
}

// BeforeUpdate 更新前钩子
func (ur *UserRole) BeforeUpdate(tx *gorm.DB) error {
	// 可以在这里添加更新前的逻辑
	return nil
}