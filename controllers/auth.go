package controllers

import (
	"context"
	"fmt"

	"ai-assistant-backend/config"
	"ai-assistant-backend/models"
	"ai-assistant-backend/utils"

	"github.com/gin-gonic/gin"
)

type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type RegisterRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
}

// Login 用户登录
func Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ValidationError(c, err.Error())
		return
	}

	// 查找用户
	var user models.User
	if err := config.DB.Where("username = ?", req.Username).First(&user).Error; err != nil {
		utils.LoginFailed(c)
		return
	}

	// 验证密码
	if !utils.CheckPassword(req.Password, user.Password) {
		utils.LoginFailed(c)
		return
	}

	// 生成JWT token并存储到Redis
	token, err := utils.GenerateToken(user.ID, user.Username, user.Role)
	if err != nil {
		utils.InternalServerError(c, "生成令牌失败")
		return
	}
	//将token存到redis
	err = config.RedisClient.Set(context.Background(), fmt.Sprintf("token:%s", token), user.ID, 0).Err()
	if err != nil {
		utils.InternalServerError(c, "保存令牌失败")
		return
	}

	utils.Success(c, gin.H{
		"token": token,
		"user": gin.H{
			"id":       user.ID,
			"username": user.Username,
			"email":    user.Email,
			"avatar":   user.Avatar,
			"role":     user.Role,
		},
	}, "登录成功")
}

// Logout 用户登出
func Logout(c *gin.Context) {
	userID := c.GetUint("user_id")

	// 使token失效
	if err := utils.InvalidateToken(userID); err != nil {
		utils.InternalServerError(c, "登出失败")
		return
	}

	utils.SuccessWithMessage(c, "登出成功")
}

// RefreshToken 刷新token
func RefreshToken(c *gin.Context) {
	user, err := utils.GetUserFromContext(c)

	if err != nil {
		utils.Unauthorized(c, "用户未登录")
		return
	}
	// 生成新token
	token, err := utils.RefreshToken(user.UserID, user.Username, user.Role)
	if err != nil {
		utils.InternalServerError(c, "刷新令牌失败")
		return
	}

	utils.Success(c, gin.H{
		"token": token,
	}, "刷新成功")
}

// Register 用户注册
func Register(c *gin.Context) {
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ValidationError(c, err.Error())
		return
	}

	// 检查用户名是否已存在
	var existingUser models.User
	if err := config.DB.Where("username = ?", req.Username).First(&existingUser).Error; err == nil {
		utils.UserAlreadyExists(c, "用户名")
		return
	}

	// 检查邮箱是否已存在
	if err := config.DB.Where("email = ?", req.Email).First(&existingUser).Error; err == nil {
		utils.UserAlreadyExists(c, "邮箱")
		return
	}

	// 加密密码
	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		utils.InternalServerError(c, "密码加密失败")
		return
	}

	// 创建用户
	user := models.User{
		Username: req.Username,
		Password: hashedPassword,
		Email:    req.Email,
	}

	if err := config.DB.Create(&user).Error; err != nil {
		utils.CreateFailed(c, "用户")
		return
	}

	utils.Success(c, gin.H{
		"id":       user.ID,
		"username": user.Username,
		"email":    user.Email,
	}, "注册成功")
}

// GetProfile 获取用户信息
func GetProfile(c *gin.Context) {
	userID := c.GetUint("user_id")

	var user models.User
	if err := config.DB.First(&user, userID).Error; err != nil {
		utils.UserNotFound(c)
		return
	}

	utils.Success(c, gin.H{
		"id":       user.ID,
		"username": user.Username,
		"email":    user.Email,
		"avatar":   user.Avatar,
		"role":     user.Role,
	}, "获取成功")
}
