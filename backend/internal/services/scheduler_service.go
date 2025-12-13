package services

import (
	"time"

	"astro-pass/internal/utils"
)

type SchedulerService struct {
	backupService *BackupService
	configService *SystemConfigService
	sloService    *SLOService
	stopChan      chan bool
}

func NewSchedulerService() *SchedulerService {
	return &SchedulerService{
		backupService: NewBackupService(),
		configService: NewSystemConfigService(),
		sloService:    NewSLOService(),
		stopChan:      make(chan bool),
	}
}

// Start 启动调度服务
func (s *SchedulerService) Start() {
	utils.Info("业务调度服务启动")

	// 启动自动备份任务
	go s.autoBackupTask()

	// 启动清理旧备份任务
	go s.cleanOldBackupsTask()

	// 启动SSO会话清理任务
	go s.cleanExpiredSSOSessionsTask()
}

// Stop 停止调度服务
func (s *SchedulerService) Stop() {
	utils.Info("业务调度服务停止")
	close(s.stopChan)
}

// autoBackupTask 自动备份任务
func (s *SchedulerService) autoBackupTask() {
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

			// 简单实现：检查是否到了备份时间（每天凌晨2点）
			now := time.Now()
			currentDate := now.Format("2006-01-02")
			currentHour := now.Hour()

			// 如果是凌晨2点且今天还没有备份过
			if currentHour == 2 && currentDate != lastBackupDate {
				utils.Info("开始执行自动备份...")
				if err := s.backupService.AutoBackup(); err != nil {
					utils.Error("自动备份失败: %v", err)
				} else {
					utils.Info("自动备份成功")
					lastBackupDate = currentDate
				}
			}

		case <-s.stopChan:
			return
		}
	}
}

// cleanOldBackupsTask 清理旧备份任务
func (s *SchedulerService) cleanOldBackupsTask() {
	ticker := time.NewTicker(24 * time.Hour) // 每天检查一次
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			// 获取保留天数配置
			retentionDays := s.configService.GetConfigInt("backup.retention_days", 30)

			utils.Info("开始清理旧备份（保留%d天）...", retentionDays)
			if err := s.backupService.CleanOldBackups(retentionDays); err != nil {
				utils.Error("清理旧备份失败: %v", err)
			} else {
				utils.Info("清理旧备份成功")
			}

		case <-s.stopChan:
			return
		}
	}
}

// cleanExpiredSSOSessionsTask 清理过期SSO会话任务
func (s *SchedulerService) cleanExpiredSSOSessionsTask() {
	ticker := time.NewTicker(6 * time.Hour) // 每6小时清理一次
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			utils.Info("开始清理过期SSO会话...")
			if err := s.sloService.CleanupExpiredSessions(); err != nil {
				utils.Error("清理过期SSO会话失败: %v", err)
			} else {
				utils.Info("清理过期SSO会话成功")
			}

		case <-s.stopChan:
			return
		}
	}
}

// RunOnce 立即执行一次所有任务（用于测试）
func (s *SchedulerService) RunOnce() {
	utils.Info("手动执行定时任务...")

	// 执行自动备份
	if err := s.backupService.AutoBackup(); err != nil {
		utils.Error("自动备份失败: %v", err)
	} else {
		utils.Info("自动备份成功")
	}

	// 清理旧备份
	retentionDays := s.configService.GetConfigInt("backup.retention_days", 30)
	if err := s.backupService.CleanOldBackups(retentionDays); err != nil {
		utils.Error("清理旧备份失败: %v", err)
	} else {
		utils.Info("清理旧备份成功")
	}

	// 清理过期SSO会话
	if err := s.sloService.CleanupExpiredSessions(); err != nil {
		utils.Error("清理过期SSO会话失败: %v", err)
	} else {
		utils.Info("清理过期SSO会话成功")
	}
}