package controllers

import (
	"astro-pass/internal/database"
	"astro-pass/internal/utils"
	"github.com/gin-gonic/gin"
)

type HealthController struct{}

func NewHealthController() *HealthController {
	return &HealthController{}
}

// Health 健康检查
func (c *HealthController) Health(ctx *gin.Context) {
	// 检查数据库连接
	var dbStatus string
	sqlDB, err := database.DB.DB()
	if err != nil {
		dbStatus = "disconnected"
	} else {
		if err := sqlDB.Ping(); err != nil {
			dbStatus = "disconnected"
		} else {
			dbStatus = "connected"
		}
	}

	status := "ok"
	if dbStatus != "connected" {
		status = "degraded"
	}

	utils.Success(ctx, gin.H{
		"status":    status,
		"database":  dbStatus,
		"service":   "星穹通行证",
		"timestamp": utils.GetCurrentTime(),
	})
}

// Ready 就绪检查（用于Kubernetes等）
func (c *HealthController) Ready(ctx *gin.Context) {
	sqlDB, err := database.DB.DB()
	if err != nil {
		utils.ErrorResponse(ctx, 503, "服务未就绪")
		return
	}

	if err := sqlDB.Ping(); err != nil {
		utils.ErrorResponse(ctx, 503, "数据库未连接")
		return
	}

	utils.Success(ctx, gin.H{
		"ready": true,
	})
}

