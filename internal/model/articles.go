package model

import "gorm.io/gorm"

type ArticlesService struct {
	gorm.Model
	UserID  uint   `json:"user_id" gorm:"index"`
	Title   string ` json:"title" gorm:"type:varchar(128);not null"`
	Context string `json:"context" gorm:"type:text"`
}
