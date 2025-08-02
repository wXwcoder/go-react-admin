package model

import (
	"encoding/json"
	"fmt"
	"time"

	"gorm.io/gorm"
)

// DynamicData 动态数据通用结构
// 这个结构用于前后端数据交互，实际存储在动态创建的表中
type DynamicData struct {
	ID        uint                   `json:"id"`
	CreatedAt time.Time              `json:"created_at"`
	UpdatedAt time.Time              `json:"updated_at"`
	DeletedAt *time.Time             `json:"deleted_at,omitempty"`
	TenantID  uint                   `json:"tenant_id"`
	Data      map[string]interface{} `json:"data"`
}

// DynamicDataRequest 动态数据请求结构
type DynamicDataRequest struct {
	TableID uint                   `json:"table_id" validate:"required"`
	Data    map[string]interface{} `json:"data" validate:"required"`
}

// DynamicDataResponse 动态数据响应结构
type DynamicDataResponse struct {
	ID        uint                   `json:"id"`
	CreatedAt time.Time              `json:"created_at"`
	UpdatedAt time.Time              `json:"updated_at"`
	TenantID  uint                   `json:"tenant_id"`
	Data      map[string]interface{} `json:"data"`
}

// DynamicDataQuery 动态数据查询结构
type DynamicDataQuery struct {
	TableID    uint                     `json:"table_id" validate:"required"`
	Conditions []DynamicQueryCondition  `json:"conditions"`
	Logic      string                   `json:"logic"` // AND, OR
	Sort       []DynamicQuerySort       `json:"sort"`
	Page       int                      `json:"page"`
	Size       int                      `json:"size"`
	Fields     []string                 `json:"fields"` // 指定返回字段
}

// DynamicQueryCondition 查询条件
type DynamicQueryCondition struct {
	Field    string      `json:"field" validate:"required"`
	Operator string      `json:"operator" validate:"required"` // =, !=, >, <, >=, <=, like, in, not_in, between
	Value    interface{} `json:"value"`
	Values   []interface{} `json:"values,omitempty"` // 用于in, not_in, between操作
}

// DynamicQuerySort 排序条件
type DynamicQuerySort struct {
	Field string `json:"field" validate:"required"`
	Order string `json:"order" validate:"oneof=asc desc"` // asc, desc
}

// DynamicDataStatistics 数据统计结构
type DynamicDataStatistics struct {
	TableID      uint                   `json:"table_id"`
	GroupBy      []string               `json:"group_by"`
	Aggregations []DynamicAggregation   `json:"aggregations"`
	Conditions   []DynamicQueryCondition `json:"conditions"`
	Logic        string                 `json:"logic"`
}

// DynamicAggregation 聚合函数
type DynamicAggregation struct {
	Field    string `json:"field" validate:"required"`
	Function string `json:"function" validate:"required,oneof=count sum avg max min"`
	Alias    string `json:"alias"`
}

// DynamicBatchOperation 批量操作结构
type DynamicBatchOperation struct {
	TableID   uint          `json:"table_id" validate:"required"`
	Operation string        `json:"operation" validate:"required,oneof=delete update"`
	IDs       []uint        `json:"ids" validate:"required"`
	Data      map[string]interface{} `json:"data,omitempty"` // 用于批量更新
}

// DynamicDataHistory 数据变更历史
type DynamicDataHistory struct {
	ID          uint           `gorm:"primaryKey" json:"id"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
	
	TableID     uint           `gorm:"index" json:"table_id"`
	DataID      uint           `gorm:"index" json:"data_id"`
	Operation   string         `gorm:"size:20" json:"operation"` // create, update, delete
	OldData     json.RawMessage `gorm:"type:json" json:"old_data"`
	NewData     json.RawMessage `gorm:"type:json" json:"new_data"`
	ChangedBy   uint           `gorm:"index" json:"changed_by"`
	ChangeReason string        `gorm:"size:500" json:"change_reason"`
}

// TableName 自定义表名
func (DynamicDataHistory) TableName() string {
	return "dynamic_data_history"
}

// ValidateCondition 验证查询条件
func (c *DynamicQueryCondition) ValidateCondition() error {
	validOperators := []string{"=", "!=", ">", "<", ">=", "<=", "like", "in", "not_in", "between", "is_null", "is_not_null"}
	
	for _, op := range validOperators {
		if c.Operator == op {
			// 检查特定操作符的值要求
			switch c.Operator {
			case "in", "not_in":
				if c.Values == nil || len(c.Values) == 0 {
					return fmt.Errorf("操作符 %s 需要提供 values 数组", c.Operator)
				}
			case "between":
				if c.Values == nil || len(c.Values) != 2 {
					return fmt.Errorf("操作符 between 需要提供包含两个值的 values 数组")
				}
			case "is_null", "is_not_null":
				// 这些操作符不需要值
			default:
				if c.Value == nil {
					return fmt.Errorf("操作符 %s 需要提供 value", c.Operator)
				}
			}
			return nil
		}
	}
	
	return fmt.Errorf("不支持的操作符: %s", c.Operator)
}

// BuildWhereClause 构建WHERE子句
func (c *DynamicQueryCondition) BuildWhereClause(tableName string) (string, []interface{}) {
	fieldName := fmt.Sprintf("`%s`.`%s`", tableName, c.Field)
	
	switch c.Operator {
	case "=", "!=", ">", "<", ">=", "<=":
		return fmt.Sprintf("%s %s ?", fieldName, c.Operator), []interface{}{c.Value}
	case "like":
		return fmt.Sprintf("%s LIKE ?", fieldName), []interface{}{fmt.Sprintf("%%%v%%", c.Value)}
	case "in":
		placeholders := ""
		for i := range c.Values {
			if i > 0 {
				placeholders += ","
			}
			placeholders += "?"
		}
		return fmt.Sprintf("%s IN (%s)", fieldName, placeholders), c.Values
	case "not_in":
		placeholders := ""
		for i := range c.Values {
			if i > 0 {
				placeholders += ","
			}
			placeholders += "?"
		}
		return fmt.Sprintf("%s NOT IN (%s)", fieldName, placeholders), c.Values
	case "between":
		return fmt.Sprintf("%s BETWEEN ? AND ?", fieldName), c.Values
	case "is_null":
		return fmt.Sprintf("%s IS NULL", fieldName), []interface{}{}
	case "is_not_null":
		return fmt.Sprintf("%s IS NOT NULL", fieldName), []interface{}{}
	default:
		return fmt.Sprintf("%s = ?", fieldName), []interface{}{c.Value}
	}
}

// GetDefaultQuery 获取默认查询参数
func GetDefaultQuery() *DynamicDataQuery {
	return &DynamicDataQuery{
		Logic: "AND",
		Page:  1,
		Size:  20,
		Sort: []DynamicQuerySort{
			{Field: "id", Order: "desc"},
		},
	}
}

// ValidateData 验证数据格式
func (d *DynamicDataRequest) ValidateData(fields []DynamicField) error {
	for _, field := range fields {
		value, exists := d.Data[field.FieldName]
		
		// 检查必填字段
		if field.IsRequired && (!exists || value == nil || value == "") {
			return fmt.Errorf("字段 %s 是必填的", field.DisplayName)
		}
		
		// 如果字段不存在且不是必填的，跳过验证
		if !exists {
			continue
		}
		
		// 根据字段类型验证数据
		if err := validateFieldValue(field, value); err != nil {
			return fmt.Errorf("字段 %s 验证失败: %v", field.DisplayName, err)
		}
	}
	
	return nil
}

// validateFieldValue 验证字段值
func validateFieldValue(field DynamicField, value interface{}) error {
	switch field.FieldType {
	case "int":
		if _, ok := value.(float64); !ok {
			if _, ok := value.(int); !ok {
				return fmt.Errorf("必须是整数")
			}
		}
	case "float":
		if _, ok := value.(float64); !ok {
			if _, ok := value.(int); !ok {
				return fmt.Errorf("必须是数字")
			}
		}
	case "boolean":
		if _, ok := value.(bool); !ok {
			return fmt.Errorf("必须是布尔值")
		}
	case "string", "text":
		if _, ok := value.(string); !ok {
			return fmt.Errorf("必须是字符串")
		}
	case "select":
		// 验证选项值
		options, err := field.GetOptions()
		if err != nil {
			return err
		}
		
		valueStr := fmt.Sprintf("%v", value)
		for _, option := range options {
			if option.Value == valueStr {
				return nil
			}
		}
		return fmt.Errorf("无效的选项值")
	case "multiselect":
		// 验证多选值
		values, ok := value.([]interface{})
		if !ok {
			return fmt.Errorf("多选字段必须是数组")
		}
		
		options, err := field.GetOptions()
		if err != nil {
			return err
		}
		
		optionMap := make(map[string]bool)
		for _, option := range options {
			optionMap[option.Value] = true
		}
		
		for _, v := range values {
			valueStr := fmt.Sprintf("%v", v)
			if !optionMap[valueStr] {
				return fmt.Errorf("无效的选项值: %s", valueStr)
			}
		}
	}
	
	return nil
}