import { useState } from 'react'
import { useNavigate, Link } from 'react-router-dom'
import { useAuthStore } from '../stores/authStore'
import Button from '../components/Button'
import Input from '../components/Input'
import Card from '../components/Card'
import './Login.css'

export default function Login() {
  const navigate = useNavigate()
  const { login, isAuthenticated } = useAuthStore()

  const [username, setUsername] = useState('')
  const [password, setPassword] = useState('')
  const [error, setError] = useState('')
  const [loading, setLoading] = useState(false)

  // 如果已登录，重定向到仪表板
  if (isAuthenticated) {
    navigate('/dashboard')
    return null
  }

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault()
    setError('')
    setLoading(true)

    try {
      await login(username, password)
      navigate('/dashboard')
    } catch (err: any) {
      setError(err.message || '登录失败，请检查您的通行证信息哦~')
    } finally {
      setLoading(false)
    }
  }

  return (
    <div className="login-page">
      <div className="login-container">
        <div className="login-header">
          <h1 className="login-title">✨ 星穹通行证</h1>
          <p className="login-subtitle">欢迎回来，请输入通行证信息哦~</p>
        </div>

        <Card className="login-card">
          <form onSubmit={handleSubmit} className="login-form">
            {error && <div className="error-message">{error}</div>}

            <Input
              label="用户名或邮箱"
              type="text"
              placeholder="请输入用户名或邮箱"
              value={username}
              onChange={(e) => setUsername(e.target.value)}
              required
            />

            <Input
              label="密码"
              type="password"
              placeholder="请输入密码"
              value={password}
              onChange={(e) => setPassword(e.target.value)}
              required
            />

            <Button type="submit" fullWidth disabled={loading}>
              {loading ? '登录中...' : '登录'}
            </Button>

            <div className="login-footer">
              <Link to="/register" className="link">
                还没有通行证？立即注册 →
              </Link>
              <Link to="/forgot-password" className="link" style={{ marginTop: '8px', display: 'block' }}>
                忘记密码？
              </Link>
            </div>
          </form>
        </Card>
      </div>
    </div>
  )
}

