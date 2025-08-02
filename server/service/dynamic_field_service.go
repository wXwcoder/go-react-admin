package service

import (
	"errors"
	"fmt"
	"strings"

	"go-react-admin/global"
	"go-react-admin/model"
)

type DynamicFieldService struct{}

// CreateField 创建动态字段
func (dfs *DynamicFieldService) CreateField(field *model.DynamicField) error {
	// 检查字段名是否已存在
	var count int64
	global.DB.Model(&model.DynamicField{}).Where("table_id = ? AND field_name = ?",
		field.TableID, field.FieldName).Count(&count)
	if count > 0 {
		return errors.New("字段名已存在")
	}

	// 验证字段配置
	if err := dfs.validateField(field); err != nil {
		return err
	}

	// 获取表信息
	var table model.DynamicTable
	if err := global.DB.First(&table, field.TableID).Error; err != nil {
		return fmt.Errorf("获取表信息失败: %v", err)
	}

	// 开启事务
	tx := global.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// 创建字段记录
	if err := tx.Create(field).Error; err != nil {
		tx.Rollback()
		return err
	}

	// 在物理表中添加字段
	if err := addFieldToPhysicalTable(table.TableName, field); err != nil {
		tx.Rollback()
		return fmt.Errorf("添加物理表字段失败: %v", err)
	}

	return tx.Commit().Error
}

// GetFieldsByTableID 根据表ID获取字段列表
func (dfs *DynamicFieldService) GetFieldsByTableID(tableID uint) ([]model.DynamicField, error) {
	var fields []model.DynamicField
	if err := global.DB.Where("table_id = ?", tableID).Order("sort_order ASC, id ASC").Find(&fields).Error; err != nil {
		return nil, err
	}
	return fields, nil
}

// GetFieldByID 根据ID获取字段
func (dfs *DynamicFieldService) GetFieldByID(id uint) (*model.DynamicField, error) {
	var field model.DynamicField
	if err := global.DB.First(&field, id).Error; err != nil {
		return nil, err
	}
	return &field, nil
}

// UpdateField 更新字段
func (dfs *DynamicFieldService) UpdateField(field *model.DynamicField) error {
	// 检查字段是否存在
	var existingField model.DynamicField
	if err := global.DB.First(&existingField, field.ID).Error; err != nil {
		return err
	}

	// 如果修改了字段名，检查是否重复
	if existingField.FieldName != field.FieldName {
		var count int64
		global.DB.Model(&model.DynamicField{}).Where("table_id = ? AND field_name = ? AND id != ?",
			field.TableID, field.FieldName, field.ID).Count(&count)
		if count > 0 {
			return errors.New("字段名已存在")
		}
	}

	// 验证字段配置
	if err := dfs.validateField(field); err != nil {
		return err
	}

	// 获取表信息
	var table model.DynamicTable
	if err := global.DB.First(&table, field.TableID).Error; err != nil {
		return fmt.Errorf("获取表信息失败: %v", err)
	}

	// 开启事务
	tx := global.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// 更新字段记录
	if err := tx.Save(field).Error; err != nil {
		tx.Rollback()
		return err
	}

	// 更新物理表字段
	if err := updateFieldInPhysicalTable(table.TableName, existingField.FieldName, field); err != nil {
		tx.Rollback()
		return fmt.Errorf("更新物理表字段失败: %v", err)
	}

	return tx.Commit().Error
}

// DeleteField 删除字段
func (dfs *DynamicFieldService) DeleteField(id uint) error {
	// 检查字段是否存在
	field, err := dfs.GetFieldByID(id)
	if err != nil {
		return err
	}

	// 获取表信息
	var table model.DynamicTable
	if err := global.DB.First(&table, field.TableID).Error; err != nil {
		return fmt.Errorf("获取表信息失败: %v", err)
	}

	// 开启事务
	tx := global.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// 删除字段记录
	if err := tx.Delete(&model.DynamicField{}, id).Error; err != nil {
		tx.Rollback()
		return err
	}

	// 从物理表中删除字段
	if err := removeFieldFromPhysicalTable(table.TableName, field.FieldName); err != nil {
		// 如果删除物理字段失败，记录错误但不回滚
		// 因为有些情况下物理字段可能已经被删除或不存在
		fmt.Printf("警告：从物理表删除字段失败: %v", err)
	}

	return tx.Commit().Error
}

// UpdateFieldOrder 更新字段排序
func (dfs *DynamicFieldService) UpdateFieldOrder(fieldIDs []uint) error {
	if len(fieldIDs) == 0 {
		return nil
	}

	// 更新每个字段的排序
	for index, fieldID := range fieldIDs {
		if err := global.DB.Model(&model.DynamicField{}).
			Where("id = ?", fieldID).
			Update("sort_order", index+1).Error; err != nil {
			return err
		}
	}

	return nil
}

// ToggleFieldStatus 切换字段状态
func (dfs *DynamicFieldService) ToggleFieldStatus(id uint) error {
	var field model.DynamicField
	if err := global.DB.First(&field, id).Error; err != nil {
		return err
	}

	// 切换状态
	if field.Status == 1 {
		field.Status = 2
	} else {
		field.Status = 1
	}

	return global.DB.Save(&field).Error
}

// BatchCreateFields 批量创建字段
func (dfs *DynamicFieldService) BatchCreateFields(fields []model.DynamicField) error {
	// 开启事务
	tx := global.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	for _, field := range fields {
		// 检查字段名是否已存在
		var count int64
		tx.Model(&model.DynamicField{}).Where("table_id = ? AND field_name = ?",
			field.TableID, field.FieldName).Count(&count)
		if count > 0 {
			tx.Rollback()
			return fmt.Errorf("字段名 %s 已存在", field.FieldName)
		}

		// 验证字段配置
		if err := dfs.validateField(&field); err != nil {
			tx.Rollback()
			return err
		}

		// 创建字段记录
		if err := tx.Create(&field).Error; err != nil {
			tx.Rollback()
			return err
		}
	}

	return tx.Commit().Error
}

// GetFieldTypes 获取支持的字段类型
func (dfs *DynamicFieldService) GetFieldTypes() []map[string]interface{} {
	return []map[string]interface{}{
		{
			"value":       "string",
			"label":       "单行文本",
			"description": "短文本输入，适用于姓名、标题等",
			"config": map[string]interface{}{
				"max_length": 255,
				"min_length": 0,
			},
		},
		{
			"value":       "text",
			"label":       "多行文本",
			"description": "长文本输入，适用于描述、备注等",
			"config": map[string]interface{}{
				"max_length": 65535,
			},
		},
		{
			"value":       "int",
			"label":       "整数",
			"description": "整数类型，适用于数量、年龄等",
			"config": map[string]interface{}{
				"max_value": 2147483647,
				"min_value": -2147483648,
			},
		},
		{
			"value":       "float",
			"label":       "小数",
			"description": "浮点数类型，适用于价格、比例等",
			"config": map[string]interface{}{
				"precision": 2,
			},
		},
		{
			"value":       "date",
			"label":       "日期",
			"description": "日期类型，格式：YYYY-MM-DD",
			"config": map[string]interface{}{
				"format": "YYYY-MM-DD",
			},
		},
		{
			"value":       "datetime",
			"label":       "日期时间",
			"description": "日期时间类型，格式：YYYY-MM-DD HH:mm:ss",
			"config": map[string]interface{}{
				"format": "YYYY-MM-DD HH:mm:ss",
			},
		},
		{
			"value":       "boolean",
			"label":       "布尔值",
			"description": "是/否选择，适用于开关状态",
			"config": map[string]interface{}{
				"true_label":  "是",
				"false_label": "否",
			},
		},
		{
			"value":       "select",
			"label":       "单选下拉",
			"description": "从预定义选项中选择一个",
			"config": map[string]interface{}{
				"options": []map[string]string{},
			},
		},
		{
			"value":       "multiselect",
			"label":       "多选下拉",
			"description": "从预定义选项中选择多个",
			"config": map[string]interface{}{
				"options":        []map[string]string{},
				"max_selections": 10,
			},
		},
		{
			"value":       "file",
			"label":       "文件上传",
			"description": "文件上传，支持多种格式",
			"config": map[string]interface{}{
				"max_size":      10485760, // 10MB
				"allowed_types": []string{"pdf", "doc", "docx", "xls", "xlsx"},
			},
		},
		{
			"value":       "image",
			"label":       "图片上传",
			"description": "图片上传，支持常见图片格式",
			"config": map[string]interface{}{
				"max_size":      5242880, // 5MB
				"allowed_types": []string{"jpg", "jpeg", "png", "gif", "webp"},
			},
		},
	}
}

// validateField 验证字段配置
func (dfs *DynamicFieldService) validateField(field *model.DynamicField) error {
	// 验证字段名
	if err := dfs.validateFieldName(field.FieldName); err != nil {
		return err
	}

	// 验证字段类型
	if !dfs.isValidFieldType(field.FieldType) {
		return errors.New("不支持的字段类型")
	}

	// 验证字段配置
	if err := dfs.validateFieldConfig(field); err != nil {
		return err
	}

	return nil
}

// validateFieldName 验证字段名
func (dfs *DynamicFieldService) validateFieldName(fieldName string) error {
	if len(fieldName) < 1 || len(fieldName) > 100 {
		return errors.New("字段名长度必须在1-100个字符之间")
	}

	// 检查是否包含非法字符
	if !isValidFieldName(fieldName) {
		return errors.New("字段名只能包含字母、数字和下划线，且必须以字母开头")
	}

	// 检查是否为保留字段名
	if isReservedFieldName(fieldName) {
		return errors.New("字段名不能使用系统保留字段")
	}

	return nil
}

// isValidFieldType 检查字段类型是否有效
func (dfs *DynamicFieldService) isValidFieldType(fieldType string) bool {
	validTypes := []string{
		"string", "text", "int", "float", "date", "datetime",
		"boolean", "select", "multiselect", "file", "image",
	}

	for _, validType := range validTypes {
		if fieldType == validType {
			return true
		}
	}

	return false
}

// addFieldToPhysicalTable 添加字段到物理表
func addFieldToPhysicalTable(tableName string, field *model.DynamicField) error {
	// 清理表名
	tableName = SanitizeTableName(tableName)

	// 构建字段定义
	columnDef, err := buildColumnDefinition(field)
	if err != nil {
		return err
	}

	// 构建SQL语句
	sql := fmt.Sprintf("ALTER TABLE `%s` ADD COLUMN %s", tableName, columnDef)

	// 执行SQL
	if err := global.DB.Exec(sql).Error; err != nil {
		return fmt.Errorf("failed to add column to physical table: %v", err)
	}

	return nil
}

// updateFieldInPhysicalTable 更新物理表字段
func updateFieldInPhysicalTable(tableName string, oldFieldName string, newField *model.DynamicField) error {
	// 清理表名
	tableName = SanitizeTableName(tableName)

	// 构建新的字段定义
	columnDef, err := buildColumnDefinition(newField)
	if err != nil {
		return err
	}

	// 构建SQL语句
	sql := fmt.Sprintf("ALTER TABLE `%s` CHANGE COLUMN `%s` %s", tableName, oldFieldName, columnDef)

	// 执行SQL
	if err := global.DB.Exec(sql).Error; err != nil {
		return fmt.Errorf("failed to update column in physical table: %v", err)
	}

	return nil
}

// removeFieldFromPhysicalTable 从物理表删除字段
func removeFieldFromPhysicalTable(tableName string, fieldName string) error {
	// 清理表名
	tableName = SanitizeTableName(tableName)

	// 构建SQL语句
	sql := fmt.Sprintf("ALTER TABLE `%s` DROP COLUMN `%s`", tableName, fieldName)

	// 执行SQL
	if err := global.DB.Exec(sql).Error; err != nil {
		return fmt.Errorf("failed to drop column from physical table: %v", err)
	}

	return nil
}

// SanitizeTableName 清理表名以防止SQL注入
// func SanitizeTableName(tableName string) string {
// 	// 移除特殊字符，只允许字母、数字和下划线
// 	var sanitized strings.Builder
// 	for _, char := range tableName {
// 		if (char >= 'a' && char <= 'z') || (char >= 'A' && char <= 'Z') ||
// 		   (char >= '0' && char <= '9') || char == '_' {
// 			sanitized.WriteRune(char)
// 		}
// 	}
// 	return sanitized.String()
// }

// buildColumnDefinition 根据字段属性构建MySQL列定义
func buildColumnDefinition(field *model.DynamicField) (string, error) {
	var columnDef strings.Builder

	// 添加字段名
	columnDef.WriteString(fmt.Sprintf("`%s` ", field.FieldName))

	// 根据字段类型确定MySQL数据类型
	mysqlType := getMySQLType(field.FieldType, field)
	columnDef.WriteString(mysqlType)

	// 处理是否必填
	if field.IsRequired {
		columnDef.WriteString(" NOT NULL")
	} else {
		columnDef.WriteString(" NULL")
	}

	// 处理默认值
	if field.DefaultValue != "" {
		// 转义默认值
		escaped := strings.ReplaceAll(field.DefaultValue, "'", "''")
		columnDef.WriteString(fmt.Sprintf(" DEFAULT '%s'", escaped))
	} else if !field.IsRequired {
		columnDef.WriteString(" DEFAULT NULL")
	}

	// 处理唯一约束
	if field.IsUnique {
		columnDef.WriteString(" UNIQUE")
	}

	// 处理注释
	if field.DisplayName != "" {
		escaped := strings.ReplaceAll(field.DisplayName, "'", "''")
		columnDef.WriteString(fmt.Sprintf(" COMMENT '%s'", escaped))
	}

	return columnDef.String(), nil
}

// getMySQLType 根据字段类型获取对应的MySQL数据类型
func getMySQLType(fieldType string, field *model.DynamicField) string {
	switch fieldType {
	case "string":
		return "VARCHAR(255)"
	case "text":
		return "TEXT"
	case "int":
		return "INT"
	case "float":
		return "DECIMAL(10,2)"
	case "date":
		return "DATE"
	case "datetime":
		return "DATETIME"
	case "boolean":
		return "BOOLEAN"
	case "select":
		return "VARCHAR(255)"
	case "multiselect":
		return "TEXT" // 存储JSON数组
	case "file":
		return "VARCHAR(500)" // 存储文件路径
	case "image":
		return "VARCHAR(500)" // 存储图片路径
	default:
		return "VARCHAR(255)"
	}
}

// validateFieldConfig 验证字段配置
func (dfs *DynamicFieldService) validateFieldConfig(field *model.DynamicField) error {
	// 根据字段类型验证配置
	switch field.FieldType {
	case "select", "multiselect":
		// 验证选项配置
		options, err := field.GetOptions()
		if err != nil {
			return fmt.Errorf("选项配置格式错误: %v", err)
		}
		if len(options) == 0 {
			return errors.New("选择类型字段必须配置选项")
		}
	case "file", "image":
		// 验证文件配置
		fileConfig, err := field.GetFileConfig()
		if err != nil {
			return fmt.Errorf("文件配置格式错误: %v", err)
		}
		if fileConfig.MaxSize <= 0 {
			return errors.New("文件大小限制必须大于0")
		}
	}

	return nil
}

// isValidFieldName 检查字段名是否有效
func isValidFieldName(name string) bool {
	if len(name) == 0 {
		return false
	}

	// 必须以字母开头
	if !((name[0] >= 'a' && name[0] <= 'z') || (name[0] >= 'A' && name[0] <= 'Z')) {
		return false
	}

	// 只能包含字母、数字和下划线
	for _, char := range name {
		if !((char >= 'a' && char <= 'z') || (char >= 'A' && char <= 'Z') ||
			(char >= '0' && char <= '9') || char == '_') {
			return false
		}
	}

	return true
}

// isReservedFieldName 检查是否为保留字段名
func isReservedFieldName(name string) bool {
	reservedFields := []string{
		"id", "created_at", "updated_at", "deleted_at", "tenant_id",
		"password", "token", "session", "admin", "root", "system",
	}

	for _, reserved := range reservedFields {
		if name == reserved {
			return true
		}
	}

	return false
}
