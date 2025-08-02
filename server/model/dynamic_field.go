package model

import (
	"encoding/json"
	"time"

	"gorm.io/gorm"
)

// DynamicField 动态字段定义
type DynamicField struct {
	ID          uint           `gorm:"primaryKey" json:"id"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
	
	TableID     uint           `gorm:"index" json:"table_id" validate:"required"`
	FieldName   string         `gorm:"size:100" json:"field_name" validate:"required,min=1,max=100"`
	DisplayName string         `gorm:"size:100" json:"display_name" validate:"required,min=1,max=100"`
	FieldType   string         `gorm:"size:50" json:"field_type" validate:"required,oneof=string text int float date datetime boolean select multiselect file image"`
	IsRequired  bool           `gorm:"default:false" json:"is_required"`
	IsUnique    bool           `gorm:"default:false" json:"is_unique"`
	IsSearchable bool          `gorm:"default:false" json:"is_searchable"`
	IsSortable  bool           `gorm:"default:false" json:"is_sortable"`
	DefaultValue string        `gorm:"size:500" json:"default_value"`
	Options     json.RawMessage `gorm:"type:json" json:"options"` // 用于select类型的选项
	Validation  json.RawMessage `gorm:"type:json" json:"validation"` // 验证规则
	SortOrder   int            `gorm:"default:0" json:"sort_order"`
	Status      int            `gorm:"default:1" json:"status" validate:"oneof=1 2"`
}

// TableName 自定义表名
func (DynamicField) TableName() string {
	return "dynamic_fields"
}

// FieldValidation 字段验证规则结构
type FieldValidation struct {
	Required    bool    `json:"required"`
	MinLength   int     `json:"min_length,omitempty"`
	MaxLength   int     `json:"max_length,omitempty"`
	MinValue    float64 `json:"min_value,omitempty"`
	MaxValue    float64 `json:"max_value,omitempty"`
	Pattern     string  `json:"pattern,omitempty"`
	CustomRules string  `json:"custom_rules,omitempty"`
}

// SelectOption 下拉选项
type SelectOption struct {
	Label string `json:"label"`
	Value string `json:"value"`
	Color string `json:"color,omitempty"`
}

// FileConfig 文件字段配置
type FileConfig struct {
	MaxSize     int      `json:"max_size"`     // 最大文件大小(MB)
	AllowedTypes []string `json:"allowed_types"` // 允许的文件类型
	Multiple    bool     `json:"multiple"`     // 是否允许多文件
}

// GetValidation 获取验证规则
func (f *DynamicField) GetValidation() (*FieldValidation, error) {
	if f.Validation == nil {
		return &FieldValidation{}, nil
	}
	
	var validation FieldValidation
	err := json.Unmarshal(f.Validation, &validation)
	return &validation, err
}

// SetValidation 设置验证规则
func (f *DynamicField) SetValidation(validation *FieldValidation) error {
	data, err := json.Marshal(validation)
	if err != nil {
		return err
	}
	f.Validation = data
	return nil
}

// GetOptions 获取选项列表
func (f *DynamicField) GetOptions() ([]SelectOption, error) {
	if f.Options == nil {
		return []SelectOption{}, nil
	}
	
	var options []SelectOption
	err := json.Unmarshal(f.Options, &options)
	return options, err
}

// SetOptions 设置选项列表
func (f *DynamicField) SetOptions(options []SelectOption) error {
	data, err := json.Marshal(options)
	if err != nil {
		return err
	}
	f.Options = data
	return nil
}

// GetFileConfig 获取文件配置
func (f *DynamicField) GetFileConfig() (*FileConfig, error) {
	if f.Options == nil {
		return &FileConfig{}, nil
	}
	
	var config FileConfig
	err := json.Unmarshal(f.Options, &config)
	return &config, err
}

// SetFileConfig 设置文件配置
func (f *DynamicField) SetFileConfig(config *FileConfig) error {
	data, err := json.Marshal(config)
	if err != nil {
		return err
	}
	f.Options = data
	return nil
}

// IsSelectType 判断是否为选择类型字段
func (f *DynamicField) IsSelectType() bool {
	return f.FieldType == "select" || f.FieldType == "multiselect"
}

// IsFileType 判断是否为文件类型字段
func (f *DynamicField) IsFileType() bool {
	return f.FieldType == "file" || f.FieldType == "image"
}

// IsNumericType 判断是否为数值类型字段
func (f *DynamicField) IsNumericType() bool {
	return f.FieldType == "int" || f.FieldType == "float"
}

// IsDateType 判断是否为日期类型字段
func (f *DynamicField) IsDateType() bool {
	return f.FieldType == "date" || f.FieldType == "datetime"
}

// GetMySQLColumnType 获取MySQL列类型
func (f *DynamicField) GetMySQLColumnType() string {
	switch f.FieldType {
	case "string":
		return "varchar(255)"
	case "text":
		return "text"
	case "int":
		return "int"
	case "float":
		return "decimal(10,2)"
	case "date":
		return "date"
	case "datetime":
		return "datetime(3)"
	case "boolean":
		return "tinyint(1)"
	case "select", "multiselect":
		return "varchar(255)"
	case "file", "image":
		return "varchar(500)"
	default:
		return "varchar(255)"
	}
}