# AI智能体后台管理系统 - 后端

## 技术栈

- **语言**: Go 1.23+
- **Web框架**: Gin
- **ORM**: GORM
- **数据库**: MySQL
- **认证**: JWT
- **密码加密**: bcrypt

## 项目结构

```
backend/
├── config/          # 配置文件
├── controllers/     # 控制器
├── middleware/      # 中间件
├── models/          # 数据模型
├── routes/          # 路由
├── utils/           # 工具函数
├── main.go          # 主程序入口
├── config.env       # 环境变量配置
└── README.md        # 项目说明
```

## 功能模块

### 1. 用户认证模块
- 用户注册
- 用户登录
- JWT令牌验证
- 获取用户信息

### 2. 智能体管理模块
- 创建智能体
- 获取智能体列表
- 更新智能体信息
- 切换智能体状态（上线/下线）
- 删除智能体
- 轮播图管理

### 3. 文档管理模块
- 文档分类管理
- 文档上传和管理
- 文档标签管理
- 文档分页查询

### 4. 常见问答模块
- 问答分类管理
- 常见问答的增删改查
- 问答内容长度验证

## API接口

### 认证接口
- `POST /api/auth/login` - 用户登录
- `POST /api/auth/register` - 用户注册
- `GET /api/auth/profile` - 获取用户信息

### 智能体接口
- `GET /api/agents` - 获取智能体列表
- `GET /api/agents/:id` - 获取单个智能体
- `POST /api/agents` - 创建智能体
- `PUT /api/agents/:id` - 更新智能体
- `PATCH /api/agents/:id/status` - 切换智能体状态
- `DELETE /api/agents/:id` - 删除智能体

### 文档管理接口
- `GET /api/documents/categories` - 获取文档分类
- `POST /api/documents/categories` - 创建文档分类
- `GET /api/documents` - 获取文档列表
- `POST /api/documents` - 创建文档
- `DELETE /api/documents/:id` - 删除文档
- `GET /api/documents/tags` - 获取标签列表
- `POST /api/documents/tags` - 创建标签

### 常见问答接口
- `GET /api/faqs/categories` - 获取问答分类
- `POST /api/faqs/categories` - 创建问答分类
- `GET /api/faqs` - 获取常见问答列表
- `POST /api/faqs` - 创建常见问答
- `PUT /api/faqs/:id` - 更新常见问答
- `DELETE /api/faqs/:id` - 删除常见问答

## 环境配置

创建 `config.env` 文件并配置以下环境变量：

```env
# 数据库配置
DB_HOST=localhost
DB_PORT=3306
DB_USER=root
DB_PASSWORD=password
DB_NAME=ai_assistant

# JWT配置
JWT_SECRET=your-secret-key
JWT_EXPIRE_HOURS=24

# 服务器配置
PORT=8080
MODE=debug

# 文件上传配置
UPLOAD_PATH=./uploads
MAX_FILE_SIZE=10485760
```

## 安装和运行

1. 安装依赖：
```bash
go mod tidy
```

2. 配置数据库：
- 创建MySQL数据库
- 修改 `config.env` 中的数据库配置

3. 运行项目：
```bash
go run main.go
```

4. 访问健康检查接口：
```
GET http://localhost:8080/health
```

## 数据库表结构

### users - 用户表
- id: 主键
- username: 用户名（唯一）
- password: 密码（加密）
- email: 邮箱（唯一）
- avatar: 头像
- created_at: 创建时间
- updated_at: 更新时间

### agents - 智能体表
- id: 主键
- name: 智能体名称
- logo: Logo图片
- status: 状态（online/offline）
- link: 体验链接
- welcome_msg: 欢迎语
- created_at: 创建时间
- updated_at: 更新时间

### agent_carousel_images - 智能体轮播图表
- id: 主键
- agent_id: 智能体ID
- image_url: 图片URL
- sort: 排序

### documents - 文档表
- id: 主键
- agent_id: 智能体ID
- category_id: 分类ID
- name: 文档名称
- format: 文档格式
- size: 文档大小
- path: 文档路径
- upload_time: 上传时间
- created_at: 创建时间
- updated_at: 更新时间

### faqs - 常见问答表
- id: 主键
- agent_id: 智能体ID
- category_id: 分类ID
- question: 问题
- answer: 回答
- created_at: 创建时间
- updated_at: 更新时间

## 注意事项

1. 所有需要认证的接口都需要在请求头中携带 `Authorization: Bearer <token>` 
2. 文档上传功能需要配合文件上传服务实现
3. 密码使用bcrypt加密存储
4. JWT令牌有效期为24小时
5. 支持CORS跨域请求 