package service

import (
	"errors"
	"fmt"
	"go-react-admin/global"
	"go-react-admin/model"
	"go-react-admin/utils"

	"golang.org/x/crypto/bcrypt"
)

type CustomerAuthService struct{}

type CustomerRegisterRequest struct {
	Username string  `json:"username" binding:"required,min=3,max=50"`
	Email    string  `json:"email" binding:"required,email,max=100"`
	Password string  `json:"password" binding:"required,min=6,max=100"`
	Phone    *string `json:"phone,omitempty" binding:"omitempty,max=20"`
	RealName *string `json:"real_name,omitempty" binding:"omitempty,max=50"`
}

type CustomerLoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type CustomerLoginResponse struct {
	Customer  *model.Customer `json:"customer"`
	Token     string          `json:"token"`
	ExpiresIn int64           `json:"expires_in"`
}

type CustomerPasswordResetRequest struct {
	Email string `json:"email" binding:"required,email"`
}

type CustomerPasswordUpdateRequest struct {
	OldPassword string `json:"old_password" binding:"required"`
	NewPassword string `json:"new_password" binding:"required,min=6,max=100"`
}

type CustomerProfileUpdateRequest struct {
	Username  *string `json:"username,omitempty" binding:"omitempty,min=3,max=50"`
	Email     *string `json:"email,omitempty" binding:"omitempty,email,max=100"`
	Phone     *string `json:"phone,omitempty" binding:"omitempty,max=20"`
	RealName  *string `json:"real_name,omitempty" binding:"omitempty,max=50"`
	AvatarURL *string `json:"avatar_url,omitempty" binding:"omitempty,url,max=500"`
}

type CustomerTokenRefreshRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

// Register 客户注册
func (s *CustomerAuthService) Register(req *CustomerRegisterRequest) (*model.Customer, error) {
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
		Status:       model.CustomerStatusActive,
	}

	if err := global.DB.Create(customer).Error; err != nil {
		return nil, fmt.Errorf("注册失败: %v", err)
	}

	return customer, nil
}

// Login 客户登录
func (s *CustomerAuthService) Login(req *CustomerLoginRequest) (*CustomerLoginResponse, error) {
	var customer model.Customer

	// 根据邮箱查找用户
	err := global.DB.Where("email = ?", req.Email).First(&customer).Error
	if err != nil {
		return nil, errors.New("用户不存在")
	}

	// 检查账号状态
	if customer.Status == model.CustomerStatusBanned {
		return nil, errors.New("账号已被封禁")
	}
	if customer.Status == model.CustomerStatusPending {
		return nil, errors.New("账号待审核")
	}

	// 验证密码
	err = bcrypt.CompareHashAndPassword([]byte(customer.PasswordHash), []byte(req.Password))
	if err != nil {
		return nil, errors.New("密码错误")
	}

	// 生成JWT token
	token, expiresIn, err := utils.GenerateCustomerToken(customer.ID)
	if err != nil {
		return nil, fmt.Errorf("生成token失败: %v", err)
	}

	// 更新登录信息
	customerService := &CustomerService{}
	customerService.UpdateLoginInfo(customer.ID)

	return &CustomerLoginResponse{
		Customer:  &customer,
		Token:     token,
		ExpiresIn: expiresIn,
	}, nil
}

// GetCustomerByID 根据ID获取客户信息
func (s *CustomerAuthService) GetCustomerByID(id uint64) (*model.Customer, error) {
	var customer model.Customer
	if err := global.DB.Where("id = ?", id).First(&customer).Error; err != nil {
		return nil, errors.New("客户不存在")
	}
	return &customer, nil
}

// UpdateCustomerProfile 更新客户个人信息
func (s *CustomerAuthService) UpdateCustomerProfile(id uint64, req *CustomerProfileUpdateRequest) (*model.Customer, error) {
	customerService := &CustomerService{}
	updateReq := &CustomerUpdateRequest{
		Username:  req.Username,
		Email:     req.Email,
		Phone:     req.Phone,
		RealName:  req.RealName,
		AvatarURL: req.AvatarURL,
	}
	return customerService.UpdateCustomer(id, updateReq)
}

// UpdatePassword 更新密码
func (s *CustomerAuthService) UpdatePassword(id uint64, req *CustomerPasswordUpdateRequest) error {
	customer, err := s.GetCustomerByID(id)
	if err != nil {
		return err
	}

	// 验证旧密码
	err = bcrypt.CompareHashAndPassword([]byte(customer.PasswordHash), []byte(req.OldPassword))
	if err != nil {
		return errors.New("旧密码错误")
	}

	// 加密新密码
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("密码加密失败: %v", err)
	}

	// 更新密码
	return global.DB.Model(&model.Customer{}).Where("id = ?", id).Update("password_hash", string(hashedPassword)).Error
}

// ResetPassword 重置密码（通过邮箱）
func (s *CustomerAuthService) ResetPassword(req *CustomerPasswordResetRequest) error {
	email := req.Email
	var customer model.Customer
	if err := global.DB.Where("email = ?", email).First(&customer).Error; err != nil {
		return errors.New("邮箱不存在")
	}

	// TODO: 实现发送重置密码邮件的逻辑
	// 这里简化处理，直接生成一个临时密码
	tempPassword := utils.GenerateRandomString(8)
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(tempPassword), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("密码加密失败: %v", err)
	}

	return global.DB.Model(&model.Customer{}).Where("id = ?", customer.ID).Update("password_hash", string(hashedPassword)).Error
}

// Logout 客户登出
func (s *CustomerAuthService) Logout(customerID uint64) error {
	// 将token加入黑名单（这里简化处理，实际应该使用Redis存储黑名单）
	return nil
}

// RefreshToken 刷新token
func (s *CustomerAuthService) RefreshToken(req *CustomerTokenRefreshRequest) (string, int64, error) {
	// 验证刷新令牌
	claims, err := utils.ParseCustomerToken(req.RefreshToken)
	if err != nil {
		return "", 0, errors.New("无效的刷新令牌")
	}

	// 生成新的token
	token, expiresIn, err := utils.GenerateCustomerToken(uint64(claims.CustomerID))
	if err != nil {
		return "", 0, fmt.Errorf("生成token失败: %v", err)
	}

	return token, expiresIn, nil
}
