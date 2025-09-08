package models

import "gorm.io/gorm"

// 用户账号
type User struct {
	ID         uint64 `gorm:"primaryKey" json:"id"`
	gorm.Model        //自动补全ID，CreateAt,UpdateAt,Delete软删除字段
	Username   string `gorm:"unique;not null"` //用户名，账号
	Password   string `gorm:"not null"`        //密码
	Email      string `gorm:"unique;not null"` //邮箱
}
