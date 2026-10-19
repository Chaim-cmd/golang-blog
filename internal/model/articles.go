package model

import "gorm.io/gorm"

type Article struct {
	gorm.Model
	UserID   uint   `json:"user_id" gorm:"index"`
	Title    string ` json:"title" gorm:"type:varchar(128);not null"`
	Content  string `json:"content" gorm:"type:text"`
	Comments []Comment
}
