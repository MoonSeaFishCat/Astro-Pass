package config

import (
	"fmt"
)

// Validate 验证配置
func (c *Config) Validate() error {
	if c.Server.Port == "" {
		return fmt.Errorf("服务器端口不能为空")
	}

	if c.Database.Host == "" {
		return fmt.Errorf("数据库主机不能为空")
	}

	if c.Database.Name == "" {
		return fmt.Errorf("数据库名称不能为空")
	}

	if c.JWT.Secret == "" || len(c.JWT.Secret) < 32 {
		return fmt.Errorf("JWT密钥长度至少为32个字符")
	}

	if c.JWT.AccessTokenExpire <= 0 {
		return fmt.Errorf("访问令牌过期时间必须大于0")
	}

	if c.JWT.RefreshTokenExpire <= 0 {
		return fmt.Errorf("刷新令牌过期时间必须大于0")
	}

	if c.JWT.RefreshTokenExpire <= c.JWT.AccessTokenExpire {
		return fmt.Errorf("刷新令牌过期时间必须大于访问令牌过期时间")
	}

	return nil
}

// ValidateDatabase 验证数据库配置
func (c *Config) ValidateDatabase() error {
	if c.Database.User == "" {
		return fmt.Errorf("数据库用户名不能为空")
	}

	return nil
}

// GetAccessTokenExpireSeconds 获取访问令牌过期秒数
func (c *Config) GetAccessTokenExpireSeconds() int {
	return int(c.JWT.AccessTokenExpire.Seconds())
}

// GetRefreshTokenExpireSeconds 获取刷新令牌过期秒数
func (c *Config) GetRefreshTokenExpireSeconds() int {
	return int(c.JWT.RefreshTokenExpire.Seconds())
}


