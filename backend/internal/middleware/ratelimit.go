package middleware

import (
	"sync"
	"time"

	"astro-pass/internal/utils"
	"github.com/gin-gonic/gin"
)

// RateLimiter 简单的内存速率限制器
type RateLimiter struct {
	visitors map[string]*Visitor
	mu       sync.RWMutex
	rate     int           // 每分钟允许的请求数
	window   time.Duration // 时间窗口
}

// Visitor 访问者信息
type Visitor struct {
	Count    int
	LastSeen time.Time
}

var globalRateLimiter *RateLimiter

func init() {
	globalRateLimiter = &RateLimiter{
		visitors: make(map[string]*Visitor),
		rate:     60, // 默认每分钟60次
		window:   time.Minute,
	}

	// 定期清理过期访问者
	go func() {
		for {
			time.Sleep(5 * time.Minute)
			globalRateLimiter.cleanup()
		}
	}()
}

// RateLimitMiddleware 速率限制中间件
func RateLimitMiddleware(rate int) gin.HandlerFunc {
	return func(c *gin.Context) {
		ip := c.ClientIP()
		
		globalRateLimiter.mu.Lock()
		visitor, exists := globalRateLimiter.visitors[ip]
		
		if !exists {
			visitor = &Visitor{
				Count:    1,
				LastSeen: time.Now(),
			}
			globalRateLimiter.visitors[ip] = visitor
			globalRateLimiter.mu.Unlock()
			c.Next()
			return
		}

		// 如果超过时间窗口，重置计数
		if time.Since(visitor.LastSeen) > globalRateLimiter.window {
			visitor.Count = 1
			visitor.LastSeen = time.Now()
			globalRateLimiter.mu.Unlock()
			c.Next()
			return
		}

		// 检查是否超过限制
		if visitor.Count >= rate {
			globalRateLimiter.mu.Unlock()
			utils.ErrorResponse(c, 429, "请求过于频繁，请稍后再试")
			c.Abort()
			return
		}

		visitor.Count++
		visitor.LastSeen = time.Now()
		globalRateLimiter.mu.Unlock()
		c.Next()
	}
}

// cleanup 清理过期的访问者记录
func (rl *RateLimiter) cleanup() {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	now := time.Now()
	for ip, visitor := range rl.visitors {
		if now.Sub(visitor.LastSeen) > 10*time.Minute {
			delete(rl.visitors, ip)
		}
	}
}

