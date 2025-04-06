package middleware

import (
	"net/http"
)

// CORSMiddleware CORS中间件
type CORSMiddleware struct {
	allowedOrigins []string
}

// NewCORSMiddleware 创建CORS中间件
func NewCORSMiddleware(allowedOrigins []string) *CORSMiddleware {
	return &CORSMiddleware{
		allowedOrigins: allowedOrigins,
	}
}

// Middleware 实现中间件
func (m *CORSMiddleware) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		origin := r.Header.Get("Origin")

		// 设置CORS头
		if origin != "" {
			// 检查是否允许该来源
			allowed := false
			if contains(m.allowedOrigins, "*") {
				allowed = true
			} else {
				for _, allowedOrigin := range m.allowedOrigins {
					if allowedOrigin == origin {
						allowed = true
						break
					}
				}
			}

			if allowed {
				w.Header().Set("Access-Control-Allow-Origin", origin)
				w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
				w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
				w.Header().Set("Access-Control-Allow-Credentials", "true")
			}
		}

		// 处理预检请求
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}

// 检查字符串切片是否包含指定字符串
func contains(slice []string, s string) bool {
	for _, item := range slice {
		if item == s {
			return true
		}
	}
	return false
} 