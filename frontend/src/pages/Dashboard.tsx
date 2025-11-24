import { useEffect } from 'react'
import { useNavigate, Link } from 'react-router-dom'
import { useAuthStore } from '../stores/authStore'
import Card from '../components/Card'
import Button from '../components/Button'
import './Dashboard.css'

export default function Dashboard() {
  const navigate = useNavigate()
  const { user, logout, isAuthenticated } = useAuthStore()

  useEffect(() => {
    if (!isAuthenticated) {
      navigate('/login')
    }
  }, [isAuthenticated, navigate])

  if (!user) {
    return null
  }

  const handleLogout = () => {
    logout()
    navigate('/login')
  }

  return (
    <div className="dashboard-page">
      <div className="dashboard-container">
        <header className="dashboard-header">
          <h1 className="dashboard-title">✨ 欢迎回来，{user.nickname || user.username}！</h1>
          <Button variant="outline" onClick={handleLogout}>
            退出登录
          </Button>
        </header>

        <div className="dashboard-grid">
          <Card title="📋 个人信息" className="dashboard-card">
            <div className="profile-info">
              <div className="info-item">
                <span className="info-label">用户名：</span>
                <span className="info-value">{user.username}</span>
              </div>
              <div className="info-item">
                <span className="info-label">邮箱：</span>
                <span className="info-value">{user.email}</span>
              </div>
              {user.nickname && (
                <div className="info-item">
                  <span className="info-label">昵称：</span>
                  <span className="info-value">{user.nickname}</span>
                </div>
              )}
            </div>
            <div style={{ display: 'flex', gap: '12px', marginTop: '16px' }}>
              <Link to="/profile" style={{ flex: 1 }}>
                <Button variant="secondary" fullWidth>
                  编辑资料
                </Button>
              </Link>
              <Link to="/change-password" style={{ flex: 1 }}>
                <Button variant="outline" fullWidth>
                  修改密码
                </Button>
              </Link>
            </div>
          </Card>

          <Card title="🔐 安全守护契约" className="dashboard-card">
            <p className="card-description">
              启用多因素认证（MFA），为您的账户添加额外的安全保护层。
            </p>
            <Link to="/mfa">
              <Button variant="primary" fullWidth style={{ marginTop: '16px' }}>
                设置MFA
              </Button>
            </Link>
          </Card>

          <Card title="📊 账户状态" className="dashboard-card">
            <div className="status-item">
              <span className="status-label">账户状态：</span>
              <span className="status-value status-active">正常</span>
            </div>
            <div className="status-item">
              <span className="status-label">MFA状态：</span>
              <span className="status-value">未启用</span>
            </div>
          </Card>

          <Card title="🔐 会话管理" className="dashboard-card">
            <p className="card-description">
              查看和管理您的所有活跃登录会话，保护账户安全。
            </p>
            <Link to="/sessions">
              <Button variant="secondary" fullWidth style={{ marginTop: '16px' }}>
                管理会话
              </Button>
            </Link>
          </Card>

          <Card title="👑 权限管理" className="dashboard-card">
            <p className="card-description">
              查看您当前拥有的角色和权限信息。
            </p>
            <Link to="/permissions">
              <Button variant="secondary" fullWidth style={{ marginTop: '16px' }}>
                查看权限
              </Button>
            </Link>
          </Card>

          <Card title="📊 审计日志" className="dashboard-card">
            <p className="card-description">
              查看您的账户操作记录和安全事件日志。
            </p>
            <Link to="/audit-logs">
              <Button variant="secondary" fullWidth style={{ marginTop: '16px' }}>
                查看日志
              </Button>
            </Link>
          </Card>

          <Card title="🔑 OAuth2 客户端" className="dashboard-card">
            <p className="card-description">
              管理您的OAuth2应用程序客户端。
            </p>
            <Link to="/oauth2-clients">
              <Button variant="secondary" fullWidth style={{ marginTop: '16px' }}>
                管理客户端
              </Button>
            </Link>
          </Card>
        </div>
      </div>
    </div>
  )
}

