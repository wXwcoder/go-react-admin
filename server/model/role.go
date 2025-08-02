package model

import (
	"time"

	"gorm.io/gorm"
)

// Role 角色模型
type Role struct {
	ID          uint           `gorm:"primaryKey" json:"id" example:"1"`
	CreatedAt   time.Time      `json:"created_at" example:"2024-01-01T00:00:00Z"`
	UpdatedAt   time.Time      `json:"updated_at" example:"2024-01-01T00:00:00Z"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty" swaggertype:"string" example:"null"`
	Name        string         `gorm:"uniqueIndex;size:50" json:"name" validate:"required,min=2,max=50" example:"管理员"`
	Description string         `gorm:"size:255" json:"description" validate:"max=255" example:"系统管理员角色"`
	Status      int            `gorm:"default:1" json:"status" validate:"oneof=1 2" example:"1"` // 1:启用 2:禁用
	TenantID    uint           `gorm:"index" json:"tenant_id" example:"1"` // 租户ID
}

// TableName 自定义表名
func (Role) TableName() string {
	return "roles"
}

// BeforeCreate 创建前钩子
func (r *Role) BeforeCreate(tx *gorm.DB) error {
	// 可以在这里添加创建前的逻辑
	return nil
}

// BeforeUpdate 更新前钩子
func (r *Role) BeforeUpdate(tx *gorm.DB) error {
	// 可以在这里添加更新前的逻辑
	return nil
}