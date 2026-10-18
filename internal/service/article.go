package service

import (
	"errors"

	"github.com/Chaim-cmd/golang-blog/internal/model"
	"gorm.io/gorm"
)

var (
	ErrArticleNotFound = errors.New("article not found")
	ErrArticleNoPerm   = errors.New("no permission for this article")
)

type Articles struct {
	db *gorm.DB
}

func NewArticleService(db *gorm.DB) *Articles {
	return &Articles{db: db}
}

func (s *Articles) Create(userID uint, title, content string) (*model.Article, error) {
	art := &model.Article{
		UserID:  userID,
		Title:   title,
		Context: content,
	}
	if err := s.db.Create(art).Error; err != nil {
		return nil, err
	}
	return art, nil
}

func (s *Articles) GetByID(id uint) (*model.Article, error) {
	var art model.Article
	if err := s.db.First(&art, id).Error; err != nil {
		return nil, err
	}
	if err := s.db.Preload("Comments").First(&art, id).Error; err != nil {
		return nil, err
	}
	return &art, nil
}

func (s *Articles) List(page, pageSize int) ([]model.Article, int64, error) {
	var count int64
	var list []model.Article
	query := s.db.Model(&model.Article{})
	if err := query.Count(&count).Error; err != nil {
		return nil, 0, err
	}
	if err := s.db.Limit(pageSize).Offset((page - 1) * pageSize).Find(&list).Error; err != nil {
		return nil, 0, err
	}
	return list, count, nil
}

func (s *Articles) Update(id, userID uint, title, content string) error {
	result := s.db.Where("id = ? AND user_id = ?", id, userID).Updates(map[string]any{
		"title":   title,
		"content": content,
	})
	rows := result.RowsAffected
	if result.Error != nil {
		return result.Error
	}
	if rows == 0 {
		var cnt int64
		s.db.Model(&model.Article{}).Where("id = ?", id).Count(&cnt)
		if cnt > 0 {
			return ErrArticleNoPerm
		}
		return ErrArticleNotFound
	}
	return nil
}

func (s *Articles) Delete(id, userID uint) error {
	result := s.db.Where("id = ? AND user_id = ?", id, userID).Delete(&model.Article{})
	rows := result.RowsAffected
	return s.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("id = ? AND user_id = ?", id, userID).Error; err != nil {
			if result.Error != nil {
				return result.Error
			}
			if rows == 0 {
				return ErrArticleNotFound
			}
			return tx.Where("article_id = ?", id).Delete(&model.Comment{}).Error
		}
		return nil
	})

}
