package utils

import (
	"errors"

	"github.com/gin-gonic/gin"
)

// 上下文键常量
const (
	UserKey = "user" // 新增：存储完整用户对象的键
)

// SetUserContextToContext 将用户上下文对象存储到上下文中
func SetUserToContext(c *gin.Context, user interface{}) {
	if user != nil {
		c.Set(UserKey, user)
	}
}

// GetUserFromContext 从上下文中获取完整的用户对象
func GetUserFromContext(c *gin.Context) (interface{}, error) {
	userInterface, exists := c.Get(UserKey)
	if !exists {
		return nil, errors.New("用户对象不存在于上下文中")
	}

	return userInterface, nil
}
