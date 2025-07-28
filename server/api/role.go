package api

import (
	"net/http"
	"strconv"

	"go-react-admin/global"
	"go-react-admin/model"

	"github.com/gin-gonic/gin"
)

// GetRoleList 获取角色列表
func GetRoleList(c *gin.Context) {
	var roles []model.Role
	// 从数据库中获取所有角色
	if err := global.DB.Find(&roles).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "获取角色列表失败",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"roles": roles,
	})
}

// CreateRole 创建角色
func CreateRole(c *gin.Context) {
	var role model.Role
	// 绑定JSON到role
	if err := c.ShouldBindJSON(&role); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "请求参数错误",
		})
		return
	}

	// 创建角色
	if err := global.DB.Create(&role).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "创建角色失败",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "角色创建成功",
		"role":    role,
	})
}

// UpdateRole 更新角色
func UpdateRole(c *gin.Context) {
	id := c.Param("id")
	var role model.Role
	// 绑定JSON到role
	if err := c.ShouldBindJSON(&role); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "请求参数错误",
		})
		return
	}

	// 更新角色
	if err := global.DB.Model(&model.Role{}).Where("id = ?", id).Updates(role).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "更新角色失败",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "角色更新成功",
	})
}

// DeleteRole 删除角色
func DeleteRole(c *gin.Context) {
	id := c.Param("id")
	roleID, err := strconv.Atoi(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "无效的角色ID",
		})
		return
	}

	// 删除角色
	if err := global.DB.Delete(&model.Role{}, roleID).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "删除角色失败",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "角色删除成功",
	})
}
