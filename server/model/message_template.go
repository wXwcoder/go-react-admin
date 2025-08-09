package model

import (
	"database/sql/driver"
	"fmt"
	"time"
)

type MessageTemplateType string

const (
	MessageTemplateTypeSystem  MessageTemplateType = "system"
	MessageTemplateTypeNotice  MessageTemplateType = "notice"
	MessageTemplateTypePrivate MessageTemplateType = "private"
)

type MessageTemplate struct {
	ID              uint64              `gorm:"primaryKey;autoIncrement" json:"id"`
	Name            string              `gorm:"type:varchar(100);not null;uniqueIndex:unique_name" json:"name"`
	TitleTemplate   string              `gorm:"type:varchar(255);not null" json:"title_template"`
	ContentTemplate string              `gorm:"type:text;not null" json:"content_template"`
	Type            MessageTemplateType `gorm:"type:enum('system','notice','private');default:'system'" json:"type"`
	Variables       JSON                `gorm:"type:json" json:"variables,omitempty"`
	IsActive        bool                `gorm:"type:tinyint(1);default:true" json:"is_active"`
	CreatedAt       time.Time           `gorm:"type:datetime;default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt       time.Time           `gorm:"type:datetime;default:CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP" json:"updated_at"`
}

func (MessageTemplate) TableName() string {
	return "message_templates"
}

// BeforeCreate 钩子函数
func (mt *MessageTemplate) BeforeCreate() error {
	if mt.Type == "" {
		mt.Type = MessageTemplateTypeSystem
	}
	return nil
}

// JSON type for GORM to handle JSON fields
type JSON []byte

func (j JSON) Value() (driver.Value, error) {
	if len(j) == 0 {
		return nil, nil
	}
	return string(j), nil
}

func (j *JSON) Scan(value interface{}) error {
	if value == nil {
		*j = nil
		return nil
	}
	var bytes []byte
	switch v := value.(type) {
	case []byte:
		bytes = v
	case string:
		bytes = []byte(v)
	default:
		return fmt.Errorf("failed to unmarshal JSON value: %v", value)
	}
	
	if len(bytes) == 0 {
		*j = nil
		return nil
	}
	
	*j = JSON(bytes)
	return nil
}