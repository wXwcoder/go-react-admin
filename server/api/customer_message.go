package api

import (
	"go-react-admin/global"
	"go-react-admin/model"
	"go-react-admin/service"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

type CustomerMessageAPI struct{}

// GetMessages 获取第三方客户消息列表
// @Summary 获取第三方客户消息列表
// @Description 获取当前登录的第三方客户的消息列表
// @Tags 第三方客户消息
// @Accept json
// @Produce json
// @Param page query int false "页码" default(1)
// @Param page_size query int false "每页数量" default(10) minimum(1) maximum(100)
// @Param is_read query bool false "是否已读"
// @Param type query string false "消息类型"
// @Param keyword query string false "关键词"
// @Success 200 {object} service.SimpleCustomerMessageResponse
// @Router /api/v1/customer/messages [get]
// @Security CustomerJWTAuth
func (api *CustomerMessageAPI) GetMessages(c *gin.Context) {
	customerID, exists := c.Get("customer_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "未授权访问"})
		return
	}

	var req service.SimpleCustomerMessageRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 设置默认值
	if req.Page == 0 {
		req.Page = 1
	}
	if req.PageSize == 0 {
		req.PageSize = 10
	}

	messageService := &service.SimpleCustomerMessageService{}
	response, err := messageService.GetCustomerMessages(uint64(customerID.(uint)), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": response})
}

// GetMessageDetail 获取消息详情
// @Summary 获取消息详情
// @Description 获取指定消息的详细信息
// @Tags 第三方客户消息
// @Accept json
// @Produce json
// @Param id path int true "消息ID"
// @Success 200 {object} model.Message
// @Router /api/v1/customer/messages/{id} [get]
// @Security CustomerJWTAuth
func (api *CustomerMessageAPI) GetMessageDetail(c *gin.Context) {
	customerID, exists := c.Get("customer_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "未授权访问"})
		return
	}

	messageIDStr := c.Param("id")
	messageID, err := strconv.ParseUint(messageIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的消息ID"})
		return
	}

	var message model.Message
	err = global.DB.Where("id = ? AND ((target_type = ? AND target_id = ?) OR target_type = ?)", 
		messageID, "customer", customerID, "all").First(&message).Error
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "消息不存在或无权限访问"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": message})
}

// MarkMessageAsRead 标记消息为已读
// @Summary 标记消息为已读
// @Description 将指定消息标记为已读
// @Tags 第三方客户消息
// @Accept json
// @Produce json
// @Param id path int true "消息ID"
// @Success 200 {object} gin.H{"message":"标记成功"}
// @Router /api/v1/customer/messages/{id}/read [put]
// @Security CustomerJWTAuth
func (api *CustomerMessageAPI) MarkMessageAsRead(c *gin.Context) {
	customerID, exists := c.Get("customer_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "未授权访问"})
		return
	}

	messageIDStr := c.Param("id")
	messageID, err := strconv.ParseUint(messageIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的消息ID"})
		return
	}

	messageService := &service.SimpleCustomerMessageService{}
	if err := messageService.MarkMessageAsRead(uint64(customerID.(uint)), messageID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "标记成功"})
}

// MarkMessagesAsReadBatch 批量标记消息为已读
// @Summary 批量标记消息为已读
// @Description 批量将消息标记为已读
// @Tags 第三方客户消息
// @Accept json
// @Produce json
// @Param message_ids body struct{MessageIDs []uint64 `json:"message_ids" binding:"required"`} true "消息ID列表"
// @Success 200 {object} gin.H{"message":"批量标记成功"}
// @Router /api/v1/customer/messages/batch-read [put]
// @Security CustomerJWTAuth
func (api *CustomerMessageAPI) MarkMessagesAsReadBatch(c *gin.Context) {
	customerID, exists := c.Get("customer_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "未授权访问"})
		return
	}

	var req struct {
		MessageIDs []uint64 `json:"message_ids" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	messageService := &service.SimpleCustomerMessageService{}
	for _, messageID := range req.MessageIDs {
		if err := messageService.MarkMessageAsRead(uint64(customerID.(uint)), messageID); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{"message": "批量标记成功"})
}

// GetUnreadCount 获取未读消息数量
// @Summary 获取未读消息数量
// @Description 获取当前登录的第三方客户的未读消息数量
// @Tags 第三方客户消息
// @Accept json
// @Produce json
// @Success 200 {object} gin.H{"count":int}
// @Router /api/v1/customer/messages/unread-count [get]
// @Security CustomerJWTAuth
func (api *CustomerMessageAPI) GetUnreadCount(c *gin.Context) {
	customerID, exists := c.Get("customer_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "未授权访问"})
		return
	}

	messageService := &service.SimpleCustomerMessageService{}
	count, err := messageService.GetUnreadCount(uint64(customerID.(uint)))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": gin.H{"count": count}})
}

// DeleteMessage 删除消息
// @Summary 删除消息
// @Description 删除指定的消息（软删除）
// @Tags 第三方客户消息
// @Accept json
// @Produce json
// @Param id path int true "消息ID"
// @Success 200 {object} gin.H{"message":"删除成功"}
// @Router /api/v1/customer/messages/{id} [delete]
// @Security CustomerJWTAuth
func (api *CustomerMessageAPI) DeleteMessage(c *gin.Context) {
	customerID, exists := c.Get("customer_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "未授权访问"})
		return
	}

	messageIDStr := c.Param("id")
	messageID, err := strconv.ParseUint(messageIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的消息ID"})
		return
	}

	// 确保消息存在且属于该客户
	var message model.Message
	if err := global.DB.Where("id = ? AND ((target_type = ? AND target_id = ?) OR target_type = ?)", 
		messageID, "customer", customerID, "all").First(&message).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "消息不存在或无权限访问"})
		return
	}

	// 软删除消息
	if err := global.DB.Delete(&model.Message{}, messageID).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "删除成功"})
}

// GetCustomerAnnouncements 获取客户公告列表
// @Summary 获取客户公告列表
// @Description 获取当前登录客户的公告列表
// @Tags 客户公告
// @Accept json
// @Produce json
// @Param page query int false "页码" default(1)
// @Param page_size query int false "每页数量" default(10) minimum(1) maximum(100)
// @Param is_read query bool false "是否已读"
// @Success 200 {object} service.CustomerAnnouncementListResponse
// @Router /api/v1/customer/announcements [get]
// @Security CustomerJWTAuth
func (api *CustomerMessageAPI) GetCustomerAnnouncements(c *gin.Context) {
	customerID, exists := c.Get("customer_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "未授权访问"})
		return
	}

	var req struct {
		Page     int  `form:"page" binding:"min=1"`
		PageSize int  `form:"page_size" binding:"min=1,max=100"`
		IsRead   *bool `form:"is_read"`
	}
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if req.Page == 0 {
		req.Page = 1
	}
	if req.PageSize == 0 {
		req.PageSize = 10
	}

	announcementService := &service.AnnouncementService{}
	response, err := announcementService.GetCustomerAnnouncements(uint64(customerID.(uint)), &service.CustomerAnnouncementListRequest{
		Page:     req.Page,
		PageSize: req.PageSize,
		IsRead:   req.IsRead,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": response})
}

// GetCustomerAnnouncementDetail 获取客户公告详情
// @Summary 获取客户公告详情
// @Description 获取指定公告的详细信息
// @Tags 客户公告
// @Accept json
// @Produce json
// @Param id path int true "公告ID"
// @Success 200 {object} service.CustomerAnnouncementDetailResponse
// @Router /api/v1/customer/announcements/{id} [get]
// @Security CustomerJWTAuth
func (api *CustomerMessageAPI) GetCustomerAnnouncementDetail(c *gin.Context) {
	customerID, exists := c.Get("customer_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "未授权访问"})
		return
	}

	announcementIDStr := c.Param("id")
	announcementID, err := strconv.ParseUint(announcementIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的公告ID"})
		return
	}

	announcementService := &service.AnnouncementService{}
	detail, err := announcementService.GetCustomerAnnouncementDetail(uint64(customerID.(uint)), announcementID)
	if err != nil {
		if err.Error() == "公告不存在" {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": detail})
}

// MarkAnnouncementRead 标记公告为已读
// @Summary 标记公告为已读
// @Description 将指定公告标记为已读
// @Tags 客户公告
// @Accept json
// @Produce json
// @Param id path int true "公告ID"
// @Success 200 {object} gin.H{"message":"标记成功"}
// @Router /api/v1/customer/announcements/{id}/read [post]
// @Security CustomerJWTAuth
func (api *CustomerMessageAPI) MarkAnnouncementRead(c *gin.Context) {
	customerID, exists := c.Get("customer_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "未授权访问"})
		return
	}

	announcementIDStr := c.Param("id")
	announcementID, err := strconv.ParseUint(announcementIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的公告ID"})
		return
	}

	announcementService := &service.AnnouncementService{}
	if err := announcementService.MarkAnnouncementRead(uint64(customerID.(uint)), announcementID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "标记成功"})
}

// MarkAnnouncementsBatchRead 批量标记公告为已读
// @Summary 批量标记公告为已读
// @Description 批量将公告标记为已读
// @Tags 客户公告
// @Accept json
// @Produce json
// @Param announcement_ids body struct{IDs []uint64} true "公告ID列表"
// @Success 200 {object} gin.H{"message":"批量标记成功"}
// @Router /api/v1/customer/announcements/batch-read [post]
// @Security CustomerJWTAuth
func (api *CustomerMessageAPI) MarkAnnouncementsBatchRead(c *gin.Context) {
	customerID, exists := c.Get("customer_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "未授权访问"})
		return
	}

	var req struct {
		IDs []uint64 `json:"ids" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	announcementService := &service.AnnouncementService{}
	if err := announcementService.MarkAnnouncementsBatchRead(uint64(customerID.(uint)), req.IDs); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "批量标记成功"})
}

// GetUnreadAnnouncementCount 获取未读公告数量
// @Summary 获取未读公告数量
// @Description 获取当前登录客户的未读公告数量
// @Tags 客户公告
// @Accept json
// @Produce json
// @Success 200 {object} gin.H{"count":int}
// @Router /api/v1/customer/announcements/unread-count [get]
// @Security CustomerJWTAuth
func (api *CustomerMessageAPI) GetUnreadAnnouncementCount(c *gin.Context) {
	customerID, exists := c.Get("customer_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "未授权访问"})
		return
	}

	announcementService := &service.AnnouncementService{}
	count, err := announcementService.GetUnreadAnnouncementCount(uint64(customerID.(uint)))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": gin.H{"count": count}})
}

// AdminCreateMessage 管理员创建消息
// @Summary 管理员创建消息
// @Description 管理员为第三方客户创建消息
// @Tags 第三方客户消息
// @Accept json
// @Produce json
// @Param message body struct{Title string `json:"title" binding:"required,max=255"`; Content string `json:"content" binding:"required"`; CustomerIDs []uint64 `json:"customer_ids" binding:"required"`; Priority string `json:"priority" binding:"omitempty,oneof=low medium high"`; ExpireTime time.Time `json:"expire_time"`} true "消息信息"
// @Success 200 {object} gin.H{"message":"创建成功"}
// @Router /api/v1/admin/customer-messages [post]
// @Security JWTAuth
func (api *CustomerMessageAPI) AdminCreateMessage(c *gin.Context) {
	var req struct {
		Title      string    `json:"title" binding:"required,max=255"`
		Content    string    `json:"content" binding:"required"`
		CustomerIDs []uint64 `json:"customer_ids" binding:"required"`
		Priority   string    `json:"priority" binding:"omitempty,oneof=low medium high"`
		ExpireTime time.Time `json:"expire_time"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	priority, _ := strconv.Atoi(req.Priority)
	
	// 为每个客户创建消息
	for _, customerID := range req.CustomerIDs {
		message := &model.Message{
			Title:      req.Title,
			Content:    req.Content,
			Type:       model.MessageTypePrivate,
			Priority:   priority,
			Status:     model.MessageStatusPublished,
			TargetType: "customer",
			TargetID:   &customerID,
			ExpiredAt:  &req.ExpireTime,
		}
		
		if err := global.DB.Create(message).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{"message": "创建成功"})
}

// AdminCreateMessagesBatch 管理员批量创建消息
// @Summary 管理员批量创建消息
// @Description 管理员为多个第三方客户批量创建消息
// @Tags 第三方客户消息
// @Accept json
// @Produce json
// @Param message body struct{Title string `json:"title" binding:"required,max=255"`; Content string `json:"content" binding:"required"`; CustomerIDs []uint64 `json:"customer_ids" binding:"required"`; Priority string `json:"priority" binding:"omitempty,oneof=low medium high"`; ExpireTime time.Time `json:"expire_time"`} true "消息信息"
// @Success 200 {object} gin.H{"message":"批量创建成功"}
// @Router /api/v1/admin/customer-messages/batch [post]
// @Security JWTAuth
func (api *CustomerMessageAPI) AdminCreateMessagesBatch(c *gin.Context) {
	var req struct {
		Title      string    `json:"title" binding:"required,max=255"`
		Content    string    `json:"content" binding:"required"`
		CustomerIDs []uint64 `json:"customer_ids" binding:"required"`
		Priority   string    `json:"priority" binding:"omitempty,oneof=low medium high"`
		ExpireTime time.Time `json:"expire_time"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	priority, _ := strconv.Atoi(req.Priority)
	
	// 批量创建消息
	messages := make([]*model.Message, len(req.CustomerIDs))
	for i, customerID := range req.CustomerIDs {
		messages[i] = &model.Message{
			Title:      req.Title,
			Content:    req.Content,
			Type:       model.MessageTypePrivate,
			Priority:   priority,
			Status:     model.MessageStatusPublished,
			TargetType: "customer",
			TargetID:   &customerID,
			ExpiredAt:  &req.ExpireTime,
		}
	}
	
	if err := global.DB.CreateInBatches(messages, 100).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "批量创建成功"})
}
