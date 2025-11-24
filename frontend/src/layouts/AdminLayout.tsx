import { useState } from 'react'
import { Link, useLocation, useNavigate } from 'react-router-dom'
import { useAuthStore } from '../stores/authStore'
import './AdminLayout.css'

interface AdminLayoutProps {
  children: React.ReactNode
}

export default function AdminLayout({ children }: AdminLayoutProps) {
  const [sidebarOpen, setSidebarOpen] = useState(true)
  const location = useLocation()
  const navigate = useNavigate()
  const { user, logout } = useAuthStore()

  const handleLogout = () => {
    logout()
    navigate('/login')
  }

  const menuItems = [
    { path: '/admin', icon: '📊', label: '仪表盘', exact: true },
    { path: '/admin/users', icon: '👥', label: '用户管理' },
    { path: '/admin/roles', icon: '👑', label: '角色管理' },
    { path: '/admin/permissions', icon: '🔐', label: '权限管理' },
    { path: '/admin/audit-logs', icon: '📋', label: '审计日志' },
    { path: '/admin/oauth2-clients', icon: '🔑', label: 'OAuth2客户端' },
    { path: '/admin/settings', icon: '⚙️', label: '系统设置' },
  ]

  const isActive = (path: string, exact?: boolean) => {
    if (exact) {
      return location.pathname === path
    }
    return location.pathname.startsWith(path)
  }

  return (
    <div className="admin-layout">
      {/* 侧边栏 */}
      <aside className={`admin-sidebar ${sidebarOpen ? 'open' : 'closed'}`}>
        <div className="sidebar-header">
          <h2 className="sidebar-logo">✨ Astro-Pass</h2>
          <button
            className="sidebar-toggle"
            onClick={() => setSidebarOpen(!sidebarOpen)}
            aria-label="切换侧边栏"
          >
            {sidebarOpen ? '◀' : '▶'}
          </button>
        </div>

        <nav className="sidebar-nav">
          {menuItems.map((item) => (
            <Link
              key={item.path}
              to={item.path}
              className={`nav-item ${isActive(item.path, item.exact) ? 'active' : ''}`}
            >
              <span className="nav-icon">{item.icon}</span>
              {sidebarOpen && <span className="nav-label">{item.label}</span>}
            </Link>
          ))}
        </nav>

        <div className="sidebar-footer">
          <Link to="/dashboard" className="nav-item">
            <span className="nav-icon">🏠</span>
            {sidebarOpen && <span className="nav-label">用户门户</span>}
          </Link>
        </div>
      </aside>

      {/* 主内容区 */}
      <div className="admin-main">
        {/* 顶部导航栏 */}
        <header className="admin-header">
          <div className="header-left">
            <h1 className="page-title">管理员后台</h1>
          </div>
          <div className="header-right">
            <div className="user-info">
              <span className="user-name">{user?.nickname || user?.username}</span>
              <button className="logout-btn" onClick={handleLogout}>
                退出登录
              </button>
            </div>
          </div>
        </header>

        {/* 内容区域 */}
        <main className="admin-content">{children}</main>
      </div>
    </div>
  )
}


