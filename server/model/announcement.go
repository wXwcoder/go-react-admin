package model

import (
	"time"
)

type AnnouncementType string

type AnnouncementStatus string

const (
	AnnouncementTypeSystem      AnnouncementType = "system"
	AnnouncementTypeNotice      AnnouncementType = "notice"
	AnnouncementTypeMaintenance AnnouncementType = "maintenance"
	AnnouncementTypeUpdate      AnnouncementType = "update"
)

const (
	AnnouncementStatusDraft     AnnouncementStatus = "draft"
	AnnouncementStatusPublished AnnouncementStatus = "published"
	AnnouncementStatusRevoked AnnouncementStatus = "revoked"
)

type Announcement struct {
	ID        uint64             `gorm:"primaryKey;autoIncrement" json:"id"`
	Title     string             `gorm:"type:varchar(255);not null" json:"title"`
	Content   string             `gorm:"type:text;not null" json:"content"`
	Type      AnnouncementType   `gorm:"type:enum('system','notice','maintenance','update');default:'notice'" json:"type"`
	Status    AnnouncementStatus `gorm:"type:enum('draft','published','revoked');default:'draft'" json:"status"`
	Priority  int                `gorm:"type:int;default:0" json:"priority"`
	ExpiredAt *time.Time         `gorm:"type:datetime" json:"expired_at,omitempty"`
	ReadCount int                `gorm:"type:int;default:0" json:"read_count"`
	CreatedAt time.Time          `gorm:"type:datetime;default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt time.Time          `gorm:"type:datetime;default:CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP" json:"updated_at"`
	DeletedAt *time.Time         `gorm:"type:datetime;index" json:"deleted_at,omitempty"`

	// 关联关系
	AnnouncementReads []AnnouncementRead `gorm:"foreignKey:AnnouncementID" json:"announcement_reads,omitempty"`
}

func (Announcement) TableName() string {
	return "announcements"
}

// BeforeCreate 钩子函数
func (a *Announcement) BeforeCreate() error {
	if a.Type == "" {
		a.Type = AnnouncementTypeNotice
	}
	if a.Status == "" {
		a.Status = AnnouncementStatusDraft
	}
	return nil
}