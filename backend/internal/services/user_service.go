package services

import (
	"errors"

	"astro-pass/internal/database"
	"astro-pass/internal/models"
	"astro-pass/internal/utils"
)

type UserService struct{}

func NewUserService() *UserService {
	return &UserService{}
}

// UpdateProfile 更新用户资料
func (s *UserService) UpdateProfile(userID uint, nickname string) (*models.User, error) {
	var user models.User
	if err := database.DB.First(&user, userID).Error; err != nil {
		return nil, errors.New("用户不存在")
	}

	user.Nickname = nickname
	if err := database.DB.Save(&user).Error; err != nil {
		return nil, errors.New("更新资料失败")
	}

	return &user, nil
}

// ChangePassword 修改密码
func (s *UserService) ChangePassword(userID uint, oldPassword, newPassword string) error {
	var user models.User
	if err := database.DB.First(&user, userID).Error; err != nil {
		return errors.New("用户不存在")
	}

	// 验证旧密码
	if !utils.CheckPassword(oldPassword, user.PasswordHash) {
		return errors.New("原密码错误")
	}

	// 加密新密码
	newPasswordHash, err := utils.HashPassword(newPassword)
	if err != nil {
		return errors.New("密码加密失败")
	}

	// 更新密码
	user.PasswordHash = newPasswordHash
	if err := database.DB.Save(&user).Error; err != nil {
		return errors.New("密码更新失败")
	}

	return nil
}

// GetUserByID 根据ID获取用户
func (s *UserService) GetUserByID(userID uint) (*models.User, error) {
	var user models.User
	if err := database.DB.First(&user, userID).Error; err != nil {
		return nil, errors.New("用户不存在")
	}
	return &user, nil
}

// GetUserByEmail 根据邮箱获取用户
func (s *UserService) GetUserByEmail(email string) (*models.User, error) {
	var user models.User
	if err := database.DB.Where("email = ?", email).First(&user).Error; err != nil {
		return nil, errors.New("用户不存在")
	}
	return &user, nil
}

// GeneratePasswordResetToken 生成密码重置令牌（简化实现，实际应该使用JWT或随机token）
func (s *UserService) GeneratePasswordResetToken(userID uint) (string, error) {
	// 这里简化处理，实际应该生成一个安全的token并存储到数据库或Redis
	token, err := utils.GenerateRefreshToken(userID)
	if err != nil {
		return "", errors.New("生成重置令牌失败")
	}
	return token, nil
}

// ResetPassword 重置密码
func (s *UserService) ResetPassword(token, newPassword string) error {
	// 这里简化处理，实际应该验证token的有效性
	// 解析token获取userID
	claims, err := utils.ParseToken(token)
	if err != nil {
		return errors.New("无效的重置令牌")
	}

	var user models.User
	if err := database.DB.First(&user, claims.UserID).Error; err != nil {
		return errors.New("用户不存在")
	}

	// 加密新密码
	newPasswordHash, err := utils.HashPassword(newPassword)
	if err != nil {
		return errors.New("密码加密失败")
	}

	// 更新密码
	user.PasswordHash = newPasswordHash
	if err := database.DB.Save(&user).Error; err != nil {
		return errors.New("密码重置失败")
	}

	return nil
}

// GetAllUsers 获取所有用户列表（管理员功能）
func (s *UserService) GetAllUsers(page, pageSize int, search string) ([]models.User, int64, error) {
	var users []models.User
	var total int64

	query := database.DB.Model(&models.User{})

	// 搜索功能
	if search != "" {
		query = query.Where("username LIKE ? OR email LIKE ? OR nickname LIKE ?",
			"%"+search+"%", "%"+search+"%", "%"+search+"%")
	}

	// 获取总数
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 分页查询
	offset := (page - 1) * pageSize
	if err := query.Preload("Roles").Offset(offset).Limit(pageSize).Order("created_at DESC").Find(&users).Error; err != nil {
		return nil, 0, err
	}

	return users, total, nil
}

// UpdateUser 更新用户信息（管理员功能）
func (s *UserService) UpdateUser(userID uint, updates map[string]interface{}) (*models.User, error) {
	var user models.User
	if err := database.DB.First(&user, userID).Error; err != nil {
		return nil, errors.New("用户不存在")
	}

	// 允许更新的字段
	allowedFields := map[string]bool{
		"nickname":       true,
		"email":          true,
		"status":         true,
		"email_verified": true,
	}

	updateData := make(map[string]interface{})
	for key, value := range updates {
		if allowedFields[key] {
			updateData[key] = value
		}
	}

	if err := database.DB.Model(&user).Updates(updateData).Error; err != nil {
		return nil, errors.New("更新用户失败")
	}

	// 重新加载用户数据
	database.DB.Preload("Roles").First(&user, userID)
	return &user, nil
}

// DeleteUser 删除用户（软删除）
func (s *UserService) DeleteUser(userID uint) error {
	var user models.User
	if err := database.DB.First(&user, userID).Error; err != nil {
		return errors.New("用户不存在")
	}

	if err := database.DB.Delete(&user).Error; err != nil {
		return errors.New("删除用户失败")
	}

	return nil
}

// AssignRoleToUser 为用户分配角色
func (s *UserService) AssignRoleToUser(userID uint, roleName string) error {
	var user models.User
	if err := database.DB.First(&user, userID).Error; err != nil {
		return errors.New("用户不存在")
	}

	var role models.Role
	if err := database.DB.Where("name = ?", roleName).First(&role).Error; err != nil {
		return errors.New("角色不存在")
	}

	if err := database.DB.Model(&user).Association("Roles").Append(&role); err != nil {
		return errors.New("分配角色失败")
	}

	return nil
}

// RemoveRoleFromUser 移除用户角色
func (s *UserService) RemoveRoleFromUser(userID uint, roleName string) error {
	var user models.User
	if err := database.DB.Preload("Roles").First(&user, userID).Error; err != nil {
		return errors.New("用户不存在")
	}

	var role models.Role
	if err := database.DB.Where("name = ?", roleName).First(&role).Error; err != nil {
		return errors.New("角色不存在")
	}

	if err := database.DB.Model(&user).Association("Roles").Delete(&role); err != nil {
		return errors.New("移除角色失败")
	}

	return nil
}

// GetUserStats 获取用户统计信息
func (s *UserService) GetUserStats() (map[string]interface{}, error) {
	var totalUsers, activeUsers, suspendedUsers, mfaEnabledUsers int64

	// 总用户数
	database.DB.Model(&models.User{}).Count(&totalUsers)

	// 活跃用户数
	database.DB.Model(&models.User{}).Where("status = ?", "active").Count(&activeUsers)

	// 暂停用户数
	database.DB.Model(&models.User{}).Where("status = ?", "suspended").Count(&suspendedUsers)

	// 启用MFA的用户数
	database.DB.Model(&models.User{}).Where("mfa_enabled = ?", true).Count(&mfaEnabledUsers)

	return map[string]interface{}{
		"total_users":        totalUsers,
		"active_users":       activeUsers,
		"suspended_users":    suspendedUsers,
		"mfa_enabled_users":  mfaEnabledUsers,
	}, nil
}

