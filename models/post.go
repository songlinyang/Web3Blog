package models

import "gorm.io/gorm"

// 博客推文模型
type Post struct {
	gorm.Model //自动补全ID，CreateAt,UpdateAt,Delete软删除字段
	ID         uint64
	Title      string `gorm:"not null"`          //标题
	Content    string `gorm:"not null"`          //内容
	UserID     uint64 `gorm:"not null"`          //与User表一对多关系，一个用户可以发多条推文
	User       User   `gorm:"foreignkey:UserID"` //表明外键为UserID，与User表的ID进行关联
}
