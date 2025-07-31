package middleware

import (
	"fmt"
	"net/http"
	"time"

	"go-react-admin/global"
	"go-react-admin/model"

	"github.com/gin-gonic/gin"
)

// Logger 日志记录中间件
func Logger() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 开始时间
		start := time.Now()
		
		// 处理请求
		c.Next()
		
		// 结束时间
		end := time.Now()
		
		// 计算响应时间（毫秒）
		responseTime := end.Sub(start).Milliseconds()
		
		// 获取用户信息
		var username string
		var userID uint
		var tenantID uint
		
		// 从JWT token中获取用户信息
		if val, exists := c.Get("username"); exists {
			username = val.(string)
		}
		if val, exists := c.Get("user_id"); exists {
			userID = val.(uint)
		}
		if val, exists := c.Get("tenant_id"); exists {
			tenantID = val.(uint)
		}
		
		// 获取客户端IP
		clientIP := c.ClientIP()
		
		// 获取User-Agent
		userAgent := c.Request.UserAgent()
		
		// 获取请求方法
		method := c.Request.Method
		
		// 获取请求路径
		path := c.Request.URL.Path
		
		// 获取状态码
		statusCode := c.Writer.Status()
		
		// 创建日志记录
		log := model.Log{
			UserID:       userID,
			Username:     username,
			IP:           clientIP,
			Method:       method,
			Path:         path,
			UserAgent:    userAgent,
			StatusCode:   statusCode,
			ResponseTime: int(responseTime),
			TenantID:     tenantID,
		}
		
		// 保存到数据库
		if err := global.DB.Create(&log).Error; err != nil {
			fmt.Printf("记录日志失败: %v\n", err)
		}
	}
}

// LoginLogger 登录日志记录中间件
func LoginLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 开始时间
		start := time.Now()
		
		// 获取客户端IP
		clientIP := c.ClientIP()
		
		// 获取User-Agent
		userAgent := c.Request.UserAgent()
		
		// 获取请求方法
		method := c.Request.Method
		
		// 获取请求路径
		path := c.Request.URL.Path
		
		// 处理请求
		c.Next()
		
		// 结束时间
		end := time.Now()
		
		// 计算响应时间（毫秒）
		responseTime := end.Sub(start).Milliseconds()
		
		// 获取状态码
		statusCode := c.Writer.Status()
		
		// 如果是登录请求且成功，记录登录日志
		if path == "/api/v1/user/login" && method == "POST" {
			var username string
			var userID uint
			var tenantID uint
			
			// 如果登录成功，从响应中获取用户信息
			if statusCode == http.StatusOK {
				// 尝试从上下文中获取用户信息（如果登录成功）
				if val, exists := c.Get("username"); exists {
					username = val.(string)
				}
				if val, exists := c.Get("user_id"); exists {
					userID = val.(uint)
				}
				if val, exists := c.Get("tenant_id"); exists {
					tenantID = val.(uint)
				}
				
				// 如果上下文中没有用户信息，尝试从请求体中获取
				if username == "" {
					var loginData struct {
						Username string `json:"username"`
					}
					if err := c.ShouldBindJSON(&loginData); err == nil {
						username = loginData.Username
					}
				}
			}
			
			// 创建登录日志记录
			log := model.Log{
				UserID:       userID,
				Username:     username,
				IP:           clientIP,
				Method:       method,
				Path:         path,
				UserAgent:    userAgent,
				StatusCode:   statusCode,
				ResponseTime: int(responseTime),
				TenantID:     tenantID,
			}
			
			// 保存到数据库
			if err := global.DB.Create(&log).Error; err != nil {
				fmt.Printf("记录登录日志失败: %v\n", err)
			}
		}
	}
}

// LogoutLogger 登出日志记录中间件
func LogoutLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 获取用户信息
		var username string
		var userID uint
		var tenantID uint
		
		// 从JWT token中获取用户信息
		if val, exists := c.Get("username"); exists {
			username = val.(string)
		}
		if val, exists := c.Get("user_id"); exists {
			userID = val.(uint)
		}
		if val, exists := c.Get("tenant_id"); exists {
			tenantID = val.(uint)
		}
		
		// 获取客户端IP
		clientIP := c.ClientIP()
		
		// 获取User-Agent
		userAgent := c.Request.UserAgent()
		
		// 获取请求方法
		method := c.Request.Method
		
		// 获取请求路径
		path := c.Request.URL.Path
		
		// 获取状态码
		statusCode := c.Writer.Status()
		
		// 创建登出日志记录
		log := model.Log{
			UserID:       userID,
			Username:     username,
			IP:           clientIP,
			Method:       method,
			Path:         path,
			UserAgent:    userAgent,
			StatusCode:   statusCode,
			ResponseTime: 0, // 登出操作响应时间通常很短
			TenantID:     tenantID,
		}
		
		// 保存到数据库
		if err := global.DB.Create(&log).Error; err != nil {
			fmt.Printf("记录登出日志失败: %v\n", err)
		}
		
		c.Next()
	}
}