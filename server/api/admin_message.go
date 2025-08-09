package api

import (
	"go-react-admin/service"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type AdminMessageAPI struct{}

// GetMessages 获取消息列表（管理员）
// @Summary 获取消息列表
// @Description 获取系统消息和公告列表
// @Tags 管理员消息管理
// @Accept json
// @Produce json
// @Param page query int false "页码" default(1)
// @Param page_size query int false "每页数量" default(10) minimum(1) maximum(100)
// @Param type query string false "消息类型" Enums(system,notice,private)
// @Param status query string false "消息状态" Enums(draft,published,revoked)
// @Param keyword query string false "关键词搜索"
// @Success 200 {object} service.MessageListResponse
// @Router /api/v1/admin/messages [get]
// @Security JWTAuth
func (api *AdminMessageAPI) GetMessages(c *gin.Context) {
	var req service.MessageListRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	messageService := &service.MessageService{}
	response, err := messageService.GetAdminMessageList(&req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": response})
}

// GetMessageDetail 获取消息详情（管理员）
// @Summary 获取消息详情
// @Description 获取指定消息的详细信息
// @Tags 管理员消息管理
// @Accept json
// @Produce json
// @Param id path int true "消息ID"
// @Success 200 {object} service.MessageDetailResponse
// @Router /api/v1/admin/messages/{id} [get]
// @Security JWTAuth
func (api *AdminMessageAPI) GetMessageDetail(c *gin.Context) {
	messageIDStr := c.Param("id")
	messageID, err := strconv.ParseUint(messageIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的消息ID"})
		return
	}

	messageService := &service.MessageService{}
	response, err := messageService.GetAdminMessageDetail(messageID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": response})
}

// CreateMessage 创建消息（管理员）
// @Summary 创建消息
// @Description 创建新的消息或公告
// @Tags 管理员消息管理
// @Accept json
// @Produce json
// @Param message body service.MessageCreateRequest true "消息信息"
// @Success 200 {object} gin.H{"message":"创建成功","data":model.Message}
// @Router /api/v1/admin/messages [post]
// @Security JWTAuth
func (api *AdminMessageAPI) CreateMessage(c *gin.Context) {
	var req service.MessageCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	messageService := &service.MessageService{}
	message, err := messageService.CreateMessage(&req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "创建成功", "data": message})
}

// UpdateMessage 更新消息（管理员）
// @Summary 更新消息
// @Description 更新现有消息信息
// @Tags 管理员消息管理
// @Accept json
// @Produce json
// @Param id path int true "消息ID"
// @Param message body service.MessageUpdateRequest true "更新信息"
// @Success 200 {object} gin.H{"message":"更新成功"}
// @Router /api/v1/admin/messages/{id} [put]
// @Security JWTAuth
func (api *AdminMessageAPI) UpdateMessage(c *gin.Context) {
	messageIDStr := c.Param("id")
	messageID, err := strconv.ParseUint(messageIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的消息ID"})
		return
	}

	var req service.MessageUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	messageService := &service.MessageService{}
	if err := messageService.UpdateMessage(messageID, &req); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "更新成功"})
}

// DeleteMessage 删除消息（管理员）
// @Summary 删除消息
// @Description 删除指定消息（软删除）
// @Tags 管理员消息管理
// @Accept json
// @Produce json
// @Param id path int true "消息ID"
// @Success 200 {object} gin.H{"message":"删除成功"}
// @Router /api/v1/admin/messages/{id} [delete]
// @Security JWTAuth
func (api *AdminMessageAPI) DeleteMessage(c *gin.Context) {
	messageIDStr := c.Param("id")
	messageID, err := strconv.ParseUint(messageIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的消息ID"})
		return
	}

	messageService := &service.MessageService{}
	if err := messageService.DeleteMessage(messageID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "删除成功"})
}

// SendMessage 发送消息（管理员）
// @Summary 发送消息
// @Description 发送消息给客户或发布为公告
// @Tags 管理员消息管理
// @Accept json
// @Produce json
// @Param id path int true "消息ID"
// @Success 200 {object} gin.H{"message":"发送成功"}
// @Router /api/v1/admin/messages/{id}/send [post]
// @Security JWTAuth
func (api *AdminMessageAPI) SendMessage(c *gin.Context) {
	messageIDStr := c.Param("id")
	messageID, err := strconv.ParseUint(messageIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的消息ID"})
		return
	}

	messageService := &service.MessageService{}
	if err := messageService.SendMessage(messageID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "发送成功"})
}

// CancelMessage 撤回消息（管理员）
// @Summary 撤回消息
// @Description 撤回已发布的消息
// @Tags 管理员消息管理
// @Accept json
// @Produce json
// @Param id path int true "消息ID"
// @Success 200 {object} gin.H{"message":"撤回成功"}
// @Router /api/v1/admin/messages/{id}/cancel [post]
// @Security JWTAuth
func (api *AdminMessageAPI) CancelMessage(c *gin.Context) {
	messageIDStr := c.Param("id")
	messageID, err := strconv.ParseUint(messageIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的消息ID"})
		return
	}

	messageService := &service.MessageService{}
	if err := messageService.CancelMessage(messageID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "撤回成功"})
}

// GetCustomerMessages 获取客户专属消息列表（管理员）
// @Summary 获取客户专属消息列表
// @Description 获取发送给客户的专属消息列表
// @Tags 管理员消息管理
// @Accept json
// @Produce json
// @Param page query int false "页码" default(1)
// @Param page_size query int false "每页数量" default(10) minimum(1) maximum(100)
// @Param customer_id query int false "客户ID"
// @Param is_read query bool false "是否已读"
// @Success 200 {object} service.CustomerMessageAdminListResponse
// @Router /api/v1/admin/customer-messages [get]
// @Security JWTAuth
func (api *AdminMessageAPI) GetCustomerMessages(c *gin.Context) {
	var req service.CustomerMessageAdminListRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	messageService := &service.MessageService{}
	response, err := messageService.GetCustomerMessageAdminList(&req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": response})
}

// SendCustomerMessage 发送客户专属消息（管理员）
// @Summary 发送客户专属消息
// @Description 向指定客户发送专属消息
// @Tags 管理员消息管理
// @Accept json
// @Produce json
// @Param message body service.CustomerMessageSendRequest true "消息信息"
// @Success 200 {object} gin.H{"message":"发送成功"}
// @Router /api/v1/admin/customer-messages/send [post]
// @Security JWTAuth
func (api *AdminMessageAPI) SendCustomerMessage(c *gin.Context) {
	var req service.CustomerMessageSendRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	messageService := &service.MessageService{}
	if err := messageService.SendCustomerMessage(&req); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "发送成功"})
}

// GetAnnouncementList 获取公告列表（管理员）
// @Summary 获取公告列表
// @Description 获取系统公告列表
// @Tags 管理员公告管理
// @Accept json
// @Produce json
// @Param page query int false "页码" default(1)
// @Param page_size query int false "每页数量" default(10) minimum(1) maximum(100)
// @Param type query string false "公告类型" Enums(system,notice,maintenance,update)
// @Param status query string false "公告状态" Enums(draft,published,revoked)
// @Param keyword query string false "关键词搜索"
// @Success 200 {object} service.AnnouncementListResponse
// @Router /api/v1/admin/announcements [get]
// @Security JWTAuth
func (api *AdminMessageAPI) GetAnnouncementList(c *gin.Context) {
	var req service.AnnouncementListRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	announcementService := &service.AnnouncementService{}
	response, err := announcementService.GetAnnouncementList(&req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": response})
}

// CreateAnnouncement 创建公告（管理员）
// @Summary 创建公告
// @Description 创建新的系统公告
// @Tags 管理员公告管理
// @Accept json
// @Produce json
// @Param announcement body service.AnnouncementCreateRequest true "公告信息"
// @Success 200 {object} gin.H{"message":"创建成功","data":model.Announcement}
// @Router /api/v1/admin/announcements [post]
// @Security JWTAuth
func (api *AdminMessageAPI) CreateAnnouncement(c *gin.Context) {
	var req service.AnnouncementCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	announcementService := &service.AnnouncementService{}
	announcement, err := announcementService.CreateAnnouncement(&req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "创建成功", "data": announcement})
}

// UpdateAnnouncement 更新公告（管理员）
// @Summary 更新公告
// @Description 更新现有公告信息
// @Tags 管理员公告管理
// @Accept json
// @Produce json
// @Param id path int true "公告ID"
// @Param announcement body service.AnnouncementUpdateRequest true "更新信息"
// @Success 200 {object} gin.H{"message":"更新成功"}
// @Router /api/v1/admin/announcements/{id} [put]
// @Security JWTAuth
func (api *AdminMessageAPI) UpdateAnnouncement(c *gin.Context) {
	announcementIDStr := c.Param("id")
	announcementID, err := strconv.ParseUint(announcementIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的公告ID"})
		return
	}

	var req service.AnnouncementUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	announcementService := &service.AnnouncementService{}
	if err := announcementService.UpdateAnnouncement(announcementID, &req); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "更新成功"})
}

// PublishAnnouncement 发布公告（管理员）
// @Summary 发布公告
// @Description 发布系统公告
// @Tags 管理员公告管理
// @Accept json
// @Produce json
// @Param id path int true "公告ID"
// @Success 200 {object} gin.H{"message":"发布成功"}
// @Router /api/v1/admin/announcements/{id}/publish [post]
// @Security JWTAuth
func (api *AdminMessageAPI) PublishAnnouncement(c *gin.Context) {
	announcementIDStr := c.Param("id")
	announcementID, err := strconv.ParseUint(announcementIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的公告ID"})
		return
	}

	announcementService := &service.AnnouncementService{}
	if err := announcementService.PublishAnnouncement(announcementID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "发布成功"})
}

// RevokeAnnouncement 撤回公告（管理员）
// @Summary 撤回公告
// @Description 撤回已发布的公告
// @Tags 管理员公告管理
// @Accept json
// @Produce json
// @Param id path int true "公告ID"
// @Success 200 {object} gin.H{"message":"撤回成功"}
// @Router /api/v1/admin/announcements/{id}/revoke [post]
// @Security JWTAuth
func (api *AdminMessageAPI) RevokeAnnouncement(c *gin.Context) {
	announcementIDStr := c.Param("id")
	announcementID, err := strconv.ParseUint(announcementIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的公告ID"})
		return
	}

	announcementService := &service.AnnouncementService{}
	if err := announcementService.RevokeAnnouncement(announcementID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "撤回成功"})
}

// DeleteAnnouncement 删除公告（管理员）
// @Summary 删除公告
// @Description 删除系统公告（软删除）
// @Tags 管理员公告管理
// @Accept json
// @Produce json
// @Param id path int true "公告ID"
// @Success 200 {object} gin.H{"message":"删除成功"}
// @Router /api/v1/admin/announcements/{id} [delete]
// @Security JWTAuth
func (api *AdminMessageAPI) DeleteAnnouncement(c *gin.Context) {
	announcementIDStr := c.Param("id")
	announcementID, err := strconv.ParseUint(announcementIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的公告ID"})
		return
	}

	announcementService := &service.AnnouncementService{}
	if err := announcementService.DeleteAnnouncement(announcementID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "删除成功"})
}