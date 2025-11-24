import { useEffect } from 'react'
import { useNavigate, Link } from 'react-router-dom'
import { useAuthStore } from '../stores/authStore'
import Card from '../components/Card'
import Button from '../components/Button'
import './UserDashboard.css'

const sections = [
  { id: 'profile', label: '信息总览', hint: '查看基础资料与联系方式' },
  { id: 'security', label: '安全中心', hint: '守护账户的每一次心跳' },
  { id: 'sessions', label: '会话守护', hint: '随时掌握登录设备' },
  { id: 'notifications', label: '信箱与通知', hint: '邮箱验证与通知中心' },
]

export default function UserDashboard() {
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

  const initials = (user.nickname || user.username || '宇')
    .trim()
    .charAt(0)
    .toUpperCase()

  const isAdmin = user.roles?.some((role: any) => role.name === 'admin') || false

  const handleLogout = () => {
    logout()
    navigate('/login')
  }

  const scrollToSection = (id: string) => {
    const target = document.getElementById(id)
    if (target) {
      target.scrollIntoView({ behavior: 'smooth', block: 'start' })
    }
  }

  return (
    <div className="user-dashboard-page">
      <div className="user-layout">
        <aside className="user-sidebar">
          <div className="user-summary-card">
            <div className="avatar-circle">{initials}</div>
            <div className="user-meta">
              <div className="meta-title">{user.nickname || user.username}</div>
              <div className="meta-sub">{user.email}</div>
            </div>
            {isAdmin && (
              <Link to="/admin" className="admin-link">
                <Button variant="secondary" fullWidth>
                  进入管理员后台
                </Button>
              </Link>
            )}
          </div>

          <nav className="sidebar-menu">
            {sections.map((item) => (
              <button key={item.id} onClick={() => scrollToSection(item.id)}>
                <div className="menu-label">{item.label}</div>
                <div className="menu-hint">{item.hint}</div>
              </button>
            ))}
          </nav>

          <div className="sidebar-footer">
            <Button variant="outline" fullWidth onClick={handleLogout}>
              退出登录
            </Button>
          </div>
        </aside>

        <main className="user-main">
          <section id="profile" className="user-section">
            <div className="section-header">
              <div>
                <h2>信息总览</h2>
                <p>更新你的称呼与联系方式，保持资料温润如初。</p>
              </div>
              <div className="section-actions">
                <Link to="/profile">
                  <Button variant="secondary">编辑资料</Button>
                </Link>
                <Link to="/change-password">
                  <Button variant="outline">修改密码</Button>
                </Link>
              </div>
            </div>
            <Card className="section-card">
              <div className="profile-info">
                <div className="info-item">
                  <span className="info-label">用户名</span>
                  <span className="info-value">{user.username}</span>
                </div>
                <div className="info-item">
                  <span className="info-label">邮箱</span>
                  <span className="info-value">{user.email}</span>
                </div>
                <div className="info-item">
                  <span className="info-label">昵称</span>
                  <span className="info-value">{user.nickname || '未填写'}</span>
                </div>
                <div className="info-item">
                  <span className="info-label">身份</span>
                  <span className="info-value">
                    {user.roles?.map((role: any) => role.display_name || role.name).join('、') ||
                      '普通用户'}
                  </span>
                </div>
                <div className="info-item">
                  <span className="info-label">邮箱状态</span>
                  <span className={`info-value ${user.email_verified ? 'status-active' : 'status-pending'}`}>
                    {user.email_verified ? '已验证' : '尚未验证'}
                  </span>
                </div>
              </div>
            </Card>
          </section>

          <section id="security" className="user-section">
            <div className="section-header">
              <div>
                <h2>安全中心</h2>
                <p>柔和守护，也要坚定可靠。开启多因素认证让安心常伴。</p>
              </div>
              <Link to="/mfa">
                <Button variant="primary">前往安全设置</Button>
              </Link>
            </div>
            <div className="section-grid">
              <Card className="section-card">
                <div className="status-item">
                  <span className="status-label">账户状态</span>
                  <span className="status-value status-active">良好</span>
                </div>
                <div className="status-item">
                  <span className="status-label">MFA 状态</span>
                  <span className="status-value">{user.mfa_enabled ? '已启用' : '尚未启用'}</span>
                </div>
                <p className="card-description">
                  建议启用 MFA、多设备登录时及时核对登录地点，温柔守护自己的数字足迹。
                </p>
              </Card>
              <Card className="section-card">
                <h3>安全小贴士</h3>
                <ul className="tips-list">
                  <li>定期更换密码，使用长一些的短句或诗句更安全。</li>
                  <li>登录陌生设备后记得在会话列表中注销。</li>
                  <li>开启浏览器储存通知，随时关注异常登录。</li>
                </ul>
              </Card>
            </div>
          </section>

          <section id="sessions" className="user-section">
            <div className="section-header">
              <div>
                <h2>会话守护</h2>
                <p>在每一次旅行前，先看看哪些设备仍握着钥匙。</p>
              </div>
              <Link to="/sessions">
                <Button variant="secondary">管理会话</Button>
              </Link>
            </div>
            <Card className="section-card">
              <p className="card-description">
                这里会汇聚你在不同设备上的登录状态。若发现陌生访客，立刻撤销即可。
              </p>
              <div className="session-highlight">
                <div>
                  <div className="session-title">当前设备</div>
                  <div className="session-meta">127.0.0.1 · 今日活跃</div>
                </div>
                <Link to="/sessions">
                  <Button variant="outline">查看全部</Button>
                </Link>
              </div>
            </Card>
          </section>

          <section id="notifications" className="user-section">
            <div className="section-header">
              <div>
                <h2>信箱与通知</h2>
                <p>留一份温柔给邮箱，也留一份警觉给通知。</p>
              </div>
            </div>
            <div className="section-grid">
              <Card className="section-card">
                <h3>邮箱验证</h3>
                <p className="card-description">
                  {user.email_verified
                    ? '邮箱已经验证，可以安心接收安全提醒。'
                    : '验证邮箱后即可接收更完整的安全提示。'}
                </p>
                <Button
                  variant={user.email_verified ? 'outline' : 'primary'}
                  fullWidth
                  onClick={() => navigate('/verify-email')}
                >
                  {user.email_verified ? '查看邮箱状态' : '前往验证邮箱'}
                </Button>
              </Card>
              <Card className="section-card">
                <h3>通知中心</h3>
                <p className="card-description">
                  所有安全、活动、系统提醒都会在这里等你查阅。
                </p>
                  <Button variant="secondary" fullWidth onClick={() => navigate('/notifications')}>
                    打开通知中心
                  </Button>
              </Card>
            </div>
          </section>
        </main>
      </div>
    </div>
  )
}

