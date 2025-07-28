package model

import (
	"time"

	"gorm.io/gorm"
)

// RoleMenu 角色菜单关联模型
type RoleMenu struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at"`
	RoleID    uint           `gorm:"index" json:"role_id"`
	MenuID    uint           `gorm:"index" json:"menu_id"`
	TenantID  uint           `gorm:"index" json:"tenant_id"` // 租户ID
}

// TableName 自定义表名
func (RoleMenu) TableName() string {
	return "role_menus"
}

// BeforeCreate 创建前钩子
func (rm *RoleMenu) BeforeCreate(tx *gorm.DB) error {
	// 可以在这里添加创建前的逻辑
	return nil
}

// BeforeUpdate 更新前钩子
func (rm *RoleMenu) BeforeUpdate(tx *gorm.DB) error {
	// 可以在这里添加更新前的逻辑
	return nil
}