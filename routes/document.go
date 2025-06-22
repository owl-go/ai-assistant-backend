package routes

import (
	"ai-assistant-backend/controllers"
	"ai-assistant-backend/middleware"

	"github.com/gin-gonic/gin"
)

func SetupDocumentRoutes(router *gin.Engine) {
	document := router.Group("/api/documents")
	document.Use(middleware.AuthMiddleware())
	{
		// 文档分类
		document.GET("/categories", controllers.GetDocumentCategories)
		document.POST("/categories", controllers.CreateDocumentCategory)

		// 文档管理
		document.GET("", controllers.GetDocuments)
		document.POST("", controllers.CreateDocument)
		document.DELETE("/:id", controllers.DeleteDocument)

		// 标签管理
		document.GET("/tags", controllers.GetTags)
		document.POST("/tags", controllers.CreateTag)
	}
}
