package initialize

import (
	"log"

	"go-react-admin/global"
	"go-react-admin/model"
)

// InitAdminUser 初始化管理员用户
func InitAdminUser() {
	// 检查是否已存在管理员用户
	var count int64
	global.DB.Model(&model.User{}).Where("username = ?", "admin").Count(&count)

	if count == 0 {
		// 创建管理员用户
		adminUser := model.User{
			Username: "admin",
			Password: "123456", // 实际项目中应该加密密码
			Nickname: "管理员",
			Email:    "admin@example.com",
			Phone:    "13800138000",
			Status:   1,
		}

		if err := global.DB.Create(&adminUser).Error; err != nil {
			log.Printf("创建管理员用户失败: %v", err)
		} else {
			log.Println("管理员用户创建成功")
		}
	} else {
		log.Println("管理员用户已存在")
	}
}
