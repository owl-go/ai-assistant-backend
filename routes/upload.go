package routes

import (
	"ai-assistant-backend/controllers"
	"ai-assistant-backend/middleware"

	"github.com/gin-gonic/gin"
)

func SetupUploadRoutes(router *gin.Engine) {
	upload := router.Group("/api/upload")
	upload.Use(middleware.AuthMiddleware())
	{
		upload.POST("/file", controllers.UploadFile)                           // 直接上传文件
		upload.GET("/presigned-upload", controllers.GetPresignedUploadURL)     // 获取预签名上传URL
		upload.GET("/presigned-download", controllers.GetPresignedDownloadURL) // 获取预签名下载URL
		upload.GET("/file-access-url", controllers.GetFileAccessURL)           // 获取文件访问URL
		upload.DELETE("/file", controllers.DeleteFile)                         // 删除文件
		upload.GET("/files", controllers.ListFiles)                            // 列出文件
		upload.GET("/file-info", controllers.GetFileInfo)                      // 获取文件信息
	}
}
