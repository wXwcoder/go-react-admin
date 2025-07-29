package middleware

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"go-react-admin/global"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
)

// CasbinMiddleware Casbin权限验证中间件
func CasbinMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 获取JWT token
		tokenString := c.GetHeader("Authorization")
		if tokenString == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"code": 401,
				"msg":  "未提供认证令牌",
			})
			c.Abort()
			return
		}

		// 移除Bearer前缀
		if strings.HasPrefix(tokenString, "Bearer ") {
			tokenString = tokenString[7:]
		}

		// 解析JWT token
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return []byte(global.GlobalConfig.Jwt.Secret), nil // 从配置文件读取JWT密钥
		})

		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{
				"code": 401,
				"msg":  "无效的认证令牌",
			})
			c.Abort()
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{
				"code": 401,
				"msg":  "无效的令牌声明",
			})
			c.Abort()
			return
		}

		// 获取用户信息
		userID, ok := claims["user_id"].(float64)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{
				"code": 401,
				"msg":  "无效的用户ID",
			})
			c.Abort()
			return
		}

		tenantID, ok := claims["tenant_id"].(float64)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{
				"code": 401,
				"msg":  "无效的租户ID",
			})
			c.Abort()
			return
		}

		// 设置用户信息到上下文
		c.Set("user_id", uint(userID))
		c.Set("tenant_id", uint(tenantID))

		// 获取请求路径和方法
		path := c.Request.URL.Path
		method := c.Request.Method

		// 构建Casbin检查参数
		sub := strconv.Itoa(int(userID))
		obj := path
		act := method
		tenant := strconv.Itoa(int(tenantID))

		// 使用Casbin进行权限检查
		enforcer := global.Enforcer
		if enforcer == nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"code": 500,
				"msg":  "权限系统未初始化",
			})
			c.Abort()
			return
		}

		// 检查权限
		allowed, err := enforcer.Enforce(sub, obj, act, tenant)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"code": 500,
				"msg":  "权限检查失败: " + err.Error(),
			})
			c.Abort()
			return
		}

		if !allowed {
			c.JSON(http.StatusForbidden, gin.H{
				"code": 403,
				"msg":  "权限不足",
			})
			c.Abort()
			return
		}

		c.Next()
	}
}

// RequirePermission 需要特定权限的中间件
func RequirePermission(resource, action string) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, exists := c.Get("user_id")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{
				"code": 401,
				"msg":  "用户未认证",
			})
			c.Abort()
			return
		}

		tenantID, exists := c.Get("tenant_id")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{
				"code": 401,
				"msg":  "租户信息缺失",
			})
			c.Abort()
			return
		}

		// 构建Casbin检查参数
		sub := strconv.Itoa(int(userID.(uint)))
		obj := resource
		act := action
		tenant := strconv.Itoa(int(tenantID.(uint)))

		// 使用Casbin进行权限检查
		enforcer := global.Enforcer
		if enforcer == nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"code": 500,
				"msg":  "权限系统未初始化",
			})
			c.Abort()
			return
		}

		allowed, err := enforcer.Enforce(sub, obj, act, tenant)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"code": 500,
				"msg":  "权限检查失败: " + err.Error(),
			})
			c.Abort()
			return
		}

		if !allowed {
			c.JSON(http.StatusForbidden, gin.H{
				"code": 403,
				"msg":  fmt.Sprintf("缺少权限: %s:%s", resource, action),
			})
			c.Abort()
			return
		}

		c.Next()
	}
}
