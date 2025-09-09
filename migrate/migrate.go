package migrate

import (
	"myblog/internal/mysqldb"

	"go.uber.org/zap"
	"gorm.io/gorm"

	"myblog/models"
	"myblog/zaplogger"
)

func InitMigrate() *gorm.DB {
	//初始化日志模块
	logger := zaplogger.InitLogger()
	logger.Debug("初始化日志模块成功")
	//初始化博客数据库模型，进行ORM映射生成表结构
	//连接数据库
	db, err := mysqldb.InitDB()
	if err != nil {
		logger.Error("error", zap.Error(err))
		panic(err)
		return nil
	}
	logger.Debug("连接数据库成功,开始迁移")
	//连接数据库成功,开始迁移
	err = db.AutoMigrate(&models.User{}, &models.Post{}, &models.Comment{})
	if err != nil {
		logger.Error("迁移数据库error", zap.Error(err))
		panic(err)
		return nil
	}
	logger.Debug("数据迁移成功，迁移表结构：", zap.String("tables", "User,Post,Comment"))
	return db
}
