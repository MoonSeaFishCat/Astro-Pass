import { useState, useEffect } from 'react'
import { useNavigate } from 'react-router-dom'
import axios from 'axios'
import Card from '../components/Card'
import Button from '../components/Button'
import Loading from '../components/Loading'
import './SSOSessions.css'

interface SSOSession {
  session_id: string
  client_id: string
  client_name: string
  created_at: string
  status: string
}

interface LogoutStatus {
  request_id: string
  status: string
  total_clients: number
  completed_clients: number
  failed_clients: number
  notifications: Array<{
    client_id: string
    status: string
    response_code: number
    attempt_count: number
    last_attempt_at: string
  }>
  created_at: string
  updated_at: string
}

export default function SSOSessions() {
  const navigate = useNavigate()
  const [sessions, setSessions] = useState<SSOSession[]>([])
  const [loading, setLoading] = useState(false)
  const [error, setError] = useState('')
  const [logoutLoading, setLogoutLoading] = useState<string | null>(null)
  const [showLogoutStatus, setShowLogoutStatus] = useState(false)
  const [logoutStatus, setLogoutStatus] = useState<LogoutStatus | null>(null)

  // 获取SSO会话列表
  const fetchSessions = async () => {
    setLoading(true)
    try {
      const response = await axios.get('/api/sso/sessions')
      if (response.data.code === 200) {
        setSessions(response.data.data.sessions || [])
      }
    } catch (error: any) {
      setError(error.response?.data?.message || '获取会话列表失败')
    } finally {
      setLoading(false)
    }
  }

  // 发起单点登出
  const handleLogout = async (sessionId: string) => {
    if (!confirm('确定要从所有应用中登出吗？这将结束您在所有已登录应用中的会话。')) {
      return
    }

    setLogoutLoading(sessionId)
    try {
      const response = await axios.post('/api/sso/logout', {
        session_id: sessionId,
      })
      if (response.data.code === 200) {
        alert('登出请求已发起')
        // 显示登出状态
        checkLogoutStatus(response.data.data.request_id)
        // 刷新会话列表
        fetchSessions()
      }
    } catch (error: any) {
      alert(error.response?.data?.message || '登出失败')
    } finally {
      setLogoutLoading(null)
    }
  }

  // 查看登出状态
  const checkLogoutStatus = async (requestId: string) => {
    try {
      const response = await axios.get(`/api/sso/logout/${requestId}/status`)
      if (response.data.code === 200) {
        setLogoutStatus(response.data.data)
        setShowLogoutStatus(true)
      }
    } catch (error: any) {
      alert(error.response?.data?.message || '获取登出状态失败')
    }
  }

  const formatDate = (dateString: string) => {
    const date = new Date(dateString)
    return date.toLocaleString('zh-CN')
  }

  const getStatusColor = (status: string) => {
    switch (status) {
      case 'active':
        return 'status-active'
      case 'logged_out':
        return 'status-inactive'
      default:
        return 'status-unknown'
    }
  }

  useEffect(() => {
    fetchSessions()
  }, [])

  if (loading) {
    return (
      <div className="sso-sessions-page">
        <div className="sso-sessions-container">
          <Loading text="加载中..." />
        </div>
      </div>
    )
  }

  return (
    <div className="sso-sessions-page">
      <div className="sso-sessions-container">
        <header className="sso-sessions-header">
          <h1 className="sso-sessions-title">🔐 SSO会话管理</h1>
          <div className="sso-sessions-actions">
            <Button variant="outline" onClick={() => navigate('/dashboard')}>
              返回
            </Button>
            <Button variant="secondary" onClick={fetchSessions} disabled={loading}>
              刷新
            </Button>
          </div>
        </header>

        <Card className="sso-sessions-card">
          <div className="sessions-info">
            <p>这里显示您在所有应用中的活跃会话。您可以选择从特定应用或所有应用中登出。</p>
          </div>

          {error && <div className="error-message">{error}</div>}

          {sessions.length === 0 ? (
            <div className="empty-state">
              <p>暂无SSO会话</p>
            </div>
          ) : (
            <div className="sessions-list">
              {sessions.map((session) => (
                <div key={session.session_id} className="session-item">
                  <div className="session-info">
                    <div className="session-header-info">
                      <span className="session-icon">🔗</span>
                      <div className="session-main-info">
                        <div className="session-name">{session.client_name || session.client_id}</div>
                        <div className="session-id">会话ID: {session.session_id.substring(0, 16)}...</div>
                        <div className="client-id">客户端ID: {session.client_id}</div>
                      </div>
                    </div>
                    <div className="session-details">
                      <div className="detail-item">
                        <span className="detail-label">状态：</span>
                        <span className={`detail-value ${getStatusColor(session.status)}`}>
                          {session.status === 'active' ? '活跃' : '已登出'}
                        </span>
                      </div>
                      <div className="detail-item">
                        <span className="detail-label">创建时间：</span>
                        <span className="detail-value">{formatDate(session.created_at)}</span>
                      </div>
                    </div>
                  </div>
                  <div className="session-actions">
                    {session.status === 'active' && (
                      <Button
                        variant="outline"
                        onClick={() => handleLogout(session.session_id)}
                        disabled={logoutLoading === session.session_id}
                      >
                        {logoutLoading === session.session_id ? '登出中...' : '登出'}
                      </Button>
                    )}
                  </div>
                </div>
              ))}
            </div>
          )}
        </Card>

        {/* 登出状态模态框 */}
        {showLogoutStatus && logoutStatus && (
          <div className="modal-overlay" onClick={() => setShowLogoutStatus(false)}>
            <div className="modal-content" onClick={(e) => e.stopPropagation()}>
              <div className="modal-header">
                <h3>登出状态</h3>
                <button className="modal-close" onClick={() => setShowLogoutStatus(false)}>
                  ×
                </button>
              </div>
              <div className="modal-body">
                <div className="logout-status">
                  <div className="status-summary">
                    <div className="status-item">
                      <span className="label">请求ID:</span>
                      <code>{logoutStatus.request_id}</code>
                    </div>
                    <div className="status-item">
                      <span className="label">状态:</span>
                      <span className={`status-badge ${logoutStatus.status}`}>
                        {logoutStatus.status === 'completed'
                          ? '已完成'
                          : logoutStatus.status === 'failed'
                          ? '失败'
                          : logoutStatus.status === 'processing'
                          ? '处理中'
                          : '等待中'}
                      </span>
                    </div>
                    <div className="status-item">
                      <span className="label">进度:</span>
                      <span>
                        {logoutStatus.completed_clients} / {logoutStatus.total_clients}
                        {logoutStatus.failed_clients > 0 && (
                          <span className="failed-count">
                            （失败: {logoutStatus.failed_clients}）
                          </span>
                        )}
                      </span>
                    </div>
                  </div>

                  <div className="notifications-section">
                    <h4>通知详情</h4>
                    <div className="notifications-list">
                      {logoutStatus.notifications.map((notification, index) => (
                        <div key={index} className="notification-item">
                          <div className="notification-client">{notification.client_id}</div>
                          <div className="notification-status">
                            <span className={`status-badge ${notification.status}`}>
                              {notification.status === 'success'
                                ? '成功'
                                : notification.status === 'failed'
                                ? '失败'
                                : notification.status === 'timeout'
                                ? '超时'
                                : '处理中'}
                            </span>
                          </div>
                          <div className="notification-details">
                            <span>响应码: {notification.response_code || '-'}</span>
                            <span>尝试次数: {notification.attempt_count}</span>
                            {notification.last_attempt_at && (
                              <span>最后尝试: {formatDate(notification.last_attempt_at)}</span>
                            )}
                          </div>
                        </div>
                      ))}
                    </div>
                  </div>
                </div>
              </div>
              <div className="modal-footer">
                <Button onClick={() => setShowLogoutStatus(false)}>关闭</Button>
              </div>
            </div>
          </div>
        )}
      </div>
    </div>
  )
}