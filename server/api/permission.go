package api

import (
	"fmt"
	"net/http"
	"strconv"

	"go-react-admin/service"

	"github.com/gin-gonic/gin"
)

type PermissionApi struct {
	permissionService *service.PermissionService
}

func NewPermissionApi() *PermissionApi {
	return &PermissionApi{
		permissionService: &service.PermissionService{},
	}
}

// AssignRolePermissions 分配角色权限
// @Summary 分配角色权限
// @Description 为角色分配菜单和API权限
// @Tags 权限管理
// @Accept json
// @Produce json
// @Param request body service.PermissionRequest true "权限分配请求"
// @Success 200 {object} map[string]interface{} "成功"
// @Failure 400 {object} map[string]interface{} "请求参数错误"
// @Failure 500 {object} map[string]interface{} "服务器错误"
// @Router /api/permissions/role [post]
func (p *PermissionApi) AssignRolePermissions(c *gin.Context) {
	var req service.PermissionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"msg":  "请求参数错误: " + err.Error(),
		})
		return
	}

	// 从上下文获取租户ID
	if tenantID, exists := c.Get("tenant_id"); exists {
		req.TenantID = tenantID.(uint)
	}

	if err := p.permissionService.AssignRolePermissions(&req); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": 500,
			"msg":  "分配权限失败: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "权限分配成功",
	})
}

// AssignUserRoles 分配用户角色
// @Summary 分配用户角色
// @Description 为用户分配角色
// @Tags 权限管理
// @Accept json
// @Produce json
// @Param request body service.UserRoleRequest true "用户角色分配请求"
// @Success 200 {object} map[string]interface{} "成功"
// @Failure 400 {object} map[string]interface{} "请求参数错误"
// @Failure 500 {object} map[string]interface{} "服务器错误"
// @Router /api/permissions/user [post]
func (p *PermissionApi) AssignUserRoles(c *gin.Context) {
	var req service.UserRoleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"msg":  "请求参数错误: " + err.Error(),
		})
		return
	}

	// 从上下文获取租户ID
	if tenantID, exists := c.Get("tenant_id"); exists {
		req.TenantID = tenantID.(uint)
	}

	if err := p.permissionService.AssignUserRoles(&req); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": 500,
			"msg":  "分配角色失败: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "角色分配成功",
	})
}

// GetRolePermissions 获取角色权限
// @Summary 获取角色权限
// @Description 获取指定角色的权限信息
// @Tags 权限管理
// @Accept json
// @Produce json
// @Param id path int true "角色ID"
// @Success 200 {object} map[string]interface{} "成功"
// @Failure 400 {object} map[string]interface{} "请求参数错误"
// @Failure 500 {object} map[string]interface{} "服务器错误"
// @Router /api/permissions/role/{id} [get]
func (p *PermissionApi) GetRolePermissions(c *gin.Context) {
	roleIDStr := c.Param("id")
	roleID, err := strconv.ParseUint(roleIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"msg":  "无效的角色ID",
		})
		return
	}

	// 从上下文获取租户ID
	tenantID, exists := c.Get("tenant_id")
	if !exists {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"msg":  "租户信息缺失",
		})
		return
	}

	permissions, err := p.permissionService.GetRolePermissions(uint(roleID), tenantID.(uint))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": 500,
			"msg":  "获取权限失败: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "获取成功",
		"data": permissions,
	})
}

// GetUserRoles 获取用户角色
// @Summary 获取用户角色
// @Description 获取指定用户的角色信息
// @Tags 权限管理
// @Accept json
// @Produce json
// @Param id path int true "用户ID"
// @Success 200 {object} map[string]interface{} "成功"
// @Failure 400 {object} map[string]interface{} "请求参数错误"
// @Failure 500 {object} map[string]interface{} "服务器错误"
// @Router /api/permissions/user/{id}/roles [get]
func (p *PermissionApi) GetUserRoles(c *gin.Context) {
	userIDStr := c.Param("id")
	userID, err := strconv.ParseUint(userIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"msg":  "无效的用户ID",
		})
		return
	}

	// 从上下文获取租户ID
	tenantID, exists := c.Get("tenant_id")
	if !exists {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"msg":  "租户信息缺失",
		})
		return
	}

	userRoles, err := p.permissionService.GetUserRoles(uint(userID), tenantID.(uint))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": 500,
			"msg":  "获取用户角色失败: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "获取成功",
		"data": userRoles,
	})
}

// GetUserPermissions 获取用户权限
// @Summary 获取用户权限
// @Description 获取指定用户的所有权限（通过角色）
// @Tags 权限管理
// @Accept json
// @Produce json
// @Param id path int true "用户ID"
// @Success 200 {object} map[string]interface{} "成功"
// @Failure 400 {object} map[string]interface{} "请求参数错误"
// @Failure 500 {object} map[string]interface{} "服务器错误"
// @Router /api/permissions/user/{id}/permissions [get]
func (p *PermissionApi) GetUserPermissions(c *gin.Context) {
	userIDStr := c.Param("id")
	userID, err := strconv.ParseUint(userIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"msg":  "无效的用户ID",
		})
		return
	}

	// 从上下文获取租户ID
	tenantID, exists := c.Get("tenant_id")
	if !exists {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"msg":  "租户信息缺失",
		})
		return
	}

	menus, apis, err := p.permissionService.GetUserPermissions(uint(userID), tenantID.(uint))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": 500,
			"msg":  "获取用户权限失败: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "获取成功",
		"data": gin.H{
			"menus": menus,
			"apis":  apis,
		},
	})
}

// CheckPermission 检查权限
// @Summary 检查权限
// @Description 检查用户是否有特定权限
// @Tags 权限管理
// @Accept json
// @Produce json
// @Param resource query string true "资源"
// @Param action query string true "操作"
// @Success 200 {object} map[string]interface{} "成功"
// @Failure 400 {object} map[string]interface{} "请求参数错误"
// @Failure 500 {object} map[string]interface{} "服务器错误"
// @Router /api/permissions/check [get]
func (p *PermissionApi) CheckPermission(c *gin.Context) {
	resource := c.Query("resource")
	action := c.Query("action")

	fmt.Printf("检查权限: 资源=%s, 操作=%s\n", resource, action)
	if resource == "" || action == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"msg":  "资源和操作参数不能为空",
		})
		fmt.Printf("检查权限: %s\n", "资源和操作参数不能为空")
		return
	}

	// 从上下文获取用户ID和租户ID
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"code": 401,
			"msg":  "用户未认证",
		})
		fmt.Printf("检查权限: %s\n", "用户未认证")
		return
	}

	tenantID, exists := c.Get("tenant_id")
	if !exists {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"msg":  "租户信息缺失",
		})
		fmt.Printf("检查权限: %s\n", "租户信息缺失")
		return
	}

	allowed, err := p.permissionService.CheckUserPermission(userID.(uint), resource, action, tenantID.(uint))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": 500,
			"msg":  "权限检查失败: " + err.Error(),
		})
		fmt.Printf("检查权限: %s\n", "权限检查失败")
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":          200,
		"msg":           "检查完成",
		"hasPermission": allowed,
	})
	fmt.Printf("检查权限: 资源=%s, 操作=%s, 结果=%v\n", resource, action, allowed)
}
