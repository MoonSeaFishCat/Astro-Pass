package services

import (
	"encoding/json"
	"time"

	"astro-pass/internal/database"
	"astro-pass/internal/models"
)

type AuditService struct{}

func NewAuditService() *AuditService {
	return &AuditService{}
}

// CreateAuditLog 创建审计日志
func (s *AuditService) CreateAuditLog(userID *uint, action, resource, resourceID, message, status string, ip, userAgent string, metadata map[string]interface{}) error {
	var metadataJSON string
	if metadata != nil {
		metadataBytes, err := json.Marshal(metadata)
		if err == nil {
			metadataJSON = string(metadataBytes)
		}
	}

	auditLog := &models.AuditLog{
		UserID:     userID,
		Action:     action,
		Resource:   resource,
		ResourceID: resourceID,
		IP:         ip,
		UserAgent:  userAgent,
		Message:    message,
		Status:     status,
		Metadata:   metadataJSON,
	}

	return database.DB.Create(auditLog).Error
}

// GetAuditLogs 查询审计日志
func (s *AuditService) GetAuditLogs(userID *uint, action, resource string, startTime, endTime *time.Time, page, pageSize int) ([]models.AuditLog, int64, error) {
	var logs []models.AuditLog
	var total int64

	query := database.DB.Model(&models.AuditLog{})

	if userID != nil {
		query = query.Where("user_id = ?", *userID)
	}

	if action != "" {
		query = query.Where("action = ?", action)
	}

	if resource != "" {
		query = query.Where("resource = ?", resource)
	}

	if startTime != nil {
		query = query.Where("created_at >= ?", *startTime)
	}

	if endTime != nil {
		query = query.Where("created_at <= ?", *endTime)
	}

	// 获取总数
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 分页查询
	offset := (page - 1) * pageSize
	if err := query.Order("created_at DESC").Offset(offset).Limit(pageSize).Find(&logs).Error; err != nil {
		return nil, 0, err
	}

	return logs, total, nil
}

// GetAuditLogByID 根据ID获取审计日志
func (s *AuditService) GetAuditLogByID(id uint) (*models.AuditLog, error) {
	var log models.AuditLog
	if err := database.DB.First(&log, id).Error; err != nil {
		return nil, err
	}
	return &log, nil
}


