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
	lastSeen   = make(map[string]time.Time)
)

// 初始化限流组件
func getLimiters(ip string, r rate.Limit, b int) *rate.Limiter {
	limitersMu.Lock()
	defer limitersMu.Unlock()
	l, ok := limiters[ip]
	if !ok {
		l = rate.NewLimiter(r, b)
		limiters[ip] = l
	}

	//每次拿到限流器最后的访问时间
	lastSeen[ip] = time.Now()
	return l

}

// 限流中间件
func RateLimitersMiddleware(r rate.Limit, b int) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		if !getLimiters(ctx.ClientIP(), r, b).Allow() {
			response.Error(ctx, http.StatusTooManyRequests, response.CodeTooManyRequest, "too many Requests")
			ctx.Abort()
			return
		}
		ctx.Next()
	}
}

// 清理
func StartCleanup(interval, ttl time.Duration) (stop func()) {
	ticker := time.NewTicker(interval)
	done := make(chan struct{})
	go func() {
		for {
			select {
			//取出时间，时间到了
			case <-ticker.C:
				limitersMu.Lock()
				for ip, t := range lastSeen {
					if time.Since(t) > ttl {
						delete(limiters, ip)
						delete(lastSeen, ip)
					}
				}
				limitersMu.Unlock()
			case <-done:
				ticker.Stop()
				return
			}
		}
	}()
	return func() { close(done) }
}
