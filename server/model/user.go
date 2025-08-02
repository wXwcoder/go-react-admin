package model

import (
	"time"

	"gorm.io/gorm"
)

// User 用户模型
type User struct {
	ID        uint           `gorm:"primaryKey" json:"id" example:"1"`
	CreatedAt time.Time      `json:"created_at" example:"2024-01-01T00:00:00Z"`
	UpdatedAt time.Time      `json:"updated_at" example:"2024-01-01T00:00:00Z"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty" swaggertype:"string" example:"null"`
	Username  string         `gorm:"uniqueIndex;size:50" json:"username" validate:"required,min=3,max=50" example:"admin"`
	Password  string         `gorm:"size:100" json:"password" validate:"required,min=6,max=100" example:"123456"`
	Nickname  string         `gorm:"size:50" json:"nickname" validate:"max=50" example:"管理员"`
	Email     string         `gorm:"size:100" json:"email" validate:"email,max=100" example:"admin@example.com"`
	Phone     string         `gorm:"size:20" json:"phone" validate:"max=20" example:"13800138000"`
	Status    int            `gorm:"default:1" json:"status" validate:"oneof=1 2" example:"1"` // 1:启用 2:禁用
	Avatar    string         `gorm:"size:255" json:"avatar" validate:"max=255" example:"https://example.com/avatar.jpg"`
	TenantID  uint           `gorm:"index" json:"tenant_id" example:"1"` // 租户ID
}

// TableName 自定义表名
func (User) TableName() string {
	return "users"
}

// BeforeCreate 创建前钩子
func (u *User) BeforeCreate(tx *gorm.DB) error {
	// 可以在这里添加创建前的逻辑
	return nil
}

// BeforeUpdate 更新前钩子
func (u *User) BeforeUpdate(tx *gorm.DB) error {
	// 可以在这里添加更新前的逻辑
	return nil
}