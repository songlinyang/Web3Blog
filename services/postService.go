package services

import (
	"myblog/models"
	Req "myblog/repository"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

// 实现文章的创建功能，只有已认证的用户才能创建文章，创建文章时需要提供文章的标题和内容。
func CreatePostService(db *gorm.DB, post models.Post) error {
	p := Req.PostRep{Db: db}
	zap.S().Debug(p)
	err := p.CreatePost(&post)
	if err != nil {
		return err
	}
	return nil
}

// 实现文章的读取功能，支持获取所有文章列表和单个文章的详细信息。
func ReadOnePostByTitleService(c *gin.Context, db *gorm.DB, post models.Post) (models.Post, error) {
	p := Req.PostRep{Db: db}
	zap.S().Debug(p)
	var postResult models.Post
	err := p.ReadPostByTitle(&post, &postResult)
	if err != nil {
		panic(err)
	}
	return postResult, err

}

func ReadPostListByUserIDService(c *gin.Context, db *gorm.DB, post models.Post) ([]models.Post, error) {
	p := Req.PostRep{Db: db}
	zap.S().Debug(p)
	var postResults []models.Post
	err := p.ReadPostByUserID(&post, &postResults)
	if err != nil {
		panic(err)
	}
	return postResults, err

}

func ReadPostListService(c *gin.Context, db *gorm.DB, post models.Post) {
	p := Req.PostRep{Db: db}
	zap.S().Debug(p)
	var postResults []models.Post
	err := p.ReadAllPostList(&post, &postResults)
	if err != nil {
		c.JSON(200, gin.H{
			"code":  200,
			"msg":   "全部文章列表获取失败",
			"error": err.Error(),
		})
		return
	}
	c.JSON(200, gin.H{
		"code": 200,
		"msg":  "全部文章列表获取成功",
		"data": postResults,
	})

}

// 实现文章的更新功能，只有文章的作者才能更新自己的文章。
func UpdatePostService(db *gorm.DB, post models.Post) {
	p := Req.PostRep{Db: db}
	zap.S().Debug(p)
	err := p.UpdatePostByUserID(&post)
	if err != nil {
		panic(err)
	}

}

// 实现文章的删除功能，只有文章的作者才能删除自己的文章。
func DeletePostService(db *gorm.DB, post models.Post) error {
	p := Req.PostRep{Db: db}
	zap.S().Debug(p)
	//查询是否存在文章
	var postResult []models.Post
	err := p.ReadPostByUserID(&post, &postResult)
	if err != nil {
		return err
	}
	if len(postResult) == 0 {
		panic("当前用户未存在文章")
	}
	//进行删除
	deleteErr := p.DeletePostByUserID(&post)
	if deleteErr != nil {
		return deleteErr
	}
	return nil

}
