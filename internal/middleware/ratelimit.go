package middleware

import (
	"net/http"
	"sync"
	"time"

	"github.com/Chaim-cmd/golang-blog/internal/response"
	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
)

var (
	limiters   = make(map[string]*rate.Limiter)
	limitersMu sync.Mutex
)

func getLimiters(ip string) *rate.Limiter {
	limitersMu.Lock()
	defer limitersMu.Unlock()
	l, ok := limiters[ip]
	if !ok {
		// 每200ms补一个令牌m,桶的容量为5
		l = rate.NewLimiter(rate.Every(200*time.Millisecond), 5)
		limiters[ip] = l
	}
	return l

}

func RateLimitersMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		if !getLimiters(ctx.ClientIP()).Allow() {
			//TODO:
			response.Error(ctx, http.StatusTooManyRequests, response.CodeForbidden, "to many Requests")
			ctx.Abort()
			return
		}
		ctx.Next()
	}
}
