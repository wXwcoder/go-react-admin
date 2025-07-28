package api

import (
	"net/http"

	"go-react-admin/global"
	"go-react-admin/model"

	"github.com/gin-gonic/gin"
)

// GetTenantList 获取租户列表
func GetTenantList(c *gin.Context) {
	var tenants []model.Tenant
	// 从数据库中获取所有租户
	if err := global.DB.Find(&tenants).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "获取租户列表失败",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"tenants": tenants,
	})
}
