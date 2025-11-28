package utils

import (
	"time"

	"astro-pass/internal/services"
)

// Scheduler 定时任务调度器
type Scheduler struct {
	backupService *services.BackupService
	configService *services.SystemConfigService
	stopChan      chan bool
}

// NewScheduler 创建调度器
func NewScheduler() *Scheduler {
	return &Scheduler{
		backupService: services.NewBackupService(),
		configService: services.NewSystemConfigService(),
		stopChan:      make(chan bool),
	}
}

// Start 启动调度器
func (s *Scheduler) Start() {
	Info("定时任务调度器启动")

	// 启动自动备份任务
	go s.autoBackupTask()

	// 启动清理旧备份任务
	go s.cleanOldBackupsTask()
}

// Stop 停止调度器
func (s *Scheduler) Stop() {
	Info("定时任务调度器停止")
	close(s.stopChan)
}

// autoBackupTask 自动备份任务
func (s *Scheduler) autoBackupTask() {
	ticker := time.NewTicker(1 * time.Hour) // 每小时检查一次
	defer ticker.Stop()

	lastBackupDate := ""

	for {
		select {
		case <-ticker.C:
			// 检查是否启用自动备份
			autoEnabled := s.configService.GetConfigBool("backup.auto_enabled", true)
			if !autoEnabled {
				continue
			}

			// 获取备份时间配置（默认每天凌晨2点）
			schedule := s.configService.GetConfigValue("backup.auto_schedule", "0 2 * * *")

			// 简单实现：检查是否到了备份时间（每天凌晨2点）
			now := time.Now()
			currentDate := now.Format("2006-01-02")
			currentHour := now.Hour()

			// 如果是凌晨2点且今天还没有备份过
			if currentHour == 2 && currentDate != lastBackupDate {
				Info("开始执行自动备份...")
				if err := s.backupService.AutoBackup(); err != nil {
					Error("自动备份失败: %v", err)
				} else {
					Info("自动备份成功")
					lastBackupDate = currentDate
				}
			}

		case <-s.stopChan:
			return
		}
	}
}

// cleanOldBackupsTask 清理旧备份任务
func (s *Scheduler) cleanOldBackupsTask() {
	ticker := time.NewTicker(24 * time.Hour) // 每天检查一次
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			// 获取保留天数配置
			retentionDays := s.configService.GetConfigInt("backup.retention_days", 30)

			Info("开始清理旧备份（保留%d天）...", retentionDays)
			if err := s.backupService.CleanOldBackups(retentionDays); err != nil {
				Error("清理旧备份失败: %v", err)
			} else {
				Info("清理旧备份成功")
			}

		case <-s.stopChan:
			return
		}
	}
}

// RunOnce 立即执行一次所有任务（用于测试）
func (s *Scheduler) RunOnce() {
	Info("手动执行定时任务...")

	// 执行自动备份
	if err := s.backupService.AutoBackup(); err != nil {
		Error("自动备份失败: %v", err)
	} else {
		Info("自动备份成功")
	}

	// 清理旧备份
	retentionDays := s.configService.GetConfigInt("backup.retention_days", 30)
	if err := s.backupService.CleanOldBackups(retentionDays); err != nil {
		Error("清理旧备份失败: %v", err)
	} else {
		Info("清理旧备份成功")
	}
}
