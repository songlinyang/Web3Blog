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
	Title   string `json:"title" binding:"required" example:"我的第一篇文章"`  // 文章标题
	Content string `json:"content" binding:"required" example:"这是文章内容"` // 文章内容
	UserId  uint64 `json:"userId" binding:"required" example:"1"`       // 用户ID
}
type PostByTitle struct {
	Title  string `form:"title" binding:"required" example:"我的第一篇文章"`   // 文章标题
	UserId uint64 `form:"userId" binding:"number,required" example:"1"` // 用户ID
}
type PostByUserId struct {
	UserId uint64 `form:"userId" binding:"number,required" example:"1"` // 用户ID
}

type UpPostByUserId struct {
	UserId  uint64 `form:"userId" binding:"number,required" example:"1"` // 用户ID
	PostId  uint64 `form:"postId" binding:"number,required" example:"1"` // 文章ID
	Title   string `form:"title" example:"更新后的标题"`                       // 文章标题
	Content string `form:"content" example:"更新后的内容"`                     // 文章内容
}

// PostCreateWeb 创建文章接口
// @Summary 创建文章
// @Description 创建新文章
// @Tags 文章管理
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param Authorization header string true "Bearer token"
// @Param post body CreatePost true "文章信息"
// @Success 200 {object} map[string]interface{} "创建成功"
// @Failure 1001 {object} map[string]interface{} "参数校验失败"
// @Failure 1002 {object} map[string]interface{} "创建失败"
// @Router /api/v1/post [post]
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

// QueryOnePostByTitleService 查询单个文章接口
// @Summary 查询单个文章
// @Description 根据标题查询单个文章
// @Tags 文章管理
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param Authorization header string true "Bearer token"
// @Param title query string true "文章标题" example("我的第一篇文章")
// @Param userId query int true "用户ID" example(1)
// @Success 200 {object} map[string]interface{} "查询成功"
// @Failure 400 {object} map[string]interface{} "参数校验失败"
// @Failure 500 {object} map[string]interface{} "查询失败"
// @Router /api/v1/post [get]
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

// QueryPostListByUserId 查询用户文章列表接口
// @Summary 查询用户文章列表
// @Description 根据用户ID查询所有文章
// @Tags 文章管理
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param Authorization header string true "Bearer token"
// @Param userId query int true "用户ID" example(1)
// @Success 200 {object} map[string]interface{} "查询成功"
// @Failure 400 {object} map[string]interface{} "参数校验失败"
// @Failure 500 {object} map[string]interface{} "查询失败"
// @Router /api/v1/post/all [get]
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

// UpdatePostByUserId 更新文章接口
// @Summary 更新文章
// @Description 根据用户ID和文章ID更新文章
// @Tags 文章管理
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param Authorization header string true "Bearer token"
// @Param post body UpPostByUserId true "文章更新信息"
// @Success 200 {object} map[string]interface{} "更新成功"
// @Failure 400 {object} map[string]interface{} "参数校验失败"
// @Failure 500 {object} map[string]interface{} "更新失败"
// @Router /api/v1/post [put]
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

// DeletePostByUserId 删除文章接口
// @Summary 删除文章
// @Description 根据用户ID和文章ID删除文章
// @Tags 文章管理
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param Authorization header string true "Bearer token"
// @Param post body UpPostByUserId true "文章删除信息"
// @Success 200 {object} map[string]interface{} "删除成功"
// @Failure 400 {object} map[string]interface{} "参数校验失败"
// @Failure 500 {object} map[string]interface{} "删除失败"
// @Router /api/v1/post [delete]
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
