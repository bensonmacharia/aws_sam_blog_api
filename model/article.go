package model

import (
	"bmacharia/aws_sam_blog_api/database"

	"gorm.io/gorm"
)

type Article struct {
	gorm.Model
	ID      uint   `gorm:"primary_key"`
	UserID  uint   `gorm:"not null" json:"user_id"`
	Title   string `gorm:"not null;unique" json:"title"`
	Content string `gorm:"not null" json:"content"`
	User    User   `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"-"`
}

// create an article
func (article *Article) Save() (*Article, error) {
	err := database.Db.Create(&article).Error
	if err != nil {
		return &Article{}, err
	}
	return article, nil
}

// get all articles
func GetArticles(Article *[]Article) (err error) {
	err = database.Db.Find(Article).Error
	if err != nil {
		return err
	}
	return nil
}

// get article by id
func GetArticle(Article *Article, id int) (err error) {
	err = database.Db.Where("id = ?", id).First(Article).Error
	if err != nil {
		return err
	}
	return nil
}

// update article
func UpdateArticle(Article *Article) (err error) {
	database.Db.Save(Article)
	return nil
}
