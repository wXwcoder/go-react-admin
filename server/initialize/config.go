package initialize

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
	"go-react-admin/global"
)

// LoadConfig 从环境变量加载配置（完全废弃YAML配置）
func LoadConfig() {
	// 自动加载.env文件（如果存在）
	if err := godotenv.Load(); err != nil {
		// .env文件不存在时使用系统环境变量
		log.Printf("未找到.env文件，使用系统环境变量: %v", err)
	} else {
		log.Println("成功加载.env文件")
	}

	config := &global.Config{}

	// 服务器配置
	config.Server.Port = getEnv("SERVER_PORT", "8080")
	config.Server.Mode = getEnv("GIN_MODE", "debug")

	// 数据库配置
	config.Mysql = global.MysqlConfig{
		Host:      getEnv("DB_HOST", "localhost"),
		Port:      getEnvAsInt("DB_PORT", 2603),
		Username:  getEnv("DB_USER", "root"),
		Password:  getEnv("DB_PASSWORD", "123456"),
		Dbname:    getEnv("DB_NAME", "go_react_admin"),
		Charset:   getEnv("DB_CHARSET", "utf8mb4"),
		ParseTime: getEnvAsBool("DB_PARSE_TIME", true),
		Loc:       getEnv("DB_LOC", "Asia/Shanghai"),
	}

	// Redis配置
	config.Redis = global.RedisConfig{
		Host:     getEnv("REDIS_HOST", "localhost"),
		Port:     getEnvAsInt("REDIS_PORT", 6379),
		Password: getEnv("REDIS_PASSWORD", ""),
		Db:       getEnvAsInt("REDIS_DB", 0),
	}

	// 日志配置
	config.Log = global.LogConfig{
		Level:  getEnv("LOG_LEVEL", "info"),
		Format: getEnv("LOG_FORMAT", "json"),
		Output: getEnv("LOG_OUTPUT", "console"),
	}

	// JWT配置
	config.Jwt = global.JwtConfig{
		Secret: getEnv("JWT_SECRET", "go-react-admin-secret"),
		Expire: getEnvAsInt("JWT_EXPIRE", 24),
	}

	// 多租户配置
	config.MultiTenant = global.MultiTenantConfig{
		Enabled: getEnvAsBool("MULTI_TENANT_ENABLED", true),
		Mode:    getEnv("MULTI_TENANT_MODE", "shared_schema"),
	}

	// 系统配置
	config.System = global.SystemConfig{
		Name:    getEnv("SYSTEM_NAME", "go-react-admin"),
		Version: getEnv("SYSTEM_VERSION", "1.0.0"),
		Theme:   getEnv("SYSTEM_THEME", "light"),
	}

	global.GlobalConfig = config

	fmt.Printf("环境变量配置加载成功:\n")
	fmt.Printf("Server: %s:%d\n", "0.0.0.0", config.Server.Port)
	fmt.Printf("Database: %s:%d/%s\n", config.Mysql.Host, config.Mysql.Port, config.Mysql.Dbname)
	fmt.Printf("Redis: %s:%d/%d\n", config.Redis.Host, config.Redis.Port, config.Redis.Db)
}

// getEnv 获取环境变量，不存在则返回默认值
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// getEnvAsInt 获取环境变量并转换为int
func getEnvAsInt(key string, defaultValue int) int {
	valueStr := getEnv(key, "")
	if valueStr == "" {
		return defaultValue
	}
	if value, err := strconv.Atoi(valueStr); err == nil {
		return value
	}
	log.Printf("警告: 环境变量 %s 不是有效的整数，使用默认值 %d", key, defaultValue)
	return defaultValue
}

// getEnvAsBool 获取环境变量并转换为bool
func getEnvAsBool(key string, defaultValue bool) bool {
	valueStr := getEnv(key, "")
	if valueStr == "" {
		return defaultValue
	}
	if value, err := strconv.ParseBool(valueStr); err == nil {
		return value
	}
	log.Printf("警告: 环境变量 %s 不是有效的布尔值，使用默认值 %t", key, defaultValue)
	return defaultValue
}
