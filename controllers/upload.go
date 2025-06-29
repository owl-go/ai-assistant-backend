package controllers

import (
	"ai-assistant-backend/config"
	"ai-assistant-backend/utils"

	"fmt"

	"github.com/gin-gonic/gin"
)

// UploadFile 文件上传
func UploadFile(c *gin.Context) {
	if config.GlobalConfig == nil {
		utils.InternalServerError(c, "配置未加载")
		return
	}

	// 获取上传的文件
	file, err := c.FormFile("file")
	if err != nil {
		utils.BadRequestWithDetail(c, "文件上传失败", err.Error())
		return
	}

	// 使用MinIO上传文件
	objectName, fileURL, err := utils.UploadFileWithValidation(
		file,
		config.GlobalConfig.Upload.AllowedTypes,
		config.GlobalConfig.Upload.MaxFileSize,
		"uploads", // 文件前缀
	)
	if err != nil {
		utils.BadRequestWithDetail(c, "文件上传失败", err.Error())
		return
	}
	// 返回文件信息
	utils.Success(c, gin.H{
		"filename":    file.Filename,
		"size":        file.Size,
		"url":         fileURL,
		"object_name": objectName,
	}, "上传成功")
}

// GetPresignedUploadURL 获取预签名上传URL
func GetPresignedUploadURL(c *gin.Context) {
	filename := c.Query("filename")
	if filename == "" {
		utils.BadRequest(c, "文件名不能为空")
		return
	}

	// 生成对象名称
	objectName := utils.GenerateObjectName(filename, "uploads")

	// 获取预签名URL
	uploader := utils.NewMinIOUploader()
	presignedURL, err := uploader.GetPresignedURL(objectName, 3600) // 1小时有效期
	if err != nil {
		utils.InternalServerError(c, "生成预签名URL失败")
		return
	}

	utils.Success(c, gin.H{
		"upload_url":  presignedURL,
		"object_name": objectName,
		"expires_in":  3600,
	}, "获取预签名URL成功")
}

// GetPresignedDownloadURL 获取预签名下载URL
func GetPresignedDownloadURL(c *gin.Context) {
	objectName := c.Query("object_name")
	if objectName == "" {
		utils.BadRequest(c, "对象名称不能为空")
		return
	}

	// 获取预签名下载URL
	uploader := utils.NewMinIOUploader()
	presignedURL, err := uploader.GetPresignedGetURL(objectName, 3600) // 1小时有效期
	if err != nil {
		utils.InternalServerError(c, "生成预签名下载URL失败")
		return
	}

	utils.Success(c, gin.H{
		"download_url": presignedURL,
		"expires_in":   3600,
	}, "获取预签名下载URL成功")
}

// DeleteFile 删除文件
func DeleteFile(c *gin.Context) {
	objectName := c.Query("object_name")
	if objectName == "" {
		utils.BadRequest(c, "对象名称不能为空")
		return
	}

	// 删除文件
	uploader := utils.NewMinIOUploader()
	err := uploader.DeleteFile(objectName)
	if err != nil {
		utils.InternalServerError(c, "删除文件失败")
		return
	}

	utils.SuccessWithMessage(c, "文件删除成功")
}

// ListFiles 列出文件
func ListFiles(c *gin.Context) {
	prefix := c.Query("prefix")
	recursive := c.Query("recursive") == "true"

	// 列出文件
	uploader := utils.NewMinIOUploader()
	files, err := uploader.ListFiles(prefix, recursive)
	if err != nil {
		utils.InternalServerError(c, "列出文件失败")
		return
	}

	// 转换为简化的文件信息
	var fileList []gin.H
	for _, file := range files {
		fileList = append(fileList, gin.H{
			"name":          file.Key,
			"size":          file.Size,
			"last_modified": file.LastModified,
			"content_type":  file.ContentType,
		})
	}

	utils.Success(c, gin.H{
		"files": fileList,
		"total": len(fileList),
	}, "获取文件列表成功")
}

// GetFileInfo 获取文件信息
func GetFileInfo(c *gin.Context) {
	objectName := c.Query("object_name")
	if objectName == "" {
		utils.BadRequest(c, "对象名称不能为空")
		return
	}

	// 获取文件信息
	uploader := utils.NewMinIOUploader()
	fileInfo, err := uploader.GetFileInfo(objectName)
	if err != nil {
		utils.NotFound(c, "文件不存在")
		return
	}

	utils.Success(c, gin.H{
		"name":          fileInfo.Key,
		"size":          fileInfo.Size,
		"content_type":  fileInfo.ContentType,
		"last_modified": fileInfo.LastModified,
		"etag":          fileInfo.ETag,
	}, "获取文件信息成功")
}

// GetFileAccessURL 获取文件访问URL
func GetFileAccessURL(c *gin.Context) {
	objectName := c.Query("object_name")
	if objectName == "" {
		utils.BadRequest(c, "对象名称不能为空")
		return
	}

	// 获取有效时间参数
	expireHoursStr := c.Query("expire_hours")
	expireHours := 0
	if expireHoursStr != "" {
		if _, err := fmt.Sscanf(expireHoursStr, "%d", &expireHours); err != nil {
			utils.BadRequest(c, "有效时间参数格式错误")
			return
		}
	}

	// 获取临时访问URL
	uploader := utils.NewMinIOUploader()
	fileURL, err := uploader.GetTemporaryFileURL(objectName, expireHours)
	if err != nil {
		utils.InternalServerError(c, "生成文件访问URL失败")
		return
	}

	// 获取文件信息
	fileInfo, err := uploader.GetFileInfo(objectName)
	if err != nil {
		utils.NotFound(c, "文件不存在")
		return
	}

	utils.Success(c, gin.H{
		"url":           fileURL,
		"object_name":   objectName,
		"size":          fileInfo.Size,
		"content_type":  fileInfo.ContentType,
		"last_modified": fileInfo.LastModified,
		"expire_hours":  expireHours,
	}, "获取文件访问URL成功")
}
