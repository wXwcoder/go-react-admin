package service

import (
	"errors"
	"go-react-admin/global"
	"go-react-admin/model"
	"strconv"
	"time"

	"gorm.io/gorm"
)

type MessageService struct{}

// MessageListRequest 消息列表请求
type MessageListRequest struct {
	Page     int    `form:"page" json:"page" binding:"min=1"`
	PageSize int    `form:"page_size" json:"page_size" binding:"min=1,max=100"`
	Type     string `form:"type" json:"type" binding:"omitempty,oneof=system notice private"`
	Status   string `form:"status" json:"status" binding:"omitempty,oneof=draft published revoked"`
	Keyword  string `form:"keyword" json:"keyword"`
}

// MessageListResponse 消息列表响应
type MessageListResponse struct {
	List  []model.Message `json:"list"`
	Total int64           `json:"total"`
}

// MessageDetailResponse 消息详情响应
type MessageDetailResponse struct {
	Message *model.Message `json:"message"`
}

// MessageCreateRequest 创建消息请求
type MessageCreateRequest struct {
	Title      string    `json:"title" binding:"required,max=255"`
	Content    string    `json:"content" binding:"required"`
	Type       string    `json:"type" binding:"required,oneof=system notice private"`
	Priority   string    `json:"priority" binding:"omitempty,oneof=low medium high"`
	TargetType string    `json:"target_type" binding:"omitempty,oneof=all customer"`
	TargetID   uint64    `json:"target_id"`
	ExpireTime time.Time `json:"expire_time"`
}

// MessageUpdateRequest 更新消息请求
type MessageUpdateRequest struct {
	Title      string    `json:"title" binding:"omitempty,max=255"`
	Content    string    `json:"content" binding:"omitempty"`
	Type       string    `json:"type" binding:"omitempty,oneof=system notice private"`
	Priority   string    `json:"priority" binding:"omitempty,oneof=low medium high"`
	ExpireTime time.Time `json:"expire_time"`
}

// CustomerMessageAdminListRequest 管理员获取客户消息列表请求
type CustomerMessageAdminListRequest struct {
	Page      int  `form:"page" json:"page" binding:"min=1"`
	PageSize  int  `form:"page_size" json:"page_size" binding:"min=1,max=100"`
	CustomerID uint64 `form:"customer_id" json:"customer_id"`
	IsRead   *bool  `form:"is_read" json:"is_read"`
}

// CustomerMessageAdminListResponse 管理员获取客户消息列表响应
 type CustomerMessageAdminListResponse struct {
	List  []model.CustomerMessage `json:"list"`
	Total int64                   `json:"total"`
}

// CustomerMessageSendRequest 发送客户消息请求
 type CustomerMessageSendRequest struct {
	CustomerID uint64 `json:"customer_id" binding:"required"`
	Title      string `json:"title" binding:"required,max=255"`
	Content    string `json:"content" binding:"required"`
	Priority   string `json:"priority" binding:"omitempty,oneof=low medium high"`
	ExpireTime time.Time `json:"expire_time"`
}

// GetAdminMessageList 获取管理员消息列表
func (s *MessageService) GetAdminMessageList(req *MessageListRequest) (*MessageListResponse, error) {
	if req.Page == 0 {
		req.Page = 1
	}
	if req.PageSize == 0 {
		req.PageSize = 10
	}

	var messages []model.Message
	query := global.DB.Model(&model.Message{})

	if req.Type != "" {
		query = query.Where("type = ?", req.Type)
	}
	if req.Status != "" {
		query = query.Where("status = ?", req.Status)
	}
	if req.Keyword != "" {
		query = query.Where("title LIKE ? OR content LIKE ?", "%"+req.Keyword+"%", "%"+req.Keyword+"%")
	}

	var total int64
	query.Count(&total)

	offset := (req.Page - 1) * req.PageSize
	if err := query.Order("created_at DESC").Limit(req.PageSize).Offset(offset).Find(&messages).Error; err != nil {
		return nil, err
	}

	return &MessageListResponse{
		List:  messages,
		Total: total,
	}, nil
}

// GetAdminMessageDetail 获取管理员消息详情
func (s *MessageService) GetAdminMessageDetail(messageID uint64) (*MessageDetailResponse, error) {
	var message model.Message
	if err := global.DB.First(&message, messageID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("消息不存在")
		}
		return nil, err
	}

	return &MessageDetailResponse{Message: &message}, nil
}

// CreateMessage 创建消息
func (s *MessageService) CreateMessage(req *MessageCreateRequest) (*model.Message, error) {
	priority, _ := strconv.Atoi(req.Priority)
	message := &model.Message{
		Title:      req.Title,
		Content:    req.Content,
		Type:       model.MessageType(req.Type),
		Priority:   priority,
		Status:     model.MessageStatusDraft,
		TargetType: req.TargetType,
		TargetID:   &req.TargetID,
		ExpiredAt:  &req.ExpireTime,
	}

	if err := global.DB.Create(message).Error; err != nil {
		return nil, err
	}

	return message, nil
}

// UpdateMessage 更新消息
func (s *MessageService) UpdateMessage(messageID uint64, req *MessageUpdateRequest) error {
	var message model.Message
	if err := global.DB.First(&message, messageID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("消息不存在")
		}
		return err
	}

	if message.Status == "published" {
		return errors.New("已发布的消息不能修改")
	}

	updates := map[string]interface{}{}
	if req.Title != "" {
		updates["title"] = req.Title
	}
	if req.Content != "" {
		updates["content"] = req.Content
	}
	if req.Type != "" {
		updates["type"] = req.Type
	}
	if req.Priority != "" {
		priority, _ := strconv.Atoi(req.Priority)
		updates["priority"] = priority
	}
	if !req.ExpireTime.IsZero() {
		updates["expired_at"] = req.ExpireTime
	}

	return global.DB.Model(&message).Updates(updates).Error
}

// DeleteMessage 删除消息
func (s *MessageService) DeleteMessage(messageID uint64) error {
	var message model.Message
	if err := global.DB.First(&message, messageID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("消息不存在")
		}
		return err
	}

	return global.DB.Delete(&message).Error
}

// SendMessage 发送消息
func (s *MessageService) SendMessage(messageID uint64) error {
	var message model.Message
	if err := global.DB.First(&message, messageID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("消息不存在")
		}
		return err
	}

	if message.Status == "published" {
		return errors.New("消息已发布")
	}

	return global.DB.Model(&message).Update("status", model.MessageStatusPublished).Error
}

// CancelMessage 撤回消息
func (s *MessageService) CancelMessage(messageID uint64) error {
	var message model.Message
	if err := global.DB.First(&message, messageID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("消息不存在")
		}
		return err
	}

	if message.Status != model.MessageStatusPublished {
		return errors.New("消息未发布")
	}

	return global.DB.Model(&message).Update("status", model.MessageStatusRevoked).Error
}

// GetCustomerMessageAdminList 获取客户消息列表（管理员）
func (s *MessageService) GetCustomerMessageAdminList(req *CustomerMessageAdminListRequest) (*CustomerMessageAdminListResponse, error) {
	if req.Page == 0 {
		req.Page = 1
	}
	if req.PageSize == 0 {
		req.PageSize = 10
	}

	var messages []model.CustomerMessage
	query := global.DB.Model(&model.CustomerMessage{})

	if req.CustomerID > 0 {
		query = query.Where("customer_id = ?", req.CustomerID)
	}
	if req.IsRead != nil {
		query = query.Where("is_read = ?", *req.IsRead)
	}

	var total int64
	query.Count(&total)

	offset := (req.Page - 1) * req.PageSize
	if err := query.Preload("Message").Preload("Customer").
		Order("created_at DESC").Limit(req.PageSize).Offset(offset).Find(&messages).Error; err != nil {
		return nil, err
	}

	return &CustomerMessageAdminListResponse{
		List:  messages,
		Total: total,
	}, nil
}

// SendCustomerMessage 发送客户消息
func (s *MessageService) SendCustomerMessage(req *CustomerMessageSendRequest) error {
	// 创建消息
	priority, _ := strconv.Atoi(req.Priority)
	message := &model.Message{
		Title:      req.Title,
		Content:    req.Content,
		Type:       model.MessageTypePrivate,
		Priority:   priority,
		Status:     model.MessageStatusPublished,
		TargetType: "customer",
		TargetID:   &req.CustomerID,
		ExpiredAt:  &req.ExpireTime,
	}

	if err := global.DB.Create(message).Error; err != nil {
		return err
	}

	// 创建客户消息关联
	customerMessage := &model.CustomerMessage{
		MessageID:  message.ID,
		CustomerID: req.CustomerID,
		IsRead:     false,
		IsDeleted:  false,
	}

	return global.DB.Create(customerMessage).Error
}