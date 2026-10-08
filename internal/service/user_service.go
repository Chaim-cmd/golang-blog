package service

import (
	"errors"

	"github.com/Chaim-cmd/golang-blog/internal/model"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// 预定义错误：handler 层用于errors.Is 精确判断
var (
	ErrUserExists    = errors.New("用户名或邮箱已被注册")
	ErrBadCredential = errors.New("用户名或密码错误")
)

type UserService struct {
	db *gorm.DB
}

func NewUserService(db *gorm.DB) *UserService {
	return &UserService{db: db}
}

// Register 注册： 查重-> bcrypt 加密 -> 入库
func (s *UserService) Register(username, email, password string) (*model.User, error) {
	var count int64
	s.db.Model(&model.User{}).Where("username = ? OR email = ? ", username, email).Count(&count)
	if count > 0 {
		return nil, ErrUserExists
	}

	//bcrypt 加密  GenerateFromPassword(密码转成本,成本参数)
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	u := &model.User{
		Username:     username,
		Email:        email,
		PasswordHash: string(hash),
		Nickname:     username,
	}
	if err := s.db.Create(u).Error; err != nil {
		return nil, err
	}

	return u, nil

}

// Login 登录：查用户 -> bcrypt 对比明文密码
func (s *UserService) Login(username, password string) (*model.User, error) {
	var u model.User
	err := s.db.Where("username = ?", username).First(&u).Error
	if err != nil {
		return nil, ErrBadCredential
	}
	if err := bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(password)); err != nil {
		return nil, ErrBadCredential
	}

	if u.Status != 1 {
		return nil, ErrBadCredential // 被禁用的账号也用同一句话术，不给额外信息
	}
	return &u, nil

}

func (s *UserService) GetByID(id uint) (*model.User, error) {
	var u model.User
	if err := s.db.First(&u, id).Error; err != nil {
		return nil, err
	}
	return &u, nil
}
