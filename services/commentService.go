package services

import (
	"myblog/models"
	Req "myblog/repository"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// 创建评论
func CreateCommentService(ctx *gin.Context, db *gorm.DB, comment models.Comment) error {
	c := Req.CommentRep{Db: db}
	u := Req.UserRep{Db: db}
	//判断当前评论用户是否认证
	userId, ok := ctx.Get("userId")
	if !ok {
		panic("当前用户无效")
	}
	//查询UserId是否存在
	queryUserErr := u.SelectUserById(userId.(uint64))
	if queryUserErr != nil {
		return queryUserErr
	}
	err := c.CreateComment(&models.Comment{UserID: userId.(uint64), PostID: comment.PostID, Content: comment.Content})
	if err != nil {
		return err
	}
	return nil
}

// 查询当前文章下的所有评论
func QueryCommentByPostIdService(db *gorm.DB, comment models.Comment) (commentReuslt models.Comment, err error) {
	c := Req.CommentRep{Db: db}
	p := Req.PostRep{Db: db}
	var comments []models.Comment
	//判断是否存在该文章
	var posts []models.Post
	postErr := p.ReadPostByUserID(&models.Post{ID: comment.PostID}, &posts)
	if postErr != nil {
		panic(postErr)
	}
	err = c.QueryCommentByPostId(&comment, &comments)
	if err != nil {
		return models.Comment{}, err
	}
	return commentReuslt, nil
}
