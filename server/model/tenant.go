package model

import (
	"time"

	"gorm.io/gorm"
)

// Tenant 租户模型
type Tenant struct {
	ID          uint           `gorm:"primaryKey" json:"id"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"deleted_at"`
	Name        string         `gorm:"uniqueIndex;size:100" json:"name"`
	Code        string         `gorm:"uniqueIndex;size:50" json:"code"`
	Description string         `gorm:"size:255" json:"description"`
	Status      int            `gorm:"default:1" json:"status"` // 1:启用 2:禁用
	AdminUserID uint           `json:"admin_user_id"` // 管理员用户ID
}

// TableName 自定义表名
func (Tenant) TableName() string {
	return "tenants"
}

// BeforeCreate 创建前钩子
func (t *Tenant) BeforeCreate(tx *gorm.DB) error {
	// 可以在这里添加创建前的逻辑
	return nil
}

// BeforeUpdate 更新前钩子
func (t *Tenant) BeforeUpdate(tx *gorm.DB) error {
	// 可以在这里添加更新前的逻辑
	return nil
}