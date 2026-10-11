package middleware

import (
	"net/http"
	"strings"

	"github.com/Chaim-cmd/golang-blog/internal/response"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func AuthMiddleware(secret []byte) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		//从header 提取token
		authHeader := ctx.GetHeader("Authorization")
		if authHeader == "" {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, response.Response{
				Code:    response.CodeUnauthorized,
				Message: "缺少 Authorization Header",
			})
			return
		}
		tokenStr := strings.TrimPrefix(authHeader, "Bearer ")
		if tokenStr == authHeader {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, response.Response{
				Code:    response.CodeUnauthorized,
				Message: "token 格式错误, 应以 Bearer 开头",
			})
			return
		}

		//Parse签发： 验签+自动校验
		token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, jwt.ErrSignatureInvalid
			}
			return secret, nil
		})
		if err != nil || !token.Valid {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, response.Response{
				Code:    response.CodeUnauthorized,
				Message: "token 过期或者无效",
			})
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, response.Response{
				Code:    response.CodeUnauthorized,
				Message: "claims 解析失败",
			})
			return
		}

		//取出 user_id
		userIDFloat, ok := claims["user_id"].(float64)
		if !ok {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, response.Response{
				Code:    response.CodeUnauthorized,
				Message: "token 缺少 user_id",
			})
			return
		}

		//塞进上下文
		ctx.Set("user_id", uint(userIDFloat))
		ctx.Next()
	}

}
