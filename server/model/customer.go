package model

import (
	"time"
)

type CustomerStatus string

const (
	CustomerStatusActive  CustomerStatus = "active"
	CustomerStatusBanned  CustomerStatus = "banned"
	CustomerStatusPending CustomerStatus = "pending"
)

type Customer struct {
	ID           uint64         `gorm:"primaryKey;autoIncrement" json:"id"`
	Username     string         `gorm:"type:varchar(50);uniqueIndex;not null" json:"username"`
	Email        string         `gorm:"type:varchar(100);uniqueIndex;not null" json:"email"`
	PasswordHash string         `gorm:"type:varchar(255);not null" json:"-"`
	Phone        *string        `gorm:"type:varchar(20)" json:"phone,omitempty"`
	RealName     *string        `gorm:"type:varchar(50)" json:"real_name,omitempty"`
	AvatarURL    *string        `gorm:"type:varchar(500)" json:"avatar_url,omitempty"`
	Status       CustomerStatus `gorm:"type:enum('active','banned','pending');default:'pending';index" json:"status"`
	LastLoginAt  *time.Time     `gorm:"type:datetime" json:"last_login_at,omitempty"`
	LoginCount   int            `gorm:"type:int;default:0" json:"login_count"`
	CreatedAt    time.Time      `gorm:"type:datetime;default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt    time.Time      `gorm:"type:datetime;default:CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP" json:"updated_at"`
	DeletedAt    *time.Time     `gorm:"type:datetime;index" json:"deleted_at,omitempty"`

	// 关联关系
}

func (Customer) TableName() string {
	return "customers"
}

// BeforeCreate 钩子函数
func (c *Customer) BeforeCreate() error {
	if c.Status == "" {
		c.Status = CustomerStatusPending
	}
	return nil
}