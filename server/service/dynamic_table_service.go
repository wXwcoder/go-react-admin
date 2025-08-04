package service

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
	"time"

	"go-react-admin/global"
	"go-react-admin/model"

	"gorm.io/gorm"
)

type DynamicTableService struct{}

// CreateTable 创建动态表
func (dts *DynamicTableService) CreateTable(table *model.DynamicTable) error {
	// 如果TableName为空，从Name生成
	if table.TableName == "" {
		table.TableName = "dyn_" + strings.ToLower(strings.ReplaceAll(table.Name, " ", "_"))
	}

	// 检查表名是否已存在
	var count int64
	global.DB.Model(&model.DynamicTable{}).Where("table_name = ?", table.TableName).Count(&count)
	if count > 0 {
		return errors.New("表名已存在")
	}

	// 验证表名
	if valid, err := dts.ValidateTableName(table.TableName); !valid {
		return err
	}

	// 开启事务
	tx := global.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// 创建表记录
	if err := tx.Create(table).Error; err != nil {
		tx.Rollback()
		return err
	}

	// 如果没有字段定义，创建一些基础字段
	if len(table.FieldDefinitions) == 0 {
		defaultFields := []model.DynamicField{
			{
				TableID:     table.ID,
				FieldName:   "name",
				FieldType:   "string",
				DisplayName: "名称",
				IsRequired:  true,
				IsUnique:    false,
				SortOrder:   1,
				Status:      1,
			},
			{
				TableID:     table.ID,
				FieldName:   "description",
				FieldType:   "text",
				DisplayName: "描述",
				IsRequired:  false,
				IsUnique:    false,
				SortOrder:   2,
				Status:      1,
			},
			{
				TableID:     table.ID,
				FieldName:   "age",
				FieldType:   "int",
				DisplayName: "年龄",
				IsRequired:  false,
				IsUnique:    false,
				SortOrder:   3,
				Status:      1,
			},
			{
				TableID:     table.ID,
				FieldName:   "email",
				FieldType:   "string",
				DisplayName: "邮箱",
				IsRequired:  false,
				IsUnique:    false,
				SortOrder:   4,
				Status:      1,
			},
			{
				TableID:     table.ID,
				FieldName:   "status",
				FieldType:   "int",
				DisplayName: "状态",
				IsRequired:  false,
				IsUnique:    false,
				SortOrder:   5,
				Status:      1,
				DefaultValue: "1",
			},
		}

		for _, field := range defaultFields {
			if err := tx.Create(&field).Error; err != nil {
				tx.Rollback()
				return err
			}
		}
		table.FieldDefinitions = defaultFields
	}

	// 创建物理表（允许空表）
	if err := dts.createPhysicalTable(tx, table); err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}

// GetTableList 获取动态表列表
func (dts *DynamicTableService) GetTableList(page, pageSize int, search string) ([]model.DynamicTable, int64, error) {
	var tables []model.DynamicTable
	var total int64

	db := global.DB.Model(&model.DynamicTable{})

	// 搜索条件
	if search != "" {
		db = db.Where("name LIKE ? OR display_name LIKE ? OR description LIKE ?",
			"%"+search+"%", "%"+search+"%", "%"+search+"%")
	}

	// 获取总数
	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 分页查询
	offset := (page - 1) * pageSize
	if err := db.Preload("FieldDefinitions").Offset(offset).Limit(pageSize).Find(&tables).Error; err != nil {
		return nil, 0, err
	}

	return tables, total, nil
}

// GetTableByID 根据ID获取动态表
func (dts *DynamicTableService) GetTableByID(id uint) (*model.DynamicTable, error) {
	var table model.DynamicTable
	if err := global.DB.Preload("FieldDefinitions").First(&table, id).Error; err != nil {
		return nil, err
	}
	return &table, nil
}

// GetTableByName 根据表名获取动态表
func (dts *DynamicTableService) GetTableByName(tableName string) (*model.DynamicTable, error) {
	var table model.DynamicTable
	if err := global.DB.Preload("FieldDefinitions").Where("table_name = ?", tableName).First(&table).Error; err != nil {
		return nil, err
	}
	return &table, nil
}

// UpdateTable 更新动态表
func (dts *DynamicTableService) UpdateTable(table *model.DynamicTable) error {
	// 检查表是否存在
	var existingTable model.DynamicTable
	if err := global.DB.Preload("FieldDefinitions").First(&existingTable, table.ID).Error; err != nil {
		return err
	}

	// 检查名称是否已存在（排除当前记录）
	var count int64
	if err := global.DB.Model(&model.DynamicTable{}).
		Where("name = ? AND id != ?", table.Name, table.ID).
		Count(&count).Error; err != nil {
		return err
	}
	if count > 0 {
		return errors.New("表名称已存在")
	}

	// 检查表名是否已存在（排除当前记录）
	if err := global.DB.Model(&model.DynamicTable{}).
		Where("table_name = ? AND id != ?", table.TableName, table.ID).
		Count(&count).Error; err != nil {
		return err
	}
	if count > 0 {
		return errors.New("表名已存在")
	}

	// 开启事务
	tx := global.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// 更新表记录 - 排除时间戳字段
	updateData := map[string]interface{}{
		"name":         table.Name,
		"display_name": table.DisplayName,
		"description":  table.Description,
		"status":       table.Status,
		"updated_at":   time.Now(),
	}
	if err := tx.Model(&model.DynamicTable{}).Where("id = ?", table.ID).Updates(updateData).Error; err != nil {
		tx.Rollback()
		return err
	}

	// 如果字段有变化，更新物理表结构
	if err := dts.updatePhysicalTable(tx, &existingTable, table); err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}

// CreateDynamicTable 创建动态表
func (dts *DynamicTableService) CreateDynamicTable(table *model.DynamicTable) (*model.DynamicTable, error) {
	// 检查表名是否已存在
	var count int64
	global.DB.Model(&model.DynamicTable{}).Where("table_name = ?", table.TableName).Count(&count)
	if count > 0 {
		return nil, errors.New("表名已存在")
	}

	// 开启事务
	tx := global.DB.Begin()

	// 创建表记录
	if err := tx.Create(table).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	// 创建物理表
	if err := dts.createPhysicalTable(tx, table); err != nil {
		tx.Rollback()
		return nil, err
	}

	return table, tx.Commit().Error
}

// GetDynamicTableList 获取动态表列表
func (dts *DynamicTableService) GetDynamicTableList(page, pageSize int, keyword string) ([]model.DynamicTable, int64, error) {
	var tables []model.DynamicTable
	var total int64

	// 构建查询条件
	db := global.DB.Model(&model.DynamicTable{})

	// 搜索条件
	if keyword != "" {
		db = db.Where("name LIKE ? OR display_name LIKE ? OR description LIKE ?",
			"%"+keyword+"%", "%"+keyword+"%", "%"+keyword+"%")
	}

	// 获取总数
	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 分页查询
	offset := (page - 1) * pageSize
	if err := db.Preload("FieldDefinitions").Offset(offset).Limit(pageSize).Find(&tables).Error; err != nil {
		return nil, 0, err
	}

	return tables, total, nil
}

// GetDynamicTableByID 根据ID获取动态表
func (dts *DynamicTableService) GetDynamicTableByID(id uint) (*model.DynamicTable, error) {
	var table model.DynamicTable
	if err := global.DB.Preload("FieldDefinitions").First(&table, id).Error; err != nil {
		return nil, err
	}
	return &table, nil
}

// GetDynamicTableByName 根据表名获取动态表
func (dts *DynamicTableService) GetDynamicTableByName(tableName string) (*model.DynamicTable, error) {
	var table model.DynamicTable
	if err := global.DB.Preload("FieldDefinitions").Where("table_name = ?", tableName).First(&table).Error; err != nil {
		return nil, err
	}
	return &table, nil
}

// UpdateDynamicTable 更新动态表
func (dts *DynamicTableService) UpdateDynamicTable(table *model.DynamicTable) (*model.DynamicTable, error) {
	var existingTable model.DynamicTable
	if err := global.DB.Preload("FieldDefinitions").First(&existingTable, table.ID).Error; err != nil {
		return nil, err
	}

	// 检查名称是否已存在（排除当前记录）
	var count int64
	if err := global.DB.Model(&model.DynamicTable{}).
		Where("name = ? AND id != ?", table.Name, table.ID).
		Count(&count).Error; err != nil {
		return nil, err
	}
	if count > 0 {
		return nil, errors.New("表名称已存在")
	}

	// 检查表名是否已存在（排除当前记录）
	if err := global.DB.Model(&model.DynamicTable{}).
		Where("table_name = ? AND id != ?", table.TableName, table.ID).
		Count(&count).Error; err != nil {
		return nil, err
	}
	if count > 0 {
		return nil, errors.New("表名已存在")
	}

	// 开启事务
	tx := global.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// 更新表记录 - 排除时间戳字段
	updateData := map[string]interface{}{
		"name":         table.Name,
		"display_name": table.DisplayName,
		"description":  table.Description,
		"status":       table.Status,
		"updated_at":   time.Now(),
	}
	if err := tx.Model(&model.DynamicTable{}).Where("id = ?", table.ID).Updates(updateData).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	// 如果字段有变化，更新物理表结构
	if err := dts.updatePhysicalTable(tx, &existingTable, table); err != nil {
		tx.Rollback()
		return nil, err
	}

	return table, tx.Commit().Error
}

// DeleteTable 删除动态表
func (dts *DynamicTableService) DeleteTable(id uint) error {
	// 获取表信息
	table, err := dts.GetTableByID(id)
	if err != nil {
		return fmt.Errorf("获取表信息失败: %v", err)
	}

	// 开启事务
	tx := global.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			fmt.Printf("删除表事务回滚: %v\n", r)
		}
	}()

	// 记录开始删除
	fmt.Printf("开始删除表: %s (ID: %d)\n", table.TableName, id)

	// 先删除相关记录（从子表开始删除，避免外键约束问题）
	
	// 1. 删除相关视图记录
	if err := tx.Where("table_id = ?", id).Delete(&model.DynamicView{}).Error; err != nil {
		tx.Rollback()
		fmt.Printf("删除视图记录失败: %v\n", err)
		return fmt.Errorf("删除视图记录失败: %v", err)
	}
	fmt.Printf("已删除表 %d 的相关视图记录\n", id)

	// 2. 删除相关权限记录
	if err := tx.Where("table_id = ?", id).Delete(&model.TablePermission{}).Error; err != nil {
		tx.Rollback()
		fmt.Printf("删除权限记录失败: %v\n", err)
		return fmt.Errorf("删除权限记录失败: %v", err)
	}
	fmt.Printf("已删除表 %d 的相关权限记录\n", id)

	// 3. 删除相关字段记录
	if err := tx.Where("table_id = ?", id).Delete(&model.DynamicField{}).Error; err != nil {
		tx.Rollback()
		fmt.Printf("删除字段记录失败: %v\n", err)
		return fmt.Errorf("删除字段记录失败: %v", err)
	}
	fmt.Printf("已删除表 %d 的相关字段记录\n", id)

	// 4. 删除表记录
	if err := tx.Delete(&model.DynamicTable{}, id).Error; err != nil {
		tx.Rollback()
		fmt.Printf("删除表记录失败: %v\n", err)
		return fmt.Errorf("删除表记录失败: %v", err)
	}
	fmt.Printf("已删除表记录: %d\n", id)

	// 5. 最后删除物理表
	if err := dts.dropPhysicalTable(tx, table.TableName); err != nil {
		tx.Rollback()
		fmt.Printf("删除物理表失败: %v\n", err)
		return fmt.Errorf("删除物理表失败: %v", err)
	}
	fmt.Printf("已删除物理表: %s\n", table.TableName)

	// 提交事务
	if err := tx.Commit().Error; err != nil {
		fmt.Printf("提交事务失败: %v\n", err)
		return fmt.Errorf("提交事务失败: %v", err)
	}

	fmt.Printf("表 %s (ID: %d) 删除成功\n", table.TableName, id)
	return nil
}

// ToggleTableStatus 切换表状态
func (dts *DynamicTableService) ToggleTableStatus(id uint) error {
	var table model.DynamicTable
	if err := global.DB.First(&table, id).Error; err != nil {
		return err
	}

	// 切换状态
	if table.Status == 1 {
		table.Status = 2
	} else {
		table.Status = 1
	}

	return global.DB.Save(&table).Error
}

// createPhysicalTable 创建物理表
func (dts *DynamicTableService) createPhysicalTable(tx *gorm.DB, table *model.DynamicTable) error {
	// 构建CREATE TABLE语句
	var columns []string

	// 添加基础字段
	columns = append(columns, "id BIGINT AUTO_INCREMENT PRIMARY KEY")
	columns = append(columns, "created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP")
	columns = append(columns, "updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP")
	columns = append(columns, "deleted_at TIMESTAMP NULL DEFAULT NULL")
	columns = append(columns, "tenant_id BIGINT DEFAULT 0")

	// 添加动态字段（如果有的话）
	if table.FieldDefinitions != nil {
		for _, field := range table.FieldDefinitions {
			if field.Status == 1 { // 只添加启用状态的字段
				columnDef := dts.buildColumnDefinition(&field)
				columns = append(columns, columnDef)
			}
		}
	}

	// 添加索引
	columns = append(columns, "INDEX idx_deleted_at (deleted_at)")
	columns = append(columns, "INDEX idx_tenant_id (tenant_id)")
	columns = append(columns, "INDEX idx_created_at (created_at)")

	// 确保表名安全
	safeTableName := SanitizeTableName(table.TableName)
	sql := fmt.Sprintf("CREATE TABLE IF NOT EXISTS `%s` (%s) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci", safeTableName, strings.Join(columns, ", "))

	return tx.Exec(sql).Error
}

// sanitizeTableName 清理表名，防止SQL注入
func SanitizeTableName(tableName string) string {
	// 只允许字母、数字、下划线，且不能以数字开头
	reg := regexp.MustCompile(`[^a-zA-Z0-9_]`)
	safeName := reg.ReplaceAllString(tableName, "")
	
	// 确保不以数字开头
	if len(safeName) > 0 && safeName[0] >= '0' && safeName[0] <= '9' {
		safeName = "_" + safeName
	}
	
	// 限制长度
	if len(safeName) > 64 {
		safeName = safeName[:64]
	}
	
	return safeName
}

// updatePhysicalTable 更新物理表结构
func (dts *DynamicTableService) updatePhysicalTable(tx *gorm.DB, oldTable, newTable *model.DynamicTable) error {
	// 这里简化处理，实际项目中需要比较字段差异并执行相应的ALTER TABLE语句
	// 为了安全起见，暂时不支持修改表结构
	return nil
}

// dropPhysicalTable 删除物理表
func (dts *DynamicTableService) dropPhysicalTable(tx *gorm.DB, tableName string) error {
	sql := fmt.Sprintf("DROP TABLE IF EXISTS %s", tableName)
	return tx.Exec(sql).Error
}

// buildColumnDefinition 构建列定义
func (dts *DynamicTableService) buildColumnDefinition(field *model.DynamicField) string {
	columnType := field.GetMySQLColumnType()

	definition := fmt.Sprintf("%s %s", field.FieldName, columnType)

	if field.IsRequired {
		definition += " NOT NULL"
	} else {
		definition += " NULL"
	}

	if field.DefaultValue != "" {
		definition += fmt.Sprintf(" DEFAULT '%s'", field.DefaultValue)
	}

	if field.IsUnique {
		definition += " UNIQUE"
	}

	return definition
}

// GetTableSchema 获取表结构信息
func (dts *DynamicTableService) GetTableSchema(tableName string) (map[string]interface{}, error) {
	// 获取表信息
	table, err := dts.GetTableByName(tableName)
	if err != nil {
		return nil, err
	}

	// 构建表结构信息
	schema := map[string]interface{}{
		"table_info": table,
		"fields":     table.FieldDefinitions,
	}

	return schema, nil
}

// ValidateTableName 验证表名
func (dts *DynamicTableService) ValidateTableName(tableName string) (bool, error) {
	// 检查表名格式
	if tableName == "" {
		return false, errors.New("表名不能为空")
	}

	// 检查表名长度
	if len(tableName) > 64 {
		return false, errors.New("表名长度不能超过64个字符")
	}

	// 检查表名格式（只允许字母、数字、下划线，且不能以数字开头）
	matched, _ := regexp.MatchString("^[a-zA-Z_][a-zA-Z0-9_]*$", tableName)
	if !matched {
		return false, errors.New("表名只能包含字母、数字、下划线，且不能以数字开头")
	}

	// 检查是否为保留字
	reservedWords := []string{
		"select", "insert", "update", "delete", "create", "drop", "alter", "table",
		"database", "index", "view", "trigger", "procedure", "function", "user",
		"group", "order", "by", "where", "having", "join", "union", "distinct",
	}

	lowerTableName := strings.ToLower(tableName)
	for _, word := range reservedWords {
		if lowerTableName == word {
			return false, errors.New("表名不能使用SQL保留字")
		}
	}

	// 检查表名是否已存在
	var count int64
	if err := global.DB.Model(&model.DynamicTable{}).Where("table_name = ?", tableName).Count(&count).Error; err != nil {
		return false, err
	}

	if count > 0 {
		return false, errors.New("表名已存在")
	}

	return true, nil
}

// isValidTableName 检查表名是否有效
func isValidTableName(name string) bool {
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

// isReservedWord 检查是否为保留字
func isReservedWord(word string) bool {
	reservedWords := []string{
		"select", "insert", "update", "delete", "create", "drop", "alter", "table",
		"database", "index", "view", "trigger", "procedure", "function", "user",
		"group", "order", "by", "where", "having", "limit", "offset", "join",
		"inner", "left", "right", "outer", "union", "distinct", "count", "sum",
		"avg", "max", "min", "and", "or", "not", "null", "true", "false",
	}

	lowerWord := strings.ToLower(word)
	for _, reserved := range reservedWords {
		if lowerWord == reserved {
			return true
		}
	}

	return false
}

// UpdateTableStatus 更新表状态
func (dts *DynamicTableService) UpdateTableStatus(id uint, status int) error {
	var table model.DynamicTable
	if err := global.DB.First(&table, id).Error; err != nil {
		return err
	}

	table.Status = status
	return global.DB.Save(&table).Error
}
