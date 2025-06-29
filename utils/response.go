package utils

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Response 统一响应结构
type Response struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
}

// Success 成功响应
func Success(c *gin.Context, data interface{}, message string) {
	c.JSON(http.StatusOK, Response{
		Code:    200,
		Message: message,
		Data:    data,
	})
}

// SuccessWithMessage 成功响应（仅消息）
func SuccessWithMessage(c *gin.Context, message string) {
	c.JSON(http.StatusOK, Response{
		Code:    200,
		Message: message,
	})
}

// Error 错误响应
func Error(c *gin.Context, statusCode int, message string) {
	c.JSON(statusCode, Response{
		Code:    statusCode,
		Message: message,
	})
}

// ErrorWithDetail 错误响应（包含详细错误信息）
func ErrorWithDetail(c *gin.Context, statusCode int, message string, errorDetail string) {
	c.JSON(statusCode, Response{
		Code:    statusCode,
		Message: message,
		Error:   errorDetail,
	})
}

// BadRequest 400 错误响应
func BadRequest(c *gin.Context, message string) {
	Error(c, http.StatusBadRequest, message)
}

// BadRequestWithDetail 400 错误响应（包含详细错误信息）
func BadRequestWithDetail(c *gin.Context, message string, errorDetail string) {
	ErrorWithDetail(c, http.StatusBadRequest, message, errorDetail)
}

// Unauthorized 401 错误响应
func Unauthorized(c *gin.Context, message string) {
	Error(c, http.StatusUnauthorized, message)
}

// Forbidden 403 错误响应
func Forbidden(c *gin.Context, message string) {
	Error(c, http.StatusForbidden, message)
}

// NotFound 404 错误响应
func NotFound(c *gin.Context, message string) {
	Error(c, http.StatusNotFound, message)
}

// InternalServerError 500 错误响应
func InternalServerError(c *gin.Context, message string) {
	Error(c, http.StatusInternalServerError, message)
}

// ValidationError 参数验证错误响应
func ValidationError(c *gin.Context, errorDetail string) {
	BadRequestWithDetail(c, "请求参数错误", errorDetail)
}

// DatabaseError 数据库错误响应
func DatabaseError(c *gin.Context, message string) {
	InternalServerError(c, message)
}

// TokenError 令牌相关错误响应
func TokenError(c *gin.Context, message string) {
	Unauthorized(c, message)
}

// UserNotFound 用户不存在错误响应
func UserNotFound(c *gin.Context) {
	NotFound(c, "用户不存在")
}

// AgentNotFound 智能体不存在错误响应
func AgentNotFound(c *gin.Context) {
	NotFound(c, "智能体不存在")
}

// DocumentNotFound 文档不存在错误响应
func DocumentNotFound(c *gin.Context) {
	NotFound(c, "文档不存在")
}

// FAQNotFound FAQ不存在错误响应
func FAQNotFound(c *gin.Context) {
	NotFound(c, "FAQ不存在")
}

// LoginFailed 登录失败错误响应
func LoginFailed(c *gin.Context) {
	Unauthorized(c, "用户名或密码错误")
}

// UserAlreadyExists 用户已存在错误响应
func UserAlreadyExists(c *gin.Context, field string) {
	BadRequest(c, field+"已存在")
}

// CreateFailed 创建失败错误响应
func CreateFailed(c *gin.Context, resource string) {
	InternalServerError(c, "创建"+resource+"失败")
}

// UpdateFailed 更新失败错误响应
func UpdateFailed(c *gin.Context, resource string) {
	InternalServerError(c, "更新"+resource+"失败")
}

// DeleteFailed 删除失败错误响应
func DeleteFailed(c *gin.Context, resource string) {
	InternalServerError(c, "删除"+resource+"失败")
}

// GetFailed 获取失败错误响应
func GetFailed(c *gin.Context, resource string) {
	InternalServerError(c, "获取"+resource+"失败")
}

// InvalidID 无效ID错误响应
func InvalidID(c *gin.Context, resource string) {
	BadRequest(c, "无效的"+resource+"ID")
}
