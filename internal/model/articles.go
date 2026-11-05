package model

import "gorm.io/gorm"

type Article struct {
	gorm.Model
	UserID    uint   `json:"user_id" gorm:"index"`
	Title     string `json:"title" gorm:"type:varchar(128);not null;index,class:FULLTEXT"`
	Content   string `json:"content" gorm:"type:text"`
	Comments  []Comment
	LikeCount int `json:"like_count" gorm:"default:0"`
}
