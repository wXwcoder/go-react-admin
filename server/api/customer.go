package api

import (
	"fmt"
	"go-react-admin/service"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type CustomerAPI struct{}

// CreateCustomer 创建第三方客户
// @Summary 创建第三方客户
// @Description 创建新的第三方客户账号
// @Tags 第三方客户管理
// @Accept json
// @Produce json
// @Param customer body service.CustomerCreateRequest true "客户信息"
// @Success 200 {object} model.Customer
// @Router /api/v1/admin/customers [post]
func (api *CustomerAPI) CreateCustomer(c *gin.Context) {
	var req service.CustomerCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	customerService := &service.CustomerService{}
	customer, err := customerService.CreateCustomer(&req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": customer})
}

// UpdateCustomer 更新第三方客户
// @Summary 更新第三方客户信息
// @Description 更新指定第三方客户的信息
// @Tags 第三方客户管理
// @Accept json
// @Produce json
// @Param id path int true "客户ID"
// @Param customer body service.CustomerUpdateRequest true "更新信息"
// @Success 200 {object} model.Customer
// @Router /api/v1/admin/customers/{id} [put]
func (api *CustomerAPI) UpdateCustomer(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的客户ID"})
		return
	}

	var req service.CustomerUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	customerService := &service.CustomerService{}
	customer, err := customerService.UpdateCustomer(id, &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": customer})
}

// DeleteCustomer 删除第三方客户
// @Summary 删除第三方客户
// @Description 删除指定的第三方客户账号（软删除）
// @Tags 第三方客户管理
// @Accept json
// @Produce json
// @Param id path int true "客户ID"
// @Success 200 {object} gin.H{"message":"删除成功"}
// @Router /api/v1/admin/customers/{id} [delete]
func (api *CustomerAPI) DeleteCustomer(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的客户ID"})
		return
	}

	customerService := &service.CustomerService{}
	if err := customerService.DeleteCustomer(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "删除成功"})
}

// GetCustomer 获取第三方客户详情
// @Summary 获取第三方客户详情
// @Description 获取指定第三方客户的详细信息
// @Tags 第三方客户管理
// @Accept json
// @Produce json
// @Param id path int true "客户ID"
// @Success 200 {object} model.Customer
// @Router /api/v1/admin/customers/{id} [get]
func (api *CustomerAPI) GetCustomer(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的客户ID"})
		return
	}

	customerService := &service.CustomerService{}
	customer, err := customerService.GetCustomerByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": customer})
}

// GetCustomerList 获取第三方客户列表
// @Summary 获取第三方客户列表
// @Description 获取第三方客户列表，支持分页、搜索和筛选
// @Tags 第三方客户管理
// @Accept json
// @Produce json
// @Param page query int false "页码" default(1)
// @Param page_size query int false "每页数量" default(10) minimum(1) maximum(100)
// @Param keyword query string false "搜索关键词"
// @Param status query string false "状态筛选" Enums(active,banned,pending)
// @Success 200 {object} service.CustomerListResponse
// @Router /api/v1/admin/customers [get]
func (api *CustomerAPI) GetCustomerList(c *gin.Context) {
	var req service.CustomerListRequest
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

	customerService := &service.CustomerService{}
	response, err := customerService.GetCustomerList(&req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	fmt.Printf("GetCustomerList:%+v \n", response)

	c.JSON(http.StatusOK, gin.H{"data": response})
}

// UpdateCustomerStatus 更新客户状态
// @Summary 更新客户状态
// @Description 更新第三方客户的状态（active/banned/pending）
// @Tags 第三方客户管理
// @Accept json
// @Produce json
// @Param id path int true "客户ID"
// @Param status body service.CustomerStatusUpdateRequest true "状态信息"
// @Success 200 {object} model.Customer
// @Router /api/v1/admin/customers/{id}/status [put]
func (api *CustomerAPI) UpdateCustomerStatus(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的客户ID"})
		return
	}

	var req service.CustomerStatusUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	customerService := &service.CustomerService{}
	customer, err := customerService.UpdateCustomerStatus(id, req.Status)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": customer})
}

// BanCustomer 封禁客户
// @Summary 封禁客户
// @Description 封禁指定的第三方客户
// @Tags 第三方客户管理
// @Accept json
// @Produce json
// @Param id path int true "客户ID"
// @Success 200 {object} model.Customer
// @Router /api/v1/admin/customers/{id}/ban [post]
func (api *CustomerAPI) BanCustomer(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的客户ID"})
		return
	}

	customerService := &service.CustomerService{}
	customer, err := customerService.BanCustomer(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": customer})
}

// UnbanCustomer 解封客户
// @Summary 解封客户
// @Description 解封指定的第三方客户
// @Tags 第三方客户管理
// @Accept json
// @Produce json
// @Param id path int true "客户ID"
// @Success 200 {object} model.Customer
// @Router /api/v1/admin/customers/{id}/unban [post]
func (api *CustomerAPI) UnbanCustomer(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的客户ID"})
		return
	}

	customerService := &service.CustomerService{}
	customer, err := customerService.UnbanCustomer(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": customer})
}
