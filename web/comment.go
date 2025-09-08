package web

import (
	"fmt"
	"myblog/models"
	"myblog/services"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"gorm.io/gorm"
)

type CommnetByPostId struct {
	UserId  int64  `json:"userId";binding:"required"`
	PostId  int64  `json:"postId";binding:"required"`
	Content string `json:"content";binding:"required"`
}

func CreateCommentByPostIdWeb(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var commentByPostId CommnetByPostId
		err := c.ShouldBindJSON(&commentByPostId)
		if err != nil {
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
		createCommentErr := services.CreateCommentService(c, db, models.Comment{PostID: uint64(commentByPostId.PostId)})
		if createCommentErr != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"code":   1002,
				"msg":    "评论失败",
				"errors": createCommentErr,
			})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"code":   200,
			"msg":    "评论成功",
			"userId": uint64(commentByPostId.UserId),
		})
	}
}
