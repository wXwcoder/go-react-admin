package api

import (
	"net/http"

	"go-react-admin/global"
	"go-react-admin/model"

	"github.com/gin-gonic/gin"
)

// GetLogList 获取日志列表
func GetLogList(c *gin.Context) {
	var logs []model.Log
	// 从数据库中获取所有日志
	if err := global.DB.Find(&logs).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "获取日志列表失败",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"logs": logs,
	})
}
