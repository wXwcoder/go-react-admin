package middleware

import (
	"net/http"
	"strings"

	"go-react-admin/utils"

	"github.com/gin-gonic/gin"
)

// JWTAuth JWT认证中间件
func JWTAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 获取Authorization header
		authHeader := c.Request.Header.Get("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "请求头中缺少Authorization字段",
			})
			c.Abort()
			return
		}

		// 检查Bearer前缀
		parts := strings.SplitN(authHeader, " ", 2)
		if !(len(parts) == 2 && parts[0] == "Bearer") {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "请求头中Authorization格式错误",
			})
			c.Abort()
			return
		}

		// 解析JWT token
		claims, err := utils.ParseToken(parts[1])
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "无效的Token",
			})
			c.Abort()
			return
		}

		// 将用户信息存储到上下文中
		c.Set("username", claims.Username)
		c.Set("user_id", claims.UserID)
		c.Set("tenant_id", claims.TenantID)

		c.Next()
	}
}

// CustomerJWTAuth 第三方客户JWT认证中间件
func CustomerJWTAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 获取Authorization header
		authHeader := c.Request.Header.Get("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "请求头中缺少Authorization字段",
			})
			c.Abort()
			return
		}

		// 检查Bearer前缀
		parts := strings.SplitN(authHeader, " ", 2)
		if !(len(parts) == 2 && parts[0] == "Bearer") {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "请求头中Authorization格式错误",
			})
			c.Abort()
			return
		}

		// 解析JWT token
		claims, err := utils.ParseCustomerToken(parts[1])
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "无效的Token",
			})
			c.Abort()
			return
		}

		// 将客户信息存储到上下文中
		c.Set("customer_id", claims.CustomerID)
		c.Set("customer_username", claims.Username)
		c.Set("customer_email", claims.Email)

		c.Next()
	}
}
