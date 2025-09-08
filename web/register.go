package web

import (
	"fmt"
	"myblog/models"
	"myblog/services"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

// user用户对象参数值
type UserRegister struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
	Email    string `json:"email" binding:"required"`
}

func Register(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var userRegister UserRegister
		if err := c.ShouldBindJSON(&userRegister); err != nil {
			errors := err.(validator.ValidationErrors)
			errorMessages := make([]string, len(errors))
			for i, e := range errors {
				errorMessages[i] = fmt.Sprintf("参数 %s 校验失败：%s", e.Field(), e.Tag())
			}
			c.JSON(http.StatusBadRequest, gin.H{
				"code":   1001,
				"msg":    "参数校验失败",
				"errors": errorMessages,
			})
			return
		}
		var user = models.User{Username: userRegister.Username, Password: userRegister.Password, Email: userRegister.Email}
		zap.S().Debug("开始注册用户", user)
		services.RegisterUserService(c, db, user)
		zap.S().Debug("注册用户成功", user)
	}
}
