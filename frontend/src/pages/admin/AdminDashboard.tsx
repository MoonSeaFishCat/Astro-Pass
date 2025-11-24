import { useState, useEffect } from 'react'
import axios from 'axios'
import Card from '../../components/Card'
import Loading from '../../components/Loading'
import './AdminDashboard.css'

interface Stats {
  total_users: number
  active_users: number
  suspended_users: number
  mfa_enabled_users: number
}

export default function AdminDashboard() {
  const [stats, setStats] = useState<Stats | null>(null)
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState('')

  useEffect(() => {
    fetchStats()
  }, [])

  const fetchStats = async () => {
    try {
      setLoading(true)
      const response = await axios.get('/api/admin/users/stats')
      setStats(response.data.data)
    } catch (err: any) {
      setError(err.response?.data?.message || '获取统计信息失败')
    } finally {
      setLoading(false)
    }
  }

  if (loading) {
    return (
      <div className="admin-dashboard">
        <Loading text="加载中..." />
      </div>
    )
  }

  if (error) {
    return (
      <div className="admin-dashboard">
        <div className="error-message">{error}</div>
      </div>
    )
  }

  return (
    <div className="admin-dashboard">
      <h2 className="dashboard-title">📊 系统概览</h2>

      <div className="stats-grid">
        <Card className="stat-card">
          <div className="stat-icon">👥</div>
          <div className="stat-content">
            <div className="stat-label">总用户数</div>
            <div className="stat-value">{stats?.total_users || 0}</div>
          </div>
        </Card>

        <Card className="stat-card">
          <div className="stat-icon">✅</div>
          <div className="stat-content">
            <div className="stat-label">活跃用户</div>
            <div className="stat-value">{stats?.active_users || 0}</div>
          </div>
        </Card>

        <Card className="stat-card">
          <div className="stat-icon">⏸️</div>
          <div className="stat-content">
            <div className="stat-label">暂停用户</div>
            <div className="stat-value">{stats?.suspended_users || 0}</div>
          </div>
        </Card>

        <Card className="stat-card">
          <div className="stat-icon">🔐</div>
          <div className="stat-content">
            <div className="stat-label">启用MFA</div>
            <div className="stat-value">{stats?.mfa_enabled_users || 0}</div>
          </div>
        </Card>
      </div>

      <div className="dashboard-actions">
        <Card className="action-card">
          <h3 className="action-title">快速操作</h3>
          <div className="action-buttons">
            <a href="/admin/users" className="action-btn">
              👥 管理用户
            </a>
            <a href="/admin/roles" className="action-btn">
              👑 管理角色
            </a>
            <a href="/admin/permissions" className="action-btn">
              🔐 管理权限
            </a>
            <a href="/admin/audit-logs" className="action-btn">
              📋 查看日志
            </a>
          </div>
        </Card>
      </div>
    </div>
  )
}


