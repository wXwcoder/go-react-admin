package api

import (
	"go-react-admin/service"
	"net/http"
	"strconv"

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
// @Param message_type query string false "消息类型"
// @Success 200 {object} service.CustomerMessageListResponse
// @Router /api/v1/customer/messages [get]
// @Security CustomerJWTAuth
func (api *CustomerMessageAPI) GetMessages(c *gin.Context) {
	customerID, exists := c.Get("customer_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "未授权访问"})
		return
	}

	var req service.CustomerMessageListRequest
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

	messageService := &service.CustomerMessageService{}
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
// @Success 200 {object} service.CustomerMessageDetailResponse
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

	messageService := &service.CustomerMessageService{}
	response, err := messageService.GetCustomerMessageDetail(customerID.(uint64), messageID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": response})
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

	messageService := &service.CustomerMessageService{}
	if err := messageService.MarkMessageAsRead(customerID.(uint64), messageID); err != nil {
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
// @Param message_ids body service.CustomerMessageBatchReadRequest true "消息ID列表"
// @Success 200 {object} gin.H{"message":"批量标记成功"}
// @Router /api/v1/customer/messages/batch-read [put]
// @Security CustomerJWTAuth
func (api *CustomerMessageAPI) MarkMessagesAsReadBatch(c *gin.Context) {
	customerID, exists := c.Get("customer_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "未授权访问"})
		return
	}

	var req service.CustomerMessageBatchReadRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	messageService := &service.CustomerMessageService{}
	if err := messageService.MarkMessagesAsReadBatch(customerID.(uint64), req.MessageIDs); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
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

	messageService := &service.CustomerMessageService{}
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

	messageService := &service.CustomerMessageService{}
	if err := messageService.DeleteCustomerMessage(customerID.(uint64), messageID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "删除成功"})
}

// AdminCreateMessage 管理员创建消息
// @Summary 管理员创建消息
// @Description 管理员为第三方客户创建消息
// @Tags 第三方客户消息
// @Accept json
// @Produce json
// @Param message body service.AdminCreateCustomerMessageRequest true "消息信息"
// @Success 200 {object} gin.H{"message":"创建成功"}
// @Router /api/v1/admin/customer-messages [post]
// @Security JWTAuth
func (api *CustomerMessageAPI) AdminCreateMessage(c *gin.Context) {
	var req service.AdminCreateCustomerMessageRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	messageService := &service.CustomerMessageService{}
	// 使用第一个客户ID创建消息，实际应该先生成消息再关联
	if len(req.CustomerIDs) > 0 {
		if err := messageService.CreateCustomerMessage(req.CustomerIDs[0], 0); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请选择接收消息的客户"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "创建成功"})
}

// AdminCreateMessagesBatch 管理员批量创建消息
// @Summary 管理员批量创建消息
// @Description 管理员为多个第三方客户批量创建消息
// @Tags 第三方客户消息
// @Accept json
// @Produce json
// @Param message body service.AdminCreateCustomerMessagesBatchRequest true "消息信息"
// @Success 200 {object} gin.H{"message":"批量创建成功"}
// @Router /api/v1/admin/customer-messages/batch [post]
// @Security JWTAuth
func (api *CustomerMessageAPI) AdminCreateMessagesBatch(c *gin.Context) {
	var req service.AdminCreateCustomerMessagesBatchRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	messageService := &service.CustomerMessageService{}
	if err := messageService.CreateCustomerMessagesBatch(req.CustomerIDs, 0); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "批量创建成功"})
}
