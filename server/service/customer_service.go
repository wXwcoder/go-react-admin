package service

import (
	"errors"
	"fmt"
	"go-react-admin/global"
	"go-react-admin/model"
	"time"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type CustomerService struct{}

type CustomerCreateRequest struct {
	Username  string  `json:"username" binding:"required,min=3,max=50"`
	Email     string  `json:"email" binding:"required,email,max=100"`
	Password  string  `json:"password" binding:"required,min=6,max=100"`
	Phone     *string `json:"phone,omitempty" binding:"omitempty,max=20"`
	RealName  *string `json:"real_name,omitempty" binding:"omitempty,max=50"`
	AvatarURL *string `json:"avatar_url,omitempty" binding:"omitempty,max=500"`
}

type CustomerUpdateRequest struct {
	Username  *string `json:"username,omitempty" binding:"omitempty,min=3,max=50"`
	Email     *string `json:"email,omitempty" binding:"omitempty,email,max=100"`
	Phone     *string `json:"phone,omitempty" binding:"omitempty,max=20"`
	RealName  *string `json:"real_name,omitempty" binding:"omitempty,max=50"`
	AvatarURL *string `json:"avatar_url,omitempty" binding:"omitempty,max=500"`
}

type CustomerStatusUpdateRequest struct {
	Status model.CustomerStatus `json:"status" binding:"required,oneof=active banned pending"`
}

type CustomerListRequest struct {
	Page     int    `form:"page" json:"page" binding:"min=1"`
	PageSize int    `form:"page_size" json:"page_size" binding:"min=1,max=100"`
	Keyword  string `form:"keyword,omitempty" json:"keyword,omitempty" binding:"max=100"`
	Status   string `form:"status,omitempty" json:"status,omitempty" binding:"omitempty,oneof=active banned pending"`
}

type CustomerListResponse struct {
	List  []model.Customer `json:"list"`
	Total int64            `json:"total"`
}

// CreateCustomer 创建第三方客户
func (s *CustomerService) CreateCustomer(req *CustomerCreateRequest) (*model.Customer, error) {
	// 检查用户名是否已存在
	var count int64
	global.DB.Model(&model.Customer{}).Where("username = ?", req.Username).Count(&count)
	if count > 0 {
		return nil, errors.New("用户名已存在")
	}

	// 检查邮箱是否已存在
	global.DB.Model(&model.Customer{}).Where("email = ?", req.Email).Count(&count)
	if count > 0 {
		return nil, errors.New("邮箱已存在")
	}

	// 加密密码
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("密码加密失败: %v", err)
	}

	customer := &model.Customer{
		Username:     req.Username,
		Email:        req.Email,
		PasswordHash: string(hashedPassword),
		Phone:        req.Phone,
		RealName:     req.RealName,
		AvatarURL:    req.AvatarURL,
		Status:       model.CustomerStatusActive,
	}

	if err := global.DB.Create(customer).Error; err != nil {
		return nil, fmt.Errorf("创建客户失败: %v", err)
	}

	return customer, nil
}

// UpdateCustomer 更新第三方客户信息
func (s *CustomerService) UpdateCustomer(id uint64, req *CustomerUpdateRequest) (*model.Customer, error) {
	customer, err := s.GetCustomerByID(id)
	if err != nil {
		return nil, err
	}

	// 检查用户名是否被其他用户使用
	if req.Username != nil && *req.Username != customer.Username {
		var count int64
		global.DB.Model(&model.Customer{}).Where("username = ? AND id != ?", *req.Username, id).Count(&count)
		if count > 0 {
			return nil, errors.New("用户名已存在")
		}
	}

	// 检查邮箱是否被其他用户使用
	if req.Email != nil && *req.Email != customer.Email {
		var count int64
		global.DB.Model(&model.Customer{}).Where("email = ? AND id != ?", *req.Email, id).Count(&count)
		if count > 0 {
			return nil, errors.New("邮箱已存在")
		}
	}

	updates := make(map[string]interface{})
	if req.Username != nil {
		updates["username"] = *req.Username
	}
	if req.Email != nil {
		updates["email"] = *req.Email
	}
	if req.Phone != nil {
		updates["phone"] = *req.Phone
	}
	if req.RealName != nil {
		updates["real_name"] = *req.RealName
	}
	if req.AvatarURL != nil {
		updates["avatar_url"] = *req.AvatarURL
	}

	if err := global.DB.Model(&model.Customer{}).Where("id = ?", id).Updates(updates).Error; err != nil {
		return nil, fmt.Errorf("更新客户失败: %v", err)
	}

	return s.GetCustomerByID(id)
}

// UpdateCustomerStatus 更新客户状态
func (s *CustomerService) UpdateCustomerStatus(id uint64, status model.CustomerStatus) (*model.Customer, error) {
	_, err := s.GetCustomerByID(id)
	if err != nil {
		return nil, err
	}

	if err := global.DB.Model(&model.Customer{}).Where("id = ?", id).Update("status", status).Error; err != nil {
		return nil, fmt.Errorf("更新客户状态失败: %v", err)
	}

	return s.GetCustomerByID(id)
}

// GetCustomerByID 根据ID获取客户信息
func (s *CustomerService) GetCustomerByID(id uint64) (*model.Customer, error) {
	var customer model.Customer
	if err := global.DB.Where("id = ?", id).First(&customer).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("客户不存在")
		}
		return nil, fmt.Errorf("查询客户失败: %v", err)
	}
	return &customer, nil
}

// GetCustomerList 获取客户列表
func (s *CustomerService) GetCustomerList(req *CustomerListRequest) (*CustomerListResponse, error) {
	var customers []model.Customer
	var total int64

	query := global.DB.Model(&model.Customer{})

	// 搜索条件
	if req.Keyword != "" {
		query = query.Where("username LIKE ? OR email LIKE ? OR real_name LIKE ? OR phone LIKE ?",
			"%"+req.Keyword+"%", "%"+req.Keyword+"%", "%"+req.Keyword+"%", "%"+req.Keyword+"%")
	}

	if req.Status != "" {
		query = query.Where("status = ?", req.Status)
	}

	// 统计总数
	query.Count(&total)

	// 分页查询
	offset := (req.Page - 1) * req.PageSize
	if err := query.Order("created_at DESC").Limit(req.PageSize).Offset(offset).Find(&customers).Error; err != nil {
		return nil, fmt.Errorf("查询客户列表失败: %v", err)
	}

	return &CustomerListResponse{
		List:  customers,
		Total: total,
	}, nil
}

// DeleteCustomer 删除客户（软删除）
func (s *CustomerService) DeleteCustomer(id uint64) error {
	_, err := s.GetCustomerByID(id)
	if err != nil {
		return err
	}

	if err := global.DB.Delete(&model.Customer{}, id).Error; err != nil {
		return fmt.Errorf("删除客户失败: %v", err)
	}

	return nil
}

// BanCustomer 封禁客户
func (s *CustomerService) BanCustomer(id uint64) (*model.Customer, error) {
	return s.UpdateCustomerStatus(id, model.CustomerStatusBanned)
}

// UnbanCustomer 解封客户
func (s *CustomerService) UnbanCustomer(id uint64) (*model.Customer, error) {
	return s.UpdateCustomerStatus(id, model.CustomerStatusActive)
}

// ValidateCustomerPassword 验证客户密码
func (s *CustomerService) ValidateCustomerPassword(customer *model.Customer, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(customer.PasswordHash), []byte(password))
	return err == nil
}

// UpdateLoginInfo 更新登录信息
func (s *CustomerService) UpdateLoginInfo(id uint64) error {
	return global.DB.Model(&model.Customer{}).Where("id = ?", id).Updates(map[string]interface{}{
		"last_login_at": time.Now(),
		"login_count":   gorm.Expr("login_count + 1"),
	}).Error
}
