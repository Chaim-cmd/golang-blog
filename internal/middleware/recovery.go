package middleware

import (
	"net/http"

	"github.com/Chaim-cmd/golang-blog/internal/response"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// RecoveryMiddleware 用 zap 记录panic,替代 gin默认的Recovery
func RecoveryMiddleware(log *zap.Logger) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				// panic 抛出来的是interface{},用zap.Any
				log.Error("panic recovered",
					zap.String("path", ctx.Request.URL.Path),
					zap.String("method", ctx.Request.Method),
					zap.Any("panic", err),
				)
				ctx.AbortWithStatusJSON(http.StatusInternalServerError, response.Response{
					Code:    response.CodeServerError,
					Message: "服务器内部错误",
				})
			}
		}()
		ctx.Next()
	}
}
