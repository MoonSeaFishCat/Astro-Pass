package services

import (
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	// HTTP请求总数
	HttpRequestsTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "Total number of HTTP requests",
		},
		[]string{"method", "endpoint", "status"},
	)

	// HTTP请求持续时间
	HttpRequestDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_request_duration_seconds",
			Help:    "HTTP request duration in seconds",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method", "endpoint"},
	)

	// 登录尝试次数
	LoginAttempts = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "login_attempts_total",
			Help: "Total number of login attempts",
		},
		[]string{"status"}, // success, failed
	)

	// 注册用户数
	UserRegistrations = promauto.NewCounter(
		prometheus.CounterOpts{
			Name: "user_registrations_total",
			Help: "Total number of user registrations",
		},
	)

	// 活跃会话数
	ActiveSessions = promauto.NewGauge(
		prometheus.GaugeOpts{
			Name: "active_sessions",
			Help: "Number of active user sessions",
		},
	)

	// OAuth2令牌颁发数
	OAuth2TokensIssued = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "oauth2_tokens_issued_total",
			Help: "Total number of OAuth2 tokens issued",
		},
		[]string{"grant_type"}, // authorization_code, client_credentials, refresh_token
	)

	// MFA验证次数
	MFAVerifications = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "mfa_verifications_total",
			Help: "Total number of MFA verifications",
		},
		[]string{"method", "status"}, // method: totp, webauthn; status: success, failed
	)

	// 数据库查询持续时间
	DatabaseQueryDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "database_query_duration_seconds",
			Help:    "Database query duration in seconds",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"operation"}, // select, insert, update, delete
	)

	// Redis操作持续时间
	RedisOperationDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "redis_operation_duration_seconds",
			Help:    "Redis operation duration in seconds",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"operation"}, // get, set, delete
	)

	// 缓存命中率
	CacheHits = promauto.NewCounter(
		prometheus.CounterOpts{
			Name: "cache_hits_total",
			Help: "Total number of cache hits",
		},
	)

	CacheMisses = promauto.NewCounter(
		prometheus.CounterOpts{
			Name: "cache_misses_total",
			Help: "Total number of cache misses",
		},
	)

	// 备份操作
	BackupOperations = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "backup_operations_total",
			Help: "Total number of backup operations",
		},
		[]string{"type", "status"}, // type: manual, auto; status: success, failed
	)

	// 审计日志记录数
	AuditLogsCreated = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "audit_logs_created_total",
			Help: "Total number of audit logs created",
		},
		[]string{"action"},
	)

	// 错误计数
	ErrorsTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "errors_total",
			Help: "Total number of errors",
		},
		[]string{"type"}, // auth, database, redis, etc.
	)
)

type MetricsService struct{}

func NewMetricsService() *MetricsService {
	return &MetricsService{}
}

// RecordHTTPRequest 记录HTTP请求
func (s *MetricsService) RecordHTTPRequest(method, endpoint, status string, duration time.Duration) {
	HttpRequestsTotal.WithLabelValues(method, endpoint, status).Inc()
	HttpRequestDuration.WithLabelValues(method, endpoint).Observe(duration.Seconds())
}

// RecordLoginAttempt 记录登录尝试
func (s *MetricsService) RecordLoginAttempt(success bool) {
	status := "failed"
	if success {
		status = "success"
	}
	LoginAttempts.WithLabelValues(status).Inc()
}

// RecordUserRegistration 记录用户注册
func (s *MetricsService) RecordUserRegistration() {
	UserRegistrations.Inc()
}

// UpdateActiveSessions 更新活跃会话数
func (s *MetricsService) UpdateActiveSessions(count float64) {
	ActiveSessions.Set(count)
}

// RecordOAuth2TokenIssued 记录OAuth2令牌颁发
func (s *MetricsService) RecordOAuth2TokenIssued(grantType string) {
	OAuth2TokensIssued.WithLabelValues(grantType).Inc()
}

// RecordMFAVerification 记录MFA验证
func (s *MetricsService) RecordMFAVerification(method string, success bool) {
	status := "failed"
	if success {
		status = "success"
	}
	MFAVerifications.WithLabelValues(method, status).Inc()
}

// RecordDatabaseQuery 记录数据库查询
func (s *MetricsService) RecordDatabaseQuery(operation string, duration time.Duration) {
	DatabaseQueryDuration.WithLabelValues(operation).Observe(duration.Seconds())
}

// RecordRedisOperation 记录Redis操作
func (s *MetricsService) RecordRedisOperation(operation string, duration time.Duration) {
	RedisOperationDuration.WithLabelValues(operation).Observe(duration.Seconds())
}

// RecordCacheHit 记录缓存命中
func (s *MetricsService) RecordCacheHit() {
	CacheHits.Inc()
}

// RecordCacheMiss 记录缓存未命中
func (s *MetricsService) RecordCacheMiss() {
	CacheMisses.Inc()
}

// RecordBackupOperation 记录备份操作
func (s *MetricsService) RecordBackupOperation(backupType string, success bool) {
	status := "failed"
	if success {
		status = "success"
	}
	BackupOperations.WithLabelValues(backupType, status).Inc()
}

// RecordAuditLog 记录审计日志创建
func (s *MetricsService) RecordAuditLog(action string) {
	AuditLogsCreated.WithLabelValues(action).Inc()
}

// RecordError 记录错误
func (s *MetricsService) RecordError(errorType string) {
	ErrorsTotal.WithLabelValues(errorType).Inc()
}

// GetCacheHitRate 获取缓存命中率
func (s *MetricsService) GetCacheHitRate() float64 {
	// 这里需要从Prometheus获取实际值，简化处理
	return 0.0
}
