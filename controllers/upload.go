package controllers

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"ai-assistant-backend/config"

	"github.com/gin-gonic/gin"
)

// UploadFile 文件上传
func UploadFile(c *gin.Context) {
	if config.GlobalConfig == nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "配置未加载",
		})
		return
	}

	// 获取上传的文件
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "文件上传失败",
			"error":   err.Error(),
		})
		return
	}

	// 检查文件大小
	if file.Size > config.GlobalConfig.Upload.MaxFileSize {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "文件大小超过限制",
		})
		return
	}

	// 检查文件类型
	ext := filepath.Ext(file.Filename)
	allowed := false
	for _, allowedType := range config.GlobalConfig.Upload.AllowedTypes {
		if ext == allowedType {
			allowed = true
			break
		}
	}

	if !allowed {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "不支持的文件类型",
		})
		return
	}

	// 创建上传目录
	uploadPath := config.GlobalConfig.Upload.Path
	if uploadPath == "" {
		uploadPath = "./uploads"
	}

	// 按日期创建子目录
	dateDir := time.Now().Format("2006-01-02")
	fullPath := filepath.Join(uploadPath, dateDir)

	if err := os.MkdirAll(fullPath, 0755); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "创建目录失败",
		})
		return
	}

	// 生成唯一文件名
	timestamp := time.Now().UnixNano()
	filename := fmt.Sprintf("%d_%s", timestamp, file.Filename)
	filePath := filepath.Join(fullPath, filename)

	// 保存文件
	if err := c.SaveUploadedFile(file, filePath); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "保存文件失败",
		})
		return
	}

	// 返回文件信息
	fileURL := fmt.Sprintf("/uploads/%s/%s", dateDir, filename)

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "上传成功",
		"data": gin.H{
			"filename": file.Filename,
			"size":     file.Size,
			"url":      fileURL,
			"path":     filePath,
		},
	})
}

// ServeFile 提供文件访问
func ServeFile(c *gin.Context) {
	filename := c.Param("filename")
	dateDir := c.Param("date")

	uploadPath := config.GlobalConfig.Upload.Path
	if uploadPath == "" {
		uploadPath = "./uploads"
	}

	filePath := filepath.Join(uploadPath, dateDir, filename)

	// 检查文件是否存在
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		c.JSON(http.StatusNotFound, gin.H{
			"code":    404,
			"message": "文件不存在",
		})
		return
	}

	c.File(filePath)
}
