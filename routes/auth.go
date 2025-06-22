package routes

import (
	"ai-assistant-backend/controllers"
	"ai-assistant-backend/middleware"

	"github.com/gin-gonic/gin"
)

func SetupAuthRoutes(router *gin.Engine) {
	auth := router.Group("/api/auth")
	{
		auth.POST("/login", controllers.Login)
		auth.POST("/register", controllers.Register)
		auth.GET("/profile", middleware.AuthMiddleware(), controllers.GetProfile)
		auth.POST("/logout", middleware.AuthMiddleware(), controllers.Logout)
		auth.POST("/refresh", middleware.AuthMiddleware(), controllers.RefreshToken)
	}
}
