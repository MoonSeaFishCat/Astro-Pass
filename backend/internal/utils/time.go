package utils

import "time"

// GetCurrentTime 获取当前时间字符串
func GetCurrentTime() string {
	return time.Now().Format(time.RFC3339)
}

// GetCurrentTimestamp 获取当前时间戳
func GetCurrentTimestamp() int64 {
	return time.Now().Unix()
}


