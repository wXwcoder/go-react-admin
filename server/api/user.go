package api

import (
	"fmt"
	"net/http"
	"strconv"

	"go-react-admin/global"
	"go-react-admin/model"
	"go-react-admin/utils"

	"github.com/gin-gonic/gin"
)

// Login 用户登录
// @Summary 用户登录
// @Description 用户登录接口，验证用户名密码并返回JWT Token
// @Tags 用户管理
// @Accept json
// @Produce json
// @Param user body model.User true "用户登录信息"
// @Success 200 {object} map[string]interface{} "{"message":"登录成功","token":"string","userId":"uint"}"
// @Failure 400 {object} map[string]interface{} "{"error":"string"}"
// @Failure 401 {object} map[string]interface{} "{"error":"用户名或密码错误"}"
// @Failure 500 {object} map[string]interface{} "{"error":"生成Token失败"}"
// @Router /api/login [post]
func Login(c *gin.Context) {
	var user model.User
	// 绑定JSON
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	// 在数据库中查找用户
	var dbUser model.User
	if err := global.DB.Where("username = ? AND password = ?", user.Username, user.Password).First(&dbUser).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "用户名或密码错误",
		})
		return
	}

	// 生成JWT Token
	token, err := utils.GenerateToken(dbUser.ID, dbUser.Username, dbUser.TenantID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "生成Token失败",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "登录成功",
		"token":   token,
		"userId":  dbUser.ID,
	})
}

// Register 用户注册
// @Summary 用户注册
// @Description 用户注册接口，创建新用户
// @Tags 用户管理
// @Accept json
// @Produce json
// @Param user body model.User true "用户注册信息"
// @Success 200 {object} map[string]interface{} "{"message":"注册成功"}"
// @Failure 400 {object} map[string]interface{} "{"error":"string"}"
// @Failure 500 {object} map[string]interface{} "{"error":"注册失败"}"
// @Router /api/register [post]
func Register(c *gin.Context) {
	var user model.User
	// 绑定JSON
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	// 在数据库中创建用户
	if err := global.DB.Create(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "注册失败",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "注册成功",
	})
}

// GetUserInfo 获取用户信息
// @Summary 获取当前用户信息
// @Description 根据JWT Token获取当前登录用户的信息
// @Tags 用户管理
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Success 200 {object} map[string]interface{} "{"user":model.User}"
// @Failure 401 {object} map[string]interface{} "{"error":"无法获取用户信息"}"
// @Failure 404 {object} map[string]interface{} "{"error":"用户不存在"}"
// @Router /api/user/info [get]
func GetUserInfo(c *gin.Context) {
	// 从上下文中获取用户ID
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "无法获取用户信息",
		})
		return
	}

	// 在数据库中查找用户
	var user model.User
	if err := global.DB.First(&user, userID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "用户不存在",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"user": user,
	})
}

// Logout 用户登出
// @Summary 用户登出
// @Description 用户登出接口
// @Tags 用户管理
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Success 200 {object} map[string]interface{} "{"message":"登出成功"}"
// @Router /api/logout [post]
func Logout(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "登出成功",
	})
}

// GetUserList 获取用户列表
// @Summary 获取用户列表
// @Description 获取所有用户的列表
// @Tags 用户管理
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Success 200 {object} map[string]interface{} "{"users":[]model.User}"
// @Failure 500 {object} map[string]interface{} "{"error":"获取用户列表失败"}"
// @Router /api/users [get]
func GetUserList(c *gin.Context) {
	var users []model.User
	// 从数据库中获取所有用户
	if err := global.DB.Find(&users).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "获取用户列表失败",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"users": users,
	})
}

// CreateUser 创建用户
// @Summary 创建用户
// @Description 创建新用户
// @Tags 用户管理
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param user body model.User true "用户创建信息"
// @Success 200 {object} map[string]interface{} "{"message":"用户创建成功","user":model.User}"
// @Failure 400 {object} map[string]interface{} "{"error":"请求参数错误"}"
// @Failure 500 {object} map[string]interface{} "{"error":"创建用户失败"}"
// @Router /api/users [post]
func CreateUser(c *gin.Context) {
	var user model.User
	// 绑定JSON到user
	if err := c.ShouldBindJSON(&user); err != nil {
		fmt.Printf("绑定JSON错误: %v\n", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "请求参数错误",
		})
		return
	}

	fmt.Printf("接收到的用户数据: %+v\n", user)

	// 创建用户
	if err := global.DB.Create(&user).Error; err != nil {
		fmt.Printf("创建用户错误: %v\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "创建用户失败",
		})
		return
	}

	fmt.Printf("用户创建成功: %+v\n", user)
	c.JSON(http.StatusOK, gin.H{
		"message": "用户创建成功",
		"user":    user,
	})
}

// UpdateUser 更新用户
// @Summary 更新用户
// @Description 根据用户ID更新用户信息
// @Tags 用户管理
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param id path int true "用户ID"
// @Param user body model.User true "用户更新信息"
// @Success 200 {object} map[string]interface{} "{"message":"用户更新成功"}"
// @Failure 400 {object} map[string]interface{} "{"error":"请求参数错误"}"
// @Failure 500 {object} map[string]interface{} "{"error":"更新用户失败"}"
// @Router /api/users/{id} [put]
func UpdateUser(c *gin.Context) {
	id := c.Param("id")
	var user model.User
	// 绑定JSON到user
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "请求参数错误",
		})
		return
	}

	// 更新用户
	if err := global.DB.Model(&model.User{}).Where("id = ?", id).Updates(user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "更新用户失败",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "用户更新成功",
	})
}

// DeleteUser 删除用户
// @Summary 删除用户
// @Description 根据用户ID删除用户
// @Tags 用户管理
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param id path int true "用户ID"
// @Success 200 {object} map[string]interface{} "{"message":"用户删除成功"}"
// @Failure 400 {object} map[string]interface{} "{"error":"无效的用户ID"}"
// @Failure 500 {object} map[string]interface{} "{"error":"删除用户失败"}"
// @Router /api/users/{id} [delete]
func DeleteUser(c *gin.Context) {
	id := c.Param("id")
	userID, err := strconv.Atoi(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "无效的用户ID",
		})
		return
	}

	// 删除用户
	if err := global.DB.Delete(&model.User{}, userID).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "删除用户失败",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "用户删除成功",
	})
}
