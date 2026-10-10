package service

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type JWTService struct {
	secret      []byte ``
	expireHours int
}

func NewJWTService(secret string, expireHours int) *JWTService {
	return &JWTService{
		secret:      []byte(secret),
		expireHours: expireHours,
	}
}

// GenerateToken登录成功签发 JWT
func (s *JWTService) GenerateToken(userID uint) (string, error) {
	now := time.Now()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": userID,                                                   //业务身份
		"iat":     now.Unix(),                                               //签发时间
		"exp":     now.Add(time.Duration(s.expireHours) * time.Hour).Unix(), //过期时间
	})
	//用secret 对 header + payload做 HMAC-SHA256 签名，产出第三段
	return token.SignedString(s.secret)
}
