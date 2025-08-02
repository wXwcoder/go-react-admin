package model

import (
	"encoding/json"
	"time"

	"gorm.io/gorm"
)

// DynamicTable 动态数据表定义
type DynamicTable struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`

	Name        string          `gorm:"uniqueIndex;size:100" json:"name" validate:"required,min=2,max=100"`
	DisplayName string          `gorm:"size:100" json:"display_name" validate:"required,min=2,max=100"`
	Description string          `gorm:"size:500" json:"description" validate:"max=500"`
	TableName   string          `gorm:"uniqueIndex;size:100" json:"table_name" validate:"required,min=2,max=100"`
	Fields      json.RawMessage `gorm:"type:json" json:"fields" validate:"required"`
	Status      int             `gorm:"default:1" json:"status" validate:"oneof=1 2"` // 1:启用 2:禁用
	TenantID    uint            `gorm:"index" json:"tenant_id"`

	// 关联字段
	FieldDefinitions []DynamicField `gorm:"foreignKey:TableID" json:"field_definitions,omitempty"`
}

// DynamicView 数据视图配置
type DynamicView struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`

	TableID   uint            `gorm:"index" json:"table_id" validate:"required"`
	ViewName  string          `gorm:"size:100" json:"view_name" validate:"required,min=1,max=100"`
	ViewType  string          `gorm:"size:50;default:table" json:"view_type" validate:"oneof=table card chart"`
	Config    json.RawMessage `gorm:"type:json" json:"config" validate:"required"`
	IsDefault bool            `gorm:"default:false" json:"is_default"`
	IsShared  bool            `gorm:"default:false" json:"is_shared"`
	CreatedBy uint            `gorm:"index" json:"created_by" validate:"required"`
}

// DynamicImportExportLog 数据导入导出记录
type DynamicImportExportLog struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`

	TableID       uint   `gorm:"index" json:"table_id" validate:"required"`
	OperationType string `gorm:"size:20" json:"operation_type" validate:"required,oneof=import export"`
	FileName      string `gorm:"size:255" json:"file_name" validate:"required"`
	FilePath      string `gorm:"size:500" json:"file_path"`
	Status        string `gorm:"size:20;default:pending" json:"status" validate:"oneof=pending processing success failed"`
	TotalRows     int    `json:"total_rows"`
	ProcessedRows int    `json:"processed_rows"`
	SuccessRows   int    `json:"success_rows"`
	FailedRows    int    `json:"failed_rows"`
	ErrorMessage  string `gorm:"type:text" json:"error_message"`
	CreatedBy     uint   `gorm:"index" json:"created_by" validate:"required"`
}

// TableName 自定义表名
func (DynamicTable) GetTableName() string {
	return "dynamic_tables"
}

// TableName 自定义表名
func (DynamicView) TableName() string {
	return "dynamic_views"
}

// TableName 自定义表名
func (DynamicImportExportLog) TableName() string {
	return "dynamic_import_export_logs"
}

// FieldConfig 字段配置结构
type FieldConfig struct {
	Type        string                 `json:"type"`
	Required    bool                   `json:"required"`
	Unique      bool                   `json:"unique"`
	MaxLength   int                    `json:"max_length,omitempty"`
	MinLength   int                    `json:"min_length,omitempty"`
	MaxValue    float64                `json:"max_value,omitempty"`
	MinValue    float64                `json:"min_value,omitempty"`
	Pattern     string                 `json:"pattern,omitempty"`
	Options     []SelectOption         `json:"options,omitempty"`
	DateFormat  string                 `json:"date_format,omitempty"`
	FileTypes   []string               `json:"file_types,omitempty"`
	MaxFileSize int                    `json:"max_file_size,omitempty"`
	CustomRules map[string]interface{} `json:"custom_rules,omitempty"`
}

// SelectOption 下拉选项
// type SelectOption struct {
// 	Label string `json:"label"`
// 	Value string `json:"value"`
// }

// ViewConfig 视图配置结构
type ViewConfig struct {
	Columns     []string               `json:"columns"`                // 显示的列
	Filters     map[string]interface{} `json:"filters"`                // 过滤条件
	Sort        ViewSort               `json:"sort"`                   // 排序
	PageSize    int                    `json:"page_size"`              // 每页显示数量
	Layout      string                 `json:"layout"`                 // 布局类型
	ChartType   string                 `json:"chart_type,omitempty"`   // 图表类型
	ChartConfig map[string]interface{} `json:"chart_config,omitempty"` // 图表配置
}

// ViewSort 视图排序
type ViewSort struct {
	Field string `json:"field"`
	Order string `json:"order"` // asc, desc
}

// TableData 通用数据格式
// 用于前端和后端的数据交互
// 实际存储时，会根据表结构转换为对应的格式
type TableData struct {
	ID        uint                   `json:"id"`
	CreatedAt time.Time              `json:"created_at"`
	UpdatedAt time.Time              `json:"updated_at"`
	Data      map[string]interface{} `json:"data"`
	TenantID  uint                   `json:"tenant_id"`
}

// GetViewConfig 获取视图配置
func (dv *DynamicView) GetViewConfig() (*ViewConfig, error) {
	if dv.Config == nil {
		return &ViewConfig{}, nil
	}

	var config ViewConfig
	err := json.Unmarshal(dv.Config, &config)
	return &config, err
}

// SetViewConfig 设置视图配置
func (dv *DynamicView) SetViewConfig(config *ViewConfig) error {
	data, err := json.Marshal(config)
	if err != nil {
		return err
	}
	dv.Config = data
	return nil
}

// GetProgress 获取导入导出进度
func (log *DynamicImportExportLog) GetProgress() float64 {
	if log.TotalRows == 0 {
		return 0
	}
	return float64(log.ProcessedRows) / float64(log.TotalRows) * 100
}

// IsCompleted 检查是否完成
func (log *DynamicImportExportLog) IsCompleted() bool {
	return log.Status == "success" || log.Status == "failed"
}

// IsSuccess 检查是否成功
func (log *DynamicImportExportLog) IsSuccess() bool {
	return log.Status == "success"
}
