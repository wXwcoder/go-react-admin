package api

import (
	"net/http"
	"strconv"

	"go-react-admin/global"
	"go-react-admin/model"

	"github.com/gin-gonic/gin"
)

// GetRoleList 获取角色列表
// @Summary 获取角色列表
// @Description 获取所有角色的列表
// @Tags 角色管理
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Success 200 {object} map[string]interface{} "{"roles":[]model.Role}"
// @Failure 500 {object} map[string]interface{} "{"error":"获取角色列表失败"}"
// @Router /api/roles [get]
func GetRoleList(c *gin.Context) {
	var roles []model.Role
	// 从数据库中获取所有角色
	if err := global.DB.Find(&roles).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "获取角色列表失败",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "获取角色列表成功",
		"roles":   roles,
	})
}

// CreateRole 创建角色
// @Summary 创建角色
// @Description 创建新角色
// @Tags 角色管理
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param role body model.Role true "角色创建信息"
// @Success 200 {object} map[string]interface{} "{"message":"角色创建成功","role":model.Role}"
// @Failure 400 {object} map[string]interface{} "{"error":"请求参数错误"}"
// @Failure 500 {object} map[string]interface{} "{"error":"创建角色失败"}"
// @Router /api/roles [post]
func CreateRole(c *gin.Context) {
	var role model.Role
	// 绑定JSON到role
	if err := c.ShouldBindJSON(&role); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "请求参数错误",
		})
		return
	}

	// 创建角色
	if err := global.DB.Create(&role).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "创建角色失败",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "角色创建成功",
		"role":    role,
	})
}

// UpdateRole 更新角色
// @Summary 更新角色
// @Description 根据角色ID更新角色信息
// @Tags 角色管理
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param id path int true "角色ID"
// @Param role body model.Role true "角色更新信息"
// @Success 200 {object} map[string]interface{} "{"message":"角色更新成功"}"
// @Failure 400 {object} map[string]interface{} "{"error":"请求参数错误"}"
// @Failure 500 {object} map[string]interface{} "{"error":"更新角色失败"}"
// @Router /api/roles/{id} [put]
func UpdateRole(c *gin.Context) {
	id := c.Param("id")
	var role model.Role
	// 绑定JSON到role
	if err := c.ShouldBindJSON(&role); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "请求参数错误",
		})
		return
	}

	// 更新角色
	if err := global.DB.Model(&model.Role{}).Where("id = ?", id).Updates(role).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "更新角色失败",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "角色更新成功",
	})
}

// DeleteRole 删除角色
// @Summary 删除角色
// @Description 根据角色ID删除角色
// @Tags 角色管理
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param id path int true "角色ID"
// @Success 200 {object} map[string]interface{} "{"message":"角色删除成功"}"
// @Failure 400 {object} map[string]interface{} "{"error":"无效的角色ID"}"
// @Failure 500 {object} map[string]interface{} "{"error":"删除角色失败"}"
// @Router /api/roles/{id} [delete]
func DeleteRole(c *gin.Context) {
	id := c.Param("id")
	roleID, err := strconv.Atoi(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "无效的角色ID",
		})
		return
	}

	// 删除角色
	if err := global.DB.Delete(&model.Role{}, roleID).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "删除角色失败",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "角色删除成功",
	})
}
