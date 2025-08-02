package api

import (
	"net/http"
	"strconv"
	"time"

	"go-react-admin/global"
	"go-react-admin/model"

	"github.com/gin-gonic/gin"
)

// GetLogList 获取日志列表
// @Summary 获取日志列表
// @Description 获取系统日志列表，支持分页、搜索和筛选
// @Tags 日志管理
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param page query int false "页码" default(1)
// @Param pageSize query int false "每页数量" default(10) minimum(1) maximum(100)
// @Param username query string false "用户名搜索"
// @Param method query string false "HTTP方法"
// @Param statusCode query int false "状态码"
// @Param startDate query string false "开始日期" format(date)
// @Param endDate query string false "结束日期" format(date)
// @Success 200 {object} map[string]interface{} "{"logs":[]model.Log,"total":int,"page":int,"pageSize":int}"
// @Failure 500 {object} map[string]interface{} "{"error":"获取日志列表失败"}"
// @Router /api/logs [get]
func GetLogList(c *gin.Context) {
	var logs []model.Log
	var total int64

	// 获取分页参数
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "10"))
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}

	// 获取搜索参数
	username := c.Query("username")
	method := c.Query("method")
	statusCode := c.Query("statusCode")
	startDate := c.Query("startDate")
	endDate := c.Query("endDate")

	// 构建查询条件
	db := global.DB.Model(&model.Log{})

	if username != "" {
		db = db.Where("username LIKE ?", "%"+username+"%")
	}
	if method != "" {
		db = db.Where("method = ?", method)
	}
	if statusCode != "" {
		if code, err := strconv.Atoi(statusCode); err == nil {
			db = db.Where("status_code = ?", code)
		}
	}
	if startDate != "" && endDate != "" {
		start, _ := time.Parse("2006-01-02", startDate)
		end, _ := time.Parse("2006-01-02", endDate)
		end = end.Add(24 * time.Hour).Add(-time.Second) // 设置到当天结束
		db = db.Where("created_at BETWEEN ? AND ?", start, end)
	}

	// 获取总数
	db.Count(&total)

	// 获取分页数据
	offset := (page - 1) * pageSize
	if err := db.Order("created_at DESC").Limit(pageSize).Offset(offset).Find(&logs).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "获取日志列表失败",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"logs": logs,
		"total": total,
		"page": page,
		"pageSize": pageSize,
	})
}
