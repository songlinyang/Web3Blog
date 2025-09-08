package repository

import "C"
import (
	"myblog/models"

	"gorm.io/gorm"
)

type CommentRep struct {
	Db *gorm.DB
}

// 对指定文章进行评论
func (c *CommentRep) CreateComment(comment *models.Comment) error {
	return c.Db.Where(models.Comment{PostID: comment.PostID}).Create(comment).Error
}

// 根据文章ID获取所有的评论列表
func (c *CommentRep) QueryCommentByPostId(comment *models.Comment, commentResults *[]models.Comment) error {
	return c.Db.Debug().Model(comment).Where("post_id=?", comment.PostID).Find(comment).Scan(commentResults).Error
}
