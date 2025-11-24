package middleware

import (
	"astro-pass/internal/utils"
	"time"

	"github.com/gin-gonic/gin"
)

// LoggerMiddleware 请求日志中间件
func LoggerMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		method := c.Request.Method

		c.Next()

		latency := time.Since(start)
		status := c.Writer.Status()
		clientIP := c.ClientIP()

		utils.Info("%s %s %d %v %s", method, path, status, latency, clientIP)

		// 记录慢请求
		if latency > time.Second {
			utils.Warn("慢请求: %s %s 耗时 %v", method, path, latency)
		}

		// 记录错误请求
		if status >= 400 {
			utils.Error("请求错误: %s %s %d %s", method, path, status, clientIP)
		}
	}
}


