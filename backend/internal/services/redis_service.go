package services

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"astro-pass/internal/config"
	"github.com/redis/go-redis/v9"
)

var redisClient *redis.Client
var ctx = context.Background()

type RedisService struct {
	client *redis.Client
}

// InitRedis 初始化Redis连接
func InitRedis() error {
	// 如果Redis未配置，跳过初始化
	if config.Cfg.Redis.Host == "" {
		return nil
	}

	redisClient = redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", config.Cfg.Redis.Host, config.Cfg.Redis.Port),
		Password: config.Cfg.Redis.Password,
		DB:       config.Cfg.Redis.DB,
	})

	// 测试连接
	_, err := redisClient.Ping(ctx).Result()
	if err != nil {
		return fmt.Errorf("Redis连接失败: %w", err)
	}

	return nil
}

// NewRedisService 创建Redis服务
func NewRedisService() *RedisService {
	return &RedisService{
		client: redisClient,
	}
}

// IsAvailable 检查Redis是否可用
func (s *RedisService) IsAvailable() bool {
	if s.client == nil {
		return false
	}
	_, err := s.client.Ping(ctx).Result()
	return err == nil
}

// Set 设置缓存
func (s *RedisService) Set(key string, value interface{}, expiration time.Duration) error {
	if !s.IsAvailable() {
		return fmt.Errorf("Redis不可用")
	}

	data, err := json.Marshal(value)
	if err != nil {
		return err
	}

	return s.client.Set(ctx, key, data, expiration).Err()
}

// Get 获取缓存
func (s *RedisService) Get(key string, dest interface{}) error {
	if !s.IsAvailable() {
		return fmt.Errorf("Redis不可用")
	}

	data, err := s.client.Get(ctx, key).Result()
	if err != nil {
		return err
	}

	return json.Unmarshal([]byte(data), dest)
}

// Delete 删除缓存
func (s *RedisService) Delete(key string) error {
	if !s.IsAvailable() {
		return nil
	}

	return s.client.Del(ctx, key).Err()
}

// Exists 检查键是否存在
func (s *RedisService) Exists(key string) bool {
	if !s.IsAvailable() {
		return false
	}

	result, err := s.client.Exists(ctx, key).Result()
	return err == nil && result > 0
}

// Expire 设置过期时间
func (s *RedisService) Expire(key string, expiration time.Duration) error {
	if !s.IsAvailable() {
		return nil
	}

	return s.client.Expire(ctx, key, expiration).Err()
}

// Incr 递增
func (s *RedisService) Incr(key string) (int64, error) {
	if !s.IsAvailable() {
		return 0, fmt.Errorf("Redis不可用")
	}

	return s.client.Incr(ctx, key).Result()
}

// SetNX 仅当键不存在时设置
func (s *RedisService) SetNX(key string, value interface{}, expiration time.Duration) (bool, error) {
	if !s.IsAvailable() {
		return false, fmt.Errorf("Redis不可用")
	}

	data, err := json.Marshal(value)
	if err != nil {
		return false, err
	}

	return s.client.SetNX(ctx, key, data, expiration).Result()
}

// GetTTL 获取键的剩余生存时间
func (s *RedisService) GetTTL(key string) (time.Duration, error) {
	if !s.IsAvailable() {
		return 0, fmt.Errorf("Redis不可用")
	}

	return s.client.TTL(ctx, key).Result()
}

// Keys 获取匹配模式的所有键
func (s *RedisService) Keys(pattern string) ([]string, error) {
	if !s.IsAvailable() {
		return nil, fmt.Errorf("Redis不可用")
	}

	return s.client.Keys(ctx, pattern).Result()
}

// FlushDB 清空当前数据库
func (s *RedisService) FlushDB() error {
	if !s.IsAvailable() {
		return nil
	}

	return s.client.FlushDB(ctx).Err()
}

// CacheUser 缓存用户信息
func (s *RedisService) CacheUser(userID uint, user interface{}) error {
	key := fmt.Sprintf("user:%d", userID)
	return s.Set(key, user, 1*time.Hour)
}

// GetCachedUser 获取缓存的用户信息
func (s *RedisService) GetCachedUser(userID uint, dest interface{}) error {
	key := fmt.Sprintf("user:%d", userID)
	return s.Get(key, dest)
}

// InvalidateUserCache 使用户缓存失效
func (s *RedisService) InvalidateUserCache(userID uint) error {
	key := fmt.Sprintf("user:%d", userID)
	return s.Delete(key)
}

// CachePermission 缓存权限检查结果
func (s *RedisService) CachePermission(userID uint, resource, action string, allowed bool) error {
	key := fmt.Sprintf("permission:%d:%s:%s", userID, resource, action)
	return s.Set(key, allowed, 5*time.Minute)
}

// GetCachedPermission 获取缓存的权限检查结果
func (s *RedisService) GetCachedPermission(userID uint, resource, action string) (bool, error) {
	key := fmt.Sprintf("permission:%d:%s:%s", userID, resource, action)
	var allowed bool
	err := s.Get(key, &allowed)
	return allowed, err
}

// AddToBlacklist 添加令牌到黑名单
func (s *RedisService) AddToBlacklist(token string, expiration time.Duration) error {
	key := fmt.Sprintf("blacklist:%s", token)
	return s.Set(key, true, expiration)
}

// IsBlacklisted 检查令牌是否在黑名单中
func (s *RedisService) IsBlacklisted(token string) bool {
	key := fmt.Sprintf("blacklist:%s", token)
	return s.Exists(key)
}

// RateLimitCheck 速率限制检查
func (s *RedisService) RateLimitCheck(key string, limit int64, window time.Duration) (bool, error) {
	if !s.IsAvailable() {
		return true, nil // Redis不可用时放行
	}

	count, err := s.Incr(key)
	if err != nil {
		return true, err
	}

	if count == 1 {
		s.Expire(key, window)
	}

	return count <= limit, nil
}

// Close 关闭Redis连接
func (s *RedisService) Close() error {
	if s.client != nil {
		return s.client.Close()
	}
	return nil
}
