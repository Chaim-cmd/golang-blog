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

type ArticlesService struct {
	db *gorm.DB
}

func NewArticleService(db *gorm.DB) *ArticlesService {
	return &ArticlesService{db: db}
}

func (s *ArticlesService) Create(userID uint, title, content string) (*model.ArticlesService, error) {
	art := &model.ArticlesService{
		UserID:  userID,
		Title:   title,
		Context: content,
	}
	if err := s.db.Create(art).Error; err != nil {
		return nil, err
	}
	return art, nil
}

func (s *ArticlesService) GetByID(id uint) (*model.ArticlesService, error) {
	var art model.ArticlesService
	if err := s.db.First(&art, id).Error; err != nil {
		return nil, err
	}
	return &art, nil
}

func (s *ArticlesService) List(page, pageSize int) ([]model.ArticlesService, int64, error) {
	var count int64
	var list []model.ArticlesService
	query := s.db.Model(&model.ArticlesService{})
	if err := query.Count(&count).Error; err != nil {
		return nil, 0, err
	}
	if err := s.db.Limit(pageSize).Offset((page - 1) * pageSize).Find(&list).Error; err != nil {
		return nil, 0, err
	}
	return list, count, nil
}

func (s *ArticlesService) Update(id, userID uint, title, content string) error {
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
		s.db.Model(&model.ArticlesService{}).Where("id = ?", id).Count(&cnt)
		if cnt == 0 {
			return ErrArticleNoPerm
		}
		return ErrArticleNotFound
	}
	return nil
}

func (s *ArticlesService) Delete(id, userID uint) error {
	result := s.db.Where("id = ? AND user_id = ?", id, userID).Delete(&model.ArticlesService{})
	rows := result.RowsAffected
	if result.Error != nil {
		return result.Error
	}
	if rows == 0 {
		var cnt int64
		s.db.Model(&model.ArticlesService{}).Where("id = ?", id).Count(&cnt)
		if cnt == 0 {
			return ErrArticleNoPerm
		}
		return ErrArticleNotFound
	}
	return nil

}
