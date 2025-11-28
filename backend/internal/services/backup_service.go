package services

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	"astro-pass/internal/config"
	"astro-pass/internal/database"
	"astro-pass/internal/models"
)

type BackupService struct{}

func NewBackupService() *BackupService {
	return &BackupService{}
}

// BackupRecord 备份记录
type BackupRecord struct {
	ID          uint      `json:"id"`
	FileName    string    `json:"file_name"`
	FilePath    string    `json:"file_path"`
	FileSize    int64     `json:"file_size"`
	BackupType  string    `json:"backup_type"` // manual, auto
	Status      string    `json:"status"`      // success, failed, in_progress
	Message     string    `json:"message"`
	CreatedBy   uint      `json:"created_by"`
	CreatedAt   time.Time `json:"created_at"`
}

// CreateBackup 创建数据库备份
func (s *BackupService) CreateBackup(userID uint, backupType string) (*BackupRecord, error) {
	// 创建备份目录
	backupDir := "./backups"
	if err := os.MkdirAll(backupDir, 0755); err != nil {
		return nil, fmt.Errorf("创建备份目录失败: %w", err)
	}

	// 生成备份文件名
	timestamp := time.Now().Format("20060102_150405")
	fileName := fmt.Sprintf("astro_pass_backup_%s.sql", timestamp)
	filePath := filepath.Join(backupDir, fileName)

	// 创建备份记录
	record := &BackupRecord{
		FileName:   fileName,
		FilePath:   filePath,
		BackupType: backupType,
		Status:     "in_progress",
		CreatedBy:  userID,
		CreatedAt:  time.Now(),
	}

	// 执行mysqldump
	cfg := config.Cfg.Database
	cmd := exec.Command("mysqldump",
		"-h", cfg.Host,
		"-P", cfg.Port,
		"-u", cfg.User,
		fmt.Sprintf("-p%s", cfg.Password),
		"--single-transaction",
		"--quick",
		"--lock-tables=false",
		cfg.Name,
	)

	// 创建输出文件
	outFile, err := os.Create(filePath)
	if err != nil {
		record.Status = "failed"
		record.Message = fmt.Sprintf("创建备份文件失败: %v", err)
		s.saveBackupRecord(record)
		return record, err
	}
	defer outFile.Close()

	cmd.Stdout = outFile
	cmd.Stderr = os.Stderr

	// 执行备份
	if err := cmd.Run(); err != nil {
		record.Status = "failed"
		record.Message = fmt.Sprintf("备份执行失败: %v", err)
		s.saveBackupRecord(record)
		return record, err
	}

	// 获取文件大小
	fileInfo, err := os.Stat(filePath)
	if err == nil {
		record.FileSize = fileInfo.Size()
	}

	record.Status = "success"
	record.Message = "备份成功"
	s.saveBackupRecord(record)

	// 记录审计日志
	auditService := NewAuditService()
	auditService.CreateAuditLog(&userID, "backup_create", "database", fileName, "数据库备份成功", "success", "", "", nil)

	return record, nil
}

// GetBackupList 获取备份列表
func (s *BackupService) GetBackupList(page, pageSize int) ([]BackupRecord, int64, error) {
	var records []BackupRecord
	var total int64

	// 从数据库获取备份记录
	offset := (page - 1) * pageSize
	if err := database.DB.Model(&models.BackupRecord{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := database.DB.Order("created_at DESC").
		Offset(offset).
		Limit(pageSize).
		Find(&records).Error; err != nil {
		return nil, 0, err
	}

	return records, total, nil
}

// DeleteBackup 删除备份
func (s *BackupService) DeleteBackup(backupID uint, userID uint) error {
	var record models.BackupRecord
	if err := database.DB.First(&record, backupID).Error; err != nil {
		return fmt.Errorf("备份记录不存在")
	}

	// 删除文件
	if err := os.Remove(record.FilePath); err != nil {
		if !os.IsNotExist(err) {
			return fmt.Errorf("删除备份文件失败: %w", err)
		}
	}

	// 删除记录
	if err := database.DB.Delete(&record).Error; err != nil {
		return fmt.Errorf("删除备份记录失败: %w", err)
	}

	// 记录审计日志
	auditService := NewAuditService()
	auditService.CreateAuditLog(&userID, "backup_delete", "database", record.FileName, "删除备份", "success", "", "", nil)

	return nil
}

// RestoreBackup 恢复备份
func (s *BackupService) RestoreBackup(backupID uint, userID uint) error {
	var record models.BackupRecord
	if err := database.DB.First(&record, backupID).Error; err != nil {
		return fmt.Errorf("备份记录不存在")
	}

	// 检查文件是否存在
	if _, err := os.Stat(record.FilePath); os.IsNotExist(err) {
		return fmt.Errorf("备份文件不存在")
	}

	// 执行恢复
	cfg := config.Cfg.Database
	cmd := exec.Command("mysql",
		"-h", cfg.Host,
		"-P", cfg.Port,
		"-u", cfg.User,
		fmt.Sprintf("-p%s", cfg.Password),
		cfg.Name,
	)

	// 读取备份文件
	inFile, err := os.Open(record.FilePath)
	if err != nil {
		return fmt.Errorf("打开备份文件失败: %w", err)
	}
	defer inFile.Close()

	cmd.Stdin = inFile
	cmd.Stderr = os.Stderr

	// 执行恢复
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("恢复执行失败: %w", err)
	}

	// 记录审计日志
	auditService := NewAuditService()
	auditService.CreateAuditLog(&userID, "backup_restore", "database", record.FileName, "恢复备份", "success", "", "", nil)

	return nil
}

// DownloadBackup 下载备份文件
func (s *BackupService) DownloadBackup(backupID uint) (string, error) {
	var record models.BackupRecord
	if err := database.DB.First(&record, backupID).Error; err != nil {
		return "", fmt.Errorf("备份记录不存在")
	}

	// 检查文件是否存在
	if _, err := os.Stat(record.FilePath); os.IsNotExist(err) {
		return "", fmt.Errorf("备份文件不存在")
	}

	return record.FilePath, nil
}

// CleanOldBackups 清理旧备份（保留最近N天）
func (s *BackupService) CleanOldBackups(days int) error {
	cutoffTime := time.Now().AddDate(0, 0, -days)

	var oldRecords []models.BackupRecord
	if err := database.DB.Where("created_at < ? AND backup_type = ?", cutoffTime, "auto").
		Find(&oldRecords).Error; err != nil {
		return err
	}

	for _, record := range oldRecords {
		// 删除文件
		os.Remove(record.FilePath)
		// 删除记录
		database.DB.Delete(&record)
	}

	return nil
}

// AutoBackup 自动备份（定时任务调用）
func (s *BackupService) AutoBackup() error {
	// 系统用户ID为0
	_, err := s.CreateBackup(0, "auto")
	return err
}

// saveBackupRecord 保存备份记录到数据库
func (s *BackupService) saveBackupRecord(record *BackupRecord) {
	dbRecord := &models.BackupRecord{
		FileName:   record.FileName,
		FilePath:   record.FilePath,
		FileSize:   record.FileSize,
		BackupType: record.BackupType,
		Status:     record.Status,
		Message:    record.Message,
		CreatedBy:  record.CreatedBy,
	}
	database.DB.Create(dbRecord)
	record.ID = dbRecord.ID
}

// GetBackupStats 获取备份统计信息
func (s *BackupService) GetBackupStats() (map[string]interface{}, error) {
	var totalCount int64
	var successCount int64
	var failedCount int64
	var totalSize int64

	database.DB.Model(&models.BackupRecord{}).Count(&totalCount)
	database.DB.Model(&models.BackupRecord{}).Where("status = ?", "success").Count(&successCount)
	database.DB.Model(&models.BackupRecord{}).Where("status = ?", "failed").Count(&failedCount)
	database.DB.Model(&models.BackupRecord{}).Select("COALESCE(SUM(file_size), 0)").Row().Scan(&totalSize)

	// 获取最近一次备份时间
	var lastBackup models.BackupRecord
	database.DB.Order("created_at DESC").First(&lastBackup)

	stats := map[string]interface{}{
		"total_count":   totalCount,
		"success_count": successCount,
		"failed_count":  failedCount,
		"total_size":    totalSize,
		"last_backup":   lastBackup.CreatedAt,
	}

	return stats, nil
}
