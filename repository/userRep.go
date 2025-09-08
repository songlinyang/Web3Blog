package repository

import (
	"myblog/models"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

// 实现User用户账号的CRUD
type UserRep struct {
	Db *gorm.DB
}

// 注册成功，插入数据
func (u *UserRep) CreateUser(user *models.User) error {
	zap.S().Debug("开始将模型数据插入数据", user)
	err := u.Db.Debug().Create(user).Error
	if err == nil {
		zap.S().Debug("插入数据成功")
	}
	return err

}

// 根据用户ID查询数据
func (u *UserRep) SelectUserById(id uint64) error {
	err := u.Db.Debug().First(&models.User{}, id).Error
	if err == nil {
		zap.S().Debug("当前用户存在，有效")
	}
	return err
}

// 根据用户名查询数据
func (u *UserRep) SelectUserByName(username string) error {
	zap.S().Debug("根据用户名查询数据,查询用户名：", zap.String("username", username))
	var user models.User
	err := u.Db.Debug().Select("id").Where("username = ?", username).First(&user)
	if err == nil {
		zap.S().Debug("注册成功,查询结果：", user)
	}
	return err.Error
}

// 根据邮箱地址查询数据
func (u *UserRep) SelectUserByEmail(email string) error {
	zap.S().Debug("根据邮箱地址查询数据：", zap.String("email", email))
	var user models.User
	err := u.Db.Debug().Select("id").Where("email = ?", email).First(&user)
	if err == nil {
		zap.S().Debug("注册成功,查询结果：", user)
	}
	return err.Error
}

// 根据用户名查询数据并返回数据模型
func (u *UserRep) SelectUserByNameScanValue(user, userResult *models.User) error {
	zap.S().Debug("根据用户名查询数据,查询用户名：", zap.String("username", user.Username))
	err := u.Db.Debug().Select("id", "username", "password").Where("username = ?", user.Username).First(&models.User{}).Scan(userResult)
	if err == nil {
		zap.S().Debug("注册成功,查询结果：", user)
	}
	return err.Error
}
