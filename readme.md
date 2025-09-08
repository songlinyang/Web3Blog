### 安装所有要求的依赖：
1. __gin__ - Web 框架 (github.com/gin-gonic/gin v1.10.1)
2. __gorm__ - ORM 库 (gorm.io/gorm v1.30.3)
3. __mysql__ - MySQL 驱动 (gorm.io/driver/mysql v1.6.0)
4. __redis__ - Redis 客户端 (github.com/go-redis/redis/v8 v8.11.5)
5. __zap__ - 日志库 (go.uber.org/zap v1.27.0)
6. __jwt__ - JWT 认证 (github.com/golang-jwt/jwt/v5 v5.3.0)
7. __fresh__ - 热加载工具 (github.com/gravityblast/fresh v0.0.0-20240621171608-8d1fef547a99)

### 启动前需验证依赖必要性
```go mod tidy ```

### 目录结构
```azure
myblog/
├── models/ #存放数据模型定义(GORM struct)
│   ├── user.go
│   ├── post.go
│   └── comment.go
├── repositories/ #gorm的数据访问层
│   ├── userRep.go
│   ├── postRep.go
│   └── commentRep.go
├── services/ #业务逻辑层
│   ├── userServices.go
│   ├── postServices.go
│   └── commentServices.go 
├── migrate/
│   └── migrate.go    # 只包含迁移逻辑
└── mysqldb/
└── db.go         # 数据库连接
```






### 接口错误响应码描述
| 错误码 | 说明         |
|--------|--------------|
| 1001   | 参数校验失败 |
| 1002   | 认证失败     |
| 2001   | 数据库错误   |
- 标准化响应格式
- 版本控制方案
- 接口文档生成