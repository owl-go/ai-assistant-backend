package config

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/go-redis/redis/v8"
)

var RedisClient *redis.Client

// InitRedis 初始化Redis连接
func InitRedis() error {
	if GlobalConfig == nil {
		return fmt.Errorf("配置未加载，请先调用 LoadConfig")
	}

	// 创建Redis客户端
	RedisClient = redis.NewClient(&redis.Options{
		Addr:         GlobalConfig.Redis.GetRedisAddr(),
		Password:     GlobalConfig.Redis.Password,
		DB:           GlobalConfig.Redis.DB,
		PoolSize:     GlobalConfig.Redis.PoolSize,
		MinIdleConns: GlobalConfig.Redis.MinIdleConns,
	})

	// 测试连接
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := RedisClient.Ping(ctx).Result()
	if err != nil {
		return fmt.Errorf("Redis连接失败: %v", err)
	}

	log.Println("Redis连接成功")
	return nil
}

// CloseRedis 关闭Redis连接
func CloseRedis() error {
	if RedisClient != nil {
		return RedisClient.Close()
	}
	return nil
}
