import { useState, useEffect } from 'react'
import axios from 'axios'
import Card from '../components/Card'
import Button from '../components/Button'
import Loading from '../components/Loading'
import './Notifications.css'

interface Notification {
  id: number
  type: string
  title: string
  message: string
  read: boolean
  created_at: string
  read_at?: string
}

export default function Notifications() {
  const [notifications, setNotifications] = useState<Notification[]>([])
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState('')
  const [unreadOnly, setUnreadOnly] = useState(false)

  useEffect(() => {
    fetchNotifications()
  }, [unreadOnly])

  const fetchNotifications = async () => {
    try {
      setLoading(true)
      const params = unreadOnly ? '?unread_only=true' : ''
      const response = await axios.get(`/api/notifications${params}`)
      setNotifications(response.data.data || [])
    } catch (err: any) {
      setError(err.response?.data?.message || '获取通知失败')
    } finally {
      setLoading(false)
    }
  }

  const handleMarkAsRead = async (id: number) => {
    try {
      await axios.put(`/api/notifications/${id}/read`)
      fetchNotifications()
    } catch (err: any) {
      alert(err.response?.data?.message || '操作失败')
    }
  }

  const handleMarkAllAsRead = async () => {
    try {
      await axios.put('/api/notifications/read-all')
      fetchNotifications()
    } catch (err: any) {
      alert(err.response?.data?.message || '操作失败')
    }
  }

  const handleDelete = async (id: number) => {
    if (!confirm('确定要删除这条通知吗？')) return
    try {
      await axios.delete(`/api/notifications/${id}`)
      fetchNotifications()
    } catch (err: any) {
      alert(err.response?.data?.message || '删除失败')
    }
  }

  const getNotificationIcon = (type: string) => {
    switch (type) {
      case 'security':
        return '🔒'
      case 'activity':
        return '📊'
      default:
        return '📢'
    }
  }

  if (loading && notifications.length === 0) {
    return (
      <div className="notifications-page">
        <Loading text="加载中..." />
      </div>
    )
  }

  const unreadCount = notifications.filter((n) => !n.read).length

  return (
    <div className="notifications-page">
      <div className="page-header">
        <h2 className="page-title">📢 通知中心</h2>
        <div className="header-actions">
          <label className="filter-toggle">
            <input
              type="checkbox"
              checked={unreadOnly}
              onChange={(e) => setUnreadOnly(e.target.checked)}
            />
            <span>仅显示未读</span>
          </label>
          {unreadCount > 0 && (
            <Button variant="primary" onClick={handleMarkAllAsRead}>
              全部标记为已读
            </Button>
          )}
        </div>
      </div>

      {error && <div className="error-message">{error}</div>}

      {notifications.length === 0 ? (
        <Card className="empty-card">
          <div className="empty-state">
            <div className="empty-icon">📭</div>
            <p>暂无通知</p>
          </div>
        </Card>
      ) : (
        <div className="notifications-list">
          {notifications.map((notification) => (
            <Card
              key={notification.id}
              className={`notification-card ${!notification.read ? 'unread' : ''}`}
            >
              <div className="notification-header">
                <div className="notification-icon">
                  {getNotificationIcon(notification.type)}
                </div>
                <div className="notification-content">
                  <h3 className="notification-title">{notification.title}</h3>
                  <p className="notification-message">{notification.message}</p>
                  <div className="notification-meta">
                    <span className="notification-time">
                      {new Date(notification.created_at).toLocaleString('zh-CN')}
                    </span>
                    {notification.type && (
                      <span className={`notification-type type-${notification.type}`}>
                        {notification.type === 'security'
                          ? '安全'
                          : notification.type === 'activity'
                          ? '活动'
                          : '系统'}
                      </span>
                    )}
                  </div>
                </div>
                <div className="notification-actions">
                  {!notification.read && (
                    <Button
                      variant="secondary"
                      size="small"
                      onClick={() => handleMarkAsRead(notification.id)}
                    >
                      标记已读
                    </Button>
                  )}
                  <Button
                    variant="outline"
                    size="small"
                    onClick={() => handleDelete(notification.id)}
                  >
                    删除
                  </Button>
                </div>
              </div>
            </Card>
          ))}
        </div>
      )}
    </div>
  )
}


