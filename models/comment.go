package models

import "gorm.io/gorm"

// 评论
type Comment struct {
	gorm.Model        //自动补全ID，CreateAt,UpdateAt,Delete软删除字段
	Content    string `gorm:"not null"`
	UserID     uint   `gorm:"not null"` //与User用户表一对多关系，一个用户下的推文有多条评论
	User       User   `gorm:"foreignkey:UserID"`
	PostID     uint   `gorm:"not null"` //与Post推文表一对多关系，一个推文下有多条评论
	Post       Post   `gorm:"foreignkey:PostID"`
}
