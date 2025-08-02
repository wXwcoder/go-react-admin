package model

import (
	"time"

	"gorm.io/gorm"
)

// Log 日志模型
type Log struct {
	ID           uint           `gorm:"primaryKey" json:"id" example:"1"`
	CreatedAt    time.Time      `json:"created_at" example:"2024-01-01T00:00:00Z"`
	UpdatedAt    time.Time      `json:"updated_at" example:"2024-01-01T00:00:00Z"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty" swaggertype:"string" example:"null"`
	UserID       uint           `json:"user_id" example:"1"`
	Username     string         `gorm:"size:50" json:"username" validate:"max=50" example:"admin"`
	IP           string         `gorm:"size:50" json:"ip" validate:"max=50" example:"192.168.1.1"`
	Method       string         `gorm:"size:10" json:"method" validate:"max=10" example:"GET"`
	Path         string         `gorm:"size:100" json:"path" validate:"max=100" example:"/api/users"`
	UserAgent    string         `gorm:"size:255" json:"user_agent" validate:"max=255" example:"Mozilla/5.0"`
	StatusCode   int            `json:"status_code" example:"200"`
	ResponseTime int            `json:"response_time" example:"150"` // 响应时间(毫秒)
	TenantID     uint           `gorm:"index" json:"tenant_id" example:"1"` // 租户ID
}

// TableName 自定义表名
func (Log) TableName() string {
	return "logs"
}

// BeforeCreate 创建前钩子
func (l *Log) BeforeCreate(tx *gorm.DB) error {
	// 可以在这里添加创建前的逻辑
	return nil
}

// BeforeUpdate 更新前钩子
func (l *Log) BeforeUpdate(tx *gorm.DB) error {
	// 可以在这里添加更新前的逻辑
	return nil
}