package model

import "time"

type User struct {
	ID           uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	Username     string    `gorm:"type:varchar(50);uniqueIndex;not null" json:"username"`
	Email        string    `gorm:"type:varchar(100);uniqueIndex;not null" json:"email"`
	PasswordHash string    `gorm:"type:varchar(100);not null" json:"-"`
	Nickname     string    `gorm:"type:varchar(50)" json:"nickname"`
	Avatar       string    `gorm:"type:varchar(255)" json:"avatar"`
	Status       int8      `gorm:"type:tinyint;default:1;comment:1正常 2禁用" json:"status"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

func (User) TableName() string {
	return "users"
}
