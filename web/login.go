package web

import (
	"fmt"
	"myblog/models"
	"myblog/services"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/go-redis/redis/v8"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type UserLogin struct {
	Username string `json:"username" binding:"required" example:"testuser"`    // 用户名
	Password string `json:"password" binding:"required" example:"password123"` // 密码
}

// Login 用户登录接口
// @Summary 用户登录
// @Description 用户登录获取token
// @Tags 用户管理
// @Accept json
// @Produce json
// @Param user body UserLogin true "用户登录信息"
// @Success 200 {object} map[string]interface{} "登录成功"
// @Failure 1001 {object} map[string]interface{} "参数校验失败"
// @Failure 1002 {object} map[string]interface{} "用户名或密码错误"
// @Router /api/login [post]
func Login(db *gorm.DB, rdb *redis.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		var loginUser UserLogin
		if err := c.ShouldBindJSON(&loginUser); err != nil {
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
		zap.S().Debug("执行用户登录业务逻辑")
		var user = models.User{Username: loginUser.Username, Password: loginUser.Password}
		zap.S().Debug("开始登录用户", user)
		services.LoginUserService(c, db, rdb, user)
		zap.S().Debug("登录用户成功", user)
	}
}
