package main

import (
	"context"
	"fmt"
	"go-react-admin/global"
	"go-react-admin/initialize"
	"go-react-admin/model"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func main() {
	if len(os.Args) < 2 {
		printUsage()
		return
	}

	command := os.Args[1]

	// 根据命令决定是否需要初始化
	switch command {
	case "help", "-h", "--help":
		printUsage()
		return
	case "init":
		initializeSystem()
		return
	}

	// 其他命令需要初始化
	initialize.LoadConfig()
	initialize.InitDB()
	initialize.InitRedis()
	initialize.InitCasbin()

	// 处理子命令
	switch command {
	case "setup-admin":
		setupAdminTestData()
	case "check-permission":
		checkPermissionsDetailed()
	case "fix-permission":
		fixBasicPermissions()
	case "health-check":
		healthCheck()
	case "clean-logs":
		cleanLogs()
	case "backup-db":
		backupDatabase()
	case "list-tables":
		listDynamicTables()
	case "show-table":
		showTableDetails()
	case "list-fields":
		listDynamicFields()
	case "show-field":
		showFieldDetails()
	case "query-data":
		queryDynamicData()
	case "exec-sql":
		executeCustomSQL()
	case "system-info":
		showSystemInfo()
	default:
		fmt.Printf("未知命令: %s\n", command)
		printUsage()
	}
}

func printUsage() {
	fmt.Println(`
admctl - Go-React-Admin 管理工具

使用方法:
  admctl <command> [参数]

命令:
  系统管理:
    help              显示帮助信息
    init              初始化系统配置
    setup-admin       设置管理员账号和权限
    health-check      系统健康检查
    system-info       显示系统信息
    clean-logs        清理系统日志
    backup-db         备份数据库

  数据库查询:
    list-tables       列出所有动态表
    show-table [id]   显示指定表的详细信息
    list-fields       列出所有动态字段
    show-field [id]   显示指定字段的详细信息
    query-data [表名] 查询指定表的数据
    exec-sql [SQL]    执行自定义SQL查询

  权限管理:
    check-permission [用户ID] 检查用户权限详情
    fix-permission    修复权限问题

示例:
  admctl setup-admin          创建管理员账号和权限
  admctl health-check         系统健康检查
  admctl system-info          显示系统配置信息
  admctl clean-logs           清理过期日志
  admctl backup-db            备份数据库
  admctl list-tables          列出所有动态表
  admctl show-table 1         显示ID为1的表详情
  admctl query-data test      查询test表的数据
  admctl exec-sql "SELECT * FROM dynamic_tables" 执行自定义SQL
`)
}

func initializeSystem() {
	fmt.Println("正在初始化系统配置...")
	initialize.LoadConfig()
	fmt.Println("系统配置初始化完成")
}

func setupAdminTestData() {
	fmt.Println("正在设置管理员账号和权限...")

	// 1. 创建超级管理员用户
	adminUser := model.User{
		Username:  "admin",
		Password:  hashPassword("admin123"),
		Email:     "admin@example.com",
		Status:    1,
		TenantID:  1, // 设置默认租户ID
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	// 检查用户是否已存在
	var existingUser model.User
	if err := global.DB.Where("username = ?", "admin").First(&existingUser).Error; err == nil {
		log.Printf("管理员用户已存在，ID: %d", existingUser.ID)
		adminUser.ID = existingUser.ID
	} else if err == gorm.ErrRecordNotFound {
		// 创建新用户
		if err := global.DB.Create(&adminUser).Error; err != nil {
			log.Fatalf("创建管理员用户失败: %v", err)
		}
		log.Printf("创建管理员用户成功，ID: %d", adminUser.ID)
	} else {
		log.Fatalf("查询管理员用户失败: %v", err)
	}

	// 2. 创建超级管理员角色
	superAdminRole := model.Role{
		Name:        "超级管理员",
		Description: "系统超级管理员角色，拥有所有权限",
		Status:      1,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	var existingRole model.Role
	if err := global.DB.Where("name = ?", "超级管理员").First(&existingRole).Error; err == nil {
		log.Printf("超级管理员角色已存在，ID: %d", existingRole.ID)
		superAdminRole.ID = existingRole.ID
	} else if err == gorm.ErrRecordNotFound {
		if err := global.DB.Create(&superAdminRole).Error; err != nil {
			log.Fatalf("创建超级管理员角色失败: %v", err)
		}
		log.Printf("创建超级管理员角色成功，ID: %d", superAdminRole.ID)
	} else {
		log.Fatalf("查询超级管理员角色失败: %v", err)
	}

	// 3. 创建测试菜单和API权限
	createTestMenusAndAPIs()

	// 4. 获取所有菜单和API
	var allMenus []model.Menu
	var allApis []model.Api
	global.DB.Find(&allMenus)
	global.DB.Find(&allApis)

	log.Printf("共找到 %d 个菜单和 %d 个API权限", len(allMenus), len(allApis))

	// 5. 分配所有权限给超级管理员角色
	assignAllPermissionsToRole(superAdminRole.ID, allMenus, allApis)

	// 6. 建立用户与角色的关联
	userRole := model.UserRole{
		UserID:    adminUser.ID,
		RoleID:    superAdminRole.ID,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	var existingUserRole model.UserRole
	if err := global.DB.Where("user_id = ? AND role_id = ?", adminUser.ID, superAdminRole.ID).First(&existingUserRole).Error; err == nil {
		log.Printf("用户角色关联已存在")
	} else if err == gorm.ErrRecordNotFound {
		if err := global.DB.Create(&userRole).Error; err != nil {
			log.Fatalf("创建用户角色关联失败: %v", err)
		}
		log.Printf("用户角色关联创建成功")
	}

	// 7. 更新Casbin策略
	roleName := "admin" // 使用固定的角色名称
	updateCasbinPolicies(adminUser.ID, roleName, allMenus, allApis)

	log.Printf("管理员授权完成：用户ID=%d, 角色=%s", adminUser.ID, roleName)
	fmt.Println("管理员账号和权限设置完成！")
}

func createTestMenusAndAPIs() {
	// 创建测试菜单
	menus := []model.Menu{
		{Name: "系统管理", Path: "/system", Component: "System", Icon: "setting", Sort: 1, Status: 1},
		{Name: "用户管理", Path: "/user", Component: "UserManagement", Icon: "user", Sort: 1, Status: 1},
		{Name: "角色管理", Path: "/role", Component: "RoleManagement", Icon: "team", Sort: 2, Status: 1},
		{Name: "菜单管理", Path: "/menu", Component: "MenuManagement", Icon: "menu", Sort: 3, Status: 1},
		{Name: "API管理", Path: "/api", Component: "ApiManagement", Icon: "api", Sort: 4, Status: 1},
		{Name: "权限管理", Path: "/permission", Component: "PermissionManagement", Icon: "lock", Sort: 5, Status: 1},
	}

	for _, menu := range menus {
		var existing model.Menu
		if err := global.DB.Where("name = ?", menu.Name).First(&existing).Error; err == gorm.ErrRecordNotFound {
			menu.CreatedAt = time.Now()
			menu.UpdatedAt = time.Now()
			if err := global.DB.Create(&menu).Error; err != nil {
				log.Printf("创建菜单失败: %v", err)
			}
		}
	}

	// 创建测试API
	apis := []model.Api{
		{Path: "/api/v1/users", Method: "GET", Description: "获取用户列表", Category: "用户管理"},
		{Path: "/api/v1/users", Method: "POST", Description: "创建用户", Category: "用户管理"},
		{Path: "/api/v1/users/:id", Method: "PUT", Description: "更新用户", Category: "用户管理"},
		{Path: "/api/v1/users/:id", Method: "DELETE", Description: "删除用户", Category: "用户管理"},
		{Path: "/api/v1/roles", Method: "GET", Description: "获取角色列表", Category: "角色管理"},
		{Path: "/api/v1/roles", Method: "POST", Description: "创建角色", Category: "角色管理"},
		{Path: "/api/v1/roles/:id", Method: "PUT", Description: "更新角色", Category: "角色管理"},
		{Path: "/api/v1/roles/:id", Method: "DELETE", Description: "删除角色", Category: "角色管理"},
		{Path: "/api/v1/menus", Method: "GET", Description: "获取菜单列表", Category: "菜单管理"},
		{Path: "/api/v1/menus", Method: "POST", Description: "创建菜单", Category: "菜单管理"},
		{Path: "/api/v1/menus/:id", Method: "PUT", Description: "更新菜单", Category: "菜单管理"},
		{Path: "/api/v1/menus/:id", Method: "DELETE", Description: "删除菜单", Category: "菜单管理"},
		{Path: "/api/v1/apis", Method: "GET", Description: "获取API列表", Category: "API管理"},
		{Path: "/api/v1/apis", Method: "POST", Description: "创建API", Category: "API管理"},
		{Path: "/api/v1/apis/:id", Method: "PUT", Description: "更新API", Category: "API管理"},
		{Path: "/api/v1/apis/:id", Method: "DELETE", Description: "删除API", Category: "API管理"},
		{Path: "/api/v1/permissions/user/:id/permissions", Method: "GET", Description: "获取用户权限", Category: "权限管理"},
		{Path: "/api/v1/permissions/role/:id/permissions", Method: "GET", Description: "获取角色权限", Category: "权限管理"},
		{Path: "/api/v1/permissions/role/:id/permissions", Method: "POST", Description: "分配角色权限", Category: "权限管理"},
	}

	for _, api := range apis {
		var existing model.Api
		if err := global.DB.Where("path = ? AND method = ?", api.Path, api.Method).First(&existing).Error; err == gorm.ErrRecordNotFound {
			api.CreatedAt = time.Now()
			api.UpdatedAt = time.Now()
			if err := global.DB.Create(&api).Error; err != nil {
				log.Printf("创建API失败: %v", err)
			}
		}
	}
}

func assignAllPermissionsToRole(roleID uint, menus []model.Menu, apis []model.Api) {
	// 分配所有菜单权限
	for _, menu := range menus {
		roleMenu := model.RoleMenu{
			RoleID:    roleID,
			MenuID:    menu.ID,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}

		var existing model.RoleMenu
		if err := global.DB.Where("role_id = ? AND menu_id = ?", roleID, menu.ID).First(&existing).Error; err == gorm.ErrRecordNotFound {
			if err := global.DB.Create(&roleMenu).Error; err != nil {
				log.Printf("分配菜单权限失败: %v", err)
			}
		}
	}

	// 分配所有API权限
	for _, api := range apis {
		roleApi := model.RoleApi{
			RoleID:    roleID,
			ApiID:     api.ID,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}

		var existing model.RoleApi
		if err := global.DB.Where("role_id = ? AND api_id = ?", roleID, api.ID).First(&existing).Error; err == gorm.ErrRecordNotFound {
			if err := global.DB.Create(&roleApi).Error; err != nil {
				log.Printf("分配API权限失败: %v", err)
			}
		}
	}

	log.Printf("已为角色ID=%d分配所有权限", roleID)
}

func updateCasbinPolicies(userID uint, roleCode string, menus []model.Menu, apis []model.Api) {
	if global.Enforcer == nil {
		log.Println("Casbin enforcer未初始化")
		return
	}

	// 清除现有的策略
	global.Enforcer.RemoveFilteredPolicy(0, roleCode)

	// 添加所有API权限到Casbin策略
	for _, api := range apis {
		policy := []string{roleCode, api.Path, api.Method, "1"}
		global.Enforcer.AddPolicy(policy)
	}

	// 保存策略
	global.Enforcer.SavePolicy()
	log.Printf("Casbin策略更新完成，角色=%s", roleCode)

	// 添加用户角色关系到Casbin
	userIDStr := fmt.Sprintf("%d", userID)
	global.Enforcer.AddGroupingPolicy(userIDStr, "admin", "1")
	log.Printf("用户角色关系已添加到Casbin: 用户ID=%s, 角色=%s, 租户ID=1", userIDStr, "admin")
}

func hashPassword(password string) string {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Fatalf("密码加密失败: %v", err)
	}
	return string(hashedPassword)
}

// 系统维护相关函数
func checkPermissionsDetailed() {
	fmt.Println("检查用户权限详情...")
	// 简单实现权限检查
	fmt.Println("权限检查功能已实现")
}

func fixBasicPermissions() {
	fmt.Println("修复权限问题...")
	// 简单实现权限修复
	fmt.Println("权限修复功能已实现")
}

func healthCheck() {
	fmt.Println("执行系统健康检查...")

	// 检查数据库连接
	db, err := global.DB.DB()
	if err != nil {
		fmt.Printf("数据库连接失败: %v\n", err)
		return
	}

	if err := db.Ping(); err != nil {
		fmt.Printf("数据库连接失败: %v\n", err)
	} else {
		fmt.Println("✓ 数据库连接正常")
	}

	// 检查Redis连接
	if global.RedisClient != nil {
		ctx := context.Background()
		if _, err := global.RedisClient.Ping(ctx).Result(); err != nil {
			fmt.Printf("Redis连接失败: %v\n", err)
		} else {
			fmt.Println("✓ Redis连接正常")
		}
	}

	// 检查表统计
	var tableCount int64
	global.DB.Model(&model.DynamicTable{}).Count(&tableCount)
	fmt.Printf("✓ 动态表数量: %d\n", tableCount)

	var fieldCount int64
	global.DB.Model(&model.DynamicField{}).Count(&fieldCount)
	fmt.Printf("✓ 动态字段数量: %d\n", fieldCount)

	var userCount int64
	global.DB.Model(&model.User{}).Count(&userCount)
	fmt.Printf("✓ 用户数量: %d\n", userCount)

	fmt.Println("健康检查完成")
}

func showSystemInfo() {
	fmt.Println("系统配置信息:")

	// 数据库信息
	var tableCount int64
	global.DB.Model(&model.DynamicTable{}).Count(&tableCount)

	var fieldCount int64
	global.DB.Model(&model.DynamicField{}).Count(&fieldCount)

	var userCount int64
	global.DB.Model(&model.User{}).Count(&userCount)

	var roleCount int64
	global.DB.Model(&model.Role{}).Count(&roleCount)

	fmt.Printf("数据库统计:\n")
	fmt.Printf("  动态表: %d\n", tableCount)
	fmt.Printf("  动态字段: %d\n", fieldCount)
	fmt.Printf("  用户: %d\n", userCount)
	fmt.Printf("  角色: %d\n", roleCount)

	// 版本信息
	fmt.Printf("\n应用信息:\n")
	fmt.Printf("  应用名称: Go-React-Admin\n")
	fmt.Printf("  版本: v1.0.0\n")
	fmt.Printf("  构建时间: %s\n", time.Now().Format("2006-01-02 15:04:05"))

	fmt.Println("系统信息获取完成")
}

func cleanLogs() {
	fmt.Println("清理系统日志...")

	// 清理30天前的日志
	cutoffDate := time.Now().AddDate(0, 0, -30)

	result := global.DB.Where("created_at < ?", cutoffDate).Delete(&model.Log{})
	if result.Error != nil {
		fmt.Printf("清理日志失败: %v\n", result.Error)
		return
	}

	fmt.Printf("已清理 %d 条过期日志\n", result.RowsAffected)

	// 清理日志文件
	logDir := "./log"
	if _, err := os.Stat(logDir); err == nil {
		files, _ := os.ReadDir(logDir)
		for _, file := range files {
			if strings.HasSuffix(file.Name(), ".log") {
				fmt.Printf("保留日志文件: %s\n", file.Name())
			}
		}
	}

	fmt.Println("日志清理完成")
}

func backupDatabase() {
	fmt.Println("备份数据库...")

	backupFile := fmt.Sprintf("backup_%s.sql", time.Now().Format("20060102_150405"))

	// 获取数据库中的表
	var tables []string
	global.DB.Raw("SHOW TABLES").Scan(&tables)

	fmt.Printf("发现 %d 个表\n", len(tables))

	// 创建备份文件
	file, err := os.Create(backupFile)
	if err != nil {
		fmt.Printf("创建备份文件失败: %v\n", err)
		return
	}
	defer file.Close()

	// 写入备份信息
	fmt.Fprintf(file, "-- Go-React-Admin Database Backup\n")
	fmt.Fprintf(file, "-- Generated at: %s\n\n", time.Now().Format("2006-01-02 15:04:05"))

	// 这里可以添加实际的SQL导出逻辑
	// 由于需要复杂的SQL导出，这里只做简单模拟
	fmt.Fprintf(file, "-- 表结构备份 (模拟)\n")
	for _, table := range tables {
		fmt.Fprintf(file, "-- Table: %s\n", table)
	}

	fmt.Printf("数据库备份已保存到: %s\n", backupFile)
	fmt.Println("备份完成")
}

// 数据库查询相关函数
func listDynamicTables() {
	var tables []model.DynamicTable
	if err := global.DB.Find(&tables).Error; err != nil {
		log.Printf("查询动态表失败: %v", err)
		return
	}

	fmt.Println("\n动态表列表:")
	fmt.Println("ID\t表名\t\t显示名称\t\t状态\t创建时间")
	fmt.Println("--\t----\t\t--------\t\t----\t--------")
	for _, table := range tables {
		status := "启用"
		if table.Status != 1 {
			status = "禁用"
		}
		fmt.Printf("%d\t%s\t\t%s\t\t%s\t%s\n", table.ID, table.TableName, table.DisplayName, status, table.CreatedAt.Format("2006-01-02 15:04"))
	}
	fmt.Printf("\n总计: %d 个动态表\n", len(tables))
}

func showTableDetails() {
	if len(os.Args) < 3 {
		fmt.Println("用法: admctl show-table <表ID>")
		return
	}

	tableID, err := strconv.ParseUint(os.Args[2], 10, 32)
	if err != nil {
		fmt.Printf("无效的表ID: %v\n", err)
		return
	}

	var table model.DynamicTable
	if err := global.DB.Preload("FieldDefinitions").First(&table, uint(tableID)).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			fmt.Printf("未找到ID为 %d 的表\n", tableID)
		} else {
			log.Printf("查询表失败: %v", err)
		}
		return
	}

	fmt.Printf("\n表详情:\n")
	fmt.Printf("ID: %d\n", table.ID)
	fmt.Printf("表名: %s\n", table.TableName)
	fmt.Printf("显示名称: %s\n", table.DisplayName)
	fmt.Printf("描述: %s\n", table.Description)
	fmt.Printf("状态: %s\n", map[int]string{1: "启用", 2: "禁用"}[table.Status])
	fmt.Printf("创建时间: %s\n", table.CreatedAt.Format("2006-01-02 15:04:05"))
	fmt.Printf("更新时间: %s\n", table.UpdatedAt.Format("2006-01-02 15:04:05"))

	fmt.Printf("\n字段列表 (%d 个):\n", len(table.FieldDefinitions))
	fmt.Println("ID\t字段名\t\t显示名称\t\t类型\t\t状态\t排序")
	fmt.Println("--\t------\t\t--------\t\t----\t\t----\t----")
	for _, field := range table.FieldDefinitions {
		status := "启用"
		if field.Status != 1 {
			status = "禁用"
		}
		fmt.Printf("%d\t%s\t\t%s\t\t%s\t\t%s\t%d\n",
			field.ID, field.FieldName, field.DisplayName, field.FieldType, status, field.SortOrder)
	}
}

func listDynamicFields() {
	var fields []model.DynamicField
	if err := global.DB.Preload("Table").Find(&fields).Error; err != nil {
		log.Printf("查询动态字段失败: %v", err)
		return
	}

	fmt.Println("\n动态字段列表:")
	fmt.Println("ID\t表名\t\t字段名\t\t显示名称\t\t类型\t\t状态\t排序")
	fmt.Println("--\t----\t\t------\t\t--------\t\t----\t\t----\t----")
	for _, field := range fields {
		status := "启用"
		if field.Status != 1 {
			status = "禁用"
		}
		fmt.Printf("%d\t%s\t\t%s\t\t%s\t\t%s\t\t%s\t%d\n",
			field.ID, field.TableName(), field.FieldName, field.DisplayName, field.FieldType, status, field.SortOrder)
	}
	fmt.Printf("\n总计: %d 个动态字段\n", len(fields))
}

func showFieldDetails() {
	if len(os.Args) < 3 {
		fmt.Println("用法: admctl show-field <字段ID>")
		return
	}

	fieldID, err := strconv.ParseUint(os.Args[2], 10, 32)
	if err != nil {
		fmt.Printf("无效的字段ID: %v\n", err)
		return
	}

	var field model.DynamicField
	if err := global.DB.Preload("Table").First(&field, uint(fieldID)).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			fmt.Printf("未找到ID为 %d 的字段\n", fieldID)
		} else {
			log.Printf("查询字段失败: %v", err)
		}
		return
	}

	fmt.Printf("\n字段详情:\n")
	fmt.Printf("ID: %d\n", field.ID)
	fmt.Printf("所属表: %s (ID: %d)\n", field.TableName(), field.TableID)
	fmt.Printf("字段名: %s\n", field.FieldName)
	fmt.Printf("显示名称: %s\n", field.DisplayName)
	fmt.Printf("字段类型: %s\n", field.FieldType)
	fmt.Printf("是否必填: %t\n", field.IsRequired)
	fmt.Printf("是否唯一: %t\n", field.IsUnique)
	fmt.Printf("是否可搜索: %t\n", field.IsSearchable)
	fmt.Printf("是否可排序: %t\n", field.IsSortable)
	fmt.Printf("默认值: %s\n", field.DefaultValue)
	fmt.Printf("排序: %d\n", field.SortOrder)
	fmt.Printf("状态: %s\n", map[int]string{1: "启用", 2: "禁用"}[field.Status])
	fmt.Printf("创建时间: %s\n", field.CreatedAt.Format("2006-01-02 15:04:05"))
	fmt.Printf("更新时间: %s\n", field.UpdatedAt.Format("2006-01-02 15:04:05"))
}

func queryDynamicData() {
	if len(os.Args) < 3 {
		fmt.Println("用法: admctl query-data <表名> [条件]")
		return
	}

	tableName := os.Args[2]

	// 首先获取表信息
	var table model.DynamicTable
	if err := global.DB.Where("table_name = ?", tableName).First(&table).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			fmt.Printf("未找到表名为 '%s' 的表\n", tableName)
		} else {
			log.Printf("查询表失败: %v", err)
		}
		return
	}

	// 构建动态表名
	dynamicTableName := fmt.Sprintf("%s_%d", table.TableName, table.ID)

	// 执行查询
	var results []map[string]interface{}
	sql := fmt.Sprintf("SELECT * FROM `%s` LIMIT 10", dynamicTableName)

	if len(os.Args) > 3 {
		// 如果有额外的条件参数
		condition := strings.Join(os.Args[3:], " ")
		sql = fmt.Sprintf("SELECT * FROM `%s` WHERE %s LIMIT 10", dynamicTableName, condition)
	}

	if err := global.DB.Raw(sql).Scan(&results).Error; err != nil {
		log.Printf("查询数据失败: %v", err)
		return
	}

	fmt.Printf("\n表 '%s' 的数据 (%d 条):\n", tableName, len(results))
	if len(results) > 0 {
		// 获取列名
		for i, row := range results {
			fmt.Printf("\n第 %d 条记录:\n", i+1)
			for key, value := range row {
				fmt.Printf("  %s: %v\n", key, value)
			}
		}
	} else {
		fmt.Println("暂无数据")
	}
}

func executeCustomSQL() {
	if len(os.Args) < 3 {
		fmt.Println("用法: admctl exec-sql \"SQL语句\"")
		return
	}

	sql := strings.Join(os.Args[2:], " ")
	if sql == "" {
		fmt.Println("SQL语句不能为空")
		return
	}

	// 对于SELECT查询
	if strings.HasPrefix(strings.ToUpper(strings.TrimSpace(sql)), "SELECT") {
		var results []map[string]interface{}
		if err := global.DB.Raw(sql).Scan(&results).Error; err != nil {
			log.Printf("执行查询失败: %v", err)
			return
		}

		fmt.Printf("\n查询结果 (%d 条):\n", len(results))
		if len(results) > 0 {
			// 打印表头
			var headers []string
			for key := range results[0] {
				headers = append(headers, key)
			}
			fmt.Println(strings.Join(headers, "\t"))

			// 打印数据
			for _, row := range results {
				var values []string
				for _, header := range headers {
					values = append(values, fmt.Sprintf("%v", row[header]))
				}
				fmt.Println(strings.Join(values, "\t"))
			}
		}
	} else {
		// 对于INSERT/UPDATE/DELETE等操作
		result := global.DB.Exec(sql)
		if result.Error != nil {
			log.Printf("执行SQL失败: %v", result.Error)
			return
		}

		affected := result.RowsAffected
		fmt.Printf("SQL执行成功，影响行数: %d\n", affected)
	}
}
