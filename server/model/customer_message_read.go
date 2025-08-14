package model

import (
	"time"
)

type CustomerMessageRead struct {
	ID         uint64     `gorm:"primaryKey;autoIncrement" json:"id"`
	CustomerID uint64     `gorm:"not null;index:idx_customer_message,priority:1;index:idx_customer,priority:1" json:"customer_id"`
	MessageID  uint64     `gorm:"not null;index:idx_customer_message,priority:2;index:idx_message,priority:1" json:"message_id"`
	IsRead     bool       `gorm:"type:tinyint(1);default:false" json:"is_read"`
	ReadTime   *time.Time `gorm:"type:datetime" json:"read_time,omitempty"`
	CreatedAt  time.Time  `gorm:"type:datetime;default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt  time.Time  `gorm:"type:datetime;default:CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP" json:"updated_at"`
}

func (CustomerMessageRead) TableName() string {
	return "customer_message_reads"
}