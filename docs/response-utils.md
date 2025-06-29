# 响应处理工具使用指南

## 概述

`utils/response.go` 提供了一套统一的API响应处理工具，用于减少控制器中的代码冗余，确保API响应格式的一致性。

## 响应结构

所有API响应都遵循统一的JSON结构：

```json
{
  "code": 200,
  "message": "操作成功",
  "data": {...},
  "error": "详细错误信息（可选）"
}
```

## 成功响应

### `Success(c *gin.Context, data interface{}, message string)`
返回成功响应，包含数据。

```go
utils.Success(c, gin.H{
    "id": 1,
    "name": "示例",
}, "创建成功")
```

### `SuccessWithMessage(c *gin.Context, message string)`
返回成功响应，仅包含消息。

```go
utils.SuccessWithMessage(c, "删除成功")
```

## 错误响应

### 通用错误响应

#### `Error(c *gin.Context, statusCode int, message string)`
返回指定状态码的错误响应。

```go
utils.Error(c, http.StatusBadRequest, "参数错误")
```

#### `ErrorWithDetail(c *gin.Context, statusCode int, message string, errorDetail string)`
返回包含详细错误信息的错误响应。

```go
utils.ErrorWithDetail(c, http.StatusBadRequest, "参数错误", err.Error())
```

### 特定状态码错误响应

#### 400 Bad Request
```go
utils.BadRequest(c, "请求参数错误")
utils.BadRequestWithDetail(c, "请求参数错误", err.Error())
```

#### 401 Unauthorized
```go
utils.Unauthorized(c, "用户未认证")
```

#### 403 Forbidden
```go
utils.Forbidden(c, "权限不足")
```

#### 404 Not Found
```go
utils.NotFound(c, "资源不存在")
```

#### 500 Internal Server Error
```go
utils.InternalServerError(c, "服务器内部错误")
```

## 业务特定错误响应

### 参数验证错误
```go
utils.ValidationError(c, err.Error())
```

### 数据库错误
```go
utils.DatabaseError(c, "数据库操作失败")
```

### 令牌相关错误
```go
utils.TokenError(c, "令牌无效")
```

### 资源不存在错误
```go
utils.UserNotFound(c)        // 用户不存在
utils.AgentNotFound(c)       // 智能体不存在
utils.DocumentNotFound(c)    // 文档不存在
utils.FAQNotFound(c)         // FAQ不存在
```

### 登录相关错误
```go
utils.LoginFailed(c)         // 登录失败
```

### 用户相关错误
```go
utils.UserAlreadyExists(c, "用户名")  // 用户名已存在
utils.UserAlreadyExists(c, "邮箱")    // 邮箱已存在
```

### 操作失败错误
```go
utils.CreateFailed(c, "用户")    // 创建用户失败
utils.UpdateFailed(c, "智能体")  // 更新智能体失败
utils.DeleteFailed(c, "文档")    // 删除文档失败
utils.GetFailed(c, "FAQ列表")    // 获取FAQ列表失败
```

### 其他错误
```go
utils.InvalidID(c, "智能体")     // 无效的智能体ID
```

## 使用示例

### 控制器中的完整示例

```go
package controllers

import (
    "ai-assistant-backend/utils"
    "github.com/gin-gonic/gin"
)

// CreateUser 创建用户
func CreateUser(c *gin.Context) {
    var req CreateUserRequest
    
    // 参数验证
    if err := c.ShouldBindJSON(&req); err != nil {
        utils.ValidationError(c, err.Error())
        return
    }
    
    // 检查用户是否已存在
    if userExists(req.Username) {
        utils.UserAlreadyExists(c, "用户名")
        return
    }
    
    // 创建用户
    user, err := createUser(req)
    if err != nil {
        utils.CreateFailed(c, "用户")
        return
    }
    
    // 返回成功响应
    utils.Success(c, user, "用户创建成功")
}

// GetUser 获取用户信息
func GetUser(c *gin.Context) {
    id := c.Param("id")
    userID, err := strconv.ParseUint(id, 10, 32)
    if err != nil {
        utils.InvalidID(c, "用户")
        return
    }
    
    user, err := getUserByID(uint(userID))
    if err != nil {
        utils.UserNotFound(c)
        return
    }
    
    utils.Success(c, user, "获取成功")
}

// DeleteUser 删除用户
func DeleteUser(c *gin.Context) {
    userID := c.GetUint("user_id")
    
    if err := deleteUser(userID); err != nil {
        utils.DeleteFailed(c, "用户")
        return
    }
    
    utils.SuccessWithMessage(c, "用户删除成功")
}
```

## 最佳实践

1. **统一使用响应工具**：避免在控制器中直接使用 `c.JSON()`，统一使用响应工具函数。

2. **选择合适的错误响应**：根据错误类型选择合适的响应函数，而不是使用通用的 `Error()` 函数。

3. **保持消息一致性**：使用预定义的错误消息，确保整个API的错误消息格式一致。

4. **包含必要的错误详情**：对于参数验证错误，使用 `ValidationError()` 包含详细的错误信息。

5. **避免重复代码**：如果发现多个地方有相似的错误处理逻辑，考虑在响应工具中添加新的专用函数。

## 扩展响应工具

如果需要添加新的响应函数，可以在 `utils/response.go` 中添加：

```go
// CustomError 自定义错误响应
func CustomError(c *gin.Context, message string) {
    Error(c, http.StatusBadRequest, message)
}

// ResourceConflict 资源冲突错误
func ResourceConflict(c *gin.Context, resource string) {
    Error(c, http.StatusConflict, resource+"已存在")
}
```

这样可以确保整个项目的API响应格式保持一致，减少代码冗余，提高代码的可维护性。 