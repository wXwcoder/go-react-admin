package model

import (
	"time"

	"gorm.io/gorm"
)

// RoleApi 角色API关联模型
type RoleApi struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at"`
	RoleID    uint           `gorm:"index;uniqueIndex:idx_role_api_tenant" json:"role_id"`
	ApiID     uint           `gorm:"index;uniqueIndex:idx_role_api_tenant" json:"api_id"`
	TenantID  uint           `gorm:"index;uniqueIndex:idx_role_api_tenant" json:"tenant_id"` // 租户ID
}

// TableName 自定义表名
func (RoleApi) TableName() string {
	return "role_apis"
}

// BeforeCreate 创建前钩子
func (ra *RoleApi) BeforeCreate(tx *gorm.DB) error {
	// 可以在这里添加创建前的逻辑
	return nil
}

// BeforeUpdate 更新前钩子
func (ra *RoleApi) BeforeUpdate(tx *gorm.DB) error {
	// 可以在这里添加更新前的逻辑
	return nil
}