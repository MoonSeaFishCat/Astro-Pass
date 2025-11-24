import { useState } from 'react'
import { useNavigate, Link } from 'react-router-dom'
import { useAuthStore } from '../stores/authStore'
import Button from '../components/Button'
import Input from '../components/Input'
import Card from '../components/Card'
import './Register.css'

export default function Register() {
  const navigate = useNavigate()
  const { register, isAuthenticated } = useAuthStore()

  const [username, setUsername] = useState('')
  const [email, setEmail] = useState('')
  const [password, setPassword] = useState('')
  const [confirmPassword, setConfirmPassword] = useState('')
  const [nickname, setNickname] = useState('')
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

    if (password !== confirmPassword) {
      setError('两次输入的密码不一致哦~')
      return
    }

    if (password.length < 6) {
      setError('密码长度至少为6位')
      return
    }

    setLoading(true)

    try {
      await register(username, email, password, nickname)
      navigate('/dashboard')
    } catch (err: any) {
      setError(err.message || '注册失败，请稍后再试')
    } finally {
      setLoading(false)
    }
  }

  return (
    <div className="register-page">
      <div className="register-container">
        <div className="register-header">
          <h1 className="register-title">✨ 加入星穹学院</h1>
          <p className="register-subtitle">创建您的通行证，开启安全之旅~</p>
        </div>

        <Card className="register-card">
          <form onSubmit={handleSubmit} className="register-form">
            {error && <div className="error-message">{error}</div>}

            <Input
              label="用户名"
              type="text"
              placeholder="请输入用户名（3-50个字符）"
              value={username}
              onChange={(e) => setUsername(e.target.value)}
              required
            />

            <Input
              label="邮箱"
              type="email"
              placeholder="请输入邮箱地址"
              value={email}
              onChange={(e) => setEmail(e.target.value)}
              required
            />

            <Input
              label="昵称"
              type="text"
              placeholder="请输入昵称（可选）"
              value={nickname}
              onChange={(e) => setNickname(e.target.value)}
            />

            <Input
              label="密码"
              type="password"
              placeholder="请输入密码（至少6位）"
              value={password}
              onChange={(e) => setPassword(e.target.value)}
              required
            />

            <Input
              label="确认密码"
              type="password"
              placeholder="请再次输入密码"
              value={confirmPassword}
              onChange={(e) => setConfirmPassword(e.target.value)}
              required
            />

            <Button type="submit" fullWidth disabled={loading}>
              {loading ? '注册中...' : '注册'}
            </Button>

            <div className="register-footer">
              <Link to="/login" className="link">
                已有通行证？立即登录 →
              </Link>
            </div>
          </form>
        </Card>
      </div>
    </div>
  )
}


