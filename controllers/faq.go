package controllers

import (
	"net/http"
	"strconv"

	"ai-assistant-backend/config"
	"ai-assistant-backend/models"

	"github.com/gin-gonic/gin"
)

type CreateFAQCategoryRequest struct {
	Name string `json:"name" binding:"required"`
}

type CreateFAQRequest struct {
	Question   string `json:"question" binding:"required"`
	Answer     string `json:"answer" binding:"required"`
	CategoryID uint   `json:"category_id"`
}

type UpdateFAQRequest struct {
	Question   string `json:"question"`
	Answer     string `json:"answer"`
	CategoryID uint   `json:"category_id"`
}

// GetFAQCategories 获取问答分类列表
func GetFAQCategories(c *gin.Context) {
	agentID := c.Query("agent_id")
	if agentID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "缺少智能体ID参数",
		})
		return
	}

	agentIDUint, err := strconv.ParseUint(agentID, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "无效的智能体ID",
		})
		return
	}

	var categories []models.FAQCategory
	if err := config.DB.Where("agent_id = ?", agentIDUint).Order("sort").Find(&categories).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "获取问答分类失败",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "获取成功",
		"data":    categories,
	})
}

// CreateFAQCategory 创建问答分类
func CreateFAQCategory(c *gin.Context) {
	var req CreateFAQCategoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "请求参数错误",
			"error":   err.Error(),
		})
		return
	}

	agentID := c.Query("agent_id")
	if agentID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "缺少智能体ID参数",
		})
		return
	}

	agentIDUint, err := strconv.ParseUint(agentID, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "无效的智能体ID",
		})
		return
	}

	category := models.FAQCategory{
		Name:    req.Name,
		AgentID: uint(agentIDUint),
		Sort:    0,
	}

	if err := config.DB.Create(&category).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "创建问答分类失败",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "创建成功",
		"data":    category,
	})
}

// GetFAQs 获取常见问答列表
func GetFAQs(c *gin.Context) {
	agentID := c.Query("agent_id")
	categoryID := c.Query("category_id")

	if agentID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "缺少智能体ID参数",
		})
		return
	}

	agentIDUint, err := strconv.ParseUint(agentID, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "无效的智能体ID",
		})
		return
	}

	query := config.DB.Where("agent_id = ?", agentIDUint)
	if categoryID != "" {
		categoryIDUint, err := strconv.ParseUint(categoryID, 10, 32)
		if err == nil {
			query = query.Where("category_id = ?", categoryIDUint)
		}
	}

	var faqs []models.FAQ
	if err := query.Order("created_at desc").Find(&faqs).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "获取常见问答失败",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "获取成功",
		"data":    faqs,
	})
}

// CreateFAQ 创建常见问答
func CreateFAQ(c *gin.Context) {
	var req CreateFAQRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "请求参数错误",
			"error":   err.Error(),
		})
		return
	}

	// 验证问题长度
	if len(req.Question) > 100 {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "问题长度不能超过100个字符",
		})
		return
	}

	// 验证回答长度
	if len(req.Answer) > 1000 {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "回答长度不能超过1000个字符",
		})
		return
	}

	agentID := c.Query("agent_id")
	if agentID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "缺少智能体ID参数",
		})
		return
	}

	agentIDUint, err := strconv.ParseUint(agentID, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "无效的智能体ID",
		})
		return
	}

	faq := models.FAQ{
		AgentID:    uint(agentIDUint),
		CategoryID: req.CategoryID,
		Question:   req.Question,
		Answer:     req.Answer,
	}

	if err := config.DB.Create(&faq).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "创建常见问答失败",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "创建成功",
		"data":    faq,
	})
}

// UpdateFAQ 更新常见问答
func UpdateFAQ(c *gin.Context) {
	id := c.Param("id")
	faqID, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "无效的问答ID",
		})
		return
	}

	var req UpdateFAQRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "请求参数错误",
			"error":   err.Error(),
		})
		return
	}

	var faq models.FAQ
	if err := config.DB.First(&faq, faqID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"code":    404,
			"message": "常见问答不存在",
		})
		return
	}

	// 验证问题长度
	if req.Question != "" && len(req.Question) > 100 {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "问题长度不能超过100个字符",
		})
		return
	}

	// 验证回答长度
	if req.Answer != "" && len(req.Answer) > 1000 {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "回答长度不能超过1000个字符",
		})
		return
	}

	// 更新问答信息
	updates := make(map[string]interface{})
	if req.Question != "" {
		updates["question"] = req.Question
	}
	if req.Answer != "" {
		updates["answer"] = req.Answer
	}
	if req.CategoryID != 0 {
		updates["category_id"] = req.CategoryID
	}

	if err := config.DB.Model(&faq).Updates(updates).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "更新常见问答失败",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "更新成功",
		"data":    faq,
	})
}

// DeleteFAQ 删除常见问答
func DeleteFAQ(c *gin.Context) {
	id := c.Param("id")
	faqID, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "无效的问答ID",
		})
		return
	}

	var faq models.FAQ
	if err := config.DB.First(&faq, faqID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"code":    404,
			"message": "常见问答不存在",
		})
		return
	}

	if err := config.DB.Delete(&faq).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "删除常见问答失败",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "删除成功",
	})
}
