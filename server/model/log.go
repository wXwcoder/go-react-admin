package model

import (
	"time"

	"gorm.io/gorm"
)

// Log 日志模型
type Log struct {
	ID          uint           `gorm:"primaryKey" json:"id"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"deleted_at"`
	UserID      uint           `json:"user_id"`
	Username    string         `gorm:"size:50" json:"username"`
	IP          string         `gorm:"size:50" json:"ip"`
	Method      string         `gorm:"size:10" json:"method"`
	Path        string         `gorm:"size:100" json:"path"`
	UserAgent   string         `gorm:"size:255" json:"user_agent"`
	StatusCode  int            `json:"status_code"`
	ResponseTime int           `json:"response_time"` // 响应时间(毫秒)
	TenantID    uint           `gorm:"index" json:"tenant_id"` // 租户ID
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