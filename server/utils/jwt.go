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

// CustomerClaims 第三方客户JWT声明
type CustomerClaims struct {
	CustomerID uint   `json:"customer_id"`
	Username   string `json:"username"`
	Email      string `json:"email"`
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

// GenerateCustomerToken 生成第三方客户JWT Token
func GenerateCustomerToken(customerID uint64) (string, int64, error) {
	// 设置过期时间
	expireTime := time.Now().Add(time.Duration(global.GlobalConfig.Jwt.Expire) * time.Hour)
	expiresIn := int64(global.GlobalConfig.Jwt.Expire * 3600)

	// 获取客户信息
	var customer struct {
		Username string `gorm:"column:username"`
		Email    string `gorm:"column:email"`
	}

	// 从数据库获取客户信息
	if err := global.DB.Table("customers").Select("username, email").Where("id = ?", customerID).First(&customer).Error; err != nil {
		return "", 0, err
	}

	// 创建声明
	claims := CustomerClaims{
		CustomerID: uint(customerID),
		Username:   customer.Username,
		Email:      customer.Email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expireTime),
			Issuer:    "go-react-admin",
		},
	}

	// 创建token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// 签名并获得完整的编码后的字符串token
	tokenString, err := token.SignedString([]byte(global.GlobalConfig.Jwt.Secret))
	return tokenString, expiresIn, err
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

// ParseCustomerToken 解析第三方客户JWT Token
func ParseCustomerToken(tokenString string) (*CustomerClaims, error) {
	// 解析token
	token, err := jwt.ParseWithClaims(tokenString, &CustomerClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(global.GlobalConfig.Jwt.Secret), nil
	})

	if err != nil {
		return nil, err
	}

	// 类型断言
	if claims, ok := token.Claims.(*CustomerClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, err
}

// GenerateRandomString 生成随机字符串
func GenerateRandomString(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	result := make([]byte, length)
	for i := range result {
		result[i] = charset[time.Now().UnixNano()%int64(len(charset))]
	}
	return string(result)
}
