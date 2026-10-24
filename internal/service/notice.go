package service

import (
	"errors"
	"fmt"
	"log"

	"github.com/Chaim-cmd/golang-blog/internal/model"
	"gorm.io/gorm"
)

type NoticeTask struct {
	LikerID   uint
	AuthorID  uint
	ArticleID uint
}

type NoticeService struct {
	db   *gorm.DB
	task chan NoticeTask
	done chan struct{}
}

func NewNotices(db *gorm.DB) *NoticeService {
	svc := &NoticeService{
		db:   db,
		task: make(chan NoticeTask, 100),
		done: make(chan struct{}),
	}
	go svc.worker()
	return svc
}

// Enqueue 生产唯一出口
func (s *NoticeService) Enqueue(t NoticeTask) {
	s.task <- t
}

// 优雅退出
func (s *NoticeService) Shutdown() {
	close(s.task)
	<-s.done //阻塞
}

// 常驻消费
func (s *NoticeService) worker() {
	for t := range s.task {
		s.process(t)
	}
	close(s.done)
}

// process 单条处理,失败只记日志
func (s *NoticeService) process(t NoticeTask) {
	if err := s.Save(t); err != nil {
		log.Printf("[notice] save failed liker=%d author=%d article=%d",
			t.LikerID, t.AuthorID, t.ArticleID,
		)

	}
}

func (s *NoticeService) Save(t NoticeTask) error {
	var user model.User
	err := s.db.First(&user, t.LikerID).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil // 点赞者已不存在（刚删号）
	}
	if err != nil {
		return err
	}
	content := fmt.Sprintf("%s 赞了你得文章", user.Username)
	notice := &model.Notice{
		UserID:    t.AuthorID,
		ArticleID: t.ArticleID,
		Content:   content,
	}
	return s.db.Create(&notice).Error
}
