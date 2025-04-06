package middleware

import (
	"net/http"
	"sync"
	"time"
)

// RateLimiter 实现简单的请求速率限制
type RateLimiter struct {
	mu      sync.Mutex
	clients map[string]*client
	limit   int // 每分钟的请求数
}

// client 跟踪单个客户端的请求
type client struct {
	count       int       // 当前周期的请求计数
	lastRequest time.Time // 上次请求的时间
}

// NewRateLimiter 创建新的速率限制器
func NewRateLimiter(limit int) *RateLimiter {
	return &RateLimiter{
		clients: make(map[string]*client),
		limit:   limit,
	}
}

// Middleware 是速率限制中间件函数
func (rl *RateLimiter) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// 从请求中获取客户端标识（IP地址）
		clientIP := r.RemoteAddr

		// 检查是否超出速率限制
		if rl.isLimited(clientIP) {
			http.Error(w, "请求过于频繁，请稍后再试", http.StatusTooManyRequests)
			return
		}

		// 继续处理请求
		next.ServeHTTP(w, r)
	})
}

// isLimited 检查客户端是否超出速率限制
func (rl *RateLimiter) isLimited(clientIP string) bool {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	now := time.Now()
	c, exists := rl.clients[clientIP]

	// 如果是新客户端或者上次请求是在一分钟前
	if !exists || now.Sub(c.lastRequest) > time.Minute {
		rl.clients[clientIP] = &client{
			count:       1,
			lastRequest: now,
		}
		return false
	}

	// 增加计数并更新上次请求时间
	c.count++
	c.lastRequest = now

	// 检查是否超出限制
	return c.count > rl.limit
}