package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
)

func TestGetLimiter(t *testing.T) {

	type CaseItem struct {
		name     string
		ipA, ipB string
		wantEq   bool
	}
	Cases := []CaseItem{
		{
			name:   "同IP服用",
			ipA:    "1.1.1.1",
			ipB:    "1.1.1.1",
			wantEq: true,
		},
		{
			name:   "不同IP服用",
			ipA:    "1.1.1.1",
			ipB:    "2.2.2.2",
			wantEq: false,
		},
	}

	for _, v := range Cases {
		t.Run(v.name, func(t *testing.T) {

			l1 := getLimiters(v.ipA, rate.Every(200*time.Millisecond), 5)
			l2 := getLimiters(v.ipB, rate.Every(200*time.Millisecond), 5)

			if (l1 == l2) != v.wantEq {
				t.Errorf("两次实例是相同的:got:%v want:%v", l1, l2)
			}
		})

	}
}

func TestRateLimiterMiddleware(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(RateLimitersMiddleware(rate.Every(time.Millisecond), 5)) //1ms快桶
	r.GET("/ping", func(ctx *gin.Context) {
		ctx.Status(http.StatusOK)
	})

	//同个IP 连打12发
	var codes []int
	for i := 0; i < 12; i++ {
		req := httptest.NewRequest(http.MethodGet, "/ping", nil) //客户端
		rec := httptest.NewRecorder()                            //接收响应
		r.ServeHTTP(rec, req)                                    //路由链
		codes = append(codes, rec.Code)
	}
	got429 := codes[5:]
	if len(got429) != 7 || codes[4] != 200 || codes[5] != 429 {
		t.Errorf("code =%v, want [200x5, 429x7]", codes)
	}

	time.Sleep(10 * time.Millisecond) //等令牌桶补满
	req := httptest.NewRequest(http.MethodGet, "/ping", nil)
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Errorf("after sleep :code = %d, want 200", rec.Code)
	}

}

func TestStartCleanup(t *testing.T) {
	testIP := "1.1.1.1"
	ttl := 20 * time.Millisecond
	interval := 10 * time.Millisecond
	limitersMu.Lock()
	limiters = make(map[string]*rate.Limiter)
	lastSeen = make(map[string]time.Time)
	limitersMu.Unlock()

	_ = getLimiters(testIP, rate.Limit(10), 20)

	limitersMu.Lock()
	_, okLim := limiters[testIP]
	_, okTs := lastSeen[testIP]
	limitersMu.Unlock()
	if !okLim || !okTs {
		t.Fatalf("预期插入后ip %s 应当存在，但是不存在", testIP)
	}

	// 启动清理协程
	stop := StartCleanup(interval, ttl)
	t.Cleanup(func() {
		stop()
	})

	//保证时间过期 + 至少两轮扫描
	waitDur := ttl + 2*interval
	time.Sleep(waitDur)

	limitersMu.Lock()
	_, okLimAfter := limiters[testIP]
	_, okTsAfter := lastSeen[testIP]
	limitersMu.Unlock()
	if okLimAfter || okTsAfter {
		t.Errorf("清理未生效：ip %s 在超时后仍然保留在map中", testIP)
	}
}
