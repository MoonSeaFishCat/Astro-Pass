package services

import (
	"encoding/json"
	"errors"
	"math/rand"
	"time"

	"astro-pass/internal/database"
	"astro-pass/internal/models"
	"github.com/pquerna/otp/totp"
)

type MFAService struct{}

func NewMFAService() *MFAService {
	return &MFAService{}
}

// GenerateTOTPSecret 生成TOTP密钥
func (s *MFAService) GenerateTOTPSecret(userID uint, email string) (string, string, error) {
	key, err := totp.Generate(totp.GenerateOpts{
		Issuer:      "星穹通行证",
		AccountName: email,
	})
	if err != nil {
		return "", "", errors.New("生成TOTP密钥失败")
	}

	// 保存密钥到用户记录（但标记为未启用）
	var user models.User
	if err := database.DB.First(&user, userID).Error; err != nil {
		return "", "", errors.New("用户不存在")
	}

	user.MFASecret = key.Secret()
	database.DB.Save(&user)

	return key.Secret(), key.URL(), nil
}

// VerifyTOTP 验证TOTP代码
func (s *MFAService) VerifyTOTP(userID uint, code string) (bool, error) {
	var user models.User
	if err := database.DB.First(&user, userID).Error; err != nil {
		return false, errors.New("用户不存在")
	}

	if user.MFASecret == "" {
		return false, errors.New("用户未启用MFA")
	}

	valid := totp.Validate(code, user.MFASecret)
	return valid, nil
}

// EnableMFA 启用MFA
func (s *MFAService) EnableMFA(userID uint, code string) error {
	// 先验证代码
	valid, err := s.VerifyTOTP(userID, code)
	if err != nil {
		return err
	}

	if !valid {
		return errors.New("验证码错误")
	}

	// 生成恢复码
	recoveryCodes := s.generateRecoveryCodes()

	// 更新用户
	var user models.User
	if err := database.DB.First(&user, userID).Error; err != nil {
		return errors.New("用户不存在")
	}

	recoveryCodesJSON, _ := json.Marshal(recoveryCodes)
	user.MFAEnabled = true
	user.MFARecoveryCodes = string(recoveryCodesJSON)
	database.DB.Save(&user)

	return nil
}

// DisableMFA 禁用MFA
func (s *MFAService) DisableMFA(userID uint, code string) error {
	// 验证代码或恢复码
	var user models.User
	if err := database.DB.First(&user, userID).Error; err != nil {
		return errors.New("用户不存在")
	}

	// 验证TOTP代码
	valid, _ := s.VerifyTOTP(userID, code)
	if !valid {
		// 尝试恢复码
		var recoveryCodes []string
		if err := json.Unmarshal([]byte(user.MFARecoveryCodes), &recoveryCodes); err == nil {
			valid = false
			for _, rc := range recoveryCodes {
				if rc == code {
					valid = true
					// 从列表中移除已使用的恢复码
					newCodes := []string{}
					for _, c := range recoveryCodes {
						if c != code {
							newCodes = append(newCodes, c)
						}
					}
					recoveryCodesJSON, _ := json.Marshal(newCodes)
					user.MFARecoveryCodes = string(recoveryCodesJSON)
					break
				}
			}
		}
	}

	if !valid {
		return errors.New("验证码或恢复码错误")
	}

	// 禁用MFA
	user.MFAEnabled = false
	user.MFASecret = ""
	user.MFARecoveryCodes = ""
	database.DB.Save(&user)

	return nil
}

// GetRecoveryCodes 获取恢复码
func (s *MFAService) GetRecoveryCodes(userID uint) ([]string, error) {
	var user models.User
	if err := database.DB.First(&user, userID).Error; err != nil {
		return nil, errors.New("用户不存在")
	}

	if !user.MFAEnabled {
		return nil, errors.New("用户未启用MFA")
	}

	var recoveryCodes []string
	if err := json.Unmarshal([]byte(user.MFARecoveryCodes), &recoveryCodes); err != nil {
		return nil, errors.New("恢复码解析失败")
	}

	return recoveryCodes, nil
}

// generateRecoveryCodes 生成恢复码
func (s *MFAService) generateRecoveryCodes() []string {
	codes := make([]string, 10)
	const charset = "ABCDEFGHJKLMNPQRSTUVWXYZ23456789" // 去除容易混淆的字符
	rand.Seed(time.Now().UnixNano())
	
	for i := 0; i < 10; i++ {
		code := make([]byte, 8)
		for j := range code {
			code[j] = charset[rand.Intn(len(charset))]
		}
		codes[i] = string(code)
	}
	return codes
}

