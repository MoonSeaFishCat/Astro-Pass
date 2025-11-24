package middleware

import "github.com/gin-gonic/gin"

// SecurityHeadersMiddleware 安全响应头中间件
func SecurityHeadersMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 防止MIME类型嗅探
		c.Header("X-Content-Type-Options", "nosniff")
		
		// 防止点击劫持
		c.Header("X-Frame-Options", "DENY")
		
		// XSS保护
		c.Header("X-XSS-Protection", "1; mode=block")
		
		// 内容安全策略
		c.Header("Content-Security-Policy", "default-src 'self'; script-src 'self' 'unsafe-inline'; style-src 'self' 'unsafe-inline'")
		
		// 推荐使用HTTPS（生产环境）
		// c.Header("Strict-Transport-Security", "max-age=31536000; includeSubDomains")
		
		// 引用策略
		c.Header("Referrer-Policy", "strict-origin-when-cross-origin")
		
		// 权限策略
		c.Header("Permissions-Policy", "geolocation=(), microphone=(), camera=()")
		
		c.Next()
	}
}


