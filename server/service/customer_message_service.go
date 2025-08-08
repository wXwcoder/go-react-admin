package service

import (
	"errors"
	"fmt"
	"go-react-admin/global"
	"go-react-admin/model"
	"time"

	"gorm.io/gorm"
)

type CustomerMessageService struct{}

type CustomerMessageListRequest struct {
	Page       int    `form:"page" json:"page" binding:"min=1"`
	PageSize   int    `form:"page_size" json:"page_size" binding:"min=1,max=100"`
	IsRead     *bool  `form:"is_read,omitempty" json:"is_read,omitempty"`
	Type       string `form:"type,omitempty" json:"type,omitempty"`
	Keyword    string `form:"keyword,omitempty" json:"keyword,omitempty" binding:"max=100"`
	CustomerID uint64 `form:"customer_id,omitempty" json:"customer_id,omitempty"`
}

type CustomerMessageListResponse struct {
	List  []CustomerMessageDetail `json:"list"`
	Total int64                   `json:"total"`
}

type CustomerMessageBatchReadRequest struct {
	MessageIDs []uint64 `json:"message_ids" binding:"required,min=1"`
}

type AdminCreateCustomerMessageRequest struct {
	Title       string   `json:"title" binding:"required,max=200"`
	Content     string   `json:"content" binding:"required,max=2000"`
	Type        string   `json:"type" binding:"required,oneof=system notice warning"`
	CustomerIDs []uint64 `json:"customer_ids" binding:"required,min=1"`
}

type AdminCreateCustomerMessagesBatchRequest struct {
	CustomerIDs []uint64 `json:"customer_ids" binding:"required,min=1"`
	Title       string   `json:"title" binding:"required,max=200"`
	Content     string   `json:"content" binding:"required,max=2000"`
	Type        string   `json:"type" binding:"required,oneof=system notice warning"`
}

type CustomerMessageDetail struct {
	ID        uint64         `json:"id"`
	MessageID uint64         `json:"message_id"`
	Title     string         `json:"title"`
	Content   string         `json:"content"`
	Type      string         `json:"type"`
	IsRead    bool           `json:"is_read"`
	ReadTime  *time.Time     `json:"read_time,omitempty"`
	CreatedAt time.Time      `json:"created_at"`
	Message   *model.Message `json:"message,omitempty"`
}

// GetCustomerMessages 获取客户的站内信列表
func (s *CustomerMessageService) GetCustomerMessages(customerID uint64, req *CustomerMessageListRequest) (*CustomerMessageListResponse, error) {
	var messages []CustomerMessageDetail
	var total int64

	query := global.DB.Table("customer_messages as cm").
		Select("cm.id, cm.message_id, m.title, m.content, m.type, cm.is_read, cm.read_time, cm.created_at").
		Joins("JOIN messages m ON cm.message_id = m.id").
		Where("cm.customer_id = ? AND cm.is_deleted = ?", customerID, false)

	// 筛选条件
	if req.IsRead != nil {
		query = query.Where("cm.is_read = ?", *req.IsRead)
	}

	if req.Type != "" {
		query = query.Where("m.type = ?", req.Type)
	}

	if req.Keyword != "" {
		query = query.Where("m.title LIKE ? OR m.content LIKE ?", "%"+req.Keyword+"%", "%"+req.Keyword+"%")
	}

	// 统计总数
	query.Count(&total)

	// 分页查询
	offset := (req.Page - 1) * req.PageSize
	if err := query.Order("cm.created_at DESC").Limit(req.PageSize).Offset(offset).Find(&messages).Error; err != nil {
		return nil, fmt.Errorf("查询消息列表失败: %v", err)
	}

	return &CustomerMessageListResponse{
		List:  messages,
		Total: total,
	}, nil
}

// GetCustomerMessageDetail 获取客户站内信详情
func (s *CustomerMessageService) GetCustomerMessageDetail(customerID uint64, messageID uint64) (*CustomerMessageDetail, error) {
	var detail CustomerMessageDetail

	err := global.DB.Table("customer_messages as cm").
		Select("cm.id, cm.message_id, m.title, m.content, m.type, cm.is_read, cm.read_time, cm.created_at, m.*").
		Joins("JOIN messages m ON cm.message_id = m.id").
		Where("cm.customer_id = ? AND cm.message_id = ? AND cm.is_deleted = ?", customerID, messageID, false).
		First(&detail).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("消息不存在或无权限访问")
		}
		return nil, fmt.Errorf("查询消息详情失败: %v", err)
	}

	return &detail, nil
}

// MarkMessageAsRead 标记消息为已读
func (s *CustomerMessageService) MarkMessageAsRead(customerID uint64, messageID uint64) error {
	var customerMessage model.CustomerMessage

	err := global.DB.Where("customer_id = ? AND message_id = ?", customerID, messageID).First(&customerMessage).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("消息不存在")
		}
		return fmt.Errorf("查询消息失败: %v", err)
	}

	if customerMessage.IsRead {
		return nil // 已经标记为已读
	}

	now := time.Now()
	return global.DB.Model(&model.CustomerMessage{}).
		Where("customer_id = ? AND message_id = ?", customerID, messageID).
		Updates(map[string]interface{}{
			"is_read":   true,
			"read_time": now,
		}).Error
}

// MarkMessagesAsReadBatch 批量标记消息为已读
func (s *CustomerMessageService) MarkMessagesAsReadBatch(customerID uint64, messageIDs []uint64) error {
	now := time.Now()
	return global.DB.Model(&model.CustomerMessage{}).
		Where("customer_id = ? AND message_id IN ? AND is_read = ?", customerID, messageIDs, false).
		Updates(map[string]interface{}{
			"is_read":   true,
			"read_time": now,
		}).Error
}

// GetUnreadCount 获取未读消息数量
func (s *CustomerMessageService) GetUnreadCount(customerID uint64) (int64, error) {
	var count int64
	err := global.DB.Model(&model.CustomerMessage{}).
		Where("customer_id = ? AND is_read = ? AND is_deleted = ?", customerID, false, false).
		Count(&count).Error
	return count, err
}

// DeleteCustomerMessage 删除客户消息（软删除）
func (s *CustomerMessageService) DeleteCustomerMessage(customerID uint64, messageID uint64) error {
	return global.DB.Model(&model.CustomerMessage{}).
		Where("customer_id = ? AND message_id = ?", customerID, messageID).
		Update("is_deleted", true).Error
}

// CreateCustomerMessage 创建客户消息关联（用于系统消息推送）
func (s *CustomerMessageService) CreateCustomerMessage(customerID uint64, messageID uint64) error {
	customerMessage := &model.CustomerMessage{
		MessageID:  messageID,
		CustomerID: customerID,
		IsRead:     false,
	}

	return global.DB.Create(customerMessage).Error
}

// CreateCustomerMessagesBatch 批量创建客户消息关联
func (s *CustomerMessageService) CreateCustomerMessagesBatch(customerIDs []uint64, messageID uint64) error {
	var customerMessages []model.CustomerMessage
	for _, customerID := range customerIDs {
		customerMessages = append(customerMessages, model.CustomerMessage{
			MessageID:  messageID,
			CustomerID: customerID,
			IsRead:     false,
		})
	}

	return global.DB.CreateInBatches(customerMessages, 100).Error
}
