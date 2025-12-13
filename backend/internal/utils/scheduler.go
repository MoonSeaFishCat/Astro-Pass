package utils

import (
	"time"
)

// Scheduler 定时任务调度器
type Scheduler struct {
	stopChan chan bool
}

// NewScheduler 创建调度器
func NewScheduler() *Scheduler {
	return &Scheduler{
		stopChan: make(chan bool),
	}
}

// Start 启动调度器
func (s *Scheduler) Start() {
	Info("定时任务调度器启动")

	// 启动清理任务（简化版本，避免循环导入）
	go s.cleanupTask()
}

// Stop 停止调度器
func (s *Scheduler) Stop() {
	Info("定时任务调度器停止")
	close(s.stopChan)
}

// cleanupTask 通用清理任务
func (s *Scheduler) cleanupTask() {
	ticker := time.NewTicker(6 * time.Hour) // 每6小时执行一次清理
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			Info("执行定时清理任务...")
			// 这里可以添加一些基础的清理逻辑
			// 具体的业务清理逻辑将在服务层单独实现

		case <-s.stopChan:
			return
		}
	}
}
