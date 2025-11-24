package services

import (
	"errors"
	"regexp"
	"unicode"

	"astro-pass/internal/database"
	"astro-pass/internal/models"
	"astro-pass/internal/utils"
)

type PasswordPolicyService struct{}

func NewPasswordPolicyService() *PasswordPolicyService {
	return &PasswordPolicyService{}
}

// ValidatePassword 验证密码是否符合策略
func (s *PasswordPolicyService) ValidatePassword(password string) error {
	// 获取密码策略（简化版，实际应该从数据库或配置读取）
	policy := s.getDefaultPolicy()

	if len(password) < policy.MinLength {
		return errors.New("密码长度至少为8位")
	}

	if policy.RequireUppercase {
		hasUpper := false
		for _, char := range password {
			if unicode.IsUpper(char) {
				hasUpper = true
				break
			}
		}
		if !hasUpper {
			return errors.New("密码必须包含至少一个大写字母")
		}
	}

	if policy.RequireLowercase {
		hasLower := false
		for _, char := range password {
			if unicode.IsLower(char) {
				hasLower = true
				break
			}
		}
		if !hasLower {
			return errors.New("密码必须包含至少一个小写字母")
		}
	}

	if policy.RequireNumber {
		hasNumber := false
		for _, char := range password {
			if unicode.IsNumber(char) {
				hasNumber = true
				break
			}
		}
		if !hasNumber {
			return errors.New("密码必须包含至少一个数字")
		}
	}

	if policy.RequireSpecialChar {
		specialCharRegex := regexp.MustCompile(`[!@#$%^&*()_+\-=\[\]{};':"\\|,.<>\/?]`)
		if !specialCharRegex.MatchString(password) {
			return errors.New("密码必须包含至少一个特殊字符")
		}
	}

	return nil
}

// CheckPasswordHistory 检查密码是否在历史记录中
func (s *PasswordPolicyService) CheckPasswordHistory(userID uint, newPassword string) error {
	policy := s.getDefaultPolicy()

	var histories []models.PasswordHistory
	database.DB.Where("user_id = ?", userID).
		Order("created_at DESC").
		Limit(policy.HistoryCount).
		Find(&histories)

	for _, history := range histories {
		if utils.CheckPassword(newPassword, history.PasswordHash) {
			return errors.New("新密码不能与最近使用的密码相同")
		}
	}

	return nil
}

// SavePasswordHistory 保存密码历史
func (s *PasswordPolicyService) SavePasswordHistory(userID uint, passwordHash string) error {
	policy := s.getDefaultPolicy()

	history := models.PasswordHistory{
		UserID:       userID,
		PasswordHash: passwordHash,
	}

	if err := database.DB.Create(&history).Error; err != nil {
		return errors.New("保存密码历史失败")
	}

	// 删除超出历史记录数量的旧密码
	var count int64
	database.DB.Model(&models.PasswordHistory{}).Where("user_id = ?", userID).Count(&count)
	if count > int64(policy.HistoryCount) {
		var oldHistories []models.PasswordHistory
		database.DB.Where("user_id = ?", userID).
			Order("created_at ASC").
			Limit(int(count - int64(policy.HistoryCount))).
			Find(&oldHistories)

		for _, oldHistory := range oldHistories {
			database.DB.Delete(&oldHistory)
		}
	}

	return nil
}

func (s *PasswordPolicyService) getDefaultPolicy() models.PasswordPolicy {
	var policy models.PasswordPolicy
	if err := database.DB.First(&policy).Error; err != nil {
		// 返回默认策略
		return models.PasswordPolicy{
			MinLength:          8,
			RequireUppercase:   true,
			RequireLowercase:   true,
			RequireNumber:      true,
			RequireSpecialChar: true,
			MaxAge:             90,
			HistoryCount:       5,
		}
	}
	return policy
}


