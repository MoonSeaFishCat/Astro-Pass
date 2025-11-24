package services

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"astro-pass/internal/database"
	"astro-pass/internal/models"
	"github.com/casbin/casbin/v2"
	gormadapter "github.com/casbin/gorm-adapter/v3"
)

type ABACService struct {
	enforcer *casbin.Enforcer
}

var globalABACService *ABACService

// ABACAttribute 属性定义
type ABACAttribute struct {
	UserID       uint      `json:"user_id"`
	Username     string    `json:"username"`
	Email        string    `json:"email"`
	Roles        []string  `json:"roles"`
	Department   string    `json:"department,omitempty"`
	IP           string    `json:"ip,omitempty"`
	TimeOfDay    string    `json:"time_of_day,omitempty"`
	ResourceOwner uint     `json:"resource_owner,omitempty"`
	CreatedAt    time.Time `json:"created_at,omitempty"`
}

// ResourceAttribute 资源属性
type ResourceAttribute struct {
	ResourceID   string    `json:"resource_id"`
	ResourceType string    `json:"resource_type"`
	OwnerID      uint      `json:"owner_id,omitempty"`
	Department   string    `json:"department,omitempty"`
	IsPublic     bool      `json:"is_public,omitempty"`
	CreatedAt    time.Time `json:"created_at,omitempty"`
}

// NewABACService 创建ABAC服务
func NewABACService() (*ABACService, error) {
	if globalABACService != nil {
		return globalABACService, nil
	}

	adapter, err := gormadapter.NewAdapterByDB(database.DB)
	if err != nil {
		return nil, fmt.Errorf("创建Casbin适配器失败: %w", err)
	}

	// 使用ABAC模型 - 从配置文件加载
	modelPath := "internal/config/abac_model.conf"

	enforcer, err := casbin.NewEnforcer(modelPath, adapter)
	if err != nil {
		return nil, fmt.Errorf("创建Casbin enforcer失败: %w", err)
	}

	// 注册自定义函数
	enforcer.AddFunction("eval", func(args ...interface{}) (interface{}, error) {
		if len(args) != 2 {
			return false, errors.New("eval函数需要2个参数")
		}

		policyEnv := args[0].(string)
		requestEnv := args[1].(string)

		// 解析环境属性
		var policyAttrs map[string]interface{}
		var requestAttrs map[string]interface{}

		if err := json.Unmarshal([]byte(policyEnv), &policyAttrs); err != nil {
			return false, err
		}
		if err := json.Unmarshal([]byte(requestEnv), &requestAttrs); err != nil {
			return false, err
		}

		// 检查属性匹配
		for key, value := range policyAttrs {
			if requestValue, ok := requestAttrs[key]; !ok || requestValue != value {
				return false, nil
			}
		}

		return true, nil
	})

	enforcer.EnableAutoSave(true)

	globalABACService = &ABACService{
		enforcer: enforcer,
	}

	return globalABACService, nil
}

// CheckPermissionWithAttributes 使用ABAC检查权限
func (s *ABACService) CheckPermissionWithAttributes(
	userID uint,
	resource, action string,
	userAttrs ABACAttribute,
	resourceAttrs ResourceAttribute,
) (bool, error) {
	// 获取用户角色
	user, err := s.getUserWithRoles(userID)
	if err != nil {
		return false, err
	}

	userAttrs.UserID = userID
	userAttrs.Username = user.Username
	userAttrs.Email = user.Email
	userAttrs.Roles = make([]string, 0, len(user.Roles))
	for _, role := range user.Roles {
		userAttrs.Roles = append(userAttrs.Roles, role.Name)
	}

	// 构建环境属性
	envAttrs := map[string]interface{}{
		"user":     userAttrs,
		"resource": resourceAttrs,
	}

	envJSON, _ := json.Marshal(envAttrs)

	// 检查每个角色的权限
	for _, role := range user.Roles {
		allowed, err := s.enforcer.Enforce(role.Name, resource, action, string(envJSON))
		if err != nil {
			return false, err
		}
		if allowed {
			return true, nil
		}
	}

	return false, nil
}

// AddABACPolicy 添加ABAC策略
func (s *ABACService) AddABACPolicy(roleName, resource, action string, envAttrs map[string]interface{}) error {
	envJSON, err := json.Marshal(envAttrs)
	if err != nil {
		return fmt.Errorf("序列化环境属性失败: %w", err)
	}

	_, err = s.enforcer.AddPolicy(roleName, resource, action, string(envJSON))
	if err != nil {
		return fmt.Errorf("添加ABAC策略失败: %w", err)
	}

	return nil
}

// RemoveABACPolicy 移除ABAC策略
func (s *ABACService) RemoveABACPolicy(roleName, resource, action string, envAttrs map[string]interface{}) error {
	envJSON, err := json.Marshal(envAttrs)
	if err != nil {
		return fmt.Errorf("序列化环境属性失败: %w", err)
	}

	_, err = s.enforcer.RemovePolicy(roleName, resource, action, string(envJSON))
	if err != nil {
		return fmt.Errorf("移除ABAC策略失败: %w", err)
	}

	return nil
}

func (s *ABACService) getUserWithRoles(userID uint) (*models.User, error) {
	var user models.User
	if err := database.DB.Preload("Roles").First(&user, userID).Error; err != nil {
		return nil, errors.New("用户不存在")
	}
	return &user, nil
}

