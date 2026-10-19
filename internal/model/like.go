package model

import "gorm.io/gorm"

type Like struct {
	gorm.Model
	UserID    uint `json:"user_id" gorm:"index:idx_user_article,unique" `
	ArticleID uint `json:"article_id" gorm:"index:idx_user_article,unique"`
}
