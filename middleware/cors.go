package middleware

import (
	"net/http"
	"strings"
)

// CorsMiddleware 处理CORS
type CorsMiddleware struct {
	allowedOrigins []string
}

// NewCorsMiddleware 创建新的CORS中间件
func NewCorsMiddleware(allowedOrigins []string) *CorsMiddleware {
	return &CorsMiddleware{
		allowedOrigins: allowedOrigins,
	}
}

// Middleware 是CORS中间件函数
func (m *CorsMiddleware) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		origin := r.Header.Get("Origin")

		// 检查来源是否在允许列表中
		allowed := false
		if origin != "" {
			for _, allowedOrigin := range m.allowedOrigins {
				if allowedOrigin == "*" || allowedOrigin == origin {
					allowed = true
					break
				}
			}
		}

		// 如果来源允许，设置CORS头
		if allowed {
			w.Header().Set("Access-Control-Allow-Origin", origin)
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
			w.Header().Set("Access-Control-Max-Age", "86400") // 24小时
		}

		// 处理预检请求
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}

		// 调用下一个处理程序
		next.ServeHTTP(w, r)
	})
}