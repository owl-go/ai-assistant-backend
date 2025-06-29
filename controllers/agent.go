package controllers

import (
	"net/http"
	"strconv"

	"ai-assistant-backend/config"
	"ai-assistant-backend/models"
	"ai-assistant-backend/utils"

	"github.com/gin-gonic/gin"
)

type CreateAgentRequest struct {
	Name           string   `json:"name" binding:"required"`
	Logo           string   `json:"logo"`
	WelcomeMsg     string   `json:"welcome_msg"`
	CarouselImages []string `json:"carousel_images"`
}

type UpdateAgentRequest struct {
	Name           string   `json:"name"`
	Logo           string   `json:"logo"`
	WelcomeMsg     string   `json:"welcome_msg"`
	CarouselImages []string `json:"carousel_images"`
}

// GetAgents 获取智能体列表
func GetAgents(c *gin.Context) {
	userInterface, err := utils.GetUserFromContext(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"code":    401,
			"message": "用户未登录",
		})
		return
	}
	user, ok := userInterface.(utils.Claims)
	userID := user.UserID

	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{
			"code":    401,
			"message": "用户未登录",
		})
		return
	}
	var agents []models.Agent
	if err := config.DB.Where("user_id = ?", userID).Find(&agents).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "获取智能体列表失败",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "获取成功",
		"data":    agents,
	})
}

// GetAgent 获取单个智能体
func GetAgent(c *gin.Context) {
	id := c.Param("id")
	agentID, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "无效的智能体ID",
		})
		return
	}

	var agent models.Agent
	if err := config.DB.First(&agent, agentID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"code":    404,
			"message": "智能体不存在",
		})
		return
	}

	// 加载轮播图
	var carouselImages []models.AgentCarouselImage
	config.DB.Where("agent_id = ?", agent.ID).Order("sort").Find(&carouselImages)

	var imageURLs []string
	for _, img := range carouselImages {
		imageURLs = append(imageURLs, img.ImageURL)
	}
	agent.CarouselImages = imageURLs

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "获取成功",
		"data":    agent,
	})
}

// CreateAgent 创建智能体
func CreateAgent(c *gin.Context) {
	var req CreateAgentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "请求参数错误",
			"error":   err.Error(),
		})
		return
	}

	// 创建智能体
	agent := models.Agent{
		Name:       req.Name,
		Logo:       req.Logo,
		WelcomeMsg: req.WelcomeMsg,
		Status:     "offline",
	}

	if err := config.DB.Create(&agent).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "创建智能体失败",
		})
		return
	}

	// 生成链接
	agent.Link = "https://example.com/agent/" + strconv.FormatUint(uint64(agent.ID), 10)
	config.DB.Model(&agent).Update("link", agent.Link)

	// 保存轮播图
	for i, imageURL := range req.CarouselImages {
		carouselImage := models.AgentCarouselImage{
			AgentID:  agent.ID,
			ImageURL: imageURL,
			Sort:     i + 1,
		}
		config.DB.Create(&carouselImage)
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "创建成功",
		"data":    agent,
	})
}

// UpdateAgent 更新智能体
func UpdateAgent(c *gin.Context) {
	id := c.Param("id")
	agentID, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "无效的智能体ID",
		})
		return
	}

	var req UpdateAgentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "请求参数错误",
			"error":   err.Error(),
		})
		return
	}

	var agent models.Agent
	if err := config.DB.First(&agent, agentID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"code":    404,
			"message": "智能体不存在",
		})
		return
	}

	// 更新智能体信息
	updates := make(map[string]interface{})
	if req.Name != "" {
		updates["name"] = req.Name
	}
	if req.Logo != "" {
		updates["logo"] = req.Logo
	}
	if req.WelcomeMsg != "" {
		updates["welcome_msg"] = req.WelcomeMsg
	}

	if err := config.DB.Model(&agent).Updates(updates).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "更新智能体失败",
		})
		return
	}

	// 更新轮播图
	if req.CarouselImages != nil {
		// 删除旧的轮播图
		config.DB.Where("agent_id = ?", agent.ID).Delete(&models.AgentCarouselImage{})

		// 添加新的轮播图
		for i, imageURL := range req.CarouselImages {
			carouselImage := models.AgentCarouselImage{
				AgentID:  agent.ID,
				ImageURL: imageURL,
				Sort:     i + 1,
			}
			config.DB.Create(&carouselImage)
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "更新成功",
		"data":    agent,
	})
}

// ToggleAgentStatus 切换智能体状态
func ToggleAgentStatus(c *gin.Context) {
	id := c.Param("id")
	agentID, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "无效的智能体ID",
		})
		return
	}

	var agent models.Agent
	if err := config.DB.First(&agent, agentID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"code":    404,
			"message": "智能体不存在",
		})
		return
	}

	// 切换状态
	newStatus := "offline"
	if agent.Status == "offline" {
		newStatus = "online"
	}

	if err := config.DB.Model(&agent).Update("status", newStatus).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "更新状态失败",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "状态更新成功",
		"data": gin.H{
			"id":     agent.ID,
			"status": newStatus,
		},
	})
}

// DeleteAgent 删除智能体
func DeleteAgent(c *gin.Context) {
	id := c.Param("id")
	agentID, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "无效的智能体ID",
		})
		return
	}

	var agent models.Agent
	if err := config.DB.First(&agent, agentID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"code":    404,
			"message": "智能体不存在",
		})
		return
	}

	// 删除相关的轮播图
	config.DB.Where("agent_id = ?", agent.ID).Delete(&models.AgentCarouselImage{})

	// 删除智能体
	if err := config.DB.Delete(&agent).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "删除智能体失败",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "删除成功",
	})
}
