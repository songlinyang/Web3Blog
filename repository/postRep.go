package repository

import (
	"myblog/models"

	"gorm.io/gorm"
)

type PostRep struct {
	Db *gorm.DB
}

func (p *PostRep) CreatePost(post *models.Post) error {
	return p.Db.Create(post).Error
}
