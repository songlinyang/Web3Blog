# MyBlog - Gin Web API 项目

一个基于 Gin 框架构建的 RESTful API 博客系统，包含用户管理、文章管理和评论功能。

## 🚀 运行环境要求

- **Go**: 1.21+ 
- **MySQL**: 5.7+
- **Redis**: 8.0+
- **Gin**: 1.10+
- **Gorm**: 1.30+

## 📦 依赖安装

### 1. 安装 Go 依赖
```bash
# 安装项目依赖
go mod tidy
go mod vendor
```

### 2. 数据库配置
在项目根目录创建 `.env` 文件，配置数据库连接：

```env
# MySQL 配置
MYSQL_ADDR=localhost:3306
MYSQL_USER=root
MYSQL_PASSWORD=your_password
MYSQL_DB=myblogdb

# Redis 配置
REDIS_ADDR=localhost:6379
REDIS_DB=0

# JWT 配置
SIGNATURE_KEY=your_secret_key
JWT_TIME=8 #token有效时长（单位8小时）
```

### 3. 数据库初始化
```bash
# 自动创建数据库表结构
go run main.go
```

## 🏃‍♂️ 启动方式

### 开发环境（热部署）
```bash
# 使用 fresh 进行热部署
fresh

# 或者直接运行
go run ./cmd/main.go
```

### 生产环境
```bash
# 编译项目
go build -o myblog ./cmd/main.go

# 运行编译后的二进制文件
./myblog
```

服务默认运行在 `http://localhost:8080`

## 📚 API 文档 (Swagger UI)

项目集成了 Swagger UI，提供完整的 API 文档。

### 访问 Swagger UI
启动服务后，在浏览器中访问：
```
http://localhost:8080/swagger/index.html
```

### Swagger 接口使用说明

#### 1. 用户注册
- **接口**: `POST /api/register`
- **参数**:
  ```json
  {
    "username": "testuser",
    "password": "password123", 
    "email": "test@example.com"
  }
  ```

#### 2. 用户登录
- **接口**: `POST /api/login`
- **参数**:
  ```json
  {
    "username": "testuser",
    "password": "password123"
  }
  ```
- **响应**: 返回 JWT token，需要在后续请求的 Header 中添加：
  ```
  Authorization: Bearer <your_token>
  ```

#### 3. 文章管理接口（需要认证）
所有文章接口都需要在 Header 中添加 JWT token。

- **创建文章**: `POST /api/v1/post`
- **查询单个文章**: `GET /api/v1/post?title=标题&userId=1`
- **查询用户所有文章**: `GET /api/v1/post/all?userId=1`
- **更新文章**: `PUT /api/v1/post`
- **删除文章**: `DELETE /api/v1/post`

#### 4. 评论管理接口（需要认证）
- **创建评论**: `POST /api/v1/comment`
- **查询文章评论**: `GET /api/v1/comment?postId=1`

## 🏗️ 项目结构

```
myblog/
├── cmd/              # 主程序入口main.go
├── models/           # 数据模型定义
│   ├── user.go
│   ├── post.go
│   └── comment.go
├── repository/       # 数据访问层
│   ├── userRep.go
│   ├── postRep.go
│   └── commentRep.go
├── services/         # 业务逻辑层
│   ├── userService.go
│   ├── postService.go
│   └── commentService.go
├── web/              # HTTP 处理层
│   ├── register.go
│   ├── login.go
│   ├── post.go
│   └── comment.go
├── middlewares/      # 中间件
│   ├── jwt.go        # JWT 认证
│   ├── cors.go       # 跨域处理
│   └── latency.go    # 性能监控
├── migrate/          # 数据库迁移
├── internal          # 内部类
│   ├── myredis/      # Redis 客户端  
│   ├── mysqldb/      # MySQL 连接          
├── tools/            # 工具函数
├── validators/       # 数据验证
└── zaplogger/        # 日志配置
```

## 🔐 认证机制

项目使用 JWT + Redis 进行身份认证：

1. **登录时**生成 JWT token 并存入 Redis
2. **请求时**验证 JWT token 并在 Redis 中校验有效性
3. **Token 过期**或 Redis 中不存在时返回 401 错误

## 📊 接口错误码

| 错误码 | 说明           | HTTP 状态码 |
|--------|----------------|-------------|
| 1001   | 参数校验失败   | 400         |
| 1002   | 认证失败       | 401         |
| 2001   | 数据库错误     | 500         |
| 2002   | 业务逻辑错误   | 400         |

## ⚡ 中间件执行顺序

```
[客户端请求]
↓
[Logger中间件]     → 记录请求开始时间
↓  
[CORS中间件]       → 处理跨域请求
↓
[JWT鉴权]          → 验证访问令牌 (Redis校验)
↓  
[业务处理]         → 核心业务逻辑
↓
[Logger中间件]     ← 记录响应耗时
```

## 🛠️ 开发工具

- **Fresh**: 热重载开发工具
- **Swagger**: API 文档生成
- **Zap**: 高性能日志库
- **GORM**: ORM 框架
- **Redis**: 缓存和会话管理

## 📝 注意事项

1. 确保 MySQL 和 Redis 服务已启动
2. 开发时使用 `.env` 文件配置环境变量
3. 生产环境建议使用环境变量而非 `.env` 文件
4. JWT token 有效期为 8 小时，Redis 中会同步过期

## 🤝 贡献指南

1. Fork 本项目
2. 创建特性分支
3. 提交更改
4. 推送到分支
5. 创建 Pull Request

## 📄 许可证

本项目采用 MIT 许可证 - 查看 [LICENSE](LICENSE) 文件了解详情。

---

如有问题，请提交 Issue 或联系开发团队。
