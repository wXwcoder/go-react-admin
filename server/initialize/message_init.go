package initialize

import (
	"go-react-admin/global"
	"go-react-admin/model"
	"log"
	"time"
)

// InitMessageSystem 初始化消息系统增强功能
func InitMessageSystem() {
	log.Println("开始初始化消息系统增强功能...")
	
	// 创建数据库索引
	createMessageIndexes()
	
	// 初始化消息模板
	initMessageTemplates()
	
	// 初始化测试数据
	initMessageTestData()
	
	// 创建触发器（由于GORM不支持直接创建触发器，这里使用原生SQL）
	createMessageTriggers()
	
	log.Println("消息系统初始化完成")
}

// 创建数据库索引
func createMessageIndexes() {
	db := global.DB
	
	// 消息表索引
	db.Exec("CREATE INDEX IF NOT EXISTS idx_type_status ON messages(type, status)")
	db.Exec("CREATE INDEX IF NOT EXISTS idx_target ON messages(target_type, target_id)")
	db.Exec("CREATE INDEX IF NOT EXISTS idx_expired_at ON messages(expired_at)")
	
	// 客户消息表索引
	db.Exec("CREATE INDEX IF NOT EXISTS idx_customer_read ON customer_messages(customer_id, is_read)")
	db.Exec("CREATE INDEX IF NOT EXISTS idx_message_customer ON customer_messages(message_id, customer_id)")
	
	// 公告表索引
	db.Exec("CREATE INDEX IF NOT EXISTS idx_announcement_type ON announcements(type)")
	db.Exec("CREATE INDEX IF NOT EXISTS idx_announcement_status ON announcements(status)")
	db.Exec("CREATE INDEX IF NOT EXISTS idx_announcement_created_at ON announcements(created_at)")
	
	// 公告阅读记录表索引
	db.Exec("CREATE INDEX IF NOT EXISTS idx_announcement_read_customer ON announcement_reads(customer_id)")
	db.Exec("CREATE INDEX IF NOT EXISTS idx_announcement_read_time ON announcement_reads(read_time)")
	
	log.Println("数据库索引创建完成")
}

// 初始化消息模板
func initMessageTemplates() {
	templates := []model.MessageTemplate{
		{
			Name:            "welcome_customer",
			TitleTemplate:   "欢迎注册，{{.CustomerName}}！",
			ContentTemplate: "尊敬的{{.CustomerName}}，欢迎注册成为我们的客户。您的账号{{.Username}}已成功创建。",
			Type:            "system",
			Variables:       []byte(`["CustomerName", "Username"]`),
			IsActive:        true,
		},
		{
			Name:            "order_notification",
			TitleTemplate:   "订单状态更新",
			ContentTemplate: "您的订单{{.OrderNumber}}状态已更新为：{{.Status}}",
			Type:            "notice",
			Variables:       []byte(`["OrderNumber", "Status"]`),
			IsActive:        true,
		},
		{
			Name:            "system_maintenance",
			TitleTemplate:   "系统维护通知",
			ContentTemplate: "系统将于{{.StartTime}}至{{.EndTime}}进行维护，期间{{.Impact}}",
			Type:            "system",
			Variables:       []byte(`["StartTime", "EndTime", "Impact"]`),
			IsActive:        true,
		},
	}
	
	for _, template := range templates {
		var existing model.MessageTemplate
		if err := global.DB.Where("name = ?", template.Name).First(&existing).Error; err != nil {
			// 如果记录不存在，则创建
			if err := global.DB.Create(&template).Error; err != nil {
				log.Printf("创建消息模板失败: %v", err)
			}
		}
	}
	
	log.Println("消息模板初始化完成")
}

// 初始化测试数据
func initMessageTestData() {
	// 初始化公告测试数据
	announcements := []model.Announcement{
		{
			Title:     "系统维护通知",
			Content:   "系统将于今晚23:00-24:00进行例行维护，期间可能无法访问。",
			Type:      "maintenance",
			Status:    "published",
			Priority:  8,
			ExpiredAt: func() *time.Time { t := time.Now().AddDate(0, 0, 7); return &t }(),
			ReadCount: 0,
		},
		{
			Title:     "新功能上线公告",
			Content:   "我们新增了客户消息功能，现在您可以接收重要通知了！",
			Type:      "update",
			Status:    "published",
			Priority:  5,
			ExpiredAt: func() *time.Time { t := time.Now().AddDate(0, 0, 30); return &t }(),
			ReadCount: 0,
		},
		{
			Title:     "欢迎使用系统",
			Content:   "欢迎注册成为我们的客户，如有任何问题请联系客服。",
			Type:      "system",
			Status:    "draft",
			Priority:  3,
			ExpiredAt: func() *time.Time { t := time.Now().AddDate(1, 0, 0); return &t }(),
			ReadCount: 0,
		},
	}
	
	for _, announcement := range announcements {
		var existing model.Announcement
		if err := global.DB.Where("title = ?", announcement.Title).First(&existing).Error; err != nil {
			// 如果记录不存在，则创建
			if err := global.DB.Create(&announcement).Error; err != nil {
				log.Printf("创建公告失败: %v", err)
			}
		}
	}
	
	// 初始化消息测试数据
	messages := []model.Message{
		{
			Title:     "欢迎消息",
			Content:   "欢迎注册成为我们的客户！",
			Type:      "private",
			Status:    "published",
			Priority:  5,
			TargetType: "all",
			ExpiredAt: func() *time.Time { t := time.Now().AddDate(1, 0, 0); return &t }(),
			ReadCount: 0,
		},
		{
			Title:     "系统更新通知",
			Content:   "系统已更新至最新版本，体验更流畅！",
			Type:      "system",
			Status:    "published",
			Priority:  7,
			TargetType: "all",
			ExpiredAt: func() *time.Time { t := time.Now().AddDate(0, 0, 30); return &t }(),
			ReadCount: 0,
		},
	}
	
	for _, message := range messages {
		var existing model.Message
		if err := global.DB.Where("title = ?", message.Title).First(&existing).Error; err != nil {
			// 如果记录不存在，则创建
			if err := global.DB.Create(&message).Error; err != nil {
				log.Printf("创建消息失败: %v", err)
			}
		}
	}
	
	log.Println("测试数据初始化完成")
}

// 创建消息触发器
func createMessageTriggers() {
	db := global.DB
	
	// 删除已存在的触发器
	db.Exec("DROP TRIGGER IF EXISTS update_announcement_read_count")
	db.Exec("DROP TRIGGER IF EXISTS update_message_statistics")
	
	// 创建公告阅读统计触发器
	triggerSQL := `
CREATE TRIGGER IF NOT EXISTS update_announcement_read_count 
AFTER INSERT ON announcement_reads
FOR EACH ROW
BEGIN
    IF NEW.is_read = 1 THEN
        UPDATE announcements 
        SET read_count = read_count + 1 
        WHERE id = NEW.announcement_id;
    END IF;
END`
	db.Exec(triggerSQL)
	
	// 创建消息统计触发器
	triggerSQL2 := `
CREATE TRIGGER IF NOT EXISTS update_message_statistics
AFTER UPDATE ON customer_messages
FOR EACH ROW
BEGIN
    IF OLD.is_read = 0 AND NEW.is_read = 1 THEN
        UPDATE messages 
        SET read_count = read_count + 1 
        WHERE id = NEW.message_id;
    END IF;
END`
	db.Exec(triggerSQL2)
	
	log.Println("消息触发器创建完成")
}