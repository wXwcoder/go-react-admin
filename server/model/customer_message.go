package model

import (
	"time"
)

type CustomerMessage struct {
	ID         uint64     `gorm:"primaryKey;autoIncrement" json:"id"`
	MessageID  uint64     `gorm:"not null;index:idx_customer_message,priority:2;index:idx_message_customer,priority:1" json:"message_id"`
	CustomerID uint64     `gorm:"not null;index:idx_customer_message,priority:1;index:idx_message_customer,priority:2" json:"customer_id"`
	IsRead     bool       `gorm:"type:tinyint(1);default:false" json:"is_read"`
	ReadTime   *time.Time `gorm:"type:datetime" json:"read_time,omitempty"`
	IsDeleted  bool       `gorm:"type:tinyint(1);default:false" json:"is_deleted"`
	CreatedAt  time.Time  `gorm:"type:datetime;default:CURRENT_TIMESTAMP" json:"created_at"`

	// 关联关系
	Message  Message  `gorm:"foreignKey:MessageID" json:"message,omitempty"`
	Customer Customer `gorm:"foreignKey:CustomerID" json:"customer,omitempty"`
}

func (CustomerMessage) TableName() string {
	return "customer_messages"
}

// BeforeUpdate 钩子函数 - 更新已读时间
func (cm *CustomerMessage) BeforeUpdate() error {
	if cm.IsRead && cm.ReadTime == nil {
		now := time.Now()
		cm.ReadTime = &now
	}
	return nil
}