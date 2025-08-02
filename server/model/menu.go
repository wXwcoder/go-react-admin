package model

import (
	"time"

	"gorm.io/gorm"
)

// Menu 菜单模型
type Menu struct {
	ID          uint           `gorm:"primaryKey" json:"id" example:"1"`
	CreatedAt   time.Time      `json:"created_at" example:"2024-01-01T00:00:00Z"`
	UpdatedAt   time.Time      `json:"updated_at" example:"2024-01-01T00:00:00Z"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty" swaggertype:"string" example:"null"`
	Name        string         `gorm:"size:50" json:"name" validate:"required,min=2,max=50" example:"用户管理"`
	Path        string         `gorm:"size:100" json:"path" validate:"max=100" example:"/user"`
	Component   string         `gorm:"size:100" json:"component" validate:"max=100" example:"User"`
	Icon        string         `gorm:"size:50" json:"icon" validate:"max=50" example:"user"`
	Title       string         `gorm:"size:50" json:"title" validate:"required,max=50" example:"用户管理"`
	ParentID    uint           `gorm:"default:0" json:"parent_id" example:"0"`
	Sort        int            `gorm:"default:0" json:"sort" example:"1"`
	Hidden      bool           `gorm:"default:false" json:"hidden" example:"false"`
	Level       int            `gorm:"default:1" json:"level" example:"1"`
	Type        string         `gorm:"size:20;default:menu" json:"type" validate:"oneof=menu group" example:"menu"` // menu:菜单项 group:菜单组
	Status      int            `gorm:"default:1" json:"status" validate:"oneof=1 2" example:"1"` // 1:启用 2:禁用
	TenantID    uint           `gorm:"index" json:"tenant_id" example:"1"` // 租户ID
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