package model

import "gorm.io/gorm"

type Comment struct {
	gorm.Model
	ArticleID uint   `json:"article_id" gorm:"index"`
	UserID    uint   `json:"user_id" `
	Content   string `json:"content" gorm:"type:text"`
}
