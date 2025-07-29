package global

import (
	"github.com/casbin/casbin/v2"
	"github.com/go-redis/redis/v8"
	"gorm.io/gorm"
)

var (
	// 数据库连接
	DB *gorm.DB

	// Redis客户端
	RedisClient *redis.Client

	// Casbin权限管理器
	Enforcer *casbin.Enforcer

	// 配置信息
	GlobalConfig *Config
)

type Config struct {
	Server     ServerConfig     `yaml:"server"`
	Mysql      MysqlConfig      `yaml:"mysql"`
	Redis      RedisConfig      `yaml:"redis"`
	Log        LogConfig        `yaml:"log"`
	Jwt        JwtConfig        `yaml:"jwt"`
	MultiTenant MultiTenantConfig `yaml:"multi_tenant"`
	System     SystemConfig     `yaml:"system"`
}

type ServerConfig struct {
	Port string `yaml:"port"`
	Mode string `yaml:"mode"`
}

type MysqlConfig struct {
	Host       string `yaml:"host"`
	Port       int    `yaml:"port"`
	Username   string `yaml:"username"`
	Password   string `yaml:"password"`
	Dbname     string `yaml:"dbname"`
	Charset    string `yaml:"charset"`
	ParseTime  bool   `yaml:"parseTime"`
	Loc        string `yaml:"loc"`
}

type RedisConfig struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	Password string `yaml:"password"`
	Db       int    `yaml:"db"`
}

type LogConfig struct {
	Level  string `yaml:"level"`
	Format string `yaml:"format"`
	Output string `yaml:"output"`
}

type JwtConfig struct {
	Secret string `yaml:"secret"`
	Expire int    `yaml:"expire"`
}

type MultiTenantConfig struct {
	Enabled bool   `yaml:"enabled"`
	Mode    string `yaml:"mode"`
}

type SystemConfig struct {
	Name    string `yaml:"name"`
	Version string `yaml:"version"`
	Theme   string `yaml:"theme"`
}