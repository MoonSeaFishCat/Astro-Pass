package utils

import (
	"encoding/json"
	"sync"
	"time"
)

// 临时内存存储（生产环境应使用Redis）
var (
	sessionStore = make(map[string]sessionDataEntry)
	sessionMutex sync.RWMutex
)

type sessionDataEntry struct {
	Data      []byte
	ExpiresAt time.Time
}

// StoreSessionData 存储会话数据（临时实现，生产环境应使用Redis）
func StoreSessionData(token string, data []byte) error {
	sessionMutex.Lock()
	defer sessionMutex.Unlock()

	sessionStore[token] = sessionDataEntry{
		Data:      data,
		ExpiresAt: time.Now().Add(5 * time.Minute), // 5分钟过期
	}

	// 清理过期数据
	go cleanupExpiredSessions()

	return nil
}

// GetSessionData 获取会话数据
func GetSessionData(token string) ([]byte, error) {
	sessionMutex.RLock()
	defer sessionMutex.RUnlock()

	entry, exists := sessionStore[token]
	if !exists {
		return nil, ErrSessionNotFound
	}

	if time.Now().After(entry.ExpiresAt) {
		delete(sessionStore, token)
		return nil, ErrSessionExpired
	}

	return entry.Data, nil
}

// DeleteSessionData 删除会话数据
func DeleteSessionData(token string) error {
	sessionMutex.Lock()
	defer sessionMutex.Unlock()

	delete(sessionStore, token)
	return nil
}

// cleanupExpiredSessions 清理过期的会话数据
func cleanupExpiredSessions() {
	sessionMutex.Lock()
	defer sessionMutex.Unlock()

	now := time.Now()
	for token, entry := range sessionStore {
		if now.After(entry.ExpiresAt) {
			delete(sessionStore, token)
		}
	}
}

var (
	ErrSessionNotFound = &SessionError{Message: "会话不存在"}
	ErrSessionExpired  = &SessionError{Message: "会话已过期"}
)

type SessionError struct {
	Message string
}

func (e *SessionError) Error() string {
	return e.Message
}

// MarshalSessionData 序列化会话数据
func MarshalSessionData(data interface{}) ([]byte, error) {
	return json.Marshal(data)
}

// UnmarshalSessionData 反序列化会话数据
func UnmarshalSessionData(data []byte, v interface{}) error {
	return json.Unmarshal(data, v)
}


