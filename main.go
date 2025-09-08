package main

import (
	"myblog/middlewares"
	"myblog/migrate"
	"myblog/myredis"
	"myblog/web"
	"myblog/zaplogger"
	"os"
	"os/signal"
	"syscall"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

//type Config struct {
//	Db  *gorm.DB
//	Rdb *myredis.Client
//}

//var config *Config
//
//func init() {
//
//
//	config = &Config{Db: db, Rdb: rdb}
//
//}

// gin程序入口文件
var (
	db  *gorm.DB
	rdb *redis.Client
)

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
	api := r.Group("/api")
	// 注册接口
	api.POST("/register", web.Register(db))
	// 登录接口
	api.POST("/login", web.Login(db, rdb))

	// 路由分版本
	v1 := r.Group("/api/v1")
	// 设置全局中间件
	v1.Use(
		middlewares.LatencyLogger(),
		middlewares.CORSMiddleware(),
		middlewares.JWTAuth(),
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
		v1.POST("/comment", web.CreateCommentByPostIdWeb(db))
	}
	err := r.Run(":8080")
	if err != nil {
		return
	}

	//fmt.Print("Hello World")
	////插入用户
	//user := models.User{Username: "admin", Password: "123456", Email: "123456@qq.com"}

	//err := u.CreateUser(&user)
	//if err != nil {
	//	logger.Error(err.Error())
	//}
	//logger.Debug("User created")

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
