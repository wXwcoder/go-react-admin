package utils

import (
	"time"

	"go-react-admin/global"

	"github.com/golang-jwt/jwt/v4"
)

// CustomClaims 自定义JWT声明
type CustomClaims struct {
	UserID   uint   `json:"user_id"`
	Username string `json:"username"`
	TenantID uint   `json:"tenant_id"`
	jwt.RegisteredClaims
}

// GenerateToken 生成JWT Token
func GenerateToken(userID uint, username string, tenantID uint) (string, error) {
	// 设置过期时间
	expireTime := time.Now().Add(time.Duration(global.GlobalConfig.Jwt.Expire) * time.Hour)

	// 创建声明
	claims := CustomClaims{
		UserID:   userID,
		Username: username,
		TenantID: tenantID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expireTime),
			Issuer:    "go-react-admin",
		},
	}

	// 创建token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// 签名并获得完整的编码后的字符串token
	tokenString, err := token.SignedString([]byte(global.GlobalConfig.Jwt.Secret))
	return tokenString, err
}

// ParseToken 解析JWT Token
func ParseToken(tokenString string) (*CustomClaims, error) {
	// 解析token
	token, err := jwt.ParseWithClaims(tokenString, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(global.GlobalConfig.Jwt.Secret), nil
	})

	if err != nil {
		return nil, err
	}

	// 类型断言
	if claims, ok := token.Claims.(*CustomClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, err
}
