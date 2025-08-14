package model

import (
	"time"
)

type MessageType string

type MessageStatus string

const (
	MessageTypeSystem  MessageType = "system"
	MessageTypeNotice  MessageType = "notice"
	MessageTypePrivate MessageType = "private"
)

const (
	MessageStatusDraft     MessageStatus = "draft"
	MessageStatusPublished MessageStatus = "published"
	MessageStatusRevoked   MessageStatus = "revoked"
)

type Message struct {
	ID         uint64        `gorm:"primaryKey;autoIncrement" json:"id"`
	Title      string        `gorm:"type:varchar(255);not null" json:"title"`
	Content    string        `gorm:"type:text;not null" json:"content"`
	Type       MessageType   `gorm:"type:enum('system','notice','private');default:'system'" json:"type"`
	Status     MessageStatus `gorm:"type:enum('draft','published','revoked');default:'draft'" json:"status"`
	Priority   int           `gorm:"type:int;default:0" json:"priority"`
	SenderID   *uint64       `gorm:"type:bigint" json:"sender_id,omitempty"`
	SenderType string        `gorm:"type:varchar(50);default:'system'" json:"sender_type"`
	TargetType string        `gorm:"type:varchar(50);default:'all'" json:"target_type"`
	TargetID   *uint64       `gorm:"type:bigint" json:"target_id,omitempty"`
	ExpiredAt  *time.Time    `gorm:"type:datetime" json:"expired_at,omitempty"`
	ReadCount  int           `gorm:"type:int;default:0" json:"read_count"`
	CreatedAt  time.Time     `gorm:"type:datetime;default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt  time.Time     `gorm:"type:datetime;default:CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP" json:"updated_at"`
	DeletedAt  *time.Time    `gorm:"type:datetime;index" json:"deleted_at,omitempty"`

	// 关联关系
}

func (Message) TableName() string {
	return "messages"
}

// BeforeCreate 钩子函数
func (m *Message) BeforeCreate() error {
	if m.Type == "" {
		m.Type = MessageTypeSystem
	}
	if m.Status == "" {
		m.Status = MessageStatusDraft
	}
	return nil
}
