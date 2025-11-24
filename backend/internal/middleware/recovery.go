package middleware

import (
	"astro-pass/internal/utils"
	"github.com/gin-gonic/gin"
)

// RecoveryMiddleware 全局错误恢复中间件
func RecoveryMiddleware() gin.HandlerFunc {
	return gin.CustomRecovery(func(c *gin.Context, recovered interface{}) {
		if utils.ErrorLogger != nil {
			utils.ErrorLogger.Printf("Panic recovered: %v", recovered)
		}
		utils.InternalError(c, "服务器内部错误，请稍后再试")
		c.Abort()
	})
}

