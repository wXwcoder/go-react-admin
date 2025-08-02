package api

import (
	"net/http"
	"strconv"

	"go-react-admin/global"
	"go-react-admin/model"

	"github.com/gin-gonic/gin"
)

// GetMenuList 获取菜单列表
// @Summary 获取菜单列表
// @Description 获取所有菜单的列表
// @Tags 菜单管理
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Success 200 {object} map[string]interface{} "{"menus":[]model.Menu}"
// @Failure 500 {object} map[string]interface{} "{"error":"获取菜单列表失败"}"
// @Router /api/menus [get]
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
// @Summary 获取用户菜单
// @Description 获取当前用户的菜单列表
// @Tags 菜单管理
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Success 200 {object} map[string]interface{} "{"code":200,"data":[]model.Menu,"message":"获取用户菜单成功"}"
// @Failure 500 {object} map[string]interface{} "{"error":"获取用户菜单失败"}"
// @Router /api/menus/user [get]
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
// @Summary 创建菜单
// @Description 创建新菜单
// @Tags 菜单管理
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param menu body model.Menu true "菜单创建信息"
// @Success 200 {object} map[string]interface{} "{"message":"菜单创建成功","menu":model.Menu}"
// @Failure 400 {object} map[string]interface{} "{"error":"请求参数错误"}"
// @Failure 500 {object} map[string]interface{} "{"error":"创建菜单失败"}"
// @Router /api/menus [post]
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
// @Summary 更新菜单
// @Description 根据菜单ID更新菜单信息
// @Tags 菜单管理
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param id path int true "菜单ID"
// @Param menu body model.Menu true "菜单更新信息"
// @Success 200 {object} map[string]interface{} "{"message":"菜单更新成功"}"
// @Failure 400 {object} map[string]interface{} "{"error":"请求参数错误"}"
// @Failure 500 {object} map[string]interface{} "{"error":"更新菜单失败"}"
// @Router /api/menus/{id} [put]
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
// @Summary 删除菜单
// @Description 根据菜单ID删除菜单
// @Tags 菜单管理
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param id path int true "菜单ID"
// @Success 200 {object} map[string]interface{} "{"message":"菜单删除成功"}"
// @Failure 400 {object} map[string]interface{} "{"error":"无效的菜单ID"}"
// @Failure 500 {object} map[string]interface{} "{"error":"删除菜单失败"}"
// @Router /api/menus/{id} [delete]
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
