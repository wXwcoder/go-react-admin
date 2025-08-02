package initialize

import (
	"fmt"
	"go-react-admin/global"
	"go-react-admin/model"
	"log"

	"gorm.io/gorm"
)

// InitDynamicTables 初始化动态数据管理平台相关表
func InitDynamicTables() {
	db := global.DB
	
	// 自动迁移动态数据管理平台相关表
	err := db.AutoMigrate(
		&model.DynamicTable{},
		&model.DynamicField{},
		&model.TablePermission{},
		&model.DynamicView{},
		&model.DynamicImportExportLog{},
	)
	if err != nil {
		log.Fatalf("动态数据管理平台表迁移失败: %v", err)
	}
	
	log.Println("动态数据管理平台表迁移成功")
	
	// 创建默认数据
	createDefaultData(db)
}

// createDefaultData 创建默认数据
func createDefaultData(db *gorm.DB) {
	// 检查是否已有数据，避免重复创建
	var count int64
	db.Model(&model.DynamicTable{}).Count(&count)
	if count > 0 {
		return
	}
	
	// 创建示例动态表
	sampleTable := &model.DynamicTable{
		Name:        "sample_users",
		DisplayName: "示例用户表",
		Description: "这是一个示例用户表，展示动态数据管理平台的功能",
		TableName:   "dyn_sample_users",
		Fields:      []byte(`[{"name":"name","display_name":"姓名","type":"string","required":true},{"name":"email","display_name":"邮箱","type":"string","required":true,"unique":true},{"name":"age","display_name":"年龄","type":"int"},{"name":"status","display_name":"状态","type":"select","options":[{"label":"激活","value":"active"},{"label":"禁用","value":"inactive"}]}]`),
		Status:      1,
		TenantID:    1,
	}
	
	if err := db.Create(sampleTable).Error; err != nil {
		log.Printf("创建示例表失败: %v", err)
		return
	}
	
	// 创建示例字段
	sampleFields := []model.DynamicField{
		{
			TableID:      sampleTable.ID,
			FieldName:    "name",
			DisplayName:  "姓名",
			FieldType:    "string",
			IsRequired:   true,
			IsSearchable: true,
			IsSortable:   true,
			SortOrder:    1,
			Status:       1,
		},
		{
			TableID:      sampleTable.ID,
			FieldName:    "email",
			DisplayName:  "邮箱",
			FieldType:    "string",
			IsRequired:   true,
			IsUnique:     true,
			IsSearchable: true,
			SortOrder:    2,
			Status:       1,
		},
		{
			TableID:     sampleTable.ID,
			FieldName:   "age",
			DisplayName: "年龄",
			FieldType:   "int",
			IsSortable:  true,
			SortOrder:   3,
			Status:      1,
		},
		{
			TableID:     sampleTable.ID,
			FieldName:   "status",
			DisplayName: "状态",
			FieldType:   "select",
			IsRequired:  true,
			Options:     []byte(`[{"label":"激活","value":"active"},{"label":"禁用","value":"inactive"}]`),
			SortOrder:   4,
			Status:      1,
		},
	}
	
	for _, field := range sampleFields {
		if err := db.Create(&field).Error; err != nil {
			log.Printf("创建示例字段失败: %v", err)
		}
	}
	
	log.Println("示例数据创建成功")
}

// CreateDynamicTable 根据配置创建物理数据表
func CreateDynamicTable(tableName string, fields []model.DynamicField) error {
	db := global.DB
	
	// 构建CREATE TABLE语句
	sql := fmt.Sprintf("CREATE TABLE IF NOT EXISTS `%s` (", tableName)
	sql += "`id` bigint unsigned NOT NULL AUTO_INCREMENT,"
	sql += "`created_at` datetime(3) DEFAULT NULL,"
	sql += "`updated_at` datetime(3) DEFAULT NULL,"
	sql += "`deleted_at` datetime(3) DEFAULT NULL,"
	sql += "`tenant_id` bigint unsigned DEFAULT NULL,"
	
	// 添加自定义字段
	for _, field := range fields {
		columnType := getColumnType(field.FieldType)
		nullable := "NULL"
		if field.IsRequired {
			nullable = "NOT NULL"
		}
		
		defaultValue := ""
		if field.DefaultValue != "" {
			defaultValue = fmt.Sprintf(" DEFAULT '%s'", field.DefaultValue)
		}
		
		sql += fmt.Sprintf("`%s` %s %s%s COMMENT '%s',", 
			field.FieldName, columnType, nullable, defaultValue, field.DisplayName)
	}
	
	sql += "PRIMARY KEY (`id`),"
	sql += "KEY `idx_deleted_at` (`deleted_at`),"
	sql += "KEY `idx_tenant_id` (`tenant_id`)"
	
	// 添加唯一索引
	for _, field := range fields {
		if field.IsUnique {
			sql += fmt.Sprintf(",UNIQUE KEY `idx_%s_%s` (`%s`)", tableName, field.FieldName, field.FieldName)
		}
	}
	
	sql += ") ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci"
	
	// 执行创建表语句
	return db.Exec(sql).Error
}

// getColumnType 根据字段类型返回MySQL列类型
func getColumnType(fieldType string) string {
	switch fieldType {
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

// DropDynamicTable 删除物理数据表
func DropDynamicTable(tableName string) error {
	db := global.DB
	sql := fmt.Sprintf("DROP TABLE IF EXISTS `%s`", tableName)
	return db.Exec(sql).Error
}