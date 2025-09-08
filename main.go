package main

import (
	"myblog/middlewares"
	"myblog/migrate"
	"myblog/myredis"
	"myblog/validators"
	"myblog/web"
	"myblog/zaplogger"
	"os"
	"os/signal"
	"syscall"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	"github.com/go-redis/redis/v8"
	"go.uber.org/zap"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// gin程序入口文件
var (
	//db  *gorm.DB
	rdb *redis.Client
)

// @title Gin Web API
// @version 1.0
// @description RESTful API 文档示例
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:8080
// @BasePath /

// @securityDefinitions.basic BasicAuth
func main() {
	//初始化logger
	//初始化日志模块
	loggerMgr := zaplogger.InitLogger()
	zap.ReplaceGlobals(loggerMgr)
	defer loggerMgr.Sync() // flushes buffer, if any
	logger := loggerMgr.Sugar()
	logger.Debug("START!")

	//初始化
	//自动迁移数据库
	db := migrate.InitMigrate()
	//初始化redis
	rdb := myredis.InitRedis()

	//初始化gin
	r := gin.Default()

	// 配置Swagger UI - 使用自定义配置指向正确的swagger.json
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler, ginSwagger.URL("/swagger.json")))

	// 配置swagger.json文件访问
	r.StaticFile("/swagger.json", "./docs/swagger.json")
	api := r.Group("/api")
	// 注册接口
	api.POST("/register", web.Register(db))
	// 登录接口
	api.POST("/login", web.Login(db, rdb))

	// 路由分版本
	v1 := r.Group("/api/v1")
	// 设置中间件
	v1.Use(
		middlewares.LatencyLogger(),
		//middlewares.CORSMiddleware(),
		middlewares.JWTAuth(rdb),
	)
	{
		//文章Restful API接口
		//新增文章
		v1.GET("/post", web.QueryOnePostByTitleService(db))
		//查询当前用户所有文章
		v1.GET("/post/all", web.QueryPostListByUserId(db))
		//查询当前用户指定文章
		v1.POST("/post", web.PostCreateWeb(db))
		//更新当前用户指定文章
		v1.PUT("/post", web.UpdatePostByUserId(db))
		//删除当前该用户指定文章
		v1.DELETE("/post", web.DeletePostByUserId(db))

		//评论Restful API接口
		//评论功能
		//添加评论
		v1.POST("/comment", web.CreateCommentByPostIdWeb(db))
		//评论查询
		v1.GET("/comment", web.QueryCommentByPostIdWeb(db))
	}
	//注册自定义验证器
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		err := v.RegisterValidation("passwordReg", validators.PasswdValidator)
		if err != nil {
			return
		}
	}
	err := r.Run(":8080")
	if err != nil {
		return
	}

	// 监听程序终止信号，进行资源释放
	go listenSignal()

}

// 监听程序终止信号，对资源进行合理释放
func listenSignal() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM) //接听2，15信号，进行资源的释放
	sign := <-c
	zap.S().Debug("signal received", zap.String("signal", sign.String()))
	if rdb != nil {
		_ = rdb.Close()
	}
	zap.S().Debug("优雅的进行了shutdown，资源已经释放完成")
	os.Exit(0)
}
