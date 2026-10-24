package model

import "gorm.io/gorm"

type Notice struct {
	gorm.Model
	UserID    uint   `json:"user_id" gorm:"index"`
	ArticleID uint   `json:"article_id" gorm:"index"`
	Content   string `json:"content" gorm:"text"`
}
