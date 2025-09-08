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
	PostId  int64  `json:"postId" form:"postId" binding:"number,required" example:"1"` // 文章ID
	Content string `json:"content" example:"这是一条评论"`                                   // 评论内容
}

// CreateCommentByPostIdWeb 创建评论接口
// @Summary 创建评论
// @Description 为文章创建评论
// @Tags 评论管理
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param Authorization header string true "Bearer token"
// @Param comment body CommnetByPostId true "评论信息"
// @Success 200 {object} map[string]interface{} "评论成功"
// @Failure 1001 {object} map[string]interface{} "参数校验失败"
// @Failure 1002 {object} map[string]interface{} "评论失败"
// @Router /api/v1/comment [post]
func CreateCommentByPostIdWeb(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var commentByPostId CommnetByPostId
		if err := c.ShouldBindJSON(&commentByPostId); err != nil {
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
		createCommentErr := services.CreateCommentService(c, db, models.Comment{PostID: uint64(commentByPostId.PostId), Content: commentByPostId.Content})
		if createCommentErr != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"code":   1002,
				"msg":    "评论失败",
				"errors": createCommentErr,
			})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"code": 200,
			"msg":  "评论成功",
		})
	}
}

// QueryCommentByPostIdWeb 查询评论接口
// @Summary 查询评论
// @Description 查询文章下的所有评论
// @Tags 评论管理
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param Authorization header string true "Bearer token"
// @Param postId query int true "文章ID" example(1)
// @Success 200 {object} map[string]interface{} "查询成功"
// @Failure 1001 {object} map[string]interface{} "参数校验失败"
// @Failure 1002 {object} map[string]interface{} "查询失败"
// @Router /api/v1/comment [get]
func QueryCommentByPostIdWeb(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var commentByPostId CommnetByPostId
		if err := c.ShouldBindQuery(&commentByPostId); err != nil {
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
		comments, queryErr := services.QueryCommentByPostIdService(db, models.Comment{PostID: uint64(commentByPostId.PostId)})
		if queryErr != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"code":   1002,
				"msg":    "查询文章对应评论失败",
				"errors": queryErr,
			})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"code": 200,
			"msg":  "查询成功",
			"data": comments,
		})
	}
}
