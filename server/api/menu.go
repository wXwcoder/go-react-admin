package api

import (
	"net/http"
	"strconv"

	"go-react-admin/global"
	"go-react-admin/model"

	"github.com/gin-gonic/gin"
)

// GetMenuList 获取菜单列表
func GetMenuList(c *gin.Context) {
	var menus []model.Menu
	// 从数据库中获取所有菜单
	if err := global.DB.Find(&menus).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "获取菜单列表失败",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"menus": menus,
	})
}

// GetUserMenus 获取用户菜单
func GetUserMenus(c *gin.Context) {
	var menus []model.Menu
	// 暂时返回所有菜单，后续可以根据用户权限过滤
	if err := global.DB.Where("status = ?", 1).Order("sort ASC").Find(&menus).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "获取用户菜单失败",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"data": menus,
		"message": "获取用户菜单成功",
	})
}

// CreateMenu 创建菜单
func CreateMenu(c *gin.Context) {
	var menu model.Menu
	// 绑定JSON到menu
	if err := c.ShouldBindJSON(&menu); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "请求参数错误",
		})
		return
	}

	// 创建菜单
	if err := global.DB.Create(&menu).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "创建菜单失败",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "菜单创建成功",
		"menu":    menu,
	})
}

// UpdateMenu 更新菜单
func UpdateMenu(c *gin.Context) {
	id := c.Param("id")
	var menu model.Menu
	// 绑定JSON到menu
	if err := c.ShouldBindJSON(&menu); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "请求参数错误",
		})
		return
	}

	// 更新菜单
	if err := global.DB.Model(&model.Menu{}).Where("id = ?", id).Updates(menu).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "更新菜单失败",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "菜单更新成功",
	})
}

// DeleteMenu 删除菜单
func DeleteMenu(c *gin.Context) {
	id := c.Param("id")
	menuID, err := strconv.Atoi(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "无效的菜单ID",
		})
		return
	}

	// 删除菜单
	if err := global.DB.Delete(&model.Menu{}, menuID).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "删除菜单失败",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "菜单删除成功",
	})
}
