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

// 桶满闲置和刚建桶未满一个令牌周期时，第 6 次获取都需等待一个完整令牌周期。
// 用生产速率 0.5（周期 2s），避免测试机短暂卡顿就补出一个新令牌导致结果抖动
func TestRequestLimiterBurstCapAfterIdle(t *testing.T) {
	tests := []struct {
		name string
		age  time.Duration
	}{
		{name: "idle", age: time.Hour + 100*time.Millisecond},
		{name: "startup", age: 100 * time.Millisecond},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := newRequestLimiter(5, tokenRefillRate)
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
			if got := due.Sub(t0); got < 1900*time.Millisecond {
				t.Fatalf("%s：第 6 次确定到期时刻距首次调用 %v，期望至少约 2s（返回等待 %v）", tt.name, got, w6)
			}
			if w6 < 1900*time.Millisecond {
				t.Fatalf("%s：第 6 次返回等待 %v，期望至少约 2s", tt.name, w6)
			}
		})
	}
}

// 桶耗尽后连续获取，等待时长应按排队位置递增（间隔约 1/refillRate），而不是人人等到同样的时长后一起放行
func TestRequestLimiterReservesSequentially(t *testing.T) {
	l := newRequestLimiter(tokenBucketSize, tokenRefillRate)

	// 耗尽初始令牌，返回值只含 0~200ms 抖动
	for i := 0; i < tokenBucketSize; i++ {
		if w := l.getToken(); w > 250*time.Millisecond {
			t.Fatalf("第 %d 次获取等待 %v，桶未空时不应超过抖动上限", i+1, w)
		}
	}

	// 桶空：第 n 个排队者要等到第 n 个补充令牌，间隔约 2s
	waits := make([]time.Duration, 3)
	for i := 0; i < 3; i++ {
		waits[i] = l.getToken()
	}
	for i, w := range waits {
		want := time.Duration(float64(i+1)/tokenRefillRate) * time.Second
		if w < want-100*time.Millisecond || w > want+250*time.Millisecond {
			t.Fatalf("第 %d 个排队者等待 %v，期望约 %v", i+1, w, want)
		}
	}
}

// 冷却期睡醒后不能绕过令牌桶直接放行
func TestWaitForNextRequestAfterCooldownUsesBucket(t *testing.T) {
	t.Cleanup(func() {
		antispiderMutex.Lock()
		defer antispiderMutex.Unlock()
		antispiderCooldown = false
		consecutiveFailures = 0
		lastRequestTime = time.Time{}
		globalLimiter = newRequestLimiter(tokenBucketSize, tokenRefillRate)
	})

	// 排空全局令牌桶
	for i := 0; i < tokenBucketSize; i++ {
		globalLimiter.getToken()
	}

	antispiderMutex.Lock()
	antispiderCooldown = true
	// 冷却期只剩约 300ms
	lastRequestTime = time.Now().Add(-cooldownTime*time.Second + 300*time.Millisecond)
	antispiderMutex.Unlock()

	start := time.Now()
	waitForNextRequest()
	elapsed := time.Since(start)

	// 修复后：冷却 ~300ms + 令牌桶排队 ~2s（含抖动）；绕过令牌桶的实现约 300ms 即放行
	if elapsed < 1500*time.Millisecond {
		t.Fatalf("冷却结束后仅 %v 即放行，说明绕过了令牌桶", elapsed)
	}
}
