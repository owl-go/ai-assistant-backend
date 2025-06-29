package config

import (
	"context"
	"fmt"
	"log"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

// MinIOClient MinIO客户端实例
var MinIOClient *minio.Client

// InitMinIO 初始化MinIO客户端
func InitMinIO() error {
	if GlobalConfig == nil {
		return fmt.Errorf("配置未加载")
	}

	// 创建MinIO客户端
	minioClient, err := minio.New(GlobalConfig.MinIO.Endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(GlobalConfig.MinIO.AccessKey, GlobalConfig.MinIO.SecretKey, ""),
		Secure: GlobalConfig.MinIO.UseSSL,
	})
	if err != nil {
		return fmt.Errorf("创建MinIO客户端失败: %v", err)
	}

	// 测试连接
	ctx := context.Background()
	_, err = minioClient.ListBuckets(ctx)
	if err != nil {
		return fmt.Errorf("MinIO连接测试失败: %v", err)
	}

	// 检查并创建默认bucket
	bucketExists, err := minioClient.BucketExists(ctx, GlobalConfig.MinIO.Bucket)
	if err != nil {
		return fmt.Errorf("检查bucket失败: %v", err)
	}
	if !bucketExists {
		// 创建bucket
		err = minioClient.MakeBucket(ctx, GlobalConfig.MinIO.Bucket, minio.MakeBucketOptions{
			Region: GlobalConfig.MinIO.Location,
		})
		if err != nil {
			return fmt.Errorf("创建bucket失败: %v", err)
		}

		log.Printf("Bucket '%s' 创建成功（私有访问）", GlobalConfig.MinIO.Bucket)
	}

	MinIOClient = minioClient
	log.Println("MinIO客户端初始化成功")
	return nil
}

// GetMinIOClient 获取MinIO客户端实例
func GetMinIOClient() *minio.Client {
	return MinIOClient
}
