package models

import "gorm.io/gorm"

// 评论模型
type Comment struct {
	gorm.Model        //自动补全ID，CreateAt,UpdateAt,Delete软删除字段
	ID         uint64 `gorm:"primaryKey;autoIncrement"`
	Content    string `gorm:"not null"`
	UserID     uint64 `gorm:"not null"` //与User用户表一对多关系，一个用户下的推文有多条评论
	User       User   `gorm:"foreignkey:UserID"`
	PostID     uint64 `gorm:"not null"` //与Post推文表一对多关系，一个推文下有多条评论
	Post       Post   `gorm:"foreignkey:PostID"`
}
