package model

import (
	"time"

	"gorm.io/gorm"
)

// Menu 菜单模型
type Menu struct {
	ID          uint           `gorm:"primaryKey" json:"id"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"deleted_at"`
	Name        string         `gorm:"size:50" json:"name"`
	Path        string         `gorm:"size:100" json:"path"`
	Component   string         `gorm:"size:100" json:"component"`
	Icon        string         `gorm:"size:50" json:"icon"`
	Title       string         `gorm:"size:50" json:"title"`
	ParentID    uint           `gorm:"default:0" json:"parent_id"`
	Sort        int            `gorm:"default:0" json:"sort"`
	Hidden      bool           `gorm:"default:false" json:"hidden"`
	Level       int            `gorm:"default:1" json:"level"`
	Status      int            `gorm:"default:1" json:"status"` // 1:启用 2:禁用
	TenantID    uint           `gorm:"index" json:"tenant_id"` // 租户ID
}

// TableName 自定义表名
func (Menu) TableName() string {
	return "menus"
}

// BeforeCreate 创建前钩子
func (m *Menu) BeforeCreate(tx *gorm.DB) error {
	// 可以在这里添加创建前的逻辑
	return nil
}

// BeforeUpdate 更新前钩子
func (m *Menu) BeforeUpdate(tx *gorm.DB) error {
	// 可以在这里添加更新前的逻辑
	return nil
}