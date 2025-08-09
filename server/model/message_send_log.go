package model

import (
	"time"
)

type SendStatus string

const (
	SendStatusPending SendStatus = "pending"
	SendStatusSent    SendStatus = "sent"
	SendStatusFailed  SendStatus = "failed"
)

type MessageSendLog struct {
	ID           uint64     `gorm:"primaryKey;autoIncrement" json:"id"`
	MessageID    uint64     `gorm:"not null;index:idx_message_id" json:"message_id"`
	CustomerID   uint64     `gorm:"not null;index:idx_customer_id" json:"customer_id"`
	SendStatus   SendStatus `gorm:"type:enum('pending','sent','failed');default:'pending'" json:"send_status"`
	SendTime     *time.Time `gorm:"type:datetime" json:"send_time,omitempty"`
	ErrorMessage string     `gorm:"type:text" json:"error_message,omitempty"`
	CreatedAt    time.Time  `gorm:"type:datetime;default:CURRENT_TIMESTAMP" json:"created_at"`

	// 关联关系
	Message  Message  `gorm:"foreignKey:MessageID" json:"message,omitempty"`
	Customer Customer `gorm:"foreignKey:CustomerID" json:"customer,omitempty"`
}

func (MessageSendLog) TableName() string {
	return "message_send_logs"
}