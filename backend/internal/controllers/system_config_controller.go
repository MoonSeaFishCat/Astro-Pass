package controllers

import (
	"astro-pass/internal/services"
	"astro-pass/internal/utils"
	"github.com/gin-gonic/gin"
)

type SystemConfigController struct {
	configService *services.SystemConfigService
}

func NewSystemConfigController() *SystemConfigController {
	return &SystemConfigController{
		configService: services.NewSystemConfigService(),
	}
}

// GetAllConfigs 获取所有配置
func (c *SystemConfigController) GetAllConfigs(ctx *gin.Context) {
	configs, err := c.configService.GetAllConfigs()
	if err != nil {
		utils.ErrorResponse(ctx, 500, "获取配置失败")
		return
	}

	utils.SuccessWithMessage(ctx, "获取成功", configs)
}

// GetConfigsByCategory 按分类获取配置
func (c *SystemConfigController) GetConfigsByCategory(ctx *gin.Context) {
	category := ctx.Param("category")

	configs, err := c.configService.GetConfigsByCategory(category)
	if err != nil {
		utils.ErrorResponse(ctx, 500, "获取配置失败")
		return
	}

	utils.SuccessWithMessage(ctx, "获取成功", configs)
}

// UpdateConfig 更新配置
func (c *SystemConfigController) UpdateConfig(ctx *gin.Context) {
	var req struct {
		Key         string `json:"key" binding:"required"`
		Value       string `json:"value" binding:"required"`
		Type        string `json:"type"`
		Category    string `json:"category"`
		Label       string `json:"label"`
		Description string `json:"description"`
	}

	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(ctx, 400, "请求参数错误")
		return
	}

	// 验证配置值
	if req.Type != "" {
		if err := c.configService.ValidateConfig(req.Key, req.Value, req.Type); err != nil {
			utils.ErrorResponse(ctx, 400, err.Error())
			return
		}
	}

	if err := c.configService.SetConfig(req.Key, req.Value, req.Type, req.Category, req.Label, req.Description); err != nil {
		utils.ErrorResponse(ctx, 500, "更新配置失败")
		return
	}

	utils.SuccessWithMessage(ctx, "更新成功", nil)
}

// GetBackupConfig 获取备份配置
func (c *SystemConfigController) GetBackupConfig(ctx *gin.Context) {
	config, err := c.configService.GetBackupConfig()
	if err != nil {
		utils.ErrorResponse(ctx, 500, "获取备份配置失败")
		return
	}

	utils.SuccessWithMessage(ctx, "获取成功", config)
}

// UpdateBackupConfig 更新备份配置
func (c *SystemConfigController) UpdateBackupConfig(ctx *gin.Context) {
	var req struct {
		AutoEnabled   bool   `json:"auto_enabled"`
		Schedule      string `json:"schedule"`
		RetentionDays int    `json:"retention_days"`
		MaxBackups    int    `json:"max_backups"`
	}

	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(ctx, 400, "请求参数错误")
		return
	}

	if err := c.configService.UpdateBackupConfig(req.AutoEnabled, req.Schedule, req.RetentionDays, req.MaxBackups); err != nil {
		utils.ErrorResponse(ctx, 500, "更新备份配置失败")
		return
	}

	utils.SuccessWithMessage(ctx, "更新成功", nil)
}

// ExportConfigs 导出配置
func (c *SystemConfigController) ExportConfigs(ctx *gin.Context) {
	jsonData, err := c.configService.ExportConfigs()
	if err != nil {
		utils.ErrorResponse(ctx, 500, "导出配置失败")
		return
	}

	ctx.Header("Content-Type", "application/json")
	ctx.Header("Content-Disposition", "attachment; filename=system_config.json")
	ctx.String(200, jsonData)
}

// ImportConfigs 导入配置
func (c *SystemConfigController) ImportConfigs(ctx *gin.Context) {
	var req struct {
		Data string `json:"data" binding:"required"`
	}

	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(ctx, 400, "请求参数错误")
		return
	}

	if err := c.configService.ImportConfigs(req.Data); err != nil {
		utils.ErrorResponse(ctx, 500, "导入配置失败: "+err.Error())
		return
	}

	utils.SuccessWithMessage(ctx, "导入成功", nil)
}
