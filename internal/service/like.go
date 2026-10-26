package service

import (
	"errors"

	"github.com/Chaim-cmd/golang-blog/internal/model"
	"github.com/go-sql-driver/mysql"
	"gorm.io/gorm"
)

var (
	ErrAlreadyLiked = errors.New("already liked")
	ErrNotLiked     = errors.New("no like")
)

type Like struct {
	db      *gorm.DB
	notices *NoticeService
}

func NewLike(db *gorm.DB, notices *NoticeService) *Like {
	return &Like{
		db:      db,
		notices: notices,
	}
}

func (n *Like) Create(userID, ArticleID uint) error {
	var authorID uint
	return n.db.Transaction(func(tx *gorm.DB) error {

		//确认文章在不在
		var article model.Article
		if err := tx.Select("id", "user_id").First(&article, ArticleID).Error; err != nil {
			return err
		}
		authorID = article.UserID

		like := model.Like{
			UserID:    userID,
			ArticleID: ArticleID,
		}
		if err := tx.Create(&like).Error; err != nil {
			// GORM v2 duplicate-key error.
			if errors.Is(err, gorm.ErrDuplicatedKey) {
				return ErrAlreadyLiked
			}

			//兜底
			var mysqlErr *mysql.MySQLError
			if errors.As(err, &mysqlErr) && mysqlErr.Number == 1062 {
				return ErrAlreadyLiked
			}
			return err
		}
		err := tx.Model(&model.Article{}).Where("id = ?", ArticleID).UpdateColumn("like_count", gorm.Expr("like_count + 1")).Error
		if err != nil {
			return err
		}
		n.notices.Enqueue(NoticeTask{
			LikerID:   userID,
			AuthorID:  authorID,
			ArticleID: ArticleID,
		})

		return nil

	})

}

func (n *Like) Delete(userID, ArticleID uint) error {
	return n.db.Transaction(func(tx *gorm.DB) error {
		like := model.Like{
			UserID:    userID,
			ArticleID: ArticleID,
		}
		result := tx.Where("user_id = ? AND article_id = ?", userID, ArticleID).Delete(&like)
		if result.Error != nil {
			return result.Error
		}
		if result.RowsAffected == 0 {
			return ErrNotLiked
		}
		if err := tx.Model(&model.Article{}).Where("id = ? ", ArticleID).UpdateColumn("like_count", gorm.Expr("GREATEST(like_count - 1,0)")).Error; err != nil {
			return err
		}
		return nil
	})
}
