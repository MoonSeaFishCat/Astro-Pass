package controllers

import (
	"strconv"

	"astro-pass/internal/services"
	"astro-pass/internal/utils"
	"github.com/gin-gonic/gin"
)

type BackupController struct {
	backupService *services.BackupService
}

func NewBackupController() *BackupController {
	return &BackupController{
		backupService: services.NewBackupService(),
	}
}

// CreateBackup 创建备份
func (c *BackupController) CreateBackup(ctx *gin.Context) {
	userID := ctx.GetUint("user_id")

	record, err := c.backupService.CreateBackup(userID, "manual")
	if err != nil {
		utils.ErrorResponse(ctx, 500, "创建备份失败: "+err.Error())
		return
	}

	utils.SuccessResponse(ctx, "备份创建成功", record)
}

// GetBackupList 获取备份列表
func (c *BackupController) GetBackupList(ctx *gin.Context) {
	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(ctx.DefaultQuery("page_size", "20"))

	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	records, total, err := c.backupService.GetBackupList(page, pageSize)
	if err != nil {
		utils.ErrorResponse(ctx, 500, "获取备份列表失败")
		return
	}

	utils.SuccessResponse(ctx, "获取成功", gin.H{
		"backups":   records,
		"total":     total,
		"page":      page,
		"page_size": pageSize,
	})
}

// DeleteBackup 删除备份
func (c *BackupController) DeleteBackup(ctx *gin.Context) {
	backupID, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		utils.ErrorResponse(ctx, 400, "无效的备份ID")
		return
	}

	userID := ctx.GetUint("user_id")

	if err := c.backupService.DeleteBackup(uint(backupID), userID); err != nil {
		utils.ErrorResponse(ctx, 500, "删除备份失败: "+err.Error())
		return
	}

	utils.SuccessResponse(ctx, "删除成功", nil)
}

// RestoreBackup 恢复备份
func (c *BackupController) RestoreBackup(ctx *gin.Context) {
	backupID, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		utils.ErrorResponse(ctx, 400, "无效的备份ID")
		return
	}

	userID := ctx.GetUint("user_id")

	if err := c.backupService.RestoreBackup(uint(backupID), userID); err != nil {
		utils.ErrorResponse(ctx, 500, "恢复备份失败: "+err.Error())
		return
	}

	utils.SuccessResponse(ctx, "恢复成功", nil)
}

// DownloadBackup 下载备份
func (c *BackupController) DownloadBackup(ctx *gin.Context) {
	backupID, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		utils.ErrorResponse(ctx, 400, "无效的备份ID")
		return
	}

	filePath, err := c.backupService.DownloadBackup(uint(backupID))
	if err != nil {
		utils.ErrorResponse(ctx, 500, "下载备份失败: "+err.Error())
		return
	}

	ctx.File(filePath)
}

// GetBackupStats 获取备份统计
func (c *BackupController) GetBackupStats(ctx *gin.Context) {
	stats, err := c.backupService.GetBackupStats()
	if err != nil {
		utils.ErrorResponse(ctx, 500, "获取统计信息失败")
		return
	}

	utils.SuccessResponse(ctx, "获取成功", stats)
}

// CleanOldBackups 清理旧备份
func (c *BackupController) CleanOldBackups(ctx *gin.Context) {
	days, _ := strconv.Atoi(ctx.DefaultQuery("days", "30"))

	if err := c.backupService.CleanOldBackups(days); err != nil {
		utils.ErrorResponse(ctx, 500, "清理失败: "+err.Error())
		return
	}

	utils.SuccessResponse(ctx, "清理成功", nil)
}
