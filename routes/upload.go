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
		upload.POST("/file", controllers.UploadFile)
	}

	// 文件访问路由（不需要认证）
	router.GET("/uploads/:date/:filename", controllers.ServeFile)
}
