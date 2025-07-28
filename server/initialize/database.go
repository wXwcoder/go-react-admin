package initialize

import (
	"fmt"
	"log"
	"net/url"

	"go-react-admin/global"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// InitDB 初始化数据库连接
func InitDB() {
	// 构建DSN
	config := global.GlobalConfig.Mysql
	// 对密码和loc参数进行URL编码以避免特殊字符问题
	encodedPassword := url.QueryEscape(config.Password)
	encodedLoc := url.QueryEscape(config.Loc)
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=%s&parseTime=%t&loc=%s",
		config.Username, encodedPassword, config.Host, config.Port, config.Dbname,
		config.Charset, config.ParseTime, encodedLoc)
	fmt.Printf("DSN: %s\n", dsn)

	// 连接数据库
	var err error
	global.DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("连接数据库失败: %v", err)
	}

	fmt.Println("数据库连接成功")
}
