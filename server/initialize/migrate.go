package initialize

import (
	"log"

	"go-react-admin/global"
	"go-react-admin/model"
)

// Migrate 数据库迁移
func Migrate() {
	// 自动迁移模型
	err := global.DB.AutoMigrate(
		&model.User{},
		&model.Role{},
		&model.Menu{},
		&model.Api{},
		&model.Log{},
		&model.Tenant{},
		&model.UserRole{},
		&model.RoleMenu{},
		&model.RoleApi{},
	)
	if err != nil {
		log.Fatalf("数据库迁移失败: %v", err)
	}

	log.Println("数据库迁移成功")
}
