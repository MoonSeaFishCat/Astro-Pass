package middleware

import (
	"astro-pass/internal/services"
	"astro-pass/internal/utils"
	"github.com/gin-gonic/gin"
)

var globalPermissionService *services.PermissionService

// initPermissionService 初始化权限服务（延迟初始化）
func initPermissionService() error {
	if globalPermissionService != nil {
		return nil
	}
	
	var err error
	globalPermissionService, err = services.NewPermissionService()
	if err != nil {
		utils.Error("权限服务初始化失败: %v", err)
		return err
	}
	return nil
}

// PermissionMiddleware 权限检查中间件
func PermissionMiddleware(resource, action string) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, exists := c.Get("user_id")
		if !exists {
			utils.Unauthorized(c, "未认证")
			c.Abort()
			return
		}

		// 延迟初始化权限服务
		if globalPermissionService == nil {
			if err := initPermissionService(); err != nil {
				utils.InternalError(c, "权限服务未初始化")
				c.Abort()
				return
			}
		}

		allowed, err := globalPermissionService.CheckPermission(userID.(uint), resource, action)
		if err != nil {
			utils.InternalError(c, "权限检查失败")
			c.Abort()
			return
		}

		if !allowed {
			utils.Forbidden(c, "权限不足")
			c.Abort()
			return
		}

		c.Next()
	}
}

