package api

import (
	"net/http"
	"strconv"

	"go-react-admin/global"
	"go-react-admin/model"

	"github.com/gin-gonic/gin"
)

// GetApiList 获取API列表
// @Summary 获取API列表
// @Description 获取所有API的列表
// @Tags API管理
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Success 200 {object} map[string]interface{} "{"apis":[]model.Api}"
// @Failure 500 {object} map[string]interface{} "{"error":"获取API列表失败"}"
// @Router /api/apis [get]
func GetApiList(c *gin.Context) {
	var apis []model.Api
	// 从数据库中获取所有API
	if err := global.DB.Find(&apis).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "获取API列表失败",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "获取API列表成功",
		"apis":    apis,
	})
}

// CreateApi 创建API
// @Summary 创建API
// @Description 创建新API接口
// @Tags API管理
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param api body model.Api true "API创建信息"
// @Success 200 {object} map[string]interface{} "{"message":"创建API成功","api":model.Api}"
// @Failure 400 {object} map[string]interface{} "{"error":"请求参数错误"}"
// @Failure 500 {object} map[string]interface{} "{"error":"创建API失败"}"
// @Router /api/apis [post]
func CreateApi(c *gin.Context) {
	var api model.Api
	if err := c.ShouldBindJSON(&api); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "请求参数错误",
		})
		return
	}

	// 保存到数据库
	if err := global.DB.Create(&api).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "创建API失败",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "创建API成功",
		"api":     api,
	})
}

// UpdateApi 更新API
// @Summary 更新API
// @Description 根据API ID更新API信息
// @Tags API管理
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param id path int true "API ID"
// @Param api body model.Api true "API更新信息"
// @Success 200 {object} map[string]interface{} "{"message":"更新API成功"}"
// @Failure 400 {object} map[string]interface{} "{"error":"请求参数错误"}"
// @Failure 500 {object} map[string]interface{} "{"error":"更新API失败"}"
// @Router /api/apis/{id} [put]
func UpdateApi(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "无效的API ID",
		})
		return
	}

	var api model.Api
	if err := c.ShouldBindJSON(&api); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "请求参数错误",
		})
		return
	}

	// 更新数据库中的API
	if err := global.DB.Model(&model.Api{}).Where("id = ?", id).Updates(api).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "更新API失败",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "更新API成功",
	})
}

// DeleteApi 删除API
// @Summary 删除API
// @Description 根据API ID删除API
// @Tags API管理
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param id path int true "API ID"
// @Success 200 {object} map[string]interface{} "{"message":"删除API成功"}"
// @Failure 400 {object} map[string]interface{} "{"error":"无效的API ID"}"
// @Failure 500 {object} map[string]interface{} "{"error":"删除API失败"}"
// @Router /api/apis/{id} [delete]
func DeleteApi(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "无效的API ID",
		})
		return
	}

	// 从数据库中删除API
	if err := global.DB.Delete(&model.Api{}, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "删除API失败",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "删除API成功",
	})
}
