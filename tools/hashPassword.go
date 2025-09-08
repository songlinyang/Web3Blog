package tools

import (
	"golang.org/x/crypto/bcrypt"
)

// 密码加密方法
func HashPassword(password string) (string, error) {
	// 加密密码
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		panic("Failed to hash password")
	}
	return string(hashedPassword), err
}
