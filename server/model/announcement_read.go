package model

import (
	"time"
)

type AnnouncementRead struct {
	ID             uint64     `gorm:"primaryKey;autoIncrement" json:"id"`
	AnnouncementID uint64     `gorm:"not null;index:idx_announcement_customer,priority:1" json:"announcement_id"`
	CustomerID     uint64     `gorm:"not null;index:idx_announcement_customer,priority:2;index:idx_customer_read,priority:1" json:"customer_id"`
	IsRead         bool       `gorm:"type:tinyint(1);default:false" json:"is_read"`
	ReadTime       *time.Time `gorm:"type:datetime" json:"read_time,omitempty"`
	CreatedAt      time.Time  `gorm:"type:datetime;default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt      time.Time  `gorm:"type:datetime;default:CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP" json:"updated_at"`

	// 关联关系
	Announcement Announcement `gorm:"foreignKey:AnnouncementID" json:"announcement,omitempty"`
	Customer     Customer    `gorm:"foreignKey:CustomerID" json:"customer,omitempty"`
}

func (AnnouncementRead) TableName() string {
	return "announcement_reads"
}

// BeforeUpdate 钩子函数 - 更新已读时间
func (ar *AnnouncementRead) BeforeUpdate() error {
	if ar.IsRead && ar.ReadTime == nil {
		now := time.Now()
		ar.ReadTime = &now
	}
	return nil
}