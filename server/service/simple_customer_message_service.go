package service

import (
	"fmt"
	"go-react-admin/global"
	"go-react-admin/model"
	"time"
)

type SimpleCustomerMessageService struct{}

type SimpleCustomerMessageRequest struct {
	Page     int    `form:"page" json:"page" binding:"min=1"`
	PageSize int    `form:"page_size" json:"page_size" binding:"min=1,max=100"`
	Type     string `form:"type,omitempty" json:"type,omitempty"`
	Keyword  string `form:"keyword,omitempty" json:"keyword,omitempty" binding:"max=100"`
	IsRead   *bool  `form:"is_read,omitempty" json:"is_read,omitempty"`
}

type SimpleCustomerMessageResponse struct {
	List  []SimpleMessageDetail `json:"list"`
	Total int64                 `json:"total"`
}

type SimpleMessageDetail struct {
	ID        uint64     `json:"id"`
	Title     string     `json:"title"`
	Content   string     `json:"content"`
	Type      string     `json:"type"`
	IsRead    bool       `json:"is_read"`
	ReadTime  *time.Time `json:"read_time,omitempty"`
	CreatedAt time.Time  `json:"created_at"`
}

// GetCustomerMessages 获取客户的站内信列表（简化版）
func (s *SimpleCustomerMessageService) GetCustomerMessages(customerID uint64, req *SimpleCustomerMessageRequest) (*SimpleCustomerMessageResponse, error) {
	var messages []SimpleMessageDetail
	var total int64

	// 创建消息阅读记录表（如果不存在）
	s.createCustomerReadTableIfNotExists()

	query := global.DB.Table("messages as m").
		Select("m.id, m.title, m.content, m.type, COALESCE(cm.is_read, false) as is_read, cm.read_time, m.created_at").
		Joins("LEFT JOIN customer_message_reads cm ON m.id = cm.message_id AND cm.customer_id = ?", customerID).
		Where("m.target_type = ? AND m.target_id = ? AND m.status = ?", "customer", customerID, "published").
		Or("m.target_type = ? AND m.status = ?", "all", "published")

	// 筛选条件
	if req.Type != "" {
		query = query.Where("m.type = ?", req.Type)
	}

	if req.Keyword != "" {
		query = query.Where("m.title LIKE ? OR m.content LIKE ?", "%"+req.Keyword+"%", "%"+req.Keyword+"%")
	}

	if req.IsRead != nil {
		if *req.IsRead {
			query = query.Where("cm.is_read = ?", true)
		} else {
			query = query.Where("(cm.is_read IS NULL OR cm.is_read = ?)", false)
		}
	}

	// 统计总数
	query.Count(&total)

	// 分页查询
	offset := (req.Page - 1) * req.PageSize
	if err := query.Order("m.created_at DESC").Limit(req.PageSize).Offset(offset).Find(&messages).Error; err != nil {
		return nil, fmt.Errorf("查询消息列表失败: %v", err)
	}

	return &SimpleCustomerMessageResponse{
		List:  messages,
		Total: total,
	}, nil
}

// MarkMessageAsRead 标记消息为已读
func (s *SimpleCustomerMessageService) MarkMessageAsRead(customerID uint64, messageID uint64) error {
	// 确保消息存在且属于该客户
	var message model.Message
	if err := global.DB.Where("id = ? AND ((target_type = ? AND target_id = ?) OR target_type = ?)", 
		messageID, "customer", customerID, "all").First(&message).Error; err != nil {
		return fmt.Errorf("消息不存在或无权限访问")
	}

	// 创建或更新阅读记录
	return global.DB.Exec(`
		INSERT INTO customer_message_reads (customer_id, message_id, is_read, read_time) 
		VALUES (?, ?, ?, ?) 
		ON DUPLICATE KEY UPDATE is_read = ?, read_time = ?
	`, customerID, messageID, true, time.Now(), true, time.Now()).Error
}

// GetUnreadCount 获取未读消息数量
func (s *SimpleCustomerMessageService) GetUnreadCount(customerID uint64) (int64, error) {
	s.createCustomerReadTableIfNotExists()

	var count int64
	err := global.DB.Raw(`
		SELECT COUNT(*) 
		FROM messages m 
		LEFT JOIN customer_message_reads cm ON m.id = cm.message_id AND cm.customer_id = ? 
		WHERE (m.target_type = ? AND m.target_id = ? OR m.target_type = ?) 
		AND m.status = ? 
		AND (cm.is_read IS NULL OR cm.is_read = ?)
	`, customerID, "customer", customerID, "all", "published", false).Scan(&count).Error

	return count, err
}

// createCustomerReadTableIfNotExists 创建客户消息阅读记录表
func (s *SimpleCustomerMessageService) createCustomerReadTableIfNotExists() {
	global.DB.Exec(`
		CREATE TABLE IF NOT EXISTS customer_message_reads (
			id BIGINT AUTO_INCREMENT PRIMARY KEY,
			customer_id BIGINT NOT NULL,
			message_id BIGINT NOT NULL,
			is_read BOOLEAN DEFAULT FALSE,
			read_time DATETIME,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
			UNIQUE KEY uk_customer_message (customer_id, message_id),
			INDEX idx_customer (customer_id),
			INDEX idx_message (message_id)
		)
	`)
}

// CreateCustomerMessage 创建客户消息
func (s *SimpleCustomerMessageService) CreateCustomerMessage(title, content, msgType string, customerID uint64) error {
	message := &model.Message{
		Title:      title,
		Content:    content,
		Type:       model.MessageType(msgType),
		Status:     model.MessageStatusPublished,
		TargetType: "customer",
		TargetID:   &customerID,
		SenderType: "system",
	}

	return global.DB.Create(message).Error
}