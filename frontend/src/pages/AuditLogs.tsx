import { useState, useEffect } from 'react'
import { useNavigate } from 'react-router-dom'
import axios from 'axios'
import Card from '../components/Card'
import Button from '../components/Button'
import Input from '../components/Input'
import Loading from '../components/Loading'
import './AuditLogs.css'

interface AuditLog {
  id: number
  action: string
  resource: string
  resource_id?: string
  ip: string
  user_agent?: string
  status: string
  message: string
  created_at: string
}

export default function AuditLogs() {
  const navigate = useNavigate()
  const [logs, setLogs] = useState<AuditLog[]>([])
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState('')
  const [page, setPage] = useState(1)
  const [pageSize] = useState(20)
  const [total, setTotal] = useState(0)
  const [filters, setFilters] = useState({
    action: '',
    resource: '',
  })

  useEffect(() => {
    fetchLogs()
  }, [page, filters])

  const fetchLogs = async () => {
    try {
      setLoading(true)
      const params = new URLSearchParams({
        page: page.toString(),
        page_size: pageSize.toString(),
      })
      
      if (filters.action) {
        params.append('action', filters.action)
      }
      if (filters.resource) {
        params.append('resource', filters.resource)
      }

      const response = await axios.get(`/api/audit/logs?${params.toString()}`)
      const data = response.data.data
      setLogs(data.logs || [])
      setTotal(data.total || 0)
    } catch (error: any) {
      setError(error.response?.data?.message || '获取审计日志失败')
    } finally {
      setLoading(false)
    }
  }

  const handleFilterChange = (key: string, value: string) => {
    setFilters({ ...filters, [key]: value })
    setPage(1) // 重置到第一页
  }

  const formatDate = (dateString: string) => {
    const date = new Date(dateString)
    return date.toLocaleString('zh-CN')
  }

  const getActionIcon = (action: string) => {
    const icons: Record<string, string> = {
      login: '🔐',
      logout: '🚪',
      register: '✨',
      update_profile: '📝',
      change_password: '🔑',
      enable_mfa: '🛡️',
      disable_mfa: '🔓',
    }
    return icons[action] || '📋'
  }

  const getStatusColor = (status: string) => {
    return status === 'success' ? 'var(--color-success)' : 'var(--color-error)'
  }

  const totalPages = Math.ceil(total / pageSize)

  return (
    <div className="audit-logs-page">
      <div className="audit-logs-container">
        <header className="audit-logs-header">
          <h1 className="audit-logs-title">📊 审计日志</h1>
          <Button variant="outline" onClick={() => navigate('/dashboard')}>
            返回
          </Button>
        </header>

        <Card className="audit-logs-card">
          <div className="filters-section">
            <h3 className="filters-title">筛选条件</h3>
            <div className="filters-grid">
              <Input
                label="操作类型"
                type="text"
                placeholder="如: login, register"
                value={filters.action}
                onChange={(e) => handleFilterChange('action', e.target.value)}
              />
              <Input
                label="资源类型"
                type="text"
                placeholder="如: user, session"
                value={filters.resource}
                onChange={(e) => handleFilterChange('resource', e.target.value)}
              />
            </div>
          </div>
        </Card>

        <Card className="audit-logs-card" style={{ marginTop: '24px' }}>
          {error && <div className="error-message">{error}</div>}

          {loading ? (
            <Loading text="加载中..." />
          ) : logs.length === 0 ? (
            <div className="empty-state">
              <p>暂无审计日志</p>
            </div>
          ) : (
            <>
              <div className="logs-list">
                {logs.map((log) => (
                  <div key={log.id} className="log-item">
                    <div className="log-header">
                      <div className="log-action">
                        <span className="action-icon">{getActionIcon(log.action)}</span>
                        <span className="action-name">{log.action}</span>
                        <span
                          className="log-status"
                          style={{ color: getStatusColor(log.status) }}
                        >
                          {log.status === 'success' ? '✓' : '✗'}
                        </span>
                      </div>
                      <div className="log-time">{formatDate(log.created_at)}</div>
                    </div>
                    <div className="log-content">
                      <div className="log-message">{log.message}</div>
                      <div className="log-details">
                        {log.resource && (
                          <span className="detail-badge">资源: {log.resource}</span>
                        )}
                        {log.ip && (
                          <span className="detail-badge">IP: {log.ip}</span>
                        )}
                      </div>
                    </div>
                  </div>
                ))}
              </div>

              {totalPages > 1 && (
                <div className="pagination">
                  <Button
                    variant="outline"
                    onClick={() => setPage(page - 1)}
                    disabled={page === 1}
                  >
                    上一页
                  </Button>
                  <span className="page-info">
                    第 {page} / {totalPages} 页（共 {total} 条）
                  </span>
                  <Button
                    variant="outline"
                    onClick={() => setPage(page + 1)}
                    disabled={page >= totalPages}
                  >
                    下一页
                  </Button>
                </div>
              )}
            </>
          )}
        </Card>
      </div>
    </div>
  )
}


