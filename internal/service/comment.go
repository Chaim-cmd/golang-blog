package service

import (
	"github.com/Chaim-cmd/golang-blog/internal/model"
	"gorm.io/gorm"
)

type Comments struct {
	db *gorm.DB
}

func NewComments(db *gorm.DB) *Comments {
	return &Comments{db: db}
}

func (c *Comments) Create(ArticleID, userID uint, content string) (*model.Comment, error) {
	var Article *model.Article
	if err := c.db.First(&Article, ArticleID).Error; err != nil {
		return nil, err
	}
	comment := &model.Comment{
		ArticleID: ArticleID,
		UserID:    userID,
		Content:   content,
	}
	if err := c.db.Create(comment).Error; err != nil {
		return nil, err
	}
	return comment, nil
}
func (c *Comments) ListByArticle(articleID uint, page, pageSize int) ([]model.Comment, int64, error) {
	var count int64
	var comment []model.Comment
	query := c.db.Model(&model.Comment{}).Where("article_id = ?", articleID)
	if err := query.Count(&count).Error; err != nil {
		return nil, 0, err
	}
	if err := c.db.Limit(pageSize).Offset((page - 1) * pageSize).Find(&comment).Error; err != nil {
		return nil, 0, err
	}
	return comment, count, nil

}
