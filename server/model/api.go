package model

import (
	"time"

	"gorm.io/gorm"
)

// Api API模型
type Api struct {
	ID          uint           `gorm:"primaryKey" json:"id" example:"1"`
	CreatedAt   time.Time      `json:"created_at" example:"2024-01-01T00:00:00Z"`
	UpdatedAt   time.Time      `json:"updated_at" example:"2024-01-01T00:00:00Z"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty" swaggertype:"string" example:"null"`
	Path        string         `gorm:"size:100" json:"path" validate:"required,max=100" example:"/api/users"`
	Method      string         `gorm:"size:10" json:"method" validate:"required,oneof=GET POST PUT DELETE" example:"GET"` // GET, POST, PUT, DELETE
	Category    string         `gorm:"size:50" json:"category" validate:"required,max=50" example:"用户管理"`
	Description string         `gorm:"size:255" json:"description" validate:"max=255" example:"获取用户列表接口"`
	Status      int            `gorm:"default:1" json:"status" validate:"oneof=1 2" example:"1"` // 1:启用 2:禁用
	TenantID    uint           `gorm:"index" json:"tenant_id" example:"1"` // 租户ID
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