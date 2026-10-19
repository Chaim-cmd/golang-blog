package service

import (
	"errors"

	"github.com/Chaim-cmd/golang-blog/internal/model"
	"gorm.io/gorm"
)

var (
	ErrAlreadyLiked = errors.New("already liked")
	ErrNotLiked     = errors.New("no like")
)

type Like struct {
	db *gorm.DB
}

func NewLike(db *gorm.DB) *Like {
	return &Like{
		db: db,
	}
}

func (n *Like) Create(userID, ArticleID uint) error {
	return n.db.Transaction(func(tx *gorm.DB) error {
		like := model.Like{
			UserID:    userID,
			ArticleID: ArticleID,
		}
		if err := tx.Create(&like).Error; err != nil {
			if errors.Is(err, ErrAlreadyLiked) {
				return ErrAlreadyLiked
			}
			return err
		}
		err := tx.Model(&model.Article{}).Where("id = ?", ArticleID).UpdateColumn("like_count", gorm.Expr("like_count + 1")).Error
		if err != nil {
			return err
		}
		return nil

	})
}

func (n *Like) Delete(userID, ArticleID uint) error {
	return n.db.Transaction(func(tx *gorm.DB) error {
		like := model.Like{
			UserID:    userID,
			ArticleID: ArticleID,
		}
		if err := tx.Delete(&like).Error; err != nil {
			if errors.Is(err, ErrNotLiked) {
				return ErrNotLiked
			}
			return err
		}
		if err := tx.Model(&model.Article{}).Where("id = ?", ArticleID).UpdateColumn("like_count", gorm.Expr("GREATEST(like_count - 1,0)")).Error; err != nil {
			return err
		}
		return nil
	})
}
