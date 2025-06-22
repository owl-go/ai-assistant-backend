# AI智能体后台管理系统 API 文档

## 基础信息

- **基础URL**: `http://localhost:8080`
- **认证方式**: JWT Bearer Token (Redis缓存)
- **请求格式**: JSON
- **响应格式**: JSON

## 通用响应格式

```json
{
  "code": 200,
  "message": "操作成功",
  "data": {}
}
```

## 认证接口

### 用户登录

**POST** `/api/auth/login`

**请求参数:**
```json
{
  "username": "admin",
  "password": "admin123"
}
```

**响应示例:**
```json
{
  "code": 200,
  "message": "登录成功",
  "data": {
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "user": {
      "id": 1,
      "username": "admin",
      "email": "admin@example.com",
      "avatar": ""
    }
  }
}
```

### 用户登出

**POST** `/api/auth/logout`

**请求头:**
```
Authorization: Bearer <token>
```

**响应示例:**
```json
{
  "code": 200,
  "message": "登出成功"
}
```

### 刷新Token

**POST** `/api/auth/refresh`

**请求头:**
```
Authorization: Bearer <token>
```

**响应示例:**
```json
{
  "code": 200,
  "message": "刷新成功",
  "data": {
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
  }
}
```

### 用户注册

**POST** `/api/auth/register`

**请求参数:**
```json
{
  "username": "newuser",
  "password": "password123",
  "email": "user@example.com"
}
```

### 获取用户信息

**GET** `/api/auth/profile`

**请求头:**
```
Authorization: Bearer <token>
```

## 智能体管理接口

### 获取智能体列表

**GET** `/api/agents`

**请求头:**
```
Authorization: Bearer <token>
```

**响应示例:**
```json
{
  "code": 200,
  "message": "获取成功",
  "data": [
    {
      "id": 1,
      "name": "招生咨询助手",
      "logo": "/uploads/2024-01-01/1234567890_logo.png",
      "status": "online",
      "link": "https://example.com/agent/1",
      "welcome_msg": "欢迎使用招生咨询助手",
      "carousel_images": [
        "/uploads/2024-01-01/1234567890_image1.jpg",
        "/uploads/2024-01-01/1234567890_image2.jpg"
      ],
      "created_at": "2024-01-01T10:00:00Z",
      "updated_at": "2024-01-01T10:00:00Z"
    }
  ]
}
```

### 创建智能体

**POST** `/api/agents`

**请求头:**
```
Authorization: Bearer <token>
```

**请求参数:**
```json
{
  "name": "新智能体",
  "logo": "/uploads/2024-01-01/1234567890_logo.png",
  "welcome_msg": "欢迎使用新智能体",
  "carousel_images": [
    "/uploads/2024-01-01/1234567890_image1.jpg",
    "/uploads/2024-01-01/1234567890_image2.jpg"
  ]
}
```

### 更新智能体

**PUT** `/api/agents/:id`

**请求头:**
```
Authorization: Bearer <token>
```

**请求参数:**
```json
{
  "name": "更新后的智能体名称",
  "logo": "/uploads/2024-01-01/1234567890_new_logo.png",
  "welcome_msg": "更新后的欢迎语",
  "carousel_images": [
    "/uploads/2024-01-01/1234567890_new_image1.jpg"
  ]
}
```

### 切换智能体状态

**PATCH** `/api/agents/:id/status`

**请求头:**
```
Authorization: Bearer <token>
```

**响应示例:**
```json
{
  "code": 200,
  "message": "状态更新成功",
  "data": {
    "id": 1,
    "status": "online"
  }
}
```

### 删除智能体

**DELETE** `/api/agents/:id`

**请求头:**
```
Authorization: Bearer <token>
```

## 文档管理接口

### 获取文档分类

**GET** `/api/documents/categories?agent_id=1`

**请求头:**
```
Authorization: Bearer <token>
```

**响应示例:**
```json
{
  "code": 200,
  "message": "获取成功",
  "data": [
    {
      "id": 1,
      "name": "招生计划",
      "agent_id": 1,
      "sort": 0
    }
  ]
}
```

### 创建文档分类

**POST** `/api/documents/categories?agent_id=1`

**请求头:**
```
Authorization: Bearer <token>
```

**请求参数:**
```json
{
  "name": "新分类"
}
```

### 获取文档列表

**GET** `/api/documents?agent_id=1&category_id=1&page=1&page_size=10`

**请求头:**
```
Authorization: Bearer <token>
```

**响应示例:**
```json
{
  "code": 200,
  "message": "获取成功",
  "data": {
    "documents": [
      {
        "id": 1,
        "agent_id": 1,
        "category_id": 1,
        "name": "2024年招生计划.pdf",
        "format": "pdf",
        "size": 1024000,
        "path": "/uploads/2024-01-01/1234567890_document.pdf",
        "upload_time": "2024-01-01T10:00:00Z",
        "created_at": "2024-01-01T10:00:00Z",
        "updated_at": "2024-01-01T10:00:00Z"
      }
    ],
    "total": 1,
    "page": 1,
    "page_size": 10
  }
}
```

### 创建文档

**POST** `/api/documents?agent_id=1`

**请求头:**
```
Authorization: Bearer <token>
```

**请求参数:**
```json
{
  "name": "新文档.pdf",
  "category_id": 1,
  "path": "/uploads/2024-01-01/1234567890_document.pdf",
  "format": "pdf",
  "size": 1024000
}
```

### 删除文档

**DELETE** `/api/documents/:id`

**请求头:**
```
Authorization: Bearer <token>
```

### 获取标签列表

**GET** `/api/documents/tags?agent_id=1`

**请求头:**
```
Authorization: Bearer <token>
```

### 创建标签

**POST** `/api/documents/tags?agent_id=1`

**请求头:**
```
Authorization: Bearer <token>
```

**请求参数:**
```json
{
  "tag_name": "重要文档"
}
```

## 常见问答接口

### 获取问答分类

**GET** `/api/faqs/categories?agent_id=1`

**请求头:**
```
Authorization: Bearer <token>
```

### 创建问答分类

**POST** `/api/faqs/categories?agent_id=1`

**请求头:**
```
Authorization: Bearer <token>
```

**请求参数:**
```json
{
  "name": "常见问题"
}
```

### 获取常见问答列表

**GET** `/api/faqs?agent_id=1&category_id=1`

**请求头:**
```
Authorization: Bearer <token>
```

**响应示例:**
```json
{
  "code": 200,
  "message": "获取成功",
  "data": [
    {
      "id": 1,
      "agent_id": 1,
      "category_id": 1,
      "question": "如何申请入学？",
      "answer": "请按照以下步骤申请入学：1. 准备相关材料 2. 提交申请 3. 等待审核",
      "created_at": "2024-01-01T10:00:00Z",
      "updated_at": "2024-01-01T10:00:00Z"
    }
  ]
}
```

### 创建常见问答

**POST** `/api/faqs?agent_id=1`

**请求头:**
```
Authorization: Bearer <token>
```

**请求参数:**
```json
{
  "question": "新问题",
  "answer": "新回答",
  "category_id": 1
}
```

### 更新常见问答

**PUT** `/api/faqs/:id`

**请求头:**
```
Authorization: Bearer <token>
```

**请求参数:**
```json
{
  "question": "更新后的问题",
  "answer": "更新后的回答",
  "category_id": 1
}
```

### 删除常见问答

**DELETE** `/api/faqs/:id`

**请求头:**
```
Authorization: Bearer <token>
```

## 文件上传接口

### 上传文件

**POST** `/api/upload/file`

**请求头:**
```
Authorization: Bearer <token>
Content-Type: multipart/form-data
```

**请求参数:**
- `file`: 文件（支持 jpg, jpeg, png, gif, pdf, doc, docx, txt）

**响应示例:**
```json
{
  "code": 200,
  "message": "上传成功",
  "data": {
    "filename": "original_name.pdf",
    "size": 1024000,
    "url": "/uploads/2024-01-01/1234567890_original_name.pdf",
    "path": "./uploads/2024-01-01/1234567890_original_name.pdf"
  }
}
```

### 访问文件

**GET** `/uploads/:date/:filename`

**说明**: 直接访问上传的文件，不需要认证

## 错误码说明

| 错误码 | 说明 |
|--------|------|
| 200 | 成功 |
| 400 | 请求参数错误 |
| 401 | 未认证或认证失败 |
| 404 | 资源不存在 |
| 500 | 服务器内部错误 |

## 注意事项

1. 所有需要认证的接口都需要在请求头中携带 `Authorization: Bearer <token>`
2. Token存储在Redis中，支持自动过期和手动失效
3. 支持token刷新功能，可以延长会话时间
4. 文件上传大小限制为10MB
5. 支持的文件类型：jpg, jpeg, png, gif, pdf, doc, docx, txt
6. 问题长度限制：100个字符
7. 回答长度限制：1000个字符
8. 轮播图最多支持4张图片
9. 系统使用YAML配置文件，支持更灵活的配置管理 