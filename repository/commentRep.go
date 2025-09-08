package repository

import (
	"myblog/models"

	"gorm.io/gorm"
)

type CommentRep struct {
	Db *gorm.DB
}

func (p *CommentRep) CreateComment(comment *models.Comment) error {
	return p.Db.Create(comment).Error
}
