package api

import (
	"net/http"
	"strconv"

	"go-react-admin/global"
	"go-react-admin/model"

	"github.com/gin-gonic/gin"
)

// GetApiList 获取API列表
func GetApiList(c *gin.Context) {
	var apis []model.Api
	// 从数据库中获取所有API
	if err := global.DB.Find(&apis).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "获取API列表失败",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"apis": apis,
	})
}

// CreateApi 创建API
func CreateApi(c *gin.Context) {
	var api model.Api
	if err := c.ShouldBindJSON(&api); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "请求参数错误",
		})
		return
	}

	// 保存到数据库
	if err := global.DB.Create(&api).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "创建API失败",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "创建API成功",
		"api":     api,
	})
}

// UpdateApi 更新API
func UpdateApi(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "无效的API ID",
		})
		return
	}

	var api model.Api
	if err := c.ShouldBindJSON(&api); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "请求参数错误",
		})
		return
	}

	// 更新数据库中的API
	if err := global.DB.Model(&model.Api{}).Where("id = ?", id).Updates(api).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "更新API失败",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "更新API成功",
	})
}

// DeleteApi 删除API
func DeleteApi(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "无效的API ID",
		})
		return
	}

	// 从数据库中删除API
	if err := global.DB.Delete(&model.Api{}, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "删除API失败",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "删除API成功",
	})
}
