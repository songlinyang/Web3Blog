package repository

import (
	"myblog/models"

	"gorm.io/gorm"
)

type PostRep struct {
	Db *gorm.DB
}

// 创建文章
func (p *PostRep) CreatePost(post *models.Post) error {
	// 插入数据
	return p.Db.Create(post).Error
}

// 读取单篇文章
func (p *PostRep) ReadPostByTitle(post *models.Post, postResult *models.Post) error {
	return p.Db.Where("title=?", post.Title).First(post).Scan(postResult).Error
}

// 读取单篇文章
func (p *PostRep) ReadPostByID(post *models.Post) error {
	return p.Db.Debug().First(post).Error
}

// 读取当前用户下所有的文章 -userID
func (p *PostRep) ReadPostByUserID(post *models.Post, postResult *[]models.Post) error {
	return p.Db.Debug().Model(&models.Post{}).Where("user_id=?", post.UserID).Find(postResult).Scan(postResult).Error
}

// 读取当前用户下指定的文章 -userID
func (p *PostRep) ReadPostByUserIDAndPostID(post *models.Post, postResult *[]models.Post) error {

	return p.Db.Debug().Model(&models.Post{}).Where(models.Post{ID: post.ID, UserID: post.UserID}).Find(postResult).Scan(postResult).Error
}

// 读取博客全部用户的文章
func (p *PostRep) ReadAllPostList(post *models.Post, postResult *[]models.Post) error {
	return p.Db.Find(post).Scan(postResult).Error
}

// 根据用户ID更新文章内容或标题
func (p *PostRep) UpdatePostByUserID(post *models.Post) error {
	return p.Db.Debug().Model(&models.Post{}).Select("title", "post").Where(models.Post{ID: post.ID, UserID: post.UserID}).Updates(post).Error
}

// 根据用户ID删除指定文章
func (p *PostRep) DeletePostByUserID(post *models.Post) error {
	return p.Db.Debug().Where(models.Post{ID: post.ID, UserID: post.UserID}).Delete(post).Error
}
