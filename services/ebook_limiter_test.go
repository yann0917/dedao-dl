package services

import (
	"sync"
	"testing"
	"time"
)

// 容量 5、每秒 10 个令牌：15 次并发获取里有 10 次要排队，最后一次放行不应早于约 1 秒
func TestRequestLimiterEnforcesRate(t *testing.T) {
	l := newRequestLimiter(5, 10)
	start := time.Now()

	var wg sync.WaitGroup
	var mu sync.Mutex
	var last time.Time
	for i := 0; i < 15; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			if w := l.getToken(); w > 0 {
				time.Sleep(w)
			}
			mu.Lock()
			if now := time.Now(); now.After(last) {
				last = now
			}
			mu.Unlock()
		}()
	}
	wg.Wait()

	if got := last.Sub(start); got < 900*time.Millisecond {
		t.Fatalf("15 次获取在 %v 内全部放行，限流未生效（期望最后一次不早于约 1s）", got)
	}
}

// 桶满闲置和刚建桶未满一个令牌周期时，第 6 次获取都需等待一个完整令牌周期
func TestRequestLimiterBurstCapAfterIdle(t *testing.T) {
	tests := []struct {
		name string
		age  time.Duration
	}{
		{name: "idle", age: time.Hour + 90*time.Millisecond},
		{name: "startup", age: 90 * time.Millisecond},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := newRequestLimiter(5, 10)
			l.lastRefillTime = time.Now().Add(-tt.age)
			t0 := time.Now()
			var w6 time.Duration
			for i := 0; i < 6; i++ {
				w6 = l.getToken()
			}
			if l.tokens != -1 {
				t.Fatalf("%s：第 6 次后令牌数为 %d，期望 -1", tt.name, l.tokens)
			}
			due := l.lastRefillTime.Add(time.Duration(float64(-l.tokens) / l.refillRate * float64(time.Second)))
			if got := due.Sub(t0); got < 90*time.Millisecond {
				t.Fatalf("%s：第 6 次确定到期时刻距首次调用 %v，期望至少 90ms（返回等待 %v）", tt.name, got, w6)
			}
			if w6 < 90*time.Millisecond {
				t.Fatalf("%s：第 6 次返回等待 %v，期望至少 90ms", tt.name, w6)
			}
		})
	}
}
