package services

import (
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/go-resty/resty/v2"
)

// resetRetryTestState 重置反爬状态、限流器与重试参数，测试结束后还原
func resetRetryTestState(t *testing.T) {
	t.Helper()
	antispiderMutex.Lock()
	consecutiveFailures = 0
	antispiderCooldown = false
	lastRequestTime = time.Time{}
	antispiderMutex.Unlock()
	globalLimiter = newRequestLimiter(tokenBucketSize, tokenRefillRate)

	oldMaxRetries, oldBackoff := maxRetries, initialBackoff
	oldCooldown := cooldownTime
	maxRetries, initialBackoff, cooldownTime = 3, time.Millisecond, 0
	t.Cleanup(func() {
		maxRetries, initialBackoff, cooldownTime = oldMaxRetries, oldBackoff, oldCooldown
		antispiderMutex.Lock()
		consecutiveFailures = 0
		antispiderCooldown = false
		lastRequestTime = time.Time{}
		antispiderMutex.Unlock()
		globalLimiter = newRequestLimiter(tokenBucketSize, tokenRefillRate)
	})
}

func newTestService(t *testing.T, handler http.HandlerFunc) *Service {
	t.Helper()
	ts := httptest.NewServer(handler)
	t.Cleanup(ts.Close)
	return &Service{client: resty.New().SetBaseURL(ts.URL)}
}

// #359A：JSON 形式的 403（h.c=0、c=null）必须被识别为 HTTP 403：
// 返回类型化错误、不重试、进入反爬冷却，而不是解出 nil 页面当成功
func TestJSON403TreatedAsStatusError(t *testing.T) {
	resetRetryTestState(t)

	var hits int32
	s := newTestService(t, func(w http.ResponseWriter, r *http.Request) {
		atomic.AddInt32(&hits, 1)
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(http.StatusForbidden)
		w.Write([]byte(`{"h":{"c":0},"c":null}`))
	})

	p, err := s.EbookPages("chapter", "token", 0, 20, 0)
	if err == nil || p != nil {
		t.Fatalf("JSON 403 被当作成功：p=%v err=%v", p, err)
	}
	var se *HTTPStatusError
	if !errors.As(err, &se) || se.Code != http.StatusForbidden {
		t.Fatalf("期望 HTTPStatusError(403)，得到 %v", err)
	}
	if got := atomic.LoadInt32(&hits); got != 1 {
		t.Fatalf("403 不可重试，期望请求 1 次，实际 %d 次", got)
	}
	antispiderMutex.Lock()
	defer antispiderMutex.Unlock()
	if !antispiderCooldown {
		t.Fatal("403 应触发反爬冷却")
	}
}

// #359A 相关：HTTP 200 但 c=null 解出 nil 页面时必须返回错误，且不能先记成功再记失败
func TestNilPageReturnedAsError(t *testing.T) {
	resetRetryTestState(t)

	s := newTestService(t, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.Write([]byte(`{"h":{"c":0},"c":null}`))
	})

	p, err := s.EbookPages("chapter", "token", 0, 20, 0)
	if err == nil || p != nil {
		t.Fatalf("c=null 应返回错误：p=%v err=%v", p, err)
	}
	if !strings.Contains(err.Error(), "页面数据为空") {
		t.Fatalf("错误信息应说明页面数据为空: %v", err)
	}
	antispiderMutex.Lock()
	defer antispiderMutex.Unlock()
	if consecutiveFailures != 3 {
		t.Fatalf("3 次失败应累计计数为 3（旧实现先记成功清零再记失败会停在 1），实际 %d", consecutiveFailures)
	}
}

// #359B：反爬判定不再对错误文本做子串匹配——同一条业务错误
// 不因响应里的时间戳恰好含 "429"/"403" 而被判成反爬
func TestBusinessErrorNotJudgedBySubstring(t *testing.T) {
	resetRetryTestState(t)

	for _, ts := range []string{"1790234293", "1790234300", "1790403496"} {
		antispiderMutex.Lock()
		consecutiveFailures, antispiderCooldown = 0, false
		antispiderMutex.Unlock()

		raw := `{"h":{"c":4000,"e":"user no legal!","s":` + ts + `,"t":0},"c":null}`
		var p *EbookPage
		err := handleJSONParse(strings.NewReader(raw), &p)

		var be *BusinessError
		if !errors.As(err, &be) || be.Code != 4000 {
			t.Fatalf("s=%s：期望 BusinessError(4000)，得到 %v", ts, err)
		}
		recordRequestFailure(err)

		antispiderMutex.Lock()
		cf, cd := consecutiveFailures, antispiderCooldown
		antispiderMutex.Unlock()
		if cf != 0 || cd {
			t.Fatalf("s=%s：业务错误不应计入反爬统计（count=%d cooldown=%v）", ts, cf, cd)
		}
	}
}

// #360：5 个章节都返回 user no legal 时，不应重试、不应触发 60 秒冷却，直接失败
func TestUserNoLegalFailsFastWithoutCooldown(t *testing.T) {
	resetRetryTestState(t)

	var hits int32
	s := newTestService(t, func(w http.ResponseWriter, r *http.Request) {
		atomic.AddInt32(&hits, 1)
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.Write([]byte(`{"h":{"c":4000,"e":"user no legal!","s":1790234300,"t":0},"c":null}`))
	})

	start := time.Now()
	var wg sync.WaitGroup
	for i := 0; i < 5; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			p, err := s.EbookPages(fmt.Sprintf("ch%d", i), "token", 0, 20, 0)
			if err == nil || p != nil {
				t.Errorf("ch%d：期望返回错误", i)
			}
			var be *BusinessError
			if !errors.As(err, &be) {
				t.Errorf("ch%d：期望 BusinessError，得到 %v", i, err)
			}
		}(i)
	}
	wg.Wait()

	if got := atomic.LoadInt32(&hits); got != 5 {
		t.Fatalf("user no legal 不可重试，期望 5 个章节共请求 5 次，实际 %d 次", got)
	}
	if elapsed := time.Since(start); elapsed > 10*time.Second {
		t.Fatalf("不应触发 60 秒冷却，实际耗时 %v", elapsed)
	}
	antispiderMutex.Lock()
	defer antispiderMutex.Unlock()
	if antispiderCooldown {
		t.Fatal("user no legal 不应触发反爬冷却")
	}
}

// HTML 形式的 403 同样返回类型化错误（原 text/html 分支）
func TestHTML403TreatedAsStatusError(t *testing.T) {
	resetRetryTestState(t)

	var hits int32
	s := newTestService(t, func(w http.ResponseWriter, r *http.Request) {
		atomic.AddInt32(&hits, 1)
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.WriteHeader(http.StatusForbidden)
		w.Write([]byte("<h2>403 Forbidden</h2>"))
	})

	_, err := s.EbookPages("chapter", "token", 0, 20, 0)
	var se *HTTPStatusError
	if !errors.As(err, &se) || se.Code != http.StatusForbidden {
		t.Fatalf("期望 HTTPStatusError(403)，得到 %v", err)
	}
	if got := atomic.LoadInt32(&hits); got != 1 {
		t.Fatalf("403 不可重试，期望请求 1 次，实际 %d 次", got)
	}
	antispiderMutex.Lock()
	defer antispiderMutex.Unlock()
	if !antispiderCooldown {
		t.Fatal("HTML 403 同样应触发反爬冷却")
	}
}

// 404 等永久状态码不再重试
func Test404NotRetried(t *testing.T) {
	resetRetryTestState(t)

	var hits int32
	s := newTestService(t, func(w http.ResponseWriter, r *http.Request) {
		atomic.AddInt32(&hits, 1)
		w.WriteHeader(http.StatusNotFound)
	})

	_, err := s.EbookPages("chapter", "token", 0, 20, 0)
	if err == nil {
		t.Fatal("期望返回错误")
	}
	if got := atomic.LoadInt32(&hits); got != 1 {
		t.Fatalf("404 不可重试，期望请求 1 次，实际 %d 次", got)
	}
	antispiderMutex.Lock()
	defer antispiderMutex.Unlock()
	if antispiderCooldown {
		t.Fatal("404 不是反爬信号，不应触发冷却")
	}
}

// 429 限流保留重试并走加长退避
func Test429RetriedWithLongerBackoff(t *testing.T) {
	resetRetryTestState(t)

	var hits int32
	s := newTestService(t, func(w http.ResponseWriter, r *http.Request) {
		atomic.AddInt32(&hits, 1)
		w.WriteHeader(http.StatusTooManyRequests)
	})

	_, err := s.EbookPages("chapter", "token", 0, 20, 0)
	if err == nil {
		t.Fatal("期望返回错误")
	}
	if got := atomic.LoadInt32(&hits); got != int32(maxRetries) {
		t.Fatalf("429 应重试满 %d 次，实际 %d 次", maxRetries, got)
	}
	antispiderMutex.Lock()
	defer antispiderMutex.Unlock()
	if !antispiderCooldown {
		t.Fatal("429 应触发反爬冷却")
	}
}

// 其他业务错误（非 user no legal）保持原有重试行为
func TestOtherBusinessErrorStillRetried(t *testing.T) {
	resetRetryTestState(t)

	var hits int32
	s := newTestService(t, func(w http.ResponseWriter, r *http.Request) {
		atomic.AddInt32(&hits, 1)
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.Write([]byte(`{"h":{"c":9999,"e":"server busy","s":1790234300,"t":0},"c":null}`))
	})

	_, err := s.EbookPages("chapter", "token", 0, 20, 0)
	if err == nil {
		t.Fatal("期望返回错误")
	}
	if got := atomic.LoadInt32(&hits); got != int32(maxRetries) {
		t.Fatalf("普通业务错误应保持重试 %d 次，实际 %d 次", maxRetries, got)
	}
	var be *BusinessError
	if !errors.As(err, &be) || be.Code != 9999 {
		t.Fatalf("最终错误应保留 BusinessError 供上层判断: %v", err)
	}
}
