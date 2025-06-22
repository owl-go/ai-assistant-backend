package utils

import (
	"context"
	"errors"
	"fmt"
	"time"

	"ai-assistant-backend/config"

	"github.com/golang-jwt/jwt/v5"
)

type Claims struct {
	UserID   uint   `json:"user_id"`
	Username string `json:"username"`
	jwt.RegisteredClaims
}

// GenerateToken 生成JWT令牌并存储到Redis
func GenerateToken(userID uint, username string) (string, error) {
	if config.GlobalConfig == nil {
		return "", errors.New("配置未加载")
	}

	// 创建声明
	claims := Claims{
		UserID:   userID,
		Username: username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Duration(config.GlobalConfig.JWT.ExpireHours) * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
		},
	}

	// 创建token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(config.GlobalConfig.JWT.Secret))
	if err != nil {
		return "", err
	}

	// 将token存储到Redis，设置过期时间
	ctx := context.Background()
	redisKey := fmt.Sprintf("token:%d", userID)
	expiration := time.Duration(config.GlobalConfig.JWT.ExpireHours) * time.Hour

	err = config.RedisClient.Set(ctx, redisKey, tokenString, expiration).Err()
	if err != nil {
		return "", fmt.Errorf("存储token到Redis失败: %v", err)
	}

	return tokenString, nil
}

// ParseToken 解析JWT令牌
func ParseToken(tokenString string) (*Claims, error) {
	if config.GlobalConfig == nil {
		return nil, errors.New("配置未加载")
	}

	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(config.GlobalConfig.JWT.Secret), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}

	return nil, errors.New("无效的token")
}

// ValidateToken 验证token是否在Redis中存在且有效
func ValidateToken(userID uint, tokenString string) (bool, error) {
	if config.RedisClient == nil {
		return false, errors.New("Redis客户端未初始化")
	}

	ctx := context.Background()
	redisKey := fmt.Sprintf("token:%d", userID)

	// 从Redis获取存储的token
	storedToken, err := config.RedisClient.Get(ctx, redisKey).Result()
	if err != nil {
		return false, fmt.Errorf("从Redis获取token失败: %v", err)
	}

	// 比较token是否一致
	if storedToken != tokenString {
		return false, errors.New("token不匹配")
	}

	return true, nil
}

// InvalidateToken 使token失效（从Redis中删除）
func InvalidateToken(userID uint) error {
	if config.RedisClient == nil {
		return errors.New("Redis客户端未初始化")
	}

	ctx := context.Background()
	redisKey := fmt.Sprintf("token:%d", userID)

	err := config.RedisClient.Del(ctx, redisKey).Err()
	if err != nil {
		return fmt.Errorf("从Redis删除token失败: %v", err)
	}

	return nil
}

// RefreshToken 刷新token
func RefreshToken(userID uint, username string) (string, error) {
	// 先使旧token失效
	if err := InvalidateToken(userID); err != nil {
		return "", fmt.Errorf("使旧token失效失败: %v", err)
	}

	// 生成新token
	return GenerateToken(userID, username)
}
