package controllers

import (
	"net/http"
	"strconv"
	"time"

	"ai-assistant-backend/config"
	"ai-assistant-backend/models"

	"github.com/gin-gonic/gin"
)

type CreateDocumentCategoryRequest struct {
	Name string `json:"name" binding:"required"`
}

type CreateDocumentRequest struct {
	Name       string `json:"name" binding:"required"`
	CategoryID uint   `json:"category_id"`
	Path       string `json:"path" binding:"required"`
	Format     string `json:"format"`
	Size       int64  `json:"size"`
}

type AddTagRequest struct {
	TagName string `json:"tag_name" binding:"required"`
}

// GetDocumentCategories 获取文档分类列表
func GetDocumentCategories(c *gin.Context) {
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

	var categories []models.DocumentCategory
	if err := config.DB.Where("agent_id = ?", agentIDUint).Order("sort").Find(&categories).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "获取文档分类失败",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "获取成功",
		"data":    categories,
	})
}

// CreateDocumentCategory 创建文档分类
func CreateDocumentCategory(c *gin.Context) {
	var req CreateDocumentCategoryRequest
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

	category := models.DocumentCategory{
		Name:    req.Name,
		AgentID: uint(agentIDUint),
		Sort:    0,
	}

	if err := config.DB.Create(&category).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "创建文档分类失败",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "创建成功",
		"data":    category,
	})
}

// GetDocuments 获取文档列表
func GetDocuments(c *gin.Context) {
	agentID := c.Query("agent_id")
	categoryID := c.Query("category_id")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))

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

	var documents []models.Document
	var total int64

	query.Model(&models.Document{}).Count(&total)
	offset := (page - 1) * pageSize
	if err := query.Offset(offset).Limit(pageSize).Order("created_at desc").Find(&documents).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "获取文档列表失败",
		})
		return
	}

	// 为每个文档加载标签
	for i := range documents {
		var tags []models.DocumentTag
		config.DB.Where("document_id = ?", documents[i].ID).Find(&tags)

		var tagNames []string
		for _, tag := range tags {
			tagNames = append(tagNames, tag.TagName)
		}
		// 这里可以添加一个字段来存储标签，但为了简单起见，我们暂时跳过
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "获取成功",
		"data": gin.H{
			"documents": documents,
			"total":     total,
			"page":      page,
			"page_size": pageSize,
		},
	})
}

// CreateDocument 创建文档
func CreateDocument(c *gin.Context) {
	var req CreateDocumentRequest
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

	document := models.Document{
		AgentID:    uint(agentIDUint),
		CategoryID: req.CategoryID,
		Name:       req.Name,
		Path:       req.Path,
		Format:     req.Format,
		Size:       req.Size,
		UploadTime: time.Now(),
	}

	if err := config.DB.Create(&document).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "创建文档失败",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "创建成功",
		"data":    document,
	})
}

// DeleteDocument 删除文档
func DeleteDocument(c *gin.Context) {
	id := c.Param("id")
	documentID, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "无效的文档ID",
		})
		return
	}

	var document models.Document
	if err := config.DB.First(&document, documentID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"code":    404,
			"message": "文档不存在",
		})
		return
	}

	// 删除相关的标签
	config.DB.Where("document_id = ?", document.ID).Delete(&models.DocumentTag{})

	// 删除文档
	if err := config.DB.Delete(&document).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "删除文档失败",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "删除成功",
	})
}

// GetTags 获取标签列表
func GetTags(c *gin.Context) {
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

	var tags []models.Tag
	if err := config.DB.Where("agent_id = ?", agentIDUint).Find(&tags).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "获取标签列表失败",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "获取成功",
		"data":    tags,
	})
}

// CreateTag 创建标签
func CreateTag(c *gin.Context) {
	var req AddTagRequest
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

	tag := models.Tag{
		Name:    req.TagName,
		AgentID: uint(agentIDUint),
	}

	if err := config.DB.Create(&tag).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "创建标签失败",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "创建成功",
		"data":    tag,
	})
}
