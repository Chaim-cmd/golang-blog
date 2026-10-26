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

	err := n.db.Transaction(func(tx *gorm.DB) error {
		var article model.Article
		if err := tx.Select("id", "user_id").First(&article, ArticleID).Error; err != nil {
			return err // 文章不存在：ErrRecordNotFound，事务回滚，什么都不会发生
		}
		authorID = article.UserID

		like := model.Like{
			UserID:    userID,
			ArticleID: ArticleID,
		}
		if err := tx.Create(&like).Error; err != nil {
			// GORM 翻译层把 1062 翻成 ErrDuplicatedKey（前提见下方"整合提示②"）
			if errors.Is(err, gorm.ErrDuplicatedKey) {
				return ErrAlreadyLiked
			}
			// 兜底：驱动层直接匹配 1062
			var mysqlErr *mysql.MySQLError
			if errors.As(err, &mysqlErr) && mysqlErr.Number == 1062 {
				return ErrAlreadyLiked
			}
			return err
		}

		if err := tx.Model(&model.Article{}).Where("id = ?", ArticleID).UpdateColumn("like_count", gorm.Expr("like_count + 1")).Error; err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return err // 事务失败（含重复点赞）：一个通知都不发
	}

	n.notices.Enqueue(NoticeTask{
		LikerID:   userID,
		AuthorID:  authorID,
		ArticleID: ArticleID,
	})
	return nil
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
