package tools

import (
	"os"

	"github.com/golang-jwt/jwt/v5"
)

// 生成JWT令牌
type MyClaims struct {
	Username      string   `json:"username"`
	Roles         []string `json:"roles"`
	Exp           int64    `json:"exp"`
	jwt.MapClaims          //这个是用来继承MapClaims的实现
}

// 生成加密串
func GenerateToken(claims *MyClaims) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(os.Getenv("SIGNATURE_KEY")))
}
