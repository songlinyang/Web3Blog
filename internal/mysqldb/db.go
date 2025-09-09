package mysqldb

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// 初始化数据库
func InitDB() (*gorm.DB, error) {
	// 加载 .env文件
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	var (
		dnsUrl = os.Getenv("dbHost") + ":" + os.Getenv("dbPort")
		user   = os.Getenv("dbUser")
		pass   = os.Getenv("dbPass")
		dbName = os.Getenv("dbName")
	)
	fmt.Println(dnsUrl, user, pass, dbName)
	var DNS = fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8&parseTime=True&loc=Local", user, pass, dnsUrl, dbName)
	fmt.Println("DNS:", DNS)
	db, err := gorm.Open(mysql.New(mysql.Config{
		DSN:                       DNS,
		DefaultStringSize:         256,   // string 类型字段的默认长度
		DisableDatetimePrecision:  true,  // 禁用 datetime 精度，MySQL 5.6 之前的数据库不支持
		DontSupportRenameIndex:    true,  // 重命名索引时采用删除并新建的方式，MySQL 5.7 之前的数据库和 MariaDB 不支持重命名索引
		DontSupportRenameColumn:   true,  // 用 `change` 重命名列，MySQL 8 之前的数据库和 MariaDB 不支持重命名列
		SkipInitializeWithVersion: false, // 根据当前 MySQL 版本自动配置
	}), &gorm.Config{})
	// 清空数据库连接 待写
	if err != nil {
		panic(err)
	}
	return db, err
}
