package routes

import (
	"ai-assistant-backend/controllers"
	"ai-assistant-backend/middleware"

	"github.com/gin-gonic/gin"
)

func SetupFAQRoutes(router *gin.Engine) {
	faq := router.Group("/api/faqs")
	faq.Use(middleware.AuthMiddleware())
	{
		// 问答分类
		faq.GET("/categories", controllers.GetFAQCategories)
		faq.POST("/categories", controllers.CreateFAQCategory)

		// 常见问答
		faq.GET("", controllers.GetFAQs)
		faq.POST("", controllers.CreateFAQ)
		faq.PUT("/:id", controllers.UpdateFAQ)
		faq.DELETE("/:id", controllers.DeleteFAQ)
	}
}
