package utils

import (
	"context"
	"fmt"
	"io"
	"mime/multipart"
	"path/filepath"
	"strings"
	"time"

	"ai-assistant-backend/config"

	"github.com/minio/minio-go/v7"
)

// MinIOUploader MinIO文件上传器
type MinIOUploader struct {
	client *minio.Client
	bucket string
}

// NewMinIOUploader 创建MinIO上传器
func NewMinIOUploader() *MinIOUploader {
	return &MinIOUploader{
		client: config.GetMinIOClient(),
		bucket: config.GlobalConfig.MinIO.Bucket,
	}
}

// UploadFile 上传文件到MinIO
func (m *MinIOUploader) UploadFile(file *multipart.FileHeader, objectName string) (string, error) {
	// 打开文件
	src, err := file.Open()
	if err != nil {
		return "", fmt.Errorf("打开文件失败: %v", err)
	}
	defer src.Close()

	// 获取文件内容类型
	contentType := file.Header.Get("Content-Type")
	if contentType == "" {
		contentType = "application/octet-stream"
	}

	// 上传文件到MinIO
	ctx := context.Background()
	_, err = m.client.PutObject(ctx, m.bucket, objectName, src, file.Size, minio.PutObjectOptions{
		ContentType: contentType,
	})
	if err != nil {
		return "", fmt.Errorf("上传文件到MinIO失败: %v", err)
	}

	// 生成文件URL
	fileURL, err := m.GetFileURL(objectName)
	if err != nil {
		return "", err
	}
	return fileURL, nil
}

// GetFileURL 获取文件URL（临时访问URL）
func (m *MinIOUploader) GetFileURL(objectName string) (string, error) {
	// 获取配置的链接有效时间
	expireHours := config.GlobalConfig.MinIO.URLExpireHours
	if expireHours <= 0 {
		expireHours = 24 // 默认24小时
	}

	// 生成预签名下载URL
	ctx := context.Background()
	url, err := m.client.PresignedGetObject(ctx, m.bucket, objectName, time.Duration(expireHours)*time.Hour, nil)
	if err != nil {
		return "", fmt.Errorf("生成文件访问URL失败: %v", err)
	}

	return url.String(), nil
}

// DeleteFile 删除文件
func (m *MinIOUploader) DeleteFile(objectName string) error {
	ctx := context.Background()
	err := m.client.RemoveObject(ctx, m.bucket, objectName, minio.RemoveObjectOptions{})
	if err != nil {
		return fmt.Errorf("删除文件失败: %v", err)
	}
	return nil
}

// FileExists 检查文件是否存在
func (m *MinIOUploader) FileExists(objectName string) (bool, error) {
	ctx := context.Background()
	_, err := m.client.StatObject(ctx, m.bucket, objectName, minio.StatObjectOptions{})
	if err != nil {
		if strings.Contains(err.Error(), "NoSuchKey") {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

// GetFileInfo 获取文件信息
func (m *MinIOUploader) GetFileInfo(objectName string) (*minio.ObjectInfo, error) {
	ctx := context.Background()
	info, err := m.client.StatObject(ctx, m.bucket, objectName, minio.StatObjectOptions{})
	if err != nil {
		return nil, fmt.Errorf("获取文件信息失败: %v", err)
	}
	return &info, nil
}

// GenerateObjectName 生成对象名称
func GenerateObjectName(filename string, prefix string) string {
	// 获取文件扩展名
	ext := filepath.Ext(filename)
	// 生成时间戳
	timestamp := time.Now().UnixNano()
	// 生成唯一文件名
	uniqueName := fmt.Sprintf("%d%s", timestamp, ext)

	// 按日期创建目录结构
	dateDir := time.Now().Format("2006-01-02")

	if prefix != "" {
		return fmt.Sprintf("%s/%s/%s", prefix, dateDir, uniqueName)
	}
	return fmt.Sprintf("%s/%s", dateDir, uniqueName)
}

// UploadFileWithValidation 上传文件并验证
func UploadFileWithValidation(file *multipart.FileHeader, allowedTypes []string, maxSize int64, prefix string) (string, string, error) {
	// 验证文件大小
	if file.Size > maxSize {
		return "", "", fmt.Errorf("文件大小超过限制")
	}

	// 验证文件类型
	ext := strings.ToLower(filepath.Ext(file.Filename))
	allowed := false
	for _, allowedType := range allowedTypes {
		if strings.ToLower(allowedType) == ext {
			allowed = true
			break
		}
	}
	if !allowed {
		return "", "", fmt.Errorf("不支持的文件类型: %s", ext)
	}

	// 生成对象名称
	objectName := GenerateObjectName(file.Filename, prefix)

	// 上传文件
	uploader := NewMinIOUploader()
	fileURL, err := uploader.UploadFile(file, objectName)
	if err != nil {
		return "", "", err
	}

	return objectName, fileURL, nil
}

// GetPresignedURL 获取预签名URL（用于直接上传）
func (m *MinIOUploader) GetPresignedURL(objectName string, expires time.Duration) (string, error) {
	ctx := context.Background()
	url, err := m.client.PresignedPutObject(ctx, m.bucket, objectName, expires)
	if err != nil {
		return "", fmt.Errorf("生成预签名URL失败: %v", err)
	}
	return url.String(), nil
}

// GetPresignedGetURL 获取预签名下载URL
func (m *MinIOUploader) GetPresignedGetURL(objectName string, expires time.Duration) (string, error) {
	ctx := context.Background()
	url, err := m.client.PresignedGetObject(ctx, m.bucket, objectName, expires, nil)
	if err != nil {
		return "", fmt.Errorf("生成预签名下载URL失败: %v", err)
	}
	return url.String(), nil
}

// CopyFile 复制文件
func (m *MinIOUploader) CopyFile(srcObjectName, dstObjectName string) error {
	ctx := context.Background()
	src := minio.CopySrcOptions{
		Bucket: m.bucket,
		Object: srcObjectName,
	}
	dst := minio.CopyDestOptions{
		Bucket: m.bucket,
		Object: dstObjectName,
	}
	_, err := m.client.CopyObject(ctx, dst, src)
	if err != nil {
		return fmt.Errorf("复制文件失败: %v", err)
	}
	return nil
}

// ListFiles 列出文件
func (m *MinIOUploader) ListFiles(prefix string, recursive bool) ([]minio.ObjectInfo, error) {
	ctx := context.Background()
	var objects []minio.ObjectInfo

	opts := minio.ListObjectsOptions{
		Prefix:    prefix,
		Recursive: recursive,
	}

	for obj := range m.client.ListObjects(ctx, m.bucket, opts) {
		if obj.Err != nil {
			return nil, fmt.Errorf("列出文件失败: %v", obj.Err)
		}
		objects = append(objects, obj)
	}

	return objects, nil
}

// DownloadFile 下载文件
func (m *MinIOUploader) DownloadFile(objectName string) (io.Reader, error) {
	ctx := context.Background()
	obj, err := m.client.GetObject(ctx, m.bucket, objectName, minio.GetObjectOptions{})
	if err != nil {
		return nil, fmt.Errorf("下载文件失败: %v", err)
	}
	return obj, nil
}

// GetTemporaryFileURL 获取临时文件访问URL
func (m *MinIOUploader) GetTemporaryFileURL(objectName string, expireHours int) (string, error) {
	if expireHours <= 0 {
		expireHours = config.GlobalConfig.MinIO.URLExpireHours
		if expireHours <= 0 {
			expireHours = 24 // 默认24小时
		}
	}

	ctx := context.Background()
	url, err := m.client.PresignedGetObject(ctx, m.bucket, objectName, time.Duration(expireHours)*time.Hour, nil)
	if err != nil {
		return "", fmt.Errorf("生成临时文件访问URL失败: %v", err)
	}

	return url.String(), nil
}
