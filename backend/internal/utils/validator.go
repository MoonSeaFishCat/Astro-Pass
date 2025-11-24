package utils

import (
	"regexp"
	"strings"
)

// ValidateEmail 验证邮箱格式
func ValidateEmail(email string) bool {
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)
	return emailRegex.MatchString(email)
}

// ValidateUsername 验证用户名格式（3-50个字符，只能包含字母、数字、下划线）
func ValidateUsername(username string) bool {
	if len(username) < 3 || len(username) > 50 {
		return false
	}
	usernameRegex := regexp.MustCompile(`^[a-zA-Z0-9_]+$`)
	return usernameRegex.MatchString(username)
}

// ValidatePassword 验证密码强度（至少6位，建议包含字母和数字）
func ValidatePassword(password string) bool {
	if len(password) < 6 {
		return false
	}
	return true
}

// SanitizeInput 清理输入（移除前后空格）
func SanitizeInput(input string) string {
	return strings.TrimSpace(input)
}

// ValidateNickname 验证昵称（1-50个字符）
func ValidateNickname(nickname string) bool {
	nickname = strings.TrimSpace(nickname)
	return len(nickname) >= 1 && len(nickname) <= 50
}


