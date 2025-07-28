package model

import (
	"time"

	"gorm.io/gorm"
)

// Api API模型
type Api struct {
	ID          uint           `gorm:"primaryKey" json:"id"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"deleted_at"`
	Path        string         `gorm:"size:100" json:"path"`
	Method      string         `gorm:"size:10" json:"method"` // GET, POST, PUT, DELETE
	Category    string         `gorm:"size:50" json:"category"`
	Description string         `gorm:"size:255" json:"description"`
	Status      int            `gorm:"default:1" json:"status"` // 1:启用 2:禁用
	TenantID    uint           `gorm:"index" json:"tenant_id"` // 租户ID
}

// TableName 自定义表名
func (Api) TableName() string {
	return "apis"
}

// BeforeCreate 创建前钩子
func (a *Api) BeforeCreate(tx *gorm.DB) error {
	// 可以在这里添加创建前的逻辑
	return nil
}

// BeforeUpdate 更新前钩子
func (a *Api) BeforeUpdate(tx *gorm.DB) error {
	// 可以在这里添加更新前的逻辑
	return nil
}