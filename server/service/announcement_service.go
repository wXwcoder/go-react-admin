package service

import (
	"errors"
	"fmt"
	"go-react-admin/global"
	"go-react-admin/model"
	"time"

	"gorm.io/gorm"
)

type AnnouncementService struct{}

type AnnouncementListRequest struct {
	Page     int    `form:"page" json:"page" binding:"min=1"`
	PageSize int    `form:"page_size" json:"page_size" binding:"min=1,max=100"`
	Type     string `form:"type,omitempty" json:"type,omitempty"`
	Status   string `form:"status,omitempty" json:"status,omitempty"`
	Keyword  string `form:"keyword,omitempty" json:"keyword,omitempty" binding:"max=100"`
}

type AnnouncementCreateRequest struct {
	Title     string     `json:"title" binding:"required,max=255"`
	Content   string     `json:"content" binding:"required,max=5000"`
	Type      string     `json:"type" binding:"required,oneof=system notice maintenance update"`
	Priority  int        `json:"priority" binding:"min=0,max=10"`
	ExpiredAt *time.Time `json:"expired_at,omitempty"`
}

type AnnouncementUpdateRequest struct {
	Title     *string     `json:"title,omitempty" binding:"omitempty,max=255"`
	Content   *string     `json:"content,omitempty" binding:"omitempty,max=5000"`
	Type      *string     `json:"type,omitempty" binding:"omitempty,oneof=system notice maintenance update"`
	Priority  *int        `json:"priority,omitempty" binding:"omitempty,min=0,max=10"`
	ExpiredAt *time.Time `json:"expired_at,omitempty"`
}

type AnnouncementListResponse struct {
	List  []AnnouncementDetail `json:"list"`
	Total int64                `json:"total"`
}

type AnnouncementDetail struct {
	ID        uint64                `json:"id"`
	Title     string                `json:"title"`
	Content   string                `json:"content"`
	Type      string                `json:"type"`
	Status    string                `json:"status"`
	Priority  int                   `json:"priority"`
	ExpiredAt *time.Time            `json:"expired_at,omitempty"`
	ReadCount int                   `json:"read_count"`
	CreatedAt time.Time             `json:"created_at"`
	UpdatedAt time.Time             `json:"updated_at"`
	Stats     *AnnouncementStats    `json:"stats,omitempty"`
}

type AnnouncementStats struct {
	TotalCustomers int64 `json:"total_customers"`
	ReadCustomers  int64 `json:"read_customers"`
	ReadRate       float64 `json:"read_rate"`
}

// AnnouncementStatsResponse 公告统计响应
type AnnouncementStatsResponse struct {
	TotalCount    int64 `json:"total_count"`    // 总公告数
	DraftCount    int64 `json:"draft_count"`    // 草稿公告数
	PublishedCount int64 `json:"published_count"` // 已发布公告数
	RevokedCount  int64 `json:"revoked_count"`  // 已撤回公告数
	SystemCount   int64 `json:"system_count"`   // 系统公告数
	NoticeCount   int64 `json:"notice_count"`   // 通知公告数
	MaintenanceCount int64 `json:"maintenance_count"` // 维护公告数
	UpdateCount   int64 `json:"update_count"`   // 更新公告数
}

type CustomerAnnouncementListRequest struct {
	Page     int    `form:"page" json:"page" binding:"min=1"`
	PageSize int    `form:"page_size" json:"page_size" binding:"min=1,max=100"`
	IsRead   *bool  `form:"is_read,omitempty" json:"is_read,omitempty"`
}

type CustomerAnnouncementListResponse struct {
	List  []CustomerAnnouncementDetail `json:"list"`
	Total int64                        `json:"total"`
}

type CustomerAnnouncementDetail struct {
	ID        uint64     `json:"id"`
	Title     string     `json:"title"`
	Content   string     `json:"content"`
	Type      string     `json:"type"`
	Priority  int        `json:"priority"`
	IsRead    bool       `json:"is_read"`
	ReadTime  *time.Time `json:"read_time,omitempty"`
	ExpiredAt *time.Time `json:"expired_at,omitempty"`
	CreatedAt time.Time  `json:"created_at"`
}

// GetAnnouncementList 获取公告列表（管理员）
func (s *AnnouncementService) GetAnnouncementList(req *AnnouncementListRequest) (*AnnouncementListResponse, error) {
	var announcements []AnnouncementDetail
	var total int64

	query := global.DB.Model(&model.Announcement{}).Where("deleted_at IS NULL")

	if req.Type != "" {
		query = query.Where("type = ?", req.Type)
	}

	if req.Status != "" {
		query = query.Where("status = ?", req.Status)
	}

	if req.Keyword != "" {
		query = query.Where("title LIKE ? OR content LIKE ?", "%"+req.Keyword+"%", "%"+req.Keyword+"%")
	}

	// 统计总数
	query.Count(&total)

	// 分页查询
	offset := (req.Page - 1) * req.PageSize
	if err := query.Order("created_at DESC").Limit(req.PageSize).Offset(offset).Find(&announcements).Error; err != nil {
		return nil, fmt.Errorf("查询公告列表失败: %v", err)
	}

	// 获取统计信息
	for i := range announcements {
		stats, err := s.getAnnouncementStats(announcements[i].ID)
		if err != nil {
			continue
		}
		announcements[i].Stats = stats
	}

	return &AnnouncementListResponse{
		List:  announcements,
		Total: total,
	}, nil
}

// GetAnnouncementDetail 获取公告详情
func (s *AnnouncementService) GetAnnouncementDetail(id uint64) (*AnnouncementDetail, error) {
	var announcement AnnouncementDetail

	if err := global.DB.Model(&model.Announcement{}).Where("id = ? AND deleted_at IS NULL", id).First(&announcement).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("公告不存在")
		}
		return nil, fmt.Errorf("查询公告详情失败: %v", err)
	}

	stats, err := s.getAnnouncementStats(id)
	if err == nil {
		announcement.Stats = stats
	}

	return &announcement, nil
}

// CreateAnnouncement 创建公告
func (s *AnnouncementService) CreateAnnouncement(req *AnnouncementCreateRequest) (*model.Announcement, error) {
	announcement := &model.Announcement{
		Title:     req.Title,
		Content:   req.Content,
		Type:      model.AnnouncementType(req.Type),
		Priority:  req.Priority,
		ExpiredAt: req.ExpiredAt,
		Status:    model.AnnouncementStatusDraft,
	}

	if err := global.DB.Create(announcement).Error; err != nil {
		return nil, fmt.Errorf("创建公告失败: %v", err)
	}

	return announcement, nil
}

// UpdateAnnouncement 更新公告
func (s *AnnouncementService) UpdateAnnouncement(id uint64, req *AnnouncementUpdateRequest) error {
	updates := make(map[string]interface{})

	if req.Title != nil {
		updates["title"] = *req.Title
	}
	if req.Content != nil {
		updates["content"] = *req.Content
	}
	if req.Type != nil {
		updates["type"] = *req.Type
	}
	if req.Priority != nil {
		updates["priority"] = *req.Priority
	}
	if req.ExpiredAt != nil {
		updates["expired_at"] = *req.ExpiredAt
	}

	if len(updates) == 0 {
		return nil
	}

	updates["updated_at"] = time.Now()

	if err := global.DB.Model(&model.Announcement{}).Where("id = ? AND deleted_at IS NULL", id).Updates(updates).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("公告不存在")
		}
		return fmt.Errorf("更新公告失败: %v", err)
	}

	return nil
}

// PublishAnnouncement 发布公告
func (s *AnnouncementService) PublishAnnouncement(id uint64) error {
	// 开始事务
	err := global.DB.Transaction(func(tx *gorm.DB) error {
		// 更新公告状态为已发布
		if err := tx.Model(&model.Announcement{}).Where("id = ? AND deleted_at IS NULL", id).Update("status", model.AnnouncementStatusPublished).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return errors.New("公告不存在")
			}
			return fmt.Errorf("发布公告失败: %v", err)
		}

		// 获取所有活跃客户
		var customers []model.Customer
		if err := tx.Where("status = ? AND deleted_at IS NULL", model.CustomerStatusActive).Find(&customers).Error; err != nil {
			return fmt.Errorf("获取客户列表失败: %v", err)
		}

		if len(customers) > 0 {
			// 为每个客户创建公告阅读记录
			var reads []model.AnnouncementRead
			for _, customer := range customers {
				reads = append(reads, model.AnnouncementRead{
					CustomerID:     customer.ID,
					AnnouncementID: id,
					IsRead:         false,
				})
			}

			// 批量创建阅读记录
			if err := tx.CreateInBatches(reads, 100).Error; err != nil {
				return fmt.Errorf("创建客户阅读记录失败: %v", err)
			}
		}

		return nil
	})

	return err
}

// RevokeAnnouncement 撤回公告
func (s *AnnouncementService) RevokeAnnouncement(id uint64) error {
	if err := global.DB.Model(&model.Announcement{}).Where("id = ? AND deleted_at IS NULL", id).Update("status", model.AnnouncementStatusRevoked).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("公告不存在")
		}
		return fmt.Errorf("撤回公告失败: %v", err)
	}
	return nil
}

// DeleteAnnouncement 删除公告（软删除）
func (s *AnnouncementService) DeleteAnnouncement(id uint64) error {
	if err := global.DB.Model(&model.Announcement{}).Where("id = ?", id).Update("deleted_at", time.Now()).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("公告不存在")
		}
		return fmt.Errorf("删除公告失败: %v", err)
	}
	return nil
}

// GetCustomerAnnouncementDetail 获取客户公告详情
func (s *AnnouncementService) GetCustomerAnnouncementDetail(customerID uint64, announcementID uint64) (*CustomerAnnouncementDetail, error) {
	var detail CustomerAnnouncementDetail
	now := time.Now()

	err := global.DB.Table("announcements as a").
		Select("a.id, a.title, a.content, a.type, a.priority, a.expired_at, a.created_at, ar.is_read, ar.read_time").
		Joins("LEFT JOIN announcement_reads ar ON a.id = ar.announcement_id AND ar.customer_id = ?", customerID).
		Where("a.id = ? AND a.status = ? AND a.deleted_at IS NULL AND (a.expired_at IS NULL OR a.expired_at > ?)",
			announcementID, model.AnnouncementStatusPublished, now).
		First(&detail).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("公告不存在")
		}
		return nil, fmt.Errorf("查询客户公告详情失败: %v", err)
	}

	return &detail, nil
}

// GetCustomerAnnouncements 获取客户公告列表
func (s *AnnouncementService) GetCustomerAnnouncements(customerID uint64, req *CustomerAnnouncementListRequest) (*CustomerAnnouncementListResponse, error) {
	var announcements []CustomerAnnouncementDetail
	var total int64

	// 获取当前时间
	now := time.Now()

	// 查询公告列表
	query := global.DB.Table("announcements as a").
		Select("a.id, a.title, a.content, a.type, a.priority, a.expired_at, a.created_at, ar.is_read, ar.read_time").
		Joins("LEFT JOIN announcement_reads ar ON a.id = ar.announcement_id AND ar.customer_id = ?", customerID).
		Where("a.status = ? AND a.deleted_at IS NULL AND (a.expired_at IS NULL OR a.expired_at > ?)", model.AnnouncementStatusPublished, now)

	if req.IsRead != nil {
		if *req.IsRead {
			query = query.Where("ar.is_read = ?", true)
		} else {
			query = query.Where("(ar.is_read IS NULL OR ar.is_read = ?)", false)
		}
	}

	// 统计总数
	query.Count(&total)

	// 分页查询
	offset := (req.Page - 1) * req.PageSize
	if err := query.Order("a.created_at DESC").Limit(req.PageSize).Offset(offset).Find(&announcements).Error; err != nil {
		return nil, fmt.Errorf("查询客户公告列表失败: %v", err)
	}

	return &CustomerAnnouncementListResponse{
		List:  announcements,
		Total: total,
	}, nil
}

// MarkAnnouncementAsRead 标记公告为已读
func (s *AnnouncementService) MarkAnnouncementAsRead(customerID uint64, announcementID uint64) error {
	var readRecord model.AnnouncementRead

	err := global.DB.Where("customer_id = ? AND announcement_id = ?", customerID, announcementID).First(&readRecord).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// 创建已读记录
			readRecord = model.AnnouncementRead{
				CustomerID:     customerID,
				AnnouncementID: announcementID,
				IsRead:         true,
			}
			return global.DB.Create(&readRecord).Error
		}
		return fmt.Errorf("查询已读记录失败: %v", err)
	}

	if !readRecord.IsRead {
		now := time.Now()
		return global.DB.Model(&model.AnnouncementRead{}).
			Where("customer_id = ? AND announcement_id = ?", customerID, announcementID).
			Updates(map[string]interface{}{
				"is_read":   true,
				"read_time": now,
			}).Error
	}

	return nil
}

// MarkAnnouncementsAsReadBatch 批量标记公告为已读
func (s *AnnouncementService) MarkAnnouncementsAsReadBatch(customerID uint64, announcementIDs []uint64) error {
	return global.DB.Transaction(func(tx *gorm.DB) error {
		for _, announcementID := range announcementIDs {
			var readRecord model.AnnouncementRead
			err := tx.Where("customer_id = ? AND announcement_id = ?", customerID, announcementID).First(&readRecord).Error
			if err != nil {
				if errors.Is(err, gorm.ErrRecordNotFound) {
					// 创建已读记录
					readRecord = model.AnnouncementRead{
						CustomerID:     customerID,
						AnnouncementID: announcementID,
						IsRead:         true,
					}
					if err := tx.Create(&readRecord).Error; err != nil {
						return err
					}
					continue
				}
				return err
			}

			if !readRecord.IsRead {
				now := time.Now()
				if err := tx.Model(&model.AnnouncementRead{}).
					Where("customer_id = ? AND announcement_id = ?", customerID, announcementID).
					Updates(map[string]interface{}{
						"is_read":   true,
						"read_time": now,
					}).Error; err != nil {
					return err
				}
			}
		}
		return nil
	})
}

// MarkAnnouncementRead 标记公告为已读（别名）
func (s *AnnouncementService) MarkAnnouncementRead(customerID uint64, announcementID uint64) error {
	return s.MarkAnnouncementAsRead(customerID, announcementID)
}

// MarkAnnouncementsBatchRead 批量标记公告为已读（别名）
func (s *AnnouncementService) MarkAnnouncementsBatchRead(customerID uint64, announcementIDs []uint64) error {
	return s.MarkAnnouncementsAsReadBatch(customerID, announcementIDs)
}

// GetUnreadAnnouncementCount 获取未读公告数量
func (s *AnnouncementService) GetUnreadAnnouncementCount(customerID uint64) (int64, error) {
	var count int64
	now := time.Now()

	err := global.DB.Table("announcements as a").
		Joins("LEFT JOIN announcement_reads ar ON a.id = ar.announcement_id AND ar.customer_id = ?", customerID).
		Where("a.status = ? AND a.deleted_at IS NULL AND (a.expired_at IS NULL OR a.expired_at > ?) AND (ar.is_read IS NULL OR ar.is_read = ?)",
			model.AnnouncementStatusPublished, now, false).
		Count(&count).Error

	return count, err
}

// getAnnouncementStats 获取公告统计信息
func (s *AnnouncementService) getAnnouncementStats(announcementID uint64) (*AnnouncementStats, error) {
	var totalCustomers int64
	var readCustomers int64

	// 获取总客户数
	if err := global.DB.Model(&model.Customer{}).Where("status = ? AND deleted_at IS NULL", model.CustomerStatusActive).Count(&totalCustomers).Error; err != nil {
		return nil, err
	}

	// 获取已读客户数
	if err := global.DB.Model(&model.AnnouncementRead{}).
		Where("announcement_id = ? AND is_read = ?", announcementID, true).
		Count(&readCustomers).Error; err != nil {
		return nil, err
	}

	readRate := float64(0)
	if totalCustomers > 0 {
		readRate = float64(readCustomers) / float64(totalCustomers) * 100
	}

	return &AnnouncementStats{
		TotalCustomers: totalCustomers,
		ReadCustomers:  readCustomers,
		ReadRate:       readRate,
	}, nil
}

// GetAnnouncementStats 获取公告统计信息
func (s *AnnouncementService) GetAnnouncementStats() (*AnnouncementStatsResponse, error) {
	var stats AnnouncementStatsResponse
	
	// 总公告数
	global.DB.Model(&model.Announcement{}).Where("deleted_at IS NULL").Count(&stats.TotalCount)
	
	// 草稿公告数
	global.DB.Model(&model.Announcement{}).Where("deleted_at IS NULL AND status = ?", model.AnnouncementStatusDraft).Count(&stats.DraftCount)
	
	// 已发布公告数
	global.DB.Model(&model.Announcement{}).Where("deleted_at IS NULL AND status = ?", model.AnnouncementStatusPublished).Count(&stats.PublishedCount)
	
	// 已撤回公告数
	global.DB.Model(&model.Announcement{}).Where("deleted_at IS NULL AND status = ?", model.AnnouncementStatusRevoked).Count(&stats.RevokedCount)
	
	// 系统公告数
	global.DB.Model(&model.Announcement{}).Where("deleted_at IS NULL AND type = ?", model.AnnouncementTypeSystem).Count(&stats.SystemCount)
	
	// 通知公告数
	global.DB.Model(&model.Announcement{}).Where("deleted_at IS NULL AND type = ?", model.AnnouncementTypeNotice).Count(&stats.NoticeCount)
	
	// 维护公告数
	global.DB.Model(&model.Announcement{}).Where("deleted_at IS NULL AND type = ?", model.AnnouncementTypeMaintenance).Count(&stats.MaintenanceCount)
	
	// 更新公告数
	global.DB.Model(&model.Announcement{}).Where("deleted_at IS NULL AND type = ?", model.AnnouncementTypeUpdate).Count(&stats.UpdateCount)
	
	return &stats, nil
}