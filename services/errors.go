package services

import (
	"errors"
	"fmt"
	"net/http"
	"strings"
	"sync"
)

// HTTPStatusError HTTP 状态码错误。
// 403/429 是反爬/限流信号（常以 JSON body 出现，如 h.c=0 且 c=null，必须在解析内容前拦截）；
// 400/401/404/496 重试也不会恢复。
type HTTPStatusError struct {
	Code int
	URL  string
	// FromHTMLPage 标记错误页来自 text/html 响应（区别于 JSON 形式），仅用于错误信息展示
	FromHTMLPage bool
}

func (e *HTTPStatusError) Error() string {
	msg := fmt.Sprintf("HTTP %d from %s", e.Code, e.URL)
	if e.FromHTMLPage {
		msg += " - Server returned HTML error page"
	}
	return msg
}

// permanent 重试也无法恢复的状态码
func (e *HTTPStatusError) permanent() bool {
	switch e.Code {
	case http.StatusBadRequest, http.StatusUnauthorized, http.StatusForbidden,
		http.StatusNotFound, 496:
		return true
	}
	return false
}

// antiSpider 视为触发了反爬/限流的状态码
func (e *HTTPStatusError) antiSpider() bool {
	return e.Code == http.StatusForbidden || e.Code == http.StatusTooManyRequests
}

// BusinessError 业务错误：服务端返回 h.c != 0（如 user no legal 的 code 4000）。
// 这类错误由账号/权限/风控状态决定，重试与否取决于具体错误。
type BusinessError struct {
	Code int
	Msg  string
	Raw  string // 原始响应体，便于排查
}

func (e *BusinessError) Error() string {
	return fmt.Sprintf("service error: %s (code %d)\nraw response: %s", e.Msg, e.Code, e.Raw)
}

// userNoLegal 账号无权访问：在其他端登录被挤、会员到期、触发风控待网页验证、无本书权限等，重试无法恢复
func (e *BusinessError) userNoLegal() bool {
	return strings.Contains(e.Msg, "user no legal")
}

// isPermanentError 重试无法恢复的错误：永久状态码或账号无权访问
func isPermanentError(err error) bool {
	var se *HTTPStatusError
	if errors.As(err, &se) && se.permanent() {
		return true
	}
	var be *BusinessError
	if errors.As(err, &be) && be.userNoLegal() {
		return true
	}
	return false
}

// isRateLimited 429 限流：保留加长退避后重试
func isRateLimited(err error) bool {
	var se *HTTPStatusError
	return errors.As(err, &se) && se.Code == http.StatusTooManyRequests
}

var userNoLegalHintOnce sync.Once

// hintAccountError user no legal 时提示可能原因（每次进程只提示一次，避免多章节刷屏）
func hintAccountError(err error) {
	var be *BusinessError
	if !errors.As(err, &be) || !be.userNoLegal() {
		return
	}
	userNoLegalHintOnce.Do(func() {
		fmt.Println("账号无权访问（user no legal）：请确认是否在其他端登录过（重新登录）、会员是否到期、是否需要到网页端完成一次验证，或这本书是否在你的权限范围内。")
	})
}
