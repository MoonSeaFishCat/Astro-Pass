package services

import (
	"errors"
	"fmt"

	"astro-pass/internal/database"
	"astro-pass/internal/models"
	"github.com/casbin/casbin/v2"
	gormadapter "github.com/casbin/gorm-adapter/v3"
)

type PermissionService struct {
	enforcer *casbin.Enforcer
}

var globalPermissionService *PermissionService

func NewPermissionService() (*PermissionService, error) {
	if globalPermissionService != nil {
		return globalPermissionService, nil
	}

	// 使用GORM适配器
	adapter, err := gormadapter.NewAdapterByDB(database.DB)
	if err != nil {
		return nil, fmt.Errorf("创建Casbin适配器失败: %w", err)
	}

	// 创建enforcer，使用RBAC模型
	enforcer, err := casbin.NewEnforcer("internal/config/rbac_model.conf", adapter)
	if err != nil {
		// 如果配置文件不存在，返回错误
		// 注意：请确保 internal/config/rbac_model.conf 文件存在
		return nil, fmt.Errorf("创建Casbin enforcer失败，请确保配置文件 internal/config/rbac_model.conf 存在: %w", err)
	}

	// 自动保存策略
	enforcer.EnableAutoSave(true)
	if err := seedDefaultPolicies(enforcer); err != nil {
		return nil, err
	}

	globalPermissionService = &PermissionService{
		enforcer: enforcer,
	}

	return globalPermissionService, nil
}

// CheckPermission 检查用户权限
func (s *PermissionService) CheckPermission(userID uint, resource, action string) (bool, error) {
	// 获取用户角色
	user, err := s.getUserWithRoles(userID)
	if err != nil {
		return false, err
	}

	// 检查每个角色的权限
	for _, role := range user.Roles {
		allowed, err := s.enforcer.Enforce(role.Name, resource, action)
		if err != nil {
			return false, err
		}
		if allowed {
			return true, nil
		}
	}

	return false, nil
}

// AssignRole 为用户分配角色
func (s *PermissionService) AssignRole(userID uint, roleName string) error {
	// 检查角色是否存在
	var role models.Role
	if err := database.DB.Where("name = ?", roleName).First(&role).Error; err != nil {
		return errors.New("角色不存在")
	}

	// 检查用户是否存在
	var user models.User
	if err := database.DB.First(&user, userID).Error; err != nil {
		return errors.New("用户不存在")
	}

	// 检查是否已经分配
	var count int64
	database.DB.Model(&user).Association("Roles").Count()
	database.DB.Model(&user).Where("roles.name = ?", roleName).Association("Roles").Count()
	if count > 0 {
		return errors.New("用户已拥有该角色")
	}

	// 分配角色
	if err := database.DB.Model(&user).Association("Roles").Append(&role); err != nil {
		return fmt.Errorf("分配角色失败: %w", err)
	}

	// 在Casbin中添加角色关系
	_, err := s.enforcer.AddGroupingPolicy(fmt.Sprintf("user_%d", userID), roleName)
	if err != nil {
		return fmt.Errorf("添加Casbin策略失败: %w", err)
	}

	return nil
}

// RemoveRole 移除用户角色
func (s *PermissionService) RemoveRole(userID uint, roleName string) error {
	var user models.User
	if err := database.DB.First(&user, userID).Error; err != nil {
		return errors.New("用户不存在")
	}

	var role models.Role
	if err := database.DB.Where("name = ?", roleName).First(&role).Error; err != nil {
		return errors.New("角色不存在")
	}

	// 移除角色关联
	if err := database.DB.Model(&user).Association("Roles").Delete(&role); err != nil {
		return fmt.Errorf("移除角色失败: %w", err)
	}

	// 从Casbin中移除角色关系
	_, err := s.enforcer.RemoveGroupingPolicy(fmt.Sprintf("user_%d", userID), roleName)
	if err != nil {
		return fmt.Errorf("移除Casbin策略失败: %w", err)
	}

	return nil
}

// CreateRole 创建角色
func (s *PermissionService) CreateRole(name, displayName, description string) (*models.Role, error) {
	role := &models.Role{
		Name:        name,
		DisplayName: displayName,
		Description: description,
	}

	if err := database.DB.Create(role).Error; err != nil {
		return nil, errors.New("角色创建失败")
	}

	return role, nil
}

// CreatePermission 创建权限
func (s *PermissionService) CreatePermission(name, displayName, resource, action, description string) (*models.Permission, error) {
	permission := &models.Permission{
		Name:        name,
		DisplayName: displayName,
		Resource:    resource,
		Action:      action,
		Description: description,
	}

	if err := database.DB.Create(permission).Error; err != nil {
		return nil, errors.New("权限创建失败")
	}

	return permission, nil
}

func seedDefaultPolicies(enforcer *casbin.Enforcer) error {
	defaultPolicies := [][]string{
		{"admin", "user", "read"},
		{"admin", "user", "write"},
		{"admin", "role", "read"},
		{"admin", "role", "write"},
		{"admin", "permission", "read"},
		{"admin", "permission", "write"},
	}

	for _, policy := range defaultPolicies {
		policyArgs := make([]interface{}, len(policy))
		for i, v := range policy {
			policyArgs[i] = v
		}

		if !enforcer.HasPolicy(policyArgs...) {
			if _, err := enforcer.AddPolicy(policyArgs...); err != nil {
				return fmt.Errorf("添加默认策略失败: %w", err)
			}
		}
	}

	return nil
}

// AssignPermissionToRole 为角色分配权限
func (s *PermissionService) AssignPermissionToRole(roleName string, resource, action string) error {
	// 检查角色是否存在
	var role models.Role
	if err := database.DB.Where("name = ?", roleName).First(&role).Error; err != nil {
		return errors.New("角色不存在")
	}

	// 在Casbin中添加策略
	_, err := s.enforcer.AddPolicy(roleName, resource, action)
	if err != nil {
		return fmt.Errorf("添加权限策略失败: %w", err)
	}

	// 在数据库中查找或创建权限
	var permission models.Permission
	if err := database.DB.Where("resource = ? AND action = ?", resource, action).First(&permission).Error; err != nil {
		// 权限不存在，创建它
		permission = models.Permission{
			Name:     fmt.Sprintf("%s:%s", resource, action),
			Resource: resource,
			Action:  action,
		}
		if err := database.DB.Create(&permission).Error; err != nil {
			return errors.New("创建权限失败")
		}
	}

	// 关联角色和权限
	if err := database.DB.Model(&role).Association("Permissions").Append(&permission); err != nil {
		return fmt.Errorf("分配权限失败: %w", err)
	}

	return nil
}

// getUserWithRoles 获取用户及其角色
func (s *PermissionService) getUserWithRoles(userID uint) (*models.User, error) {
	var user models.User
	if err := database.DB.Preload("Roles").First(&user, userID).Error; err != nil {
		return nil, errors.New("用户不存在")
	}
	return &user, nil
}

// GetUserRoles 获取用户角色列表
func (s *PermissionService) GetUserRoles(userID uint) ([]models.Role, error) {
	user, err := s.getUserWithRoles(userID)
	if err != nil {
		return nil, err
	}
	return user.Roles, nil
}

// GetRolePermissions 获取角色权限列表
func (s *PermissionService) GetRolePermissions(roleName string) ([]string, error) {
	// 从Casbin获取策略
	policies := s.enforcer.GetPermissionsForUser(roleName)
	permissions := make([]string, 0, len(policies))
	for _, policy := range policies {
		if len(policy) >= 3 {
			permissions = append(permissions, fmt.Sprintf("%s:%s", policy[1], policy[2]))
		}
	}
	return permissions, nil
}

// GetAllRoles 获取所有角色列表
func (s *PermissionService) GetAllRoles() ([]models.Role, error) {
	var roles []models.Role
	if err := database.DB.Preload("Permissions").Find(&roles).Error; err != nil {
		return nil, errors.New("获取角色列表失败")
	}
	return roles, nil
}

// GetAllPermissions 获取所有权限列表
func (s *PermissionService) GetAllPermissions() ([]models.Permission, error) {
	var permissions []models.Permission
	if err := database.DB.Find(&permissions).Error; err != nil {
		return nil, errors.New("获取权限列表失败")
	}
	return permissions, nil
}

// UpdateRole 更新角色信息
func (s *PermissionService) UpdateRole(roleID uint, displayName, description string) (*models.Role, error) {
	var role models.Role
	if err := database.DB.First(&role, roleID).Error; err != nil {
		return nil, errors.New("角色不存在")
	}

	role.DisplayName = displayName
	role.Description = description

	if err := database.DB.Save(&role).Error; err != nil {
		return nil, errors.New("更新角色失败")
	}

	return &role, nil
}

// DeleteRole 删除角色
func (s *PermissionService) DeleteRole(roleID uint) error {
	var role models.Role
	if err := database.DB.First(&role, roleID).Error; err != nil {
		return errors.New("角色不存在")
	}

	if err := database.DB.Delete(&role).Error; err != nil {
		return errors.New("删除角色失败")
	}

	return nil
}

// UpdatePermission 更新权限信息
func (s *PermissionService) UpdatePermission(permissionID uint, displayName, description string) (*models.Permission, error) {
	var permission models.Permission
	if err := database.DB.First(&permission, permissionID).Error; err != nil {
		return nil, errors.New("权限不存在")
	}

	permission.DisplayName = displayName
	permission.Description = description

	if err := database.DB.Save(&permission).Error; err != nil {
		return nil, errors.New("更新权限失败")
	}

	return &permission, nil
}

// DeletePermission 删除权限
func (s *PermissionService) DeletePermission(permissionID uint) error {
	var permission models.Permission
	if err := database.DB.First(&permission, permissionID).Error; err != nil {
		return errors.New("权限不存在")
	}

	if err := database.DB.Delete(&permission).Error; err != nil {
		return errors.New("删除权限失败")
	}

	return nil
}

