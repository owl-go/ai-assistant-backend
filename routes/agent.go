package routes

import (
	"ai-assistant-backend/controllers"
	"ai-assistant-backend/middleware"

	"github.com/gin-gonic/gin"
)

func SetupAgentRoutes(router *gin.Engine) {
	agent := router.Group("/api/agents")
	agent.Use(middleware.AuthMiddleware())
	{
		agent.GET("", controllers.GetAgents)
		agent.GET("/:id", controllers.GetAgent)
		agent.POST("", controllers.CreateAgent)
		agent.PUT("/:id", controllers.UpdateAgent)
		agent.PATCH("/:id/status", controllers.ToggleAgentStatus)
		agent.DELETE("/:id", controllers.DeleteAgent)
	}
}
