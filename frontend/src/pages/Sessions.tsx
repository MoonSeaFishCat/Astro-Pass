import { useState, useEffect } from 'react'
import { useNavigate } from 'react-router-dom'
import axios from 'axios'
import Card from '../components/Card'
import Button from '../components/Button'
import Loading from '../components/Loading'
import './Sessions.css'

interface Session {
  id: number
  ip: string
  user_agent: string
  device: string
  location?: string
  last_activity: string
  created_at: string
}

export default function Sessions() {
  const navigate = useNavigate()
  const [sessions, setSessions] = useState<Session[]>([])
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState('')
  const [revoking, setRevoking] = useState<number | null>(null)

  useEffect(() => {
    fetchSessions()
  }, [])

  const fetchSessions = async () => {
    try {
      setLoading(true)
      const response = await axios.get('/api/session/list')
      setSessions(response.data.data || [])
    } catch (error: any) {
      setError(error.response?.data?.message || '获取会话列表失败')
    } finally {
      setLoading(false)
    }
  }

  const handleRevokeSession = async (sessionId: number) => {
    if (!confirm('确定要撤销这个会话吗？')) {
      return
    }

    try {
      setRevoking(sessionId)
      await axios.delete(`/api/session/${sessionId}`)
      await fetchSessions()
    } catch (error: any) {
      alert(error.response?.data?.message || '撤销会话失败')
    } finally {
      setRevoking(null)
    }
  }

  const handleRevokeAll = async () => {
    if (!confirm('确定要撤销所有其他会话吗？当前会话将保持活跃。')) {
      return
    }

    try {
      await axios.delete('/api/session/all')
      await fetchSessions()
      alert('所有其他会话已撤销')
    } catch (error: any) {
      alert(error.response?.data?.message || '撤销失败')
    }
  }

  const formatDate = (dateString: string) => {
    const date = new Date(dateString)
    return date.toLocaleString('zh-CN')
  }

  const getDeviceIcon = (device: string) => {
    switch (device) {
      case 'mobile':
        return '📱'
      case 'tablet':
        return '📱'
      default:
        return '💻'
    }
  }

  if (loading) {
    return (
      <div className="sessions-page">
        <div className="sessions-container">
          <Loading text="加载中..." />
        </div>
      </div>
    )
  }

  return (
    <div className="sessions-page">
      <div className="sessions-container">
        <header className="sessions-header">
          <h1 className="sessions-title">🔐 活跃会话</h1>
          <div className="sessions-actions">
            <Button variant="outline" onClick={() => navigate('/dashboard')}>
              返回
            </Button>
            {sessions.length > 1 && (
              <Button variant="secondary" onClick={handleRevokeAll}>
                撤销所有其他会话
              </Button>
            )}
          </div>
        </header>

        <Card className="sessions-card">
          {error && <div className="error-message">{error}</div>}

          {sessions.length === 0 ? (
            <div className="empty-state">
              <p>暂无活跃会话</p>
            </div>
          ) : (
            <div className="sessions-list">
              {sessions.map((session) => (
                <div key={session.id} className="session-item">
                  <div className="session-info">
                    <div className="session-header-info">
                      <span className="device-icon">{getDeviceIcon(session.device)}</span>
                      <div className="session-main-info">
                        <div className="session-device">{session.device === 'desktop' ? '桌面设备' : session.device === 'mobile' ? '移动设备' : '平板设备'}</div>
                        <div className="session-ip">{session.ip}</div>
                      </div>
                    </div>
                    <div className="session-details">
                      <div className="detail-item">
                        <span className="detail-label">最后活动：</span>
                        <span className="detail-value">{formatDate(session.last_activity)}</span>
                      </div>
                      <div className="detail-item">
                        <span className="detail-label">创建时间：</span>
                        <span className="detail-value">{formatDate(session.created_at)}</span>
                      </div>
                      {session.user_agent && (
                        <div className="detail-item">
                          <span className="detail-label">用户代理：</span>
                          <span className="detail-value user-agent">{session.user_agent}</span>
                        </div>
                      )}
                    </div>
                  </div>
                  <div className="session-actions">
                    <Button
                      variant="outline"
                      onClick={() => handleRevokeSession(session.id)}
                      disabled={revoking === session.id}
                    >
                      {revoking === session.id ? '撤销中...' : '撤销'}
                    </Button>
                  </div>
                </div>
              ))}
            </div>
          )}
        </Card>
      </div>
    </div>
  )
}


