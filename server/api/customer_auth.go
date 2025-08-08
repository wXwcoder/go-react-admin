package api

import (
	"fmt"
	"go-react-admin/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

type CustomerAuthAPI struct{}

// Register 第三方客户注册
// @Summary 第三方客户注册
// @Description 第三方客户注册新账号
// @Tags 第三方客户认证
// @Accept json
// @Produce json
// @Param register body service.CustomerRegisterRequest true "注册信息"
// @Success 200 {object} service.CustomerAuthResponse
// @Router /api/v1/customer/register [post]
func (api *CustomerAuthAPI) Register(c *gin.Context) {
	var req service.CustomerRegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	authService := &service.CustomerAuthService{}
	response, err := authService.Register(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": response})
}

// Login 第三方客户登录
// @Summary 第三方客户登录
// @Description 第三方客户登录系统
// @Tags 第三方客户认证
// @Accept json
// @Produce json
// @Param login body service.CustomerLoginRequest true "登录信息"
// @Success 200 {object} service.CustomerAuthResponse
// @Router /api/v1/customer/login [post]
func (api *CustomerAuthAPI) Login(c *gin.Context) {
	var req service.CustomerLoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	authService := &service.CustomerAuthService{}
	response, err := authService.Login(&req)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": response})
}

// GetProfile 获取第三方客户个人信息
// @Summary 获取第三方客户个人信息
// @Description 获取当前登录的第三方客户的个人信息
// @Tags 第三方客户认证
// @Accept json
// @Produce json
// @Success 200 {object} model.Customer
// @Router /api/v1/customer/profile [get]
// @Security CustomerJWTAuth
func (api *CustomerAuthAPI) GetProfile(c *gin.Context) {
	customerID, exists := c.Get("customer_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "未授权访问"})
		return
	}
	fmt.Printf("customerID: %+v\n", customerID)
	authService := &service.CustomerAuthService{}
	customer, err := authService.GetCustomerByID(uint64(customerID.(uint)))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "客户不存在"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": customer})
}

// UpdateProfile 更新第三方客户个人信息
// @Summary 更新第三方客户个人信息
// @Description 更新当前登录的第三方客户的个人信息
// @Tags 第三方客户认证
// @Accept json
// @Produce json
// @Param profile body service.CustomerProfileUpdateRequest true "个人信息"
// @Success 200 {object} model.Customer
// @Router /api/v1/customer/profile [put]
// @Security CustomerJWTAuth
func (api *CustomerAuthAPI) UpdateProfile(c *gin.Context) {
	customerID, exists := c.Get("customer_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "未授权访问"})
		return
	}

	var req service.CustomerProfileUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	authService := &service.CustomerAuthService{}
	customer, err := authService.UpdateCustomerProfile(uint64(customerID.(uint)), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": customer})
}

// UpdatePassword 更新第三方客户密码
// @Summary 更新第三方客户密码
// @Description 更新当前登录的第三方客户的密码
// @Tags 第三方客户认证
// @Accept json
// @Produce json
// @Param password body service.CustomerPasswordUpdateRequest true "密码信息"
// @Success 200 {object} gin.H{"message":"密码更新成功"}
// @Router /api/v1/customer/password [put]
// @Security CustomerJWTAuth
func (api *CustomerAuthAPI) UpdatePassword(c *gin.Context) {
	customerID, exists := c.Get("customer_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "未授权访问"})
		return
	}

	var req service.CustomerPasswordUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	authService := &service.CustomerAuthService{}
	if err := authService.UpdatePassword(uint64(customerID.(uint)), &req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "密码更新成功"})
}

// ResetPassword 重置第三方客户密码
// @Summary 重置第三方客户密码
// @Description 通过邮箱重置第三方客户的密码
// @Tags 第三方客户认证
// @Accept json
// @Produce json
// @Param reset body service.CustomerPasswordResetRequest true "重置信息"
// @Success 200 {object} gin.H{"message":"重置邮件已发送"}
// @Router /api/v1/customer/reset-password [post]
func (api *CustomerAuthAPI) ResetPassword(c *gin.Context) {
	var req service.CustomerPasswordResetRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	authService := &service.CustomerAuthService{}
	if err := authService.ResetPassword(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "重置邮件已发送"})
}

// Logout 第三方客户登出
// @Summary 第三方客户登出
// @Description 第三方客户登出系统
// @Tags 第三方客户认证
// @Accept json
// @Produce json
// @Success 200 {object} gin.H{"message":"登出成功"}
// @Router /api/v1/customer/logout [post]
// @Security CustomerJWTAuth
func (api *CustomerAuthAPI) Logout(c *gin.Context) {
	customerID, exists := c.Get("customer_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "未授权访问"})
		return
	}

	authService := &service.CustomerAuthService{}
	if err := authService.Logout(uint64(customerID.(uint))); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "登出成功"})
}

// RefreshToken 刷新第三方客户令牌
// @Summary 刷新第三方客户令牌
// @Description 使用刷新令牌获取新的访问令牌
// @Tags 第三方客户认证
// @Accept json
// @Produce json
// @Param refresh body service.CustomerTokenRefreshRequest true "刷新令牌"
// @Success 200 {object} service.CustomerAuthResponse
// @Router /api/v1/customer/refresh-token [post]
func (api *CustomerAuthAPI) RefreshToken(c *gin.Context) {
	var req service.CustomerTokenRefreshRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	authService := &service.CustomerAuthService{}
	response, expiresIn, err := authService.RefreshToken(&req)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": gin.H{"token": response, "expires_in": expiresIn}})
}
