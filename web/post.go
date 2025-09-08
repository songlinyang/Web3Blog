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

// Post对象参数值
type CreatePost struct {
	Title   string `json:"title" binding:"required"`
	Content string `json:"content" binding:"required"`
	UserId  uint64 `json:"userId" binding:"required"`
}
type PostByTitle struct {
	Title  string `form:"title" binding:"required"`
	UserId uint64 `form:"userId" binding:"number,required"`
}
type PostByUserId struct {
	UserId uint64 `form:"userId" binding:"number,required"`
}

type UpPostByUserId struct {
	UserId  uint64 `form:"userId" binding:"number,required"`
	PostId  uint64 `form:"postId" binding:"number,required"`
	Title   string `form:"title"`
	Content string `form:"content"`
}

// 新增文章
func PostCreateWeb(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var createPost CreatePost
		err := c.ShouldBindJSON(&createPost)
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
		zap.S().Debug("开始新增文章，文章Title: ", createPost.Title)
		err = services.CreatePostService(db, models.Post{Title: createPost.Title, Content: createPost.Content, UserID: createPost.UserId})
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"code":   1002,
				"msg":    "新增文章失败",
				"errors": err,
			})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"code":   200,
			"userId": createPost.UserId,
			"msg":    "新增文章成功",
		})
	}
}

// 查询单个文章

func QueryOnePostByTitleService(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var postByTitle PostByTitle
		err := c.ShouldBindQuery(&postByTitle)
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
		postResult, readError := services.ReadOnePostByTitleService(c, db, models.Post{Title: postByTitle.Title, UserID: postByTitle.UserId})
		if readError != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"code":   1002,
				"msg":    "查询单个文章失败",
				"errors": err,
			})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"code":    200,
			"msg":     "查询成功",
			"title":   postResult.Title,
			"Content": postResult.Content,
		})
	}
}

// 根据用户id查询多篇文章
func QueryPostListByUserId(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var postByUserId PostByUserId
		err := c.ShouldBindQuery(&postByUserId)
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
		postReuslts, postError := services.ReadPostListByUserIDService(c, db, models.Post{UserID: postByUserId.UserId})
		if postError != nil {
			c.JSON(200, gin.H{
				"code":   200,
				"userID": postByUserId.UserId,
				"msg":    "文章列表获取失败",
				"error":  postError.Error(),
			})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"code":   200,
			"msg":    "查询成功",
			"userId": postByUserId.UserId,
			"data":   postReuslts,
		})
	}
}

// 根据用户ID更新文章
func UpdatePostByUserId(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var upPostByUserId UpPostByUserId
		err := c.ShouldBindJSON(&upPostByUserId)
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
		services.UpdatePostService(db, models.Post{UserID: upPostByUserId.UserId, ID: upPostByUserId.PostId, Title: upPostByUserId.Title, Content: upPostByUserId.Content})
		c.JSON(http.StatusOK, gin.H{
			"code":   200,
			"msg":    "更新成功",
			"userId": upPostByUserId.UserId,
		})
	}
}

// 根据用户删除数据
func DeletePostByUserId(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var deleteUser UpPostByUserId
		err := c.ShouldBindJSON(&deleteUser)
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
		zap.S().Debug("开始删除数据，通过用户ID：", deleteUser.UserId)
		deleteErr := services.DeletePostService(db, models.Post{UserID: deleteUser.UserId, ID: deleteUser.PostId})
		if deleteErr != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"code":   1002,
				"msg":    "删除失败",
				"userId": deleteUser.UserId,
				"postId": deleteUser.PostId,
				"errors": deleteErr,
			})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"code":   200,
			"msg":    "删除成功",
			"userId": deleteUser.UserId,
			"postId": deleteUser.PostId,
		})
		zap.S().Debug("删除成功用户ID,文章ID：", deleteUser.UserId, deleteUser.PostId)

	}
}
