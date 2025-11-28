package services

import (
	"encoding/json"
	"errors"
	"strconv"

	"astro-pass/internal/database"
	"astro-pass/internal/models"
)

type SystemConfigService struct{}

func NewSystemConfigService() *SystemConfigService {
	return &SystemConfigService{}
}

// GetConfig 获取配置
func (s *SystemConfigService) GetConfig(key string) (*models.SystemConfig, error) {
	var config models.SystemConfig
	if err := database.DB.Where("key = ?", key).First(&config).Error; err != nil {
		return nil, err
	}
	return &config, nil
}

// GetConfigValue 获取配置值
func (s *SystemConfigService) GetConfigValue(key string, defaultValue string) string {
	config, err := s.GetConfig(key)
	if err != nil {
		return defaultValue
	}
	return config.Value
}

// GetConfigInt 获取整数配置
func (s *SystemConfigService) GetConfigInt(key string, defaultValue int) int {
	value := s.GetConfigValue(key, "")
	if value == "" {
		return defaultValue
	}
	intValue, err := strconv.Atoi(value)
	if err != nil {
		return defaultValue
	}
	return intValue
}

// GetConfigBool 获取布尔配置
func (s *SystemConfigService) GetConfigBool(key string, defaultValue bool) bool {
	value := s.GetConfigValue(key, "")
	if value == "" {
		return defaultValue
	}
	boolValue, err := strconv.ParseBool(value)
	if err != nil {
		return defaultValue
	}
	return boolValue
}

// SetConfig 设置配置
func (s *SystemConfigService) SetConfig(key, value, configType, category, label, description string) error {
	var config models.SystemConfig
	err := database.DB.Where("key = ?", key).First(&config).Error

	if err != nil {
		// 创建新配置
		config = models.SystemConfig{
			Key:         key,
			Value:       value,
			Type:        configType,
			Category:    category,
			Label:       label,
			Description: description,
		}
		return database.DB.Create(&config).Error
	}

	// 更新配置
	config.Value = value
	if configType != "" {
		config.Type = configType
	}
	if category != "" {
		config.Category = category
	}
	if label != "" {
		config.Label = label
	}
	if description != "" {
		config.Description = description
	}

	return database.DB.Save(&config).Error
}

// GetConfigsByCategory 按分类获取配置
func (s *SystemConfigService) GetConfigsByCategory(category string) ([]models.SystemConfig, error) {
	var configs []models.SystemConfig
	if err := database.DB.Where("category = ?", category).Find(&configs).Error; err != nil {
		return nil, err
	}
	return configs, nil
}

// GetAllConfigs 获取所有配置
func (s *SystemConfigService) GetAllConfigs() ([]models.SystemConfig, error) {
	var configs []models.SystemConfig
	if err := database.DB.Find(&configs).Error; err != nil {
		return nil, err
	}
	return configs, nil
}

// DeleteConfig 删除配置
func (s *SystemConfigService) DeleteConfig(key string) error {
	return database.DB.Where("key = ?", key).Delete(&models.SystemConfig{}).Error
}

// InitDefaultConfigs 初始化默认配置
func (s *SystemConfigService) InitDefaultConfigs() error {
	defaultConfigs := []models.SystemConfig{
		// 备份配置
		{
			Key:         "backup.auto_enabled",
			Value:       "true",
			Type:        "boolean",
			Category:    "backup",
			Label:       "启用自动备份",
			Description: "是否启用自动备份功能",
		},
		{
			Key:         "backup.auto_schedule",
			Value:       "0 2 * * *",
			Type:        "string",
			Category:    "backup",
			Label:       "自动备份时间",
			Description: "Cron表达式，默认每天凌晨2点",
		},
		{
			Key:         "backup.retention_days",
			Value:       "30",
			Type:        "number",
			Category:    "backup",
			Label:       "备份保留天数",
			Description: "自动备份保留天数，超过将被清理",
		},
		{
			Key:         "backup.max_backups",
			Value:       "10",
			Type:        "number",
			Category:    "backup",
			Label:       "最大备份数量",
			Description: "保留的最大备份数量",
		},
		// 安全配置
		{
			Key:         "security.session_timeout",
			Value:       "168",
			Type:        "number",
			Category:    "security",
			Label:       "会话超时时间（小时）",
			Description: "用户会话超时时间",
		},
		{
			Key:         "security.password_min_length",
			Value:       "8",
			Type:        "number",
			Category:    "security",
			Label:       "密码最小长度",
			Description: "用户密码最小长度要求",
		},
		{
			Key:         "security.login_max_attempts",
			Value:       "5",
			Type:        "number",
			Category:    "security",
			Label:       "最大登录尝试次数",
			Description: "连续登录失败次数限制",
		},
		{
			Key:         "security.account_lock_duration",
			Value:       "30",
			Type:        "number",
			Category:    "security",
			Label:       "账户锁定时长（分钟）",
			Description: "账户被锁定的时长",
		},
		// 邮件配置
		{
			Key:         "email.welcome_enabled",
			Value:       "true",
			Type:        "boolean",
			Category:    "email",
			Label:       "启用欢迎邮件",
			Description: "新用户注册时发送欢迎邮件",
		},
		{
			Key:         "email.notification_enabled",
			Value:       "true",
			Type:        "boolean",
			Category:    "email",
			Label:       "启用邮件通知",
			Description: "安全事件邮件通知",
		},
	}

	for _, config := range defaultConfigs {
		var existing models.SystemConfig
		err := database.DB.Where("key = ?", config.Key).First(&existing).Error
		if err != nil {
			// 不存在则创建
			if err := database.DB.Create(&config).Error; err != nil {
				return err
			}
		}
	}

	return nil
}

// UpdateBackupConfig 更新备份配置
func (s *SystemConfigService) UpdateBackupConfig(autoEnabled bool, schedule string, retentionDays, maxBackups int) error {
	configs := map[string]string{
		"backup.auto_enabled":   strconv.FormatBool(autoEnabled),
		"backup.auto_schedule":  schedule,
		"backup.retention_days": strconv.Itoa(retentionDays),
		"backup.max_backups":    strconv.Itoa(maxBackups),
	}

	for key, value := range configs {
		if err := s.SetConfig(key, value, "", "", "", ""); err != nil {
			return err
		}
	}

	return nil
}

// GetBackupConfig 获取备份配置
func (s *SystemConfigService) GetBackupConfig() (map[string]interface{}, error) {
	autoEnabled := s.GetConfigBool("backup.auto_enabled", true)
	schedule := s.GetConfigValue("backup.auto_schedule", "0 2 * * *")
	retentionDays := s.GetConfigInt("backup.retention_days", 30)
	maxBackups := s.GetConfigInt("backup.max_backups", 10)

	return map[string]interface{}{
		"auto_enabled":   autoEnabled,
		"schedule":       schedule,
		"retention_days": retentionDays,
		"max_backups":    maxBackups,
	}, nil
}

// ExportConfigs 导出配置（JSON格式）
func (s *SystemConfigService) ExportConfigs() (string, error) {
	configs, err := s.GetAllConfigs()
	if err != nil {
		return "", err
	}

	data, err := json.MarshalIndent(configs, "", "  ")
	if err != nil {
		return "", err
	}

	return string(data), nil
}

// ImportConfigs 导入配置（JSON格式）
func (s *SystemConfigService) ImportConfigs(jsonData string) error {
	var configs []models.SystemConfig
	if err := json.Unmarshal([]byte(jsonData), &configs); err != nil {
		return err
	}

	for _, config := range configs {
		if err := s.SetConfig(config.Key, config.Value, config.Type, config.Category, config.Label, config.Description); err != nil {
			return err
		}
	}

	return nil
}

// ValidateConfig 验证配置值
func (s *SystemConfigService) ValidateConfig(key, value, configType string) error {
	switch configType {
	case "number":
		if _, err := strconv.Atoi(value); err != nil {
			return errors.New("配置值必须是数字")
		}
	case "boolean":
		if _, err := strconv.ParseBool(value); err != nil {
			return errors.New("配置值必须是布尔值")
		}
	case "json":
		var js json.RawMessage
		if err := json.Unmarshal([]byte(value), &js); err != nil {
			return errors.New("配置值必须是有效的JSON")
		}
	}
	return nil
}
