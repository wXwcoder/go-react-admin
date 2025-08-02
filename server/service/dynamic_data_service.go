package service

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"go-react-admin/global"
	"go-react-admin/model"
)

type DynamicDataService struct{}

// CreateData 创建动态数据
func (dds *DynamicDataService) CreateData(tableName string, data map[string]interface{}) (map[string]interface{}, error) {
	// 获取表定义
	table, err := (&DynamicTableService{}).GetTableByName(tableName)
	if err != nil {
		return nil, fmt.Errorf("表不存在: %v", err)
	}

	// 检查物理表是否存在，如果不存在则创建
	if err := dds.ensurePhysicalTableExists(table); err != nil {
		return nil, fmt.Errorf("确保物理表存在失败: %v", err)
	}

	// 验证数据
	if err := dds.validateData(table, data); err != nil {
		return nil, err
	}

	// 处理数据
	processedData := dds.processDataForInsert(table, data)

	// 构建插入SQL
	columns := []string{"created_at", "updated_at"}
	values := []interface{}{time.Now(), time.Now()}

	for key, value := range processedData {
		columns = append(columns, key)
		values = append(values, value)
	}

	// 构建SQL语句
	placeholders := strings.Repeat("?,", len(values))
	placeholders = placeholders[:len(placeholders)-1]

	sql := fmt.Sprintf("INSERT INTO %s (%s) VALUES (%s)",
		tableName, strings.Join(columns, ","), placeholders)

	// 执行插入
	result := global.DB.Exec(sql, values...)
	if result.Error != nil {
		return nil, result.Error
	}

	// 获取插入的ID
	var insertID int64
	global.DB.Raw("SELECT LAST_INSERT_ID()").Scan(&insertID)

	// 返回创建的数据
	return dds.GetDataByID(tableName, uint(insertID))
}

// GetDataList 获取动态数据列表
func (dds *DynamicDataService) GetDataList(tableName string, page, pageSize int, filters map[string]interface{}, orderBy string) ([]map[string]interface{}, int64, error) {
	// 构建查询
	db := global.DB.Table(tableName)

	// 添加软删除条件
	db = db.Where("deleted_at IS NULL")

	// 添加过滤条件
	for field, value := range filters {
		if value != nil && value != "" && field != "_t" {
			db = db.Where(fmt.Sprintf("%s = ?", field), value)
		}
	}

	// 获取总数
	var total int64
	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 添加排序
	if orderBy != "" {
		db = db.Order(orderBy)
	} else {
		db = db.Order("id DESC")
	}

	// 分页
	if page > 0 && pageSize > 0 {
		offset := (page - 1) * pageSize
		db = db.Offset(offset).Limit(pageSize)
	}

	// 构建查询SQL
	var whereClauses []string
	var whereValues []interface{}
	
	whereClauses = append(whereClauses, "deleted_at IS NULL")
	
	for field, value := range filters {
		if value != nil && value != "" && field != "_t" {
			whereClauses = append(whereClauses, fmt.Sprintf("%s = ?", field))
			whereValues = append(whereValues, value)
		}
	}
	
	whereClause := strings.Join(whereClauses, " AND ")
	
	// 构建排序
	orderClause := "id DESC"
	if orderBy != "" {
		orderClause = orderBy
	}
	
	// 构建分页
	limitClause := ""
	if page > 0 && pageSize > 0 {
		offset := (page - 1) * pageSize
		limitClause = fmt.Sprintf(" LIMIT %d OFFSET %d", pageSize, offset)
	}
	
	sql := fmt.Sprintf("SELECT * FROM %s WHERE %s ORDER BY %s%s", tableName, whereClause, orderClause, limitClause)
	
	rows, err := global.DB.Raw(sql, whereValues...).Rows()
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	columns, err := rows.Columns()
	if err != nil {
		return nil, 0, err
	}

	var results []map[string]interface{}
	for rows.Next() {
		values := make([]interface{}, len(columns))
		valuePtrs := make([]interface{}, len(columns))
		for i := range columns {
			valuePtrs[i] = &values[i]
		}

		if err := rows.Scan(valuePtrs...); err != nil {
			return nil, 0, err
		}

		result := make(map[string]interface{})
		for i, col := range columns {
			val := values[i]
			if b, ok := val.([]byte); ok {
				result[col] = string(b)
			} else {
				result[col] = val
			}
		}
		results = append(results, result)
	}

	return results, total, nil
}

// GetDataByID 根据ID获取动态数据
func (dds *DynamicDataService) GetDataByID(tableName string, id uint) (map[string]interface{}, error) {
	sql := fmt.Sprintf("SELECT * FROM %s WHERE id = ? AND deleted_at IS NULL LIMIT 1", tableName)
	
	rows, err := global.DB.Raw(sql, id).Rows()
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	if !rows.Next() {
		return nil, errors.New("记录不存在")
	}

	columns, err := rows.Columns()
	if err != nil {
		return nil, err
	}

	values := make([]interface{}, len(columns))
	valuePtrs := make([]interface{}, len(columns))
	for i := range columns {
		valuePtrs[i] = &values[i]
	}

	if err := rows.Scan(valuePtrs...); err != nil {
		return nil, err
	}

	result := make(map[string]interface{})
	for i, col := range columns {
		val := values[i]
		if b, ok := val.([]byte); ok {
			result[col] = string(b)
		} else {
			result[col] = val
		}
	}

	return result, nil
}

// UpdateData 更新动态数据
func (dds *DynamicDataService) UpdateData(tableName string, id uint, data map[string]interface{}) (map[string]interface{}, error) {
	// 检查数据是否存在
	var count int64
	global.DB.Table(tableName).Where("id = ? AND deleted_at IS NULL", id).Count(&count)
	if count == 0 {
		return nil, errors.New("数据不存在")
	}

	// 添加更新时间
	data["updated_at"] = time.Now()

	// 执行更新
	if err := global.DB.Table(tableName).Where("id = ?", id).Updates(data).Error; err != nil {
		return nil, err
	}

	// 返回更新后的数据
	return dds.GetDataByID(tableName, id)
}

// DeleteData 删除动态数据（软删除）
func (dds *DynamicDataService) DeleteData(tableName string, id uint) error {
	// 检查数据是否存在
	var count int64
	global.DB.Table(tableName).Where("id = ? AND deleted_at IS NULL", id).Count(&count)
	if count == 0 {
		return errors.New("数据不存在")
	}

	// 软删除
	return global.DB.Table(tableName).Where("id = ?", id).Update("deleted_at", time.Now()).Error
}

// BatchDeleteData 批量删除动态数据
func (dds *DynamicDataService) BatchDeleteData(tableName string, ids []uint) error {
	if len(ids) == 0 {
		return nil
	}

	// 软删除
	return global.DB.Table(tableName).Where("id IN ?", ids).Update("deleted_at", time.Now()).Error
}

// GetDataStatistics 获取动态数据统计信息
func (dds *DynamicDataService) GetDataStatistics(tableName string) (map[string]interface{}, error) {
	stats := make(map[string]interface{})

	// 总记录数
	var totalCount int64
	global.DB.Table(tableName).Where("deleted_at IS NULL").Count(&totalCount)
	stats["total"] = totalCount

	// 今日新增
	var todayCount int64
	today := time.Now().Format("2006-01-02")
	global.DB.Table(tableName).Where("deleted_at IS NULL AND DATE(created_at) = ?", today).Count(&todayCount)
	stats["today"] = todayCount

	// 本周新增
	var weekCount int64
	weekStart := time.Now().AddDate(0, 0, -int(time.Now().Weekday())).Format("2006-01-02")
	global.DB.Table(tableName).Where("deleted_at IS NULL AND DATE(created_at) >= ?", weekStart).Count(&weekCount)
	stats["week"] = weekCount

	// 本月新增
	var monthCount int64
	monthStart := time.Now().AddDate(0, 0, -time.Now().Day()+1).Format("2006-01-02")
	global.DB.Table(tableName).Where("deleted_at IS NULL AND DATE(created_at) >= ?", monthStart).Count(&monthCount)
	stats["month"] = monthCount

	return stats, nil
}

// validateData 验证数据
func (dds *DynamicDataService) validateData(table *model.DynamicTable, data map[string]interface{}) error {
	// 获取字段定义
	fields := table.FieldDefinitions
	if fields == nil {
		return nil
	}

	// 遍历字段定义进行验证
	for _, field := range fields {
		if field.Status != 1 {
			continue // 跳过禁用字段
		}

		value, exists := data[field.FieldName]

		// 检查必填字段
		if field.IsRequired && (!exists || value == nil || value == "") {
			return fmt.Errorf("字段 %s 是必填的", field.DisplayName)
		}

		// 如果字段不存在且不是必填的，跳过验证
		if !exists {
			continue
		}

		// 验证字段值
		if err := dds.validateFieldValue(&field, value); err != nil {
			return fmt.Errorf("字段 %s 验证失败: %v", field.FieldName, err)
		}
	}

	return nil
}

// validateFieldValue 验证单个字段值
func (dds *DynamicDataService) validateFieldValue(field *model.DynamicField, value interface{}) error {
	if value == nil {
		return nil
	}

	strValue := fmt.Sprintf("%v", value)

	// 根据字段类型进行验证
	switch field.FieldType {
	case "string", "text":
		// 字符串长度验证
		validation, _ := field.GetValidation()
		if validation.MaxLength > 0 && len(strValue) > validation.MaxLength {
			return fmt.Errorf("长度不能超过 %d 个字符", validation.MaxLength)
		}
		if validation.MinLength > 0 && len(strValue) < validation.MinLength {
			return fmt.Errorf("长度不能少于 %d 个字符", validation.MinLength)
		}
	case "int":
		// 整数验证
		if _, err := strconv.ParseInt(strValue, 10, 64); err != nil {
			return errors.New("必须是整数")
		}
		validation, _ := field.GetValidation()
		if validation.MaxValue != 0 {
			if val, _ := strconv.ParseFloat(strValue, 64); val > validation.MaxValue {
				return fmt.Errorf("值不能大于 %f", validation.MaxValue)
			}
		}
		if validation.MinValue != 0 {
			if val, _ := strconv.ParseFloat(strValue, 64); val < validation.MinValue {
				return fmt.Errorf("值不能小于 %f", validation.MinValue)
			}
		}
	case "float":
		// 浮点数验证
		if _, err := strconv.ParseFloat(strValue, 64); err != nil {
			return errors.New("必须是数字")
		}
		validation, _ := field.GetValidation()
		if validation.MaxValue != 0 {
			if val, _ := strconv.ParseFloat(strValue, 64); val > validation.MaxValue {
				return fmt.Errorf("值不能大于 %f", validation.MaxValue)
			}
		}
		if validation.MinValue != 0 {
			if val, _ := strconv.ParseFloat(strValue, 64); val < validation.MinValue {
				return fmt.Errorf("值不能小于 %f", validation.MinValue)
			}
		}
	case "date":
		// 日期验证
		if _, err := time.Parse("2006-01-02", strValue); err != nil {
			return errors.New("日期格式不正确，应为 YYYY-MM-DD")
		}
	case "datetime":
		// 日期时间验证
		if _, err := time.Parse("2006-01-02 15:04:05", strValue); err != nil {
			return errors.New("日期时间格式不正确，应为 YYYY-MM-DD HH:MM:SS")
		}
	case "boolean":
		// 布尔值验证
		if strValue != "true" && strValue != "false" && strValue != "1" && strValue != "0" {
			return errors.New("必须是布尔值")
		}
	case "select", "multiselect":
		// 选项验证
		options, _ := field.GetOptions()
		validValues := make(map[string]bool)
		for _, opt := range options {
			validValues[opt.Value] = true
		}
		
		if field.FieldType == "select" {
			if !validValues[strValue] {
				return errors.New("无效的选择值")
			}
		} else {
			// 多选验证
			values := strings.Split(strValue, ",")
			for _, v := range values {
				if !validValues[strings.TrimSpace(v)] {
					return fmt.Errorf("无效的选择值: %s", v)
				}
			}
		}
	case "file", "image":
		// 文件路径验证
		if len(strValue) > 500 {
			return errors.New("文件路径过长")
		}
	}

	// 唯一性验证
	if field.IsUnique {
		// 这里需要检查数据库中是否已存在相同值
		// 实际实现时需要根据具体情况处理
	}

	return nil
}

// processDataForInsert 处理插入数据
func (dds *DynamicDataService) processDataForInsert(table *model.DynamicTable, data map[string]interface{}) map[string]interface{} {
	processed := make(map[string]interface{})
	
	for key, value := range data {
		// 处理特殊类型的数据
		if value != nil {
			switch v := value.(type) {
			case []interface{}:
				// 数组类型转为JSON字符串
				if jsonBytes, err := json.Marshal(v); err == nil {
					processed[key] = string(jsonBytes)
				} else {
					processed[key] = value
				}
			case map[string]interface{}:
				// 对象类型转为JSON字符串
				if jsonBytes, err := json.Marshal(v); err == nil {
					processed[key] = string(jsonBytes)
				} else {
					processed[key] = value
				}
			default:
				processed[key] = value
			}
		} else {
			processed[key] = value
		}
	}
	
	return processed
}

// processDataForUpdate 处理更新数据
func (dds *DynamicDataService) processDataForUpdate(table *model.DynamicTable, data map[string]interface{}) map[string]interface{} {
	return dds.processDataForInsert(table, data)
}

// ensurePhysicalTableExists 确保物理表存在
func (dds *DynamicDataService) ensurePhysicalTableExists(table *model.DynamicTable) error {
	// 检查表是否存在
	var count int64
	checkSQL := fmt.Sprintf("SELECT COUNT(*) FROM information_schema.tables WHERE table_schema = DATABASE() AND table_name = '%s'", table.TableName)
	if err := global.DB.Raw(checkSQL).Scan(&count).Error; err != nil {
		return fmt.Errorf("检查表存在性失败: %v", err)
	}

	// 如果表不存在，创建它
	if count == 0 {
		if err := dds.createPhysicalTable(table); err != nil {
			return fmt.Errorf("创建物理表失败: %v", err)
		}
	}

	return nil
}

// createPhysicalTable 创建物理表
func (dds *DynamicDataService) createPhysicalTable(table *model.DynamicTable) error {
	// 构建CREATE TABLE语句
	var columns []string

	// 添加基础字段
	columns = append(columns, "id INT AUTO_INCREMENT PRIMARY KEY")
	columns = append(columns, "created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP")
	columns = append(columns, "updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP")
	columns = append(columns, "deleted_at TIMESTAMP NULL")
	columns = append(columns, "tenant_id INT DEFAULT 0")

	// 添加动态字段（如果有的话）
	if table.FieldDefinitions != nil {
		for _, field := range table.FieldDefinitions {
			columnDef := dds.buildColumnDefinition(&field)
			columns = append(columns, columnDef)
		}
	}

	// 添加索引
	columns = append(columns, "INDEX idx_deleted_at (deleted_at)")
	columns = append(columns, "INDEX idx_tenant_id (tenant_id)")

	// 确保表名安全
	safeTableName := strings.ReplaceAll(table.TableName, "`", "")
	sql := fmt.Sprintf("CREATE TABLE `%s` (%s)", safeTableName, strings.Join(columns, ", "))

	return global.DB.Exec(sql).Error
}

// buildColumnDefinition 构建列定义
func (dds *DynamicDataService) buildColumnDefinition(field *model.DynamicField) string {
	columnType := field.GetMySQLColumnType()

	definition := fmt.Sprintf("`%s` %s", field.FieldName, columnType)

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

// CreateView 创建数据视图
func (dds *DynamicDataService) CreateView(view *model.DynamicView) error {
	return global.DB.Create(view).Error
}

// GetViewList 获取视图列表
func (dds *DynamicDataService) GetViewList(tableID uint) ([]model.DynamicView, error) {
	var views []model.DynamicView
	err := global.DB.Where("table_id = ?", tableID).Order("sort_order ASC").Find(&views).Error
	return views, err
}

// GetViewByID 根据ID获取视图
func (dds *DynamicDataService) GetViewByID(id uint) (*model.DynamicView, error) {
	var view model.DynamicView
	err := global.DB.First(&view, id).Error
	if err != nil {
		return nil, err
	}
	return &view, nil
}

// UpdateView 更新视图
func (dds *DynamicDataService) UpdateView(view *model.DynamicView) error {
	return global.DB.Save(view).Error
}

// DeleteView 删除视图
func (dds *DynamicDataService) DeleteView(id uint) error {
	return global.DB.Delete(&model.DynamicView{}, id).Error
}

// ApplyView 应用视图
func (dds *DynamicDataService) ApplyView(viewID uint, params map[string]interface{}) (interface{}, error) {
	// 获取视图配置
	view, err := dds.GetViewByID(viewID)
	if err != nil {
		return nil, fmt.Errorf("视图不存在: %v", err)
	}

	// 解析视图配置
	config, err := view.GetViewConfig()
	if err != nil {
		return nil, fmt.Errorf("解析视图配置失败: %v", err)
	}

	// 获取表定义
	table, err := (&DynamicTableService{}).GetTableByID(view.TableID)
	if err != nil {
		return nil, fmt.Errorf("获取表定义失败: %v", err)
	}

	// 根据视图配置构建查询
	filters := make(map[string]interface{})
	if config.Filters != nil {
		for field, value := range config.Filters {
			if paramValue, ok := params[field]; ok {
				filters[field] = paramValue
			} else {
				filters[field] = value
			}
		}
	}

	// 获取排序配置
	orderBy := "id DESC"
	if config.Sort.Field != "" {
		orderBy = fmt.Sprintf("%s %s", config.Sort.Field, config.Sort.Order)
	}

	// 获取分页配置
	page := 1
	pageSize := 10
	if p, ok := params["page"]; ok {
		if pFloat, ok := p.(float64); ok {
			page = int(pFloat)
		}
	}
	if ps, ok := params["pageSize"]; ok {
		if psFloat, ok := ps.(float64); ok {
			pageSize = int(psFloat)
		}
	}

	// 执行查询
	data, total, err := dds.GetDataList(table.TableName, page, pageSize, filters, orderBy)
	if err != nil {
		return nil, fmt.Errorf("查询数据失败: %v", err)
	}

	// 返回结果
	return map[string]interface{}{
		"data":     data,
		"total":    total,
		"page":     page,
		"pageSize": pageSize,
		"config":   config,
	}, nil
}
