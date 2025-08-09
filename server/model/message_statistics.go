package model

import (
	"time"
)

type MessageStatistics struct {
	ID         uint64    `gorm:"primaryKey;autoIncrement" json:"id"`
	Date       time.Time `gorm:"type:date;not null;uniqueIndex:unique_date_type,priority:1" json:"date"`
	MessageType string    `gorm:"type:varchar(50);not null;uniqueIndex:unique_date_type,priority:2" json:"message_type"`
	TotalCount int       `gorm:"type:int;default:0" json:"total_count"`
	ReadCount  int       `gorm:"type:int;default:0" json:"read_count"`
	CreatedAt  time.Time `gorm:"type:datetime;default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt  time.Time `gorm:"type:datetime;default:CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP" json:"updated_at"`
}

func (MessageStatistics) TableName() string {
	return "message_statistics"
}