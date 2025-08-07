package api

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"

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
			"success": false,
			"message": err.Error(),
		})
		return
	}

	// 在数据库中查找用户
	var dbUser model.User
	if err := global.DB.Where("username = ? AND password = ?", user.Username, user.Password).First(&dbUser).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"message": "用户名或密码错误",
		})
		return
	}

	// 生成JWT Token
	token, err := utils.GenerateToken(dbUser.ID, dbUser.Username, dbUser.TenantID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "生成Token失败",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
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
			"success": false,
			"message": err.Error(),
		})
		return
	}

	// 在数据库中创建用户
	if err := global.DB.Create(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "注册失败",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
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
			"success": false,
			"message": "无法获取用户信息",
		})
		return
	}

	// 在数据库中查找用户
	var user model.User
	if err := global.DB.First(&user, userID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"message": "用户不存在",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "获取用户信息成功",
		"user":    user,
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
		"success": true,
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
			"success": false,
			"message": "获取用户列表失败",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "获取用户列表成功",
		"users":   users,
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
	var requestData struct {
		Username string `json:"username"`
		Email    string `json:"email"`
		Password string `json:"password"`
		TenantID uint   `json:"tenant_id"`
		Status   int    `json:"status"`
		RoleIDs  []uint  `json:"role_ids"`
	}
	
	// 绑定JSON到requestData
	if err := c.ShouldBindJSON(&requestData); err != nil {
		fmt.Printf("绑定JSON错误: %v\n", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "请求参数错误",
		})
		return
	}

	// 验证密码
	if requestData.Password == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "密码不能为空",
		})
		return
	}

	// 创建用户对象
	user := model.User{
		Username: requestData.Username,
		Email:    requestData.Email,
		Password: requestData.Password,
		TenantID: requestData.TenantID,
		Status:   requestData.Status,
	}

	fmt.Printf("接收到的用户数据: %+v\n", user)

	// 创建用户
	if err := global.DB.Create(&user).Error; err != nil {
		fmt.Printf("创建用户错误: %v\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "创建用户失败",
		})
		return
	}

	// 分配角色
	if len(requestData.RoleIDs) > 0 {
		var userRoles []model.UserRole
		for _, roleID := range requestData.RoleIDs {
			userRoles = append(userRoles, model.UserRole{
				UserID:   user.ID,
				RoleID:   roleID,
				TenantID: user.TenantID,
			})
		}
		
		if err := global.DB.Create(&userRoles).Error; err != nil {
			fmt.Printf("分配角色错误: %v\n", err)
			// 如果角色分配失败，不返回错误，只记录日志
			fmt.Printf("用户创建成功但角色分配失败: %v\n", err)
		}
	}

	fmt.Printf("用户创建成功: %+v\n", user)
	c.JSON(http.StatusOK, gin.H{
		"success": true,
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
	var requestData struct {
		Username string `json:"username"`
		Email    string `json:"email"`
		Password string `json:"password"`
		TenantID uint   `json:"tenant_id"`
		Status   int    `json:"status"`
		RoleIDs  []uint  `json:"role_ids"`
	}
	
	// 绑定JSON到requestData
	if err := c.ShouldBindJSON(&requestData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "请求参数错误",
		})
		return
	}

	// 如果密码为空，则不更新密码
	updateData := map[string]interface{}{
		"username":  requestData.Username,
		"email":     requestData.Email,
		"status":    requestData.Status,
		"tenant_id": requestData.TenantID,
	}

	// 如果有密码，则添加密码更新
	if requestData.Password != "" {
		updateData["password"] = requestData.Password
	}

	// 更新用户
	if err := global.DB.Model(&model.User{}).Where("id = ?", id).Updates(updateData).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "更新用户失败",
		})
		return
	}

	// 更新角色关联
	if requestData.RoleIDs != nil {
		// 先删除现有的角色关联
		if err := global.DB.Where("user_id = ?", id).Delete(&model.UserRole{}).Error; err != nil {
			fmt.Printf("删除旧角色关联错误: %v\n", err)
		}

		// 如果有角色ID，则添加新的角色关联
		if len(requestData.RoleIDs) > 0 {
			var userRoles []model.UserRole
			userID, _ := strconv.Atoi(id)
			for _, roleID := range requestData.RoleIDs {
				userRoles = append(userRoles, model.UserRole{
					UserID:   uint(userID),
					RoleID:   roleID,
					TenantID: requestData.TenantID,
				})
			}
			
			if err := global.DB.Create(&userRoles).Error; err != nil {
				fmt.Printf("分配新角色错误: %v\n", err)
				// 如果角色分配失败，不返回错误，只记录日志
			}
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
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
		"success": true,
		"message": "用户删除成功",
	})
}

// UploadAvatar 上传头像
// @Summary 上传用户头像
// @Description 上传并更新用户头像
// @Tags 用户管理
// @Accept multipart/form-data
// @Produce json
// @Security ApiKeyAuth
// @Param avatar formData file true "头像文件"
// @Success 200 {object} map[string]interface{} "{"message":"头像上传成功","avatar_url":"string"}"
// @Failure 400 {object} map[string]interface{} "{"error":"头像文件不能为空"}"
// @Failure 500 {object} map[string]interface{} "{"error":"头像上传失败"}"
// @Router /api/user/upload-avatar [post]
func UploadAvatar(c *gin.Context) {
	// 从上下文中获取用户ID
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"message": "无法获取用户信息",
		})
		return
	}

	// 获取上传的文件
	file, err := c.FormFile("avatar")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "头像文件不能为空",
		})
		return
	}

	// 验证文件类型
	allowedTypes := map[string]bool{
		"image/jpeg": true,
		"image/png":  true,
		"image/jpg":  true,
	}
	if !allowedTypes[file.Header.Get("Content-Type")] {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "只允许上传JPG、PNG格式的图片",
		})
		return
	}

	// 验证文件大小 (最大2MB)
	if file.Size > 2*1024*1024 {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "头像文件大小不能超过2MB",
		})
		return
	}

	// 生成文件名
	ext := ".jpg"
	if file.Header.Get("Content-Type") == "image/png" {
		ext = ".png"
	}
	filename := fmt.Sprintf("avatar_%d_%d%s", userID, time.Now().Unix(), ext)

	// 确保上传目录存在
	uploadDir := "./uploads/avatars"
	if err := os.MkdirAll(uploadDir, 0755); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "创建上传目录失败",
		})
		return
	}

	// 保存文件
	filepath := filepath.Join(uploadDir, filename)
	if err := c.SaveUploadedFile(file, filepath); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "保存头像文件失败",
		})
		return
	}

	// 生成访问URL
	avatarURL := fmt.Sprintf("/uploads/avatars/%s", filename)

	// 更新用户头像信息
	if err := global.DB.Model(&model.User{}).Where("id = ?", userID).Update("avatar", avatarURL).Error; err != nil {
		// 删除已上传的文件
		os.Remove(filepath)
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "更新用户头像信息失败",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "头像上传成功",
		"avatar_url": avatarURL,
	})
}
