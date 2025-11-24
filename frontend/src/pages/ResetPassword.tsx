import { useState, useEffect } from 'react'
import { useNavigate, useSearchParams, Link } from 'react-router-dom'
import axios from 'axios'
import Card from '../components/Card'
import Button from '../components/Button'
import Input from '../components/Input'
import './ResetPassword.css'

export default function ResetPassword() {
  const navigate = useNavigate()
  const [searchParams] = useSearchParams()
  const token = searchParams.get('token')

  const [password, setPassword] = useState('')
  const [confirmPassword, setConfirmPassword] = useState('')
  const [error, setError] = useState('')
  const [loading, setLoading] = useState(false)
  const [success, setSuccess] = useState(false)

  useEffect(() => {
    if (!token) {
      setError('缺少重置令牌，请通过邮件中的链接访问')
    }
  }, [token])

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault()
    setError('')

    if (!token) {
      setError('缺少重置令牌')
      return
    }

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
      await axios.post('/api/auth/reset-password', {
        token,
        new_password: password,
      })

      setSuccess(true)
      setTimeout(() => {
        navigate('/login')
      }, 2000)
    } catch (error: any) {
      setError(error.response?.data?.message || '密码重置失败')
    } finally {
      setLoading(false)
    }
  }

  return (
    <div className="reset-password-page">
      <div className="reset-password-container">
        <div className="reset-password-header">
          <h1 className="reset-password-title">🔑 重置密码</h1>
          <p className="reset-password-subtitle">请设置您的新密码</p>
        </div>

        <Card className="reset-password-card">
          {success ? (
            <div className="success-content">
              <div className="success-icon">✨</div>
              <h2>密码重置成功！</h2>
              <p>正在跳转到登录页面...</p>
            </div>
          ) : (
            <form onSubmit={handleSubmit} className="reset-password-form">
              {error && <div className="error-message">{error}</div>}

              <Input
                label="新密码"
                type="password"
                placeholder="请输入新密码（至少6位）"
                value={password}
                onChange={(e) => setPassword(e.target.value)}
                required
              />

              <Input
                label="确认新密码"
                type="password"
                placeholder="请再次输入新密码"
                value={confirmPassword}
                onChange={(e) => setConfirmPassword(e.target.value)}
                required
              />

              <Button type="submit" fullWidth disabled={loading || success || !token}>
                {loading ? '重置中...' : '重置密码'}
              </Button>

              <div className="reset-password-footer">
                <Link to="/login" className="link">
                  返回登录 →
                </Link>
              </div>
            </form>
          )}
        </Card>
      </div>
    </div>
  )
}


